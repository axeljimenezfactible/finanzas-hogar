# Revisión U1 1.2 · ronda 2 (flujo P1c, 11-oct-2026)

Revisor: **revisor-U1-r2** (independiente; no escribí nada de lo revisado y no corrijo). Evalúo solo contra `docs/criterios/U1.md` 1.2 y `decisiones.md` §K (K35, K36, K38, K43 y K44). Reejecuté todo con mis propios scripts en `datos/tmp/revisor-u1-r2/` (ya borrados). Revisé `git diff HEAD -- docs/design/U1-arquitectura.md`, el verificador `docs/design/tools/verificar_u1.py`, el reporte `docs/estado/p1c/U1-1.2.md` y la corrección `docs/estado/p1c/correccion-U1-r1.md`.

## Veredicto

**NO APROBADO**: 0 hallazgos altos, **1 medio** y 7 bajos.

`U1-arquitectura.md` y el verificador ya corrigieron H-01: el pago próximo usa `days_before` [1, 0] y `advance_days` queda solo para precios viejos. Lo comprobé por mi cuenta. El medio que queda (H-r2-01) es la otra mitad de H-01: la sección «Para C1» de `docs/estado/p1c/U1-1.2.md` sigue con la viñeta vieja de N-297, N-298 y N-300. La versión corregida está en `correccion-U1-r1.md` como solicitud 1, pero nadie la aplicó. El verificador manda a C1 justo a ese archivo. Para cerrarla basta cambiar una viñeta. Todo lo demás cumple.

## Tabla de criterios

| Criterio | Estado | Evidencia (propia) |
|---|---|---|
| (a) Verificador corre; PENDIENTE = N nuevas de «Para C1» | CUMPLE | `python3 -I docs/design/tools/verificar_u1.py`: sale con 0, 71 líneas, «RESULTADO: TODO OK» y una sola línea PENDIENTE: N-320 a N-327. Mi script cruza las 211 N de U1 con `x-necesidades` de `openapi.yaml` 1.1.0 (200 N con operación; Local: N-003, N-301 y N-317). Sin operación ni Local quedan **exactamente** N-320 a N-327. Son las 8 «Nuevas» de «Para C1» en `U1-1.2.md` (las extraje con un script). No hay N de openapi que falten en U1 ni N Local con operación |
| (b) K38 FIRME en §0.2 y §7, sin marcas | CUMPLE | §0.2 línea 81: «(**FIRME**, K38; «Cambiar semana» siempre disponible)». §7 línea 1202: FIRME (K38), con el bloqueo `U1-lista-vigente-domingo` como resuelto (su archivo tiene «Respuesta de Axel»). Mi búsqueda no encuentra ninguna marca PROVISIONAL en el cuerpo |
| (c) K36: Salud y Ajustes » Salud y metas, 0 y 0, «el de la casa» | CUMPLE | P-21 en §2.2 (línea 269) y §2.6 (línea 425). N-327 (línea 712): «Valor inicial: el de la casa (N-326; 0 y 0 kcal al empezar)», vacío = el de la casa, «Restablecer» vacía el campo. P-26.13 (línea 294) y N-326 (línea 810): «valores iniciales 0 kcal y 0 kcal; K36, K44». N-230 y N-311: «mantener» no ajusta (RN-17). Coincide con 00003: líneas 58-59 (columnas por persona, ≥ 0) y 67-68 (claves de la casa en 0) |
| (d) K44: sección propuesta, valor inicial, validación, «restablecer», sin sección 15 | CUMPLE | Mi recuento de la primera tabla de §2.3 da 14 secciones, P-26.01 a P-26.14. La tabla «Parámetros operativos» (líneas 308-322) cubre las 12 claves `@param` de `00003_puerta1b.sql` (líneas 3-14 del encabezado) en la sección que propone K44: alertas en Notificaciones, frescura en Tiendas y precios, intentos y sesión en Cuenta y roles, comparador en Créditos, adjuntos en Datos y respaldos, kcal en Salud y metas. Suma `days_before` [1, 0], `send_time` y `advance_days`. Las reglas comunes (líneas 300-304) cubren valor inicial visible, «Sin valor inventado», «Restablecer al valor inicial» y «Sin valor decidido». N-320 a N-323, N-325 y N-326 remiten a ellas |
| (e) K43: N-023, N-064, N-226, N-297, N-300 FIRMES | CUMPLE | N-023 (línea 463): «Definición **FIRME** (K43)… es una regla, no un parámetro». N-064: «10 MB (**FIRME**, K43)». N-226: «5 s (**FIRME**, K43…)» con «Comparador no disponible». N-297 y N-300: claves `send_time` y `advance_days` FIRMES. Ninguna lleva marca |
| (f) 3 toques intactos; conteos §3.29 | CUMPLE | Ningún bloque del diff toca §4. En F-1 (líneas 904-907) los toques acumulados van 1, 2 y **3**. Mi recuento: 211 N (138 L + 73 E), sin duplicados, la última es N-327 y no hay referencias a N inexistentes. La matriz §3.29 tiene 42 filas, todas con las mismas columnas, y contiene las 211 N (expandí los rangos «a»). Pantallas: 27 P-xx y 14 P-26.xx. Las cifras declaradas en las líneas 18 y 816 coinciden |
| (g) PROVISIONAL restantes citan solo bloqueos abiertos | CUMPLE (ver H-05) | El cuerpo no tiene marcas. Hay 4 apariciones en las listas de cambios (líneas 22, 37, 38 y 51) y 2 usos comunes de la palabra (líneas 249 y 328). Bloqueos abiertos (sin «Respuesta de Axel»): `S1-parametros-avisos` y `criterios-nxx`. Crucé los 13 bloqueos que cita U1: `S1-parametros-avisos` aparece como «abierto» (líneas 790, 794 y 1208) y los demás como «resuelto», lo que coincide con su archivo |
| (h) «Para C1» completa y coherente con las N nuevas | **NO CUMPLE → H-r2-01** | Las 8 nuevas tienen claves y valores iguales a 00003 y a 00001. La viñeta de N-297, N-298 y N-300 de `U1-1.2.md` (línea 60) contradice a U1 1.2 corregida (ver H-r2-01). La viñeta corregida de `correccion-U1-r1.md` sí es correcta: comparé las dos secciones con `diff` y solo difiere esa línea |
| (i) Ningún valor inventado | CUMPLE | Saqué por script las cifras nuevas del diff: 0 y 0 kcal, [1, 0], 1 GB, 10 MB, 12 caracteres, 15 min, 24 h, 5 s, 30 y 90 días. Todas vienen de K6, K11, K22, K36, K40, K43 o K44, o de la decisión F (00001 línea 1244). Ningún tope: «≥ 1», «365», «no menos» y «como mínimo» dan 0. «≥ 0» sale de K36 (CHECK de 00003) y de K43. `send_time` y `advance_days` dicen «Sin valor inicial decidido (K31)», como la propuesta aprobada (`C1-valores-k31` línea 6: «los valores por defecto los propone S7») y 00003 líneas 45-46 |
| U1-01, U1-02, U1-08 | CUMPLE | 9 entradas en la barra. 27 P-xx únicas y 14 P-26.xx con nombre y orden exactos (verificador y mi script) |
| U1-03 | CUMPLE (PENDIENTE C1 v1.2) | Ver (a). Solo faltan las 8 nuevas, que esperan a C1 v1.2 |
| U1-04, U1-05, U1-15 | CUMPLE | Sin cambios en §4. El verificador da OK en F-1 a F-8, con las variantes «Sin conexión» de F-1 y F-6 |
| U1-13 | CUMPLE (ver H-06) | Nombres en español. El verificador revisa 77 nombres visibles y no encuentra palabras en inglés |
| U1-14 | CUMPLE (ver H-04) | §7 lleva los temas resueltos y el bloqueo abierto `S1-parametros-avisos` |
| U1-16 | CUMPLE | Una sola definición (línea 81). La regla del domingo solo aparece en §0.2 y §7 |
| K35 | CUMPLE | §7 línea 1209: la Pista 3 la ejecuta Axel con el paquete de `docs/design/entrada/`, que se entrega al cierre de WF-1. U1 no entra como archivo en el paquete |

