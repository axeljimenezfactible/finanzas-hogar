# Revisión de criterios 1.1 · ronda 3 de 3 (revisor independiente `revisor-criterios-r3`)

11-oct-2026 · Flujo «P1 · Aplicar respuestas de la puerta 1» · Alcance: los 18 archivos de `docs/criterios/` y su `INDICE.md`, evaluados **solo** contra `INDICE.md` §1, `docs/plan/decisiones.md` §K y el hallazgo R2-07 de `docs/estado/wf0/REVISION-r3.md`. No corregí nada; todas las verificaciones las reejecuté con scripts propios en `datos/tmp/rv3/` (ya borrados) y con `docs/criterios/tools/verif_final.sh`.

## Veredicto
**APROBADO.** 0 hallazgos de severidad alta, 0 de severidad media y 4 de severidad baja (se registran, no bloquean). El hallazgo R2-M01 de la ronda 2 queda **cerrado**.

## 1. Verificaciones pedidas (a–f)
| # | Qué | Cómo lo reejecuté | Resultado | Estado |
|---|---|---|---|---|
| a | 5 secciones por archivo; cada criterio con «Verifica» y «Evidencia» y presente en su §5 | `cuenta_rv3.py` propio (regex de definición `**ID (sev)`, bloques hasta el siguiente criterio o `## `, IDs y rangos «a/…/–» en §5) y, aparte, `bash docs/criterios/tools/verif_final.sh` | 18 archivos con §1–§5; **257 criterios (197 alta / 60 media / 0 baja)**; 0 sin Verifica, 0 sin Evidencia, 0 fuera de §5, 0 IDs repetidos. `traza.py` y `refs.py` dan lo mismo (257; 0 citados inexistentes). Contra `HEAD`: 248 → 257, **+9** (C1-16…C1-21, U1-15 y U5-13 alta; U1-16 media), 0 quitados, 0 severidades cambiadas | Cumple |
| b | Cada K de §K y R2-07 aplicado, con línea | `kmap.py` propio (menciones `K<n>` por archivo y línea) + lectura del texto de cada criterio citado | K1–K32 con criterio y línea (tabla §2); K33 y K34 son de proceso y no requieren criterio. R2-07: los 7 puntos aplicados (tabla §3) | Cumple |
| c | Sin rutas `/tmp` | `grep -nP '(?<!datos)/tmp\b'` en los 18 archivos y en `tools/` | 0 rutas operativas; quedan 11 menciones descriptivas en las listas de cambios («pasa de `/tmp` a …datos/tmp/») y 1 patrón de `grep` en `verif_final.sh` | Cumple |
| d | INDICE: filas S y WF4 coherentes; filas U/C/P3 | conteo propio contra la tabla de `INDICE.md` §2 | S1 17 (12/5), S2 16 (13/3), S3 20 (17/3), S4 15 (11/4), S5 13 (10/3), S6 16 (12/4), S7 14 (12/2), S8 11 (8/3), WF4 12 (10/2): **iguales**. U/C/P3 (conteo correcto): C1 21 (19/2), C2 12 (8/4), C3 11 (9/2), P3 13 (10/3), U1 16 (10/6), U2 14 (10/4), U3 12 (8/4), U4 11 (8/3), U5 13 (10/3): la tabla ya los trae; total 257 (197/60) igual. Queda texto viejo en §1, §5 y §6 (hallazgo B-01) | Cumple con observación |
| e | Versión 1.1 con lista de cambios en los 18 | línea 3 y bloque «Cambios de la versión 1.1» antes de `## 1.` | 18 de 18 con «Versión 1.1 (11-oct-2026…)» y lista de 2 a 10 viñetas (C1 9, C2 2, C3 2, P3 2, S1 9, S2 8, S3 7, S4 6, S5 7, S6 10, S7 6, S8 6, U1 7, U2 6, U3 3, U4 4, U5 2, WF4 7) | Cumple |
| f | Ningún criterio reabre una decisión de §A–§K | lectura de los criterios que citan cada K y búsqueda de valores viejos (gasto 25°, `go 1.26`, `SameSite=Lax`, 374.93 como tope, lista de compras local o PROVISIONAL, marcas PROVISIONAL sobre K5/K6/K7) | 0 contradicciones. Los dos bloqueos nuevos (`criterios-promedios-redondeo`, `criterios-frescura-extremos`) preguntan lo que §K no fija, usan la lectura literal mientras tanto y marcan PROVISIONAL solo esa parte; formato del protocolo §2 completo. U2-02 admite 135° como propuesta, que decisiones I prevé («o U2 propone 135°») | Cumple |

