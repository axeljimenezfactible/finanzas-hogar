// P-11 Recetario (U1 §2.1–§2.2, §3.13 N-135 a N-137): galería*, tabla y tablero por tiempo de comida.
import { useDeferredValue, useState } from 'react'
import { BookOpen, Sparkles } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import type { Esquemas } from '../../api/tipos'
import { pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { aCentavos, dinero } from '../../api/formatos'
import { navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ChipSeleccionable } from '../../componentes/Chip'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { useAvisos } from '../../componentes/Avisos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { BarraAcciones, MenuMas, Selector, filasDe, kcal, mensajeError, usarCatalogo } from './comun'
import { guardarBorrador } from './P12Receta'
import { textosCC as t } from './textos'

type Receta = Esquemas['Recipe']

export default function P11Recetario(_: PropsPantalla) {
  const conexion = useConexion()
  const avisar = useAvisos()
  const tiempos = [...(usarCatalogo('tiempos') ?? [])].sort((a, b) => a.sort_order - b.sort_order)
  const [texto, setTexto] = useState('')
  const [tiempo, setTiempo] = useState('')
  const [origen, setOrigen] = useState('')
  const [costoMax, setCostoMax] = useState('')
  const [kcalMax, setKcalMax] = useState('')
  const [ia, setIa] = useState(false)
  const q = useDeferredValue(texto.trim())
  const c = usarConsulta('get', '/api/v1/recipes', { consulta: { q: q || undefined, meal_slot_id: tiempo ? Number(tiempo) : undefined, origin: (origen || undefined) as Receta['origin'] | undefined, limit: 200 } })
  const sinRed = conexion ? undefined : textos.conexion.requiere
  const tope = aCentavos(costoMax)
  const topeKcal = Number(kcalMax) || null
  const nombresTiempos = (r: Receta) => (r.meal_slot_ids ?? []).map((id) => tiempos.find((x) => x.id === id)?.name).filter(Boolean).join(', ')
  return (
    <div className="flex flex-col gap-4">
      <BarraAcciones acciones={<>
        <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => navegar(rutaDe('P-12', { id: 'nueva' }))}>{t.nuevaReceta}</Boton>
        <MenuMas acciones={[{ texto: t.sugerirIa, alHacer: () => setIa(true), deshabilitadoPor: sinRed }]} />
      </>}>
        <Campo etiqueta={t.buscar} ocultarEtiqueta type="search" placeholder={t.buscarReceta} value={texto} onChange={(e) => setTexto(e.target.value)} />
      </BarraAcciones>
      <div className="flex flex-col gap-3 tablet:flex-row tablet:flex-wrap tablet:items-end">
        <div role="group" aria-label={t.tiempo} className="flex flex-wrap gap-2">
          <ChipSeleccionable seleccionado={!tiempo} onClick={() => setTiempo('')}>{t.todos}</ChipSeleccionable>
          {tiempos.map((x) => <ChipSeleccionable key={x.id} seleccionado={tiempo === String(x.id)} onClick={() => setTiempo(String(x.id))}>{x.name}</ChipSeleccionable>)}
        </div>
        <Selector id="p11-origen" etiqueta={t.origen} valor={origen} alCambiar={setOrigen} vacio={t.todos} opciones={Object.entries(t.origenes).map(([valor, texto]) => ({ valor, texto }))} />
        <Campo etiqueta={t.costoMax} inputMode="decimal" value={costoMax} onChange={(e) => setCostoMax(e.target.value)} />
        <Campo etiqueta={t.kcalMax} inputMode="numeric" value={kcalMax} onChange={(e) => setKcalMax(e.target.value)} />
      </div>
      <ConDatos consulta={c}>{(d) => {
        const filas = filasDe<Receta>(d).filter((r) => (tope == null || (r.cost_per_serving_cents ?? 0) <= tope) && (topeKcal == null || (r.kcal_per_serving ?? 0) <= topeKcal))
        if (!filas.length) return <Vacio icono={<BookOpen aria-hidden className="size-6" size={ICONO_MD} />} texto={t.sinRecetas} />
        return (
          <VistaBD def={{
            clave: 'P-11-recetas', filas, idFila: (r) => r.id, vistas: ['galeria', 'tabla', 'tablero'], titulo: (r) => r.name,
            alAbrir: (r) => navegar(rutaDe('P-12', { id: r.id })),
            columnasTablero: { propiedad: 'tiempos', columnas: tiempos.map((x) => x.name) },
            propiedades: [
              { id: 'nombre', nombre: t.nombre, valor: (r) => r.name, editar: { tipo: 'texto', textoInicial: (r) => r.name, guardar: async (r, v) => { await pedir('patch', '/api/v1/recipes/{recipe_id}', { ruta: { recipe_id: r.id }, cuerpo: { name: v } }); c.recargar() } } },
              { id: 'tiempos', nombre: t.tiempos, valor: nombresTiempos, grupo: (r) => tiempos.find((x) => x.id === r.meal_slot_ids?.[0])?.name ?? textos.vistas.sinValor },
              { id: 'porciones', nombre: t.porcionesBase, valor: (r) => r.base_servings, alinear: 'derecha' },
              { id: 'tiempoTotal', nombre: t.tiempoTotal, valor: (r) => (r.prep_minutes ?? 0) + (r.cook_minutes ?? 0), pintar: (r) => t.minutos((r.prep_minutes ?? 0) + (r.cook_minutes ?? 0)), alinear: 'derecha' },
              { id: 'costo', nombre: t.costoPorcion, valor: (r) => r.cost_per_serving_cents ?? null, pintar: (r) => (r.cost_per_serving_cents != null ? dinero(r.cost_per_serving_cents) : '—'), alinear: 'derecha' },
              { id: 'kcal', nombre: t.kcalPorcion, valor: (r) => r.kcal_per_serving ?? null, pintar: (r) => kcal(r.kcal_per_serving), alinear: 'derecha' },
              { id: 'origen', nombre: t.origen, valor: (r) => t.origenes[r.origin] },
              { id: 'usada', nombre: t.usada8, valor: (r) => r.times_used_8w ?? 0, pintar: (r) => t.veces(r.times_used_8w ?? 0), alinear: 'derecha' },
            ],
          }} />
        )
      }}</ConDatos>
      {ia && <HojaIa cerrar={() => setIa(false)} alError={(e) => avisar({ texto: mensajeError(e), rol: 'error' })} />}
    </div>
  )
}

/** Sugerir receta con IA (N-137): borrador que se revisa en P-12 antes de guardar. Si la IA está apagada, C1 responde 409 y se muestra su mensaje. */
function HojaIa({ cerrar, alError }: { cerrar: () => void; alError: (e: unknown) => void }) {
  const [desc, setDesc] = useState('')
  const [ocupado, setOcupado] = useState(false)
  async function pedirIa() {
    setOcupado(true)
    try {
      const r = await pedir('post', '/api/v1/recipes/ai-suggestions', { cuerpo: { description: desc.trim() } })
      guardarBorrador(r)
      navegar(rutaDe('P-12', { id: 'nueva' }))
    } catch (e) { alError(e) } finally { setOcupado(false) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.sugerirIa} pie={<Boton tam="hoja" variante="primario" disabled={ocupado || !desc.trim()} onClick={pedirIa}><Sparkles aria-hidden className="size-5" />{t.sugerirIa}</Boton>}>
      <Campo etiqueta={t.describeIa} value={desc} onChange={(e) => setDesc(e.target.value)} />
    </Hoja>
  )
}
