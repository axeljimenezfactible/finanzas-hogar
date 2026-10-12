# correccion-U4-r1 · P2c · corrección de los hallazgos medios de REVISION-U4-r1

Agente `correccion-U4-r1` (Opus) · flujo P2c · 11-oct-2026, noche · árbol de trabajo sobre 1ccb468, con C1 v1.3 ya regenerado en el repositorio (`openapi.yaml` `version: "1.3.0"`, 212 operaciones, 75 de Ajustes). Manda: `docs/estado/p2c/REVISION-U4-r1.md` (M-1 a M-5), `docs/estado/puerta2/RESPUESTAS.md`, los bloqueos que cito y `docs/estado/p2c/U4-puerta2.md`. No toqué B-1 ni B-2, que son bajos y no están en el encargo, ni nada «sin recomendación».

## Hecho
Versiones: parte a 1.5 → **1.6**, parte b 1.4 → **1.5**, U3 1.3 → **1.4**. El unido se rehízo igual que antes (encabezado de 38 líneas + a + línea en blanco + b; solo cambian en el encabezado las versiones y el estado de revisión). Respaldos en `/tmp/corr-u4-r1/bak/`.

**M-1 · preselección sin origen** (`U4-ajustes-b.md` l. 247 y 392).
- P-26.12 ya no da la preselección como decidida. Dice: «**La preselección es PROVISIONAL** (bloqueo abierto `U4-reparto-preseleccion`, se pregunta en la puerta 4)». En uso queda «Solo a los gastos nuevos», y se nombran las otras dos opciones.
- Bloqueo nuevo: `docs/estado/bloqueos/U4-reparto-preseleccion.md`. Separa HECHO, pregunta y propuesta. Opciones: A) «solo los nuevos», B) «también los ya capturados» y C) ninguna, con «Confirmar» inactivo hasta elegir. Yo no recomiendo ninguna; dejo anotado que `U4-puerta2` propuso A.
- Supuesto 6 reescrito: preselección PROVISIONAL. La edición posterior de un gasto ya no es «PREGUNTA abierta» sin archivo: ahora cita el bloqueo abierto `S1-reparto-edicion-gasto`.

**M-2 · dato falso sobre C1 v1.2** (b l. 247). El texto ahora dice: «C1 v1.2 no recibe alcance y ya se comporta como "solo los nuevos": las partes de los gastos ya guardados están congeladas en `expense_share`; "también los ya capturados" necesita `scope` = `all` de C1 v1.3». Ya no aparece «aplica las partes nuevas a todos» (0 apariciones).

**M-3 · bloqueos RESUELTOS citados como pendientes** (parte a l. 290, 295 y 296; parte b l. 261, 303, 308 y 383).
- Vista previa 4 (c): sigue PROVISIONAL, pero ahora cita el bloqueo abierto `S3-formula-tarjeta-pdf`. Explica que con Q2-29 a Q2-31 (a) la fórmula sale de los estados de cuenta, que todavía no se entregan. No muestra cifra. Lo FIRME ($0.00 si se paga el total a tiempo, meses sin intereses y abonos sin interés) cita el «Impacto» de ese bloqueo.
- Vista previa 9: **FIRME** con Q2-32 A. El corte cae en la primera fecha ≥ hoy, así que si hoy es día de corte, el corte es hoy. El pago cae en la primera fecha > corte. En un mes sin día 29, 30 o 31, los dos se recortan al último día. Se quitó el «Vista previa pendiente». Ejemplo comprobado con Python: hoy 11-oct-2026, corte día 31 y pago día 20 → 31-oct-2026, 20-nov-2026 y siguiente corte 30-nov-2026.
- Vista previa 10: **FIRME** con Q2-28 (a):
  - fórmula: r = anual ÷ 12, 24 o 52; con IVA, r × (1 + IVA); pago fijo y cada interés con ROUND_HALF_UP, y el último pago ajusta;
  - datos: capital, plazo, frecuencia y primera fecha salen de `listCreditPurchases` (campos comprobados en C1 1.3);
  - se calcula en el navegador, porque `Amortize` no usa saldo diario ni saldo del corte;
  - el IVA de un préstamo con `vat_included` vacío queda PROVISIONAL con el bloqueo abierto `S3-iva-lectura` (pregunta 1, opción en uso).
  - Las citas de S3 usan la numeración de la 1.5: §5.3 es la tarjeta y §5.4 el préstamo.
