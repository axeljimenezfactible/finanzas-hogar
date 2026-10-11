# Revisión de criterios 1.1 · ronda 1 (revisor-criterios-r1)
11-oct-2026 · Flujo «P1 · Aplicar respuestas de la puerta 1» · Revisor independiente (no escribió nada de lo revisado y no corrigió nada).
Base de evaluación: `docs/criterios/INDICE.md` §1, `docs/plan/decisiones.md` §K completa (K1–K34) y el hallazgo R2-07 de `docs/estado/wf0/REVISION-r3.md`.
Objeto: `git diff HEAD -- docs/criterios/` (19 archivos, +403 −190), `docs/criterios/tools/` y los reportes `docs/estado/p1/criterios-S-1.1.md` y `criterios-UC-1.1.md`.

## 1. Veredicto
**NO APROBADO.** 0 hallazgos de severidad alta, **6 de severidad media** y 7 de severidad baja.
- La estructura de los 18 archivos está completa: 256 criterios (196 alta, 60 media, 0 baja), cada uno con «Verifica» y «Evidencia» y con su línea en la §5. Las 7 partes de R2-07 están resueltas. Ninguna cifra difiere del oráculo. No quedan rutas `/tmp` en los comandos.
- Lo que impide aprobar: varias decisiones de §K no llegaron a los criterios que las tienen que verificar. Los criterios de interfaz U1-09, U4-03 y U4-04 no tienen K2, K3, K4, K12, K13, K16, K18, K21 ni K23. Ningún criterio de WF-3 verifica la cola de la PWA de K9. S1-09 y U3-05 no tienen los límites de quincena de K8. S8-06 no recoge el alcance de K25. Además, U2-06 conserva el gasto en 25° (lo contradice K26).

## 2. Comprobaciones pedidas (a–f)
| # | Comprobación | Estado | Evidencia (reejecutada por el revisor) |
|---|---|---|---|
| a | 5 secciones por archivo; cada criterio con Verifica y Evidencia y presente en §5 | Cumple | `rev.py` propio, revisa criterio por criterio y no solo cuántos hay en total: 18/18 archivos con §1–§5 en orden; 256/256 criterios con al menos un «- Verifica» y un «- Evidencia» en su bloque; 256/256 en su §5 (contando rangos); los ID van seguidos, sin huecos ni repetidos. Reejecución de `tools/verif_final.sh`: `cuenta.py` 134 (105/29), «incompletos: ninguno»; `traza.py` 256, 0 sin §5; `refs.py` 0 ID inexistentes |
| b | Cada K y R2-07 aplicado, con línea | Parcial | `kgrep.py` propio, con la tabla §3 y la §4. Faltan K2, K3, K4, K12, K13, K16, K18, K21 y K23 en los criterios de interfaz (H-02 y H-03), K9 en WF-3 (H-04), K8 en S1-09 y U3-05 (H-05), K25 en S8-06 (H-06) y K26 en U2-06 (H-01) |
| c | Sin rutas `/tmp` | Cumple | `rev.py` sobre los 18 archivos: 0 rutas de comando con `/tmp` fuera de `datos/tmp`. Las 10 apariciones son la viñeta de la lista de cambios que dice «Toda ruta temporal pasa de `/tmp` a …/datos/tmp/» (S1–S8 y WF4:7) y C2.md:6. `tools/`: solo `datos/tmp` (s109_sql.sh:5, 13 y 42). Paso 10 de `verif_final.sh`: 0 |
| d | INDICE: filas S y WF4 coherentes con el conteo | Cumple | Las 9 filas S1–S8 y WF4 son iguales al conteo real (tabla §5). Las filas U, C y P3 quedan para el integrador: C1 dice 15 (13/2) y debe decir **21 (19/2)**; U1 dice 14 (9/5) y debe decir **16 (10/6)**; las otras 7 son correctas. Total real **256 (196 / 60)**; la fila Total dice 248 (189 / 59) (H-12, baja) |
| e | Versión 1.1 con lista de cambios en los 18 archivos | Cumple | `rev.py`: 18/18 con «Versión 1.1 … 11-oct-2026» y «Cambios de la versión 1.1» antes de la §1 (de 1 a 10 viñetas; la de U4 dice «sin cambios» y deja una nota de pendientes) |
| f | Ningún criterio reabre una decisión de §A–§K | Cumple con observación | Revisados S5-04 y su bloqueo de frontera (H-08, baja), S2-08 y su bloqueo de redondeo (legítimo: K20 no fija en qué paso se redondea el diario y el protocolo §6 guarda el dinero en centavos enteros), C1-19 (la `Idempotency-Key` obligatoria ya está en C1 v1.0: `components.parameters.IdempotencyKey.required: true`), S3-14 (4 × 0.25 viene de decisiones C, «pesos iguales normalizados»), S4-07 (constantes publicadas de Mifflin-St Jeor; los factores intermedios llevan fuente o bloqueo) y U2-02 (135° como propuesta, lo pide R2-07) |

