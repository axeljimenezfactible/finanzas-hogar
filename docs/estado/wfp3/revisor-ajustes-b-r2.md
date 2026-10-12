# revisor-ajustes-b-r2 (vuelta 2 de 3)

Fecha: 12-oct-2026. Revisor independiente de ajustes-b (P-26.08 a P-26.14). No corrige. Informe completo: `docs/estado/wfp3/P3-iter-2-ajustes-b.md`.

## Hecho
- Revisé los 14 artboards de `docs/design/maquetas/lienzo/project/ajustes/` (s08 a s14, `-m` y `-e`) y `fragmentos/ajustes-b.json` contra P3-07, P3-08 y P3-12, U5-04, U5-06 y U5-07, DESIGN.md, U1 §1.3, §2.1, §2.2 y §2.3, U4 §0 y P-26.08 a P-26.14, y `referencia-tipo/format.md`.
- Comprobé los 12 hallazgos de la vuelta 1: 11 cerrados y 1 reabierto en parte (H-11, la mitad K7, que depende de `datos-sinteticos.json`, sin cambios desde el 11-oct 23:55).
- Hallazgos nuevos: R2-A1 (alta), contraste de 3.72:1 en claro por `opacity:.85` en s13-m y s13-e (ya estaba en iter-1); R2-M1 (media), la antigua B-02, que sube de severidad por P3-08 punto 5; y 9 bajas (B-11 a B-19).
- **Veredicto: NO APROBADO** (1 alta, 2 medias, 15 bajas).

## Verificación
Scripts propios con `python3 -I` en `datos/tmp/rev-ajb-r2/`, borrados al final:
- Forma (`forma.py`): 14/14 con las líneas obligatorias; `data-props` válido y `$preview` igual al fragmento; `html.parser` sin etiquetas sin cerrar, sin cierres huérfanos y sin ids repetidos; 90 stems únicos en el lienzo; fragmento con `boards` = `order` = los 14 archivos y `expand:"fill"` solo en escritorio.
- Tokens: 0 hex, `rgb()` o `hsl()` fuera de `data-props`; todos los `var(--x)` existen en `tokens.json` y en `tokens.css`.
- Contraste (`cascade.py` y `contraste.py`): cascada CSS propia con herencia, fondo efectivo, `data-theme` local y opacidad. Revisé 1,620 nodos por tema y 22 pares: 0 pares de color bajo 4.5:1, 1 texto con opacidad a 3.72:1 y 0 pares de UI bajo 3:1 (perilla 5.77 / 10.47 y 3.71 / 7.96).
- Táctil (`tactil.py`): 142 objetivos en móvil, 0 bajo 44 px.
- P3-12 y datos (`forma.py`, `datos.py`): 0 cifras canon (21 patrones, con y sin coma, y en centavos) y 0 credenciales. Coincidencias con el calendario real: solo 30-oct-2026 y 28-feb-2027, sin monto. Los 65 montos y 46 fechas de texto salen del JSON o se derivan de él con aritmética comprobada. Las fechas de nacimiento no están en el JSON y tampoco coinciden con datos reales.
- Accesibilidad (`a11y.py`, `tabular.py`): nombres accesibles y labels completos; 30 interruptores con `aria-checked` y nombre; 7 montos sin `tabular-nums`; 0 emoji, 0 ISO visibles, 0 `<code>`, 0 `type="date"`.
- Enlaces (`links.py`): 242 `href`, 30 a `#`, 0 rotos y 0 de otro ancho.

## Supuestos
- Los botones de `.seg` sin borde propio no se cuentan como falla de UI, porque el contenedor lleva `ui-border-control` (3.47 / 6.09).
- B-02 se reclasifica a media (R2-M1) porque P3-08 punto 5 pone las cifras tabulares entre las decisiones. La vuelta 1 la había dejado en baja.
- B-11 (rol de Axel) y B-16 (formato de avance y estados de metas) quedan en baja porque los documentos se contradicen entre sí (U1 frente a U4, DESIGN frente a U4): decide el orquestador.
- Las alturas del fragmento no se volvieron a medir: este revisor no usa navegador.

## Bloqueos
- H-11 (K7) no puede cerrar sin que correccion-datos cambie `inventario.articulos` y `resumen` en `datos-sinteticos.json`. Si en la vuelta 3 no se corrige, el orquestador debe registrarlo en `docs/estado/bloqueos/`.

## Pendientes
- Para el autor (vuelta 3): R2-A1, R2-M1 y las bajas baratas B-03, B-05, B-07, B-12, B-13 y B-14.
- Para correccion-datos u orquestador: H-11 (K7), B-11, B-15, B-16 y B-19.
- Copia de la vuelta 2 en `docs/design/salida/iter-2/` (P3-06): la hace el orquestador.

## Consumo aproximado
~185 KB leídos (format.md completo; U4 §0 y P-26.08 a P-26.14; DESIGN §3.3 y §7.9 a §7.12; U1 §2.2 y §2.3 por filas; informes de la vuelta 1; los artboards a través de scripts), ~24 KB escritos (este reporte y el informe P3-iter-2), unas 55 acciones. 0 acciones de navegador, 0 Artifact, 0 red, 0 git, 0 Telegram.
