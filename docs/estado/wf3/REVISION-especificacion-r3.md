# REVISION-especificacion-r3 · Revisión de la especificación del tablero (WF-3c, K54), ronda 3 de 3 · 12-oct-2026

Revisor independiente `revisor-especificacion-r3`. No corrijo. Reviso el trabajo de `u3-v2` y `s1-c1-v15`, con las correcciones `correccion-especificacion-r1` y `-r2`. Los evalúo solo contra:
- K54 (`finanzas-hogar-orq/k54-tablero.md`);
- `docs/criterios/U3.md`, `U2.md` (DESIGN), `C1.md` y `S1.md`;
- las decisiones §A–§K.

Volví a correr cada verificación con mis propios comandos, en `datos/tmp/rev-esp-r3/`. Esa carpeta ya está borrada.

**Veredicto: NO APROBADO. Quedan 0 altos, 1 medio nuevo (M-R3-1) y bajos.**
- Los dos hallazgos de la ronda 2 están **cerrados**:
  - A-1 / H-1 (alto): `make verify-all` da EXIT 0.
  - M-N1 (medio): la fila M-01 de C1-datos §9 ya trae los 8 campos.
- El medio nuevo es una inconsistencia entre U3 G-6 y C1 M-21. Se corrige con una línea en U3 o con un campo en C1.

## 1. Resultado por punto de la tarea

| Punto | Resultado | Evidencia propia |
|---|---|---|
| (a) K54: palabras de Axel y las 4 decisiones | **CUMPLE** (sigue el bajo B-r1-1) | Un script compara cada cita «…» de `k54-tablero.md` con la línea K54 de `decisiones.md`. Las dos citas de Axel (658 y 293 caracteres) son idénticas. Las decisiones (1) a (4) están. `git diff --numstat` de `decisiones.md` da 1 línea agregada y 0 borradas, así que §G se conserva. No se registra la cita de la lentitud (92 caracteres). |
| (b) U3 2.0 cubre lo que pidió Axel; cada bloque tiene su operación en C1 1.5 con campos coherentes | **CUMPLE con M-R3-1** | Lo que pidió Axel está: ingreso con extra (K-2), créditos globales (K-3 y G-2), recordatorio de N pagos (K-4), barómetro con dígito y umbrales 30/40 (B-1) y 6 gráficas (G-1 a G-6), más el desglose en cadena (§2.2) y la señal subir/recortar (G-3). Un script cruza los campos de cada bloque con los extractos de sus operaciones y no encuentra faltantes. La excepción es G-6, que usa `scale_max_pct`: ese campo solo existe en `DashboardDebtCapacity` (B-1), no en `DashboardMonthlyOverview` (ver M-R3-1). |
| (c) Ninguna cifra inventada | **CUMPLE** | Las 3 cifras nuevas de U3 frente a HEAD (1,482.57, 1,767.51 y 2,232.43) están en `oraculo/model.json`, en `periodos[2]`. Las de M-01, M-20 y M-21 de C1-datos también: 1,197.64, 598.82, 1,197.64 / 1,197.64 y 1,482.57. Recalculé los ejemplos ilustrativos: RN-30 da 29.94 → 30, 37.06 → 37, 117.43 y 55.81 → 56; RN-31 da 113, 112, 117 y 110. El 27 % de DESIGN y U3 lleva la marca SINTÉTICO. Los 6 bloqueos tienen qué falta, opciones e impacto. |
| (d) DESIGN §8.9: solo tokens existentes y contraste en los dos temas | **CUMPLE** (sigue el bajo B-N5) | Los 51 tokens que citan las líneas nuevas de DESIGN y U3 existen en `C2-tokens.json`. No hay hex, `rgb()` ni `oklch()` nuevos. Calculé con mi propia implementación los 14 pares C-159 a C-172: resuelvo las referencias de C2, aplico WCAG 2.x y trunco a 2 decimales. Coinciden exactamente en claro y en oscuro; el menor de interfaz es C-170 con 3.46 en claro. Ninguno queda bajo su mínimo. |
| (e) Migración 00005 | **CUMPLE** | Con goose v3 en una base temporal: `up` 0 (5 claves 30/40/5/3/10 y 2 disparadores) → `down` 0 (0 claves, 0 disparadores) → `up` 0 → `down-to 0` 0 (0 objetos) → `up` 0 (versión 5). El disparador rechaza `0` y `21` y acepta `12`. Los 5 `@param` tienen sus 7 campos, más `@param_rango`. `gen_dict.py` a un archivo temporal sale idéntico (`cmp`) a `diccionario-datos.md`, versión 1.5. |
| (f) `make verify-all`, conteos de C1 y `verify_c1 --db` | **CUMPLE** | `make verify-all GOOSE=…` da **EXIT 0** («Verificación ampliada OK»). C1: los 22 archivos coinciden byte a byte y `verify_c1 --db` da «comprobaciones OK: 33, TODO OK». Mis conteos sobre `openapi.yaml` y los extractos: 221 operaciones, 176 rutas (174 `/api` + `/health` + `/openapi.yaml`) y 267 esquemas, igual que `INDICE.md`. La versión es `1.5.0`, en la línea 7. `make -k verify-wf3-contrato verify-wf3-secretos` da rc 0 («221 operaciones en 176 rutas, 0 diferencias»). |
| (g) Criterios | **CUMPLE** | Un script contra `git show HEAD` cuenta 261 criterios (201 altos / 60 medios / 0 bajos) y 0 archivos con cambios de id o de severidad. `make verify-criterios` da rc 0. «6 tarjetas» solo aparece en citas históricas o de trazabilidad (§G, plan §3.5, listas de cambios): ninguna lo pone como regla vigente. |