### Reglas de INDICE §1
| Regla | Estado | Evidencia |
|---|---|---|
| Un archivo por ítem (18 + índice), 5 secciones | Cumple | `rev.py`, `traza.py` |
| Criterio numerado: enunciado, Verifica, Evidencia y severidad | Cumple | 256/256, criterio por criterio |
| Severidad alta o media | Cumple | 196 alta / 60 media / 0 baja |
| Trazabilidad: cada criterio al plan, a decisiones o al protocolo; lo que no deriva de una frase se marca «convención propia del directorio» | Cumple con observación | `citas_rev.py`: 252 citas con fuente nombrada, 0 fragmentos que no estén textuales; 5 citas juntan fragmentos de dos fuentes o fuera de orden (H-07). C1-21 marca el tope de 80 KB como convención (H-09); el valor por omisión de C1-16 no lleva esa marca (H-10); falta la línea de K1 en la §5 de C1 (H-11) |
| Rutas absolutas, sin `jq`, `python3 -I` | Cumple | En las 403 líneas agregadas: 0 `jq`, 0 `python3` sin `-I` y 0 rutas relativas en comandos «Verifica» (la única que se detectó es prosa en S1.md:25) |
| Cifras del oráculo exactas (sha256 1a5706…bfd4, 596 hojas) | Cumple | `sha256sum` igual; `canon.py` 178 comprobadas, 0 diferencias; cifras nuevas recalculadas por el revisor (§6) |
| Identificadores `P-xx` / `N-xxx` | Cumple | Residuos de `N-xx`: 0 |

