#!/usr/bin/env python3
"""Verifica README, componentes, portada e índice del sistema de diseño (project/).
Uso: python3 -I docs/design/sistema/verificar_brandbook.py [ruta a project/]
Contraste: fórmula WCAG 2.x con los hex de tokens.json (alias resueltos por tema); cascada CSS
mínima (tipo, .clase, #id, [attr], [attr="v"], descendiente y >; :hover/:active/... como estados)."""
import json, re, sys, pathlib
from html.parser import HTMLParser

ROOT = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path(__file__).resolve().parent / 'project'
T = json.loads((ROOT / 'tokens.json').read_text())
COL = {t['name']: t['value'] for t in T['color']['tokens']}
THEMES = [t['id'] for t in T['color']['themes']]
NAMES = set(COL) | {'font-' + k for k in T['type'].get('families', {})}
for k, v in T.items():
    if k != 'color' and isinstance(v, dict) and isinstance(v.get('tokens'), list):
        NAMES |= {t['name'] for t in v['tokens']}
VAR = re.compile(r'var\(--([A-Za-z0-9_.-]+)\)')
DECOR = re.compile(r'^(ui-surface-|ui-skeleton-|ui-border-divider$|ui-viz-grid$|ui-viz-track$|ui-overlay-scrim$|ui-role-[a-z]+-surface$|ui-critical-surface$)')
VOID = {'area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'source', 'track', 'wbr'}
SHAPES = {'path', 'rect', 'circle', 'polygon', 'polyline', 'line', 'ellipse'}
fallas = []
avisos = set()
# Pares del propio DESIGN.md que no llegan al umbral: se dejan exactos y se reportan (bloqueo P3-contraste-chip-hover).
EXCEPCIONES = {('ChipPersonaTienda', 'ui-person-karen', 'ui-surface-hover'): 'DESIGN §7.4 pide hover ui.surface.hover en el chip seleccionable; la marca acompaña al nombre (aria-hidden)'}

def hexof(n, th):
    v = COL[n]; v = v.get(th, v.get(THEMES[0])) if isinstance(v, dict) else v
    m = re.fullmatch(r'\{(.+)\}', v)
    return hexof(m.group(1), th) if m else v

def lum(h):
    h = h.lstrip('#')
    if len(h) in (3, 4): h = ''.join(c * 2 for c in h)
    f = lambda c: c / 12.92 if c <= 0.04045 else ((c + 0.055) / 1.055) ** 2.4
    r, g, b = (f(int(h[i:i + 2], 16) / 255) for i in (0, 2, 4))
    return 0.2126 * r + 0.7152 * g + 0.0722 * b

def ratio(a, b, th):
    x, y = lum(hexof(a, th)), lum(hexof(b, th))
    return (max(x, y) + 0.05) / (min(x, y) + 0.05)

class El:
    def __init__(s, tag, attrs, parent):
        s.tag, s.attrs, s.parent, s.kids, s.text = tag, dict(attrs), parent, [], ''
        s.cls = set((s.attrs.get('class') or '').split())
        if parent: parent.kids.append(s)

class Doc(HTMLParser):
    def __init__(s):
        super().__init__(convert_charrefs=True); s.root = El('#doc', {}, None); s.cur = s.root; s.stack = []; s.err = []; s.css = ''
    def handle_starttag(s, tag, attrs):
        e = El(tag, attrs, s.cur)
        if tag not in VOID: s.stack.append(tag); s.cur = e
    def handle_startendtag(s, tag, attrs): El(tag, attrs, s.cur)
    def handle_endtag(s, tag):
        if tag in VOID: return
        if not s.stack or s.stack[-1] != tag: s.err.append(f'cierre </{tag}> sin abrir (pila {s.stack[-3:]})'); return
        s.stack.pop(); s.cur = s.cur.parent
    def handle_data(s, d):
        if s.cur.tag == 'style': s.css += d
        elif s.cur.tag not in ('script', 'title'): s.cur.text += d

def parse_css(css):
    css = re.sub(r'/\*.*?\*/', '', css, flags=re.S); out = []; i = 0
    while (j := css.find('{', i)) >= 0:
        sel, depth, k = css[i:j].strip(), 1, j + 1
        while depth:
            depth += {'{': 1, '}': -1}.get(css[k], 0); k += 1
        if not sel.startswith('@'):
            decl = {p.strip().lower(): v.strip() for p, _, v in (d.partition(':') for d in css[j + 1:k - 1].split(';')) if v}
            for s in sel.split(','):
                out.append(selector(s.strip(), decl, len(out)))
        i = k
    return out

TOK = re.compile(r'(::?[\w-]+)|\.([\w-]+)|#([\w-]+)|\[([\w-]+)(?:="([^"]*)")?\]|([a-zA-Z][\w-]*|\*)')
def selector(s, decl, order):
    parts = re.split(r'(\s*>\s*|\s+)', s); chain = []; spec = [0, 0, 0]; pseudo = []
    for n, p in enumerate(parts):
        if n % 2: chain.append('>' if '>' in p else ' '); continue
        c = {'tag': None, 'cls': set(), 'id': None, 'at': []}
        for m in TOK.finditer(p):
            if m.group(1): pseudo.append(m.group(1)); spec[1 if not m.group(1).startswith('::') else 2] += 1
            elif m.group(2): c['cls'].add(m.group(2)); spec[1] += 1
            elif m.group(3): c['id'] = m.group(3); spec[0] += 1
            elif m.group(4): c['at'].append((m.group(4), m.group(5))); spec[1] += 1
            elif m.group(6) != '*': c['tag'] = m.group(6); spec[2] += 1
        chain.append(c)
    return {'chain': chain, 'decl': decl, 'key': (tuple(spec), order), 'pseudo': pseudo, 'txt': s}

def mcomp(c, e):
    return (not c['tag'] or c['tag'] == e.tag) and c['cls'] <= e.cls and (not c['id'] or e.attrs.get('id') == c['id']) \
        and all(a in e.attrs and (v is None or e.attrs[a] == v) for a, v in c['at'])

def match(ch, i, e):
    if e is None or e.tag == '#doc' or not mcomp(ch[i], e): return False
    if i == 0: return True
    if ch[i - 1] == '>': return match(ch, i - 2, e.parent)
    p = e.parent
    while p is not None:
        if match(ch, i - 2, p): return True
        p = p.parent
    return False

def cval(v):
    if v is None: return None
    m = VAR.search(v)
    if m: return m.group(1)
    v = v.strip().lower()
    return 'none' if v.startswith(('none', 'transparent', 'url')) else ('currentcolor' if v == 'currentcolor' else None)

def own(e, rules, extra=None):
    d = {}
    for r in sorted((r for r in rules if not r['pseudo'] and match(r['chain'], len(r['chain']) - 1, e)), key=lambda r: r['key']):
        d.update(r['decl'])
    for kv in (e.attrs.get('style') or '').split(';'):
        p, _, v = kv.partition(':')
        if v: d[p.strip().lower()] = v.strip()
    for a in ('fill', 'stroke'):
        if a in e.attrs and a not in d: d[a] = e.attrs[a]
    if extra: d.update(extra)
    return d

def border(d):
    for p in ('border-color', 'border'):
        if p in d:
            v = cval(d[p]) if VAR.search(d[p]) else ('none' if 'transparent' in d[p] or d[p] == '0' else None)
            return v
    return None

def walk(e, ctx, rules, pares, extra=None, indefs=False):
    d = own(e, rules, (extra or {}).get(id(e)))
    color = cval(d.get('color')) or ctx['color']
    fill = cval(d.get('fill')) or ctx['fill']; stroke = cval(d.get('stroke')) or ctx['stroke']
    fill = color if fill == 'currentcolor' else fill; stroke = color if stroke == 'currentcolor' else stroke
    bgo = cval(d.get('background') or d.get('background-color'))
    bg = bgo if bgo and bgo != 'none' else ctx['bg']
    indefs = indefs or e.tag in ('defs', 'pattern')
    if e.text.strip() and not indefs:
        pares.add(('texto', fill if e.tag in ('text', 'tspan') else color, bg, e.tag))
    if e.tag in ('input', 'select') and e.attrs.get('type') not in ('checkbox', 'radio', 'hidden'):
        pares.add(('texto', color, bg, e.tag))
    if bgo and bgo not in ('none',) and not DECOR.match(bgo): pares.add(('grafico', bgo, ctx['bg'], e.tag))
    b = border(d)
    if b and b != 'none' and not DECOR.match(b): pares.add(('grafico', b, ctx['bg'], e.tag + ' borde'))
    if e.tag in SHAPES and not indefs:
        for v in (fill, stroke):
            if v and v not in ('none',) and v in COL and not DECOR.match(v): pares.add(('grafico', v, bg, e.tag))
    sub = {'color': color, 'fill': fill if 'fill' in d else ctx['fill'], 'stroke': stroke if 'stroke' in d else ctx['stroke'], 'bg': bg}
    sub['fill'], sub['stroke'] = fill, stroke
    e.ctx = ctx
    for k in e.kids: walk(k, sub, rules, pares, extra, indefs)

def check_preview(f):
    s = f.read_text(); rel = f.relative_to(ROOT); lines = s.split('\n')
    cover = f.parent.name == 'Cover'
    m = re.fullmatch(r'<!-- @dsCard(?: group="[^"]+")? height=(\d+)(?: subtitle="[^"]*")? -->', lines[0])
    if not m: fallas.append(f'{rel}: línea 1 sin marcador @dsCard')
    elif not (240 <= int(m.group(1)) <= 360 if cover else 40 <= int(m.group(1)) <= 4000): fallas.append(f'{rel}: height fuera de rango')
    doc = Doc(); doc.feed(s); doc.close()
    for e in doc.err: fallas.append(f'{rel}: {e}')
    if doc.stack: fallas.append(f'{rel}: etiquetas sin cerrar {doc.stack}')
    if re.search(r'<(iframe|frame|object|embed|portal|noscript)\b', s): fallas.append(f'{rel}: elemento prohibido')
    if re.search(r'https?://', s): fallas.append(f'{rel}: URL externa')
    if re.search(r'font-family\s*:[^;}]*Inter', s): fallas.append(f'{rel}: «Inter» escrito como pila (usar var(--font-sans))')
    rules = parse_css(doc.css); pares = set()
    for e in doc.root.kids: walk(e, {'color': 'ui-text-default', 'fill': None, 'stroke': None, 'bg': 'ui-surface-page'}, rules, pares)
    # estados (:hover, :active, :focus-within…) sobre el elemento que los recibe y su subárbol
    nodes = []
    def coll(e):
        nodes.append(e); [coll(k) for k in e.kids]
    coll(doc.root)
    for r in rules:
        if '::placeholder' in r['pseudo']:
            for e in nodes:
                if e.tag == 'input' and 'placeholder' in e.attrs and match(r['chain'], len(r['chain']) - 1, e):
                    pares.add(('texto', cval(r['decl'].get('color')), e.ctx['bg'], 'placeholder'))
    # estados: el elemento se pone en :hover, :active… y se recalcula su cascada completa (base + reglas de ese estado)
    COLP = ('color', 'background', 'background-color', 'border', 'border-color', 'fill', 'stroke')
    for P in sorted({p for r in rules for p in r['pseudo'] if not p.startswith('::')}):
        srules = [r for r in rules if r['pseudo'] and set(r['pseudo']) == {P}]
        for e in nodes:
            if e.tag == '#doc': continue
            ms = [r for r in srules if match(r['chain'], len(r['chain']) - 1, e)]
            if not any(k in r['decl'] for r in ms for k in COLP): continue
            d = {}
            for r in sorted([r for r in rules if not r['pseudo'] and match(r['chain'], len(r['chain']) - 1, e)] + ms, key=lambda r: r['key']):
                d.update(r['decl'])
            walk(e, e.ctx, rules, pares, {id(e): d})
    n = {'texto': 0, 'grafico': 0}; peor = {}
    for kind, fg, bg, tag in sorted(pares, key=str):
        if not fg or fg not in COL or not bg or bg not in COL: continue
        if cover and kind == 'grafico': continue
        n[kind] += 1
        for th in THEMES:
            r = ratio(fg, bg, th); lim = 4.5 if kind == 'texto' else 3.0
            peor[(kind, th)] = min(peor.get((kind, th), 99), r)
            if r < lim and (f.parent.name, fg, bg) in EXCEPCIONES: avisos.add(f'{rel}: {kind} {fg} sobre {bg} {th} {r:.2f}:1 < {lim} · excepción declarada: {EXCEPCIONES[(f.parent.name, fg, bg)]}')
            elif r < lim: fallas.append(f'{rel}: {kind} {fg} sobre {bg} ({tag}) {th} {r:.2f}:1 < {lim}')
    for e in nodes:
        if e.tag in ('button', 'a') and not ''.join(x.text for x in sub_nodes(e)).strip() and not (e.attrs.get('aria-label') or e.attrs.get('aria-labelledby')):
            fallas.append(f'{rel}: <{e.tag}> sin texto ni aria-label')
    mins = ' · '.join(f'{k[0]} {k[1]} mín {v:.2f}' for k, v in sorted(peor.items()))
    print(f'{rel}: pares texto {n["texto"]}, gráficos {n["grafico"]} · {mins}')

def sub_nodes(e):
    out = [e]
    for k in e.kids: out += sub_nodes(k)
    return out

previews = sorted(ROOT.glob('components/*/preview.html'))
for f in previews: check_preview(f)
textos = [ROOT / 'README.md'] + sorted(ROOT.glob('components/*/README.md')) + sorted(ROOT.glob('*.md')) + previews
usados = set()
for f in set(textos):
    s = f.read_text(); rel = f.relative_to(ROOT)
    for v in VAR.findall(s):
        usados.add(v)
        if v not in NAMES: fallas.append(f'{rel}: var(--{v}) no existe en tokens.json')
    if re.search('[\U0001F000-\U0001FAFF☀-➿️]', s): fallas.append(f'{rel}: emoji o pictograma')
    if re.search(r'\$[0-9]{4,}\.[0-9]{2}', s): fallas.append(f'{rel}: monto de 1,000 o más sin separador de miles')
    if re.search(r'-\$[0-9]', s): fallas.append(f'{rel}: guion corto antes de $ (usar U+2212)')
for d in sorted(p for p in (ROOT / 'components').iterdir() if p.is_dir()):
    tiene = {p.name for p in d.iterdir()}
    falta = ({'preview.html'} if d.name == 'Cover' else {'preview.html', 'README.md'}) - tiene
    if falta: fallas.append(f'components/{d.name}: falta {falta}')
    if d.name == 'Cover' and tiene - {'preview.html'}: fallas.append('components/Cover: la carpeta debe quedar sola')
if (ROOT / 'README.md').read_text().startswith('# '): fallas.append('README.md: no lleva título')
canon = re.compile(r'3,?094\.50|800\.50|598\.82|374\.93|5,?034\.74|1,?258\.68|883\.75|6,?293\.44|973\.75|1,?258\.70|883\.77|2,?395\.28|1,?499\.72|2,?616\.87|1,?767\.51|849\.36|8,?094\.50|14,?362\.48|13,?103\.78|25\.8685|17\.27')
hits = [f'{p.relative_to(ROOT)}:{m.group(0)}' for p in ROOT.rglob('*') if p.is_file() and p.suffix in ('.md', '.html', '.json') for m in canon.finditer(p.read_text())]
fallas += [f'cifra canon (P3-12): {h}' for h in hits]
idx = ROOT / 'design-system.json'
if idx.exists():
    j = json.loads(idx.read_text())
    for k in ('v', 'layout', 'createdOnFiles', 'title', 'namespace', 'libraries', 'sections', 'groups', 'assetGroups', 'blobs', 'docs', 'lastChange'):
        if k not in j: fallas.append(f'design-system.json: falta «{k}»')
else: fallas.append('design-system.json: falta el archivo')
print(f'previews {len(previews)} · variables distintas usadas {len(usados)} (todas en tokens.json: {not any("no existe" in x for x in fallas)}) · cifras canon {len(hits)}')
print('\n'.join(sorted(avisos)))
print('\n'.join(fallas) if fallas else f'RESULTADO: OK · 0 fallas · {len(avisos)} avisos de excepción declarada')
sys.exit(1 if fallas else 0)
