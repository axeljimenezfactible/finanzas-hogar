package recetas

import "math/big"

// Salud: edad, IMC y meta de calorías (§10; RN-17, K7, K23, K36; S4-07).

var activityFactor = map[string]*big.Rat{
	"sedentary": frac(6, 5), "light": frac(11, 8), "moderate": frac(31, 20), "active": frac(69, 40), "very_active": frac(19, 10),
}

// Estado de la meta que rige (§10.4).
const (
	TargetConfirmed = "confirmed"
	TargetSuggested = "suggested"
	TargetNone      = "none"
)

// Health es la salida de health_profile.
type Health struct {
	Person              string   `json:"person"`
	Age                 *int64   `json:"age"`
	BMITenths           *int64   `json:"bmi_tenths"`
	BMIClass            *string  `json:"bmi_class"`
	SuggestedKcalExact  *Frac    `json:"suggested_kcal_exact"`
	SuggestedKcalTarget *int64   `json:"suggested_kcal_target"`
	DailyKcalTarget     *int64   `json:"daily_kcal_target"`
	TargetStatus        string   `json:"target_status"`
	MissingInputs       []string `json:"missing_inputs"`
}

// BMIClass clasifica sobre bmi_tenths, el valor que se muestra (K7: 18.5 / 25 / 30; PROVISIONAL, bloqueo
// S4-reglas-lista, pregunta 3, opción A).
func BMIClass(tenths int64) string {
	switch {
	case tenths < 185:
		return "underweight"
	case tenths < 250:
		return "normal"
	case tenths < 300:
		return "overweight"
	}
	return "obesity"
}

// SuggestedKcal = (MSJ × factor ± ajuste × 1000) / 1000 en kcal exactas; la meta es su HALF_UP (§10.2).
// MSJ en milésimas: 10 × weight_g + 625 × height_mm − 5000 × edad + 5000 (hombre) o − 161000 (mujer).
func SuggestedKcal(weightG, heightMM, age int64, sex, activity, goal string, adjLose, adjGain int64) (*big.Rat, error) {
	f, ok := activityFactor[activity]
	if !ok {
		return nil, errf(ErrInvalidArgument, "actividad %q", activity)
	}
	msj := 10*weightG + 625*heightMM - 5000*age
	switch sex {
	case "male":
		msj += 5000
	case "female":
		msj -= 161000
	default:
		return nil, errf(ErrInvalidArgument, "sexo %q", sex)
	}
	adj := int64(0)
	switch goal {
	case "lose":
		adj = -adjLose
	case "gain":
		adj = adjGain
	case "maintain":
	default:
		return nil, errf(ErrInvalidArgument, "objetivo %q", goal)
	}
	return quoR(addR(mulR(r64(msj), f), r64(adj*1000)), r64(1000)), nil
}

// HealthFor calcula edad, IMC, clase, meta sugerida y meta que rige de una persona con plan.
func (e *Engine) HealthFor(person string) (out *Health, err error) {
	defer catch(&err)
	p, err := e.person(person)
	if err != nil {
		return nil, err
	}
	if !p.HasMealPlan {
		return nil, errf(ErrPersonWithoutPlan, "%q", person)
	}
	out = &Health{Person: person, TargetStatus: TargetNone, MissingInputs: []string{}}
	var hp HealthProfile
	for _, x := range e.Profiles {
		if x.Person == person {
			hp = x
		}
	}
	var w *WeightEntry
	for i := range e.Weights {
		x := &e.Weights[i]
		if x.Person == person && (w == nil || !x.Date.Before(w.Date)) {
			w = x
		}
	}
	if hp.BirthDate != nil {
		out.Age = ptr(hp.BirthDate.YearsUntil(e.Today))
	}
	if hp.HeightMM != nil && w != nil {
		out.BMITenths = ptr(roundHalfUp(frac(w.WeightG*10000, *hp.HeightMM**hp.HeightMM)))
		out.BMIClass = ptr(BMIClass(*out.BMITenths))
	}
	for _, m := range []struct {
		name    string
		missing bool
	}{{"sex", hp.Sex == nil}, {"birth_date", hp.BirthDate == nil}, {"height_mm", hp.HeightMM == nil},
		{"weight", w == nil}, {"activity", hp.Activity == nil}, {"goal", hp.Goal == nil}} {
		if m.missing {
			out.MissingInputs = append(out.MissingInputs, m.name)
		}
	}
	if len(out.MissingInputs) == 0 {
		lose, gain := e.Settings.HouseAdjustLoseKcal, e.Settings.HouseAdjustGainKcal
		if hp.AdjustLoseKcal != nil {
			lose = *hp.AdjustLoseKcal
		}
		if hp.AdjustGainKcal != nil {
			gain = *hp.AdjustGainKcal
		}
		s, err := SuggestedKcal(w.WeightG, *hp.HeightMM, *out.Age, *hp.Sex, *hp.Activity, *hp.Goal, lose, gain)
		if err != nil {
			return nil, err
		}
		out.SuggestedKcalExact, out.SuggestedKcalTarget = ptr(fracOf(s)), ptr(roundHalfUp(s))
	}
	switch {
	case hp.DailyKcalTarget != nil:
		out.DailyKcalTarget, out.TargetStatus = hp.DailyKcalTarget, TargetConfirmed
	case out.SuggestedKcalTarget != nil:
		out.DailyKcalTarget, out.TargetStatus = out.SuggestedKcalTarget, TargetSuggested
	}
	return out, nil
}
