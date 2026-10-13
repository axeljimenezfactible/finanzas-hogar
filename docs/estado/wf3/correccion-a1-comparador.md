# correccion-a1-comparador · WF-3c bis · 12-oct-2026

Corrige el hallazgo **A-1** (alto) de `docs/estado/wf3/REVISION-backend-cierre-r2.md`: desde M-1, `compareScenarios` recibía todas las líneas, también las archivadas. Por eso generaba, puntuaba y **recomendaba** escenarios en una línea archivada. No toqué motores, frontend, `docs/design/`, git ni el servidor de desarrollo.

## Hecho
- **`backend/internal/api/creditos/motor.go`.** Agregué `(*snapshot).activeStates()`, que arma los `CreditState` solo de las líneas con `credit.is_active` = 1. `states` no cambia: sigue contando todas las líneas para `getCreditTotals` (K-3, P-17) y `getDebtTrend` (G-5), como pide RN-30 (PROVISIONAL, decisión del orquestador). En su comentario dice que el comparador usa `activeStates`.
- **`backend/internal/api/creditos/condiciones.go:676`.** `comparador.Request.Credits` pasa de `s.states(nil)` a `s.activeStates()`. Así cumple con el contrato `S3-interfaz.md` §8 (`Request.Credits`, «créditos activos»), con `diccionario-datos.md:776` (`credit.is_active`: «deja de ofrecerse en capturas nuevas») y con U1 N-053. Ninguna fórmula del motor cambia, solo qué líneas recibe.
- **Prueba nueva `TestOp_compareScenarios_LineaArchivada`** (`backend/internal/api/creditos/creditos_test.go`). Reproduce el caso del revisor con dos líneas revolventes con límite, las dos con la despensa del 10-oct (R-13). La segunda tiene una tasa a 4 quincenas más baja (40000/309450).
  1. Con solo Kueski, guarda el resultado base.
  2. Con las dos activas, comprueba que el comparador ofrece la segunda y la recomienda, como en el caso del revisor.
  3. Archiva la segunda y comprueba tres cosas. Ningún escenario `valid` ni `rejected` tiene su `credit_id`. La recomendación no está en ella. El resultado es igual al base: mismo `evaluated_count` (3), misma recomendación con el mismo `score_bp` y la respuesta idéntica byte por byte.
  4. Con la línea ya archivada, comprueba que `getCreditTotals` sí la cuenta: `live_debt_cents` 779000 (2 × 389500), igual a `getLiveDebt` del 12-oct, y `month_total_cents` 194750 (2 × 97375).
- **Contraprueba.** Volví a poner `s.states(nil)` en `condiciones.go`, solo para este chequeo. La prueba nueva falló con «la línea archivada 2 sigue en 3 escenarios (válidos o descartados)». Después dejé otra vez la corrección.

### Qué operaciones usan `states` (lo pidió el revisor)
Lo revisé con `grep` en `backend/internal/api/creditos/`. El comparador solo se llama en `condiciones.go`.

| Operación | Qué líneas usa | Criterio |
|---|---|---|
| `getCreditTotals` | `states(nil)`: todas | RN-21 y RN-30, M-1 (PROVISIONAL `S6-tablero-lineas-archivadas`, punto 5) |
| `getDebtTrend` | `states(ids)`: todas, o las nombradas aunque estén archivadas | igual que la anterior |
| `PaymentsByMonth` (barómetro) | su propio ciclo sobre todas las líneas | RN-30, opción a1 del mismo bloqueo |
| `compareScenarios` | **`activeStates()`: solo activas** | S3-interfaz §8 `Request.Credits`; `diccionario-datos.md` `credit.is_active` |
| `previewCreditPurchase`, `previewAvailability`, `getAvailabilitySeries`, `getCardStatus`, `listCredits` y `getCredit` | **no usan `states`**: cada una trabaja con el crédito que se pide por id (`s.state(c)`) | sin cambio (ver Supuestos) |

## Verificación
Todo desde `/home/deepsec/finanzas-hogar`, después de la corrección.

```
$ cd backend && gofmt -l .                 → (sin salida), exit 0
$ go vet ./...                             → (sin salida), exit 0
$ go test ./internal/api/creditos/... ./internal/api/tablero/... -count=1 -cover
ok  finanzas-hogar/internal/api/creditos  2.827s  coverage: 82.6% of statements
ok  finanzas-hogar/internal/api/tablero   1.175s  coverage: 87.8% of statements
$ go test ./internal/api/creditos/ -run TestOp_compareScenarios -count=1 -v
--- PASS: TestOp_compareScenarios (0.20s)
--- PASS: TestOp_compareScenarios_LineaArchivada (0.10s)

$ make verify-wf3-go                       → EXIT 0
   go build ./... y go vet ./...: sin errores
   gofmt -l .: sin archivos por formatear
   go mod tidy -diff: sin cambios
   go test ./...: 33 paquetes OK, 5 sin pruebas

$ make verify-wf3-contrato                 → EXIT 0
   RESULTADO: los 22 archivos coinciden byte a byte con el generador
   S6-03: 221 operaciones en 176 rutas, 0 diferencias
   --- PASS: TestRutasIgualesAlContrato
   --- PASS: TestContratoSinSesionYEsqueleto
```

La cobertura de `api/creditos` pasa de 82.5 % (revisión r2) a 82.6 %. No hay endpoints ni tipos nuevos, y `openapi.yaml` y `types_gen.go` no cambian.

## Supuestos
1. **El comparador ya no ve las compras vivas de las líneas archivadas.** S3-interfaz §8.4 calcula la serie de pago parejo con «compras existentes no canceladas de todos los créditos de `Credits`». Como `Credits` ahora tiene solo las líneas activas, los pagos que aún quedan en una línea archivada no entran en esa serie. Tampoco entran en `max_fortnight_payment_cents` ni en el pico por persona (`OverCap`). Seguí el contrato tal como está escrito («créditos activos»), que es la misma conducta que había antes de M-1. Si Axel quiere que esos pagos sí cuenten en el pago parejo sin que la línea se ofrezca, hay que cambiar S3-interfaz para separar «líneas que se ofrecen» de «compras que pesan en la serie». Lo dejo como pregunta y no lo implementé.
2. **`previewCreditPurchase`, `previewAvailability` y `getAvailabilitySeries` quedan como estaban.** No usan `states` y trabajan con el `credit_id` que se pide. Leer el disponible o la serie de una línea archivada es consultar su historia, y el diccionario dice que la línea «se conserva para la historia». Que `previewCreditPurchase` (y `createCreditPurchase`, que el revisor ya anotó en su §4) rechacen una línea archivada no está claro: ni C1/`openapi.yaml` ni S3-interfaz lo dicen, y el diccionario solo dice «deja de ofrecerse» (selectores y comparador). Por eso no agregué ninguna guarda.
3. La prueba usa una tasa más baja en la línea archivada para que el comparador la recomiende mientras está activa (el caso de la revisión). No es una cifra canon.

## Pendiente (fuera del alcance de esta tarea)
- El revisor sugiere que `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md:55` diga qué operaciones usan `states`. No lo edité. La tabla de arriba trae el dato para quien lo integre.
- Reiniciar el servidor de `:8443` con el binario nuevo (sigue pendiente desde `correccion-backend-cierre-r1.md`).
