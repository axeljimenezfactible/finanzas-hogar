# U4 · Ajustes estilo Notion · documento unido (partes a y b)

Unión de `docs/design/U4-ajustes-a.md` (parte a, versión 1.6: convenciones comunes §0 y P-26.01 a P-26.07) y `docs/design/U4-ajustes-b.md` (parte b, versión 1.5: P-26.08 a P-26.14). La hizo el agente `integrador-wf1` el 11-oct-2026 (flujo WF-1) y la rehicieron los agentes `U4-puerta2` y `correccion-U4-r1` el mismo día (P2c) con el mismo procedimiento (`docs/estado/wf1/INTEGRACION.md` §2.5). Debajo del índice va el texto de cada parte **sin cambios**, primero la a y luego la b; las partes siguen siendo las fuentes que se editan y este archivo se rehace concatenándolas (`make verify-wf1` comprueba que coincida byte a byte).

Estado de revisión: al cierre de WF-1, U4 **no quedó aprobado** tras 5 rondas (`docs/estado/wf1/REVISION-U4-r5.md`). En la puerta 2 Axel autorizó corregir R5-1 (Q2-09 A) y respondió los bloqueos que tocan este documento (`docs/estado/puerta2/RESPUESTAS.md`); P2c los aplicó; la ronda 1 de su revisor no lo aprobó (`docs/estado/p2c/REVISION-U4-r1.md`, M-1 a M-5), `correccion-U4-r1` corrigió esos hallazgos y su revisor decide. `docs/estado/bloqueos/S1-parametros-avisos.md` quedó RESUELTO (Q2-11 B, K47): P-26.08 y P-26.09 leen y escriben las claves de `app_setting` de la migración 00004, con una sola variante.

## Índice del documento unido

