#!/usr/bin/env python3 -I
"""U5-03 / U5-04 / U5-07 / U5-08 / U5-10 · Resumen por combinación a partir de los JSON de docs/design/capturas-u5/.
Uso: python3 -I frontend/pruebas-u5/resumen_u5.py"""
import json, glob, os, re, collections
D = '/home/deepsec/finanzas-hogar/docs/design/capturas-u5'
COMB = ['android-claro', 'android-oscuro', 'escritorio-claro', 'escritorio-oscuro']
def comb(f): return re.search(r'-(android|escritorio)-(claro|oscuro)\.json$', f).expand(r'\1-\2')
axe = collections.defaultdict(lambda: collections.Counter()); serias = []
for f in sorted(glob.glob(f'{D}/axe/*.json')):
    j = json.load(open(f)); c = comb(f); axe[c]['pantallas'] += 1
    axe[c]['violaciones'] += len(j['violations']); axe[c]['incompletos'] += len(j['incomplete'])
    for v in j['violations']:
        if v['impact'] in ('serious', 'critical'):
            axe[c]['serias_o_criticas'] += 1
            serias.append((j['pantalla'], c, v['id'], v['impact'], [n['target'] for n in v['nodes']][:3]))
print('U5-03 axe-core (wcag2a, wcag2aa, wcag21a, wcag21aa, wcag22aa)')
for c in COMB: print(f'  {c:18} pantallas {axe[c]["pantallas"]:3}  serias/críticas {axe[c]["serias_o_criticas"]:3}  otras violaciones {axe[c]["violaciones"]-axe[c]["serias_o_criticas"]:3}  incompletos {axe[c]["incompletos"]}')
print('  total serias/críticas:', len(serias))
for s in serias: print('   ', s)
con = collections.defaultdict(lambda: collections.Counter()); bajos = []
for f in sorted(glob.glob(f'{D}/contraste/*.json')):
    j = json.load(open(f)); c = comb(f)
    con[c]['pantallas'] += 1; con[c]['evaluados'] += j['evaluados']; con[c]['bajo'] += j['bajo_umbral']; con[c]['indet'] += j['indeterminados']
    for b in j['bajos']: bajos.append((j['pantalla'], c, b['tipo']+':'+str(b.get('subtipo')), b['razon'], b['umbral'], b['fg'], b['bg'], b['nodo'][-70:], b['texto']))
print('U5-04 contraste en el DOM')
for c in COMB: print(f'  {c:18} pantallas {con[c]["pantallas"]:3}  pares (elementos) evaluados {con[c]["evaluados"]:6}  bajo umbral {con[c]["bajo"]:4}  indeterminados {con[c]["indet"]}')
print('  total bajo umbral:', len(bajos))
agr = collections.Counter((b[0], b[2], b[5], b[6], b[3]) for b in bajos)
print('  bajo umbral, agrupado (pantalla, tipo, fg, bg, razón) -> elementos, sobre todas las combinaciones:')
for k, v in sorted(agr.items())[:60]: print('   ', k, v)
H = json.load(open(f'{D}/hosts.json')); print('U5-08 hosts únicos:', H['unicos'], '· capturas con registro:', H['capturas_con_registro'])
M = json.load(open(f'{D}/evidencia/manifiesto.json'))['capturas']
malos = [k for k, v in M.items() if not v['tema']['coincide'] or v['tema']['data_theme'] is not None]
print('U5-10 capturas cuyo fondo del body != token C2 o con data-theme:', len(malos), malos[:5])
print('Capturas con aviso: cargando al capturar:', [k for k, v in M.items() if v['cargando_al_capturar']], '· error de pantalla:', [k for k, v in M.items() if v['error_de_pantalla']], '· h1 distinto del nombre:', [(k, v['h1']) for k, v in M.items() if v['h1_coincide'] is False][:10])
cif = collections.defaultdict(set)
for k, v in M.items():
    for c in v['cifras_visibles']: cif[c].add(v['pantalla'])
print('U5-02 cifras visibles de la semilla (pantallas):'); [print('  ', c, sorted(p)) for c, p in sorted(cif.items())]
u7 = collections.defaultdict(lambda: collections.Counter()); u7det = []
for k, v in M.items():
    r = v.get('u5_07')
    if r is None: u7[v['combinacion']]['sin_dato'] += 1; continue
    c = v['combinacion']; u7[c]['pantallas'] += 1; u7[c]['importes'] += r['importes']
    for campo in ('importes_fuera_de_formato', 'fechas_iso', 'ingles', 'claves_internas'):
        u7[c][campo] += len(r[campo])
        if r[campo]: u7det.append((k, campo, r[campo][:5]))
print('U5-07 texto visible (innerText de cada captura)')
for c in COMB: print(f'  {c:18} pantallas {u7[c]["pantallas"]:3}  importes {u7[c]["importes"]:5}  fuera de formato {u7[c]["importes_fuera_de_formato"]}  fechas ISO {u7[c]["fechas_iso"]}  inglés {u7[c]["ingles"]}  claves internas {u7[c]["claves_internas"]}  sin dato {u7[c]["sin_dato"]}')
for d in u7det: print('   ', d)
g = [(k, v.get('puntero_grueso')) for k, v in M.items() if v['combinacion'].startswith('android')]
print('Android con (any-pointer: coarse) al capturar:', sum(1 for _, x in g if x is True), 'de', len(g), '· sin puntero grueso:', [k for k, x in g if x is not True][:10])
