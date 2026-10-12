# REVISION-U4-r4 · Ajustes estilo Notion (partes a y b) · ronda 4

Revisor independiente `revisor-U4-r4` · 11-oct-2026 · Evalúa solo contra `docs/criterios/U4.md` 1.2 y `docs/plan/decisiones.md` §A–§K. Artefactos: `docs/design/U4-ajustes-a.md` 1.3 (405 líneas) y `docs/design/U4-ajustes-b.md` 1.2 (388 líneas). No escribí ni corregí nada de lo revisado. Todas las verificaciones son mías, con scripts propios `python3 -I` en `datos/tmp/revU4r4/`, que borré al terminar.

## Veredicto
**NO APROBADO**: 0 hallazgos altos, 1 medio (R4-1) y 7 bajos. El hallazgo medio de la ronda 3 (R3-1) quedó **cerrado**. Pero la misma vista previa del gasto fijo deja fuera dos caminos que cambian el mismo dinero: el interruptor «activo» y la edición múltiple con `Ctrl+/`. En P-26.13 pasa lo mismo con el interruptor «Activa» de una meta (R4-1).

## Estado de los hallazgos de la ronda 3
| Hallazgo | Estado | Evidencia propia |
|---|---|---|
| R3-1 (media, U4-05): gasto fijo sin vista previa de dinero | **Cerrado** | «Edición en el lugar» de P-26.02 nombra los 8 campos y la vista previa 3 trae (a) a (f). Lo comprobé: la regla de fechas sale de los comentarios de `planned_item` en 00001 (líneas 530 a 548: ancla recortada a fin de mes, `fortnightly` en los días 15 y 30, `once` solo en `start_date`). Las partes siguen RN-03 y K14 (HALF_UP(monto × `pct_bp` ÷ 10000), y la persona del resto toma lo demás). Mi cálculo: 19,900 × 5,000 ÷ 10,000 = 9,950, así que Karen $99.50 y Hogar $99.50. `remainder_person_id`, `pct_bp`, `fixed_cents`, `second_half_amount_cents` y `accrual_months` existen en los extractos `ajustes-a` y `ajustes-b`. La excepción de §0.2 para las listas coincide con P-26.02. |
| R3-2 (baja): calendarios sin «próximo evento» | Sigue abierto (baja) | `grep -c 'Próximo respaldo:\|Próxima visita'` da 0 en las dos partes. |
| R3-3 (baja): restricciones alimentarias y §7 | Sigue abierto (baja) | P-26.13 edita «Restricciones alimentarias» (`updateHealthProfile`), aunque «Fronteras» de la parte a dice que se editan en P-26.07. La línea 3 de la parte b sigue diciendo «la §7 de componentes aún no existe». |
| R3-4 (baja): vistas de P-26.14 y P-26.07 contra U1 §2.5.a | Sigue abierto (baja) | U1 §2.5.a da a `inventario_item` un tablero por Estado y galería «Sí». P-26.14 agrupa el tablero por modo (como U1 §2.3) y no tiene galería. P-26.07 agrega «por origen». |
| R3-5 (baja): interruptor K21 fuera de Créditos | Sigue abierto (baja) | En el bloque P-26.05, `awk` + `grep -c 'meses 7 a 12'` da 0. |
| R3-6 (baja): encabezados y tablas con formatos distintos | Sigue abierto (baja) | La parte a usa `## P-26.01 · Cuenta y roles` y tablas de 6 columnas. La parte b usa `## P-26.09 Notificaciones` y tablas de 7 columnas. |