## 2. Hallazgos de la ronda 2

| Id r2 | Estado | Comprobación propia |
|---|---|---|
| A-1 / H-1 (alto) · `verify-all` roto, 5 operaciones en 501 | **CERRADO** | Existen `backend/internal/api/tablero/` con `Register` y las 5 operaciones, `areas.go` con `tablero.Register` y el archivo nuevo `creditos/tablero.go`. Las 4 pruebas pasan con `go test -count=1 -v`: `TestAreas_SinOperacionesEn501` («221 de 221»), `TestContrato_Operaciones` («221 de 221 operaciones cumplen»), `TestContrato_Rutas` («sin implementar 0») y `TestPermisos_MismasOperacionesAmbasCuentas`. `go test -cover` da 71.2 % en `api/tablero` y 80.4 % en `api/creditos`. `make verify-all`: EXIT 0. |
| M-N1 (medio) · K-1 sin la fórmula de los 8 campos en C1-datos §9 | **CERRADO** | La línea 404 de `C1-datos.md` (fila M-01) ya trae las fórmulas de `cash_cents`, `previous_*`, `*_difference_cents` y los dos conteos, con el periodo (mes y mes anterior), las tablas y el ejemplo del oráculo. Comprobé cada cifra en `model.json`: `periodos[0]` da 1197.64 / 598.82 y `periodos[1]` da 1197.64 / 1197.64. U3 K-1 cita esa fila. |

**Cifras fuera del motor.** El código nuevo no recalcula cifras de motor. El ingreso sale de `mf.FortnightFlow`, el devengado de `mf.PlanVsActual`, el efectivo y el % usado de `mf.BudgetMonth`, y los pagos de `mc.UnifiedCalendar`. `PaymentsByMonth` suma las partes igual que `getCreditTotals` (`creditos.go:270-285`). Las divisiones nuevas de RN-30 usan `dominio.PercentOf` y `MulDivRound`, con ROUND_HALF_UP. Es lo que C1 y U3 §2.6 declaran («Motor: S2 (ingreso) y S3 (pagos)»).

## 3. Hallazgos

### Medio
- **M-R3-1 · G-6 usa un campo que su operación no trae (U3-02, C1-09; punto (b)).**
  - U3 G-6 (línea 229) dice: «Eje Y de 0 % al mayor entre `scale_max_pct` y el mayor `debt_ratio_pct` a la vista».
  - G-6 se alimenta solo de `getDashboardMonthlyOverview`. Su esquema `DashboardMonthlyOverview` (`extractos/inicio.yaml`, de 1499 a 1573) trae solo `warning_pct`, `limit_pct` y `months[]`. `scale_max_pct` existe únicamente en `DashboardDebtCapacity` (B-1, `inicio.yaml:1469`).
  - Además, U3 §2 pide que cada bloque haga su propia petición: G-1 y G-6 comparten una sola, y si B-1 falla, G-6 no debe depender de él. También pide que la interfaz no calcule (la regla 1.5 × tope sería un cálculo).
  - Tal como está, la PWA tendría que pedir B-1 o calcular la escala.
  - **Corrección, una de dos:**
    - (A, sin tocar el contrato) En U3 G-6, cambiar a «Eje Y de 0 % al mayor entre `limit_pct` y el mayor `debt_ratio_pct` a la vista».
    - (B, recomendada porque mantiene la misma escala que B-1) Agregar `scale_max_pct` a `DashboardMonthlyOverview`:
      - en `c1gen` (`gen2_schemas.py` / `c1_plantilla.md`), como entero requerido, `x-provisional`, con la regla de M-20;
      - regenerar con `build.py`;
      - llenarlo en `tablero.go` con la misma expresión de `getDashboardDebtCapacity`, mejor en una función compartida;
      - y citarlo en U3 §2.6 y en M-21.

