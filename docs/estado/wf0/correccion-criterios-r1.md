# Reporte WF-0 · correccion-criterios-r1 (corrección del artefacto «criterios», ronda 1)

2026-10-10 · Estado: **terminado** (cerrado por tope: no) · Origen: `docs/estado/wf0/REVISION-r1.md` §5, hallazgos H-06, H-07 y H-08 (severidad media). Los criterios pasan a la versión 1.1 con el mismo conteo: 248 (189 alta, 59 media); no cambia ningún umbral numérico ni ninguna cifra canon.

## Hecho
- **H-06 · identificadores (pantalla P-xx, necesidad N-xxx).** Se reescribieron los criterios que usaban «N-xx» como identificador de pantalla:
  - `U1.md`: objetivo, U1-02 (pantalla `P-xx`, sección de Ajustes `P-26.xx`, necesidad `N-xxx`, componente `G-xx`, flujo `F-x`; el conteo se declara en pantallas `P-xx` y secciones `P-26.xx`, y su suma es T), U1-03 (matriz pantalla, necesidad y endpoint; cada `N-xxx` con operación o Local; cada `P-xx` con operación o Local o sin datos del servidor; una `N-xxx` bloqueada o PROVISIONAL sin endpoint no cuenta como cubierta), U1-04, U1-05, U1-07, U1-08 (`P-26.01` a `P-26.14` en orden), U1-09, «No aceptable» y umbrales.
  - `C1.md`: C1-02 cruza `N-xxx` con `operationId` y exige que cada `P-xx` tenga al menos una operación o la declaración Local; salida esperada con cinco conteos; C1-10 (cada vista con su pantalla `P-xx` de U1 §2.5 y §3.30); «No aceptable» y umbrales.
  - `C3.md`: C3-01 (filas por `P-xx`; el `grep` pasa de `^| N-` a `^| P-`), C3-02, C3-08 y umbrales. `U5.md`: U5-01, U5-03 y umbrales (el conteo `N` pasa a `T` para no chocar con `N-xxx`). `U4.md`: U4-01. `P3-claude-design.md`: P3-06.
  - `INDICE.md`: bullet de identificadores en §1, supuesto de §8 y fila de `criterios-nxx` en §6.
  - Bloqueo `docs/estado/bloqueos/criterios-nxx.md` **cerrado con la opción (b)** (sección «Resolución» al final; el texto original queda intacto).
- **H-07 · periodos.** `S1.md` S1-09 ya no pide insertar P1–P6 (ventanas del 10 al 9): pide (a) los 6 meses calendario de oct-2026 a mar-2027 con sus 12 quincenas, la quincena 2 de febrero 2027 en el día 28 y el rechazo por el `CHECK` de la ventana `2026-10-10` a `2026-11-09`, y (b) la equivalencia con el oráculo (efectivo y devengado por mes calendario contra `periodos[0..5]`: 12 de 12 cifras) más la condición de la que depende (0 eventos entre los días 1 y 9). Se agregó la nota de dependencia de la semilla en `S2.md` (S2-09, nota (e) de §6) y `S3.md` (nota (d) de §6). Consecuencias del mismo problema que también se alinearon: S2-06 (etiqueta de mes calendario en lugar de `10-oct-2026 a 9-nov-2026`), S2-03, S2-07 y §4 de S2, C1-09, C3-09 y S7-02 (aclaran que P1–P6 son oct-2026…mar-2027). «No aceptable» de S1 y trazabilidad de S1-09.
- **H-08 · trazabilidad.** Los 35 criterios sin línea en su §5 la tienen ahora (`traza.py`: 248 de 248). Cada línea cita el plan, decisiones o el protocolo; las 88 citas «…» agregadas se verificaron contra esas fuentes (`citas.py`, 0 sin fuente). 16 de los 35 llevan además una parte marcada **«convención propia del directorio»** porque no deriva de una frase de las fuentes (C1-04, C2-09, C2-10, C2-11, C3-03, C3-11, S1-15, S4-06, S6-08, S6-10, S6-15, S7-14, S8-08, U4-09, U4-10, U5-06); los otros 19 se trazan completos.
- `INDICE.md`: versión 1.1, §5 punto 12 (ventanas frente a mes calendario) y nueva §9 con el resumen de cambios por hallazgo. `docs/estado/wf0/criterios.md` (reporte original) actualizado: el supuesto 1 y la mención de C1-02 quedaron corregidos y se agregó la sección final de correcciones.
- No se aplicaron H-09 a H-12 (baja), salvo la cláusula «una necesidad bloqueada no cuenta como cubierta» de C1-02 y U1-03, que forma parte de la reescritura de H-06.

## Archivos
Modificados (todos bajo `/home/deepsec/finanzas-hogar/`):
- `docs/criterios/`: `INDICE.md`, `U1.md`, `C1.md`, `C3.md`, `U5.md`, `U4.md`, `P3-claude-design.md`, `S1.md`, `S2.md`, `S3.md`, `S7.md`, `C2.md`, `S4.md`, `S5.md`, `S6.md`, `S8.md`, `U3.md`. Sin cambios: `U2.md`, `WF4-operacion.md`. Diferencia contra la copia previa (guardada en el scratchpad): 118 líneas agregadas y 70 quitadas o reemplazadas; neto +48 (1,948 → 1,996 líneas en los 19 archivos).
- `docs/estado/wf0/criterios.md` (reporte original) y este reporte `docs/estado/wf0/correccion-criterios-r1.md`.
- `docs/estado/bloqueos/criterios-nxx.md` (estado y resolución; ver Supuestos, punto 1).
- Temporal creado y borrado: `datos/tmp/correccion-criterios-r1/cr1.db`; `datos/` vuelve a tener solo `semilla/`.
- Scripts de verificación (en el scratchpad de la sesión, efímero; fuentes en el Anexo): `traza.py`, `refs.py`, `citas.py`, `canon.py`, `s109_sql.sh`, `s109_equiv.py`, `matriz_c102.py`, `verif_final.sh`. Ruta: `/tmp/claude-1000/-home-deepsec/ad9b2909-7e27-4a0a-9e87-9e6a5395adce/scratchpad/cr1/`.