- P-26.13: se aplican Q2-16 A, Q2-17 C y Q2-18 A (`S2-metas` RESUELTO).
  - **Tablero de metas:** tres columnas, «En curso», «Cumplida» y «Vencida».
  - **Fila de solo lectura:** el avance va con un decimal (`progress_tenths_pct`). Los meses que faltan cuentan del mes de hoy al mes de la meta, los dos incluidos.
  - **Vista previa (1):** usa un decimal. Ejemplo de RN-18: 314,672 × 1,000 ÷ 629,344 = 500 décimas = 50.0 %.
  - **Vista previa de aportación y estado:** FIRME y calculada en el navegador. Ejemplo comprobado: 3 meses → $1,048.91 al mes, con la última de $1,048.90; con 2 meses → $1,573.36.
  - **Meta vencida sin cumplir:** su aportación queda PROVISIONAL con el bloqueo abierto `S2-metas-vencida`. Ahí la vista previa dice «vista previa pendiente».
- l. 383 (cobertura U4-05 y U4-11): `S2-metas` aparece como RESUELTO. Lista los bloqueos abiertos que se citan: `C1-hora-reportes-job`, `C1-reparto-alcance`, `S2-metas-vencida`, `S3-formula-tarjeta-pdf`, `S3-iva-lectura` y el nuevo `U4-reparto-preseleccion`.

**M-4 · U3** (`U3-dashboard.md` l. 94, 109, 160, 166, 170, 186, 190, 246 y 474 a 479).
- Tarjeta 1: FIRME (Q2-19 B y Q2-20 A).
- Tarjeta 2: FIRME con la regla de Q2-14 A. Aclara que `S2-flujo-quincenal` sigue abierto solo por Q2-12.
- Tarjeta 6:
  - métrica, insignias y marca con los tres estados de Q2-18 A: `in_progress` es «En curso», `achieved` es «Cumplida» y `overdue` es «Vencida»;
  - `progress_pct` sigue entero (Q2-17 C, tarjeta de Inicio);
  - no depende de `S2-metas-vencida`.
- D-6: aportación con Q2-16 A, estados de Q2-18 A, y PROVISIONAL solo la meta vencida (`S2-metas-vencida`).
- T-2: meses 7 a 12 sin PROVISIONAL (Q2-21 a Q2-24 C). La nota pasa a «Proyección con despensas repetidas», sin «por confirmar».
- §15: lista de bloqueos rehecha (abiertos y RESUELTOS por separado), más una solicitud a C1 sobre el enum de `SavingsGoal.status`.
- Fuera de las líneas del hallazgo: también corregí D-2 (l. 186). Marcaba PROVISIONAL el ahorro planificado, que Q2-13 B resolvió, y sin el cambio contradecía la lista nueva de §15. El gasto variable estimado sigue PROVISIONAL (Q2-12).

**M-5 · C1 v1.3** (b l. 26 a 28, 35, 41 a 45, 53, 373 y 393; parte a l. 381 y 383).
- P-26.08 lee y escribe día y hora con `getReportSettings` y `updateReportSettings`: 3 filas de campos, 5 filas de parámetros, el párrafo «Escritura», la matriz y el supuesto 7. P-26.09 sigue con `getAlertSettings` y `updateAlertSettings`.
- P-26.08 cita como PROVISIONAL el bloqueo abierto `C1-hora-reportes-job`. Pasa en dos lugares: en «Escritura» (qué manda, las claves o `job.schedule`, y si `next_run_at` se recalcula) y en la vista previa del «Próximo envío».
- P-26.12 escribe `scope` (`onward` o `all`; obligatorio con `parts` y sin valor por omisión) y cita como PROVISIONAL `C1-reparto-alcance` en «También a los gastos ya capturados». C1 1.3 no da el «{n}» ni la fecha, así que siguen como «Vista previa pendiente».
- Matriz de la parte a: suma `getReportSettings` y `updateReportSettings` (P-26.08) y `listCreditPurchases` (vista previa 10).

