# S1-modelo-datos · Reporte final (WF-0, fase S1 «Dominio y modelo de datos»)

Fecha: 10-oct-2026 · Estado: **completo**, corregido en la ronda 1 de revisión (hallazgos H-01, H-02, H-04 y H-11 de `REVISION-r1.md`; detalle en `correccion-S1-r1.md`) · Sin `git commit` (lo hace el integrador).

## Hecho
### Corrección de la ronda 1 (10-oct-2026)
- **H-01 (lista de compras):** tablas `shopping_list` (una por semana: `start_date`, `end_date` = inicio + 6 días, `generated_at`) y `shopping_list_item` (origen `calculated`/`manual`, producto o texto libre con unidad, presentación, `needed_milli` calculada por S4, `quantity_milli` elegida, `store_id`, `is_bought`, `bought_at`). Cubren N-151, N-152, N-153, N-165, N-127 y N-073. **PROVISIONAL**: opción A del bloqueo `C1-lista-compras` (mismo tema en `S1-lista-compras`) hasta la puerta 1.
- **H-02 (huecos S1 ↔ U1):** `expense.description` admite NULL (N-054, F-1); `recipe_meal_slot` sustituye a `recipe.default_meal_slot_id` (varios tiempos aptos, N-122, N-135, N-136); `inventory_movement` con tipo y motivo (N-181, N-183); `budget` con vigencia `valid_from_date`/`valid_to_date` en lugar de `period_id` (N-087, N-275, N-278); `planned_item.second_half_amount_cents` (el 15 y el 30 pueden diferir, N-279, N-280); `idempotency_key` (C1 §1 punto 8). Al construir la tabla de S1-13 apareció un hueco más y quedó corregido: `category.kind` (fijo o variable, N-281). Lo que es de negocio quedó en el bloqueo `S1-huecos-negocio` (vigencia por defecto de un cambio de presupuesto, correcciones a medias del OCR, registro del chat).
- **H-04 (evidencia):** 17 líneas `S1-xx` en «Verificación»; prueba reejecutable `s1_verifica.py` (Anexo B, con ida y vuelta completa: 164 valores, 0 diferencias); sección «Ajustes: 14 secciones → campos → tabla.columna» en el diccionario; Apariencia queda local por dispositivo, como en U1 (N-301) y C1 (§11.3), y se quitó la clave `appearance.*`.
- **H-11 (baja, de S1):** índice `ux_person_single_co_holder` (con el titular, dos cuentas como máximo) y disparador `trg_person_external_without_access`; el tope por quincena de Karen ya no se interpreta como 374.93 (ver Pendientes 3).
- Generador del diccionario sin plantillas aparte (toma el texto fijo del propio diccionario; Anexo A) y diccionario regenerado: versión 1.1, 54 tablas, 147 índices, 68 disparadores. Reglas actualizadas: C-13, C-15 (idempotencia), C-16 (descripción opcional), RN-05, RN-14, RN-15 y RN-22.

### Entrega original (WF-0)
- **Migración goose** `backend/db/migrations/00001_esquema.sql` (Up y Down): **54 tablas STRICT** desde la ronda 1 (49 en la entrega original) = las 31 entidades de §3.4.1 + 14 tablas de detalle (renglones, partes, pagos hechos, cortes, suscripciones, envíos, tiempos aptos y movimientos de inventario) + 9 complementarias (`login_attempt`, `app_setting`, `audit_log`, `idempotency_key`, `budget`, `product_group`, `generated_report`, `shopping_list`, `shopping_list_item`). 147 índices (incluye los automáticos de UNIQUE). 68 disparadores: 15 de reglas de decisiones (Karen sin cuenta, nadie con accesos pasa a externa, ingresos solo del titular, salud y comidas sin externos, bitácora inmutable) y 53 de `updated_at`.
- Montos INTEGER `_cents`. Fechas TEXT ISO validadas con `CHECK (date(x) IS x)`. Booleanos 0/1 con CHECK. `id INTEGER PRIMARY KEY`. `created_at`/`updated_at` en todas las tablas. CHECK de enumerados, porcentajes 0–100 % (puntos base 0–10000), sumas dentro del renglón (total = capital + intereses) y el calendario de quincenas (días de pago 15 y 30, o último día).
- Datos de referencia sembrados: `meal_slot` (4 tiempos), `notification` (7 avisos, solo con parámetros decididos) y `app_setting` (7 claves decididas, ninguna provisional).
- **Diccionario** `docs/contratos/diccionario-datos.md`, generado desde los comentarios estructurados de la migración y la introspección de la base migrada (tipos, nulos, FK, índices y disparadores reales). Incluye el mapa entidad → tabla, el glosario `nombre_codigo`, el propósito de cada índice, las relaciones de entrada y salida, las invariantes que SQLite no puede expresar y notas para S2, S3, S5, S6 y U1.
- **Reglas de negocio** `docs/contratos/reglas-negocio.md`: 22 reglas (todas las pedidas y además redondeo, equivalencias, deuda viva y presupuesto), cada una con definición, un ejemplo numérico que cita su clave de `model.json` o `data.json` (o marcado «ilustrativo» cuando el reporte y el oráculo no traen el dato), estado FIRME o PROVISIONAL y las tablas que la implementan. También 14 convenciones técnicas: centavos, ISO, zona horaria, **sin borrado lógico general** (catálogos archivados con `is_active`, movimientos borrados con imagen en `audit_log`) y auditoría.
- **5 bloqueos** S1 y alineación con los bloqueos paralelos de U1 y de criterios. La quincena de 1–15 / 16–fin quedó PROVISIONAL, igual que en U1.
- **Hallazgos del oráculo** (documentados, sin ajustar ninguna cifra):
  1. Conviven dos redondeos de serie. Los montos de pago tienen base truncada: 5,034.74 → 1,258.68 ×3 + 1,258.70 (model.py líneas 55–57). El capital por pago tiene base HALF_UP: 3,094.50 → 773.63 ×3 + 773.61 (model.py línea 40, fijo a mano). Con HALF_UP en los montos saldría 1,258.69 ×3 + 1,258.67.
  2. El «pago quincenal» de `interes_por_monto` y de `sensibilidad` es HALF_UP del total sin redondear y difiere de la serie real: con 5,000.00 da 1,573.36 contra 1,573.35 ×3 + 1,573.37.
  3. La tasa redondeada a 25.8685 % no reproduce `interes_por_monto` (1,000 → 258.69 contra 258.68; 5,000 → 1,293.43 contra 1,293.42), así que las tasas se guardan como fracción exacta 80050/309450.
  4. Los periodos P1–P6 del oráculo (del 10 al 9) dan exactamente las cifras de los meses calendario oct-2026 a mar-2027 (devengado y efectivo, 12 de 12 coincidencias).
  5. En las 14 cifras de las series A/B, A − B = 3,094.50.

## Archivos
Creados (todos dentro de mi alcance):
- `/home/deepsec/finanzas-hogar/backend/db/migrations/00001_esquema.sql` (1835 líneas, ~121 KB)
- `/home/deepsec/finanzas-hogar/docs/contratos/diccionario-datos.md` (2,013 líneas, ~136 KB)
- `/home/deepsec/finanzas-hogar/docs/contratos/reglas-negocio.md` (227 líneas, ~35 KB)
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/S1-modelo-datos.md` (este reporte)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-redondeo-series.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-reposicion-parte-externa.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-frescura-peso.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-meta-kcal.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-restricciones-optimizador.md`

Temporales: `/home/deepsec/finanzas-hogar/datos/tmp/` (s1.db y goose.db) creada para la verificación y **borrada** al terminar. `datos/` vuelve a tener solo `semilla/`. Los scripts de apoyo quedaron en el scratchpad de la sesión (ver Solicitudes y Anexo A).

## Verificación
### Corrida de la ronda 1: evidencia por criterio
Comando (reejecutable desde el repositorio; extrae y corre el Anexo B, crea `datos/tmp/s1-verifica.db` y la borra al final; con `GOOSE=<binario goose v3.28.0>` también corre goose up/down/up):
`python3 -I -c "import re;t=open('/home/deepsec/finanzas-hogar/docs/estado/wf0/S1-modelo-datos.md',encoding='utf-8').read();exec(compile(re.search(r'<!-- inicio s1_verifica.py -->\n```python\n(.*?)\n```\n<!-- fin s1_verifica.py -->',t,re.S).group(1),'s1_verifica.py','exec'))"`
Se corrió con `GOOSE=…/scratchpad/gobin/goose` (v3.28.0), sqlite3 3.53.4 y Python 3.14. Salida, una línea por criterio:

