# Revisión U1 1.1 · ronda 2 de 3

Revisor: `revisor-U1-r2` (independiente; no escribió nada de lo que revisa y no corrige) · Flujo P1 · 11-oct-2026.
Revisado: `docs/design/U1-arquitectura.md` (versión 1.1 más la corrección r1), `docs/design/tools/verificar_u1.py`, `docs/estado/p1/U1-1.1.md` y `docs/estado/p1/correccion-U1-r1.md`. Se evalúa contra `docs/criterios/U1.md` 1.1, las decisiones K2–K5, K8–K10, K12, K13, K16, K21 y K23, y los hallazgos R2-05 y R3-01 de `docs/estado/wf0/REVISION-r3.md`.

**Veredicto: NO APROBADO.** Hay **1 hallazgo de severidad alta** (A1), 0 de severidad media y 5 de severidad baja (B1 a B4 vienen de la ronda 1 y siguen abiertos; B5 es nuevo).
- Los tres hallazgos medios de la ronda 1 (M1, M2 y M3) quedan **cerrados**.
- `verificar_u1.py` da TODO OK y mis recuentos propios coinciden con él.
- El hallazgo alto es nuevo en esta ronda: U1-09 (criterio alto) pide la capacidad de K18, y U1 1.1 no la tiene. La ronda 1 marcó U1-09 «Cumple» sin revisar K18.

## 1. Método
- Corrí `python3 -I docs/design/tools/verificar_u1.py`.
- Escribí mis propios scripts en `datos/tmp/revisor-u1-r2/` (ya borrados):
  - `recuento.py`: lee U1 y, con PyYAML, `openapi.yaml`, y saca pantallas, secciones, N-xxx, tipos, la matriz §3.29, las operaciones y `period`;
  - `extra.py`: toques de F-1, apariciones de «vigente» y la sección «Para C1» del reporte;
  - `ctx.py`: grep con línea y contexto para cada K.
- U1 solo lo leí por rangos (`sed -n`) y con grep. No leí entero `openapi.yaml`: lo leyó el script.

## 2. Verificaciones pedidas