| Sección (encabezado de segundo nivel) | Parte |
|---|---|
| 0. Convenciones comunes (rigen las 14 secciones, P-26.01 a P-26.14) | a |
| P-26.01 · Cuenta y roles | a |
| P-26.02 · Presupuesto | a |
| P-26.03 · Ingresos | a |
| P-26.04 · Categorías | a |
| P-26.05 · Créditos y condiciones (incluye promociones) | a |
| P-26.06 · Tiendas y precios (incluye mercados y membresía Sam's) | a |
| P-26.07 · Recetas | a |
| Cobertura de operaciones de la etiqueta Ajustes (parte a) | a |
| Parámetros operativos de la parte a (K44) | a |
| Fronteras con la parte b | a |
| SUPUESTOS | a |
| Convenciones de la parte b (formato; lo común rige desde la parte a §0) | b |
| P-26.08 Reportes y frecuencias | b |
| P-26.09 Notificaciones | b |
| P-26.10 Apariencia | b |
| P-26.11 Datos y respaldos | b |
| P-26.12 Personas y reparto | b |
| P-26.13 Salud y metas de ahorro | b |
| P-26.14 Inventario | b |
| Parámetros operativos de la parte b (K44, U4-12) | b |
| Campos de la puerta 1: dónde quedan | b |
| Matriz de operaciones de la parte b (U4-02) | b |
| Cobertura de criterios U4 en la parte b | b |
| SUPUESTOS | b |

---

# U4 · Ajustes estilo Notion · parte a (secciones 1 a 7)

Versión 1.6 (11-oct-2026, P2c, agente `correccion-U4-r1`: hallazgos M-3 y M-5 de `docs/estado/p2c/REVISION-U4-r1.md`; P-26.05 aplica las respuestas de la puerta 2 en las vistas previas 4 (c), 9 y 10: Q2-29 a Q2-31 (a) dejan 4 (c) PROVISIONAL por el bloqueo abierto `S3-formula-tarjeta-pdf`, y Q2-32 A y Q2-28 (a) dejan 9 y 10 FIRMES, con el IVA del préstamo sin capturar PROVISIONAL por `S3-iva-lectura`; la matriz de operaciones suma `getReportSettings` y `updateReportSettings` de C1 v1.3 y `listCreditPurchases`). Versión 1.5 (11-oct-2026, P2c, agente `U4-puerta2`: decisiones de la puerta 2 en `docs/estado/puerta2/RESPUESTAS.md`; Q2-09 A cierra R5-1 de `docs/estado/wf1/REVISION-U4-r5.md` con la regla general de §0.2 para las 14 secciones (listas, interruptores, `Ctrl+/` y arrastre de tablero con vista previa) y la vista previa de P-26.04; bajos R5-2 (P-26.04: gastos fijos sin reparto propio), R5-3 (P-26.05: efecto de las promociones) y R4-2 («cada cuánto» de 1 a 24); Q2-67 opción 1 (P-26.05: «¿la tasa incluye IVA?» en `credit.vat_included` y 5 datos en «Falta capturar»); Q2-69 (a) (comisiones solo como dato); Q2-74, Q2-75 y Q2-77 (K50: rangos en §0.5 y P-26.06, `optimizer.date_window_days` = 7); mensajes V-23 a V-25; esquema 1.4, migración 00004). Versión 1.4 (11-oct-2026, WF-1, agente `correccion-U4-r4`: hallazgo R4-1 de `docs/estado/wf1/REVISION-U4-r4.md`; en §0.2, la excepción abarca los interruptores con vista previa en línea («activo» de un gasto fijo y «Activa» de una meta de ahorro) y la edición múltiple con `Ctrl+/` de esos campos; en P-26.02, «activo» entra a los campos con vista previa 3 (puntos g y h) y `Ctrl+/` muestra el bloque sumado antes de aplicar). Versión 1.3 (11-oct-2026, WF-1, agente `correccion-U4-r3`: hallazgo R3-1 de `docs/estado/wf1/REVISION-U4-r3.md`; en P-26.02, vista previa en línea del gasto fijo con la partida en el mes, los gastos fijos del mes, las partes por persona, la persona a quien se atribuye y el aviso de que P-07 se recalcula; en §0.2, la excepción de las listas de opciones con vista previa en línea). Versión 1.2 (11-oct-2026, WF-1, agente `correccion-U4-r2`: hallazgos R2-1 (puntos 2 a 4) y R2-2 de `docs/estado/wf1/REVISION-U4-r2.md`; en P-26.05, vista previa de las condiciones de tarjeta o tienda también al cambiar un dato ya capturado, con el pago mínimo PROVISIONAL, y vistas previas nuevas 9 (días de corte y de pago) y 10 (tasa del préstamo)). Versión 1.1 (11-oct-2026, WF-1, agente `correccion-U4-r1`: hallazgos H-1, H-2 y H-3 de `docs/estado/wf1/REVISION-U4-r1.md`; §0 pasa a regir las 14 secciones, vista previa en línea y una sola regla de reparto). Versión 1.0 del mismo día, agente `U4a-ajustes`. El integrador une esta parte con `U4-ajustes-b.md` (secciones 8 a 14) en `docs/design/U4-ajustes.md`.

Fuentes: `docs/plan/decisiones.md` §A–§K (K35–K44 incluidas), `docs/criterios/U4.md` 1.2, `docs/design/U1-arquitectura.md` 1.2 (§2.3, §2.5, §2.6, §3.28, §5), extractos de C1 v1.2 (`docs/contratos/extractos/`), `docs/contratos/C1-datos.md` §11.5 y §11.7, inventario `@param` de la migración 00003 (esquema 1.3), `docs/design/DESIGN.md` §2.5 y §4. `DESIGN.md` §7 (componentes) no existía en la 1.0; en la 1.1 ya existe y §0 cita §7.0, §7.12, §7.13, §7.14 y §7.15. Los controles de cada sección se nombran por su función y el integrador los enlaza a §7.

Alcance: P-26.01 Cuenta y roles · P-26.02 Presupuesto · P-26.03 Ingresos · P-26.04 Categorías · P-26.05 Créditos y condiciones (incluye promociones) · P-26.06 Tiendas y precios (incluye mercados y membresía Sam's) · P-26.07 Recetas. Conteo de esta parte: `grep -c '^## P-26\.' docs/design/U4-ajustes-a.md` = 7.

## 0. Convenciones comunes (rigen las 14 secciones, P-26.01 a P-26.14)

Este es el **único** bloque de convenciones de guardado, error, sin conexión, carga, diálogo crítico, vista previa, táctil y teclado para las 14 secciones. La parte b (`U4-ajustes-b.md`) no lo repite ni lo cambia: solo declara su formato de tablas, sus atajos y sus puntos de quiebre.

### 0.1 Cómo se lee cada tabla de campos
Columnas: **Campo** (nombre en pantalla) · **Tipo** · **Valor por defecto (origen)**: solo el decidido en §A–§K o el sembrado por la migración con su decisión; si no hay, dice «sin valor: lo captura Axel» · **Validación → mensaje** (los códigos V-xx están en §0.6) · **Lee / escribe** (`operationId` de C1; «—» solo cuando el dato es calculado o del dispositivo) · **N-xxx** de U1 §3.28.

«Semilla» es un dato precargado (compra del 10-oct-2026, crédito Kueski del oráculo) y **no** es un valor por defecto: un registro nuevo nace sin él.

### 0.2 Edición en el lugar, guardado y vista previa (consistente en las 14)
- Como en Notion (U1 §2.6, DESIGN §7.3): se toca el valor o se pulsa `Enter`, se edita y se guarda solo al salir del campo o con `Enter`; `Esc` deshace la edición en curso y vuelve al valor anterior. Interruptores y listas de opciones **sin** vista previa guardan al elegir. **Regla general de las 14 secciones para lo que tiene vista previa** (Q2-09 A; R5-1 de `docs/estado/wf1/REVISION-U4-r5.md`):
  1. Toda lista de opciones o interruptor con vista previa en línea muestra la vista previa al elegir o al tocar y no guarda todavía: guarda con `Enter` o al salir del campo, y `Esc` vuelve a la opción o al estado guardado. Son, entre otros: periodicidad, persona, reparto, categoría y «activo» de un gasto fijo (P-26.02); reparto por defecto de una categoría (P-26.04); sexo, actividad y objetivo del perfil de salud y «Activa» de una meta de ahorro (P-26.13); modo de inventario (P-26.14).
  2. La edición múltiple con `Ctrl+/` y el arrastre de una tarjeta entre columnas de un tablero agrupado por un campo con vista previa muestran esa vista previa **antes de aplicar**. Con `Ctrl+/`, sumada sobre las filas seleccionadas: `Enter` aplica a todas y `Esc` cancela sin cambiar ninguna. Al soltar una tarjeta, queda en la columna nueva con el bloque de vista previa y el botón «Aplicar» (o `Enter`); `Esc` la devuelve a su columna sin guardar. El arrastre que solo cambia el orden no tiene vista previa.
  3. Qué muestra cada caso lo dice su sección: P-26.02 (vista previa 3, puntos g y h), P-26.04 (el ejemplo de $199.00 con el reparto nuevo y los gastos fijos que cambian), P-26.13 (vista previa (3) y la del perfil) y P-26.14 (estado antes y después por fila y, con varias filas, cuántas pasan a cada estado).
- **Vista previa en línea** (U4-05): en todo campo con efecto numérico que guarda solo, mientras el valor tecleado difiere del guardado aparece debajo del campo un bloque «Antes → Después» con las cifras que cambian. `Enter` o salir del campo guarda (esa es la confirmación); `Esc` descarta y el bloque desaparece. Si la vista previa necesita al servidor y no se puede calcular, el bloque dice «No se pudo calcular el efecto · Reintentar» y el campo no guarda hasta que la vista previa se vea (`Esc` vuelve al valor anterior). Cada sección dice qué muestra su vista previa y de dónde sale cada cifra. Una vista previa marcada PROVISIONAL (por un bloqueo o una solicitud a C1 que cita) muestra lo que sí se puede calcular, rotula lo pendiente como «Vista previa pendiente» y no impide guardar.
- **Excepciones con botón «Confirmar»** (muestran primero su efecto en un diálogo): el interruptor A/B (P-26.05), el alcance de un cambio de presupuesto (P-26.02), el alcance de un cambio de reparto (P-26.12, Q2-68 C), las acciones críticas (§0.4) y los formularios de alta («Nuevo crédito», «Nueva tienda o mercado», etc.), que guardan con «Guardar» o `Ctrl+Enter`.
- Lo calculado (deuda viva, disponible, rendimiento, punto de equilibrio, conteos) nunca se edita en el lugar.

### 0.3 Estados comunes
Un solo token y un solo texto por estado en las 14 secciones (U1 §2.4; DESIGN §7.0, §7.13, §7.14 y §7.15). Las secciones solo nombran qué se carga o qué falta.

| Estado | Cómo se ve |
|---|---|
| Vacío | Una frase de qué falta y la acción primaria que lo llena («Sin créditos · Nuevo crédito»); nunca una pantalla en blanco. |
| Carga | Esqueleto con la forma final del formulario o de la tabla, sin saltos de diseño y sin cifras provisionales; a los 10 s, «Esto está tardando más de lo normal · Reintentar» (U1 §2.4). |
| Guardando | «Guardando…» junto al control mientras dura el envío. |
| Guardado | Junto al control, icono `circle-check` y el texto «Guardado» en `color.role.success.light` (claro) o `color.role.success.dark` (oscuro) durante 2 s (DESIGN §7.13); texto e icono, nunca solo color; el valor nuevo queda en su sitio y «Guardado» se anuncia con `aria-live="polite"`. |
| Restablecido | Tras «Restablecer al valor inicial» (§0.5), el mismo «Guardado» y además el aviso breve «Valor restablecido · Deshacer» (DESIGN §7.15; «Deshacer» reenvía el valor anterior). |
| Error de validación | El valor no se guarda; borde de error, icono `circle-alert` y el mensaje del catálogo §0.6 debajo del campo en `color.role.error.light` o `color.role.error.dark` (DESIGN §7.0); un valor fuera de rango dice el rango. |
| Error al guardar | El campo vuelve a su valor anterior y debajo dice, con `circle-alert` en `color.role.error.light` o `color.role.error.dark`: «No se pudo guardar. Revisa la conexión con el servidor de casa · Reintentar». Lo capturado en un formulario no se pierde. |
| Error al cargar | Franja en la página con `circle-alert` sobre `color.role.error.surface-light` o `color.role.error.surface-dark`: «No se pudo cargar {qué} · Reintentar» (DESIGN §7.14). |
| Sin conexión | Franja fija «Sin conexión con el servidor de casa» con `wifi-off` (U1 §2.4, DESIGN §7.14); lo ya cargado se ve con su hora de carga y en solo lectura; toda escritura (edición en el lugar, interruptores, acciones, «Restablecer») queda deshabilitada con la razón en texto: «Requiere conexión con el servidor de casa». Apariencia (P-26.10) sigue funcionando porque es local (DESIGN §7.13). |

### 0.4 Acción crítica
Diálogo de confirmación con el énfasis de acción crítica de DESIGN §2.5 y §7.12: encabezado con el icono `shield-alert` en `color.critical.light` (claro) o `color.critical.dark` (oscuro), texto que dice qué pasará antes de confirmar y botón de confirmar con `color.critical.700` y texto `color.base.white` en claro; `color.critical.300` con texto `color.critical.950` en oscuro. Una sola vez por diálogo, con verbo explícito que dice qué pasará («Cerrar 3 sesiones»), y un botón «Cancelar» con el foco inicial (nunca se confirma con `Enter` por omisión). En móvil el diálogo se vuelve hoja inferior (DESIGN §7.12). Nunca se usa para error ni advertencia.

### 0.5 Parámetros operativos (K44)
Cada parámetro es una fila con: valor vigente editable, el texto «Valor inicial: …» en la misma fila, unidad, validación de tipo y signo (entero ≥ 0), y la acción «Restablecer al valor inicial», activa solo cuando el vigente difiere del inicial (`is_initial = false`). Lectura: el `GET` de settings de la sección trae `parameters` con `value`, `initial_value`, `is_initial`, `value_type` y `unit`. Cambio: el `PATCH` de la misma sección con `parameters`; «Restablecer» manda la clave en `reset_parameters` (C1-datos §11.7). Rangos aprobados en la puerta 2 (K50; esquema 1.4, migración 00004): `comparator.timeout_seconds` de 1 a 30 (Q2-74); `price.stale_after_days` de 1 en adelante y `price.zero_weight_after_days` igual o mayor que `price.stale_after_days` (Q2-77; esta regla entre dos claves no está en el esquema: la validan el motor de precios, `Settings.Validate`, y la API); día de la semana de 1 a 7 y día del mes de 1 a 31 (claves `reports.*` y `alerts.*`, Q2-11 B, K47); las horas son texto «HH:MM» de 24 h, hora del centro de México, de 00:00 a 23:59 (tipo de hora que C1 v1.3 agrega a `value_type`, solicitud de S1 en `docs/estado/p2c/S1-puerta2.md`). Las demás claves: entero ≥ 0, sin otros mínimos ni máximos de seguridad. Con los valores iniciales el sistema reproduce las cifras canon.

### 0.6 Catálogo de mensajes de validación (es-MX, sin códigos técnicos)
| Id | Regla | Mensaje |
|---|---|---|
| V-01 | Obligatorio | «Escribe un valor para guardar.» |
| V-02 | Entero ≥ 0 | «Escribe un número entero, de 0 en adelante.» |
| V-03 | Dinero ≥ 0, dos decimales | «Escribe una cantidad en pesos, por ejemplo $1,234.56.» |
| V-04 | Dinero > 0 | «La cantidad debe ser mayor que $0.00.» |
| V-05 | Porcentaje con hasta 4 decimales, ≥ 0 | «Escribe un porcentaje de 0 en adelante, con hasta 4 decimales (por ejemplo 25.8685).» |
| V-06 | Porcentaje con hasta 2 decimales (puntos base) | «Escribe un porcentaje con hasta 2 decimales.» |
| V-07 | Fecha válida | «Elige una fecha del calendario.» (en pantalla `10-oct-2026`; al servidor va en ISO) |
| V-08 | Fin ≥ inicio | «La fecha final no puede ser anterior a la inicial.» |
| V-09 | Texto de 1 a 120 caracteres | «Escribe un nombre de 1 a 120 caracteres.» |
| V-10 | Día del mes 1 a 31 | «Elige un día del 1 al 31.» |
| V-11 | Día de la semana | «Elige un día de la semana.» |
| V-12 | Nombre repetido | «Ya existe uno con ese nombre.» |
| V-13 | Contraseña actual incorrecta | «La contraseña actual no es correcta.» |
| V-14 | Contraseña nueva corta | «La contraseña nueva debe tener al menos {mínimo vigente} caracteres.» |
| V-15 | Confirmación distinta | «Las dos contraseñas nuevas no coinciden.» |
| V-16 | Partes con porcentaje de una regla de reparto: suman a lo más 100 %; la persona del resto recibe lo que falta (P-26.12, RN-03, S1-08) | «Las partes con porcentaje suman {suma} %. No pueden pasar de 100 %: lo que falta lo recibe {persona del resto}.» |
| V-17 | Traslape de vigencia | «Esta categoría ya tiene un monto en esos meses. Ajusta el que existe.» |
| V-18 | En uso | «No se puede borrar porque hay registros que la usan. Puedes desactivarla.» |
| V-19 | Sin veredicto de términos | «Este sitio no tiene veredicto de términos favorable; la lectura asistida sigue apagada.» |
| V-20 | Archivo demasiado grande | «El archivo pasa de {máximo vigente}. Elige uno más chico.» |
| V-21 | Plazo entero ≥ 1 | «Escribe un plazo de 1 pago en adelante.» |
| V-22 | Uno de la lista | «Elige una opción de la lista.» |
| V-23 | Entero dentro de un rango (p. ej. «cada cuánto» de 1 a 24; tiempo máximo del comparador de 1 a 30 s) | «Escribe un número entero del {mínimo} al {máximo}.» |
| V-24 | Entero ≥ 1 | «Escribe un número entero, de 1 en adelante.» |
| V-25 | Peso cero ≥ precio viejo (Q2-77) | «Escribe {días para precio viejo} días o más: el peso de un precio no puede llegar a cero antes de que sea viejo.» |

### 0.7 Estructura, barra lateral, vistas y atajos
- **Escritorio (1440 px)**: Ajustes es una página de la barra lateral colapsable (U1 §1); dentro, una columna interna con las 14 secciones y el contenido a la derecha; cada sección es una página con bloques (encabezado, párrafo de propósito, propiedades, bases de datos en línea). **Móvil (412 px)**: lista de las 14 secciones; cada sección abre como página con «‹ Ajustes» arriba; las tablas se muestran como lista de filas apiladas y cada fila abre una hoja; controles de al menos 44 px.
- **Vistas de base de datos** (U1 §2.5; C1 §8, `queryView`): se nombran en cada sección; la vista elegida, filtros y orden se recuerdan por dispositivo.
- **Menú «/»** (U1 §5.2): `/ir` con el nombre de la página abre Ajustes; `/ingreso`, `/credito`, `/receta` y `/salir` llevan a las hojas de P-08, P-19, P-12 y al cierre de sesión. U4 no agrega comandos.
- **Atajos** (solo los de U1 §5.3; U4 no agrega ninguno, así no hay duplicados): `Alt+N` nuevo registro en la base de datos de la sección; `Ctrl+Enter` guardar el formulario abierto o abrir la fila seleccionada; `Ctrl+/` editar a la vez las propiedades de las filas seleccionadas; `↑`/`↓` y `Enter` moverse entre filas y abrir; `Esc` cerrar hoja o deshacer la edición; `Ctrl+K` buscar una sección; `Ctrl+\` barra lateral; `Ctrl+Shift+L` claro y oscuro; `?` lista de atajos. En Android no hay atajos.
- **Claro y oscuro**: todos los colores por token de C2; ningún dato se distingue solo por color (texto, icono o forma de refuerzo).
- **Táctil, teclado y foco** (DESIGN §7.0 y §4.6): todo control tiene área táctil de al menos `size.touch.min` (44 × 44 px); altura `size.control.coarse` (44 px) con puntero grueso y `size.control.fine` (36 px) solo con puntero fino; `space.2` entre objetivos contiguos; **foco visible en todo control** con `:focus-visible` (anillo `border.width.focus` en `ui.border.focus`, separado `focus.offset`), que nunca se quita; la edición en el lugar se abre con `Enter` o con un toque y se cancela con `Esc`.

## P-26.01 · Cuenta y roles

**Propósito.** Ver las dos cuentas (Axel y Ximena, co-titulares con los mismos permisos), dejar claro que Karen participa sin cuenta, cambiar la contraseña propia, ver y cerrar sesiones por dispositivo, preparar «Este dispositivo» y ajustar los parámetros de acceso (K44). No hay roles distintos ni alta de cuentas (decisión A).

### Campos y acciones
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Usuarios: nombre visible, usuario, papel «Co-titular», último acceso | Tabla de solo lectura (2 filas) | Axel y Ximena, co-titulares (decisión A) | — (no editable) | `listUsers` / — | N-270 |
| Karen · «Participante externa, sin cuenta» | Fila informativa | Sin cuenta (decisión A) | — (no hay acción para crearle cuenta) | `listPersons` (`has_account = false`) / — | N-270, N-309 |
| Contraseña actual | Contraseña | — | V-01; incorrecta → V-13 | — / `changeOwnPassword` (`current_password`) | N-273 |
| Contraseña nueva | Contraseña | — | V-01; menor al mínimo vigente `auth.password_min_length` → V-14 | — / `changeOwnPassword` (`new_password`) | N-273 |
| Repite la contraseña nueva | Contraseña (solo en pantalla) | — | Distinta → V-15 | — (se compara en el dispositivo) | N-273 |
| Sesiones: dispositivo, usuario, inicio, último uso, vence, «Esta sesión» | Base de datos (tabla) | Sesiones abiertas de ambos usuarios, último uso descendente (U1 §2.5.b) | — | `listSessions` (o `queryView` `device_session`) / — | N-271 |
| Cerrar esta sesión de la lista | Acción por fila | — | Confirmación simple | — / `revokeSession` | N-272 |
| Cerrar sesiones de un usuario salvo la actual · cerrar todas salvo la actual · cerrar todas | Acción con alcance (3 opciones) | Ninguna preseleccionada | V-22; acción crítica (§0.4) | — / `revokeSessions` (`scope` = `user_except_current`, `all_except_current`, `all`; `user_account_id` con el primero) | N-272 |
| Salir de este dispositivo | Acción | — | Confirmación simple | — / `logout` | N-006 |
| Este dispositivo: «App instalada» sí o no | Indicador (dato del navegador) | — | — | — (Local) | N-274 |
| Certificado raíz de la casa: descargar e instrucciones | Acción de descarga | — | — | `getCaCertificate` / — | N-274 |
| Parámetros de acceso (4) | Ver tabla de abajo | Valores iniciales sembrados | V-02 | `getSecuritySettings` / `updateSecuritySettings` | N-320 |

### Parámetros de acceso (K44)
| Parámetro en pantalla | Clave de `app_setting` | Valor inicial | Validación | «Restablecer» | Decisión |
|---|---|---|---|---|---|
| Intentos fallidos permitidos en la ventana | `auth.max_failed_attempts` | 5 intentos | Entero ≥ 0 (V-02) | Vuelve a 5 | Decisión A, K11 |
| Ventana de los intentos (deslizante, por usuario y dirección) | `auth.lockout_window_minutes` | 15 minutos | Entero ≥ 0 (V-02) | Vuelve a 15 | Decisión A, K11 |
| Duración de la sesión, renovable con el uso | `session.max_age_days` | 90 días | Entero ≥ 0 (V-02) | Vuelve a 90 | K11 |
| Largo mínimo de una contraseña nueva | `auth.password_min_length` | 12 caracteres | Entero ≥ 0 (V-02) | Vuelve a 12 | K6 |

Se leen y guardan juntos (`getSecuritySettings` / `updateSecuritySettings` con `parameters` y `reset_parameters`). Rigen desde el siguiente intento de inicio de sesión, cambio de contraseña o renovación (C1).

**Vista de base de datos.** Sesiones: tabla (por defecto), agrupable por usuario (`device_session`, `group_by = user_account_id`). Usuarios: tabla de 2 filas.
**Edición en el lugar.** Solo los 4 parámetros (guardado al salir del campo) y el formulario de contraseña (botón «Cambiar contraseña»). Usuarios y sesiones no se editan.
**Vista previa del efecto.** Parámetros: «Con 5 intentos en 15 minutos, el sexto intento fallido desde el mismo dispositivo espera hasta que salga de la ventana»; sesión: «Las sesiones nuevas o renovadas vencerán {días} días después de su último uso». Cerrar sesiones: el diálogo cuenta cuántas sesiones y de quién («Cerrar 3 sesiones de Ximena; tendrán que volver a entrar»).
**Estados.** Vacío: no aplica a usuarios (siempre 2); sesiones vacías no ocurre (al menos la actual). Carga: esqueleto de la tabla. Error: el parámetro vuelve a su valor anterior (§0.3). Guardado: «Guardado»; tras cambiar la contraseña: «Contraseña cambiada. Tus otras sesiones siguen abiertas.»
**Acciones críticas.** «Cerrar todas las sesiones» y «Cerrar todas salvo la actual» (§0.4; U4-08).
**Atajos.** `↑`/`↓` y `Enter` en la tabla de sesiones; `Ctrl+Enter` cambia la contraseña con el formulario lleno.
**Móvil (412 px).** Bloques apilados: Usuarios, Karen, Contraseña, Sesiones (filas apiladas con «Cerrar» de 44 px), Este dispositivo, Parámetros de acceso. **Escritorio (1440 px).** Usuarios y sesiones en tablas a todo el ancho del contenido; contraseña y parámetros en dos columnas de propiedades.

## P-26.02 · Presupuesto

**Propósito.** Fijar el monto mensual por categoría y el total (periodo base: mes calendario; la quincena del 1 al 15 y del 16 al último día es subperiodo, decisión B y K8), los gastos fijos recurrentes, la despensa bimestral programada y el interruptor de la proyección de los meses 7 a 12 (K21).

### Campos y acciones
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Periodo base | Texto informativo «Mes calendario; quincenas del 1 al 15 y del 16 al último día» | Mes calendario (decisión B, K8) | — (no editable) | — | N-275 |
| Monto del mes en curso por categoría | Dinero en tabla | Sin valor: lo captura Axel | V-03 | `listBudgets` / `updateBudget` (monto, `scope`, `month` = mes en curso) o `createBudget` si la categoría no tiene monto ese mes | N-275, N-278 |
| Alcance del cambio de monto | Selección de 2 opciones en el diálogo «Confirmar» | **«De este mes en adelante»** preseleccionado (K2) | V-22; ninguna otra opción | — / `updateBudget` `scope` = `onward` o `this_month`; en alta, `createBudget` con vigencia de un mes (`this_month`) o sin fin (`onward`) | N-278 |
| Marca «Solo este mes» y monto que rige el mes siguiente | Etiqueta y dinero de solo lectura | — | — (calculado) | `listBudgets` (`is_single_month`, `following_month_amount_cents`) / — | N-275 |
| Presupuesto total del mes | Dinero (renglón sin categoría) | Sin valor: lo captura Axel | V-03 | `listBudgets` / `updateBudget` o `createBudget` (sin `category_id`) | N-275, N-278 |
| Quitar monto de una categoría | Acción por fila | — | Acción crítica (§0.4): «Quitar el monto de {categoría}» | — / `deleteBudget` | N-278 |
| Si se traslapa con otro monto | — | — | V-17 | `createBudget` responde conflicto | N-278 |
| Gasto fijo: concepto | Texto | Sin valor: lo captura Axel (o el CSV, plan §8.2) | V-01, V-09 | `listPlannedItems` (`kind = fixed_expense`) / `createPlannedItem`, `updatePlannedItem` (`description`) | N-276, N-278 |
| Gasto fijo: categoría | Relación con categorías | Sin valor: lo captura Axel | V-22 | idem (`category_id`) | N-276 |
| Gasto fijo: monto por vez | Dinero | Sin valor: lo captura Axel | V-04 | idem (`amount_cents`) | N-276 |
| Gasto fijo: periodicidad y cada cuánto | Opción (semanal, quincenal, mensual, anual, una vez) y entero | Sin valor: lo captura Axel | V-22; cada cuánto entero de 1 a 24 (V-23; el esquema exige `interval_count` de 1 a 24, 00001) | idem (`frequency`, `interval_count`) | N-276 |
| Gasto fijo: primera fecha (su día es el día del mes) y fecha final opcional | Fechas | Sin valor: lo captura Axel; fin vacío = sin fin | V-07, V-08 | idem (`start_date`, `end_date`) | N-276 |
| Gasto fijo: persona, reparto, forma de pago | Relaciones y opción | Vacío = hogar en general y reparto de la categoría o el de por defecto (C1) | V-22 | idem (`person_id`, `cost_share_rule_id`, `payment_method`) | N-276 |
| Gasto fijo: activo | Interruptor con vista previa en línea (§0.2: no guarda al tocar; vista previa 3, punto g) | Encendido al crear | — | idem (`is_active`; desactivar en vez de borrar) | N-276, N-278 |
| Despensa programada: monto base | Dinero | Sin valor: lo captura Axel | V-04 | `listPlannedItems` / `createPlannedItem`, `updatePlannedItem` (`amount_cents`) | N-277, N-278 |
| Despensa programada: frecuencia | Texto «Bimestral» (mensual cada 2) | Bimestral (decisión B) | — | idem (`frequency = monthly`, `interval_count = 2`) | N-277 |
| Despensa programada: devengo | Texto «Mitad por mes» | 2 meses (decisión B: «mitad por mes») | — | idem (`accrual_months = 2`) | N-277 |
| Despensa programada: día y próxima fecha | Fecha y dato calculado | Sin valor: lo captura Axel | V-07 | idem (`start_date`; `next_date` calculado) | N-277 |
| Despensa programada: línea y plazo con que se paga | Relación con créditos y entero (quincenas) | Sin valor: lo captura Axel | V-22, V-21 | idem (`payment_method = credit`, `credit_id`, `term_count`) | N-277 |
| Repetir la despensa bimestral y el papel (meses 7 a 12) | Interruptor | **Apagado** (K21; `projection.repeat_cadence_enabled = false`, sembrado) | — | `getProjectionSettings` / `updateProjectionSettings` | N-318, N-319 |

**Vista de base de datos.** Presupuesto: tabla (una fila por categoría activa más el total). Gastos fijos: tabla ordenada por día (U1). Despensa programada: bloque de propiedades (una partida).
**Edición en el lugar.** Monto: se toca, se escribe y al salir del campo se abre el diálogo de alcance con «De este mes en adelante» preseleccionado y botón «Confirmar»; `Esc` cancela sin cambiar. Gastos fijos y despensa: guardado automático por campo, después de su vista previa en línea (§0.2) en los campos que cambian el dinero: del gasto fijo, monto por vez, periodicidad, cada cuánto, primera fecha, fecha final, persona, reparto, activo y la categoría cuando la partida no tiene reparto propio (vista previa 3); de la despensa, monto base, día, línea y plazo (vista previa 4). En los de lista (periodicidad, persona, reparto y categoría), elegir una opción muestra la vista previa y no guarda todavía: `Enter` o salir del campo guarda y `Esc` vuelve a la opción guardada. El interruptor «activo» tampoco guarda al tocar: muestra la vista previa (punto g); `Enter` o salir del interruptor guarda y `Esc` lo devuelve a como estaba. Edición múltiple (`Ctrl+/`, categoría o activo de varios gastos fijos): antes de aplicar muestra la vista previa 3 sumada sobre las partidas seleccionadas (punto h); `Enter` aplica a todas y `Esc` cancela sin cambiar ninguna. El interruptor de la proyección de los meses 7 a 12 guarda al tocar (vista previa 2).
**Vista previa del efecto.** (1) Cambio de monto: el diálogo dice «De este mes en adelante: {mes en curso} y los siguientes quedan en {monto nuevo}» o «Solo este mes: {mes en curso} queda en {monto nuevo} y {mes siguiente} vuelve a {monto actual}», con el total del mes antes y después. (2) Interruptor de proyección: «Encendido: los meses 7 a 12 de Flujo y proyección agregan la despensa bimestral a crédito y el papel alterno de $199.00. Apagado: solo compromisos contratados, gastos fijos e ingresos del CSV» (K21); P-07 se recalcula al volver a abrirla. (3) **Gasto fijo** (vista previa en línea, §0.2, en monto por vez, periodicidad, cada cuánto, primera fecha, fecha final, persona, reparto, categoría sin reparto propio y activo, y en la edición múltiple con `Ctrl+/`): «Antes → Después» calculado en el navegador, sin llamar al servidor y en centavos enteros, con lo que la sección ya cargó (`listPlannedItems` con `kind = fixed_expense`, `listCategories` y `listCostShareRules` con sus `parts`). Muestra: (a) **sus próximas 3 fechas**, con la regla de fechas del esquema (00001, comentarios de `planned_item.frequency`, `interval_count`, `start_date` y `end_date`): quincenal, en los días de pago 15 y 30 o último día del mes (RN-06) desde la primera fecha; semanal, mensual y anual, cada {cada cuánto} semanas, meses o años desde la primera fecha, con el día recortado al último del mes en meses cortos; «una vez», solo la primera fecha; ninguna después de la fecha final; (b) **la partida en el mes**: «En {mes en curso}: {n} veces × {monto por vez} = {total de la partida}», antes y después; en quincenal con monto distinto en la quincena del 30 (`second_half_amount_cents`), la suma de los montos de las fechas del mes; si ni antes ni después cae una fecha en el mes en curso, el bloque usa el primer mes con alguna fecha y lo nombra; (c) **«Gastos fijos de {mes}: {antes} → {después}»**: la suma de las partidas activas de la tabla de gastos fijos con fecha en ese mes; cada fecha cuenta su monto completo en su mes porque el gasto fijo devenga en un mes (`accrual_months` = 1, valor del esquema que esta tabla no cambia); una partida con otro devengo cuenta la porción del mes que fija RN-02; (d) **las partes por persona en {mes}**: cada fecha será un gasto y se reparte como tal (RN-03), con la regla de la partida o, si no tiene, la de su categoría o la de por defecto (por eso la vista previa aparece también al cambiar la categoría de una partida sin reparto propio); en cada fecha, una parte por porcentaje vale ROUND_HALF_UP(monto × `pct_bp` ÷ 10,000), una parte fija vale `fixed_cents`, y la persona del resto (`remainder_person_id`, el titular) toma lo que falta y absorbe los centavos (K14); luego se suman las fechas del mes: «{persona de cada parte} {antes} → {después} · Hogar (el resto) {antes} → {después}». Ejemplo: $199.00 por vez con Conjunto 50/50 da Karen $99.50 y Hogar $99.50 por vez, como en P-26.04. Si una parte fija pasa del monto por vez, en lugar de las partes dice «La parte fija de {persona} ({monto fijo}) pasa del monto por vez ({monto}); así no se puede repartir el gasto (RN-03).»; (e) **la persona**: «Se atribuye a {persona antes} → {persona después}» (vacío = hogar en general), con los gastos fijos de {mes} de cada una antes y después; el total del mes no cambia; (f) **el aviso**: «Flujo y proyección (P-07) se recalcula al volver a abrirla: los gastos fijos entran en sus 12 meses con el interruptor de los meses 7 a 12 encendido o apagado» (K21, RN-28); (g) **al tocar «activo»**: (c) con la partida sumada (al encender) o quitada (al apagar) de los gastos fijos del mes que nombra (b), (d) con sus partes sumadas o quitadas, y (f); (a) no cambia y (b) dice «No cuenta mientras esté inactiva» del lado apagado; si la partida ya no tiene fechas por delante, el bloque dice «Sin fechas por delante: los gastos fijos no cambian»; (h) **edición múltiple con `Ctrl+/`** (categoría o activo de las {n} partidas seleccionadas): antes de aplicar, el mismo bloque sumado sobre las {n} en el mes en curso: (c) «Gastos fijos de {mes}: {antes} → {después}», (d) las partes por persona con todas las fechas de las {n} en ese mes, y (f); al cambiar la categoría, (c) no cambia y solo cambian las partes de las partidas sin reparto propio, y el bloque lo dice: «{k} de {n} partidas toman el reparto de su categoría; las demás conservan el suyo». `Enter` aplica a las {n} (un `updatePlannedItem` por partida, igual que al editar una sola) y `Esc` cancela sin cambiar ninguna; una partida que no se guarde vuelve a su valor anterior con el error de §0.3. Mientras falten la periodicidad o la primera fecha, el bloque muestra lo que se puede calcular y dice «Completa la periodicidad y la primera fecha para ver las cifras». (4) **Despensa programada** (monto base, día, línea o plazo; vista previa en línea, §0.2): «Antes → Después» con (a) las próximas 3 fechas de compra, calculadas en el navegador igual que las del gasto fijo (mensual cada 2 desde el día elegido); (b) con la línea y el plazo elegidos, el calendario de pagos de la próxima compra: pago por quincena, fechas de pago y total con intereses, de `previewCreditPurchase` (no guarda nada; cuerpo: `credit_id` = línea, `kind = installment_plan`, `purchase_date` = próxima fecha, `principal_cents` = monto base, `term_count` = plazo, `interest_source = term_table`; el motor S3 aplica la tabla de tasas, RN-08, y la serie de RN-07); (c) «El disponible de {línea} baja {monto base} el día de la compra» (RN-09: cada compra resta su capital). Con la semilla (Kueski, $3,094.50 a 4 quincenas con 25.8685 %, compra del 10-oct-2026): pagos el 30-oct, 15-nov, 30-nov y 15-dic-2026 de $973.75 cada uno, total $3,895.00 = $3,094.50 de capital + $800.50 de intereses (oráculo: `despensas[0].pagos`, `capital` 3094.5, `intereses` 800.5). Sin línea o sin plazo, el bloque muestra solo (a) y dice «Elige la línea y el plazo para ver los pagos».
**Estados.** Vacío: «Sin montos este mes · Fijar monto» y «Sin gastos fijos · Nuevo gasto fijo» (o importar el CSV desde Datos y respaldos). Carga, error y guardado según §0.3.
**Atajos.** `Alt+N` nuevo gasto fijo; `Ctrl+/` cambiar a la vez categoría o activo de varios gastos fijos, con la vista previa 3 (punto h) antes de aplicar.
**Móvil (412 px).** Lista de categorías con su monto (toque abre teclado numérico y luego el diálogo de alcance a pantalla completa); gastos fijos en filas apiladas; el interruptor arriba de todo con su explicación. **Escritorio (1440 px).** Tabla de presupuesto y tabla de gastos fijos lado a lado; la despensa programada como bloque de propiedades debajo.

## P-26.03 · Ingresos

**Propósito.** Capturar el ingreso quincenal esperado de Axel (quincena del 15 y del 30 o último día), desde cuándo rige, y recordar que los ingresos extra se registran a mano en Ingresos (P-08). Ximena no aporta ingresos al modelo (decisión B).

### Campos y acciones
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Persona | Texto fijo «Axel» | Axel, titular (decisión B; el esquema solo admite al titular en ingresos) | — (no editable) | `listPlannedItems` (`kind = income`, `person_id`) / — | N-279 |
| Frecuencia | Texto fijo «Quincenal: días 15 y 30 (o último día)» | Quincenal (decisión B) | — | idem (`frequency = fortnightly`) | N-279 |
| Monto de la quincena del 15 | Dinero | Sin valor: lo captura Axel (o el CSV) | V-04 | `listPlannedItems` / `createPlannedItem`, `updatePlannedItem` (`amount_cents`) | N-279, N-280 |
| Monto de la quincena del 30 (si es distinto) | Dinero opcional | Vacío = igual a la del 15 (C1) | V-04 si se escribe | idem (`second_half_amount_cents`) | N-279, N-280 |
| Rige desde | Fecha | Sin valor: lo captura Axel | V-01, V-07 | idem (`start_date`) | N-279, N-280 |
| Ximena | Fila informativa «No aporta ingresos al modelo» | Decisión B | — (sin campo) | — | N-279 |
| Registrar ingreso extra | Enlace a la hoja de P-08 (`/ingreso`) | — | Las de P-08 | — / `createIncome` (en P-08) | N-101, N-102 |

**Vista de base de datos.** Tabla (una fila por partida de ingreso, normalmente una).
**Edición en el lugar.** Los dos montos y la fecha, con guardado automático.
**Vista previa del efecto.** Antes de guardar: «Desde {fecha}: quincena del 15 {monto}, quincena del 30 {monto}; ingreso esperado del mes {suma}» y aviso de que Flujo y proyección (P-07) y Presupuesto (P-06) se recalculan.
**Estados.** Vacío: «Sin ingreso esperado · Capturar ingreso quincenal». Carga, error y guardado según §0.3.
**Atajos.** `Ctrl+Enter` guarda la partida abierta.
**Móvil (412 px).** Dos campos de dinero grandes y la fecha; el enlace a ingresos extra debajo. **Escritorio (1440 px).** Bloque de propiedades con los dos montos lado a lado.

## P-26.04 · Categorías

**Propósito.** Mantener las categorías de gasto: nombre, jerarquía, tipo (fijo o variable), reparto por defecto, orden y archivo. Incluye «Comida fuera» (decisión B) y «Despensa», que usan los cálculos.

### Campos y acciones
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Nombre | Texto | Sin valor: lo captura Axel; semilla: «Despensa» y «Comida fuera» (decisión B) | V-01, V-09, V-12 | `listCategories` (o `queryView` `category`) / `createCategory`, `updateCategory` (`name`) | N-281, N-282 |
| Clave que usan los cálculos | Texto de solo lectura (`groceries` = Despensa, `eating_out` = Comida fuera) | Vacía en las categorías de la casa (C1) | — | `listCategories` (`code`) / — | N-281 |
| Categoría padre | Relación con categorías | Vacía = raíz | V-22; no puede ser ella misma | idem (`parent_id`) | N-281, N-282 |
| Tipo | Opción: Fijo, Variable | Vacío = sin clasificar (C1) | V-22 | idem (`kind`) | N-281, N-282 |
| Reparto por defecto | Relación con las reglas de reparto (Hogar, Karen, Conjunto 50/50 y las que existan) | Vacío = la regla por defecto (C1 `is_default`); Despensa: Conjunto 50/50 (decisión B) | V-22 | `listCostShareRules` / `createCategory`, `updateCategory` (`default_cost_share_rule_id`) | N-281, N-282 |
| Orden | Arrastrar la fila | El orden actual | — | — / `reorderCategories` (lista completa de ids) | N-282 |
| Activa | Interruptor; apagar = archivar | Encendida al crear | Una categoría con gastos no se borra, se archiva (C1) | idem (`is_active`) | N-281, N-282 |
| Número de gastos | Entero de solo lectura | — | — (calculado) | `listCategories` (`expense_count`) / — | N-281 |

**Vista de base de datos.** Tabla (por defecto, orden manual) y tablero agrupado por reparto por defecto (`queryView` `category`, `group_by = default_cost_share_rule_id`). Filtro por defecto: activas (U1 §2.5.b).
**Edición en el lugar.** Todas las propiedades editables en la celda, guardado automático. Arrastrar una fila cambia el orden; arrastrar una tarjeta a otra columna del tablero cambia su reparto por defecto (`updateCategory`) después de la vista previa (§0.2, punto 2). El porcentaje del Conjunto (50/50 por defecto) se edita en Personas y reparto (P-26.12, `updateCostShareRule`).
**Vista previa del efecto.** Al elegir el reparto (lista con vista previa: no guarda hasta `Enter` o al salir del campo, §0.2), al soltar la tarjeta en otra columna del tablero o con `Ctrl+/` sobre el reparto de varias categorías, antes de aplicar: «Ejemplo con un gasto de $199.00 (papel del oráculo): Hogar $99.50 · Karen $99.50» con Conjunto 50/50; con otro porcentaje, las partes de personas externas se redondean ROUND_HALF_UP y el titular toma el resto (K14). Aviso: «Aplica a los gastos que registres desde ahora; los guardados conservan sus partes.» Si la categoría tiene gastos fijos sin reparto propio (`listPlannedItems`), además: «{n} gastos fijos sin reparto propio cambian sus partes: {persona} {antes} → {después}», con la regla de la vista previa 3 (d) de P-26.02 (R5-2). Con `Ctrl+/` sobre {n} categorías: «{n} categorías pasan a {reparto}», el mismo ejemplo de $199.00 y los gastos fijos que cambian. Archivar: «{n} gastos conservan esta categoría; ya no aparecerá al registrar.»
**Estados.** Vacío: «Sin categorías · Nueva categoría». Carga, error y guardado según §0.3. Archivar «Despensa» o «Comida fuera» muestra antes: «Esta categoría la usan los cálculos de despensa o de comida fuera.»
**Atajos.** `Alt+N` nueva categoría; `Ctrl+/` cambiar a la vez tipo, reparto o activa de varias (el reparto, con la vista previa sumada antes de aplicar, §0.2).
**Móvil (412 px).** Filas apiladas con nombre, tipo y reparto; reordenar con el asa de 44 px. **Escritorio (1440 px).** Tabla con todas las columnas; selector tabla/tablero arriba a la derecha.

## P-26.05 · Créditos y condiciones (incluye promociones)

**Propósito.** Dar de alta y mantener las líneas de los 4 tipos (línea revolvente tipo Kueski, tarjeta bancaria, préstamo personal a plazo fijo, crédito de tienda con abonos; decisión C) con sus condiciones; el interruptor A/B con vista previa; la tabla editable de tasas por plazo; el calendario de promociones; las restricciones duras y los pesos del comparador (K16, K17) y su tiempo máximo (K43, K44). La tasa real de cada compra se captura en la compra (P-19, decisión C): Kueski cotiza total e intereses al pagar y sus intereses ya incluyen IVA.

### Bloque 1 · Líneas de crédito
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Nombre | Texto | Sin valor: lo captura Axel; semilla «Kueski Pay» | V-01, V-09, V-12 | `listCredits`, `getCredit` (o `queryView` `credit`) / `createCredit`, `updateCredit` (`name`) | N-283, N-289 |
| Tipo | Opción de 4: Línea revolvente, Tarjeta bancaria, Préstamo a plazo fijo, Crédito de tienda | Se elige al crear; sin preselección | V-22 | idem (`kind`) | N-283, N-289 |
| Institución | Texto | Sin valor: lo captura Axel | V-09 | idem (`lender`) | N-283 |
| Titular | Texto fijo «Axel» | Axel (decisión C: créditos del titular) | — | idem (`holder_person_id`) | N-283 |
| Límite | Dinero | Sin valor: lo captura Axel; semilla Kueski $10,000.00 (oráculo R-63) | V-03 | idem (`credit_limit_cents`) | N-283, N-289 |
| Frecuencia de pagos | Opción: semanal, quincenal (días 15 y 30 o último), mensual | Sin valor: lo captura Axel | V-22 | idem (`payment_frequency`) | N-283 |
| Días mínimos entre compra y primer pago | Entero | Sin valor: lo captura Axel | V-02 | idem (`first_payment_min_days`) | N-283 |
| Comisión anual, por pago tardío y por apertura | Dinero (3 campos) | Sin valor: lo captura Axel | V-03 | idem (`annual_fee_cents`, `late_fee_cents`, `opening_fee_cents`) | N-283 |
| **Revolvente:** lectura del disponible (interruptor A/B) | Dos opciones con explicación: «B: el límite es el tope total» y «A: el límite está libre hoy y los pagos suman encima» | **B** (decisión C) | V-22; confirma con vista previa | `getCredit`, `previewAvailability`, `getAvailabilitySeries` / `updateCredit` (`availability_reading`) | N-284, N-285, N-289 |
| **Revolvente:** «Pendiente confirmar con el acreedor» | Aviso en `color.role.warning` con icono y texto, y casilla «Ya lo confirmé» | Sin confirmar: aviso visible (decisión C) | — | `getCredit` / `updateCredit` (`availability_confirmed`) | N-284 |
| **Revolvente:** fecha de arranque de la lectura A | Fecha | Sin valor: lo captura Axel; semilla Kueski 10-oct-2026 (oráculo) | V-07 | idem (`reading_a_anchor_date`) | N-283 |
| **Tarjeta y tienda:** tasa anual, como se capturó | Porcentaje con hasta 4 decimales; la pantalla lo muestra con los decimales tecleados (K12) | Sin valor: lo captura Axel de su estado de cuenta (K18) | V-05 | idem (`annual_rate_num`, `annual_rate_den`, `annual_rate_decimals`; se lee `annual_rate_pct`) | N-283, N-289 |
| **Tarjeta y tienda:** base de la tasa | Opción: diaria, mensual | Sin valor: lo captura Axel (K18) | V-22 | idem (`rate_base`) | N-283, N-289 |
| **Tarjeta y tienda:** ¿la tasa incluye IVA? y, si no, IVA que se suma | Opción: «Sí, ya incluye IVA» o «No, se suma aparte»; con «No», el porcentaje que se suma | Sin valor: lo captura Axel (K18); mientras esté vacío, «Falta capturar» nombra el IVA (Q2-67 opción 1; bloqueo `docs/estado/bloqueos/U4-iva-credito.md` RESUELTO) | V-22; con «No», porcentaje V-06; al elegir «Sí» el porcentaje se vacía (el esquema no admite porcentaje con «Sí»); «No» sin porcentaje: PROVISIONAL, bloqueo `docs/estado/bloqueos/S3-iva-lectura.md` (pregunta 2) | idem (`vat_included`: 1 = sí, 0 = no, vacío = falta capturar; `vat_on_interest_bp` solo con 0); `vat_included` es campo nuevo de C1 v1.3 (solicitud de S1, `docs/estado/p2c/S1-puerta2.md`) | N-283, N-289 |
| **Tarjeta y tienda:** porcentaje de pago mínimo | Porcentaje con hasta 2 decimales | Sin valor: lo captura Axel (K18) | V-06 | idem (`min_payment_pct_bp`) | N-283, N-289 |
| **Tarjeta y tienda:** mínimo fijo | Dinero | Sin valor: lo captura Axel (K18) | V-03 | idem (`min_payment_fixed_cents`) | N-283, N-289 |
| **Tarjeta:** día de corte y día límite de pago | Día del mes (2 campos) | Sin valor: lo captura Axel | V-10 | idem (`statement_day`, `payment_due_day`) | N-283 |
| **Préstamo:** tasa anual, como se capturó | Porcentaje con hasta 4 decimales (K12) | Sin valor: lo captura Axel | V-05 | idem (`annual_rate_*`) | N-283 |
| «Falta capturar: …» | Etiqueta en `color.role.warning` con icono, que nombra los datos faltantes | Visible mientras falte alguno (K18) | — (calculado) | `listCredits`, `getCredit` (`missing_terms`: 5 datos con `vat`, RN-26 1.4 y C1 v1.3) / — | N-283 |
| Activa | Interruptor | Encendida al crear | — | idem (`is_active`) | N-283 |
| Deuda viva, disponible, próximo pago, compras vivas | Cifras de solo lectura | — | — (motor S3) | `listCredits` (`live_debt_cents`, `available_cents`, `next_payment_date`, `next_payment_cents`, `live_purchase_count`) / — | N-283 |

### Bloque 2 · Tasas por plazo (líneas revolventes)
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Plazo (quincenas) | Entero | Sin valor en filas nuevas | V-21; plazo repetido → V-12 | `listCreditTermRates` / `createCreditTermRate`, `updateCreditTermRate` (`term_count`) | N-286, N-289 |
| Tasa por todo el plazo, como se capturó | Porcentaje con hasta 4 decimales; se muestra con sus decimales (K12) | Semilla Kueski (decisión C): **2 quincenas = 17.27 %**, **4 quincenas = 25.8685 %**; filas nuevas sin valor | V-05 | idem (`rate_num`, `rate_den`, `rate_decimals`; se lee `rate_pct`); los cálculos usan la fracción exacta | N-286, N-289 |
| Origen | Opción: cotización real, tecleada | Tecleada en filas nuevas | V-22 | idem (`source`) | N-286 |
| Fecha de la cotización | Fecha | Sin valor | V-07 | idem (`quote_date`) | N-286 |
| Nota | Texto | Vacía | — | idem (`notes`) | N-286 |
| Eliminar tasa | Acción por fila | — | Acción crítica (§0.4): «Eliminar la tasa de {n} quincenas» | — / `deleteCreditTermRate` | N-289 |

### Bloque 3 · Calendario de promociones
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Comercio | Texto (como lo anuncia la promoción) | Sin valor: lo captura Axel (decisión C, captura a mano) | V-01, V-09 | `listPromotions` (o `queryView` `promotion`) / `createPromotion`, `updatePromotion` (`merchant_name`) | N-287, N-289 |
| Tienda registrada | Relación con tiendas, opcional | Vacía | V-22 | idem (`store_id`) | N-287 |
| Línea | Relación con créditos | Vacía = cualquiera (C1) | V-22 | idem (`credit_id`) | N-287 |
| Plazo (pagos) | Entero | Sin valor | V-21 | idem (`term_count`) | N-287 |
| Tasa promocional por todo el plazo | Porcentaje como se capturó; 0 se muestra «Sin intereses» (K12) | Sin valor | V-05 | idem (`rate_num`, `rate_den`, `rate_decimals`; se lee `rate_pct`) | N-287 |
| Compra mínima | Dinero | Sin valor | V-03 | idem (`min_purchase_cents`) | N-287 |
| Vigencia: inicio y fin (inclusive) | Fechas | Sin valor | V-01, V-07, V-08 | idem (`valid_from`, `valid_to`) | N-287 |
| Fuente o nota | Texto | Vacía | — | idem (`source_note`) | N-287 |
| Activa | Interruptor | Encendida al crear | — | idem (`is_active`) | N-287 |
| Vigente hoy | Etiqueta de solo lectura | — | — (calculado) | `listPromotions` (`is_current`) / — | N-287 |
| Eliminar promoción | Acción | — | Acción crítica (§0.4); si hay compras que la usan → V-18 y ofrece desactivar | — / `deletePromotion` | N-289 |

### Bloque 4 · Comparador: restricciones duras, pesos y tiempo máximo
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Colchón mínimo de disponible (en la línea que se usa) | Dinero | **Sin valor: lo captura Axel** (K16) | V-03 | `getOptimizerSettings` / `updateOptimizerSettings` (`min_available_buffer_cents`) | N-288, N-289 |
| Máximo de compras vivas a la vez (en el crédito que se usa) | Entero | **Sin valor: lo captura Axel** (K16) | V-02 | idem (`max_live_purchases`) | N-288, N-289 |
| Deuda en cero antes de (en el crédito que se usa) | Fecha | **Sin valor: lo captura Axel** (K16) | V-07 | idem (`zero_debt_by_date`) | N-288, N-289 |
| Ventana de fechas de los escenarios («fecha ± días») | Parámetro K44 (tabla de abajo) | 7 días (Q2-75, K50; clave `optimizer.date_window_days`, 00004) | V-02 | idem (`date_window_days`; desde C1 v1.3 como clave en `parameters`, solicitud de S1) | N-288, N-289 |
| «Faltan restricciones: …» | Aviso en `color.role.warning` con icono que nombra las faltantes; el comparador no recomienda mientras falte alguna (K16) | Visible al inicio (nada capturado) | — (calculado) | `getOptimizerSettings` (`missing_constraints`) / — | N-288 |
| Pesos de los 4 objetivos: menores intereses, pago por quincena parejo y bajo tope, máximo disponible, comprar más por peso | 4 enteros; la pantalla muestra además el peso normalizado en % | 1, 1, 1, 1: iguales y normalizados (decisión C; sembrado en `optimizer.weights`) | V-02; los 4 en cero → «Al menos un peso debe ser mayor que 0.» | idem (`weights`: `lowest_interest`, `even_low_payment`, `max_available`, `most_value_per_peso`) | N-288, N-289 |
| Pesos iguales | Acción | — | — | — / `updateOptimizerSettings` (`weights` = 1, 1, 1, 1) | N-289 |
| Tope de pago por quincena por persona | Enlace de solo lectura a Personas y reparto | Sin valor (K13); no es restricción dura (decisión C) | — | `listPersons` (`fortnight_payment_cap_cents`) / se edita en P-26.12 | N-309 |
| Tiempo máximo del comparador | Parámetro K44 (tabla de abajo) | 5 segundos (K43; rango Q2-74, K50) | V-23, de 1 a 30 s | `getOptimizerSettings` / `updateOptimizerSettings` | N-321 |

| Parámetro en pantalla | Clave de `app_setting` | Valor inicial | Validación | «Restablecer» | Decisión |
|---|---|---|---|---|---|
| Tiempo máximo del comparador (si se pasa, responde «no disponible») | `comparator.timeout_seconds` | 5 segundos | Entero de 1 a 30 (V-23) | Vuelve a 5 (`reset_parameters`) | K43; rango Q2-74 (K50) |
| Ventana de fechas de los escenarios («fecha ± días») | `optimizer.date_window_days` | 7 días | Entero ≥ 0 (V-02) | Vuelve a 7 (`reset_parameters`) | Q2-75 (K50) |

Las restricciones son globales y se ajustan en cada comparación desde P-20 sin cambiar las de aquí (K16).

**Vista de base de datos.** Líneas: tabla (por defecto en Ajustes), galería (tarjeta con deuda, disponible y próximo pago) y tablero por tipo (`queryView` `credit`, `group_by = kind`). Tasas: tabla ordenada por plazo (`credit_term_rate`). Promociones: **calendario por defecto** con franja de inicio a fin, tabla ordenada por inicio y tablero por línea (`queryView` `promotion`, `date_field = valid_from`, `group_by = credit_id`); filtro por defecto: vigentes y futuras.
**Edición en el lugar.** Condiciones, tasas, promociones, restricciones y pesos: guardado automático por campo; límite, frecuencia de pagos, días mínimos al primer pago, comisiones, tasas por plazo, condiciones de tarjeta o tienda (tasa, base, IVA, porcentaje de pago mínimo y mínimo fijo), días de corte y de pago de la tarjeta y tasa del préstamo guardan después de su vista previa en línea (§0.2; puntos 2 y 4 a 10 de abajo). Una vista previa rotulada «Vista previa pendiente» no impide guardar. Las tasas se guardan con los decimales tecleados (K12). El interruptor A/B **no** guarda al tocar: abre la vista previa y pide «Confirmar» (U1 §2.6).
**Vista previa del efecto.**
1. **Interruptor A/B (obligatoria).** Antes de confirmar se ven las dos lecturas lado a lado (`previewAvailability` y `getAvailabilitySeries`). Con la semilla del oráculo (Kueski, 10-oct-2026):

| Lectura | Disponible de hoy | Mínimo del horizonte | Serie del disponible (14 valores) |
|---|---|---|---|
| A | $10,000.00 (R-64) | $8,094.50 (R-61) | $10,773.63 · $11,547.26 · $12,320.89 · $8,320.89 · $9,094.50 · $10,094.50 · $11,094.50 · $12,094.50 · $8,094.50 · $9,094.50 · $10,094.50 · $11,094.50 · $12,094.50 · $8,094.50 (R-65) |
| B | $6,905.50 (R-64) | $5,000.00 (R-62) | $7,679.13 · $8,452.76 · $9,226.39 · $5,226.39 · $6,000.00 · $7,000.00 · $8,000.00 · $9,000.00 · $5,000.00 · $6,000.00 · $7,000.00 · $8,000.00 · $9,000.00 · $5,000.00 (R-66) |

   Texto del botón: «Usar la lectura A» o «Usar la lectura B»; el aviso «pendiente confirmar» sigue visible mientras la casilla «Ya lo confirmé» esté apagada.
2. **Tasa por plazo.** Al editar una tasa: «Con {tasa} a {n} quincenas, una compra de $3,094.50 pagaría {intereses} de intereses» (la cotización del oráculo: $800.50 con 25.8685 % a 4 quincenas), calculado con la fracción exacta y ROUND_HALF_UP.
3. **Restricciones y pesos.** Muestra cuántas faltan y, con todo capturado, el peso normalizado de cada objetivo (con 1, 1, 1, 1: 25 % cada uno).
4. **Condiciones de tarjeta o tienda** (tasa anual, base, IVA, porcentaje de pago mínimo y mínimo fijo; vista previa en línea, §0.2), tanto al completarlas como al cambiar una que ya estaba capturada. «Antes» usa los valores guardados y «Después» los tecleados: (a) **«Falta capturar»**: los datos que faltan, antes → después, con la regla de `missing_terms` (K18 y RN-26 1.4: tasa anual, base, IVA, porcentaje de pago mínimo y mínimo fijo) calculada en el navegador; el IVA falta mientras «¿la tasa incluye IVA?» esté vacío (Q2-67 opción 1); al completar el último, «Falta capturar» desaparece; (b) **la tasa**, antes → después, como se capturó (K12); (c) **pago mínimo e intereses del periodo: PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/S3-formula-tarjeta-pdf.md`). En la puerta 2 Axel eligió (a) en Q2-29 a Q2-31 (`docs/estado/bloqueos/S3-formula-tarjeta.md`, preguntas 1 a 3, RESUELTO): S3 escribe la fórmula de intereses y de pago mínimo de la tarjeta y del crédito de tienda a partir de un estado de cuenta real de cada producto, con un ejemplo comprobable. Esos estados de cuenta todavía no se entregan, y mientras tanto S3 usa una fórmula general marcada PROVISIONAL (`docs/contratos/motores/S3-interfaz.md` 1.5 §5.3). Mientras ese bloqueo siga abierto, el bloque **no muestra cifra**: rotula «Pago mínimo e intereses: vista previa pendiente» (§0.2) y no impide guardar. Cuando S3 tenga la fórmula de cada producto, el bloque muestra el pago mínimo y los intereses antes → después con ella (`CardMinimumPayment` y `CardPeriodInterest` del motor S3, `S3-interfaz.md` §5.3), en el navegador si la fórmula lo permite o con la solicitud SC-U4-4 a C1 (en `docs/estado/wf1/correccion-U4-r2.md`). Lo que ya está decidido no depende de estos datos y el bloque lo dice: «Si pagas el total del corte a tiempo, los intereses son $0.00; los meses sin intereses y los abonos sin interés no cambian» (FIRME: `S3-formula-tarjeta-pdf`, «Impacto»).
5. **Límite** (vista previa en línea, §0.2; sin operación nueva). Sea Δ = límite tecleado − límite guardado. **Línea revolvente:** el disponible de hoy y el mínimo del horizonte de las dos lecturas cambian exactamente Δ, porque las dos parten del límite y las compras y los pagos no dependen de él (RN-09: B = límite − capital vivo; A = límite en la fecha ancla, menos las compras posteriores, más el capital repuesto). «Antes» sale de `previewAvailability` (`reading_a` y `reading_b`: `available_today_cents`, `horizon_min_cents`) y de `getAvailabilitySeries` (14 valores por lectura); «Después» es cada cifra + Δ, calculada en el navegador en centavos enteros. Se muestra primero la lectura vigente y debajo la otra. Con la semilla (límite $10,000.00): B hoy $6,905.50 → {$6,905.50 + Δ} y mínimo $5,000.00 → {$5,000.00 + Δ}; A hoy $10,000.00 → {$10,000.00 + Δ} y mínimo $8,094.50 → {$8,094.50 + Δ} (R-64, R-61, R-62); la serie del disponible de P-17 se desplaza Δ en sus 14 valores (R-65, R-66). **Tarjeta:** disponible de hoy = límite − saldo (C1, `available_cents` de `listCredits`), así que pasa de `available_cents` a `available_cents` + Δ. **Préstamo y crédito de tienda:** C1 no les calcula disponible (`available_cents` ausente); el bloque dice «El límite de este tipo de crédito no cambia ninguna cifra calculada».
6. **Frecuencia de pagos** (vista previa en línea). Una compra de ejemplo con `previewCreditPurchase` (no guarda nada), dos veces: con la frecuencia guardada («Antes») y con la tecleada en `payment_frequency` («Después»; `CreditPurchaseCreate` acepta la frecuencia de cada compra). Muestra fechas de los pagos, monto de cada pago y total. Compra de ejemplo: capital y plazo de la última compra de la línea, editables en el mismo bloque; si la línea no tiene compras, los dos campos nacen vacíos y el bloque dice «Escribe un monto y un plazo de ejemplo» (no se inventa un monto). Con la semilla (quincenal, $3,094.50 a 4 quincenas, compra del 10-oct-2026): pagos el 30-oct, 15-nov, 30-nov y 15-dic-2026 de $973.75, total $3,895.00 (oráculo `despensas[0].pagos`; RN-06).
7. **Días mínimos entre compra y primer pago** (vista previa en línea). «Primer pago de una compra hecha hoy: {fecha con el valor guardado} → {fecha con el valor tecleado}», más el calendario de la compra de ejemplo del punto 6, con `previewCreditPurchase`. **PROVISIONAL** (solicitud SC-U4-1 a C1, en `docs/estado/wf1/correccion-U4-r1.md`): el cuerpo de `previewCreditPurchase` no recibe los días mínimos, que el motor toma de la línea guardada; la solicitud agrega el campo opcional `first_payment_min_days`, solo para la vista previa. Mientras no exista, el bloque muestra el calendario con el valor guardado y rotula el «Después» como «Vista previa pendiente» (§0.2). Ejemplo de RN-06 con 15 días: compra del 10-oct-2026 → el pago del 15-oct-2026 queda a 5 días y no cuenta; el primero es el 30-oct-2026.
8. **Comisiones anual, por pago tardío y por apertura.** Hoy ninguna regla de negocio las usa (`docs/contratos/reglas-negocio.md` 1.4: RN-06 a RN-11 no las incluyen), así que no cambian ninguna cifra calculada; la vista previa en línea lo dice: «Se guarda como dato de la línea; hoy no cambia el costo de las compras, la deuda ni el disponible». Así lo decidió Axel en la puerta 2: las comisiones solo se guardan como dato (Q2-69 (a); bloqueo `docs/estado/bloqueos/U4-vista-previa-formulas.md`, pregunta 1, RESUELTO).
9. **Tarjeta: día de corte y día límite de pago** (vista previa en línea, §0.2). «Próximo corte: {fecha con los días guardados} → {fecha con los tecleados} · Pago de ese corte: {fecha} → {fecha}», calculadas en el navegador sin operación nueva con el calendario de la tarjeta del motor S3 (`docs/contratos/motores/S3-interfaz.md` 1.5 §5.3), FIRME desde la puerta 2 (Q2-32 A; bloqueo `docs/estado/bloqueos/S3-formula-tarjeta.md`, pregunta 4, RESUELTO): corte = primera fecha igual o posterior a hoy cuyo día es `statement_day` (si hoy es día de corte, el corte es hoy); pago = primera fecha posterior al corte cuyo día es `payment_due_day`; en un mes que no tiene el día 29, 30 o 31, el corte y el pago caen en el último día del mes. Ejemplo ilustrativo (no canon), hoy 11-oct-2026, con corte el día 31 y pago el día 20: próximo corte 31-oct-2026 y su pago 20-nov-2026; el corte siguiente es el 30-nov-2026, porque noviembre no tiene día 31. Debajo: «El próximo pago de la tarjeta se recalcula al guardar» (el «próximo pago» de solo lectura del bloque, `next_payment_date`, lo calcula el motor S3). No impide guardar.
10. **Préstamo: tasa anual** (vista previa en línea, §0.2). La tasa cambia el pago de cada periodo y los intereses de la tabla de amortización que el motor S3 arma con ella (tabla francesa, `Amortize`, `docs/contratos/motores/S3-interfaz.md` 1.5 §5.4). El bloque muestra la tasa antes → después como se capturó (K12) y, debajo, el pago por periodo y el total de intereses antes → después con la convención FIRME desde la puerta 2 (Q2-28 (a); bloqueo `docs/estado/bloqueos/S3-convencion-prestamo.md` RESUELTO): tasa del periodo r = tasa anual ÷ 12 (mensual), ÷ 24 (quincenal) o ÷ 52 (semanal), según la frecuencia de pagos, y, si el préstamo suma IVA, r × (1 + IVA); pago fijo = ROUND_HALF_UP(capital × r ÷ (1 − (1 + r)^−n)) con fracciones exactas; interés de cada periodo = ROUND_HALF_UP(saldo × r); el último pago ajusta los centavos. «Antes» usa la tasa guardada y «Después» la tecleada, con el capital, el plazo, la frecuencia y la primera fecha de pago del préstamo guardado (`listCreditPurchases`: `principal_cents`, `term_count`, `payment_frequency`, `first_payment_date`). Se calcula en el navegador: no necesita saldo diario ni saldo del corte, así que la solicitud condicional SC-U4-4 no aplica a este bloque. Si el préstamo no tiene su compra registrada, el bloque muestra solo la tasa antes → después. Si el préstamo suma IVA lo dice `vat_included` de `listCredits` (S3 §5.4: 1 = la tasa ya incluye IVA y no se suma; 0 = se suma el de `vat_on_interest_bp`). Esta sección no captura ese dato del préstamo (supuesto 9), y con `vat_included` vacío el motor calcula sin sumar IVA, la lectura en uso: **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/S3-iva-lectura.md`, pregunta 1). No impide guardar.
11. **Promociones** (tasa promocional, plazo, compra mínima y vigencia; guardan solos, §0.2). Debajo del campo, mientras el valor tecleado difiere del guardado: «Se guarda como dato de la promoción; cambia los escenarios del comparador la próxima vez que compares; las compras ya registradas conservan su tasa» (R5-3).

**Estados.** Vacío: «Sin créditos · Nuevo crédito»; «Sin tasas · Nueva tasa»; «Sin promociones · Nueva promoción». Carga, error y guardado según §0.3; el error de la vista previa A/B deja la lectura vigente sin cambio.
**Acciones críticas.** Eliminar tasa y eliminar promoción (§0.4). Cambiar A/B usa «Confirmar» con vista previa (no es borrado, no usa el énfasis crítico).
**Atajos.** `Alt+N` nuevo registro en el bloque con foco (línea, tasa o promoción); `←`/`→` mes anterior o siguiente en el calendario de promociones; `Ctrl+/` editar a la vez varias promociones.
**Móvil (412 px).** Bloques plegables en este orden: Líneas (tarjetas), Interruptor A/B (vista previa a pantalla completa con las dos lecturas una sobre otra), Tasas, Promociones (calendario de lista por mes), Comparador. **Escritorio (1440 px).** Líneas en tabla; al abrir una, su página con propiedades a la izquierda y tasas y A/B a la derecha; promociones en calendario mensual; comparador como bloque de propiedades con los avisos arriba.

## P-26.06 · Tiendas y precios (incluye mercados y membresía Sam's)

**Propósito.** Mantener cadenas (Walmart, Bodega Aurrera, Sam's Club), sucursales y mercados o tianguis (nombre, día de visita y qué se compra ahí; decisión D), la membresía de Sam's (no la tienen: se evalúa con punto de equilibrio; decisión E, K7), las fuentes de precio (decisión E, K22, K37) y los parámetros de precios (K44).

### Bloque 1 · Tiendas y mercados
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Tipo | Opción: cadena, sucursal, mercado o tianguis | Se elige al crear | V-22 | `listStores` (o `queryView` `store`) / `createStore`, `updateStore` (`kind`) | N-290, N-292 |
| Nombre | Texto | Sin valor: lo captura Axel; semilla Walmart, Bodega Aurrera y Sam's Club (decisión D) | V-01, V-09, V-12 | idem (`name`) | N-290, N-292 |
| Cadena de la sucursal | Relación con cadenas (solo sucursal) | Sin valor | V-22 | idem (`parent_id`) | N-290 |
| Nombre en PROFECO | Texto opcional | Vacío | V-09 | idem (`profeco_name`) | N-290 |
| Sitio de la cadena | Texto (solo cadena) | Vacío | — | idem (`website_url`) | N-290 |
| Dirección, municipio o alcaldía y entidad | Texto (3 campos) | Vacío | — | idem (`address`, `municipality`, `state`) | N-290 |
| Día de visita (mercado o tianguis) | Día de la semana; obligatorio en mercado | Sin valor: lo captura Axel (decisión D) | V-01, V-11 | idem (`visit_weekday`, lunes = 1 … domingo = 7) | N-290, N-292 |
| Qué se compra ahí | Selección múltiple de grupos de producto (verduras, carne, etc.) | Sin valor: lo captura Axel (decisión D) | V-22 | `listProductGroups` / `setStoreProductGroups` | N-290, N-292 |
| Exige membresía | Interruptor | Sam's Club: sí (decisión E); las demás: no | — | idem (`requires_membership`) | N-290 |
| La casa tiene la membresía | Interruptor | Sam's Club: **no** (decisión E: «no tienen membresía; evaluarla») | — | idem (`has_membership`) | N-290 |
| Costo anual de la membresía | Dinero | **Sin valor: lo captura Axel** | V-03 | idem (`membership_annual_fee_cents`) | N-290, N-292 |
| Punto de equilibrio de la membresía | Cifra de solo lectura: compras al año que cubren la membresía, ahorro por canasta y anual | — (calculado, K7); sin costo capturado dice «Captura el costo anual para evaluar» | — | `getBasketComparison` (`membership`) / — | N-164, N-290 |
| Lectura asistida del sitio (por cadena) | Interruptor con el veredicto de términos visible (permitido, negado, sin veredicto) | **Apagada** (decisión E) | Sin veredicto favorable en `S5-terminos.md` → V-19 y el interruptor queda bloqueado; activar = acción crítica (§0.4) | `listStores` (`assisted_reading_verdict`) / `updateStore` (`assisted_reading_enabled`) | N-291, N-292 |
| Activa y notas | Interruptor y texto | Activa al crear | — | idem (`is_active`, `notes`) | N-290 |

### Bloque 2 · Fuentes de precio
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Captura manual | Indicador «Siempre activa»; compite con el mismo peso que el ticket (K37) | Activa (decisión E) | — | `getPriceSources` (`manual_enabled`) / — | N-291 |
| Ticket con OCR local | Indicador con el motor (Tesseract 5.5.3, español), «siempre con tu confirmación» | Activo, siempre con confirmación humana (decisión E) | — | `getPriceSources` (`ocr_engine`) / — | N-291 |
| IA para leer tickets | Estado de solo lectura con enlace a Datos y respaldos | **Apagada** (decisión H) | — | `listAiProviders` / se cambia en P-26.11 (`updateAiProvider`) | N-306 |
| Lectura asistida de sitios (interruptor general) | Interruptor; aunque se encienda, cada cadena necesita su veredicto (bloque 1); solo es referencia, no compite (K37) | **Apagada** (decisión E; `price.assisted_reading_enabled = false`, sembrado) | Activar = acción crítica (§0.4) | `getPriceSettings`, `getPriceSources` / `updatePriceSettings` (`assisted_reading_enabled`) | N-291, N-292 |
| PROFECO: último corte, registros y cobertura por cadena y municipio | Datos de solo lectura; «Solo referencia: nunca compite» (K22) | — | — | `getPriceSources` (`profeco_last_import_at`, `profeco_price_count`, `profeco_coverage`), `listJobs` / — | N-291 |
| Importar corte PROFECO | Acción con selector de archivo local (RAR o CSV de datos abiertos) | — | Mayor al máximo vigente → V-20 | — / `importProfeco` (devuelve registros y cobertura) | N-293 |

### Bloque 3 · Parámetros de precios (K44)
| Parámetro en pantalla | Clave de `app_setting` | Valor inicial | Validación | «Restablecer» | Decisión |
|---|---|---|---|---|---|
| Días para que un precio sea viejo | `price.stale_after_days` | 30 días | Entero ≥ 1 (V-24; Q2-77, K50) | Vuelve a 30 | Decisión E, K22 |
| Días en que el peso de un precio viejo llega a cero y deja de competir | `price.zero_weight_after_days` | 90 días | Entero igual o mayor que «Días para que un precio sea viejo» (V-25; Q2-77, K50: lo validan el motor de precios y la API, no el esquema) | Vuelve a 90 | K22, K40 |
| Tamaño máximo de un corte PROFECO | `price.profeco_import_max_gb` | 1 GB | Entero ≥ 0 (V-02) | Vuelve a 1 | K6 |

Lectura y cambio: `getPriceSettings` / `updatePriceSettings` (`parameters`, `reset_parameters`); S5 aplica RN-12 con los valores vigentes (N-322).

**Vista de base de datos.** Tiendas: tabla (por defecto), galería, tablero por tipo (cadena, sucursal, mercado) y calendario por día de visita, recurrente (`queryView` `store`, `group_by = kind`, `date_field = visit_weekday`). Filtro por defecto: activas.
**Edición en el lugar.** Propiedades de cada tienda y parámetros con guardado automático; los interruptores de lectura asistida no guardan al tocar: abren la confirmación crítica.
**Vista previa del efecto.** (1) Frescura con los valores iniciales: «Un precio del 10-oct-2026 queda viejo el 9-nov-2026 (día 30, peso 0.9836); el día 60 pesa 0.4918; el 8-ene-2027 (día 90) su peso llega a 0 y desde el día siguiente solo se muestra marcado» (K22, K40; peso = (90 − días) ÷ 61). Al cambiar los valores, el texto se recalcula con los nuevos y dice cuántos precios cambian de fresco a viejo hoy. (2) Membresía: al capturar el costo anual se muestra el punto de equilibrio. (3) Importar PROFECO: antes de cargar, nombre y tamaño del archivo contra el máximo vigente.
**Estados.** Vacío: «Sin tiendas · Nueva tienda o mercado»; PROFECO sin cortes: «Sin cortes importados · Importar corte PROFECO». Carga, error y guardado según §0.3; la importación muestra su avance como tarea (`listJobs`).
**Acciones críticas.** Activar la lectura asistida (general o de un sitio) (§0.4; U4-08), con texto «Activar la lectura asistida de {sitio}: el sistema leerá su sitio con baja frecuencia según su veredicto de términos».
**Atajos.** `Alt+N` nueva tienda o mercado; `←`/`→` semana anterior o siguiente en el calendario de visitas.
**Móvil (412 px).** Tiendas como tarjetas; fuentes de precio como lista de estados; parámetros al final. **Escritorio (1440 px).** Tabla de tiendas arriba; fuentes y parámetros en dos columnas debajo.

## P-26.07 · Recetas

**Propósito.** Ver y mantener el recetario (recetario inicial escrito por Claude más las recetas de Axel y Ximena, todas editables; decisión D), los tiempos de comida y las porciones base, las bases nutricionales locales y las restricciones alimentarias por persona (el campo existe y va vacío; decisión D).

### Campos y acciones
| Campo | Tipo | Valor por defecto (origen) | Validación → mensaje | Lee / escribe | N-xxx |
|---|---|---|---|---|---|
| Porciones base | Entero de solo lectura «2 porciones; se escala en cada receta» | 2 (decisión D) | — (cada receta cambia las suyas en P-12) | `getRecipeSettings` (`default_base_servings`) / — | N-294 |
| Tiempos de comida | Catálogo fijo de 4: desayuno, comida, cena, colación (opcional) | Los 4 (decisión D) | — (no editable) | `listMealSlots`, `getRecipeSettings` (`meal_slots`) / — | N-294 |
| Recetas por origen | Conteo: recetario inicial y de la casa | — (calculado) | — | `getRecipeSettings` (`recipe_counts`) / — | N-294 |
| Recetario: nombre y tiempos aptos | Base de datos en línea; texto y selección múltiple | Recetario inicial de Claude más las de la casa (decisión D) | V-01, V-09; al menos un tiempo apto → «Elige al menos un tiempo de comida.» | `listRecipes` (o `queryView` `recipe`) / `updateRecipe` | N-142, N-294 |
| Nueva receta · duplicar receta | Acciones (abren P-12) | — | Las de P-12 | — / `createRecipe`, `duplicateRecipe` | N-142 |
| Eliminar receta | Acción | — | Acción crítica (§0.4) con su uso en el plan | `getRecipeUsage` / `deleteRecipe` | N-142 |
| Bases nutricionales: fuente, fecha de carga y número de alimentos | Lista de solo lectura (USDA Foundation Foods 04/2026, Open Food Facts filtrado, tablas mexicanas, etiquetas, manual) | — | — | `getRecipeSettings` (`nutrition_sources`) / — | N-294 |
| Recargar base nutricional | Acción por fuente (USDA u Open Food Facts) desde el archivo local | — | V-22 | — / `reloadNutritionBase` (`source`; devuelve el conteo) | N-295 |
| Restricciones alimentarias de Axel y de Ximena | Texto opcional por persona | **Vacías** (decisión D) | — | `listHealthProfiles` / `updateHealthProfile` (`dietary_restrictions`) | N-311, N-313 |
| Sugerir recetas con IA | Estado de solo lectura y enlace a Datos y respaldos | **Apagada** (decisión H); con IA apagada la acción no aparece | — | `listAiProviders` / se cambia en P-26.11; la sugerencia se pide en P-11 o P-12 (`suggestRecipeWithAi`) | N-306 |

**Vista de base de datos.** Recetario en línea: tabla (por defecto en Ajustes), tablero por tiempo de comida o por origen y galería (`queryView` `recipe`, `group_by = meal_slot_id` u `origin`). Bases nutricionales: lista.
**Edición en el lugar.** Nombre y tiempos aptos en la tabla; restricciones al salir del campo. Ingredientes y pasos se editan en P-12.
**Vista previa del efecto.** Recargar base: «Se reemplazará la base {fuente} cargada el {fecha} ({n} alimentos) por la del archivo local; las kcal de las recetas se recalculan». Eliminar receta: «Se usa en {n} casillas del plan de comidas» antes de confirmar.
**Estados.** Vacío: «Sin recetas · Nueva receta»; base sin cargar: «Sin datos de {fuente} · Recargar base nutricional». Carga, error y guardado según §0.3; la recarga muestra avance y el conteo final.
**Acciones críticas.** Eliminar receta (§0.4).
**Atajos.** `Alt+N` nueva receta; `Ctrl+/` cambiar a la vez los tiempos aptos de varias.
**Móvil (412 px).** Recetario en galería de una columna; datos de solo lectura en lista; restricciones como dos campos. **Escritorio (1440 px).** Recetario en tabla a todo el ancho; bases y restricciones en dos columnas debajo.

## Cobertura de operaciones de la etiqueta Ajustes (parte a)
- **Con campo o acción en las secciones 1 a 7 (36):** `listUsers`, `listSessions`, `revokeSession`, `revokeSessions`, `changeOwnPassword`, `logout`, `getCaCertificate`, `getSecuritySettings`, `updateSecuritySettings` (P-26.01); `getProjectionSettings`, `updateProjectionSettings`, `listPlannedItems`, `createPlannedItem`, `updatePlannedItem` (P-26.02 y P-26.03); `listCategories`, `createCategory`, `updateCategory`, `reorderCategories`, `listCostShareRules` (P-26.04); `createCredit`, `updateCredit`, `createCreditTermRate`, `updateCreditTermRate`, `deleteCreditTermRate`, `createPromotion`, `updatePromotion`, `deletePromotion`, `getOptimizerSettings`, `updateOptimizerSettings` (P-26.05); `createStore`, `updateStore`, `setStoreProductGroups`, `getPriceSources`, `getPriceSettings`, `updatePriceSettings`, `importProfeco` (P-26.06); `getRecipeSettings`, `reloadNutritionBase` (P-26.07). También se leen aquí `listPersons`, `listAiProviders` y `updateHealthProfile`, cuyo lugar principal está en la parte b.
- **Lugar principal en la parte b (secciones 8 a 14):** `previewPlannedItemsImport`, `importPlannedItems`, `listBackups`, `createBackup`, `restoreBackup`, `getSeedStatus`, `getDataSettings`, `updateDataSettings`, `listAiProviders`, `updateAiProvider`, `getSystemStatus` (P-26.11); `createSavingsGoal`, `updateSavingsGoal`, `createHealthProfile`, `updateHealthProfile`, `getHealthSettings`, `updateHealthSettings` (P-26.13); `listPersons`, `updatePerson`, `updateCostShareRule` (P-26.12); `updateInventoryItem` (P-26.14); `updateJob` (P-26.11, hora del respaldo); `listNotifications`, `updateNotification`, `setNotificationRecipients`, `sendTestNotification`, `listNotificationRecipients`, `updateNotificationRecipient`, `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate`, `getAlertSettings`, `updateAlertSettings` (P-26.09); `getReportSettings`, `updateReportSettings` (P-26.08, C1 v1.3).
- **Fuera de las 14 secciones:** `login` (P-00) y `getOpenApiDocument` (excepción técnica, C1-datos §11.2).
- **De otras etiquetas, usadas aquí:** `queryView`, `listBudgets`, `createBudget`, `updateBudget`, `deleteBudget`, `createIncome`, `listCredits`, `getCredit`, `listCreditPurchases` (vista previa 10 de P-26.05), `previewAvailability`, `getAvailabilitySeries`, `listCreditTermRates`, `listPromotions`, `listStores`, `listProductGroups`, `getBasketComparison`, `listJobs`, `listMealSlots`, `listRecipes`, `createRecipe`, `updateRecipe`, `duplicateRecipe`, `deleteRecipe`, `getRecipeUsage`, `suggestRecipeWithAi`, `listHealthProfiles`.

## Parámetros operativos de la parte a (K44)
| Parámetro | Sección | Clave de `app_setting` | Valor inicial | Validación | «Restablecer» | Decisión |
|---|---|---|---|---|---|---|
| Intentos fallidos en la ventana | Cuenta y roles | `auth.max_failed_attempts` | 5 | Entero ≥ 0 | Sí | Decisión A, K11 |
| Ventana de los intentos | Cuenta y roles | `auth.lockout_window_minutes` | 15 min | Entero ≥ 0 | Sí | Decisión A, K11 |
| Duración de la sesión | Cuenta y roles | `session.max_age_days` | 90 días | Entero ≥ 0 | Sí | K11 |
| Largo mínimo de contraseña | Cuenta y roles | `auth.password_min_length` | 12 caracteres | Entero ≥ 0 | Sí | K6 |
| Tiempo máximo del comparador | Créditos y condiciones | `comparator.timeout_seconds` | 5 s | Entero de 1 a 30 | Sí | K43; Q2-74 (K50) |
| Ventana de fechas de los escenarios | Créditos y condiciones | `optimizer.date_window_days` | 7 días | Entero ≥ 0 | Sí | Q2-75 (K50) |
| Precio viejo | Tiendas y precios | `price.stale_after_days` | 30 días | Entero ≥ 1 | Sí | Decisión E, K22; Q2-77 (K50) |
| Peso cero | Tiendas y precios | `price.zero_weight_after_days` | 90 días | Entero ≥ precio viejo | Sí | K22, K40; Q2-77 (K50) |
| Máximo de un corte PROFECO | Tiendas y precios | `price.profeco_import_max_gb` | 1 GB | Entero ≥ 0 | Sí | K6 |

Los otros 7 parámetros sembrados (`attachment.max_size_mb`, `idempotency.retention_hours`, `health.goal_adjust_lose_kcal`, `health.goal_adjust_gain_kcal`, `alerts.budget_warning_pct`, `alerts.budget_limit_pct`, `alerts.payment_due_days_before`) las 8 claves `reports.*` y `alerts.*` de la 00004 (hora, anticipación, día de la semana y día del mes; Q2-11 B, K47; bloqueo `docs/estado/bloqueos/S1-parametros-avisos.md` RESUELTO) y `ocr.timeout_seconds` (Q2-76, K50) pertenecen a las secciones 8, 9, 11 y 13 (parte b). El ajuste de kcal de la casa (K36) va en Salud y metas de ahorro y el de cada persona en su perfil (P-26.13 y P-21): parte b.

## Fronteras con la parte b
- Tope de pago por quincena por persona (K13) y porcentaje del Conjunto (50/50): se editan en P-26.12; aquí solo se leen (P-26.05 y P-26.04).
- Interruptores de IA (decisión H): se cambian en P-26.11; aquí se muestra su estado (P-26.06 y P-26.07).
- Importar `gastos_fijos_ingresos.csv`: P-26.11; aquí solo se enlaza desde los estados vacíos de P-26.02 y P-26.03.
- Restricciones alimentarias: se editan aquí (P-26.07) sobre el perfil de salud; el resto del perfil, en P-26.13.

## SUPUESTOS
1. La despensa programada es una partida de `planned_item` de gasto con categoría Despensa, mensual cada 2 y devengo de 2 meses (decisión B); C1 no fija si su `kind` es `fixed_expense` o `upcoming_expense`, y la pantalla no lo muestra.
2. El monto y el plazo de un préstamo personal se registran como compra a crédito en P-19 (`createCreditPurchase`), porque `CreditCreate` no trae capital ni plazo.
3. El OCR local (Tesseract) no tiene interruptor en C1 (`ocr_engine` es solo lectura); lo desactivable es la IA que ayuda a leer tickets (decisión H), que vive en P-26.11.
4. Las porciones base (2) se muestran en solo lectura: C1 no tiene operación para cambiarlas en Ajustes; cada receta cambia las suyas en P-12.
5. La regla «al menos un peso mayor que 0» del comparador es técnica (la normalización de K17 divide entre la suma), no un mínimo de seguridad.
6. U1 N-288 llama «opcional» a la fecha de deuda en cero, pero C1 la lista en `missing_constraints` y K16 dice que sin restricciones el comparador no recomienda: se sigue a C1 y K16 (solicitud a U1 en el reporte del agente).
7. Las cifras de ejemplo vienen del oráculo (R-61 a R-66, papel de $199.00, cotización $800.50 sobre $3,094.50) o de decisiones (K22, K40); las demás son marcadores entre llaves.
8. Vista previa del gasto fijo (P-26.02, vista previa 3): la regla de fechas sale de los comentarios de `planned_item` en el esquema; con «cada cuánto» mayor que 1 en quincenal se lee como «uno de cada {cada cuánto} días de pago, desde el primero que no es anterior a la primera fecha». El reparto por fecha aplica RN-03 a cada gasto que generará la partida. Ninguna de las dos lecturas toca cifras canon (el oráculo no trae gastos fijos), y la vista previa usa la misma regla de fechas que el motor de S2; si S2 la precisa distinto, manda S2.
9. El préstamo a plazo fijo no captura «¿la tasa incluye IVA?» (como en la 1.4). Qué hace S3 con un préstamo con `vat_included` vacío es PROVISIONAL, en uso la opción A (no suma IVA), bloqueo `docs/estado/bloqueos/S3-iva-lectura.md` (pregunta 1, sin recomendación); si Axel elige B, P-26.05 suma ese campo al préstamo y «Falta capturar» lo nombra.
10. Los créditos que ya existían al migrar a 1.4 sin porcentaje de IVA se ven con «Sí, ya incluye IVA»: la 00004 los pasa a `vat_included` = 1 (supuesto 3 de `docs/estado/p2c/S1-puerta2.md`, que S1 pide confirmar a Axel). La pantalla solo muestra lo guardado.

# U4 · Ajustes estilo Notion — parte b (secciones 8 a 14)

Versión 1.5 (11-oct-2026) · Flujo P2c, agente `correccion-U4-r1` (hallazgos M-1, M-2, M-3 y M-5 de `docs/estado/p2c/REVISION-U4-r1.md`): P-26.08 lee y escribe día y hora con `getReportSettings` y `updateReportSettings` de C1 v1.3 y cita el bloqueo abierto `C1-hora-reportes-job`; P-26.12 escribe `scope` (`onward` o `all`) de C1 v1.3, corrige lo que hace C1 v1.2, marca PROVISIONAL la preselección del diálogo (bloqueo nuevo `U4-reparto-preseleccion`) y cita `C1-reparto-alcance`; P-26.13 aplica Q2-16 A, Q2-17 C y Q2-18 A (`S2-metas` RESUELTO) y cita el bloqueo abierto `S2-metas-vencida`. Versión 1.4 (11-oct-2026) · Flujo P2c, agente `U4-puerta2` (decisiones de la puerta 2 en `docs/estado/puerta2/RESPUESTAS.md`: Q2-11 B (K47), P-26.08 y P-26.09 leen y escriben las claves `reports.*` y `alerts.*` de `app_setting` de la migración 00004, con una sola variante; Q2-72 y Q2-73 (K50) dan sus valores iniciales; Q2-66 B (en un mes corto sale el último día); Q2-68 C (P-26.12 pregunta el alcance); Q2-70 (a) (P-26.14: días restantes y consumo por semana); Q2-09 A (R5-1: P-26.13 y P-26.14 siguen la regla general de la parte a §0.2); Q2-76 (K50) `ocr.timeout_seconds` en P-26.11; Q2-82 (K51) interruptor «Atajos de una tecla» en P-26.10; R4-3 (citas en la misma línea)). Versión 1.3 (11-oct-2026) · Flujo WF-1, agente `correccion-U4-r4` (hallazgo R4-1 de `docs/estado/wf1/REVISION-U4-r4.md`: en P-26.13, el interruptor «Activa» de una meta no guarda al tocar: muestra antes la vista previa (3) del total reservado para metas, como la excepción de la parte a §0.2, y `Ctrl+/` la muestra sumada sobre las metas seleccionadas). Versión 1.2 (11-oct-2026) · Flujo WF-1, agente `correccion-U4-r2` (hallazgos R2-1 punto 1 y R2-3 de `docs/estado/wf1/REVISION-U4-r2.md`: en P-26.13, vista previa del IMC y de la meta sugerida con cualquier dato del perfil, y metas que citan `S2-metas` para los meses que faltan, el formato del avance y los estados). Versión 1.1 (11-oct-2026) · Flujo WF-1, agente `correccion-U4-r1` (hallazgos H-1, H-2 y H-3 de `docs/estado/wf1/REVISION-U4-r1.md`: lo común pasa a la parte a §0, vista previa en P-26.13 y P-26.14, una sola regla de reparto en P-26.12). Versión 1.0 del mismo día, agente `U4b-ajustes`. El integrador une esta parte con `U4-ajustes-a.md` (secciones 1 a 7) en `docs/design/U4-ajustes.md`. Fuentes: `docs/plan/decisiones.md` §A–§K, `docs/criterios/U4.md` 1.2, `docs/design/U1-arquitectura.md` 1.2 (§2.3, §2.4, §2.5, §2.6, §3.28, §5), extractos de C1 v1.2 (`docs/contratos/extractos/ajustes-a.yaml` a `ajustes-e.yaml`, `reportes.yaml`, `salud.yaml`, `presupuesto-a.yaml`, `presupuesto-b.yaml`, `compras-tiendas-a.yaml`, `compras-tiendas-b.yaml`, `inicio.yaml`), esquema 1.4 (`app_setting` e inventario `@param` de 00003 y 00004) y `docs/design/DESIGN.md` §2 y §4.6 (la §7 de componentes aún no existe; aquí solo se citan nombres de token ya publicados).

## Convenciones de la parte b (formato; lo común rige desde la parte a §0)

- **Encabezado de sección**: `## P-26.xx Nombre`. Conteo: `grep -c '^## P-26\.' docs/design/U4-ajustes-b.md` = 7 (14 al unir con la parte a).
- **Tabla de campos** de cada sección: Campo · Tipo · Validación y mensaje · Valor por defecto · Lectura · Escritura · Necesidad (`N-xxx` de U1 §3.28). «Valor por defecto» lleva solo lo decidido en §A–§K (con su cita) o, si no hay decisión, «Sin valor: lo captura Axel»; lo que no se edita dice «Calculado» o «Fijo». Las operaciones son `operationId` de C1 v1.2; «—» en Escritura = no se edita aquí.
- **Parámetros operativos (K44, U4-12)**: tabla Parámetro · Clave · Valor actual e inicial · Validación y mensaje · Restablecer · Decisión · Lectura y escritura. Reglas de U1 §2.3: el valor vigente y «Valor inicial: …» van en la misma fila; «Restablecer al valor inicial» aparece solo cuando `is_initial` es falso y manda `reset_parameters: [clave]`; un valor vacío o de otro tipo no se guarda (422 `validation_failed`) y el campo vuelve a su valor anterior; con los rangos aprobados en la puerta 2 (parte a §0.5; K50): día de la semana de 1 a 7, día del mes de 1 a 31 y horas «HH:MM» de 00:00 a 23:59; los demás, entero ≥ 0.
- **Lo común a las 14 secciones** está en un solo bloque, `U4-ajustes-a.md` §0, y esta parte no lo repite ni lo cambia: edición en el lugar, guardado y vista previa en línea (§0.2), estados vacío, carga, guardando, guardado, restablecido, error y sin conexión con un solo token y un solo texto cada uno (§0.3), diálogo de acción crítica (§0.4), parámetros operativos (§0.5), catálogo de mensajes V-xx (§0.6) y táctil, teclado y foco visible (§0.7). Lo reversible se guarda sin preguntar; las acciones de U4-08 de esta parte (restaurar un respaldo, activar una IA, descartar un chat detectado, que borra) usan el diálogo de §0.4.
- **Atajos**: solo los de U1 §5.3, sin atajos nuevos (0 duplicados): `↑`/`↓` y `Enter` (filas), `Esc` (cancelar edición o cerrar hoja), `Ctrl+Enter` (guardar hoja o formulario), `Alt+N` (nuevo registro en la base de la página), `Ctrl+/` (editar a la vez filas seleccionadas), `Ctrl+K` (buscar una sección: «Notificaciones», «Respaldos»…), `Ctrl+Shift+L` (tema), `?` (lista). Menú «/» de U1 §5.2: `/reporte` (P-26.08), `/respaldo` (P-26.11), `/tema` (P-26.10) e `/ir`.
- **Puntos de quiebre**: móvil 412 px = lista de secciones y cada sección como página; las tablas se ven como lista de tarjetas con sus propiedades apiladas y cada fila abre una hoja inferior con botón primario de `size.control.primary-sheet`. Escritorio 1440 px = columna interna con las 14 secciones y el contenido a la derecha (U1 §2.3).
- **Formatos** (decisión B): dinero `$1,234.56` con cifras tabulares; fechas `10-oct-2026` (con día de la semana donde importa: `Dom 18 oct 2026`); en el envío, ISO. Las cifras de ejemplo vienen del oráculo, de los criterios o son marcadores (`$T`, `n`, `hh:mm`).

## P-26.08 Reportes y frecuencias

**Propósito.** Programar los tres reportes automáticos de la decisión F (resumen diario en la noche, resumen semanal el domingo, reporte mensual en PDF el día 1), elegir quién los recibe y lanzar a mano lo que va bajo demanda: el reporte mensual ahora, el estado de cuenta de Karen y la vista web mensual (N-296, N-297).

**Vista de base de datos.** Tabla (única vista, U1 §2.3): una fila por reporte programado (`daily_summary`, `weekly_summary`, `monthly_report` de `listNotifications`), con su tarea programada de `listJobs`. Sin tablero, calendario ni galería: son tres filas fijas (no se crean ni se borran reportes). Debajo, el bloque «Bajo demanda».

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Reporte | Texto, solo lectura | — | Resumen diario · Resumen semanal · Reporte mensual en PDF (fijo, sembrado en 00001) | `listNotifications` (`name`) | — | N-296 |
| Activo | Interruptor | — | Encendido (decisión F: son los reportes por defecto; `is_enabled` = 1 al sembrar) | `listNotifications` (`is_enabled`) | `updateNotification` (`is_enabled`) | N-297 |
| Frecuencia | Texto, solo lectura | — | Diario · semanal · mensual (fijo por tipo, decisión F) | `listNotifications` (`code`) | — | N-296 |
| Día de la semana (solo semanal) | Selector de lunes a domingo (1 a 7, ISO 8601) | Entero de 1 a 7 · «Elige un día de lunes a domingo» | Domingo (7), decisión F; clave `reports.weekly_weekday` (Q2-11 B, K47) | `getReportSettings` (`reports.weekly_weekday`) | `updateReportSettings` (ver «Escritura de día y hora») | N-297 |
| Día del mes (solo mensual) | Número entero | Entero de 1 a 31 · «Escribe un día del mes, del 1 al 31»; con 29, 30 o 31, en los meses que no tienen ese día sale el último día del mes (Q2-66 B; bloqueo `docs/estado/bloqueos/U4-dia-mes-reporte.md` RESUELTO) | 1, decisión F; clave `reports.monthly_day_of_month` (Q2-11 B, K47) | `getReportSettings` (`reports.monthly_day_of_month`) | `updateReportSettings` (ver «Escritura de día y hora») | N-297 |
| Hora | Hora `HH:MM` de 24 h, hora del centro de México (America/Mexico_City) | De 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Diario 21:30, semanal 19:00 y mensual 08:00 (Q2-72, K50); claves `reports.daily_send_time`, `reports.weekly_send_time` y `reports.monthly_send_time` (Q2-11 B, K47). Debajo del campo, la pista de la decisión F: «Noche», «Domingo» o «Día 1» | `getReportSettings` (la clave de la fila) | `updateReportSettings` (ver «Escritura de día y hora») | N-297 |
| Destinatarios | Selección múltiple: personas con chat confirmado (Axel, Ximena) | Solo chats confirmados en P-26.09; Karen no aparece (sin cuenta, decisión A) · «Primero confirma el chat de esta persona en Notificaciones» | Sin valor: lo captura Axel (decisión A: destinatarios configurables por aviso) | `listNotifications` (`recipient_ids`) y `listNotificationRecipients` | `setNotificationRecipients` | N-297 |
| Próximo envío | Fecha y hora, calculada | — | Calculado por la tarea programada | `listJobs` (`next_run_at`) | — | N-296 |
| Último resultado | Estado y fecha, calculado | — | Calculado («Correcto», «Con error: …», «Omitido») | `listJobs` (`last_status`, `last_run_at`, `last_error`) | — | N-296 |

Reglas fijas que se leen en la sección (no son campos): canal Telegram en texto plano, sin formato (decisión F); el reporte mensual se envía por Telegram **y** se guarda en el Centro de reportes (decisión F).

**Escritura de día y hora (Q2-11 B, K47; bloqueo `docs/estado/bloqueos/S1-parametros-avisos.md` RESUELTO).** Día y hora de cada reporte son claves de `app_setting` del hogar (migración 00004; ya no están en `notification.params_json`, que queda `{}` en los tres reportes). Se leen con `getReportSettings` (`parameters`: `value`, `initial_value`, `is_initial`) y se cambian o restablecen por clave con `updateReportSettings` (`parameters` y `reset_parameters`; parte a §0.5), la sección «Reportes y frecuencias» que C1 v1.3 agregó para las 5 claves `reports.*` (GET `/api/v1/settings/reports`) (PUT `/api/v1/settings/reports`) (N-296 y N-297). Las claves `alerts.*` siguen en `getAlertSettings` y `updateAlertSettings` (P-26.09). `updateNotification` solo cambia «Activo». La tarea programada de cada reporte (`listJobs`) tiene además su `schedule`, que C1 deja cambiar con `updateJob`. Qué manda, las claves o el `schedule` de la tarea, y por lo tanto si «Próximo envío» (`next_run_at`) se recalcula al guardar el día o la hora, es **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/C1-hora-reportes-job.md`; en uso, `updateJob` acepta `schedule` en todas las tareas). Esta sección no usa `updateJob`.

**Parámetros operativos de la sección (K44).**

| Parámetro | Clave | Valor actual e inicial | Validación y mensaje | Restablecer | Decisión | Lectura y escritura |
|---|---|---|---|---|---|---|
| Hora del resumen diario | `reports.daily_send_time` | Actual: `value` · Inicial: 21:30 | `HH:MM` de 24 h, de 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Vuelve a 21:30 | Q2-72 (K50); lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |
| Día de la semana del resumen semanal | `reports.weekly_weekday` | Actual: `value` · Inicial: 7 (domingo) | Entero de 1 a 7 · «Elige un día de lunes a domingo» | Vuelve a 7 | Decisión F; lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |
| Hora del resumen semanal | `reports.weekly_send_time` | Actual: `value` · Inicial: 19:00 | `HH:MM` de 24 h, de 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Vuelve a 19:00 | Q2-72 (K50); lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |
| Día del mes del reporte mensual | `reports.monthly_day_of_month` | Actual: `value` · Inicial: 1 | Entero de 1 a 31; en un mes más corto sale el último día (Q2-66 B) · «Escribe un día del mes, del 1 al 31» | Vuelve a 1 | Decisión F, Q2-66 B; lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |
| Hora del reporte mensual | `reports.monthly_send_time` | Actual: `value` · Inicial: 08:00 | `HH:MM` de 24 h, de 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Vuelve a 08:00 | Q2-72 (K50); lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |

**Bajo demanda.**
- «Generar reporte mensual ahora» (acción primaria; también `/reporte`): `generateMonthlyReport` con el mes (selector con los meses de `listMonthlyReportMonths`) y la casilla «Enviarlo también por Telegram a los destinatarios del mensual» (`send_telegram`, sin preselección: la persona decide). Respuesta aceptada: «Reporte en preparación · Ver en el Centro de reportes» (P-23).
- «Enviar ahora a…» en cada fila: `sendTestNotification` con la persona (Axel o Ximena); sirve de disparo manual del resumen.
- «Generar estado de cuenta de Karen» (decisión A: bajo demanda; Axel lo reenvía a mano): `generateExternalStatement` con Karen fija y el periodo (inicio y fin; «La fecha final no puede quedar antes de la inicial»). Al terminar abre P-25.
- «Abrir la vista web mensual»: enlace a P-24 con el último mes de `listMonthlyReportMonths`. Siempre disponible (decisión F); no se configura.

**Vista previa del efecto.** Al cambiar día u hora, debajo de la fila: «Próximo envío: Dom dd mmm aaaa, hh:mm», calculado en el teléfono o navegador con el día y la hora tecleados antes de guardar; al guardar, lo reemplaza `next_run_at` de `listJobs` (**PROVISIONAL**, bloqueo abierto `C1-hora-reportes-job`: ver «Escritura de día y hora»). Con día del mes 29, 30 o 31: «En los meses que no tienen el día {d}, sale el último día del mes (por ejemplo, el 28-feb-2027)» (Q2-66 B). Con el interruptor apagado: «No se enviará mientras esté apagado». Sin destinatarios: «Sin destinatarios: no llegará a nadie» (insignia del rol advertencia de DESIGN §2.9 con icono y texto).

**Estados.** Vacío: no aplica (las tres filas vienen sembradas); si `listNotifications` no trae alguna, «Falta el reporte en la base · Avisar a Axel» (error de instalación). Carga: esqueleto de tres filas. Error al cargar: «No se pudo cargar la configuración de reportes · Reintentar». Error al guardar: el de la parte a §0.3 («No se pudo guardar. Revisa la conexión con el servidor de casa · Reintentar»; el campo vuelve al valor anterior). Guardado: automático, «Guardado».

**Edición en el lugar.** Activo, día, hora y destinatarios se editan en la celda (o en la hoja en móvil). No se editan: nombre, frecuencia, próximo envío y último resultado (calculados o fijos).

**Atajos.** `↑`/`↓` entre filas, `Enter` abre la celda o la hoja, `Esc` cancela, `Ctrl+Enter` guarda la hoja; `/reporte` genera el mensual.

**Móvil y escritorio.** 412 px: tres tarjetas (nombre, activo, «Próximo envío») que abren una hoja con día, hora y destinatarios; el bloque «Bajo demanda» como botones de ancho completo. 1440 px: tabla con columnas Reporte · Activo · Día · Hora · Destinatarios · Próximo envío · Último resultado, y el bloque «Bajo demanda» debajo.

## P-26.09 Notificaciones

**Propósito.** Configurar las 4 alertas de la decisión F (pago próximo, presupuesto al 80 % y al 100 %, precios viejos antes del día de despensa, plan de comidas sin armar), quién recibe cada una, sus umbrales y anticipación (K44), y registrar los chats de Telegram de Axel y Ximena en el bot MultiX: el bot detecta el chat cuando la persona le escribe y ella lo confirma aquí o lo descarta (K4). Necesidades N-298, N-299, N-300, N-316, N-323 y N-324.

**Vista de base de datos.** Tres bloques en tabla (U1 §2.3: tabla): «Alertas» (4 filas fijas de `listNotifications`: `payment_due`, `budget_threshold`, `stale_prices`, `meal_plan_missing`), «Chats confirmados» (una fila por persona con cuenta) y «Chats detectados por confirmar» (lista que se vacía al confirmar o descartar). Los avisos ya enviados no viven aquí: están en P-23 (tabla, tablero por tipo y calendario por fecha de envío, U1 §2.5).

**Alertas.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Alerta | Texto, solo lectura | — | Pago próximo · Presupuesto al 80 % y al 100 % · Precios viejos antes del día de despensa · Plan de comidas sin armar (fijo, 00001) | `listNotifications` (`name`) | — | N-298 |
| Activa | Interruptor | — | Encendida (decisión F; `is_enabled` = 1 al sembrar) | `listNotifications` (`is_enabled`) | `updateNotification` (`is_enabled`) | N-300 |
| Destinatarios | Selección múltiple: personas con chat confirmado | Solo chats confirmados · «Primero confirma el chat de esta persona» | Sin valor: lo captura Axel (decisión A) | `listNotifications` (`recipient_ids`) y `listNotificationRecipients` | `setNotificationRecipients` | N-300 |
| Días antes del pago (pago próximo) | Lista de días en fichas: 1 = el día anterior, 0 = el mismo día | Uno o más enteros ≥ 0 · «Escribe al menos un día; 0 es el mismo día del pago» | [1, 0] (decisión F, K44) | `getAlertSettings` (`alerts.payment_due_days_before`) | `updateAlertSettings` | N-323 |
| Primer aviso del presupuesto | Porcentaje entero | Entero ≥ 0 · «Escribe un porcentaje entero, sin decimales» | 80 % (decisión F, K44) | `getAlertSettings` (`alerts.budget_warning_pct`) | `updateAlertSettings` | N-323 |
| Aviso de tope del presupuesto | Porcentaje entero | Entero ≥ 0 · «Escribe un porcentaje entero, sin decimales» | 100 % (decisión F, K44) | `getAlertSettings` (`alerts.budget_limit_pct`) | `updateAlertSettings` | N-323 |
| Días antes del día de despensa (precios viejos) | Número entero de días | Entero ≥ 0 · «Escribe cuántos días antes; 0 es el mismo día» | 1 día (Q2-73, K50); clave `alerts.stale_prices_advance_days` (Q2-11 B, K47) | `getAlertSettings` (`alerts.stale_prices_advance_days`) | `updateAlertSettings` | N-300 |
| Día de revisión (plan sin armar) | Selector de lunes a domingo (1 a 7) | Entero de 1 a 7 · «Elige un día de lunes a domingo» | Domingo (7), decisión F; clave `alerts.meal_plan_missing_weekday` (Q2-11 B, K47) | `getAlertSettings` (`alerts.meal_plan_missing_weekday`) | `updateAlertSettings` | N-300 |
| Hora (avisos con hora) | `HH:MM` de 24 h, hora del centro | De 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | 09:00 (Q2-73, K50); clave `alerts.send_time` (Q2-11 B, K47) | `getAlertSettings` (`alerts.send_time`) | `updateAlertSettings` | N-300 |

Reglas fijas que se leen en la fila (no son campos): el aviso de presupuesto mira **cada categoría y el total** (decisión F); «casillas sin armar» = pares día × tiempo obligatorio (desayuno, comida, cena) en que alguna persona no tiene receta ni «fuera», máximo 21 por semana (K43: es una regla, no un parámetro). A qué alertas aplica la hora lo fija S7 (ver SUPUESTOS).

**Chats confirmados** (registro de chat, N-299; nunca se muestran el `chat_id` ni el token del bot, K30).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Persona | Texto, solo lectura | — | Axel y Ximena (las dos personas con cuenta; Karen no recibe avisos, decisión A) | `listNotificationRecipients` (`person_id`) y `listPersons` | — | N-299 |
| Chat registrado | Sí o no, con fecha de comprobación | — | Calculado («No» hasta confirmar) | `listNotificationRecipients` (`is_registered`, `verified_at`) | — | N-299 |
| Etiqueta | Texto opcional | Texto libre; vacío borra la etiqueta | Sin valor: lo captura Axel | `listNotificationRecipients` (`label`) | `updateNotificationRecipient` (`label`) | N-299 |
| Recibe avisos | Interruptor | — | El que deja la confirmación del chat (`is_active` del destinatario creado) | `listNotificationRecipients` (`is_active`) | `updateNotificationRecipient` (`is_active`) | N-299 |

**Chats detectados por confirmar** (K4).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Nombre en Telegram | Texto, solo lectura | — | El que da Telegram | `listTelegramChatCandidates` (`telegram_name`) | — | N-299 |
| Detectado | Fecha y hora, solo lectura | — | Cuando el bot recibió el mensaje | `listTelegramChatCandidates` (`detected_at`) | — | N-299 |
| Terminación del chat | Texto «•••» con los últimos 4 caracteres | — | Calculado, para distinguir dos chats con el mismo nombre | `listTelegramChatCandidates` (`chat_id_masked`) | — | N-299 |
| Confirmar como | Selector de persona (Axel o Ximena) + botón «Confirmar» | Un chat ya confirmado responde 409 · «Este chat ya estaba confirmado; recarga la lista» | La persona del usuario de la sesión (K4, `TelegramChatCandidateConfirm`) | — | `confirmTelegramChatCandidate` (con `Idempotency-Key`) | N-316 |
| Descartar | Botón (acción crítica) | — | — | — | `dismissTelegramChatCandidate` | N-324 |

Texto de ayuda sobre la lista: «Para recibir avisos, escribe un mensaje al bot desde tu Telegram y confírmalo aquí». El texto exacto lo fija S7 (WF-3) según cómo detecte el chat (K4); es solo redacción, no datos.

**Parámetros operativos de la sección (K44).**

| Parámetro | Clave | Valor actual e inicial | Validación y mensaje | Restablecer | Decisión | Lectura y escritura |
|---|---|---|---|---|---|---|
| Primer aviso del presupuesto | `alerts.budget_warning_pct` | Actual: `value` · Inicial: 80 % | Entero ≥ 0 · «Escribe un porcentaje entero, sin decimales» | Vuelve a 80 % | Decisión F, K44 | `getAlertSettings` · `updateAlertSettings` |
| Aviso de tope del presupuesto | `alerts.budget_limit_pct` | Actual: `value` · Inicial: 100 % | Entero ≥ 0 · «Escribe un porcentaje entero, sin decimales» | Vuelve a 100 % | Decisión F, K44 | `getAlertSettings` · `updateAlertSettings` |
| Días antes del pago próximo | `alerts.payment_due_days_before` | Actual: `value` · Inicial: [1, 0] | Lista de uno o más enteros ≥ 0 · «Escribe al menos un día; 0 es el mismo día del pago» | Vuelve a [1, 0] | Decisión F, K44 | `getAlertSettings` · `updateAlertSettings` |
| Días antes del día de despensa (precios viejos) | `alerts.stale_prices_advance_days` | Actual: `value` · Inicial: 1 día | Entero ≥ 0 · «Escribe cuántos días antes; 0 es el mismo día» | Vuelve a 1 | Q2-73 (K50), K43; lugar Q2-11 B (K47) | `getAlertSettings` · `updateAlertSettings` |
| Hora de los avisos con hora | `alerts.send_time` | Actual: `value` · Inicial: 09:00 | `HH:MM` de 24 h, de 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Vuelve a 09:00 | Q2-73 (K50), K43; lugar Q2-11 B (K47) | `getAlertSettings` · `updateAlertSettings` |
| Día de revisión del plan sin armar | `alerts.meal_plan_missing_weekday` | Actual: `value` · Inicial: 7 (domingo) | Entero de 1 a 7 · «Elige un día de lunes a domingo» | Vuelve a 7 | Decisión F; lugar Q2-11 B (K47) | `getAlertSettings` · `updateAlertSettings` |

**Acciones.** «Enviar aviso de prueba» (acción primaria): elige alerta y persona y llama `sendTestNotification`; respuesta: «Aviso de prueba enviado a Ximena» o el error en palabras simples («El bot no pudo entregar el mensaje · Revisa que el chat esté confirmado»). «Descartar» un chat detectado abre la acción crítica: «Descartar el chat de “nombre en Telegram” (•••nnnn): se borra de esta lista y el bot no le mandará avisos. Si era de alguien de la casa, tendrá que volver a escribirle al bot. · Descartar chat». Los chats ya confirmados no se tocan (N-324).

**Vista previa del efecto.**
- Umbrales: «Con el presupuesto total de este mes ($T), el primer aviso llega al gastar $A y el de tope al gastar $B», con $T de `getBudgetMonth` y $A = $T × porcentaje ÷ 100 en centavos con ROUND_HALF_UP (decisión B); igual por categoría al tocar una categoría de ejemplo.
- Días antes del pago: «El próximo pago (Vie dd mmm aaaa) avisaría el Jue dd mmm y el Vie dd mmm», con la fecha de `getDashboardNextPayment` y un renglón por día de la lista.
- Precios viejos y plan sin armar: «Llegará n días antes del día de despensa» y «Se revisa cada domingo» (o el día elegido); la hora, en los avisos donde S7 la aplique, sale de `alerts.send_time` (09:00 al inicio).

**Estados.** Vacío: «Chats detectados por confirmar» sin filas = «No hay chats nuevos. Escribe al bot desde Telegram y vuelve aquí» (con «Actualizar»); «Chats confirmados» sin chat = «Ximena aún no confirma su chat» en la fila. Carga: esqueleto de las tres tablas. Error: «No se pudieron cargar los avisos · Reintentar»; al guardar, el campo vuelve a su valor. Guardado: automático, «Guardado»; confirmar un chat lo pasa a «Chats confirmados» con «Chat confirmado».

**Edición en el lugar.** Activa, destinatarios, días, porcentajes, hora, etiqueta y «Recibe avisos». No se editan: nombre de la alerta, chat registrado, datos del chat detectado (vienen del bot).

**Atajos.** `↑`/`↓` y `Enter` en las tres tablas; `Ctrl+/` activa o desactiva a la vez las alertas seleccionadas; `Esc` cancela.

**Móvil y escritorio.** 412 px: tarjetas por alerta (nombre, activa y resumen «1 día antes y el mismo día») que abren su hoja; los chats detectados arriba cuando los hay, con «Confirmar» de ancho completo. 1440 px: las tres tablas una debajo de otra; parámetros en la fila de su alerta.

## P-26.10 Apariencia

**Propósito.** Elegir el tema de este dispositivo: según el sistema (por defecto), claro u oscuro (decisión G; N-301), y encender o apagar los atajos de una tecla (Q2-82 A, K51; WCAG 2.1.4, DESIGN §9.6). Dos controles, los dos del dispositivo: el tema (U1 §2.3) y el interruptor que ratificó la puerta 2 (U1 aún no lo lista: solicitud en `docs/estado/p2c/U4-puerta2.md`).

**Vista de base de datos.** Ninguna (un solo control).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Tema | Control segmentado: Sistema · Claro · Oscuro | Uno de los tres | Sistema (decisión G) | Local del dispositivo (U1 §0.3: no usa C1) | Local del dispositivo (también `Ctrl+Shift+L` y `/tema`) | N-301 |
| Atajos de una tecla | Interruptor | — | Encendido (Q2-82 A, K51; DESIGN §9.6) | Local del dispositivo (no usa C1) | Local del dispositivo | N-301 |

Valores fijos del sistema de diseño que se leen en la sección (no son campos; DESIGN §1.2 P4 y P6, §2.1): «Saturación: apagada» y «Densidad: cómoda» (decisión G). No se configuran en esta versión (ver SUPUESTOS).

**Vista previa del efecto.** Atajos de una tecla: al apagarlo, `/`, `?` y las flechas dejan de actuar como atajos fuera de los campos de texto (DESIGN §9, criterio 2.1.4; U1 §5.4); los atajos con `Ctrl` o `Alt` siguen. Tema: el cambio es su propia vista previa: la página cambia de tema al instante (sin transición si el sistema pide movimiento reducido, DESIGN §6.3). **Estados.** Sin vacío, carga ni error (no hay red); guardado inmediato en el dispositivo, «Guardado». Sin conexión funciona igual. **Edición en el lugar:** el control mismo. **Atajos:** `Ctrl+Shift+L`, `/tema`. **Móvil y escritorio:** 412 px, tres botones de 44 px a todo el ancho; 1440 px, el control a la izquierda del contenido, con una muestra de tarjeta en claro y oscuro.

## P-26.11 Datos y respaldos

**Propósito.** Ver y operar los respaldos locales (copia diaria versionada, `integrity_check`, restauración probada; decisión H), importar `gastos_fijos_ingresos.csv`, revisar los datos semilla, encender o apagar la IA opcional (desactivada por defecto, decisión H), ver el estado del servicio y ajustar los parámetros de datos (K44). Necesidades N-302 a N-308 y N-325.

**Vista de base de datos.** Tabla de respaldos (vista por defecto, U1 §2.5; filtro por periodo, orden por fecha descendente). Los demás bloques son fichas: Importar, Datos semilla (3 filas), IA opcional (2 filas), Estado del servicio y Parámetros.

**Respaldos.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Fecha y hora | Fecha y hora, solo lectura | — | Calculado | `listBackups` (`started_at`, `finished_at`) | — | N-302 |
| Tipo | Etiqueta: Diario · Manual · Previo a una migración | — | Calculado | `listBackups` (`kind`) | — | N-302 |
| Resultado | Etiqueta: En curso · Correcto · Falló (con icono) | — | Calculado | `listBackups` (`status`, `error`) | — | N-302 |
| Tamaño | Tamaño en MB o GB | — | Calculado | `listBackups` (`size_bytes`) | — | N-302 |
| Integridad | Texto: «Íntegro» o el error de `integrity_check` | — | Calculado | `listBackups` (`integrity_check`) | — | N-302 |
| Restauración probada | Fecha y resultado | — | Calculado | `listBackups` (`restore_tested_at`, `restore_test_result`) | — | N-302 |
| Ubicación y huella | Ruta en la carpeta o disco de respaldos y huella abreviada | — | Calculado | `listBackups` (`file_path`, `sha256`) | — | N-302 |
| Próximo respaldo diario | Fecha y hora, calculada | — | Calculado | `listJobs` (`next_run_at` de la tarea del respaldo) | — | N-302 |
| Hora del respaldo diario | `HH:MM` de 24 h | De 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Sin valor inicial: lo proponen S6 y WF-4 (K31 punto 3, puerta 4), PROVISIONAL | `listJobs` (`schedule`) | `updateJob` (`schedule` de la tarea del respaldo), PROVISIONAL (ver SOLICITUDES del reporte) | N-302 |

Encabezado de la tabla: «Último respaldo correcto: dd-mmm-aaaa hh:mm · Próximo: … · Última restauración probada: …».

**Acciones de respaldo.** «Respaldar ahora» (acción primaria; también `/respaldo`): `createBackup`; al aceptar, la fila nueva aparece «En curso» y cambia sola a «Correcto» o «Falló». «Restaurar» en cada fila correcta: `restoreBackup`, con **doble confirmación de acción crítica** (U4-08, N-303): paso 1, vista previa del efecto; paso 2, escribir `RESTAURAR` (`BackupRestoreRequest.confirmation`). Texto: «Restaurar el respaldo del dd-mmm-aaaa hh:mm: la base vuelve a como estaba en ese momento y se pierden los cambios hechos después en los dos teléfonos y en la computadora. Antes se toma un respaldo de hoy por si necesitas regresar. · Restaurar respaldo». Al terminar: «Respaldo restaurado · Íntegro · Respaldo previo guardado» (`RestoreResult`: `integrity_check`, `pre_restore_backup_id`) y todas las sesiones recargan.

**Importar gastos fijos e ingresos.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Archivo | Archivo CSV (plantilla de `plan.md` §8.2) | Solo `.csv` con las columnas de la plantilla · «El archivo no tiene las columnas de la plantilla: …» (renglón y columna en palabras) | Sin valor: lo elige Axel | — | `previewPlannedItemsImport` (vista previa) | N-304 |
| Lo que se creará | Tabla de vista previa, solo lectura | Renglones con error marcados con icono y texto | Calculado | `previewPlannedItemsImport` | — | N-304 |
| Importar | Botón primario | Se activa solo sin errores · «Corrige los renglones marcados antes de importar» | — | — | `importPlannedItems` | N-304 |

Resultado: «n partidas creadas» (`ImportResult.created_count`) con enlace a Presupuesto (P-26.02).

**Datos semilla.** Tres filas fijas de `getSeedStatus` (N-305): compra del 10-oct-2026, crédito Kueski con sus despensas y pagos, CSV de gastos fijos e ingresos; cada una con «Cargada» o «Pendiente» y su número de registros. Solo lectura.

**IA opcional** (decisión H; la IA solo se usa para OCR de tickets y sugerencias de recetas; el núcleo es determinista).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Proveedor | Texto: «Modelo local (Ollama)» · «API de Claude» | — | Fijo (2 filas) | `listAiProviders` (`provider`) | — | N-306 |
| Activa | Interruptor | Encender pide acción crítica | **Apagada** los dos (decisión H; `ai.local_model_enabled` y `ai.claude_api_enabled` = false en 00001) | `listAiProviders` (`is_enabled`) | `updateAiProvider` (`is_enabled`) | N-307 |
| Usos permitidos | Etiquetas: OCR de tickets · Sugerencias de recetas | — | Fijo (decisión H) | `listAiProviders` (`allowed_uses`) | — | N-306 |
| Disponible | Sí o no (solo si está activa) | — | Calculado | `listAiProviders` (`is_available`) | — | N-306 |
| Llave | «Configurada» o «Sin configurar» (solo API de Claude); la llave nunca se muestra | — | Calculado | `listAiProviders` (`is_configured`) | — | N-306 |

Encender abre la acción crítica: modelo local, «Activar el modelo local: el OCR de tickets y las sugerencias de recetas usarán el modelo instalado en esta máquina. Siempre pedirá confirmación antes de guardar un ticket. · Activar modelo local»; API de Claude, «Activar la API de Claude: las fotos de tickets y los pedidos de sugerencias de recetas saldrán de la red de casa hacia el servicio de Claude. Siempre pedirá confirmación antes de guardar un ticket. · Activar API de Claude». Apagar es inmediato y sin diálogo. El OCR local con Tesseract no es IA: vive en Tiendas y precios › Fuentes de precio (P-26.06, parte a).

**Estado del servicio.** `getSystemStatus` (N-308): versión, en marcha desde y tamaño de la base de datos. Solo lectura.

**Parámetros operativos de la sección (K44).**

| Parámetro | Clave | Valor actual e inicial | Validación y mensaje | Restablecer | Decisión | Lectura y escritura |
|---|---|---|---|---|---|---|
| Tamaño máximo de un adjunto de gasto | `attachment.max_size_mb` | Actual: `value` · Inicial: 10 MB | Entero ≥ 0, en MB · «Escribe un número entero de MB» | Vuelve a 10 MB | K43, K44 | `getDataSettings` · `updateDataSettings` |
| Retención de las claves de idempotencia | `idempotency.retention_hours` | Actual: `value` · Inicial: 24 h | Entero ≥ 0, en horas · «Escribe un número entero de horas» | Vuelve a 24 h | K6, K44 | `getDataSettings` · `updateDataSettings` |
| Tiempo máximo del OCR de un ticket | `ocr.timeout_seconds` | Actual: `value` · Inicial: 30 s | Entero ≥ 0, en segundos · «Escribe un número entero de segundos» | Vuelve a 30 s | Q2-76 (K50); sección: la de su `@param` en la 00004 (PROPUESTA de S1 que la puerta 2 no preguntó, `docs/estado/p2c/S1-puerta2.md`) | `getDataSettings` · `updateDataSettings` (C1 v1.3 la expone; solicitud de S1) |

Parámetros de K31 de esta sección **sin clave ni operación todavía** (no son campos; se piden en SOLICITUDES): cuántas versiones del respaldo diario se guardan (K31 punto 3) y los permisos de los archivos de la base y de los respaldos (K31 punto 8).

**Vista previa del efecto.** Restaurar: paso 1 con fecha, tamaño, integridad del respaldo elegido y la advertencia de lo que se pierde (arriba). Importar: tabla de lo que se creará antes de confirmar. Adjunto: «Con n MB, un ticket o adjunto más pesado se rechaza con “El archivo pesa más de n MB”». Retención: «La API recuerda cada captura repetida n horas para no duplicarla; los gastos y pesos de la cola del teléfono nunca se duplican porque su clave queda en el registro (K9)». IA: el texto del diálogo de activación. OCR: «Tiempo máximo para leer un ticket: {antes} s → {después} s».

**Estados.** Vacío: sin respaldos = «Aún no hay respaldos · Respaldar ahora». Carga: esqueleto de la tabla y de las fichas. Error: «No se pudo cargar la lista de respaldos · Reintentar»; un respaldo que falla queda en la tabla como «Falló» con su error resumido; una restauración fallida dice «No se restauró; la base sigue como estaba» (`integrity_check` distinto de ok). Guardado: automático en parámetros y en el interruptor de IA (tras confirmar).

**Edición en el lugar.** Parámetros, hora del respaldo (PROVISIONAL: sin valor inicial hasta K31 punto 3; `updateJob` según SOLICITUDES del reporte y supuesto 4) e interruptores de IA. No se editan: filas de respaldos, semillas, estado del servicio (calculados).

**Atajos.** `/respaldo`; `↑`/`↓` y `Enter` en la tabla de respaldos (Enter abre el detalle con «Restaurar»); `Esc` cierra el diálogo.

**Móvil y escritorio.** 412 px: «Respaldar ahora» arriba a todo el ancho, respaldos como tarjetas (fecha, resultado, integridad) y las demás fichas apiladas; «Restaurar» dentro del detalle, nunca en la tarjeta. 1440 px: tabla completa y fichas en dos columnas.

## P-26.12 Personas y reparto

**Propósito.** Ver quién participa en el dinero del hogar y cómo: Axel (titular, con cuenta), Ximena (co-titular, mismos permisos) y Karen (participante externa, sin cuenta ni acceso, paga directo a Kueski su parte; decisión A); fijar el reparto por defecto del conjunto (50/50, decisión B) y, si Axel quiere, un tope de pago por quincena por persona (K13). Necesidades N-309 y N-310.

**Vista de base de datos.** Tabla de personas (única vista, U1 §2.5; orden Axel, Ximena, Karen; sin filtros). Debajo, tabla «Repartos reutilizables» (las reglas con `is_template` de `listCostShareRules`: Hogar, Karen, Conjunto 50/50). No se crean ni se borran personas ni reglas aquí (C1 v1.2 no tiene esas operaciones).

**Personas.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Nombre visible | Texto | No vacío · «Escribe un nombre» | Axel · Ximena · Karen (semilla) | `listPersons` (`display_name`) | `updatePerson` (`display_name`) | N-310 |
| Papel | Etiqueta, solo lectura: Titular · Co-titular · Participante externa | — | Fijo (decisión A) | `listPersons` (`participation`) | — | N-309 |
| Cuenta | Sí o no, solo lectura | — | Axel y Ximena sí; Karen nunca (decisión A) | `listPersons` (`has_account`) | — | N-309 |
| Paga su parte directo al acreedor | Sí o no, solo lectura | — | Karen sí (paga a Kueski; su parte es informativa y repone la línea, decisión A y K15) | `listPersons` (`pays_share_directly`) | — | N-309 |
| Qué hace en el sistema | Texto fijo según el papel | — | Titular y co-titular: «Mismos permisos: captura, consulta y ajustes». Karen: «Sin cuenta. Solo aparece en la despensa en común y en sus compras con Kueski; recibe su estado de cuenta en PDF, que Axel reenvía» | `listPersons` (`participation`) | — | N-309 |
| Tope de pago por quincena | Dinero, opcional | Vacío o monto ≥ $0.00 con dos decimales · «Escribe un monto como $1,234.56 o deja el campo vacío» | **Vacío** (sin tope, K13: Axel lo captura cuando quiera); lo usa solo el comparador como objetivo, no como restricción dura (decisión C) | `listPersons` (`fortnight_payment_cap_cents`) | `updatePerson` (`fortnight_payment_cap_cents`; vacío = `null`) | N-310 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listPersons` (`notes`) | `updatePerson` (`notes`) | N-310 |

**Repartos reutilizables.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Nombre | Texto | No vacío · «Escribe un nombre para el reparto» | Hogar · Karen · Conjunto 50/50 (semilla) | `listCostShareRules` (`name`) | `updateCostShareRule` (`name`) | N-310 |
| Parte de cada persona | Por persona: porcentaje (dos decimales) o monto fijo | Cada parte es porcentaje (0.00 % a 100.00 %, dos decimales, V-06) o monto fijo (≥ $0.00, V-03), nunca los dos (00001, `cost_share_rule_part`). Las partes con porcentaje suman a lo más 100 % · V-16: «Las partes con porcentaje suman {suma} %. No pueden pasar de 100 %: lo que falta lo recibe {persona del resto}.» | Conjunto: Karen 50 % y el hogar 50 % (decisión B); Hogar: todo al hogar; Karen: todo a Karen | `listCostShareRules` (`parts`: `pct_bp`, `fixed_cents`) | `updateCostShareRule` (`parts`) | N-310 |
| Quién toma el resto y los centavos | Texto, solo lectura | — | El titular (K14: las partes de personas externas con ROUND_HALF_UP y el titular toma el resto) | `listCostShareRules` (`remainder_person_id`) | — | N-309 |
| Reparto por defecto | Marca, solo lectura | — | «Hogar» (se aplica cuando ni el gasto ni su categoría traen reparto) | `listCostShareRules` (`is_default`) | — | N-309 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listCostShareRules` (`notes`) | `updateCostShareRule` (`notes`) | N-310 |

**Regla única del reparto (RN-03, S1-08, `CostShareRule` de C1).** Se editan solo las partes explícitas (`parts`: `pct_bp` o `fixed_cents` por persona). La persona del resto (`remainder_person_id`, el titular en las tres reglas sembradas) no tiene campo: su parte se muestra al final, calculada y de solo lectura, como «el resto» = 100 % menos las partes con porcentaje (al aplicarse a un gasto, además menos los montos fijos; RN-03: «El resto recibe la suma menos las partes»). Así cada regla suma 100 % por construcción y no hay que cuadrarla a mano; el único error posible es que las partes con porcentaje pasen de 100 % (V-16), y entonces no se guarda. Un monto fijo mayor que el gasto se rechaza al aplicar la regla al gasto (RN-03), no aquí. En pantalla: «Conjunto 50/50: Karen 50.00 % · el resto 50.00 % (Axel, hogar)». Qué reparto usa cada categoría se elige en Categorías (P-26.04, parte a) y cada gasto en P-03 y P-04.

**Vista previa del efecto.** Al cambiar una parte, antes de guardar: «Con un gasto de $M: Karen $K (M × p con ROUND_HALF_UP) y el hogar $H (el resto)», con $M marcador que la persona puede cambiar en la misma vista previa. Al salir del último campo de la regla, antes de guardar, el diálogo de alcance (Q2-68 C; bloqueo `docs/estado/bloqueos/U4-reparto-retroactivo.md` RESUELTO; como el del presupuesto en P-26.02, K2) pregunta «¿A qué gastos aplica el reparto nuevo?» con dos opciones: «Solo a los gastos nuevos» (los gastos ya capturados conservan las partes que se congelaron al guardarlos y los estados de cuenta ya entregados a Karen no cambian) y «También a los gastos ya capturados» («{n} gastos ya capturados con esta regla cambian sus partes, desde el {fecha del más antiguo}; incluye meses anteriores»), con «Confirmar» y «Cancelar»; `Esc` cancela sin cambiar (parte a §0.2). **La preselección es PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/U4-reparto-preseleccion.md`, se pregunta en la puerta 4): ni Q2-68 C ni K46 la fijan y C1 v1.3 no da valor por omisión. En uso, «Solo a los gastos nuevos» preseleccionada; las otras opciones son preseleccionar «También a los gastos ya capturados» o no preseleccionar ninguna, con «Confirmar» inactivo hasta elegir. Escritura: `updateCostShareRule` con `parts` y `scope` (C1 v1.3, `CostShareRuleUpdate.scope`, obligatorio con `parts` y sin valor por omisión): `onward` = «Solo a los gastos nuevos» y `all` = «También a los gastos ya capturados». C1 v1.2 no recibe alcance y ya se comporta como «solo los nuevos»: las partes de los gastos ya guardados están congeladas en `expense_share`; «también los ya capturados» necesita `scope` = `all` de C1 v1.3. Qué gastos abarca `all` (solo los que nombran la regla, o también los que la toman por su categoría o por ser la de por defecto) es **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/C1-reparto-alcance.md`), y con él el «{n}» y la fecha de «También a los gastos ya capturados». C1 v1.3 no da ese conteo ni esa fecha (solicitud SC-U4-P2-1 en `docs/estado/p2c/U4-puerta2.md`): mientras no los dé, «{n}» y la fecha se rotulan «Vista previa pendiente» (parte a §0.2). Tope por quincena: «El comparador usará $X como tope de esta persona al evaluar el objetivo “pago por quincena parejo y bajo tope”; no descarta escenarios» (K17: pico contra el tope cuando hay tope; decisión C: no es restricción dura).

**Estados.** Vacío: no aplica (las 3 personas y las 3 reglas vienen sembradas). Carga: esqueleto de las dos tablas. Error al guardar: el de la parte a §0.3, con las partes de vuelta a su valor. Guardado: automático al salir del campo; las partes se guardan juntas al confirmar el alcance (Q2-68 C).

**Edición en el lugar.** Nombre visible, tope, notas, nombre y partes de las reglas. No se editan: papel, cuenta, paga directo, persona del resto, reparto por defecto (decididos o fijos).

**Atajos.** `↑`/`↓` y `Enter` en las dos tablas; `Esc` cancela la edición de una parte.

**Móvil y escritorio.** 412 px: tres tarjetas de persona (nombre, papel, tope) y tres de reparto con la barra de partes etiquetada (nombre y porcentaje en texto, nunca solo color: K41 y DESIGN §2.6, Karen en su color categórico con su etiqueta). 1440 px: las dos tablas con la vista previa a la derecha de la regla que se edita.

## P-26.13 Salud y metas de ahorro

**Propósito.** Los perfiles de salud de Axel y Ximena (plan de comidas y salud son solo de ellos dos, decisión A): estatura, nacimiento, sexo, actividad, objetivo, IMC y la meta de calorías sugerida que la persona confirma o cambia (decisión D, K23); el ajuste de la meta por objetivo de la casa y de cada persona (K36); las metas de ahorro con fecha y la aportación mensual planificada al fondo (decisión B). Necesidades N-311, N-312, N-313, N-326 y N-327.

**Vista de base de datos.** Perfiles: una ficha por persona (2), sin vista de lista (U1 §2.5). Metas de ahorro: **tabla** por defecto en esta sección; también tablero por estado (`status`, motor S2) con tres columnas, «En curso», «Cumplida» y «Vencida» (Q2-18 A; bloqueo `docs/estado/bloqueos/S2-metas.md`, pregunta 3, RESUELTO; `in_progress`, `achieved` y `overdue` en `docs/contratos/motores/S2-interfaz.md` 1.3 §3.13). C1 v1.3 todavía declara `status` con `on_track`, `behind` y `achieved`, y U1 §2.5 nombra las columnas en ruta, atrasada y cumplida: los dos tienen que seguir a Q2-18 A (solicitudes en `docs/estado/p2c/correccion-U4-r1.md`), calendario por fecha meta y galería con el avance (U1 §2.5; la galería es la vista por defecto de P-09). Filtro por defecto: activas; orden: fecha meta ascendente.

**Perfil de salud (una ficha por persona).**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Estatura | Metros con dos decimales (se guarda en mm) | Mayor que 0 · «Escribe la estatura en metros, como 1.75» | Sin valor: lo captura la persona | `listHealthProfiles` (`height_mm`) | `updateHealthProfile` o `createHealthProfile` (`height_mm`) | N-311, N-313 |
| Fecha de nacimiento | Fecha (`10-oct-2026` en pantalla) | No después de hoy · «La fecha de nacimiento no puede ser futura» | Sin valor: lo captura la persona | `listHealthProfiles` (`birth_date`; edad en `age_years`) | `updateHealthProfile` (`birth_date`) | N-311, N-313 |
| Sexo | Selector: Mujer · Hombre | Uno de los dos (lo pide la fórmula) · «Elige una opción para calcular la meta» | Sin valor: lo captura la persona | `listHealthProfiles` (`sex`) | `updateHealthProfile` (`sex`) | N-311, N-313 |
| Actividad | Selector de 5 niveles con su factor: sedentaria 1.2 · ligera 1.375 · moderada 1.55 · activa 1.725 · muy activa 1.9 | Uno de los cinco | Sin valor: lo captura la persona (K23: 5 niveles de 1.2 a 1.9) | `listHealthProfiles` (`activity_level`) | `updateHealthProfile` (`activity_level`) | N-311, N-313 |
| Objetivo | Selector: Bajar · Mantener · Subir de peso | Uno de los tres o vacío | **Vacío** («todavía no lo elige», K23) | `listHealthProfiles` (`goal`) | `updateHealthProfile` (`goal`) | N-311, N-313 |
| Peso actual | Kg con fecha, solo lectura | — | Calculado del historial (se registra en P-22) | `listHealthProfiles` (`current_weight_g`, `current_weight_date`) | — | N-311 |
| IMC | Número con un decimal y clase OMS (bajo peso, normal, sobrepeso, obesidad) | — | Calculado (K7: cortes 18.5, 25 y 30) | `listHealthProfiles` (`bmi_tenths`, `bmi_class`) | — | N-311 |
| Meta sugerida | Kcal por día con desglose, solo lectura | — | Calculado (RN-17: Mifflin-St Jeor × actividad − ajuste de bajar o + ajuste de subir; un redondeo ROUND_HALF_UP al final) | `listHealthProfiles` (`suggested_kcal_target`) | — | N-311 |
| Meta confirmada | Kcal por día (entero) | Entero mayor que 0 · «Escribe las kilocalorías por día, sin decimales» | **Vacía** hasta que la persona confirme la sugerida o escriba otra (K23) | `listHealthProfiles` (`daily_kcal_target`) | `updateHealthProfile` (`daily_kcal_target`); botón «Usar la sugerida» | N-313 |
| Mi ajuste para bajar | Kcal (entero), opcional | Vacío o entero ≥ 0 · «Escribe un número entero de kcal, 0 o más, o deja vacío para usar el de la casa» | **Vacío = usa el de la casa** y lo dice: «El de la casa: n kcal» (K36) | `listHealthProfiles` (`goal_adjust_lose_kcal`, `household_goal_adjust_lose_kcal`) | `updateHealthProfile` (`goal_adjust_lose_kcal`; «Usar el de la casa» = `null`) | N-327 |
| Mi ajuste para subir | Kcal (entero), opcional | Igual que el anterior | **Vacío = usa el de la casa** (K36) | `listHealthProfiles` (`goal_adjust_gain_kcal`, `household_goal_adjust_gain_kcal`) | `updateHealthProfile` (`goal_adjust_gain_kcal`; «Usar el de la casa» = `null`) | N-327 |
| Restricciones alimentarias | Texto opcional | — | **Vacío** (decisión D: el campo existe y va vacío) | `listHealthProfiles` (`dietary_restrictions`) | `updateHealthProfile` (`dietary_restrictions`) | N-311, N-313 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listHealthProfiles` (`notes`) | `updateHealthProfile` (`notes`) | N-313 |

