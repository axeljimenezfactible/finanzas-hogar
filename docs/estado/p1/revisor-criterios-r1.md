# Reporte P1 · revisor-criterios-r1 (11-oct-2026)
Revisor independiente, ronda 1 de 3, de los criterios 1.1 (`docs/criterios/`). No escribí nada de lo revisado y no corregí nada. Informe: `docs/estado/p1/REVISION-criterios-r1.md`. Cerrado por tope: no.

## Hecho
- Veredicto **NO APROBADO**: 0 hallazgos altos, 6 medios (H-01 a H-06) y 7 bajos (H-07 a H-13), cada uno con ruta y corrección sugerida.
- (a) Estructura: 18/18 archivos con §1–§5. 256 criterios (196 alta / 60 media / 0 baja), revisados uno por uno: todos tienen Verifica y Evidencia y su línea en la §5, con ID seguidos.
- (b) Tabla K1–K34 con línea por decisión: 20 aplicadas, 1 que no aplica (K33) y 13 parciales. Las parciales son K2, K3, K4, K8, K9, K12, K13, K16, K18, K21, K23, K25 y K26: les falta su criterio en U1-09, U4-03/U4-04, U5/C3, S1-09/U3-05, S8-06 o U2-06. R2-07: las 7 partes resueltas, con un residuo bajo (H-07).
- (c) 0 rutas `/tmp` en comandos (18 archivos y `tools/`). (d) Las filas S1–S8 y WF4 del INDICE son iguales al conteo real. Conteos correctos para el integrador: C1 21 (19/2), U1 16 (10/6), total 256 (196/60). (e) 18/18 en versión 1.1 con lista de cambios. (f) Ningún criterio reabre §A–§K; hay una observación baja sobre S5-04 (H-08).
- Recalculé desde `oraculo/model.json` y `datos/semilla/data.json` todas las cifras nuevas: S2-07 (7 filas), S2-08, S4-07, S5-02, S5-09 y S3-14. 0 diferencias.

## Archivos
- Creados: `docs/estado/p1/REVISION-criterios-r1.md` y este reporte.
- Temporales: `datos/tmp/revisor-criterios-r1/` (`rev.py`, `kgrep.py`, `citas_rev.py`), borrada. `s109_sql.sh` crea y borra su propia carpeta.

## Verificación
- `python3 -I rev.py` (propio): 18 archivos con 5 secciones, versión 1.1 y lista de cambios; 256 criterios (196/60/0) con Verifica y Evidencia en cada bloque; 0 fuera de la §5; 0 rutas `/tmp` en comandos. En el INDICE, S y WF4 coinciden, C1 y U1 difieren y el total dice 248 contra 256 reales. En las 403 líneas agregadas: 0 `jq` y 0 `python3` sin `-I`. Las cifras que no son hoja del oráculo son derivadas o ilustrativas y las recalculé.
- `python3 -I kgrep.py` (propio): K1–K34 y R2-07 con archivo:línea, más una búsqueda de restos obsoletos. Encontró «gasto 25°» en U2-06 y las menciones de 374.93 y «10 al 9» (estas en su contexto correcto).
- `python3 -I citas_rev.py` (propio): 252 citas de la §5 con fuente nombrada. 0 fragmentos que no estén textuales. 5 citas juntan dos fuentes o van fuera de orden (C1.md:130, C2.md:80 y 83, U2.md:95, U5.md:83).
- `bash docs/criterios/tools/verif_final.sh` (reejecutado a las 11:46):
  - `cuenta.py`: 134 (105/29), incompletos: ninguno. `traza.py`: 256 criterios, 0 sin §5. `refs.py`: 0 ID inexistentes.
  - `citas.py`: 253 citas, 5 «sin fuente», las mismas 5 citas combinadas. Residuos de N-xx: 0.
  - `canon.py`: 178 comprobadas, 0 diferencias, sha256 1a5706…bfd4. `s109_sql.sh`: UP OK, 54 tablas, 18 filas, la ventana 10→9 rechazada por el CHECK, febrero paga el 28. `s109_equiv.py`: 12 de 12.
  - `matriz_c102.py`: 203 N-xxx y 197 operaciones; 3 N-xxx sin operación (N-316, N-318 y N-319), porque otro agente está editando U1 1.1. Rutas `/tmp` en S y WF4: 0.
- Otras comprobaciones:
  - `grep` en `openapi.yaml`: el parámetro `IdempotencyKey` ya es `required: true` en C1 v1.0, así que C1-19 no inventa la obligatoriedad.
  - En decisiones C (decisiones.md:23) dice «pesos iguales normalizados», de donde sale el 0.25 de S3-14.
  - `docs/contratos/extractos/` no existe todavía (C1-21 pide algo que aún no se genera).

## Supuestos
- «Cada K de la lista» en (b) lo evalué contra §K completa, que es la base que fija la orden, y no solo contra las K que cada agente tenía asignadas.
- Que falte una decisión de §K en un criterio que debe verificarla lo califiqué como **media** (deja inconsistencia entre artefactos) y no como alta: no hay cifra ni esquema incorrecto y las verificaciones que existen pasan.
- `criterios-frescura-extremos` y `criterios-promedios-redondeo` los consideré bloqueos legítimos según el protocolo §2, porque son una frontera y una etapa de redondeo que §K no fija. La primera lleva una observación baja (H-08).
- Las secciones del INDICE fuera de las filas S y WF4 (total, §1, §5, §6, §8 y §9) le tocan al integrador. Las anoté como baja (H-12) para que no se pierdan.

## Bloqueos
Ninguno nuevo.

## Pendientes
Ninguno dentro del alcance. Para la ronda 2: reverificar H-01 a H-06 (y las bajas que se atiendan) con `rev.py`, `kgrep.py`, `citas_rev.py` y `tools/verif_final.sh`.

## Solicitudes fuera de alcance
- Agentes de criterios: aplicar las correcciones de H-01 a H-11 y H-13 en las rutas que da el informe (§7): `docs/criterios/U2.md`, `U4.md`, `U1.md`, `U5.md`, `C3.md`, `S1.md`, `U3.md`, `S8.md`, `S5.md` y `C1.md`, y la §5 de S1–S8 y WF4.
- Integrador:
  - `docs/criterios/INDICE.md` (H-12): filas C1 21 (19/2) y U1 16 (10/6), total 256 (196/60), sumando los criterios que se agreguen por H-04; notas 8 a 11; §6; §8; §9.
  - `docs/estado/bloqueos/INDICE.md`: agregar `criterios-promedios-redondeo` y `criterios-frescura-extremos`.

## Paquetes instalados
Ninguno. `verif_final.sh` usa `uv run --with pyyaml` desde la caché de uv.

## Consumo aproximado
34 acciones de herramienta, contando la salida final. Leído: unos 200 KB (protocolo, §K, INDICE, 2 reportes, `tools/`, líneas sueltas de criterios, colas de bloqueos y salidas de los scripts). Escrito: unos 45 KB (3 scripts temporales de unos 15 KB, informe de 22 KB, este reporte y el de progreso).
