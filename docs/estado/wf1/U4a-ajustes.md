# U4a-ajustes · reporte final (WF-1)

Agente U4a-ajustes (Opus), 11-oct-2026. Alcance: `docs/design/U4-ajustes-a.md` (secciones P-26.01 a P-26.07 de Ajustes) y este reporte.

## Hecho
- Especificación estilo Notion de las 7 secciones de la parte a: propósito, tabla de campos (nombre, tipo, valor por defecto con origen, validación → mensaje, `operationId` de lectura y escritura, N-xxx), vistas de base de datos, edición en el lugar, vista previa del efecto, estados, acciones críticas, atajos, móvil (412 px) y escritorio (1440 px).
- §0 de convenciones comunes (tabla de campos, guardado, estados, acción crítica con tokens `color.critical.*`, patrón K44 con `parameters`/`reset_parameters`, catálogo de 22 mensajes es-MX, barra lateral, vistas, menú «/» y atajos solo de U1 §5.3).
- Campos nuevos de la puerta 1 en la parte a: alcance K2, interruptor de proyección K21, interruptor A/B con vista previa del oráculo, aviso «pendiente confirmar», tasas como se capturaron K12, datos K18 con «falta capturar», restricciones duras K16 sin valor con «faltan restricciones», ventana ± días sin valor, pesos iguales, tope por persona (solo lectura, se edita en P-26.12), promociones, mercados, membresía de Sam's con punto de equilibrio, lectura asistida bloqueada sin veredicto, estado de la IA (OCR y recetas) con enlace a P-26.11.
- 8 parámetros K44 en su sección (Cuenta y roles 4, Créditos 1, Tiendas y precios 3) con valor inicial, validación y «restablecer».
- Bloqueo nuevo: `docs/estado/bloqueos/U4-iva-credito.md` (C1 no distingue «IVA no capturado» de «incluye IVA»).

