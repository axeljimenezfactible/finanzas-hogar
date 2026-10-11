# Revisión U1 1.1 · ronda 1 de 3

Revisor: `revisor-U1-r1` (independiente; no escribió nada de lo revisado) · Flujo P1 · 11-oct-2026

**Veredicto: NO APROBADO.** 0 hallazgos de severidad alta, **3 de severidad media** y 4 de severidad baja. La verificación automática pasa (`TODO OK`). Las 12 decisiones de la lista y los hallazgos R2-05 y R3-01 están aplicados en las líneas que cita el reporte. Aun así, la aplicación de K21, K16 y de la regla PROVISIONAL del domingo dejó tres inconsistencias dentro del documento.

- Artefactos revisados: `docs/design/U1-arquitectura.md` 1.1 (diff contra HEAD: 115 líneas añadidas y 81 quitadas), `docs/design/tools/verificar_u1.py`, `docs/estado/bloqueos/U1-lista-vigente-domingo.md` y el reporte `docs/estado/p1/U1-1.1.md`.
- Referencia: `docs/criterios/U1.md` 1.1; `docs/plan/decisiones.md` §K (K2–K5, K8–K10, K12, K13, K16, K21, K23); R2-05 y R3-01 de `docs/estado/wf0/REVISION-r3.md`.

## 1. Método
- Volví a correr `python3 -I docs/design/tools/verificar_u1.py` y comparé su salida con la del reporte.
- Usé scripts propios en `datos/tmp/rev-u1/` (ya borrados), con un método distinto al del agente:
  - el diff de HEAD contra la copia de trabajo, línea por línea;
  - el recuento de N-xxx en HEAD y en la versión actual;
  - el tipo de cada N-xxx contra su columna en §3.29, y su pantalla de definición contra su fila de la matriz;
  - la lectura de `openapi.yaml` con PyYAML, no con expresiones regulares;
  - la cobertura de §2.4 sin conexión;
  - una búsqueda (grep) por cada K y de «PROVISIONAL», «vigente», «px», «cola», «domingo», «tope» y «familia».
- Leí solo las líneas añadidas del diff (32 KB) y las secciones citadas. No abrí entero ni U1 ni `openapi.yaml`.

## 2. Verificaciones pedidas

