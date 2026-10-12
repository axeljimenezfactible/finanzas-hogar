# revisor-ajustes-a-r1 · reporte

Fecha: 12-oct-2026. Pista 3, fase 2 (lienzo). Revisión independiente, vuelta 1 de 3, del área `ajustes`, parte a: P-26 índice y secciones 01 a 07, 16 artboards y el fragmento `ajustes-a.json`. No corrige nada.

## Hecho
- Informe de revisión (P3-07): `docs/estado/wfp3/P3-iter-1-ajustes-a.md`, con la tabla de verificación, las líneas P3-07, P3-08, P3-12, U5-04, U5-06 y U5-07, la fidelidad por pantalla y, por artboard, las reglas falladas y la corrección para la vuelta siguiente.
- **Veredicto: no aprobado (0 altos, 8 medios, 7 bajos).**
  - Medios: M-01 enlace a `s08-reportes-frecuencias` roto (el archivo real es `s08-reportes`); M-02 `#instrucciones` sin destino; M-03 cinco objetivos táctiles de menos de 44 px en móvil; M-04 falta el «Error al guardar» de U1 §2.2 y U4 §0.3, s01 tiene un vacío que contradice U4 y faltan 4 vacíos; M-05 no hay ningún error de validación de campo dibujado; M-06 fechas con cero a la izquierda (17); M-07 la barra lateral no tiene colapso `panel-left` ni la anatomía de §7.9; M-08 faltan el calendario por defecto de promociones, «Eliminar promoción», «Quitar monto», el interruptor «Activa» y los bloques plegables de s05-m.
  - Bajos: título del índice, concordancia «No se pudieron cargar», leyenda «Enter…» fuera de lugar, asa decorativa, alturas estimadas, acción primaria no única, relleno crítico fuera del diálogo.

## Verificación
Hecha con scripts propios (`python3 -I`) en `datos/tmp/rev-aj-a/`, borrados al terminar:
- Forma: 16/16 tienen las líneas obligatorias exactas. `data-props` es JSON válido y `$preview` coincide con el fragmento. `html.parser` no encuentra etiquetas sin cerrar ni ids repetidos. Nombres y stems son válidos y únicos, y el fragmento lista exactamente los 16 archivos (`expand: "fill"` solo en `-e`).
- Tokens: 0 hex fuera de `data-props`; los 68 `var(--x)` existen en `tokens.json`.
- Contraste: 69 pares × 2 temas = 138 combinaciones, 0 bajo umbral. El mínimo en texto es 5.02; en UI, 3.46 en claro y 4.05 en oscuro.
- P3-12: 0 cifras canon (con y sin coma y en centavos) y 0 credenciales. 55 de 59 montos y 24 de 25 fechas salen literales de `datos-sinteticos.json`; el resto son sumas o ejemplos que el autor declara. La única fecha que coincide con el oráculo, 30-oct-2026, va con un monto sintético.
- Accesibilidad estática: 0 emoji, 0 negativos con guion, 0 fechas ISO, 0 SVG sin `aria-hidden`, 0 botones sin nombre, 0 campos sin etiqueta, 0 controles dentro de `<a>`.
- Fidelidad: comparé el texto de cada `-e` y la estructura de los `-m` con U1 §2.1, §2.2 y §2.3 y con U4 P-26.01 a P-26.07, incluido §0.

## Supuestos
1. El contraste es estático: una cascada propia del CSS del `<helmet>` con los tokens de los dos temas, sin navegador (K45). La medición con estilos computados la hace U5-04.
2. `--font-sans` se da por token de `tokens.json` (`type.families.sans`, según la convención de `meta.note`).
3. Los ejemplos tecleados y las sumas que el autor declara en su reporte (supuestos 2 a 7) cuentan como derivados del JSON y no como cifras inventadas.
4. Fechas: por U4 («9-nov-2026», «8-ene-2027») y por las 50 fechas del resto del lienzo, el formato es el día sin cero a la izquierda.
5. Las claves `groceries` y `eating_out` visibles son un choque entre U4 P-26.04 y U5-07, no un error del autor: queda señalado para el orquestador.

## Bloqueos
Ninguno.

## Pendientes
- Vuelta 2: el autor corrige M-01 a M-08 (y los bajos que pueda) y se vuelve a revisar.
- Copiar el informe a `docs/design/salida/iter-1/` (P3-06): le toca al orquestador.
- Medir las alturas reales (`h`) al publicar el lienzo (B-05).

## Consumo aproximado
- Leído: unos 200 KB. Criterios P3 y U5 por secciones; U4 §0 y P-26.01 a P-26.07, resumidos por script; U1 §2.2 y §2.3, solo las filas de P-26; DESIGN §7.9, §9.3, §9.4 y el encabezado de §10.5; format.md entero; el bloque `ajustes` y la estructura de `datos-sinteticos.json`; esquemas de texto de los 16 artboards generados por script; reporte y fragmento del autor.
- Escrito: unos 17 KB (informe y este reporte), más unos 20 KB de scripts temporales ya borrados.
- Unas 40 acciones de herramienta.
