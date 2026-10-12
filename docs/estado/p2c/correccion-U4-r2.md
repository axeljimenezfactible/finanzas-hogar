# correccion-U4-r2 · P2c · corrección de la ronda 2 de U4

Agente `correccion-U4-r2` · flujo P2c · 11-oct-2026, noche · árbol de trabajo sobre HEAD 4079cee. Corrige los dos hallazgos medios (M-1 y M-2) de `docs/estado/p2c/REVISION-U4-r2.md`. Esa revisión no tiene hallazgos altos. Los bajos (B-1, B-2 y B-3) no estaban en el encargo y no los toqué. No hice `git add` ni `git commit`.

## Hecho
Archivo tocado: `docs/design/U3-dashboard.md`, de la 1.4 a la **1.5**. Cambian 6 líneas y se agregan 2. No toqué `U4-ajustes-a.md`, `U4-ajustes-b.md`, `U4-ajustes.md` ni `DESIGN.md`: su md5 antes y después es el mismo.

1. **M-1 · D-2, ahorro planificado** (l. 186). Antes decía «ahorro planificado FIRME (Q2-13 B…)». Ahora dice:
   - el reparto mitad y mitad es FIRME (Q2-13 B: HALF_UP en la quincena 1 y el resto en la 2);
   - el valor es **PROVISIONAL** y cita el bloqueo abierto `docs/estado/bloqueos/S2-ahorro-quincenal.md`, con lo que falta: la base, qué metas cuentan y si resta del saldo;
   - el motor S2 lo entrega en `null` y no lo resta de `balance_cents` ni de `cumulative_balance_cents` (`S2-interfaz.md` 1.3 §3.11, HECHO);
   - mientras llegue `null`, la fila muestra «—» y nunca un monto en cero (regla «Sin datos no es cero» de U3).
   - En §15 «Abiertos» (l. 478) entra `S2-ahorro-quincenal`. En «RESUELTOS» (l. 479) sale «ahorro planificado de D-2»: queda «tarjeta 2; Q2-13 B fija solo el reparto del ahorro planificado, cuyo valor en D-2 sigue abierto».
2. **M-2 · D-6, avance con un decimal** (l. 190).
   - Cambio:
     - D-6 pide `progress_tenths_pct` en lugar de `progress_pct`. Muestra el % con un decimal (décimas con ROUND_HALF_UP, C1 v1.3).
     - Cita Q2-17 C literal.
     - Sin `target_cents`, el campo llega ausente y se muestra «—» (descripción de C1).
     - En §15, la fila de `S2-metas` dice «tarjeta 6 con el avance entero y D-6 con un decimal».
   - Lectura que lo sostiene. Usé la opción A de la revisión, porque D-6 sí cae en «Metas»:
     - D-6 cita la métrica **M-16 «Metas» (N-111, N-027)** de C1-datos §9.
     - El bloqueo resuelto `S2-metas.md` dice en «Qué cambia» (l. 22): «Alinear RN-18 (un decimal) con M-07 y M-16 (entero): un decimal en Metas (`progress_tenths_pct`) y entero en la tarjeta de Inicio (`progress_pct`)». M-16 queda del lado del decimal y M-07 (la tarjeta 6) del lado del entero.
     - La descripción de `progress_tenths_pct` en C1 1.3 y `S2-interfaz.md` 1.3 l. 295 dejan el entero solo en la tarjeta.
     - Por eso no abrí `U3-avance-detalle.md` (ver Supuestos).
3. **Nota de versión 1.5** (l. 3). Resume M-1 y M-2 y dice que ninguna cifra de canon cambia.
4. **§15, dos solicitudes nuevas a C1** (l. 475 y 476). Las piden M-1 y M-2. Solo piden; no deciden.
   - C1-datos M-16 todavía dice «% de avance entero». Hay que alinearla con Q2-17 C, como pide el «Qué cambia» de `S2-metas`.
   - `rows[].planned_savings_cents` de `getCashFlowFortnights` está declarado en C1 1.3 como `MoneyCents` obligatorio. La descripción de `balance_cents` es «income − egresos − ahorro», que da por contestada la pregunta 3 de `S2-ahorro-quincenal`. La solicitud pide marcarlo PROVISIONAL con ese bloqueo, admitir `null` y no dar por contestada la pregunta 3.