## Verificación
Salida literal de `verif_final.sh` (2026-10-10 23:57, árbol de trabajo compartido; los artefactos S1, C1 y U1 se estaban corrigiendo en paralelo por otros agentes):

```
== fecha: 2026-10-10 23:57:23
== 1. estructura y trazabilidad (traza.py)
archivos: 18 | criterios definidos: 248
secciones 1-5 faltantes: []
archivos con Verifica != criterios: [] | con Evidencia != criterios: []
criterios sin linea en su seccion 5 (con rangos expandidos): 0
== 2. referencias cruzadas entre IDs (refs.py)
criterios definidos: 248 | IDs distintos citados: 248 | citados que no existen: {}
== 3. citas de las lineas de trazabilidad AGREGADAS contra plan/decisiones/protocolo (citas.py)
citas verificadas: 88 | sin fuente: 0
== 4. residuos de N-xx de 2 equis fuera de la explicacion del INDICE
0
== 5. cifras canon R-01..R-85 (S2 y S3) y K-01..K-08 (S5) (canon.py)
sha256 model.json: 1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4
S2.md filas R verificadas acumuladas: 85
S3.md filas R verificadas acumuladas: 170
S5 filas K en el archivo: 8 | K-01..K-08 contra data.json: OK
TOTAL comprobadas: 178 | con diferencias: 0
== 6. S1-09 (a): SQL en base temporal (s109_sql.sh)
UP OK
tablas: 54
ventana 10->9 como mes: RECHAZADA por CHECK -> CHECK constraint failed: -- @rule Mes: del día 1 al último. 
filas insertadas: 18 (6 meses + 12 quincenas = 18)
recuperadas iguales a lo insertado: True
febrero 2027, quincena 2: ('2027-02-16', '2027-02-28', '2027-02-28') -> paga el 28: True
foreign_key_check filas: 0 | integrity_check: ok
mi carpeta temporal borrada: si | datos/ contiene: semilla 
== 7. S1-09 (b): equivalencia calendario vs ventanas del oraculo (s109_equiv.py)
eventos fechados: 28 | entre los dias 1 y 9 de un mes: 0
  P1 <-> 2026-10 efectivo  calendario=   973.75 oraculo=   973.75 OK
  P1 <-> 2026-10 devengado calendario= 1,947.50 oraculo= 1,947.50 OK
  P2 <-> 2026-11 efectivo  calendario= 1,947.50 oraculo= 1,947.50 OK
  P2 <-> 2026-11 devengado calendario= 1,947.50 oraculo= 1,947.50 OK
  P3 <-> 2026-12 efectivo  calendario= 2,232.43 oraculo= 2,232.43 OK
  P3 <-> 2026-12 devengado calendario= 2,517.37 oraculo= 2,517.37 OK
  P4 <-> 2027-01 efectivo  calendario= 2,716.36 oraculo= 2,716.36 OK
  P4 <-> 2027-01 devengado calendario= 2,716.37 oraculo= 2,716.37 OK
  P5 <-> 2027-02 efectivo  calendario= 2,517.38 oraculo= 2,517.38 OK
  P5 <-> 2027-02 devengado calendario= 2,517.37 oraculo= 2,517.37 OK
  P6 <-> 2027-03 efectivo  calendario= 2,716.36 oraculo= 2,716.36 OK
  P6 <-> 2027-03 devengado calendario= 2,716.37 oraculo= 2,716.37 OK
cifras comparadas: 12 | diferencias: 0
== 8. C1-02/U1-03: consistencia de la matriz sobre los artefactos actuales (matriz_c102.py)
U1: necesidades N-xxx definidas = 199 | pantallas P-xx = 27 | secciones P-26.xx = 14 | Local declaradas = ['N-003', 'N-301']
openapi: operaciones = 193 | N-xxx citadas en x-necesidades = 193
N-xxx sin operacion ni Local: 4 ['N-151', 'N-152', 'N-153', 'N-165']
P-xx (con secciones P-26.xx y Global) sin operacion ni Local: 0 []
operaciones sin N-xxx: 1 ['GET /openapi.yaml']
N-xxx citadas por el contrato que no existen en U1: []
matriz C1 (seccion 10): filas = 220 | N-xxx distintas = 199 | N-xxx de U1 ausentes en la matriz C1: []
estados en la matriz C1: {'BLOQUEO': 4, 'Local': 2, 'OK': 169, 'PROVISIONAL': 24}
```

