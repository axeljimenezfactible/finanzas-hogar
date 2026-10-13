// Unidad (`unit`) que manda HOY el servidor para cada parámetro de K44: la tercera columna de la anotación `@param` de las
// migraciones, que es lo que lee `parseParamDefs` (backend/internal/dominio/ajustes.go, la última anotación gana). Los datos falsos
// de P-26.02 la toman de aquí para no inventar unidades (M-1 de docs/estado/wf3/REVISION-frontend-cierre-r1.md); para las 5 claves
// de K54 es una descripción larga (bloqueo docs/estado/bloqueos/C1-unidad-parametros-tablero.md, PROVISIONAL).
import { readdirSync, readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
const MIGRACIONES = fileURLToPath(new URL('../../../backend/db/migrations/', import.meta.url))
export function unidadesDelServidor(): Record<string, string> {
  const r: Record<string, string> = {}
  for (const archivo of readdirSync(MIGRACIONES).filter((n) => n.endsWith('.sql')).sort()) {
    for (const renglon of readFileSync(MIGRACIONES + archivo, 'utf8').split('\n')) {
      const linea = renglon.trim()
      if (!linea.startsWith('-- @param ')) continue
      const f = linea.slice('-- @param '.length).split('|').map((x) => x.trim())
      if (f.length === 7 && !f[0].startsWith('<')) r[f[0]] = f[2]
    }
  }
  return r
}
