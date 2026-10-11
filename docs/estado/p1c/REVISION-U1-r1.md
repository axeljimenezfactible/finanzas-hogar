# Revisión U1 1.2 · ronda 1 (flujo P1c, 11-oct-2026)

Revisor independiente: **revisor-U1-r1**. No escribí nada de lo revisado y no corregí nada. Evalué solo contra `docs/criterios/U1.md` 1.2 y `docs/plan/decisiones.md` K35, K36, K38, K43 y K44. Todas las comprobaciones las volví a correr con mis propios scripts en `datos/tmp/rev-u1-r1/`, que ya borré.

## Veredicto

**NO APROBADO**: 0 hallazgos altos, **1 medio** y 6 bajos.

El único hallazgo medio (H-01) está en la anticipación del aviso de pago próximo. U1 la asigna a `advance_days` (un solo entero), pero el esquema la guarda como `days_before` = [1, 0]. Además, «Para C1» no pide que C1 exponga ese valor inicial para «Restablecer». Todo lo demás cumple.

## Tabla de criterios

| Criterio | Estado | Evidencia (propia) |
|---|---|---|
| (a) Verificador corre; PENDIENTE = nuevas de «Para C1» | CUMPLE | `python3 -I docs/design/tools/verificar_u1.py`: sale con 0, «RESULTADO: TODO OK» y 1 línea PENDIENTE con N-320 a N-327. Mi script cruza las 211 N-xxx de U1 con `x-necesidades` de `openapi.yaml`. Sin operación ni Local quedan exactamente N-320 a N-327, que son las 8 de «Para C1» en `docs/estado/p1c/U1-1.2.md`. No hay N-xxx en openapi que falten en U1 |
| (b) K38 FIRME en §0.2 y §7, sin marcas | CUMPLE | §0.2 línea 80: «**FIRME**, K38; «Cambiar semana» siempre disponible». En §7, la fila «Semana por defecto los domingos» está FIRME y su bloqueo figura como resuelto. En el cuerpo no queda ninguna marca PROVISIONAL (ver g) |
| (c) K36: Salud y Ajustes » Salud y metas, 0 y 0, «el de la casa» | CUMPLE | N-326 (P-26.13): «valores iniciales 0 kcal y 0 kcal; K36, K44». N-327 (§3.23 P-21 Salud y ficha de P-26.13): «Valor inicial: el de la casa», vacío = el de la casa, «Restablecer» vacía el campo. N-230 y N-311 aplican RN-17 FIRME («mantener» no ajusta). La matriz tiene N-327 en P-21 y en P-26.13 |
| (d) K44: sección propuesta, valor inicial, validación, «restablecer», sin sección 15 | CUMPLE con H-01 | La tabla de §2.3 (líneas 298-322) tiene 15 filas en P-26.01, .05, .06, .08, .09, .11 y .13, que coinciden con la lectura del orquestador. Siguen siendo 14 secciones (mi recuento de la primera tabla de §2.3 da 14, con nombres y orden exactos). Las reglas comunes están escritas: valor inicial visible, «Sin valor inventado», «Restablecer al valor inicial» y «Sin valor decidido». Cada necesidad de grupo (N-320 a N-323, N-325 y N-326) remite a esas reglas. Los 12 valores de K44 y RN-29 coinciden con las claves que siembra `00003_puerta1b.sql` |
| (e) K43: N-023, N-064, N-226, N-297, N-300 FIRMES | CUMPLE | N-023 «Definición **FIRME** (K43)… es una regla, no un parámetro»; N-064 «10 MB (**FIRME**, K43)»; N-226 «5 s (**FIRME**, K43)» con «Comparador no disponible»; N-297 y N-300 con `send_time` y `advance_days` FIRMES. Ninguna lleva marca |
| (f) 3 toques intactos; conteos §3.29 | CUMPLE | Ningún bloque del diff toca §4 (F-1, líneas 899-922). La tabla de F-1 da 1 + 1 + 1 = **3**. Mi recuento: 211 N-xxx (138 L + 73 E), sin duplicados, la última es N-327; las 211 están en la matriz §3.29 (N-001 a N-006 en la fila «Global»). Pantallas: 27 en §2.1 y 14 secciones P-26.xx. Las cifras declaradas en las líneas 18 y 815 coinciden. Las tablas del documento no tienen filas con columnas de más o de menos |
| (g) PROVISIONAL restantes citan solo bloqueos abiertos | CUMPLE (ver H-05) | `grep`: el cuerpo no tiene marcas. Hay 4 apariciones en las listas de cambios (líneas 22, 36, 37 y 50) y 2 usos comunes de «provisional» (líneas 248 y 327, «nunca cifras provisionales»). Bloqueos abiertos hoy: `S1-parametros-avisos` y `criterios-nxx`. N-323 cita el primero sin marca, y es correcto: la pantalla y el valor no dependen de él, solo el lugar donde se guarda |
| (h) «Para C1» completa y coherente con las N nuevas | PARCIAL → H-01 | Las 8 N nuevas están con sus claves, que coinciden con RN-29 y con `00003_puerta1b.sql`. Están los cambios de N-230, N-311, N-297, N-298, N-300, N-064, N-226, N-023, N-120, N-299 y N-289. Falta el valor inicial de `days_before` del pago próximo para «Restablecer» (H-01) |
| (i) Ningún valor inventado | CUMPLE | Las cifras de las líneas añadidas son 5, 15 min, 90 días, 12 caracteres, 5 s, 30 días, 90 días, 1 GB, 80 % y 100 %, 1 día, 10 MB, 24 h y 0 kcal. Todas vienen de §A, §F, K6, K11, K22, K36, K40, K43 o K44. No hay topes mínimos ni máximos: «≥ 1», «365» y «no menos» dan 0 coincidencias. `send_time` y `advance_days` sin valor dicen «Sin valor inicial decidido (K31)» |
| U1-01, U1-02, U1-08 | CUMPLE | Barra lateral de 9 entradas. 27 pantallas P-xx únicas y 14 P-26.xx en orden |
| U1-03 | CUMPLE (PENDIENTE C1 v1.2) | Las 3 listas del verificador están vacías, salvo las 8 nuevas, que esperan a C1 v1.2 |
| U1-04, U1-05, U1-15 | CUMPLE | No cambian (diff). El verificador da OK en F-1 a F-8 y en «Sin conexión» |
| U1-07, U1-10, U1-11, U1-12 | CUMPLE | No cambian. El verificador da OK |
| U1-13 | CUMPLE (ver H-06) | Los nombres de pantallas y secciones no cambian. Los nombres nuevos están en español |
| U1-14 | CUMPLE (ver H-04) | §7 lleva los temas resueltos y el bloqueo abierto. U1 decide «Descartar chat detectado» por encargo de `INTEGRACION.md` §7, punto 7 |
| U1-16 | CUMPLE en U1 | Una sola definición (línea 80). Que la regla quede FIRME en `getMealPlanWeek` depende de C1 v1.2 |
| K35 | CUMPLE | La fila de §7 dice que la Pista 3 la ejecuta Axel con `docs/design/entrada/` al cierre de WF-1 y que U1 no entra como archivo en el paquete, lo que coincide con la lista de K35 |

