// U5 · Capturas, axe y contraste contra la APLICACIÓN REAL (servidor Go con la semilla), no contra el API falso.
// Uso (desde frontend/):  npx playwright test -c playwright.u5.config.ts            (todo: 164 capturas + interruptor)
//                         npx playwright test -c playwright.u5.config.ts --grep axe        (U5-03)
//                         npx playwright test -c playwright.u5.config.ts --grep contraste  (U5-04)
//                         npx playwright test -c playwright.u5.config.ts --grep Tendencias (P-01 ▸ Tendencias, U5-06 #14)
// Requiere el servidor de desarrollo (bash ops/dev-run.sh) en https://finanzas.local:8443 (FH_U5_URL lo cambia).
// 4 proyectos = 4 combinaciones: Android 412×915 con emulación táctil y escritorio 1440×900, cada uno en claro y oscuro
// (colorScheme emulado: el tema sigue al sistema, el interruptor NO se toca en las capturas). Escala declarada: 1.
import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { defineConfig } from '@playwright/test'

function chromiumLocal(): string | undefined {
  if (process.env.FH_CHROMIUM) return process.env.FH_CHROMIUM
  const base = path.join(os.homedir(), '.cache', 'ms-playwright')
  try {
    const revs = fs.readdirSync(base).filter((n) => /^chromium_headless_shell-\d+$/.test(n)).sort((a, b) => Number(b.split('-')[1]) - Number(a.split('-')[1]))
    for (const r of revs) { const exe = path.join(base, r, 'chrome-headless-shell-linux64', 'chrome-headless-shell'); if (fs.existsSync(exe)) return exe }
  } catch { /* sin caché: el de Playwright */ }
  return undefined
}
const comunes = {
  browserName: 'chromium' as const,
  baseURL: process.env.FH_U5_URL ?? 'https://finanzas.local:8443',
  locale: 'es-MX', timezoneId: 'America/Mexico_City', deviceScaleFactor: 1,
  testIdAttribute: 'data-prueba', trace: 'off' as const,
  launchOptions: { executablePath: chromiumLocal() },
}
const android = { viewport: { width: 412, height: 915 }, hasTouch: true, isMobile: true }
const escritorio = { viewport: { width: 1440, height: 900 } }
export default defineConfig({
  testDir: 'pruebas-u5',
  outputDir: '../datos/tmp/playwright-u5-capturas',
  workers: 1,
  timeout: 120_000,
  reporter: [['list']],
  globalTeardown: './pruebas-u5/cierre.ts',
  projects: [
    { name: 'android-claro', use: { ...comunes, ...android, colorScheme: 'light' }, testMatch: /(capturas|tema|tendencias)\.spec\.ts/ },
    { name: 'android-oscuro', use: { ...comunes, ...android, colorScheme: 'dark' }, testMatch: /(capturas|tendencias)\.spec\.ts/ },
    { name: 'escritorio-claro', use: { ...comunes, ...escritorio, colorScheme: 'light' }, testMatch: /(capturas|tema|tendencias)\.spec\.ts/ },
    { name: 'escritorio-oscuro', use: { ...comunes, ...escritorio, colorScheme: 'dark' }, testMatch: /(capturas|tendencias)\.spec\.ts/ },
  ],
})
