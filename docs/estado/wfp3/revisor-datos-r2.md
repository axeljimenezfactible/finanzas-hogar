# revisor-datos-r2 · Reporte del agente

Fecha: 11-oct-2026, noche. Rol: revisor independiente, ronda 2 de 3, de los datos sintéticos de la Pista 3. No corrigió nada. Informe completo: `docs/estado/wfp3/REVISION-datos-r2.md`.

## Hecho
- Repetí el verificador del autor (51 OK) y la regeneración determinista (salida a `datos/tmp/revdatos-r2/`, idéntica con `cmp`).
- Repetí con comandos propios el grep de P3-12, los patrones de credencial y el grep de P3-08 sobre `docs/design/maquetas/`: 0 en los tres.
- Corrí 3 scripts propios en `datos/tmp/revdatos-r2/`:
  - `chk.py`: 26 sumas y porcentajes.
  - `cmp.py`: montos, fechas, productos y marcas contra `data.json` y `model.json`.
  - `cov.py`: campos por necesidad de U1.
- Contrasté las 6 tarjetas con los esquemas de C1 (YAML parseado), U1 §2.1, §2.3 y §3, U3 §2, U4 P-26.05/.09/.11 y K50.
- Comprobé los 4 hallazgos medios de r1: los 4 están **cerrados**.
- Hallazgos nuevos: **M-1** (resumen diario del 20-oct: «3 gastos» contra 2 en `gastos`), **M-2** (cobertura: P-08, P-02/P-06 y P-15 en Tendencias, P-12 N-141 y P-26.11 número de registros) y B-1 a B-5.
- Veredicto: **no aprobado** (0 altos, 2 medios, 5 bajos).

## Verificación (comandos y salidas)
```
python3 -I docs/design/tools/ds/verificar_sinteticos.py                      -> 51 OK, RESULTADO: OK
python3 -I docs/design/tools/ds/generar_sinteticos.py datos/tmp/revdatos-r2/regen.json; cmp -> idéntica
grep -rnoE "<21 cifras P3-12, con y sin coma>" docs/design/maquetas/ | wc -l  -> 0
grep -rn -i -E "<patrones de credencial P3-12>" docs/design/maquetas/ | wc -l -> 0
grep -rnoE '\$[0-9]{4,}\.[0-9]{2}' docs/design/maquetas/ | wc -l              -> 0
python3 -I -c '<comando de fechas de P3-12>' model.json data.json datos-sinteticos.json
   -> 15 fechas (día 15, 30 y 28-feb), las mismas que lista el LEEME
python3 -I datos/tmp/revdatos-r2/chk.py  -> 25 OK + S9 (falso positivo por redondeo al par; con ROUND_HALF_UP, OK)
python3 -I datos/tmp/revdatos-r2/cmp.py  -> 0 cifras canon en *_cents; 14 coincidencias de valores redondos entre
   conceptos distintos; 0 nombres de producto idénticos; 5 marcas en común con precio y presentación distintos;
   fecha coincidente con monto de model.json: solo $9,000.00 (B-3)
grep -n "Resumen del" docs/design/tools/ds/generar_sinteticos.py -> línea 540: «en 3 gastos» fijo en el código (M-1)
```

## Supuestos
- Si un campo sin cifras falta pero se puede dibujar con «—» o con el texto fijo de U1 (auditoría de P-04, estado diario de P-10, «Este dispositivo»), lo califiqué **baja**. Si obliga al autor a inventar una serie o un conteo (K45.1), lo califiqué **media** (cobertura incompleta).
- Según el plan, se dibuja 1 artboard por pantalla y ancho (`pista3-plan.md` §2.2). Eso puede dejar sin dibujar las Tendencias que no son la vista por defecto. Aun así, M-2 se queda en media: N-141 está en la página principal de P-12 y el LEEME no dice qué hacer con esas vistas.
- El $9,000.00 de B-3 lo califiqué baja: el par fecha-monto no es el real y el valor es redondo. Ya estaba como B-4 de r1.

## Bloqueos
Ninguno.

## Pendientes
- Lo corrige el agente de corrección: M-1 (generador, línea 540, y una comprobación nueva) y M-2 (series y campos, o una regla en el LEEME). Las bajas quedan a su criterio.
- La ronda 3 debe repetir el grep de P3-12 sobre `docs/design/salida/iter-*/project` cuando exista.

## Consumo aproximado
- ~44 acciones de herramienta.
- ~190 KB leídos: JSON por script, U1 §2 y §3 por partes, U3 §2, U4 por grep, criterios, LEEME y reportes.
- ~19 KB escritos: este reporte, el informe y 3 scripts temporales, ya borrados.
- `cerrado_por_tope` = false.
