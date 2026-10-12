# Revisión U3 · ronda 3 de 3 (WF-1)

Revisor independiente `revisor-U3-r3` · 11-oct-2026 · Artefacto: `docs/design/U3-dashboard.md` v1.2 (478 líneas, 79.3 KB) · Evalúo contra `docs/criterios/U3.md` v1.2 (completo), decisiones §G y §K. No escribí nada de lo que reviso y no corrijo nada. Rehice cada verificación con scripts propios (`python3 -I`, en `datos/tmp/rev-u3-r3/`, ya borrados) y con cálculos propios sobre `oraculo/model.json` (Decimal).

## Veredicto

**APROBADO.** Hay 0 hallazgos altos y 0 medios. Los tres hallazgos de la ronda 2 (H-1 alto, H-2 y H-3 medios) están cerrados y los comprobé con cálculo propio. Siguen abiertos 6 hallazgos bajos de rondas anteriores (L-1 a L-6) y anoto 3 bajos nuevos (N-1 a N-3). Ninguno bloquea.

## Tabla de criterios

| Criterio | Sev. | Estado | Evidencia propia |
|---|---|---|---|
| U3-01 · ≤ 6 tarjetas, las del plan | alta | Cumple | `grep -c "^### Tarjeta"` = 6 (líneas 82, 97, 112, 129, 144 y 159), en el orden del plan y con nombres equivalentes. Cada una tiene título, valor principal, valor secundario y una acción primaria única. Tarjeta 4: el contenido es la semana en curso (N-023) y el destino P-10 sigue K38. Cita U1 §0.2. Recalculé las semanas ISO: W41 = 5 al 11 oct y W42 = 12 al 18 oct 2026. Caso domingo (Dom 11 oct): la tarjeta va en W41 y P-10 abre W42. Caso entre semana (Mié 14 oct): las dos en W42. A-03 lo prueba. |
| U3-02 · métricas de C1 y endpoint | alta | Cumple | `check.py`: 28 identificadores camelCase; 27 son `operationId` de los extractos y el 28.º es la palabra `operationId`. Las 13 rutas GET citadas existen en los extractos. M-01 a M-11 y M-13 a M-17 están en C1-datos §9, 0 faltantes. 147 identificadores en minúsculas: todos están en los extractos salvo 6 iconos y `tnum`, los 2 campos que U3 **pide** a C1 (`cumulative_accrued_cents` y `cumulative_cash_cents`, solo en la variante (a) PROVISIONAL), y `weekday` y `day_of_month`, que §15 cita solo para decir que U3 no los usa. Los enums citados coinciden con `inicio.yaml`: `status` de M-01, `availability_reading`, `basket`, `event_type`, el estado de `CalendarEvent` y el de la meta. D-4 y D-5 citan N-121 y N-162, que C1-datos asigna a P-01 con su endpoint (líneas 386 y 388); U3 escribe su fórmula y periodo y los pide a C1 (§15). |
| U3-03 · vacío, carga y error | alta | Cumple | §7 tiene 18 filas (6 tarjetas, D-1 a D-7, T-1 a T-4 y C-1) × 4 columnas (vacío, carga, error y sin conexión): 0 celdas vacías. Cada bloque de tarjeta trae su vacío y su error con «Reintentar». §1.3 trae las reglas comunes y el aviso a los 10 s. |
| U3-04 · móvil primero 412 | alta | Cumple | §8 tiene las columnas 360/412, 768/1024 y 1440, sin desplazamiento horizontal. La rejilla coincide con DESIGN §4.2: 4, 8 y 12 columnas, márgenes `space.4`, `space.6`, `space.8` y `space.12`. Táctil `size.touch.min` 44 × 44 (§1.2, §8). Capturas a 412 × 915 y 1440 × 900 (§12). |
| U3-05 · Detalle, Tendencias y Calendario | alta | Cumple | Las 3 subsecciones cubren propósito, contenido, métrica, filtros y estados: 15 de 15. Patrones: barras horizontales en D-1, columnas agrupadas en T-1 (DESIGN §8.2 «pocas series»), área apilada en T-2, líneas en T-3 y T-4; 0 donas. Recalculé T-2 con Decimal y coincide (ver H-1). Quincenas de K8, día de pago aparte y semana de lunes a domingo (§5). «Dom 28 feb 2027» existe como pago en `model.json` (P5, 2027-02-28, 1,258.68). |
| U3-06 · Karen | alta | Cumple | 31 líneas con «Karen», en: reglas del inicio, §1.4 (Parcial), tarjeta 3, D-3, D-7 (§3.1), tabla de T-2 (partes informativas), Calendario, §6 (fila D-7), §7 (D-7), §9, A-05, A-13, A-17, §12 y §13. Ninguna habla de las finanzas de su familia. No tiene cuenta y el PDF es una acción de Axel desde P-25. |
| U3-07 · próximo pago y A/B | alta | Cumple | Muestra solo la lectura activa (B por defecto); la otra se ve en D-3 con `previewAvailability`. Lleva la insignia «Pendiente confirmar». Ejemplo: $973.75 el Vie 30 oct 2026 (R-13), «Faltan 20 días» (del 10 al 30 oct), B $6,905.50 y A $10,000.00 (R-64). D-3 (b): mínimos $8,094.50 y $5,000.00 (R-61, R-62); C1 los entrega en `horizon_min_cents`. |
| U3-08 · cifras trazables | alta | Cumple | 57 montos, 20 distintos. 18 aparecen literales en `model.json` o `data.json`. $138.50 es K-03 de `S5.md`. $2,921.25 es el efectivo de R-45 más el de R-46 (la suma está escrita en el documento) y la cifra también aparece en `S2.md` S2-07. R-01, R-11, R-13 a R-16, R-45, R-46, R-50 a R-55, R-61, R-62 y R-64 coinciden con las tablas de S2 y S3. |
| U3-09 · formatos y lenguaje | media | Cumple (bajo L-1) | 0 fechas ISO. 30 fechas con día de la semana, 0 con el día equivocado (comprobado con `datetime`). 478 textos de interfaz entre «» (260 distintos): 0 en inglés, salvo los marcadores `[covered]` e `[items]` (L-1). |
| U3-10 · filtros y periodos | media | Cumple (bajo L-3) | §6 tiene 18 × 6 celdas, 0 vacías. Cada «S» corresponde a un parámetro real de la operación: D-3 (`person_id`, `family`, `period`, `credit_id`), D-4 (`period`, `person_id`), D-1 (`period`, con quincena válida en `Period`), T-1 (`category_id`), T-3 (`credit_id`, arreglo), T-4 (`person_id`) y C-1 (`period`, `credit_id`, `family`). K5 se escribe con sus tres valores incluyentes en §5 y coincide con el `Family` de `componentes-comunes.yaml` (`household`, `external`, `joint`). |
| U3-11 · tokens y gráficas accesibles | media | Cumple (bajos L-2, L-4, N-1) | 0 hex y 0 `rgb(`. 58 tokens citados, todos definidos en `C2-tokens.json`, `C2-colores.json` o `tokens.css`. `ui.surface.card` = `color.neutral.100` y `ui.text.default` = `color.neutral.900`, lo mismo que U3 §1.2. Cada gráfica tiene orden fijo, marcador o patrón, etiquetas en texto, «Ver como tabla» y `<figure>`/`<figcaption>` con su resumen (4 bloques «Resumen»). |
| U3-12 · implementación y supuestos | media | Cumple | §11 (F2 en WF-3, componentes, A-01 a A-17), §14 (9 supuestos, ninguno toca cifras ni esquema), §15. `ls docs/estado/bloqueos/U3-*` → `U3-t2-acumulados.md`, con el formato del protocolo §2 (contexto, qué falta, 3 opciones, impacto y cómo siguió). |