«Mantener» **no ajusta y no tiene campo** (K36): con ese objetivo, los dos ajustes se ven atenuados con «Con “mantener” no hay ajuste». El campo que corresponde al objetivo elegido va primero y resaltado. «Usar el de la casa» es el «restablecer» del ajuste de la persona (U4-13). Si la persona aún no tiene perfil: «Sin perfil de salud · Crear perfil» (`createHealthProfile` con `person_id`). El mismo ajuste por persona también se edita en P-21 (N-327).

**Ajuste por objetivo de la casa (K36, K44).**

| Parámetro | Clave | Valor actual e inicial | Validación y mensaje | Restablecer | Decisión | Lectura y escritura |
|---|---|---|---|---|---|---|
| Ajuste de la casa para bajar | `health.goal_adjust_lose_kcal` | Actual: `value` · Inicial: 0 kcal | Entero ≥ 0 · «Escribe un número entero de kcal, 0 o más» | Vuelve a 0 kcal | K36, K44 | `getHealthSettings` · `updateHealthSettings` |
| Ajuste de la casa para subir | `health.goal_adjust_gain_kcal` | Actual: `value` · Inicial: 0 kcal | Entero ≥ 0 · «Escribe un número entero de kcal, 0 o más» | Vuelve a 0 kcal | K36, K44 | `getHealthSettings` · `updateHealthSettings` |

