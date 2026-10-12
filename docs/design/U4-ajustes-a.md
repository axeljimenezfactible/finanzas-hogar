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
