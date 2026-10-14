// Última carga de cada lectura (GET) en el teléfono: «lo ya cargado sigue visible con su hora de carga» (U1 §2.4).
// Se borra al cerrar sesión. Las escrituras nunca se guardan aquí (solo gastos y peso van a la cola, K9).
import { leer, poner, quitar } from '../cola/idb'
export interface Lectura { url: string; datos: unknown; cargadoEn: number }
export async function guardarLectura(url: string, datos: unknown) {
  try { await poner<Lectura>('lecturas', { url, datos, cargadoEn: Date.now() }) } catch { /* sin IndexedDB: se sigue sin copia */ }
}
export async function leerLectura(url: string): Promise<Lectura | undefined> {
  try { return await leer<Lectura>('lecturas', url) } catch { return undefined }
}
export const borrarLecturas = () => quitar('lecturas').catch(() => undefined)