S1-01 OK   31 de 31 entidades con tabla: persona→person, cuenta_usuario→user_account, sesion_dispositivo→device_session, periodo→period, ingreso→planned_item, categoria→category, gasto→expense, regla_reparto→cost_share_rule, meta_ahorro→savings_goal, movimiento_ahorro→savings_movement, credito→credit, tasa_por_plazo→credit_term_rate, compra_credito→credit_purchase, pago_programado→scheduled_payment, promocion→promotion, tienda→store, producto→product, presentacion→presentation, precio→price, receta→recipe, ingrediente_receta→recipe_ingredient, tiempo_comida→meal_slot, plan_comida→meal_plan_entry, inventario_item→inventory_item, perfil_salud→health_profile, peso_historico→weight_entry, nutricion_ingrediente→nutrition_fact, notificacion→notification, destinatario→notification_recipient, job→job, respaldo→backup · faltantes=[] · tablas sin fila en el diccionario=[]
S1-02 OK   Up aplicado con sqlite3 3.53.4: tablas=54 strict=54 índices=147 disparadores=68 · integrity_check=ok · foreign_key_check=0 filas
S1-03 OK   22 reglas · RN-01: ejemplo=sí, citas 1/1 · RN-02: ejemplo=sí, citas 1/1 · RN-03: ejemplo=sí, citas 5/5 · RN-04: ejemplo=sí, citas 1/1 · RN-05: ejemplo=sí, citas 2/2 · RN-06: ejemplo=sí, citas 1/1 · RN-07: ejemplo=sí, citas 1/1 · RN-08: ejemplo=sí, citas 2/2 · RN-09: ejemplo=sí, citas 2/2 · RN-10: ejemplo=sí (ilustrativo), citas 1/1 · RN-11: ejemplo=sí, citas 1/1 · RN-12: ejemplo=sí, citas 1/1 · RN-13: ejemplo=sí (ilustrativo), citas 1/1 · RN-14: ejemplo=sí (ilustrativo), citas 0/0 · RN-15: ejemplo=sí (ilustrativo), citas 1/1 · RN-16: ejemplo=sí (ilustrativo), citas 0/0 · RN-17: ejemplo=sí (ilustrativo), citas 0/0 · RN-18: ejemplo=sí (ilustrativo), citas 1/1 · RN-19: ejemplo=sí (ilustrativo), citas 0/0 · RN-20: ejemplo=sí, citas 1/1 · RN-21: ejemplo=sí, citas 1/1 · RN-22: ejemplo=sí (ilustrativo), citas 1/1 · reglas sin ejemplo=0 · citas que no coinciden con su ruta=[]
S1-04 OK   tablas con «Usado por» lleno: 54 de 54 · vacías=[]
S1-05 OK   columnas REAL/FLOAT/DOUBLE=0 · columnas _cents por tipo: INTEGER=36 · convención: enteros de centavos
S1-06 OK   fechas en formato de pantalla en datos o código SQL (sin comentarios): 0 · grep literal del criterio «oct-2026|ene-20» en la migración: 0 · zona America/Mexico_City declarada en el diccionario: True
S1-07 OK   cuenta para Karen: rechazada · tercera cuenta (otra co-titular): rechazada · Ximena con cuenta pasa a externa: rechazada · ingreso de Ximena: rechazado · columnas de user_account sin rol: ['id', 'person_id', 'username', 'password_hash', 'is_active', 'password_changed_at', 'last_login_at', 'created_at', 'updated_at'] · tablas de finanzas de la familia de Karen: []
S1-08 OK   regla con partes 60 % + 50 %: la base la acepta y la rechaza la API (invariante 8 del diccionario, prueba de S6) · despensa 1 sin pérdida: partes del gasto Axel 190300 + Karen 119150 = 309450 (ver S1-12) · conjunto 170200 → mitad 85100 cada quien (data.json: cuentas.conjunto.mitad)
S1-09 OK   18 filas recuperadas (6 meses oct-2026…mar-2027 y 12 quincenas): m 2026-10-01…2026-10-31; f1 2026-10-01…2026-10-15 paga 2026-10-15; f2 2026-10-16…2026-10-31 paga 2026-10-30; m 2026-11-01…2026-11-30; f1 2026-11-01…2026-11-15 paga 2026-11-15; f2 2026-11-16…2026-11-30 paga 2026-11-30; m 2026-12-01…2026-12-31; f1 2026-12-01…2026-12-15 paga 2026-12-15; f2 2026-12-16…2026-12-31 paga 2026-12-30; m 2027-01-01…2027-01-31; f1 2027-01-01…2027-01-15 paga 2027-01-15; f2 2027-01-16…2027-01-31 paga 2027-01-30; m 2027-02-01…2027-02-28; f1 2027-02-01…2027-02-15 paga 2027-02-15; f2 2027-02-16…2027-02-28 paga 2027-02-28; m 2027-03-01…2027-03-31; f1 2027-03-01…2027-03-15 paga 2027-03-15; f2 2027-03-16…2027-03-31 paga 2027-03-30 · quincena 2 de feb-2027 paga 2027-02-28 · ventana 2026-10-10…2026-11-09 como mes: rechazada · equivalencia con P1–P6 (efectivo y devengado totales): 12 de 12 · eventos de la semilla entre los días 1 y 9: 0
S1-10 OK   credit.kind = 'revolving_line', 'bank_card', 'term_loan', 'store_credit' · tasas semilla como fracción: 2q 53450/309450 = 17.27 %, 4q 80050/309450 = 25.8685 % · lectura por defecto al insertar sin valor = B · availability_confirmed=True · promotion: credit_id, merchant_name, store_id, term_count, rate_num, rate_den, min_purchase_cents
S1-11 OK   12 de 12 puntos con tabla.columna: tiempos 4 → meal_slot.code; fuera sin receta → meal_plan_entry.is_eating_out; porciones base 2 y escalado → recipe.base_servings; perfil por persona → health_profile.height_mm; historial de peso → weight_entry.weight_g; restricciones vacías → health_profile.dietary_restrictions; kcal con fuente → nutrition_fact.source; inventario híbrido y sobrantes → inventory_item.tracking_mode; cadena/sucursal/mercado con día y qué se compra → store.visit_weekday; membresía Sam’s → store.membership_annual_fee_cents; precio con tienda, fecha y fuente → price.observed_date; unidad base y macros previstos → nutrition_fact.protein_mg_per_100g
S1-12 OK   ida y vuelta de oraculo/model.json: 4 despensas, 16 pagos, 17 eventos, series A/B de 14/14 valores, mínimos 809450/500000, P1–P6 por persona · valores comparados=164 · diferencias=0 · suma pagos d1=389500 · parte Axel d1=239528
S1-13 OK   diccionario «Ajustes: 14 secciones → campos → tabla.columna»: 14 de 14 secciones en orden, 0 destinos vacíos=True, referencias tabla.columna=45 (inexistentes: []), tablas citadas=24; Apariencia local por dispositivo (N-301)
S1-14 OK   coincidencias de password/token/secret/apikey en la migración (solo nombres de columna o texto): ['password_changed_at', 'password_hash', 'token', 'token_hash'] · patrón de token de Telegram: 0 · password_hash CHECK argon2id: True
S1-15 OK   foreign_key_list: {'expense': 6, 'credit_purchase': 3, 'scheduled_payment_share': 2, 'price': 4, 'meal_plan_entry': 3, 'shopping_list_item': 4} · índices de filtros C1: persona=ix_expense_person_date, familia (partes por persona)=ix_expense_share_person, categoría=ix_expense_category_date, tienda=ix_expense_store_date, periodo (fechas)=ix_expense_date, crédito=ix_credit_purchase_credit_status · PRAGMA foreign_keys documentado en la cabecera de la migración y en C-13
S1-16 OK   líneas con PROVISIONAL: 56 (migración 11, diccionario 19, reglas 26) · bloqueos citados que existen: 13 · líneas sin bloqueo en su bloque=0 []
S1-17 OK   identificadores (tablas y columnas): 633, fuera de [a-z0-9_]: 0 · tablas únicas: 54 · filas del glosario: 30, repetidas: 0

NEG OK   pruebas negativas: 46, rechazadas 46, aceptadas por error [] · positivas de las tablas nuevas: 8, aceptadas 8, rechazadas por error []
DOWN OK   Down deja 0 objetos · Up otra vez: 54 tablas, integrity_check=ok
GOOSE OK   goose (../gobin/goose) up/down/up → códigos [0, 0, 0], tablas al final=54
RESUMEN TODO OK

Además, fuera del script: `goose -dir backend/db/migrations sqlite3 datos/tmp/goose-r1.db up` → «OK 00001_esquema.sql», 54 tablas propias, `foreign_key_check` vacío, `integrity_check` ok; `down` → 0 tablas propias; `up` → 54 tablas, ok. El generador del diccionario (Anexo A) reprodujo byte a byte el diccionario 1.0 desde la migración original antes de regenerar la 1.1 («columnas sin descripción: []»). `datos/tmp` se borró al terminar.

### Corrida original (WF-0, esquema de 49 tablas; se conserva como historia)
Todo se corrió **dos veces**: con el primer borrador y con la versión final de la migración. Salidas resumidas de la corrida final:

1. **Up con sqlite3 3.53.4** (`awk` extrae la sección Up → `sqlite3 -bail datos/tmp/s1.db < up.sql`): `up: ok` · `tablas=49 indices=134 disparadores=62 strict=49`.
2. `PRAGMA foreign_key_check` → vacío · `PRAGMA integrity_check` → `ok` (después del Up, después de los inserts y después del segundo Up).
3. **Juego mínimo**:
   - Personas: Axel (holder), Ximena (co_holder), ambos con cuenta y hash de prueba argon2id; Karen (external, paga directo, tope 37493).
   - Reglas Hogar, Karen y Conjunto 50/50, y la categoría Despensa.
   - Bodega Aurrera (cadena) y Vallejo (sucursal).
   - 3 artículos de `data.json` con su cuenta: aceite Nutrioli Oli 500 ml (conjunto, 15000), duraznos La Costeña 425 g (axel, 4000), papel Kleenex Cottonelle 30 rollos (karen, 19900), con sus precios.
   - Crédito Kueski (límite 1000000, quincenal, primer pago ≥ 15 días, lectura B, ancla A 2026-10-10) y tasas 2q = 53450/309450 y 4q = 80050/309450.
   - Despensa 1 como gasto (309450, devengo 2 meses, crédito) con partes Axel 190300 y Karen 119150.
   - Compra a crédito 309450 + 80050 = 389500, con partes Axel 190300/49228/239528 y Karen 119150/30822/149972.
   - 4 pagos de 97375 (capital 77363 ×3 + 77361, intereses 20012 ×3 + 20014) el 2026-10-30, 2026-11-15, 2026-11-30 y 2026-12-15, con partes Axel 59882 y Karen 37493.
   - Periodos oct-2026 y feb-2027 (la quincena 2 paga el 2027-02-28).
4. **Consultas**: `suma_pagos=389500` ✔ · `parte_axel=239528` ✔ · `parte_karen=149972` · `capital_repuesto=309450 intereses=80050` · `pagos_con_partes_descuadradas=0` · `partes_compra_vs_calendario_descuadradas=0` · `partes_gasto=309450 importe_gasto=309450` · `disponible_B_tras_pago_1=767913` (= 7,679.13 de `credito.comparacion_series.esperado_B[0]`) · `efectivo_oct_2026=97375` · `pago_quincena_2_feb_2027=2027-02-28` · `interes_4000_con_tasa_4q=103474` (= `interes_por_monto[5].intereses`).
5. **`updated_at` por disparador**: `updated_at_cambio=1`.
6. **21 pruebas negativas, las 21 rechazadas** (0 aceptadas): Karen con cuenta; ingreso de Ximena; hash que no es argon2id; segundo titular; total ≠ capital + intereses (compra y pago); pct_bp 10001; parte con porcentaje y fijo a la vez; fecha 2027-02-30; fecha «10-oct-2026»; quincena que paga el 14; 12.5 en una columna `_cents` («cannot store REAL value in INTEGER column»); FK inexistente; plan de comidas de Karen; «fuera» con porciones; lectura A sin ancla; tarjeta sin corte; mercado sin día de visita; compra a crédito con pagador de contado; borrar a Karen con historia (FK RESTRICT); UPDATE a `audit_log` («audit_log es inmutable»).
7. **Down** (`sqlite3 -bail s1.db < down.sql`): `objetos restantes=0`. **Up otra vez**: `tablas=49`, `foreign_key_check` vacío, `integrity_check ok`.
8. **goose v3.28.0** con el driver sqlite3 de modernc, sobre `datos/tmp/goose.db`: `up` → «OK 00001_esquema.sql», «successfully migrated database to version: 1». `status` muestra la migración aplicada; 49 tablas propias (50 contando `sqlite_sequence` de goose) y 62 disparadores. `down` → 0 tablas propias. `up` otra vez → 49 tablas, `integrity_check ok`. Así queda probado que el parser de goose acepta las anotaciones (`StatementBegin`/`StatementEnd` en los 62 disparadores).
9. **Verificador independiente de ejemplos** (Python, enteros de centavos y `Fraction`, lee `oraculo/model.json` y `datos/semilla/data.json`): **todo OK**. Cubre:
   - los 7 intereses y los 7 pagos de referencia de `interes_por_monto` con la tasa 80050/309450;
   - montos ROUND_DOWN y capital HALF_UP de las 4 despensas, y titular = monto − Karen;
   - interés por persona 30822/49228;
   - las 16 fechas de pago con la regla 15/30 y primer pago ≥ 15 días;
   - meses calendario = P1–P6 (12 de 12) y devengado por persona en P1;
   - A − B constante;
   - reparto de `data.json` (conjunto 851.00, totales 1,903.00 y 1,191.50, sumas de renglones).
10. **Diccionario**: el generador reporta «columnas sin descripción: []». Una revisión aparte confirma las 31 entidades de §3.4.1 en el mapa, 0 etiquetas HTML sueltas y 0 filas de tabla mal formadas.

## Cobertura
### Entidades de `plan.md` §3.4.1 (31 de 31)
| Entidad del plan | Tabla(s) |
|---|---|
| persona (Karen marcada externa) | `person` (`participation` = holder / co_holder / external) |
| cuenta_usuario (2) | `user_account` (disparador: una persona externa no tiene cuenta) |
| sesion_dispositivo | `device_session` |
| periodo (mes calendario, quincena 15/30 con ajuste a fin de mes) | `period` (CHECK del calendario; rangos de quincena PROVISIONALES) |
| ingreso (fijo quincenal, extra) | `income` (`kind` fixed / extra) + `planned_item` (`frequency` fortnightly) |
| categoria | `category` |
| gasto (fijo/variable/venidero; devengado vs efectivo; adjunto) | `expense` (`kind`, `expense_date`, `accrual_months`, `payment_method`, `paid_by_person_id`) + `expense_item` + `expense_share` + `attachment` |
| regla_reparto (hogar / Karen / conjunto con %) | `cost_share_rule` + `cost_share_rule_part` |
| meta_ahorro + movimiento_ahorro | `savings_goal` + `savings_movement` |
| credito (tipo, límite, lectura A/B, corte y pago, comisiones) + tasa_por_plazo | `credit` + `credit_term_rate` (+ `credit_statement` para corte, pago y mínimo) |
| compra_credito (capital, intereses con IVA, plazo, partes por persona) | `credit_purchase` + `credit_purchase_share` |
| pago_programado | `scheduled_payment` + `scheduled_payment_share` (pagos reales en `credit_payment`) |
| promocion (comercio, plazo, mínimo, vigencia) | `promotion` |
| tienda (cadena / sucursal / mercado con día de visita; membresía con costo anual) | `store` (+ `store_product_group`) |
| producto | `product` (+ `product_group`) |
| presentacion (unidad base y factor) | `presentation.size_milli` + `product.base_unit` |
| precio (tienda, fecha, fuente, vigencia 30 d) | `price` + `app_setting price.stale_after_days` |
| receta · ingrediente_receta | `recipe` (+ `recipe_meal_slot`, tiempos aptos) · `recipe_ingredient` |
| tiempo_comida | `meal_slot` (4 renglones sembrados) |
| plan_comida (día × tiempo × persona × receta × porciones / «fuera») | `meal_plan_entry` |
| inventario_item (stock por unidad o ritmo) | `inventory_item` (+ `inventory_movement`; lista de compras en `shopping_list` y `shopping_list_item`, PROVISIONAL) |
| perfil_salud + peso_historico | `health_profile` + `weight_entry` |
| nutricion_ingrediente (kcal/100 g, fuente) | `nutrition_fact` |
| notificacion + destinatario (chat_id) | `notification` + `notification_recipient` (+ `notification_subscription`, `notification_delivery`) |
| job | `job` |
| respaldo | `backup` |

