# correccion-c3-r2 · Corrección de C3, ronda 2 (WF-3b)

Fecha: 2026-10-12 · Agente: `correccion-c3-r2` (Opus) · Base: HEAD `6d691d0` con el árbol de trabajo de WF-3b · Revisión que se corrige: `docs/estado/wf3/REVISION-c3-r2.md`, hallazgo R2-1 (alta, C3-08; es H-1 de la ronda 1, reabierto). Sin `git add` ni commit (eso lo hace el integrador).

## Hecho
- **R2-1 queda abierto en parte.** No lo puedo cerrar completo sin Axel, pero ahora depende de una sola de sus dos partes:
  - **Punto 3 de `C3-flujos` (cifra de la ruta B): ya hay decisión.** Encontré que K46 aprobó Q2-47 = A en la puerta 2 (bloqueo `S5-carrito-comparativa`, «Respuesta de Axel»). Con `basket = expense`, cada otra tienda se cotiza con los productos y las cantidades de la compra base, y «la API **no** reproduce +138.50 con la semilla». `S5-interfaz.md` §7.2 lo deja FIRME y descarta la opción B. La pantalla da +$264.50 y la API lo mismo; los idénticos dan +$100.50, igual que U1. Los $3,233.00 (+$138.50) de U1 1.2 §4 son el canon carrito contra carrito (`CompareCarts`, fixture `s5-comparativa-canon`), y U1 1.2 (11-oct) es anterior a la puerta 2. Cambié la prueba de la ruta B (`flujos.spec.ts`) para que compruebe que la pantalla da lo mismo que la API y que los idénticos dan +$100.50, citando K46. **La ruta B pasa** en 5 toques. Falta corregir el ejemplo de U1 §4 F-3 y §3.16 y el «Ejemplo trazable» de la tarjeta 5 de U3; eso no pide una decisión nueva, pero esos archivos no están en mi alcance.
  - **Punto 2 (ruta A): sigue abierto y es lo único que impide C3-08.** Lo dividí en dos caras en `C3-flujos`. (2a) Tras F-2, la lista no tiene precios de ingredientes. Por S5 §7.4 (Q2-60 = A), una tienda incompleta no puede ser la más barata, así que la tarjeta lleva a «Productos y precios», como pide el «Vacío» de U3. (2b) Es el síntoma de F-3A0: con la lista vacía, la tarjeta compara la compra del 10-oct-2026 (N-024), pero su acción abre la semana 42 vacía. Así lo escribe U3 en la tarjeta 5: «abre la semana según K38». No es un error de código sino un hueco entre U1 y U3. Dejé opciones con la PROVISIONAL primero y una recomendación. No inventé precios para la semilla ni cambié el enlace de la tarjeta, porque las dos cosas serían asumir.
- `docs/contratos/C3-integracion.md` pasa a **v1.2, sin firma**. Cambié el encabezado, la fila de excepciones de §3 (F-3 por la ruta A; punto 2, y el punto 3 contestado por K46), la introducción de §4, la fila F-3 y el «Resultado de C3-08», una fila nueva en §6 y el historial de §8. Las referencias a `flujos.spec.ts` quedaron al día porque mi cambio movió líneas: F-4 `:222`, F-5 `:254` y F-6 `:276` (antes `:204`, `:236` y `:258`). F-3 sigue en `:122` y F-3A0 en `:70`.
- `docs/estado/bloqueos/C3-flujos.md`: agregué la actualización de `correccion-c3-r2` (punto 3 contestado, punto 2 con sus caras 2a y 2b, observación sobre Q2-49) y la anoté en la columna de impacto de los puntos 2 y 3.

## Archivos
- `/home/deepsec/finanzas-hogar/frontend/pruebas-u5/flujos/flujos.spec.ts` (ruta B de F-3: el paso «Cifra esperada por U1» pasa a «Cifra FIRME (K46, Q2-47 = A)», con pantalla = API, y `regB` ya no se tolera a medias: `expect(regB.ok)`)
- `/home/deepsec/finanzas-hogar/frontend/pruebas-u5/flujos/resultados/flujos.json` (registro por paso de mi corrida en una copia nueva)
- `/home/deepsec/finanzas-hogar/frontend/pruebas-u5/flujos/capturas/*.jpg` y `u5-13-peticiones.json` (las regeneraron mis corridas, cada archivo en su copia; mismos nombres)
- `/home/deepsec/finanzas-hogar/docs/contratos/C3-integracion.md` (v1.2)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C3-flujos.md`
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/correccion-c3-r2.md` (este reporte)

