#!/usr/bin/env python3
"""verificar_ds.py · Verifica docs/design/sistema/project/tokens.json contra C2 y tokens.css.

Comprueba: conteo de tokens por familia (C2 contra tokens.json), cada hex de claro y oscuro contra
tokens.css (resolviendo alias en las dos listas), nombres unicos y validos, alias existentes y sin ciclos,
0 valores vacios, fuentes presentes en fonts/, ningun var( ni color-mix( en los valores, topes de la
gramatica del tipo y 0 cifras canon (P3-12). Termina con «RESULTADO: OK» o «RESULTADO: FALLA».
Uso: python3 -I verificar_ds.py
"""
import json, re, sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[3]            # docs/
C2 = RAIZ / 'contratos' / 'C2-tokens.json'
CSS = RAIZ / 'contratos' / 'tokens.css'
PROY = RAIZ / 'design' / 'sistema' / 'project'
NOMBRE_OK = re.compile(r'^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$')
HEX = re.compile(r'^#(?:[0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})$')
ALIAS = re.compile(r'^\{([A-Za-z0-9][A-Za-z0-9_.-]{0,63})\}$')
SIMPLE = re.compile(r'^[A-Za-z0-9 #%(),./+_-]{1,200}$')
CANON = re.compile(r'3,?094\.50|800\.50|598\.82|374\.93|5,?034\.74|1,?258\.68|883\.75|6,?293\.44|973\.75|1,?258\.70|883\.77|2,?395\.28|1,?499\.72|2,?616\.87|1,?767\.51|849\.36|8,?094\.50|14,?362\.48|13,?103\.78|25\.8685|17\.27')
fallas, avisos, info = [], [], []
def ok(c, msg):
    if not c:
        fallas.append(msg)
    return c

# ---------------- carga ----------------
c2 = json.loads(C2.read_text(encoding='utf-8'))
ds = json.loads((PROY / 'tokens.json').read_text(encoding='utf-8'))
ds_texto = (PROY / 'tokens.json').read_text(encoding='utf-8')
css_texto = CSS.read_text(encoding='utf-8')

hojas = []
def recorre(n, ruta):
    if not isinstance(n, dict):
        return
    if '$value' in n:
        hojas.append((ruta, n)); return
    for k, v in n.items():
        if k == '$root' and '$value' in v:
            hojas.append((ruta, v))
        elif not k.startswith('$'):
            recorre(v, ruta + [k])
recorre(c2, [])
por_raiz = {}
for r, n in hojas:
    por_raiz[r[0]] = por_raiz.get(r[0], 0) + 1
n_c2 = len(hojas)

# ---------------- tokens.css: bloques ----------------
def parsea_css(texto):
    bloques, actual, nombre, buf = {}, None, None, ''
    for linea in texto.splitlines():
        t = linea.strip()
        if actual is None:
            if t == ':root {':
                actual = 'root'; bloques[actual] = {}
            elif t == ':root[data-theme="dark"] {':
                actual = 'dark'; bloques[actual] = {}
            elif t.startswith('@media (prefers-color-scheme: dark)'):
                actual = 'media_pre'
            continue
        if actual == 'media_pre':
            if t.startswith(':root:not([data-theme="light"]) {'):
                actual = 'media'; bloques[actual] = {}
            continue
        if t == '}':
            if buf.strip():
                raise RuntimeError('declaracion sin cerrar: ' + buf)
            actual = None if actual != 'media' else 'media_post'
            continue
        if actual == 'media_post':
            if t == '}':
                actual = None
            continue
        if not t or t.startswith('/*'):
            continue
        buf += ' ' + t
        if t.endswith(';'):
            m = re.match(r'\s*--([\w-]+):\s*(.*);$', buf)
            if m:
                bloques[actual][m.group(1)] = re.sub(r'\s+', ' ', m.group(2)).strip()
            buf = ''
    return bloques
