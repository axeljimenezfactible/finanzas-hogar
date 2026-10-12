# Generador C1 · parte 8: backend/internal/api/types_gen.go desde openapi.json
#   - un tipo Go por esquema de components/schemas (structs con etiquetas json, alias para los escalares)
#   - Operations: lo que el contrato declara de cada operación y que usa el servidor (sesión, idempotencia, cuerpo, códigos 2xx)
# Reglas de traducción (CONVENCIONES-backend.md §3): requerido → valor; opcional → puntero u omitzero; tipo con null → Nullable[T].
# La alineación de los structs reproduce la de gofmt (columnas con un espacio de separación), así gofmt -l no lista el archivo.
import json, os, re
from c1paths import OUT, WORK
d = json.load(open(os.path.join(WORK, 'openapi.json'), encoding='utf-8'))
S = d['components']['schemas']
SKIP = {'Error': 'ErrorBody', 'FieldError': 'FieldError'}  # ya declarados en routes.go (gen5_go.py)
INIT = {'id': 'ID', 'ids': 'IDs', 'url': 'URL', 'api': 'API', 'json': 'JSON', 'sha256': 'SHA256', 'http': 'HTTP', 'ip': 'IP',
        'pdf': 'PDF', 'ocr': 'OCR', 'ai': 'AI', 'csv': 'CSV', 'bmi': 'BMI', 'mb': 'MB', 'gb': 'GB'}
SCALAR = {'integer': 'int64', 'string': 'string', 'boolean': 'bool', 'number': 'float64'}


def goname(s):
    parts = [p for p in re.split(r'[^A-Za-z0-9]+', s) if p]
    n = ''.join(INIT.get(p.lower(), p[:1].upper() + p[1:]) for p in parts)
    return 'N' + n if n[:1].isdigit() else n


names = {SKIP.get(k, goname(k)) for k in S}
out_types = []  # (name, kind, payload)
uses_raw = [False]


def unique(name):
    base, i = name, 2
    while name in names:
        name, i = f'{base}{i}', i + 1
    names.add(name)
    return name


def resolve(s):
    while '$ref' in s:
        s = S[s['$ref'].split('/')[-1]]
    return s


def merged(s):
    props, req = {}, []
    for part in s.get('allOf', []):
        p, r = merged(resolve(part))
        props.update(p)
        req += r
    props.update(s.get('properties', {}))
    req += s.get('required', [])
    return props, req


def raw():
    uses_raw[0] = True
    return 'json.RawMessage'


def gotype(s, ctx):
    """Devuelve (tipo Go, admite null)."""
    if '$ref' in s:
        n = s['$ref'].split('/')[-1]
        return SKIP.get(n, goname(n)), False
    if 'anyOf' in s or 'oneOf' in s:
        alts = s.get('anyOf') or s.get('oneOf')
        nn = [a for a in alts if a.get('type') != 'null']
        if len(nn) == 1:
            return gotype(nn[0], ctx)[0], len(nn) < len(alts)
        return raw(), False
    t = s.get('type')
    nullable = isinstance(t, list) and 'null' in t
    if isinstance(t, list):
        base = [x for x in t if x != 'null']
        t = base[0] if len(base) == 1 else None
    if t is None and ('allOf' in s or 'properties' in s):
        t = 'object'
    if t is None and 'const' in s:
        t = {str: 'string', bool: 'boolean', int: 'integer'}.get(type(s['const']))
    if t in SCALAR:
        return SCALAR[t], nullable
    if t == 'array':
        if 'items' not in s:
            return '[]' + raw(), nullable
        return '[]' + gotype(s['items'], ctx + 'Item')[0], nullable
    if t == 'object':
        if s.get('properties') or s.get('allOf'):
            name = unique(ctx)
            emit_struct(name, s, f'objeto anidado de {ctx}')
            return name, nullable
        ap = s.get('additionalProperties')
        if isinstance(ap, dict) and ap:
            return 'map[string]' + gotype(ap, ctx + 'Value')[0], nullable
        return 'map[string]any', nullable
    return raw(), nullable


def emit_struct(name, s, note):
    props, req = merged(s)
    slot = len(out_types)
    out_types.append(None)
    rows, seen = [], set()
    for p, ps in props.items():
        f = goname(p) or 'Field'
        while f in seen:
            f += '_'
        seen.add(f)
        t, nullable = gotype(ps, name + f)
        required = p in req
        if nullable:
            t, tag = f'Nullable[{t}]', p if required else p + ',omitzero'
        elif required:
            tag = p
        else:
            if not (t.startswith('[]') or t.startswith('map[') or t == 'json.RawMessage'):
                t = '*' + t
            tag = p + ',omitzero'
        rows.append((f, t, f'`json:"{tag}"`'))
    out_types[slot] = (name, 'struct', (rows, note))


