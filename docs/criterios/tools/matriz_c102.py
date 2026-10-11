#!/usr/bin/env python3
"""Comprobacion del texto nuevo de C1-02 / U1-03 sobre los artefactos actuales (solo lectura).
Uso: uv run --no-project --with pyyaml python3 -I matriz_c102.py"""
import re, sys, yaml
R = '/home/deepsec/finanzas-hogar/'
u1 = open(R + 'docs/design/U1-arquitectura.md', encoding='utf-8').read()
c1 = open(R + 'docs/contratos/C1-datos.md', encoding='utf-8').read()
api = yaml.safe_load(open(R + 'docs/contratos/openapi.yaml', encoding='utf-8'))

# --- U1: necesidades definidas (seccion 3, antes de 3.29) y cuales declara Local
s3 = u1[u1.index('\n## 3. '):u1.index('### 3.29 ')]
defs = re.findall(r'^\| (N-\d{3}) \| ([LE]) \|(.*)$', s3, re.M)
N = {d[0] for d in defs}
local = {d[0] for d in defs if d[2].rstrip().rstrip('|').strip().endswith('Local')}
# --- U1 3.29: pantalla -> necesidades (se expanden rangos "N-020 a N-029")
def expandir(txt):
    ids = set(re.findall(r'N-\d{3}', txt))
    for a, b in re.findall(r'(N-\d{3}) a (N-\d{3})', txt):
        ids |= {'N-%03d' % k for k in range(int(a[2:]), int(b[2:]) + 1)} & N
    return ids
m329 = u1[u1.index('### 3.29 '):u1.index('### 3.30 ')]
pant = {}
for fila in re.findall(r'^\| (Global[^|]*|P-\d\d(?:\.\d\d)?) \|(.*)$', m329, re.M):
    clave = 'Global' if fila[0].startswith('Global') else fila[0]
    pant.setdefault(clave, set()).update(expandir(fila[1]))
# P-26 (contenedor) = union de sus secciones
pant['P-26'] = set().union(*[v for k, v in pant.items() if k.startswith('P-26.')])
P_inv = set(re.findall(r'^\| (P-\d\d) \| ', u1[u1.index('### 2.1 '):u1.index('### 2.2 ')], re.M))
S_inv = set(re.findall(r'^\| (P-26\.\d\d) \| ', u1[u1.index('### 2.3 '):u1.index('### 2.4 ')], re.M))
# --- openapi: operaciones y x-necesidades
ops = {}
for ruta, item in api['paths'].items():
    for met, op in item.items():
        if met in ('get', 'post', 'put', 'patch', 'delete'):
            ops[op['operationId']] = (met.upper() + ' ' + ruta, set(op.get('x-necesidades') or []))
con_op = set().union(*[v[1] for v in ops.values()])
# --- C1 matriz (seccion 10): filas pantalla | N | tipo | texto | operaciones | estado
m10 = c1[c1.index('## 10. '):c1.index('## 11. ')]
filas = re.findall(r'^\| (Global|P-\d\d(?:\.\d\d)?) \| (N-\d{3}) \| ([LE]) \|.*\| ([^|]*) \|$', m10, re.M)
estado = {}
for p, n, t, e in filas:
    estado.setdefault(n, set()).add(e.strip().split(' ')[0])
sin_op_ni_local = sorted(N - con_op - local)
print(f'U1: necesidades N-xxx definidas = {len(N)} | pantallas P-xx = {len(P_inv)} | secciones P-26.xx = {len(S_inv)} | Local declaradas = {sorted(local)}')
print(f'openapi: operaciones = {len(ops)} | N-xxx citadas en x-necesidades = {len(con_op)}')
print(f'N-xxx sin operacion ni Local: {len(sin_op_ni_local)} {sin_op_ni_local}')
pant_mal = []
for p in sorted(P_inv | S_inv | {'Global'}):
    necs = pant.get(p, set())
    if not necs:
        pant_mal.append((p, 'sin necesidades en la matriz de U1'))
    elif not (necs & con_op) and not necs <= local:
        pant_mal.append((p, 'ninguna de sus N-xxx tiene operacion'))
print(f'P-xx (con secciones P-26.xx y Global) sin operacion ni Local: {len(pant_mal)} {pant_mal}')
sin_n = sorted(v[0] for k, v in ops.items() if not v[1])
print(f'operaciones sin N-xxx: {len(sin_n)} {sin_n}')
print(f'N-xxx citadas por el contrato que no existen en U1: {sorted(con_op - N)}')
c1_n = {n for _, n, _, _ in filas}
print(f'matriz C1 (seccion 10): filas = {len(filas)} | N-xxx distintas = {len(c1_n)} | N-xxx de U1 ausentes en la matriz C1: {sorted(N - c1_n)}')
print('estados en la matriz C1:', {e: sum(1 for v in estado.values() if e in v) for e in sorted({x for v in estado.values() for x in v})})
