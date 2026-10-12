#!/usr/bin/env python3
"""Arma docs/design/maquetas/lienzo/project/canvas.json desde los fragmentos por área.
Uso: python3 -I armar_canvas.py [--ds-url URL] [--ds-version ID] [--at RFC3339]
Filas por área (orden fijo), 80 px entre artboards y 120 entre filas, nota title1 por fila, páginas por área."""
import json, sys, os, glob
ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '..', 'maquetas', 'lienzo')
ROOT = os.path.abspath(ROOT)
args = sys.argv[1:]
def opt(name, default=None):
    return args[args.index(name) + 1] if name in args else default
ds_url = opt('--ds-url', 'https://claude.ai/artifact/D6r3woEKMBvpWiK6gzr3Ra')
ds_version = opt('--ds-version', None)
at = opt('--at', '2026-10-12T07:30:00Z')
ORDEN = ['shell', 'inicio', 'gastos', 'presupuesto', 'comidas', 'compras', 'creditos', 'salud-reportes', 'ajustes-a', 'ajustes-b']
NOMBRES = {'shell': 'Shell: entrar, barra lateral, hojas', 'inicio': 'Inicio', 'gastos': 'Gastos', 'presupuesto': 'Presupuesto', 'comidas': 'Comidas', 'compras': 'Compras y tiendas', 'creditos': 'Créditos', 'salud-reportes': 'Salud y Reportes', 'ajustes-a': 'Ajustes 1 a 7', 'ajustes-b': 'Ajustes 8 a 14'}
PAGINA = {'ajustes-a': 'ajustes', 'ajustes-b': 'ajustes'}
boards, order, notes, pages = {}, [], {}, []
y = 0
vistos = set()
for area in ORDEN:
    f = os.path.join(ROOT, 'fragmentos', area + '.json')
    frag = json.load(open(f, encoding='utf-8'))
    page = PAGINA.get(area, area)
    if page not in vistos:
        pages.append({'id': page, 'name': NOMBRES.get(page, page) if page != 'ajustes' else 'Ajustes'})
        vistos.add(page)
    notes[f'titulo-{area}'] = {'x': 0, 'y': y - 300, 'text': NOMBRES[area], 'kind': 'title1', 'maxW': 1840, 'page': page}
    x = 0
    fila_h = 0
    for ruta in frag['order']:
        b = dict(frag['boards'][ruta])
        if not os.path.exists(os.path.join(ROOT, 'project', ruta)):
            print('FALTA archivo', ruta); sys.exit(1)
        if ruta in boards:
            print('DUPLICADO', ruta); sys.exit(1)
        b['x'] = x; b['y'] = y; b['page'] = page
        b.setdefault('is_interactive', False)
        boards[ruta] = b
        order.append(ruta)
        x += int(b['w']) + 80
        fila_h = max(fila_h, int(b['h']))
    y += fila_h + 120 + 300  # espacio para el título de la fila siguiente
entrada = 'shell/shell-barra-lateral-e.dc.html'
if entrada not in boards:
    print('FALTA entrada', entrada); sys.exit(1)
order.remove(entrada); order.insert(0, entrada)
canvas = {
    'v': 3,
    'createdOnFiles': {'v': 1, 'at': at},
    'title': 'Finanzas Hogar · Pantallas',
    'launch': {'view': 'canvas', 'page': 'shell'},
    'pages': pages,
    'boards': boards,
    'order': order,
    'notes': notes,
    'designSystems': [{'title': 'Finanzas Hogar', 'namespace': 'finanzas-hogar', 'artifact': ds_url, 'version': ds_version, 'copiedAt': at}],
}
# comprobaciones del tipo: stems únicos, nombres válidos
import re
stems = [os.path.splitext(os.path.basename(r).replace('.dc', ''))[0] for r in order]
assert len(stems) == len(set(stems)), 'stems repetidos'
for r in order:
    for seg in r.split('/'):
        assert re.match(r'^[A-Za-z0-9_][A-Za-z0-9_.\-]*$', seg), f'segmento inválido {seg}'
archivos = set(os.path.relpath(p, os.path.join(ROOT, 'project')) for p in glob.glob(os.path.join(ROOT, 'project', '**', '*.dc.html'), recursive=True))
sin_listar = archivos - set(order)
if sin_listar: print('ARTBOARDS SIN LISTAR', sorted(sin_listar)); sys.exit(1)
out = os.path.join(ROOT, 'project', 'canvas.json')
json.dump(canvas, open(out, 'w', encoding='utf-8'), ensure_ascii=False, indent=1)
print(f'canvas.json: {len(order)} artboards, {len(pages)} páginas, {len(notes)} notas; alto total {y}px')
