# Revisión U1 1.2 · ronda 3 de 3 (flujo P1c, 11-oct-2026)

Revisor: **revisor-U1-r3** (independiente; no escribió nada de lo revisado y no corrige). Evalúa solo contra `docs/criterios/U1.md` 1.2 y `decisiones.md` K35, K36, K38, K43 y K44. Objeto: `docs/design/U1-arquitectura.md` 1.2 (`git diff HEAD`, solo lo cambiado), `docs/design/tools/verificar_u1.py`, `docs/estado/p1c/U1-1.2.md` y la sección «Para C1 (corregida r2)» de `docs/estado/p1c/correccion-U1-r2.md`.

## Veredicto
**APROBADO.** Hay 0 hallazgos de severidad alta y 0 de severidad media. H-r2-01 (media) queda cerrado por la «otra salida» que admitió la ronda 2: el verificador y la fila 9 de los cambios 1.2 mandan a C1 v1.2 a la sección corregida, y la comprobación nueva rechaza la viñeta vieja (lo probé con una prueba negativa). Quedan 10 hallazgos bajos: 3 nuevos y 7 que vienen de la ronda 2. Ninguno cambia cifras, esquema ni pantallas.

## Verificación propia (reejecutada; temporales en `datos/tmp/rev-u1-r3/`, ya borrados)
- `python3 -I docs/design/tools/verificar_u1.py` sale con 0, da 63 OK y 0 FALLA y termina en «RESULTADO: TODO OK». La única línea PENDIENTE es `PENDIENTE C1 v1.2: 8 N-xxx nuevas … («Para C1 (corregida r2)» de docs/estado/p1c/correccion-U1-r2.md): N-320 … N-327`.
- **Recuento propio** (script `check.py`, que no usa el verificador), sobre §3, sin contar §3.29 ni §3.30:
  - **211** N-xxx únicas: 138 de lectura (L) y 73 de escritura (E), sin duplicados; la más alta es N-327.
  - Las nuevas son N-320 a N-327, todas E. Las Local son N-003, N-301 y N-317.
  - Matriz §3.29: 41 filas con N-xxx. Coincide con §3 en los dos sentidos (faltan 0, sobran 0), y la columna L/E concuerda en el 100 % de los casos.
  - La línea 817 declara «211 · 138 · 73» y coincide con el recuento.
  - Pantallas en §2.1: 27, sin duplicados.
  - Secciones de Ajustes: 14, de P-26.01 a P-26.14, con los nombres exactos de U1-08. `P-26.15` aparece 0 veces y «15 secciones» también 0.
- **Cruce propio con openapi.yaml** (`grep` de las listas `x-necesidades`, sin leer el archivo): hay 200 N-xxx con operación.
  - N-xxx de U1 sin operación: N-003, N-301 y N-317 (las tres Local) y N-320 a N-327.
  - N-xxx citadas por operaciones que no existen en U1: ninguna.
  - Las «Nuevas» de «Para C1 (corregida r2)» son exactamente N-320 a N-327, igual que la línea PENDIENTE.
- **Cruce con el esquema 1.3**: las 15 líneas `@param` de `backend/db/migrations/00003_puerta1b.sql` (líneas 35 a 49) coinciden clave por clave y valor por valor con la tabla de §2.3 (líneas 309 a 323) y con N-320 a N-327:
  - 5, 15, 12 y 90;
  - 30, 90 y 1;
  - 5; 10 y 24; 0 y 0;
  - 80, 100 y [1, 0].

  Las secciones de RN-29 en `reglas-negocio.md` (líneas 291 a 305) son las mismas que las de U1. `send_time` y `advance_days` son `@param_sin_valor` (líneas 53 y 54), y U1 dice «Sin valor inicial decidido (K31)», sin «Restablecer».
- **Bloqueos**: de los que hay en `docs/estado/bloqueos/`, solo `S1-parametros-avisos` sigue abierto sin «Respuesta de Axel». `criterios-nxx` aparece así en el listado porque su resolución usa otro rótulo; su encabezado dice RESUELTO. U1 cita `S1-parametros-avisos` 3 veces, como abierto. Los resueltos que cita (`U1-lista-vigente-domingo`, `C1-valores-k31`, `S1-ajuste-objetivo-kcal`, `criterios-umbrales-ausentes`, …) los nombra como resueltos.
- **PROVISIONAL**: `grep` da 4 apariciones (líneas 22, 38, 39 y 52), todas en el historial de cambios. En el cuerpo no queda ninguna.
- **Pruebas negativas** (copia de `docs/`, `backend/db/`, `oraculo/` y `datos/semilla/data.json`):
  - V0, copia sin cambios: sale con 0, «TODO OK».
  - V1, «Para C1 (corregida r2)» con la viñeta vieja: FALLA H-r2-01 y PENDIENTE «sin fuente coherente».
  - V2, PROVISIONAL en §0.2 con un bloqueo resuelto: FALLA en la marca y en K38.
  - V3, ajuste de la casa en 100 kcal: FALLA K44.
  - V4, sección P-26.15: FALLA U1-08 y la matriz.
  - V5, umbral del 85 %: FALLA K44.

  Las 5 alteraciones salen con 1.