## Hallazgos

### H-01 (media) · Anticipación del pago próximo: `advance_days` en U1, `days_before` [1, 0] en el esquema
- **Ruta:** `docs/design/U1-arquitectura.md` líneas 789 (N-298), 791 (N-300) y 317-318 (tabla K44); «Para C1» en `docs/estado/p1c/U1-1.2.md` (viñeta N-297, N-298 y N-300).
- **Qué pasa:** N-298 y N-300 dicen que la anticipación de los avisos es `advance_days` («días enteros ≥ 0»). Pero en el esquema el pago próximo usa `days_before` = [1, 0] en `notification.params_json` (`00001_esquema.sql` línea 1244, RN-29). Además, `00003_puerta1b.sql` línea 46 y la propuesta que aprobó K43 (`C1-valores-k31`, punto 3) definen `advance_days` como la anticipación del aviso de precios viejos. Un solo entero no puede guardar «día anterior y mismo día» (decisión F). La propia tabla de §2.3 separa el pago próximo de «los demás avisos (`advance_days`…)», pero no nombra su clave. «Para C1» dice de N-297, N-298 y N-300 que «hoy no tienen valor sembrado» y no pide exponer el valor inicial [1, 0] de `days_before` para «Restablecer», que K44 exige porque ese valor sí existe.
- **Corrección sugerida:** en N-298 y N-300 nombrar `days_before` (lista de días, valor inicial [1, 0], decisión F) para el pago próximo y dejar `advance_days` solo para precios viejos, como en 00003. En la fila 317 citar la clave. En «Para C1», pedir que `listNotifications` y `updateNotification` devuelvan el valor inicial de `days_before` y de `thresholds_pct` para «Restablecer», sujeto a la opción del bloqueo abierto `S1-parametros-avisos`.

