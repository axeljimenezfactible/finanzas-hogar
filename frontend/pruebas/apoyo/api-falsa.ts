// API falsa para las pruebas base: context.route sobre /api/v1 y /health, con estado en memoria.
// Respeta C1: Page {items,next_cursor}, Error {code,message,request_id}, Idempotency-Key (200 + Idempotent-Replayed).
import type { BrowserContext, Request } from '@playwright/test'
export const SESION = { user_account_id: 1, username: 'axel', person_id: 1, display_name: 'Axel', participation: 'holder', device_session_id: 7, expires_at: '2027-01-10T09:00:00-06:00', last_seen_at: '2026-10-12T09:00:00-06:00' }
const t = '2026-10-01T09:00:00-06:00'
export const CATEGORIAS = ['Despensa', 'Comida fuera', 'Transporte', 'Servicios', 'Salud', 'Hogar', 'Ropa'].map((name, i) => ({ id: i + 1, name, sort_order: i + 1, is_active: true, recent_usage_count: 7 - i, created_at: t, updated_at: t }))
export const PERSONAS = [['axel', 'Axel', 'holder'], ['ximena', 'Ximena', 'co_holder'], ['karen', 'Karen', 'external']].map(([code, display_name, participation], i) => ({ id: i + 1, code, display_name, participation, pays_share_directly: false, is_active: true, created_at: t, updated_at: t }))
export interface Captura { clave: string; cuerpo: Record<string, unknown> }
export class ApiFalsa {
  conSesion = true
  gastos: Captura[] = []
  pesos: Captura[] = []
  peticiones: { metodo: string; url: string; encabezados: Record<string, string>; cuerpo: string | null }[] = []
  fallarSiguienteGasto = false
  /** Mientras sea true, cada POST de gasto responde 500 (servidor con falla), sin guardar nada. */
  fallaServidorGastos = false
  limitarLogin = false
  async instalar(ctx: BrowserContext) {
    await ctx.route(/\/(api\/v1|health)(\/|\?|$)/, async (route) => {
      const pet = route.request()
      const u = new URL(pet.url())
      this.peticiones.push({ metodo: pet.method(), url: u.pathname + u.search, encabezados: pet.headers(), cuerpo: pet.postData() })
      const json = (estado: number, cuerpo: unknown, enc: Record<string, string> = {}) => route.fulfill({ status: estado, contentType: 'application/json', headers: { 'X-Request-Id': 'prueba', ...enc }, body: JSON.stringify(cuerpo) })
      const error = (estado: number, code: string, message: string, extra = {}) => json(estado, { code, message, request_id: 'prueba', ...extra })
      const ruta = `${pet.method()} ${u.pathname}`
      if (ruta === 'GET /health') return json(200, { status: 'ok' })
      if (ruta === 'POST /api/v1/auth/login') {
        if (this.limitarLogin) return error(429, 'too_many_attempts', 'Demasiados intentos.', { retry_at: '2026-10-12T21:45:00-06:00' })
        const b = pet.postDataJSON()
        if (b.username === 'axel' && b.password === 'contrasena-de-prueba') { this.conSesion = true; return json(200, SESION) }
        return error(422, 'invalid_credentials', 'Usuario o contraseña incorrectos.')
      }
      if (!this.conSesion) return error(401, 'unauthorized', 'Inicia sesión.')
      if (ruta === 'GET /api/v1/auth/me') return json(200, SESION)
      if (ruta === 'POST /api/v1/auth/logout') { this.conSesion = false; return route.fulfill({ status: 204 }) }
      if (ruta === 'GET /api/v1/categories') return json(200, { items: CATEGORIAS, next_cursor: null })
      if (ruta === 'GET /api/v1/persons') return json(200, { items: PERSONAS, next_cursor: null })
      if (ruta === 'POST /api/v1/expenses' && this.fallaServidorGastos) return error(500, 'internal_error', 'Ocurrió un error en el servidor de casa.')
      if (ruta === 'POST /api/v1/expenses') return this.capturar(pet, this.gastos, json, error, () => this.fallarSiguienteGasto && !(this.fallarSiguienteGasto = false))
      if (ruta === 'POST /api/v1/weight-entries') return this.capturar(pet, this.pesos, json, error, () => false)
      if (pet.method() === 'GET') return json(200, { items: [], next_cursor: null })
      return error(501, 'not_implemented', 'No implementado en la API falsa.')
    })
  }
  private capturar(pet: Request, lista: Captura[], json: (e: number, c: unknown, h?: Record<string, string>) => Promise<void>, error: (e: number, c: string, m: string) => Promise<void>, fallar: () => boolean) {
    const clave = pet.headers()['idempotency-key'] ?? ''
    if (clave.length < 16 || clave.length > 128 || !/^[A-Za-z0-9_-]+$/.test(clave)) return error(400, 'bad_request', 'Falta la clave de idempotencia.')
    if (fallar()) return error(422, 'validation_failed', 'La categoría ya no existe.')
    const cuerpo = pet.postDataJSON()
    const previo = lista.findIndex((c) => c.clave === clave)
    if (previo >= 0) return json(200, { id: previo + 1, ...lista[previo].cuerpo, created_at: t, updated_at: t }, { 'Idempotent-Replayed': 'true' })
    lista.push({ clave, cuerpo })
    return json(201, { id: lista.length, ...cuerpo, created_at: t, updated_at: t })
  }
}
