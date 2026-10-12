# correccion-U2-r2 · Sistema de diseño (DESIGN.md + C2) · WF-1 · 11-oct-2026

Corrige el único hallazgo alto o medio de `docs/estado/wf1/REVISION-U2-r2.md` (N-1, media). No se tocan los bajos B-1 a B-7, que no estaban en la tarea, ni ninguna decisión de §A–§K. Ninguna cifra ni ningún token cambian.

## Hecho
- **N-1 (media) · U2-12 / C2-04**: en `docs/design/DESIGN.md` §7.10, línea 893 (estado Activo de la barra inferior), cambié «píldora de 32 × 56 px» por «píldora de 56 × 32 px (`size.bottombar.indicator-width` × `size.bottombar.indicator-height`)». Así sigue la convención ancho × alto del resto del documento (asa 32 × 4, línea 923; pista 36 × 20, línea 675; Android 412 × 915, línea 1117) y coincide con C2-tokens.json (`indicator-width` = 56px, `indicator-height` = 32px) y con §10.6, línea 1372. La píldora mide 32 px de alto, menos que la barra de 56 px, y 56 px de ancho, menos que el ancho mínimo de 64 px del destino.
- Volví a correr toda la verificación del artefacto. Los 5 archivos generados quedaron byte a byte iguales a los anteriores.

## Archivos
- Modificado: `/home/deepsec/finanzas-hogar/docs/design/DESIGN.md` (solo la línea 893).
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-U2-r2.md` (este reporte).
- `npm run check` regeneró sin cambios `docs/contratos/C2-tokens.json`, `docs/contratos/tokens.css` y `docs/contratos/tailwind.theme.css`, con el mismo sha256 antes y después.

## Verificación
- **N-1**:
  - `grep -c '32 × 56' docs/design/DESIGN.md` da 0.
  - La línea 893 dice «56 × 32 px (`size.bottombar.indicator-width` × `size.bottombar.indicator-height`)».
  - Hice un cruce con `python3 -I` entre la expresión de §7.10 y C2-tokens.json. Resultado: ('56', '32', `size.bottombar.indicator-width`, `size.bottombar.indicator-height`), «coincide con C2: True». Ancho 56 ≤ `item.min-width` 64: True. Alto 32 < `layout.bottombar.height` 56: True.
- **U2-01 / U2-10 / C2-04**: `python3 -I docs/design/fuentes/scripts/verificar.py .` termina con código 0 y **RESULTADO: OK**.
  - 11 encabezados `##` en orden y 0 referencias a terceros.
  - Sprite: 1,877 símbolos. Iconos citados en §5.4: 55, y no falta ninguno.
  - 162 hex citados, 0 fuera de C2.
  - Rutas de token citadas: 958, y 0 no existen. Son 2 más que antes porque la línea 893 ahora cita los dos tokens.
  - 158 pares token→hex, 0 diferencias. Los 19 contrastes de §2.8 coinciden con el recálculo.
- **C2-01 a C2-10**: `cd docs/design/tools/tokens && npm run check` termina con código 0.
  - generate: escribe 100073 B, idéntico al anterior.
  - build: Style Dictionary sin advertencias.
  - verify: **RESULTADO: OK**.
    - C2-02: 445 tokens; 0 sin tipo, 0 sin `$value` y 0 referencias rotas.
    - C2-03: 76 semánticos con claro y oscuro, 0 incompletos.
    - C2-04: 190 colores con 0 diferencias de hex; 158 pares de DESIGN con 0 diferencias; 242 tokens citados, 0 inexistentes; 0 grupos citados como token.
    - C2-05: 170 combinaciones, 0 bajo umbral.
    - C2-06: `--ui-*` 78/78/78, prefers-color-scheme 1, data-theme 2 y prefers-reduced-motion 1.
    - C2-08: 0 URL a terceros; las fuentes y la licencia existen.
    - C2-10: 0 hex repetidos sin alias, 0 semánticos literales y 18 grupos de medida informativos (sin cambio).
  - check:sync: tokens.css `2b9429bc…708d` y tailwind.theme.css `f66eb866…c88a`, iguales a los versionados.
  - test:tailwind: tailwindcss v4.3.3, «Done».
- **U2-04 / C2-05**: `python3 -I docs/design/tools/contraste_componentes.py --check` termina con código 0.
  - «CSV en disco IDENTICO al regenerado».
  - 158 pares y 316 combinaciones: 304 con umbral y 12 decorativas exentas. **0 bajo umbral**.
  - Cruce con contraste.csv: igual 266, nuevo 50 y difiere 0. Cruces secundarios: 0.
- `--independiente`: 10 razones con segundo cálculo, todas coinciden, código 0.
- Ningún archivo generado cambió: `sha256sum -c` da OK en C2-tokens.json (`570ff553…`), tokens.css (`2b9429bc…`), tailwind.theme.css (`f66eb866…`), contraste-componentes.csv (`df8cf761…`) y contraste-componentes.md (`997a278b…`).
- `python3 -I -m py_compile` de los 3 scripts del alcance: OK.
- **U2-12**: N-1 está cerrado. Los 12 de 12 grupos de §7.16 no cambian.
- Borré el temporal `datos/tmp/c2-prueba.css`, que crea `test:tailwind`. `datos/tmp/c2-sync` no quedó. Las carpetas `datos/tmp/rev-s2-r3` y `datos/tmp/revU3r2` son de otros agentes y no las toqué.

## Supuestos
- (ninguno) La corrección es la que sugirió el revisor, palabra por palabra. Ninguna cifra ni ningún token cambian.

## Bloqueos
- Ninguno nuevo. El único bloqueo abierto del proyecto, `docs/estado/bloqueos/S1-parametros-avisos.md`, no afecta a U2.

## Pendientes
- Quedan abiertos los bajos de REVISION-U2-r2, que están fuera de esta tarea porque solo se pidieron los altos y medios:
  - B-1: §1.2 P5 cita §3.4, y debe ser §3.3.
  - B-2: §2.7 debe aclarar «mínimo en gris y en las tres simulaciones».
  - B-3: §2.2 dice ±30°; los ángulos reales son −25° y +35°.
  - B-4: §2.5, línea 116: decir «relleno 5.00:1 sobre tarjeta; texto blanco 5.59:1».
  - B-5: `verificar_c2.py` debe comparar los valores no de color de las tablas de DESIGN.
  - B-6: §7.2–§7.15 deben citar `ui.*` o dar una tabla de equivalencias en §7.0.
  - B-7: quitar «en construcción» del encabezado.

  Todos están en `docs/design/DESIGN.md`, salvo B-5, que está en `docs/design/tools/tokens/verificar_c2.py`. El revisor recomienda corregir B-1 a B-4 y B-7 antes del paquete de la Pista 3.

## Solicitudes fuera de alcance
- Ninguna.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna. Esta corrección no toca valores de K31.

## Consumo aproximado
- Acciones de herramienta: 15.
- Leído: unos 85 KB (protocolo, §K de decisiones, criterios U2, REVISION-U2-r2, extractos de DESIGN.md, C2-tokens.json, correccion-U2-r1 y salidas de verificación).
- Escrito: unos 5 KB (1 línea de DESIGN.md y este reporte). Los archivos que regeneró `npm run check` salieron idénticos.
