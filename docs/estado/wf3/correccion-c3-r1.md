# correccion-c3-r1 · Corrección de C3, ronda 1 (WF-3b)

Fecha: 2026-10-12 · Agente: `correccion-c3-r1` (Opus) · Base: HEAD `6d691d0` con el árbol de trabajo de WF-3b · Revisión que se corrige: `docs/estado/wf3/REVISION-c3-r1.md` (H-1, H-2 y H-3, los tres de severidad alta, criterio C3-08).

## Hecho
- **H-2 (F-7, obligatorio (b) de U1-05):** escribí la prueba `frontend/pruebas-u5/flujos/f7-f8.spec.ts:16`. Corre contra el servidor real en una copia aislada (`servidor-prueba.sh`, `127.0.0.1:8444`), en el proyecto `movil-real` (412 × 915, táctil) y sin API simulada. **Pasa.**
  - Parte A, con la semilla tal cual: P-01 → ☰ → G-01 → Créditos → P-17 → «Simular compra» → P-20, monto 3094.50 → «Comparar escenarios». La API responde `missing_constraints` (faltan las 3 restricciones duras de P-26.05) y la pantalla muestra «Faltan restricciones · Sin capturar: …» con «Capturar», sin recomendación y sin «Usar escenario recomendado» (K16). «Capturar» abre P-26.05.
  - Parte B, el resto de la cadena de U1: tecleé restricciones de prueba solo para esta comparación (colchón 500.00, 3 compras vivas, deuda en cero el 2027-12-31). El comparador evalúa 16 escenarios válidos y da «Escenario recomendado · puntaje 88.01 de 100». «Usar escenario recomendado» abre P-19 prellenada sin guardar nada (4 compras antes y 4 después; vista previa $3,094.50 + $534.50 = $3,629.00). «Guardar compra» crea la compra (201, 2 pagos programados). Al final comprobé que las restricciones globales siguen en null.
- **H-3 (F-8, obligatorio (d) de U1-05):** escribí la prueba `f7-f8.spec.ts:139`: P-01 → «+» → P-03 → Ticket → «Capturar ticket» → P-05 → «Tomar foto» con la imagen sintética `datos/semilla/tickets-sinteticos/sintetico-01.png`. El OCR local (Tesseract spa, job `ticket_ocr`) propone 8 renglones. La prueba corrige 1 con el catálogo (HUEVO BLANCO 12 PZA → Huevo blanco, cartón 12 piezas), quita los 2 renglones del encabezado, elige la tienda y la categoría y toca «Confirmar ticket» → P-04. Después comprueba el precio guardado con `listProductPrices`: presentación 116, 4590 centavos, fuente `ticket_ocr`, Bodega Aurrera, 2026-10-12. También comprueba el gasto de 25153 centavos con `getExpense`. **Pasa.**
- **H-1 (F-3):** no lo puedo cerrar yo. La corrección pide que **Axel decida** los puntos 2 (precios de los ingredientes en la semilla, o una aclaración de U1) y 3 (cifra de Walmart de U1 o la del motor S5) de `docs/estado/bloqueos/C3-flujos.md`. Inventar precios para la semilla o elegir una de las dos cifras sería asumir. Lo que comprobé:
  - Reejecuté `flujos.spec.ts` en una copia nueva: 6 de 7, F-3 falla en el mismo paso (la tarjeta «Mejor canasta» lleva a «Productos y precios»).
  - Ese destino es el estado vacío que U1 §2.4 pide para P-01 («cada tarjeta se vacía por separado con texto guía y la acción que la llena»). No es un defecto del frontend que yo pueda corregir sin la decisión.
- `docs/contratos/C3-integracion.md` pasa a **v1.1, sin firma**. Cambié §3 (fila de flujos), §4 (introducción, filas F-7 y F-8 con pasos, toques, pantallas, prueba y resultado, y la línea de resultado: **7 de 8**), §6 (nueva ejecución) y §8 (historial). De paso cambié «sin mocks» por «sin API simulada» en la frase de §4 que reescribí (era H-8, de severidad baja). No toqué la matriz de §2.
- `docs/estado/bloqueos/C3-flujos.md`: el punto 1 queda resuelto con la opción (b) y agregué 2 puntos nuevos:
  - 6: la fecha de deuda en cero es «opcional» según U1 y la pantalla, pero obligatoria según C1 y el servidor.
  - 7: los toques de F-8 son más que los que pide U1.
