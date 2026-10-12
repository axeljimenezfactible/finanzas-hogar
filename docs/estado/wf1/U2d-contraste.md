# U2d-contraste · reporte (11-oct-2026)

## Hecho
- Script `docs/design/tools/contraste_componentes.py`: lee la lista de pares de DESIGN.md §10.5 y §7.1, los hex solo de `docs/contratos/C2-colores.json`, calcula WCAG 2.x con `Decimal` (60 dígitos) sin importar `colorimetria.py`; modos: escribir, `--check`, `--independiente`.
- `docs/design/contraste-componentes.csv`: 158 pares × 2 temas = 316 filas (304 con umbral, 12 decorativas exentas).
- `docs/design/contraste-componentes.md`: resumen, 15 pares más justos, hallazgos de fragilidad F-1 a F-3, decorativos, 10 razones con segundo cálculo.
- Resultado: 0 pares bajo umbral; ningún token cambiado; ningún bloqueo nuevo.

## Archivos
- docs/design/contraste-componentes.csv (nuevo)
- docs/design/tools/contraste_componentes.py (nuevo)
- docs/design/contraste-componentes.md (nuevo)
- docs/estado/wf1/U2d-contraste.md (este)

## Verificación
- `python3 -I docs/design/tools/contraste_componentes.py` → «pares: 158 · combinaciones: 316 · con umbral: 304 · decorativas exentas: 12 · bajo umbral: 0», cruce con contraste.csv igual 266 / nuevo 50 / difiere 0; incoherencias §7.1/§10.5/C2: 0; exit 0.
- Regeneración idéntica: dos corridas seguidas dan el mismo sha256 (`df8cf761…6d2f`); `--check` imprime «CSV en disco IDENTICO», exit 0.
- `--independiente`: 10 razones recalculadas con otra fórmula (flotante, umbral 0.03928, producto punto) coinciden; diferencia máxima +0.0095 (el CSV trunca); exit 0.
- U2-04 (parte de componentes): 316 combinaciones, 0 bajo umbral, salida 0. U2-07: C-001 claro 5.77 y oscuro 10.47; C-002 claro 5.53 y oscuro 10.37; blanco sobre marca 600 no aparece en ninguna fila de texto. U5-04: el CSV queda como tabla de referencia para el DOM.
- Las 6 combinaciones sin tipo de umbral son exactamente C-087, C-105, C-116, C-143, C-144, C-158 (12 combinaciones), como declara §10.5.

## Supuestos
- La lista de pares es la de §10.5 tal cual (sin añadir pares de §2.9 o §7 que no estén ahí); si U2b agrega componentes, hay que actualizar §10.5 y volver a correr el script (el conteo esperado 158/316/304/12 está en el script y sale con código 1 si cambia).
- Los decorativos se calculan pero no cuentan; «texto grande» solo se usa para la cifra de 24 px (C-096), como dice §10.5.
- `cumple` se decide con la razón exacta; `razon` se imprime truncada a 2 decimales como `contraste.csv`.

## Bloqueos
Ninguno nuevo. Bloqueo abierto único (`docs/estado/bloqueos/S1-parametros-avisos.md`) no afecta a este trabajo.

## Pendientes
- Ninguno de mi alcance. Opcional para U2/U5: decidir si F-1 (C-138 claro 4.51) merece más holgura; U5-04 puede usar el CSV como referencia al medir el DOM.

## Solicitudes fuera de alcance
- `docs/design/DESIGN.md` §10.5: cuando se cierre, agregar la ruta de `contraste-componentes.csv` y `contraste_componentes.py` como la tabla oficial (el texto dice «U2d la calcula» sin ruta). Motivo: trazabilidad U2-04.
- `docs/design/DESIGN.md` §10.5 (hallazgo F-1, solo si U2 quiere holgura): C-138 claro 4.51 contra 4.5; corrección posible en el menú «/» (alias y atajo en la fila resaltada con `ui.text.default`).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (no hay valores K31 en este alcance).

## Consumo aproximado
- Acciones de herramienta: 16.
- Leído: ≈ 115 KB (protocolo 6, sección K ≈ 9, U2.md 20, DESIGN.md §10.5 y §7.1 ≈ 45, C2 extractos y CSV ≈ 6, salidas de verificación ≈ 25).
- Escrito: ≈ 36 KB (script 14, CSV 56 KB generado por el script, md 9, reporte 3; escrito por herramienta ≈ 26 KB, el CSV lo genera el script).