## Verificación
- Respaldo previo: `/tmp/correccion-U4-r2/U3-dashboard.md.bak` (md5 811093cd…). Cambios aplicados con `python3 -I /tmp/correccion-U4-r2/edit.py`, siete reemplazos exactos, cada uno con 1 aparición comprobada: «OK: 7 reemplazos».
- Verificación propia, `python3 -I /tmp/correccion-U4-r2/verif.py`:
  - Forma (U3-11): `hex: 0 · rgb(): 0 · fechas ISO: 0`.
  - Bloqueos citados en D-2 y D-6:
    - `l.186 PROVISIONAL -> S2-ahorro-quincenal=ABIERTO, S2-flujo-quincenal=ABIERTO`
    - `l.190 PROVISIONAL -> S2-metas=RESUELTO, S2-metas-vencida=ABIERTO`
    - En D-6, `S2-metas` aparece como RESUELTO y lo PROVISIONAL es `S2-metas-vencida`.
    - Las otras líneas con la palabra son la bitácora (l. 3 y 6), la solicitud a C1 (l. 476) y el encabezado de la lista de bloqueos (l. 477).
  - Ahorro planificado:
    - `«ahorro planificado FIRME»: 0`
    - `D-2 cita S2-ahorro-quincenal: True · PROVISIONAL: True · «—»: True`
    - `§15 Abiertos nombra S2-ahorro-quincenal: True · Resueltos con «ahorro planificado de D-2»: False`
  - Avance de D-6:
    - `D-6 progress_pct: 1 (solo como referencia a M-07) · progress_tenths_pct: 2 · un decimal: True`
    - `progress_tenths_pct en presupuesto-b.yaml (SavingsGoal de listSavingsGoals): True`
    - `planned_savings_cents en presupuesto-a.yaml: True`
  - U3-02: `operaciones en extractos: 212 · citadas en U3: 26 · inexistentes: []` y `métricas citadas: 16 · inexistentes en C1-datos: []`.
- Alcance sin tocar: md5 de `U4-ajustes-a.md` c926252b…, `U4-ajustes-b.md` 52c5dc53…, `U4-ajustes.md` 18135d33… y `DESIGN.md` 81405a85…, iguales antes y después.
- Global: `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → **rc=0** (22.1 s). Salidas:
  - «U4-ajustes.md termina con parte a + parte b sin cambios: True · secciones P-26: 14 de 14»
  - «RESULTADO: los 21 archivos coinciden byte a byte con el generador»
  - «RESUMEN TODO OK»
  - «s2 · 32 fixtures, 32 coinciden, 0 diferencias»; S3 a S5: «oráculo contra motor» con todas las hojas coincidiendo
  - «== WF-1: verificación completa OK…» y «== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»
  - Bitácora: `/tmp/correccion-U4-r2/verify-all.log`

## Supuestos
1. **D-6 queda dentro de «Metas»** (M-2), por M-16 y por el «Qué cambia» de `S2-metas` (Hecho 2). No es una decisión nueva: es la lectura literal de Q2-17 C. Si Axel quisiera el entero en todo Inicio, D-6 se vuelve a `progress_pct` y la solicitud a C1-datos M-16 se retira. En ese caso se abriría `U3-avance-detalle` con A (un decimal) y B (entero).
2. **«—» sin texto de motivo en D-2.** La revisión pidió solo «—». No inventé un texto de motivo nuevo; si hace falta, se fija con DESIGN o en WF-3.
3. **D-2 dice que el saldo llega sin restar el ahorro.** Es HECHO del motor (`S2-interfaz.md` 1.3 §3.11 y el «Cómo seguí» del bloqueo), no una respuesta a la pregunta 3.

## Bloqueos
- No abrí bloqueos nuevos. Siguen ABIERTOS los que ya estaban: `S2-ahorro-quincenal` (ahora citado en D-2 y en §15), `S2-flujo-quincenal` (Q2-12) y `S2-metas-vencida`.

## Pendientes
- Los bajos de la r2 quedaron fuera del encargo, para la ronda 3 o el integrador:
  - B-1: falta una coma en `U4-ajustes-a.md` l. 398.
  - B-2: `DESIGN.md` l. 933 todavía cita como PROVISIONAL `S1-parametros-avisos`, que ya está RESUELTO.
  - B-3: la sección `ocr.timeout_seconds` de `U4-ajustes-b.md` l. 203 y 387 no tiene marca PROVISIONAL ni una pregunta para la puerta 4.
- Si se corrige B-1 en la parte a, hay que regenerar el unido.
- Línea «Fuentes que rigen» de U3: sigue diciendo «contrato C1 v1.2». No estaba en los hallazgos y no la toqué.

## Propuestas para Axel
- Contestar `S2-ahorro-quincenal` por id. Son tres preguntas, A o B cada una, sin recomendación. Mientras no se contesten, D-2 muestra «—» en el ahorro planificado.

## Consumo aproximado
- 27 acciones de herramienta.
- KB leídos: unos 75, entre la revisión r2, partes de U3, el bloqueo, S2-interfaz §3.13, extractos C1 y C1-datos M-07 y M-16.
- KB escritos: unos 3 netos en U3 (el script reescribe el archivo de 83 KB) y unos 6 de este reporte, más 7 KB de scripts y bitácora en `/tmp/correccion-U4-r2/`, fuera del repo.