- El registro de pasos de la última corrida quedó en `frontend/pruebas-u5/flujos/resultados/f7-f8.json`, con 6 capturas `F7-*.jpg` y `F8-*.jpg` en `frontend/pruebas-u5/flujos/capturas/` (revisé F7-faltan-restricciones y F8-gasto-creado: ninguna muestra la contraseña).

## Archivos
- `frontend/pruebas-u5/flujos/f7-f8.spec.ts` (nuevo, ~14 KB)
- `frontend/pruebas-u5/flujos/resultados/f7-f8.json` (nuevo; lo genera la prueba)
- `frontend/pruebas-u5/flujos/capturas/F7-faltan-restricciones.jpg`, `F7-recomendacion.jpg`, `F7-compra-prellenada.jpg`, `F7-compra-guardada.jpg`, `F8-renglones-propuestos.jpg` y `F8-gasto-creado.jpg` (nuevos). Mi corrida de `flujos.spec.ts` volvió a escribir las capturas `F1-*` a `F6-*`: son las mismas pantallas.
- `docs/contratos/C3-integracion.md` (v1.0 → v1.1, sin firma)
- `docs/estado/bloqueos/C3-flujos.md` (actualización, puntos 6 y 7)
- `docs/estado/wf3/correccion-c3-r1.md` (este reporte)

No cambié código de la aplicación (`frontend/src`, `backend/`). `npx tsc --noEmit -p frontend` sale sin errores. La configuración de Playwright de `frontend/` (`testDir: 'pruebas'`) no recoge `pruebas-u5/`, así que el conteo de C3-04 no cambia.

## Verificación
- **C3-08:** la cuento como no cumplida mientras F-3 falle.
  - Comando: `bash frontend/pruebas-u5/flujos/servidor-prueba.sh arrancar && cd frontend && FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts f7-f8`.
  - Resultado: **2 passed (33 s)** en 3 corridas seguidas, cada una en una copia nueva.
  - Toques: F-7A 5 (☰, Créditos, Simular compra, Comparar escenarios, Capturar); F-7B 4 más (Restricciones duras, Comparar escenarios, Usar escenario recomendado, Guardar compra); F-8 15 (con 1 corrección con el catálogo, 2 renglones quitados, tienda y categoría).
  - `flujos.spec.ts` en copia nueva: **6 passed, 1 failed** (F-3, ruta A, paso «Toque 1: tarjeta Mejor canasta», P-14 no visible; F-1 = 3 toques).
  - Total: **7 de 8** flujos de U1 pasan, con los 4 obligatorios de U1-05: (a) F-1, (b) F-7, (c) F-2 y (d) F-8. Falta F-3 (decide Axel).
- **C3-01:** `grep -c "^| P-" docs/contratos/C3-integracion.md` = **41** (sin cambio); el encabezado de §1 no se tocó.
- **C3-06:** `grep -c "| no |" docs/contratos/C3-integracion.md` = **1** (la fila Global, con bloqueo citado; sin cambio).
- **C3-11:** `grep -n "v1" docs/contratos/C3-integracion.md` → línea 3 «Versión v1.1 · 2026-10-12 · sin firma» y líneas 162 y 163 (historial v1.0 y v1.1). `docs/estado/bloqueos/C3-flujos.md` sigue existiendo y está al día.
- **C3-02 a C3-05, C3-09 y C3-10:** no los reejecuté porque no cambié la matriz de §2, el frontend, el backend ni las cifras de §5. Valen las corridas del revisor `revisor-c3-r1`.
- **C3-07:** no aplica a este agente; la firma la agrega el revisor.
- **U5-08:** las pruebas solo llaman a `https://127.0.0.1:8444`, con la CA de mkcert y sin red externa. No mandé nada por Telegram (el servidor de prueba corre con HOME aislado).
- **U5-09:** F-1 = 3 toques en mi reejecución de `flujos.spec.ts:13`.
- **Base de desarrollo intacta:** `datos/finanzas-dev.db` tenía antes y después 6 gastos, 4 compras a crédito y 72 precios (`sqlite3 -readonly`).

