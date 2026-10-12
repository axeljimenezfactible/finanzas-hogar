# P2c · contratos-puerta2 (Opus) · 11-oct-2026, noche

Alcance: `docs/contratos/motores/S2-interfaz.md` (1.2 → 1.3), `S4-interfaz.md` (1.4 → 1.5), `S5-interfaz.md` (1.4 → 1.5), `S3-interfaz.md` (solo marcas de Q2-28 a Q2-38; versión sin cambio), `docs/estado/bloqueos/S2-ventanas.md` (pregunta 4) y 4 bloqueos nuevos. Fuentes: `docs/estado/puerta2/RESPUESTAS.md` (modo «lo recomendado + excepciones», sin excepciones), `PUERTA-2.md` §2, §3 y §8, `REVISION-S2-r3.md` (M1), `REVISION-S5-r3.md` (H-r3-01 a 05), `REVISION-contrato-S4S5-r5.md` (M-r5-1 a 3) y `correccion-S5-r2.md` (solicitudes 1 a 6). No hice `git add` ni `commit`; no toqué la Pista 3.

## Hecho
- **S2 1.3.** §2.12: ventanas «contiguas, ascendentes y disjuntas», con 4 reglas: O1, A y H1 FIRMES (Q2-27, Q2-25, Q2-26) y la contigüidad PROVISIONAL (`S2-ventanas` pregunta 4, caso 1: toda lista con hueco es `invalid_argument`), con las cifras A + C. §3.2: el caso 2 (ventanas que no duran un mes) queda escrito como regla PROVISIONAL con el comportamiento actual y sus cifras (97375 contra 292125 al 2026-11-09); no lo decidí. FIRMES: Q2-13 y Q2-14 (§3.11), Q2-15 a Q2-18 (§3.13), Q2-19 y Q2-20 (§3.8) y Q2-21 a Q2-24 (§2.10 y §3.12). La marca `provisional` de §3.11 y §3.12 queda en `null`. Fila nueva de `s2-ventanas-sin-cubrir` en §5. §8 rehecha.
- **`S2-ventanas.md`:** pregunta 4 completa: caso 1 con opciones 1A a 1D y caso 2 con 2A a 2C, cifras del informe y sin recomendación. Ajusté 2 líneas de la sección «sin respuesta» que decían que la pregunta aún no estaba escrita.
- **Bloqueos nuevos (ABIERTOS):** `S2-ahorro-quincenal` (base del ahorro por quincena y si resta del saldo; el reparto de Q2-13 = B ya es FIRME), `S2-metas-vencida` (aportación de una meta vencida sin cumplir: motor y oráculo difieren), `S4-presentacion-respaldo` (M-r5-2, variantes i, ii y iii) y `S5-carritos-granel` (H-r3-02).
- **S4 1.5.** M-r5-1: `quotes` con `choice_store` y `choice_presentation`; §8.1 paso 1 toma exactamente la de esa elección (ejemplo del arroz resuelto: Walmart 1 kg a 2800). M-r5-2: el respaldo sin cotización queda PROVISIONAL con la regla que el motor aplica hoy (leí `shopping.go` y `catalog.go`): primer ingrediente no nulo, luego la presentación típica y luego la primera del catálogo. M-r5-3: «Orden de las listas de salida» en §12.2. Q2-39 a Q2-46 FIRMES (§6.2, §6.4, §6.5, §10.1, §12.1 y §12.2). La variante sin campo nuevo de C1 queda descartada.
- **S5 1.5.** Solicitudes 1 a 3 de `correccion-S5-r2` (§11 archivos y base, filtro A, `rejections[]` y `source_ref`; §7.2 conteos, sustituto, regla de 3 canastas o más y orden; §12 con tickets-README y tickets-sinteticos). H-r3-02: §7.2 paso 8 PROVISIONAL. H-r3-03: «Lectura de un registro» con la gramática del precio, la lista cerrada de motivos y la cadena vacía como rechazo (regla técnica); la liga normalizada queda PROVISIONAL por `S5-profeco-filtro`; `loaded_proposed` en §11 y §15.2. H-r3-04 y M-r5-3: un solo significado de `cheapest_store` y `difference_vs_cheapest_cents` en los dos modos (§7.4) y orden de cada lista (§15.2). H-r3-05: opciones de Tesseract. FIRMES: Q2-43, Q2-47, Q2-49, Q2-53 a Q2-55, Q2-59 a Q2-63 y Q2-76. Siguen PROVISIONAL «en uso A»: Q2-48, Q2-50 a Q2-52 y Q2-56 a Q2-58.
- **S3 (solo marcas).** FIRMES: Q2-28, Q2-32, Q2-33 a Q2-36, Q2-37 y Q2-38, en §5.1, §5.3, §5.4, §5.5, §5.8, §6, §8.2, §8.3, §8.6, fixtures, pruebas, §14 y §14.1. Las listas `Provisional` ya no llevan los 4 bloqueos resueltos. La fórmula de tarjeta y tienda sigue PROVISIONAL (`S3-formula-tarjeta`, pendiente de los PDF de Axel). Consecuencia de Q2-38 = (a): la línea pasa de RN-07 literal a capital = monto.
- `oraculo/fixtures/FORMATO.md`: sin cambio; ninguna regla de orden lo exige, porque el orden vive en cada contrato.