Comprobaciones mínimas del orquestador:
- (a) 6 tarjetas, cada una con una operación que existe y su fórmula: sí.
- (b) Vacío, carga, error y sin conexión en cada tarjeta y sección: sí (§7 y §1.3).
- (c) Móvil 412 y escritorio 1440: sí (§8 y §12).
- (d) Filtros con K5: sí (§5 y §6).
- (e) Karen solo donde corresponde: sí.
- (f) Cada U3-xx con evidencia: sí (tabla).
- (g) es-MX sin inglés: sí, salvo L-1.

## Estado de los hallazgos de la ronda 2

| Id | Sev. | Estado | Evidencia propia |
|---|---|---|---|
| H-1 · T-2 con dos orígenes | alta | **Cerrado** | T-2 exige que las dos capas tengan el origen de `cumulative_difference_cents` y prohíbe sumar `cash_cents` desde el primer mes visible (línea 238). C1 v1.2 no trae los acumulados, así que las capas quedan PROVISIONAL con su bloqueo y dos variantes. Recalculé con `model.json`. **Variante (a):** P1 973.75 / 1,947.50; P2 2,921.25 / 3,895.00; P6 13,103.78 / 14,362.48 (R-52, R-51). En los 6 periodos, devengado acumulado − efectivo acumulado = `diferencia_acumulada` (aserción cumplida). **Variante (b)**, ventana desde P2: P6 12,130.03 / 12,414.98, que es lo que dice A-11 (R-52 − 973.75 y R-51 − 1,947.50). La cifra errónea de la v1.1, 13,388.73, ya no aparece (grep = 0). A-11 suma la corrida con reloj en nov-2026. |
| H-2 · mes calendario contra periodo | media | **Cerrado** | La frase «no son meses calendario» ya no aparece (grep = 0). La regla de ejemplos (línea 26) cita tres fuentes: C1-datos M-01 «P1 = oct-2026» (línea 338), `is_canonical` («oct-2026 a mar-2027», `presupuesto-a.yaml`) y la nota de equivalencia de `S2.md` S2-09 (línea 55). El contrato ya lo decide, así que no hace falta bloqueo. La tarjeta 1, D-7, T-2, A-06 y A-11 usan la misma lectura. |
| H-3 · resumen de DESIGN §7.7 | media | **Cerrado** | §4 tiene la regla de `<figure>`/`<figcaption>` con título y resumen en `ui.text.secondary`, sin variaciones calculadas. T-1, T-2, T-3 y T-4a/b traen su texto («Resumen** (DESIGN §7.7)» = 4). El de T-2 no rotula deuda y A-11 lo prueba. La desviación frente al ejemplo «subió 12 %» está declarada (§4) y pedida a DESIGN (§15, solicitud 5). Supuesto 9. |

