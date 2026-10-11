# REVISION-criterios-r3 · Criterios 1.2 (flujo P1c, puerta 1b) · ronda 3 de 3 · 11-oct-2026

Revisor independiente `revisor-criterios-r3`. Evalúo solo contra `docs/criterios/INDICE.md` §1 y `docs/plan/decisiones.md` §K (K35 a K44). No escribí nada de lo revisado y no corregí nada. Todas las comprobaciones las reejecuté con comandos y scripts propios en `datos/tmp/rev3c/` (borrados al terminar).

## Veredicto
**NO APROBADO.** 1 hallazgo alto, 2 medios y 4 bajos. Los cuatro hallazgos medios de la ronda 2 quedan cerrados (dos con un residuo bajo). El hallazgo alto es nuevo: la corrección de S1 (H-01 de `REVISION-S1-r1`), que se hizo en paralelo (00003 modificado a las 15:23:44, `C1.md` y `U4.md` a las 15:24), llevó los umbrales 80 % y 100 % y los días antes del pago a claves `alerts.*` de `app_setting`, como manda K44, y redujo el bloqueo `S1-parametros-avisos` a `send_time` y `advance_days`. Los criterios C1-22, U4-12, S7 y el `INDICE.md` §6 siguen diciendo lo contrario.

## Comprobaciones pedidas
| # | Comprobación | Estado | Evidencia (reejecutada) |
|---|---|---|---|
| a | Estructura: 5 secciones, Verifica y Evidencia en cada criterio, cada ID en §5 | Cumple | Script propio: 18 archivos, 261 criterios (201 alta / 60 media / 0 baja); 0 sin Verifica, 0 sin Evidencia, 0 fuera de §5 (con rangos expandidos), 0 repetidos. S2, S3 y S5 traen además un §6 anexo de cifras canon desde WF-0 (aprobado en la puerta 1, K34). `verif_final.sh` pasos 1–4: 261/261, 305 citas, 0 sin fuente. |
| b | Cada K con línea | Cumple, salvo K44 en avisos (A) | K39 `S2.md:51` (219.88, 629.35); K40 y K37 `S5.md:36`; K36 `S4.md:45`, `U4.md:72`, `C1.md:116`; K41 `U2.md:27`, `U2.md:51`; K42 `U2.md:51`; K44 `U4.md:68`, `C1.md:116`, `S3.md:66`, `S6.md:61`; K43 `C1.md:120`, `S3.md:66`, `S6.md:61`; K38 `U1.md:83`; K35 `P3-claude-design.md:3,6,22,27` y §5. |
| c | Cifras | Cumple | `Decimal`: 973.75 / 31 × 7 = **219.88** (doble redondeo 219.87); 2,517.38 / 28 × 7 = **629.35** (629.37); (90 − 30) / 61 = **0.9836**; (90 − 60) / 61 = **0.4918**; día 90 = 0; día 91 no compite. 973.75 y 2,517.38 están en `oraculo/model.json` (P1 y P5 en efectivo, `s109_equiv.py`). Fixture de S4-07: 1,748.75 × 1.2 = 2,098.50; 1,330.25 × 1.2 = 1,596.30; − 300 = 1,798.50; + 250 = 1,846.30. `canon.py` 178 / 0 diferencias. |
| d | P3 1.2 sin navegador del orquestador; paquete completo y carpetas de salida | Cumple con B, C y F | P3-01, P3-03 y P3-04 prohíben el navegador al orquestador; P3-02 lista exactamente los 22 elementos de K35 (11 archivos + `reporte-p01..p11.png`) y `MANIFIESTO.md` con sha256 (21 hashes); P3-11 exige las carpetas `salida/iter-N/` y `salida/final/`. |
| e | 0 PROVISIONAL que citen bloqueos resueltos | Cumple en lo literal; ver A | Script propio: 42 bloqueos con «Respuesta de Axel»; las líneas con PROVISIONAL y uno de ellos son avisos de que se quitó (S2:6, S2:51, S3:6, S4:6, S5:6, S5:36, U1:6, U2:6) o Verifica que buscan ausencia (C1:89, C1:121). Abiertos: `S1-parametros-avisos` (y `criterios-nxx`, cerrado por su «Resolución» de WF-0). Pero las marcas sobre umbrales y días antes del pago citan un bloqueo que ya no los cubre (A). |
| f | INDICE coherente con el conteo | Cumple, salvo §6 (A) | Las 18 filas de §2 y §10 y el total **261 (201 / 60)** coinciden con mi conteo; `conteo.py`: 0 diferencias. §6 describe mal el alcance del bloqueo abierto (A). |
| g | Ningún criterio reabre §A–§K ni inventa valores | No cumple (A) | C1-22, U4-12 y S7 declaran «no decidido» el lugar de los umbrales 80/100 % y días antes del pago, que K44 fija en `app_setting`. Sin valores inventados: el 413 ya existe en C1 (`extractos/gastos.yaml:671`), M-05 existe (`C1-datos.md:329`), 80 000 bytes lo declara `extractos/INDICE.md:10`. |

