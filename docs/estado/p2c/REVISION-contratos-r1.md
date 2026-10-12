# P2c · REVISION-contratos-r1 · revisor independiente (ronda 1 de 3) · 11-oct-2026, noche

Objeto: lo que escribió `contratos-puerta2` (reporte `docs/estado/p2c/contratos-puerta2.md`): `S2-interfaz.md` 1.3, `S4-interfaz.md` 1.5, `S5-interfaz.md` 1.5, marcas de `S3-interfaz.md`, pregunta 4 de `S2-ventanas.md` y los 4 bloqueos nuevos. Contra qué: `docs/criterios/S2.md`, `S3.md`, `S4.md` y `S5.md` (1.2), `PUERTA-2.md` §2 y §3, `RESPUESTAS.md` y `decisiones.md` §B–§E y §K (K46). No escribí nada de lo revisado ni corregí nada. Cada cifra y cada estado de abajo lo saqué con mis comandos (scripts en `datos/tmp/revisor-contratos-r1/`, borrados al final).

## Veredicto

**NO APROBADO: 2 altos, 4 medios y 5 bajos.** Las correcciones de Q2-01 a Q2-08 están bien aplicadas, ninguna pregunta «sin recomendación» se decidió y las cifras canon no cambian. Fallan tres cosas: (c), porque S3, S4 y S5 tienen marcas PROVISIONAL que citan bloqueos ya RESUELTOS; (e), porque la lista de cambios para los motores no incluye las fixtures de S4 y S5 ni Q2-67 en S3; y la verificación global, porque `make verify-all` se cae en `verify-c1`. Esa última falla viene de fuera de los contratos.

## Verificaciones (evidencia propia)

