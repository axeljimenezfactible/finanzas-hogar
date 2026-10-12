#!/usr/bin/env python3 -I
"""U5-01 · Lee los P-xx y P-26.xx de docs/design/U1-arquitectura.md, comprueba los 4 PNG por pantalla (existe, tamaño != 0,
ancho = viewport declarado, desviación estándar de píxeles > 0) y escribe el resumen. Usa ImageMagick `identify` para leer píxeles.
Uso: python3 -I frontend/pruebas-u5/verificar_capturas.py"""
import re, os, subprocess, sys
R = '/home/deepsec/finanzas-hogar'
md = open(f'{R}/docs/design/U1-arquitectura.md', encoding='utf8').read()
ids = sorted(set(re.findall(r'\bP-\d{2}(?:\.\d{2})?\b', md)))
sec = [i for i in ids if '.' in i]; pan = [i for i in ids if '.' not in i]
T = len(ids)
print(f'pantallas P-xx: {len(pan)} · secciones P-26.xx: {len(sec)} · T = {T} · 4T = {4*T}')
VP = {'android': 412, 'escritorio': 1440}
falt = []; mal = []; ok = 0
for i in ids:
    for d in VP:
        for t in ('claro', 'oscuro'):
            f = f'{R}/docs/design/capturas-u5/{i}-{d}-{t}.png'
            if not os.path.isfile(f) or os.path.getsize(f) == 0: falt.append(os.path.basename(f)); continue
            w, h, sd = subprocess.run(['identify', '-format', '%w %h %[fx:standard_deviation]', f], capture_output=True, text=True).stdout.split()
            w, h, sd = int(w), int(h), float(sd)
            alto_min = 915 if d == 'android' else 900
            if w != VP[d] or h < alto_min or sd <= 0: mal.append(f'{os.path.basename(f)} {w}x{h} sd={sd}')
            else: ok += 1
print(f'PNG correctos: {ok} de {4*T} · faltantes: {len(falt)} · con dimensiones o píxeles mal: {len(mal)}')
for x in falt[:20]: print('  FALTA', x)
for x in mal[:20]: print('  MAL', x)
extra = [n for n in os.listdir(f'{R}/docs/design/capturas-u5') if n.endswith('.png')]
print(f'PNG en la carpeta: {len(extra)}')
sys.exit(0 if ok == 4*T and not falt and not mal else 1)