| # | Verificación | Resultado | Evidencia |
|---|---|---|---|
| a | `verificar_u1.py` corre y coincide con el reporte | Cumple | El script sale con 0 y «RESULTADO: TODO OK». Da 27 pantallas, 14 secciones, 203 N-xxx (138 L y 65 E) y Local = N-003, N-301 y N-317. F-1 = 3 toques y hay 8 flujos. Marcas PROVISIONAL en 65, 419, 1147 y 1148; 30 cifras en $. `openapi.yaml` 1.0.0 tiene 197 operaciones, 196 con `x-necesidades`. Línea `PENDIENTE C1 v1.1: ['N-316', 'N-318', 'N-319']`. Todo es igual a la §Verificación del reporte. |
| b | Cada K y hallazgo aplicado, con su línea | Cumple, con M1 y M2 | Ver §3 |
| c | Conteos de pantallas y N-xxx coherentes con §3.29 (recuento propio) | Cumple | Hay 203 N-xxx definidas (138 L + 65 E), sin duplicados; la última es N-319. HEAD tenía 199 y la última era N-315, así que «antes 199» es correcto y N-316 a N-319 van después de N-315. No se quitó ninguna N-xxx ni cambió el tipo o la pantalla de ninguna. Las líneas 18 y 764 declaran 203 / 138 / 65. Hay 27 P-xx + 14 P-26.xx, T = 41. La matriz §3.29 tiene 41 filas (26 pantallas + 14 secciones + Global). Ninguna N-xxx falta en la matriz, ninguna tiene un tipo distinto al de su columna y no hay referencias a P-xx inexistentes. Única diferencia: N-242 se define en §3.24 (P-22) y la matriz la pone en P-21, que es donde se usa («desde el historial de P-21»). Ya estaba así en v1.0 (fila 761 de HEAD) y no tiene efecto. |
| d | Regla de 3 toques intacta (F-1) | Cumple | Líneas 846–869: la regla K10 está escrita. Pasos: abrir (1), teclear (sigue en 1), categoría (2), guardar (**3**). Declara los valores por defecto: fecha, persona y forma de pago; el reparto sale de la categoría. El resumen da 3 y la variante sin conexión (870) también. Frente a v1.0 solo se agregaron «abrir, elegir, guardar» y «los valores por defecto». |
| e | «Lista vigente» definida y `period` explícito en N-150 a N-153 y N-165 | Cumple en U1, con M3 | La definición única está en la línea 65. N-150 (571), N-151 (572), N-152 (573), N-153 (574) y N-165 (587) remiten a §0.2 con la semana explícita. También lo hacen N-024 (420), N-073 (482), N-127 (543), N-160 (582) y F-3 (893, 902). Las demás apariciones de «vigente» (619, 645, 731) se refieren a promociones vigentes: es otro sentido. |
| f | Objetivo táctil ≥ 44 px | Cumple | Líneas 35 y 1071: «al menos 44 × 44 px»; los 24 × 24 px quedan solo como piso de WCAG 2.5.8. El documento no declara ningún otro tamaño táctil en px; los demás px son anchos de pantalla. |
| g | Las marcas PROVISIONAL restantes citan bloqueos abiertos | Cumple, con M3 | 65 y 1148 citan `U1-lista-vigente-domingo`, abierto y con el formato completo del protocolo §2. 419 y 1147 citan `criterios-umbrales-ausentes` (RESUELTO) más K31. K31 manda dejar el valor PROVISIONAL «con la propuesta escrita» hasta que Axel lo apruebe, y ese bloqueo incluye la definición de «plan de comidas sin armar» (S7-08), así que la marca es correcta. Pero la regla del domingo aparece también sin marca en las líneas 353 y 876 (M3). |
| h | «Para C1» lista exactamente las N-xxx sin operación en C1 v1.0 | Cumple | Leí `openapi.yaml` con PyYAML. Las N-xxx sin operación y no Local son N-316, N-318 y N-319: lo mismo que la línea PENDIENTE del script y que las «Nuevas» de «Para C1» que piden operación (N-317 figura como Local y sin operación). Las 43 N-xxx cuya fila cambió en el diff están todas citadas en «Para C1», con o sin cambio de contrato. Sobre N-075 y N-299, ver B2. |

## 3. Decisiones §K y hallazgos de r3

| Decisión | Aplicada | Líneas | Comentario |
|---|---|---|---|
| K2 | Sí | 27, 238, 498, 717, 720 | Alcance «Solo este mes» / «De este mes en adelante» (preseleccionado) en P-06, N-087, N-275 y N-278 |
| K3 | Sí | 28, 237, 480, 484, 485, 973 | Borrador en N-075; N-071 lo devuelve y N-076 lo indica; F-8 |
| K4 | Sí, con B1 | 29, 274, 743, 745 | El bot detecta el chat y la persona lo confirma en P-26.09; N-316 nueva |
| K5 | Sí | 24, 59, 441 | Tres valores incluyentes. Las demás menciones de «familia» (200, 332, 341, 346, 347, 445, 619, 629) usan el filtro sin redefinirlo |
| K8 | Sí | 22, 57, 58, 65, 536, 1144 | Marcado FIRME; no queda ningún «domingo a sábado» |
| K9 | Sí | 25, 186, 288, 398, 400, 401, 460, 461, 675, 870, 941, 1145 | §2.4 cubre las 27 pantallas: 21 en solo lectura, 4 con aviso (P-00, P-05, P-19, P-20) y captura en P-03 (pestaña Gasto) y P-22, sin traslapes. «Cola» solo aparece con N-317, N-054 y N-240 |
| K10 | Sí | 26, 846, 852, 853, 866 | — |
| K12 | Sí | 30, 83, 631, 730, 731, 733 | No queda ningún «25.87» |
| K13 | Sí | 31, 277, 655, 732, 754, 755 | 374.93 solo aparece como parte de Karen por quincena (649, 912), no como tope |
| K16 | Parcial | 32, 252, 270, 655, 656, 732, 733 | Falta F-7 (M2) |
| K21 | Parcial | 33, 267, 504, 505, 507, 721, 722, 1150 | N-090 contradice a N-318 y a N-091 (M1) |
| K23 | Sí | 34, 253, 278, 664, 756, 758 | — |
| R2-05 | Sí | 35, 1071 | — |
| R3-01 | Sí en U1, con M3 | 23, 65, 420, 482, 543, 571–574, 582, 587, 893, 902, 1148 | El lado de C1 sigue pendiente (§6) |