## Hallazgos

### H-r2-01 (media) · «Para C1» de `U1-1.2.md` sigue con la viñeta vieja de N-297, N-298 y N-300 (resto de H-01)
- **Ruta:** `docs/estado/p1c/U1-1.2.md` línea 60. Relacionadas: `docs/design/U1-arquitectura.md` líneas 790, 792 y 318, y `docs/design/tools/verificar_u1.py`, cuya línea PENDIENTE dice «van en «Para C1» de docs/estado/p1c/U1-1.2.md».
- **Qué pasa:** la viñeta dice que las tres N «Llevan su valor inicial para «Restablecer» cuando lo haya (K44); hoy no tienen valor sembrado (RN-29 1.3)». U1 1.2 corregida dice otra cosa:
  - la anticipación del pago próximo es `days_before`, con valor inicial [1, 0] sembrado (00001 línea 1244, decisión F);
  - los umbrales `thresholds_pct` [80, 100] también están sembrados (00001 línea 1245).

  K44 exige mostrar ese valor inicial y ofrecer «Restablecer». C1 v1.2, que lee `U1-1.2.md` por indicación del verificador, recibiría la instrucción vieja. La viñeta corregida ya existe en `docs/estado/p1c/correccion-U1-r1.md`, sección «Para C1 (corregida)», con la solicitud 1, pero nadie la aplicó: `U1-1.2.md` no cambió desde las 15:06 y la corrección es de las 15:18.
- **Corrección sugerida:** el integrador cambia la línea 60 de `U1-1.2.md` por la viñeta de «Para C1 (corregida)» de `correccion-U1-r1.md`. El `diff` muestra que es la única línea distinta y la revisé contra 00001 y 00003. Otra salida: que la línea PENDIENTE de `verificar_u1.py` apunte a la sección corregida.