## Verificación
- Reemplazos con conteo exacto (`/tmp/p2c-contratos/*.py`, fuera del repo): S2 25 + 1 + 1, S4 24, S5 39 y S3 38 + 7 + 3. Cada script aborta sin escribir si un texto no aparece exactamente las veces esperadas.
- Revisión de PROVISIONAL contra el estado de los bloqueos (script en Python): S2 cita `S2-ventanas`, `S2-flujo-quincenal`, `S2-ahorro-quincenal` y `S2-metas-vencida`, todos ABIERTOS; S4 cita `S4-presentacion-respaldo` (ABIERTO); S5 cita `S5-carrito-comparativa`, `S5-carritos-detalle`, `S5-carritos-granel` y `S5-profeco-filtro`, todos ABIERTOS. En S3 quedan `S3-formula-tarjeta` y `U4-iva-credito`, que el agente de bloqueos marcó RESUELTO (ver Bloqueos). Las demás coincidencias son textos de bitácora o definiciones de «Marcas».
- Q2 aplicadas (grep): S2 Q2-01, 11 a 21 y 24 a 27 (Q2-22 y Q2-23 van en el rango «Q2-21 a Q2-24»); S3 Q2-11 y 28 a 38 (Q2-30 en «Q2-29 a Q2-31»); S4 Q2-06 a 08 y 39 a 46; S5 Q2-02 a 05, 08, 39, 42, 43, 47 a 63 y 76.
- `git diff --stat`: S2 +67/−, S3 +88/−, S4 +57/−, S5 +81/−, S2-ventanas +61; 4 archivos nuevos en `docs/estado/bloqueos/`.
- `make verify-wf1-cobertura`: K25 OK en todos los paquetes. `make verify-criterios`: falló con `grep: datos/tmp/verify-wf0/criterios.log: No such file or directory` antes de evaluar. No toqué criterios ni el script; la carpeta temporal la comparten otros agentes. Queda para el integrador.
- No corrí `verify-wf1-oraculos`: no cambié código ni fixtures, y los cambios de motor de abajo se hacen en WF-3.

## Supuestos
1. Donde motor y oráculo difieren y nada está decidido (granel, presentación de respaldo, meta vencida), la regla PROVISIONAL es la del motor, siguiendo la instrucción de M-r5-2.
2. H-r3-03: tomé la gramática del precio, los motivos y la cadena vacía como regla técnica (FIRME), con la forma del motor; 0 rechazos en los 37,548,200 registros reales. La liga normalizada no es técnica: depende de Q2-57.
3. `difference_vs_cheapest_cents` `null` en tiendas incompletas: apliqué M-r5-3 (Q2-08 = A, la ronda más reciente) sobre la solicitud 3 de `correccion-S5-r2` («puede ser negativa»). Las dos venían aprobadas en bloque y chocan.
4. Q2-18 = A, «tres estados»: un fondo sin monto o sin fecha es `in_progress`, y `achieved` depende solo del monto.
5. Q2-13 = B no dice la base. Por eso `savings_planned_cents` sigue en `null` y abrí `S2-ahorro-quincenal`, en vez de elegir una base.
6. Los órdenes que el revisor no propuso (`leftovers[]`, `references[]`, `stale[]`, `coverage[]`, `missing_inputs[]`, `days[]`) los fijé yo como regla técnica.
7. S3: un corte `computed` sigue PROVISIONAL (`S3-formula-tarjeta`) y uno `bank` queda en `[]`. No subí la versión de S3, porque la orden era «solo marcas».
8. La cifra de la opción 2B de la pregunta 4 (97375 para `[d, r]`) sale de la regla; ningún motor la implementa.

## Bloqueos
- `S3-formula-tarjeta` dice «RESUELTO», pero Q2-29 a Q2-31 = (a) deja la fórmula esperando los estados de cuenta de Axel. Las 14 marcas PROVISIONAL de S3 lo citan. Tiene que seguir ABIERTO (preguntas 1 a 3) hasta que lleguen los PDF.
- `U4-iva-credito` está «RESUELTO», pero S3 todavía lo cita en 12 marcas PROVISIONAL. Q2 de U4, fuera de mi rango Q2-28 a Q2-38: no lo toqué.
- Abiertos que pasan a WF-3: `S2-ventanas` pregunta 4, `S2-flujo-quincenal` (Q2-12), `S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo`, `S5-carritos-granel`, `S5-carrito-comparativa` (Q2-48), `S5-carritos-detalle` (Q2-50 a 52) y `S5-profeco-filtro` (Q2-56 a 58).

## Pendientes
- Alta de los 4 bloqueos nuevos en `docs/estado/bloqueos/INDICE.md`: es del integrador o del agente de bloqueos.
- §5 de S2 dice «Total: 27 fixtures», pero `S2-ventanas` habla de 30. No lo verifiqué ni lo cambié.
- Los cambios de motor y oráculo de abajo (WF-3).