## 3. Decisiones §K contra los criterios
| K | Estado | Líneas (archivo:línea, criterio) |
|---|---|---|
| K1 lista compartida | Aplicada | C1.md:77 (C1-15, sin marca PROVISIONAL); U1.md:78 (U1-16); falta su línea en la §5 de C1 (H-11) |
| K2 alcance del presupuesto | Parcial | C1.md:81 (C1-16). Falta en U1-09 y U4 (H-02, H-03) |
| K3 borrador del ticket | Parcial | C1.md:85 (C1-17). Falta en U1-09 (H-03) |
| K4 chat de Ximena | Parcial | C1.md:89 (C1-18), S7.md:42 (S7-07), S7.md:50 (S7-09), WF4-operacion.md:35 (WF4-05). U4-03 sigue diciendo «registro del `chat_id`» (H-02) |
| K5 `family` | Aplicada | C1.md:41 (C1-06); U3.md:50 (U3-10) remite a los filtros de C1 |
| K6 valores técnicos | Aplicada | C1.md:65 (C1-12), C1.md:77 (C1-15), C1.md:93 (C1-19), S1.md:80 (S1-16), S6.md:34 (S6-04) |
| K7 campos calculados | Aplicada | S1.md:80 (S1-16), S4.md:30 (S4-04), S4.md:42 (S4-07), S5.md:23 (S5-02), S5.md:51 (S5-09), C1.md:77 |
| K8 periodos | Parcial | C1.md:41 (C1-06), C1.md:45 (C1-07), U1.md:78 (U1-16). Falta en S1.md:52 (S1-09) y U3.md:30 (U3-05) (H-05) |
| K9 captura sin conexión | Parcial | C1.md:93 (C1-19), U1.md:74 (U1-15), S6.md:34 (S6-04/05). La cola de la PWA no tiene criterio en WF-3 (H-04) |
| K10 toque | Aplicada | U1.md:30 (U1-04), U5.md:45 (U5-09), C3.md:42 (C3-08) |
| K11 login y sesión | Aplicada | S6.md:42 (S6-06), S6.md:46 (S6-07), S1.md:72 (S1-14), C1.md:65 (C1-12) |
| K12 tasas como se capturaron | Parcial | S1.md:56 (S1-10), C1.md:45 (C1-07, «Z»), C1.md:97 (C1-20). Falta mostrarlas en pantalla en U4 (H-02) |
| K13 tope de Karen | Parcial | S3.md:67 (S3-13). Falta el campo en U4 y U1-09 (H-02, H-03) |
| K14 redondeo de series | Aplicada | S2.md:35 (S2-05), los cinco modos textuales |
| K15 reposición al completar | Aplicada | S3.md:35 (S3-05) |
| K16 restricciones duras | Parcial | S3.md:63-64 (S3-12, sin valor por defecto y con caso «faltan restricciones»). Falta en U4 y U1-09 (H-02, H-03) |
| K17 puntaje | Aplicada | S3.md:67 (S3-13), S3.md:71-72 (S3-14; 100/200/400 → 1, 2/3, 0, verificado) |
| K18 tarjeta y crédito de tienda | Parcial | S3.md:47 (S3-08), S3.md:55 (S3-10). Falta la captura por producto en U4 y U1-09 (H-02, H-03) |
| K19 deuda viva | Aplicada | S2.md:43-44 (S2-07); las 7 filas, verificadas (§6) |
| K20 promedios | Aplicada | S2.md:47-48 (S2-08); el semanal queda PROVISIONAL con bloqueo legítimo |
| K21 proyección de los meses 7 a 12 | Parcial | S2.md:51 (S2-09), S2.md:88 (§3). Falta el interruptor en U4-03 y U4-04 (H-02) |
| K22 precio viejo | Aplicada con observación | S5.md:31-32 (S5-04), S1.md:28 (S1-03); bloqueo de frontera (H-08) |
| K23 meta de calorías | Parcial | S4.md:42-43 (S4-07). Falta el objetivo en U4 y U1-09 (H-02, H-03) |
| K24 prueba del OCR | Aplicada | S5.md:39 (S5-06) |
| K25 cobertura | Parcial | S6.md:26 (S6-02), S4.md:74 (S4-15), S5.md:67 (S5-13). S8.md:37 (S8-06) sin el alcance (H-06) |
| K26 color de gasto | Parcial | U2.md:21 (U2-02), U2.md:45 (U2-08), C2.md:7. U2.md:37 (U2-06) sigue en 25° (H-01) |
| K27 paleta de gráficas | Aplicada | U2.md:49 (U2-09) |
| K28 Go 1.27 | Aplicada | S6.md:22 (S6-01) |
| K29 solo LAN | Aplicada | S6.md:74 (S6-14), WF4-operacion.md:27 (WF4-03) |
| K30 token por ACL | Aplicada | S7.md:38 (S7-06), WF4-operacion.md:19 (WF4-01) |
| K31 ocho valores operativos | Aplicada | S3.md:63 (S3-12, puerta 2), S6.md:54 (S6-09), S6.md:58 (S6-10), S7.md:30 (S7-04), S7.md:46 (S7-08), S7.md:54 (S7-10), S8.md:21 (S8-02), WF4-operacion.md:31 (WF4-04), WF4-operacion.md:47 (WF4-08): PROVISIONAL con propuesta escrita |
| K32 extractos | Aplicada con observación | C1.md:101-102 (C1-21); el tope de 80 KB es convención (H-09) |
| K33 consumo | No aplica | Decisión de proceso del orquestador; no pide criterio de producto |
| K34 puerta 1 aprobada | Aplicada | 18/18 archivos en versión 1.1 |