| # | Qué | Resultado | Evidencia |
|---|---|---|---|
| a1 | Q2-01 (M1 de S2-r3) | CUMPLE | §2.12 trae las 4 reglas: O1, A y H1 FIRMES; la contigüidad, PROVISIONAL por `S2-ventanas` pregunta 4. §3.2 escribe el caso 2 como PROVISIONAL 2A. La pregunta 4 trae las opciones 1A–1D y 2A–2C, sin recomendación (l. 90–115). Cifras con el oráculo (`fin_devengo.periodos`): `[P1]` da 97375; `[d, r]` da 194750 y **292125**; semilla sin crédito + 100000 de contado el 2026-10-10 a 2 meses: A + C deja C con devengado **50000** y diferencia 0; A + B + C deja C en **0**; en meses, dic-2026 da **0**. Coinciden con el contrato y con `REVISION-S2-r3.md` l. 41–45 |
| a2 | Q2-02 (H-r3-01) | CUMPLE | §7.2 pasos 2, 3 y 7, §11 (3 RAR, `raw/<corte>-csv/`, `db/qqp.sqlite`), §12 (tickets-README y tickets-sinteticos) y §15.2. Las solicitudes 4 a 6 quedan anotadas para el integrador |
| a3 | Q2-03 (H-r3-02) | CUMPLE | §7.2 paso 8, PROVISIONAL por el bloqueo nuevo `S5-carritos-granel` (ABIERTO). Las cifras 101/103, 135/137/721/736 y 2933/2800 son las de `REVISION-S5-r3.md` l. 52–53. Revisé el motor: `compare.go` l. 309–336 redondea una sola vez con HALF_UP sobre la suma exacta |
| a4 | Q2-04 (H-r3-03) | CUMPLE | La lista cerrada de motivos de §11 coincide con los literales de `backend/internal/precios/profeco/*.go`: «precio vacío», «precio con comas mal puestas», «precio con más de dos decimales», «precio ilegible», «precio no positivo», «fecha vacía/ilegible/futura», «producto vacío», «cadena comercial vacía» y «columnas: %d en vez de %d». La liga normalizada es PROVISIONAL por `S5-profeco-filtro` (ABIERTO) |
| a5 | Q2-05 y Q2-08 (H-r3-04, M-r5-3) | CUMPLE | §7.4 da un solo significado: completas primero, `null` en las incompletas. Agregados el orden de §15.2 de S5 y de §12.2 de S4. Fixtures con una tienda incompleta y diferencia no nula: solo `s5-carritos-presentaciones`, que ya está en la lista de cambios |
| a6 | Q2-06 (M-r5-1) | CUMPLE | `quotes` trae `choice_store` y `choice_presentation`, y §8.1 paso 1 toma la elección exacta. El ejemplo (Bodega 750 g a 2300 incluido) es el de `REVISION-contrato-S4S5-r5.md` l. 87. Lo recalculé: la recomendada es Bodega 500 g (2700/kg); (walmart, null) da 1 kg a 2800; (null, 750 g) da Walmart 2250 < 2300 |
| a7 | Q2-07 (M-r5-2) | CUMPLE | La variante queda sin decidir. Lo PROVISIONAL es la regla del motor, como manda `RESPUESTAS.md`. La confirmé en `catalog.go` l. 66–77 (`fallbackPresentation`) y en `shopping.go` l. 83–137. Ver B-5 |
| b | «Sin recomendación» no decididas | CUMPLE | Q2-01 caso 2, Q2-07 variante, Q2-12, Q2-48, Q2-50–52 y Q2-56–58: todas PROVISIONAL «en uso», con bloqueo ABIERTO. Coincide con K46 |
| c | Cada PROVISIONAL cita un bloqueo abierto | **NO CUMPLE** | Conté por segmento con `prov.py`. **S2**: `S2-ahorro-quincenal`, `S2-flujo-quincenal`, `S2-metas-vencida` y `S2-ventanas`, todos ABIERTOS. **S4**: `S4-asignacion-mercado` 4, `S4-confirmacion-empates` 2 y `S4-reglas-lista` 1, todos RESUELTOS (fuera de la bitácora: l. 292, 341, 345 y 361). **S5**: `S5-reglas-comparativa` 3, `S5-profeco-almacen` 1 y los de S4 (fuera de la bitácora: l. 234 y 243), RESUELTOS. **S3**: `S3-formula-tarjeta` 11, `U4-iva-credito` 8, `S3-serie-sin-intereses` 2 y `S3-origen-cargo-normal` 1, todos RESUELTOS |
| d | Cifras canon | CUMPLE | Ninguna cifra salió del diff. Las cifras nuevas tienen origen: S2, el oráculo y PUERTA-2 Q2-21/24; S4, el r5; S5, H-r3-03. `make verify-wf1-oraculos`: S2 30/30, S3 27, S4 26/26 y S5 29/29, con 0 diferencias. `model.json` y `verify3.json` tienen el sha256 intacto. `verify-criterios`: 261 criterios con 0 diferencias |
| e | «Cambios que deben implementar los motores» contra el diff | **PARCIAL** | S2 cubre todo, incluidas las fixtures `s2-metas`, `s2-liga-provisional` y `s2-proyeccion-7-12-cadencia`. A S3 le falta Q2-67 (A-1). En S4 y S5 faltan 14 fixtures cuyo `provisional` cita bloqueos resueltos y los comentarios PROVISIONAL del código (M-4) |
| g | Verificación global | **FALLA** | `make verify-all … TMP=<propia>` pasa WF-0 (SQL con la 00004, Go, color y OpenAPI) y se cae en `verify-c1`: `gen2: tipo de parámetro desconocido: hora HH:MM` (A-2). Con `-k` corrí `verify-criterios`, `verify-wf1-cobertura` (K25 OK: motor_fin 89.2 %, motor_cred 93.7 %, recetas 93.2 %, precios 96.1 %) y `verify-wf1-oraculos`: todo OK |

## Hallazgos

