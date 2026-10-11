# Revisión de los criterios 1.1 · ronda 2 de 3 (P1) · 11-oct-2026

Revisor independiente: `revisor-criterios-r2`. Objeto: `docs/criterios/` (18 archivos e `INDICE.md`) después de `criterios-S-1.1`, `criterios-UC-1.1` y `correccion-criterios-r1`. Base de evaluación: `INDICE.md` §1, `decisiones.md` §K (K1 a K34) y el hallazgo R2-07 de `docs/estado/wf0/REVISION-r3.md`. No corregí nada: todo lo reejecuté con scripts propios en `datos/tmp/rev-crit-r2/` (borrados al terminar) y con `docs/criterios/tools/verif_final.sh`.

## 1. Veredicto
**NO APROBADO.** Hay 0 hallazgos altos, **1 medio** y 6 bajos nuevos. Siguen abiertos 6 bajos de la ronda 1 y uno queda a medias.
- Los 6 medios de la ronda 1 (H-01 a H-06) están **cerrados** (§5).
- El medio nuevo es **R2-M01**: S2-06 (alta) no aplica K8 (FIRME). Describe la quincena por sus días de pago y dice que en un mes de 31 días «la quincena de fin de mes es el 30, no el 31». Es el mismo defecto que H-05, que solo se corrigió en S1-09 y U3-05.
- La estructura, los conteos, el INDICE, las versiones, las rutas y las cifras cumplen.

## 2. Comprobaciones pedidas (a)–(f)
| | Comprobación | Cómo la reejecuté | Resultado |
|---|---|---|---|
| a | 5 secciones; cada criterio con Verifica, Evidencia y línea en su §5 | `estructura.py` (propio). Revisa las secciones 1 a 5 en orden y que **cada** criterio tenga exactamente 1 Verifica y 1 Evidencia. Expande rangos en §5 y busca huecos o duplicados de ID. Corrí también `bash docs/criterios/tools/verif_final.sh` | 18 de 18 archivos con las 5 secciones en orden. **257 criterios (197 alta / 60 media / 0 baja)**, todos con 1 Verifica y 1 Evidencia. 0 sin línea en su §5, 0 huecos, 0 duplicados, 0 cabeceras fuera de §2. Herramientas: `cuenta.py` 134 (105/29), incompletos ninguno; `traza.py` 257, 0 sin traza; `refs.py` 257/257, 0 inexistentes; `canon.py` 178 cifras, 0 diferencias, sha256 `1a5706…bfd4`; `s109_sql.sh` con los límites de K8 bien (q2 de oct `2026-10-16`→`31`, feb `16`→`28`, 0 quincenas fuera, semana `2026-10-12`→`18`); `s109_equiv.py` 12 de 12; `citas.py` 275 citas, 5 sin fuente (H-07) |
| b | Cada K y R2-07 aplicado, con línea | `kmat.py` (propio): matriz K → archivo:línea(criterio). `kw.py` (propio): 64 comprobaciones de palabras clave por criterio. Leí completos 40 criterios | 32 de las 33 decisiones K que piden criterio están aplicadas (§3.2). **K8 falta en S2-06** (R2-M01). R2-07 está aplicado en sus 6 partes (§3.3) |
| c | Sin rutas `/tmp` | expresión `(?<!datos)/tmp\b` sobre los 18 archivos, el INDICE y `tools/` | 0 rutas en comandos. Las 11 apariciones son la viñeta «Toda ruta temporal pasa de `/tmp` a …datos/tmp» y el patrón grep de `verif_final.sh:16` |
| d | INDICE coherente con el conteo | `estructura.py` lee la tabla §2 del INDICE | Las 18 filas y el total (257 / 197 / 60) coinciden. S y WF4: **134 (105/29)**. U, C y P3: **123 (92/31)**: C1 21 (19/2), C2 12 (8/4), C3 11 (9/2), P3 13 (10/3), U1 16 (10/6), U2 14 (10/4), U3 12 (8/4), U4 11 (8/3), U5 13 (10/3). Queda texto viejo en INDICE.md:7 («Hoy los 248») y en §5, §6, §7, §8 y §9 (R2-B01, para el integrador) |
| e | Versión 1.1 con lista de cambios en los 18 | `estructura.py` revisa el texto antes de «## 1.» | 18 de 18 con «Versión 1.1 (11-oct-2026…)» y «Cambios de la versión 1.1» (2 a 10 viñetas). U3.md:6 contradice U3.md:7 (R2-B02) |
| f | Ningún criterio reabre §A–§K | lectura de S2-08, S5-04, U2-02, S4-07, S3-12, S2-09, S1-16, U1-16 y los bloqueos nuevos | Ninguno cambia un valor decidido. S2-08 marca PROVISIONAL el semanal con bloqueo (protocolo §2: el redondeo no está en K20). S5-04 usa la lectura literal de K22 (ver H-08). El 135° de U2-02 lo permite decisiones I. El «mismo peso» de S5-04 va marcado PROPUESTA, igual que en K22. S2-06 no reabre K8, pero lo contradice en el texto (R2-M01) |