## 4. R2-07 (REVISION-r3) contra los criterios
| Parte | Estado | Evidencia |
|---|---|---|
| U2-02 admite 135° como propuesta | Resuelto | U2.md:21 |
| U1-05 «≥ 6 flujos con los 4 obligatorios» | Resuelto | U1.md:34; también C3.md:42 (C3-08) |
| C1-01 con el comando `uvx` primero | Resuelto | C1.md:22 |
| S1-03 admite ejemplos «ilustrativos» rotulados | Resuelto | S1.md:28-29 (rótulo, conteo aparte y 0 ilustrativos sin rótulo) |
| `/tmp` → `datos/tmp` (protocolo §7) | Resuelto | Comprobación (c) |
| Citar `verify_c1.py --db` y `build.py --check` | Resuelto | C1.md:26 (C1-02), C1.md:30 (C1-03), S6.md:31 (S6-03) |
| Paráfrasis entre comillas | Resuelto con residuo bajo | 0 fragmentos que no estén textuales en 252 citas; quedan 5 citas que juntan fuentes (H-07) |

## 5. Conteo real para el INDICE (para el integrador)
| Archivo | INDICE | Real | | Archivo | INDICE | Real |
|---|---|---|---|---|---|---|
| S1 | 17 (12/5) | 17 (12/5) | | U2 | 14 (10/4) | 14 (10/4) |
| C1 | 15 (13/2) | **21 (19/2)** | | C2 | 12 (8/4) | 12 (8/4) |
| U1 | 14 (9/5) | **16 (10/6)** | | U3 | 12 (8/4) | 12 (8/4) |
| S2 | 16 (13/3) | 16 (13/3) | | U4 | 11 (8/3) | 11 (8/3) |
| S3 | 20 (17/3) | 20 (17/3) | | P3 | 13 (10/3) | 13 (10/3) |
| S4 | 15 (11/4) | 15 (11/4) | | S6 | 16 (12/4) | 16 (12/4) |
| S5 | 13 (10/3) | 13 (10/3) | | S7 | 14 (12/2) | 14 (12/2) |
| U5 | 12 (9/3) | 12 (9/3) | | C3 | 11 (9/2) | 11 (9/2) |
| S8 | 11 (8/3) | 11 (8/3) | | WF4 | 12 (10/2) | 12 (10/2) |
| **Total** | **248 (189/59)** | **256 (196/60)** | | | | |

## 6. Cifras de las líneas nuevas, recalculadas por el revisor
Script propio (`rev.py`): de las cifras con 2 decimales en las 403 líneas agregadas, las que no son hoja de `model.json` o `data.json` son todas derivadas o ilustrativas, y las recalculé:
- **S2-07** (S2.md:43), desde `periodos[i].diferencia_acumulada`, `devengado` y `despensas[k].total`: P1 3,895.00 − 973.75 = **2,921.25** = 973.75 + 1,947.50; P2 973.75 + 0.00; P3 1,258.69 + (5,034.74 − 2,517.37) = **3,776.06**; P4 1,258.70 + 0.00; P5 1,258.69 + 2,517.37 = 3,776.06; P6 1,258.70 + 0.00; 10-abr-2027 1,258.70 + 5,034.74 = **6,293.44** (`horizonte.deuda_viva_10_abr`). Coinciden las 7.
- **S2-08** (S2.md:47): el efectivo de oct, nov y feb es `periodos[0|1|4].efectivo.total` = 973.75, 1,947.50 y 2,517.38. Los diarios dan 31.41, 64.92 y 89.91. El semanal de octubre da 219.88 (diario exacto × 7) contra 219.87 (diario redondeado × 7), y el de febrero 629.35 contra 629.37. R-56 = 2,183.96 y R-57 = 2,393.75. Coinciden.
- **S4-07** (S4.md:43): 10 × 80 + 6.25 × 175 − 5 × 30 + 5 = 1,748.75; × 1.2 = 2,098.50. 10 × 60 + 6.25 × 165 − 5 × 28 − 161 = 1,330.25; × 1.2 = 1,596.30. Coinciden (son personas ficticias).
- **S5-02** (S5.md:23), con `data.json` `sustitutos[0]` (arroz: 66.00 y 84.00) y `sustitutos[2]` (leche: 135.00 y 170.00): 29.33/kg, 28.00/kg, 33.75/l y 28.33/l, iguales a las notas de la semilla.
- **S5-09** (S5.md:52): 1,000.00 ÷ 125.00 = 8, rotulado «ejemplo ilustrativo». **S3-14** (S3.md:72): la normalización mín-máx de 100, 200 y 400 da 1, 2/3 y 0.
- Los valores 0.10, 0.15 y 0.25, 1.27 y 2.62 no son dinero: son ΔL, pesos, la versión de Go y la versión de `@redocly/cli`.

