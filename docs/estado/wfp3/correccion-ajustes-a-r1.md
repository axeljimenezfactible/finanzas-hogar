# Corrección ajustes-a · vuelta 2 (agente `correccion-ajustes-a-r1`)

- Fecha: 12-oct-2026. Alcance: los 16 artboards de la parte a en `docs/design/maquetas/lienzo/project/ajustes/` (índice y s01 a s07, `-m` y `-e`) y `fragmentos/ajustes-a.json`.
- Copia de la vuelta 1 (P3-06): `docs/design/salida/iter-1/ajustes-a/` (16 artboards, el fragmento y `artboards-ajustes-a.md`).
- Fuente de hallazgos: `P3-iter-1-ajustes-a.md` (0 altos, 8 medios, 7 bajos).

## Hecho

| Hallazgo | Qué se hizo |
|---|---|
| M-01 enlaces a s08 | Slug `s08-reportes` en los 16 archivos (10 enlaces); 0 `href` a `frecuencias`. |
| M-02 `#instrucciones` | Bloque `id="instrucciones"` (3 pasos) bajo el certificado raíz en s01 `-e` y `-m`. |
| M-03 objetivos táctiles | «Ver instrucciones» (s01-m) y «Ir a Personas y reparto» / «Ir a Datos y respaldos» (s05-m, s06-m, s07-m) con `btn t` (44 px); casilla «Ya lo confirmé» (s05) en fila de 44 px con `label` que ocupa toda la fila. |
| M-04 estados | Tarjeta «Error al guardar» (circle-alert + texto de U4 §0.3) en los 8 `-e` (en el índice dice «No aplica: no tiene campos»). s01: vacío «No aplica: siempre está la sesión actual». Índice: vacío «No aplica: las 14 secciones siempre existen». Añadidos: «Sin gastos fijos · Nuevo gasto fijo» (s02), «Sin tasas · Nueva tasa» y «Sin promociones · Nueva promoción» (s05), «Sin cortes importados · Importar corte PROFECO» (s06). |
| M-05 validaciones | Un campo con `ui-border-error`, circle-alert y mensaje literal en cada sección, `-m` y `-e`: s01 V-15 (campo real «Repite la contraseña nueva»), s02 V-03 (Ocio, en la tabla), s03 V-04, s04 V-12, s05 V-23 (1 a 30), s06 V-25, s07 V-09. En s03 a s07 va como bloque «Así se ve un valor no válido» (mismo patrón que «Así se ve el diálogo»), para no contradecir los campos en edición con su vista previa. |
| M-06 fechas | 16 fechas sin cero a la izquierda (1-oct-2026, 3-oct-2026, 6-nov-2026, 1-ene-2027, etc.). El revisor contó 17 apariciones; el grep final da 0. |
| M-07 barra lateral | En los 8 `-e`: encabezado con avatar y menú (Cuenta, Apariencia, Cerrar sesión), botón `panel-left` «Contraer barra lateral» (Ctrl+\), grupo «Tema» (Sistema, Claro, Oscuro) al pie; `Ctrl+\ contraer la barra lateral` en la línea de atajos. |
| M-08 (1) calendario | s05-e: vista Calendario activa (franjas de inicio a fin en cuadrícula de 35 días, del 19-oct al 22-nov-2026, semanas rotuladas, «Vigente hoy» con franja verde). La tabla queda como vista alternativa del selector. |
| M-08 (2) Eliminar promoción | Botón por promoción en s05 `-e` y `-m` y diálogo crítico (Cancelar primero, confirmar en crítico) «Eliminar la promoción de Walmart». |
| M-08 (3) Quitar monto | s02: «Quitar el monto de {categoría}» por fila (`-e` en columna Acción, `-m` en cada tarjeta) y diálogo crítico «Quitar el monto de Ocio». |
| M-08 (4) Activa | s04: `.sw role="switch"` con `aria-checked` y nombre «Activa: {categoría}» en las 13 filas de `-e` y `-m`. |
| M-08 (5) plegables | s05-m: 7 bloques con `<h2><button aria-expanded aria-controls>` de 44 px; abiertos: Líneas, Condiciones de Kueski Pay, Lectura A o B, Promociones y Comparador; plegados: Condiciones de Banorte Oro y Tasas por plazo. |
| Bajos resueltos de paso | B-01 título «Ajustes»; B-02 «No se pudieron cargar…»; B-04 asa de s04 como botón «Mover {categoría}» de 44 px; B-06 parcial en s01 (la primaria pasa a «Cerrar otras sesiones») y B-07 (el disparador «Cerrar sesiones» ya no lleva relleno crítico fuera del diálogo). |

