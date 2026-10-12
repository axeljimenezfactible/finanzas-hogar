# revisor-shell-r2 · Pista 3, fase 2 (lienzo) · área shell · vuelta 2 de 3

Fecha: 12-oct-2026. Revisor independiente: no corrige. El informe P3-07 está en `docs/estado/wfp3/P3-iter-2-shell.md`.

## Hecho
- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/shell/` y el fragmento `fragmentos/shell.json`. Usé como referencia P3-08 y P3-12, U5-04/06/07, DESIGN.md, U1 (§1.2, §2.1, §2.2, §2.4, §3.5, F-1 y F-6) y `referencia-tipo/format.md`.
- Comparé con `diff` contra la copia de la vuelta 1 (`docs/design/salida/iter-1/shell/`): hay 8 líneas cambiadas en 7 archivos, todas de la corrección.
- Comprobé los 2 medios de la vuelta 1: los dos están **cerrados**.
  - S1-01: los campos móviles miden 44 px o más.
  - S1-02: los 4 enlaces de Ajustes en escritorio van a `ajustes-p26-s01-cuenta-roles-e`.
- Las 7 bajas de la vuelta 1 (S1-03 a S1-09) siguen abiertas.
- Hay 1 baja nueva, S2-01. En P-00 m, el campo mide 60 px y su caja 62 px. La corrección supuso que había un reinicio a border-box, pero no existe.
- **Veredicto: APROBADO** (0 altos, 0 medios, 8 bajos).

## Verificación (scripts propios, `python3 -I`, en `datos/tmp/rev-shell-r2/`, borrados al final)
- **(a) Forma:** los 8 tienen las líneas obligatorias exactas, el JSON de `data-props` parsea y `$preview` coincide con el fragmento. `html.parser` da 0 errores. `lang` es es-MX, cada uno tiene su título y no se repite. Los ids y sus referencias ARIA y `for` son válidos. Los 90 stems del lienzo son únicos y tienen nombre válido. El fragmento lista exactamente los 8 archivos.
- **(b) Tokens:** 0 hex, 0 `rgb()` y 0 `hsl()` fuera de `data-props`. Todos los `var(--x)` existen (`--font-sans` está en `type.families.sans` y en tokens.css).
- **(c) Contraste:** cascada mínima sobre el `<helmet>` y los estilos en línea, en reposo y en hover, con los hex de tokens.json en claro y oscuro. Salen 65 pares, 0 bajo umbral. Quedan exentos los divisores, que DESIGN §7.0 declara decorativos.
- **(d) P3-12:** el grep canon da 0 y el de credenciales da 0. Ninguna fecha del shell está entre las 15 coincidencias de fechas. Revisé 10 valores contra `datos-sinteticos.json` y coinciden (08:10 sigue siendo derivado: S1-07).
- **(e) Fidelidad:** comparé con U1 pantalla por pantalla. En escritorio están todos los estados de §2.2 y §2.4.
- **(f) Accesibilidad:**
  - Controles reales con nombre accesible.
  - 76 controles móviles, 0 bajo 44 px.
  - Sin emoji, sin barra de estado ni teclado falsos.
  - «Error» sale como falso positivo de la lista de U5-07.
  - Dinero bien formado y tabular.
- **(g) Anchos:** las raíces móviles son de 412 px con border-box. Las de escritorio son fluidas, con `flex-wrap` y `999 1 560px`, y llevan `expand:"fill"` en el fragmento.
- **(h) Enlaces:** 116 `href`, de los que 24 son «#». Ninguno está roto ni apunta a otro ancho.
- **Prueba del verificador:** lo corrí sobre copias alteradas y detectó todo lo inyectado: `$preview` distinto, regla `body`, etiqueta desparejada, hex, `var` inexistente, `883.75`, enlace a otro ancho y campo de 40 px.

## Supuestos
- Las alturas se miden por estilo declarado, sin navegador (K45). Para ello tomé los valores por omisión de Chromium y Firefox:
  - `<input type="text">` y `"password"` usan content-box.
  - `type="search"` usa border-box.
  - Sin estilos declarados, las alturas son las del interlineado de la clase de texto en tokens.css más el relleno.
- Igual que en la vuelta 1, la palabra «Error» en las tarjetas «Estados» es español y es nota del lienzo, no texto de la app.
- Igual que en la vuelta 1, `--font-sans` es un token válido.

## Bloqueos
Ninguno.

## Pendientes
- **Para el autor (opcional, no bloquea):** S2-01 y S1-03 a S1-09, todas bajas.
- **Para el orquestador:**
  - El grep de P3-12 busca en `salida/iter-*/project`, pero las copias están en `salida/iter-N/<area>/`, así que no las recorre.
  - Falta copiar los informes de revisión en `salida/iter-N/` (P3-06).
  - En `datos-sinteticos.json`, Axel figura como «Titular» en `personas` y como «Co-titular» en `ajustes.usuarios`.

## Consumo aproximado
~30 acciones de herramienta, ~140 KB leídos (artboards, fuentes por secciones y la revisión anterior) y ~33 KB escritos (los 2 reportes y los scripts temporales, ya borrados).
