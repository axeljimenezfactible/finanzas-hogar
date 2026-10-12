# WF-3b · u5-flujos · Reporte

Agente `u5-flujos` (Sonnet) · 12-oct-2026 · HEAD 6d691d0 «[WF-3a bis]» (sin commit; solo el integrador hace `git add/commit`). Criterios: C3-08, U5-09, U5-13, C3-09. No se corrigió la aplicación: las fallas van con pantalla y paso.

## Hecho
- Pruebas Playwright contra el servidor Go real (sin mocks), 412×915 táctil, `frontend/pruebas-u5/flujos/`: `flujos.spec.ts` (F-1…F-6 de U1 §4, con toques contados solo en controles, K10), `captura-sin-conexion.spec.ts` (U5-13, 7 pasos + 4 capturas), `c3-09.spec.ts` (111 cifras de pantalla contra `oraculo/model.json`), `apoyo.ts`, `playwright.config.ts`, `servidor-prueba.sh`.
- **F-1…F-6 y U5-13 corrieron en una COPIA aislada** de la base de desarrollo (`servidor-prueba.sh`: `sqlite3 .backup` → `datos/tmp/u5-flujos/finanzas.db`, mismo binario y frontend, `127.0.0.1:8444`, HOME aislado = sin Telegram) porque escriben datos. **C3-09 (solo lectura) corrió contra el servidor de desarrollo `https://finanzas.local:8443`.** La base de desarrollo quedó intacta (6 gastos, 0 pesos antes y después). La copia, su servidor y `datos/tmp/playwright-u5-flujos` se borraron. Contraseña leída de `~/.config/finanzas-hogar/dev-semilla.txt`; nunca impresa ni en capturas (las capturas son de Inicio, hojas y franjas).
- Comando: `bash frontend/pruebas-u5/flujos/servidor-prueba.sh arrancar && cd frontend && FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts flujos.spec` (y, con base nueva, `captura-sin-conexion`); `... c3-09` sin `FH_BASE`.

## Tabla flujo → pasos → P-xx → resultado (C3-08)
| Flujo | Pasos (toques acumulados) | P-xx | Resultado |
|---|---|---|---|
| F-1 Registrar gasto | «+» (1) → teclear 185.50 → chip Despensa (2) → Guardar (3) | P-01→P-03→P-01, P-02 | **OK, 3 toques**. Inicio: Despensa $1,197.64→$1,290.39 de $2,000.00 (+$92.75 = parte del hogar del reparto 50/50 de $185.50); tabla P-02: fila 12-oct-2026 · Despensa · $185.50 · Axel · Hogar $92.75 / Karen $92.75; API: 1 gasto de 18550 |
| F-2 Planear la semana | Comidas (1) → casilla Lun·Desayuno (2) → receta Avena con manzana y canela (3) → Generar lista (4) | P-01→P-10→P-13 | **OK, 4 toques** (= fórmula de U1: 1+1×(1+0+1)+1). Lista: 31 artículos, total estimado $618.00. No ejecuté «Copiar semana anterior» (opcional) |
| F-3 ruta B (compra del 10-oct) | Compras (1) → ⋯ (2) → Comparar tiendas (3) → selector de canasta (4) → «Compra del 10-oct-2026» (5) | P-01→P-13→P-14 | **Pasos OK, 5 toques** (= U1). Cifra: Vallejo $3,094.50 contra Walmart **$3,359.00 (+$264.50; +$100.50 en idénticos)**; U1 dice $3,233.00 (+$138.50): hallazgo H-4, ya conocido por semilla-c3 |
| F-3 ruta A (tarjeta Mejor canasta) | tarjeta (1) → Detalle (2) → Usar tienda (3) | P-01→P-14→P-13 | **FALLA**, no se llegó a P-14 con la lista (H-2 y H-3) |
| F-4 Crédito y calendario | tarjeta Próximo pago (1) → Calendario (2) → tocar el pago (3) | P-01→P-18→G-06 | **OK, 3 toques**. Resumen: disponible A $10,000.00 / B $6,905.50 «pendiente de confirmar», próximo pago $973.75; primera fila «Vie 30 oct 2026 · $973.75 · Hogar $598.82 · Karen $374.93 · pendiente»; vista previa: capital que se repone $773.63, intereses $200.12, partes Hogar y Karen pendientes. No ejecuté el paso 4 opcional (Marcar pago) |
| F-5 Reporte | ☰ (1) → Reportes (2) → «Generar reporte mensual» (3) | P-01→G-01→P-23 | **PARCIAL**: pasos 1–3 OK. No existe la fila «Mensual · sep-2026» (el centro dice «Aún no hay reportes (el primero se genera el día 1)»), así que el paso 3 de U1 se hizo con «Generar»; no verifiqué la generación ni la descarga del PDF (pasos 4–5) |
| F-6 Capturar peso | «+» (1) → pestaña Peso (2) → teclear 78.4 → Guardar (3) | P-01→P-03→P-22→P-01, P-21 | **OK, 3 toques**. «Peso guardado · Deshacer»; Salud ▸ Detalle: Axel · 12-oct-2026 · 78.4 kg · IMC «Sin dato» (no hay perfil de salud; Resumen pide «Completa el perfil de salud») |