| # | Qué | Resultado | Evidencia |
|---|---|---|---|
| a | `verificar_u1.py` corre y coincide con lo reportado | **Sí** | Sale con 0 y dice «RESULTADO: TODO OK». Coincide con `correccion-U1-r1.md`: 9 entradas, 27 pantallas, 14 secciones, 203 N-xxx (138 L y 65 E), Local N-003, N-301 y N-317, F-1 = 3, 8 flujos (F-7 con 7 pasos), PROVISIONAL en 65, 419, 1149 y 1150, 30 cifras en $ (21 distintas), 197 operaciones (196 con `x-necesidades`) y `PENDIENTE C1 v1.1: ['N-316', 'N-318', 'N-319']`. Los números de línea de `U1-1.1.md` (1147 y 1148) se recorrieron 2 líneas por la nota nueva de F-7 (ver B5) |
| b | Cada K y hallazgo aplicado, con línea | **Sí** | Ver §3 |
| c | Conteos coherentes con §3.29 (recuento propio) | **Sí** | 27 `P-xx` (P-00 a P-26, sin duplicados) y 14 `P-26.xx` con los nombres exactos, así que T = 41. 203 N-xxx (138 L, 65 E), sin duplicados, y la última es N-319. Lo declarado en §3.29 (línea 764) y en la lista de cambios (línea 18) da lo mismo: 203 = 138 + 65. En la matriz: 0 N-xxx inexistentes, 0 definidas fuera de la matriz, 0 pantallas sin fila y 0 lecturas o escrituras en la columna equivocada. N-317 está en Global, P-03 y P-22; N-313 en P-21 y P-26.13; N-316 en P-26.09; N-318 y N-319 en P-26.02. N-242 está en la fila de P-21 y no en la de P-22, pero su definición dice «Desde el historial de P-21» (igual que en v1.0): no es defecto |
| d | Regla de 3 toques intacta (F-1) | **Sí** | Toques acumulados por paso: 1 → 1 → 2 → 3. Solo suben en los pasos «Toca»; teclear no cuenta (líneas 852 y 853). La regla K10 está en 846 y en 866. Los valores por defecto declarados son fecha, persona, forma de pago y el reparto de la categoría |
| e | «Lista vigente» definida una vez y `period` explícito | **Sí en U1** | Hay una sola definición, en la línea 65 (§0.2). Las apariciones de «vigente» sobre la lista (420, 571, 582, 587, 893 y 1150) remiten a §0.2; las de 619, 645 y 731 son de vigencia de promociones. N-150 (571): «Lista de la semana que manda la pantalla». N-151 a N-153 (572 a 574): «semana explícita, §0.2». N-165 (587): «semana explícita; por defecto la vigente, §0.2». Del lado de C1: en `openapi.yaml` 1.0.0, `period` sigue siendo opcional y sin valor por omisión en `getShoppingList`, `generateShoppingList`, `createShoppingListItem` y `assignShoppingListStore` (las cuatro con `x-provisional`). Queda para C1 v1.1 |
| f | Objetivo táctil ≥ 44 px | **Sí** | La línea 1073 dice «al menos 44 × 44 px»; 24 × 24 px queda solo como piso. Busqué medidas de 10 a 43 px y solo salen las líneas 35 y 1073, que son esa misma regla |
| g | Las marcas PROVISIONAL citan un bloqueo abierto | **Sí** | Las líneas 65 y 1150 citan `U1-lista-vigente-domingo`, que sigue abierto (sin «Respuesta de Axel», con el formato completo del protocolo §2). Las líneas 419 y 1149 citan K31 y `criterios-umbrales-ausentes`, que está RESUELTO con K31: «Hasta entonces, PROVISIONAL con la propuesta escrita», y aplica a S7-08, «plan sin armar». Las líneas 22, 23 y 36 son texto de la lista de cambios |
| h | «Para C1» lista exactamente las N-xxx sin operación en C1 v1.0 | **Sí** | Mi recuento contra `openapi.yaml` 1.0.0 da como N-xxx sin operación y que no son Local: N-316, N-318 y N-319, lo mismo que la línea PENDIENTE del script. «Para C1 · Nuevas» lista N-316, N-317 (Local, sin operación), N-318 y N-319. Ninguna N-xxx Local tiene operación y las 45 N-xxx que cita «Para C1» existen en U1 |

## 3. Decisiones §K y hallazgos de r3

| Decisión | Estado | Líneas (grep propio) |
|---|---|---|
| K2 alcance del cambio de presupuesto, «De este mes en adelante» preseleccionado | Cumple | 27, 238, 498, 717, 720 |
| K3 «Guardar para después» como borrador | Cumple | 28, 237, 480, 484, 485, 975 |
| K4 chat detectado por el bot y confirmado | Cumple (B1 abierto) | 29, 274, 743, 745 |
| K5 familia con tres valores incluyentes | Cumple | 24, 59, 441 |
| K8 semana lunes a domingo; quincenas 1–15 y 16–último | Cumple | 22, 57, 58, 65, 536, 1146 |
| K9 cola sin conexión solo para gastos y peso | Cumple | 25, 186, 288, 401, 460, 675, 870, 941, 1147 |
| K10 regla de toques (a) | Cumple | 26, 846, 852, 853, 866 |
| K12 tasas como se capturaron | Cumple | 30, 83, 631, 730, 731, 733 |
| K13 tope por persona opcional, sin valor fijo | Cumple | 31, 277, 655, 732, 754, 755. 374.93 solo aparece como pago por quincena (649 y 912) |
| K16 restricciones sin valor por defecto, «Faltan restricciones» | Cumple | 32, 252, 270, 655, 656, 732, 733, 953, 955, 959, 961 |
| K21 interruptor de los meses 7 a 12 | Cumple | 33, 267, 504, 505, 507, 721, 1152 |
| K23 objetivo y meta sugerida que la persona confirma | Cumple | 34, 253, 278, 664, 756, 758 |
| R2-05 objetivo táctil | Cerrado | 35, 1073 |
| R3-01 lista vigente y `period` | Cerrado en U1; abierto en C1 | 23, 65, 543 y las remisiones de §2 e |

