// D-7 Participación de Karen (U3 §3.1, §9): solo lo que comparte; Karen no tiene cuenta (decisiones §A).
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia } from '../../api/formatos'
import { rutaDe } from '../../app/rutas'
import { MarcaPersona } from '../../componentes/Chip'
import { EnlaceAccion, EsqueletoFilas, Estados, IconoOk, Insignia, Panel, karen, useInicio } from './comun'
import { textosInicio as t } from './textos'
import { CircleDashed } from 'lucide-react'

export default function Karen({ aviso }: { aviso?: string[] }) {
  const c = useInicio()
  const k = karen(c)
  const q = usarConsulta('get', '/api/v1/persons/{person_id}/participation', { ruta: { person_id: k?.id ?? 0 } }, k !== undefined)
  const th = 'py-2 pr-3 text-left text-caption font-medium text-text-secondary', td = 'py-2 pr-3 tabular-nums'
  return (
    <Panel etiqueta={t.d7.titulo} className="desktop:col-span-2" titulo={<><MarcaPersona codigo="karen" />{t.d7.titulo}</>}>
      {aviso?.map((a) => <p key={a} className="text-caption text-text-secondary">{a}</p>)}
      <Estados consulta={q} falla={t.d7.falla} esqueleto={<EsqueletoFilas n={4} />}>{(d) => {
        const sinNada = d.purchases.length === 0 && d.payment_shares.length === 0 && d.accrued_groceries_cents === 0
        if (sinNada) return <p className="text-body text-text-secondary">{t.d7.vacio}</p>
        return (
          <div className="flex flex-col gap-4">
            <p className="text-body tabular-nums">{t.d7.despensa(dinero(d.accrued_groceries_cents))}</p>
            {d.purchases.length > 0 && <div className="overflow-x-auto"><table className="w-full text-small"><caption className="mb-1 text-left text-small font-medium text-text-default">{t.d7.compras}</caption>
              <thead><tr><th scope="col" className={th}>{t.d7.comercio}</th><th scope="col" className={th}>{t.d7.fecha}</th><th scope="col" className={th + ' text-right'}>{t.d7.total}</th><th scope="col" className={th}>{t.d7.pagosDe(0, 0).split(' ')[0]}</th><th scope="col" className={th + ' text-right'}>{t.d7.saldo}</th></tr></thead>
              <tbody>{d.purchases.map((p) => <tr key={p.credit_purchase_id} className="border-t border-border-divider"><td className={td}>{p.description ?? t.comun.sinDato}</td><td className={td}>{fecha(p.purchase_date)}</td><td className={td + ' text-right'}>{dinero(p.total_cents)}</td><td className={td}>{t.d7.pagosDe(p.payments_made, p.payments_total)}</td><td className={td + ' text-right'}>{dinero(p.balance_cents)}</td></tr>)}</tbody>
            </table></div>}
            {d.payment_shares.length > 0 && <div className="flex flex-col gap-2"><h3 className="text-small font-medium">{t.d7.pagos}</h3>
              <ul className="flex flex-col">{d.payment_shares.map((p) => <li key={p.scheduled_payment_id} className="flex flex-wrap items-center justify-between gap-2 border-t border-border-divider py-2 text-small tabular-nums">
                <span>{fechaConDia(p.due_date)}</span><span>{dinero(p.amount_cents)}</span>
                {p.status === 'paid' ? <Insignia rol="exito" icono={<IconoOk />}>{t.pago.pagado}</Insignia> : <Insignia rol="neutro" icono={<CircleDashed aria-hidden />}>{t.pago.pendiente}</Insignia>}
              </li>)}</ul></div>}
            <p className="text-body tabular-nums">{t.d7.debe(dinero(d.live_debt_cents))}</p>
          </div>
        )
      }}</Estados>
      <p className="text-small text-text-secondary">{t.d7.nota}</p>
      <EnlaceAccion a={rutaDe('P-25')} texto={t.d7.accion} />
    </Panel>
  )
}
