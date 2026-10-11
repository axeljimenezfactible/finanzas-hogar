# Generador C1 · parte 1: emisor YAML, componentes base y esquemas de entidad desde S1
import json, re, os
from c1paths import ROOT, WORK
SCHEMA = json.load(open(os.path.join(WORK, 'schema.json'), encoding='utf-8'))

# ---------- emisor YAML (estilo bloque, cadenas entre comillas dobles) ----------
PLAIN_KEY = re.compile(r'^[A-Za-z_$/][A-Za-z0-9_./{}$-]*$')
def ykey(k):
    k = str(k)
    return k if PLAIN_KEY.match(k) and k not in ('true', 'false', 'null', 'yes', 'no', 'on', 'off') else json.dumps(k, ensure_ascii=False)
def yscalar(v):
    if v is True: return 'true'
    if v is False: return 'false'
    if v is None: return 'null'
    if isinstance(v, (int, float)): return str(v)
    return json.dumps(v, ensure_ascii=False)
def emit(o, ind=0):
    out = []; sp = ' ' * ind
    if isinstance(o, dict):
        for k, v in o.items():
            if isinstance(v, dict) and v:
                out.append(f'{sp}{ykey(k)}:'); out += emit(v, ind + 2)
            elif isinstance(v, list) and v:
                out.append(f'{sp}{ykey(k)}:'); out += emit(v, ind + 2)
            elif isinstance(v, dict): out.append(f'{sp}{ykey(k)}: {{}}')
            elif isinstance(v, list): out.append(f'{sp}{ykey(k)}: []')
            else: out.append(f'{sp}{ykey(k)}: {yscalar(v)}')
    elif isinstance(o, list):
        for it in o:
            if isinstance(it, (dict, list)) and it:
                sub = emit(it, ind + 2)
                out.append(f'{sp}- {sub[0].lstrip()}'); out += sub[1:]
            else: out.append(f'{sp}- {yscalar(it) if not isinstance(it,(dict,list)) else ("{}" if isinstance(it,dict) else "[]")}')
    return out

def ref(name, kind='schemas'): return {'$ref': f'#/components/{kind}/{name}'}

# ---------- tipos base ----------
BASE = {
 'MoneyCents': {'type': 'integer', 'description': 'Importe exacto en centavos de MXN (entero, nunca punto flotante). 97375 = $973.75. Redondeo ROUND_HALF_UP; en una serie de pagos el último ajusta centavos. En pantalla: `$` y dos decimales.', 'examples': [97375]},
 'Date': {'type': 'string', 'format': 'date', 'pattern': '^[0-9]{4}-[0-9]{2}-[0-9]{2}$', 'description': 'Fecha de negocio ISO 8601 (AAAA-MM-DD) del calendario de America/Mexico_City. En pantalla se muestra 10-oct-2026.', 'examples': ['2026-10-30']},
 'DateTime': {'type': 'string', 'format': 'date-time', 'pattern': '^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\\.[0-9]{1,6})?(Z|[+-][0-9]{2}:[0-9]{2})$', 'description': 'Instante ISO 8601 con desplazamiento explícito de America/Mexico_City (hoy -06:00). La base lo guarda en UTC con «Z» y la API lo convierte al responder.', 'examples': ['2026-10-10T21:30:00-06:00']},
 'BasisPoints': {'type': 'integer', 'minimum': 0, 'maximum': 10000, 'description': 'Porcentaje en puntos base: 10000 = 100 %, 5000 = 50 %.', 'examples': [5000]},
 'Milli': {'type': 'integer', 'description': 'Cantidad en milésimas de la unidad base (g, ml o pieza): 1000 = 1 g, 1 ml o 1 pieza; en porciones, 1000 = 1 porción.', 'examples': [1000]},
 'Month': {'type': 'string', 'pattern': '^[0-9]{4}-(0[1-9]|1[0-2])$', 'description': 'Mes calendario AAAA-MM.', 'examples': ['2026-10']},
 'Error': {'type': 'object', 'additionalProperties': False, 'required': ['code', 'message', 'request_id'],
   'description': 'Esquema único de error de toda la API (C1 §3). Lo usan todas las respuestas 4xx y 5xx.',
   'properties': {
     'code': {'type': 'string', 'pattern': '^[a-z][a-z0-9_]*$', 'description': 'Código estable en inglés (tabla de códigos en C1-datos.md §3).', 'examples': ['validation_failed']},
     'message': {'type': 'string', 'description': 'Mensaje para la persona usuaria, en español de México.', 'examples': ['Revisa los campos marcados.']},
     'fields': {'type': 'array', 'description': 'Detalle por campo cuando aplica (validación).', 'items': ref('FieldError')},
     'request_id': {'type': 'string', 'description': 'Identificador de la petición; es el mismo valor que el encabezado X-Request-Id y el de audit_log.request_id.', 'examples': ['01JABCDEF0123456789']},
     'retry_at': {'$ref': '#/components/schemas/DateTime', 'description': 'Solo en 429 (límite de intentos de inicio de sesión): hora a partir de la cual se puede intentar de nuevo.'}}},
 'FieldError': {'type': 'object', 'additionalProperties': False, 'required': ['field', 'code', 'message'], 'properties': {
     'field': {'type': 'string', 'description': 'Nombre de la propiedad (nombre_codigo) o ruta JSON del campo con error.', 'examples': ['amount_cents']},
     'code': {'type': 'string', 'pattern': '^[a-z][a-z0-9_]*$', 'description': 'Motivo estable en inglés: required, invalid_format, out_of_range, not_found, not_allowed, conflict.', 'examples': ['out_of_range']},
     'message': {'type': 'string', 'description': 'Explicación en español de México.', 'examples': ['El monto debe ser mayor que cero.']}}},
 'Page': {'type': 'object', 'required': ['items', 'next_cursor'], 'description': 'Envoltura única de toda colección paginada (C1 §5). `next_cursor` = null indica que no hay más.', 'properties': {
     'items': {'type': 'array', 'description': 'Registros de la página, en el orden pedido.'},
     'next_cursor': {'type': ['string', 'null'], 'description': 'Cursor opaco para pedir la página siguiente; null cuando no hay más registros.'}}},
}