## 4. Criterios de `docs/criterios/U1.md` 1.1

| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| U1-01 Barra lateral | alta | Cumple | Script: 9 entradas, en orden |
| U1-02 Inventario | alta | Cumple | Recuento propio: 27 + 14 = 41 y 203 N-xxx, sin duplicados (§2 c) |
| U1-03 Matriz completa | alta | Cumple en U1; pendiente C1 v1.1 | N-316, N-318 y N-319 aún no tienen operación en C1 1.0.0. Es lo esperado y están en «Para C1». 0 operaciones sin N-xxx, salvo `getOpenApiDocument` (excepción de C1 §11.2) |
| U1-04 3 toques | alta | Cumple | §2 d |
| U1-05 Flujos | alta | Cumple | 8 flujos con los 4 obligatorios. F-7 tiene la rama de K16 (953 a 961) |
| U1-06 Vistas por entidad | alta | Cumple | §2.5 tiene 31 entidades, como S1-01. plan_comida remite a §0.2 (353) |
| U1-07 Menú «/» y atajos | media | Cumple | 27 comandos con destino y 15 atajos sin duplicados |
| U1-08 Ajustes | alta | Cumple | 14 secciones con los nombres y el orden exactos |
| U1-09 Cobertura de decisiones | alta | **No cumple** | A1: falta la capacidad de K18 |
| U1-10 Karen sin acceso | alta | Cumple | La única pantalla de Karen es «Estado de cuenta de Karen» |
| U1-11 Móvil y táctil | media | Cumple | §2 f |
| U1-12 Consistencia con S1 y C1 | media | Cumple con observación | A1 también deja sin necesidad las columnas K18 de `credit` del esquema 1.2. Lo demás queda alineado, incluido N-090 con RN-28 |
| U1-13 Idioma | media | Cumple | Script: 77 nombres revisados, 0 en inglés. Los textos nuevos están en español de México |
| U1-14 Supuestos y bloqueos | media | Cumple con observación | `U1-lista-vigente-domingo` está completo. El supuesto 2 del reporte quedó viejo (B5) |
| U1-15 Sin conexión | alta | Cumple | 186, 288, 870 y 941. Solo gastos y peso entran a la cola. N-054 y N-240 llevan clave y fecha de captura |
| U1-16 Lista vigente | media | Cumple en U1; C1 pendiente | §2 e |

## 5. Hallazgos

### A1 · alta · U1 no tiene la capacidad de K18 (tarjeta bancaria y crédito de tienda)
- **Criterio:** U1-09 (`docs/criterios/U1.md:9`, `:51` y `:118`) pide una `P-xx` y sus `N-xxx` para los «parámetros de la tarjeta bancaria y del crédito de tienda, con el aviso «falta capturar» (K18)». También afecta a U1-12 y a «Para C1».
- **Ruta:** `docs/design/U1-arquitectura.md:270` (P-26.05), `:727` (N-283), `:733` (N-289) y `:632` (N-205, ciclo de la tarjeta).
- **Qué pasa:**
  - U1 1.1 no menciona K18 ni una vez.
  - N-283 dice «condiciones según el tipo (límite, tasa, días de corte y de pago, comisiones, plazo, abonos)» y N-289, «Líneas…». No están la tasa anual tomada del estado de cuenta, la base diaria o mensual, si incluye IVA, el porcentaje de pago mínimo ni el mínimo fijo.
  - El único «Falta capturar» del documento (línea 732) es de las restricciones del comparador (K16).
  - N-205 muestra el «pago mínimo» de la tarjeta, pero no dice qué pasa cuando faltan esos datos.
  - El esquema 1.2 ya agrega `credit.rate_base`, `credit.min_payment_pct_bp` y `credit.min_payment_fixed_cents` (`backend/db/migrations/00002_puerta1.sql:86-92`), con la nota «@used_new Ajustes «Créditos y condiciones» (Axel los captura de su estado de cuenta, K18). S3 (fórmula general; sin datos muestra «falta capturar»)». U4-03 también los pide.
  - «Para C1» (`docs/estado/p1/U1-1.1.md`) no los incluye, así que C1 v1.1 no tendría ninguna necesidad que los exponga.
