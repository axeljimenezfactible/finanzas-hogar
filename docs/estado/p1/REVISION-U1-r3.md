# Revisión U1 · ronda 3 de 3 (`revisor-U1-r3`)

Flujo P1 · 11-oct-2026. Artefacto revisado: `docs/design/U1-arquitectura.md` versión 1.1 (141,969 bytes, 1,161 líneas) y `docs/design/tools/verificar_u1.py` (21,744 bytes). Criterios: `docs/criterios/U1.md` 1.1, decisiones §K (K2 a K5, K8 a K10, K12, K13, K16, K21 y K23, más K18 por el hallazgo A1 de la ronda 2) y los hallazgos R2-05 y R3-01 de `docs/estado/wf0/REVISION-r3.md`.

**Veredicto: APROBADO.** No hay hallazgos de severidad alta ni media. A1 (K18), el único hallazgo alto de la ronda 2, queda **cerrado** en U1. Hay 7 hallazgos de severidad baja: B1 a B5 siguen abiertos desde rondas anteriores y B6 y B7 son nuevos. Los dos nuevos tocan la entrega a C1 («Para C1»), no el documento U1.

## 1. Método
- Lecturas: el protocolo, §K, `U1.md` 1.1, R2-05 y R3-01, `U1-1.1.md`, `correccion-U1-r2.md` y `REVISION-U1-r2.md` (§2, encabezados de §5 y §8). Del diff contra HEAD leí las 130 líneas agregadas completas; de las 91 quitadas, solo los primeros 150 caracteres. También leí `verificar_u1.py` completo, el apartado «Respuesta de Axel» de 12 bloqueos y el bloqueo `U1-lista-vigente-domingo` completo.
- Ejecuté `verificar_u1.py` y dos scripts propios (`datos/tmp/revisor-u1-r3/rev.py` y `rev2.py`, ya borrados) que usan otro método:
  - asignan cada N-xxx a la pantalla de su sección §3.k (en P-26, a su sección `P-26.xx`);
  - leen la matriz §3.29 con sus columnas L y E y desarrollan los rangos;
  - cargan `openapi.yaml` con PyYAML, no línea por línea.
- No leí enteros U1, `openapi.yaml`, `C1-datos.md`, `diccionario-datos.md` ni las migraciones. De ellos solo usé `grep` y lo que imprimieron los scripts.

