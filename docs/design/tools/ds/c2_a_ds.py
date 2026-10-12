#!/usr/bin/env python3
"""c2_a_ds.py · Convierte C2-tokens.json (DTCG, mapa) al formato LISTA del tipo
«Design System» (docs/design/sistema/referencia-tipo/format.md, «tokens.json grammar»).

Entrada : docs/contratos/C2-tokens.json (1.0, aprobado Q2-96 A)  [+ tokens.css solo para nombres]
Salida  : docs/design/sistema/project/tokens.json  y  project/fonts/*.woff2
No cambia ningun valor: solo cambia la forma. Uso:  python3 -I c2_a_ds.py
"""
import json, re, shutil, sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[3]            # docs/
C2 = RAIZ / 'contratos' / 'C2-tokens.json'
CSS = RAIZ / 'contratos' / 'tokens.css'
FUENTES = RAIZ / 'design' / 'fuentes' / 'inter' / 'web'
OUT = RAIZ / 'design' / 'sistema' / 'project'
NOMBRE_OK = re.compile(r'^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$')

c2 = json.loads(C2.read_text(encoding='utf-8'))
css_vars = set(re.findall(r'^\s*--([\w-]+):', CSS.read_text(encoding='utf-8'), re.M))

# ---------- 1. aplanar C2: hojas con su ruta, nodo y descripciones de los ancestros ----------
hojas = []          # (ruta, nodo, [descripciones de ancestros, del mas cercano al mas lejano])
def recorre(n, ruta, anc):
    if not isinstance(n, dict):
        return
    if '$value' in n:
        hojas.append((ruta, n, anc))
        return
    anc2 = ([n['$description']] if '$description' in n else []) + anc
    for k, v in n.items():
        if k == '$root':
            if '$value' in v:
                hojas.append((ruta, v, anc2))
        elif not k.startswith('$'):
            recorre(v, ruta + [k], anc2)
recorre(c2, [], [])
POR_RUTA = {'.'.join(r): n for r, n, _ in hojas}

def es_tema(ruta):
    return ruta[-1] in ('light', 'dark')

def nombre_css(ruta):
    """Nombre de la variable CSS de tokens.css (sin --) que corresponde a una ruta C2."""
    r = list(ruta)
    if r[0] == 'ui' and es_tema(r):
        r = r[:-1]
    if r[0] == 'overlay':
        return 'ui-overlay-scrim'
    if r[0] == 'motion':
        r = [{'duration': 'duration', 'easing': 'ease'}[r[1]]] + r[2:]
    return '-'.join(r)

# ---------- 2. resolucion de valores ----------
REF = re.compile(r'^\{([^}]+)\}$')
def literal(valor, hasta=0):
    """Resuelve alias DTCG a su valor final (cadena o estructura)."""
    if hasta > 20:
        raise RuntimeError('ciclo de alias')
    if isinstance(valor, str):
        m = REF.match(valor)
        if m:
            return literal(POR_RUTA[m.group(1)]['$value'], hasta + 1)
    return valor

def hex_min(h):
    return h.lower()

def hex_rgba(h, a):
    h = h.lstrip('#')
    r, g, b = int(h[0:2], 16), int(h[2:4], 16), int(h[4:6], 16)
    return r, g, b, a

def num(x):
    s = ('%.4f' % x).rstrip('0').rstrip('.')
    return s or '0'

def color_dtcg_a_ds(valor, tema):
    """Alias DTCG '{a.b.c}' -> alias del tipo '{nombre-aplanado}' si el destino existe; si no, hex literal."""
    m = REF.match(valor) if isinstance(valor, str) else None
    if not m:
        return hex_min(valor)
    destino = m.group(1).split('.')
    if destino[0] == 'ui' and es_tema(destino):
        assert destino[-1] == tema, ('alias de ui a otro tema', valor)
    nombre = nombre_css(destino)
    if nombre in NOMBRES_COLOR:
        return '{' + nombre + '}'
    return hex_min(literal(valor))

# ---------- 3. descripciones (usage) ----------
def etiqueta(anc):
    """Primeros dos segmentos de la descripcion del grupo: 'Ingreso · matiz 135°'."""
    for d in anc:
        if ' · matiz' in d or d.startswith('Marca') or d.startswith('Neutro'):
            return ' · '.join(d.split(' · ')[:2])
    return None

def recorta(s, n=1000):
    return s if len(s) <= n else s[:n - 1].rstrip() + '…'