## Verificación
- `python3 -I /tmp/corr-u4-r1/scripts/verif.py` (script mío; lee los extractos de C1 1.3):
  - `openapi version: 1.3.0` · 212 operaciones en los extractos, 75 de Ajustes · el U4 unido cita 108 que existen · **inexistentes: [] · Ajustes sin cita: []**;
  - P-26.08: `getReportSettings` 9 y `updateReportSettings` 9. Queda 1 `getAlertSettings`: es la frase que remite las claves `alerts.*` a P-26.09. `C1-hora-reportes-job` aparece 2 veces;
  - P-26.09: `getAlertSettings` 12 y `getReportSettings` 0;
  - P-26.12: `U4-reparto-preseleccion` 1, `C1-reparto-alcance` 1, «aplica las partes nuevas a todos» 0, `onward` 1 y `all` 3;
  - frases viejas («`S3-formula-tarjeta.md`, preguntas 1 a 3)», «pregunta 4)», «`S3-convencion-prestamo.md`:», «no está decidido», «`S2-metas.md`, preguntas 1 y 3)» y las PROVISIONAL de `S2-metas`): 0 en las dos partes;
  - bloqueos citados en U4:
    - ABIERTOS: `S1-reparto-edicion-gasto`, `S2-metas-vencida`, `S3-formula-tarjeta-pdf`, `S3-iva-lectura`, `S5-terminos` y `U4-reparto-preseleccion`;
    - RESUELTOS y citados como tales: `S2-metas`, `S3-convencion-prestamo`, `S3-formula-tarjeta`, `U4-*` y `S1-parametros-avisos`;
    - `C1-hora-reportes-job` y `C1-reparto-alcance` no tienen una línea «Estado», pero los dos están abiertos (C1-puerta2, C1-datos §13);
  - campos de C1 1.3 citados que existen: `progress_tenths_pct`, `required_monthly_cents`, `vat_included`, `vat_on_interest_bp`, `principal_cents`, `term_count`, `first_payment_date` y `payment_frequency`. `CostShareRuleUpdate.scope` = `onward`/`all`: True;
  - parte a: 0 hex, 0 ISO, 214 filas de tabla y 0 filas con columnas distintas · parte b: 0 hex, 0 ISO, 153 filas y 0 filas con columnas distintas · secciones P-26: 14.
- U3 (script en línea): 6 tarjetas, 0 ISO y 0 hex.
  - Monto nuevo: solo $0.00, el «plan $0.00» de Q2-20 A.
  - U4 b suma $1,048.91, $1,048.90 y $1,573.36 (ejemplo de RN-18; los calculé en centavos con HALF_UP entero: 104891, 104890 y 157336).
  - No se quitó ningún monto.
- Unido: `U == encabezado(38) + a + "\n" + b` → True, 849 líneas. `make verify-wf1-u4` → «U4-ajustes.md termina con parte a + parte b sin cambios: True · secciones P-26: 14 de 14».
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → **rc=0** (registro en `/tmp/corr-u4-r1/verify-all.log`):
  - «WF-0: verificación completa OK»;
  - C1 «los 21 archivos coinciden byte a byte con el generador» y «comprobaciones OK: 29»;
  - S1 «RESUMEN TODO OK»;
  - U1 «comprobaciones OK: 63»;
  - Criterios con 0 diferencias;
  - oráculos S2 a S5 con 0 diferencias;
  - «WF-1: verificación completa OK» y «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».