Lectura de la salida:
1. Estructura y trazabilidad: 18 archivos, 248 criterios, cada uno con «Verifica» y «Evidencia»; antes de la corrección había 35 sin línea en §5 (el script reproduce exactamente la lista del revisor), ahora 0.
2. Referencias cruzadas: los 248 IDs citados existen; ninguna referencia a un ID inexistente.
3. Citas: 88 de 88 de las líneas agregadas existen en `plan.md`, `decisiones.md` o `protocolo-agentes.md` (fragmentos separados por «…» en orden; se ignoran asteriscos y comillas invertidas).
4. Residuos: 0 «N-xx» de dos equis fuera de la explicación histórica del INDICE.
5. Cifras canon: las tablas R-01 a R-85 (S2 y S3, idénticas entre sí y byte a byte iguales a las de antes de la corrección) y K-01 a K-08 (S5) coinciden con `oraculo/model.json` y `datos/semilla/data.json`: 178 comprobaciones, 0 diferencias; el sha256 de `model.json` sigue siendo el declarado en el INDICE.
6. S1-09 (a), sobre la migración vigente (54 tablas en este momento): los 6 meses y 12 quincenas se insertan y se recuperan iguales; febrero 2027 quincena 2 paga el 28; la ventana `2026-10-10` a `2026-11-09` como mes la rechaza el `CHECK`; `foreign_key_check` vacío e `integrity_check` ok.
7. S1-09 (b): 12 de 12 cifras iguales entre los meses calendario oct-2026…mar-2027 y `periodos[0..5]`, con 28 eventos fechados y ninguno entre los días 1 y 9. Contraejemplo de dependencia (hipotético): un gasto del 5-nov-2026 cae en P1 (`2026-10-10` a `2026-11-09`) en el oráculo y en noviembre en el mes calendario.
8. C1-02 y U1-03 (texto nuevo) sobre los artefactos de este momento: 199 `N-xxx`, 193 operaciones, 220 filas en la matriz de C1, 0 `P-xx` sin operación ni Local, una operación sin `N-xxx` (`GET /openapi.yaml`, excepción estática justificada en C1 §11.2) y **4 `N-xxx` sin operación ni Local (N-151, N-152, N-153, N-165)**. Es el hallazgo H-01, que corresponde a las rondas de S1, C1 y U1 (S1 ya empezó: escribió `S1-lista-compras.md`); los criterios ya no lo admiten como excepción. Reejecutar `matriz_c102.py` cuando C1 se regenere debe dar 0.

No aplica en esta corrección (no se tocó el artefacto): `go build` y `go vet` (`backend/`), `colorimetria.py` (`docs/design/tools/`), validador OpenAPI (`docs/contratos/openapi.yaml`). Los criterios solo cambian de redacción; las verificaciones propias del artefacto «criterios» son la estructura, las referencias, la trazabilidad, las cifras canon y las pruebas de S1-09 y C1-02 de arriba.

## Supuestos
1. **Alcance del bloqueo.** El hallazgo H-06 lista `docs/estado/bloqueos/criterios-nxx.md` en su ruta y es un bloqueo de la fase «criterios» (lo escribió el agente de criterios), pero mi lista de rutas de escritura no nombra la carpeta de bloqueos. Lo leí como parte del artefacto y solo agregué una línea de estado y la sección «Resolución» (el texto original no se tocó). Si el integrador lo considera fuera de alcance, basta quitar esa línea y esa sección.
2. **T.** En U1-02, U5-01 y C3-01, T = número de pantallas `P-xx` (Ajustes cuenta como una) más las 14 secciones `P-26.xx`; U1 entregó 27 y 14. Es la lectura que conserva el criterio original, donde cada «pantalla o nodo» tenía identificador.
3. **P-26.** Para U1-03 y C1-02, la pantalla contenedora `P-26` cumple por las `N-xxx` de sus 14 secciones (U1 §3.28 las lista bajo P-26 y §3.29 las reparte por sección).
4. **Fixtures con eventos entre los días 1 y 9.** Leí la sugerencia del revisor («comparar contra las ventanas del oráculo en las fixtures») como: el valor esperado sale del oráculo `fin_*.py` con las ventanas que la propia fixture declara, nunca de R-45…R-50. No agregué una fixture obligatoria nueva.
5. **Versión 1.1 y «convención propia».** Subir la versión de los 17 archivos tocados y usar la etiqueta «convención propia del directorio» son convenciones de este directorio; no cambian cifras, esquema, endpoints ni tokens.
6. **Temporal compartido.** Mis dos primeras pruebas de S1-09 usaron `datos/tmp/cr1.db` y borraron `datos/tmp` completo (la carpeta no existía antes de mi prueba). Después vi que la fecha de modificación de `datos/` cambió a las 23:56:42 sin que yo ejecutara nada, así que parece que otro agente también usa `datos/tmp`; desde la tercera ejecución uso un subdirectorio propio y borro solo mi archivo y, con `rmdir`, las carpetas vacías (`rmdir` no elimina carpetas con contenido ajeno).
Ninguno cambia cifras, esquema, endpoints ni tokens.

## Bloqueos
- Ninguno nuevo. Cerrado: `docs/estado/bloqueos/criterios-nxx.md` (opción b).
- Siguen abiertos los otros 15 `criterios-*.md` (`INDICE.md` §6) y, de otras fases, los que afectan a la puerta 1 (`C1-lista-compras`, `S1-lista-compras`, `C1-huecos-u1-s1`, `S1-huecos-negocio`).

## Pendientes
- **H-10 (baja), no asignado a esta ronda:** U2-02 sin la propuesta de 135° para ingreso; U1-05 espera exactamente 6 flujos (U1 entregó 6 + 2: F-7 y F-8); C1-01 usa un módulo que no está en el Python del sistema (equivalente: `uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator <archivo>`); S1-03 sin regla para reglas sin datos en el oráculo. H-09, H-11 y H-12 no son de criterios.
- **Citas abreviadas del texto original:** 25 de las 190 citas «…» de las secciones 5 (versión original, no las agregadas) son resúmenes o abreviaturas y no citas literales (p. ej. `C1.md` «Dashboard (≤ 6 tarjetas)», `S7.md` «Reportes por defecto: …», `U1.md` «Axel y Ximena co-titulares, mismos permisos»). El contenido es coherente con decisiones, pero no se tocaron porque no son hallazgo del revisor. `python3 -I citas.py` sin argumento (versión del Anexo) las lista, junto con las agregadas.
- Reejecutar `matriz_c102.py` y la verificación de `traza.py` después de las rondas de S1, C1 y U1 (la primera debe dar 0 `N-xxx` sin operación ni Local).
- El agente de la corrección de S1 debe evidenciar S1-09 con el texto nuevo (18 filas, rechazo de la ventana 10 al 9, 12 de 12 cifras): el script `s109_sql.sh` y `s109_equiv.py` del Anexo lo hacen.

