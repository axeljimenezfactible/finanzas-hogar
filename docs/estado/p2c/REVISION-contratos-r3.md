# P2c · REVISION-contratos-r3 · 11-oct-2026, noche

Revisor independiente, ronda 3 de 3 (`revisor-contratos-r3`). No escribí nada de lo revisado y no corregí nada.
Objeto: `S2-interfaz.md` 1.3, `S3-interfaz.md` 1.5, `S4-interfaz.md` 1.5 y `S5-interfaz.md` 1.5 (secciones cambiadas según `git diff HEAD`); `S2-ventanas.md` pregunta 4; los bloqueos nuevos de la cadena de contratos; los reportes `contratos-puerta2.md`, `correccion-contratos-r1.md` y `correccion-contratos-r2.md`.
Contra qué: `docs/criterios/S2.md` a `S5.md`; PUERTA-2 §2, §3 y §8; decisiones §B a §E y §K (K46); `RESPUESTAS.md`.
Volví a correr cada verificación con mis scripts en `datos/tmp/revisor-contratos-r3/` (`r3.py`, `hunks.py`, `hk.py` y `code3.py`). Al terminar los borré.

## Veredicto
**NO APROBADO.** Queda 1 hallazgo alto, 1 medio y 5 bajos.
- **A-1 (alta).** `make verify-all` sigue roto en `verify-c1`. Es el A-1 de la ronda 2 y está fuera del alcance de los contratos.
- **M-1 (media).** El registro de bloqueos sigue sin los 4 que abrieron las correcciones de P2c. Es el M-1 de la ronda 2. La corrección está preparada en `correccion-contratos-r2.md`, pero nadie la aplicó.
- **Cerrado: M-2 de la ronda 2.** Ya está la lista para WF-3, con los 7 comentarios, `card.go` l. 89, las pruebas, las constantes y una receta nueva. La comprobé renglón por renglón.
- **Los contratos están igual que en la ronda 2.** El sha256 de los 9 archivos del alcance coincide con `datos/tmp/correccion-contratos-r2/sha-antes.txt` (`sha256sum -c`: 9 de 9 OK).
- **Lo que no se movió.** Ninguna cifra canon cambió y ninguna pregunta «sin recomendación» quedó decidida. Cada marca PROVISIONAL activa cita un bloqueo ABIERTO.