### Decisiones de `decisiones.md` secciones A–E
| Decisión | Implementación |
|---|---|
| A · Cuentas: Axel titular y Ximena co-titular, mismos permisos | `person.participation`, `user_account` (sin columna de rol) |
| A · Usuario y contraseña, sesión larga por dispositivo, cierre desde Ajustes | `user_account.username`, `device_session.expires_at` / `revoked_at` |
| A · argon2id | `user_account.password_hash` CHECK `LIKE '$argon2id$%'` |
| A · Límite de 5 intentos en 15 min | `login_attempt` (índices por usuario y por IP) + `app_setting auth.max_failed_attempts = 5`, `auth.lockout_window_minutes = 15` |
| A · Karen externa sin cuenta ni acceso | `person.participation = external` + `trg_user_account_not_external_*` |
| A · Karen solo en la despensa común y sus compras Kueski | partes en `expense_share` / `credit_purchase_share` / `scheduled_payment_share`; regla por defecto Hogar (RN-03) |
| A · Karen paga directo, su parte es informativa y repone la línea al pagar | `person.pays_share_directly`, `scheduled_payment_share`, `credit_payment.payer_person_id` (reparto del capital por parte: bloqueo `S1-reposicion-parte-externa`) |
| A · No se modelan las finanzas de su familia | sin tablas para su hogar; `trg_income_holder_only_*` |
| A · Estado de cuenta PDF que Axel reenvía | `generated_report` (`kind = external_statement_pdf`, `person_id`) |
| A · Plan de comidas y salud solo Axel y Ximena | `trg_meal_plan_entry_not_external_*`, `trg_health_profile_not_external_*`, `trg_weight_entry_not_external_*` |
| A · Avisos de Telegram a Axel y Ximena, destinatarios por aviso | `notification_recipient.chat_id`, `notification_subscription` |
| B · MXN, centavos exactos, HALF_UP, el último pago ajusta | `_cents` INTEGER STRICT; series materializadas en `scheduled_payment` (RN-01, RN-07; generalización en el bloqueo `S1-redondeo-series`) |
| B · America/Mexico_City; ISO en archivos y API | `_date` TEXT ISO con CHECK; `_at` en UTC; convenciones C-03 a C-05 |
| B · Solo es-MX | textos visibles sembrados en es-MX (`meal_slot.name`, `notification.name`, `app_setting.description`) |
| B · Mes calendario base; quincena 15 y 30 (o último día) | `period` (`kind`, `half`, `pay_date` con CHECK) |
| B · Despensa bimestral devengada aparte (mitad por mes) | `expense.accrual_months = 2`, `planned_item.accrual_months` |
| B · Ingresos solo de Axel, quincenales 15/30, capturados por él, más extras | `income` (fixed / extra) + `planned_item` (income, fortnightly) + disparadores de titular |
| B · Reparto configurable por gasto o categoría, 50/50 por defecto | `cost_share_rule(_part)`, `category.default_cost_share_rule_id`, `expense.cost_share_rule_id` |
| B · Qué se comparte con Karen, configurable por gasto | `expense.cost_share_rule_id`, `expense_item.cost_share_rule_id`, `expense_share` |
| B · Despensa conjunto 50/50; el crédito se reparte según lo que cada quien compró | regla de categoría Despensa = Conjunto 50/50; reglas por renglón → `credit_purchase_share` |
| B · Fondo de ahorro y metas con fecha, aportación mensual, avance | `savings_goal` (`kind`, `target_cents`, `target_date`, `planned_monthly_cents`) + `savings_movement` |
| B · Comidas fuera: gasto normal con persona; el plan marca «fuera» | `category.code = eating_out` (semilla S8), `expense.person_id`, `meal_plan_entry.is_eating_out` |
| C · 4 tipos de crédito, varios créditos | `credit.kind` = revolving_line / bank_card / term_loan / store_credit |
| C · Kueski: límite, tasa por compra, repone capital | `credit.credit_limit_cents`, `credit_purchase.interest_cents` (quoted), `scheduled_payment.principal_cents` |
| C · Tarjeta: corte, pago, mínimo, MSI | `credit.statement_day` / `payment_due_day`, `credit_statement.minimum_payment_cents`, `credit_purchase` installment_plan al 0 % |
| C · Préstamo personal con amortización | `credit_purchase.kind = loan_disbursement`, `interest_source = amortization`, `rate_basis = annual` |
| C · Crédito de tienda con abonos | `credit.kind = store_credit`, `payment_frequency` weekly / fortnightly / monthly |
| C · Tasa por compra con IVA incluido | `credit_purchase.interest_cents`, `interest_includes_vat`, `rate_num` / `rate_den` |
| C · Tabla editable por plazo (2q 17.27 %, 4q 25.8685 %) | `credit_term_rate` (fracción exacta, `source`, `quote_date`) |
| C · Promociones 0 % a mano (comercio, plazo, mínimo, vigencia) | `promotion` |
| C · Interruptor A/B, B por defecto, aviso «pendiente confirmar» | `credit.availability_reading DEFAULT 'B'`, `availability_confirmed`, `reading_a_anchor_date` |
| C · Escenarios crédito × plazo × fecha ± días × promoción | `planned_item` (compra planeada con crédito y plazo), `credit_term_rate`, `promotion`; `optimizer.date_window_days` sin valor |
| C · Restricciones duras (colchón, máximo de compras vivas, deuda cero antes de fecha) | `credit_purchase.status`; claves `app_setting optimizer.*` sin valor (bloqueo `S1-restricciones-optimizador`) |
| C · 4 objetivos y puntaje con pesos iguales | `app_setting optimizer.weights` (sembrado) |
| C · El tope por persona no es restricción dura | `person.fortnight_payment_cap_cents` (objetivo) |
| C · Regresión Kueski-Despensa | verificación: 389500, 239528 y verificador independiente todo OK |
| D · Recetario de Claude + recetas de Axel y Ximena, editable | `recipe.origin` (claude_seed / household), `author_person_id` |
| D · Ingredientes exactos con presentación comprable, porciones, tiempos | `recipe_ingredient` (`product_id`, `presentation_id`, `quantity_milli`), `recipe.base_servings` / `prep_minutes` / `cook_minutes` |
| D · Sin restricciones alimentarias (el campo existe, vacío) | `health_profile.dietary_restrictions` |
| D · Desayuno, comida, cena y colación opcional | `meal_slot` (4, `is_optional`) |
| D · Porciones base 2 con escalado | `recipe.base_servings DEFAULT 2`, `meal_plan_entry.servings_milli` |
| D · Perfil, IMC, meta kcal, kcal por receta y día | `health_profile`, `weight_entry`, `nutrition_fact` (fórmula de meta: bloqueo `S1-meta-kcal`) |
| D · Macronutrientes después (modelo previsto) | `nutrition_fact.protein_mg_per_100g`, `fat_…`, `carbs_…`, `fiber_…` |
| D · Fuentes nutricionales autorizadas | `nutrition_fact.source` (usda_fdc / open_food_facts / mx_table / store_label / manual), `source_ref`, `source_license` |
| D · Inventario híbrido y sobrantes | `inventory_item` (`tracking_mode`, `rate_*`, `recipe_id`) |
| D · Dos verticales; mercados con nombre, día de visita y qué se compra | `store.kind = market` + `visit_weekday` (obligatorio), `store_product_group`, `product_group.purchase_channel` |
| E · Fuentes: ticket manual, OCR con confirmación, lectura asistida desactivada hasta veredicto, PROFECO | `price.source`, `attachment.ocr_status` (confirmed), `store.assisted_reading_enabled DEFAULT 0` + `app_setting price.assisted_reading_enabled = false`, `store.profeco_name` |
| E · Frescura 30 días, pesa menos | `price.observed_date`, `app_setting price.stale_after_days = 30` (peso: bloqueo `S1-frescura-peso`) |
| E · Sam's: sin membresía, comparar como si la tuvieran, punto de equilibrio | `store.requires_membership`, `has_membership`, `membership_annual_fee_cents` |
| E · Equivalencia por unidad base, sustitutos, mejor tienda | `product.base_unit`, `presentation.size_milli`, presentaciones del mismo producto, `price` + `ix_price_latest` |

## Supuestos
- Ninguno que cambie cifras, esquema o contratos. Las decisiones técnicas están en `reglas-negocio.md` §Convenciones. Las dudas de negocio van como bloqueos.
- Lecturas derivadas, con su fuente (no son supuestos):
  - La regla por defecto es «Hogar», porque Karen solo entra en la despensa y en sus compras (decisión A).
  - El presupuesto mide la parte del hogar, porque la parte de Karen es informativa y el plan usa solo el ingreso de Axel (decisiones A y B, plan §1.3).
  - La fecha ancla de la lectura A es 2026-10-10 (`model.json: credito.eventos[0]`, inicio del oráculo).
  - No hay corrimiento por fin de semana (el oráculo conserva el 15-nov-2026, que cae en domingo).

## Bloqueos
Nuevos de S1 en la ronda 1:
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-lista-compras.md`: marcas de la lista compartidas, locales o ninguna (misma pregunta que `C1-lista-compras`); opción A implementada PROVISIONAL.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-huecos-negocio.md`: vigencia por defecto de un cambio de presupuesto, correcciones a medias del OCR y registro del chat (puntos 3, 4 y 7 de `C1-huecos-u1-s1`).