- **Corrección sugerida:**
  - En §2.3 (P-26.05), N-283 y N-289, por tarjeta bancaria y crédito de tienda: tasa anual tal como se capturó (K12), base diaria o mensual, si incluye IVA, porcentaje de pago mínimo y mínimo fijo. Los captura Axel desde su estado de cuenta y no tienen valor por defecto.
  - En P-18 (N-205 y la necesidad que muestre intereses o pago mínimo): el aviso «Falta capturar» (abre P-26.05) cuando faltan esos datos (K18).
  - Agregar una fila K18 a «Cambios de la versión 1.1», las N-xxx cambiadas a «Para C1» y una comprobación de K18 en `verificar_u1.py`.
  - Si algún campo no tiene lugar claro en la interfaz, abrir el bloqueo `U1-<tema>` en vez de suponer.

### B1 · baja · (de la ronda 1, sigue abierto) N-316 asigna un chat a otra persona
- **Ruta:** `docs/design/U1-arquitectura.md:745`.
- **Qué pasa:** sigue diciendo «Asigna un chat detectado (N-299) a una persona (por defecto el usuario de la sesión)». K4 dice que «ella confirma».
- **Corrección sugerida:** confirmar solo para el usuario de la sesión, o dejar pendiente la asignación a otra persona hasta que esa persona la confirme.

### B2 · baja · (de la ronda 1, sigue abierto en parte) `verificar_u1.py` no comprueba los cambios de contrato
- **Ruta:** `docs/design/tools/verificar_u1.py:235-238` (PENDIENTE solo para las N-xxx nuevas).
- **Qué pasa:**
  - N-075 apunta a `createTicketScan` y N-299 a `listNotificationRecipients`, `createNotificationRecipient` y `updateNotificationRecipient`. Las dos cuentan como cubiertas.
  - Después de C1 v1.1, el script daría TODO OK aunque falten la operación de borrador (K3) y la de candidatos (K4).
  - La parte de la regla del domingo ya quedó cerrada (líneas 145 a 148).
- **Corrección sugerida:** comprobar los cambios de «Para C1» operación por operación, por lo menos N-075 y N-299.

### B3 · baja · (de la ronda 1, sigue abierto) «Para C1» pide `period` para N-151 y N-153
- **Ruta:** `docs/estado/p1/U1-1.1.md:54`; en U1, líneas 572 y 574.
- **Qué pasa:** las dos usan `updateShoppingListItem` (PATCH por id), y el artículo ya pertenece a una semana.
- **Corrección sugerida:** limitar `period` a las operaciones sobre la lista, o pedir que C1 rechace los artículos de otra semana.

### B4 · baja · (de la ronda 1, sigue abierto) Los domingos, las tarjetas de Inicio muestran semanas distintas
- **Ruta:** `docs/design/U1-arquitectura.md:419` (N-023, «Semana en curso») y `:420` (N-024, lista vigente de §0.2).
- **Qué pasa:** con la regla PROVISIONAL del domingo, las dos tarjetas muestran semanas distintas y ninguna dice cuál.
- **Corrección sugerida:** rotular la semana en las dos tarjetas, o usar en N-023 la semana por defecto de §0.2. Conviene resolverlo junto con la respuesta del bloqueo.