## Verificaciones (evidencia propia)
| # | Verificación | Resultado | Evidencia |
|---|---|---|---|
| a | Q2-01 a Q2-08 aplicadas como dice su corrección sugerida | CUMPLE | **Q2-01:** S2 l. 151 y l. 155 (regla 4, «contiguas», PROVISIONAL, `S2-ventanas` pregunta 4, caso 1); l. 177 (caso 2, PROVISIONAL, «no está decidido»); l. 353 (A + C y `[d, r]` en `s2-ventanas-sin-cubrir`). `S2-ventanas.md` «Pregunta 4» trae 1A–1D, 2A–2C, las cifras 50000, 97375 y 292125, e impacto. **Q2-02:** S5 l. 10 y l. 125–126 (idéntico = misma presentación en **todas** las canastas). **Q2-03:** S5 l. 131 (paso 8, PROVISIONAL, `S5-carritos-granel`). **Q2-04:** S5 l. 186–191 (gramática del precio, «1,021.00» = 102100, lista cerrada de motivos, cadena vacía); la liga normalizada es PROVISIONAL (`S5-profeco-filtro` pregunta 2 = Q2-57, sin recomendación). **Q2-05 y Q2-08:** S5 l. 166–167 y l. 291–292. **Q2-06:** S4 l. 10, l. 254–255 y l. 388 (`choice_store`, `choice_presentation`). **Q2-07:** S4 l. 185, PROVISIONAL, `S4-presentacion-respaldo`; los pasos (b) y (c) coinciden con `recetas/catalog.go` l. 66–77 (`fallbackPresentation`: la típica y, si no hay, la primera del catálogo). **Q2-08:** S4 l. 408 |
| b | Ninguna «sin recomendación» decidida | CUMPLE | Las que quedan quedan PROVISIONALES con la opción en uso y su bloqueo ABIERTO:<br>– Q2-01 caso 2 → 2A en uso (S2 l. 177).<br>– Variante de Q2-07 → regla actual del motor (S4 l. 185).<br>– Q2-12 («C o A») → `null` PROVISIONAL, como antes (palabra a palabra, S2 l. 258).<br>– Q2-48 → S5 l. 74. Q2-50 a Q2-52 → S5 l. 130. Q2-56 a Q2-58 → S5 l. 194.<br>– Q2-90 → S4 l. 293: valores `propuesto` PROVISIONALES, `S4-empates-sin-confirmar`.<br>De Q2-85 a Q2-94, los contratos solo mencionan Q2-90 (`grep`: S4 l. 14 y l. 293). Ver B-4 |
| c | Cada PROVISIONAL cita un bloqueo abierto | CUMPLE en los contratos; el registro queda en M-1 | `r3.py c` toma el estado de cada bloqueo de la última línea «Estado del bloqueo (puerta 2)» o «- **Estado:**». Resultado: 79 archivos, 62 RESUELTO, 14 ABIERTO y 3 sin esa línea (`criterios-nxx` y los 2 de la Pista 3).<br>Renglones con PROVISIONAL: S2 10, S3 18, S4 8, S5 15.<br>Los que no citan ningún bloqueo ABIERTO son de dos clases, ninguno una marca activa:<br>– definiciones o encabezados: S2 l. 7 y 458; S3 l. 612, 658 y 884; S4 l. 7; S5 l. 7;<br>– bitácora: S3 l. 893 y 896; S4 l. 20, 26 y 31; S5 l. 22, 28 y 34. Ver B-5.<br>Citas a bloqueos ABIERTOS: `S3-formula-tarjeta-pdf` 21, `S3-iva-lectura` 15, `S4-empates-sin-confirmar` 4, `S4-presentacion-respaldo` 3, `S5-carritos-granel` 3, `S2-ventanas` 7, `S2-ahorro-quincenal` 3 y `S2-metas-vencida` 3 |
| d | Las cifras canon no cambian | CUMPLE | `r3.py d` compara `HEAD` con el árbol de trabajo de S2 a S5. **No desaparece ninguna cifra.** Las nuevas son de ejemplos o de versión: S2 100000, 199.00 y 4,000.00; S3 00004; S4 1350, 2250, 2300 y 1.5; S5 0.00, 1,021.00, 102100 y 1.5.<br>`git status` no muestra cambios en `oraculo/` ni en `backend/internal/`.<br>En mi `verify-all`, todos los oráculos con 0 diferencias:<br>– S2: 30/30 fixtures, 2320 hojas;<br>– S3: 906/906, canon 437/437;<br>– S4: 26/26, 8451 hojas;<br>– S5: 29/29, 747 hojas.<br>`model.json` sha256 `1a57066689dd…` intacto |
| e | «Cambios que deben implementar los motores» cubre lo que cambió | CUMPLE (B-1) | **Contratos.** `hunks.py` y `hk.py` (diff por palabra) recorren los 101 bloques del diff (S2 21, S3 34, S4 23, S5 23).<br>– Cada cambio de comportamiento tiene su renglón en las listas de `contratos-puerta2.md`, `correccion-contratos-r1.md` o `correccion-contratos-r2.md`. Algunos ejemplos: regla 4 de §2.12, `provisional` null y `metas` (S2); capital = monto, `VATIncluded`, §14.1 y pruebas (S3); `choice_*`, respaldo del oráculo y orden (S4); §7.4, granel, PROFECO y orden (S5).<br>– Los cambios que solo pasan de PROVISIONAL a FIRME (S2 §3.8, §3.11 y §3.13; S4 §6.4, §6.5 y §10.1; S5 §6, §7.1, §8 y §14) no cambian el cálculo: son las listas de comentarios.<br>**Código.** `code3.py` encuentra 52 renglones con PROVISIONAL y un bloqueo RESUELTO a ±2 renglones.<br>– Los 40 de motores y oráculos están en alguna de las listas, con archivo y renglón.<br>– Los otros 12 son de `api/c1gen` (C1, fuera del objeto; ver la observación).<br>**Recetas de r2.** La receta 1 da 40 y la 2 da 106, como dice el reporte. Los 106 renglones están en las listas.<br>**Los 7 sitios de M-2 y `card.go` l. 89.** Los imprimí y están en esos renglones con ese texto. Sus destinos coinciden con PUERTA-2 §8: Q2-46 = A, Q2-60 = A, Q2-55 = A, Q2-53 = A y Q2-32 = A; l. 89 pasa a `-pdf` por Q2-29 a Q2-31 = (a).<br>**Pruebas** (`engine_test.go` l. 171 y 183, `comparador_test.go` l. 394, `timeline_test.go` l. 460–461): los destinos coinciden con S3 §14.1 (l. 905–921) |
| f | `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | FALLA (A-1) | **rc=2.**<br>– Pasan: WF-0; S1 «RESUMEN TODO OK»; U1 63 OK; criterios con 178 comprobadas y 0 diferencias; INDICE de criterios con 0 diferencias; WF-1 «verificación completa OK».<br>– K25: motor_cred 93.7 %, motor_fin 89.2 %, recetas 93.2 %, precios 96.1 %.<br>– **Falla `verify-c1`: «gen2: tipo de parámetro desconocido: hora HH:MM»; `make: *** [Makefile:130: verify-c1] Error 1`** |

## Hallazgos
**A-1 (alta) · la verificación global sigue rota.** `make -k verify-all` termina con rc=2. `verify-c1` falla con `gen2: tipo de parámetro desconocido: hora HH:MM`, que provoca el tipo nuevo de `00004_puerta2.sql` y `diccionario-datos.md`: `backend/internal/api/c1gen/gen2_schemas.py` (`_pv`) no lo conoce. `git status` no muestra cambios en `backend/internal/api/c1gen/`. Es el A-1 de la ronda 2 (A-2 de la ronda 1), fuera del alcance de los contratos.
- Ruta: `backend/internal/api/c1gen/gen2_schemas.py`
- Corrección: que un agente con alcance en C1, o el integrador, aplique la solicitud 1 de `S1-puerta2.md`. Es el tipo «hora HH:MM» como string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`. Después, correr `build.py --check` y `make verify-all`.