## Hallazgos (para el integrador; no corregidos)
- **H-1 (alta) Inicio ▸ Detalle a 412 px desborda a 475 px.** La etiqueta «Reposición de capital pendiente de confirmar con Kueski» (`span.inline-flex w-max`, 442 px) ensancha la página; en el móvil emulado el viewport pasa a 475×1055, la barra inferior queda en y=1000 (fuera de los 915 visibles) y el «+» no se puede tocar sin desplazar. Lo vi al tocar «Detalle» y volver a «+». Pantalla: `/inicio?sub=detalle`. Arreglo sugerido: que la etiqueta haga salto de línea.
- **H-2 (media) F-3A con lista vacía.** La tarjeta dice «Vallejo $3,094.50 · $264.50 menos que Acueducto de Guadalupe · Tu lista está vacía; se compara la compra del 10-oct-2026», pero al tocarla P-14 abre con «Lista de compras» y dice «La canasta no tiene artículos»; no abre la compra que anuncia (prueba `F-3A0`).
- **H-3 (alta) F-3A tras F-2.** Con la lista generada del plan (31 artículos), `getDashboardBestBasket?week=2026-W42` da `covered_count 0 de 31`; la tarjeta dice «Aún no hay precios para comparar» y apunta a `/compras/productos`; P-14 dice «Sin precios para comparar». Al agregar a la lista (API, base de prueba) 3 productos con `best_price` seguía en `0 de 34`. Mientras, P-13 muestra «Total estimado $618.00» y totales por tienda: P-13 y P-14 no coinciden. Sin diagnosticar (puede ser cobertura de la lista o dato); impide «Usar tienda» (3 toques).
- **H-4 (decisión de Axel) cifra de F-3.** +$264.50 en pantalla contra +$138.50 de U1 (+$100.50 en idénticos sí coincide); método de S5, ya reportado por semilla-c3.
- **H-5 (alta) U5-13 paso 2 en pantalla.** Recargar sin red muestra «Iniciar sesión · Requiere conexión con el servidor de casa»: no se ve «Sin conexión» ni el contador, y no se puede capturar más hasta tener red. Las 2 capturas SÍ siguen en IndexedDB (`finanzas-hogar/cola`: gasto y peso). Incumple «recarga la página y el contador sigue en 2» (U5-13 paso 2; U1-15).
- **H-6 (baja) envío doble.** Al volver la red la cola hizo 2 `POST /expenses` con la misma clave (y 1 de peso); el servidor respondió con `Idempotent-Replayed` y quedó 1 gasto: 0 duplicados en BD, pero la petición sobra.
- **H-7 (info)** Plan: costo $0.00 y kcal «—» en las recetas (sin precios de ingredientes ni nutrición confirmada; ya declarado en la semilla).

