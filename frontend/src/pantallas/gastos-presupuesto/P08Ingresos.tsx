// P-08 Ingresos (U1 §2.1, §2.2, §3.10: N-100 a N-105). Resumen · Detalle* · Tendencias · Calendario. Solo ingresos de Axel (decisión B).
import { useState } from 'react'
import { Plus, Trash2, Wallet } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { conConsulta, navegar } from '../../app/enrutador'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, dinero, fecha, hoyISO, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Grafica } from '../../componentes/Grafica'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { textos } from '../../textos/es-MX'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { Calendario } from './P06Presupuesto'
import { Bloque, Confirmar, HojaCaptura, SelectorMes, Selector, centavosATexto, mensajeDe, periodoDe, rejillaTarjetas, useSinRed } from './comun'
import { t } from './textos'

type Ingreso = Esquemas['Income']
type Quincena = { pay_date: string; status: 'expected' | 'captured'; expected_cents: number; income_id: number; amount_cents: number }

export default function P08Ingresos({ sub, consulta }: PropsPantalla) {
  const periodo = periodoDe(consulta)
  const im = usarConsulta('get', '/api/v1/incomes/month', { consulta: { period: periodo } })
  const sinRed = useSinRed()
  const avisar = useAvisos()
  const hoja = consulta.get('hoja')
  const [concepto, setConcepto] = useState('')
  const [fechaQ, setFechaQ] = useState('')
  const [borrar, setBorrar] = useState<number | null>(null)
  const abrir = (h: string | null) => navegar(conConsulta({ hoja: h }))
  const quincenas = (im.datos?.fortnights ?? []) as Quincena[]
  const esperadas = quincenas.filter((q) => q.status === 'expected')
  const elegida = esperadas.find((q) => q.pay_date === fechaQ) ?? esperadas[0]
  const parchar = async (id: number, cuerpo: Esquemas['IncomeUpdate']) => { await pedir('patch', '/api/v1/incomes/{income_id}', { ruta: { income_id: id }, cuerpo }); im.recargar() }
  const editarMonto = (id: (x: Ingreso) => number) => ({ tipo: 'dinero' as const, textoInicial: (x: Ingreso) => centavosATexto(x.amount_cents),
    vistaPrevia: async (x: Ingreso, v: string) => { const c = aCentavos(v); if (c == null || c <= 0) throw new Error(textos.registro.montoInvalido); return [{ etiqueta: t.comun.monto, antes: dinero(x.amount_cents), despues: dinero(c) }] },
    guardar: (x: Ingreso, v: string) => parchar(id(x), { amount_cents: aCentavos(v)! }) })
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-2">
        <SelectorMes valor={periodo} />
        <div className="ml-auto flex flex-wrap gap-2">
          <Boton deshabilitadoPor={sinRed} onClick={() => abrir('ingreso')}><Plus aria-hidden className="size-5" />{t.p08.extra}</Boton>
          <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => abrir('quincena')}><Wallet aria-hidden className="size-5" />{t.p08.capturar}</Boton>
        </div>
      </div>
      {sub === 'calendario' ? <Calendario periodo={periodo} tipos={['income']} clave="P-08-calendario" /> : sub === 'tendencias' ? <Tendencias /> : (
        <ConDatos consulta={im}>{(d) => {
          if (!quincenas.length && !d.extras.length) return <Vacio texto={t.p08.vacio} accion={<Boton onClick={() => navegar(rutaDe('P-26', { seccion: '03' }))}>{t.p08.irAjustes}</Boton>} />
          if (sub === 'resumen') return (
            <div className={rejillaTarjetas}>
              <TarjetaMetrica titulo={t.p08.capturado} rol="income" valor={dinero(d.captured_cents)} alAbrir={() => navegar(conConsulta({ sub: 'detalle' }))} />
              <TarjetaMetrica titulo={t.p08.esperadoPendiente} valor={dinero(d.expected_pending_cents)} alAbrir={() => navegar(conConsulta({ sub: 'detalle' }))} />
              <TarjetaMetrica titulo={t.p08.extras} rol="income" valor={dinero(d.extras_cents)} alAbrir={() => navegar(conConsulta({ sub: 'detalle' }))} />
              <TarjetaMetrica titulo={t.p08.total} valor={dinero(d.total_cents)} nota={mes(periodo)} />
            </div>)
          return (
            <div className="flex flex-col gap-6">
              <Bloque titulo={t.p08.quincenas}>
                <VistaBD def={{ clave: 'P-08-quincenas', filas: quincenas, idFila: (q) => q.pay_date, vistas: ['tabla'], titulo: (q) => t.p08.quincenaDel(fecha(q.pay_date)),
                  propiedades: [
                    { id: 'pay_date', nombre: t.comun.fecha, valor: (q) => q.pay_date, pintar: (q) => fecha(q.pay_date) },
                    { id: 'status', nombre: t.comun.estado, valor: (q) => t.p08.estados[q.status] },
                    { id: 'amount_cents', nombre: t.comun.monto, valor: (q) => (q.status === 'captured' ? q.amount_cents : q.expected_cents), alinear: 'derecha',
                      pintar: (q) => (q.status === 'captured' ? dinero(q.amount_cents, 'ingreso') : <span className="text-text-secondary">{dinero(q.expected_cents)}</span>) },
                    { id: 'accion', nombre: '', valor: () => '', pintar: (q) => (q.status === 'expected'
                      ? <Boton variante="enlace" deshabilitadoPor={sinRed} onClick={(e) => { e.stopPropagation(); setFechaQ(q.pay_date); abrir('quincena') }}>{t.p08.capturar}</Boton>
                      : <Boton tam="icono" variante="fantasma" aria-label={t.comun.eliminar} deshabilitadoPor={sinRed} onClick={(e) => { e.stopPropagation(); setBorrar(q.income_id) }}><Trash2 aria-hidden className="size-6" size={ICONO_MD} /></Boton>) },
                  ] }} />
              </Bloque>
              <Bloque titulo={t.p08.extras}>
                <VistaBD def={{ clave: 'P-08-extras', filas: d.extras, idFila: (x) => x.id, vistas: ['tabla', 'calendario'], fecha: (x) => x.income_date, titulo: (x) => x.description, vacio: <Vacio texto={t.p08.sinExtras} />,
                  propiedades: [
                    { id: 'income_date', nombre: t.comun.fecha, valor: (x) => x.income_date, pintar: (x) => fecha(x.income_date), editar: { tipo: 'fecha', textoInicial: (x) => x.income_date, guardar: (x, v) => parchar(x.id, { income_date: v }) } },
                    { id: 'description', nombre: t.comun.concepto, valor: (x) => x.description, editar: { tipo: 'texto', textoInicial: (x) => x.description, guardar: (x, v) => parchar(x.id, { description: v }) } },
                    { id: 'amount_cents', nombre: t.comun.monto, valor: (x) => x.amount_cents, pintar: (x) => dinero(x.amount_cents, 'ingreso'), alinear: 'derecha', editar: editarMonto((x) => x.id) },
                    { id: 'notes', nombre: t.comun.nota, valor: (x) => x.notes ?? '', editar: { tipo: 'texto', textoInicial: (x) => x.notes ?? '', guardar: (x, v) => parchar(x.id, { notes: v || null }) } },
                    { id: 'accion', nombre: '', valor: () => '', pintar: (x) => <Boton tam="icono" variante="fantasma" aria-label={t.comun.eliminar} deshabilitadoPor={sinRed} onClick={(e) => { e.stopPropagation(); setBorrar(x.id) }}><Trash2 aria-hidden className="size-6" size={ICONO_MD} /></Boton> },
                  ] }} />
              </Bloque>
            </div>)
        }}</ConDatos>)}
      <HojaCaptura titulo={t.p08.capturar} abierta={hoja === 'quincena'} alCerrar={() => abrir(null)} montoInicial={elegida ? centavosATexto(elegida.expected_cents) : ''} fechaInicial={elegida?.pay_date ?? hoyISO()}
        alGuardar={async (c) => { await pedir('post', '/api/v1/incomes', { idempotencia: nuevaClave(), cuerpo: { kind: 'fixed', description: t.p08.conceptoQuincena, amount_cents: c.monto, income_date: c.fecha, notes: c.nota || undefined } }); avisar({ texto: t.p08.capturada, rol: 'exito' }); im.recargar() }}
        extra={esperadas.length > 1 ? <Selector etiqueta={t.p08.quincenas} valor={elegida?.pay_date ?? ''} alCambiar={setFechaQ} opciones={esperadas.map((q) => ({ valor: q.pay_date, texto: t.p08.quincenaDel(fecha(q.pay_date)) }))} /> : undefined} />
      <HojaCaptura titulo={t.p08.extra} abierta={hoja === 'ingreso'} alCerrar={() => { abrir(null); setConcepto('') }} fechaInicial={hoyISO()}
        extra={<Campo etiqueta={t.comun.concepto} value={concepto} onChange={(e) => setConcepto(e.target.value)} />}
        alGuardar={async (c) => { if (!concepto.trim()) throw new Error(textos.errores.revisaCampos); await pedir('post', '/api/v1/incomes', { idempotencia: nuevaClave(), cuerpo: { kind: 'extra', description: concepto.trim(), amount_cents: c.monto, income_date: c.fecha, notes: c.nota || undefined } }); avisar({ texto: t.p08.registrado, rol: 'exito' }); im.recargar() }} />
      <Confirmar abierta={borrar !== null} alCambiar={(a) => !a && setBorrar(null)} titulo={t.comun.eliminar} texto={t.comun.confirmarEliminar} accion={t.comun.eliminar}
        alConfirmar={async () => { try { await pedir('delete', '/api/v1/incomes/{income_id}', { ruta: { income_id: borrar! } }); avisar({ texto: t.comun.eliminado, rol: 'exito' }); setBorrar(null); im.recargar() } catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) } }} />
    </div>
  )
}

function Tendencias() {
  const r = usarConsulta('get', '/api/v1/incomes/monthly-trend')
  return (
    <ConDatos consulta={r}>{(d) => d.months.length === 0 ? <Vacio texto={textos.vacio.grafica} /> : (
      <Grafica titulo={t.p08.tendencia} resumen={t.p08.resumenTendencia} tipo="barras" categorias={d.months.map((m) => mes(m.month ?? ''))}
        series={[{ nombre: t.p08.fijo, valores: d.months.map((m) => m.fixed_cents ?? 0) }, { nombre: t.p08.extras, valores: d.months.map((m) => m.extra_cents ?? 0) }]} />)}
    </ConDatos>
  )
}
