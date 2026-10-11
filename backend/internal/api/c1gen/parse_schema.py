import re, json, sqlite3, sys, os
from c1paths import ROOT, WORK
SQL = os.path.join(ROOT, 'backend/db/migrations/00001_esquema.sql')
txt = open(SQL, encoding='utf-8').read()
up = txt.split('-- +goose Down')[0].split('-- +goose Up', 1)[1]
up_exec = up.replace('-- +goose StatementBegin', '').replace('-- +goose StatementEnd', '')
db = sqlite3.connect(':memory:')
db.executescript(up_exec)
tables = {}
# structured comments
blocks = re.split(r'\n(?=-- @table )', up)
for b in blocks:
    m = re.match(r'-- @table (\w+)', b)
    if not m: continue
    t = m.group(1)
    meta = {k: (re.search(r'-- @%s (.*)' % k, b).group(1).strip() if re.search(r'-- @%s (.*)' % k, b) else '') for k in ('title', 'entity', 'desc', 'used')}
    cm = re.search(r'CREATE TABLE %s \((.*?)\n\) STRICT;' % t, b, re.S)
    body = cm.group(1)
    colcomments = {}; enums = {}
    for line in body.split('\n'):
        line = line.strip()
        mm = re.match(r'(\w+) (INTEGER|TEXT)\b(.*)', line)
        if not mm: continue
        name, typ, rest = mm.groups()
        comment = ''
        if '--' in rest:
            rest, comment = rest.split('--', 1)
        colcomments[name] = comment.strip()
        e = re.search(r"CHECK \(%s IN \(([^)]*)\)" % name, rest)
        if e:
            vals = [v.strip().strip("'") for v in e.group(1).split(',')]
            enums[name] = vals
    cols = []
    for cid, name, typ, notnull, dflt, pk in db.execute(f'PRAGMA table_info({t})'):
        cols.append(dict(name=name, type=typ, notnull=bool(notnull), default=dflt, pk=bool(pk), comment=colcomments.get(name, ''), enum=enums.get(name)))
    fks = [dict(col=r[3], ref=r[2], refcol=r[4], on_delete=r[6]) for r in db.execute(f'PRAGMA foreign_key_list({t})')]
    tables[t] = dict(meta=meta, cols=cols, fks=fks)
json.dump(tables, open(os.path.join(WORK, 'schema.json'), 'w', encoding='utf-8'), ensure_ascii=False, indent=1)
print(len(tables), 'tablas')
for t, d in tables.items():
    print(f"\n## {t} [{d['meta']['entity']}] {d['meta']['title']}")
    for c in d['cols']:
        if c['name'] in ('created_at', 'updated_at'): continue
        en = ('{' + '|'.join(c['enum']) + '}') if c['enum'] else ''
        fk = [f"->{f['ref']}" for f in d['fks'] if f['col'] == c['name']]
        print(f"  {c['name']}:{c['type'][0]}{'!' if c['notnull'] else ''}{en}{''.join(fk)}" + (f" ={c['default']}" if c['default'] and 'strftime' not in str(c['default']) else ''))