B = parsea_css(css_texto)
ok(set(B) >= {'root', 'dark', 'media'}, 'tokens.css: faltan bloques :root / dark / media')
ok(B['dark'] == B['media'], 'tokens.css: el bloque @media oscuro difiere del [data-theme="dark"]')
RAIZ_V, OSC_V = B['root'], B['dark']
info.append('tokens.css: %d variables en :root, %d reasignadas en oscuro' % (len(RAIZ_V), len(OSC_V)))

def css_valor(nombre, tema, prof=0):
    """Valor de una variable CSS con var(--x) resuelto (el oscuro reasigna ui-*)."""
    if prof > 16:
        raise RuntimeError('ciclo css ' + nombre)
    v = (OSC_V if tema == 'dark' and nombre in OSC_V else RAIZ_V)[nombre]
    return re.sub(r'var\(--([\w-]+)\)', lambda m: css_valor(m.group(1), tema, prof + 1), v)

def a_rgb(h):
    h = h.lstrip('#')
    return int(h[0:2], 16), int(h[2:4], 16), int(h[4:6], 16)

MEZCLA = re.compile(r'color-mix\(in srgb, (#[0-9A-Fa-f]{6}) (\d+)%, transparent\)')
def mezcla_rgba(m):
    r, g, b = a_rgb(m.group(1)); a = int(m.group(2)) / 100
    return 'rgba(%d, %d, %d, %s)' % (r, g, b, ('%.4f' % a).rstrip('0').rstrip('.'))

# ---------------- DS: resolucion de alias ----------------
tokens_color = {t['name']: t for t in ds['color']['tokens']}
def ds_hex(nombre, tema, prof=0):
    if prof > 16:
        fallas.append('alias de más de 16 niveles o ciclo: ' + nombre); return None
    if nombre not in tokens_color:
        fallas.append('alias a token inexistente: ' + nombre); return None
    v = tokens_color[nombre]['value']
    v = v[tema] if isinstance(v, dict) else v
    m = ALIAS.match(v)
    return ds_hex(m.group(1), tema, prof + 1) if m else v

# ---------------- 1. estructura y topes ----------------
claves_lista = [k for k in ds if k not in ('version', 'meta', 'color', 'type')]
ok(ds.get('meta', {}).get('source') == {'kind': 'c2', 'file': 'docs/contratos/C2-tokens.json', 'version': '1.0', 'synced': '2026-10-11'}, 'meta.source no coincide')
ok([t['id'] for t in ds['color']['themes']] == ['light', 'dark'], 'themes: light primero, luego dark')
ok(len(tokens_color) == len(ds['color']['tokens']), 'nombres de color repetidos')
ok(len(tokens_color) <= 600, 'más de 600 colores')
ok(len(claves_lista) + 2 <= 12 + 2 and len(claves_lista) <= 12, 'más de 12 familias')
for k in claves_lista:
    ok(len(ds[k]['tokens']) <= 60, 'familia %s con más de 60 tokens' % k)
    ok(len(ds[k].get('note', '')) <= 400, 'nota larga en ' + k)
ok(len(ds['color'].get('note', '')) <= 400, 'nota larga en color')
ok(len(ds['type']['groups']) <= 12 and sum(len(g['styles']) for g in ds['type']['groups']) <= 80, 'topes de tipografía')

# ---------------- 2. nombres unicos y validos (una sola namespace salvo type) ----------------
vistos = {}
total_tokens = 0
for fam in ['color'] + claves_lista:
    for t in ds[fam]['tokens']:
        total_tokens += 1
        nom = t['name']
        ok(NOMBRE_OK.match(nom) is not None, 'nombre inválido: %r' % nom)
        ok(nom not in vistos, 'nombre repetido entre familias: %s (%s y %s)' % (nom, vistos.get(nom), fam))
        vistos[nom] = fam
        ok(bool(str(t.get('usage', '')).strip()), 'usage vacío: ' + nom)
        ok(len(t.get('usage', '')) <= 1000, 'usage > 1000: ' + nom)