## 4. Criterios de `docs/criterios/U1.md`

| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| U1-01 Barra lateral | alta | Cumple | Script: 9 entradas, nombres y orden exactos |
| U1-02 Inventario | alta | Cumple | 27 P-xx (≈ 25, aceptado en WF-0), 14 P-26.xx, T = 41, 203 N-xxx únicas, todas con L o E (§2 c) |
| U1-03 Matriz completa | alta | Cumple en U1; cobertura de C1 pendiente | Ninguna N-xxx queda sin operación, salvo N-316, N-318 y N-319: son nuevas, están en «Para C1» y esperan C1 v1.1. Ninguna operación queda sin N-xxx, salvo `GET /openapi.yaml` (excepción de C1 §11.2). Las Local (N-003, N-301, N-317) no tienen operación |
| U1-04 3 toques | alta | Cumple | §2 d |
| U1-05 Flujos | alta | Cumple | 8 flujos con los 4 obligatorios; cada paso tiene una pantalla existente. F-7 queda inconsistente con K16 (M2) |
| U1-06 Vistas por entidad | alta | Cumple (sin cambios en v1.1) | Las 31 entidades de S1-01 tienen fila en las tablas de U1. Sobre la fila de plan_comida, ver M3 |
| U1-07 Menú «/» y atajos | media | Cumple | 27 comandos, todos con destino; 15 atajos sin duplicados |
| U1-08 Ajustes | alta | Cumple | 14 secciones con nombre, orden e identificador |
| U1-09 Cobertura | alta | Cumple | Las capacidades de WF-0 siguen intactas. Lo nuevo (borrador del ticket, chat detectado, interruptor de proyección, objetivo de peso, cola) tiene su P-xx y su N-xxx |
| U1-10 Karen sin acceso | alta | Cumple | Las menciones nuevas de Karen están en el filtro de familia, el reparto, el tope del comparador y los créditos. No hay pantalla de acceso para ella |
| U1-11 Móvil y acción primaria | media | Cumple | 44 × 44 px (§2 f). P-06, P-20 y P-21 siguen con una sola acción primaria |
| U1-12 Consistencia | media | **No cumple** | M1 (N-090 contra N-318 y N-091) y M2 (F-7 contra N-225 y P-20) |
| U1-13 Idioma | media | Cumple | Los textos nuevos están en español de México. El script no encuentra palabras en inglés y las fechas son válidas |
| U1-14 Supuestos y bloqueos | media | **No cumple** | M3: la regla PROVISIONAL del domingo aparece sin marca en 353 y 876. M1: el supuesto 2 del reporte cambia cifras de N-090 y debió ser bloqueo |
| U1-15 Sin conexión | alta | Cumple | F-1 (870) y F-6 (941) con P-03 y P-22: «n en cola», clave y fecha de captura, falla sin pérdida, solo gastos y peso, 3 toques. N-054 y N-240 están en la matriz |
| U1-16 Lista vigente | media | **Parcial** | La definición única (65) y las escrituras con semana explícita cumplen, pero hay dos repeticiones sin remisión a §0.2 (M3). El lado de C1 sigue pendiente |

## 5. Hallazgos

### M1 · media · N-090 contradice K21 y a N-318 sobre los meses 7 a 12
- **Ruta:** `docs/design/U1-arquitectura.md:504` (N-090). Comparar con 505 (N-091) y 721 (N-318), y con el supuesto 2 de `docs/estado/p1/U1-1.1.md:72`.
- **Qué pasa:** K21 dice que, por defecto, los meses 7 a 12 llevan «solo compromisos contratados más gastos fijos e ingresos del CSV».
  - N-318 aplica esa regla a «los meses 7 a 12 de P-07» y N-091 la repite.
  - En cambio, N-090 (el flujo por quincena de P-07, 24 quincenas) sigue listando en todas las quincenas el gasto variable estimado, la aportación planificada a ahorro y los gastos venideros. Solo la despensa bimestral y el papel dependen del interruptor.