Debajo: «Lo usan las personas que dejan vacío su propio ajuste: hoy Axel y Ximena» (según `goal_adjust_*` nulos de `listHealthProfiles`).

**Metas de ahorro y fondo.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Nombre | Texto | No vacío · «Escribe un nombre para la meta» | Sin valor: lo captura Axel | `listSavingsGoals` (`name`) | `createSavingsGoal` / `updateSavingsGoal` (`name`) | N-312, N-313 |
| Tipo | Selector: Fondo de ahorro · Meta con fecha | Obligatorio · «Elige fondo o meta» | Sin valor: lo elige Axel | `listSavingsGoals` (`kind`) | `createSavingsGoal` / `updateSavingsGoal` (`kind`) | N-312, N-313 |
| Monto objetivo | Dinero | Obligatorio en una meta; monto > $0.00 · «Escribe el monto de la meta, como $1,234.56» | Sin valor: lo captura Axel | `listSavingsGoals` (`target_cents`) | ídem (`target_cents`) | N-312, N-313 |
| Fecha meta | Fecha | Obligatoria en una meta; no antes del inicio · «La fecha meta no puede quedar antes del inicio» | Sin valor: lo captura Axel | `listSavingsGoals` (`target_date`) | ídem (`target_date`) | N-312, N-313 |
| Aportación mensual planificada | Dinero | Vacío o ≥ $0.00 · «Escribe un monto como $1,234.56» | Sin valor: lo captura Axel (decisión B: aportación mensual planificada) | `listSavingsGoals` (`planned_monthly_cents`) | ídem (`planned_monthly_cents`) | N-312, N-313 |
| Inicio | Fecha | Obligatoria · «Escribe la fecha de inicio» | Sin valor: lo captura Axel | `listSavingsGoals` (`start_date`) | ídem (`start_date`) | N-313 |
| Activa | Interruptor con vista previa en línea (parte a §0.2: no guarda al tocar; vista previa 3 de metas) | — | Encendida (valor de la base al crear; `is_active` de `savings_goal`) | `listSavingsGoals` (`is_active`) | `updateSavingsGoal` (`is_active`; «Desactivar», no se borra) | N-313 |
| Ahorrado, avance, aportación necesaria y estado | Solo lectura | — | Calculados (motor S2; `S2-metas` RESUELTO): avance con un decimal (`progress_tenths_pct` de C1 v1.3, décimas con ROUND_HALF_UP; Q2-17 C: un decimal en Metas y entero solo en la tarjeta de Inicio); aportación necesaria = (monto meta − ahorrado) ÷ meses que faltan (C1-datos M-16), con los meses del mes de hoy al de la fecha meta, los dos incluidos (Q2-16 A), base ROUND_HALF_UP y la última aportación del plan ajusta el centavo (`S2-interfaz.md` 1.3 §3.13); estado con tres valores, en curso, cumplida y vencida (Q2-18 A). La aportación de una meta vencida sin cumplir es **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/S2-metas-vencida.md`; en uso, todo el restante en el mes de hoy) | `listSavingsGoals` (`saved_cents`, `progress_tenths_pct`, `required_monthly_cents`, `status`) | — | N-312 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listSavingsGoals` (`notes`) | ídem (`notes`) | N-313 |