**M-1 (media) · el registro de bloqueos sigue sin los 4 que abrieron las correcciones de P2c.**
- `S3-formula-tarjeta-pdf`, `S3-iva-lectura`, `S4-empates-sin-confirmar` y `S1-reparto-edicion-gasto` tienen **0 menciones** en `INDICE.md` y 0 en `decisiones.md` (K46).
- El INDICE sigue en 75 / 63 / 12, con la fila P2c en 4 / 0 / 4. En disco hay 79 archivos: 63 resueltos (62 con línea RESUELTO más `criterios-nxx`) y 16 abiertos (14 más los 2 de la Pista 3).
- `S3-formula-tarjeta.md`, `S4-confirmacion-empates.md` y `U4-iva-credito.md` no traen la remisión (`grep Remisión`: 0).
- La corrección está escrita: `correccion-contratos-r2.md`, solicitud 2 y Anexo, más `datos/tmp/correccion-contratos-r2/aplica-registro.py`. No se aplicó.
- Si WF-3 toma el INDICE como la lista de lo abierto, no ve 4 de los 16.
- Ruta: `docs/estado/bloqueos/INDICE.md`
- Corrección: que el integrador aplique la solicitud 2 de `correccion-contratos-r2.md`:
  - el script en ensayo y, si da «11 todos con el conteo esperado», con `--aplicar`; si no, el Anexo a mano;
  - después, comprobar que la fila P2c dé 8 / 0 / 8 y el total, 79 / 63 / 16, y que K46 y las 3 remisiones estén en su lugar.