### Cifras reproducidas desde `oraculo/model.json` (sha256 `1a57066689dd…bfd4`)
- **S2-07 (K19)**, deuda viva = (devengado acumulado − efectivo acumulado) + saldo no devengado, con compras a crédito y papel de contado por fecha: P1 2,921.25 = 973.75 + 1,947.50; P2 973.75 = 973.75 + 0.00; P3 3,776.06 = 1,258.69 + 2,517.37; P4 1,258.70 = 1,258.70 + 0.00; P5 3,776.06 = 1,258.69 + 2,517.37; P6 1,258.70 = 1,258.70 + 0.00; 10-abr-2027 6,293.44 = 1,258.70 + 5,034.74 (= `horizonte.deuda_viva_10_abr`). Diferencias acumuladas iguales a `periodos[i].diferencia_acumulada`. **7 de 7 exactas.**
- **S2-08 (K20)**: diario oct 973.75 ÷ 31 = 31.41; nov 1,947.50 ÷ 30 = 64.92; feb 2,517.38 ÷ 28 = 89.91. Semanal sin redondear / con diario redondeado: 219.88 / 219.87 y 629.35 / 629.37 (lo que dice el criterio y su bloqueo). **Exactas.**
- **S2-09 (K21)**: pagos de la despensa 4 = 2027-04-30, 2027-05-15, 2027-05-30, 2027-06-15. **Iguales.**
- **K8 / K22**: 2026-10-12 lunes, 14 miércoles, 18 domingo, 31 sábado; 10-oct-2026 + 30 días = 9-nov-2026. **Iguales a S1-03, S1-09, S2-06 y S5-04.**
- `verif_final.sh`: `canon.py` 178 comprobadas, 0 diferencias; `s109_sql.sh` 18 filas, quincena 2 de oct `2026-10-16`–`2026-10-31`, 0 fuera de límites, ventana 10→9 rechazada; `s109_equiv.py` 12 de 12.