## 2. Verificaciones pedidas
| # | Qué | Resultado | Evidencia |
|---|---|---|---|
| a | `verificar_u1.py` corre y coincide con lo reportado | **Sí** | `python3 -I docs/design/tools/verificar_u1.py` sale con 0 y dice «RESULTADO: TODO OK»: 52 OK, 0 FALLA, 9 INFO y 2 PENDIENTE. Coincide con `correccion-U1-r2.md`: 52 OK (45 de antes y 7 de K18), PROVISIONAL en 66, 421, 1151 y 1152, 30 cifras en $ (21 distintas), 77 nombres revisados y 197 operaciones (196 con `x-necesidades`). Las dos líneas PENDIENTE son N-316, N-318 y N-319, y los 4 campos de K18 en `Credit`, `CreditCreate` y `CreditUpdate`. `U1-1.1.md` quedó viejo: cita las PROVISIONAL en 65, 419, 1147 y 1148 y una sola línea PENDIENTE (B5) |
| b | Cada K y hallazgo aplicado, con línea | **Sí** | Ver §3 |
| c | Conteos coherentes con §3.29 (recuento propio) | **Sí** | **Pantallas:** 27 `P-xx` (P-00 a P-26, sin duplicados, todas con nombre y entrada de barra) y 14 `P-26.xx` con los ids y nombres del criterio, así que T = 41. **Necesidades:** 203 N-xxx (138 L y 65 E); 0 duplicadas, 0 con otro tipo, 0 citadas sin definir; la última es N-319. **Nuevas:** N-316 (E, P-26.09), N-317 (L, Global, Local), N-318 (L, P-26.02) y N-319 (E, P-26.02). **Lo declarado coincide:** la línea 18 dice 203 (antes 199) = 138 + 65 y la 766 dice 203 hasta N-319 = 138 + 65. **Matriz:** 41 filas (Global, 26 pantallas y 14 secciones); 0 N-xxx fuera, 0 inexistentes, 0 E en la columna L, 0 L en la columna E y 0 pantallas sin fila. Solo N-242, definida en §3.24 (P-22), aparece únicamente en la fila de P-21; se usa desde el historial de P-21 y viene de v1.0 (§7) |
| d | Regla de 3 toques intacta (F-1) | **Sí** | Mi cuenta: 1 → 1 → 2 → 3; solo suman los pasos «Toca» y teclear no cuenta (855). La regla K10 está en 848, 854 y 868. Valores por defecto declarados: fecha, persona y forma de pago (854) y el reparto de la categoría (N-051). La variante «Sin conexión» (872) también da 3, igual que el resumen de toques (979) |
| e | «Lista vigente» definida y `period` explícito en N-150 a N-153 y N-165 | **Sí en U1** | **Definición:** una sola, en 66 (§0.2). **Remisiones:** las apariciones de «vigente» sobre la lista (422, 573, 584, 589, 895 y 1152) remiten a §0.2; las de 621, 647 y 733 hablan de promociones vigentes. **Necesidades:** N-150 (573) toma «la semana que manda la pantalla»; N-151 a N-153 (574 a 576) dicen «semana explícita, §0.2»; N-165 (589), «semana explícita; por defecto la vigente, §0.2». N-073, N-127 y N-160 también remiten a §0.2. **Domingo:** la regla solo aparece en 66 y 1152 (y en la fila 2 de la lista de cambios). **Lado de C1:** en C1 1.0.0 no encuentro un parámetro `period` obligatorio en las operaciones de la lista; queda para C1 v1.1 |
| f | Objetivo táctil ≥ 44 px | **Sí** | La línea 1075 dice «al menos 44 × 44 px». Busqué toda medida menor de 44 px y solo salen 35 y 1075, que son el piso de 24 × 24 px de WCAG 2.5.8. §6.4 cubre las 27 pantallas |
| g | Las marcas PROVISIONAL restantes citan su bloqueo | **Sí** | **66 y 1152:** citan `U1-lista-vigente-domingo`, que está abierto (sin «Respuesta de Axel») y tiene los 5 apartados del protocolo §2. **421 y 1151:** citan `criterios-umbrales-ausentes` y K31. Ese bloqueo está RESUELTO con K31, que dice «Hasta entonces, PROVISIONAL con la propuesta escrita» y aplica a S7-08 («plan de comidas sin armar (domingo)», `docs/criterios/S7.md:10` y `:46`). **22, 23 y 36:** son texto de la lista de cambios, no marcas |
| h | «Para C1» lista exactamente las N-xxx sin operación en C1 v1.0 | **Sí** | Contra `openapi.yaml` 1.0.0, las N-xxx sin operación y que no son Local son N-316, N-318 y N-319, lo mismo que la línea PENDIENTE del script. «Para C1 · Nuevas» de `U1-1.1.md` lista N-316, N-317 (Local, sin operación), N-318 y N-319. Ninguna N-xxx Local tiene operación y todas las N-xxx que cita «Para C1» existen en U1. La única operación sin N-xxx es `getOpenApiDocument`, la excepción de C1 §11.2. Los cambios de contrato de K18 no están aquí sino en la adenda de `correccion-U1-r2.md` (B6) |

## 3. Decisiones §K y hallazgos, con línea
Búsqueda con `\bKn(?!\d)` sobre U1 1.1:

| Origen | Líneas | Qué dice | Estado |
|---|---|---|---|
| K2 | 27, 239, 500, 719, 722 | Alcance «Solo este mes» o «De este mes en adelante», con el segundo preseleccionado. En P-06, «este mes» es el mes que se muestra; en P-26.02, el mes en curso | Aplicado |
| K3 | 28, 238, 482, 486, 487, 977 | «Guardar para después» guarda las correcciones como borrador; N-071 lo devuelve y N-076 lo indica | Aplicado |
| K4 | 29, 275, 745, 747 | El bot detecta el chat y la persona lo confirma en P-26.09 (N-299 y N-316) | Aplicado (B1 sigue abierto) |
| K5 | 24, 60, 443 | Filtro «familia» con tres valores: Hogar y Karen incluyen lo conjunto | Aplicado |
| K8 | 22, 58, 59, 66, 538, 1148 | Semana de lunes a domingo; quincenas del 1 al 15 y del 16 al último día. FIRME | Aplicado |
| K9 | 25, 290, 403, 462, 677, 872, 943, 1149 (y «Sin conexión» en 187 y 400) | Cola en el teléfono solo para gastos y peso, con clave de idempotencia y fecha de captura; las 27 pantallas declaran qué hacen sin conexión (mi recuento sobre §2.4: 0 sin declarar) | Aplicado |
| K10 | 26, 848, 854, 855, 868 | Cuentan solo los toques en controles | Aplicado |
| K12 | 30, 38, 84, 271, 629, 633, 729, 732, 733, 735 | Tasas tal como se capturaron; los cálculos usan la fracción | Aplicado |
| K13 | 31, 278, 657, 734, 756, 757 | Tope por persona, opcional, en P-26.12. 374.93 solo aparece como la parte de Karen por quincena (651 y 914), nunca como tope | Aplicado |
| K16 | 32, 253, 271, 657, 658, 734, 735, 955, 957, 961, 963 | Restricciones globales sin valor por defecto, ajustables en cada comparación; «Faltan restricciones» y sin recomendación, también en F-7 | Aplicado |
| K18 | 38, 251, 271, 289, 629, 630, 634, 729, 735, 1155 | Ver §6, A1 | **Aplicado (A1 cerrado)** |
| K21 | 33, 268, 506, 507, 509, 723, 1154 | Interruptor N-318, apagado por defecto; N-090 y N-091 siguen su regla en los meses 7 a 12 | Aplicado |
| K23 | 34, 254, 279, 666, 758, 760 | Objetivo (bajar, mantener o subir) y meta sugerida que la persona confirma o cambia | Aplicado |
| R2-05 | 7, 18, 35, 1075 | ≥ 44 × 44 px; 24 × 24 px solo como piso | Aplicado |
| R3-01 | 7, 18, 23, 66, 545 (y remisiones a §0.2 en N-073, N-150 a N-153, N-160, N-165, F-2 y F-3) | Lista vigente definida una sola vez; la semana va explícita | Aplicado en U1; C1 pendiente |

S1 1.2 respalda las necesidades nuevas (`00002_puerta1.sql`):
- `attachment.ocr_draft_json` (K3), `telegram_chat_candidate` (K4), `idempotency_key` y `captured_at` en `expense` y `weight_entry` (K9);
- en `credit`: `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents` y `annual_rate_decimals` (K18 y K12);
- `rate_decimals` (K12), `health_profile.goal` (K23) y la clave `projection.repeat_cadence_enabled` (K21, RN-28).

