# Reporte P1 · criterios-UC-1.1 (U1 a U5, C1 a C3 y P3: criterios 1.1)

11-oct-2026 · Estado: **terminado** · Cerrado por tope: no · Origen: `docs/plan/decisiones.md` §K (puerta 1) y hallazgos bajos de `docs/estado/wf0/REVISION-r3.md` (R2-05, R2-06, R2-07, R3-01, R3-04).

## Hecho
- **U1.md**: U1-04 con la regla K10 (opción a, confirmada); U1-05 pide al menos 6 flujos con los 4 obligatorios (R2-07; §1, §3 y §4 ajustados); U1-11 con objetivo táctil ≥ 44 px (R2-05); **nuevos** U1-15 (alta, estado «Sin conexión» con cola de gastos y peso, K9) y U1-16 (media, «lista de compras vigente» definida y `period` explícito, R3-01); tres citas de §5 vueltas textuales y cinco líneas nuevas de trazabilidad.
- **U2.md** (de versión 1 a 1.1): U2-02 con gasto en 40–45° (K26), Karen y Bodega Aurrera reubicadas y 135° admitido como propuesta para el ingreso; U2-08 con la marca (700 en claro, 400 en oscuro) en el barrido de pares cercanos y sin la colisión error–gasto de 5° (R3-04); U2-09 con K27 (series consecutivas; los 28 pares solo informativos); U2-11 con táctil ≥ 44 px; §3, §4 y §5 ajustados; la cita de decisiones I pasa a «error va en rojo 20°».
- **U5.md**: U5-09 con K10 (cuentan solo los toques en controles; `fill`/`type` y valores por defecto no cuentan).
- **C3.md**: C3-08 ejecuta los flujos de U1 (≥ 6) con K10; §1 y §4 ajustados.
- **C1.md**: C1-01 con el comando `uvx` primero; C1-02 y C1-03 citan `verify_c1.py --db` y `build.py --check` (bases temporales en `datos/tmp`); C1-06 con `family` de tres valores incluyentes (K5) y semana/quincenas de K8; C1-07 con «Z» solo en entradas (K12) y K8; C1-12 con `SameSite=Strict` (K6) y K11; C1-13 pasa a C1 v1.1; C1-15 ya no exige marca a lo aprobado en K6 y K7 (ni a lo resuelto por K1, K5 y K12); **nuevos** C1-16 (K2, `scope`), C1-17 (K3, borrador del ticket), C1-18 (K4, candidatos y confirmación del chat), C1-19 (K9, `Idempotency-Key` obligatoria y `captured_at`, reintento idempotente), C1-20 (K12, decimales capturados) y C1-21 (extractos por módulo: cada operación una vez, ≤ 80 KB, regenerables con `build.py --extractos`); §3, §4 y §5 ajustados y dos citas vueltas textuales.
- **C2.md**: C2-07 usa `datos/tmp` en lugar de `/tmp`.
- **U3.md, U4.md, P3-claude-design.md**: versión 1.1 con lista de cambios; en U3 y P3 una cita de §5 vuelta textual; sin cambios de fondo.
- Cada uno de los 9 archivos lleva «Versión 1.1 (11-oct-2026)» y la lista «Cambios de la versión 1.1» antes de §1.
- Las paráfrasis entre comillas de §5 se citan textuales o sin comillas: 11 fragmentos no textuales al inicio, 0 al final.

## Archivos
Modificados en `/home/deepsec/finanzas-hogar/docs/criterios/`: `U1.md`, `U2.md`, `U3.md`, `U4.md`, `U5.md`, `C1.md`, `C2.md`, `C3.md`, `P3-claude-design.md`. Creado: este reporte. No se tocó `INDICE.md` ni `docs/criterios/tools/`.

## Verificación
- `traza.py` (anexo de `docs/estado/wf0/correccion-criterios-r1.md`, copiado a `datos/tmp`): «archivos: 18 | criterios definidos: 256; secciones 1-5 faltantes: []; archivos con Verifica != criterios: [] | con Evidencia != criterios: []; criterios sin línea en su sección 5: 0».
- `refs.py` (mismo anexo): «criterios definidos: 256 | IDs distintos citados: 256 | citados que no existen: {}».
- `citas.py` (anexo de este reporte; compara cada fragmento entre « » de la §5 con `plan.md`, `decisiones.md` y `protocolo-agentes.md`, sin distinguir mayúsculas ni `**`, comillas inversas o « »; parte las citas por «…»): primera corrida 139 fragmentos y 11 no textuales; final **162 fragmentos, 0 no textuales**.
- Conteo alta/media por archivo (para el INDICE), por script sobre `^\*\*ID (alta|media)`:

| Archivo | Criterios | Alta | Media |
|---|---|---|---|
| C1.md | 21 | 19 | 2 |
| C2.md | 12 | 8 | 4 |
| C3.md | 11 | 9 | 2 |
| P3-claude-design.md | 13 | 10 | 3 |
| U1.md | 16 | 10 | 6 |
| U2.md | 14 | 10 | 4 |
| U3.md | 12 | 8 | 4 |
| U4.md | 11 | 8 | 3 |
| U5.md | 12 | 9 | 3 |

  Total de los 18 archivos con los S1–S8 y WF4 sin cambios: **256 criterios, 196 alta y 60 media** (antes 248, 189 y 59: +8 = U1-15 alta, U1-16 media, C1-16 a C1-21 alta). Cada criterio de los 9 archivos conserva su «Verifica» y su «Evidencia» (traza.py: 0 diferencias) y está en su §5.
- `grep -n "/tmp"` en los 9 archivos, quitando `datos/tmp`: 0 líneas.
- Criterios tocados, con su línea (`grep -n "^\*\*ID ("`): U1.md:30 (U1-04), 34 (U1-05), 58 (U1-11), 74 (U1-15), 78 (U1-16); U2.md:21 (U2-02), 45 (U2-08), 49 (U2-09), 57 (U2-11); U5.md:45 (U5-09); C3.md:42 (C3-08); C2.md:38 (C2-07); C1.md:21 (C1-01), 25 (C1-02), 29 (C1-03), 41 (C1-06), 45 (C1-07), 65 (C1-12), 69 (C1-13), 77 (C1-15), 81 (C1-16), 85 (C1-17), 89 (C1-18), 93 (C1-19), 97 (C1-20), 101 (C1-21).

## Supuestos
1. Los encabezados de ocho archivos ya decían «Versión 1.1 (10-oct-2026)» por la ronda 1 de WF-0 (U2 decía «Versión 1»). La orden pide «Versión 1.1»: se dejó 1.1 con fecha 11-oct-2026 y la lista de cambios; el encabezado conserva el origen (ronda 1 de WF-0). Si el integrador prefiere 1.2 para esos ocho, es un cambio de una línea por archivo.
2. Cambios de consistencia no enumerados en la orden, derivados de §K: C1-13 pasa a «C1 v1.1» (el generador produce v1.1); C1-12 recoge `SameSite=Strict` (K6) y K11; C1-06 y C1-07 agregan K8; C1-15 extiende la regla a K1, K5 y K12; U1, C3 y U5 ajustan «6 flujos» a «al menos 6» y la regla K10 en sus apartados de umbrales; U4 solo lleva una nota de pendientes en su encabezado.
3. C1-21 (extractos) se marcó **alta**: sin ellos los agentes no pueden leer el contrato (protocolo §4). El tope «80 KB» se interpretó como 81 920 bytes; es convención del directorio (no está en K32) y así se dice en §5.
4. U1-16 no fija qué semana se muestra los domingos: esa regla (supuesto 4 del reporte de U1 de WF-0) no está en `decisiones.md`. El criterio pide una sola definición, `period` explícito y, si la regla de los domingos no está decidida, marca PROVISIONAL con bloqueo.
5. Los nombres en inglés de `scope`, del campo del borrador, de los decimales (`rate_decimals`) y de las operaciones nuevas no se fijaron en los criterios; los fija C1 (el criterio exige que existan y sean únicos). Solo se citaron `scope`, `Idempotency-Key` y `captured_at`, que aparecen en §K y en los bloqueos.
6. Admitir 135° para el ingreso (R2-07) es convención propia del directorio: `decisiones.md` no menciona 135°. Se dijo así en la §5 de U2.
7. La comprobación de citas es textual salvo normalización de mayúsculas, `**`, comillas inversas y « »; no valida citas sin « ».

## Bloqueos
Ninguno nuevo. `docs/estado/bloqueos/criterios-toques.md` y `criterios-deltal.md` ya traen «RESUELTO» con la respuesta de Axel y quedan aplicados en U1-04, U5-09, C3-08 y U2-09.