## 3. Tabla de criterios de esta revisión (estado y evidencia)
### 3.1 Reglas de INDICE §1
| Regla | Estado | Evidencia |
|---|---|---|
| Un archivo por ítem, con las secciones 1–5 | Cumple | `estructura.py`: 18/18 en orden |
| Criterios `<ID>-nn` con enunciado, Verifica, Evidencia y severidad | Cumple | 257, cada uno con 1 Verifica y 1 Evidencia |
| Severidades alta y media | Cumple | 197 / 60 / 0 |
| Evidencia con el ID al inicio de la línea | Cumple | las Evidencias nuevas (U1-15, U1-16, U5-13, C1-16 a C1-21) piden «línea `ID`» |
| Lo que depende de Axel no se declara cumplido | Cumple | S5-06 «sin medir» sin las 5 fotos (S5.md:39); S7-07 «pendiente de acción» (S7.md:42); WF4-05 «no ejecutado» con motivo (WF4-operacion.md:35) |
| Rutas absolutas, `python3 -I`, sin `jq` | Cumple con observación | 0 `jq`, 0 `/tmp`; C1-02 llama a los scripts con ruta relativa (R2-B06) |
| Identificadores `P-xx`, `N-xxx`, `F-x` | Cumple | 0 residuos de N-xx; U1-16 usa P-13, N-150 y N-165 |
| Cifras del oráculo exactas | Cumple | `canon.py` 178/0; recalculé las cifras nuevas de S2-07, S2-08 y S2-09 (§6) |
| Trazabilidad en §5 | Cumple con observación | 257/257; 5 citas juntan dos fuentes (H-07); el 135° tiene la fuente mal atribuida en U2 §5 (R2-B03) |

### 3.2 Decisiones §K
| K | Dónde se aplica (archivo:línea, criterio) | Estado |
|---|---|---|
| K1 | C1.md:77 (C1-15); S1-16 por la regla general (S1.md:81-82) | Cumple; C1 §5 sin línea K1 (H-11) |
| K2 | C1.md:81 (C1-16); U1.md:51 (U1-09); U4.md:24, :28 (U4-03, U4-04) | Cumple |
| K3 | C1.md:85 (C1-17); U1.md:51 | Cumple |
| K4 | C1.md:89 (C1-18); S7.md:42, :50 (S7-07, S7-09); U1.md:51; U4.md:24; WF4-operacion.md:35 (WF4-05) | Cumple |
| K5 | C1.md:41 (C1-06), :77 | Cumple |
| K6 | C1.md:65, :77, :93 (C1-12, C1-15, C1-19); S1.md:81 (S1-16); S6.md:34, :38 (S6-04, S6-05) | Cumple |
| K7 | C1.md:77; S1.md:81; S4.md:30, :42 (S4-04, S4-07); S5.md:23, :51 (S5-02, S5-09) | Cumple |
| K8 | S1.md:53-55 (S1-09); C1.md:41, :45 (C1-06, C1-07); U3.md:31 (U3-05); U1.md:79 (U1-16) | **No cumple en S2.md:39-40 (S2-06)**; S1.md:29 (S1-03) no trae la regla de periodos → R2-M01 |
| K9 | C1.md:93 (C1-19); S6.md:34 (S6-04); U1.md:75 (U1-15); U5.md:38, :62 (U5-07, U5-13) | Cumple |
| K10 | U1.md:31 (U1-04); U5.md:46 (U5-09); C3.md:42 (C3-08) | Cumple |
| K11 | C1.md:65 (C1-12); S1.md:73 (S1-14); S6.md:42, :46 (S6-06, S6-07) | Cumple |
| K12 | C1.md:45, :97 (C1-07, C1-20); S1.md:57 (S1-10); U4.md:24, :40 (U4-03, U4-07) | Cumple |
| K13 | S3.md:67 (S3-13); U1.md:51; U4.md:24, :28 | Cumple |
| K14 | S2.md:35 (S2-05) | Cumple |
| K15 | S3.md:35 (S3-05) | Cumple |
| K16 | S3.md:63 (S3-12); U1.md:51; U4.md:24, :28 | Cumple |
| K17 | S3.md:67, :71 (S3-13, S3-14) | Cumple |
| K18 | S3.md:47, :55 (S3-08, S3-10); U1.md:51; U4.md:24, :28 | Cumple |
| K19 | S2.md:43 (S2-07) | Cumple; 7 cierres recalculados (§6) |
| K20 | S2.md:47 (S2-08) | Cumple; semanal PROVISIONAL con bloqueo `criterios-promedios-redondeo` (formato del protocolo: 5 de 5 campos) |
| K21 | S2.md:51 (S2-09); U1.md:51; U4.md:24, :28 | Cumple; las fechas de la despensa 4 coinciden con `model.json` |
| K22 | S5.md:31 (S5-04); S1.md:29 (S1-03) | Cumple con observación (H-08) |
| K23 | S4.md:42 (S4-07); U1.md:51; U4.md:24, :28 | Cumple; S1-11 no nombra la columna (R2-B04) |
| K24 | S5.md:39 (S5-06) | Cumple |
| K25 | S4.md:74 (S4-15); S5.md:67 (S5-13); S6.md:26 (S6-02); S8.md:38 (S8-06) | Cumple |
| K26 | U2.md:22, :38, :46 (U2-02, U2-06, U2-08); C2.md:7 | Cumple |
| K27 | U2.md:50 (U2-09) | Cumple |
| K28 | S6.md:22 (S6-01) | Cumple |
| K29 | S6.md:74 (S6-14); WF4-operacion.md:27 (WF4-03) | Cumple |
| K30 | S7.md:38, :42 (S7-06, S7-07); WF4-operacion.md:19 (WF4-01, con `getfacl` y `namei`) | Cumple |
| K31 | S1.md:81; S3.md:63; S6.md:54, :58; S7.md:30, :46, :54; S8.md:22; WF4-operacion.md:31, :47 | Cumple: PROVISIONAL con la propuesta escrita hasta la puerta 2 o la 4 |
| K32 | C1.md:101 (C1-21) | Cumple (el tope de 80 KB es convención: H-09) |
| K33 | proceso | No aplica |
| K34 | C1.md:69 (C1-13 pide C1 v1.1) | Cumple |

