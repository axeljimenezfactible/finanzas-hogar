# semilla-dev · WF-3a bis (Sonnet) · 12-oct-2026

Overlay de **semilla de desarrollo** (C3-matriz punto 5, opción b): `finanzas semilla --desarrollo`.

## Hecho
- **HECHO** `backend/internal/semilla/desarrollo.go`: overlay idempotente y separado de la semilla real (request_id de audit_log `semilla-desarrollo`, etiquetas «[desarrollo]» en el reporte y «(desarrollo)» en nombres y notas). Exige la semilla real ya cargada.
- **HECHO** Los montos **no están escritos en el código**: se leen en ejecución de las fixtures del oráculo S2 y de ninguna otra parte:

| Dato en la base | Valor | Fixture (parche de `entradas.datos`) |
|---|---|---|
| Sueldo quincenal (ingreso, Axel, desde 15-oct) | 12,000.00 la del 15; 12,500.00 la del 30 | `s2-proyeccion-csv` · `planned_items[1]` (y `s2-liga-provisional`) |
| Presupuesto Despensa desde 1-oct-2026 | 2,000.00 | `s2-presupuesto-scope` · `budgets[1]` |
| Meta de ahorro «Meta de ahorro (desarrollo)» | 6,293.44 al 10-abr-2027, plan 449.53 al mes, inicio 1-oct | `s2-metas` · `savings_goals[2]` |
| Aportación a la meta (5-oct) | 3,146.72 (50.0 %, RN-18) | `s2-metas` · `savings_movements` de la meta 2 |
| Renta (solo con `--con-gasto-fijo`) | 8,000.00 mensual desde 20-oct, categoría Vivienda | `s2-proyeccion-csv` · `planned_items[2]` y `categories` |

- **HECHO** Precios de la compra del 10-oct ya sembrada: 36 de Bodega Aurrera (precio unitario exacto de cada renglón del ticket, fuente `ticket_manual`, ligados al renglón) y 36 de Walmart (carrito comparado de `datos/semilla/data.json`, fuente `manual`). Total 72, todos del 10-oct.
- **HECHO** Plan de comidas (dato de desarrollo, **sin fixture**): lunes 12 a sábado 17-oct, tiempo «comida», Axel y Ximena, 1 porción, 6 recetas del recetario (Espagueti con atún y elote; Arroz con pollo y zanahoria; Sopa de fideo con jitomate; Rajas poblanas con crema y elote; Espagueti con calabacita y queso manchego; Arroz rojo con huevo estrellado y plátano macho): las de tiempo «comida» con más ingredientes con precio. 12 entradas.
- **HECHO** Candado «nunca en producción» (PROVISIONAL, bloqueo punto 2) y bandera `--desarrollo` (más `--con-gasto-fijo`) en `backend/cmd/finanzas`. El único cambio fuera del paquete es `semilla()` en `main.go` y `devSeedRunner` en `cmd/finanzas/semilla.go`; en `semilla.go` del paquete solo se agregó `Options.FixturesDir`, `Options.DevFixedExpense` y el `request_id` configurable.
- **DECISIÓN PROVISIONAL** El gasto fijo (renta) **no** va por omisión: rompe las cifras P1 y P2 de model.json (+8,000.00 por mes en devengado y efectivo). Ver bloqueo punto 1.
- **HALLAZGO** (no del overlay) P3 a P6 de `getAccrualVsCash` no son las de model.json ni con la semilla real sola. Ver bloqueo punto 4.

## Archivos
- `backend/internal/semilla/desarrollo.go` (nuevo) · `backend/internal/semilla/desarrollo_test.go` (nuevo) · `backend/internal/semilla/semilla.go` (3 campos)
- `backend/cmd/finanzas/main.go` (banderas `--desarrollo`, `--con-gasto-fijo`) · `backend/cmd/finanzas/semilla.go`
- `docs/estado/bloqueos/WF3a-bis-semilla-dev.md` · este reporte
- Base de desarrollo `datos/finanzas-dev.db` rehecha desde cero (la anterior ya traía la renta de una primera carga; quedó en `datos/tmp/semdev/prev/` y se borró al cerrar). Sesiones previas se perdieron (solo desarrollo).

## Verificación
Reaplicado: `bash ops/dev-stop.sh` → `finanzas migrate` → `finanzas semilla` (6 cotejos «cotejo OK» del oráculo) → `finanzas semilla --desarrollo` (nuevos: 1 presupuesto, 1 partida, 1 meta, 1 aportación, 36+36 precios, 12 comidas) → segunda carga (todo «ya estaba», 0 nuevos) → `bash ops/dev-run.sh`. Curl con la CA de mkcert y sesión de Axel (contraseña leída del archivo fuera del repo, sin imprimir):