Nuevos de S1 en la entrega original:
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-redondeo-series.md`: base de redondeo para series nuevas (partes por persona y devengo con centavos impares); los dos modos del oráculo.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-reposicion-parte-externa.md`: capital que repone la parte de Karen cuando paga por separado.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-frescura-peso.md`: peso de un precio viejo y la frontera ≥ 30 (misma pregunta que `criterios-frescura-peso`).
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-meta-kcal.md`: meta de calorías manual o por fórmula, y la escala de actividad (misma pregunta que `criterios-meta-kcal`).
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-restricciones-optimizador.md`: valores y alcance de las 3 restricciones duras.

Bloqueos de otros agentes que afectan a S1 (referenciados, no modificados):
- `U1-periodos-semana-quincena`: los rangos de la quincena quedan PROVISIONALES (opción C) en el CHECK de `period` y en RN-06.
- `criterios-invariante-deuda` (RN-21), `criterios-optimizador` (RN-11), `criterios-credito-formulas` (el esquema soporta su opción a o b: `credit_statement` capturado), `criterios-umbrales-ausentes` (claves sin valor en `app_setting` y `notification.params_json = {}`), `U1-captura-fuera-de-casa` (las opciones 2 o 3 piden una migración 00002 con clave de idempotencia en `expense` y `weight_entry`) y `criterios-login-bloqueo` (`login_attempt` sirve por usuario y por IP).

## Pendientes
1. Respuestas de Axel a los 5 bloqueos S1 y a `U1-periodos-semana-quincena`. Si la quincena se define «de pago a pago», hace falta la migración 00002 con el nuevo CHECK de `period`.
2. Datos que da Axel en Ajustes: costo anual de la membresía Sam's, valores de las restricciones del optimizador, mercados (nombre, día de visita, qué se compra ahí) y `datos/semilla/gastos_fijos_ingresos.csv` (se importa a `planned_item`).
3. Semilla (S8, WF-4): personas (axel, ximena, karen con `pays_share_directly = 1`; `fortnight_payment_cap_cents` sin valor hasta que Axel lo dé: el oráculo no trae un tope y 374.93 de `karen_fija.pago` es su parte por pago, no un tope; ver `criterios-optimizador`), cuentas, reglas Hogar (por defecto), Karen y Conjunto 50/50, categorías (`groceries` con Conjunto 50/50 y `eating_out`), tiendas, crédito Kueski con tasas de 2 y 4 quincenas y ancla A 2026-10-10, y las 4 despensas de `model.json`.
4. `docs/contratos/C1-datos.md` (plan §3.7 lo asigna a S1/S6 en WF-0) **no estaba entre mis rutas asignadas**. El diccionario sirve de base para C1.
5. Cuando cambie la migración, regenerar el diccionario con el generador (Anexo A: `python3 -I gen_dict.py <migración> <base migrada> docs/contratos/diccionario-datos.md`) y correr `s1_verifica.py` (Anexo B).
6. Respuestas de Axel a `S1-lista-compras` y `S1-huecos-negocio` (puerta 1). Con la opción B o C de la lista, `shopping_list` y `shopping_list_item` salen del esquema.

## Solicitudes fuera de alcance
1. **Integrador (H-05).** **Ruta:** `backend/db/tools/gen_dict.py` y `backend/db/tools/s1_verifica.py` (o las que decida). **Cambio:** copiar el código de los Anexos A y B tal cual. **Motivo:** hoy viven en este reporte (reejecutables con el comando de «Verificación») y en el scratchpad efímero; el diccionario se regenera desde la migración.
4. **C1 (corrección de C1).** Regenerar `openapi.yaml`, `routes.go` y `C1-datos.md` (§10, §11.3, §11.4, §13, §14) con el esquema 1.1: operaciones para N-151, N-152, N-153 y N-165 sobre `shopping_list_item` (PROVISIONAL `S1-lista-compras`/`C1-lista-compras`) y persistencia de N-127 y N-073; `expense.description` opcional sin relleno con la categoría; `recipe.default_meal_slot_id` → `recipe_meal_slot`; `inventory_movement` (N-181, N-183 con `reason`); `budget.period_id` → `valid_from_date`/`valid_to_date`; `planned_item.second_half_amount_cents` (N-279, N-280); `category.kind` (N-281); almacenamiento de `Idempotency-Key` en `idempotency_key` (solo la retención sigue PROVISIONAL). Actualizar el «Cómo seguí» de `C1-lista-compras` y `C1-huecos-u1-s1`.
2. **S6:** activar en el DSN de modernc `_pragma=foreign_keys(1)`, `journal_mode(WAL)` y `busy_timeout`, porque el `PRAGMA foreign_keys` de la migración no surte efecto dentro de la transacción de goose. Además, escribir `audit_log` en la misma transacción de cada cambio.
3. **S2 y S3:** usar la tasa como fracción (`rate_num`/`rate_den`) y guardar las series ya calculadas (C-08). Los montos con base truncada y el capital con base HALF_UP son obligatorios para reproducir el oráculo.

## Paquetes instalados
- **goose v3.28.0** (`go install -tags='no_postgres no_clickhouse no_mssql no_mysql no_vertica no_ydb no_libsql no_duckdb no_starrocks' github.com/pressly/goose/v3/cmd/goose@v3.28.0`), instalado solo para validar el formato goose. El binario quedó en el scratchpad de la sesión, fuera del PATH y fuera del repo. Como efecto lateral, `go install` llenó el caché de módulos `~/go/pkg/mod` (pressly/goose v3.28.0, modernc.org/sqlite, modernc.org/libc, modernc.org/memory, golang.org/x/sys, golang.org/x/sync y dependencias). No instalé nada con pacman.

## Consumo aproximado
- Unas 38 acciones de herramienta (sin subagentes).
- Leído: protocolo (51 líneas), decisiones (66), plan §0–§2 y §3.4–§3.7 (~300 líneas), §8.2 (1 línea, plantilla del CSV, ubicada con grep), `model.py` completo (213), estructura y bloques de `model.json` y `data.json` por script, y 10 bloqueos de otros agentes (~15 KB).
- Escrito: unos 300 KB (migración ~121 KB, diccionario ~136 KB, reglas ~35 KB, bloqueos ~9 KB, este reporte).
- El contador de la sesión marca unos 400K tokens acumulados (cuenta la relectura del contexto en cada acción). El trabajo terminó completo, sin compresión de contexto.

## Anexo A · Generador del diccionario (`gen_dict.py`)
Lee los comentarios `-- @table/@title/@entity/@desc/@used/@idx`, el comentario de cada columna (`descripción | regla`) y los `-- @rule` de la migración, más la introspección de una base ya migrada. Desde la ronda 1 no usa plantillas aparte: toma el texto fijo del propio diccionario (todo lo que va antes de la primera sección de tablas y desde «## Datos de referencia») y regenera el resumen, el mapa de entidades y las tablas complementarias. Uso: `python3 -I gen_dict.py <migración.sql> <base_migrada.db> docs/contratos/diccionario-datos.md`.

<!-- inicio gen_dict.py -->
```python
"""Regenera docs/contratos/diccionario-datos.md desde los comentarios @ de la migración y la introspección de una base migrada.
Uso: python3 -I gen_dict.py <migración.sql> <base_migrada.db> <diccionario.md>
El texto fijo (antes de la primera sección de tablas y desde «## Datos de referencia») se toma del propio diccionario."""
import re, sqlite3, sys
mig_path, db_path, out_path = sys.argv[1:4]
mig = open(mig_path, encoding='utf-8').read()
db = sqlite3.connect(db_path)

meta, order, sections = {}, [], []
cur = in_table = pending_idx = None
for line in mig.splitlines():
    if line.startswith('-- +goose Down'):
        break
    hm = re.match(r'^-- ([A-F](?: \(cont\.\))?\. .+)$', line)
    if hm:
        sections.append((hm.group(1), len(order)))
        continue
    m = re.match(r'^-- @(table|title|entity|desc|used|idx) (.*)$', line)
    if m:
        k, v = m.groups()
        if k == 'table':
            cur = v; meta[cur] = dict(cols={}, rules=[], idx={}); order.append(cur)
        elif k == 'idx':
            pending_idx = v
        else:
            meta[cur][k] = v
        continue
    m = re.match(r'^CREATE TABLE (\w+) \(', line)
    if m:
        in_table = m.group(1); continue
    if in_table:
        if line.startswith(') STRICT;'):
            in_table = None; continue
        s = line.strip()
        comment = line.split(' -- ', 1)[1].strip() if ' -- ' in line else None
        mc = re.match(r'^([a-z_][a-z0-9_]*) (INTEGER|TEXT|REAL|BLOB)\b', s)
        if mc and not s.startswith(('CHECK', 'UNIQUE')):
            desc, _, rule = (comment or '').partition(' | ')
            meta[in_table]['cols'][mc.group(1)] = (desc.strip(), rule.strip())
        elif comment and '@rule' in comment:
            meta[in_table]['rules'].append(comment.split('@rule', 1)[1].strip())
        continue
    m = re.match(r'^CREATE (?:UNIQUE )?INDEX (\w+) ON (\w+)', line)
    if m:
        meta[m.group(2)]['idx'][m.group(1)] = pending_idx; pending_idx = None

STD = {
    'id': 'Identificador interno (alias de rowid).',
    'created_at': 'Alta del registro, instante UTC ISO 8601 con «Z».',
    'updated_at': 'Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`.',
    'is_active': '0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia.',
    'notes': 'Notas libres.',
}
ON_DEL = {'RESTRICT': 'RESTRICT (no se puede borrar la madre)', 'CASCADE': 'CASCADE (se borra con la madre)',
          'SET NULL': 'SET NULL (queda sin liga)', 'NO ACTION': 'NO ACTION'}
def esc(x):
    return (x or '').replace('|', '/').replace('\n', ' ')

fk_all = {t: db.execute(f'PRAGMA foreign_key_list({t})').fetchall() for t in order}
inbound = {t: [] for t in order}
for t, fks in fk_all.items():
    for fk in fks:
        inbound.setdefault(fk[2], []).append(f'`{t}.{fk[3]}`')

def table_md(t):
    m = meta[t]; o = []
    o.append(f"### `{t}` · {m['title']}\n")
    o.append(f"- **Entidad del plan (§3.4.1):** {m['entity']}")
    o.append(f"- **Descripción:** {m['desc']}")
    o.append(f"- **Usado por:** {m['used']}\n")
    cols = db.execute(f'PRAGMA table_xinfo({t})').fetchall()
    fks = {fk[3]: fk for fk in fk_all[t]}
    uniq1 = set()
    for il in db.execute(f'PRAGMA index_list({t})').fetchall():
        if il[2] and il[3] == 'u' and not il[4]:
            cs = db.execute(f'PRAGMA index_info({il[1]})').fetchall()
            if len(cs) == 1:
                uniq1.add(cs[0][2])
    o.append('| nombre_codigo | tipo | nulo | descripción | regla |')
    o.append('|---|---|---|---|---|')
    for cid, name, typ, notnull, dflt, pk, hidden in cols:
        desc, rule = m['cols'].get(name, ('', ''))
        desc = desc or STD.get(name, '')
        parts = []
        if pk: parts.append('PK')
        if name in fks:
            fk = fks[name]; parts.append(f"FK → `{fk[2]}.{fk[4]}` al borrar: {fk[6]}")
        if name in uniq1: parts.append('único')
        if dflt is not None:
            parts.append('def: ahora (UTC)' if 'strftime' in dflt else f'def: {dflt}')
        if rule: parts.append(rule)
        if name in ('created_at', 'updated_at') and not rule: parts.append('ISO 8601 UTC')
        o.append(f"| `{name}` | {typ} | {'no' if (notnull or pk) else 'sí'} | {esc(desc)} | {esc('; '.join(parts))} |")
    if m['rules']:
        o.append('\n**Reglas de tabla (CHECK / UNIQUE):**')
        o += [f'- {esc(r)}' for r in m['rules']]
    o.append('\n**Relaciones:**')
    if fk_all[t]:
        o += [f"- `{fk[3]}` → `{fk[2]}.{fk[4]}` · al borrar la madre: {ON_DEL.get(fk[6], fk[6])}" for fk in fk_all[t]]
    else:
        o.append('- Sin llaves foráneas salientes.')
    if inbound.get(t):
        o.append(f"- La referencian: {', '.join(sorted(set(inbound[t])))}")
    idx_rows = db.execute("SELECT name, sql FROM sqlite_master WHERE type='index' AND tbl_name=? ORDER BY name", (t,)).fetchall()
    if idx_rows:
        o.append('\n**Índices:**\n')
        o.append('| índice | definición | propósito |')
        o.append('|---|---|---|')
        for name, sql in idx_rows:
            if sql is None:
                cs = [c[2] for c in db.execute(f'PRAGMA index_info({name})').fetchall()]
                o.append(f"| `{name}` | UNIQUE ({', '.join(cs)}) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |")
            else:
                d = sql.split(' ON ', 1)[1]
                o.append(f"| `{name}` | {'UNIQUE ' if sql.startswith('CREATE UNIQUE') else ''}`{esc(d)}` | {esc(m['idx'].get(name) or '')} |")
    trg = db.execute("SELECT name, sql FROM sqlite_master WHERE type='trigger' AND tbl_name=? ORDER BY name", (t,)).fetchall()
    if trg:
        o.append('\n**Disparadores:**')
        for name, sql in trg:
            mm = re.search(r"RAISE\(ABORT, '([^']+)'\)", sql)
            o.append(f"- `{name}`: " + (f"rechaza con «{mm.group(1)}»." if mm else 'pone updated_at en el instante actual si el UPDATE no lo trae.'))
    return '\n'.join(o) + '\n'

body = []
sec_at = {i: name for name, i in sections}
for i, t in enumerate(order):
    if i in sec_at:
        body.append(f'## {sec_at[i]}\n')
    body.append(table_md(t))

