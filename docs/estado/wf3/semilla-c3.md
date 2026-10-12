# Reporte · semilla-c3 (WF-3b, 12-oct-2026)

Agente semilla-c3 (Sonnet). Alcance: `backend/internal/semilla/`, su prueba de cotejo, `docs/estado/bloqueos/WF3a-bis-semilla-dev.md` (puntos 4 y 7) y este reporte. Decisiones del orquestador PROVISIONALES hasta la puerta 4.

## Hecho
- **Punto 4 = (a).** La semilla real agrega los gastos 2 a 6 de `oraculo/fixtures/s2/_datos/semilla-oraculo.json` (`proyectados.go`): despensas 2 a 4 (4,000.00 a crédito, devengo a 2 meses) y 2 papeles de contado de 199.00 de Karen (categoría nueva «Papel»), y liga `credit_purchase` 2 a 4 a su gasto (`expense_id`). Cada gasto se coteja con `oraculo/model.json` (HALLAZGO y la carga se detiene si algo difiere). El esquema no distingue «proyectado» (`expense.kind` = fixed/variable/upcoming; `upcoming` significa «ya ocurrió o quedó comprometido»): va `variable` como en la fixture, con fecha futura, descripción «… (proyectado)» y notas «PROYECTADO … bloqueo WF3a-bis-semilla-dev, punto 4».
- **Punto 7 = sí.** Los 36 precios de Walmart de `data.json` (fuente `manual`, 10-oct-2026) entran a la semilla real (`precios.go`); el overlay de desarrollo ya no los carga (conserva los de Bodega por ticket).
- Idempotente: gasto por (fecha, categoría, importe, medio de pago); la liga solo si falta (audit_log `update`). Una base de la forma anterior se repara sola.
- Pruebas nuevas en `cotejo_test.go`; pruebas viejas ajustadas a las cifras nuevas; bloqueo actualizado; `finanzas migrate` + `semilla` + `semilla --desarrollo` sobre `datos/finanzas-dev.db` (borrada y recreada); servidor de desarrollo arrancado con `bash ops/dev-run.sh`.
- Puntos 1, 3 y 5 no se tocaron (puerta 4, pregunta a Axel).

## Archivos
- Nuevos: `backend/internal/semilla/proyectados.go`, `backend/internal/semilla/precios.go`, `backend/internal/semilla/cotejo_test.go`, `docs/estado/wf3/semilla-c3.md`.
- Modificados: `backend/internal/semilla/semilla.go` (pasos `proyectados` y `precios`, `Options.OracleSeedPath`, `audit()`), `datos.go` (`papel` de model.json), `desarrollo.go` (sin precios de Walmart, avisos), `semilla_test.go`, `desarrollo_test.go`, `docs/estado/bloqueos/WF3a-bis-semilla-dev.md`.
- **Fuera de alcance, necesario:** `backend/internal/api/comprasa/mejor_canasta_test.go` (la prueba `…SemillaSinPrecios` suponía la semilla sin precios; ahora `…SemillaConPreciosDeWalmart`). Ver Propuestas.

