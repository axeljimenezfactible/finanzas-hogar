// Edición en grupo con Ctrl+/ (U4 §0.2 punto 2 y §0.7; pieza de la base `def.editarVarias`, H-07): apoyo común de 02, 04, 05 y 07.
import type { Cambio } from '../../vistas-bd/tipos'
import { t } from './textos'
/** El bloque «Antes → Después» de la base solo pinta renglones `Cambio`; un aviso de U4 sin cifra va como renglón «— → —».
 *  PROVISIONAL (bloqueo WF3a-bis-ajustes-a §1): mientras la base no acepte notas en `vistaPrevia`. */
export const nota = (texto: string): Cambio => ({ etiqueta: texto, antes: '—', despues: '—' })
/** Valor «Antes» de varias filas: el valor si todas lo comparten; si no, cuántos valores distintos hay. */
export const resumen = (valores: string[]) => { const u = [...new Set(valores)]; return u.length === 1 ? u[0] : t.grupo.varios(u.length) }
export const siNo = [{ valor: 'si', texto: t.comun.si }, { valor: 'no', texto: t.comun.no }]
/** Un PATCH por fila, igual que al editar una sola (U4 P-26.02 punto h). La que falla conserva su valor anterior; al final se
 *  recarga y, si alguna falló, se lanza el primer error para que el bloque lo muestre y siga abierto. */
export async function guardarVarias<T>(filas: T[], una: (fila: T) => Promise<unknown>, recargar: () => void) {
  let error: unknown = null
  for (const f of filas) { try { await una(f) } catch (e) { error ??= e } }
  recargar()
  if (error) throw error
}