La aportación planificada al fondo general es la fila de tipo «Fondo de ahorro» con su aportación mensual. Acción primaria: «Nueva meta» (`Alt+N`).

**Vista previa del efecto.** **Perfil de salud** (vista previa en línea, parte a §0.2, calculada en el navegador con RN-17; al guardar manda el servidor): al cambiar cualquier dato del perfil que entra a una fórmula, el bloque «Antes → Después» bajo el campo muestra las cifras que cambian. **Estatura** → IMC con su clase OMS (K7 y RN-17: `weight_g` × 1000 ÷ `height_mm`², con el último peso registrado, un decimal con ROUND_HALF_UP; cortes 18.5, 25 y 30) y meta sugerida. **Fecha de nacimiento** (edad en años cumplidos a hoy, America/Mexico_City), **sexo** (+ 5 en hombres, − 161 en mujeres), **actividad**, **objetivo** y **mi ajuste para bajar o subir** → meta sugerida (el IMC no depende de ellos y no se muestra). La meta sugerida va con su desglose: basal de Mifflin-St Jeor × factor de actividad, menos el ajuste de bajar o más el de subir (el de la persona o, si lo deja vacío, el de la casa), con un solo redondeo ROUND_HALF_UP al final (RN-17, K36). Si falta un insumo (sexo, fecha de nacimiento, estatura, peso, actividad u objetivo), el lado sin cálculo dice «Falta {dato}» (RN-17: sin sugerencia mientras falte). Ejemplos ilustrativos de RN-17 (no canon): hombre de 35 años, 72.5 kg, actividad moderada y objetivo mantener; al capturar la estatura de 1.75 m: «IMC: Falta la estatura → 23.7 (normal) · Meta sugerida: Falta la estatura → 1,648.75 × 1.55 = 2,555.5625 → 2,556 kcal»; al cambiar el objetivo a bajar con su ajuste de 500 kcal: «Meta sugerida: 2,556 → 2,056 kcal» (2,555.5625 − 500 = 2,055.5625). Con los valores de prueba del hombre de S4-07, al cambiar el ajuste para bajar: «Basal 1,748.75 × 1.2 = 2,098.50 − ajuste para bajar 300 = 1,798.50 → meta sugerida 1,799 kcal» (con los valores iniciales 0 y 0 la meta queda en 2,098.50 antes de redondear). La meta confirmada no cambia sola: la persona confirma con «Usar la sugerida» o escribe otra (K23). Ajuste de la casa: «Con n kcal, la meta sugerida de Axel pasaría de a a b» por cada persona que usa el de la casa (fórmula RN-17 en el navegador; al guardar manda el servidor). **Metas** (vista previa en línea, parte a §0.2, calculada en el navegador): (1) al cambiar el monto objetivo, «Avance: {ahorrado} de {objetivo} = {p} %» antes → después, con `saved_cents` de `listSavingsGoals` y la fracción exacta ahorrado ÷ monto objetivo (RN-18), redondeada en décimas con ROUND_HALF_UP, en el mismo formato que la columna «avance» de la tabla de metas (`progress_tenths_pct`, un decimal; Q2-17 C, bloqueo `docs/estado/bloqueos/S2-metas.md` pregunta 2, RESUELTO), para que «Antes» sea igual a lo que ya se ve. Ejemplo ilustrativo de RN-18: meta de $6,293.44 con $3,146.72 ahorrados = 50.0 % (314,672 × 1,000 ÷ 629,344 = 500 décimas); (2) al cambiar la aportación mensual planificada, «El plan de {mes en curso} reserva {aportación nueva} para esta meta (antes {aportación actual}); el total reservado para metas en el mes pasa de {suma actual} a {suma nueva}» (RN-18: la aportación planificada se reserva en el plan del mes; la suma es la de `planned_monthly_cents` de las metas activas); (3) al tocar «Activa» (interruptor con vista previa en línea, excepción de la parte a §0.2): «Total reservado para metas en {mes en curso}: {antes} → {después}», la misma suma de `planned_monthly_cents` de las metas activas con esta meta sumada (al encender) o quitada (al apagar); una meta sin aportación planificada no cambia el total y el bloque lo dice: «Sin aportación planificada: el total reservado no cambia»; el interruptor no guarda al tocar: `Enter` o salir de él guarda y `Esc` lo devuelve a como estaba. Con `Ctrl+/` sobre varias metas, el mismo total sumado sobre las seleccionadas antes de aplicar. **Aportación necesaria y estado** (`required_monthly_cents`, `status`, motor S2; FIRMES desde la puerta 2, bloqueo `docs/estado/bloqueos/S2-metas.md`, preguntas 1 y 3, RESUELTO): al cambiar el monto objetivo o la fecha meta, «Aportación necesaria: {antes} → {después} al mes · Estado: {antes} → {después}», calculado en el navegador en centavos enteros con la regla de `docs/contratos/motores/S2-interfaz.md` 1.3 §3.13 (sin la solicitud SC-U4-2 a C1, que solo hacía falta si no se podía calcular aquí): meses que faltan = del mes de hoy al mes de la fecha meta, los dos incluidos (Q2-16 A; una fecha meta de un mes anterior da 0); aportación necesaria = ROUND_HALF_UP((monto meta − ahorrado) ÷ meses que faltan) (C1-datos M-16; la última aportación del plan ajusta el centavo); estado (Q2-18 A) = «Cumplida» si lo ahorrado alcanza el monto meta (aportación $0.00); si no, «Vencida» si la fecha meta ya pasó; si no, «En curso». Ejemplo ilustrativo (no canon), hoy 11-oct-2026, con la meta de $6,293.44 y $3,146.72 ahorrados: con fecha meta en diciembre de 2026 faltan 3 meses y la aportación necesaria es $1,048.91 al mes (314,672 ÷ 3 = 104,890.67 → 104,891 centavos; la última, $1,048.90); al cambiar la fecha meta a noviembre de 2026, $1,048.91 → $1,573.36 al mes, estado «En curso» → «En curso». **PROVISIONAL** solo la aportación de una meta vencida sin cumplir (meses que faltan 0 y restante mayor que $0.00; bloqueo abierto `docs/estado/bloqueos/S2-metas-vencida.md`; en uso, la del motor: todo el restante en el mes de hoy): ahí el bloque rotula «Aportación necesaria: vista previa pendiente; se recalcula al guardar» (parte a §0.2) y muestra el estado «Vencida».