### Bajos
- **B-R3-1 · Cuenta errónea en la justificación de la escala del arco.** `docs/estado/bloqueos/U3-tablero-barometro.md`, punto 4, opción (a), dice que la escala a 60 % «da a la zona roja el mismo ancho que la verde y la amarilla juntas hasta 40». Pero el rojo va de 40 a 60, que son 20 puntos, y el verde con el amarillo van de 0 a 40, que son 40. Axel decide en la puerta 4 con ese argumento. **Corrección:** «da a la zona roja la mitad del ancho de la verde y la amarilla juntas», o quitar la frase.
- **B-R3-2 · Ningún manejador llena los 14 campos opcionales nuevos de C1 v1.5.** `grep` en `backend/` fuera de `types_gen.go` da 0 resultados para `PreviousActualCents`, `CategoriesExceededCount`, `FortnightsDifferenceCents`, `MonthPaymentCount` y `MonthPendingHouseholdCents`.
  - El contrato pasa porque los campos son opcionales, y la corrección r2 lo lista en Pendientes.
  - Pero con el servidor actual K-1, el indicador principal de K54, no mostrará el efectivo ni la comparación con el mes anterior. K-2 y K-3 también pierden sus líneas nuevas.
  - Está fuera de U3, C1 y S1: le toca al backend (S6), antes del frontend P-01 y de las capturas U5.
- **Bajos de la ronda 2, siguen abiertos (no se corrigieron en esta ronda):**
  - B-N1: U3, línea 28, dice «M-01 a M-17». §15 sigue dando como pendiente la viñeta de U1, que ya se hizo.
  - B-N2: el nombre accesible de G-3 lleva «3 meses» y «10 %» fijos (línea 213).
  - B-N3: K-3 no cita M-19.
  - B-N4: «esquema 1.4» en `Makefile` (líneas 8, 66 y 149) y en `s1_verifica.py:190`. También sale en la salida de `verify-all`.
  - B-N5: en DESIGN, línea 1082, la aguja va «del centro al arco de escala» y ese par no está en §10.5.
- **Bajos de la ronda 1, siguen abiertos:**
  - B-r1-1: K54 no registra el «segundo problema» (la lentitud) ni la cita de Axel.
  - B-r1-3: `diccionario-datos.md` no tiene «Cambios en 1.5» y su línea 3 sigue con 11-oct-2026.
  - B-r1-4: `docs/estado/bloqueos/INDICE.md` no lista los 6 bloqueos `U3-tablero-*` y `C1-tablero`.

## 4. Criterios evaluados (resumen)
- **U3:**
  - CUMPLEN U3-01 (tablero con 4 K, B-1, 6 G, §2.2 y §2.3), U3-03 (estados de §7), U3-04 (§2.4, 412 y 1440), U3-05, U3-06, U3-07 (lectura activa B y «Pendiente confirmar»), U3-08, U3-09, U3-10 (§6 con G-6), U3-11 y U3-12.
  - U3-02 cumple **salvo M-R3-1**.
- **U2 (DESIGN 1.1):** cumple. Hay tokens existentes, 0 colores nuevos y los 14 pares están calculados en los dos temas. Quedan abiertos B-N5 y la numeración §8.9 contra el §8.5 de K54; esa diferencia está declarada PROVISIONAL en `U3-tablero-barometro`, punto 6.
- **C1:**
  - C1-02 cumple: `verify-all` y el contrato pasan.
  - C1-09 cumple: M-01 está ampliada y M-18 a M-22 existen.
  - Sale M-R3-1: M-21 no da la escala que G-6 usa.
  - C1-13 cumple: 1.5.0.
- **S1:** cumple. La 00005 hace up/down/up, los `@param` están completos, el diccionario sale idéntico y `verify-s1` da TODO OK. Queda el bajo B-N4.