# ---------- 4. color ----------
color_hojas = [(r, n, a) for r, n, a in hojas if r[0] in ('color', 'ui')]
NOMBRES_COLOR = set()
for r, n, a in hojas:
    if r[0] in ('color', 'ui'):
        NOMBRES_COLOR.add(nombre_css(r))
NOMBRES_COLOR.add('ui-overlay-scrim')

tokens_color = {}      # nombre -> dict(value={light,dark}, usage, ruta)
orden = []
for r, n, a in color_hojas:
    nom = nombre_css(r)
    if r[0] == 'color':
        v = color_dtcg_a_ds(n['$value'], 'light')
        tokens_color[nom] = {'value': {'light': v, 'dark': color_dtcg_a_ds(n['$value'], 'dark')}, 'r': r, 'n': n, 'a': a}
        orden.append(nom)
    else:   # ui: se juntan .light y .dark
        t = r[-1]
        e = tokens_color.setdefault(nom, {'value': {}, 'r': r[:-1], 'n': n, 'a': a})
        e['value'][t] = color_dtcg_a_ds(n['$value'], t)
        if nom not in orden:
            orden.append(nom)
for nom, e in tokens_color.items():
    if nom.startswith('ui-'):
        assert set(e['value']) == {'light', 'dark'}, nom
        e['value'] = {'light': e['value']['light'], 'dark': e['value']['dark']}

# velo (overlay.scrim): neutral.950 con transparencia, resuelto a #rrggbbaa por tema
def con_alfa(hexv, alfa):
    r, g, b, _ = hex_rgba(hexv, alfa)
    return '#%02x%02x%02x%02x' % (r, g, b, round(alfa * 255))
vs = {}
for t in ('light', 'dark'):
    nodo = POR_RUTA['overlay.scrim.' + t]
    vs[t] = con_alfa(literal(nodo['$value']), nodo['$extensions']['finanzas']['alpha'])
tokens_color['ui-overlay-scrim'] = {'value': vs, 'r': ['overlay', 'scrim'], 'n': None,
    'a': [c2['overlay']['scrim']['$description'], c2['overlay']['$description']]}
orden.append('ui-overlay-scrim')

# quien referencia a quien (para el usage de los primitivos)
refs = {}
for nom, e in tokens_color.items():
    for t in ('light', 'dark'):
        v = e['value'][t]
        if v.startswith('{'):
            refs.setdefault(v[1:-1], {}).setdefault(nom, set()).add(t)
TEMA = {'light': 'claro', 'dark': 'oscuro'}
def usado_por(nom):
    if nom not in refs:
        return ''
    trozos = []
    for quien, ts in sorted(refs[nom].items()):
        if len(ts) == 2:
            trozos.append(quien)
        else:
            trozos.append('%s (%s)' % (quien, TEMA[next(iter(ts))]))
    s = 'Lo usan: ' + ', '.join(trozos) + '.'
    return s

ROL_PASO = {'700': 'texto e icono en claro', '400': 'texto e icono en oscuro',
            '100': 'fondo de insignia en claro', '900': 'fondo de insignia en oscuro'}
def usage_color(nom, e):
    r, n, a = e['r'], e['n'], e['a']
    partes = []
    if nom == 'ui-overlay-scrim':
        return ('Velo bajo diálogos, hojas y el panel de la barra lateral: neutral.950 al 40 % en claro y al 60 % en oscuro (DESIGN §4.5). '
                'Valor resuelto a #rrggbbaa porque el tipo no admite la mezcla color-mix de tokens.css.')
    if n is not None and '$description' in n:
        partes.append(n['$description'].rstrip('.') + '.')
    elif nom.startswith('ui-') and a:
        partes.append(a[0].rstrip('.') + '.')
    et = etiqueta(a)
    if nom.startswith('color-'):
        paso = r[-1]
        if paso.isdigit():
            ext = (n or {}).get('$extensions', {}).get('oklch', {})
            base = '%s, paso %s' % (et or '.'.join(r[:-1]), paso)
            if 'l' in ext:
                base += ' (L %s)' % num(ext['l'])
            partes.insert(0, base + '.')
            if r[1] == 'role' and paso in ROL_PASO:
                partes.append('Del rol: %s (DESIGN §2.4).' % ROL_PASO[paso])
        elif et and not partes[0].startswith(et):
            partes.insert(0, et + '.')
        elif not partes:
            partes.append((et or '.'.join(r)) + '.')
    u = usado_por(nom)
    if u:
        partes.append(u)
    return recorta(' '.join(partes))