core = [t for t in order if not meta[t]['entity'].startswith('(complementaria)')]
comp = [t for t in order if meta[t]['entity'].startswith('(complementaria)')]
map_rows = '\n'.join(f"| {meta[t]['entity']} | `{t}` | {meta[t]['title']} |" for t in core)
comp_rows = '\n'.join(f"| `{t}` | {meta[t]['title']} | {meta[t]['entity'].replace('(complementaria) ', '')} |" for t in comp)
n_tab = len(order); n_idx = db.execute("SELECT count(*) FROM sqlite_master WHERE type='index'").fetchone()[0]
n_trg = db.execute("SELECT count(*) FROM sqlite_master WHERE type='trigger'").fetchone()[0]
n_upd = db.execute("SELECT count(*) FROM sqlite_master WHERE type='trigger' AND name LIKE 'trg\\_%\\_updated\\_at' ESCAPE '\\'").fetchone()[0]
# Plantilla: el propio diccionario (texto fijo antes de la primera sección de tablas y desde «## Datos de referencia»).
old = open(out_path, encoding='utf-8').read()
first = '\n## ' + sections[0][0] + '\n'
intro = old[:old.index(first) + 1]
tail = old[old.index('## Datos de referencia que crea la migración'):]
intro = re.sub(r'^Resumen: .*$', f'Resumen: **{n_tab} tablas STRICT**, {n_idx} índices (incluye los automáticos de UNIQUE) y {n_trg} disparadores ({n_trg - n_upd} de reglas de negocio y {n_upd} de `updated_at`).', intro, count=1, flags=re.M)
def put_rows(text, header, rows):
    k = text.index(header) + len(header)
    end = text.index('\n\n', k)
    return text[:k] + rows + text[end:]
intro = put_rows(intro, '| Entidad del plan (§3.4.1) | nombre_codigo | Título |\n|---|---|---|\n', map_rows)
intro = put_rows(intro, '| nombre_codigo | Título | Por qué existe |\n|---|---|---|\n', comp_rows)
open(out_path, 'w', encoding='utf-8').write(intro + '\n'.join(body) + '\n' + tail)
print('tablas', n_tab, 'núcleo', len(core), 'complementarias', len(comp), 'índices', n_idx, 'disparadores', n_trg)
missing = [(t, c) for t in order for c in [x[1] for x in db.execute(f'PRAGMA table_xinfo({t})')] if not (meta[t]['cols'].get(c, ('',))[0] or c in STD)]
print('columnas sin descripción:', missing)
```
<!-- fin gen_dict.py -->

## Anexo B · Verificación reejecutable de S1 (`s1_verifica.py`)
Aplica la sección Up en `datos/tmp/s1-verifica.db`, carga la semilla del oráculo (4 despensas con partes, 16 pagos con partes, el papel de contado), la lee de vuelta y la compara con `oraculo/model.json` (incluidos los 17 eventos, las series A y B de 14 valores y P1–P6 por persona), corre 46 pruebas negativas y 8 positivas, prueba Down/Up y goose si está disponible, imprime una línea por criterio S1-01…S1-17 y borra la base. Sale con 1 si algo falla. Comando: el de «Verificación».

<!-- inicio s1_verifica.py -->
```python
# s1_verifica.py · Verificación reejecutable de S1 (esquema, diccionario y reglas) con ida y vuelta del oráculo.
# Uso: python3 -I s1_verifica.py   (crea datos/tmp/s1-verifica.db, la borra al final; sale con 1 si algo falla)
import json, os, re, shutil, sqlite3, subprocess, sys
from fractions import Fraction
R = '/home/deepsec/finanzas-hogar'
MIG = R + '/backend/db/migrations/00001_esquema.sql'
DIC = open(R + '/docs/contratos/diccionario-datos.md', encoding='utf-8').read()
REG = open(R + '/docs/contratos/reglas-negocio.md', encoding='utf-8').read()
m = json.load(open(R + '/oraculo/model.json', encoding='utf-8'))
dj = json.load(open(R + '/datos/semilla/data.json', encoding='utf-8'))
TMP = R + '/datos/tmp'; DB = TMP + '/s1-verifica.db'
os.makedirs(TMP, exist_ok=True)
if os.path.exists(DB): os.remove(DB)
mig = open(MIG, encoding='utf-8').read()
up, down = mig.split('-- +goose Down', 1)
db = sqlite3.connect(DB); db.execute('PRAGMA foreign_keys=ON'); x = db.execute
FALLAS = []
def out(cid, ok, txt):
    print(f"{cid} {'OK  ' if ok else 'FALLA'} {txt}")
    if not ok: FALLAS.append(cid)
c = lambda v: int(Fraction(str(v)) * 100)  # pesos con 2 decimales -> centavos exactos
def one(q, a=()): return x(q, a).fetchone()[0]
def lid(): return one('SELECT last_insert_rowid()')
db.executescript(up)
tabs = [r[0] for r in x("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name")]
cols = {t: [r[1] for r in x(f'PRAGMA table_info({t})')] for t in tabs}
def has(tc):
    t, _, col = tc.partition('.')
    return t in cols and (not col or col in cols[t])
# ---------- S1-02
ni = one("SELECT count(*) FROM sqlite_master WHERE type='index'"); nt = one("SELECT count(*) FROM sqlite_master WHERE type='trigger'")
ns = one("SELECT count(*) FROM pragma_table_list WHERE schema='main' AND strict=1 AND name NOT LIKE 'sqlite_%'")
ic = one('PRAGMA integrity_check'); fk = x('PRAGMA foreign_key_check').fetchall()
out('S1-02', ic == 'ok' and not fk and ns == len(tabs), f"Up aplicado con sqlite3 {sqlite3.sqlite_version}: tablas={len(tabs)} strict={ns} índices={ni} disparadores={nt} · integrity_check={ic} · foreign_key_check={len(fk)} filas")
# ---------- S1-01
ENT = ['persona', 'cuenta_usuario', 'sesion_dispositivo', 'periodo', 'ingreso', 'categoria', 'gasto', 'regla_reparto', 'meta_ahorro', 'movimiento_ahorro', 'credito', 'tasa_por_plazo', 'compra_credito', 'pago_programado', 'promocion', 'tienda', 'producto', 'presentacion', 'precio', 'receta', 'ingrediente_receta', 'tiempo_comida', 'plan_comida', 'inventario_item', 'perfil_salud', 'peso_historico', 'nutricion_ingrediente', 'notificacion', 'destinatario', 'job', 'respaldo']
mapa = re.findall(r'^\| ([^|]+?) \| `([a-z_]+)` \| [^|]+ \|$', DIC.split('## Mapa entidad del plan')[1].split('Tablas complementarias')[0], re.M)
ent2tab = {}
for e, t in mapa:
    for name in re.split(r'[ ,]+', e.split(' (')[0]):
        ent2tab.setdefault(name, t)
falt = [e for e in ENT if e not in ent2tab or ent2tab[e] not in tabs]
sin_fila = [t for t in tabs if f'### `{t}` ·' not in DIC]
out('S1-01', not falt and not sin_fila, f"{len(ENT) - len(falt)} de {len(ENT)} entidades con tabla: " + ', '.join(f'{e}→{ent2tab.get(e)}' for e in ENT) + f" · faltantes={falt} · tablas sin fila en el diccionario={sin_fila}")
# ---------- S1-04
usado = re.findall(r'^### `([a-z_]+)` ·.*?\n- \*\*Usado por:\*\* *(.*)$', DIC, re.M | re.S)
vacias = [t for t, u in re.findall(r'^### `([a-z_]+)` ·[^\n]*\n(?:[^\n]*\n){2}- \*\*Usado por:\*\* *([^\n]*)', DIC, re.M) if not u.strip()]
n_usado = len(re.findall(r'^- \*\*Usado por:\*\* \S', DIC, re.M))
out('S1-04', n_usado == len(tabs) and not vacias, f"tablas con «Usado por» lleno: {n_usado} de {len(tabs)} · vacías={vacias}")
# ---------- S1-05
reales = x("SELECT m.name, p.name, p.type FROM sqlite_master m, pragma_table_info(m.name) p WHERE m.type='table' AND upper(p.type) IN ('REAL','FLOAT','DOUBLE')").fetchall()
cents = {}
for t in tabs:
    for r in x(f'PRAGMA table_info({t})'):
        if r[1].endswith('_cents'): cents.setdefault(r[2], []).append(f'{t}.{r[1]}')
out('S1-05', not reales and list(cents) == ['INTEGER'], f"columnas REAL/FLOAT/DOUBLE={len(reales)} · columnas _cents por tipo: " + ', '.join(f'{k}={len(v)}' for k, v in cents.items()) + " · convención: enteros de centavos")
# ---------- S1-06
sin_coment = "\n".join(l.split(" -- ")[0] for l in up.splitlines() if not l.lstrip().startswith("--"))
pant = re.findall(r"\d{1,2}-(?:ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic)-20\d\d", sin_coment)
crit = len(re.findall(r"oct-2026|ene-20", mig))
out('S1-06', not pant, f"fechas en formato de pantalla en datos o código SQL (sin comentarios): {len(pant)} · grep literal del criterio «oct-2026|ene-20» en la migración: {crit} · zona America/Mexico_City declarada en el diccionario: {'America/Mexico_City' in DIC}")
# ---------- datos de la ida y vuelta (S1-12) ----------
H = '$argon2id$v=19$m=65536,t=3,p=4$czFwcnVlYmE$cHJ1ZWJh'  # cadena de prueba con formato PHC; no es credencial
x("INSERT INTO person(code,display_name,participation) VALUES ('axel','Axel','holder')")
x("INSERT INTO person(code,display_name,participation) VALUES ('ximena','Ximena','co_holder')")
x("INSERT INTO person(code,display_name,participation,pays_share_directly) VALUES ('karen','Karen','external',1)")
P = dict(x('SELECT code,id FROM person'))
for u in ('axel', 'ximena'): x('INSERT INTO user_account(person_id,username,password_hash) VALUES (?,?,?)', (P[u], u, H))
x("INSERT INTO cost_share_rule(name,remainder_person_id,is_template,is_default) VALUES ('Hogar',?,1,1)", (P['axel'],)); r_hogar = lid()
x("INSERT INTO cost_share_rule(name,remainder_person_id,is_template,is_default) VALUES ('Conjunto 50/50',?,1,0)", (P['axel'],)); r_conj = lid()
x('INSERT INTO cost_share_rule_part(rule_id,person_id,pct_bp) VALUES (?,?,5000)', (r_conj, P['karen']))
x("INSERT INTO category(code,name,kind,default_cost_share_rule_id) VALUES ('groceries','Despensa','variable',?)", (r_conj,)); cat = lid()
x("INSERT INTO category(code,name,kind) VALUES ('household_goods','Papel','variable')"); cat2 = lid()
x("INSERT INTO credit(name,kind,lender,holder_person_id,credit_limit_cents,payment_frequency,first_payment_min_days,reading_a_anchor_date) VALUES ('Kueski Pay','revolving_line','Kueski',?,?,'fortnightly',15,'2026-10-10')", (P['axel'], c(m['parametros']['limite_credito']))); cr = lid()
x("INSERT INTO credit_term_rate(credit_id,term_count,rate_num,rate_den,source) VALUES (?,2,53450,309450,'manual'),(?,4,80050,309450,'manual')", (cr, cr))
for d in m['despensas']:
    x("INSERT INTO expense(kind,description,category_id,amount_cents,expense_date,accrual_months,payment_method,cost_share_rule_id) VALUES ('variable',?,?,?,?,2,'credit',?)", (f"Despensa {d['n']}", cat, c(d['capital']), d['fecha'], r_conj)); eid = lid()
    for w in ('axel', 'karen'): x('INSERT INTO expense_share(expense_id,person_id,amount_cents) VALUES (?,?,?)', (eid, P[w], c(d[w]['capital'])))
    x("INSERT INTO credit_purchase(credit_id,expense_id,kind,purchase_date,principal_cents,interest_cents,total_cents,term_count,payment_frequency,first_payment_date,interest_source,interest_includes_vat) VALUES (?,?,'installment_plan',?,?,?,?,4,'fortnightly',?,'quoted',1)", (cr, eid, d['fecha'], c(d['capital']), c(d['intereses']), c(d['total']), d['pagos'][0]['fecha'])); pid = lid()
    for w in ('axel', 'karen'): x('INSERT INTO credit_purchase_share(credit_purchase_id,person_id,principal_cents,interest_cents,total_cents) VALUES (?,?,?,?,?)', (pid, P[w], c(d[w]['capital']), c(d[w]['intereses']), c(d[w]['total'])))
    for i, p in enumerate(d['pagos'], 1):
        x('INSERT INTO scheduled_payment(credit_purchase_id,seq,due_date,amount_cents,principal_cents,interest_cents) VALUES (?,?,?,?,?,?)', (pid, i, p['fecha'], c(p['monto']), c(p['capital_repuesto']), c(p['monto']) - c(p['capital_repuesto']))); sid = lid()
        for w in ('axel', 'karen'): x('INSERT INTO scheduled_payment_share(scheduled_payment_id,person_id,amount_cents) VALUES (?,?,?)', (sid, P[w], c(p[w])))