## Verificación (`python3 -I`, script propio, `RESULTADO: OK`)

- Forma de los 16: doctype, `lang="es-MX"`, charset, `./support.js`, `../ds/finanzas-hogar/tokens.css`, `x-dc`, `helmet`, `data-dc-script`, `class Component extends DCLogic`, `renderVals()`.
- `data-props` parsea y `$preview` = `w` x `h` del fragmento; `expand: "fill"` solo en `-e`; raíz móvil con `width: var(--breakpoint-phone)`.
- `html.parser`: 0 etiquetas sin cerrar o mal anidadas; 0 ids repetidos; todo `for`, `aria-describedby`, `aria-controls` y `aria-labelledby` apunta a un id existente; 0 anclas `#x` sin destino; 0 enlaces rotos a archivos del lienzo propio (incluido s08-b).
- 0 hex literales fuera de `data-props`; todos los `var(--x)` existen en `tokens.json`.
- 0 cifras prohibidas de P3-12 (con y sin coma, y en centavos); 0 credenciales (el campo de confirmación lleva doce `x` como ejemplo).
- 0 fechas con cero a la izquierda; 0 negativos con guion; 0 emoji; 0 iframe, img, innerHTML, red.
- Fragmento: JSON válido, 16 `boards` y 16 `order`, exactamente los archivos de la parte a.
- Contraste de lo nuevo (tokens resueltos, claro y oscuro): `ui-role-error` y `ui-border-error` sobre página, tarjeta, campo y superficie elevada: mínimo 5.40 (claro) y 5.82 (oscuro); `ui-role-success` sobre su superficie 4.81 y 6.19; texto por defecto sobre superficies de franja 10.4 o más.

## Supuestos

1. Las alturas del fragmento y de `$preview` de los artboards tocados son estimaciones (sin navegador, K45): suman el contenido añadido; s05-m baja a 7,600 por los dos bloques plegados. Hay que medirlas al publicar (B-05 sigue abierto).
2. Los textos de «Instrucciones» del certificado son de redacción propia (U4 no los fija).
3. La hoja de «Hoy es 21-oct-2026» sigue la fecha de la pantalla de shell; las franjas del calendario se calculan de las vigencias de `datos-sinteticos.json`.
4. $26,686.00 (diálogo «Quitar monto de Ocio») es 27,586.00 menos 900.00, suma derivada del JSON.
5. `ui-role-success` y `ui-border-selected` se usan como borde de franja; ambos existen en `tokens.json`.

## Bloqueos

Ninguno.

## Pendientes

- B-03 (leyenda «Enter o salir del campo guarda» bajo vistas previas que no son de un campo, s02 y s07) y B-06 restante (s02 con 3 primarios, s05 con 2; `.tbar` móvil sin la primaria): no tocados, son bajos.
- B-05: medir las alturas con navegador.
- La tabla de promociones de s05-e ya no se dibuja (solo se llega por el selector); si el revisor la quiere visible junto al calendario, es un añadido sin conflicto.

## Consumo aproximado

Leídos unos 330 KB, escritos unos 60 KB (lo tecleado; los archivos editados por script suman más), 34 acciones. Sin cierre por tope.
