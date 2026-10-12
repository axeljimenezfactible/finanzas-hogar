import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { defineConfig } from '@playwright/test'

// WF-3b · u5-flujos (C3-08, U5-09, U5-13, C3-09): pruebas contra el servidor real YA EN MARCHA (ops/dev-run.sh,
// https://finanzas.local:8443) con la semilla; sin mocks. No entran en «npm run verify».
//   cd frontend && npx playwright test -c pruebas-u5/flujos/playwright.config.ts
// FH_BASE: otra dirección (https://127.0.0.1:8443). FH_DEV_SEMILLA: archivo de contraseñas (fuera del repositorio; nunca se imprime).
// FH_CHROMIUM: ejecutable de chromium-headless-shell. TLS: la CA de mkcert está en ~/.pki/nssdb; no se ignoran errores de certificado.
function chromiumLocal(): string | undefined {
  if (process.env.FH_CHROMIUM) return process.env.FH_CHROMIUM
  const base = path.join(os.homedir(), '.cache', 'ms-playwright')
  try {
    const revs = fs.readdirSync(base).filter((n) => /^chromium_headless_shell-\d+$/.test(n)).sort((a, b) => Number(b.split('-')[1]) - Number(a.split('-')[1]))
    for (const r of revs) { const exe = path.join(base, r, 'chrome-headless-shell-linux64', 'chrome-headless-shell'); if (fs.existsSync(exe)) return exe }
  } catch { /* sin caché: el de Playwright */ }
  return undefined
}

export default defineConfig({
  testDir: '.',
  testMatch: /.*\.spec\.ts/,
  outputDir: '../../../datos/tmp/playwright-u5-flujos',
  workers: 1,
  retries: 0,
  reporter: [['list']],
  timeout: 120_000,
  expect: { timeout: 8_000 },
  use: {
    baseURL: process.env.FH_BASE ?? 'https://finanzas.local:8443',
    locale: 'es-MX', timezoneId: 'America/Mexico_City', testIdAttribute: 'data-prueba',
    actionTimeout: 15_000,
    trace: 'off', screenshot: 'off', video: 'off',
    launchOptions: { executablePath: chromiumLocal() },
  },
  projects: [{ name: 'movil-real', use: { browserName: 'chromium', viewport: { width: 412, height: 915 }, hasTouch: true, isMobile: true } }],
})