### 3.3 R2-07 (REVISION-r3)
| Parte | Evidencia | Estado |
|---|---|---|
| U2-02: admitir 135° como propuesta | U2.md:22, :46, :91 | Cumple (fuente mal atribuida en §5: R2-B03) |
| U1-05: ≥ 6 flujos con los 4 obligatorios | U1.md:35, :86, :99; C3.md:70 | Cumple |
| C1-01: el comando `uvx` va primero | C1.md:22 | Cumple |
| S1-03: ejemplos «ilustrativos» rotulados | S1.md:29-30 | Cumple |
| `/tmp` → `datos/tmp` | 0 rutas (comprobación c) | Cumple |
| Citar `verify_c1.py --db` y `build.py --check` | C1.md:26, :30; S6.md:31 | Cumple (C1-02 con ruta relativa: R2-B06) |
| Paráfrasis entre comillas | `citas.py`: 5 sin fuente | Abierto como H-07 (baja) |

## 4. Hallazgos de esta ronda
| ID | Sev. | Criterio | Ruta | Descripción | Corrección sugerida |
|---|---|---|---|---|---|
| R2-M01 | media | §K K8 (FIRME) / coherencia S1↔S2 | `docs/criterios/S2.md:39-40` (S2-06, alta); secundario `docs/criterios/S1.md:29` (S1-03) | S2-06 no cambió en esta ronda (0 líneas en el diff) y S2 no cita K8 en ninguna parte. Dice «quincena (15 y 30, o último día del mes) como subperiodo» y, en el caso de un mes de 31 días, «la quincena de fin de mes es el 30, no el 31». Confunde el día de pago con el periodo. Leído como periodo, deja el día 31 fuera de la quincena 2, contra K8 («el 31 cae en la segunda») y contra S1-09 1.1 («el día de pago … no es el fin del periodo»). La Verifica solo compara con el oráculo, que no tiene eventos el día 31: un motor con la quincena 16→30 pasaría. S1-03 lista «quincena 15 y 30 con ajuste a fin de mes» como regla mínima y no trae la regla de periodos de K8 | S2-06: «quincena como subperiodo con los límites de K8 (del 1 al 15 y del 16 al último día; el 31 en la segunda); el día de pago (15 y 30, o el último día; el 28 en febrero de 2027) es otro dato». Caso de 31 días: «el pago de fin de mes cae el 30 (oráculo) y un gasto de fixture del 31-oct-2026 cuenta en la quincena 2 de octubre». Ese caso entra en la Verifica, con traza de K8 en §5 y viñeta en la lista de cambios. S1-03: «día de pago 15 y 30 …» y una regla de periodos de K8 con ejemplo ilustrativo (31-oct-2026 en la quincena 2; semana del 14-oct-2026 = del 12 al 18 de octubre) |
| R2-B01 | baja | Coherencia del INDICE (integrador; sigue H-12) | `docs/criterios/INDICE.md:7`, :53-66, :67-86, :87-100, :105, :107 | La §1 dice «Hoy los 248 criterios» y el total es 257. Las notas 8 a 11 de la §5 ya las resolvieron K26, K27, K30 y K28. La §6 no marca los bloqueos que resolvió §K ni trae `criterios-promedios-redondeo` y `criterios-frescura-extremos`. La §7 no tiene fila de §K. La §8 conserva «≥ 2 meses», cuando S6-07 ya pide 90 días (K11). La §9 no registra el flujo P1 | Reconciliar al integrar: 257, marcar RESUELTO con su K, agregar los 2 bloqueos nuevos, fila K en §7, quitar el supuesto de §8 y una sección de cambios P1 |
| R2-B02 | baja | INDICE §1 (lista de cambios) | `docs/criterios/U3.md:6-7` | La viñeta «Sin cambios en los criterios U3-01 a U3-12» contradice la siguiente, que cambia U3-05 | «Sin cambios en U3-01 a U3-04 y U3-06 a U3-12» |
| R2-B03 | baja | INDICE §1 (trazabilidad) | `docs/criterios/U2.md:100` (§5) | La §5 dice que admitir 135° para el ingreso es «convención propia del directorio». Decisiones I sí lo dice: «ingreso 150° (comparte familia con éxito: se distingue por icono y contexto, o U2 propone 135°)» (`decisiones.md:63`) | Citar decisiones I textual → U2-02, U2-08 |
| R2-B04 | baja | §K K3, K9, K23 (esquema 1.2) y K4, K13, K16, K18, K21 | `docs/criterios/S1.md:61` (S1-11), `:69` (S1-13) | La lista de 12 puntos de S1-11 y los ejemplos de S1-13 no nombran lo que §K agrega al esquema 1.2: objetivo (K23), borrador del ticket (K3), claves de idempotencia de `expense` y `weight_entry` (K9), tope por persona, restricciones, parámetros por producto, interruptor de los meses 7–12 y chats candidatos. Hoy los cubren C1-17, C1-19, C1-20, S4-07 y U4-03, así que no queda ninguna K sin verificar | Agregarlos a S1-11 y S1-13 (o a S1-02) para que el revisor del esquema 1.2 los compruebe con S1 |
| R2-B05 | baja | INDICE §1 (medible) | `docs/criterios/U1.md:51` (U1-09) | Pide «≥ 28 capacidades», pero su lista separada por «;» tiene 27. Viene de antes: HEAD tenía 19 con «≥ 20» | «Las 27 capacidades de la lista», o separar «reporte mensual» y «vista web mensual» |
| R2-B06 | baja | INDICE §1 (rutas absolutas) / R2-07 | `docs/criterios/C1.md:26` (C1-02), `:30` (C1-03) | Llama a `backend/internal/api/c1gen/verify_c1.py` y a `build.py --check` con ruta relativa. S6-03 (S6.md:31) ya usa la absoluta | Ruta absoluta o `cd /home/deepsec/finanzas-hogar &&` antes del comando |

