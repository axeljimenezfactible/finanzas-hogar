# revisor-c3-r3 · Reporte (WF-3b, ronda 3 de 3 de C3)

Fecha: 2026-10-12 · Agente: `revisor-c3-r3` (revisor independiente; no redacté la matriz ni implementé nada) · Base: HEAD `6d691d0` con el árbol de trabajo de WF-3b · Objeto: `docs/contratos/C3-integracion.md` v1.2 (SHA-256 `4d730b0cd34fd208c77d04e50053b3e09d02e9eaddab39db9101e4bf6ce4284e`, sin cambios) · Criterios: `docs/criterios/C3.md` 1.1. Sin `git add` ni commit.

## Hecho
- Revisé C3 v1.2 contra `C3.md` completo, con mis propios scripts y ejecuciones. **Veredicto: devuelto, sin firma.** Hay 1 alto (R3-1: C3-08, F-3 ruta A con la semilla), 0 medios y 6 bajos que ya venían de r2. Detalle en `docs/estado/wf3/REVISION-c3-r3.md`.
- Comprobé los hallazgos de r2. **R2-1 queda reabierto en parte:** la ruta B queda cerrada, porque K46 (Q2-47 = A) fija la cifra y la prueba pasa; la ruta A sigue abierta y está en el punto 2 de `C3-flujos`. Los 6 bajos siguen abiertos.
- Como es la ronda 3 de 3, lo escalo a Axel: creé `docs/estado/bloqueos/C3-ronda3-escalamiento.md` con opciones y una recomendación.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/REVISION-c3-r3.md` (informe de la revisión)
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/revisor-c3-r3.md` (este reporte)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C3-ronda3-escalamiento.md` (nuevo)
- No toqué la matriz, las capturas ni las pruebas. Los temporales de `datos/tmp/revisor-c3-r3/` los borré al terminar. Respaldé y restauré `frontend/pruebas-u5/flujos/capturas/` y `resultados/`.

## Verificación (comandos reejecutados y salidas)
- **C3-01 · cumple.** `grep -c "^| P-" docs/contratos/C3-integracion.md` da 41 = T de U1 (27 + 14), más la fila Global; las 3 definiciones están en §1.
- **C3-02 · cumple.** `python3 -I c302.py` (script propio): U1 da 27 `P-xx` y 14 `P-26.xx`; hay 211 N definidas, que son las 211 de §3.29; 0 N sin operación ni Local; 0 N de C1 fuera de U1; 0 pantallas sin fila y 0 filas sin pantalla; 0 filas con operaciones distintas de `x-necesidades` ∩ §3.29; 5 operaciones sin pantalla, las 5 en la lista de excepciones; `routes.go` = `openapi.yaml` (216).
- **C3-03 · cumple.** `python3 -I c303.py` (script propio): 116 archivos, 166 literales y 208 pares (método, ruta); 0 fuera de `openapi.yaml` y 0 sin registrar; solo `listPeriods` (PROVISIONAL) y las 5 excepciones sin llamada. `cd frontend && npm run build` sale con 0 (PWA con 82 entradas) y `git status` queda igual.
- **C3-04 · cumple.** `cd backend && go test ./... -count=1 -v` sale con 0: 32 `ok`, 0 FAIL, 1126 PASS, 0 SKIP, 216 de 216 `TestContrato_Operaciones`. `cd frontend && npx playwright test --reporter=json`: 245 ejecutadas, 241 pasan, 4 omitidas, 0 fallan y 0 inestables (311.7 s). Crucé los 60 nombres citados y todos pasan; `:146` se omite en móvil a propósito.
- **C3-05 · cumple.** `python3 -I c305.py`: 41 filas en sí, 4N = 164, 0 faltantes, 0 diferencias con la evidencia citada y PNG de 412 o 1440 px; la carpeta tiene 164 PNG.
- **C3-06 · cumple.** `grep -c "| no |"` da 1 (la fila Global), con el bloqueo `C3-capturas-global.md` citado.
- **C3-07 · no se firma** (veredicto devuelto). Al terminar, `sha256sum -c` de la matriz da OK: no cambió.
- **C3-08 · no cumple.** Cada archivo en una copia nueva (`bash pruebas-u5/flujos/servidor-prueba.sh arrancar`, luego `FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts pruebas-u5/flujos/<archivo>`): `flujos.spec.ts` da 6 passed y 1 failed (`:122`, F-3A, toque 1 → `/compras/productos`); `pendientes.spec.ts` da 4 passed; `f7-f8.spec.ts` da 2 passed. Son 7 de 8 flujos. F-1 se registra en 3 toques (K10). La base de desarrollo quedó en 6 gastos y 0 pesos antes y después.
- **C3-09 · cumple.** `npx playwright test -c pruebas-u5/flujos/playwright.config.ts pruebas-u5/flujos/c3-09.spec.ts` contra `finanzas.local:8443` da 1 passed, «111 cifras… diferencias != 0: 0»; `npx playwright test -c playwright.integracion.config.ts humo-inicio` da 1 passed.
- **C3-10 · cumple.** Leí 11 filas y §3 a §8: el texto está en es-MX y las fechas de evidencia en ISO.
- **C3-11 · cumple.** `grep -n "v1"` encuentra las líneas 3 y 163 a 165; `ls docs/estado/bloqueos/C3-*` da 3 bloqueos (más el de escalamiento que agregué). Los supuestos están en los reportes de los redactores.
- **U5-09 (de paso):** F-1 se registra en 3 toques en la copia nueva.

## Supuestos
1. Tomo K46 como respuesta de Axel al punto 3 de `C3-flujos`. K46 aprueba en bloque las opciones recomendadas de la puerta 2, y Q2-47 = A dice literalmente: «La API no reproduce +138.50 con la semilla». Por eso considero válida la prueba de la ruta B que compara la pantalla con la API.
2. Cuento como «pasa» la prueba citada `comidas-compras/correccion-bis.escritorio.spec.ts:146`, que pasa en escritorio y se omite a propósito en el proyecto móvil (`test.skip(isMobile, …)`).
3. Con N = 8 (supuesto 4 de la matriz), C3-08 exige 8 de 8. La corrida con precios de prueba agregados por API (F-3A3) no cuenta como «con la semilla».

## Bloqueos
- `docs/estado/bloqueos/C3-flujos.md` punto 2 (decide Axel), sin cambios.
- `docs/estado/bloqueos/C3-ronda3-escalamiento.md` (nuevo): no hay ronda 4 y decide Axel.

## Pendientes
- Después de la decisión de Axel: corregir R3-1 y hacer una verificación acotada de C3-08 con la firma (C3-07), si Axel la autoriza.
- Bajos R3-2 a R3-7 en el documento: cita `:141`, `getLiveDebt`, lista completa de pruebas, `expect(reg.ok)` en F-3A0, F-5 (punto 4) y `listPeriods`.
- Fuera de C3: corregir el ejemplo de F-3 en U1 §4 y §3.16 y en U3 tarjeta 5. También revisar el renglón duplicado «Bodega Aurrera · Vallejo» en «Totales por tienda» de P-13, que vi en F-3A3.
- El integrador debe agregar `frontend/pruebas-u5/flujos/` y `docs/design/capturas-u5/` al commit.

## Propuestas para Axel
| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. Opción (c) del punto 2: U1 aclara que la ruta A pide una lista con precios completos en al menos una tienda. Después, verificación acotada y firma | C3 se firma sin otra ronda completa | Bajo | Bajo. La tarjeta anuncia una canasta y abre otra (2b) hasta que se arregle |
| B. Opción (b): la semilla trae los precios de los ingredientes | F-3 pasa con la semilla tal cual | Medio | Hace falta una fuente para los precios |
| C. Firmar con F-3 ruta A como excepción con fecha | Se llega a la puerta 4 | Ninguno ahora | C3-08 queda con una excepción explícita |

Recomendación: **A**.

## Consumo aproximado
- Leído: unos 110 KB (C3.md, REVISION y correccion r2, extractos de la matriz, U1, `flujos.spec.ts`, bloqueos y salidas de las pruebas). Escrito: unos 17 KB. Acciones: unas 40. Tiempo de ejecución: unos 12 min (Playwright completo 5.2 min, flujos 1.5 min y Go menos de 1 min).
