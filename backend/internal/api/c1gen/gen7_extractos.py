# Generador C1 · parte 7: extractos por etiqueta de docs/contratos/openapi.yaml (protocolo-agentes §4, K32).
# Un archivo por etiqueta de la barra lateral con SOLO sus rutas y los esquemas que esas rutas usan (resueltos de forma
# transitiva y copiados en su components), más componentes-comunes.yaml e INDICE.md. Tope: 80 KB por archivo; una etiqueta
# que no cabe se divide en partes consecutivas (-a, -b, ...) por orden de rutas. Cada operación queda en un solo extracto.
import json, os, re
from c1paths import OUT, WORK
from gen1_base import emit
LIMIT = 80 * 1000  # 80 KB contados en bytes decimales: cumple también la lectura de 80 KiB
DOC = json.load(open(os.path.join(WORK, 'openapi.json'), encoding='utf-8'))
COMP = DOC['components']
DIR = os.path.join(OUT, 'docs/contratos/extractos')
COMMON = 'componentes-comunes.yaml'
SLUGS = {'Inicio': 'inicio', 'Gastos': 'gastos', 'Presupuesto': 'presupuesto', 'Comidas': 'comidas', 'Compras y tiendas': 'compras-tiendas',
         'Créditos': 'creditos', 'Salud': 'salud', 'Reportes': 'reportes', 'Ajustes': 'ajustes'}
COMMON_SCHEMAS = ['Error', 'Page', 'MoneyCents', 'Date', 'DateTime']
REF = re.compile(r'^#/components/(\w+)/([^/]+)$')
VER = DOC['info']['version']


def refs_in(o, acc):
    if isinstance(o, dict):
        r = o.get('$ref')
        if isinstance(r, str) and REF.match(r):
            m = REF.match(r)
            acc.add((m.group(1), m.group(2)))
        for v in o.values():
            refs_in(v, acc)
    elif isinstance(o, list):
        for v in o:
            refs_in(v, acc)
    return acc


def closure(seed):
    seen, todo = set(), list(seed)
    while todo:
        k = todo.pop()
        if k in seen:
            continue
        seen.add(k)
        todo += list(refs_in(COMP[k[0]][k[1]], set()))
    return seen


COMMON_SET = closure({('parameters', n) for n in COMP['parameters']} | {('headers', n) for n in COMP['headers']}
                     | {('responses', n) for n in COMP['responses']} | {('schemas', n) for n in COMMON_SCHEMAS})


def rewrite(o):
    # Las referencias a componentes comunes apuntan al archivo de comunes; las demás quedan internas al extracto.
    if isinstance(o, dict):
        out = {}
        for k, v in o.items():
            m = REF.match(v) if k == '$ref' and isinstance(v, str) else None
            out[k] = COMMON + v if m and (m.group(1), m.group(2)) in COMMON_SET else rewrite(v)
        return out
    if isinstance(o, list):
        return [rewrite(v) for v in o]
    return o


def head(what):
    return (f'# Finanzas Hogar · Contrato C1 v{VER} · {what}. Archivo generado por backend/internal/api/c1gen/build.py: no se edita a mano.\n'
            '# Fuente única: docs/contratos/openapi.yaml. Índice y regla de uso: docs/contratos/extractos/INDICE.md.\n')


def render_extract(tag, label, paths):
    used = closure(refs_in(paths, set()))
    own = {k for k in used if k not in COMMON_SET}
    assert all(k[0] == 'schemas' for k in own), own
    schemas = {n: rewrite(s) for n, s in COMP['schemas'].items() if ('schemas', n) in own}
    comps = {'securitySchemes': COMP['securitySchemes']}
    if schemas:
        comps['schemas'] = schemas
    doc = {'openapi': DOC['openapi'],
           'info': {'title': f"{DOC['info']['title']} · extracto {label}", 'version': VER,
                    'description': f'Extracto de la etiqueta «{tag}» del contrato C1 (docs/contratos/openapi.yaml): sus rutas y los esquemas que usan. Parámetros, encabezados, respuestas y esquemas base viven en {COMMON}.'},
           'security': DOC['security'], 'tags': [t for t in DOC['tags'] if t['name'] == tag], 'paths': rewrite(paths), 'components': comps}
    return head(f'extracto «{label}»') + '\n'.join(emit(doc)) + '\n', len(schemas)


bytag = {t['name']: [] for t in DOC['tags']}
for p, item in DOC['paths'].items():
    for m, o in item.items():
        bytag[o['tags'][0]].append((p, m, o))