for fam, k in (('type.families', ds['type']['families']),):
    for nom in k:
        ok(NOMBRE_OK.match(nom) is not None, 'clave de familia inválida: ' + nom)
        vistos['font-' + nom] = 'type'
estilos = [s for g in ds['type']['groups'] for s in g['styles']]
nombres_estilo = [s['name'] for s in estilos]
ok(len(set(nombres_estilo)) == len(nombres_estilo), 'estilos de tipo repetidos')
for s in estilos:
    ok(NOMBRE_OK.match(s['name']) is not None, 'estilo inválido ' + s['name'])
    ok(len(s.get('sample', '')) <= 200 and len(s.get('usage', '')) <= 1000, 'sample/usage largo: ' + s['name'])

# ---------------- 3. valores vacios, var(, color-mix( ----------------
def valores(o):
    if isinstance(o, dict):
        for v in o.values():
            yield from valores(v)
    elif isinstance(o, list):
        for v in o:
            yield from valores(v)
    else:
        yield o
vacios = 0
prohibidos = 0
for fam in ['color'] + claves_lista:
    for t in ds[fam]['tokens']:
        for v in valores(t['value']):
            if v is None or str(v).strip() == '':
                vacios += 1; fallas.append('valor vacío: ' + t['name'])
            if re.search(r'var\(|color-mix\(|url\(', str(v)):
                prohibidos += 1; fallas.append('valor con var(/color-mix(/url(: ' + t['name'])
for s in estilos:
    for k in ('fontSize', 'lineHeight', 'fontWeight', 'letterSpacing'):
        if s.get(k) in (None, ''):
            vacios += 1; fallas.append('estilo sin %s: %s' % (k, s['name']))
ok(not re.search(r'var\(|color-mix\(', ds_texto), 'tokens.json contiene «var(» o «color-mix(» en alguna parte')

# ---------------- 4. color: formato, alias y contra tokens.css ----------------
n_hex = 0
for t in ds['color']['tokens']:
    nom, v = t['name'], t['value']
    ok(isinstance(v, dict) and set(v) == {'light', 'dark'}, 'color sin light/dark: ' + nom)
    for tema in ('light', 'dark'):
        x = v[tema]
        ok(HEX.match(x) is not None or ALIAS.match(x) is not None, 'color no válido %s/%s: %r' % (nom, tema, x))
    for tema in ('light', 'dark'):
        h = ds_hex(nom, tema)
        if h is None:
            continue
        # a) contra tokens.css (resuelto)
        if nom == 'ui-overlay-scrim':
            esp = MEZCLA.search(css_valor(nom, tema))
            r, g, b = a_rgb(esp.group(1)); a = round(int(esp.group(2)) / 100 * 255)
            esperado = '#%02x%02x%02x%02x' % (r, g, b, a)
        else:
            esperado = css_valor(nom, tema).lower()
        ok(h == esperado, 'hex distinto de tokens.css: %s/%s ds=%s css=%s' % (nom, tema, h, esperado))
        n_hex += 1
# b) contra C2: cada primitivo de color
for r, n in hojas:
    if r[0] == 'color' and isinstance(n['$value'], str) and n['$value'].startswith('#'):
        nom = '-'.join(r)
        for tema in ('light', 'dark'):
            ok(ds_hex(nom, tema) == n['$value'].lower(), 'hex distinto de C2: ' + nom)

# ---------------- 5. familias de valor simple ----------------
def px(v):
    v = v.strip()
    m = re.fullmatch(r'(-?[\d.]+)rem', v)
    if m:
        v = ('%g' % (float(m.group(1)) * 16)) + 'px'
    v = re.sub(r'(?<![\d.])0(?![\d.%a-z])', '0px', v) if re.fullmatch(r'0', v) else v
    return re.sub(r',\s*', ', ', v)