## Tabla de criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| U4-01 (alta) | Cumple | `secs.py`: 14 secciones, P-26.01 a P-26.14 en orden, con nombres idénticos a los de la decisión G. Las 14 traen Propósito, Vista de BD, Edición en el lugar, Vista previa, Estados, Atajos y Móvil. 0 filas con celdas vacías. |
| U4-02 (alta) | Cumple | `ops.py`: hay 210 `operationId` en los extractos, los mismos que en `openapi.yaml`. La parte a cita 98 y la b 45, con 0 inexistentes. Las 68 operaciones con etiqueta Ajustes están citadas (`login` y `logout`, que mi expresión regular no captaba, están en P-00 y P-26.01). |
| U4-03 (alta) | Cumple | Revisé a mano los ≈ 53 elementos, todos presentes. K21 está en P-26.02, como U1 §2.3 (R3-5, baja). K4 tiene candidatos, «Confirmar como» y «Descartar». K16 «Faltan restricciones» y K18 «Falta capturar» están, con el IVA PROVISIONAL `U4-iva-credito`. |
| U4-04 (alta) | Cumple | `defaults.py` encontró 140 valores por defecto con contenido. Tomé al azar 15 (semilla 20261011) y los 15 tienen su origen en §A–§K: ninguna preselección en «cerrar sesiones», mes calendario (B, K8), «de este mes en adelante» (K2), devengo de 2 meses (B), 17.27 % y 25.8685 % (C), lectura asistida apagada (E), 15 min (A, K11), alertas fijas (F), [1, 0] (F), domingo 7 (F), meta confirmada vacía (K23), ajuste vacío = el de la casa (K36), 0 kcal (K36), «Activa» al crear (estado de la base, no es un parámetro) y 24 h (K6). Además: A/B = B, IA apagada, tema según el sistema, 50/50, 5 en 15 y 7 a 12 apagado. |
| U4-05 (alta) | Cumple, con R4-1 (media) | Caso A/B: las series A y B (`calculado_A` y `calculado_B` de `model.json`, 14 valores) aparecen literales en la fila, con mínimos de $8,094.50 y $5,000.00. Mis cálculos: intereses de 3,094.50 × 25.8685 % = 800.50, total 3,895.00, pago 973.75 (igual que `despensas[0]`). Frescura: 9-nov-2026 con 0.9836, día 60 con 0.4918 y 8-ene-2027. RN-17: 1,648.75 × 1.55 = 2,555.5625, que da 2,556, y menos 500 da 2,056; IMC 23.7; S4-07: 2,098.50 − 300, que da 1,799. RN-15: 1,486.571. Hueco: R4-1. |
| U4-06 (alta) | Cumple | `misc.py`: los atajos de las dos partes están todos en U1 §5.3, sin duplicados en §0.7 ni en las convenciones de la parte b. Los comandos `/ir`, `/ingreso`, `/credito`, `/receta`, `/salir`, `/reporte`, `/respaldo` y `/tema` están en U1 §5.2. Las vistas por lista coinciden con U1 §2.5.a, salvo R3-4 (baja). La edición en el lugar es coherente con U1 §2.6 (fila P-26.01 a P-26.14: guardado automático, y A/B y restaurar confirman después de mostrar su efecto). |
| U4-07 (alta) | Cumple, con R4-2 (baja) | Catálogo V-01 a V-22 en es-MX. Tasas con 4 decimales y mostradas como se capturaron (K12). Fechas `10-oct-2026`. V-16 para las partes. |
| U4-08 (alta) | Cumple | §0.4 con `color.critical.*` (tokens existentes). Cubre: cerrar todas las sesiones, quitar un monto, eliminar tasa, promoción o receta, activar la lectura asistida, restaurar con doble confirmación, activar la IA y descartar un chat. |
| U4-09 (media) | Cumple | 0 hex, 0 `rgb(`. 21 tokens citados, todos en `C2-colores.json` o `entrada/C2-tokens.json`. Las 14 secciones tienen 412 px y 1440 px. 44 px y foco visible en §0.7. |
| U4-10 (media) | Cumple | Bloque «Estados» en las 14 y regla única en §0.3. |
| U4-11 (media) | Cumple, con R4-3 (baja) | 0 fechas ISO. Hay SUPUESTOS en las dos partes. Los 8 bloqueos citados existen. Quedan 4 líneas de la parte b con PROVISIONAL sin cita en la misma línea (104, 209, 361 y 362). |
| U4-12 (alta) | Cumple | `params.py`: las 15 claves `@param` de 00003 están en la sección que dice el esquema, con su valor inicial en la fila y «Vuelve a…». Las 20 claves sembradas están citadas. `send_time`, `advance_days`, `weekday` y `day_of_month` van con `notification.params_json`, lugar PROVISIONAL y `S1-parametros-avisos`. 00001 siembra `weekday` 7 en `weekly_summary` y `meal_plan_missing`, y `day_of_month` 1 en `monthly_report`, y U4 los da como FIRMES. `send_time` y `advance_days` dicen «sin valor inicial: lo propone S7 (K31)». «Mínimo de seguridad» solo aparece en §0.5, en las convenciones de la parte b y en SUPUESTOS, siempre negado. |
| U4-13 (alta) | Cumple | Están `health.goal_adjust_lose_kcal` y `health.goal_adjust_gain_kcal` con 0 y «Vuelve a 0 kcal». Por persona, vacío = el de la casa, y «Usar el de la casa» sirve de restablecer. «Mantener» no tiene campo. La vista previa da 1,799. Necesidades N-326 y N-327. Usa `getHealthSettings`, `updateHealthSettings` y `updateHealthProfile`. |