files = []  # (archivo, etiqueta, operaciones, rutas, esquemas propios, texto)
for tag, ops in bytag.items():
    slug = SLUGS[tag]
    units = []
    for p, m, o in ops:
        if not units or units[-1][0] != p:
            units.append((p, {}))
        units[-1][1][m] = o
    for k in range(1, len(units) + 1):
        chunks = [units[i * len(units) // k:(i + 1) * len(units) // k] for i in range(k)]
        names = [f'{slug}.yaml'] if k == 1 else [f'{slug}-{chr(97 + i)}.yaml' for i in range(k)]
        out = []
        for name, ch in zip(names, chunks):
            txt, nsch = render_extract(tag, name[:-5], {p: ms for p, ms in ch})
            out.append((name, tag, sum(len(ms) for _, ms in ch), len(ch), nsch, txt, [o['operationId'] for _, ms in ch for o in ms.values()]))
        if all(len(x[5].encode('utf-8')) <= LIMIT for x in out):
            files += out
            break
    else:
        raise SystemExit(f'gen7: la etiqueta {tag} no cabe en {LIMIT} bytes ni con una ruta por archivo')
common_doc = {'openapi': DOC['openapi'],
              'info': {'title': f"{DOC['info']['title']} · componentes comunes", 'version': VER,
                       'description': 'Componentes que comparten todos los extractos: parámetros, encabezados, respuestas, seguridad y los esquemas base (Error, Page, MoneyCents, Date, DateTime y los que ellos referencian).'},
              'paths': {},
              'components': {'securitySchemes': COMP['securitySchemes'], 'parameters': COMP['parameters'], 'headers': COMP['headers'], 'responses': COMP['responses'],
                             'schemas': {n: s for n, s in COMP['schemas'].items() if ('schemas', n) in COMMON_SET}}}
common_txt = head('componentes comunes') + '\n'.join(emit(common_doc)) + '\n'
assert len(common_txt.encode('utf-8')) <= LIMIT, 'componentes-comunes.yaml pasa de 80 KB'
os.makedirs(DIR, exist_ok=True)
for f in os.listdir(DIR):
    if f.endswith('.yaml') or f == 'INDICE.md':
        os.remove(os.path.join(DIR, f))
for name, tag, nop, npath, nsch, txt, ids in files:
    open(os.path.join(DIR, name), 'w', encoding='utf-8').write(txt)
open(os.path.join(DIR, COMMON), 'w', encoding='utf-8').write(common_txt)


def kb(n):
    return f'{n / 1024:.1f} KB ({n} bytes)'


nops = sum(x[2] for x in files)
total_ops = sum(len(v) for v in DOC['paths'].values())
assert nops == total_ops, (nops, total_ops)
split = {}
for name, tag, *_ in files:
    split.setdefault(tag, []).append(name)
divididas = [f'«{t}» en {", ".join("`" + n + "`" for n in ns)}' for t, ns in split.items() if len(ns) > 1]
csch = sum(1 for k in COMMON_SET if k[0] == 'schemas')
lines = ['# Extractos del contrato C1 · índice', '',
         'Generado por `backend/internal/api/c1gen/build.py` (sin opciones o con `--extractos`); `build.py --check` los compara byte a byte con el generador. No se editan a mano.', '',
         '**Regla (protocolo-agentes §4, K32):** estos archivos son los **únicos** archivos del contrato que leen los agentes de WF-1 y WF-3. '
         '`docs/contratos/openapi.yaml` es la fuente única y no se lee entero; si falta algo que no está en un extracto, se pide al integrador.', '',
         f'- Contrato C1 v{VER}: {total_ops} operaciones en {len(DOC["paths"])} rutas y {len(COMP["schemas"])} esquemas.',
         '- Cada operación aparece en exactamente un extracto. Cada extracto trae sus rutas y solo los esquemas que esas rutas usan, resueltos de forma transitiva y copiados en su `components`.',
         f'- Parámetros, encabezados, respuestas, seguridad y los esquemas `Error`, `Page`, `MoneyCents`, `Date`, `DateTime` (y los que ellos referencian) están en `{COMMON}`; los extractos los citan con `$ref: {COMMON}#/components/...`.',
         f'- Tope por archivo: 80 KB (80 000 bytes). Etiquetas divididas en partes consecutivas por orden de rutas: {"; ".join(divididas) if divididas else "ninguna"}.', '',
         '| Archivo | Etiqueta | Operaciones | Rutas | Esquemas propios | Tamaño |', '|---|---|---|---|---|---|']
for name, tag, nop, npath, nsch, txt, ids in files:
    lines.append(f'| `{name}` | {tag} | {nop} | {npath} | {nsch} | {kb(len(txt.encode("utf-8")))} |')
lines.append(f'| `{COMMON}` | (comunes) | — | — | {csch} esquemas, {len(COMP["parameters"])} parámetros, {len(COMP["headers"])} encabezados, {len(COMP["responses"])} respuestas | {kb(len(common_txt.encode("utf-8")))} |')
lines.append(f'| **Total** | | **{nops}** | | | |')
lines += ['', '## Operaciones por extracto', '']
for name, tag, nop, npath, nsch, txt, ids in files:
    lines.append(f'- `{name}`: ' + ', '.join(f'`{i}`' for i in ids))
open(os.path.join(DIR, 'INDICE.md'), 'w', encoding='utf-8').write('\n'.join(lines) + '\n')
print(f'escrito docs/contratos/extractos: {len(files)} extractos + {COMMON} + INDICE.md, {nops} operaciones, mayor {max(len(x[5].encode("utf-8")) for x in files)} bytes')
