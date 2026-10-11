# Reporte P1 · revisor-criterios-r3 (11-oct-2026)

Revisor independiente, ronda 3 de 3, de los criterios 1.1 (`docs/criterios/`). No escribí ni corregí criterios. Cerrado por tope: no.

## Hecho
- Evalué los 18 archivos y `INDICE.md` solo contra `INDICE.md` §1, `decisiones.md` §K y R2-07 de `docs/estado/wf0/REVISION-r3.md`.
- Reejecuté por mi cuenta: conteo y estructura (script propio), mapa K → archivo y línea (script propio), citas de §5 fragmento por fragmento (script propio), cifras de S2-07, S2-08 y S2-09 desde `oraculo/model.json` (script propio con Decimal y ROUND_HALF_UP), días de la semana de K8 y K22, comparación de IDs y severidades contra `HEAD`, y `docs/criterios/tools/verif_final.sh` completo.
- Veredicto: **APROBADO** (0 alta, 0 media, 4 baja). R2-M01 cerrado. Informe: `docs/estado/p1/REVISION-criterios-r3.md`.

## Archivos
- Creados: `docs/estado/p1/REVISION-criterios-r3.md`, `docs/estado/p1/revisor-criterios-r3.md` (este).
- Temporales: `datos/tmp/rv3/` (`cuenta_rv3.py`, `kmap.py`, `cifras_rv3.py`, `citas_rv3.py`, `kfrases.py`), borrados al terminar.

## Verificación
- `python3 -I datos/tmp/rv3/cuenta_rv3.py`: 18 archivos con §1–§5; 257 criterios (197 alta / 60 media / 0 baja); 0 sin Verifica, 0 sin Evidencia, 0 fuera de §5, 0 duplicados; versión 1.1 y lista de cambios en los 18.
- Comparación con `git show HEAD:`: 248 → 257; agregados C1-16…C1-21, U1-15, U5-13 (alta) y U1-16 (media); 0 quitados; 0 severidades cambiadas.
- `bash docs/criterios/tools/verif_final.sh`: cuenta S/WF4 134 (105/29); traza y refs 257, 0 faltantes; citas 276 verificadas, 5 «sin fuente» (C1, C2 ×2, U2, U5); canon 178 comprobadas, 0 diferencias; S1-09 (a) 18 filas, K8 0 fuera de límites, ventana 10→9 rechazada; S1-09 (b) 12 de 12; `/tmp` en S/WF4: 0.
- `python3 -I datos/tmp/rv3/citas_rv3.py`: en las 5 citas señaladas todos los fragmentos existen textualmente; solo cambia el orden (B-04).
- `python3 -I datos/tmp/rv3/cifras_rv3.py`: deuda viva 2,921.25 / 973.75 / 3,776.06 / 1,258.70 / 3,776.06 / 1,258.70 / 6,293.44 con saldo no devengado 1,947.50 / 0.00 / 2,517.37 / 0.00 / 2,517.37 / 0.00 / 5,034.74 (7 de 7 iguales a S2-07); diarios 31.41, 64.92, 89.91 y semanales 219.88/219.87, 629.35/629.37 (iguales a S2-08); pagos de la despensa 4 iguales a S2-09; 2026-10-12 lunes, 14 miércoles, 18 domingo, 31 sábado; 10-oct + 30 días = 9-nov.
- `grep -nP '(?<!datos)/tmp\b'`: 0 rutas operativas (11 menciones descriptivas en listas de cambios).
- `grep` de valores viejos (gasto 25°, `go 1.26`, `SameSite=Lax`, 374.93 como tope, lista de compras local o PROVISIONAL, PROVISIONAL sobre K5/K6): 0 en criterios (solo la nota 8 de `INDICE.md` §5, B-01).
- Existencia de los comandos citados: `backend/internal/api/c1gen/build.py --check` y `verify_c1.py --db` existen.

## Supuestos
- Leí «paráfrasis entre comillas» como parte de R2-07 porque `REVISION-r3.md` (líneas 115-116) la atribuye a ese hallazgo.
- Las diferencias de texto que el integrador reconcilia en `INDICE.md` (§1, §5, §6) las registro como baja, según la orden (las filas U/C/P3 ya coinciden con mi conteo).

## Bloqueos
Ninguno nuevo.

## Pendientes
Ninguno dentro de mi alcance. Las correcciones baja quedan en el informe (B-01 a B-04) para el integrador o el agente que toque esos archivos.

## Solicitudes fuera de alcance
- `docs/criterios/INDICE.md` (integrador): 257 en §1 línea 7; notas 8-10 de §5 como resueltas por K26, K27 y K30; filas de `criterios-promedios-redondeo` y `criterios-frescura-extremos` en §6 y RESUELTO en lo cerrado por §K (B-01).
- `docs/criterios/C1.md:41` y `:45`: «día de pago 15 y 30» en lugar de «quincena 15/30» (B-02).
- `docs/criterios/U2.md:100`: citar decisiones I «o U2 propone 135°» (B-03).
- `docs/criterios/C1.md`, `C2.md`, `U2.md:96`, `U5.md` §5: reordenar los fragmentos de 5 citas (B-04).

## Paquetes instalados
Ninguno. `verif_final.sh` usa `uv run --no-project --with pyyaml` (caché de uv, sin paquetes del sistema).

## Consumo aproximado
~48 acciones de herramienta; ~185 KB leídos (protocolo, §K, §1-§2 y fragmentos del INDICE, 2 reportes, fragmentos de criterios por `grep`/`awk` y salidas de scripts); ~20 KB escritos (scripts temporales, informe y este reporte).
