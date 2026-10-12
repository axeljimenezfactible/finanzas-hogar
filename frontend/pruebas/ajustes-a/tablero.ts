// Apoyo de pruebas del tablero de Ajustes (U4 §0.2 punto 2): arrastre con teclado. Lo usan ajustes-a y ajustes-b.
import type { Page } from '@playwright/test'
/** Arrastra con teclado (dnd-kit: Espacio toma, flechas mueven 25 px, Espacio suelta) una tarjeta hasta la columna indicada.
 *  Avanza flecha por flecha hasta que la región en vivo de dnd-kit anuncia la columna (cada paso se calcula con el rectángulo ya pintado). */
export async function moverTarjetaA(page: Page, tarjeta: string, columna: string) {
  await page.getByRole('button', { name: tarjeta, exact: true }).focus()
  await page.keyboard.press('Space')
  const anuncio = page.locator('[id^="DndLiveRegion"]').first()
  for (let n = 0; n < 40 && !((await anuncio.textContent()) ?? '').includes(`over droppable area ${columna}`); n++) { await page.keyboard.press('ArrowRight'); await page.waitForTimeout(40) }
  await page.keyboard.press('Space')
}
