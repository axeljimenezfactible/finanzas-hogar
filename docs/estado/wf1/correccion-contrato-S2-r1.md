# WF-1 · correccion-contrato-S2-r1 · Reporte final
Agente: correccion-contrato-S2-r1 · 11-oct-2026 · Alcance: `docs/contratos/motores/S2-interfaz.md` (1.0 → 1.1), `oraculo/fixtures/FORMATO.md` (1.0 → 1.1) y este reporte. Corrige los 10 hallazgos medios (M1–M10) de `docs/estado/wf1/REVISION-contrato-S2-r1.md`. Ninguna cifra canon cambia y no se reabrió ninguna decisión de §A–§K.

## Hecho
- **M1** (FORMATO §2.1, §3, §4; S2 §3.6, §3.14, §5): regla única del `id` de una consulta: ruta de `model.json` (10 raíces) o `data.json:<ruta>` → vista con claves literales, y la hoja de `model.json` es `id` + ruta relativa; cualquier otro `id` es nombre libre `^[a-z0-9][a-z0-9_-]*$`. Cruces con vistas propias: `vista_sensibilidad` (ids `sensibilidad[0..4]`, 15 hojas) y `vista_interes_por_monto` (ids `interes_por_monto[0..6]`, 28 hojas). El recorte deja las listas completas (los elementos de más salen `"__sobra__"`), recorta solo las claves de los objetos y marca con `"__tipo__"` un tipo que no casa. Por eso `.#` sí detecta largos distintos.
- **M2** (FORMATO §4.1, §5, §6; S2 encabezado): `backend/cmd/motores/main.go` y `main_test.go` los escribe el agente Go de S2, con la API `Query`, `Evaluator` y `register`; cada motor agrega solo `cmd/motores/<motor>.go`; si `main.go` aún no existe, nadie escribe un sustituto. `oraculo/compara.py` lo escribe el agente Python de S2: es genérico y trae la tabla de ramas de s2, s3, s4 y s5. Cada motor escribe su propio ejecutor del oráculo y no hay módulo Python compartido. S3 adopta FORMATO sin excepción (`s3 dump` va en `cmd/motores/s3.go`); `cred_compara.py` y `--canon` quedan solo como verificación adicional; `$model` no entra en `fixture/1`.
- **M3** (S2 §5.1): esquema `insumo-s3/1`. Las 55 filas S3 van llave por fila con rutas de hoja de `model.json` o `data.json:` y se cubren todos los R-id: R-05 lleva `{num, den, decimales}`, R-65 y R-66 son listas de 14, y R-67 son 85 llaves aplanadas. `TestCanon` aplica tres pasos a cada fila S3.
- **M4** (S2 §3.2, §3.12, §8): fórmulas por mes. Las ocurrencias de gasto no ligadas entran a `accrued` y `cash` como un gasto de contado del titular, con devengo RN-02 y las partes de `share_rule`. `income_cents` suma los registrados y los planificados. `planned_expense_cents` es el monto total (ya está dentro de `cash`). `net_cash_cents = income − cash.by_person[holder]`. `provisional` pasa a ser una lista e incluye `S2-flujo-quincenal` en los meses donde la liga decide.
- **M5** (S2 §1.4, §3.11, §5.2): tabla con la fórmula y la base de cada campo del flujo, incluidos `income.total_cents` (nuevo), `provisional` por quincena y `variable_estimated_cents` y `savings_planned_cents` en null mientras el bloqueo siga abierto. El acumulado arranca en 0 en `from`. En `s2-mes-31-dias`, el parche fija `today` 2026-10-31 y se comparan `variable_spent_cents` y `spent_cents` (50000 en q2 y 0 en q1) y `periodos` 2026-10 (147375).
- **M6** (S2 §3.11, §5, §7): `s2-flujo-quincenal` cubre 16 quincenas (2026-10-q2 a 2027-06-q1) con los pagos de `despensas[0..3]`. Contraste con el valor estable en `TestFlujoQuincenal`, sin campo nuevo: en los ciclos de las despensas 2–4, Σ hogar = 2 × `valor_estable.mes.axel` y Σ Karen = 2 × `karen_mitad_despensa`, ambos con diferencia 0; por quincena se informa la diferencia (+1 ×9, −1 ×3) y no se corrige.
- **M7** (S2 §2.7): `share_rule` llega como la regla efectiva que quien llama ya resolvió (partida → categoría → «Hogar», RN-03). Es obligatoria en los gastos y va en null en los ingresos. Ejemplo: Despensa → «Conjunto 50/50».
- **M8** (S2 §3.2–§3.13): contenido y orden de `periodos`, `payments_by_purchase`, `by_purchase`, `components`, `papel` (desempate de B6), `categories`, `rows`, `deleted_row_ids`, `plan_contra_real`, `goals` y `plan`.
- **M9** (S2 §5, §5.2): fixture `s2-interes-por-monto`, un cruce que toma como entrada la serie que entregaría S3 y confirma `total` por devengado y por efectivo. Con ella S2-03 se cumple tal como está escrito; la solicitud 1 de `contrato-S2.md` deja de hacer falta.
- **M10** (S2 §2.4, §2.6, §2.7): `person` en `Income` (obligatorio) y en `PlannedItem`; en un gasto a crédito, `amount_cents` = `principal_cents`, y si difiere es `sum_mismatch`.
- Fixtures: de 24 a 27 (nuevas `s2-interes-por-monto`, `s2-proyeccion-csv` y `s2-liga-provisional`); la línea de cruces de `compara.py` pasa de 15 a 43.