## Verificación
- **C3-08 · no se cumple todavía (7 de 8 flujos; F-3 falla solo por la ruta A).** Corrí `bash pruebas-u5/flujos/servidor-prueba.sh arrancar && FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts pruebas-u5/flujos/flujos.spec.ts`, cada vez en una copia nueva. Resultado: **6 passed, 1 failed** (`:122`). Registros: F-1 OK con 3 toques, F-2 OK con 4, **F-3B OK con 5** (pantalla = API: $3,094.50 contra $3,359.00, +$264.50, +$100.50 en idénticos), F-3A FALLA (la tarjeta lleva a `/compras/productos`), F-3A0 FALLA en su registro (la prueba pasa: R2-5), F-4 OK con 3, F-5 OK con 3 y F-6 OK con 3. Con el mismo método, cada archivo en su copia: `pendientes.spec.ts` dejó todos sus registros OK (F-2b 3 toques, F-3A2, F-3A3 en 3 toques con lista completa en precios, F-4b, F-5b), y `f7-f8.spec.ts` también (F-7A, F-7B y F-8 con 15 toques). La base de desarrollo tenía 6 gastos, 0 pesos y 72 precios antes y después (`sqlite3 -readonly datos/finanzas-dev.db`).
- **C3-08 · cifra de la ruta B contra la API (lectura propia, solo lectura, `127.0.0.1:8443`, sesión de axel; la contraseña se leyó del archivo fuera del repositorio y no se imprimió):** `getDashboardBestBasket?week=2026-W42` dio `basket=expense`, `expense_id` 1, tienda 2, 309450, `difference_vs_second_cents` 26450 y 36 de 36. `getBasketComparison?basket=expense&expense_id=1` dio tienda 2 con 309450 (idénticos 279450) y tienda 4 con 335900 (+26450; idénticos 289500, +10050), 33 idénticos y 3 sustitutos. `shopping-list?week=2026-W42` tiene 0 artículos (es el caso 2b).
- **C3-08 · registrar gasto:** F-1 en 3 toques (K10).
- **C3-09 · sin cambio, se cumple:** `npx playwright test -c pruebas-u5/flujos/playwright.config.ts pruebas-u5/flujos/c3-09.spec.ts` contra `finanzas.local:8443` dio 1 passed, «111 cifras leídas en pantalla · diferencias != 0: 0».
- **C3-07:** v1.2 sin firma; el bloque lo agrega el revisor.
- **C3-10:** el texto nuevo está en es-MX y las fechas de evidencia en ISO.
- **C3-11:** `grep -n "v1" docs/contratos/C3-integracion.md` encuentra la línea 3 («Versión v1.2 · 2026-10-12») y las líneas 163 a 165 del historial. Bloqueos citados: `C3-flujos` (actualizado). Los supuestos y bloqueos están en este reporte.
- **C3-01 a C3-06:** no toqué §2. `docs/design/capturas-u5/` sigue con 164 PNG y ninguno es más nuevo que mi corrida.
- **U5-09:** F-1 en 3 toques en la copia nueva. **U5-13:** `captura-sin-conexion.spec.ts` dio OK en su copia, con 6 toques (corrida secundaria; regeneró sus capturas en `frontend/pruebas-u5/flujos/capturas/`).
- **Incidente de la corrida, sin efecto en la evidencia final:** la primera vez pasé el filtro `flujos`, que coincide con toda la carpeta, y corrieron los 6 archivos en **una sola** copia: 21 pruebas, 17 pasan y 4 fallan por datos cruzados entre archivos. Esa corrida no cuenta. La repetí archivo por archivo, cada uno en su copia nueva, y luego otra vez `flujos.spec.ts`, así que las capturas que quedan salen de corridas limpias. Los `resultados/*.json` no se tocan en una corrida porque `Flujo.guardar()` escribe en `datos/tmp/u5-flujos/`. Solo copié `flujos.json` de la última corrida limpia. La base de desarrollo no se tocó.