## 5. Estado de los hallazgos anteriores (`REVISION-criterios-r1.md`)
| Hallazgo | Sev. | Estado | Evidencia reejecutada |
|---|---|---|---|
| H-01 U2-06 con gasto 25° | media | **Cerrado** | U2.md:38: «El gasto va en terracota 40–45° (K26) … no un valor de partida de 25°». U2.md:100 traza K26 a U2-06 |
| H-02 U4 sin §K | media | **Cerrado** | U4.md:24 (U4-03) trae K2, K4, K12, K13, K16, K18, K21 y K23, ≈ 53 elementos. U4.md:28 (U4-04) trae los valores por defecto (alcance, interruptor apagado, sin valor por defecto en restricciones, tarjeta, tope y objetivo). U4.md:40 (U4-07) trae K12. Trazas en U4.md:86-93 y lista de cambios en U4.md:5-9 |
| H-03 U1-09 sin §K | media | **Cerrado** | U1.md:51 suma K2, K3, K4, K13, K16, K18, K21 y K23, con su `P-xx` y sus `N-xxx`. Trazas en U1.md:113-120. Queda la nota del umbral (R2-B05) |
| H-04 WF-3 sin cola | media | **Cerrado** | U5.md:62-64: U5-13 (alta), Playwright con la red cortada en 7 pasos, `Idempotency-Key` y `captured_at`, 0 duplicados. U5.md:38-39: U5-07 con «Sin conexión». Traza en U5.md:97. INDICE U5 13 (10/3) |
| H-05 K8 en S1-09 y U3-05 | media | **Cerrado** en S1-09 y U3-05 | S1.md:53-55 trae los límites y la semana separados del día de pago. U3.md:31 igual. `s109_sql.sh` los comprueba (§2 a). **El mismo defecto sigue en S2-06 → R2-M01** |
| H-06 S8-06 sin K25 | media | **Cerrado** | S8.md:38: «(1) cobertura de `go test` ≥ 85 % en cada uno de los 5 paquetes de K25 … los demás … sin umbral» |
| H-07 citas con dos fuentes | baja | Abierto | `citas.py`: 5 sin fuente (C1, C2 ×2, U2, U5) |
| H-08 S5-04, día 30 | baja | Abierto | S5.md:31-32 sin cambios. La Verifica no fija peso(30) |
| H-09 C1-21, 80 KB | baja | Abierto | C1.md:144. C1.md no cambió desde la r1 (11:31 contra 11:53) |
| H-10 C1-16 `default` | baja | Abierto | C1.md:81-82, :136 sin marca de convención |
| H-11 K1 en §5 de C1 | baja | Abierto | `kmat.py`: K1 solo en C1.md:12 y :77 |
| H-12 INDICE | baja | Parcial | Filas y total ya en 257. Lo demás pasa a R2-B01 |
| H-13 número de K en §5 de S | baja | Abierto | Las §5 de S2 a S7 citan «decisiones K «…»» sin número (por ejemplo, la de S2) |