## 7. Hallazgos
| ID | Sev. | Criterio de evaluación | Ruta | Descripción | Corrección sugerida |
|---|---|---|---|---|---|
| H-01 | media | §K K26 / coherencia interna | `docs/criterios/U2.md:37` (U2-06, alta) | U2-06 sigue diciendo «4.80 (gasto 25°)» como valor de partida del paso 600 de gasto. U2-02 (U2.md:21) y K26 ponen el gasto en 40–45°, así que el mismo archivo da dos matices para el gasto | Quitar «4.80 (gasto 25°)» o cambiarlo por «gasto 40–45° (K26): su contraste lo mide el script de colorimetría 1.1», y agregar la línea de K26 para U2-06 en la §5 |
| H-02 | media | §K K2, K4, K12, K13, K16, K18, K21 y K23 | `docs/criterios/U4.md:5-7`, `:22` (U4-03), `:26` (U4-04) | U4 1.1 dice «Sin cambios en los criterios U4-01 a U4-11» y deja las decisiones K como pendiente para el integrador. Las listas de U4-03 y U4-04 contra las que revisa U4 no traen el alcance del presupuesto, la confirmación del chat de Ximena (sigue el «registro del `chat_id`»), las tasas como se capturaron, el tope por persona sin valor fijo, las restricciones globales sin valor por defecto, los parámetros por producto de la tarjeta y el crédito de tienda, el interruptor de la proyección ni el objetivo de la meta de calorías | U4-03: Presupuesto (alcance «solo este mes» / «de este mes en adelante», K2); Notificaciones (candidatos y confirmación del chat, K4); Créditos y condiciones (tasas como se capturaron, K12; restricciones duras globales y ajustables con aviso «faltan restricciones», K16; parámetros por producto con «falta capturar», K18; interruptor de los meses 7–12, K21); Personas y reparto (tope por persona opcional, K13); Salud (objetivo y meta sugerida que la persona confirma, K23). U4-04: alcance = «de este mes en adelante», interruptor = apagado, restricciones y tope = sin valor. Cada uno con su línea en la §5 |
| H-03 | media | §K K2, K3, K4, K13, K16, K18, K21 y K23 / coherencia C1↔U1 | `docs/criterios/U1.md:50` (U1-09, alta) | La lista de capacidades de U1-09 sigue siendo la de §A–§J. C1-17 (C1.md:86) y C1-18 (C1.md:90) ya dan por hecho que U1 tiene las necesidades «Guardar para después» y las de listar y confirmar el chat, pero ningún criterio de U1 pide esas necesidades | Agregar a U1-09: alcance del cambio de presupuesto (K2), «Guardar para después» del ticket (K3), candidatos y confirmación del chat (K4), tope por persona (K13), restricciones del comparador con «faltan restricciones» (K16), captura de parámetros de tarjeta y crédito de tienda (K18), interruptor de la proyección (K21) y objetivo con meta sugerida y confirmada (K23), cada uno con su `P-xx` y sus `N-xxx`. Trazas en la §5 |
| H-04 | media | §K K9 («la PWA implementa la cola en WF-3») | `docs/criterios/U5.md`, `docs/criterios/C3.md` | Ningún criterio de WF-3 verifica la cola del teléfono. U1-15 la exige en la arquitectura, C1-19 en el contrato y S6-04/05 en el servidor. La lista de textos de U5-07 no incluye «Sin conexión» | Nuevo criterio en U5 (Playwright con la red cortada: capturar un gasto y un peso, ver el contador de pendientes, volver a la red y comprobar que hay un solo registro por captura, con su `Idempotency-Key` y su `captured_at`) o en C3 (F-1 y F-6 sin conexión de extremo a extremo). Agregar «Sin conexión» a U5-07, trazar a K9 y actualizar las filas del INDICE |
| H-05 | media | §K K8 (FIRME) | `docs/criterios/S1.md:52-53` (S1-09, alta), `docs/criterios/U3.md:30` (U3-05) | S1-09 describe la quincena como «(15 y 30, o último día del mes)», sin los límites de K8: del 1 al 15 y del 16 al último día, con el 31 en la segunda. Leída como fin de periodo, deja fuera el día 31. `tools/s109_sql.sh:29` sí usa 16 → último día, pero el texto que lee el revisor no lo dice. La semana de lunes a domingo no aparece en ningún criterio S. U3-05 repite «quincena 15/30 (o último día)» | S1-09: «quincena 1 del día 1 al 15 y quincena 2 del 16 al último día del mes (el 31 cae en la segunda); semana de lunes a domingo (K8); el día de pago es el 15 y el 30 o el último día». La Verifica (a) comprueba inicio y fin, p. ej. 2026-10-16…2026-10-31 y 2027-02-16…2027-02-28. U3-05: separar el día de pago de los límites de K8. Línea de K8 en la §5 |
| H-06 | media | §K K25 | `docs/criterios/S8.md:37` (S8-06, alta) | El punto (1) dice «`go test` ≥ 85 % y oráculo 100 %» sin el alcance de K25. Leído como cobertura global, contradice K25 y S6-02 (S6.md:26), que miden sin umbral `api`, `auth`, `jobs`, `reportes` y `telegram` | «(1) cobertura ≥ 85 % en cada uno de los 5 paquetes de K25 (S6-02); los demás se miden y reportan sin umbral; oráculo 100 %», con la línea de K25 en la §5 |
| H-07 | baja | INDICE §1 trazabilidad (residuo de R2-07) | `docs/criterios/C1.md:130`, `C2.md:80`, `C2.md:83`, `U2.md:95`, `U5.md:83` | 5 citas «…» juntan fragmentos de dos fuentes, o los ponen fuera de orden, en una sola cita. Cada fragmento es textual (`citas_rev.py`: 0 que no estén textuales), pero `tools/citas.py` las marca «SIN FUENTE» | Separar en una cita por fuente y en el orden del texto, hasta que `python3 -I docs/criterios/tools/citas.py` dé «sin fuente: 0» |
| H-08 | baja | §K K22 / no reabrir | `docs/criterios/S5.md:31-32` (S5-04), `docs/estado/bloqueos/criterios-frescura-extremos.md` | El bloqueo pregunta cuánto pesa el día 30. Su opción (a) es el texto literal de K22, peso = (90 − días) / 60, y es la lectura que usa el criterio. La Verifica solo prueba lo que comparten las dos lecturas, así que una implementación con /61 también pasaría | El integrador cierra el bloqueo con la lectura literal o lo lleva a Axel. Mientras siga abierto, la Verifica prueba la lectura declarada con la marca PROVISIONAL: peso(30) = 1, peso(60) = 0.5, peso(89) = 1/60 |
| H-09 | baja | INDICE §1 (convención) / no inventar valores | `docs/criterios/C1.md:101-102`, `:144` (C1-21) | El tope de 80 KB (81 920 bytes) por extracto no viene de K32 ni del protocolo §4. Está marcado como convención, pero condiciona un criterio alto sobre un entregable que todavía no existe (`docs/contratos/extractos/` no existe hoy) | Justificarlo en la §5 con el presupuesto de lectura de 400 KB del protocolo §4 y rotularlo PROPUESTA hasta que Axel lo vea, o pedir solo que el generador declare el tamaño de cada extracto |
| H-10 | baja | INDICE §1 (convención) | `docs/criterios/C1.md:81-82`, `:136` (C1-16) | K2 habla de lo que la pantalla deja preseleccionado. C1-16 además exige `default` en la API, que no viene de ninguna frase de las fuentes y no está marcado como convención | Marcar en la §5 «valor por omisión en la API: convención propia del directorio», o exigir solo el `enum` de 2 valores y el preseleccionado de la pantalla |
| H-11 | baja | INDICE §1 trazabilidad | `docs/criterios/C1.md:77` (C1-15) y §5 | C1-15 cita K1, pero la §5 de C1 no tiene línea de K1 (sí la tienen K2–K12 y K32) | Agregar «decisiones K1 «compartida en el servidor (`shopping_list`, `shopping_list_item`, 5 operaciones). Deja de ser PROVISIONAL» → C1-15» |
| H-12 | baja | Coherencia del INDICE (lo reconcilia el integrador) | `docs/criterios/INDICE.md:3`, `:7`, `:35`, `:61-64`, `:105`, §6 y §9 | La línea 3 dice «totales actualizados», pero la fila Total sigue en 248 (189 / 59) y la §1 dice «Hoy los 248». El total real es 256 (196 / 60). Las notas 8 a 11 de la §5 ya las resolvieron K26, K27, K30 y K28. La §6 no marca los bloqueos resueltos por §K ni lista los 2 nuevos. La §8 sigue con «≥ 2 meses», que K11 cambió a 90 días | Lo hace el integrador: filas C1 y U1 y total de la tabla §5 de este informe; notas 8 a 11 marcadas como resueltas; §6 con «RESUELTO (Kn)» y las filas de `criterios-promedios-redondeo` (S2-08) y `criterios-frescura-extremos` (S5-04), y lo mismo en `docs/estado/bloqueos/INDICE.md`; §8 con 90 días; §9 con este flujo |
| H-13 | baja | INDICE §1 trazabilidad (estilo) | §5 de `S1.md` a `S8.md` y `WF4-operacion.md` (p. ej. la §5 de S2) | Las trazas citan «decisiones K «…»» sin el número de la decisión, mientras que C1 y U1–U5 escriben «K2», «K9»… La cita es textual, pero cuesta ubicarla | Escribir el número: «decisiones K19 «…»» |