## Pendientes
1. **`INDICE.md`** (integrador): conteos 256 / 196 / 60 y los de la tabla anterior por archivo; fila de versión (1.1 del 11-oct-2026); sacar de «Hoy los 248 criterios» la cifra vieja; actualizar las notas 8 y 9 (colisión error–gasto y ΔL entre series, ya resueltas por K26 y K27) y marcar `criterios-toques` y `criterios-deltal` como resueltos.
2. **Decisiones §K sin criterio nuevo en mi alcance (no estaban en la orden):** U4 (K2 alcance del presupuesto en el campo de monto; K4 confirmación del chat de Ximena en Notificaciones; K12 decimales de tasa que la pantalla respeta; K13 tope por persona sin valor fijo; K16 restricciones del comparador en Créditos con el aviso «faltan restricciones»; K23 objetivo de calorías), U3-05 y U3-10 (quincena 1–15 y 16–último día, K8), U5/C3 (prueba Playwright de la cola «Sin conexión», K9, y su flujo) y U5-07 si se agrega «Sin conexión» a la lista de textos. Rutas: `/home/deepsec/finanzas-hogar/docs/criterios/U4.md`, `U3.md`, `U5.md`, `C3.md`.
3. Reejecutar `verify_c1.py --db` y `build.py --check` cuando exista C1 v1.1 (los criterios C1-16 a C1-21 dependen de la regeneración y de `build.py --extractos`, que aún no existe en `backend/internal/api/c1gen/build.py`).

## Solicitudes fuera de alcance
- `/home/deepsec/finanzas-hogar/docs/criterios/INDICE.md`: reconciliar con mis conteos (Pendiente 1). Motivo: la orden lo reserva al integrador.
- `/home/deepsec/finanzas-hogar/docs/criterios/S1.md` (4 líneas), `S2.md`, `S3.md`, `S4.md`, `S5.md`, `S6.md` (2), `S7.md` y `WF4-operacion.md`: reemplazar rutas `/tmp` por `/home/deepsec/finanzas-hogar/datos/tmp` (R2-07). Motivo: la regla de la puerta 1 es «toda ruta /tmp pasa a datos/tmp» y esos archivos no son de mi alcance.
- `/home/deepsec/finanzas-hogar/docs/criterios/S1.md` a `S7.md`: R2-07 también pide admitir ejemplos «ilustrativos» marcados (S1-03) y citar `verify_c1.py --db` y `build.py --check` donde aplique (S6-03); no los toqué.

## Paquetes instalados
Ninguno. Los scripts corrieron con `python3 -I` y la biblioteca estándar.

## Consumo aproximado
- Acciones de herramienta: 34 (con este reporte). Leído: unos 165 KB (protocolo, §K, 9 archivos de criterios completos, colas de 9 bloqueos, 6 líneas de REVISION-r3 y el anexo de scripts). Escrito: unos 85 KB (parches por `str.replace` con una sola aparición, 9 archivos de criterios con +25 KB netos, este reporte y los scripts temporales, que se borraron).

## Anexo · `citas.py`
Se ejecuta con `python3 -I citas.py <archivos .md>` desde `docs/criterios`; imprime los fragmentos de §5 que no están en el plan, las decisiones o el protocolo.
```python
import re,sys,glob,os
B='/home/deepsec/finanzas-hogar/'
def norm(s):
    s=s.replace('«','').replace('»','').replace('*','').replace('`','').replace('“','"').replace('”','"')
    s=re.sub(r'\s+',' ',s).lower()
    return s
corp=''
for f in ['docs/plan/plan.md','docs/plan/decisiones.md','docs/plan/protocolo-agentes.md']:
    corp+=norm(open(B+f,encoding='utf-8').read())+' '
def quotes(line):
    out=[];st=[]
    for i,ch in enumerate(line):
        if ch=='«': st.append(i)
        elif ch=='»' and st:
            j=st.pop(); out.append(line[j+1:i])
    return out
def sec5(t):
    m=re.search(r'^## 5\. .*?$',t,re.M); return t[m.end():]
mal=0;tot=0
for f in sys.argv[1:]:
    t=open(f,encoding='utf-8').read()
    for ln,line in enumerate(sec5(t).split('\n'),1):
        for q in quotes(line):
            q=q.replace('«','').replace('»','')
            parts=[p.strip(' .,;:') for p in re.split(r'…|\.\.\.',q)]
            for p in parts:
                if len(p)<4: continue
                tot+=1
                if norm(p) not in corp:
                    mal+=1; print(os.path.basename(f),'|',p[:140])
print('fragmentos',tot,'no textuales',mal)
```
