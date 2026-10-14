import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { defineConfig } from '@playwright/test'

// Humo real del integrador frontend/C3 (WF-3a): corre contra el servidor de desarrollo YA EN MARCHA
// (ops/dev-run.sh, https://finanzas.local:8443), sin mocks ni «vite preview». No entra en «npm run verify».
//   npx playwright test -c playwright.integracion.config.ts
// FH_BASE: otra dirección (p. ej. https://127.0.0.1:8443). FH_DEV_SEMILLA: archivo de contraseñas de desarrollo
// (por omisión ~/.config/finanzas-hogar/dev-semilla.txt, fuera del repositorio; nunca se imprime).
// TLS: la CA de mkcert está en el almacén NSS del usuario (~/.pki/nssdb), que es el que usa Chromium en Linux,
// así que NO se ignoran errores de certificado. FH_IGNORAR_TLS=1 solo como último recurso (y se documenta).
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
  testDir: 'pruebas-integracion',
  outputDir: '../datos/tmp/playwright-integracion',
  workers: 1,
  retries: 0,
  reporter: [['list']],
  timeout: 60_000,
  use: {
    baseURL: process.env.FH_BASE ?? 'https://finanzas.local:8443',
    ignoreHTTPSErrors: process.env.FH_IGNORAR_TLS === '1',
    locale: 'es-MX', timezoneId: 'America/Mexico_City', testIdAttribute: 'data-prueba',
    trace: 'off', screenshot: 'off', video: 'off',
    launchOptions: { executablePath: chromiumLocal() },
  },
  projects: [{ name: 'movil-real', use: { browserName: 'chromium', viewport: { width: 412, height: 915 }, hasTouch: true, isMobile: true } }],
})