**Estados.** Vacío: sin perfil = «Sin perfil de salud · Crear perfil»; sin metas = «Sin metas de ahorro · Nueva meta»; sin peso = «Sin peso registrado · Registrar peso» (IMC y meta sugerida dicen «Falta el peso»). Carga: esqueleto de las dos fichas y de la tabla. Error al guardar: el de la parte a §0.3 (el campo vuelve a su valor anterior). Guardado: automático.

**Edición en el lugar.** Todos los campos del perfil y de las metas que no son calculados; los que tienen vista previa (los datos del perfil que entran a RN-17, mis ajustes, el ajuste de la casa, el monto objetivo, la aportación planificada y el interruptor «Activa», que no guarda al tocar) guardan después de ella (regla general de la parte a §0.2, Q2-09 A: sexo, actividad y objetivo son listas con vista previa); en la tabla de metas, arrastrar una tarjeta entre columnas del tablero no cambia el estado (es calculado). No se editan: edad, peso actual, IMC, meta sugerida, ahorrado, avance y estado.

**Atajos.** `Alt+N` nueva meta; `↑`/`↓` y `Enter` en la tabla; `Ctrl+/` cambiar a la vez «Activa» de varias metas, con la vista previa (3) sumada sobre las seleccionadas antes de aplicar (`Enter` aplica a todas y `Esc` cancela sin cambiar ninguna); `←`/`→` mes anterior o siguiente en el calendario de metas.

