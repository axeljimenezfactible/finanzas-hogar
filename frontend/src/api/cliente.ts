// Cliente del API C1. Única puerta de red del frontend: pedir(metodo, ruta-del-contrato, opciones).
// Convenciones: errores C1 §3 → ErrorApi; sin respuesta → ErrorSinConexion; paginación C1 §5 (cursor);
// Idempotency-Key obligatoria por tipo en los POST de captura; dinero en centavos; fechas ISO.
import { ErrorApi, ErrorSinConexion } from './errores'
import { marcarConServidor, marcarSinServidor } from './conexion'
import { guardarLectura } from './lecturas'
import type { Consulta, Cuerpo, ItemsDe, Metodo, ParamsRuta, RequiereClave, Respuesta, RutasCon } from './tipos'

export type Opciones<P extends RutasCon<M>, M extends Metodo> = {
  ruta?: ParamsRuta<P, M>
  consulta?: Consulta<P, M>
  cuerpo?: Cuerpo<P, M>
  senal?: AbortSignal
} & (RequiereClave<P, M> extends true ? { idempotencia: string } : { idempotencia?: string })

/** Clave de idempotencia (16–128, [A-Za-z0-9_-]): UUID v4, o 32 hex si el navegador no tiene randomUUID. */
export function nuevaClave(): string {
  if (typeof crypto.randomUUID === 'function') return crypto.randomUUID()
  return [...crypto.getRandomValues(new Uint8Array(16))].map((b) => b.toString(16).padStart(2, '0')).join('')
}

export function armarUrl(plantilla: string, ruta?: object, consulta?: object): string {
  let url = plantilla.replace(/\{(\w+)\}/g, (_, k: string) => encodeURIComponent(String((ruta as Record<string, unknown>)?.[k] ?? '')))
  const q = new URLSearchParams()
  for (const [k, v] of Object.entries(consulta ?? {})) {
    if (v === undefined || v === null || v === '') continue
    q.set(k, Array.isArray(v) ? v.join(',') : String(v)) // style=form, explode=false (C1)
  }
  const s = q.toString()
  return s ? `${url}?${s}` : url
}

export async function pedir<M extends Metodo, P extends RutasCon<M>>(metodo: M, ruta: P, opciones?: Opciones<P, M>): Promise<Respuesta<P, M>> {
  const o = (opciones ?? {}) as { ruta?: object; consulta?: object; cuerpo?: unknown; senal?: AbortSignal; idempotencia?: string }
  const url = armarUrl(ruta as string, o.ruta, o.consulta)
  const encabezados: Record<string, string> = { Accept: 'application/json' }
  let cuerpo: BodyInit | undefined
  if (o.cuerpo instanceof FormData) cuerpo = o.cuerpo
  else if (o.cuerpo !== undefined) { cuerpo = JSON.stringify(o.cuerpo); encabezados['Content-Type'] = 'application/json' }
  if (o.idempotencia) encabezados['Idempotency-Key'] = o.idempotencia
  let res: Response
  try {
    res = await fetch(url, { method: metodo.toUpperCase(), headers: encabezados, body: cuerpo, credentials: 'same-origin', signal: o.senal, cache: 'no-store' })
  } catch (e) {
    if ((e as Error).name === 'AbortError') throw e
    marcarSinServidor()
    throw new ErrorSinConexion()
  }
  marcarConServidor()
  const tipo = res.headers.get('Content-Type') ?? ''
  const datos = res.status === 204 ? undefined : tipo.includes('json') ? await res.json() : await res.blob()
  if (!res.ok) {
    const e = (datos ?? {}) as { code?: string; message?: string; fields?: []; request_id?: string; retry_at?: string }
    if (res.status === 401 && ruta !== '/api/v1/auth/login') window.dispatchEvent(new Event('fh:sesion-vencida'))
    throw new ErrorApi(res.status, e.code ?? 'internal_error', e.message ?? 'Algo salió mal. Intenta de nuevo.', e.fields ?? [], e.request_id ?? res.headers.get('X-Request-Id'), e.retry_at ?? null)
  }
  if (metodo === 'get' && tipo.includes('json')) void guardarLectura(url, datos)
  return datos as Respuesta<P, M>
}

/** Una página (C1 §5): `limit` 50 por omisión, máximo 200 (K6); `next_cursor` null = no hay más. */
export function pedirPagina<P extends RutasCon<'get'>>(metodo: 'get', ruta: P, opciones?: Opciones<P, 'get'>) {
  return pedir(metodo, ruta, opciones)
}

/** Todas las páginas de una colección, en orden (para listas cortas: categorías, personas, tiendas). */
export async function todas<P extends RutasCon<'get'>>(metodo: 'get', ruta: P, opciones?: Opciones<P, 'get'>): Promise<ItemsDe<Respuesta<P, 'get'>>[]> {
  const items: ItemsDe<Respuesta<P, 'get'>>[] = []
  let cursor: string | null = null
  do {
    const consulta = { ...((opciones as { consulta?: object })?.consulta ?? {}), limit: 200, ...(cursor ? { cursor } : {}) }
    const pagina = (await pedir(metodo, ruta, { ...(opciones as object), consulta } as unknown as Opciones<P, 'get'>)) as unknown as { items: ItemsDe<Respuesta<P, 'get'>>[]; next_cursor: string | null }
    items.push(...pagina.items)
    cursor = pagina.next_cursor
  } while (cursor)
  return items
}
