#!/usr/bin/env python3
"""Reverifica las tablas R-01..R-85 (S2 y S3) y K-01..K-08 (S5) contra oraculo/model.json y datos/semilla/data.json.
Cada celda 'Valor esperado' se descompone en fechas y numeros; cada uno debe existir entre las hojas del nodo citado."""
import json, re, hashlib
from decimal import Decimal as D
R = '/home/deepsec/finanzas-hogar/'
M = json.load(open(R + 'oraculo/model.json')); DJ = json.load(open(R + 'datos/semilla/data.json'))
print('sha256 model.json:', hashlib.sha256(open(R + 'oraculo/model.json', 'rb').read()).hexdigest())
def hojas(x, out):
    if isinstance(x, dict):
        for v in x.values(): hojas(v, out)
    elif isinstance(x, list):
        for v in x: hojas(v, out)
    elif isinstance(x, bool) or x is None: pass
    elif isinstance(x, (int, float)): out.append(('n', D(repr(x))))
    else: out.append(('s', str(x)))
    return out
def resolver(raiz, ruta):
    nodo = raiz
    for parte in re.findall(r'[A-Za-z_0-9]+|\[[^\]]+\]', ruta):
        if parte.startswith('['):
            idx = parte[1:-1]
            if '..' in idx:
                a, b = idx.split('..'); nodo = nodo[int(a):int(b) + 1]
            else: nodo = nodo[int(idx)]
        else: nodo = nodo[parte]
    return nodo
def toks(cel):
    t = []
    for m in re.finditer(r'\d{4}-\d{2}-\d{2}|\d[\d,]*\.\d+|\d+', cel):
        s = m.group(0)
        t.append(('f', s) if '-' in s else ('n', D(s.replace(',', ''))))
    return t
def chequea(fila_id, ruta, esp):
    raiz = M; ruta = ruta.strip('`').strip()
    if ruta.startswith('data.json:'): raiz = DJ; ruta = ruta[len('data.json:'):]
    partes = [p.strip() for p in ruta.split(' / ')]
    nodos = []
    base = partes[0]
    nodos.append(resolver(raiz, base))
    for p in partes[1:]:
        if p.startswith('.'):
            padre = base.rsplit('.', 1)[0]; nodos.append(resolver(raiz, padre + p))
        else: nodos.append(resolver(raiz, p))
    h = []
    for n in nodos: hojas(n, h)
    nums = [v for k, v in h if k == 'n']; strs = {v for k, v in h if k == 's'}
    malos = []
    for d in re.findall(r'despensa (\d)', esp):
        if not any('despensa %s' % d in s for s in strs): malos.append('despensa ' + d)
    esp = re.sub(r'despensa \d', '', esp)
    for k, v in toks(esp):
        if k == 'f':
            if v not in strs: malos.append(v)
        else:
            if not any(abs(v - n) <= D('0.0000000001') for n in nums): malos.append(str(v))
    return malos
tot = bad = 0
for archivo in ('S2.md', 'S3.md'):
    for ln in open(R + 'docs/criterios/' + archivo, encoding='utf-8'):
        m = re.match(r'\| (R-\d\d) \| (.*?) \| (`[^|]*`) \| (.*?) \| (S2|S3) \|', ln)
        if not m: continue
        rid, desc, ruta, esp, motor = m.groups()
        tot += 1
        if rid == 'R-67':                               # 17 eventos
            ok = len(M['credito']['eventos']) == 17; malos = [] if ok else ['n eventos']
        elif rid == 'R-04':                             # razon sin redondear
            malos = [] if abs(D(repr(M['parametros']['tasa_pct'])) - D('25.8684763290')) < D('1e-9') else ['tasa_pct']
        else:
            malos = chequea(rid, ruta, esp)
        if malos: bad += 1; print('  DIFIERE', archivo, rid, ruta, malos)
    print(archivo, 'filas R verificadas acumuladas:', tot)
# K-01..K-08 (S5) contra data.json
T = DJ['totales']
K = {'K-01': T['bodega'], 'K-02': T['walmart'], 'K-03': round(T['walmart'] - T['bodega'], 2), 'K-04': T['bodega_identicos'],
     'K-05': T['walmart_identicos'], 'K-06': round(T['walmart_identicos'] - T['bodega_identicos'], 2)}
esperado = {'K-01': 3094.50, 'K-02': 3233.00, 'K-03': 138.50, 'K-04': 2773.50, 'K-05': 2874.00, 'K-06': 100.50}
for k, v in esperado.items():
    ok = abs(K[k] - v) < 1e-9; tot += 1
    if not ok: bad += 1; print('  DIFIERE', k, K[k], v)
for k, ok in (('K-07', (T['articulos_bodega'], T['articulos_walmart'], T['articulos_identicos']) == (61, 60, 55)),
              ('K-08', (len(DJ['items_identicos']), len(DJ['sustitutos'])) == (33, 3))):
    tot += 1
    if not ok: bad += 1; print('  DIFIERE', k)
ks = [l for l in open(R + 'docs/criterios/S5.md', encoding='utf-8') if re.match(r'\| K-\d\d \|', l)]
print('S5 filas K en el archivo:', len(ks), '| K-01..K-08 contra data.json: OK' if bad == 0 else '')
print(f'TOTAL comprobadas: {tot} | con diferencias: {bad}')