## Hallazgos de esta ronda

No hay altos ni medios.

### Bajos que siguen abiertos de rondas anteriores (fuera del alcance de la corrección r2)
- **L-1** · línea 151: «[covered] de [items] artículos con precio» usa nombres en inglés en un texto de interfaz. → «[con precio] de [artículos] artículos con precio».
- **L-2** · T-3 (línea 255): con filtro de crédito, las series se asignan «por línea elegida», y un filtro podría repintar las demás, contra DESIGN §8.6 y §8.8. Nombra 4 marcadores y solo 3 trazos. → Asignar la serie por la posición en `listCredits` completo y citar la tabla de DESIGN §8.6 (serie 4: rombo y trazo 8 · 4 · 2 · 4).
- **L-3** · §6, fila D-3: la parte (b), `previewAvailability`, solo admite `credit_id` en la ruta, y a esa parte le falta la línea «Este bloque no usa el filtro [nombre]» para persona y familia. → Escribirla para (b).
- **L-4** · D-6 (línea 190): pinta la barra de avance con `color.role.savings.*`, pero DESIGN §7.6 pide `ui.viz.meter`. La desviación no está declarada. → Declararla con su razón y pedirla en §15.
- **L-5** · D-1 (línea 185): ordena por % usado. DESIGN §8.2 ordena por monto y pinta el plan como marca vertical de 2 px en `ui.viz.baseline`. → Citar §8.2 y declarar el orden por %.
- **L-6** · §1.4 (línea 68) y §5: (a) «Por vencer» usa pertenencia a la lista (`uno de los de`), igual que la descripción de C1 `AlertSettings`; DESIGN §7.8 dice «dentro de los días de aviso». Con [1, 0] dan lo mismo. (b) El Calendario pide `status` = `pending` para unir «Parcial», así que un pago vencido con una parte pagada no muestra «Falta la parte de Karen». → Alinear con DESIGN (solicitud 2) y pedir también `overdue` en la unión del Calendario.

### Bajos nuevos
- **N-1** · `docs/design/U3-dashboard.md` T-2, variante (b) (línea 243), y §15. Si la puerta 2 elige B, la cima deja de ser el devengado acumulado de DESIGN §8.3. El bloqueo lo dice, pero U3 no declara esa desviación en el lugar ni la pide a DESIGN en §15 (la solicitud 3 solo trata el subtítulo). Además, §11.1 no la lista entre las desviaciones. → En la variante (b), agregar «U3 se aparta aquí de DESIGN §8.3 (cima = devengado en la ventana)», y en §15, solicitud 3, agregar «si la puerta 2 elige B, declarar esa desviación en §8.3».
- **N-2** · T-2, «Patrón» (línea 238): nombra las capas y la etiqueta «Devengado acumulado» sin condición, y el PROVISIONAL llega en el renglón siguiente. Quien lea solo «Patrón» puede tomar esos rótulos como firmes. → Agregar «(variante (a); PROVISIONAL, `U3-t2-acumulados`)» junto a esos rótulos.
- **N-3** · A-11 (línea 411): «con un mes de diferencia negativa la capa 2 vale cero con "Pagado por adelantado"» no dice que solo aplica a la variante (a). En (b), el texto es «En estos meses se pagó más de lo devengado» y la capa 2 no sale de `cumulative_difference_cents`. → Escribir «en la variante (a)…; en la (b), con un tramo en que Σ efectivo visible > Σ devengado visible, la capa 2 vale cero con "En estos meses se pagó más de lo devengado"».

## Notas
- `docs/estado/bloqueos/INDICE.md` todavía dice que hay un solo bloqueo abierto. En la carpeta están abiertos, entre otros, `U3-t2-acumulados`, `S2-presupuesto-total`, `S2-flujo-quincenal`, `S2-metas` y `S2-proyeccion-interruptor`, que U3 cita con su marca PROVISIONAL (ninguno trae «Respuesta de Axel»). No es alcance de U3: el integrador ya lo tiene como solicitud 4 de `correccion-U3-r2.md`.
- Las cifras de capas de la variante (a) son exactas al oráculo. Que $3,895.00 coincida con `horizonte.pagos_despensa_1` es casualidad: U3 la cita bien, como suma de los devengados de R-45 y R-46.