## Solicitudes fuera de alcance
1. **Integrador** · `docs/criterios/tools/` (o `ops/tools/criterios/`) · copiar los scripts del Anexo (`traza.py`, `refs.py`, `citas.py`, `canon.py`, `s109_sql.sh`, `s109_equiv.py`, `matriz_c102.py`, `verif_final.sh`) · motivo: el scratchpad es efímero y los revisores reejecutan los comandos (mismo motivo que H-05). Mi alcance de escritura no incluye subcarpetas de `docs/criterios/`.
2. **Integrador** · `docs/estado/bloqueos/criterios-toques.md` · cerrarlo con la opción (a), como pidió U1 (solicitud 2 de su reporte): la regla de conteo ya está escrita en U1 §4 y U1-04 la exige. No lo toqué porque no es parte de H-06, H-07 ni H-08.
3. **Corrección de S1 (si aún corre)** · avisar que S1-09 cambió (ver Hecho) para que su evidencia use las pruebas nuevas.
4. **WF-1, motor financiero** · al escribir las fixtures de S2: declarar ventanas por fixture (S2-09) y no reutilizar R-45…R-50 como esperado cuando haya eventos entre los días 1 y 9.

## Paquetes instalados
Ninguno. Usé una vez `uv run --no-project --with pyyaml` (entorno efímero de la caché de uv, sin red) antes de comprobar que el Python del sistema ya trae `yaml`; los scripts finales corren con `python3 -I` y la biblioteca estándar más `yaml`.

## Consumo aproximado
- Unas 95 acciones de herramienta (sin subagentes): 1 bloque de lecturas (protocolo, decisiones, REVISION-r1, reportes y bloqueos), 18 archivos de criterios, secciones de U1, C1 y del plan con `sed` y `grep`; 11 scripts de edición con comprobación de una sola aparición por reemplazo y varios cambios puntuales; 9 scripts de verificación.
- Leído: unos 280 KB (el cuerpo de `U1-arquitectura.md`, `C1-datos.md`, `openapi.yaml`, `model.json` y la migración se leyeron por secciones o por script).
- Escrito: 118 líneas agregadas y 70 quitadas o reemplazadas en los criterios, 1 resolución de bloqueo, 1 sección en el reporte original y este reporte.
- Contexto acumulado: del orden de 400K tokens según el contador de la sesión, por encima de los 300K de planeación; el trabajo asignado (H-06, H-07 y H-08) quedó completo y consistente, por lo que no se marca cierre por tope.

## Anexo · fuentes de los scripts de verificación
Se ejecutan con `python3 -I <script>` (o `bash` para los `.sh`) desde cualquier carpeta; no escriben en el repositorio salvo la base temporal propia de `s109_sql.sh`.

### traza.py
```python
#!/usr/bin/env python3
"""Verifica la estructura de docs/criterios y la trazabilidad (seccion 5) de cada criterio.
Uso: python3 -I traza.py [--lista]   (sin argumentos: resumen; --lista: imprime los IDs sin traza)"""
import re, sys, glob, os
BASE = '/home/deepsec/finanzas-hogar/docs/criterios'
ARCH = [f for f in sorted(glob.glob(BASE + '/*.md')) if os.path.basename(f) != 'INDICE.md']
ID = re.compile(r'\b((?:WF4|P3|[SUC][0-9])-[0-9]{2})\b')

def seccion(txt, n):
    m = re.search(r'^## %d\. .*?$' % n, txt, re.M)
    if not m:
        return ''
    ini = m.end()
    m2 = re.search(r'^## \d+\. ', txt[ini:], re.M)
    return txt[ini:ini + m2.start()] if m2 else txt[ini:]

def expandir(txt):
    """IDs citados en un texto, expandiendo rangos 'X-01...X-04', 'X-01 a X-04' y 'X-01-X-04'."""
    ids = set(ID.findall(txt))
    for a, b in re.findall(r'\b((?:WF4|P3|[SUC][0-9])-[0-9]{2})\s*(?:…|\.\.\.|a|–|—)\s*((?:WF4|P3|[SUC][0-9])-[0-9]{2})\b', txt):
        pa, na = a.rsplit('-', 1); pb, nb = b.rsplit('-', 1)
        if pa == pb:
            for k in range(int(na), int(nb) + 1):
                ids.add('%s-%02d' % (pa, k))
    return ids

total = 0; sin = []; estruct_mal = []; sin_ver = []; sin_evi = []
for f in ARCH:
    t = open(f, encoding='utf-8').read()
    nom = os.path.basename(f)
    for n in range(1, 6):
        if not re.search(r'^## %d\. ' % n, t, re.M):
            estruct_mal.append((nom, n))
    s2 = seccion(t, 2); s5 = seccion(t, 5)
    defs = re.findall(r'^\*\*((?:WF4|P3|[SUC][0-9])-[0-9]{2}) \((alta|media|baja)\)', s2, re.M)
    ids = [d[0] for d in defs]
    total += len(ids)
    ver = len(re.findall(r'^- Verifica:', s2, re.M)); evi = len(re.findall(r'^- Evidencia:', s2, re.M))
    if ver != len(ids): sin_ver.append((nom, len(ids), ver))
    if evi != len(ids): sin_evi.append((nom, len(ids), evi))
    cit = expandir(s5)
    for i in ids:
        if i not in cit:
            sin.append((nom, i))
print('archivos:', len(ARCH), '| criterios definidos:', total)
print('secciones 1-5 faltantes:', estruct_mal)
print('archivos con Verifica != criterios:', sin_ver, '| con Evidencia != criterios:', sin_evi)
print('criterios sin linea en su seccion 5 (con rangos expandidos):', len(sin))
if '--lista' in sys.argv:
    por = {}
    for nom, i in sin: por.setdefault(nom, []).append(i)
    for k, v in por.items(): print(' ', k, ' '.join(v))
```