## 4. Criterios de `docs/criterios/U1.md` 1.1
| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| U1-01 Barra lateral | alta | Cumple | 9 secciones en orden. «Buscar» y «Registrar gasto» son las dos acciones fijas (§1.2 punto 1, 154), no secciones |
| U1-02 Inventario | alta | Cumple | §2 c: 27 + 14 = 41 y 203 N-xxx, sin duplicados; lo declarado coincide |
| U1-03 Matriz completa | alta | Cumple en U1; C1 v1.1 pendiente | 0 N-xxx sin operación salvo las nuevas N-316, N-318 y N-319, que están en «Para C1»; 0 operaciones sin N-xxx salvo la excepción de C1 §11.2. Las nuevas cuentan como cubiertas cuando C1 v1.1 tenga sus operaciones |
| U1-04 3 toques | alta | Cumple | §2 d |
| U1-05 Flujos | alta | Cumple | 8 flujos, cada paso con una pantalla existente. Los 4 obligatorios: (a) F-1, (b) F-7, (c) F-2, cuyo paso 8 genera la lista con N-127, y (d) F-8 |
| U1-06 Vistas por entidad | alta | Cumple | §2.5 tiene 31 entidades, de `persona` a `respaldo`, como S1-01. `plan_comida` remite a §0.2 (355) |
| U1-07 Menú «/» y atajos | media | Cumple | 27 comandos, todos con destino; 15 atajos, sin duplicados (script y recuento propio) |
| U1-08 Ajustes | alta | Cumple | 14 secciones de `P-26.01` a `P-26.14` con los nombres y el orden del criterio |
| U1-09 Cobertura de decisiones | alta | **Cumple** | Recorrí 29 capacidades y todas tienen pantalla y necesidades (lista abajo). K18 ya está |
| U1-10 Karen sin acceso | alta | Cumple | La única pantalla de Karen es P-25, «Estado de cuenta de Karen». No hay inicio de sesión para ella. En P-26.12 aparece como externa, sin cuenta |
| U1-11 Móvil y táctil | media | Cumple | §6.4 cubre las 27 pantallas; táctil en §2 f |
| U1-12 Consistencia con S1 y C1 | media | Cumple con observación | **S1:** S1 1.2 respalda todo lo nuevo (§3). **Diccionario:** coinciden 10 nombres que elegí: `sesion_dispositivo`, `regla_reparto`, `tasa_por_plazo`, `pago_programado`, `peso_historico`, `inventario_item`, `tiempo_comida`, `meta_ahorro`, `fortnight_payment_cap_cents` y `daily_kcal_target`. **Observaciones:** B7 y la representación del IVA en S1 (§7.1) |
| U1-13 Idioma | media | Cumple | El script revisa 77 nombres y no halla ninguno en inglés. Los textos nuevos están en español de México: «Sin conexión con el servidor de casa · n en cola», «Falta capturar · Capturar», «Faltan restricciones · Capturar» y «Meta sin confirmar» |
| U1-14 Supuestos y bloqueos | media | Cumple con observación | `U1-lista-vigente-domingo` está completo y abierto, y lo provisional va marcado. El reporte `U1-1.1.md` quedó viejo (B5). El bloqueo nuevo aún no está en `INDICE.md` (§7.3) |
| U1-15 Sin conexión | alta | Cumple | 187, 290, 400, 403, 462, 677, 872, 943 y 1149. Solo gastos y peso entran a la cola. N-054 y N-240 llevan clave y fecha de captura y están en las filas de P-03 y P-22 de la matriz |
| U1-16 Lista vigente | media | Cumple en U1; C1 pendiente | §2 e |

**Recorrido de U1-09:**
1. Inicio con Resumen, Detalle, Tendencias y Calendario: P-01.
2. Karen en la despensa en común: P-01, P-02 y P-03 (N-022, N-026, N-040, N-041, N-052, N-054).
3. Karen en el crédito: P-19 (N-215) y P-26.12 (N-309).
4. Comparador con 4 objetivos y recomendación: P-20 (N-225, N-226).
5. Calendario de promociones: P-26.05 (N-287).
6. Interruptor A/B con «pendiente confirmar»: P-17 (N-193), P-26.05 (N-284) y P-01 (N-022).
7. Plan de comidas con 4 tiempos y «fuera»: P-10 (N-120).
8. Lista de compras: P-13 (N-150 a N-153).
9. Recetas editables: P-12 (N-142 a N-144) y P-26.07 (N-295).
10. Perfil de salud, IMC y kcal: P-21 (N-230) y P-26.13 (N-311).
11. Precios con frescura: P-15 (N-170).
12. Membresía de Sam's y su punto de equilibrio: P-14 (N-164).
13. Mercados: P-26.06 (N-290, N-292) y P-13 (N-150).
14. Ticket con OCR: P-05 (N-070 a N-076).
15. Estado de cuenta de Karen en PDF: P-25 (N-261).
16. Reporte mensual: P-24 (N-255 a N-257).
17. Vista web mensual: P-24 (223; N-255 y N-256; F-5, 926).
18. Metas de ahorro: P-09 (N-110, N-111, N-113, N-115) y P-26.13 (N-312, N-313).
19. Ingresos y extras: P-08 (N-100 a N-105).
20. Tema: P-26.10 (N-301, Local).
21. Sesiones por dispositivo: P-26.01 (N-271, N-274).
22. K2, alcance del cambio de presupuesto: P-06 (N-087) y P-26.02 (N-275, N-278).
23. K3, borrador del ticket: P-05 (N-071, N-075, N-076).
24. K4, chat detectado y confirmado: P-26.09 (N-299, N-316).
25. K13, tope por persona: P-26.12 (N-309, N-310).
26. K16, restricciones duras: P-26.05 (N-288, N-289) y P-20 (N-225, N-226).
27. K18, datos de tarjeta y crédito de tienda: P-26.05 (N-283, N-289) y P-18 (N-200, N-201, N-205).
28. K21, interruptor de la proyección: P-26.02 (N-318, N-319) y P-07 (N-090, N-091, N-093).
29. K23, objetivo y meta de kcal: P-21 (N-230, N-313) y P-26.13 (N-311, N-313).