## Verificación
- **C3-09 (cifras = API = oráculo), semilla real sola, vía API con sesión:** `TestCotejoAccrualVsCashP1aP6ConModelJSON` — `getAccrualVsCash` oct-2026…mar-2027 = `periodos[0..5]` de model.json (R-45…R-50), devengado/efectivo en centavos: 194750/97375, 194750/194750, 251737/223243, 271637/271636, 251737/251738, 271637/271636; `is_canonical` true; diferencia acumulada P6 = 125870. Antes de este cambio P3 a P6 no coincidían.
- **C3-09 deuda viva:** `TestCotejoDeudaVivaAl10Abril` — `getLiveDebt?date=2027-04-10` = 629344 (6,293.44) y `Verify` del motor de crédito igual.
- **Servidor en marcha (curl con sesión, https://finanzas.local:8443 y https://127.0.0.1:8443):** `getAccrualVsCash` devuelve las mismas 6 filas y `getLiveDebt` total_cents 629344 (compras 3 y 4: 125870 + 503474). Contraseña leída de `~/.config/finanzas-hogar/dev-semilla.txt` sin imprimirla.
- **Comparativa de canasta (`TestCotejoComparativaCanastaBodegaVsWalmart`):** idénticos Walmart − Bodega = +100.50 (= data.json; ambos con la sexta lata de atún, +21.00). **Canasta completa: el motor S5 da +264.50, no +138.50** (HALLAZGO, no coincide con lo pedido): 138.50 + 21.00 (sexta lata de atún solo del lado de Bodega en data.json) + 105.00 (el motor compra los paquetes que cubren la cantidad: dos Coffee Mate de 640 g contra el pack de 780 g; data.json pone uno). La prueba deja escritas las dos cifras. 36 precios de Walmart sembrados.
- `go test ./internal/semilla/...` ok; `go test ./...` en backend: todo ok; `gofmt -l` y `go vet` limpios.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT 0 («Verificación ampliada OK»).
- Idempotencia: `TestSembrarDosVecesNoDuplica` (0 nuevos, mismas filas), `TestProyectadosSeLigan` (reparación y segunda carga sin nada nuevo); CLI: segunda corrida de `finanzas semilla` sin renglones nuevos.
- U5: sin evidencia propia (este agente no hizo capturas).

## Supuestos
1. «Proyectado» = `variable` + fecha futura + marca textual (el esquema no tiene campo). No se usó `upcoming` por su significado.
2. Categoría «Papel» (código `papel`, kind `variable`, regla por defecto Karen 100 %, orden 3) creada en la semilla real; el nombre sale de la fixture.
3. Papeles de contado: `payment_method` cash, `paid_by` Karen, parte única de Karen 199.00, como en la fixture.
4. Solo se cargan los precios de Walmart (lo que decidió el orquestador); los de Bodega siguen en el overlay porque la comparativa usa el ticket como canasta base.
5. Al sembrar de cero la compra a crédito se inserta sin gasto y luego se liga: deja 3 filas `update` en audit_log (aceptado).
6. `Options.OracleSeedPath` es obligatoria fuera del repositorio (carpeta copiada a `datos/tmp`): sin el archivo la carga falla en lugar de omitir los gastos.

## Bloqueos
- `docs/estado/bloqueos/WF3a-bis-semilla-dev.md`: puntos 4 y 7 anotados con la decisión del orquestador y el hallazgo de la comparativa; puntos 1, 3 y 5 siguen como PREGUNTA a Axel en la puerta 4; punto 2 (candado) sigue PROVISIONAL.

## Pendientes
- Axel decide cuál es la cifra canon de la comparativa del 10-oct (motor S5 +264.50 o manual +138.50) y si la lista de Gastos debe mostrar las 5 filas proyectadas.
- Las capturas U5 de Inicio y Gastos se toman con la semilla nueva; la tarjeta «Mejor canasta» ya no sale vacía.
- Reconstruir `backend/bin/finanzas` cambia solo textos desde el arranque (avisos del overlay); no hace falta reiniciar el servidor para las cifras.

## Propuestas para Axel
- Si quiere distinguir lo proyectado en datos y no en texto, agregar una columna (p. ej. `expense.is_projected`) en una migración de WF-4; hoy se reconoce por fecha futura y nota.
- Decidir la regla de comparación de la despensa (paquetes que cubren la cantidad, como hace S5, o la comparación manual por artículo) para que C3-09 y la pantalla digan la misma cifra.
- La prueba de `comprasa` (`mejor_canasta_test.go`) la toca este agente por consecuencia de la decisión 7; pedir al integrador que la revise.

## Consumo aproximado
~60 acciones de herramienta; ~200 KB leídos (código de la semilla, pruebas, esquema, fixtures y extractos); ~45 KB escritos (código, pruebas, bloqueo y reporte). Temporales de `datos/tmp` borrados.
