# REVISIÓN U3 · ronda 1 de 3

Revisor: `revisor-U3-r1` (independiente; no escribió U3 y no corrige nada) · 11-oct-2026 · Artefacto: `docs/design/U3-dashboard.md` v1.0 (421 líneas, 58.5 KB) · Evaluado **solo** contra `docs/criterios/U3.md` v1.2, `docs/plan/decisiones.md` §G y §K.

## Veredicto
**NO APROBADO.** 0 hallazgos altos, **7 medios** y 3 bajos. Los 12 criterios se cumplen en lo que pide su «Verifica»; los medios son incoherencias con DESIGN §7 y §8 (que ya existen y U3 todavía trata como pendientes), incoherencias internas de D-3 y T-2, y dos paneles de Detalle sin métrica de C1. En cada uno, quien implemente F2 tendría que adivinar.

## Cómo verifiqué (comandos propios, en `datos/tmp/revU3/`, ya borrados)
- `ops.py` (con `uv run --with pyyaml python -I`): índice de las **210 operaciones** de los 17 extractos de `docs/contratos/extractos/` con parámetros (con `$ref` resueltos) y campos de la respuesta 2xx.
- `check.py`: identificadores camelCase entre comillas invertidas citados en U3 → 28; **27 operaciones existen** (el 28.º es la palabra `operationId`). Por cada línea que cita una operación, los campos y parámetros citados están en esa operación; lo único sin coincidencia son valores de enumeración, que comprobé aparte y existen (`normal|warning|exceeded`, `on_track|behind|achieved`, `shopping_list|expense`, los 8 estados y los 12 tipos de `listCalendarEvents`).
- `nums.py`: montos `$n,nnn.nn` contra `oraculo/model.json` y `datos/semilla/data.json`; fechas ISO; hex y `rgb(`; fechas con día de la semana contra `datetime`; semanas ISO; días que faltan.
- `tok.py`: 57 nombres de token citados contra `C2-colores.json` y `C2-tokens.json` (0 sin definir); 22 iconos contra DESIGN §5 (0 faltantes).
- Lecturas puntuales: C1-datos §9 (líneas 336–354), `componentes-comunes.yaml` (parámetros `Period` y `Family`), parámetros `status` de `listScheduledPayments` y `listSavingsGoals`, descripciones de `days_left`, U1 N-020 a N-029, DESIGN §4.2, §7.6 a §7.8, §7.14 y §8.

## Tabla de criterios