# ---------- mapeo columna S1 -> JSON Schema ----------
DATE_COLS = {'valid_from', 'valid_to'}
def col_schema(t, c, nullable_ok=False):
    n = c['name']; desc = (c['comment'].split('|')[0].strip() or n)
    if c['enum'] and set(c['enum']) == {'0', '1'}: s = {'type': 'boolean'}
    elif n.endswith('_cents'): s = {'$ref': '#/components/schemas/MoneyCents'}
    elif n.endswith('_date') or n in DATE_COLS: s = {'$ref': '#/components/schemas/Date'}
    elif n.endswith('_at'): s = {'$ref': '#/components/schemas/DateTime'}
    elif n.endswith('_bp'): s = {'$ref': '#/components/schemas/BasisPoints'}
    elif n.endswith('_milli') and not n.startswith('kcal'): s = {'$ref': '#/components/schemas/Milli'}
    elif c['enum']:
        if c['type'] == 'INTEGER': s = {'type': 'integer', 'enum': [int(v) for v in c['enum']]}
        else: s = {'type': 'string', 'enum': c['enum']}
    elif c['type'] == 'INTEGER': s = {'type': 'integer'}
    else: s = {'type': 'string'}
    if n == 'params_json': s = {'type': 'object', 'additionalProperties': True}
    if n == 'value_json': s = {}
    s = dict(s); s['description'] = desc
    if 'PROVISIONAL' in c['comment']:
        # El bloqueo se toma del propio comentario de S1 (antes se citaba siempre S1-meta-kcal, corrección r2).
        s['x-provisional'] = True
        bl = re.search(r'\b((?:S|C|U)\d+[a-z]?-[a-z0-9]+(?:-[a-z0-9]+)*|criterios-[a-z0-9]+(?:-[a-z0-9]+)*)\b', c['comment'])
        if not ('PROVISIONAL' in desc and bl and bl.group(1) in desc):
            s['description'] += f" (PROVISIONAL: ver docs/estado/bloqueos/{bl.group(1) if bl else 'S1-meta-kcal'}.md)"
    if nullable_ok and not c['notnull']:
        extra = {k: s[k] for k in ('x-provisional',) if k in s}
        d2 = s['description'] + ' · null borra el valor'
        inner = {k: v for k, v in s.items() if k not in ('description', 'x-provisional')}
        if '$ref' in inner: s = {'anyOf': [inner, {'type': 'null'}], 'description': d2, **extra}
        elif 'type' in inner: s['type'] = [inner['type'], 'null']; s['description'] = d2
        if 'enum' in s and isinstance(s.get('type'), list): s['enum'] = list(s['enum']) + [None]
    return s

HIDDEN = {'user_account': {'password_hash', 'password_changed_at'}, 'device_session': {'token_hash', 'ip_address'}, 'attachment': {'file_path'},
          'generated_report': {'file_path'}, 'notification_recipient': {'chat_id'}}
SYSTEM_MANAGED = {'id', 'created_at', 'updated_at'}
def pascal(t): return ''.join(p.capitalize() for p in t.split('_'))
def cols(t): return [c for c in SCHEMA[t]['cols']]
def entity_read(t, computed=None, hide=()):
    props = {}; req = []
    for c in cols(t):
        if c['name'] in HIDDEN.get(t, set()) or c['name'] in hide: continue
        props[c['name']] = col_schema(t, c)
        if c['notnull'] or c['pk']: req.append(c['name'])
    for k, v in (computed or {}).items():
        v = dict(v); v['x-calculado'] = True; props[k] = v
    m = SCHEMA[t]['meta']
    return {'type': 'object', 'description': f"{m['title']} (tabla `{t}`). {m['desc']} Campos ausentes = sin valor (NULL).", 'x-tabla': t, 'required': req, 'properties': props}
def entity_write(t, include=None, exclude=(), required=None, update=False, extra=None, desc=''):
    props = {}; req = []
    for c in cols(t):
        n = c['name']
        if n in SYSTEM_MANAGED or n in exclude or (include is not None and n not in include): continue
        props[n] = col_schema(t, c, nullable_ok=update)
        if not update and c['notnull'] and c['default'] is None: req.append(n)
    for k, v in (extra or {}).items(): props[k] = v
    if required is not None: req = required
    m = SCHEMA[t]['meta']
    s = {'type': 'object', 'additionalProperties': False, 'description': (desc or ('Cambios parciales (JSON Merge Patch, RFC 7396) de ' if update else 'Alta de ') + f"{m['title'].lower()} (tabla `{t}`)."), 'x-tabla': t, 'properties': props}
    if req and not update: s['required'] = req
    if update: s['minProperties'] = 1
    return s