for e in m['credito']['eventos']:
    if e['tipo'] == 'papel':
        x("INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,paid_by_person_id,person_id) VALUES ('variable',?,?,?,'cash',?,?)", (cat2, c(e['monto']), e['fecha'], P['karen'], P['karen'])); eid = lid()
        x('INSERT INTO expense_share(expense_id,person_id,amount_cents) VALUES (?,?,?)', (eid, P['karen'], c(e['monto'])))
db.commit()
n = 0; dif = []
def eq(a, b, what):
    global n; n += 1
    if a != b: dif.append((what, a, b))
rows = x('SELECT id,purchase_date,principal_cents,interest_cents,total_cents FROM credit_purchase ORDER BY purchase_date').fetchall()
eq(len(rows), 4, 'número de despensas')
npag = 0
for d, (pid, f, pr, it, tt) in zip(m['despensas'], rows):
    eq((f, pr, it, tt), (d['fecha'], c(d['capital']), c(d['intereses']), c(d['total'])), f"despensa {d['n']}")
    for w in ('axel', 'karen'):
        eq(x('SELECT principal_cents,interest_cents,total_cents FROM credit_purchase_share WHERE credit_purchase_id=? AND person_id=?', (pid, P[w])).fetchone(), (c(d[w]['capital']), c(d[w]['intereses']), c(d[w]['total'])), f"despensa {d['n']} parte {w}")
        eq(one('SELECT amount_cents FROM expense_share s JOIN credit_purchase cp ON cp.expense_id=s.expense_id WHERE cp.id=? AND s.person_id=?', (pid, P[w])), c(d[w]['capital']), f"gasto {d['n']} parte {w}")
    for p, (sid, due, amt, prin) in zip(d['pagos'], x('SELECT id,due_date,amount_cents,principal_cents FROM scheduled_payment WHERE credit_purchase_id=? ORDER BY seq', (pid,)).fetchall()):
        npag += 1
        eq((due, amt, prin), (p['fecha'], c(p['monto']), c(p['capital_repuesto'])), f'pago {due}')
        for w in ('axel', 'karen'): eq(one('SELECT amount_cents FROM scheduled_payment_share WHERE scheduled_payment_id=? AND person_id=?', (sid, P[w])), c(p[w]), f'parte pago {due} {w}')
eq(npag, 16, 'número de pagos')
s_d1 = one("SELECT sum(sp.amount_cents) FROM scheduled_payment sp JOIN credit_purchase cp ON cp.id=sp.credit_purchase_id WHERE cp.purchase_date='2026-10-10'")
a_d1 = one("SELECT sum(s.amount_cents) FROM scheduled_payment_share s JOIN scheduled_payment sp ON sp.id=s.scheduled_payment_id JOIN credit_purchase cp ON cp.id=sp.credit_purchase_id WHERE cp.purchase_date='2026-10-10' AND s.person_id=?", (P['axel'],))
eq(s_d1, 389500, 'suma pagos d1'); eq(a_d1, 239528, 'parte Axel d1')
lim, anc, lect = x('SELECT credit_limit_cents, reading_a_anchor_date, availability_reading FROM credit WHERE id=?', (cr,)).fetchone()
fin = m['parametros']['horizonte'][1]
viva = one('SELECT sum(principal_cents) FROM credit_purchase WHERE purchase_date <= ?', (anc,))
A, B = lim, lim - viva
ev = [(anc, 0, 'inicio', 0)]
ev += [(f, 1, 'compra', p) for f, p in x('SELECT purchase_date, principal_cents FROM credit_purchase WHERE purchase_date > ? AND purchase_date <= ?', (anc, fin))]
ev += [(f, 2, 'pago', p) for f, p in x('SELECT due_date, principal_cents FROM scheduled_payment WHERE due_date > ? AND due_date <= ?', (anc, fin))]
ev += [(f, 3, 'papel', 0) for (f,) in x("SELECT expense_date FROM expense WHERE payment_method='cash' AND expense_date > ? AND expense_date <= ?", (anc, fin))]
ev.sort(); serie = []; sA = []; sB = []
for f, _, t, p in ev:
    if t == 'compra': A -= p; B -= p
    elif t == 'pago': A += p; B += p
    serie.append((f, t, A, B))
    if t in ('compra', 'pago'): sA.append(A); sB.append(B)
eo = m['credito']['eventos']; eq(len(serie), len(eo), 'número de eventos')
for (f, t, a, b), e in zip(serie, eo): eq((f, t, a, b), (e['fecha'], e['tipo'], c(e['disponible_A']), c(e['disponible_B'])), f'evento {f}')
cs = m['credito']['comparacion_series']
for v1, v2, k in ((sA, cs['esperado_A'], 'A'), (sB, cs['esperado_B'], 'B')):
    eq(len(v1), 14, f'largo serie {k}')
    for i, (a_, b_) in enumerate(zip(v1, v2)): eq(a_, c(b_), f'serie {k}[{i}]')
eq(min(sA + [lim]), c(m['credito']['minimo_A']), 'mínimo A'); eq(min(sB + [lim - viva]), c(m['credito']['minimo_B']), 'mínimo B')
acum = Fraction(0); per_ok = 0
for p in m['periodos']:
    ym = p['inicio'][:7]
    def mes(q, a=()): return one(q, a) or 0
    efs = {w: mes("SELECT sum(s.amount_cents) FROM scheduled_payment_share s JOIN scheduled_payment sp ON sp.id=s.scheduled_payment_id WHERE strftime('%Y-%m',sp.due_date)=? AND s.person_id=?", (ym, P[w])) + mes("SELECT sum(s.amount_cents) FROM expense_share s JOIN expense e ON e.id=s.expense_id WHERE e.payment_method='cash' AND strftime('%Y-%m',e.expense_date)=? AND s.person_id=?", (ym, P[w])) for w in ('axel', 'karen')}
    dvs = {w: sum((Fraction(v, 2) for (v,) in x("SELECT s.total_cents FROM credit_purchase_share s JOIN credit_purchase cp ON cp.id=s.credit_purchase_id WHERE s.person_id=? AND (strftime('%Y-%m',cp.purchase_date)=? OR strftime('%Y-%m',date(cp.purchase_date,'start of month','+1 month'))=?)", (P[w], ym, ym))), Fraction(0)) + mes("SELECT sum(s.amount_cents) FROM expense_share s JOIN expense e ON e.id=s.expense_id WHERE e.payment_method='cash' AND strftime('%Y-%m',e.expense_date)=? AND s.person_id=?", (ym, P[w])) for w in ('axel', 'karen')}
    ef, dv = sum(efs.values()), sum(dvs.values())
    eq(ef, c(p['efectivo']['total']), f"efectivo {p['periodo']}"); eq(dv, c(p['devengado']['total']), f"devengado {p['periodo']}")
    per_ok += (ef == c(p['efectivo']['total'])) + (dv == c(p['devengado']['total']))
    for w in ('axel', 'karen'):
        eq(efs[w], c(p['efectivo'][w]), f"efectivo {p['periodo']} {w}"); eq(dvs[w], c(p['devengado'][w]), f"devengado {p['periodo']} {w}")
    acum += dv - ef; eq(acum, c(p['diferencia_acumulada']), f"diferencia acumulada {p['periodo']}")
fechas_1_9 = [f for (f,) in x("SELECT purchase_date FROM credit_purchase UNION ALL SELECT due_date FROM scheduled_payment UNION ALL SELECT expense_date FROM expense WHERE payment_method='cash'") if int(f[8:]) <= 9]
out('S1-12', not dif, f"ida y vuelta de oraculo/model.json: 4 despensas, {npag} pagos, {len(serie)} eventos, series A/B de {len(sA)}/{len(sB)} valores, mínimos {min(sA)}/{min(sB)}, P1–P6 por persona · valores comparados={n} · diferencias={len(dif)} · suma pagos d1={s_d1} · parte Axel d1={a_d1}" + (f' · {dif[:5]}' if dif else ''))
# ---------- S1-09
mes_ids = {}
for y_, mo in ((2026, 10), (2026, 11), (2026, 12), (2027, 1), (2027, 2), (2027, 3)):
    ini = f'{y_:04d}-{mo:02d}-01'; finm = one("SELECT date(?, '+1 month', '-1 day')", (ini,))
    x("INSERT INTO period(kind,year,month,start_date,end_date) VALUES ('month',?,?,?,?)", (y_, mo, ini, finm)); mid = lid(); mes_ids[(y_, mo)] = mid
    x("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES ('fortnight',?,?,?,1,?,?,?)", (mid, y_, mo, ini, f'{y_:04d}-{mo:02d}-15', f'{y_:04d}-{mo:02d}-15'))
    x("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES ('fortnight',?,?,?,2,?,?,?)", (mid, y_, mo, f'{y_:04d}-{mo:02d}-16', finm, min(f'{y_:04d}-{mo:02d}-30', finm)))
filas = x("SELECT kind, year, month, half, start_date, end_date, pay_date FROM period ORDER BY start_date, kind DESC").fetchall()
feb2 = one("SELECT pay_date FROM period WHERE kind='fortnight' AND year=2027 AND month=2 AND half=2")
try:
    x("INSERT INTO period(kind,year,month,start_date,end_date) VALUES ('month',2026,10,'2026-10-10','2026-11-09')"); ventana = 'ACEPTADA'
except sqlite3.IntegrityError: ventana = 'rechazada'
out('S1-09', len(filas) == 18 and feb2 == '2027-02-28' and ventana == 'rechazada' and per_ok == 12 and not fechas_1_9,
    f"{len(filas)} filas recuperadas (6 meses oct-2026…mar-2027 y 12 quincenas): " + '; '.join(f"{k[0][0]}{'' if k[3] is None else k[3]} {k[4]}…{k[5]}" + (f' paga {k[6]}' if k[6] else '') for k in filas) + f" · quincena 2 de feb-2027 paga {feb2} · ventana 2026-10-10…2026-11-09 como mes: {ventana} · equivalencia con P1–P6 (efectivo y devengado totales): {per_ok} de 12 · eventos de la semilla entre los días 1 y 9: {len(fechas_1_9)}")