| Criterio | Sev. | Estado | Evidencia propia |
|---|---|---|---|
| U3-01 · ≤ 6 tarjetas, las 6 del plan | alta | **Cumple** | `grep -c "^### Tarjeta"` = 6 (Mes contra plan, Flujo de la quincena, Próximo pago y disponible, Comidas de la semana, Mejor canasta, Ahorro y metas); cada una con título, valor principal, secundario y un solo enlace. Tarjeta 4 cita K38 y U1 §0.2 (líneas 102 y 106): `isocalendar` da 11-oct-2026 = W41, día 7 → tarjeta en W41 y P-10 en W42; 14-oct-2026 = W42 en los dos. |
| U3-02 · Métricas de C1 | alta | **Cumple, con observación (H-7)** | 27 de 27 operaciones existen; 0 campos inventados. M-01 a M-11 y M-13 a M-17 existen en C1-datos §9 con el mismo endpoint. D-4 y D-5 citan N-121 y N-162 (están en C1-datos, líneas 386 y 388) pero no tienen métrica M-xx con fórmula y periodo. |
| U3-03 · Estados vacío, carga, error | alta | **Cumple** | §7: 18 bloques (6 tarjetas + D-1 a D-7 + T-1 a T-4 + C-1) × vacío, carga, error y sin conexión, 0 celdas vacías; errores en es-MX con «Reintentar»; reglas comunes en §1.3. |
| U3-04 · Móvil primero 412 | alta | **Cumple, con observación (H-1)** | §8 con columnas 360/412, 768/1024 y 1440; 412 × 915 y 1440 × 900 en §12; `size.touch.min` 44 px. Rejilla igual a DESIGN §4.2 (4/8/12 columnas; márgenes `space.4`, `space.6`, `space.8`, `space.12`). |
| U3-05 · Detalle, Tendencias, Calendario | alta | **Cumple, con observaciones (H-2, H-6)** | 3 subsecciones × (propósito, contenido, métrica, filtros, estados) = 15 de 15; barras, líneas, área apilada y 0 donas; quincenas de K8, día de pago aparte, semana de lunes a domingo, «Dom 28 feb 2027» (domingo, verificado). |
| U3-06 · Karen | alta | **Cumple, con observación (H-5)** | 33 apariciones en 27 líneas, todas en la tarjeta 3, D-3, D-7, el Calendario, el filtro K5, §9, F2 y la trazabilidad; sin cuenta, no ve Inicio, Axel reenvía el PDF, ningún indicador de su familia. La línea 209 (tabla de T-2) queda fuera de §9 y de A-13. |
| U3-07 · Próximo pago y A/B | alta | **Cumple** | Muestra solo la lectura activa (B por defecto) con su etiqueta; la otra lectura está en D-3 (`previewAvailability`) y en P-18; insignia «Pendiente confirmar». $973.75 el «Vie 30 oct 2026» (R-13 = `despensas[0].pagos[0]`); A $10,000.00 y B $6,905.50 = `credito.eventos[0].disponible_A/_B` (R-64). «Faltan 20 días» = 30 − 10 (definición de `days_left` de `next_payment`). |
| U3-08 · Cifras trazables | alta | **Cumple, con observaciones (H-3, H-6)** | 42 montos, 18 distintos: 17 están en el oráculo o la semilla (p. ej. `periodos[0].devengado.axel` = 1197.64, `.karen` = 749.86, `data.totales.walmart` = 3233.0, `horizonte.deuda_viva_10_abr` = 6293.44); $138.50 = K-03 de `S5.md` (3,233.00 − 3,094.50, recalculado). 0 montos sin origen. |
| U3-09 · Formatos y lenguaje | media | **Cumple (bajo L-1)** | 0 fechas ISO; 25 fechas con día de la semana, 0 con el día equivocado; 220 textos de interfaz entre «», 0 en inglés; dinero con `tnum`, signo menos e icono. |
| U3-10 · Filtros y periodos | media | **Cumple, con observaciones (H-3, H-4)** | §6: 18 × 6, 0 celdas vacías. Cada «S» coincide con un parámetro real: D-3 (`person_id`, `family`, `period`, `credit_id`), D-4 (`person_id`, `period`), T-1 (`category_id`), T-3 (`credit_id`), T-4 (`person_id`), C-1 (`family`, `period`, `credit_id`; sin `person_id`, por eso «—»). La familia usa los 3 valores de K5 (`household|external|joint` en `componentes-comunes.yaml`) y §5 los rotula como K5. |
| U3-11 · Tokens y gráficas accesibles | media | **Cumple, con observación (H-2)** | 0 hex y 0 `rgb(`; 57 nombres de token, 0 sin definir; series en orden fijo con forma, patrón y etiqueta; «Ver como tabla» en cada gráfica. |
| U3-12 · Implementación y supuestos | media | **Cumple** | §11 (F2 en WF-3, componentes, A-01 a A-15), §14 (5 supuestos), §15; `ls docs/estado/bloqueos/` sin `U3-*`; los 4 bloqueos S2 que cita existen. |

Comprobaciones mínimas que pidió el orquestador: (a) 6 tarjetas, cada una con una operación que existe y su fórmula: sí. (b) Vacío, carga, error y sin conexión en cada tarjeta y sección: sí (§7 y §1.3). (c) Móvil 412 y escritorio 1440: sí (§8 y §12). (d) Filtros con K5: sí (§5 y §6). (e) Karen solo donde corresponde: sí, salvo la incoherencia H-5. (f) Cada criterio con evidencia: tabla de arriba. (g) es-MX sin inglés: sí (L-1 es un marcador, no un texto de interfaz).