Comprobación del esquema (base temporal, `goose up` de 00001 a 00003): `app_setting` de la casa tiene `alerts.budget_warning_pct` = 80, `alerts.budget_limit_pct` = 100 y `alerts.payment_due_days_before` = [1, 0]; `notification.params_json` queda `payment_due` = `{}` y `budget_threshold` = `{"scopes": ["category", "total"]}`. Las 23 claves sembradas citan una decisión según la regex de C1-22 (`decisión [A-Z]` o `K[0-9]+`).

## Hallazgos

### A (alta) · Parámetros de avisos: los criterios contradicen K44 y el esquema 1.3
- **Ruta:** `docs/criterios/C1.md:9` y `:116-117` (C1-22, enunciado y Verifica); `docs/criterios/U4.md:7` y `:68-69` (U4-12); `docs/criterios/S7.md:7`, S7-04 (`:35`), S7-08 (`:51`), S7-10 (`:59`) y `:109`; `docs/criterios/INDICE.md:68`, `:96` y `:174`.
- **Qué pasa:** K44 nombra los «umbrales de alertas 80 % y 100 %» y fija «cada parámetro es una clave de `app_setting` con valor inicial sembrado». El esquema 1.3 ya lo cumple (00003, corrección S1-r1 H-01: claves `alerts.*` con cita «Decisión F, K44», `thresholds_pct` y `days_before` quitados de `params_json`). El bloqueo `S1-parametros-avisos` ahora solo pregunta dónde viven `send_time` y `advance_days`. Pero C1-22 dice que los umbrales y `days_before` «viven hoy en `notification.params_json`», que su lugar «no está decidido» y que «no se siembran claves `alerts.*` en `app_setting`»; su Verifica exige `budget_threshold` y `payment_due` en `params_json` con [80, 100] y [1, 0], y con el esquema 1.3 falla. U4-12 manda poner `notification.params_json` en la columna «clave» de esos tres valores y los cruza contra `params_json`, lo que también falla. S7 y el INDICE §6 («opción A, la que usa hoy el esquema 1.3») repiten esa lectura. El resultado: se reabre K44 para valores ya decididos, quedan marcas PROVISIONAL sobre algo decidido con la cita de un bloqueo que ya no los cubre, y dos criterios altos tienen una Verifica que falla con el esquema correcto.
- **Corrección sugerida:** en C1-22 y U4-12, pasar los umbrales y los días antes del pago a la regla general (claves `alerts.budget_warning_pct` = 80, `alerts.budget_limit_pct` = 100 y `alerts.payment_due_days_before` = [1, 0], sembradas en 00003, con cita a la decisión F y K44, expuestas por settings con valor actual e inicial). Quitar «no se siembran claves `alerts.*`» y la comprobación contra `params_json`. Dejar la excepción PROVISIONAL solo para `send_time` y `advance_days`: lugar según `S1-parametros-avisos` y valor inicial según K31 (S7). En el enunciado de U4-12 (línea 68), decir que esos dos no tienen valor inicial y que su «restablecer» se activa cuando S7 proponga el valor; hoy dice que «siguen mostrando el valor inicial» (residuo del hallazgo 2 de r2). En S7 (lista de cambios, S7-04, S7-08, S7-10 y §5), decir que los umbrales y los días antes se leen de `app_setting` `alerts.*`. En INDICE §6 (líneas 68 y 96) y §11 (línea 174), reducir el alcance del bloqueo a `send_time` y `advance_days`.

### B (media) · P3-06: el conteo por iteración falla con el proceso de K35
- **Ruta:** `docs/criterios/P3-claude-design.md:43-44`.
- **Qué pasa:** la Verifica pide «por iteración, conteo de archivos = número de artboards de `prompt.md`». K35 pone `siguiente-prompt.md` dentro de `docs/design/salida/iter-N/`, así que cada iteración revisada, salvo la última, tiene artboards + 1 archivos. Además el enunciado exige a Axel «una imagen o archivo por artboard» con el nombre de la pantalla `P-xx`. K35 solo dice que Axel «guarda lo que descargue», y Claude Design puede exportar un solo PDF o HTML con todos los artboards. En la versión 1.1 decía «conteo de capturas», que hacía el orquestador.
- **Corrección sugerida:** contar solo lo descargado, sin `siguiente-prompt.md`. Que el agente Sonnet asigne cada archivo o página a su artboard en `P3-iter-N.md` (un PDF de N páginas cuenta N). El nombre por `P-xx` queda como sugerencia del `LEEME.md`, no como condición del criterio.

### C (media) · P3-04: la Verifica no tiene ruta y da coincidencias con un proceso correcto
- **Ruta:** `docs/criterios/P3-claude-design.md:36`.
- **Qué pasa:** `grep -rn "claude.ai/design"` «en los registros de acciones del orquestador» no nombra ningún archivo (INDICE §1 exige rutas absolutas). Si se corre sobre el registro del flujo, encuentra la escritura de `LEEME.md`, que por P3-03 debe decir «abrir claude.ai/design». Lo que importa es que haya 0 acciones de navegador, no que aparezca o no el texto.
- **Corrección sugerida:** comprobar que en el registro del flujo WF-1 (con ruta absoluta) hay 0 llamadas a herramientas de navegador (`mcp__claude-in-chrome__`, `navigate`), o que el reporte del orquestador declare «0 acciones de navegador» y el revisor lo contraste con la lista de herramientas usadas.