### refs.py
```python
import re, glob, os
B = '/home/deepsec/finanzas-hogar/docs/criterios/'
ID = re.compile(r'\b((?:WF4|P3|[SUC][0-9])-[0-9]{2})\b')
definidos = set(); citados = {}
for f in sorted(glob.glob(B + '*.md')):
    t = open(f, encoding='utf-8').read()
    for m in re.finditer(r'^\*\*((?:WF4|P3|[SUC][0-9])-[0-9]{2}) \((?:alta|media|baja)\)', t, re.M):
        definidos.add(m.group(1))
    for i in ID.findall(t):
        citados.setdefault(i, set()).add(os.path.basename(f))
# IDs de otros dominios que no son criterios: R-nn, K-nn no coinciden con el patron; U1-02 etc. si
inexistentes = {i: sorted(v) for i, v in citados.items() if i not in definidos}
print('criterios definidos:', len(definidos), '| IDs distintos citados:', len(citados), '| citados que no existen:', inexistentes)
```

### citas.py (argumento: carpeta con la copia previa de `docs/criterios`; sin argumento revisa todas las citas)
```python
#!/usr/bin/env python3
"""Comprueba que cada cita «...» de las lineas de trazabilidad (seccion 5) AGREGADAS en esta ronda
existe en plan.md, decisiones.md o protocolo-agentes.md (fragmentos separados por … en orden)."""
import re, glob, os, sys
R = '/home/deepsec/finanzas-hogar/docs/'
norm = lambda s: re.sub(r'\s+', ' ', s.replace('*', '').replace('`', '').replace(' ', ' ')).strip().lower()
FUENTES = {n: norm(open(R + p, encoding='utf-8').read()) for n, p in
           (('plan', 'plan/plan.md'), ('decisiones', 'plan/decisiones.md'), ('protocolo', 'plan/protocolo-agentes.md'))}

def citas(linea):
    out, depth, ini = [], 0, None
    for i, ch in enumerate(linea):
        if ch == '«':
            if depth == 0: ini = i + 1
            depth += 1
        elif ch == '»' and depth:
            depth -= 1
            if depth == 0: out.append(linea[ini:i])
    return out

def seccion5(t):
    m = re.search(r'^## 5\. Trazabilidad\s*$', t, re.M)
    ini = m.end(); m2 = re.search(r'^## \d+\. ', t[ini:], re.M)
    return t[ini:ini + m2.start()] if m2 else t[ini:]

ORIG = sys.argv[1] if len(sys.argv) > 1 else None
total = bad = 0
for f in sorted(glob.glob(R + 'criterios/*.md')):
    nom = os.path.basename(f)
    if nom == 'INDICE.md': continue
    nuevo = seccion5(open(f, encoding='utf-8').read()).splitlines()
    viejo = set(seccion5(open(ORIG + '/' + nom, encoding='utf-8').read()).splitlines()) if ORIG else set()
    for ln in nuevo:
        if ln in viejo or not ln.startswith('- '): continue
        for c in citas(ln):
            partes = [norm(p) for p in re.split(r'…|\.\.\.', c) if norm(p)]
            # solo se verifican citas de las fuentes del plan (las de nombres internos como «Local» o «sin precio» no)
            if not re.search(r'plan\.md|decisiones|protocolo', ln): continue
            total += 1
            ok_fuentes = []
            for n, src in FUENTES.items():
                pos = 0; ok = True
                for p in partes:
                    k = src.find(p, pos)
                    if k < 0: ok = False; break
                    pos = k + len(p)
                if ok: ok_fuentes.append(n)
            if not ok_fuentes:
                bad += 1; print(f'  SIN FUENTE en {nom}: «{c[:100]}»')
print(f'citas verificadas: {total} | sin fuente: {bad}')
```

### canon.py
```python
#!/usr/bin/env python3
"""Reverifica las tablas R-01..R-85 (S2 y S3) y K-01..K-08 (S5) contra oraculo/model.json y datos/semilla/data.json.
Cada celda 'Valor esperado' se descompone en fechas y numeros; cada uno debe existir entre las hojas del nodo citado."""
import json, re, hashlib
from decimal import Decimal as D
R = '/home/deepsec/finanzas-hogar/'
M = json.load(open(R + 'oraculo/model.json')); DJ = json.load(open(R + 'datos/semilla/data.json'))
print('sha256 model.json:', hashlib.sha256(open(R + 'oraculo/model.json', 'rb').read()).hexdigest())
def hojas(x, out):
    if isinstance(x, dict):
        for v in x.values(): hojas(v, out)
    elif isinstance(x, list):
        for v in x: hojas(v, out)
    elif isinstance(x, bool) or x is None: pass
    elif isinstance(x, (int, float)): out.append(('n', D(repr(x))))
    else: out.append(('s', str(x)))
    return out
