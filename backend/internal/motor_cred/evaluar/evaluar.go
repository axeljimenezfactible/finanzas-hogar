// Package evaluar despacha las consultas de las fixtures de S3 (S3-interfaz.md §9.2 y §9.3)
// sobre el núcleo motor_cred. Vive en su propio paquete para que el binario común
// (cmd/motores/s3.go) no cree un ciclo con el comparador.
package evaluar

import (
	"bytes"
	"encoding/json"

	mc "finanzas-hogar/internal/motor_cred"
	"finanzas-hogar/internal/motor_cred/comparador"
)

type creditArg struct {
	CreditID int64 `json:"credit_id"`
}

type windowArg struct {
	From mc.Date `json:"from"`
	To   mc.Date `json:"to"`
}

// Evaluate decodifica input con DecodeInput (§5.7), despacha por fn y devuelve la salida
// (dinero en enteros de centavos) o un código de §9.1. Una fn desconocida devuelve "invalid".
func Evaluate(input json.RawMessage, fn string, args json.RawMessage) (out any, code string) {
	in := mc.Input{}
	if t := bytes.TrimSpace(input); len(t) > 0 && string(t) != "null" {
		var err error
		if in, err = mc.DecodeInput(bytes.NewReader(t)); err != nil {
			return nil, mc.Code(err)
		}
	}
	out, err := dispatch(in, fn, args)
	if err != nil {
		return nil, mc.Code(err)
	}
	return out, ""
}

func decode(args json.RawMessage, v any) error {
	if t := bytes.TrimSpace(args); len(t) == 0 || string(t) == "null" {
		args = json.RawMessage("{}")
	}
	return mc.DecodeStrict(bytes.NewReader(args), v)
}

type handler func(in mc.Input, args json.RawMessage) (any, error)