def capas(s):
    """box-shadow -> lista de tuplas normalizadas (offx, offy, blur, spread, rgba)."""
    s = MEZCLA.sub(mezcla_rgba, s)
    partes = re.split(r'(?<=\)),\s*', s)
    sal = []
    for p in partes:
        m = re.match(r'^(\S+) (\S+) (\S+) (\S+) (rgba\(.*\))$', p.strip())
        sal.append(tuple([px(x) if x != '0' else '0px' for x in m.groups()[:4]] + [re.sub(r'\s+', '', m.group(5))]))
    return sal

for fam in claves_lista:
    for t in ds[fam]['tokens']:
        nom, v = t['name'], t['value']
        if fam == 'shadow':
            ok(set(v) == {'light', 'dark'}, 'sombra sin light/dark: ' + nom)
            css_v = css_valor(nom, 'light')
            if css_v == 'none':
                ok(v['light'] == 'none', 'sombra %s debía ser none' % nom)
            else:
                ok(capas(v['light']) == capas(css_v), 'sombra distinta de tokens.css: %s' % nom)
            ok(v['dark'] == 'none', 'sombra oscura debía ser none: ' + nom)
            continue
        ok(isinstance(v, str), 'valor no es cadena: ' + nom)
        ok(SIMPLE.match(v) is not None and v.count('(') == v.count(')'), 'valor fuera de la gramática: %s=%r' % (nom, v))
        if nom not in RAIZ_V:
            fallas.append('sin variable en tokens.css: ' + nom); continue
        ok(px(v) == px(css_valor(nom, 'light')), 'valor distinto de tokens.css: %s ds=%s css=%s' % (nom, v, css_valor(nom, 'light')))

# ---------------- 6. tipografia ----------------
fam_css = RAIZ_V['font-sans']
ok(ds['type']['families'].get('sans') == fam_css, 'pila sans distinta de tokens.css')
ok(len(ds['type']['families']['sans']) <= 200 and not re.search(r'[;{}<>\\()]', ds['type']['families']['sans']), 'pila sans fuera de gramática')
ESCALA = ['caption', 'small', 'body', 'title-3', 'title-2', 'title-1']
ok(sorted(nombres_estilo) == sorted(ESCALA), 'estilos distintos de los 6 de DESIGN §3.2: %s' % nombres_estilo)
pesos = {400: 'regular', 500: 'medium', 600: 'semibold', 700: 'bold'}
usados_css = {'font-sans', 'font-letter-spacing', 'font-numeric-tabular'}
for s in estilos:
    n = s['name']
    ok(px(s['fontSize']) == px(RAIZ_V['font-size-' + n]), 'fontSize distinto en ' + n)
    ok(px(s['lineHeight']) == px(RAIZ_V['line-height-' + n]), 'lineHeight distinto en ' + n)
    ok(px(s['letterSpacing']) == px(RAIZ_V['font-letter-spacing']), 'letterSpacing distinto en ' + n)
    ok(s['fontWeight'] in pesos, 'peso no está en C2: %s' % n)
    usados_css |= {'font-size-' + n, 'line-height-' + n}
    ok(isinstance(s['fontWeight'], int) and str(s['fontWeight']) == RAIZ_V['font-weight-' + pesos[s['fontWeight']]], 'peso distinto en ' + n)
    ok('tabular-nums' in s['usage'] or n == 'title-1', 'sin nota de cifras tabulares: ' + n)
