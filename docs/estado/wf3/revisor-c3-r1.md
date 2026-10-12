# WF-3b · revisor-c3-r1 · Reporte

Fecha: 2026-10-12 · Agente: `revisor-c3-r1` (revisor independiente, ronda 1 de 3) · Objeto: `docs/contratos/C3-integracion.md` v1.0 · Criterios: `docs/criterios/C3.md` 1.1 · Informe: `docs/estado/wf3/REVISION-c3-r1.md`.

## Hecho
- Reejecuté C3-01 a C3-11 con comandos propios (abajo). **Veredicto: devuelto, sin firma** (3 altos en C3-08, 0 medios, 5 bajos). No modifiqué la matriz ni agregué bloque de firma.
- Script propio en Python (`python3 -I`, temporal en `datos/tmp/revisor-c3-r1/`, borrado al final) que lee U1 §2.1, §2.3, §3.1–3.28 y §3.29, `openapi.yaml` (`operationId`, ruta, método, `x-necesidades`, `x-agente-externo`), `routes.go`, `frontend/src` (literales de ruta y llamadas `pedir`/`todas`/`usarConsulta` con método; cierre de importaciones por archivo de entrada), la matriz entre sus marcadores, el registro de `go test` y los PNG (cabecera IHDR).
- Consulta de solo lectura a la API de desarrollo (sesión de `axel`, contraseña leída del archivo que indica `semilla.md`, nunca impresa; cierre de sesión al final) para el supuesto 6 de C3.

## Archivos
- Escritos: `docs/estado/wf3/REVISION-c3-r1.md`, `docs/estado/wf3/revisor-c3-r1.md`.
- Sin cambios: `docs/contratos/C3-integracion.md` (SHA-256 del archivo revisado `09e6b63c3fd9932c0bc6deef82563522e35ff4adb534f6859371c52195cad502`).
- Efectos laterales de las pruebas (no son míos ni los edité): `npm run build` regeneró `src/api/esquema.gen.ts`, `src/mocks/*.gen.json` y `src/estilos/*` con el mismo contenido (git sin cambios); las corridas de flujos reescribieron las capturas de paso `.jpg` de `frontend/pruebas-u5/flujos/capturas/` (no las `U5-13-*`), `datos/tmp/u5-flujos/` (copia de la base y `c3-09.json`) y `datos/tmp/playwright-u5-flujos/`, carpetas temporales de `u5-flujos` que dejo como están.

## Verificación
- **C3-01** · Cumple · `grep -c "^| P-"` = 41 = T (27 + 14) + fila Global; §1 trae las 3 definiciones.
- **C3-02** · Cumple · script propio: 41 pantallas, 0 sin fila; 211 `N-xxx` (U1 §3 = §3.29), 0 sin operación ni Local, 0 de C1 fuera de U1; operaciones de las 42 filas = `x-necesidades` ∩ U1 (0 diferencias); 216 operaciones, 5 sin pantalla con excepción una por una; openapi = routes.go 216/216.
- **C3-03** · Cumple · script propio: 166 rutas / 210 (método, ruta) invocadas, 0 fuera de `openapi.yaml`, 0 sin registrar; 0 necesidades sin llamada en el cierre de cada pantalla; `npm run build` exit 0, 0 errores.
- **C3-04** · Cumple · `go test ./... -count=1 -v`: 32 ok, 0 FAIL, 1126 PASS, 216/216 contrato; `npx playwright test`: 245 / 241 pasan / 4 omitidas / 0 fallan; 69 de 69 pruebas nombradas en la matriz pasaron; script del redactor con mis registros: matriz = generada.
- **C3-05** · Cumple · 164 de 164 PNG (4 × 41), ancho 412/1440, nombres = citados; U5-revision r3 164/164.
- **C3-06** · Cumple · `grep -c "| no |"` = 1 = filas con bloqueo citado (Global → `C3-capturas-global.md`).
- **C3-07** · No se firma (devuelto por C3-08).
- **C3-08** · **No cumple** · `flujos.spec.ts` 6/7 (✘ F-3 `:122`), `pendientes.spec.ts` 4/4, en copias aisladas `127.0.0.1:8444`; F-1 = 3 toques; F-7 y F-8 sin prueba; 5 de 8; base de desarrollo intacta (6 gastos, 0 pesos).
- **C3-09** · Cumple · `c3-09.spec.ts` 1/1, 111 cifras, 0 diferencias; `humo-inicio.spec.ts` 1/1; API: 97375 el 2026-10-30, 59882/37493, B 690550, `getLiveDebt` 2027-04-10 → 629344; §5 = `model.json`.
- **C3-10** · Cumple · 10 filas en es-MX, fechas ISO; anglicismo «mocks» en §4 (bajo).
- **C3-11** · Cumple · «Versión v1.0 · 2026-10-12»; 3 bloqueos `C3-*`; `c3-matriz.md` con §Supuestos y §Bloqueos.