var handlers = map[string]handler{
	"cotizar_compra": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			Purchase mc.PurchaseInput `json:"purchase"`
		}
		return withCredit(in, raw, &a, &a.creditArg, func(st mc.CreditState) (any, error) {
			a.Purchase.Credit = st.Credit
			return mc.QuotePurchase(a.Purchase)
		})
	},
	"promociones_aplicables": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			Date      mc.Date  `json:"date"`
			Principal mc.Cents `json:"principal_cents"`
			TermCount int      `json:"term_count"`
			StoreID   *int64   `json:"store_id"`
			Merchant  string   `json:"merchant"`
		}
		return withCredit(in, raw, &a, &a.creditArg, func(st mc.CreditState) (any, error) {
			ids := []int64{}
			for _, p := range mc.ApplicablePromotions(st.Credit, a.Date, a.Principal, a.TermCount, a.StoreID, a.Merchant) {
				ids = append(ids, p.ID)
			}
			return map[string]any{"promotion_ids": ids}, nil
		})
	},
	"linea_tiempo": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			Today   mc.Date `json:"today"`
			From    mc.Date `json:"from"`
			To      mc.Date `json:"to"`
			Reading *string `json:"reading"`
		}
		return withCredit(in, raw, &a, &a.creditArg, func(st mc.CreditState) (any, error) {
			if a.Reading != nil {
				st.Credit.Reading = mc.Reading(*a.Reading)
			}
			return mc.AvailabilityTimeline(mc.TimelineInput{State: st, Cash: in.Cash, Today: a.Today, From: a.From, To: a.To})
		})
	},
	"calendario_unificado": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a windowArg
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		return mc.UnifiedCalendar(in.Credits, a.From, a.To)
	},
	"efectivo_ventana": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a windowArg
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		cal, err := mc.UnifiedCalendar(in.Credits, a.From, a.To)
		if err != nil {
			return nil, err
		}
		return mc.CashFlow(cal, in.Cash, a.From, a.To), nil
	},
	"deuda_viva": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			AsOf mc.Date `json:"as_of"`
		}
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		return mc.LiveDebt(in.Credits, a.AsOf)
	},
	"compras_vivas": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			At mc.Date `json:"at"`
		}
		return withCredit(in, raw, &a, &a.creditArg, func(st mc.CreditState) (any, error) {
			refs, err := mc.LivePurchases(st, a.At)
			if err != nil {
				return nil, err
			}
			return map[string]any{"refs": refs}, nil
		})
	},
	"cotizar_monto": func(in mc.Input, raw json.RawMessage) (any, error) {
		a, err := amountArgs(raw)
		if err != nil {
			return nil, err
		}
		if len(a.Amounts) > 0 {
			return nil, mc.ErrInvalid
		}
		return mc.QuoteByAmount(a.Amount, a.Rate, a.TermCount), nil
	},
	"sensibilidad": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a mc.SensitivityInput
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		rows, err := mc.Sensitivity(a)
		if err != nil {
			return nil, err
		}
		return map[string]any{"rows": rows}, nil
	},
	"interes_tarjeta": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			Statement mc.StatementInput `json:"statement"`
			Paid      mc.Cents          `json:"paid_cents"`
			PaidDate  mc.Date           `json:"paid_date"`
		}
		return withCredit(in, raw, &a, &a.creditArg, func(st mc.CreditState) (any, error) {
			v, prov, err := st.Credit.CardInterest(a.Statement, a.Paid, a.PaidDate)
			if err != nil {
				return nil, err
			}
			return map[string]any{"interest_cents": v, "provisional": prov}, nil
		})
	},
	"minimo_tarjeta": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			Statement mc.StatementInput `json:"statement"`
		}
		return withCredit(in, raw, &a, &a.creditArg, func(st mc.CreditState) (any, error) {
			v, prov, err := st.Credit.CardMinimum(a.Statement)
			if err != nil {
				return nil, err
			}
			return map[string]any{"minimum_cents": v, "provisional": prov}, nil
		})
	},
	"amortizar": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			Principal mc.Cents `json:"principal_cents"`
			Annual    mc.Rate  `json:"annual_rate"`
			VatBP     *int     `json:"vat_bp"`
			TermCount int      `json:"term_count"`
			Frequency string   `json:"frequency"`
			First     mc.Date  `json:"first_date"`
		}
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		return mc.Amortize(a.Principal, a.Annual, a.VatBP, a.TermCount, a.Frequency, a.First)
	},
	// «comparar» (§9.2): timeout_seconds y request; today y credits salen de datos. El comparador
	// no importa evaluar, así que no hay ciclo (corrige REVISION-S3-r1 H2).
	"comparar": func(in mc.Input, raw json.RawMessage) (any, error) {
		res, err := comparador.Query(in, raw)
		if err != nil {
			return nil, err
		}
		return res, nil
	},
	"vista_parametros": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a creditArg
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		return mc.ViewParametros(in, a.CreditID)
	},
	"vista_despensa":    viewN(mc.ViewDespensa),
	"vista_disposicion": viewN(mc.ViewDisposicion),
	"vista_credito": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			Today mc.Date `json:"today"`
			From  mc.Date `json:"from"`
			To    mc.Date `json:"to"`
		}
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		return mc.ViewCredito(in, a.CreditID, a.Today, a.From, a.To)
	},
	"vista_sensibilidad": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a mc.SensitivityInput
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		return mc.ViewSensibilidad(a)
	},
	"vista_interes_por_monto": func(in mc.Input, raw json.RawMessage) (any, error) {
		a, err := amountArgs(raw)
		if err != nil {
			return nil, err
		}
		return mc.ViewInteresPorMonto(a.Amounts, a.Rate, a.TermCount)
	},
	"vista_dos_quincenas": func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			Purchase mc.PurchaseInput `json:"purchase"`
		}
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		return mc.ViewDosQuincenas(in, a.CreditID, a.Purchase)
	},
}

type amountArg struct {
	Amount    mc.Cents   `json:"amount_cents"`
	Amounts   []mc.Cents `json:"amounts_cents"`
	Rate      mc.Rate    `json:"rate"`
	TermCount int        `json:"term_count"`
}

func amountArgs(raw json.RawMessage) (amountArg, error) {
	var a amountArg
	err := decode(raw, &a)
	return a, err
}

func viewN(f func(mc.Input, int64, int) (map[string]any, error)) handler {
	return func(in mc.Input, raw json.RawMessage) (any, error) {
		var a struct {
			creditArg
			N int `json:"n"`
		}
		if err := decode(raw, &a); err != nil {
			return nil, err
		}
		return f(in, a.CreditID, a.N)
	}
}

func withCredit(in mc.Input, raw json.RawMessage, a any, c *creditArg, f func(mc.CreditState) (any, error)) (any, error) {
	if err := decode(raw, a); err != nil {
		return nil, err
	}
	st, err := in.CreditByID(c.CreditID)
	if err != nil {
		return nil, err
	}
	return f(st)
}

func dispatch(in mc.Input, fn string, args json.RawMessage) (any, error) {
	h, ok := handlers[fn]
	if !ok {
		return nil, mc.ErrInvalid
	}
	return h(in, args)
}