## Cambios que deben implementar los motores
**S2 (`backend/internal/motor_fin`) y su oráculo (`oraculo/fin_*.py`)**
- §2.12 regla 4: rechazar con `invalid_argument` toda lista con hueco (`start[i] ≠ end[i−1] + 1 día`), aunque no tenga eventos (`checkCovered` y `sin_cubrir`).
- `s2-ventanas-sin-cubrir`: agregar la consulta A + C (`invalid_argument`) y `[d = 2026-10-10, r = 2026-10-11 a 2026-11-09]` (292125 al 2026-11-09, `provisional` = `S2-ventanas` pregunta 4).
- `provisional` siempre `null` en `flujo_quincenal` y `proyeccion` (Q2-14 y Q2-21 a 24). Las fixtures `s2-liga-provisional`, `s2-proyeccion-7-12-cadencia` y `s2-metas` pasan a `provisional` null.
- `metas`: `state` solo `in_progress`, `achieved` u `overdue` (el motor quita `open`; `achieved` solo exige `target_cents`). `months_left` = `max(0, mes(target) − mes(hoy) + 1)`: el motor da hoy 0 si `target_date` < hoy dentro del mes. Meta vencida: el oráculo pasa a `plan` = [mes de hoy, restante] y `needed` = restante. Las fixtures ya pueden comparar `state`.
- Quitar los comentarios PROVISIONAL de A, H1, O1, la liga, la cadencia, las metas y el presupuesto.

**S3 (`motor_cred`) y su oráculo**
- Línea con promoción al 0 % y n ≥ 2: capital = monto (Q2-38 = (a)); hoy usa RN-07 literal.
- Listas `Provisional` según la tabla nueva de §14.1: tarjeta, préstamo con IVA, línea, MSI y tienda en `[]`; corte `computed` con `S3-formula-tarjeta` y corte `bank` en `[]`; `Scenario` y `Result` sin `S3-objetivos-medida`. Ajustar las fixtures de la tabla de §14.1 y `TestProvisionalListas`.

**S4 (`backend/internal/recetas`) y su oráculo (`rec_*.py`)**
- `quotes` con `choice_store` y `choice_presentation`, más la nueva unicidad; §8.1 paso 1 por elección exacta. Actualizar `s4-costo-cotizacion-elegida` con el ejemplo r5 y `TestCostoCotizacionQueRige`.
- Oráculo: respaldo igual al del motor (primer no nulo, luego `presentacion_tipica`, luego la primera del catálogo). Fixture nueva con los casos (a), (b) y (c).
- Verificar el «Orden de las listas de salida» de §12.2. Ojo con `items[].recipes[]`: el motor las agrega en orden del plan y el oráculo las ordena.

**S5 (`backend/internal/precios`) y su oráculo (`pre_*.py`)**
- §7.4: completas primero y luego incompletas, con diferencia `null`; `cheapest_store` `null` si ninguna está completa (`rankCarts`, `compare.go:356`; `pre_calculo.py:251`). Regenerar `s5-carritos-presentaciones`: `tercera` va al final con `null`.
- Granel (PROVISIONAL A): el oráculo deja de redondear por línea. Fixture con los casos del revisor: totales 101 y 103; precios por unidad base 135, 137, 721 y 736; sustituto 2933 y 2800.
- PROFECO: el oráculo rechaza 0.00 y la cadena vacía, acepta «1,021.00» = 102100, liga normalizado y usa los textos cerrados. El motor emite «CSV mal formado» sin detalle. Los dos emiten `loaded_proposed`. Fixture con los 17 registros de `rv-bordes` (producción: rechazados 4, filtrados 1, saltados 1, cargados 10).
- Verificar el orden de §15.2.

## Solicitudes a C1
- `ShoppingListItem.assigned_store_id` (campo calculado nuevo) = `items[].store`.
- `groups[].store_id` = «la tienda asignada (`assigned_store_id` de sus artículos)».
- `recommended_store_id` = «recomendación de S5 por precio; rige precio y presentación mientras `store_id` esté ausente; en vertical `market` no rige la agrupación».
- `estimated_price_cents` = «costo con la cotización que rige (S4 §8.1)».
- `groups[].visit_weekday` = día ISO de `visit_date` (opcional: `visit_date`).
- `groups[].total_cents` y `ShoppingList.total_cents` = «Σ de los artículos con precio; los que no tienen precio no suman y se reconocen por `estimated_price_cents` ausente» (Q2-42 = A).
- `BasketComparison`: `cheapest_store_id` y `difference_vs_cheapest_cents` admiten `null` (§7.4). `identical_count` = «Renglones idénticos».

## Propuestas para Axel
Contestar por id: `S2-ventanas` pregunta 4 (casos 1 y 2), `S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo`, `S5-carritos-granel` y Q2-12 («C o A»). Entregar los PDF de un estado de cuenta de tarjeta y uno de crédito de tienda (Q2-29 a Q2-31 = (a)).

## Consumo aproximado
59 acciones de herramienta · unos 300 KB leídos · unos 95 KB escritos (contratos, 5 bloqueos y este reporte).