color_list = []
for nom in orden:
    e = tokens_color[nom]
    assert NOMBRE_OK.match(nom), nom
    color_list.append({'name': nom, 'value': e['value'], 'usage': usage_color(nom, e)})

# ---------- 5. familias {tokens:[...]} (valores CSS simples) ----------
def px(v):
    return v if isinstance(v, str) else num(v)

def valor_simple(r, n):
    v = literal(n['$value'])
    if r[0] == 'viz':
        if v == 'solid':
            return 'none'                                   # tokens.css: --viz-series-1-dash: none
        return ' '.join(x.replace('px', '') for x in v['dashArray'])
    if r[0] == 'motion' and r[1] == 'easing':
        return 'cubic-bezier(%s)' % ', '.join(num(x) for x in v)
    if isinstance(v, (int, float)):
        return num(v) + n.get('$extensions', {}).get('finanzas', {}).get('unit', '')   # p. ej. 85 + «%»
    return str(v)

def usage_simple(r, n, a):
    partes = []
    if '$description' in n:
        partes.append(n['$description'])
    elif a:
        partes.append(a[0])
    raw = n['$value']
    if isinstance(raw, str) and REF.match(raw):
        partes.append('Alias de %s en C2.' % nombre_css(REF.match(raw).group(1).split('.')))
    ext = n.get('$extensions', {}).get('finanzas', {})
    if r[0] == 'viz':
        pat = ext.get('pattern', {})
        pd = ', '.join('%s %s' % (k, pat[k]) for k in pat if k != 'kind')
        partes.append('Marca: %s. Patrón: %s%s. Valor de stroke-dasharray (none = continuo).' % (
            ext.get('marker'), pat.get('kind'), (' (' + pd + ')') if pd else ''))
    return recorta(' '.join(p.rstrip() for p in partes))

FAMILIAS = [  # clave del tipo, raiz C2, nota
    ('spacing', 'space', 'Retícula de 4 pt: siete valores y el cero (DESIGN §4.1). En Tailwind 4, --spacing: 0.25rem.'),
    ('radius', 'radius', 'Un radio interior nunca es mayor que el de su contenedor (DESIGN §4.4).'),
    ('border', 'border', 'Grosores de borde (DESIGN §4.5). El anillo de foco va con la separación de la familia focus.'),
    ('focus', 'focus', 'Separación entre el control y el anillo de foco (DESIGN §4.5, §9.2).'),
    ('size', 'size', 'Medidas de control, icono y componente (DESIGN §4.6, §5.2, §10.6). Los alias de C2 se resuelven a su valor.'),
    ('layout', 'layout', 'Medidas de estructura (DESIGN §4.3, §10.6).'),
    ('icon', 'icon', 'Trazo de iconos Lucide (DESIGN §5.2): 1.5 px dibujados; en unidades del viewBox de 24, 1.5 (24 px) y 1.8 (20 px).'),
    ('motion', 'motion', 'Movimiento (DESIGN §6): 150–250 ms; la curva enfática solo para entrar.'),
    ('breakpoint', 'breakpoint', 'Puntos de quiebre con min-width, móvil primero (DESIGN §9.4). Como variables no sirven dentro de @media.'),
    ('viz', 'viz', 'Trazo, marca y patrón de cada serie de gráficas (DESIGN §8.6, K27): el color nunca va solo.'),
]
familias = {}
nombres_vistos = {t['name'] for t in color_list}
for clave, raiz, nota in FAMILIAS:
    lista = []
    for r, n, a in hojas:
        if r[0] != raiz:
            continue
        nom = nombre_css(r)
        if r == ['icon', 'stroke']:
            nom = 'icon-stroke'
        assert NOMBRE_OK.match(nom), nom
        assert nom not in nombres_vistos, ('nombre repetido', nom)
        nombres_vistos.add(nom)
        lista.append({'name': nom, 'value': valor_simple(r, n), 'usage': usage_simple(r, n, a)})
    familias[clave] = {'note': nota, 'tokens': lista}