## Supuestos
1. Tomo K46 (Q2-47 = A) como respuesta del punto 3 porque su texto lo dice literalmente: «la API no reproduce +138.50 con la semilla». El cálculo que da +$264.50 (las cantidades de la compra base y los paquetes que las cubren, S5 §7.3) es el que describe la misma opción A. Si el revisor o Axel lo leen distinto, el punto 3 vuelve a abrirse y la prueba de la ruta B se cambia de nuevo.
2. La ruta B comprueba que la tienda más barata da $3,094.50 (= U1) y que en idénticos hay +$100.50 (= U1). El resto lo compara contra la API, no contra una cifra escrita a mano.
3. Usé `expense_id = 1` (la compra del 10-oct-2026 de la semilla), igual que el `selectOption('1')` que ya tenía la prueba.

## Bloqueos
- `docs/estado/bloqueos/C3-flujos.md` **punto 2** (caras 2a y 2b): decide Axel. Opciones: (a) PROVISIONAL, dejarlo como está; (b) que la tarjeta abra P-14 con la compra que anuncia; (c) que U1 aclare que la ruta A pide una lista con todos los precios en al menos una tienda. También tiene que decir, junto con el revisor, si para C3-08 vale una corrida con precios de prueba agregados a la semilla. Sin esa decisión, C3-08 no se cumple.

## Pendientes
- Corregir el ejemplo de F-3 en U1 1.2 (§4 «Resultado esperado» y §3.16 «Cifras canon que esta pantalla debe poder mostrar») y en U3 (tarjeta 5, «Ejemplo trazable»). Con K46 (Q2-47 = A), la pantalla muestra la cotización con las cantidades de la compra base y no +$138.50. Lo debe encargar el orquestador al dueño de U1 y U3; yo no puedo tocar esos archivos.
- Cuando Axel decida el punto 2, volver a correr `flujos.spec.ts` en una copia nueva y actualizar §4.
- Hallazgos bajos de r2 que no estaban en mi encargo y quedan como estaban: R2-2 (F-8 está en `f7-f8.spec.ts:141`, no en `:139`), R2-3 (`getLiveDebt` en lugar de `getDebtTrend` y supuesto 6), R2-4 (lista de pruebas truncada), R2-5 (F-3A0 pasa con su paso en FALLA; falta `expect(reg.ok)` o sacarla de la cuenta), R2-6 (F-5 con reserva) y R2-7 (`listPeriods`).
- El integrador debe agregar al commit `frontend/pruebas-u5/flujos/` y `docs/estado/bloqueos/C3-flujos.md` (hoy `??`). También conviene anotar en `docs/estado/bloqueos/INDICE.md` que el punto 3 está contestado; no lo edité porque pesa 80 KB.

## Propuestas para Axel
| Decisión | Qué pasa con cada opción | Recomendación |
|---|---|---|
| Ruta A de «Comparar canasta» con la semilla (`C3-flujos` punto 2) | (a) Como está: con la semilla la ruta A no termina y C3-08 no se cumple. (b) La tarjeta abre la compra que anuncia: se acaba la contradicción visible, pero «Usar tienda» sigue pidiendo una lista. (c) U1 aclara que la ruta A pide una lista con precios completos: con eso pasa en 3 toques. | (c), y además (b) para que la tarjeta y la pantalla muestren lo mismo |
| Cifra de la compra del 10-oct en pantalla | Ya la decidiste en la puerta 2 (Q2-47 = A): la pantalla da +$264.50, porque cotiza Walmart con las cantidades de Bodega. Si quieres ver +$138.50, sería una excepción (opción B o C, que cambia C1) | Dejar A y corregir el ejemplo de U1 y U3 |
| Precios de Walmart en la semilla | Q2-49 = A dice que la app arranca sin precios de Walmart, pero el orquestador los metió de forma PROVISIONAL. Sin ellos, «Comparar canasta» no tiene contra qué comparar la compra del 10-oct | Confirmar en la puerta 4 lo que hizo el orquestador (precios de Walmart sí) |

## Consumo aproximado
Unos 150 KB leídos (partes de U1, U3, S5-interfaz, la revisión, los bloqueos y la salida de las pruebas), unos 22 KB escritos (C3, bloqueo, prueba y este reporte) y unas 47 acciones. Archivos temporales en `datos/tmp/correccion-c3-r2/`, `datos/tmp/u5-flujos/` y `datos/tmp/playwright-u5-flujos/`, todos borrados al terminar.