## Supuestos
1. F-7 parte B teclea restricciones de prueba **solo para esta comparación**: colchón 500.00, 3 compras vivas y deuda en cero el 2027-12-31. U1 N-225 permite cambiarlas así sin tocar las globales, y la prueba comprueba al final que las globales siguen en null. No son valores de Axel ni se guardan; con la semilla tal cual, la parte A prueba lo que U1 pide sin restricciones («Faltan restricciones · Capturar», K16).
2. El monto de F-7 es 3094.50, el total de la compra del 10-oct-2026 de la semilla.
3. En P-19 dejé vacíos los intereses. La pantalla dice «Vacío: se usa la tasa de la tabla por plazo de la línea», y la semilla trae la de 2 quincenas ($534.50). Así no invento una cotización de Kueski.
4. La «foto de ticket de prueba» de F-8 es la imagen sintética de `datos/semilla/tickets-sinteticos/`. `datos/semilla/tickets/` sigue sin las 5 fotos reales de Axel (S5-06), así que la prueba no mide la precisión del OCR.
5. En F-8 la tienda (Bodega Aurrera), la categoría (Despensa), la fecha (2026-10-12) y el total (251.53, la suma de los importes impresos del CSV esperado) son datos de captura de la prueba, porque el OCR no los propuso.
6. «Tomar foto» y aceptar la foto en la cámara cuentan 2 toques, como en U1 («toma la foto y la acepta»). En la prueba, la cámara es el selector de archivos de Playwright.

## Bloqueos
- `docs/estado/bloqueos/C3-flujos.md`:
  - **Puntos 2 y 3 (F-3):** siguen abiertos y deciden C3-08. Decide Axel.
  - **Punto 6 (nuevo):** fecha de deuda en cero, opcional (U1 N-225 y etiqueta de P-20) u obligatoria (C1 y servidor).
  - **Punto 7 (nuevo):** toques de F-8 por encima de U1. U1 pide 5; la pantalla real necesita como mínimo 8.
  - **Punto 1:** resuelto con (b).
  - Los puntos 6 y 7 no cambian el resultado de F-7 ni de F-8.

## Pendientes
- Volver a correr `flujos.spec.ts` (F-3) en una copia aislada cuando Axel decida los puntos 2 y 3. Si Axel elige que la semilla traiga los precios de los ingredientes, falta saber de dónde salen esos precios (semilla o agente externo, K53).
- Para el integrador: `frontend/pruebas-u5/flujos/` sigue **sin seguimiento en git**, con `f7-f8.spec.ts`, `resultados/f7-f8.json` y las capturas `F7-*` y `F8-*`. Hay que agregarlo al commit para que la evidencia de §4 quede en el repositorio.
- Observación sin verificar (no es hallazgo de C3): al confirmar el ticket, la respuesta trae `inventory_item_ids: []`, aunque U1 F-8 habla de «entradas de inventario». No revisé la regla que decide cuándo se crean.
- Los hallazgos bajos H-4, H-5, H-6 y H-7 de la revisión no estaban en mi encargo y quedan como estaban.

## Propuestas para Axel
| Decisión | Opción recomendada | Por qué |
|---|---|---|
| F-3 ruta A (punto 2) | Que la semilla no invente precios y que U1 aclare que la ruta A pide una lista con precios en al menos una tienda | El estado vacío que muestra hoy la tarjeta es el que pide U1 §2.4. Los precios reales llegan con el agente externo (K53). |
| F-3 ruta B (punto 3) | Elegir entre la cifra de U1 ($3,233.00, +$138.50) y la del motor S5 (+$264.50), y corregir el otro documento | Es la única diferencia de cifras que queda en los flujos |
| Fecha de deuda en cero (punto 6) | Decidir si es opcional u obligatoria y alinear U1, C1, servidor y etiqueta | Hoy la pantalla dice «opcional» y el servidor no recomienda sin ella |
| Toques de F-8 (punto 7) | Aceptar «pasa con diferencia» por ahora y encargar después «Ticket abre la cámara» y la categoría propuesta | No bloquea C3-08, que solo fija el tope de toques para registrar gasto |

## Consumo aproximado
- Leído: ~150 KB (revisión, bloqueo, C3 §1 y §3–§8, U1 §4 extracto, C3.md extracto, apoyo, configuración y partes de `flujos.spec.ts`, P20/P05/P17/P19 por partes, salidas de las pruebas, 3 capturas).
- Escrito: ~40 KB (prueba de 14 KB, JSON generado de 10 KB, cambios en C3 y en el bloqueo de ~7 KB, este reporte de ~9 KB).
- Acciones: ~52.