## Archivos
- `/home/deepsec/finanzas-hogar/docs/design/U4-ajustes-a.md` (nuevo, 54.7 KB)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U4-iva-credito.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/U4a-ajustes.md` (este reporte)

## Verificación
Script `python3 -I` (temporal en `datos/tmp/u4a/verify.py`, borrado al terminar; copia abajo). Salida:
```
U4-02 operationIds citados: 99 | inexistentes: []
U4-02 ops etiqueta Ajustes: 73 | sin mención en parte a: []
U4-09 hex: 0 | rgb(: 0
U4-11 fechas ISO: []
U4-01 secciones ## P-26.: 7
U4-01/U4-12 filas de tabla de 6-7 columnas: 155 | celdas vacías: 0
U4-12 claves parte a esperadas: 8 | en tabla: 8 | sin clave o distintas: [] | claves sin @param: []
U4-05 serie A igual al oráculo: True | serie B: True | mín A True | mín B True | n 14
R-64 evento 0: {'disponible_A': 10000.0, 'disponible_B': 6905.5}
U4-06 atajos en §0.7: 11 | no están en U1 §5.3: []
mínimo de seguridad (singular): [385]   (línea de SUPUESTOS)
```
Cifras de ejemplo recalculadas con Decimal y ROUND_HALF_UP: cotización 80050/309450 = 25.8685 %; $3,094.50 × 25.8685 % = $800.50; peso K40 día 30/60/90 = 0.9836/0.4918/0; 10-oct-2026 + 30/60/90 días = 9-nov-2026, 9-dic-2026, 8-ene-2027; $199.00 al 50 % = $99.50.

Criterios de U4 (parte a):
- **U4-01**: 7 de 14 secciones (las 7 restantes son de la parte b); cada campo con tipo, valor por defecto, validación, `operationId` y N-xxx; 0 celdas vacías; encabezados `## P-26.0x`.
- **U4-02**: 0 `operationId` inexistentes; las 73 operaciones de la etiqueta Ajustes figuran con su lugar (36 con campo en la parte a, el resto asignado a la parte b o a excepciones, sección «Cobertura»). El cruce en ambos sentidos sobre el documento unido lo hace el integrador.
- **U4-03** (parte a): Cuenta y roles (co-titulares mismos permisos, sesiones y cierre, Karen sin cuenta) · Presupuesto (mes calendario, quincena subperiodo, alcance con «De este mes en adelante» preseleccionado, interruptor K21) · Ingresos (15 y 30, extras en P-08, Ximena sin ingresos) · Categorías (Comida fuera) · Créditos (4 tipos, 17.27 % y 25.8685 %, tasa por compra en P-19, promociones con comercio/plazo/compra mínima/vigencia, A/B = B con aviso, K12, K16, K18) · Tiendas (mercados con nombre, día y qué se compra; membresía y costo; frescura 30; lectura asistida bloqueada sin veredicto) · Recetas (editables, restricciones vacías) · IA apagada (estado visible). El dato «incluye IVA» de K18 queda PROVISIONAL (bloqueo U4-iva-credito).
- **U4-04** (parte a): B; IA apagada; lectura asistida apagada; 30 y 90 días; 5 en 15 min; 90 días de sesión; tasas 17.27 % y 25.8685 %; alcance «de este mes en adelante»; meses 7 a 12 apagado; restricciones, K18 y costo de membresía sin valor; comparador 5 s; pesos 1, 1, 1, 1 (decisión C).
- **U4-05**: A/B con R-61, R-62, R-64, R-65 y R-66 (verificado contra `oraculo/model.json`); además alcance de presupuesto, proyección, ingreso, reparto ($199.00), tasa ($800.50), restricciones y pesos, frescura, membresía, PROFECO, recarga nutricional.
- **U4-06**: tabla, tablero, calendario y galería usados (sesiones, presupuesto, gastos fijos, categorías, créditos, tasas, promociones con calendario por defecto, tiendas, recetario); menú «/» y 11 atajos de U1 §5.3 sin duplicados ni nuevos; edición en el lugar por sección; claro y oscuro por token.
- **U4-07**: catálogo V-01 a V-22 más 3 mensajes en línea, es-MX; dinero `$1,234.56`; tasas con hasta 4 decimales mostradas como se capturaron.
- **U4-08**: cerrar todas las sesiones (dos alcances), quitar monto, eliminar tasa, eliminar promoción, activar lectura asistida (general y por sitio), eliminar receta: diálogo con `color.critical.700`/`color.critical.300`/`color.critical.950`.
- **U4-09**: 412 px y 1440 px en las 7; 44 px; 0 hex.
- **U4-10**: vacío, carga, error y guardado en las 7, con la regla común de §0.2.
- **U4-11**: 0 fechas ISO; SUPUESTOS presentes; bloqueo U4-iva-credito.
- **U4-12** (parte a): 8 parámetros cruzados con `@param` de 00003 (0 diferencias). Los 7 restantes y los 4 de `notification.params_json` (PROVISIONAL, `S1-parametros-avisos`) son de la parte b.
- **U4-13**: no aplica a la parte a (Salud y metas de ahorro es la sección 13, parte b).

```python
import re, json, glob, sys
doc = open('docs/design/U4-ajustes-a.md', encoding='utf-8').read()
ops, tag = {}, {}
for f in glob.glob('docs/contratos/extractos/*.yaml'):
    for m in re.finditer(r'operationId:\s*"(\w+)"', open(f, encoding='utf-8').read()):
        ops[m.group(1)] = f.split('/')[-1]
aj = {o for o, f in ops.items() if f.startswith('ajustes-')}
cited = set(re.findall(r'`([a-z][A-Za-z]+)`', doc))
cand = {c for c in cited if re.match(r'^(list|get|create|update|delete|revoke|change|reorder|set|import|reload|preview|duplicate|suggest|logout|login|query|mark|send|confirm|dismiss|restore|generate|download|upload|compare)[A-Z]?', c) and c not in ('logout','login') or c in ('logout','login')}
missing = sorted(c for c in cand if c not in ops)
print('U4-02 operationIds citados:', len(cand), '| inexistentes:', missing)
print('U4-02 ops etiqueta Ajustes:', len(aj), '| sin mención en parte a:', sorted(aj - cited))
print('U4-09 hex:', len(re.findall(r'#[0-9A-Fa-f]{3,8}\b', doc)), '| rgb(:', doc.count('rgb('))
print('U4-11 fechas ISO:', re.findall(r'\d{4}-\d{2}-\d{2}', doc))
print('U4-01 secciones ## P-26.:', len(re.findall(r'^## P-26\.', doc, re.M)))
# celdas vacías en tablas de campos (6 columnas)
empty = 0; rows = 0
for line in doc.split('\n'):
    if line.startswith('|') and not line.startswith('|---'):
        cells = [c.strip() for c in line.strip().strip('|').split('|')]
        if len(cells) in (6, 7):
            rows += 1; empty += sum(1 for c in cells if c == '')
print('U4-01/U4-12 filas de tabla de 6-7 columnas:', rows, '| celdas vacías:', empty)
# K44: cruce con @param de 00003
params = {}
for l in open(glob.glob('backend/db/migrations/00003*.sql')[0], encoding='utf-8'):
    m = re.match(r'--\s*@param (\S+) \| [^|]+\| [^|]+\| ([^|]+)\| ([^|]+)\| ([^|]+)\|', l)
    if m: params[m.group(1)] = (m.group(2).strip(), m.group(4).strip())
seca = {'Cuenta y roles', 'Créditos y condiciones', 'Tiendas y precios'}
want = {k for k, v in params.items() if v[1] in seca}
sec = doc.split('## Parámetros operativos de la parte a (K44)')[1].split('\n## ')[0]
got = {}
for line in sec.split('\n'):
    m = re.match(r'\| [^|]+\| ([^|]+)\| `([a-z_.]+)` \| ([^|]+)\|', line)
    if m: got[m.group(2)] = (m.group(3).strip(), m.group(1).strip())
bad = [(k, got.get(k), params[k]) for k in want if k not in got or got[k][1] != params[k][1] or not got[k][0].startswith(params[k][0])]
print('U4-12 claves parte a esperadas:', len(want), '| en tabla:', len(got), '| sin clave o distintas:', bad, '| claves sin @param:', sorted(set(got) - set(params)))
# oráculo A/B
m = json.load(open('oraculo/model.json'))['credito']
fmt = lambda x: '$' + f'{x:,.2f}'
sa = ' · '.join(fmt(v) for v in m['comparacion_series']['esperado_A'])
sb = ' · '.join(fmt(v) for v in m['comparacion_series']['esperado_B'])
print('U4-05 serie A igual al oráculo:', sa in doc, '| serie B:', sb in doc, '| mín A', fmt(m['minimo_A']) in doc, '| mín B', fmt(m['minimo_B']) in doc, '| n', len(m['comparacion_series']['esperado_A']))
ev = m['eventos'][0]; print('R-64 evento 0:', {k: ev[k] for k in ev if 'disp' in k.lower() or k in ('A','B')})
# atajos: todos en U1 §5.3
u1 = open('docs/design/U1-arquitectura.md', encoding='utf-8').read().split('### 5.3')[1].split('### 5.4')[0]
ks = set(re.findall(r'`([^`]+)`', doc.split('**Atajos**')[1].split('\n')[0]))
print('U4-06 atajos en §0.7:', len(ks), '| no están en U1 §5.3:', sorted(k for k in ks if '`'+k+'`' not in u1))
print('mínimo de seguridad (singular):', [i for i, l in enumerate(doc.split('\n'), 1) if 'mínimo de seguridad' in l.lower()])
```

## Supuestos
Los 7 de la sección SUPUESTOS de `docs/design/U4-ajustes-a.md` (despensa programada como `planned_item` sin fijar su `kind`; capital y plazo del préstamo en P-19; OCR local sin interruptor; porciones base solo lectura; «al menos un peso > 0» como regla técnica; fecha de deuda en cero obligatoria según C1 y K16; cifras de ejemplo del oráculo o marcadores).

## Bloqueos
- `docs/estado/bloqueos/U4-iva-credito.md` (nuevo).
- `docs/estado/bloqueos/S1-parametros-avisos.md` (abierto; no afecta a las secciones 1 a 7).

## Pendientes
- Integrador: unir `docs/design/U4-ajustes-a.md` y `U4-ajustes-b.md` en `docs/design/U4-ajustes.md` y correr el cruce U4-02 en ambos sentidos y el conteo U4-01 (14) sobre el documento unido.
- Enlazar los controles a `docs/design/DESIGN.md` §7 cuando U2b-2 la escriba.

## Solicitudes fuera de alcance
1. `docs/contratos/openapi.yaml` (`Credit`, `CreditCreate`, `CreditUpdate`, `missing_terms`) y migración: resolver el bloqueo U4-iva-credito (K18 pide «falta capturar» también para el IVA).
2. `docs/design/U1-arquitectura.md` §3.28 N-288: dice «fecha de deuda en cero (opcional)», pero C1 la lista en `missing_constraints` y K16 exige todas; alinear el texto.
3. `docs/design/U1-arquitectura.md` §2.3 y §2.5.a: «días de visita» en plural; C1 y el esquema guardan un solo `visit_weekday` (decisión D dice «día de visita»); alinear o pedir lista.
4. `docs/design/U1-arquitectura.md` §2.3 P-26.07 lista «Porciones base» como ajuste, pero C1 no tiene operación para cambiar `default_base_servings`; confirmar que es solo lectura o agregar la operación.
5. `docs/contratos/reglas-negocio.md` RN-12 y motor S5: con frescura y peso cero configurables (K44), escribir la fórmula general del peso (propuesta abajo) para que la vista previa de P-26.06 no dependa de 30 y 90.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- (Puerta 2, K44) Validación de dependencia en Tiendas y precios: `price.zero_weight_after_days` mayor que `price.stale_after_days`; si no, el peso de un precio viejo no tiene sentido.
- (Puerta 2, K44) Fórmula general del peso con S = días para precio viejo y Z = días para peso cero: peso = 1 si días < S; (Z − días) ÷ (Z − S + 1) si S ≤ días ≤ Z; no compite después de Z. Con 30 y 90 da exactamente K40 ((90 − días) ÷ 61).
- (Puerta 2, K44; la propuesta formal de mínimos la hace S1) Topes técnicos para no dejar el sistema inservible: `auth.max_failed_attempts` ≥ 3 (el ejemplo de K44); `session.max_age_days` entre 1 y 365 (el ejemplo de K44); `auth.lockout_window_minutes` ≥ 1; `auth.password_min_length` ≥ 12 (bajarlo debilita lo decidido en K6); `comparator.timeout_seconds` ≥ 1 (con 0 el comparador siempre responde «no disponible»); `price.profeco_import_max_gb` ≥ 1 (con 0 no se puede importar).
- (Puerta 2) Bloqueo U4-iva-credito: se recomienda la opción 1 (columna `vat_included` que admite vacío), porque cumple K18 sin cambiar el significado de lo que ya hay.

## Consumo aproximado
- Acciones de herramienta: 44.
- Leído: ~205 KB (protocolo 6, decisiones 26, criterios U4 27, U1 por partes ~45, extractos por operación y esquema ~70, C1-datos §11.5 y §11.7 10, otros ~20).
- Escrito: ~78 KB (documento 54.7, este reporte ~12, bloqueo 1.7, scripts temporales ~8, progreso 0.6).