### B5 · baja · (nuevo) El reporte `U1-1.1.md` quedó viejo después de la corrección r1
- **Ruta:** `docs/estado/p1/U1-1.1.md:72` (supuesto 2) y las líneas que cita en «Hecho» y «Verificación».
- **Qué pasa:**
  - El supuesto 2 sigue diciendo que «U1 solo dice que la despensa bimestral y el papel dependen del interruptor» y que S2 fija qué más lleva N-090 en los meses 7 a 12. Ya no es así: la línea 504 lo fija con K21 y RN-28.
  - Las líneas citadas después de F-7 se recorrieron 2: R2-05 está en 1073, K21 de §7 en 1152 y las PROVISIONAL en 1149 y 1150.
  - Quien lea solo ese reporte (S2, C1) encuentra abierto un punto que ya se cerró.
- **Corrección sugerida:** que el integrador aplique la solicitud 1 de `correccion-U1-r1.md`: marcar el supuesto 2 como sustituido y remitir a ese reporte para las líneas.

## 6. Estado de los hallazgos anteriores

| Hallazgo | Sev. | Estado | Evidencia |
|---|---|---|---|
| M1 · N-090 contra K21 y N-318 | media | **Cerrado** | La línea 504 tiene dos tramos. En los meses 7 a 12 rige N-318: con el interruptor apagado no entran el gasto variable estimado, la aportación a ahorro ni los gastos venideros. Coincide con N-091 (505), N-318 (721), §7 (1152) y RN-28 (`docs/contratos/reglas-negocio.md:274`, «solo los compromisos contratados…, los gastos fijos y los ingresos del CSV»). La fila 12 (33) lo dice. Lo único que queda es el supuesto viejo del reporte (B5) |
| M2 · F-7 contra K16 | media | **Cerrado** | El paso 3 (953) ya no dice «por defecto» para las restricciones y tiene «Faltan restricciones · Capturar» (abre P-26.05). La recomendación depende de las restricciones en el paso 5 (955), en la nota (959) y en el atajo (961). La fila 11 (32) incluye F-7. Coincide con P-20 (252), N-225 y N-226 (655 y 656). Los pesos iguales y normalizados por defecto son decisión C (`decisiones.md:23`) |
| M3 · regla del domingo repetida sin marca | media | **Cerrado** | La línea 353 dice «Semana por defecto (§0.2)» y la 876, «semana por defecto (§0.2; N-120, N-121)». Con grep propio de «domingo», «la siguiente», «semana siguiente» y «próxima semana», la regla solo aparece en 65 (§0.2) y 1150 (§7), las dos PROVISIONAL y con su bloqueo. La fila 2 (23) incluye §2.5 y F-2, y el script lo vigila (145 a 148) |
| B1 a B4 | baja | Abiertos | §5 |
| R2-05 | baja | Cerrado | 35 y 1073 |
| R3-01 | baja | Cerrado en U1; abierto en C1 | §2 e |

## 7. Observaciones sin hallazgo
- **F-7 (959):** U1 no dice si, con las restricciones globales vacías, capturarlas solo para esa comparación permite recomendar (N-225: «se pueden cambiar para esta comparación»). El supuesto 6 de `U1-1.1.md` lo deja a S3 y a RN-11. No contradice K16.
- **N-090:** en los meses 7 a 12 deja fuera los gastos venideros, y eso sigue al pie de la letra la lista cerrada de K21 y RN-28 (`planned_item.kind = upcoming_expense` no entra). Si Axel los quiere ahí, es un cambio de K21, no de U1.
- **§3.30** sigue sin las tablas y columnas nuevas de S1 1.2. El reporte lo deja como pendiente opcional.

## 8. Qué revisar en la ronda 3
1. A1: K18 en P-26.05, N-283, N-289 y P-18 con «Falta capturar», una fila en la lista de cambios y su entrada en «Para C1». `verificar_u1.py` debe dar TODO OK con una comprobación de K18.
2. Que M1, M2 y M3 no vuelvan a aparecer (comprobaciones de las líneas 141 a 162 del script).
3. Opcional: B1 a B5.