## 2. Respuestas de §K: dónde quedó cada una
| K | Criterio y línea (`docs/criterios/`) | Lectura |
|---|---|---|
| K1 | C1.md:77 (C1-15) | lista de compras sin marca; 0 menciones de lista local o PROVISIONAL |
| K2 | C1.md:81 (C1-16); U4.md:24 (U4-03), 28 (U4-04); U1.md:51 (U1-09) | `scope`, «de este mes en adelante» preseleccionado |
| K3 | C1.md:85 (C1-17); U1.md:51 | borrador en columna de 00002 y operación nueva |
| K4 | C1.md:89 (C1-18); S7.md:42 (S7-07), 50 (S7-09); WF4-operacion.md:35 (WF4-05); U4.md:24 | detecta, ella confirma; nota de MultiX con bloqueo si no hay vía limpia; Ximena «pendiente de acción» |
| K5 | C1.md:41 (C1-06) | tres valores incluyentes; sin `x-provisional` |
| K6 | C1.md:65 (C1-12), 77 (C1-15), 93 (C1-19); S1.md:81 (S1-16); S6.md:34 (S6-04), 38 (S6-05) | los 9 valores iguales a §K; 0 marcas PROVISIONAL sobre ellos |
| K7 | C1.md:77; S1.md:81; S4.md:30 (S4-04), 42 (S4-07); S5.md:23 (S5-02), 51 (S5-09) | IMC 18.5/25/30, unidad base HALF_UP, equilibrio de Sam's, inventario |
| K8 | S1.md:29 (S1-03), 53 (S1-09); S2.md:40 (S2-06); C1.md:41 (C1-06), 45 (C1-07); U1.md:79 (U1-16); U3.md:31 (U3-05) | 1–15 y 16–último día, el 31 en la segunda; lunes a domingo; día de pago aparte (ver B-02 sobre C1) |
| K9 | C1.md:93 (C1-19); S6.md:34 (S6-04); U1.md:75 (U1-15); U5.md:62 (U5-13) | cola solo de gastos y peso; clave y `captured_at`; sin duplicar |
| K10 | U1.md:31 (U1-04); U5.md:46 (U5-09); C3.md:42 (C3-08) | solo toques en controles; 3 toques |
| K11 | C1.md:65 (C1-12); S1.md:73 (S1-14); S6.md:42 (S6-06), 46 (S6-07) | 15 min deslizante por usuario e IP; 90 días renovable |
| K12 | C1.md:45 (C1-07), 97 (C1-20); S1.md:57 (S1-10); U4.md:24 | decimales capturados; «Z» solo en entradas |
| K13 | S3.md:67 (S3-13); U4.md:28 (U4-04); U1.md:51 | sin tope fijo; 374.93 solo como cifra del oráculo, nunca como tope |
| K14 | S2.md:36 (S2-05) | los dos modos del oráculo |
| K15 | S3.md:35 (S3-05) | capital repuesto solo al completar el pago |
| K16 | S3.md:63 (S3-12); U4.md:24, 28 | globales y ajustables; sin valor por defecto; «faltan restricciones» |
| K17 | S3.md:67 (S3-13), 71 (S3-14) | mín-máx; desviación estándar; unidades base por peso (100/200/400 → 1, 2/3, 0 comprobado) |
| K18 | S3.md:47 (S3-08), 55 (S3-10); U4.md:24 | parámetros por producto; «falta capturar» |
| K19 | S2.md:44 (S2-07) | identidad exacta; 7 cifras reproducidas |
| K20 | S2.md:48 (S2-08) | días reales; semanal PROVISIONAL por bloqueo |
| K21 | S2.md:52 (S2-09); U4.md:24, 28 | interruptor apagado por defecto |
| K22 | S1.md:29 (S1-03); S5.md:31 (S5-04) | 30 días; peso 1 hasta el 29; PROFECO no compite; propuesta del orquestador rotulada |
| K23 | S4.md:42 (S4-07); U4.md:28 | Mifflin-St Jeor × 5 niveles con ajuste por objetivo; la persona confirma |
| K24 | S5.md:39 (S5-06) | 5 tickets de Axel; «sin medir» si faltan |
| K25 | S4.md:74 (S4-15); S5.md:67 (S5-13); S6.md:26 (S6-02); S8.md:38 (S8-06) | 85 % en 5 paquetes; el resto sin umbral |
| K26 | U2.md:22 (U2-02), 38 (U2-06), 46 (U2-08); C2.md:7 | gasto 40–45°, error 20°, Karen y Bodega reubicadas, marca en el barrido |
| K27 | U2.md:50 (U2-09) | ΔL ≥ 0.10 entre consecutivas; pares informativos |
| K28 | S6.md:22 (S6-01) | Go 1.27 |
| K29 | S6.md:74 (S6-14); WF4-operacion.md:27 (WF4-03) | solo LAN; Tailscale no se abre |
| K30 | S7.md:38 (S7-06), 42 (S7-07); WF4-operacion.md:19 (WF4-01) | ACL de lectura; token en ejecución |
| K31 | S3.md:63 (S3-12, puerta 2); S6.md:54 (S6-09), 58 (S6-10); S7.md:30 (S7-04), 46 (S7-08), 54 (S7-10); S8.md:22 (S8-02); WF4-operacion.md:31 (WF4-04), 47 (WF4-08); S1.md:81 | PROVISIONAL con propuesta escrita hasta su puerta |
| K32 | C1.md:101 (C1-21) | extractos por módulo; 80 KB marcado como convención propia del directorio (C1.md:144) |
| K33, K34 | — (proceso) | K34 se refleja en la versión 1.1 de los 18 archivos |

## 3. Hallazgo R2-07 (`docs/estado/wf0/REVISION-r3.md`)
| Punto | Dónde | Estado |
|---|---|---|
| U2-02: admitir 135° como propuesta | U2.md:22 | Cerrado (traza imprecisa en §5: B-03) |
| U1-05: «≥ 6 flujos con los 4 obligatorios» | U1.md:35-36; C3.md:42-43 | Cerrado |
| C1-01: comando `uvx` primero | C1.md:22 | Cerrado |
| S1-03: ejemplos «ilustrativos» rotulados | S1.md:29 | Cerrado |
| `/tmp` → `datos/tmp` | 18 archivos | Cerrado (0 rutas operativas) |
| Citar `verify_c1.py --db` y `build.py --check` | C1.md:26 (C1-02), 30 (C1-03); S6.md:31 (S6-03) | Cerrado (los dos comandos existen en `backend/internal/api/c1gen/`) |
| Paráfrasis entre comillas en §5 | C1, C2, U2, U5 | Casi cerrado: quedan 5 citas con fragmentos textuales unidos por «…» fuera del orden de la fuente (B-04) |