**A-1 (alta) · S3 no aplica Q2-67 = 1.** `U4-iva-credito` está RESUELTO con la opción 1 (K46). La 00004 agrega `credit.vat_included` (NULL = falta capturar) y RN-26 1.4 ya pide 5 datos. Pero `S3-interfaz.md` sigue leyendo `vatBP` nil como «la tasa ya trae IVA», con la marca PROVISIONAL (`U4-iva-credito`), y deja la lista de K18 en 4 datos (§5.4, §5.5, §5.8, §14 y §14.1: 8 marcas). La lista «Cambios que deben implementar los motores» de S3 no lo menciona. El reporte lo deja fuera porque cae fuera de su rango Q2-28–38. El contrato que recibe WF-3 contradice a Axel y al esquema.
- Ruta: `docs/contratos/motores/S3-interfaz.md`.
- Corrección: aplicar Q2-67 = 1 en S3. Entrada `VATIncluded` (`vat_included`): nil = falta capturar (`ErrMissingCapture` con `vat`), 1 = no suma IVA y 0 = suma `VATOnInterestBP`. `missing_terms` con 5 datos. Quitar `U4-iva-credito` de las listas de §14.1 y de las fixtures. Agregar el cambio a la lista de S3 y a «Solicitudes a C1» (`Credit.vat_included`). Si no es de este agente, el orquestador se lo asigna a otro antes del commit.

**A-2 (alta) · Verificación global rota (fuera de los contratos).** `verify-c1` falla con `gen2: tipo de parámetro desconocido: hora HH:MM` (`backend/internal/api/c1gen/gen2_schemas.py` l. 244). Lo provoca el tipo nuevo que traen `00004_puerta2.sql` y `diccionario-datos.md` (S1-puerta2).
- Ruta: `backend/db/migrations/00004_puerta2.sql`.
- Corrección (para S1 o el integrador): enseñar a `gen2` el tipo «hora HH:MM» (string con patrón `^([01]\d|2[0-3]):[0-5]\d$`) o usar un tipo que ya conozca. Después, `build.py --check` y `make verify-all`.

**M-1 (media) · S4 y S5: Q2 aprobadas sin aplicar; quedan PROVISIONAL citando bloqueos resueltos.**
- S4: l. 292, §9.1, «Quién confirma y cuándo: PROVISIONAL (`S4-confirmacion-empates`)», cuando Q2-43 = A; el reporte dice que la aplicó. L. 341, §11, «vertical `market`: PROVISIONAL (`S4-asignacion-mercado`)». L. 345, §12, `provisional` con `S4-reglas-lista`. L. 361, la fila `s4-canal-mercados`.
- S5: l. 234, §15, «los que tocan `S5-reglas-comparativa` o `S5-profeco-almacen` llevan `provisional`». L. 243, la fila `s5-cotizacion-recomendada`.
- Corrección: «FIRME (puerta 2, Q2-nn)». En §12 y §15, que solo lleven `provisional` las fixtures de bloqueos abiertos.

**M-2 (media) · S3 §5.8, tabla de despacho (l. 309, fila `bank_card`), con marcas viejas.** Dice «capital = monto (PROVISIONAL `S3-serie-sin-intereses`)» y «`quoted` PROVISIONAL (`S3-origen-cargo-normal`)». Q2-37 y Q2-38 son FIRMES, y el propio contrato lo dice en §5.3 y §5.5. La l. 594 todavía llama PROVISIONALES a los casos `s3-comparador-*`.
- Corrección: «FIRME (puerta 2, Q2-37 / Q2-38)» y quitar «PROVISIONALES» de la l. 594.

**M-3 (media) · `S3-formula-tarjeta` dice RESUELTO, pero S3 lo cita en 11 marcas PROVISIONAL.** Q2-29 a Q2-31 = (a) deja la fórmula esperando los PDF de Axel. Mientras no lleguen, sigue en uso la opción (b). La regla de conteo de `INDICE.md` l. 30 lo dice: «una decisión tomada pero sin resultado ya no cuenta como resuelto». Choca el bloqueo (más INDICE y K46, que lo cuentan resuelto) contra el contrato.
- Ruta: `docs/estado/bloqueos/S3-formula-tarjeta.md`.
- Corrección (bloqueos o integrador): preguntas 1 a 3 en «ABIERTO: respuesta (a), en espera de los estados de cuenta» y estado del bloqueo ABIERTO; la pregunta 4 sigue RESUELTA. Ajustar INDICE y K46 a 7 abiertos. Otra salida: un bloqueo nuevo que S3 cite.