## Archivos
- `docs/contratos/motores/S2-interfaz.md` (modificado: 44.4 KB → 63.5 KB)
- `oraculo/fixtures/FORMATO.md` (modificado: 7.1 KB → 14.8 KB)
- `docs/estado/wf1/correccion-contrato-S2-r1.md` (este reporte)

## Verificación
El artefacto es un contrato: todavía no existen `backend/internal/motor_fin`, `backend/cmd/motores`, `oraculo/fin_*.py` ni `oraculo/compara.py` (`ls backend/internal` → `api`; `ls oraculo` → `model.py`, `verify3.py`, `model.json`, `verify3.json`, `fixtures`). Por eso no aplican aún las pruebas Go con cobertura, el oráculo ni `compara.py`. Corrí un script propio de validación (`python3 -I datos/tmp/corr-s2/verifica.py`, con `Fraction` y `Decimal`, ya borrado): **31 OK y 0 FALLA**.
- V1 (S2-02): la tabla §4 se expande (con llaves anidadas) a 20 renglones, 159 rutas y 159 distintas, sin hojas sin dueño ni rutas sobrantes. `model.json` tiene 596 hojas numéricas.
- V2 (S2-01, M3): `insumo-s3.json` de §5.1 tiene 55 filas, exactamente las S3 (R-01…R-36, R-61…R-67, R-74…R-85). Todas las llaves se resuelven en `model.json` o `data.json` y su conversión a centavos es exacta. R-04: el binary64 de 100·80050/309450 es igual a `tasa_pct`. R-05: HALF_UP(100·53450/309450) a 2 decimales = 17.27. R-67 tiene 85 llaves, 17 eventos y los tipos compra, inicio, pago y papel; R-65 y R-66 tienen 14 valores cada una. El ejemplo JSON de §5.1 coincide con la conversión.
- V3 (S2-03, M1): `vista_sensibilidad` reproduce 15 de 15 hojas (`costo_mensual_total`, `axel_mensual` y `karen_mensual` de las 5 filas).
- V4 (S2-03, M9): `vista_interes_por_monto` reproduce 28 de 28 hojas. En las 7 series el capital y los intereses por pago son ≥ 0, y Σ devengado = Σ efectivo = total. El monto 3,094.50 reproduce el pago 97375.
- V5 (S2-10, M6): 16 quincenas consecutivas de 2026-10-q2 a 2027-06-q1, una por pago. En los ciclos de las despensas 2–4, el hogar suma 353502 = 2 × 176751 y Karen 149972 = 2 × 74986. Valor estable − pago: +1 en 9 quincenas y −1 en 3. Despensa 1: 59882 ×4; Karen: 37493 en las 16.
- V6 (S2-06, M5): 2026-10-31 cae en 2026-10-q2. Efectivo de octubre: 97375 + 50000 = 147375 (octubre no tiene papel de contado).
- V7: 27 fixtures en §5; FORMATO dice «27 fixtures» y «43 de 43»; no quedan restos («8 quincenas», «24 fixtures», «15 de 15», «null = todo al titular», «despensas[0..1].pagos», «texto libre o la ruta»).
- V8: existen las referencias a FORMATO §2.1, §3, §4, §4.1 y §5.
- V9: todas las tablas de los dos archivos tienen el mismo número de columnas en cada renglón.
- V10 (S2-04): existen los 4 bloqueos de S2 citados. `sha256sum oraculo/model.json` da `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`, el mismo valor antes y después.
- Criterios cubiertos por el contrato tras la corrección: S2-01 (§5.1, §7), S2-02 (§4, FORMATO §2.1 y §5), S2-03 (§3.14, §5, FORMATO §2.1, §4 y §4.1; los cruces incluyen interés por monto), S2-06 (§3.11, §5.2), S2-09 (§3.12, §5.2), S2-10 (§3.11, §5) y S2-16 (§8). Los demás quedan como en la 1.0.

