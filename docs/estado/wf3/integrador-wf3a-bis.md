# integrador-wf3a-bis · WF-3a bis · 12-oct-2026

Agente: `integrador-wf3a-bis` (Opus).

Alcance:
- `docs/estado/wf3/INTEGRACION-wf3a.md`, sección nueva §15 «WF-3a bis»;
- `README.md`, estado;
- `docs/estado/bloqueos/INDICE.md`, sección «Bloqueos nuevos de WF-3a bis»;
- este reporte y el commit.

No toqué código, pruebas, contratos, motores ni la Pista 3. No envié Telegram. No hay credenciales en ningún archivo.

## Hecho
- **HECHO.** Corrí la verificación completa del frontend y del backend, con el detalle de comandos y salidas en INTEGRACION §15.2:
  - build;
  - Playwright completo y `cola.spec.ts` ×3;
  - rutas y hosts;
  - `go test ./...`;
  - `make verify-wf3` y `make verify-all`;
  - C3 regenerada sin `--escribir`.
- **HECHO.** Reinicié el servidor de desarrollo con la base re-sembrada: `migrate`, después `semilla` y al final `semilla --desarrollo`, todo idempotente. Queda **corriendo**: PID 2419156, en 127.0.0.1:8443 y 192.168.50.56:8443.
  - `https://finanzas.local:8443/` sirve el `dist` nuevo, igual byte a byte.
  - Las 6 tarjetas de Inicio responden 200 con sesión; 5 de las 6 traen montos distintos de 0.
- **HECHO.** Escribí INTEGRACION §15: hallazgos cerrados y abiertos con su ronda, C3 v0.95, semilla de desarrollo y bloqueos nuevos. Actualicé el README y el INDICE de bloqueos (4 nuevos; 129 archivos).
- **HECHO.** Hice el commit «[WF-3a bis] …». No entran `.db`, certificados, `node_modules`, `dist` ni `frontend/tsconfig.tsbuildinfo`.

## Archivos
- `docs/estado/wf3/INTEGRACION-wf3a.md` (§15 nueva)
- `README.md` (párrafo de estado y «Cómo sigue» 6)
- `docs/estado/bloqueos/INDICE.md` (sección nueva al final)
- `docs/estado/wf3/integrador-wf3a-bis.md` (este archivo)

## Verificación
- **C3-03** · `npm run build` rc 0. `node scripts/rutas-api.mjs` da «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». El script C3 cuenta 210 rutas, con 0 fuera del contrato.
- **U5-08** · `node scripts/hosts-externos.mjs dist`: «hosts externos que se pedirían = 0».
- **U5-13** · `npx playwright test pruebas/cola.spec.ts --repeat-each=3`: 12 passed.
- **Suite completa (U5-06, U5-07, U5-09 con mocks)** · `npx playwright test`: 230 passed, 4 skipped y **1 failed**, en `gastos-presupuesto/correccion-bis.escritorio.spec.ts:79` [escritorio] (R5-05).
- **S6-01** · `go test ./... -count=1`: 32 paquetes OK, 5 sin pruebas y 0 fallas. Dentro de `make verify-wf3`: build, vet, gofmt y `go mod tidy -diff` limpios.
- **S6-02** · K25: dominio 87.0 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 % y precios 88.8 %.
- **S6-03** · «los 23 archivos coinciden byte a byte con el generador». «216 operaciones en 172 rutas, 0 diferencias». `TestContrato_Operaciones`: 216 PASS.
- **`make -k verify-wf3`** · **rc 2**. Solo falla `verify-wf3-frontend`, por la misma prueba.
- **`make verify-all GOOSE=…`** · **rc 0**: «Verificación ampliada OK».
- **C3-01, C3-02, C3-04 y C3-06** · El script C3 sin `--escribir` da:
  - 42 filas (41 `P-` y la fila Global), 0 pantallas sin fila y 5 operaciones sin pantalla, todas justificadas;
  - implementado 42/42, probado 42/42 y capturado 0/42;
  - no llama 1 (`listPeriods`, con excepción y cobertura);
  - «Matriz del documento = matriz generada: sí».
- **S6-11 (parcial)** · `curl` con la CA de mkcert contra `https://finanzas.local:8443/` da 200. No volví a correr `make verify-wf3-servidor`.
- **S6-12** · `/sw.js` y los recursos de `/assets/` son iguales a los de `frontend/dist`.
- **S6-14** · `ss -ltn :8443` muestra solo `127.0.0.1:8443` y `192.168.50.56:8443`.
- **Humo de Inicio (centavos)** · month-vs-plan 119,764 / 200,000; fortnight-flow 1,200,000; best-basket 309,450 (−26,450); savings 314,672 / 44,953; next-payment 97,375 (59,882 + 37,493) y disponible B 690,550; week-meals responde con costo 0 (bloqueo `WF3a-bis-semilla-dev` punto 5).
- **S6-13 en lo que entra al commit** · Mi script, que solo imprime conteos, encontró 0 llaves privadas, 0 tokens de bot, 0 llaves de API, 0 hashes argon2 y 0 contraseñas de desarrollo.

## Supuestos
- `frontend/tsconfig.tsbuildinfo` es un archivo generado por `tsc`: apareció a las 06:46 y no lo crea `npm run build`. Lo dejé fuera del commit con `':!frontend/tsconfig.tsbuildinfo'`. Pido agregarlo a `.gitignore` (solicitud 4 de §15.8).
- «Las 6 tarjetas devuelven cifras»: lo cuento como 6 respuestas 200 con su forma de C1. La de comidas trae costo 0 por un bloqueo conocido; las otras 5 traen montos distintos de 0.
- El humo abrió 2 sesiones de Axel («humo integrador-wf3a-bis») en la base de **desarrollo**.

## Bloqueos
- Sin bloqueos nuevos de este integrador. Di de alta los 4 de la ronda en el INDICE.
- `make verify-wf3` queda en rojo por R5-05, que está fuera de mi alcance (dueño: front-gastos-presupuesto).

## Pendientes
- Decisión del orquestador sobre una ronda mínima: core R5-01, áreas R5-01 y R5-05 (INTEGRACION §15.8).
- Solicitudes fuera de alcance 1 a 4 de §15.8: `ops/dev-run.sh`, `ops/runbook.md`, `frontend/package.json` y `.gitignore`.
- WF-3b: capturas U5, C3-08, C3-09 y firma de C3.

## Propuestas para Axel
- R5-02: con «A crédito», deshabilitar «Guardar» y mostrar el motivo «La compra a crédito se termina en Créditos».
- N-053: listar solo las líneas revolventes y de tienda, como dice U1. Es la opción (a) de R5-03; hoy la opción PROVISIONAL lista todas.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 58 acciones, unos 180 KB leídos y unos 30 KB escritos. Los temporales estaban en `datos/tmp/bis/` y los borré al cerrar.