## Verificación (una línea por criterio)
- **C3-08**: 6 flujos ejecutados (F-1, F-2, F-3, F-4, F-5, F-6). Pasan F-1, F-2, F-3B, F-4, F-6; **F-3A falla** (H-2, H-3); **F-5 parcial**. No cumple limpio: salida de `flujos.spec` = 6 pruebas pasan y 1 falla (F-3, por F-3A); detalle por paso en `frontend/pruebas-u5/flujos/resultados/flujos.json`.
- **U5-09**: 3 toques (`1. «+»; 2. chip Despensa; 3. Guardar`; teclear 185.50 no cuenta) ≤ 3; el gasto aparece en Inicio (tarjeta y Detalle) y en la tabla de P-02. Capturas `capturas/F1-hoja-antes-de-guardar.jpg`, `F1-inicio-con-el-gasto.jpg`, `F1-tabla-todos-los-gastos.jpg`.
- **U5-13**: 6 de 7 pasos literales; el paso 2 queda en 2a (cola persistente: OK) y 2b (contador visible tras recargar: **FALLA**, H-5). Paso 1: red cortada, «Sin conexión con el servidor de casa · 2 en cola», gasto en 3 toques y peso. Paso 3: `POST /api/v1/expenses` `Idempotency-Key: 0790b374-9129-47b6-ae5f-4ed8a55992bc` (36 car.) `captured_at: 2026-10-12T09:34:12-06:00`; `POST /api/v1/weight-entries` `Idempotency-Key: ff4761eb-ad09-4265-9aba-b1655365e253` `captured_at: 2026-10-12T09:34:12-06:00`; ambos anteriores al envío (`capturas/u5-13-peticiones.json`). Paso 4: contador a 0 y franja fuera. Paso 5: antes 1 gasto y 0 pesos, después 2 y 1 (1 nuevo de cada uno); BD: `expense.idempotency_key` y `weight_entry.idempotency_key` guardadas con su `captured_at`. Paso 6: misma clave y cuerpo → 200, `Idempotent-Replayed: true`, mismo id (7), sigue 1. Paso 7: con 500 simulado en `POST /expenses`, aviso «No se pudo enviar · 1 en cola · Revisar», la captura sigue en cola, el servidor no la recibió; al quitar el fallo llega 1 sola vez. 0 duplicados. 0 colas en otras pantallas: P-10 sin red tiene «Generar lista de compras» deshabilitado y el único almacén con capturas es `finanzas-hogar/cola` (gasto, peso). Capturas del estado «Sin conexión» en las 4 combinaciones: `capturas/U5-13-sin-conexion-{movil-claro,movil-oscuro,escritorio-claro,escritorio-oscuro}.jpg` (fondo claro rgb(250,250,251), oscuro rgb(34,35,37)); más `U5-13-sin-conexion-movil-claro-contador-2.jpg`, `U5-13-fallo-de-envio-aviso.jpg`, `U5-13-tras-recargar-sin-red.jpg`, `U5-13-otra-pantalla-sin-conexion.jpg`. Base de prueba: copia en `datos/tmp/u5-flujos/` (6 gastos·0 pesos al inicio; 8·1 al final), borrada al terminar.
- **C3-09**: **111 cifras leídas en pantalla, 0 diferencias** (`resultados/c3-09.json`; `c3-09.spec.ts` pasa en 8.7 s contra `finanzas.local:8443`): disponible A $10,000.00 / B $6,905.50 y mínimos $8,094.50 / $5,000.00 (Inicio ▸ Detalle); próximo pago 30-oct-2026 $973.75 con Hogar $598.82 y Karen $374.93; **P1–P6** (oct-2026…mar-2027) devengado, efectivo, diferencia acumulada, y partes hogar/Karen, más acumulados $14,362.48 / $13,103.78 (Inicio ▸ Tendencias); tasa a 4 quincenas 25.8685 %; las 4 despensas (capital, intereses, total: 3,094.50 + 800.50 = 3,895.00; y 4,000.00 + 1,034.74 = 5,034.74 en las despensas 2, 3 y 4); la serie de disponible A y B de los 15 eventos del crédito; **deuda viva al 10-abr-2027 $6,293.44** (Flujo y proyección, campo «Fecha de deuda viva»). Todo coincide con `model.json` al centavo.