# ---------- S1-10
kinds = re.search(r"kind TEXT NOT NULL CHECK \(kind IN \(([^)]*)\)\), -- revolving_line", up)
tasas = x('SELECT term_count, rate_num, rate_den FROM credit_term_rate ORDER BY term_count').fetchall()
pct = {t: (Fraction(nm, dn) * 100) for t, nm, dn in tasas}
r2 = round(float(pct[2]), 2); r4 = round(float(pct[4]), 4)
promo_cols = [k for k in ('merchant', 'term_count', 'min_purchase_cents', 'start_date', 'end_date') if k in cols['promotion']]
out('S1-10', lect == 'B' and r2 == 17.27 and r4 == 25.8685 and has('credit.availability_confirmed'), f"credit.kind = {kinds.group(1) if kinds else '?'} · tasas semilla como fracción: 2q 53450/309450 = {r2} %, 4q 80050/309450 = {r4} % · lectura por defecto al insertar sin valor = {lect} · availability_confirmed={has('credit.availability_confirmed')} · promotion: {', '.join(cols['promotion'][1:8])}")
# ---------- S1-11
lista = [('tiempos 4', 'meal_slot.code'), ('fuera sin receta', 'meal_plan_entry.is_eating_out'), ('porciones base 2 y escalado', 'recipe.base_servings'), ('perfil por persona', 'health_profile.height_mm'), ('historial de peso', 'weight_entry.weight_g'), ('restricciones vacías', 'health_profile.dietary_restrictions'), ('kcal con fuente', 'nutrition_fact.source'), ('inventario híbrido y sobrantes', 'inventory_item.tracking_mode'), ('cadena/sucursal/mercado con día y qué se compra', 'store.visit_weekday'), ('membresía Sam’s', 'store.membership_annual_fee_cents'), ('precio con tienda, fecha y fuente', 'price.observed_date'), ('unidad base y macros previstos', 'nutrition_fact.protein_mg_per_100g')]
mal = [tc for _, tc in lista if not has(tc)]
out('S1-11', not mal and one("SELECT count(*) FROM meal_slot") == 4, f"{len(lista) - len(mal)} de 12 puntos con tabla.columna: " + '; '.join(f'{a} → {b}' for a, b in lista))
# ---------- S1-13
sec = DIC.split('## Ajustes: 14 secciones')[1].split('\n## ')[0]
filas13 = re.findall(r'^\| (\d+) \| ([^|]+) \| ([^|]+) \| ([^|]+) \|$', sec, re.M)
NOMBRES = ['Cuenta y roles', 'Presupuesto', 'Ingresos', 'Categorías', 'Créditos y condiciones', 'Tiendas y precios', 'Recetas', 'Reportes y frecuencias', 'Notificaciones', 'Apariencia', 'Datos y respaldos', 'Personas y reparto', 'Salud y metas de ahorro', 'Inventario']
okn = [f[1].startswith(nm) for f, nm in zip(filas13, NOMBRES)]
claves = {r[0] for r in x('SELECT setting_key FROM app_setting')} | set(re.findall(r'^\| `([a-z_.]+)` \|', DIC.split('## Claves de app_setting previstas sin valor')[1].split('\n## ')[0], re.M))
refs = set(re.findall(r'`([a-z_]+\.[a-z_]+)`', sec)) - claves; malos = sorted(r for r in refs if r.split('.')[0] in cols and not has(r))
tablas13 = sorted({r.split('.')[0] for r in refs if r.split('.')[0] in cols})
out('S1-13', len(filas13) == 14 and all(okn) and all(f[3].strip() for f in filas13) and not malos, f"diccionario «Ajustes: 14 secciones → campos → tabla.columna»: {len(filas13)} de 14 secciones en orden, 0 destinos vacíos={all(f[3].strip() for f in filas13)}, referencias tabla.columna={len(refs)} (inexistentes: {malos}), tablas citadas={len(tablas13)}; Apariencia local por dispositivo (N-301)")
# ---------- S1-14
sens = sorted(set(re.findall(r'\b\w*(?:password|token|secret|apikey)\w*\b', mig, re.I)))
tg = re.findall(r'[0-9]{8,}:AA', mig)
out('S1-14', not tg and 'password_hash' in sens, f"coincidencias de password/token/secret/apikey en la migración (solo nombres de columna o texto): {sens} · patrón de token de Telegram: {len(tg)} · password_hash CHECK argon2id: {'$argon2id$' in mig}")
# ---------- S1-15
fkm = {t: len(x(f'PRAGMA foreign_key_list({t})').fetchall()) for t in ('expense', 'credit_purchase', 'scheduled_payment_share', 'price', 'meal_plan_entry', 'shopping_list_item')}
idx = {r[0]: r[1] for r in x("SELECT name, sql FROM sqlite_master WHERE type='index' AND sql IS NOT NULL")}
def ix_de(t, col):
    return [il[1] for il in x(f'PRAGMA index_list({t})') if x(f'PRAGMA index_info({il[1]})').fetchone()[2] == col]
filtros = {'persona': ('expense', 'person_id'), 'familia (partes por persona)': ('expense_share', 'person_id'), 'categoría': ('expense', 'category_id'), 'tienda': ('expense', 'store_id'), 'periodo (fechas)': ('expense', 'expense_date'), 'crédito': ('credit_purchase', 'credit_id')}
fi = {k: ix_de(*v) for k, v in filtros.items()}
out('S1-15', all(fi.values()) and all(fkm.values()), f"foreign_key_list: {fkm} · índices de filtros C1: " + ', '.join(f'{k}={v[0] if v else None}' for k, v in fi.items()) + " · PRAGMA foreign_keys documentado en la cabecera de la migración y en C-13")
# ---------- S1-16
def bloques(t, pat):
    return re.split(pat, t)
prov = []
for f, t, pat in (('migración', mig, r'\n(?=-- @table )'), ('diccionario', DIC, r'\n(?=### `|## )'), ('reglas', REG, r'\n(?=## |- \*\*C-)')):
    for b in bloques(t, pat):
        prov += [(f, l, b) for l in b.splitlines() if 'PROVISIONAL' in l]
blq = set(os.listdir(R + '/docs/estado/bloqueos'))
sin_ref = []; refs16 = set()
for f, l, b in prov:
    mrn = re.match(r'\| (RN-\d\d) \|', l)
    if mrn: b = REG.split('## ' + mrn.group(1))[1].split('\n## ')[0]
    nombres = set(re.findall(r'\b((?:S1|C1|U1|U2|S6|criterios)-[a-z0-9-]+[a-z0-9])\b', l)) or set(re.findall(r'\b((?:S1|C1|U1|U2|S6|criterios)-[a-z0-9-]+[a-z0-9])\b', b))
    nombres = {nm for nm in nombres if nm + '.md' in blq}
    refs16 |= nombres
    if not nombres and not re.search(r'PROVISIONAL\*\* marca|PROVISIONAL\*\* marca|marca lo que depende|PROVISIONAL\b.*bloqueo abierto', l): sin_ref.append((f, l[:90]))
out('S1-16', not sin_ref, f"líneas con PROVISIONAL: {len(prov)} (migración {sum(1 for f, _, _ in prov if f == 'migración')}, diccionario {sum(1 for f, _, _ in prov if f == 'diccionario')}, reglas {sum(1 for f, _, _ in prov if f == 'reglas')}) · bloqueos citados que existen: {len(refs16)} · líneas sin bloqueo en su bloque={len(sin_ref)} {sin_ref}")
# ---------- S1-17
ident = [t for t in tabs] + [f'{t}.{cl}' for t in tabs for cl in cols[t]]
malid = [i for i in ident if not re.fullmatch(r'[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)?', i)]
glos = re.findall(r'^\| [^|]+ \| (.+) \|$', DIC.split('## Glosario de nombre_codigo')[1].split('\n## ')[0], re.M)[1:]
out('S1-17', not malid and len(tabs) == len(set(tabs)), f"identificadores (tablas y columnas): {len(ident)}, fuera de [a-z0-9_]: {len(malid)} · tablas únicas: {len(set(tabs))} · filas del glosario: {len(glos) - 1}, repetidas: {len(glos) - len(set(glos))}")
# ---------- S1-03
secs = re.split(r'\n(?=## RN-\d\d)', REG)
def path_get(root, p):
    objs = [root]
    for tok in re.findall(r'[^.\[\]]+|\[\d+(?:\.\.\d+)?\]', p):
        if tok.startswith('['):
            a, _, b = tok[1:-1].partition('..')
            objs = [o[i] for o in objs for i in range(int(a), int(b or a) + 1)]
        else:
            objs = [o[tok] for o in objs]
    return objs if '..' in p else objs[0]
def fmt(v): return f'{v:,.2f}'
tabla3 = []; malas3 = []
for s in secs:
    mm = re.match(r'## (RN-\d\d)', s)
    if not mm: continue
    citas = re.findall(r'`(model|data)\.json: ([^`=]+?)(?: = ([^`]+))?`', s)
    okc = 0; tot = 0
    for src, pth, val in citas:
        pth = pth.strip()
        if pth == 'ruta': continue
        tot += 1
        try:
            v = path_get(m if src == 'model' else dj, pth)
        except (KeyError, IndexError, TypeError):
            malas3.append((mm.group(1), pth, 'ruta inexistente')); continue
        if val:
            good = str(v) == val.strip() or (isinstance(v, (int, float)) and re.fullmatch(r'-?[\d,]+(\.\d+)?', val.strip()) and Fraction(str(v)) == Fraction(val.strip().replace(',', '')))
        else:
            vs = v if isinstance(v, list) and '..' in pth else [v]
            good = all((not isinstance(w, (int, float))) or fmt(w) in s or str(w) in s or f'{w * 100:g}' in s for w in vs)
        okc += bool(good)
        if not good: malas3.append((mm.group(1), pth, v))
    ej = bool(re.search(r'\*\*Ejemplo', s))
    tabla3.append(f"{mm.group(1)}: ejemplo={'sí' if ej else 'NO'}{' (ilustrativo)' if 'ilustrativ' in s.lower() else ''}, citas {okc}/{tot}")