## Supuestos
1. **Nombres de los estados.** C1 v1.3 **no** trae los nombres de Q2-18 A: `SavingsGoal.status` sigue con `on_track`, `behind` y `achieved` («en ruta, atrasada o cumplida»; sale de `GSTAT` en `backend/internal/api/c1gen/gen2_schemas.py`, l. 139). Por eso U3 y U4 usan los de `S2-interfaz.md` 1.3 §3.13 (`in_progress`, `achieved` y `overdue`) y dicen que C1 todavía no los tiene. No es una decisión nueva (los tres estados ya los fijó Q2-18 A): es una solicitud a C1.
2. **Insignias.** «Vencida» toma el rol advertencia y el icono `triangle-alert` que tenía «Atrasada»; «En curso» toma el rol información y el icono `info` que tenía «En ruta». Es diseño, no negocio.
3. **Un decimal.** P-26.13 cuenta como «Metas» para Q2-17 C, como pide la corrección del revisor. D-6 de U3 no cambió su porcentaje, porque eso no estaba en el hallazgo.
4. **Vista previa 4 (c).** Sigue sin mostrar cifra, aunque S3 tiene una fórmula general PROVISIONAL. Así queda como lo diseñó la versión anterior: no muestra cifras provisionales.
5. **Vista previa 10.** El préstamo sin compra registrada solo muestra la tasa antes → después. Ya no se pide SC-U4-4 para este bloque, porque la condición de esa solicitud (saldo diario o saldo del corte) no se cumple.
6. **D-2 de U3.** Lo corregí por coherencia con M-4 (mismo defecto y mismo bloqueo), aunque la l. 186 no estaba en la lista del revisor.

## Bloqueos
- **Nuevo:** `docs/estado/bloqueos/U4-reparto-preseleccion.md` (ABIERTO, se pregunta en la puerta 4).
- **Siguen abiertos y citados:**
  - `C1-hora-reportes-job` y `C1-reparto-alcance`;
  - `S2-metas-vencida`;
  - `S2-flujo-quincenal` (solo Q2-12);
  - `S3-formula-tarjeta-pdf` y `S3-iva-lectura`;
  - `S1-reparto-edicion-gasto`.

## Pendientes
- Revisor de P2c, ronda 2: U4 a 1.6, b 1.5 y U3 1.4.
- Bajos sin tocar: B-1 (falta una coma en la parte a, l. 398) y B-2 (`DESIGN.md` l. 933). También siguen abiertos R3-2 a R3-6.
- Cuando C1 alinee el enum de `SavingsGoal.status`, hay que repetir el cruce de campos de U3 y U4.
- Cuando se cierre `S3-formula-tarjeta-pdf`, la vista previa 4 (c) muestra cifras.

## Solicitudes fuera de alcance
1. **C1 (`c1gen`, `gen2_schemas.py`, `GSTAT`):** cambiar el enum de `SavingsGoal.status` (inicio, ajustes-b y presupuesto-b) a `in_progress`, `achieved` y `overdue`, con su descripción, según Q2-18 A y `S2-interfaz.md` 1.3 §3.13.
2. **U1 §2.5:** las columnas del tablero de metas pasan de «en ruta, atrasada y cumplida» a «En curso, Cumplida y Vencida» (Q2-18 A).
3. **Integrador (`docs/estado/bloqueos/INDICE.md`):** registrar `U4-reparto-preseleccion` (ABIERTO, U4, puerta 4).
4. **C1 (SC-U4-P2-1, sigue sin atender en 1.3):** el conteo y la fecha del gasto más antiguo que cambiaría «También a los gastos ya capturados». Depende de `C1-reparto-alcance`.
5. **Dueño de C1:** agregar la línea «Estado: ABIERTO» a `C1-hora-reportes-job.md` y a `C1-reparto-alcance.md`, como tienen los demás bloqueos.

## Propuestas para Axel
- En la puerta 4, contestar `U4-reparto-preseleccion` (A, B o C). No hay recomendación de este agente.

## Consumo aproximado
- 49 acciones de herramienta.
- Unos 150 KB leídos, por partes con `sed -n`, `grep` y `awk`.
- Unos 30 KB de texto nuevo neto:
  - partes a y b: unos 14 KB;
  - U3: unos 4 KB;
  - bloqueo: 2.8 KB;
  - este reporte: unos 9 KB.
- Además, unos 20 KB de scripts en `/tmp/corr-u4-r1/scripts/`, y el unido (unos 158 KB), que se rehízo por concatenación sin texto propio.
