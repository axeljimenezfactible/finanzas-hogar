#!/usr/bin/env python3
"""Cuenta los criterios por archivo y severidad, y revisa que cada uno tenga «Verifica» y «Evidencia»
y que aparezca en la seccion 5 (rangos expandidos). Solo lectura.
Uso: python3 -I cuenta.py [archivo.md ...]   (sin argumentos: S1 a S8 y WF4-operacion)
Codigo de salida: 0 si todo completo; 1 si falta algo."""
import re, sys, os
BASE = '/home/deepsec/finanzas-hogar/docs/criterios/'
DEF = ['S1.md', 'S2.md', 'S3.md', 'S4.md', 'S5.md', 'S6.md', 'S7.md', 'S8.md', 'WF4-operacion.md']
ID = r'(?:WF4|P3|[SUC][0-9])-[0-9]{2}'
def seccion(t, n):
    m = re.search(r'^## %d\. .*?$' % n, t, re.M)
    if not m: return ''
    ini = m.end(); m2 = re.search(r'^## \d+\. ', t[ini:], re.M)
    return t[ini:ini + m2.start()] if m2 else t[ini:]
def citados(t):
    ids = set(re.findall(r'\b(%s)\b' % ID, t))
    for a, b in re.findall(r'\b(%s)\s*(?:…|\.\.\.|a|–|—)\s*(%s)\b' % (ID, ID), t):
        pa, na = a.rsplit('-', 1); pb, nb = b.rsplit('-', 1)
        if pa == pb: ids |= {'%s-%02d' % (pa, k) for k in range(int(na), int(nb) + 1)}
    return ids
archivos = [os.path.basename(a) for a in sys.argv[1:]] or DEF
tot = {'alta': 0, 'media': 0, 'baja': 0}; mal = []
print('%-20s %5s %5s %5s %5s %5s %5s %5s' % ('archivo', 'total', 'alta', 'media', 'baja', 'Verif', 'Evid', 'sin§5'))
for nom in archivos:
    t = open(BASE + nom, encoding='utf-8').read()
    s2 = seccion(t, 2); s5 = seccion(t, 5)
    defs = re.findall(r'^\*\*(%s) \((alta|media|baja)\)' % ID, s2, re.M)
    sev = {k: sum(1 for d in defs if d[1] == k) for k in tot}
    for k in tot: tot[k] += sev[k]
    v = len(re.findall(r'^- Verifica:', s2, re.M)); e = len(re.findall(r'^- Evidencia:', s2, re.M))
    cit = citados(s5); sin = [d[0] for d in defs if d[0] not in cit]
    secs = [n for n in range(1, 6) if not re.search(r'^## %d\. ' % n, t, re.M)]
    print('%-20s %5d %5d %5d %5d %5d %5d %5d' % (nom, len(defs), sev['alta'], sev['media'], sev['baja'], v, e, len(sin)))
    if v != len(defs) or e != len(defs) or sin or secs: mal.append((nom, len(defs), v, e, sin, secs))
print('TOTAL %d (alta %d / media %d / baja %d)' % (sum(tot.values()), tot['alta'], tot['media'], tot['baja']))
print('incompletos:', mal if mal else 'ninguno')
sys.exit(1 if mal else 0)