## 5. Hallazgos nuevos

### B6 · baja · La adenda K18 de «Para C1» no está en `U1-1.1.md`
- **Criterio:** U1-12 (consistencia con C1) y la parte «Para C1» de A1 (ronda 2).
- **Ruta:** `docs/estado/p1/U1-1.1.md`, sección «Para C1», que no menciona K18 ni N-200, N-201, N-205 o N-283. La adenda está en `docs/estado/p1/correccion-U1-r2.md`, «Para C1 (adenda K18 · corrección r2)», con N-200, N-201, N-205, N-283 y N-289.
- **Escenario:** si C1 v1.1 trabaja solo con «Para C1» de `U1-1.1.md`:
  - `Credit`, `CreditCreate` y `CreditUpdate` se quedan sin `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents` y `annual_rate_decimals`;
  - N-201 y N-205 no aplican la regla «sin valor y con la indicación «falta capturar», nunca en cero».
  Nada lo detendría: `docs/criterios/C1.md` 1.1 no menciona K18 (0 coincidencias), y `verificar_u1.py` lo marca como PENDIENTE, no como falla.
- **Corrección sugerida:** el integrador suma la adenda a «Para C1» de `U1-1.1.md` (solicitud 1 de `correccion-U1-r2.md`) o entrega ambos archivos al agente de C1 v1.1.

### B7 · baja · «Para C1» no pide quitar `updatePerson` de N-289
- **Criterio:** U1-03 y U1-12 (la columna de endpoint de la matriz debe ser exacta).
- **Ruta:**
  - `docs/estado/p1/U1-1.1.md`, «Para C1», entrada de N-289;
  - `docs/contratos/openapi.yaml`, `PATCH /api/v1/persons/{person_id}` (`updatePerson`), con `x-necesidades` [N-289, N-310];
  - U1 734 (N-288: «El tope por persona está en P-26.12 (N-309, K13)») y 735 (N-289 ya no edita personas).
- **Escenario:** en v1.0, N-288 leía el tope por persona y por eso N-289 llegaba a `updatePerson`. El cambio 10 (K13) movió el tope a P-26.12 (N-309 y N-310). Si C1 v1.1 deja `updatePerson` en N-289, la matriz dirá que P-26.05 edita personas. La cobertura no se rompe, porque N-310 conserva la operación.
- **Corrección sugerida:** sumar a «Para C1»: «N-289: quitar `updatePerson`; el tope por persona lo edita N-310 (K13)».

## 6. Estado de los hallazgos anteriores
| Hallazgo | Sev. | Estado | Evidencia |
|---|---|---|---|
| A1 (r2) · K18 en U1 | alta | **Cerrado** | Ver el detalle después de la tabla |
| M1 (r1) · N-090 contra K21 | media | Sigue cerrado | N-090 (506) y N-091 (507) aplican N-318 en los meses 7 a 12; el script da OK |
| M2 (r1) · F-7 sin K16 | media | Sigue cerrado | Rama «Faltan restricciones» en 955, 957, 961 y 963; el script da OK |
| M3 (r1) · regla del domingo repetida | media | Sigue cerrado | Mi búsqueda de «domingo»: la regla solo está en 23 (lista de cambios), 66 (§0.2) y 1152 (§7). Las demás apariciones son «lunes a domingo» o el resumen semanal |
| B1 · N-316 asigna el chat a otra persona | baja | Abierto | 747 sin cambios: «a una persona (por defecto el usuario de la sesión)». K4 dice «ella confirma» |
| B2 · el script no revisa los cambios de contrato | baja | Abierto en parte | Ahora revisa los campos de K18 (líneas 270 a 281), pero no `scope`, `family`, `Idempotency-Key` y fecha de captura, ni `period` |
| B3 · «Para C1» pide `period` para N-151 y N-153 | baja | Abierto | `U1-1.1.md` sin cambios |
| B4 · los domingos, las tarjetas de Inicio muestran semanas distintas | baja | Abierto | N-023 (421) usa la «Semana en curso» y N-024 (422), la «Lista de compras vigente», que el domingo es la siguiente |
| B5 · `U1-1.1.md` viejo | baja | Abierto, y más viejo | Cita PROVISIONAL en 65, 419, 1147 y 1148; dice 137.7 KB, 15.6 KB del script y 81 → 115 líneas; no menciona K18 |