- **Valores inventados**: `grep` de `365`, `3 intentos`, `≥ 1`, «no menos de» y «como máximo n» no encuentra ningún tope en U1. Los «≥ 0» de N-326, N-327 y `advance_days` salen de K36, K43 y 00003 (línea 108). Las validaciones de coherencia (por ejemplo, ≥ 1) solo están en «Propuestas para Axel» de `U1-1.2.md`.

## Criterios
| Criterio | Estado | Evidencia |
|---|---|---|
| (a) Verificador corre; PENDIENTE = N nuevas de «Para C1» | CUMPLE | Sale con 0 y TODO OK. PENDIENTE = N-320 a N-327, igual que las «Nuevas» de «Para C1 (corregida r2)», según el cruce propio con openapi |
| (b) K38 FIRME en §0.2 y §7, sin marcas | CUMPLE | §0.2, línea 82: «la que empieza el lunes (**FIRME**, K38; «Cambiar semana» siempre disponible)». §7, línea 1203: FIRME, «bloqueo resuelto». V2 lo detecta si se rompe |
| (c) K36 en Salud y en Ajustes › Salud y metas: 0 y 0 y «el de la casa» | CUMPLE | P-21: §2.2 línea 270, §2.6 línea 426, N-327 línea 713 («Valor inicial: el de la casa (N-326; 0 y 0 kcal al empezar)», «Restablecer» vacía el campo). P-26.13: línea 295 y N-326, línea 811 (0 kcal y 0 kcal). N-230 y N-311: «menos el ajuste de bajar o más el de subir; «mantener» no ajusta; RN-17». §7: línea 1207 |
| (d) K44: cada parámetro en su sección, con valor inicial, validación y «Restablecer»; sin sección 15 | CUMPLE | §2.3, líneas 300 a 325: 4 reglas comunes y 15 filas. Las secciones son las de la lectura PROPUESTA de K44: alertas en Notificaciones, frescura en Tiendas y precios, acceso en Cuenta y roles, comparador en Créditos, adjuntos en Datos y kcal en Salud. Hay una N-xxx E por grupo (N-320 a N-323, N-325, N-326) con «valor inicial» y «restablece». 14 secciones, P-26.15 = 0 |
| (e) K43: N-023, N-064, N-226, N-297 y N-300 FIRMES | CUMPLE | Líneas 464, 516, 701, 790 y 793 con «**FIRME**, K43»: 10 MB, 5 s, «Comparador no disponible», `send_time` y `advance_days` |
| (f) 3 toques intactos; conteos de §3.29 | CUMPLE | §4 no aparece en el diff. F-1, líneas 905 a 919: 1 + 1 + 1 = 3. Recuento propio: 211 / 138 / 73 y matriz biyectiva |
| (g) PROVISIONAL solo con bloqueos abiertos | CUMPLE | No hay marcas en el cuerpo. Las 4 del historial ya estaban en la ronda 2 (H-05, baja) |
| (h) «Para C1» completa y coherente | CUMPLE (ver H-r3-01, baja) | «Para C1 (corregida r2)»: 8 nuevas con claves y valores iguales a 00003 y a U1. N-323 con las tres `alerts.*`. Viñeta N-297/298/300 coherente con N-298, N-300 y N-323. El verificador y la fila 9 de cambios 1.2 apuntan ahí |
| (i) Ningún valor inventado | CUMPLE | Ver «Valores inventados» arriba. Hay 2 parámetros sin valor (`send_time`, `advance_days`), declarados así y sin «Restablecer» |
| U1-01, U1-02, U1-05, U1-07, U1-10, U1-11, U1-13 | CUMPLE | No están tocados por el diff. Líneas OK del verificador |
| U1-03 (alta) | CUMPLE con PENDIENTE esperado | 0 N-xxx «reales» sin operación. Las 8 nuevas quedan PENDIENTE hasta C1 v1.2, como pide el encargo (a). 0 operaciones sin N-xxx |
| U1-04 (alta) | CUMPLE | Regla K10 escrita; F-1 = 3 |
| U1-08 (alta) | CUMPLE | 14 de 14, en orden y con nombres exactos |
| U1-09 (alta) | CUMPLE | K18 OK. Cubiertos K36 (P-21, P-26.13) y K44 (§2.3) |
| U1-12 (media) | CUMPLE | Claves = 00003 y RN-29. `getMealPlanWeek` usa `period` en C1 (`extractos/comidas.yaml`), igual que en U1 §0.2 |
| U1-14 (media) | CUMPLE | `U1-1.2.md` §Supuestos (5) y §Bloqueos. Solo queda un bloqueo abierto que afecta a U1, `S1-parametros-avisos`, y está citado |
| U1-15, U1-16 (media) | CUMPLE | U1-16 1.2: regla del domingo FIRME y sin marca |

