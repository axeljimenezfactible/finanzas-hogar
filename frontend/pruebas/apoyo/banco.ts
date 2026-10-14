// Vite de desarrollo para el banco de pruebas de la base (pruebas/banco/): monta las piezas de src/ sin pasar por
// dist/ ni por una pantalla de área. Sin red: solo 127.0.0.1 y los paquetes ya instalados.
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { createServer, type ViteDevServer } from 'vite'
const RAIZ = fileURLToPath(new URL('../..', import.meta.url))
export interface Banco { url: (caso: string) => string; cerrar: () => Promise<void> }
export async function iniciarBanco(): Promise<Banco> {
  const srv: ViteDevServer = await createServer({ root: RAIZ, configFile: path.join(RAIZ, 'vite.config.ts'), logLevel: 'error', clearScreen: false, server: { host: '127.0.0.1', port: 5181, strictPort: false, hmr: false } })
  await srv.listen()
  const base = srv.resolvedUrls!.local[0].replace(/\/$/, '')
  return { url: (caso) => `${base}/pruebas/banco/index.html?caso=${caso}`, cerrar: () => srv.close() }
}
