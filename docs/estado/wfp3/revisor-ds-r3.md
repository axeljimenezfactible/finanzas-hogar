# Reporte · revisor-ds-r3 (Pista 3, fase 1, ronda 3 de 3)

Fecha: 11-oct-2026, noche. Informe completo: `docs/estado/wfp3/REVISION-ds-r3.md`.

## Hecho
- Volví a verificar el sistema de diseño de `docs/design/sistema/project/` contra DESIGN.md 1.0, C2-tokens.json 1.0, P3-claude-design.md 1.3, format.md, cover.md y craft.md, con scripts propios en `datos/tmp/r3/`, que borré al terminar.
- Comprobé los hallazgos de la r2. A-01 (iconos de 24 px en la barra superior) está cerrado. M-01 (BarraVista y Galeria) está cerrado, y el calendario genérico queda como pendiente declarado.
- Dictamen: **APROBADO**, con 0 altos, 0 medios y 5 bajos (B-01 a B-05).

## Verificación (comandos y salidas)
- `python3 -I docs/design/tools/ds/verificar_ds.py` → `RESULTADO: OK · 445 tokens C2 -> 355 en listas de tokens.json + 6 estilos de tipo (18 font.*) + 2 fuentes; 0 fallas`.
- `v_tokens.py` (propio) → C2 445; 355 nombres, 0 duplicados, 0 inválidos; `var(` 0, `color-mix(` 0; primer tema light. C2 → json: 0 diferencias en 340 colores. json → tokens.css: 530 iguales más 2 equivalentes (velo con alfa 0x66/0x99). Muestra de 20 al azar igual en claro y oscuro. Fuentes con sha256 idéntico a `docs/design/fuentes/inter/web/`.
- `v_prev.py` (propio, motor de cascada) → 21 vistas previas con marcador válido, 0 errores de HTML, 0 `<script>`, 0 red, 0 etiquetas prohibidas, 0 `var` sin token. Contraste: 3,818 evaluaciones. Solo fallan pares exentos y la excepción ya declarada de Karen (2.83 sobre `ui-surface-hover`).
- `v_readme.py` (propio) → todos los tokens nombrados en los README existen. Los contrastes que cita el README coinciden con los recalculados (truncados como DESIGN §2.8).
- Grep P3-12 (21 cifras canon) sobre `docs/design/sistema/project` → `0`. Grep de credenciales → `0`.

## Supuestos
- «Solos en barras» (§5.2) se lee con §7.2: `size-icon-md` solo en los botones de la barra superior y de la barra inferior.
- Los pares exentos de 3:1 son el esqueleto, el velo, el anillo de los marcadores, la retícula y la portada, como están listados en el informe.

## Bloqueos
- Ninguno nuevo. Siguen abiertos P3-contraste-chip-hover y P3-huecos-componentes (puntos 1 a 12), los dos en espera de Axel.

## Pendientes
- B-01: referencia cruzada a P3-contraste-chip-hover en BarraVista/README.
- B-02: corregir en el punto 12 el motivo de `ui-text-disabled` (C-093 cumple: 5.40 y 7.96).
- B-03: calendario genérico §7.5.4, según lo que responda Axel.
- B-04: al publicar, confirmar que el `tokens.css` compilado declara `--duration-*`, `--ease-*`, `--font-weight-*`, `--viz-*` y `--breakpoint-*`.
- B-05: los montos de las vistas previas vienen casi todos de fuera de `datos-sinteticos.json`; avisar a los autores del lienzo.

## Consumo aproximado
- Lectura: unos 200 KB (tope 400). Escritura: unos 16 KB, contando los scripts temporales ya borrados (tope 150). Acciones: unas 50 (tope 60). cerrado_por_tope=false.
