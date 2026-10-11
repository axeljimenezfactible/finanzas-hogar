import re, glob, os
B = '/home/deepsec/finanzas-hogar/docs/criterios/'
ID = re.compile(r'\b((?:WF4|P3|[SUC][0-9])-[0-9]{2})\b')
definidos = set(); citados = {}
for f in sorted(glob.glob(B + '*.md')):
    t = open(f, encoding='utf-8').read()
    for m in re.finditer(r'^\*\*((?:WF4|P3|[SUC][0-9])-[0-9]{2}) \((?:alta|media|baja)\)', t, re.M):
        definidos.add(m.group(1))
    for i in ID.findall(t):
        citados.setdefault(i, set()).add(os.path.basename(f))
# IDs de otros dominios que no son criterios: R-nn, K-nn no coinciden con el patron; U1-02 etc. si
inexistentes = {i: sorted(v) for i, v in citados.items() if i not in definidos}
print('criterios definidos:', len(definidos), '| IDs distintos citados:', len(citados), '| citados que no existen:', inexistentes)