## Supuestos
1. F-1…F-6 y U5-13 se corren en una copia de la base de desarrollo, no en la de desarrollo, para no ensuciarla (CONVENCIONES no fija una base de pruebas de flujos). Mismo binario y frontend.
2. Elegí la categoría «Despensa» en F-1 (de las 3 sugeridas) por afectar el presupuesto visible; el gasto se reparte 50/50 y el presupuesto cuenta la parte del hogar (verificado, no es defecto).
3. F-3A: la preparación de 3 artículos con precio se hizo por API (base de prueba) y está marcada como tal; no cambió el resultado.
4. En F-5 sustituí la fila «Mensual · sep-2026» (inexistente) por «Generar reporte mensual».
5. Capturas en JPEG calidad 45–50 para ocupar poco.

## Bloqueos
Ninguno nuevo en `docs/estado/bloqueos/`; H-3 queda para el integrador (¿cobertura de la lista o dato?) y H-4 ya depende de Axel.

## Pendientes
- Diagnosticar H-3 y repetir F-3A; verificar generación y descarga del PDF (F-5 pasos 4–5) y el paso 4 de F-4 (Marcar pago); probar «Copiar semana anterior».
- Tras corregir H-5, repetir `captura-sin-conexion` (el paso 2b debe pasar).
- Pasar el estado «capturado» de U5-11 no era mi alcance.

## Propuestas para Axel
- Decidir cuál cifra de F-3 vale (+$138.50 de U1 o +$264.50 del motor) para cerrar H-4.
- Pedir a WF-3 que corrija H-1 y H-5 antes de la revisión: con ellos la captura sin conexión y el «+» fallan en el teléfono.

## Paquetes instalados
Ninguno.

## Consumo aproximado
≈ 70 acciones de herramienta (rebasé el tope de 60 por depurar el desborde H-1, el servidor aislado y el reintento de U5-13) · ≈ 150 KB leídos · ≈ 75 KB escritos de código y reporte (más 22 capturas JPEG de Playwright, ver `du`).

---

## Continuación 1 (12-oct-2026) · solo lo pendiente del tope anterior

Hecho sin rehacer lo anterior (el reporte de arriba se conserva). Todo corrió en una COPIA nueva de la base de desarrollo (`servidor-prueba.sh arrancar`, `127.0.0.1:8444`, HOME aislado, Telegram apagado, casilla «Enviarlo por Telegram» verificada sin marcar). Al terminar: servidor detenido, `datos/tmp/u5-flujos` y `datos/tmp/playwright-u5-flujos` borrados; la base de desarrollo sigue en 6 gastos y 0 pesos. Archivos nuevos: `frontend/pruebas-u5/flujos/pendientes.spec.ts` (F-5, F-4 paso 4, F-2 paso 2, F-3 ruta A repetida), `h3-diagnostico.spec.ts` (sonda de H-3), `resultados/pendientes.json`, `resultados/h3-diagnostico.json` y 9 capturas (`F5b-*`, `F4b-*`, `F2b-*`, `F3A2-*`, `F3A3-*`). Cambio en `apoyo.ts`: `entrar()` ahora espera a que desaparezca el h1 «Iniciar sesión» (el login también es nivel 1; sin eso la llamada siguiente a la API podía salir con 401). Comando: `bash servidor-prueba.sh arrancar && cd frontend && FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts pendientes` → 4 pruebas pasan.

