#!/usr/bin/env python3
"""Comprueba que cada cita «...» de las lineas de trazabilidad (seccion 5) AGREGADAS en esta ronda
existe en plan.md, decisiones.md o protocolo-agentes.md (fragmentos separados por … en orden)."""
import re, glob, os, sys
R = '/home/deepsec/finanzas-hogar/docs/'
norm = lambda s: re.sub(r'\s+', ' ', s.replace('*', '').replace('`', '').replace(' ', ' ')).strip().lower()
FUENTES = {n: norm(open(R + p, encoding='utf-8').read()) for n, p in
           (('plan', 'plan/plan.md'), ('decisiones', 'plan/decisiones.md'), ('protocolo', 'plan/protocolo-agentes.md'))}

def citas(linea):
    out, depth, ini = [], 0, None
    for i, ch in enumerate(linea):
        if ch == '«':
            if depth == 0: ini = i + 1
            depth += 1
        elif ch == '»' and depth:
            depth -= 1
            if depth == 0: out.append(linea[ini:i])
    return out

def seccion5(t):
    m = re.search(r'^## 5\. Trazabilidad\s*$', t, re.M)
    ini = m.end(); m2 = re.search(r'^## \d+\. ', t[ini:], re.M)
    return t[ini:ini + m2.start()] if m2 else t[ini:]

ORIG = sys.argv[1] if len(sys.argv) > 1 else None
total = bad = 0
for f in sorted(glob.glob(R + 'criterios/*.md')):
    nom = os.path.basename(f)
    if nom == 'INDICE.md': continue
    nuevo = seccion5(open(f, encoding='utf-8').read()).splitlines()
    viejo = set(seccion5(open(ORIG + '/' + nom, encoding='utf-8').read()).splitlines()) if ORIG else set()
    for ln in nuevo:
        if ln in viejo or not ln.startswith('- '): continue
        for c in citas(ln):
            partes = [norm(p) for p in re.split(r'…|\.\.\.', c) if norm(p)]
            # solo se verifican citas de las fuentes del plan (las de nombres internos como «Local» o «sin precio» no)
            if not re.search(r'plan\.md|decisiones|protocolo', ln): continue
            total += 1
            ok_fuentes = []
            for n, src in FUENTES.items():
                pos = 0; ok = True
                for p in partes:
                    k = src.find(p, pos)
                    if k < 0: ok = False; break
                    pos = k + len(p)
                if ok: ok_fuentes.append(n)
            if not ok_fuentes:
                bad += 1; print(f'  SIN FUENTE en {nom}: «{c[:100]}»')
print(f'citas verificadas: {total} | sin fuente: {bad}')
