// Apoyo común de las pruebas de flujos contra el servidor real (WF-3b u5-flujos). La contraseña se lee de un archivo
// fuera del repositorio (modo 600) y nunca se imprime ni se guarda en capturas.
import { readFileSync } from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { expect, type Page } from '@playwright/test'

/** Raíz del repositorio (frontend/pruebas-u5/flujos → ../../..). */
export const RAIZ = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../..')

export function contrasenaDesarrollo(usuario: string): string {
  const archivo = process.env.FH_DEV_SEMILLA ?? path.join(os.homedir(), '.config', 'finanzas-hogar', 'dev-semilla.txt')
  for (const renglon of readFileSync(archivo, 'utf8').split('\n')) {
    const i = renglon.indexOf('=')
    if (!renglon.startsWith('#') && i > 0 && renglon.slice(0, i).trim() === usuario) return renglon.slice(i + 1).trim()
  }
  throw new Error(`No hay contraseña de desarrollo para «${usuario}»`)
}

/** Entra como «usuario» (sin capturar nada con la contraseña) y llega a /inicio. */
export async function entrar(page: Page, usuario = 'axel', destino = '/inicio') {
  await page.goto(destino)
  await expect(page.getByRole('heading', { name: 'Iniciar sesión' })).toBeVisible()
  await page.getByLabel('Usuario').fill(usuario)
  await page.getByLabel('Contraseña').fill(contrasenaDesarrollo(usuario))
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page.getByRole('heading', { name: 'Iniciar sesión' })).toBeHidden() // la sesión ya quedó abierta (el h1 del login también es nivel 1)
  await expect(page.getByRole('heading', { level: 1 })).toBeVisible()
}

export const texto = async (page: Page, sel = 'main') => (await page.locator(sel).first().innerText()).replace(/\s+/g, ' ').trim()

/** Cuenta de toques: SOLO pulsaciones en controles (K10). Teclear (fill) y los valores por defecto no cuentan. */
export class Toques {
  n = 0
  pasos: string[] = []
  async tocar(loc: ReturnType<Page['locator']>, nombre: string) {
    await loc.tap()
    this.n++
    this.pasos.push(`${this.n}. ${nombre}`)
  }
}

// ── Registro de flujos (tabla flujo → pasos → P-xx → resultado) ─────────────────────────────────────────────────────
import { existsSync, mkdirSync, readFileSync as leer, writeFileSync } from 'node:fs'
export const DIR_TMP = path.join(RAIZ, 'datos/tmp/u5-flujos')
export const DIR_CAPTURAS = path.join(RAIZ, 'frontend/pruebas-u5/flujos/capturas')
export interface PasoReg { n: number; accion: string; pantalla: string; toques: number | null; ok: boolean; resultado: string }
export interface FlujoReg { id: string; nombre: string; toquesTotal: number; ok: boolean; pasos: PasoReg[]; capturas: string[] }

export class Flujo {
  readonly pasos: PasoReg[] = []
  readonly capturas: string[] = []
  readonly t = new Toques()
  private fallo = false
  constructor(readonly id: string, readonly nombre: string, private page: Page) {}
  /** Ejecuta un paso; si falla lo registra (con el texto visible de la pantalla) y detiene el flujo. */
  async paso(accion: string, pantalla: string, f: () => Promise<string | void>, contarToque = false, continuar = false) {
    if (this.fallo) { this.pasos.push({ n: this.pasos.length + 1, accion, pantalla, toques: null, ok: false, resultado: 'no se ejecutó (paso previo falló)' }); return }
    try {
      const antes = this.t.n
      const r = await f()
      void antes
      this.pasos.push({ n: this.pasos.length + 1, accion, pantalla, toques: contarToque || this.t.n > 0 ? this.t.n : 0, ok: true, resultado: r || 'ok' })
    } catch (e) {
      if (!continuar) this.fallo = true
      const vis = await this.page.locator('body').innerText().then((s) => s.replace(/\s+/g, ' ').slice(0, 500)).catch(() => '?')
      this.pasos.push({ n: this.pasos.length + 1, accion, pantalla, toques: this.t.n, ok: false, resultado: `FALLA: ${String((e as Error).message).split('\n')[0].slice(0, 300)} · pantalla: ${vis}` })
    }
  }
  async captura(nombre: string) {
    mkdirSync(DIR_CAPTURAS, { recursive: true })
    const f = `${nombre}.jpg`
    await this.page.screenshot({ path: path.join(DIR_CAPTURAS, f), type: 'jpeg', quality: 50 })
    this.capturas.push(`frontend/pruebas-u5/flujos/capturas/${f}`)
  }
  guardar(): FlujoReg {
    const reg: FlujoReg = { id: this.id, nombre: this.nombre, toquesTotal: this.t.n, ok: this.pasos.every((p) => p.ok), pasos: this.pasos, capturas: this.capturas }
    mkdirSync(DIR_TMP, { recursive: true })
    const arch = path.join(DIR_TMP, 'flujos.json')
    const todos: FlujoReg[] = existsSync(arch) ? JSON.parse(leer(arch, 'utf8')) : []
    writeFileSync(arch, JSON.stringify([...todos.filter((x) => x.id !== this.id), reg].sort((a, b) => a.id.localeCompare(b.id)), null, 1))
    console.log(`\n${this.id} ${this.nombre} · toques ${reg.toquesTotal} · ${reg.ok ? 'OK' : 'FALLA'}`)
    for (const p of this.pasos) console.log(`  ${p.n}. [${p.ok ? 'ok' : 'XX'}] ${p.pantalla} · ${p.accion} · ${p.resultado.slice(0, 700)}`)
    return reg
  }
}

/** Llamada a la API desde la propia página (misma sesión y mismo origen; el Node de Playwright no usa el almacén NSS de mkcert). */
export async function api<T = any>(page: Page, url: string, init?: { method?: string; headers?: Record<string, string>; body?: unknown }): Promise<{ estado: number; json: T }> {
  return page.evaluate(async ({ url, init }) => {
    const r = await fetch(url, { method: init?.method ?? 'GET', credentials: 'same-origin', headers: { ...(init?.body !== undefined ? { 'Content-Type': 'application/json' } : {}), ...(init?.headers ?? {}) }, body: init?.body !== undefined ? JSON.stringify(init.body) : undefined })
    let json: unknown = null
    try { json = await r.json() } catch { /* sin cuerpo */ }
    return { estado: r.status, json }
  }, { url, init }) as Promise<{ estado: number; json: T }>
}