- **Consecuencia:** la misma pantalla tiene dos definiciones de lo que entra en los meses 7 a 12. Además, el supuesto 2 del reporte («qué otras partidas lleva N-090 en esos meses lo fija S2») cambia cifras. Según el protocolo §2, eso es un bloqueo, no un supuesto.
- **Corrección sugerida:** decir en N-090 que en las quincenas de los meses 7 a 12 rige N-318, igual que en N-091:
  - apagado: solo compromisos contratados, gastos fijos e ingresos del CSV;
  - encendido: lo mismo más la despensa bimestral a crédito y el papel.

  Si U1 no puede afirmar que K21 abarca el gasto variable y la aportación de ahorro de N-090, debe abrir `docs/estado/bloqueos/U1-proyeccion-quincenas-7-12.md` y marcar esa parte de N-090 como PROVISIONAL, citando el bloqueo.

### M2 · media · F-7 no aplica K16
- **Ruta:** `docs/design/U1-arquitectura.md:953` (F-7, paso 3), 955–956 (pasos 5 y 6) y 959 (atajo desde una canasta).
- **Qué pasa:** F-7, el flujo obligatorio (b), no se actualizó con K16.
  - El paso 3 abre el comparador «con restricciones, tope y pesos por defecto (N-225)», pero N-225 1.1 dice «No tienen valor por defecto».
  - Los pasos 5 y 6 y el atajo dan por hecha una recomendación («Usar escenario recomendado»).
  - Hoy Axel no ha capturado las restricciones, así que el primer uso real termina en «Faltan restricciones · Capturar», sin recomendación (P-20, línea 252). El flujo no lo menciona y el cambio 11 de la lista no incluye F-7.
- **Corrección sugerida:**
  - Paso 3: «con las restricciones globales de P-26.05 (N-225), el tope por persona de P-26.12 y los pesos».
  - Agregar la condición: si faltan restricciones, P-20 muestra «Faltan restricciones · Capturar» (abre P-26.05) y no recomienda (K16).
  - Hacer lo mismo en el atajo y sumar F-7 al cambio 11.

### M3 · media · La regla del domingo se repite sin la marca PROVISIONAL
- **Ruta:** `docs/design/U1-arquitectura.md:353` (vista de plan_comida) y 876 (F-2, paso 1). La definición está en la línea 65 y en §7, línea 1148.
- **Qué pasa:** §0.2 se presenta como la «única definición» de la semana por defecto de P-10, P-13 y P-14, y marca la regla del domingo como PROVISIONAL (bloqueo abierto `U1-lista-vigente-domingo`). Sin embargo, la línea 353 («Semana en curso (domingos, la siguiente)») y la 876 («semana en curso (los domingos, la siguiente)») repiten la regla como hecho:
  - sin la marca PROVISIONAL;
  - sin citar el bloqueo;
  - sin remitir a §0.2.
- **Consecuencia:** contradice el cambio 2 («definidas en un solo lugar») y el «Cómo seguí» del bloqueo. Si Axel elige la opción 2 o la 3, esas dos líneas quedan desfasadas. El script no lo detecta porque solo busca la palabra PROVISIONAL y la definición en negritas.
- **Corrección sugerida:** sustituir ambos textos por «semana por defecto (§0.2)».

### B1 · baja · N-316 permite asignar un chat a otra persona
- **Ruta:** `docs/design/U1-arquitectura.md:745` (N-316); comparar con 274 (P-26.09).
- **Qué pasa:** K4 dice que el bot detecta el chat cuando la persona le escribe «y ella confirma», y P-26.09 lo dice igual. Pero N-316 «asigna un chat detectado a una persona (por defecto el usuario de la sesión)». Así, un titular podría asignarle a la otra persona un chat que ella no confirmó. Como el bot es el mismo de MultiX, un chat ajeno detectado por error podría recibir los avisos del hogar.
- **Corrección sugerida:** que N-316 confirme el chat solo para el usuario de la sesión, o que la asignación a otra persona quede pendiente hasta que ella la confirme.

### B2 · baja · `verificar_u1.py` no detecta dos casos
- **Ruta:** `docs/design/tools/verificar_u1.py:212` (PENDIENTE solo para las N-xxx nuevas), 141 (busca PROVISIONAL por palabra) y 132 (U1-16 solo busca la definición en negritas).
- **Qué pasa:**
  - K3 pide una operación nueva para N-075, que en `openapi.yaml` 1.0.0 apunta a `createTicketScan`.
  - K4 pide leer los candidatos pendientes en N-299, hoy `listNotificationRecipients`.
  - Las dos cuentan como cubiertas, así que después de C1 v1.1 el script daría TODO OK aunque falten.
  - Tampoco detecta repeticiones sin marca como las de M3.