**Detalle de A1:**
- **Captura:** P-26.05 (271) pide los 5 datos sin valor por defecto, con «Falta capturar» (K12 y K18). N-283 (729) los lee y dice cuáles faltan. N-289 (735) los edita sin valor por defecto.
- **Lectura en P-18:** la fila de P-18 en §2.2 (251) muestra «Falta capturar · Capturar» y abre P-26.05. N-200 (629) muestra los datos; N-201 (630) y N-205 (634) no calculan con valores supuestos.
- **Regla común y registro:** está en §2.4 (289), en §7 (1155) y en la fila 17 de la lista de cambios (38).
- **Respaldo:** `00002_puerta1.sql:112-115` («NULL = falta capturar») y RN-26. El script tiene 7 comprobaciones de K18 y dan OK.
- **«Para C1»:** la adenda existe, pero falta fusionarla (B6).

## 7. Observaciones sin hallazgo
1. **IVA (S1 y RN-26, fuera de U1).** `00001_esquema.sql:448` y RN-26 (`reglas-negocio.md:266`) definen `vat_on_interest_bp` NULL como «la tasa ya trae IVA».
   - Así, en `bank_card` y `store_credit`, «si incluye IVA» nunca puede quedar como «falta capturar».
   - K18, y U1 con él (N-283 en 729 y §2.4 en 289), dicen que falta capturar si falta cualquiera de los 5 datos.
   - U1 sigue a K18, así que el arreglo es de S1 1.2 y RN-26 (solicitud 2 de `correccion-U1-r2.md`). Conviene resolverlo antes de que C1 v1.1 exponga ese campo.
2. **Comparador con una línea en «Falta capturar».** N-226 (658) no dice si los escenarios con esa línea se descartan. La regla común de §2.4 (289) ya cubre que no se calculan intereses con valores supuestos; el tratamiento exacto lo fija S3.
3. **Índice de bloqueos.** `U1-lista-vigente-domingo` aún no está en `docs/estado/bloqueos/INDICE.md`. Es la solicitud de `U1-1.1.md` al integrador; sin ella, Axel podría no verlo en la próxima puerta.
4. **N-242.** Está definida en §3.24 (P-22), pero solo se usa desde P-21, y la matriz la lista ahí. Es orden del documento, viene de v1.0 y no afecta la cobertura.
5. **Pendiente para C1 v1.1** (no es de U1):
   - operaciones para N-316, N-318 y N-319;
   - los campos de K18;
   - `period` en las operaciones de la lista (U1-16);
   - `scope` (K2), `family` (K5), y `Idempotency-Key` con fecha de captura (K9).
   Después hay que volver a correr `verificar_u1.py`: las líneas PENDIENTE deben desaparecer.

## 8. Para el integrador
1. Sumar la adenda K18 de `correccion-U1-r2.md` a «Para C1» de `U1-1.1.md`, o entregar ambos archivos a C1 v1.1 (B6). `C1.md` no menciona K18.
2. Agregar a «Para C1»: «N-289: quitar `updatePerson`» (B7).
3. Registrar `U1-lista-vigente-domingo` en `INDICE.md` como abierto.
4. Pedir a S1 que decida cómo representar «falta capturar» del IVA en `bank_card` y `store_credit` (§7.1).