# sombras: por tema ('none' en oscuro, C2 finanzas.dark / DESIGN §4.5)
sombra = []
for r, n, a in hojas:
    if r[0] != 'shadow':
        continue
    ext = n['$extensions']['finanzas']
    capas = []
    for capa, alfa in zip(n['$value'], ext['alpha']):
        rr, gg, bb, _ = hex_rgba(literal(capa['color']), 0)
        capas.append('%s %s %s %s rgba(%d, %d, %d, %s)' % (
            capa['offsetX'], capa['offsetY'], capa['blur'], capa['spread'], rr, gg, bb, num(alfa)))
    claro = 'none' if ext.get('css') == 'none' else ', '.join(capas)
    oscuro = 'none'
    sombra.append({'name': nombre_css(r), 'value': {'light': claro, 'dark': oscuro},
                   'usage': recorta(n['$description'] + ('' if 'oscuro' in n['$description'] else ' En oscuro no hay sombra: se usa borde de 1 px de ui-border-divider.'))})
    assert NOMBRE_OK.match(sombra[-1]['name'])
    assert sombra[-1]['name'] not in nombres_vistos
    nombres_vistos.add(sombra[-1]['name'])
familias['shadow'] = {'note': 'Dos niveles en claro; en oscuro, ninguna y borde (DESIGN §4.5). Capas con rgba de neutral.950 (#222325) en lugar de la mezcla color-mix de tokens.css.',
                      'tokens': sombra}

# ---------- 6. tipografia ----------
def f(*ruta):
    return POR_RUTA['.'.join(ruta)]
ESCALA = [  # nombre, peso, muestra (sintetica, K45.1)
    ('title-1', 'bold', 'Resumen de octubre'),
    ('title-2', 'semibold', '$4,826.37'),
    ('title-3', 'semibold', 'Despensa de la quincena'),
    ('body', 'regular', 'Ximena registró el pago del súper del viernes.'),
    ('small', 'regular', 'Walmart · 12-oct-2026 · −$1,037.40'),
    ('caption', 'regular', 'Vence el 30-oct-2026 · Karen'),
]
TABULAR = ' Cifras tabulares en todo dinero, porcentaje y cantidad: font-variant-numeric: tabular-nums (rasgo tnum de Inter; clase .money, DESIGN §3.3).'
estilos = []
for nom, peso, muestra in ESCALA:
    desc = f('font', 'size', nom)['$description']
    desc = re.sub(r'^Tamaño \d+ px: ', '', desc)
    desc = desc[0].upper() + desc[1:]
    desc = re.sub(r'\s*\(DESIGN §3\.2\)\.?$', '', desc)
    w = f('font', 'weight', peso)
    uso = '%s. Peso %s: %s.' % (desc.rstrip('.'), w['$value'], f('font', 'weight', peso)['$description'].split(' (DESIGN')[0].rstrip('.').lower())
    if nom == 'caption':
        uso += ' 500 (medium) en etiquetas (DESIGN §3.2).'
    estilos.append({
        'name': nom,
        'fontSize': f('font', 'size', nom)['$value'],
        'lineHeight': f('font', 'line-height', nom)['$value'],
        'fontWeight': w['$value'],
        'letterSpacing': f('font', 'letter-spacing')['$value'],
        'sample': muestra,
        'usage': recorta(uso + ('' if nom == 'title-1' else TABULAR)),
    })
fam_c2 = f('font', 'family', 'sans')
pila = re.search(r'^\s*--font-sans:\s*(.+);', CSS.read_text(encoding='utf-8'), re.M).group(1)   # misma pila que tokens.css
assert [x.strip().strip('"') for x in pila.split(',')] == fam_c2['$value'], 'la pila de tokens.css difiere de C2'
fuentes = []
for cara in fam_c2['$extensions']['finanzas']['faces']:
    archivo = Path(cara['src']).name
    fuentes.append({'family': 'Inter', 'file': 'fonts/' + archivo, 'weight': cara['weight'], 'style': cara['style']})
    OUT.joinpath('fonts').mkdir(parents=True, exist_ok=True)
    shutil.copyfile(FUENTES / archivo, OUT / 'fonts' / archivo)
tipo = {
    'fonts': fuentes,
    'families': {'sans': pila},
    'groups': [{'name': 'Interfaz', 'family': 'sans',
                'note': 'Inter 4.1 variable (SIL OFL 1.1, autoalojada), seis tamaños sobre retícula de 4 pt; la jerarquía la dan el peso y el color, no muchos tamaños (DESIGN §3).',
                'styles': estilos}],
}

