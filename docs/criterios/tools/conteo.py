#!/usr/bin/env python3
"""Conteo propio de docs/criterios (P1c, criterios 1.2): criterios por archivo y severidad, version de cada
archivo y comparacion con las tablas de INDICE.md (§2 y §10) y con las cifras de su §1. Solo lectura."""
import re, glob, os, sys
D = '/home/deepsec/finanzas-hogar/docs/criterios/'
ORDEN = ['S1', 'S2', 'S3', 'S4', 'S5', 'S6', 'S7', 'S8', 'WF4-operacion', 'U1', 'U2', 'U3', 'U4', 'U5', 'C1', 'C2', 'C3', 'P3-claude-design']
ARCH = {f: os.path.join(D, f + '.md') for f in ORDEN}
real, errores = {}, []
for f, p in ARCH.items():
    t = open(p, encoding='utf-8').read()
    ids = re.findall(r'^\*\*([A-Z0-9]+-\d+) \((alta|media|baja)\)', t, re.M)
    if len({i for i, _ in ids}) != len(ids): errores.append(f'{f}: IDs repetidos')
    sev = {k: sum(1 for _, s in ids if s == k) for k in ('alta', 'media', 'baja')}
    m = re.search(r'^Versi[oó]n (\d+\.\d+)', t, re.M)
    real[f] = (len(ids), sev['alta'], sev['media'], sev['baja'], m.group(1) if m else '?')
    if not re.search(r'Cambios de la versi[oó]n ' + re.escape(real[f][4]), t): errores.append(f'{f}: falta la lista «Cambios de la versión {real[f][4]}»')
tot = [sum(real[f][i] for f in ORDEN) for i in range(4)]
print('archivo                 total alta media baja version')
for f in ORDEN: print(f'{f:<22} {real[f][0]:>5} {real[f][1]:>5} {real[f][2]:>5} {real[f][3]:>4}   {real[f][4]}')
print(f'TOTAL {tot[0]} (alta {tot[1]} / media {tot[2]} / baja {tot[3]}) | archivos: {len(ORDEN)}')
ind = open(D + 'INDICE.md', encoding='utf-8').read()
for f in ORDEN:
    n, a, me, b, v = real[f]
    fila = re.search(r'^\| `' + re.escape(f + '.md') + r'` \|.*\| (\d+) \((\d+) / (\d+)\) \|\s*$', ind, re.M)
    if not fila: errores.append(f'INDICE §2: sin fila de {f}'); continue
    if tuple(map(int, fila.groups())) != (n, a, me): errores.append(f'INDICE §2 {f}: {fila.groups()} contra {(n, a, me)}')
    f10 = re.search(r'^\| `' + re.escape(f + '.md') + r'` \| (\d+) \((\d+) / (\d+) / (\d+)\) \| (\S+) \|', ind, re.M)
    if not f10: errores.append(f'INDICE §10: sin fila de {f}'); continue
    if (tuple(map(int, f10.groups()[:4])), f10.group(5)) != ((n, a, me, b), v): errores.append(f'INDICE §10 {f}: {f10.groups()} contra {(n, a, me, b, v)}')
if f'**{tot[0]} ({tot[1]} / {tot[2]})**' not in ind: errores.append(f'INDICE §2: total distinto de {tot[0]} ({tot[1]} / {tot[2]})')
if f'los {tot[0]} criterios de la tabla' not in ind or f'{tot[1]} alta y {tot[2]} media' not in ind: errores.append('INDICE §1: texto de totales distinto')
print('INDICE contra conteo propio: ' + ('0 diferencias' if not errores else f'{len(errores)} diferencias'))
for e in errores: print('  DIFERENCIA', e)
sys.exit(1 if errores else 0)