## Hallazgos

### H-1 (media) · U3 contradice a DESIGN §7 (componentes), que ya está escrito
- **Ruta**: `docs/design/U3-dashboard.md` §1.2 (líneas 30 y 36), §7 (líneas 299 a 302), tarjeta 3 (línea 92), D-3 (línea 157), §5 (línea 233), §15 (línea 418).
- **Qué pasa**: U3 §15 y su reporte dan DESIGN §7 y §8 por pendientes («los escribe U2b-2»), pero ya existen (`DESIGN.md` líneas 554 a 1069) y dicen otra cosa en puntos concretos:
  1. Título de la tarjeta: U3 pone «icono de 20 px + título `font.size.title-3`, `font.weight.semibold`»; DESIGN §7.6 pone «Etiqueta (`font.size.small`, `ui.text.secondary`, sin dos puntos)».
  2. Relleno de la tarjeta: U3 usa `space.4` por debajo de 1024 px y `space.6` desde 1024; DESIGN §7.6 usa `space.6` desde 768. A 768–1023 px no coinciden.
  3. Carga de las gráficas: U3 §7 dice «marco de gráfica con ejes sin datos»; DESIGN §7.7 dice «esqueleto del alto final […], **sin ejes falsos**».
  4. Pago vencido: U3 pone «Sin marcar como pagado» con `triangle-alert`; DESIGN §7.8 pone «Vencido» con `calendar-x`. DESIGN también define «Por vencer» (`alerts.*`) y «Parcial · Falta la parte de Karen» (K15), que U3 no tiene en D-3 ni en el Calendario.
  5. Destino de la tarjeta: DESIGN §7.6 dice que «la tarjeta entera es enlace al Detalle»; en U3 cada tarjeta lleva a su pantalla destino (P-06, P-07, P-18, P-10, P-14, P-09).
- La salvedad de U3 §11.1 («si DESIGN los nombra distinto, rige DESIGN») cubre los nombres, no estas diferencias de contenido. Quien implemente F2 no sabe cuál seguir.
- **Corrección sugerida**: alinear §1.2, §5, §7 y D-3 con DESIGN §7.6 a §7.8 (tipografía y relleno de la tarjeta, carga sin ejes, estados de pago con los iconos y textos de DESIGN, y el estado «Parcial»). Donde U3 tenga una razón para apartarse (p. ej. «sin marcar como pagado» viene de C1 M-03 y U1 N-022), dejar una solicitud a DESIGN con esa razón en lugar de contradecirlo en silencio. Actualizar §15 para que ya no diga que DESIGN §7 y §8 están pendientes.

### H-2 (media) · U3 §4 contradice a DESIGN §8 (patrones de visualización)
- **Ruta**: `docs/design/U3-dashboard.md` §4 (línea 192), T-1 (línea 200), T-2 (líneas 207 y 208), T-4 (línea 223), §14 supuesto 5.
- **Qué pasa**:
  1. T-2 con diferencia negativa: U3 cambia **toda** la gráfica a dos líneas («Se pagó por adelantado en algún mes; se muestran líneas»). DESIGN §8.3 deja la capa 2 en cero en ese tramo y lo dice con «pagado por adelantado» («nunca se dibuja una capa negativa»). DESIGN §8.3 también pone la línea «Hoy», que U3 no tiene.
  2. T-4b: U3 da a Ximena triángulo y trazo continuo. DESIGN §8.7 dice que la segunda persona usa el marcador y el trazo de la serie 2 (cuadrado, 8 · 4).
  3. Etiquetas: U3 dice «desde 768 px, sobre cada punto o barra si caben». DESIGN §8.5 dice «Etiquetado selectivo: nunca un número en cada punto», y en líneas solo el último valor y los extremos.
  4. T-1: U3 dibuja el «Plan» con trama diagonal y contorno; DESIGN §8.1 dice «nunca un borde dibujado alrededor de una marca».