- **Corrección sugerida:**
  - Comprobar los cambios de contrato de «Para C1» operación por operación: al menos que N-075 tenga una operación distinta de `createTicketScan` y que N-299 devuelva los candidatos.
  - Hacer que el script falle si aparece «domingos, la siguiente» fuera de §0.2 y §7.

### B3 · baja · «Para C1» pide `period` para N-151 y N-153
- **Ruta:** `docs/estado/p1/U1-1.1.md:54`; en U1, líneas 572 y 574.
- **Qué pasa:** N-151 y N-153 usan `updateShoppingListItem` (PATCH por id de artículo), que no lleva `period` porque el artículo ya pertenece a una semana. Volver `period` obligatorio ahí agrega un dato redundante que puede contradecir la semana del artículo.
- **Corrección sugerida:** limitar `period` a las operaciones sobre la lista: `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `assignShoppingListStore` y las de N-024 y N-160. Si C1 acepta `period` en N-151 y N-153, que rechace los artículos de otra semana.

### B4 · baja · Los domingos, las tarjetas de Inicio muestran semanas distintas
- **Ruta:** `docs/design/U1-arquitectura.md:419` (N-023, «Semana en curso») y 420 (N-024, «Lista de compras vigente»).
- **Qué pasa:** con la regla PROVISIONAL del domingo, ese día «Comidas de la semana» muestra la semana que termina y «Mejor canasta» compara la lista de la semana siguiente. Ninguna de las dos tarjetas dice qué semana muestra.
- **Corrección sugerida:** rotular la semana en ambas tarjetas (por ejemplo, «Semana del 12-oct-2026») o usar en N-023 la semana por defecto de §0.2. Conviene resolverlo junto con la respuesta del bloqueo.

## 6. Estado de los hallazgos anteriores

| Hallazgo | Sev. | Estado | Evidencia |
|---|---|---|---|
| R2-05 (U1 §6.1 punto 4: 24 × 24 px) | baja | **Cerrado** | Líneas 35 y 1071: el objetivo es 44 × 44 px y los 24 px quedan solo como piso |
| R3-01 (lista vigente y `period`) | baja | **Cerrado en U1** (con M3); **abierto en C1** | En `openapi.yaml` 1.0.0, `period` sigue opcional (query) en `getShoppingList`, `generateShoppingList`, `createShoppingListItem` y `assignShoppingListStore`. Lo corrige C1 v1.1 con «Para C1» |

## 7. Observaciones sin hallazgo
- **U1-03:** N-316, N-318 y N-319 no tienen operación en C1 1.0.0. Es lo esperado hasta C1 v1.1 y están en «Para C1». Al cerrar C1 v1.1 hay que volver a correr el script: la línea PENDIENTE debe desaparecer.
- **U1-15 en C1:** `createExpense` y `createWeightEntry` ya aceptan `Idempotency-Key` en 1.0.0. La fecha de captura llega con C1 v1.1.
- **§3.30:** no lista las tablas nuevas de S1 1.2. El reporte lo deja como pendiente opcional y queda fuera de la lista de K de esta revisión.
- **Supuestos 1, 3, 4, 5 y 6 del reporte:** tratan de la ubicación en la interfaz o de su comportamiento, sin cifras ni esquema; son aceptables. El supuesto 6 (fecha de deuda en cero opcional) ya estaba en v1.0 (N-225 y N-288).
- **Bloqueo `U1-lista-vigente-domingo`:** tiene el formato completo del protocolo §2 y está justificado. `decisiones.md` solo menciona el domingo para el resumen semanal y para el aviso de plan sin armar, no para la semana por defecto.

## 8. Qué revisar en la ronda 2
1. Que se hayan corregido M1, M2 y M3 (las bajas son opcionales) y que `verificar_u1.py` dé TODO OK.
2. Que N-090 use la misma regla que N-318 y N-091, o que exista el bloqueo nuevo con su marca.
3. Que F-7 ya no hable de «restricciones por defecto» y tenga la rama «Faltan restricciones».
4. Que «la siguiente)» solo aparezca en §0.2 y §7, o con remisión a §0.2.
