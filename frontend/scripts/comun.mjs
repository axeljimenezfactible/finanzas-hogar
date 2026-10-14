// Utilidades de los scripts del frontend. Rutas relativas a frontend/.
import { fileURLToPath } from 'node:url'
import path from 'node:path'
export const RAIZ_FRONT = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')
export const RAIZ_REPO = path.resolve(RAIZ_FRONT, '..')
export const desdeRepo = (...p) => path.join(RAIZ_REPO, ...p)
export const desdeFront = (...p) => path.join(RAIZ_FRONT, ...p)