### Resultado por pendiente
| Pendiente | Pasos (toques acumulados) | P-xx | Resultado |
|---|---|---|---|
| F-5 pasos 3-5 (generar y descargar) | ☰ (1) → Reportes (2) → «Generar reporte mensual» (3) → «Generar» (4) → fila «Mensual PDF · oct-2026» (5) → «Descargar PDF» (6) | P-01→G-01→P-23→P-24 | **OK, 6 toques** (U1: 5 «para generar y descargar»; el toque extra es confirmar «Generar» en la hoja). Fila nueva «Mensual PDF · oct-2026 · Sin envío»; vista web `/reportes/mensual?mes=2026-10`; se descargó `reporte-mensual-2026-10.pdf`, 101,331 bytes, cabecera `%PDF-1.7` (se lee el archivo descargado). U1 pide la fila «Mensual · sep-2026»: no existe (el primero se genera el día 1); se generó oct-2026 (parcial) |
| F-4 paso 4 (Marcar pago) | tarjeta Próximo pago (1) → Calendario (2) → pago del 30-oct (3) → «Marcar pagada la parte de Hogar» (4) | P-01→P-18→G-06 | **OK, 4 toques** (U1: 6, con «Guardar»; aquí el marcado es inmediato con «Deshacer», más corto). Fila → «pagada / pendiente»; «Deshacer» la regresa a «pendiente». El disponible A $10,000.00 / B $6,905.50, la deuda viva $18,999.22 y el próximo pago NO cambian tras marcar la parte del hogar (ni recargando): ver H-10 |
| F-2 paso 2 (Copiar semana anterior) | Comidas (1) → ⋯ (2) → «Copiar semana anterior» (3) | P-01→P-10 | **OK, 3 toques** (= 1 + 2 de U1). `POST /meal-plan/copy-previous-week` → 201 `{"created_count":0}` (la semana anterior no tiene plan en la semilla, y la semana vigente ya tiene casillas); no se ve aviso en pantalla (H-11) |
| F-3 ruta A repetida (H-3) | tarjeta Mejor canasta (1) → Detalle (2) → «Usar Bodega Aurrera · Vallejo» (3) | P-01→P-14→P-13 | **OK, 3 toques** con una lista donde una tienda tiene todos los precios. Tarjeta: «Vallejo $1,309.00 · 26 de 26 artículos con precio»; P-13 queda con «Totales por tienda» y «Total estimado $1,309.00». Preparación por API, solo en la copia: 18 precios de $50.00 en Vallejo (ver diagnóstico) |

### H-3 diagnosticado (no es un defecto del servidor; sí hay dos de presentación)
- Con la lista generada del plan de la semilla (26 artículos; con la receta de F-2 son 31) **solo 8 de 26 productos tienen precio** (productos 1, 2, 10, 16, 17, 22, 23, 34, todos de la compra del 10-oct en Vallejo y Acueducto); **ningún precio de ingrediente de receta existe en la semilla** (H-7). `basket-comparison` da Vallejo $409.00 con 18 sin precio, Acueducto $415.00 con 18 sin precio, Bodega Aurrera y Walmart $0.00 con 26 sin precio: **ninguna tienda está completa** y S5 §7.4 define la ganadora como la primera tienda completa, así que `cheapest_store` no existe, `best-basket` devuelve `covered_count 0` y la tarjeta dice «Aún no hay precios para comparar». Con precios en todos los artículos (18 altas en la copia) todo funciona: F-3 ruta A en 3 toques. La causa de «0 de 31» es dato (falta de precios de ingredientes), no cálculo. Evidencia: `resultados/h3-diagnostico.json`. El agente anterior probó con 3 productos con precio, que no completan una tienda; por eso seguía en 0.
- **P-13 vs P-14 «no coinciden»**: P-13 muestra un «Total estimado» parcial (suma del mejor precio por artículo de los 8 con precio, ≈ $406 a $412 según paquetes; $618.00 con la receta de F-2) sin avisar que faltan 18 precios; P-14 no declara ganadora. No hay contradicción de cifras, pero P-13 presenta como total un monto parcial (H-3b).
- La tarjeta, con la lista sin tienda completa, lleva a «Productos y precios» (`/compras/productos`), no a P-14: es el diseño de «Registrar precios» (U3), así que la ruta A de U1 solo aplica con una lista completa en precios.
- **H-9 (media, nuevo)** En P-14 Resumen las tiendas sin ningún precio se muestran con «$0.00 · −$1,309.00 contra la más barata · −$1,287.00 en idénticos» (Bodega Aurrera y Walmart, «0 de 26 con precio · cobertura 0 %»): un $0.00 con diferencia negativa parece la más barata. Debería decir «sin precios». Pantalla `capturas/F3A3-p14-resumen.jpg`.