- El subtítulo de DESIGN §8.3, «Deuda viva $6,293.44 = …», choca con la regla de U3, «No se rotula “deuda viva”» (K19, N-091, M-11). En ese punto U3 sigue a C1 y el ajuste le toca a DESIGN.
- **Corrección sugerida**: adoptar en U3 las reglas de DESIGN §8.1, §8.3, §8.5 y §8.7 (y retirar o reescribir el supuesto 5). Para el rótulo de la deuda viva en §8.3, dejar una solicitud a DESIGN (U2b) que cite K19 y N-091.

### H-3 (media) · El ejemplo de D-3 no cabe en su propio rango
- **Ruta**: `docs/design/U3-dashboard.md` §3, ejemplo de D-3 (a) (línea 163), contra §6, fila D-3 (línea 265: «S · de hoy a hoy + 60»), y C1 M-13 («hoy..+60»).
- **Qué pasa**: hoy + 60 días el 10-oct-2026 es el **9-dic-2026** (calculado con `datetime`). El ejemplo pone cuatro pagos y el último es «Mar 15 dic 2026», que queda fuera del rango. Con el oráculo (`despensas[*].pagos`), el 10-oct-2026 D-3 trae solo 3 pagos: 30-oct, 15-nov y 30-nov. La siguiente compra es del 10-dic-2026, así que todavía no hay más pagos.
- **Corrección sugerida**: dejar el ejemplo de D-3 con tres pagos («Vie 30 oct 2026», «Dom 15 nov 2026», «Lun 30 nov 2026»), o quitar el rango de D-3 (los próximos 5 sin límite de fecha) y cambiar la celda de §6. El ejemplo del Calendario (octubre a diciembre) sí está bien.

### H-4 (media) · El filtro `status` = `pending` de D-3 deja fuera los estados que el panel promete
- **Ruta**: `docs/design/U3-dashboard.md` §3, fila D-3 (línea 157).
- **Qué pasa**: D-3 pide `listScheduledPayments` con `status` pendiente. En C1, `status` acepta un solo valor: «pending (alguna parte sin pagar), paid u overdue». Con ese filtro y el rango «de hoy en adelante», nunca llegan pagos vencidos (`overdue`) ni pagados del todo, pero el panel dice que muestra «Pagado» y «Sin marcar como pagado». Además, la tarjeta 3 enseña primero el pago vencido más antiguo (M-03, `is_overdue`), y D-3, que debería explicar de dónde sale la cifra de la tarjeta, no lo traería.
- **Corrección sugerida**: definir la petición sin ambigüedad. Por ejemplo, dos peticiones (`status` = `overdue` sin rango y `status` = `pending` de hoy a hoy + 60) o una sin `status` ordenada por fecha. Quitar los estados que no pueden llegar o definir el estado parcial (K15, ver H-1.4).

### H-5 (media) · Karen aparece en la tabla de T-2, fuera de la lista de §9 y de A-13
- **Ruta**: `docs/design/U3-dashboard.md` T-2 (línea 209: «La tabla trae también las partes del hogar y de Karen de cada mes»), contra §9 (línea 331) y §11.2 A-13 (línea 371).
- **Qué pasa**: §9 enumera dónde aparece Karen en Inicio (tarjeta 3, D-3, pagos del Calendario, D-7 y el filtro K5) y A-13 lo convierte en prueba («Karen: solo en tarjeta 3, D-3, D-7, Calendario y el filtro de familia»). La tabla de T-2 con columnas de Karen (`accrued_external_cents`, `cash_external_cents`) haría fallar A-13. El contenido sí es válido para U3-06 (su parte de la despensa en común); lo que falla es la coherencia interna.
- **Corrección sugerida**: agregar la tabla de T-2 (partes informativas de Karen) a §9 y a A-13, o quitar esas columnas de la tabla.

