# WF-3c · u5-cierre · Reporte

Agente `u5-cierre` (Sonnet) · 12-oct-2026 · HEAD `0197dad` + árbol de trabajo. Solo regeneré evidencia: no cambié código de aplicación, ni criterios, ni contratos. No hice `git add` ni commit.

## Hecho
HECHO:
1. `npm run build` (frontend, PWA con 80 entradas) y `bash ops/dev-stop.sh` + `bash ops/dev-run.sh` (recompila `backend/bin/finanzas`; la base de desarrollo ya estaba en el esquema 1.5, versión de goose 5). `/health` responde `ok` en `https://finanzas.local:8443`.
2. `make verify-u5-servidor`: 174 pruebas pasan (164 capturas con axe y contraste, 4 de tema, 4 de Tendencias, 2 del interruptor), ≈ 9 min. Regeneró los 164 PNG, `axe/`, `contraste/`, `hosts.json` y `evidencia/`.
3. `make verify-wf3` (incluye `verify-u5`): verde. 271 pruebas de interfaz pasan, 4 se omiten solas por proyecto. `make verify-all`: verde.
4. Cifras del tablero contra la API (P-01 Resumen, 4 combinaciones, `evidencia/tablero-cifras/*.json`, corrida de solo lectura con un spec de `pruebas-u5/tablero-cifras.spec.ts` y una configuración temporal que borré): ver Verificación.

## Archivos
- `docs/design/capturas-u5/` regenerado: 164 PNG, `axe/` 164 JSON, `contraste/` 164 JSON, `hosts.json`, `evidencia/{manifiesto,cifras-oraculo,tema-*}.json`, `evidencia/tendencias/`, `evidencia/inicio-subsecciones/`, `evidencia/tablero-cifras/` (4 JSON nuevos de esta corrida).
- `frontend/dist/` recompilado y `backend/bin/finanzas` recompilado (ambos ignorados por git).
- Este reporte. Temporales de `datos/tmp/` (registros `u5-cierre-*.log`) borrados.

## Verificación
- U5-01: `python3 -I frontend/pruebas-u5/verificar_capturas.py` da `T = 41 · 4T = 164 · PNG correctos: 164 de 164 · faltantes: 0 · con dimensiones o píxeles mal: 0` (412x915 Android y 1440x900 escritorio, anchos exactos). `ls docs/design/capturas-u5 | wc -l` = 168 (164 PNG + `axe`, `contraste`, `evidencia`, `hosts.json`): sigue abierto el bloqueo `U5-conteo-carpeta`.
- U5-03: 164 JSON de axe, **0 serias o críticas** y 0 de otro impacto en las 4 combinaciones (incompletos: 17 · 17 · 8 · 8).
- U5-04: 164 JSON, **0 pares bajo umbral**, 0 indeterminados, 25,818 pares evaluados (5,659 · 5,680 · 7,229 · 7,250). Sin falsos positivos que descartar: la lista agrupada de bajo umbral está vacía.
- U5-08: 1 host único, `finanzas.local:8443`, en 164 capturas; 0 hosts externos. Las pruebas de interfaz (`login.spec.ts`) también: externos 0.
- U5-10: 0 capturas con el fondo del body distinto del token de C2 o con `data-theme`; `tema.spec.ts` pasa en Android y escritorio; `shell.spec.ts` U5-10 pasa.
- U5-02: las cifras de la semilla se ven ($1,191.50; $10,000.00; $3,094.50; $6,293.44; $6,905.50; $973.75). Ninguna captura quedó cargando ni en pantalla de error.
- U5-07 (informativo): 0 fechas ISO, 0 palabras en inglés, 0 claves internas, 0 «sin dato». Android con puntero grueso: 82 de 82.
- U5-11: 41 filas P- de C3 con «capturado» = sí y sus 4 rutas en disco.
- `make verify-wf3`: `== WF-3: verificación completa OK (Go, cobertura K25, contrato C1, interfaz y evidencia U5)`.
- `make verify-all`: `== Verificación ampliada OK`.
- Tablero contra la API (fecha del servidor 12-oct-2026):
  - Deuda viva de hoy: `GET /credits/totals` da `live_debt_cents` = 389500 → pantalla «Créditos $3,895.00 deuda viva total». Coincide con el oráculo ($3,895.00).
  - Próximo pago: `GET /dashboard/next-payment` da `amount_cents` = 97375, vence 2026-10-30, hogar 59882 y Karen 37493 → pantalla «Vie 30 oct 2026 · Kueski · $973.75 · Hogar $598.82 · Karen $374.93». Coincide ($973.75).
  - Gasto del mes 119764 de 200000 → «$1,197.64 · 60 % de $2,000.00 · te quedan $802.36»; capacidad de endeudamiento 2 % (pagos 59882 de 2450000), «Te caben $9,201.18»; pagos de octubre «$973.75 · 1 pago», falta pagar 59882 → «Falta pagar $598.82»; disponible hoy 690550 → «$6,905.50 · Lectura B». Todo coincide con la pantalla en las 4 combinaciones.

## Supuestos
1. La «deuda viva de hoy» se tomó de `getCreditTotals`, que es la cifra que pinta K-3 (la API del tablero no la trae). Rige la opción PROVISIONAL (a) del bloqueo `S6-deuda-viva-hoy`.
2. No re-sembré la base: ya estaba migrada al esquema 1.5 y con la semilla de desarrollo de la corrida anterior.

## Bloqueos
Ninguno nuevo. Siguen abiertos `U5-conteo-carpeta`, `U5-nombres-propios`, `U5-tipografia-tres-tamanos` y `S6-deuda-viva-hoy` (puerta 4).

## Pendientes
- Observación (no es falla de criterio): `GET /api/v1/debt/live` sin parámetro de fecha respondió `date` 2027-09-30 con `total_cents` 0. Conviene que backend confirme cuál es la fecha por omisión (el contrato dice «hoy» para `getCreditTotals`). Ruta: `backend/internal/api/creditos/`.
- Observación: U5-07 registra 2 importes «$0» por combinación en P-01 y P-03 (etiqueta del eje Y de las gráficas); §7.7 permite el compacto en ejes, así que no son aserción.
- La revisión visual (U5-05, U5-06) y los flujos de C3-08/U5-09/U5-13 en copia aislada no se tocaron.

## Propuestas para Axel
Ninguna.

## Paquetes instalados
Ninguno.

## Consumo aproximado
≈ 22 acciones, ≈ 60 KB leídos, ≈ 6 KB escritos.