### D (baja) · P3-01, P3-03 y P3-12: residuo de falsos positivos en las regex
- **Ruta:** `docs/criterios/P3-claude-design.md:23,32,68`.
- **Qué pasa:** la prueba propia confirma que «No escribas tu contraseña…» y `C2-tokens.json` dan 0 (hallazgo 4 de r2 cerrado). Pero una advertencia en formato de etiqueta, como «**Contraseña:** nunca la pegues aquí» (P3-01 y P3-03) o «Paso clave: exporta…» (P3-03), coincide.
- **Corrección sugerida:** exigir un valor después del separador (`(password|contraseña|clave) *[:=] *[^ ]{6,}` sin palabras de advertencia) o aceptar las coincidencias que el revisor del paquete justifique por escrito.

### E (baja) · U4-12 ubica `send_time` y `advance_days` como si los ubicara K44
- **Ruta:** `docs/criterios/U4.md:68`.
- **Qué pasa:** la propuesta de K44 solo ubica las alertas en Notificaciones. «Anticipación de avisos» y la hora del resumen diario no tienen sección, y por la propia regla de U4-12 deberían llevar «la sección que U4 proponga», que también podría ser «Reportes y frecuencias». Hoy U4-12 los atribuye a la propuesta de K44.
- **Corrección sugerida:** pasar `send_time` y `advance_days` a la frase «los demás… llevan la sección que U4 proponga, con su justificación en SUPUESTOS».

### F (baja) · Las carpetas de salida no se comprueban en la puerta 3
- **Ruta:** `docs/criterios/P3-claude-design.md:63-65`; `docs/criterios/INDICE.md` §3.
- **Qué pasa:** P3-11 exige que las carpetas de salida existan al entregar el paquete, pero no está entre los criterios de la puerta 3 (solo P3-01 y P3-02), y su Verifica solo hace `ls -d` de `salida/`, no de `salida/final/`.
- **Corrección sugerida:** sumar `ls -d …/docs/design/salida/final/` a la Verifica de P3-11 y meter P3-11 en la fila de la puerta 3 (o pasar esa comprobación a P3-01).

### G (baja) · S6-10: «10 MB + 1 byte» sin definir el MB
- **Ruta:** `docs/criterios/S6.md:61-62`.
- **Qué pasa:** la prueba exige aceptar 10 MB y rechazar 10 MB + 1 byte, pero no dice si 10 MB son 10,000,000 o 10,485,760 bytes. K43 no lo fija.
- **Corrección sugerida:** tomar el número de bytes que declare C1 v1.2 en `uploadExpenseAttachment` (o `attachment.max_size_mb` × el factor que fije C1) y citarlo.

## Estado de los hallazgos de la ronda 2
| Hallazgo r2 | Estado | Evidencia |
|---|---|---|
| 1 (media) · P3-02 pedía evidencia en `P3-iter-1.md` | **Cerrado** | `P3-claude-design.md:29`: evidencia en el reporte de `paquete-claude-design` y de su revisor (`docs/estado/wf1/`); aclara que `P3-iter-1.md` no es evidencia. |
| 2 (media) · U4-12 pedía valor inicial de `send_time` y `advance_days` | **Cerrado en la Verifica**, con residuo | `U4.md:69`: «sin valor inicial: lo propone S7 (K31), PROVISIONAL», la celda no cuenta como vacía. El enunciado (línea 68) todavía dice que los de avisos «siguen mostrando el valor inicial»: se corrige con A. |
| 3 (media) · C1-22 y U4-12 exigían la siembra en 00003 y la cita de un K-n | **Cerrado** (ahora lo supera A) | Las dos nombran las 5 claves de 00001 y 00002 y aceptan `decisión [A-Z]` o `K[0-9]+`; en la base temporal, las 23 claves sembradas cumplen esa regex. |
| 4 (media) · regex de P3 con la palabra suelta | **Cerrado**, con residuo D | Prueba propia: un LEEME correcto que nombra `C2-tokens.json` y advierte en prosa da 0; `password: abc`, `bot_token=12` y un token de Telegram dan 3. |

## Verificación reejecutada
- `bash docs/criterios/tools/verif_final.sh`: conteo 261 (201/60/0); traza 0 faltantes; refs 261/261; citas 305/0; residuos N-xx 0; canon 178/0; S1-09 SQL OK (54 tablas, K8) y equivalencia 12/0; rutas `/tmp` 0; conteo contra INDICE 0 diferencias. El paso 9 sigue con N-320…N-327 sin operación (B-01 de r1, en `C1-datos.md`, fuera de los criterios).
- Script propio `datos/tmp/rev3c/chk.py`: estructura, PROVISIONAL contra bloqueos resueltos y cifras (salidas arriba).
- `goose` v3.28.0 `up` en una base temporal y consulta de `app_setting` y `notification`.
- Regex de P3 contra dos textos de muestra propios.