def resolver(raiz, ruta):
    nodo = raiz
    for parte in re.findall(r'[A-Za-z_0-9]+|\[[^\]]+\]', ruta):
        if parte.startswith('['):
            idx = parte[1:-1]
            if '..' in idx:
                a, b = idx.split('..'); nodo = nodo[int(a):int(b) + 1]
            else: nodo = nodo[int(idx)]
        else: nodo = nodo[parte]
    return nodo
def toks(cel):
    t = []
    for m in re.finditer(r'\d{4}-\d{2}-\d{2}|\d[\d,]*\.\d+|\d+', cel):
        s = m.group(0)
        t.append(('f', s) if '-' in s else ('n', D(s.replace(',', ''))))
    return t
def chequea(fila_id, ruta, esp):
    raiz = M; ruta = ruta.strip('`').strip()
    if ruta.startswith('data.json:'): raiz = DJ; ruta = ruta[len('data.json:'):]
    partes = [p.strip() for p in ruta.split(' / ')]
    nodos = []
    base = partes[0]
    nodos.append(resolver(raiz, base))
    for p in partes[1:]:
        if p.startswith('.'):
            padre = base.rsplit('.', 1)[0]; nodos.append(resolver(raiz, padre + p))
        else: nodos.append(resolver(raiz, p))
    h = []
    for n in nodos: hojas(n, h)
    nums = [v for k, v in h if k == 'n']; strs = {v for k, v in h if k == 's'}
    malos = []
    for d in re.findall(r'despensa (\d)', esp):
        if not any('despensa %s' % d in s for s in strs): malos.append('despensa ' + d)
    esp = re.sub(r'despensa \d', '', esp)
    for k, v in toks(esp):
        if k == 'f':
            if v not in strs: malos.append(v)
        else:
            if not any(abs(v - n) <= D('0.0000000001') for n in nums): malos.append(str(v))
    return malos
tot = bad = 0
for archivo in ('S2.md', 'S3.md'):
    for ln in open(R + 'docs/criterios/' + archivo, encoding='utf-8'):
        m = re.match(r'\| (R-\d\d) \| (.*?) \| (`[^|]*`) \| (.*?) \| (S2|S3) \|', ln)
        if not m: continue
        rid, desc, ruta, esp, motor = m.groups()
        tot += 1
        if rid == 'R-67':                               # 17 eventos
            ok = len(M['credito']['eventos']) == 17; malos = [] if ok else ['n eventos']
        elif rid == 'R-04':                             # razon sin redondear
            malos = [] if abs(D(repr(M['parametros']['tasa_pct'])) - D('25.8684763290')) < D('1e-9') else ['tasa_pct']
        else:
            malos = chequea(rid, ruta, esp)
        if malos: bad += 1; print('  DIFIERE', archivo, rid, ruta, malos)
    print(archivo, 'filas R verificadas acumuladas:', tot)
# K-01..K-08 (S5) contra data.json
T = DJ['totales']
K = {'K-01': T['bodega'], 'K-02': T['walmart'], 'K-03': round(T['walmart'] - T['bodega'], 2), 'K-04': T['bodega_identicos'],
     'K-05': T['walmart_identicos'], 'K-06': round(T['walmart_identicos'] - T['bodega_identicos'], 2)}
esperado = {'K-01': 3094.50, 'K-02': 3233.00, 'K-03': 138.50, 'K-04': 2773.50, 'K-05': 2874.00, 'K-06': 100.50}
for k, v in esperado.items():
    ok = abs(K[k] - v) < 1e-9; tot += 1
    if not ok: bad += 1; print('  DIFIERE', k, K[k], v)
for k, ok in (('K-07', (T['articulos_bodega'], T['articulos_walmart'], T['articulos_identicos']) == (61, 60, 55)),
              ('K-08', (len(DJ['items_identicos']), len(DJ['sustitutos'])) == (33, 3))):
    tot += 1
    if not ok: bad += 1; print('  DIFIERE', k)
ks = [l for l in open(R + 'docs/criterios/S5.md', encoding='utf-8') if re.match(r'\| K-\d\d \|', l)]
print('S5 filas K en el archivo:', len(ks), '| K-01..K-08 contra data.json: OK' if bad == 0 else '')
print(f'TOTAL comprobadas: {tot} | con diferencias: {bad}')
```

### s109_sql.sh
```bash
#!/bin/bash
# Verificacion SQL de S1-09 (texto nuevo): meses calendario y quincenas oct-2026..mar-2027 en base temporal bajo datos/tmp
set -u
R=/home/deepsec/finanzas-hogar
D=$R/datos/tmp/correccion-criterios-r1
DB=$D/cr1.db
mkdir -p $D && rm -f $DB
awk '/^-- \+goose Up/{f=1;next} /^-- \+goose Down/{f=0} f' $R/backend/db/migrations/00001_esquema.sql > /tmp/claude-1000/-home-deepsec/ad9b2909-7e27-4a0a-9e87-9e6a5395adce/scratchpad/cr1/up.sql
sqlite3 -bail $DB < /tmp/claude-1000/-home-deepsec/ad9b2909-7e27-4a0a-9e87-9e6a5395adce/scratchpad/cr1/up.sql && echo "UP OK"
echo "tablas: $(sqlite3 $DB "select count(*) from sqlite_master where type='table' and name not like 'sqlite_%'")"
python3 -I - <<'PY'
import sqlite3, calendar, datetime
DB='/home/deepsec/finanzas-hogar/datos/tmp/correccion-criterios-r1/cr1.db'
con=sqlite3.connect(DB); con.execute('PRAGMA foreign_keys=ON')
# 1) la ventana del oraculo (10 -> 9) como mes calendario debe rechazarse (antes de insertar los meses, sin conflicto de unicidad)
try:
    con.execute("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES('month',NULL,2026,10,NULL,'2026-10-10','2026-11-09',NULL)")
    print('VENTANA 10->9 ACEPTADA (inesperado)')