## Supuestos
1. Evalué la versión v1.0 que está en el árbol de trabajo (la tarea mencionaba el borrador v0.9 de WF-3a; el historial de C3 muestra que v0.95 y v1.0 la sustituyeron).
2. Para C3-08 tomé el comando declarado por la matriz (`-c pruebas-u5/flujos/playwright.config.ts`, copias aisladas), que C3-08 admite («o el comando declarado»), y N = 8 como declara la matriz; aun con N = 6, F-7 y F-8 son los únicos que cubren los obligatorios (b) y (d) de U1-05.
3. «Probado» con la regla PROVISIONAL de `C3-matriz` punto 4 a (título de la prueba y API simulada): la evalué tal cual; no la cuento como hallazgo porque tiene bloqueo y la definición de C3-01 se cumple.
4. C3-02 «coincide con la matriz de C1»: tomé como matriz de C1 las `x-necesidades` de `openapi.yaml`.

## Bloqueos
Ninguno nuevo. Los abiertos que impiden firmar: `docs/estado/bloqueos/C3-flujos.md` (puntos 1 a 4). Siguen PROVISIONALES: `C3-capturas-global.md` punto 1 y `C3-matriz.md` puntos 1, 3 c y 4.

## Pendientes
- Para la ronda 2: pruebas de F-7 y F-8 contra el servidor real, decisión sobre F-3 (y su reejecución), y regenerar la matriz (H-4, H-6, H-8 opcionales). Luego reejecutar C3-02 a C3-09 y firmar.
- Integrador: `frontend/pruebas-u5/flujos/` está sin seguimiento en git y la matriz cita evidencia de ahí.

## Propuestas para Axel
- **C3-08**: la opción (c) de `C3-flujos` punto 1 («fijar N = 6, F-1 a F-6») no basta por sí sola: C3-08 pide «los 4 obligatorios» de U1-05 y (b) compra a crédito y (d) ticket con OCR solo los cubren F-7 y F-8. Recomiendo la opción (b): encargar las 2 pruebas en copia aislada (costo bajo; riesgo de que F-7 muestre «Faltan restricciones» con la semilla, que es lo correcto según K16).
- **F-3**: decidir si la semilla trae precios de los ingredientes (para la ruta A) y qué cifra de Walmart vale (U1 $3,233.00 o la del motor S5).
- **Deuda viva**: la cifra correcta en la API sale de `getLiveDebt` con fecha 2027-04-10 (629344 centavos); conviene citarla así en C3.

## Consumo aproximado
~125 KB leídos (C3.md, partes de la matriz, bloqueos, extractos de U1, openapi y salidas), ~13 KB escritos, 37 acciones. Pruebas: `go test` (~1 min), `npm run build` (~1 min), Playwright completo (5.2 min), flujos (55 s), C3-09 y humo (15 s).
