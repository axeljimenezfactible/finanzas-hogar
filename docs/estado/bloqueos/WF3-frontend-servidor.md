# Bloqueo WF3-frontend-servidor · Cómo entrega el servidor Go la PWA

- Fase: WF-3a, agente `frontend-base` (12-oct-2026). Estado: **PROVISIONAL**; lo cierra `backend-base` (CONVENCIONES-backend.md) o el integrador.
- **Qué falta (PREGUNTA):** dónde y cómo sirve el servidor `frontend/dist/` en `https://finanzas.local:8443`.
- **Lo que el frontend necesita (PROPUESTA):**
  1. Mismo origen para la app y el API (`/api/v1/*`, `/health`): la cookie `fh_session` es `SameSite=Strict` y el cliente usa rutas relativas.
  2. Respaldo de SPA: toda ruta GET que no sea `/api/`, `/health` ni un archivo de `dist/` responde `dist/index.html` (rutas de `WF3-frontend-rutas`).
  3. `sw.js`, `index.html` y `manifest.webmanifest` con `Cache-Control: no-cache`; `assets/*` (con hash) con `max-age=31536000, immutable`; `manifest.webmanifest` con `Content-Type: application/manifest+json`.
  4. Cabeceras (K31, puerta 4): el frontend no usa scripts en línea ni recursos de terceros, así que admite `Content-Security-Policy: default-src 'self'; img-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; connect-src 'self'; worker-src 'self'; frame-ancestors 'none'` (`'unsafe-inline'` en estilos por los `style=` que ponen dnd-kit y ECharts).
  5. En desarrollo, Vite reenvía `/api` y `/health` a `FH_API` (por omisión `https://127.0.0.1:8443`).
- **Opciones:** (A) `embed` de `dist/` en el binario; (B) carpeta `frontend/dist` leída en ejecución. **Recomendación:** A para producción (un solo binario), B en desarrollo.
- **Impacto:** sin el punto 2, recargar una página que no sea `/` da 404 fuera de la PWA instalada.
