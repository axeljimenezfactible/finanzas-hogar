# correccion-ds-r1 · Pista 3 fase 1 (sistema de diseño) · 11-oct-2026, noche

Agente `correccion-ds-r1`. Corrige los hallazgos alto y medios de `REVISION-ds-r1.md` (A-01, M-01, M-02, M-03). Los bajos (B-01 a B-05) no estaban en el encargo y no se tocaron. Alcance: `docs/design/sistema/project/`, `docs/design/tools/ds/`, `docs/estado/bloqueos/P3-*` y este reporte. Sin git, sin red, sin navegador, sin Artifact. Respaldo previo en `/tmp/correccion-ds-r1-respaldo/` (project/, ds/ y los dos bloqueos P3).

## Hecho
1. **A-01 · contraste del chip en hover.** `components/ChipPersonaTienda/README.md` declara la excepción con los 4 valores (Karen 2.83 y Sam's Club 2.74 en claro; Axel 2.14 y Mercado 2.11 en oscuro, sobre `ui-surface-hover`) y nombra el bloqueo `P3-contraste-chip-hover.md` con la opción A recomendada. **No se cambió ni el color ni la regla `button.chip:hover`**: se hará cuando Axel responda. El bloqueo lleva una nota de actualización.
2. **M-01 · inventario de componentes.** 5 componentes nuevos, cada uno con `README.md` y `preview.html` (marcador en la línea 1, `lang="es-MX"`, solo `var(--token)`, iconos copiados del sprite de Lucide 1.55.0, sin red ni iframes):
   - `BarraSuperior` (§7.10, G-02 y pestañas): móvil 56 px con «‹ Créditos», primaria y «⋯»; desde 768 con 48 px, migas completas con `chevron-right`, Buscar, «+ Registrar» y «⋯»; título `title-1`, propiedades y pestañas Resumen · Detalle · Tendencias · Calendario (`role="tablist"`, indicador de 2 px en `ui-border-selected`).
   - `BarraInferior` (§7.10, G-05): 5 destinos en orden fijo, píldora de 56 × 32 en `ui-surface-selected`, «+» de 40 px en `ui-action-primary-bg`, `env(safe-area-inset-bottom)`.
   - `HojaDialogo` (§7.12): hoja inferior con asa, `ui-overlay-scrim`, `radius-lg` solo arriba, `shadow-2` y primaria de 48 px; diálogo corto de confirmación («¿Descartar cambios?»); diálogo de acción crítica con `shield-alert` en `ui-critical`, «Cancelar» con el foco inicial y «Borrar 12 gastos» (el ejemplo de DESIGN §7.2).
   - `VistaPrevia` (§7.12, G-06): `x`, título editable en el lugar (hover de solo borde), 4 propiedades con el gasto `−$318.45` en `ui-role-expense` y `arrow-up-right`, «Mostrar 2 más» y cuerpo; `layout-preview-width-*`, `radius-md`, `shadow-1`.
   - `PanelAjustes` (§7.13): las 14 secciones con sus iconos de §5.4 (Notificaciones abierta) y tres filas: parámetro operativo con «Valor inicial: 80 %», «Restablecer al valor inicial» (nombre accesible con lo que restablece) y «Guardado» con `circle-check`; otra igual a su valor inicial (sin «Restablecer»); hora `HH:MM`.
   - Ahora quedan cubiertas todas las filas de la matriz de DESIGN §7.16.
3. **M-02 · `font.weight.medium` como valor.** `c2_a_ds.py` agrega la familia de valores simples `fontWeight` con `font-weight-regular` 400, `font-weight-medium` 500, `font-weight-semibold` 600 y `font-weight-bold` 700. Los nombres son los de `tokens.css` (el script lo comprueba con `assert`) y el `usage` es la `$description` de C2 (DESIGN §3.2). `tokens.json` se regeneró: la diferencia con el anterior es **solo** el bloque `fontWeight` (9 líneas añadidas, 0 cambiadas). `verificar_ds.py` exige ahora los 4 pesos como valor y cuenta la fila `font` desde lo que encuentra (antes la cifra 18 estaba fija). Además se hizo el cambio opcional: los 16 `font-weight:400|500|600|700` de los 14 previews existentes pasan a `var(--font-weight-*)`. Se comprobó que en esos archivos no cambió nada más. El README del sistema (§Tipografía) nombra las 4 variables.
4. **M-03 · supuestos marcados.** Se agregó una línea «Supuesto provisional (P3-huecos-componentes, punto N): …» en BarraLateral (punto 1), InterruptorOpciones (puntos 2 y 3), CalendarioPagos y ChipPersonaTienda (punto 4) y MenuBarra (punto 5). Los 5 componentes nuevos llevan la suya (puntos 7, 8 y 9, nuevos en el bloqueo).
5. `P3-huecos-componentes.md`: el punto 6 se corrigió en parte (los pesos ya son variables; el 2 px y tamaños e interlíneas siguen como estaban) y se agregaron los puntos 7, 8 y 9 con su pregunta. En `design-system.json` se actualizó `lastChange` (`at` 2026-10-12T05:20:19Z y la nota de esta corrección).

## Verificación (comandos y salidas)
```
$ python3 -I docs/design/tools/ds/c2_a_ds.py          # antes del cambio: cmp con el respaldo = idéntico
tokens.json: 70593 bytes · color 266 · otras familias 89 · tipos 6 estilos · fuentes 2   (fontWeight 4)
$ python3 -I docs/design/tools/ds/verificar_ds.py
  font        C2  18  ->  tokens.json  18   (6 tamaños + 6 interlíneas + 1 espaciado en los estilos de type; 4 pesos como valor en fontWeight; 1 pila)
C2: 445 tokens · tokens.json: 355 tokens en listas (color 266 + otras 89) + 6 estilos de tipo + 1 pila + 2 fuentes
Comparaciones de color contra tokens.css (token × tema): 532
RESULTADO: OK · 445 tokens C2 -> 355 en listas de tokens.json + 6 estilos de tipo (18 font.*) + 2 fuentes; 0 fallas
# prueba negativa: font-weight-medium 500 -> 510 en tokens.json
FALLAS (3): valor distinto de tokens.css: font-weight-medium ds=510 css=500 · font.weight.medium no está como valor en fontWeight · conteo font: C2=18, esperado 18, tokens.json=17
RESULTADO: FALLA · 3 fallas      -> se regeneró: cmp idéntico y RESULTADO: OK
$ python3 -I docs/design/sistema/verificar_brandbook.py
previews 19 · variables distintas usadas 161 (todas en tokens.json: True) · cifras canon 0
BarraInferior texto mín 5.40/7.97, gráfico mín 5.15/7.94 · BarraSuperior texto 5.40/7.42, gráfico 3.71/5.50 · HojaDialogo texto 5.40/5.04, gráfico 3.71/6.09 · PanelAjustes texto 5.08/5.04, gráfico 3.73/3.73 · VistaPrevia texto 5.40/5.04, gráfico 4.51/5.50  (claro/oscuro)
components/ChipPersonaTienda/preview.html: grafico ui-person-karen sobre ui-surface-hover light 2.83:1 < 3.0 · excepción declarada …
RESULTADO: OK · 0 fallas · 1 avisos de excepción declarada
```
- Una primera corrida marcó «línea 1 sin marcador @dsCard» en HojaDialogo y PanelAjustes, porque llevaban `width=`, que format.md admite pero el verificador no. Se quitó `width` y quedó en OK.
- Contraste del chip recalculado con los hex de `tokens.json`: hover claro `#e5e6e9`: Karen 2.83 y Sam's Club 2.74; hover oscuro `#4d4e53`: Axel 2.14 y Mercado 2.11. Coincide con el revisor y con el README.
- P3-12 (grep de las 21 cifras canon del criterio) sobre `docs/design/sistema/project` = **0**; patrón de credenciales = **0**. Fechas nuevas (5-nov-2026 ya estaba en CalendarioPagos; 6-nov-2026): 0 coincidencias en `oraculo/model.json`. Montos nuevos: `$318.45` y `−$318.45` (no son canon).
- Previews nuevos: 0 hex, 0 `rgb(`, 0 URL `http(s)://`, 0 `iframe/object/embed/portal/noscript`; HTML bien cerrado (lo comprueba el verificador con `html.parser`).
- Previews existentes: al revertir `var(--font-weight-*)` al número, cada archivo queda idéntico al respaldo (14 de 14 «solo pesos»).
- Líneas «Supuesto provisional (P3-huecos-componentes…» presentes en 10 README (InterruptorOpciones con 2).

## Supuestos
- La clave de familia se llama `fontWeight`, como el ejemplo `zIndex` de format.md. El título de la sección («Font weight») lo genera el tipo a partir de la clave.
- Tamaños, interlíneas, `font-letter-spacing` y `font-numeric-tabular` no se agregaron como variables: el encargo pedía solo los pesos, y ya hay 12 familias de lista, el tope que comprueba `verificar_ds.py`. Siguen dentro de los 6 estilos de Tipografía (clases `.caption` … `.title-1`). Esto queda dicho en el punto 6 del bloqueo.
- Lo que DESIGN §7.10, §7.12 y §7.13 no fija va como supuesto PROVISIONAL en `P3-huecos-componentes.md`, puntos 7 a 9:
  - «+ Registrar» como secundario;
  - alto y relleno de la pestaña;
  - el indicador de 2 px como número;
  - «‹» dibujado con `chevron-left`;
  - migas subrayadas;
  - «+» de la barra inferior sin nombre visible;
  - rellenos de hoja, diálogo y vista previa;
  - título de la vista previa en `title-2` y de la sección de Ajustes en `title-3`;
  - texto «Guardado» en `ui-role-success`.
- Los textos de ejemplo de los previews nuevos son inventados y sintéticos (página «Crédito de despensa», «Súper del viernes», 85 %, 08:30, «Borrar 12 gastos» tomado de DESIGN §7.2). Así lo dicen sus README.
- En la vista previa los diálogos no llevan `aria-modal`, porque se muestran tres a la vez. El README lo exige en el producto.

## Bloqueos
- `P3-contraste-chip-hover.md` (abierto, pregunta a Axel): sin su respuesta no se cambia el hover del chip.
- `P3-huecos-componentes.md` (abierto, provisional): ahora con 9 puntos. No detiene la Pista 3.

## Pendientes
- **Axel:** responder el bloqueo del chip. Con la opción A, cambiar `button.chip:hover` a solo borde `ui-border-control-hover` sobre `ui-surface-raised` y quitar la excepción del README. Confirmar los puntos 1 a 9 de P3-huecos-componentes.
- **Fuera de mi alcance:** `docs/design/sistema/verificar_brandbook.py` declara como excepción solo el par de Karen y su cascada solo detecta ese par. Los otros 3 los calculé aparte (arriba). Conviene que el dueño de ese script lo extienda.
- **Orquestador, al publicar:**
  - comprobar que la familia `fontWeight` se compila en `:root` (`--font-weight-*`). Si no, los previews vuelven a peso 400; no se pierde color;
  - comprobar que el grupo nuevo «Hojas y diálogos» aparece;
  - revisar, con la comprobación de render, los 5 previews nuevos en los dos temas. Aquí no hubo navegador.
- **Hallazgos bajos no atendidos** (no estaban en el encargo):
  - B-01: fecha 30-oct-2026 en la muestra de `caption`;
  - B-02: ubicación de `design-system.json`;
  - B-03: `usage` de los colores de texto sin sus fondos;
  - B-04: `rx` literal en Cover;
  - B-05: ruta `tools/ds` contra `tools/tokens` y licencia OFL.
- **Ronda 2 de revisión:** contrastar los 5 componentes nuevos con DESIGN §7.10, §7.12 y §7.13.

## Consumo aproximado
- **Leído:** ≈ 115 KB:
  - revisión y reportes, 35 KB;
  - DESIGN §3.2, §7.0, §7.1 y §7.9 a §7.16, más filas C-129 a C-141, ≈ 30 KB;
  - format.md, 5 KB;
  - partes de los scripts, ≈ 22 KB;
  - README, previews y salidas, ≈ 23 KB.
- **Escrito:** ≈ 68 KB:
  - 5 componentes nuevos, 44.6 KB;
  - ediciones de README y previews, ≈ 5 KB;
  - scripts, ≈ 4 KB;
  - `tokens.json`, +0.7 KB generado;
  - bloqueos, ≈ 3.5 KB;
  - ayudante de iconos en `/tmp`, 1 KB;
  - este reporte, ≈ 9 KB.
- **Acciones de herramienta:** ≈ 35.