# pesos: desde la revisión ds-r1 (M-02) van como valor en la familia fontWeight, no solo en los estilos
ok('fontWeight' in ds and [t['name'] for t in ds['fontWeight']['tokens']] == ['font-weight-' + w for w in ('regular', 'medium', 'semibold', 'bold')], 'familia fontWeight ausente o con otros nombres')
valor_peso = {t['name']: t['value'] for t in ds.get('fontWeight', {}).get('tokens', [])}
font_como_valor = {'size': 0, 'line-height': 0, 'weight': 0, 'letter-spacing': 0, 'family': 0}
# los 18 tokens font.* de C2 estan en estilos/pila
for r, n in hojas:
    if r[0] != 'font':
        continue
    if r[1] == 'size':
        if ok(any(s['name'] == r[2] and s['fontSize'] == n['$value'] for s in estilos), 'font.size.%s no está en estilos' % r[2]): font_como_valor['size'] += 1
    elif r[1] == 'line-height':
        if ok(any(s['name'] == r[2] and s['lineHeight'] == n['$value'] for s in estilos), 'font.line-height.%s no está en estilos' % r[2]): font_como_valor['line-height'] += 1
    elif r[1] == 'weight':
        if ok(valor_peso.get('font-weight-' + r[2]) == str(n['$value']), 'font.weight.%s no está como valor en fontWeight' % r[2]): font_como_valor['weight'] += 1
    elif r[1] == 'letter-spacing' or r == ['font', 'letter-spacing']:
        if ok(all(s['letterSpacing'] == n['$value'] for s in estilos), 'letter-spacing distinto'): font_como_valor['letter-spacing'] += 1
    elif r[1] == 'family':
        if ok(ds['type']['families'].get(r[2]) == fam_css, 'font.family.%s distinta' % r[2]): font_como_valor['family'] += 1
# fuentes
fonts_css = re.findall(r'@font-face\s*{[^}]*?src: url\("[^"]*?([^/"]+)"\)[^}]*?font-weight: ([\d ]+);\s*font-style: (\w+);', css_texto)
fuentes_ok = 0
for f in ds['type']['fonts']:
    ruta = PROY / f['file']
    if ok(ruta.is_file(), 'fuente no está en fonts/: ' + f['file']):
        b = ruta.read_bytes()
        ok(b[:4] == b'wOF2', 'firma de fuente no es woff2: ' + f['file'])
        ok(len(b) <= 1 << 20, 'fuente > 1 MB: ' + f['file'])
        fuentes_ok += 1
        ok((Path(f['file']).name, f['weight'], f['style']) in fonts_css, 'fuente distinta de @font-face de tokens.css: %s' % f['file'])
        ok(not re.search(r'["\']', f['family']), 'family con comillas')
ok(fuentes_ok == len(ds['type']['fonts']) == len(fonts_css), 'fuentes listadas, presentes y en tokens.css no coinciden')

# ---------------- 7. conteos por familia: C2 contra tokens.json ----------------
def cuenta(fam):
    return len(ds[fam]['tokens'])
