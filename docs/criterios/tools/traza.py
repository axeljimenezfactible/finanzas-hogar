#!/usr/bin/env python3
"""Verifica la estructura de docs/criterios y la trazabilidad (seccion 5) de cada criterio.
Uso: python3 -I traza.py [--lista]   (sin argumentos: resumen; --lista: imprime los IDs sin traza)"""
import re, sys, glob, os
BASE = '/home/deepsec/finanzas-hogar/docs/criterios'
ARCH = [f for f in sorted(glob.glob(BASE + '/*.md')) if os.path.basename(f) != 'INDICE.md']
ID = re.compile(r'\b((?:WF4|P3|[SUC][0-9])-[0-9]{2})\b')

def seccion(txt, n):
    m = re.search(r'^## %d\. .*?$' % n, txt, re.M)
    if not m:
        return ''
    ini = m.end()
    m2 = re.search(r'^## \d+\. ', txt[ini:], re.M)
    return txt[ini:ini + m2.start()] if m2 else txt[ini:]

def expandir(txt):
    """IDs citados en un texto, expandiendo rangos 'X-01...X-04', 'X-01 a X-04' y 'X-01-X-04'."""
    ids = set(ID.findall(txt))
    for a, b in re.findall(r'\b((?:WF4|P3|[SUC][0-9])-[0-9]{2})\s*(?:…|\.\.\.|a|–|—)\s*((?:WF4|P3|[SUC][0-9])-[0-9]{2})\b', txt):
        pa, na = a.rsplit('-', 1); pb, nb = b.rsplit('-', 1)
        if pa == pb:
            for k in range(int(na), int(nb) + 1):
                ids.add('%s-%02d' % (pa, k))
    return ids

total = 0; sin = []; estruct_mal = []; sin_ver = []; sin_evi = []
for f in ARCH:
    t = open(f, encoding='utf-8').read()
    nom = os.path.basename(f)
    for n in range(1, 6):
        if not re.search(r'^## %d\. ' % n, t, re.M):
            estruct_mal.append((nom, n))
    s2 = seccion(t, 2); s5 = seccion(t, 5)
    defs = re.findall(r'^\*\*((?:WF4|P3|[SUC][0-9])-[0-9]{2}) \((alta|media|baja)\)', s2, re.M)
    ids = [d[0] for d in defs]
    total += len(ids)
    ver = len(re.findall(r'^- Verifica:', s2, re.M)); evi = len(re.findall(r'^- Evidencia:', s2, re.M))
    if ver != len(ids): sin_ver.append((nom, len(ids), ver))
    if evi != len(ids): sin_evi.append((nom, len(ids), evi))
    cit = expandir(s5)
    for i in ids:
        if i not in cit:
            sin.append((nom, i))
print('archivos:', len(ARCH), '| criterios definidos:', total)
print('secciones 1-5 faltantes:', estruct_mal)
print('archivos con Verifica != criterios:', sin_ver, '| con Evidencia != criterios:', sin_evi)
print('criterios sin linea en su seccion 5 (con rangos expandidos):', len(sin))
if '--lista' in sys.argv:
    por = {}
    for nom, i in sin: por.setdefault(nom, []).append(i)
    for k, v in por.items(): print(' ', k, ' '.join(v))
