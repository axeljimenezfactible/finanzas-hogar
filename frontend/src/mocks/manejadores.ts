// Manejadores de msw para TODAS las operaciones de C1 (solo «npm run dev:mocks»; nunca en el build).
// Genéricos: la respuesta de éxito de cada operación con los ejemplos del contrato (operaciones.gen.json).
// Con estado: sesión, personas, categorías, gastos y pesos (idempotentes por Idempotency-Key, K9).
import { http, HttpResponse, type HttpHandler } from 'msw'
import operaciones from './operaciones.gen.json'
import sinteticos from './sinteticos.gen.json'
type Op = { id: string; metodo: string; ruta: string; estado: number; ejemplo: unknown; binario: boolean }
const ahora = () => '2026-10-21T09:00:00-06:00'
const sesion = { user_account_id: 1, username: 'axel', person_id: 1, display_name: 'Axel', participation: 'holder', device_session_id: 1, expires_at: '2027-01-19T09:00:00-06:00', last_seen_at: ahora() }
const pagina = (items: unknown[]) => HttpResponse.json({ items, next_cursor: null })
const gastos: Record<string, unknown>[] = []
const pesos: Record<string, unknown>[] = []
const porClave = new Map<string, Record<string, unknown>>()
function captura(lista: Record<string, unknown>[]) {
  return async ({ request }: { request: Request }) => {
    const clave = request.headers.get('Idempotency-Key')
    if (!clave || clave.length < 16 || clave.length > 128) return HttpResponse.json({ code: 'validation_failed', message: 'Falta la clave de idempotencia.', request_id: 'mock' }, { status: 422 })
    const previo = porClave.get(clave)
    if (previo) return HttpResponse.json(previo, { status: 200, headers: { 'Idempotent-Replayed': 'true' } })
    const nuevo = { id: lista.length + 1, ...(await request.json() as object), created_at: ahora(), updated_at: ahora() }
    lista.unshift(nuevo)
    porClave.set(clave, nuevo)
    return HttpResponse.json(nuevo, { status: 201 })
  }
}
const conEstado: HttpHandler[] = [
  http.post('/api/v1/auth/login', () => HttpResponse.json(sesion)),
  http.get('/api/v1/auth/me', () => HttpResponse.json(sesion)),
  http.post('/api/v1/auth/logout', () => new HttpResponse(null, { status: 204 })),
  http.get('/api/v1/persons', () => pagina(sinteticos.personas)),
  http.get('/api/v1/categories', () => pagina(sinteticos.categorias)),
  http.get('/api/v1/expenses', () => pagina(gastos)),
  http.post('/api/v1/expenses', captura(gastos)),
  http.get('/api/v1/weight-entries', () => pagina(pesos)),
  http.post('/api/v1/weight-entries', captura(pesos)),
  http.get('/health', () => HttpResponse.json({ status: 'ok' })),
]
const genericos = (operaciones as Op[]).map((o) =>
  http[o.metodo as 'get' | 'post' | 'put' | 'patch' | 'delete'](o.ruta.replace(/\{(\w+)\}/g, ':$1'), () =>
    o.estado === 204 ? new HttpResponse(null, { status: 204 }) : o.binario ? new HttpResponse('%PDF-1.7 ejemplo', { status: o.estado, headers: { 'Content-Type': 'application/pdf' } }) : HttpResponse.json(o.ejemplo as object, { status: o.estado }),
  ),
)
// Cada área puede sumar manejadores con estado en src/pantallas/<area>/mocks.ts (export const manejadores); van primero.
const deAreas = Object.values(import.meta.glob<{ manejadores: HttpHandler[] }>('../pantallas/*/mocks.ts', { eager: true })).flatMap((m) => m.manejadores)
export const manejadores = [...deAreas, ...conEstado, ...genericos]