c_ui = len([1 for t in ds['color']['tokens'] if t['name'].startswith('ui-') and t['name'] != 'ui-overlay-scrim'])
c_color = len([1 for t in ds['color']['tokens'] if t['name'].startswith('color-')])
tabla = [  # familia C2, tokens C2, tokens DS, regla
    ('color', por_raiz['color'], c_color, 'uno a uno'),
    ('ui', por_raiz['ui'], c_ui, '2 (light, dark) -> 1 con {light, dark}'),
    ('overlay', por_raiz['overlay'], 1, '2 (light, dark) -> ui-overlay-scrim'),
    ('font', por_raiz['font'], sum(font_como_valor.values()), '%(size)d tamaños + %(line-height)d interlíneas + %(letter-spacing)d espaciado en los estilos de type; %(weight)d pesos como valor en fontWeight; %(family)d pila' % font_como_valor),
    ('space', por_raiz['space'], cuenta('spacing'), 'spacing'),
    ('radius', por_raiz['radius'], cuenta('radius'), 'radius'),
    ('shadow', por_raiz['shadow'], cuenta('shadow'), 'shadow'),
    ('border', por_raiz['border'], cuenta('border'), 'border'),
    ('focus', por_raiz['focus'], cuenta('focus'), 'focus'),
    ('size', por_raiz['size'], cuenta('size'), 'size'),
    ('layout', por_raiz['layout'], cuenta('layout'), 'layout'),
    ('icon', por_raiz['icon'], cuenta('icon'), 'icon (con $root)'),
    ('motion', por_raiz['motion'], cuenta('motion'), 'motion'),
    ('breakpoint', por_raiz['breakpoint'], cuenta('breakpoint'), 'breakpoint'),
    ('viz', por_raiz['viz'], cuenta('viz'), 'viz'),
]
esperado_ds = {'ui': por_raiz['ui'] // 2, 'overlay': 1}
filas = []
for fam, a, b, regla in tabla:
    esp = esperado_ds.get(fam, a)
    if fam == 'font':
        esp = 18
    ok(b == esp, 'conteo %s: C2=%d, esperado %d, tokens.json=%d' % (fam, a, esp, b))
    filas.append('  %-11s C2 %3d  ->  tokens.json %3d   (%s)' % (fam, a, b, regla))
ok(sum(t[1] for t in tabla) == n_c2, 'suma de familias C2 distinta del total de hojas')
ok(n_c2 == 445, 'C2 trae %d tokens $value, se esperaban 445' % n_c2)

# ---------------- 8. cobertura de tokens.css ----------------
cubiertas = set(vistos) | usados_css
sin_cubrir = sorted(set(RAIZ_V) - cubiertas)
CONOCIDAS = {'ui-shadow-1', 'ui-shadow-2'}          # copias temáticas de shadow-1/2: en oscuro valen none (ya en shadow-1/2.dark)
ok(set(sin_cubrir) == CONOCIDAS, 'variables de tokens.css sin token en tokens.json: %s' % sin_cubrir)
sobrantes = sorted(n for n in vistos if n not in RAIZ_V and not n.startswith('font-'))
ok(not sobrantes, 'tokens sin variable en tokens.css: %s' % sobrantes)
ok(ds_texto.count('"name"') >= 1, 'tokens.json sin nombres')

# ---------------- 9. P3-12: cifras canon ----------------
canon = CANON.findall(ds_texto)
ok(not canon, 'cifras canon del reporte en tokens.json: %s' % canon)
canon_estilos = [s['sample'] for s in estilos if CANON.search(s['sample'])]
ok(not canon_estilos, 'cifras canon en samples')

# ---------------- salida ----------------
print('Verificación de %s contra C2 %s y tokens.css' % ((PROY / 'tokens.json').relative_to(RAIZ.parent), C2.name))
print('Conteo por familia (tokens $value de C2 contra tokens.json):')
print('\n'.join(filas))
print('C2: %d tokens · tokens.json: %d tokens en listas (color %d + otras %d) + %d estilos de tipo + 1 pila + 2 fuentes' % (
    n_c2, total_tokens, len(tokens_color), total_tokens - len(tokens_color), len(estilos)))
print('Comparaciones de color contra tokens.css (token × tema): %d' % n_hex)
print('Valores vacíos: %d · con var(/color-mix(/url(: %d · fuentes presentes: %d/%d' % (vacios, prohibidos, fuentes_ok, len(ds['type']['fonts'])))
print('Variables de tokens.css: %d (cubiertas por tokens.json o tipografía: %d; copias temáticas declaradas: %s)' % (
    len(RAIZ_V), len(RAIZ_V) - len(sin_cubrir), ', '.join(sorted(CONOCIDAS))))
for i in info:
    print(i)
if fallas:
    print('FALLAS (%d):' % len(fallas))
    for f in fallas[:60]:
        print('  - ' + f)
    print('RESULTADO: FALLA · %d fallas' % len(fallas))
    sys.exit(1)
print('RESULTADO: OK · %d tokens C2 -> %d en listas de tokens.json + %d estilos de tipo (18 font.*) + 2 fuentes; 0 fallas' % (n_c2, total_tokens, len(estilos)))