**M-4 (media) · Faltan fixtures y código de S4 y S5 en la lista de cambios.** Hay 14 fixtures cuyo `provisional` cita bloqueos RESUELTOS:
- S4: `s4-canal-mercados`, `s4-kcal-recetas-autor`, `s4-kcal-sin-dato`, `s4-lista-reposicion-ritmo` y `s4-salud-imc-franja`.
- S5: `s5-canasta-regla-unica`, `s5-cotizacion-recomendada`, `s5-mejor-tienda-empate-orden`, `s5-mejor-tienda-empate-recomendados`, `s5-mejor-tienda-empate-viejos`, `s5-mejor-tienda-incompleta`, `s5-membresia-anual`, `s5-promocion-vigencia` y `s5-profeco-muestra`.

Tampoco se listan los comentarios PROVISIONAL del código, por ejemplo `shopping.go` (reposición, `S4-reglas-lista`). Sin eso, WF-3 tendría que adivinar.
- Corrección: agregar a la lista: `provisional` en null en esas 14 o, en `s5-profeco-muestra`, re-citar `S5-profeco-filtro` si depende del filtro; quitar los comentarios PROVISIONAL de S4 y S5 resueltos.

**B-1 (baja).** S2 §5, l. 355: la fila de `s2-liga-provisional` dice «(PROVISIONAL, `S2-flujo-quincenal`)», pero §3.11 y §5.2 la dan FIRME con `provisional` null. Corrección: «(FIRME, Q2-14 = A; `provisional` null)».

**B-2 (baja).** S2 §5, l. 359: dice «Total: **27 fixtures**», pero hay 30 (`oraculo/fixtures/s2`; `verify-wf1-oraculos` dice «30 fixtures»), y el encabezado 1.3 dice «§5 y §8 al día». Corrección: 30, con la lista.

**B-3 (baja).** Reporte, «Cambios» de S4: «el motor las agrega en orden del plan» (`items[].recipes[]`) es falso: `shopping.go` hace `sort.Strings(n.recipes)`. Corrección: corregir la nota.

**B-4 (baja).** El supuesto 3 elige M-r5-3 (Q2-08) sobre la solicitud 3 de `correccion-S5-r2` (Q2-02). Las dos venían aprobadas en bloque y chocan. Corrección: que entre en «Propuestas para Axel» como posible excepción por id.

**B-5 (baja).** S4 §6.2 paso 5 (a): el motor solo mira los ingredientes no opcionales de las entradas que no son «come fuera» ni quedan cubiertas por sobrantes (`shopping.go` l. 83–99). Corrección: agregar ese matiz.

## Lo que está bien (sin hallazgo)
- Q2-13 a Q2-27 en S2. Q2-12 sigue abierta porque su recomendada nombra dos opciones (K46 dice «— (C o A)»).
- Los bloqueos nuevos `S2-ahorro-quincenal` y `S2-metas-vencida` traen qué falta, opciones, impacto y «sin recomendación». Revisé el HECHO de metas vencidas en `projection.go` l. 168–198 y en `fin_metas.py` l. 20–36: el motor da todo el restante en el mes de hoy y el oráculo da `plan` vacío y `null`. Los cambios de `months_left`, `open` y `achieved` sí están en la lista.
- `S4-presentacion-respaldo` y `S5-carritos-granel` están ABIERTOS y completos.
- Q2-39 a Q2-46 en §6.2, §6.4, §6.5, §10.1 y §12.2 de S4, salvo M-1.
- Q2-47, Q2-49, Q2-53 a Q2-55, Q2-59 a Q2-63 y Q2-76 en S5, salvo M-1.

## Estado de los hallazgos anteriores (WF-1)
- M1 de `REVISION-S2-r3`: aplicado.
- H-r3-01 a H-r3-05 de `REVISION-S5-r3`: aplicados. H-r3-00 no es de este agente (Q2-64).
- M-r5-1 a M-r5-3 de `REVISION-contrato-S4S5-r5`: aplicados, salvo las marcas de M-1 y M-4.

Es la primera ronda de esta revisión, así que no hay hallazgos P2c anteriores.