## 6. Cifras recalculadas por el revisor (`oraculo/model.json`, sha256 `1a5706…bfd4`)
- **S2-07 (K19).** Diferencia acumulada (Σ `devengado.total` − Σ `efectivo.total` de `periodos[0..5]`): 973.75, 973.75, 1,258.69, 1,258.70, 1,258.69, 1,258.70. Saldo no devengado (total de la despensa menos las mitades devengadas): 1,947.50 / 0.00 / 2,517.37 / 0.00 / 2,517.37 / 0.00, y al 10-abr-2027 5,034.74 (despensa 4, `2027-04-10`). Deuda viva: 2,921.25 / 973.75 / 3,776.06 / 1,258.70 / 3,776.06 / 1,258.70 / **6,293.44**. Es igual a `horizonte.deuda_viva_10_abr`, y `deuda_viva_detalle` es 1,258.70 + 5,034.74. **Coincide** con S2-07.
- **S2-08 (K20).** 973.75 / 31 = 31.41; 1,947.50 / 30 = 64.92; 2,517.38 / 28 = 89.91 (ROUND_HALF_UP). Por 7: 219.88 contra 219.87, y 629.35 (629.345 exacto) contra 629.37. **Coincide**. El efectivo por mes calendario es el de `s109_equiv.py` (12/12).
- **S2-09 (K21).** `despensas[3].pagos`: `2027-04-30`, `2027-05-15`, `2027-05-30`, `2027-06-15`. **Coincide**.
- **S4-07 (K23).** Fórmula de Mifflin-St Jeor: 10×80 + 6.25×175 − 5×30 + 5 = 1,748.75 → ×1.2 = 2,098.50; 10×60 + 6.25×165 − 5×28 − 161 = 1,330.25 → ×1.2 = 1,596.30. Correcto (personas ficticias). **S5-09:** 1,000.00 ÷ 125.00 = 8 (ilustrativo).
- **Semana (K8):** el 14-oct-2026 es miércoles, el 12 lunes y el 18 domingo.

## 7. Notas (no son hallazgos de los criterios)
- `matriz_c102.py`: U1 ya define 203 `N-xxx`. N-316, N-318 y N-319 no tienen operación, y N-316 a N-319 faltan en la matriz de C1. Es el artefacto C1 v1.1, que aún no se regenera, no los criterios: C1-02 seguirá fallando hasta regenerarlo.
- C3 §4 (C3.md:71) dice «dominio y motores» y remite a S6-02, que trae los 5 paquetes de K25. Es coherente.