## Supuestos
Ninguno cambia cifras canon; todos definen salidas que no son canon a partir de documentos existentes.
- El efectivo del hogar en el flujo es lo que sale del bolsillo del titular: el gasto de contado completo cuando `paid_by` es null o el titular, y su parte en los pagos programados. Es la misma base que el efectivo por persona de §3.2 y el «contado + partes del hogar pagadas» de M-02.
- Una ocurrencia planificada de gasto sale del bolsillo del titular, como un `Expense` con `paid_by` null, porque `planned_item` no tiene la columna «quién paga».
- `cumulative_balance_cents` arranca en 0 en `from`: el modelo no tiene saldo inicial de cuentas, y M-10 empieza en la Q1 del mes en curso.
- `variable_estimated_cents` y `savings_planned_cents` van en null, no en 0, mientras `S2-flujo-quincenal` siga abierto, y no entran a ningún saldo. Así nada provisional entra en cifras.
- En un gasto a crédito, `amount_cents` = `principal_cents`, según el ejemplo de RN-03 (importe del gasto 3,094.50 = capital, R-01).
- La serie de `s2-interes-por-monto` se arma con la regla de S3 (RN-07: el último pago ajusta; K14: capital con base HALF_UP). Es un dato de entrada de la fixture; S2 no la calcula.
- `presupuesto_mes.categories` y `metas.goals` solo listan elementos activos (M-01; las partidas inactivas ya se ignoraban en §2.7).

## Bloqueos
No abrí bloqueos nuevos. Siguen vigentes `docs/estado/bloqueos/S2-proyeccion-interruptor.md`, `S2-flujo-quincenal.md` (ahora también cubre la marca `provisional` de §3.11 y §3.12 y los dos campos en null), `S2-metas.md` y `S2-presupuesto-total.md`. `S1-parametros-avisos` no afecta a S2.

## Pendientes
- Hallazgos bajos sin tocar (fuera de este encargo): B1 (`plan_contra_real` y `promedios {month, basis}` sin fixture), B2 (módulos Python sin nombre en §3.3–§3.7, §3.9 y §3.14), B3 (redacción de S3 §7), B4 (sigue sin nombrarse el mes de `s2-mes-30-dias` y falta el error de un evento fuera de las ventanas; la función desconocida ya quedó en §3.15) y B5 (aserción de 0 eventos entre los días 1 y 9 en `TestProyeccionDefecto`). B6 quedó resuelto con M8 (desempate de `papel`).
- Implementación (otros agentes): `backend/cmd/motores/main.go` y `s2.go`, `backend/internal/motor_fin/*.go`, `oraculo/fin_*.py`, `oraculo/fin_run_fixtures.py`, `oraculo/compara.py`, las 27 fixtures de `oraculo/fixtures/s2/` y `_datos/semilla-oraculo.json`, `_datos/insumo-s3.json` e `_datos/interes-por-monto.json`.
- Revisión r2 del contrato S2.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S3-interfaz.md` §2 y §9–§11. Cambio: alinear a FORMATO 1.1 §6. Fixtures en `oraculo/fixtures/s3/` con `fixture/1`; `s3 dump` en `backend/cmd/motores/s3.go` con `register` (no en `backend/cmd/finanzas`); salida `salida/1`; `esperado` con valores escritos en lugar de `$model`; R-05 y R-06 con `id` `data.json:financiamiento.dos_quincenas` y razón con `decimales`; retirar el SUPUESTO de §10. Motivo: M2 de esta revisión y H1 de `REVISION-contrato-S3-r1.md`.
2. `docs/contratos/motores/S4-interfaz.md` §12 y `S5-interfaz.md` §15. Cambio: citar FORMATO 1.1: `id` de §2.1 (nombres libres sin puntos, corchetes ni dos puntos), recorte con listas completas y `cmd/motores/<motor>.go` con `register`. Motivo: la 1.1 cambia reglas que esos contratos heredan.
3. `docs/criterios/S2.md` S2-03 (solicitud 1 de `contrato-S2.md`). Cambio: ya no hace falta, porque la interfaz cubre interés por monto como cruce (`s2-interes-por-monto`, 28 hojas que siguen siendo de S3); el integrador puede cerrarla sin tocar el criterio. Motivo: M9.
4. `docs/estado/bloqueos/S2-flujo-quincenal.md`, «Cómo seguí». Cambio: anotar que `variable_estimated_cents` y `savings_planned_cents` salen en null (no en 0) y que la marca `provisional` se pone donde la liga decide (S2 §3.11 y §3.12, fixture `s2-liga-provisional`). Motivo: M4 y M5.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: no hay valores operativos de K31 en este alcance. Las propuestas sobre los 4 bloqueos de S2 siguen siendo las de `contrato-S2.md`.

## Consumo aproximado
38 acciones de herramienta; ~215 KB leídos (protocolo 6, decisiones §B/§C/§K 17, criterios S2 34, revisión 17, interfaz 44, FORMATO 7, reportes y bloqueos 12, extractos de esquema, reglas y S3 ~25, salidas de verificación ~10); ~75 KB escritos (scripts de edición de FORMATO ~11 y de la interfaz ~38, verificación ~10, progreso y reporte ~12).