for key in sorted(S):
    if key in SKIP:
        continue
    s, name = S[key], goname(key)
    if s.get('type') == 'object' or 'allOf' in s or 'properties' in s:
        if s.get('properties') or s.get('allOf'):
            emit_struct(name, s, f'esquema {key} de C1')
            continue
    t, _ = gotype({k: v for k, v in s.items() if k != '$ref'}, name + 'Value') if '$ref' not in s else gotype(s, name)
    out_types.append((name, 'alias', (t, f'esquema {key} de C1')))


def fmt_struct(name, rows, note):
    lines = [f'// {name} es el {note}.', f'type {name} struct {{']
    if rows:
        w1 = max(len(r[0]) for r in rows) + 1
        w2 = max(len(r[1]) for r in rows) + 1
        lines += [f'\t{f.ljust(w1)}{t.ljust(w2)}{tag}' for f, t, tag in rows]
    lines.append('}')
    return '\n'.join(lines)


ops = []
for path, items in d['paths'].items():
    for m, o in items.items():
        if not isinstance(o, dict) or 'operationId' not in o:
            continue
        sec = o.get('security')
        auth = 'public' if sec == [] else 'service' if sec and any('serviceToken' in x for x in sec) else 'session'
        idem = any(p.get('$ref', '').endswith('/IdempotencyKey') for p in o.get('parameters', []))
        content = (o.get('requestBody') or {}).get('content', {})
        ctype = next(iter(content), '')
        sch = (content.get(ctype) or {}).get('schema', {}) if ctype else {}
        sname = sch.get('$ref', '').split('/')[-1] if '$ref' in sch else ''
        codes = sorted(int(c) for c in o.get('responses', {}) if c.isdigit() and 200 <= int(c) < 300)
        ops.append((o['operationId'], auth, idem, ctype, sname, codes))
ops.sort()

body = []
for name, kind, payload in out_types:
    if kind == 'struct':
        body.append(fmt_struct(name, *payload))
    else:
        body.append(f'// {name} es el {payload[1]}.\ntype {name} = {payload[0]}')
oplines = [f'\t{{OperationID: "{oid}", Auth: "{a}", Idempotent: {"true" if i else "false"}, RequestType: "{c}", RequestSchema: "{sn}", Statuses: []int{{{", ".join(map(str, cs))}}}}},'
           for oid, a, i, c, sn, cs in ops]
nstruct = sum(1 for t in out_types if t[1] == 'struct')
src = ('// Code generated by backend/internal/api/c1gen (gen8_types.py) from docs/contratos/openapi.yaml (contrato C1 v' + d['info']['version'] + '). DO NOT EDIT.\n'
       '// Regenerar con «python3 -I backend/internal/api/c1gen/build.py» desde la raíz del repositorio; «build.py --check» lo compara byte a byte.\n'
       f'// {len(out_types)} tipos ({nstruct} structs) de components/schemas y {len(ops)} operaciones. Error y FieldError están en routes.go (ErrorBody, FieldError).\n\n'
       'package api\n\n' + ('import "encoding/json"\n\n' if uses_raw[0] else '') +
       '// OperationMeta resume lo que el contrato declara de una operación y que usa el servidor.\n'
       'type OperationMeta struct {\n'
       '\tOperationID   string\n'
       '\tAuth          string // "session" (cookie fh_session), "public" (sin sesión) o "service" (Authorization: Bearer del agente externo, K53)\n'
       '\tIdempotent    bool   // exige el encabezado Idempotency-Key\n'
       '\tRequestType   string // tipo de contenido del cuerpo: "", "application/json" o "multipart/form-data"\n'
       '\tRequestSchema string // esquema del cuerpo en components/schemas, o ""\n'
       '\tStatuses      []int  // códigos 2xx declarados\n'
       '}\n\n'
       f'// Operations es la lista de las {len(ops)} operaciones de docs/contratos/openapi.yaml con sus metadatos.\n'
       'var Operations = []OperationMeta{\n' + '\n'.join(oplines) + '\n}\n\n' + '\n\n'.join(body) + '\n')
dst = os.path.join(OUT, 'backend/internal/api/types_gen.go')
os.makedirs(os.path.dirname(dst), exist_ok=True)
open(dst, 'w', encoding='utf-8').write(src)
print(f'types_gen.go: {len(out_types)} tipos ({nstruct} structs) y {len(ops)} operaciones → {dst}')