### H-02 (baja) · Columna «Decisión» de la tabla K44 incompleta
- **Ruta:** U1 líneas 307, 316 y 317.
- **Qué pasa:** 5 intentos, 80 % y 100 % y «1 día antes y el mismo día» citan solo K44 (o K11). Esos valores ya estaban decididos en §A y §F (decisiones línea 41), y RN-29 los cita así («A, K11»; «F»).
- **Corrección sugerida:** sumar «A» o «F» en la columna «Decisión» de esas filas.

### H-03 (baja) · §7 dice «Decidido» para `send_time` y `advance_days`, y deja a U4 sus valores por defecto
- **Ruta:** U1 líneas 1206 y 1210.
- **Qué pasa:** la línea 1206 dice «**Decidido** (K43): … `send_time` y `advance_days`», pero K43 decidió las claves, no sus valores. La línea 1210 deja a U4 los «valores por defecto (hora del resumen diario, anticipación de alertas)». Según K31, RN-29 y `00003` (línea 46), esos valores los propone S7 y los aprueba Axel.
- **Corrección sugerida:** en la línea 1206 decir «claves decididas; valor inicial sin decidir, lo propone S7 (K31)». En la línea 1210 decir que U4 muestra la propuesta que haga S7.

### H-04 (baja) · Ubicación de tres parámetros que K44 no ubicó: aparece como lectura del orquestador
- **Ruta:** U1 línea 298 y filas 310, 314 y 320 (largo de contraseña, corte PROFECO y retención de idempotencia).
- **Qué pasa:** el párrafo atribuye toda la tabla a la «lectura PROPUESTA del orquestador en K44». La lectura de K44 no ubica esos tres. Los ubican S1 y U1, y RN-29 los marca «(propuesta S1)». Tampoco están en §7 como supuesto de U1 (U1-14).
- **Corrección sugerida:** marcar esas tres filas «(propuesta S1 y U1, puerta 2)» y sumar en §7 una fila con esa propuesta y con la decisión de U1 sobre «Descartar chat detectado».

### H-05 (baja) · El historial de la 1.1 sigue diciendo «queda PROVISIONAL»
- **Ruta:** U1 líneas 37 y 50.
- **Qué pasa:** son filas de la lista de cambios de la 1.1 y no son marcas vigentes. La línea 37 cita un bloqueo que ya está resuelto (`U1-lista-vigente-domingo`). Un `grep` de marcas las cuenta.
- **Corrección sugerida:** añadir «(en v1.1; FIRME en v1.2, K38/K43)» o dejar la palabra sin negritas para que no se lea como marca.

### H-06 (baja) · Redacción
- **Ruta:** U1 línea 761 (N-320) y línea 320.
- **Qué pasa:** a N-320 le falta una coma: «(valor inicial 15 min) duración de la sesión». En la línea 320, «claves de idempotencia» es jerga en un rótulo que ve Ximena (U1-13 pide no usar jerga innecesaria).
- **Corrección sugerida:** poner la coma. Proponer a U4 un rótulo llano, por ejemplo «Tiempo para no duplicar capturas reenviadas», con el término técnico en la ayuda.

### H-07 (baja) · Días de los reportes y avisos fuera del inventario K44
- **Ruta:** U1 línea 288 (P-26.08) y tabla de §2.3; `docs/contratos/reglas-negocio.md` RN-29 (S1).
- **Qué pasa:** el día 1 del reporte mensual y el domingo del semanal y del aviso «plan sin armar» (decisión F) no aparecen como parámetros configurables ni en U1 ni en RN-29. K44 dice «todo valor numérico operativo». U1 ya deja dicho que suma lo que S1 inventaríe, así que el inventario es de S1.
- **Corrección sugerida:** que S1 confirme si son parámetros de K44 (o lo pregunte en un bloqueo). Si lo son, U1 los ubica en P-26.08 y P-26.09.

## Estado de los anteriores
Es la ronda 1 del flujo P1c: no hay hallazgos anteriores de este flujo. Los hallazgos de las rondas del flujo P1 (r1 a r3) quedaron cerrados en U1 1.1. El verificador sigue dando OK en los criterios que los cubrían: R2-05 (44 px), R3-01 (lista vigente) y K18 (falta capturar).