### H-r2-02 (baja, nuevo) · Alcance de `send_time` en los reportes
- **Ruta:** U1 línea 316 («Hora de cada reporte (`send_time`)») y N-297 (línea 789). Contra: `00003_puerta1b.sql` línea 45 y `C1-valores-k31` línea 6.
- **Qué pasa:** U1 pone `send_time` en cada reporte (diario, semanal y mensual). S1 lo define para «el resumen diario y los avisos con hora», y la propuesta que aprobó K43 lo define solo «para el resumen diario». Ninguna cifra depende de esto, porque no hay valor decidido.
- **Corrección sugerida:** que U1 y S1 (00003 y RN-29) digan lo mismo: o S1 extiende la descripción a los reportes semanal y mensual, o U1 aclara qué reportes llevan `send_time`. Si hay dos lecturas, que S7 lo incluya en su propuesta.

### H-02 (baja, sigue) · Columna «Decisión» de la tabla K44
- **Ruta:** U1 líneas 308 (5 intentos: solo K44), 317 (80 % y 100 %: solo K44) y 318 (`days_before`: K44 en la columna; la decisión F solo aparece en la columna del valor).
- **Corrección sugerida:** sumar «A» o «F» en la columna «Decisión», como hace RN-29.

### H-03 (baja, sigue) · §7 dice «Decidido» para `send_time` y `advance_days`
- **Ruta:** U1 líneas 1207 y 1211.
- **Qué pasa:** la línea 1207 sigue diciendo «hora y anticipación de avisos… **Decidido** (K43)… `send_time` y `advance_days`». K43 decidió las claves, no los valores, y desde H-01 la anticipación del pago próximo es `days_before` (decisión F), no `advance_days`. La línea 1211 deja a U4 los valores por defecto, que propone S7 (K31).
- **Corrección sugerida:** en la 1207, «claves decididas (K43); valor inicial sin decidir, lo propone S7 (K31); pago próximo: `days_before` [1, 0], decisión F». En la 1211, decir que U4 muestra lo que proponga S7.

### H-04 (baja, sigue) · Tres ubicaciones que K44 no da aparecen como lectura del orquestador
- **Ruta:** U1 línea 299 y filas 311 (largo de contraseña), 315 (corte PROFECO) y 321 (retención de idempotencia).
- **Corrección sugerida:** marcar esas filas «(propuesta S1 y U1, puerta 2)» y sumar la propuesta en §7.

### H-05 (baja, sigue) · El historial de la 1.1 dice «queda PROVISIONAL» y cita un bloqueo resuelto
- **Ruta:** U1 líneas 38 (`U1-lista-vigente-domingo`, ya resuelto) y 51.
- **Qué pasa:** son filas de historial, no marcas vigentes (la fila 1 de la 1.2 dice que se quitaron), pero un `grep` de marcas las cuenta.
- **Corrección sugerida:** añadir «(en v1.1; FIRME en v1.2, K38/K43)» o quitar la palabra en mayúsculas.

### H-06 (baja, sigue) · Redacción
- **Ruta:** U1 línea 762 (N-320: falta la coma en «(valor inicial 15 min) duración de la sesión») y línea 321 («claves de idempotencia» en un rótulo que ve Ximena).
- **Corrección sugerida:** poner la coma. Proponer a U4 un rótulo llano, con el término técnico en la ayuda.

### H-07 (baja, sigue; con evidencia nueva) · Días de los reportes y del aviso «plan sin armar»
- **Ruta:** U1 N-297 (línea 789), N-300 (línea 792) y P-26.08 (línea 289). Contra: `00001_esquema.sql` líneas 1242, 1243 y 1247.
- **Qué pasa:** en la misma `notification.params_json` donde está `days_before`, 00001 siembra `weekday` 7 (resumen semanal y «plan sin armar») y `day_of_month` 1 (reporte mensual). Desde la corrección de H-01, U1 trata `days_before` como parámetro K44, con valor inicial y «Restablecer». Pero N-297 y N-300 no permiten editar el día ni muestran su valor inicial. El inventario de 00003 tampoco los incluye.
- **Corrección sugerida:** que S1 confirme si son parámetros de K44 o lo pregunte en un bloqueo. Si lo son, U1 los suma a N-297 y N-300 con valor inicial (domingo y día 1, decisión F) y «Restablecer».

## Estado de los anteriores

| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H-01 (media) | **Cerrado en U1 y en el verificador; abierto en «Para C1» → H-r2-01** | N-298 (línea 790) y N-300 (línea 792) nombran `days_before` con [1, 0] y la decisión F. La fila 318 de §2.3 cita la clave y su valor. La fila 319 deja `advance_days` solo para precios viejos (como 00003 línea 46). Las filas 2 y 8 de los cambios 1.2 lo registran. La comprobación nueva del verificador lee 00001 y 00003 por su cuenta, y yo la repetí con mis scripts. En cambio, `U1-1.2.md` línea 60 no cambió |
| H-02 a H-07 (bajas) | Siguen abiertas | No estaban en el encargo de `correccion-U1-r1` (lo dice su reporte). Las rutas actualizadas están arriba |