except sqlite3.IntegrityError as e:
    print('ventana 10->9 como mes: RECHAZADA por CHECK ->', str(e).split(chr(10))[0][:60])
meses=[(2026,10),(2026,11),(2026,12),(2027,1),(2027,2),(2027,3)]
esperado=[]
for (y,m) in meses:
    ult=calendar.monthrange(y,m)[1]
    ini=f'{y:04d}-{m:02d}-01'; fin=f'{y:04d}-{m:02d}-{ult:02d}'
    cur=con.execute("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES('month',NULL,?,?,NULL,?,?,NULL)",(y,m,ini,fin))
    pid=cur.lastrowid
    q1=(f'{y:04d}-{m:02d}-01',f'{y:04d}-{m:02d}-15',f'{y:04d}-{m:02d}-15')
    q2=(f'{y:04d}-{m:02d}-16',fin,min(f'{y:04d}-{m:02d}-30',fin))
    con.execute("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES('fortnight',?,?,?,1,?,?,?)",(pid,y,m)+q1)
    con.execute("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES('fortnight',?,?,?,2,?,?,?)",(pid,y,m)+q2)
    esperado.append(('month',y,m,None,ini,fin,None)); esperado.append(('fortnight',y,m,1)+q1); esperado.append(('fortnight',y,m,2)+q2)
