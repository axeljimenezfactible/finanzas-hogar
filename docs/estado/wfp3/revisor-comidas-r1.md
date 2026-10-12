# revisor-comidas-r1 · Pista 3, fase 2 (lienzo), vuelta 1 de 3

Fecha: 12-oct-2026 · Agente: `revisor-comidas-r1` (revisor independiente; no corrige) · Área: `comidas`

## Hecho

- Revisé los 6 artboards de `docs/design/maquetas/lienzo/project/comidas/` (P-10, P-11 y P-12, en móvil y escritorio) y el fragmento `docs/design/maquetas/lienzo/fragmentos/comidas.json`. Los comparé con P3 1.3 (P3-04, P3-06, P3-07, P3-08, P3-12), U5 (U5-04, U5-06, U5-07), DESIGN.md, U1, U3 y `referencia-tipo/format.md` y `craft.md`.
- Informe P3-07 con la tabla artboard × regla, la fidelidad a U1, el contraste, la lista de P3-08 y las líneas de evidencia: `docs/estado/wfp3/P3-iter-1-comidas.md`.
- **Veredicto: no aprobado.** 0 altos, 4 medios y 12 bajos.
  - M-1: en móvil, los enlaces a la receta van al artboard de escritorio (5 en p10-m y 12 en p11-m).
  - M-2: el enlace «‹ Página padre» de la barra superior móvil mide unos 20 px (< 44).
  - M-3: falta el selector de recetas con la fila de días de P-10 (U1 §2.1, F-2).
  - M-4: falta la barra de la vista de P-10 en escritorio (Calendario*, Tabla, Tablero, Galería y filtro por persona).

## Verificación

Scripts propios con `python3 -I` en `datos/tmp/rev-comidas/` (un motor de cascada CSS mínimo, el verificador y el volcado del plan; borrados al terminar):
- Forma: líneas obligatorias exactas y en orden; JSON de `data-props` válido, con `$preview` igual al fragmento; `html.parser` sin etiquetas sin cerrar; `lang="es-MX"`; títulos propios; nombres válidos y stems únicos en todo el lienzo. Resultado: 6 de 6.
- Tokens: 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`. Todos los `var(--x)` existen en `tokens.css`; en `tokens.json` también, salvo `font-sans`, que es `type.families.sans`.
- Contraste: 169 pares, cada uno en claro y en oscuro, 0 bajo el umbral. 12 pares contrastados con `contraste.csv` y `contraste-componentes.csv` coinciden al ±0.01.
- P3-12: el grep de las 21 cifras canon da 0, el de centavos 0 y el de credenciales 0. Ninguna fecha de los artboards coincide con las del modelo real. Una muestra de 10 valores y las 28 casillas de la cuadrícula salen de `datos-sinteticos.json` sin diferencias.
- Accesibilidad: 0 controles falsos, 0 botones sin nombre, 0 inputs sin label y 177 SVG con `aria-hidden`. Las alturas se midieron resolviendo `min-height` y `height` con los tokens: solo falla el enlace de la barra superior (M-2). El falso positivo de los enlaces de la galería se descartó porque `.rn::after` cubre toda la tarjeta.
- Enlaces: todos apuntan a la lista del lienzo o a «#». Los 17 que cruzan de ancho son M-1.
- Fragmento: lista exactamente los 6 archivos; `expand: "fill"` solo en escritorio; `is_interactive: false`; sin x ni y.

## Supuestos

1. «Pantalla principal», en la escala de severidad, es el elemento central de la pantalla (la cuadrícula de P-10), que sí está. Por eso el selector de recetas, que es una hoja que se abre con un toque, y la barra de la vista cuentan como «elemento secundario ausente» (media).
2. «Plan», «Total» y «Error» son palabras del español y no cuentan como inglés aunque «Error» esté en la lista de U5-07.
3. `font-sans` cumple «existe en tokens.json» porque es `type.families.sans` con su nombre de CSS, tal como dice `meta.note` de `tokens.json`.
4. Los altos de contenido son estimaciones sin navegador (K45). Por eso B-8 es baja y pide medir en el lienzo.
5. El informe va donde lo pidió el encargo (`docs/estado/wfp3/`), no en `docs/estado/` como dice el texto de P3-07.

## Bloqueos

Ninguno.

## Pendientes

- Vuelta 2: que el autor corrija M-1 a M-4 (y los bajos que quiera) y que otro revisor vuelva a pasar los mismos scripts.
- Orquestador: copiar este informe a `docs/design/salida/iter-1/` (P3-06) y decidir `is_interactive` si el prototipo debe recorrerse en Play.

## Consumo aproximado

- Leído: unos 230 KB (secciones de U1 y DESIGN, P3, U5, format.md entero, partes de craft.md, cuerpos de los 6 artboards sin SVG, partes de `datos-sinteticos.json`, `tokens.json` por script y los CSV de contraste).
- Escrito: unos 30 KB (scripts temporales, ya borrados, este reporte y el informe P3-iter-1-comidas.md).
- Acciones de herramienta: unas 50.