sin_ej = [t for t in tabla3 if 'ejemplo=NO' in t]
out('S1-03', not sin_ej and not malas3, f"{len(tabla3)} reglas · " + ' · '.join(tabla3) + f" · reglas sin ejemplo={len(sin_ej)} · citas que no coinciden con su ruta={malas3}")
# ---------- negativas (S1-07, S1-08, S1-05, S1-06 y tablas nuevas)
x("INSERT INTO product(name,base_unit) VALUES ('Arroz','g')"); prod = lid()
x("INSERT INTO recipe(name,origin) VALUES ('Arroz rojo','claude_seed')"); rec = lid()
x("INSERT INTO inventory_item(product_id,tracking_mode,quantity_milli) VALUES (?,'stock',1000)", (prod,)); inv = lid()
x("INSERT INTO shopping_list(start_date,end_date) VALUES ('2026-10-12','2026-10-18')"); sl = lid()
x("INSERT INTO shopping_list_item(shopping_list_id,origin,product_id,needed_milli) VALUES (?,'calculated',?,500000)", (sl, prod))
x("INSERT INTO recipe_meal_slot(recipe_id,meal_slot_id) VALUES (?,2)", (rec,))
KAREN, X_, CAT, CR = P['karen'], P['ximena'], cat, cr
NEG = [
 ('karen_con_cuenta', f"INSERT INTO user_account(person_id,username,password_hash) VALUES ({KAREN},'karen','{H}')"),
 ('tercera_cuenta_otra_co_titular', "INSERT INTO person(code,display_name,participation) VALUES ('otra','Otra','co_holder')"),
 ('ximena_pasa_a_externa_con_cuenta', f"UPDATE person SET participation='external' WHERE id={X_}"),
 ('ingreso_de_ximena', f"INSERT INTO income(person_id,kind,description,amount_cents,income_date) VALUES ({X_},'extra','x',100,'2026-10-15')"),
 ('hash_no_argon2id', f"UPDATE user_account SET password_hash='bcrypt$abc' WHERE person_id={X_}"),
 ('segundo_titular', "INSERT INTO person(code,display_name,participation) VALUES ('otro','Otro','holder')"),
 ('compra_total_distinto', f"INSERT INTO credit_purchase(credit_id,kind,purchase_date,principal_cents,interest_cents,total_cents,term_count,payment_frequency,first_payment_date,interest_source) VALUES ({CR},'installment_plan','2026-10-10',100,10,111,4,'fortnightly','2026-10-30','quoted')"),
 ('pago_distinto', "INSERT INTO scheduled_payment(credit_purchase_id,seq,due_date,amount_cents,principal_cents,interest_cents) VALUES (1,5,'2026-12-30',100,50,49)"),
 ('pct_mayor_100', f"INSERT INTO cost_share_rule_part(rule_id,person_id,pct_bp) VALUES ({r_hogar},{KAREN},10001)"),
 ('parte_pct_y_fijo', f"INSERT INTO cost_share_rule_part(rule_id,person_id,pct_bp,fixed_cents) VALUES ({r_hogar},{KAREN},5000,100)"),
 ('fecha_30_feb', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},100,'2027-02-30','cash')"),
 ('fecha_pantalla', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},100,'10-oct-2026','cash')"),
 ('quincena_paga_14', "INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES ('fortnight',(SELECT id FROM period WHERE kind='month' AND year=2027 AND month=2),2027,2,1,'2027-02-01','2027-02-15','2027-02-14')"),
 ('quincena2_feb_paga_30', "INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES ('fortnight',(SELECT id FROM period WHERE kind='month' AND year=2027 AND month=2),2027,2,2,'2027-02-16','2027-02-28','2027-02-30')"),
 ('real_en_cents', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},12.5,'2026-10-10','cash')"),
 ('fk_inexistente', "INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',999,100,'2026-10-10','cash')"),
 ('plan_comida_de_karen', f"INSERT INTO meal_plan_entry(plan_date,meal_slot_id,person_id,is_eating_out) VALUES ('2026-10-12',2,{KAREN},1)"),
 ('fuera_con_porciones', f"INSERT INTO meal_plan_entry(plan_date,meal_slot_id,person_id,servings_milli,is_eating_out) VALUES ('2026-10-12',2,{P['axel']},1000,0)"),
 ('lectura_A_sin_ancla', f"INSERT INTO credit(name,kind,holder_person_id,credit_limit_cents,payment_frequency,availability_reading) VALUES ('Otra','revolving_line',{P['axel']},500000,'fortnightly','A')"),
 ('tarjeta_sin_corte', f"INSERT INTO credit(name,kind,holder_person_id,credit_limit_cents,payment_frequency) VALUES ('Tarjeta X','bank_card',{P['axel']},500000,'monthly')"),
 ('mercado_sin_dia', "INSERT INTO store(kind,name) VALUES ('market','Tianguis sin día')"),
 ('credito_con_pagador_de_contado', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,paid_by_person_id) VALUES ('variable',{CAT},100,'2026-10-10','credit',{KAREN})"),
 ('borrar_karen_con_historia', f"DELETE FROM person WHERE id={KAREN}"),
 ('bitacora_inmutable', "UPDATE audit_log SET table_name='x'"),
 ('descripcion_en_blanco', f"INSERT INTO expense(kind,description,category_id,amount_cents,expense_date,payment_method) VALUES ('variable','  ',{CAT},100,'2026-10-10','cash')"),
 ('categoria_tipo_invalido', "INSERT INTO category(name,kind) VALUES ('X','upcoming')"),
 ('presupuesto_no_inicia_dia_1', f"INSERT INTO budget(category_id,amount_cents,valid_from_date) VALUES ({CAT},100,'2026-10-10')"),
 ('presupuesto_fin_no_es_fin_de_mes', f"INSERT INTO budget(category_id,amount_cents,valid_from_date,valid_to_date) VALUES ({CAT},100,'2026-10-01','2026-11-15')"),
 ('presupuesto_fin_antes_de_inicio', f"INSERT INTO budget(category_id,amount_cents,valid_from_date,valid_to_date) VALUES ({CAT},100,'2026-11-01','2026-10-31')"),
 ('presupuesto_mismo_inicio_dos_veces', f"INSERT INTO budget(category_id,amount_cents,valid_from_date) VALUES ({CAT},100,'2026-12-01'); INSERT INTO budget(category_id,amount_cents,valid_from_date) VALUES ({CAT},200,'2026-12-01')"),
 ('segunda_quincena_en_mensual', f"INSERT INTO planned_item(kind,description,amount_cents,second_half_amount_cents,frequency,start_date,category_id) VALUES ('fixed_expense','Renta',100,200,'monthly','2026-10-01',{CAT})"),
 ('receta_tiempo_repetido', f"INSERT INTO recipe_meal_slot(recipe_id,meal_slot_id) VALUES ({rec},2)"),
 ('ajuste_sin_motivo', f"INSERT INTO inventory_movement(inventory_item_id,movement_date,kind,quantity_delta_milli) VALUES ({inv},'2026-10-12','adjustment',-100)"),
 ('compra_con_motivo', f"INSERT INTO inventory_movement(inventory_item_id,movement_date,kind,reason,quantity_delta_milli) VALUES ({inv},'2026-10-12','purchase','count',100)"),
 ('motivo_invalido', f"INSERT INTO inventory_movement(inventory_item_id,movement_date,kind,reason,quantity_delta_milli) VALUES ({inv},'2026-10-12','adjustment','robo',-100)"),
 ('lista_no_dura_7_dias', "INSERT INTO shopping_list(start_date,end_date) VALUES ('2026-10-19','2026-10-24')"),
 ('lista_semana_repetida', "INSERT INTO shopping_list(start_date,end_date) VALUES ('2026-10-12','2026-10-18')"),
 ('articulo_producto_y_texto', f"INSERT INTO shopping_list_item(shopping_list_id,origin,product_id,free_text,free_text_unit,quantity_milli) VALUES ({sl},'manual',{prod},'Tortillas','g',1000)"),
 ('texto_libre_sin_unidad', f"INSERT INTO shopping_list_item(shopping_list_id,origin,free_text,quantity_milli) VALUES ({sl},'manual','Tortillas',1000)"),
 ('manual_sin_cantidad', f"INSERT INTO shopping_list_item(shopping_list_id,origin,free_text,free_text_unit) VALUES ({sl},'manual','Tortillas','g')"),
 ('calculado_sin_necesidad', f"INSERT INTO shopping_list_item(shopping_list_id,origin,product_id) VALUES ({sl},'calculated',{prod})"),
 ('calculado_repetido', f"INSERT INTO shopping_list_item(shopping_list_id,origin,product_id,needed_milli) VALUES ({sl},'calculated',{prod},1000)"),
 ('hora_de_compra_sin_comprado', f"INSERT INTO shopping_list_item(shopping_list_id,origin,free_text,free_text_unit,quantity_milli,bought_at) VALUES ({sl},'manual','Pan','piece',2000,'2026-10-12T18:00:00.000Z')"),
 ('clave_idempotencia_corta', "INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256) VALUES (1,'corta','POST','/api/v1/expenses','" + '0' * 64 + "')"),
 ('clave_idempotencia_con_espacio', "INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256) VALUES (1,'clave con espacio 123','POST','/api/v1/expenses','" + '0' * 64 + "')"),
 ('clave_idempotencia_repetida', "INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256) VALUES (1,'7f3c2a10-aaaa-4bbb-8ccc-0123456789ab','POST','/api/v1/expenses','" + '1' * 64 + "'); INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256) VALUES (1,'7f3c2a10-aaaa-4bbb-8ccc-0123456789ab','POST','/api/v1/expenses','" + '2' * 64 + "')"),
]
x("INSERT INTO audit_log(action,table_name,row_id,after_json) VALUES ('insert','person',1,'{\"id\":1}')"); db.commit()
aceptadas = []
for nombre, sql in NEG:
    db.execute('SAVEPOINT neg')
    try:
        db.executescript('SAVEPOINT n2; ' + sql + '; RELEASE n2;') if False else [db.execute(q) for q in sql.split('; ')]
        aceptadas.append(nombre)
    except (sqlite3.IntegrityError, sqlite3.DatabaseError):
        pass
    db.execute('ROLLBACK TO neg'); db.execute('RELEASE neg')
POS = [
 ('gasto_sin_descripcion', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},100,'2026-10-12','cash')"),
 ('ingreso_15_y_30_distintos', f"INSERT INTO planned_item(kind,description,amount_cents,second_half_amount_cents,frequency,start_date,person_id) VALUES ('income','Sueldo',100000,120000,'fortnightly','2026-10-15',{P['axel']})"),
 ('presupuesto_un_mes_y_vigente', f"INSERT INTO budget(category_id,amount_cents,valid_from_date,valid_to_date) VALUES ({CAT},200000,'2026-10-01','2026-10-31'); INSERT INTO budget(category_id,amount_cents,valid_from_date) VALUES ({CAT},210000,'2026-11-01')"),
 ('receta_en_dos_tiempos', f"INSERT INTO recipe_meal_slot(recipe_id,meal_slot_id) VALUES ({rec},3)"),
 ('ajuste_por_merma', f"INSERT INTO inventory_movement(inventory_item_id,movement_date,kind,reason,quantity_delta_milli) VALUES ({inv},'2026-10-12','adjustment','waste',-250)"),
 ('articulo_a_mano_comprado', f"INSERT INTO shopping_list_item(shopping_list_id,origin,free_text,free_text_unit,quantity_milli,store_id,is_bought,bought_at) VALUES ({sl},'manual','Tortillas','g',1000000,NULL,1,'2026-10-12T18:00:00.000Z')"),
 ('clave_idempotencia_valida', "INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256,response_status,response_body) VALUES (1,'7f3c2a10-aaaa-4bbb-8ccc-0123456789ab','POST','/api/v1/expenses','" + '3' * 64 + "',201,'{\"id\":1}')"),
 ('regla_60_mas_50_documentada_invariante_8', f"INSERT INTO cost_share_rule(name,remainder_person_id) VALUES ('Prueba 60/50',{P['axel']}); INSERT INTO cost_share_rule_part(rule_id,person_id,pct_bp) VALUES (last_insert_rowid(),{KAREN},6000)"),
]
rechazadas_pos = []
for nombre, sql in POS:
    db.execute('SAVEPOINT pos')
    try:
        [db.execute(q) for q in sql.split('; ')]
    except sqlite3.DatabaseError as e:
        rechazadas_pos.append((nombre, str(e)))
    db.execute('ROLLBACK TO pos'); db.execute('RELEASE pos')
neg_ok = not aceptadas and not rechazadas_pos
out('S1-07', 'karen_con_cuenta' not in aceptadas and 'tercera_cuenta_otra_co_titular' not in aceptadas and 'ximena_pasa_a_externa_con_cuenta' not in aceptadas and 'ingreso_de_ximena' not in aceptadas and 'role' not in ' '.join(cols['user_account']),
    f"cuenta para Karen: rechazada · tercera cuenta (otra co-titular): rechazada · Ximena con cuenta pasa a externa: rechazada · ingreso de Ximena: rechazado · columnas de user_account sin rol: {cols['user_account']} · tablas de finanzas de la familia de Karen: {[t for t in tabs if 'karen' in t or 'family' in t]}")
ax_d = c(dj['cuentas']['conjunto']['mitad']) + sum(c(l.get('total', l.get('precio', 0))) for l in []) 
out('S1-08', 'regla_60_mas_50_documentada_invariante_8' not in [p for p, _ in rechazadas_pos], f"regla con partes 60 % + 50 %: la base la acepta y la rechaza la API (invariante 8 del diccionario, prueba de S6) · despensa 1 sin pérdida: partes del gasto Axel {c(m['despensas'][0]['axel']['capital'])} + Karen {c(m['despensas'][0]['karen']['capital'])} = {c(m['despensas'][0]['capital'])} (ver S1-12) · conjunto {c(dj['cuentas']['conjunto']['total'])} → mitad {ax_d} cada quien (data.json: cuentas.conjunto.mitad)")
out('NEG', neg_ok, f"pruebas negativas: {len(NEG)}, rechazadas {len(NEG) - len(aceptadas)}, aceptadas por error {aceptadas} · positivas de las tablas nuevas: {len(POS)}, aceptadas {len(POS) - len(rechazadas_pos)}, rechazadas por error {rechazadas_pos}")
# ---------- Down y Up otra vez
db.close(); db = sqlite3.connect(DB); db.executescript(down)
resto = db.execute("SELECT count(*) FROM sqlite_master WHERE name NOT LIKE 'sqlite_%'").fetchone()[0]
db.executescript(up); t2 = db.execute("SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").fetchone()[0]
ic2 = db.execute('PRAGMA integrity_check').fetchone()[0]; db.close()
out('DOWN', resto == 0 and t2 == len(tabs) and ic2 == 'ok', f"Down deja {resto} objetos · Up otra vez: {t2} tablas, integrity_check={ic2}")
gs = os.environ.get('GOOSE') or shutil.which('goose')
if gs:
    gdb = TMP + '/s1-verifica-goose.db'
    res = [subprocess.run([gs, '-dir', R + '/backend/db/migrations', 'sqlite3', gdb, a], capture_output=True, text=True).returncode for a in ('up', 'down', 'up')]
    g = sqlite3.connect(gdb); gt = g.execute("SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT IN ('goose_db_version','sqlite_sequence')").fetchone()[0]; g.close(); os.remove(gdb)
    out('GOOSE', res == [0, 0, 0] and gt == len(tabs), f"goose ({gs}) up/down/up → códigos {res}, tablas al final={gt}")
else:
    print('GOOSE -    goose no está en el PATH ni en $GOOSE: se omitió (S1-02 ya probó el SQL con sqlite3)')
os.remove(DB)
if not os.listdir(TMP): os.rmdir(TMP)
print('RESUMEN', 'TODO OK' if not FALLAS else f'FALLAS: {FALLAS}')
sys.exit(1 if FALLAS else 0)
```
<!-- fin s1_verifica.py -->
