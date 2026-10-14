import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { defineConfig } from '@playwright/test'

// Chromium local sin red: FH_CHROMIUM si viene; si no, el chrome-headless-shell más nuevo que ya esté en la caché de
// Playwright (así `make verify-wf3-frontend` → `npm run verify` corre sin variables extra). Sin ninguno, el de Playwright.
function chromiumLocal(): string | undefined {
  if (process.env.FH_CHROMIUM) return process.env.FH_CHROMIUM
  const base = path.join(os.homedir(), '.cache', 'ms-playwright')
  try {
    const revs = fs.readdirSync(base).filter((n) => /^chromium_headless_shell-\d+$/.test(n)).sort((a, b) => Number(b.split('-')[1]) - Number(a.split('-')[1]))
    for (const r of revs) { const exe = path.join(base, r, 'chrome-headless-shell-linux64', 'chrome-headless-shell'); if (fs.existsSync(exe)) return exe }
  } catch { /* sin caché: el de Playwright */ }
  return undefined
}

// Pruebas base (frontend-base). Corren contra «vite preview» del build real (con service worker);
// FH_CHROMIUM: ruta a un Chromium local si el de Playwright no está descargado (sin red en ejecución).
// el API se simula con context.route (pruebas/apoyo/api-falsa.ts). Sin red externa.
export default defineConfig({
  testDir: 'pruebas',
  outputDir: '../datos/tmp/playwright-frontend',
  workers: 1,
  reporter: [['list']],
  timeout: 45_000,
  use: { baseURL: 'http://localhost:4173', locale: 'es-MX', timezoneId: 'America/Mexico_City', trace: 'off', testIdAttribute: 'data-prueba', launchOptions: { executablePath: chromiumLocal() } },
  projects: [
    { name: 'movil', use: { browserName: 'chromium', viewport: { width: 412, height: 915 }, hasTouch: true, isMobile: true } },
    { name: 'escritorio', use: { browserName: 'chromium', viewport: { width: 1440, height: 900 } }, testMatch: /(login|shell|\.escritorio)\.spec\.ts/ },
  ],
  webServer: { command: 'npm run preview', url: 'http://localhost:4173', reuseExistingServer: true, timeout: 60_000 },
})