### H-6 (media) · La fila «Total» de T-2 no puede valer $14,362.48 y $13,103.78 en una ventana de 12 meses
- **Ruta**: `docs/design/U3-dashboard.md` T-2, ejemplo (línea 210).
- **Qué pasa**: T-2 cubre del mes en curso a 11 meses después (M-11), pero U3 dice que los totales de los 6 primeros periodos (R-51, R-52) van «en la fila “Total” de la tabla». Una fila Total de 12 meses no da esas cifras. Por ejemplo, el pago de la despensa 3 del 15-abr-2027 por $1,258.70 (`despensas[2].pagos[3]`) cae en el mes 7, fuera de P1–P6. No queda claro si la fila suma los 12 meses o solo el canon.
- **Corrección sugerida**: que la fila «Total» sume las 12 filas visibles sin cifra de canon. Si se quiere mostrar R-51 y R-52, ponerlos en una fila aparte rotulada «oct-2026 a mar-2027 (6 periodos del canon)».

### H-7 (media) · D-4 y D-5 no tienen métrica de C1 con fórmula y periodo
- **Ruta**: `docs/design/U3-dashboard.md` §3, filas D-4 y D-5 (líneas 158 y 159); `docs/contratos/C1-datos.md` §9 (líneas 336 a 354).
- **Qué pasa**: U3-02 pide que cada panel cite su métrica de C1 (identificador, fórmula, periodo y filtros) y su endpoint, «100 %». D-4 y D-5 citan necesidades U1 (N-121 y N-162) y su operación, pero C1-datos §9 no tiene métrica para ellos: M-09, M-13 y M-16 cubren N-027 y ninguna cubre N-121 ni N-162. El «Verifica» del criterio pasa (los identificadores están en C1-datos, líneas 386 y 388, y las operaciones existen), pero en sentido estricto son 16 de 18 bloques. Ni U3 ni su reporte lo señalan ni piden nada a C1.
- **Corrección sugerida**: en el reporte de U3, una solicitud a C1-datos §9 para dos métricas nuevas, «Comidas de hoy» (N-121 vía N-027: casillas con `plan_date` = hoy, periodo semana en curso, filtro `person_id`) y «Diferencias de precio» (N-162 vía N-027: los 5 renglones con mayor `difference_cents` de `getBasketComparison`, periodo lista vigente). Mientras tanto, escribir en D-4 y D-5 la fórmula, el periodo y los filtros con la marca de esa dependencia.

### L-1 (baja) · Marcadores en inglés
- **Ruta**: `docs/design/U3-dashboard.md` tarjeta 5 (línea 121): «[covered] de [items] artículos con precio». No es texto de interfaz, pero el marcador usa nombres de campo en inglés. **Corrección**: «[con precio] de [artículos] artículos con precio».

### L-2 (baja) · T-3 enumera 4 marcadores y solo 3 trazos
- **Ruta**: `docs/design/U3-dashboard.md` T-3 (línea 215). **Corrección**: citar DESIGN §8.6 (serie 4: rombo y trazo 8 · 4 · 2 · 4) en lugar de una lista incompleta.

### L-3 (baja) · D-3 (b) no admite los filtros de persona ni de familia
- **Ruta**: `docs/design/U3-dashboard.md` §6, fila D-3 (línea 265). `previewAvailability` solo recibe `credit_id`. La fila marca persona y familia con «S» sin aclarar que el bloque (b) no los usa. **Corrección**: decir en D-3 que (b) muestra «Este bloque no usa el filtro [nombre]» (regla de §3).

## Estado de los hallazgos anteriores
No aplica: es la ronda 1.

## Solicitudes derivadas (para el integrador; no son correcciones a U3)
- `docs/design/DESIGN.md` §8.3 (U2b): quitar del subtítulo el rótulo «Deuda viva …» en devengado contra efectivo, o justificarlo frente a K19, N-091 y M-11 («no se rotula deuda viva»).
- `docs/contratos/C1-datos.md` §9: las métricas de H-7 si el integrador las acepta.