## Hallazgos (todos bajos)

### H-r3-01 (baja, resto de H-r2-01) · «Para C1» de `U1-1.2.md` sigue con 2 viñetas viejas
- **Ruta:** `docs/estado/p1c/U1-1.2.md`, líneas 41 a 60: la viñeta de N-323 y la línea 60, «hoy no tienen valor sembrado».
- **Qué pasa:** la sección corregida está en `correccion-U1-r2.md`, y U1 (fila 9) y el verificador mandan ahí. Aun así quedan dos «Para C1» distintas en la misma carpeta, y quien abra `U1-1.2.md` sin pasar por el verificador lee la versión vieja.
- **Corrección sugerida:** que el integrador aplique la solicitud 1 de `correccion-U1-r2.md`, o que ponga al principio de «Para C1» de `U1-1.2.md` la nota «Sustituida por "Para C1 (corregida r2)" de correccion-U1-r2.md». Al aplicarla, el verificador cambia solo de fuente (prueba B de correccion-U1-r2).

### H-r3-02 (baja, nuevo) · N-300 da por hecho `params_json`
- **Ruta:** U1, línea 793 (N-300), frente a las líneas 791 (N-298), 1209 (§7) y 30 (fila 9).
- **Qué pasa:** N-300 dice «`advance_days` … y `send_time`, en `params_json`». N-298 y §7 dicen que la ubicación la decide el bloqueo abierto `S1-parametros-avisos`, y S1 la marca PROVISIONAL (00003, líneas 53 y 54). No cambia la pantalla ni la necesidad.
- **Corrección sugerida:** «en `params_json` (K43; la ubicación final la decide el bloqueo abierto `S1-parametros-avisos`)», igual que «Para C1 (corregida r2)».

### H-r3-03 (baja, nuevo, y H-06 que sigue) · Redacción
- **Ruta:** U1, línea 791: «…K43; K44). los días antes…», con minúscula después de punto. Línea 763 (N-320): sigue faltando la coma en «(valor inicial 15 min) duración de la sesión» (H-06). Línea 322: «claves de idempotencia» en un rótulo que ve Ximena (H-06).
- **Corrección sugerida:** escribir «Los» con mayúscula, poner la coma y proponer a U4 un rótulo llano.

### Bajos que vienen de la ronda 2 y siguen abiertos (sin cambio en 1.2)
- **H-r2-02**: `send_time` «de cada reporte» (línea 317) frente a 00003 (línea 53: «resumen diario y avisos con hora»). Corrección: que U1, S1 y S7 digan lo mismo.
- **H-02**: la columna «Decisión» de las líneas 309 y 318 no nombra A ni F, como sí hace 00003 («Decisión A, K11»; «Decisión F, K44»). Corrección: sumarlas.
- **H-03**: la línea 1208 de §7 dice «**Decidido** (K43) … `send_time` y `advance_days`»; K43 decidió las claves, no el valor. Corrección: «claves decididas (K43); valor sin decidir, lo propone S7 (K31)».
- **H-04**: las filas 312 (contraseña), 316 (PROFECO) y 322 (idempotencia) son ubicaciones que K44 no da y no dicen que son propuesta. Corrección: marcarlas «(propuesta S1 y U1, puerta 2)».
- **H-05**: el historial de la 1.1 (líneas 39 y 52) dice «queda PROVISIONAL». Corrección: añadir «(en v1.1; FIRME en v1.2)».
- **H-07**: `weekday` 7 y `day_of_month` 1 de 00001 no figuran como parámetros de K44 ni se pueden editar en N-297 y N-300. Corrección: que S1 confirme si son parámetros de K44 o lo pregunte en un bloqueo.

## Estado de los hallazgos anteriores
| Hallazgo | Estado | Evidencia |
|---|---|---|
| H-r2-01 (media), resto de H-01 r1: «Para C1» completa y coherente; valor inicial y «Restablecer» (K44) | **Cerrado** (por la «otra salida» de la ronda 2), con resto bajo → H-r3-01 | La línea PENDIENTE y la comprobación H-r2-01 apuntan a «Para C1 (corregida r2)». Esa sección expone [1, 0], 80 y 100 para «Restablecer» por N-323, dice que `send_time` y `advance_days` no tienen valor ni «Restablecer», y cita el bloqueo abierto. La V1 negativa demuestra que la viñeta vieja da FALLA. La fila 9 de los cambios 1.2 de U1 nombra la sección corregida |
| H-01 (media) r1 | Cerrado (ya en r2) | N-298, N-300, N-323 y las filas 318 y 319 usan `alerts.*` de 00003 |
| H-r2-02, H-02, H-03, H-04, H-05, H-06 y H-07 (bajos) | Siguen abiertos (bajos) | Ver arriba. No bloquean la aprobación |