## 4. Estado de los hallazgos anteriores
| ID | Sev. | Estado | Evidencia |
|---|---|---|---|
| R2-M01 (S2-06 y S1-03 contra K8) | media | **Cerrado** | S2.md:40 aplica K8 (quincena 2 del 16 al último día, el 31 en la segunda; día de pago «otro dato y no es el fin del periodo») y S2.md:41 trae la fixture del 2026-10-31 con valor esperado a mano: «un motor que cierre la quincena 2 el día 30 falla la prueba». S2.md:107 traza a K8. S1.md:29 cambia «quincena 15 y 30» por «día de pago 15 y 30» y agrega la regla de periodos de K8 con ejemplo ilustrativo correcto. 0 restos de «no el 31» o «16 al 30» en los 18 archivos |

## 5. Hallazgos de esta ronda (todos baja)
| ID | Sev. | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| B-01 | baja | INDICE §1 y coherencia | La tabla de §2 ya suma 257 (197/60), pero §1 dice «Hoy los 248 criterios de la tabla» (línea 7); la nota 8 de §5 sigue con «gasto (25°)… a 5°» (línea 61) y la 10 con el token como problema, ya resueltos por K26 y K30; §6 no lista los bloqueos nuevos `criterios-promedios-redondeo` y `criterios-frescura-extremos` (citados por S2-08 y S5-04) ni marca los resueltos por §K | `docs/criterios/INDICE.md:7`, `:61-63`, §6 (líneas 67-86) | Integrador: 257 en §1; notas 8-10 como resueltas por K26, K27 y K30; agregar las dos filas de bloqueo y marcar RESUELTO lo cerrado por §K |
| B-02 | baja | §K K8 (redacción) | C1-06 conserva «quincena 15/30 (o último día del mes)» y C1-07 «las quincenas (15 y 30, o último día)» junto a la regla de K8: es la misma mezcla de día de pago y periodo que R2-M01 quitó de S1 y S2. No contradice K8 porque la frase trae K8 explícito, pero las Verifica no comprueban los límites | `docs/criterios/C1.md:41`, `:45` | Decir «día de pago 15 y 30 (o último día del mes)» y dejar «quincena» solo para los periodos de K8; en C1-07, que la Verifica busque en `C1-datos.md` la declaración 1–15 y 16–último día |
| B-03 | baja | INDICE §1 (traza) | U2 §5 marca «se admita 135°» como «convención propia del directorio», pero decisiones I lo prevé textualmente: «ingreso 150° (comparte familia con éxito: se distingue por icono y contexto, o U2 propone 135°)» | `docs/criterios/U2.md:100` | Citar decisiones I «o U2 propone 135°» → U2-02 y quitar la marca de convención |
| B-04 | baja | R2-07 (citas) | 5 citas de §5 unen fragmentos textuales con «…» en un orden distinto al de la fuente; `tools/citas.py` las reporta «SIN FUENTE». Comprobé que cada fragmento existe tal cual: falla solo el orden | `docs/criterios/C1.md` §5 («Moneda MXN…»), `C2.md` §5 (2: «tokens compilan…», «Style Dictionary 5…»), `U2.md:96`, `U5.md` §5 («Android (Axel y Ximena)…») | Reordenar los fragmentos según la fuente o partirlos en citas separadas, hasta que `citas.py` dé 0 |

## 6. Observaciones sin hallazgo
- `matriz_c102.py` reporta hoy N-316, N-318 y N-319 sin operación: es estado de C1 (U1 1.1 ya tiene 203 necesidades y C1 v1.1 no está regenerado) y no de los criterios; C1-02 lo detectará. `docs/contratos/extractos/` y `build.py --extractos` aún no existen (C1-21 lo exige a C1 v1.1).
- S5-04: la lectura literal de K22 (peso (90 − días)/60) es la que se usa; el bloqueo pide confirmar el día 30 y la Verifica prueba solo lo común a las dos lecturas. No reabre K22.
- S2-06: la Verifica dice «los 3 casos» y la Evidencia «4 casos»; el cuarto es la semana. Coherente.