# ---------- 6b. familia fontWeight (revision ds-r1, M-02) ----------
# Los 4 pesos de C2 (font.weight.*) como valores simples, con los mismos nombres que tokens.css
# (--font-weight-regular/medium/semibold/bold). format.md: cualquier otra clave de primer nivel
# {tokens:[...]} con valores CSS simples se compila en :root. Los estilos de type conservan su peso.
CSS_RAIZ = dict(re.findall(r'^\s*--(font-weight-[a-z]+):\s*([^;]+);', CSS.read_text(encoding='utf-8'), re.M))
pesos_fam = []
for clave_peso in ('regular', 'medium', 'semibold', 'bold'):
    w = f('font', 'weight', clave_peso)
    nom = 'font-weight-' + clave_peso
    assert NOMBRE_OK.match(nom) and nom not in nombres_vistos, nom
    assert CSS_RAIZ.get(nom) == str(w['$value']), ('peso distinto de tokens.css', nom)
    nombres_vistos.add(nom)
    pesos_fam.append({'name': nom, 'value': str(w['$value']), 'usage': recorta(w['$description'])})
familias['fontWeight'] = {'note': 'Pesos de Inter como variables --font-weight-* (DESIGN §3.2; mismos nombres que tokens.css). Los estilos de Tipografía llevan además su propio peso.',
                          'tokens': pesos_fam}

# ---------- 7. salida ----------
salida = {
    'version': 1,
    'meta': {
        'name': 'Finanzas Hogar',
        'source': {'kind': 'c2', 'file': 'docs/contratos/C2-tokens.json', 'version': '1.0', 'synced': '2026-10-11'},
        'note': 'Generado por docs/design/tools/ds/c2_a_ds.py; no se edita a mano. Nombres = variables de tokens.css sin «--».',
    },
    'color': {
        'themes': [{'id': 'light', 'name': 'Claro'}, {'id': 'dark', 'name': 'Oscuro'}],
        'note': 'Primitivos iguales en los dos temas (color-*) y capa semántica por tema (ui-*): los componentes usan solo ui-* (DESIGN §7.1).',
        'tokens': color_list,
    },
    'type': tipo,
}
for clave in ('spacing', 'radius', 'shadow', 'fontWeight', 'border', 'focus', 'size', 'layout', 'icon', 'motion', 'breakpoint', 'viz'):
    salida[clave] = familias[clave]

def vuelca(o):
    """JSON legible: un token por linea."""
    lineas = ['{']
    claves = list(o)
    for i, k in enumerate(claves):
        v = o[k]
        coma = ',' if i < len(claves) - 1 else ''
        if k in ('version',) or k == 'meta':
            lineas.append('  %s: %s%s' % (json.dumps(k), json.dumps(v, ensure_ascii=False), coma))
        elif k == 'type':
            lineas.append('  "type": ' + json.dumps(v, ensure_ascii=False, indent=1).replace('\n', '\n  ') + coma)
        else:
            cab = {kk: vv for kk, vv in v.items() if kk != 'tokens'}
            toks = v['tokens']
            lineas.append('  %s: {' % json.dumps(k))
            for kk, vv in cab.items():
                if kk == 'themes':
                    lineas.append('    "themes": %s,' % json.dumps(vv, ensure_ascii=False))
                else:
                    lineas.append('    %s: %s,' % (json.dumps(kk), json.dumps(vv, ensure_ascii=False)))
            lineas.append('    "tokens": [')
            for j, t in enumerate(toks):
                lineas.append('      ' + json.dumps(t, ensure_ascii=False) + (',' if j < len(toks) - 1 else ''))
            lineas.append('    ]')
            lineas.append('  }' + coma)
    lineas.append('}')
    return '\n'.join(lineas) + '\n'

texto = vuelca(salida)
json.loads(texto)                                   # valida
OUT.mkdir(parents=True, exist_ok=True)
(OUT / 'tokens.json').write_text(texto, encoding='utf-8')
total = len(color_list) + sum(len(v['tokens']) for k, v in salida.items() if k not in ('version', 'meta', 'color', 'type'))
print('tokens.json: %d bytes · color %d · otras familias %d · tipos %d estilos · fuentes %d' % (
    len(texto.encode()), len(color_list), total - len(color_list), len(estilos), len(fuentes)))
for k in ('spacing', 'radius', 'shadow', 'fontWeight', 'border', 'focus', 'size', 'layout', 'icon', 'motion', 'breakpoint', 'viz'):
    print('  %-10s %d' % (k, len(salida[k]['tokens'])))