**Móvil y escritorio.** 412 px: dos fichas apiladas (Axel, Ximena) con el desglose de la meta bajo el campo, luego el ajuste de la casa y las metas como tarjetas con barra de avance etiquetada. 1440 px: las dos fichas lado a lado, el ajuste de la casa debajo y la tabla de metas con su selector de vista.

## P-26.14 Inventario

**Propósito.** Decidir cómo se lleva cada producto en el inventario híbrido (decisión D): existencia por unidad (básicos y no perecederos) o ritmo de consumo (el resto), con su nivel de reposición. Regla fija: los sobrantes de recetas entran como existencia (decisión D). Necesidades N-314 y N-315.

**Vista de base de datos.** Tabla (por defecto) y tablero agrupado por modo (U1 §2.3); filtros: categoría, modo, estado y texto; orden: agotado, reponer, suficiente y luego nombre (U1 §2.5). En el tablero, arrastrar una tarjeta de «Existencia por unidad» a «Ritmo de consumo» cambia su modo después de la vista previa (parte a §0.2, punto 2).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Producto | Texto, solo lectura (o sobrante de receta, en porciones) | — | — | `listInventoryItems` (`product_id`, `recipe_id`) y `listProducts` | — | N-314 |
| Modo | Selector: Existencia por unidad · Ritmo de consumo | Uno de los dos | Sin valor: lo elige Axel por producto (los sobrantes, existencia por unidad, decisión D) | `listInventoryItems` (`tracking_mode`) | `updateInventoryItem` (`tracking_mode`) | N-315 |
| Unidad | Texto, solo lectura: g, ml o pieza (porción en sobrantes) | — | La unidad base del producto | `listProducts` | — | N-314 |
| Nivel de reposición | Cantidad en la unidad base (hasta 3 decimales) | Vacío o ≥ 0 · «Escribe una cantidad en g, ml o piezas, o deja vacío para no avisar» | **Vacío = sin aviso** (C1 `reorder_level_milli`) | `listInventoryItems` (`reorder_level_milli`) | `updateInventoryItem` (`reorder_level_milli`) | N-315 |
| Ritmo: cantidad | Cantidad en la unidad base | Con modo «ritmo», > 0 · «Escribe cuánto se consume» | Sin valor: lo captura Axel | `listInventoryItems` (`rate_quantity_milli`) | `updateInventoryItem` (`rate_quantity_milli`) | N-315 |
| Ritmo: cada cuántos días | Entero de días | Con modo «ritmo», entero > 0 · «Escribe en cuántos días se consume esa cantidad» | Sin valor: lo captura Axel | `listInventoryItems` (`rate_period_days`) | `updateInventoryItem` (`rate_period_days`) | N-315 |
| Caducidad | Fecha opcional | — | Sin valor | `listInventoryItems` (`expires_date`) | `updateInventoryItem` (`expires_date`) | N-315 |
| Existencia, estado y días restantes | Solo lectura | — | Calculados (K7: reponer si existencia ≤ nivel; agotado si 0) | `listInventoryItems` (`quantity_milli`, `estimated_quantity_milli`, `status`, `days_left`) | — (la existencia se ajusta en P-16) | N-314 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listInventoryItems` (`notes`) | `updateInventoryItem` (`notes`) | N-315 |

Acción secundaria: «Agregar producto al inventario» (`Alt+N`): `createInventoryItem` con producto, modo y existencia inicial. Con modo «ritmo» y sin cantidad o días, la fila dice «Falta el ritmo» y no muestra existencia estimada.

**Vista previa del efecto.** Al cambiar nivel o modo, antes de guardar (el modo es una lista con vista previa: muestra el bloque al elegir y guarda con `Enter` o al salir, parte a §0.2): por fila, el estado antes → después, «Con nivel n g y existencia m g: Suficiente → Reponer» (regla K7, en el navegador), y al pasar a «Ritmo de consumo» también la existencia estimada de hoy (RN-15; sin ritmo capturado, «Falta el ritmo»). Al soltar una tarjeta en la otra columna del tablero, el mismo bloque de su fila; con `Ctrl+/` sobre {n} filas (modo o nivel), «{k} pasan a Reponer, {j} a Agotado y {s} a Suficiente» y la lista de las filas que cambian; `Enter` o «Aplicar» guarda todas y `Esc` cancela sin cambiar ninguna (parte a §0.2); al cambiar el ritmo (cantidad o cada cuántos días), vista previa en línea (parte a §0.2) calculada en el navegador con RN-15: «Existencia estimada hoy: {antes} → {después}» = existencia contada − HALF_UP(cantidad del ritmo × días desde el conteo ÷ días del ritmo), nunca menor que 0, y su estado con la regla K7 (reponer si la existencia ≤ nivel de reposición; agotado si es 0), antes → después. Ejemplo de RN-15 (ritmo ilustrativo): 3 botellas de 946 ml = 2,838 ml contadas; con una botella cada 21 días, a los 30 días la existencia estimada es 2,838 − 1,351.429 = 1,486.571 ml. **Días restantes y consumo por semana** (`days_left`, `weekly_consumption_milli`, motor S4), antes → después, en el navegador con la regla de Q2-70 (a) (bloqueo `docs/estado/bloqueos/U4-vista-previa-formulas.md`, pregunta 3, RESUELTO; `docs/contratos/motores/S4-interfaz.md`): días restantes = piso(existencia estimada × días del ritmo ÷ cantidad del ritmo) y consumo por semana = HALF_UP(cantidad del ritmo × 7 ÷ días del ritmo), en milésimas de la unidad base; solo en modo «ritmo». Con el ejemplo de RN-15: piso(1,486.571 × 21 ÷ 946) = 32 días (el cociente es 32.99999…) y HALF_UP(946 × 7 ÷ 21) = 315.333 ml por semana.

**Estados.** Vacío: «Aún no hay productos en el inventario · Agregar producto». Carga: esqueleto de la tabla. Error al guardar: el de la parte a §0.3 (la celda vuelve a su valor anterior). Guardado: automático.

**Edición en el lugar.** Modo, nivel, ritmo, caducidad y notas (celda o arrastre en el tablero; el modo y el arrastre, con la vista previa antes de guardar, parte a §0.2). No se editan: producto, unidad, existencia (P-16), estado y días restantes.

**Atajos.** `Alt+N`; `Ctrl+/` cambia a la vez el modo o el nivel de las filas seleccionadas, con la vista previa sumada antes de aplicar (`Enter` aplica a todas, `Esc` cancela); `↑`/`↓` y `Enter`.

**Móvil y escritorio.** 412 px: tarjetas con producto, modo, estado (icono y texto) y nivel; el tablero se ve como dos listas plegables. 1440 px: tabla completa o tablero de dos columnas.

## Parámetros operativos de la parte b (K44, U4-12)

Filas de la lista mínima de U4-12 que viven en estas secciones (las demás, 6 a 11, van en la parte a: precio viejo y peso cero en P-26.06, sesión, intentos y ventana en P-26.01, comparador en P-26.05).

| # U4-12 | Parámetro | Sección | Clave | Inicial | Decisión |
|---|---|---|---|---|---|
| 1 | Primer aviso del presupuesto | P-26.09 | `alerts.budget_warning_pct` | 80 % | Decisión F, K44 |
| 2 | Aviso de tope | P-26.09 | `alerts.budget_limit_pct` | 100 % | Decisión F, K44 |
| 3 | Días antes del pago próximo | P-26.09 | `alerts.payment_due_days_before` | [1, 0] | Decisión F, K44 |
| 4 | `send_time` | P-26.08 (reportes) y P-26.09 (avisos con hora) | `reports.daily_send_time`, `reports.weekly_send_time`, `reports.monthly_send_time` y `alerts.send_time` (Q2-11 B, K47) | 21:30, 19:00, 08:00 y 09:00 | Q2-72 y Q2-73 (K50), K43 |
| 5 | `advance_days` | P-26.09 | `alerts.stale_prices_advance_days` (Q2-11 B, K47) | 1 día | Q2-73 (K50), K43 |
| 12 | Adjunto máximo | P-26.11 | `attachment.max_size_mb` | 10 MB | K43 |
| — | Retención de claves de idempotencia | P-26.11 | `idempotency.retention_hours` | 24 h | K6, K44 |
| — | Tiempo máximo del OCR | P-26.11 | `ocr.timeout_seconds` | 30 s | Q2-76 (K50) |
| 13 | Ajuste de la casa para bajar | P-26.13 | `health.goal_adjust_lose_kcal` | 0 kcal | K36 |
| 14 | Ajuste de la casa para subir | P-26.13 | `health.goal_adjust_gain_kcal` | 0 kcal | K36 |
| 15 | `weekday` | P-26.08 (`weekly_summary`) y P-26.09 (`meal_plan_missing`) | `reports.weekly_weekday` y `alerts.meal_plan_missing_weekday` (Q2-11 B, K47) | 7 (domingo) | Decisión F |
| 16 | `day_of_month` | P-26.08 (`monthly_report`) | `reports.monthly_day_of_month` (Q2-11 B, K47) | 1 | Decisión F, Q2-66 B |

## Campos de la puerta 1: dónde quedan

Tope por persona (K13) → P-26.12 · objetivo de salud y meta confirmada (K23) y ajuste por persona y de la casa (K36) → P-26.13 · candidatos de chat (K4) → P-26.09 · IA desactivable (decisión H) → P-26.11 · restricciones duras (K16), interruptor A/B (decisión C), promociones, tasas como se capturaron (K12) y datos de tarjeta (K18) → P-26.05 · proyección de los meses 7 a 12 (K21) → P-26.02 · mercados, membresía de Sam's y OCR local como fuente de precio → P-26.06 (parte a).

## Matriz de operaciones de la parte b (U4-02)

- P-26.08: `listNotifications`, `updateNotification`, `setNotificationRecipients`, `listNotificationRecipients`, `listJobs`, `sendTestNotification`, `generateMonthlyReport`, `listMonthlyReportMonths`, `generateExternalStatement`, `getReportSettings`, `updateReportSettings` (día y hora, Q2-11 B; C1 v1.3).
- P-26.09: `listNotifications`, `updateNotification`, `setNotificationRecipients`, `sendTestNotification`, `listNotificationRecipients`, `updateNotificationRecipient`, `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate`, `getAlertSettings`, `updateAlertSettings`, `listPersons`; vista previa: `getBudgetMonth`, `getDashboardNextPayment`.
- P-26.10: sin operación (Local, U1 §0.3, N-301).
- P-26.11: `listBackups`, `createBackup`, `restoreBackup`, `listJobs`, `updateJob`, `previewPlannedItemsImport`, `importPlannedItems`, `getSeedStatus`, `listAiProviders`, `updateAiProvider`, `getSystemStatus`, `getDataSettings`, `updateDataSettings`.
- P-26.12: `listPersons`, `updatePerson`, `listCostShareRules`, `updateCostShareRule`.
- P-26.13: `listHealthProfiles`, `createHealthProfile`, `updateHealthProfile`, `getHealthSettings`, `updateHealthSettings`, `listSavingsGoals`, `createSavingsGoal`, `updateSavingsGoal`.
- P-26.14: `listInventoryItems`, `listProducts`, `updateInventoryItem`, `createInventoryItem`.

## Cobertura de criterios U4 en la parte b

U4-01: 7 secciones `## P-26.xx` con tabla de campos de 7 columnas. U4-02: matriz de arriba. U4-03: decisión F en P-26.08 y P-26.09; K4 en P-26.09; decisión G en P-26.10; decisión H en P-26.11; decisión B y K13 en P-26.12; decisión D, K23 y metas en P-26.13; decisión D en P-26.14. U4-04: valores por defecto con su cita en cada tabla. U4-05: «Vista previa del efecto» en las 7; en P-26.13, IMC y meta sugerida con cualquier dato del perfil que entra a RN-17, y en metas la aportación necesaria y el estado siguen Q2-16 A, Q2-17 C y Q2-18 A (`S2-metas`, RESUELTO), con PROVISIONAL solo la aportación de una meta vencida sin cumplir (`S2-metas-vencida`, abierto); el interruptor «Activa» de una meta y `Ctrl+/` muestran el total reservado para metas antes de guardar (R4-1). U4-06: vistas, edición en el lugar y atajos por sección. U4-07: columna «Validación y mensaje». U4-08: restaurar, activar IA y descartar chat. U4-09: «Móvil y escritorio» por sección y tokens de DESIGN. U4-10: «Estados» por sección. U4-11: SUPUESTOS y bloqueos `U4-dia-mes-reporte`, `U4-reparto-retroactivo` y `U4-vista-previa-formulas` (RESUELTOS en la puerta 2), `U4-reparto-preseleccion` (abierto, P-26.12) y `S3-iva-lectura` (abierto, parte a); las vistas previas y lecturas que dependen de un bloqueo abierto de otro artefacto lo citan: `C1-hora-reportes-job` en P-26.08, `C1-reparto-alcance` en P-26.12, `S2-metas-vencida` en P-26.13, y `S3-formula-tarjeta-pdf` y `S3-iva-lectura` en P-26.05 (parte a). U4-12: tabla de parámetros. U4-13: P-26.13.