## 8. Estado de los hallazgos anteriores
| Hallazgo | Severidad original | Estado | Evidencia |
|---|---|---|---|
| R2-07 (REVISION-r3) · criterios | baja | **Resuelto**; queda el residuo H-07 (baja) | Tabla §4: U2.md:21, U1.md:34, C1.md:22, S1.md:28, C1.md:26 y 30, S6.md:31; 0 rutas `/tmp` en comandos; 0 paráfrasis |
| R2-05, R2-06, R3-01, R3-04 | baja | No evaluados (fuera del alcance de esta revisión) | Los tocó el agente UC en U1-11, U2-11, U1-16 y U2-08; aquí solo se revisó que no contradigan §K |

## 9. Notas de la verificación (no son hallazgos)
- `matriz_c102.py`, reejecutado: hoy U1 declara 203 N-xxx y quedan 3 sin operación ni Local (N-316, N-318 y N-319). El reporte S, que corrió antes, decía 199 y 0. La diferencia viene de U1 1.1, que otro agente del flujo está editando en paralelo. Es un estado del entregable, no de los criterios.
- 27 líneas «Evidencia» no repiten el ID de su criterio (p. ej. C1.md:79, S1.md:82 y siete de P3). Es formato heredado de WF-0 y la regla de INDICE §1 se aplica a la línea del reporte del agente, no al texto del criterio.
- La etiqueta «1.1» queda dos veces en la historia (ronda 1 de WF-0, 10-oct, y puerta 1, 11-oct). Las listas de cambios lo explican y la orden pedía 1.1.
- `datos/tmp/u2a11/` es de otro agente del flujo y no se tocó. Los temporales de este revisor (`datos/tmp/revisor-criterios-r1/`) se borraron al terminar.