| Operación | Respuesta |
|---|---|
| `getDashboardMonthVsPlan` | oct-2026: real 119,764 de plan 200,000, 60 %, `normal`, quedan 80,236, 20 días |
| `getDashboardFortnightFlow` | 1–15 oct: ingreso 1,200,000 (esperado), gastado 0, comprometido 0, saldo 1,200,000 |
| `getDashboardWeekMeals` | 12–18 oct: 2 personas con 6 días de plan, 15 casillas sin planear, `cost_cents` 0 y `avg_kcal_per_day` 0 (bloqueo punto 5) |
| `getDashboardBestBasket` | compra del 10-oct: Bodega (tienda 2) 3,094.50, 36 de 36 con precio, 264.50 menos que Walmart (tienda 4) |
| `getDashboardSavings` | ahorrado 314,672; aportado en el mes 314,672; plan 44,953; meta 50.0 % (500 décimas), `in_progress` |
| `getDashboardNextPayment` | **30-oct $973.75** (97,375; hogar 598.82 y Karen 374.93), disponible B **$6,905.50** (690,550) sin cambios |
| `getAccrualVsCash` | idéntico a la línea base previa al overlay (comparación JSON exacta: `True`) |

Pruebas (`go test ./internal/semilla -run Desarrollo -v`): TestDesarrolloIdempotente, ExigeLaSemillaReal, SinFixturesNoInventa, Candado, MontosSalenDeLasFixtures, CotejoConElOraculo, ConGastoFijoMueveLosMeses, PlanDeComidas: 8 de 8 PASS. El cotejo con el oráculo comprueba: devengado y efectivo de los 12 meses iguales a los de la semilla real sola; P1 y P2 = model.json; deuda viva al 10-abr-2027 **6,293.44** = `horizonte.deuda_viva_10_abr`; tarjetas con las cifras de las fixtures. Paquete `semilla` 79.0 % de cobertura (la semilla real ya existente no cambió; las pruebas de ella pasan).
- `go test ./...` (backend): todo OK. `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: exit 0 («Verificación ampliada OK»). `make verify-wf3-go verify-wf3-secretos` y `verify-wf3-servidor`: exit 0 (secretos: 0 coincidencias; el servidor de desarrollo quedó corriendo en `:8443` con la base sembrada).
- No se envió Telegram; no hay credenciales en código, pruebas ni reportes.

## Supuestos
- «Hoy» de las tarjetas es el reloj del servidor (12-oct-2026); el overlay fija octubre de 2026 (plan de comidas del 12 al 17-oct, aportación del 5-oct), así que fuera de esa fecha las tarjetas de semana quedan vacías.
- Categoría «Despensa» de las fixtures = `groceries` de la base. «Vivienda» (fija) solo existe con `--con-gasto-fijo`.
- Reparto de la renta de la fixture (`parts: []`, resto al titular) = regla por defecto «Hogar» (`cost_share_rule_id` NULL).
- Walmart: fuente `manual` (la app no lee sitios, K53).

## Bloqueos
`docs/estado/bloqueos/WF3a-bis-semilla-dev.md` (7 puntos). Los que requieren a Axel o al orquestador: 1 (renta contra P1–P6), 3 (presupuesto de otras categorías), 4 (P3 a P6 de la semilla real), 5 (precios y nutrición para costo y kcal), 7 (Walmart en la semilla real).

## Pendientes y solicitudes fuera de alcance
- `ops/runbook.md`: agregar «`finanzas semilla --desarrollo` solo en desarrollo; la unidad systemd de producción lleva `Environment=FINANZAS_ENTORNO=produccion` y el comando se niega».
- `ops/dev-run.sh`: tras `finanzas semilla` en la rama `--sembrar`, agregar `"$BIN" semilla --desarrollo -db "$DB" -dir "$ROOT/datos/semilla"` (hoy solo se aplica a mano).
- WF-3b / integrador: decidir el punto 4 antes de C3-09 completo (P1–P6).

## Propuestas para Axel
- Dar presupuesto a Comida fuera y Vivienda (montos tuyos) para que «Mes contra plan» tenga más de una categoría.
- Que el agente externo cargue precios de los ingredientes de 6 recetas completas, para que la tarjeta de comidas muestre costo.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Acciones de herramienta: ~48. KB leídos: ~190. KB escritos: ~45 (código y pruebas ~33, documentos ~12).
