import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { VitePWA } from 'vite-plugin-pwa'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import colores from './src/estilos/colores-manifest.json' with { type: 'json' }

// Solo en «npm run dev:mocks»: sirve el worker de msw sin copiarlo a public/ (nunca llega a dist/).
function workerMocks(): Plugin {
  const archivo = fileURLToPath(new URL('./node_modules/msw/lib/mockServiceWorker.js', import.meta.url))
  return {
    name: 'fh-worker-mocks',
    apply: 'serve',
    configureServer(servidor) {
      servidor.middlewares.use('/mockServiceWorker.js', (_pet, res) => {
        res.setHeader('Content-Type', 'text/javascript')
        res.end(readFileSync(archivo))
      })
    },
  }
}

// PROVISIONAL (bloqueo WF3-frontend-servidor): en desarrollo el API se pide al servidor Go por proxy.
const API = process.env.FH_API ?? 'https://127.0.0.1:8443'

export default defineConfig(({ mode }) => ({
  base: '/',
  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  plugins: [
    react(),
    tailwindcss(),
    mode === 'mocks' ? workerMocks() : null,
    VitePWA({
      registerType: 'autoUpdate',
      injectRegister: false,
      includeAssets: ['favicon.svg', 'tema-inicial.js', 'fonts/*.woff2', 'icons/*.png'],
      manifest: {
        id: '/',
        name: 'Finanzas Hogar',
        short_name: 'Finanzas',
        description: 'Gastos, presupuesto, comidas, compras y créditos del hogar.',
        lang: 'es-MX',
        dir: 'ltr',
        start_url: '/inicio',
        scope: '/',
        display: 'standalone',
        theme_color: colores.tema,
        background_color: colores.fondo,
        icons: [
          { src: '/icons/icon-192.png', sizes: '192x192', type: 'image/png' },
          { src: '/icons/icon-512.png', sizes: '512x512', type: 'image/png' },
          { src: '/icons/icon-maskable-512.png', sizes: '512x512', type: 'image/png', purpose: 'maskable' },
        ],
        shortcuts: [{ name: 'Registrar gasto', short_name: 'Gasto', url: '/inicio?hoja=gasto' }],
      },
      workbox: {
        // Caché de la app (cascarón). Los datos NO pasan por el service worker: la cola de gastos
        // y peso vive en IndexedDB (src/cola) y la última carga de lecturas en src/api/lecturas.ts.
        globPatterns: ['**/*.{js,css,html,woff2,png,svg,webmanifest}'],
        navigateFallback: '/index.html',
        navigateFallbackDenylist: [/^\/api\//, /^\/health/],
        cleanupOutdatedCaches: true,
        clientsClaim: true,
        skipWaiting: true,
        maximumFileSizeToCacheInBytes: 3 * 1024 * 1024,
      },
      devOptions: { enabled: false },
    }),
  ],
  server: {
    host: '127.0.0.1',
    port: 5173,
    proxy: mode === 'mocks' ? undefined : { '/api': { target: API, secure: false }, '/health': { target: API, secure: false } },
  },
  preview: { host: '127.0.0.1' },
  build: { target: 'es2022', sourcemap: false, chunkSizeWarningLimit: 1200 },
}))