### Otros hallazgos nuevos (sin corregir)
- **H-8 (info/media)** P-24 Reporte mensual oct-2026: «Primera quincena … Real $1,197.64» y «Segunda quincena (16 al 31-oct) … Real $1,197.64»; el único gasto real es el del 10-oct (1.ª quincena), así que la 2.ª repite el total del mes. No diagnosticado.
- **H-10 (decisión)** U1 F-4 dice que «Marcar pago» «recalcula el disponible»; con la lectura B activa y «reposición pendiente de confirmar con Kueski», marcar la parte del hogar no mueve el disponible ni la deuda viva. Puede ser el comportamiento previsto (la reposición espera la confirmación); confirmar con Axel.
- **H-11 (baja)** «Copiar semana anterior» responde 201 con `created_count 0` y no se alcanza a ver aviso en pantalla (la prueba lo busca tras 0.5 s); con 0 casillas copiadas conviene un mensaje.
- **H-12 (info)** `shopping-list/items/{id}` solo admite PATCH: no hay forma de quitar un artículo de la lista (ni API ni la prueba lo encontraron).
- **H-1, H-2, H-4, H-5 y H-6 siguen igual**: el código del frontend/servidor de WF-3 no cambió desde la primera pasada (no se repitió `captura-sin-conexion` porque H-5 no está corregida).

### Verificación (continuación)
- **C3-08**: ahora F-1, F-2 (incl. copiar semana), F-3 (ruta A con lista completa en precios y ruta B), F-4 (incl. marcar pago), F-5 (generar, abrir, descargar PDF; fila sep-2026 no existe) y F-6 pasan en pantalla. Reservas: F-3 ruta A solo con precios completos (H-3); la cifra de F-3 es +$264.50 y no +$138.50 de U1 (H-4, decide Axel); F-5 usa oct-2026 en lugar de sep-2026. `resultados/pendientes.json` (por paso) y `resultados/flujos.json`.
- **U5-09**: sin cambio, 3 toques.
- **U5-13**: sin cambio, 6 de 7; el paso 2b sigue en falla (H-5), a repetir cuando se corrija.
- **C3-09**: vuelto a correr contra `finanzas.local:8443` (solo lectura): «111 cifras leídas en pantalla · diferencias != 0: 0», pasa en 8.3 s (disponible A 10,000.00 / B 6,905.50, pago 30-oct-2026 973.75, P1–P6, deuda viva 10-abr-2027 6,293.44).

### Supuestos (continuación)
1. Para completar F-3 ruta A se registraron en la COPIA 18 precios ficticios de $50.00 en Vallejo; son datos de prueba por API, no de la semilla ni de desarrollo.
2. Marcar pago se probó solo con la parte del hogar y se deshizo; no se marcó la de Karen.
3. El PDF se verificó por cabecera `%PDF-` y tamaño; no se inspeccionó su contenido página por página.

### Bloqueos
Ninguno nuevo; H-10 y H-4 dependen de Axel.

### Pendientes
- Repetir `captura-sin-conexion` y el paso 2b de U5-13 cuando WF-3 corrija H-5; repetir F-1 tras corregir H-1 (el «+» fuera de pantalla en Detalle).
- Revisar el contenido del PDF (fuera de alcance).

### Propuestas para Axel
- Decidir H-10 (¿marcar el pago del hogar debe mover el disponible?) y H-4 (+$138.50 de U1 o +$264.50 del motor).
- Cargar precios de los ingredientes de las recetas (agente externo, K53) o la tarjeta «Mejor canasta» quedará siempre en «Aún no hay precios» con la lista del plan.

### Consumo aproximado (continuación)
≈ 30 acciones de herramienta · ≈ 60 KB leídos · ≈ 20 KB escritos (pruebas y reporte) · sin paquetes instalados.
