// G-03 menú «/» y G-04 búsqueda (U1 §5): un panel; si el texto empieza con «/», comandos; si no, búsqueda (GET /api/v1/search).
import { useEffect, useMemo, useState } from 'react'
import { Command } from 'lucide-react'
import { Dialog } from 'radix-ui'
import { textos } from '../textos/es-MX'
import { usarConsulta } from '../api/consulta'
import { conConsulta, navegar } from '../app/enrutador'
import { PANTALLAS, SECCIONES_AJUSTES, rutaDe } from '../app/rutas'
import { useSesion } from '../app/sesion'
import { leerTema, ponerTema, siguienteTema } from '../app/tema'
import { estilosControl } from '../componentes/Campo'
import { cn } from '../componentes/utils'
import { abrirRegistro, cambiar, useShell } from './estado'
interface Comando { cmd: string; alias: string; texto: string; hacer: () => void }
const ir = (id: string, hoja?: string) => () => navegar(rutaDe(id, { seccion: '01' }) + (hoja ? `?hoja=${hoja}` : ''))
const norm = (s: string) => s.normalize('NFD').replace(/\p{Diacritic}/gu, '').toLowerCase()
const DESTINO: Record<string, (id: string) => string> = { expense: (id) => rutaDe('P-04', { id }), recipe: (id) => rutaDe('P-12', { id }), credit: (id) => rutaDe('P-18', { id }), settings_section: (id) => `/ajustes/${id}`, income: () => rutaDe('P-08'), product: (id) => `${rutaDe('P-15')}?producto=${id}`, store: () => '/ajustes/06', credit_purchase: () => rutaDe('P-17'), savings_goal: () => rutaDe('P-09') }
const GRUPO: Record<string, string> = { expense: 'Gastos', income: 'Ingresos', recipe: 'Recetas', product: 'Productos', store: 'Tiendas', credit: 'Créditos', credit_purchase: 'Compras a crédito', savings_goal: 'Metas de ahorro', settings_section: 'Ajustes' }
export function MenuComandos() {
  const { menu } = useShell()
  const { salir } = useSesion()
  const [texto, setTexto] = useState('')
  useEffect(() => { if (menu !== null) setTexto(menu) }, [menu])
  const comandos: Comando[] = useMemo(() => [
    { cmd: 'gasto', alias: 'g gastar', texto: 'Registrar gasto', hacer: () => abrirRegistro('gasto') },
    { cmd: 'ticket', alias: 't foto ocr', texto: 'Capturar ticket', hacer: ir('P-05') },
    { cmd: 'ingreso', alias: 'extra', texto: 'Registrar ingreso extra', hacer: ir('P-08', 'ingreso') },
    { cmd: 'quincena', alias: 'q sueldo', texto: 'Capturar quincena de Axel', hacer: ir('P-08', 'quincena') },
    { cmd: 'credito', alias: 'kueski msi compra a credito', texto: 'Registrar compra a crédito', hacer: ir('P-17', 'compra') },
    { cmd: 'pago', alias: 'abono', texto: 'Marcar pago de crédito', hacer: ir('P-17', 'pago') },
    { cmd: 'ahorro', alias: 'aportar retiro', texto: 'Registrar movimiento de ahorro', hacer: ir('P-09', 'movimiento') },
    { cmd: 'peso', alias: 'p bascula', texto: 'Registrar peso', hacer: () => abrirRegistro('peso') },
    { cmd: 'precio', alias: '', texto: 'Capturar precio', hacer: ir('P-15', 'precio') },
    { cmd: 'receta', alias: '', texto: 'Nueva receta', hacer: () => navegar(rutaDe('P-12', { id: 'nueva' })) },
    { cmd: 'comida', alias: 'plan menu', texto: 'Asignar receta al plan', hacer: ir('P-10', 'asignar') },
    { cmd: 'lista', alias: 'comprar', texto: 'Agregar artículo a la lista de compras', hacer: ir('P-13', 'articulo') },
    { cmd: 'existencia', alias: 'stock sobrante', texto: 'Ajustar existencia o registrar sobrante', hacer: ir('P-16', 'existencia') },
    { cmd: 'promocion', alias: 'promo', texto: 'Nueva promoción', hacer: () => navegar('/ajustes/05?hoja=promocion') },
    { cmd: 'comparar', alias: 'canasta', texto: 'Comparar canasta', hacer: ir('P-14') },
    { cmd: 'simular', alias: 'escenarios', texto: 'Comparador de escenarios', hacer: ir('P-20') },
    { cmd: 'reporte', alias: 'pdf', texto: 'Generar reporte mensual', hacer: ir('P-23', 'generar') },
    { cmd: 'karen', alias: 'estado de cuenta', texto: 'Estado de cuenta de Karen', hacer: ir('P-25') },
    { cmd: 'respaldo', alias: 'backup', texto: 'Respaldar ahora', hacer: () => navegar('/ajustes/11?hoja=respaldo') },
    ...(['tabla', 'tablero', 'galeria'] as const).map((v) => ({ cmd: v, alias: '', texto: `Ver como ${textos.vistas[v].toLowerCase()}`, hacer: () => navegar(conConsulta({ vista: v })) })),
    ...(['resumen', 'detalle', 'tendencias', 'calendario'] as const).map((s) => ({ cmd: s, alias: '', texto: `Ir a ${textos.subsecciones[s]}`, hacer: () => navegar(conConsulta({ sub: s })) })),
    ...(['filtrar', 'ordenar', 'agrupar'] as const).map((v) => ({ cmd: v, alias: '', texto: textos.vistas[v], hacer: () => navegar(conConsulta({ controles: v })) })),
    { cmd: 'tema', alias: 'claro oscuro', texto: textos.tema.cambiar, hacer: () => ponerTema(siguienteTema(leerTema())) },
    { cmd: 'barra', alias: '', texto: 'Mostrar u ocultar la barra lateral', hacer: () => cambiar({ barraPlegada: !(document.querySelector('aside') !== null) }) },
    { cmd: 'atajos', alias: 'ayuda', texto: 'Ver la lista de atajos', hacer: () => navegar(conConsulta({ hoja: 'atajos' })) },
    { cmd: 'salir', alias: 'cerrar sesion', texto: textos.barra.cerrarSesion, hacer: () => void salir() },
    ...PANTALLAS.filter((p) => p.seccion && !p.ruta.includes(':')).map((p) => ({ cmd: 'ir', alias: norm(p.nombre), texto: `Ir a ${p.nombre}`, hacer: () => navegar(p.ruta) })),
    ...SECCIONES_AJUSTES.map((s) => ({ cmd: 'ir', alias: norm(`ajustes ${s.nombre}`), texto: `Ir a Ajustes › ${s.nombre}`, hacer: () => navegar(`/ajustes/${s.id}`) })),
  ], [salir])
  const esComando = texto.startsWith('/')
  const q = norm(texto.replace(/^\//, '').trim())
  const filtrados = esComando ? comandos.filter((c) => !q || norm(`${c.cmd} ${c.alias} ${c.texto}`).includes(q)).slice(0, 12) : []
  const busqueda = usarConsulta('get', '/api/v1/search', { consulta: { q: texto.trim() } }, !esComando && texto.trim().length >= 2)
  const cerrar = () => cambiar({ menu: null })
  const ejecutar = (f: () => void) => { cerrar(); f() }
  return (
    <Dialog.Root open={menu !== null} onOpenChange={(a) => !a && cerrar()}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-overlay-scrim" />
        <Dialog.Content className="fixed inset-x-0 top-0 z-50 flex max-h-dvh flex-col gap-2 bg-surface-raised p-3 tablet:top-24 tablet:left-1/2 tablet:max-w-xl tablet:-translate-x-1/2 tablet:rounded-lg tablet:border tablet:border-border-divider tablet:shadow-2">
          <Dialog.Title className="sr-only">{textos.menu.titulo}</Dialog.Title>
          <Dialog.Description className="sr-only">{textos.menu.placeholder}</Dialog.Description>
          <div className="flex items-center gap-2"><Command aria-hidden className="size-5 text-icon-default" />
            <input autoFocus value={texto} onChange={(e) => setTexto(e.target.value)} placeholder={textos.menu.placeholder} aria-label={textos.menu.titulo} className={cn(estilosControl, 'flex-1')}
              onKeyDown={(e) => { if (e.key === 'Enter' && filtrados[0]) { e.preventDefault(); ejecutar(filtrados[0].hacer) } }} />
          </div>
          <ul role="listbox" aria-label={esComando ? textos.menu.comandos : textos.menu.resultados} className="flex flex-col overflow-y-auto">
            {filtrados.map((c, i) => (
              <li key={`${c.cmd}-${c.texto}`} role="option" aria-selected={i === 0}>
                <button type="button" onClick={() => ejecutar(c.hacer)} className="flex h-9 w-full items-center gap-2 rounded-sm px-2 text-left text-small hover:bg-surface-hover coarse:h-11">
                  <span className="w-24 shrink-0 font-medium text-text-secondary">/{c.cmd}</span>{c.texto}
                </button>
              </li>
            ))}
            {!esComando && busqueda.datos?.groups?.map((g) => g.items.map((it) => (
              <li key={`${g.type}-${it.id}`} role="option" aria-selected={false}>
                <button type="button" onClick={() => ejecutar(() => navegar(DESTINO[g.type]?.(it.id) ?? '/inicio'))} className="flex w-full flex-col rounded-sm px-2 py-1 text-left text-small hover:bg-surface-hover">
                  <span><span className="text-text-secondary">{GRUPO[g.type]} · </span>{it.title}</span>{it.subtitle && <span className="text-caption text-text-secondary">{it.subtitle}</span>}
                </button>
              </li>
            )))}
            {((esComando && !filtrados.length) || (!esComando && busqueda.datos && !busqueda.datos.groups?.some((g) => g.items.length))) && <li className="px-2 py-3 text-small text-text-secondary">{textos.menu.sinResultados}</li>}
          </ul>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