## SUPUESTOS

1. Sección de cada clave: la que da su `@param` (00003 y 00004), que coincide con la que proponía esta parte: hora y día de los reportes en **Reportes y frecuencias**; hora, anticipación y día del plan sin armar en **Notificaciones**; retención de claves de idempotencia en **Datos y respaldos** (como U1 §2.3). `ocr.timeout_seconds` en **Datos y respaldos** es PROPUESTA de S1 que la puerta 2 no preguntó (`docs/estado/p2c/S1-puerta2.md`, «Propuestas para Axel»; la otra opción es Tiendas y precios).
2. A qué alertas aplica la hora lo fija S7; el campo aparece donde S7 lo confirme (el aviso de presupuesto salta al registrar el gasto que cruza el umbral y no necesita hora).
3. Saturación y densidad son propiedades fijas del sistema de diseño (DESIGN §1.2 P4 y P6) y no se configuran (U1: Apariencia tiene un solo control).
4. La hora del respaldo diario se edita con `updateJob` sobre la tarea del respaldo; C1 hoy liga `updateJob` solo a N-297.
5. Rangos: los aprobados en la puerta 2 (K50; parte a §0.5). No hay otros mínimos de seguridad (p. ej. no menos de 3 intentos).
6. La preselección del diálogo de alcance de P-26.12 no la fija ninguna fuente: queda **PROVISIONAL** con la opción en uso, «Solo a los gastos nuevos», y la decide Axel en la puerta 4 (bloqueo abierto `docs/estado/bloqueos/U4-reparto-preseleccion.md`). Con «solo los nuevos», con qué partes se recalcula un gasto anterior cuyo importe se edita después es PREGUNTA abierta para C1 y el motor de reparto (bloqueo abierto `docs/estado/bloqueos/S1-reparto-edicion-gasto.md`, WF-3), no de esta pantalla.
7. Las claves `reports.*` se leen y escriben con `getReportSettings`/`updateReportSettings`, la sección «Reportes y frecuencias» de C1 v1.3 (GET `/api/v1/settings/reports`) (PUT `/api/v1/settings/reports`); las `alerts.*` siguen en `getAlertSettings`/`updateAlertSettings` (P-26.09). Qué manda entre esas claves y el `schedule` de la tarea del reporte es PROVISIONAL (bloqueo abierto `docs/estado/bloqueos/C1-hora-reportes-job.md`).
8. El interruptor «Atajos de una tecla» (P-26.10) usa la necesidad de Apariencia (N-301) porque U1 §3.28 no tiene una propia; U1 lo agrega (solicitud del reporte).
