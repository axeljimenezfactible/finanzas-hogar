// Pantallas P-00..P-26 (U1 §2.1) y su ruta en el navegador. Las rutas son PROVISIONALES
// (bloqueo docs/estado/bloqueos/WF3-frontend-rutas.md): U1 define pantallas, no direcciones.
import { textos } from '../textos/es-MX'
export type Subseccion = 'resumen' | 'detalle' | 'tendencias' | 'calendario'
export type IdSeccion = keyof typeof textos.secciones
export interface Pantalla { id: string; nombre: string; seccion: IdSeccion | null; ruta: string; subs?: Subseccion[]; sub?: Subseccion; padre?: string }
const R = 'resumen', D = 'detalle', T = 'tendencias', C = 'calendario'
export const PANTALLAS: Pantalla[] = [
  { id: 'P-00', nombre: 'Iniciar sesión', seccion: null, ruta: '/entrar' },
  { id: 'P-01', nombre: 'Inicio', seccion: 'inicio', ruta: '/inicio', subs: [R, D, T, C], sub: R },
  { id: 'P-05', nombre: 'Capturar ticket', seccion: 'gastos', ruta: '/gastos/ticket' },
  { id: 'P-04', nombre: 'Gasto', seccion: 'gastos', ruta: '/gastos/:id', padre: 'P-02' },
  { id: 'P-02', nombre: 'Todos los gastos', seccion: 'gastos', ruta: '/gastos', subs: [R, D, T, C], sub: D },
  { id: 'P-07', nombre: 'Flujo y proyección', seccion: 'presupuesto', ruta: '/presupuesto/flujo', subs: [R, D, T], sub: R },
  { id: 'P-08', nombre: 'Ingresos', seccion: 'presupuesto', ruta: '/presupuesto/ingresos', subs: [R, D, T, C], sub: D },
  { id: 'P-09', nombre: 'Ahorro y metas', seccion: 'presupuesto', ruta: '/presupuesto/ahorro', subs: [R, D, T, C], sub: R },
  { id: 'P-06', nombre: 'Presupuesto del mes', seccion: 'presupuesto', ruta: '/presupuesto', subs: [R, D, T, C], sub: R },
  { id: 'P-12', nombre: 'Receta', seccion: 'comidas', ruta: '/comidas/recetas/:id', padre: 'P-11' },
  { id: 'P-11', nombre: 'Recetario', seccion: 'comidas', ruta: '/comidas/recetas' },
  { id: 'P-10', nombre: 'Plan de la semana', seccion: 'comidas', ruta: '/comidas', subs: [R, D, T, C], sub: D },
  { id: 'P-14', nombre: 'Comparar canasta', seccion: 'compras', ruta: '/compras/comparar', subs: [R, D], sub: R },
  { id: 'P-15', nombre: 'Productos y precios', seccion: 'compras', ruta: '/compras/productos', subs: [R, D, T], sub: D },
  { id: 'P-16', nombre: 'Inventario', seccion: 'compras', ruta: '/compras/inventario', subs: [R, D], sub: D },
  { id: 'P-13', nombre: 'Lista de compras', seccion: 'compras', ruta: '/compras' },
  { id: 'P-20', nombre: 'Comparador de escenarios', seccion: 'creditos', ruta: '/creditos/comparador' },
  { id: 'P-19', nombre: 'Compra a crédito', seccion: 'creditos', ruta: '/creditos/:id/compras/:compra', padre: 'P-18' },
  { id: 'P-18', nombre: 'Crédito', seccion: 'creditos', ruta: '/creditos/:id', subs: [R, D, T, C], sub: R, padre: 'P-17' },
  { id: 'P-17', nombre: 'Todas las líneas', seccion: 'creditos', ruta: '/creditos', subs: [R, D, T, C], sub: R },
  { id: 'P-21', nombre: 'Salud', seccion: 'salud', ruta: '/salud', subs: [R, D, T], sub: R },
  { id: 'P-24', nombre: 'Reporte mensual', seccion: 'reportes', ruta: '/reportes/mensual' },
  { id: 'P-25', nombre: 'Estado de cuenta de Karen', seccion: 'reportes', ruta: '/reportes/karen' },
  { id: 'P-23', nombre: 'Centro de reportes', seccion: 'reportes', ruta: '/reportes' },
  { id: 'P-26', nombre: 'Ajustes', seccion: 'ajustes', ruta: '/ajustes/:seccion' },
  { id: 'P-26', nombre: 'Ajustes', seccion: 'ajustes', ruta: '/ajustes' },
]
/** Las 14 secciones de Ajustes (U1 §2.3), en orden. */
export const SECCIONES_AJUSTES = ['Cuenta y roles', 'Presupuesto', 'Ingresos', 'Categorías', 'Créditos y condiciones', 'Tiendas y precios', 'Recetas', 'Reportes y frecuencias', 'Notificaciones', 'Apariencia', 'Datos y respaldos', 'Personas y reparto', 'Salud y metas de ahorro', 'Inventario'].map((nombre, i) => ({ id: String(i + 1).padStart(2, '0'), nombre }))
/** Barra lateral (U1 §1.1–§1.2): orden fijo de las 9 secciones y sus subpáginas. */
export const SECCIONES: { id: IdSeccion; hijos: string[] }[] = [
  { id: 'inicio', hijos: [] },
  { id: 'gastos', hijos: ['P-02', 'P-05'] },
  { id: 'presupuesto', hijos: ['P-06', 'P-07', 'P-08', 'P-09'] },
  { id: 'comidas', hijos: ['P-10', 'P-11'] },
  { id: 'compras', hijos: ['P-13', 'P-14', 'P-15', 'P-16'] },
  { id: 'creditos', hijos: ['P-17', 'P-20'] },
  { id: 'salud', hijos: [] },
  { id: 'reportes', hijos: ['P-23', 'P-24', 'P-25'] },
  { id: 'ajustes', hijos: [] },
]
export const PRIMERA: Record<IdSeccion, string> = { inicio: 'P-01', gastos: 'P-02', presupuesto: 'P-06', comidas: 'P-10', compras: 'P-13', creditos: 'P-17', salud: 'P-21', reportes: 'P-23', ajustes: 'P-26' }
const compiladas = PANTALLAS.map((p) => ({ p, re: new RegExp('^' + p.ruta.replace(/:(\w+)/g, '(?<$1>[^/]+)') + '/?$') }))
export function resolver(ruta: string): { pantalla: Pantalla; params: Record<string, string> } | null {
  for (const { p, re } of compiladas) { const m = re.exec(ruta); if (m) return { pantalla: p, params: { ...m.groups } } }
  return null
}
export const pantalla = (id: string) => PANTALLAS.find((p) => p.id === id)!
export function rutaDe(id: string, params: Record<string, string | number> = {}): string {
  return pantalla(id).ruta.replace(/:(\w+)/g, (_, k: string) => encodeURIComponent(String(params[k] ?? '')))
}