con.commit()
filas=con.execute("select kind,year,month,half,start_date,end_date,pay_date from period order by year,month,kind='fortnight',half").fetchall()
print('filas insertadas:', len(filas), '(6 meses + 12 quincenas = 18)')
print('recuperadas iguales a lo insertado:', sorted(filas)==sorted(esperado))
feb=con.execute("select start_date,end_date,pay_date from period where kind='fortnight' and year=2027 and month=2 and half=2").fetchone()
print('febrero 2027, quincena 2:', feb, '-> paga el 28:', feb[2]=='2027-02-28')
print('foreign_key_check filas:', len(con.execute('PRAGMA foreign_key_check').fetchall()), '| integrity_check:', con.execute('PRAGMA integrity_check').fetchone()[0])
con.close()
PY
rm -f $DB; rmdir $D; rmdir $R/datos/tmp 2>/dev/null
echo "mi carpeta temporal borrada: $([ -e $D ] && echo NO || echo si) | datos/ contiene: $(ls $R/datos | tr '\n' ' ')"
```

### s109_equiv.py
```python
#!/usr/bin/env python3
"""Equivalencia P1-P6 del oraculo = meses calendario oct-2026..mar-2027 (S1-09, texto nuevo).
Lee solo oraculo/model.json (no escribe). Dinero en centavos enteros."""
import json, datetime as dt
from decimal import Decimal as D, ROUND_HALF_UP
M = json.load(open('/home/deepsec/finanzas-hogar/oraculo/model.json'))
c = lambda x: int((D(str(x)) * 100).quantize(D(1), rounding=ROUND_HALF_UP))
F = dt.date.fromisoformat
# eventos fechados de la semilla
pagos = [(F(p['fecha']), c(p['monto'])) for d in M['despensas'] for p in d['pagos']]
papel_contado = [(F(x['compra']), c(x['monto'])) for x in M['papel'] if 'contado' in x['forma']]
compras = [F(d['fecha']) for d in M['despensas']] if 'fecha' in M['despensas'][0] else []
eventos = [e[0] for e in pagos + papel_contado] + compras + [F(x['compra']) for x in M['papel']]
print('eventos fechados:', len(eventos), '| entre los dias 1 y 9 de un mes:', sum(1 for e in eventos if e.day <= 9))
meses = [(2026, 10), (2026, 11), (2026, 12), (2027, 1), (2027, 2), (2027, 3)]
dif = 0; n = 0
for i, (y, m) in enumerate(meses):
    P = M['periodos'][i]
    ef = sum(v for d, v in pagos if (d.year, d.month) == (y, m)) + sum(v for d, v in papel_contado if (d.year, d.month) == (y, m))
    papel_m = sum(v for d, v in papel_contado if (d.year, d.month) == (y, m))
    desp = M['despensas'][i // 2]              # cada despensa cubre 2 meses (bimestral)
    dev = c(D(str(desp['total'])) / 2) + papel_m
    for nombre, mio, ora in (('efectivo', ef, c(P['efectivo']['total'])), ('devengado', dev, c(P['devengado']['total']))):
        n += 1
        ok = mio == ora
        dif += (not ok)
        print(f"  {P['periodo']} <-> {y}-{m:02d} {nombre:9s} calendario={mio/100:>9,.2f} oraculo={ora/100:>9,.2f} {'OK' if ok else 'DIFIERE'}")
print(f'cifras comparadas: {n} | diferencias: {dif}')
```

### matriz_c102.py
```python
#!/usr/bin/env python3
"""Comprobacion del texto nuevo de C1-02 / U1-03 sobre los artefactos actuales (solo lectura).
Uso: uv run --no-project --with pyyaml python3 -I matriz_c102.py"""
import re, sys, yaml
R = '/home/deepsec/finanzas-hogar/'
u1 = open(R + 'docs/design/U1-arquitectura.md', encoding='utf-8').read()
c1 = open(R + 'docs/contratos/C1-datos.md', encoding='utf-8').read()
api = yaml.safe_load(open(R + 'docs/contratos/openapi.yaml', encoding='utf-8'))

# --- U1: necesidades definidas (seccion 3, antes de 3.29) y cuales declara Local
s3 = u1[u1.index('\n## 3. '):u1.index('### 3.29 ')]
defs = re.findall(r'^\| (N-\d{3}) \| ([LE]) \|(.*)$', s3, re.M)
N = {d[0] for d in defs}
local = {d[0] for d in defs if d[2].rstrip().rstrip('|').strip().endswith('Local')}
# --- U1 3.29: pantalla -> necesidades (se expanden rangos "N-020 a N-029")
def expandir(txt):
    ids = set(re.findall(r'N-\d{3}', txt))
    for a, b in re.findall(r'(N-\d{3}) a (N-\d{3})', txt):
        ids |= {'N-%03d' % k for k in range(int(a[2:]), int(b[2:]) + 1)} & N
    return ids
m329 = u1[u1.index('### 3.29 '):u1.index('### 3.30 ')]
pant = {}
for fila in re.findall(r'^\| (Global[^|]*|P-\d\d(?:\.\d\d)?) \|(.*)$', m329, re.M):
    clave = 'Global' if fila[0].startswith('Global') else fila[0]
    pant.setdefault(clave, set()).update(expandir(fila[1]))
# P-26 (contenedor) = union de sus secciones
pant['P-26'] = set().union(*[v for k, v in pant.items() if k.startswith('P-26.')])
P_inv = set(re.findall(r'^\| (P-\d\d) \| ', u1[u1.index('### 2.1 '):u1.index('### 2.2 ')], re.M))
S_inv = set(re.findall(r'^\| (P-26\.\d\d) \| ', u1[u1.index('### 2.3 '):u1.index('### 2.4 ')], re.M))
# --- openapi: operaciones y x-necesidades
ops = {}
for ruta, item in api['paths'].items():
    for met, op in item.items():
        if met in ('get', 'post', 'put', 'patch', 'delete'):
            ops[op['operationId']] = (met.upper() + ' ' + ruta, set(op.get('x-necesidades') or []))
con_op = set().union(*[v[1] for v in ops.values()])
# --- C1 matriz (seccion 10): filas pantalla | N | tipo | texto | operaciones | estado
m10 = c1[c1.index('## 10. '):c1.index('## 11. ')]
filas = re.findall(r'^\| (Global|P-\d\d(?:\.\d\d)?) \| (N-\d{3}) \| ([LE]) \|.*\| ([^|]*) \|$', m10, re.M)
estado = {}
for p, n, t, e in filas:
    estado.setdefault(n, set()).add(e.strip().split(' ')[0])
sin_op_ni_local = sorted(N - con_op - local)
print(f'U1: necesidades N-xxx definidas = {len(N)} | pantallas P-xx = {len(P_inv)} | secciones P-26.xx = {len(S_inv)} | Local declaradas = {sorted(local)}')
print(f'openapi: operaciones = {len(ops)} | N-xxx citadas en x-necesidades = {len(con_op)}')
print(f'N-xxx sin operacion ni Local: {len(sin_op_ni_local)} {sin_op_ni_local}')
pant_mal = []
for p in sorted(P_inv | S_inv | {'Global'}):
    necs = pant.get(p, set())
    if not necs:
        pant_mal.append((p, 'sin necesidades en la matriz de U1'))
    elif not (necs & con_op) and not necs <= local:
        pant_mal.append((p, 'ninguna de sus N-xxx tiene operacion'))
print(f'P-xx (con secciones P-26.xx y Global) sin operacion ni Local: {len(pant_mal)} {pant_mal}')
sin_n = sorted(v[0] for k, v in ops.items() if not v[1])
print(f'operaciones sin N-xxx: {len(sin_n)} {sin_n}')
print(f'N-xxx citadas por el contrato que no existen en U1: {sorted(con_op - N)}')
c1_n = {n for _, n, _, _ in filas}
print(f'matriz C1 (seccion 10): filas = {len(filas)} | N-xxx distintas = {len(c1_n)} | N-xxx de U1 ausentes en la matriz C1: {sorted(N - c1_n)}')
print('estados en la matriz C1:', {e: sum(1 for v in estado.values() if e in v) for e in sorted({x for v in estado.values() for x in v})})
```

### verif_final.sh
```bash
#!/bin/bash
# Verificacion final de correccion-criterios-r1 (solo lectura salvo la base temporal propia bajo datos/tmp)
S=/tmp/claude-1000/-home-deepsec/ad9b2909-7e27-4a0a-9e87-9e6a5395adce/scratchpad/cr1
cd $S
echo "== fecha: $(date '+%F %T')"
echo "== 1. estructura y trazabilidad (traza.py)"; python3 -I traza.py
echo "== 2. referencias cruzadas entre IDs (refs.py)"; python3 -I refs.py
echo "== 3. citas de las lineas de trazabilidad AGREGADAS contra plan/decisiones/protocolo (citas.py)"; python3 -I citas.py $S/orig/criterios
echo "== 4. residuos de N-xx de 2 equis fuera de la explicacion del INDICE"; grep -n -E "N-xx([^x]|$)" /home/deepsec/finanzas-hogar/docs/criterios/*.md | grep -v "INDICE.md" | wc -l
echo "== 5. cifras canon R-01..R-85 (S2 y S3) y K-01..K-08 (S5) (canon.py)"; python3 -I canon.py
echo "== 6. S1-09 (a): SQL en base temporal (s109_sql.sh)"; ./s109_sql.sh
echo "== 7. S1-09 (b): equivalencia calendario vs ventanas del oraculo (s109_equiv.py)"; python3 -I s109_equiv.py
echo "== 8. C1-02/U1-03: consistencia de la matriz sobre los artefactos actuales (matriz_c102.py)"; python3 -I matriz_c102.py
```