**B-1 (baja) · `correccion-contratos-r1.md` conserva la receta vieja sin remitir a r2.** La lista y la receta de r2 la sustituyen, pero r1 tiene 0 menciones de `correccion-contratos-r2`. Quien lea solo r1 usa una receta que no encuentra los 7 sitios. Es el resto del M-2 de la ronda 2.
- Corrección: la solicitud 3 de r2. Unir las tres listas en la INTEGRACION de P2c o poner una línea de remisión en r1.

**B-2 (baja) · `S2-ventanas.md` l. 85.** «Qué falta» de Q2-01 sigue diciendo que las opciones «las redacta quien aplique Q2-01» y que «esta línea se sustituye por la pregunta completa». La pregunta ya está desde l. 90. Es el B-1 de la ronda 2.
- Corrección: remitir a «Pregunta 4» (l. 90).

**B-3 (baja) · bajos de S2 de la ronda 1, sin atender.**
- S2 l. 359 dice «27 fixtures», y en `oraculo/fixtures/s2` hay 30 (`ls … | wc -l`).
- La fila de S2 l. 354 está junto a la marca de `S2-flujo-quincenal`.
- `contratos-puerta2.md` l. 57.
- Corrección: la de `REVISION-contratos-r1.md` l. 55–59.

**B-4 (baja) · `S3-iva-lectura` pregunta 2.** «Cómo seguí» (l. 28) no dice que no había opción en uso. `ErrMissingCapture` se eligió por RN-26 porque ese estado no existía en 1.3. Es razonable y cita el bloqueo ABIERTO. Es el B-3 de la ronda 2.
- Corrección: agregar «no hay opción en uso; se tomó la que no da número (RN-26), PROVISIONAL».

**B-5 (baja) · bitácora en presente con bloqueos ya resueltos.**
- S4 l. 20 dice «PROVISIONAL y con una variante sin campo nuevo», pero rige «Sin variante» (l. 446). Es el B-4 de la ronda 2.
- S3 l. 893 (Hallazgo 5) dice «Las tres series quedan PROVISIONALES con el bloqueo … hasta que Axel decida», y l. 885 y l. 191 dicen FIRME por Q2-38 = (a).
- Corrección: agregar «(superado en la puerta 2: …)» a las dos.

**Observación fuera del objeto (no cuenta para el veredicto; para el integrador o el agente de C1).** `gen2_schemas.py` l. 364 y l. 374 dicen todavía que dónde viven `send_time` y `advance_days` (`params_json` o `app_setting`) es PROVISIONAL, con el bloqueo `S1-parametros-avisos`. Ese bloqueo está RESUELTO: Q2-11 = B, todo en `app_setting`. `code3.py` señala también `gen6_c1.py` l. 46 y 200. Conviene alinearlo cuando se aplique A-1, y que lo revise quien revise C1.

## Estado de los hallazgos de la ronda 2
| Ronda 2 | Estado | Evidencia |
|---|---|---|
| A-1 (alta) · `verify-c1` | SIGUE ABIERTO → A-1 | mismo error en mi corrida; `gen2_schemas.py` sin cambios |
| M-1 (media) · registro de 3 (+1) bloqueos | SIGUE ABIERTO → M-1 | 0 menciones en INDICE y en `decisiones.md`; totales 75 / 63 / 12; sin remisiones; la corrección está preparada y no se aplicó |
| M-2 (media) · 7 comentarios fuera de la lista para WF-3 | CERRADO (el resto pasa a B-1) | `correccion-contratos-r2.md` «Lista para WF-3»: A (7 sitios), B (`card.go` l. 89), C (pruebas), D (constantes y generadores) y E (receta). Sitios y destinos comprobados; recetas 40 y 106; 106 de 106 en las listas; 40 de 40 sitios PROVISIONAL de motores y oráculos en las listas |
| B-1 · `S2-ventanas.md` l. 85 | SIGUE → B-2 | archivo con el mismo sha256 |
| B-2 · bajos de S2 de la ronda 1 | SIGUE → B-3 | S2 l. 359 «27 fixtures» contra 30 |
| B-3 · `S3-iva-lectura` pregunta 2 | SIGUE → B-4 | l. 28 sin cambio |
| B-4 · S4 l. 20 | SIGUE → B-5 (y S3 l. 893) | sin cambio |