## Hallazgos de la ronda 4

**R4-1 · media · U4-05 · dos caminos del gasto fijo y uno de las metas cambian dinero sin vista previa.** Rutas: `docs/design/U4-ajustes-a.md`, P-26.02 (fila «Gasto fijo: activo», frase «Edición en el lugar», «Atajos»); `docs/design/U4-ajustes-b.md`, P-26.13 (fila «Activa» de metas).
- La vista previa 3 (c) define «Gastos fijos de {mes}» como «la suma de las partidas **activas**». Apagar o encender «activo» cambia esa suma, las partes por persona (d) y la proyección de P-07 (f). Aun así, el interruptor guarda al tocar y no tiene vista previa (§0.2: «Interruptores… guardan al elegir»).
- La frase «guardado automático por campo, después de su vista previa en línea… en los campos que cambian el dinero: [8 campos]» deja fuera «activo», aunque sí cambia el dinero.
- «Atajos» ofrece `Ctrl+/` para «cambiar a la vez categoría o activo de varios gastos fijos» sin vista previa. Cambiar la categoría de una partida sola, en cambio, sí pide vista previa antes de guardar, porque cambia el reparto (d). Quien implemente tendría que adivinar si la edición múltiple muestra el bloque.
- Pasa lo mismo en P-26.13: la vista previa (2) de metas dice que el total reservado del mes es la suma de `planned_monthly_cents` de las metas **activas**, pero el interruptor «Activa» guarda sin vista previa.

**Corrección sugerida:** tratar «activo» del gasto fijo y «Activa» de la meta como las listas de la excepción de §0.2: al tocar, mostrar «Gastos fijos de {mes}: {antes} → {después}», las partes (d) y el aviso (f), o «Total reservado para metas en {mes}: {antes} → {después}»; guardar con `Enter` o al salir; `Esc` revierte. Para `Ctrl+/`, mostrar antes de aplicar el mismo bloque sumado sobre las partidas seleccionadas, o sacar categoría y activo de la edición múltiple. Agregar «activo» a la lista de campos de «Edición en el lugar».

**R4-2 · baja · U4-07 · «cada cuánto» con un mensaje que no le corresponde y sin tope.** Ruta: `U4-ajustes-a.md` P-26.02, fila «Gasto fijo: periodicidad y cada cuánto». La fila usa V-21, cuyo mensaje dice «Escribe un plazo de 1 pago en adelante». Además, el esquema limita `interval_count BETWEEN 1 AND 24` (00001, línea 533): un 25 lo rechaza el servidor y no hay mensaje definido. **Corrección:** un V-xx propio, «Escribe cada cuántas veces se repite, de 1 a 24».

**R4-3 · baja · U4-11/U4-12 · PROVISIONAL sin cita en la misma línea.** Ruta: `U4-ajustes-b.md`, líneas 209 («hora del respaldo (PROVISIONAL)»), 361 y 362 (filas 15 y 16 de la tabla resumen: «lugar PROVISIONAL» sin `S1-parametros-avisos`) y 104 (redacción, cita K4). **Corrección:** agregar la cita en cada línea (`S1-parametros-avisos` en 361 y 362; K31 y SOLICITUDES en 209).

R3-2 a R3-6 siguen abiertos como bajos; la tabla de arriba trae su evidencia y su corrección está en `REVISION-U4-r3.md`.

## Verificaciones (comandos propios, resumen)
- `python3 -I datos/tmp/revU4r4/ops.py` dio: ops 210 = 210; inexistentes a 0 y b 0; Ajustes sin cita: solo `login` y `logout`, que sí están citados (falso positivo de mi expresión regular).
- `python3 -I datos/tmp/revU4r4/params.py` dio: 20 claves sembradas y 15 `@param`, todas con `sec_ok=True en_fila=True restablecer=True`; `params_json` de 00001: `weekday` 7 y `day_of_month` 1; 00003 quita `days_before` y `thresholds_pct` de `params_json`, así que no queda doble fuente con `alerts.*`.
- `python3 -I datos/tmp/revU4r4/misc.py` dio: atajos y comandos «/» dentro de U1 §5.2 y §5.3, sin duplicados; 0 hex; 0 ISO; tokens existentes; 8 bloqueos existentes.
- `python3 -I datos/tmp/revU4r4/oraculo.py` dio: las series A y B literales, igual que `calculado_*`; intereses 800.50; pagos 973.75; frescura, RN-17, RN-15 y reparto reproducidos.
- `python3 -I datos/tmp/revU4r4/defaults.py` y `secs.py`: muestra de 15 sin diferencias; 14 secciones con 0 celdas vacías.
