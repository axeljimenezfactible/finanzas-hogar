#!/usr/bin/env python3
"""C3 · matriz pantalla → operación → estado (integrador frontend/C3, WF-3a).

Fuentes (todas locales, sin red): docs/design/U1-arquitectura.md (§2.1, §2.3, §3.29), docs/contratos/openapi.yaml
(paths, operationId, x-necesidades, x-agente-externo), backend/internal/api/routes.go (rutas registradas),
frontend/src (llamadas del cliente: pedir/pedirPagina/todas/encolar/usarConsulta, armarUrl y fetch literales; sin mocks
ni *.gen.ts), un reporte JSON de Playwright, la salida -v de «go test -run TestContrato_Operaciones» y las capturas de U5
(docs/design/capturas-u5/<P-xx>-{android,escritorio}-{claro,oscuro}.png; columna «capturado», C3-05 / U5-11; WF-3b c3-matriz).

Uso (desde la raíz del repositorio; siempre con python3 -I):
  python3 -I frontend/scripts/c3-integracion.py --pw-json R.json --go-log G.log --escribir   # regenera la matriz
  python3 -I frontend/scripts/c3-integracion.py --pw-json R.json --go-log G.log              # C3-02/C3-03 + compara
Sale con 1 si hay rutas fuera del contrato, sin registrar, pantallas u operaciones faltantes o la matriz difiere.
"""
import argparse, json, re, struct, sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[2]
FRONT = RAIZ / 'frontend'
DOC = RAIZ / 'docs/contratos/C3-integracion.md'
INI, FIN = '<!-- c3:matriz:inicio (generado por frontend/scripts/c3-integracion.py; no editar a mano) -->', '<!-- c3:matriz:fin -->'
METODOS = ('get', 'post', 'put', 'patch', 'delete')
LOCALES = {'N-003', 'N-301', 'N-317'}  # U1 §3.29: viven en el teléfono, sin API
# Excepciones técnicas PROVISIONALES (bloqueo C3-matriz, punto 3c; decisión del orquestador 2026-10-12): operaciones con
# pantalla en U1 que el frontend no llama. Solo se admiten si TODAS sus necesidades las cubre otra operación que la
# pantalla sí llama (el script lo comprueba) y con la justificación escrita aquí. Una operación no llamada y no listada
# aquí hace fallar el script.
# Capturas de U5 (U5-01): 4 por pantalla, PNG no vacío, ancho del viewport con factor de escala 1 (412 Android, 1440 escritorio).
COMBOS = (('android', 'claro'), ('android', 'oscuro'), ('escritorio', 'claro'), ('escritorio', 'oscuro'))
ANCHO = {'android': 412, 'escritorio': 1440}
# Fila Global: U5-01 solo pide capturas de las P-xx; G-01..G-07 no tienen las 4 propias. PROVISIONAL «no» con bloqueo citado.
GLOBAL_CAPTURAS = ('G-01, G-02 y G-05 salen en las capturas de todas las pantallas (barra lateral y superior en escritorio, barra inferior en '
                   'Android; p. ej. las 4 de P-01) y G-07 «Sin conexión» en `frontend/pruebas-u5/flujos/capturas/U5-13-sin-conexion-*.jpg` '
                   '(4 combinaciones); G-03, G-04 y G-06 no tienen captura propia · bloqueo `docs/estado/bloqueos/C3-capturas-global.md`')
EXCEPCIONES_NO_LLAMADA = {
    'listPeriods': 'el mes sale del selector de la pantalla; sus necesidades las cubren getBudgetMonth y getBudgetFortnights (P-06) y getIncomeMonth (P-08), que traen las quincenas y su estado; no hay lista de periodos que mostrar',
}


# Operaciones con pantalla en U1 que el frontend todavía NO llama y para las que nadie autorizó una excepción: se muestran en la
# tabla con su motivo y el script sigue saliendo con 1 (no se admiten como excepción técnica; WF-3c c3-tablero, bloqueo C3-ajustes-tablero).
PENDIENTES_PANTALLA = {
    'updateDashboardSettings': 'PANTALLA PENDIENTE, no es excepción: U1 §3.3 (fila «Parámetros del tablero», N-275 y N-278, que §3.29 cuelga de P-26.02) pide leer y cambiar las cinco claves de K54 en Ajustes, pero ninguna pantalla de `frontend/src` llama PATCH `/api/v1/settings/dashboard` (solo Inicio lee `getDashboardSettings`) · bloqueo `docs/estado/bloqueos/C3-ajustes-tablero.md`',
}


def leer(p):
    return Path(p).read_text(encoding='utf-8')


# ---------- U1 ----------
def u1():
    texto = leer(RAIZ / 'docs/design/U1-arquitectura.md').split('\n')
    nombres, filas, seccion = {}, [], None
    for r in texto:
        if r.startswith('### '):
            seccion = r.split()[1]
        m = re.match(r'^\| (P-\d\d(?:\.\d\d)?) \| ([^|]+?) \|', r)
        if m and seccion in ('2.1', '2.3') and m.group(1) not in nombres:
            nombres[m.group(1)] = m.group(2).strip()
        if seccion == '3.29' and r.startswith('| ') and not r.startswith('| Pantalla') and not r.startswith('|---'):
            c = [x.strip() for x in r.strip('|').split('|')]
            filas.append((c[0].split(' ')[0], c[1] + ' ; ' + c[2]))
    necesidades = {}
    for ident, celdas in filas:
        ns, via = set(), {}
        for parte in re.split(r'[,;]', celdas):
            m = re.search(r'N-(\d{3})\s+a\s+N-(\d{3})', parte)
            nums = range(int(m.group(1)), int(m.group(2)) + 1) if m else [int(x) for x in re.findall(r'N-(\d{3})', parte)]
            v = re.search(r'vía (P-\d\d)', parte)
            for n in nums:
                ns.add(f'N-{n:03d}')
                if v:
                    via[f'N-{n:03d}'] = v.group(1)
        necesidades[ident] = (ns, via)
    return nombres, necesidades


# ---------- C1 y servidor ----------
def contrato():
    import yaml
    doc = yaml.safe_load(leer(RAIZ / 'docs/contratos/openapi.yaml'))
    ops = {}
    for ruta, it in doc['paths'].items():
        for m, op in it.items():
            if m in METODOS and isinstance(op, dict) and 'operationId' in op:
                ops[op['operationId']] = {'clave': f'{m.upper()} {ruta}', 'n': set(op.get('x-necesidades') or []),
                                         'externo': bool(op.get('x-agente-externo')),
                                         'tecnica': (op.get('description') or '').split('\n')[0].strip() if not op.get('x-necesidades') and (op.get('description') or '').startswith('Excepción técnica') else ''}
    return ops


def registradas():
    return {f'{a} {b}' for a, b in re.findall(r'Method:\s*"(\w+)",\s*Pattern:\s*"([^"]+)"', leer(RAIZ / 'backend/internal/api/routes.go'))}


# ---------- frontend ----------
RE_CLIENTE = re.compile(r"\b(?:pedir|pedirPagina|todas|encolar|usarConsulta)\(\s*'(get|post|put|patch|delete)'\s*,\s*'([^']+)'")
RE_GET = re.compile(r"\b(?:armarUrl|fetch)\(\s*'(/[^']+)'")
RE_IMPORT = re.compile(r"""(?:from\s+|import\(\s*|^import\s+)'(\.{1,2}/[^']+)'""", re.M)


def archivos_front():
    return [p for p in (FRONT / 'src').rglob('*') if p.suffix in ('.ts', '.tsx') and 'mocks' not in p.parts and not p.name.endswith('.gen.ts')]


def llamadas(p):
    t = leer(p)
    return {f'{m.upper()} {r}' for m, r in RE_CLIENTE.findall(t)} | {f'GET {r}' for r in RE_GET.findall(t)}


def resolver(desde, rel):
    base = (desde.parent / rel).resolve()
    for c in (base, base.with_suffix('.ts'), base.with_suffix('.tsx'), base / 'index.ts', base / 'index.tsx'):
        if c.is_file():
            return c
    return None


def cierre(raices):
    vistos, pila = set(), [r for r in raices if r]
    while pila:
        f = pila.pop()
        if f in vistos or 'mocks' in f.parts or f.name.endswith('.gen.ts') or f.name == 'registro.ts':
            continue
        vistos.add(f)
        for rel in RE_IMPORT.findall(leer(f)):
            r = resolver(f, rel)
            if r:
                pila.append(r)
    return vistos


def entradas():
    """P-xx → archivo de entrada: indice.ts / secciones-ajustes.ts de cada área + los de la base (U1 §2.1)."""
    e = {'Global': FRONT / 'src/main.tsx', 'P-00': FRONT / 'src/app/IniciarSesion.tsx',
         'P-03': FRONT / 'src/shell/HojaRegistro.tsx', 'P-22': FRONT / 'src/shell/HojaRegistro.tsx'}
    for ind in (FRONT / 'src/pantallas').glob('*/indice.ts'):
        for ident, rel in re.findall(r"'(P-\d\d)':\s*lazy\(\(\)\s*=>\s*import\('([^']+)'\)", leer(ind)):
            e[ident] = resolver(ind, rel)
    for sec in (FRONT / 'src/pantallas').glob('*/secciones-ajustes.ts'):
        for num, rel in re.findall(r"'(\d\d)':\s*lazy\(\(\)\s*=>\s*import\('([^']+)'\)", leer(sec)):
            e[f'P-26.{num}'] = resolver(sec, rel)
    return e


# ---------- pruebas ----------
def pruebas_pw(ruta):
    if not ruta:
        return {}, (0, 0, 0)
    d = json.loads(leer(ruta))
    por, tot = {}, [0, 0, 0]

    def rec(s, titulos):
        for sp in s.get('specs', []):
            estados = [t.get('status') for t in sp.get('tests', [])]
            ok = sp.get('ok') and 'expected' in estados
            tot[0] += len(estados); tot[1] += estados.count('expected'); tot[2] += estados.count('skipped')
            texto = ' '.join(titulos + [sp['title']])
            if ok:
                ids = set(re.findall(r'P-\d\d(?:\.\d\d)?|G-\d\d', texto))
                # Reglas declaradas (no hay P-xx en el título): F-1 = hoja P-03 (U1 F-1); F-6 o «Registrar peso» = hoja P-22
                # (U1 F-6); shell.spec.ts y cola.spec.ts = componentes globales G-01…G-07 (barra, menú «/», franja y cola).
                if re.search(r'\bF-1\b|registrar gasto', texto, re.I): ids.add('P-03')
                if re.search(r'\bF-6\b|registrar peso', texto, re.I): ids.add('P-22')
                if sp['file'] in ('shell.spec.ts', 'cola.spec.ts'): ids.add('G-01')
                ids |= {'P-26' for i in ids if i.startswith('P-26.')} | ({'Global'} if any(i.startswith('G-') for i in ids) else set())
                for i in ids:
                    por.setdefault(i, set()).add(f"{sp['file']}:{sp['line']}")
        for h in s.get('suites', []):
            rec(h, titulos + ([h['title']] if h.get('title') and not h['title'].endswith('.ts') else []))
    for s in d['suites']:
        rec(s, [])
    return por, tuple(tot)


def pruebas_go(ruta):
    if not ruta:
        return set()
    return set(re.findall(r'--- PASS: TestContrato_Operaciones/(\w+)', leer(ruta)))


# ---------- capturas (C3-05) ----------
def capturas(i, carpeta):
    if i == 'Global':
        return 'no', GLOBAL_CAPTURAS
    rutas, malas = [], []
    for disp, tema in COMBOS:
        f = carpeta / f'{i}-{disp}-{tema}.png'
        rutas.append('`' + (f.relative_to(RAIZ).as_posix() if f.is_relative_to(RAIZ) else str(f)) + '`')
        if not f.is_file() or f.stat().st_size == 0:
            malas.append(f'{f.name}: no existe o está vacía')
            continue
        b = f.read_bytes()[:24]
        ancho = struct.unpack('>I', b[16:20])[0] if b[:8] == b'\x89PNG\r\n\x1a\n' else None
        if ancho != ANCHO[disp]:
            malas.append(f'{f.name}: ancho {ancho} px (U5-01 pide {ANCHO[disp]})')
    if malas:
        return 'no', f'{4 - len(malas)}/4: ' + '; '.join(malas) + ' · ' + ', '.join(rutas)
    extra = ''
    if i == 'P-01':
        sub = RAIZ / 'docs/design/capturas-u5/evidencia/inicio-subsecciones'
        n = sum(1 for s in ('resumen', 'detalle', 'tendencias', 'calendario') for d, tm in COMBOS if (sub / f'P-01-{s}-{d}-{tm}.png').is_file())
        extra = f' · además las 4 subsecciones (Resumen, Detalle, Tendencias, Calendario) × 4 combinaciones: {n}/16 en `docs/design/capturas-u5/evidencia/inicio-subsecciones/P-01-<subsección>-<dispositivo>-<tema>.png` (U5-06; subcarpeta que no cuenta en U5-01)'
    return 'sí', '4/4: ' + ', '.join(rutas) + extra


# ---------- matriz ----------
ORDEN = lambda i: (0, '') if i == 'Global' else (1, i)


def construir(a):
    nombres, nec = u1()
    ops, regs, ent = contrato(), registradas(), entradas()
    llam_archivo = {p: llamadas(p) for p in archivos_front()}
    por_pw, tot_pw = pruebas_pw(a.pw_json)
    go_ok = pruebas_go(a.go_log)
    carpeta = (RAIZ / a.capturas) if not Path(a.capturas).is_absolute() else Path(a.capturas)
    ids = sorted(set(nombres) | set(nec) | {'Global'}, key=ORDEN)
    cierres = {i: cierre([ent[i]]) for i in ent if ent[i]}
    filas, informe = [], {'sin_operacion': [], 'no_impl': [], 'no_prob': []}
    usadas_por_fila = {}
    for i in ids:
        ns, via = nec.get(i, (set(), {}))
        local = i == 'P-26' or (ns and ns <= LOCALES)
        fila_ops = sorted(o for o, v in ops.items() if v['n'] & ns)
        usadas_por_fila[i] = fila_ops
        nombre = 'Componentes globales (G-01 a G-07)' if i == 'Global' else nombres.get(i, '¿?')
        if not fila_ops and not local:
            informe['sin_operacion'].append(i)
        # implementado
        if i not in ent or not ent[i]:
            impl = 'no · sin pantalla registrada en frontend/src'
            informe['no_impl'].append((i, 'sin pantalla'))
        else:
            rel = ent[i].relative_to(RAIZ).as_posix()
            if local and not fila_ops:
                impl = f'sí · `{rel}` (Local, sin API)'
            else:
                # Por necesidad (C1-02): cada N-xxx de la fila con operación en C1 debe tener AL MENOS una de sus
                # operaciones registrada en routes.go y llamada por la pantalla (su cierre de importaciones, o el de
                # la pantalla «vía» que declara U1 §3.29). Las N-xxx sin operación en C1 se listan aparte.
                nec_api = sorted(n for n in ns - LOCALES if any(n in ops[o]['n'] for o in fila_ops))
                nec_sin_op = sorted(n for n in ns - LOCALES if not any(n in v['n'] for v in ops.values()))
                faltan = []
                for n in nec_api:
                    alcance = set(cierres[i]) | (cierres.get(via[n], set()) if n in via else set())
                    cands = [o for o in fila_ops if n in ops[o]['n']]
                    if not any(ops[o]['clave'] in regs and any(ops[o]['clave'] in llam_archivo.get(f, ()) for f in alcance) for o in cands):
                        faltan.append(f'{n} ({"/".join(cands)})')
                k = len(nec_api)
                extra = f'; N sin operación en C1: {", ".join(nec_sin_op)}' if nec_sin_op else ''
                if faltan or nec_sin_op:
                    impl = f'no · `{rel}`: {k - len(faltan)}/{k} necesidades; sin llamada: {", ".join(faltan) or "—"}{extra}'
                    informe['no_impl'].append((i, faltan + nec_sin_op))
                else:
                    impl = f'sí · `{rel}`: {k}/{k} necesidades con llamada y ruta registrada'
        # probado
        pw = sorted(por_pw.get(i, set()))
        go_falta = [o for o in fila_ops if o not in go_ok]
        go_txt = 'Go: no aplica (Local)' if not fila_ops else (f'Go: {len(fila_ops) - len(go_falta)}/{len(fila_ops)} `TestContrato_Operaciones/<operationId>`')
        pw_txt = f'PW: {len(pw)} ({", ".join("`" + x + "`" for x in pw[:2])}{", …" if len(pw) > 2 else ""})' if pw else 'PW: 0'
        if pw and not go_falta:
            prob = f'sí · {go_txt} · {pw_txt}'
        else:
            prob = f'no · {go_txt}{" (faltan: " + ", ".join(go_falta) + ")" if go_falta else ""} · {pw_txt}'
            informe['no_prob'].append(i)
        opcol = ('Local (contenedor de P-26.01…P-26.14; PROVISIONAL, decisión del orquestador, bloqueo C3-matriz punto 1)' if i == 'P-26' else
                 'Local (N-301, sin API)' if local else ', '.join(fila_ops))
        # Cada estado en dos columnas: «sí»/«no» sola (para `grep -c "| no |"` de C3-06) y su evidencia.
        est = lambda x: ' | '.join(x.split(' · ', 1))
        cap, cap_ev = capturas(i, carpeta)
        filas.append(f'| {i} | {nombre} | {opcol} | {est(impl)} | {est(prob)} | {cap} | {cap_ev} |')
    # operaciones sin pantalla
    en_filas = {o for v in usadas_por_fila.values() for o in v}
    sin_pantalla = sorted(o for o in ops if o not in en_filas)
    exc = [(o, 'x-agente-externo (K53): la usa el agente externo por la API; la app no la llama') for o in sin_pantalla if ops[o]['externo']]
    exc += [(o, ops[o]['tecnica']) for o in sin_pantalla if not ops[o]['externo'] and ops[o]['tecnica']]
    sin_justificar = [o for o in sin_pantalla if not ops[o]['externo'] and not ops[o]['tecnica']]
    # C3-03
    invocadas = set().union(*llam_archivo.values())
    claves = {v['clave'] for v in ops.values()}
    fuera = sorted(invocadas - claves)
    no_llamadas = sorted((o, [i for i in ids if o in usadas_por_fila[i]]) for o in en_filas if ops[o]['clave'] not in invocadas)
    # Cobertura de cada operación no llamada: por necesidad, las otras operaciones del contrato que la cubren y que sí se llaman.
    cobertura, no_llamadas_malas = {}, []
    for o, _ in no_llamadas:
        cub = {n: sorted(x for x, v in ops.items() if x != o and n in v['n'] and v['clave'] in invocadas and v['clave'] in regs) for n in sorted(ops[o]['n'])}
        cobertura[o] = cub
        if o not in EXCEPCIONES_NO_LLAMADA or not all(cub.values()):
            no_llamadas_malas.append(o)
    sin_reg = sorted(invocadas - regs)
    # C3-02 por necesidad: N-xxx de U1 §3.29 sin operación en C1 ni Local, y N-xxx de C1 sin fila en U1.
    n_u1 = set().union(*(v[0] for v in nec.values())) if nec else set()
    n_c1 = set().union(*(v['n'] for v in ops.values()))
    n_sin_op = sorted(n_u1 - n_c1 - LOCALES)
    n_sin_fila = sorted(n_c1 - n_u1)
    return dict(n_u1=n_u1, n_c1=n_c1, n_sin_op=n_sin_op, n_sin_fila=n_sin_fila, filas=filas, exc=exc, sin_justificar=sin_justificar, informe=informe, ids=ids, ops=ops, regs=regs,
                invocadas=invocadas, fuera=fuera, no_llamadas=no_llamadas, cobertura=cobertura, no_llamadas_malas=no_llamadas_malas, sin_reg=sin_reg, tot_pw=tot_pw, go_ok=go_ok, nombres=nombres)


def tabla(r):
    enc = ['| Id | Pantalla | Operaciones de C1 (`operationId`) | Implementado | Evidencia (implementado) | Probado | Evidencia (probado) | Capturado | Evidencia (capturado) |',
           '|---|---|---|---|---|---|---|---|---|']
    exc = ['', '**Operaciones sin pantalla (excepciones técnicas, una por una):**', '',
           '| Operación | Motivo |', '|---|---|'] + [f'| {o} | {m} |' for o, m in r['exc']]
    exc += ['', '**Operaciones con pantalla en U1 que el frontend no llama en ninguna parte** (excepciones técnicas PROVISIONALES, una por una; decisión del orquestador 2026-10-12, bloqueo C3-matriz punto 3c; solo se admiten si otra operación llamada cubre cada una de sus necesidades):', '',
            '| Operación | Filas donde aparece | Necesidades cubiertas por otra operación llamada | Motivo |', '|---|---|---|---|']
    for o, f in r['no_llamadas']:
        cub = '; '.join(f'{n}: {", ".join(v) or "SIN COBERTURA"}' for n, v in r['cobertura'][o].items())
        exc.append(f"| {o} | {', '.join(f)} | {cub} | {EXCEPCIONES_NO_LLAMADA.get(o) or PENDIENTES_PANTALLA.get(o, 'SIN JUSTIFICAR (ver bloqueo C3-matriz)')} |")
    if not r['no_llamadas']:
        exc.append('| (ninguna) | — | — | — |')
    if r['sin_justificar']:
        exc += [f'| {o} | SIN JUSTIFICAR (ver bloqueo C3-matriz) |' for o in r['sin_justificar']]
    return '\n'.join([INI] + enc + r['filas'] + exc + [FIN])


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('--pw-json'); ap.add_argument('--go-log'); ap.add_argument('--escribir', action='store_true')
    ap.add_argument('--capturas', default='docs/design/capturas-u5', help='carpeta de las capturas de U5 (C3-05)')
    a = ap.parse_args()
    r = construir(a)
    print(f"C3-03 frontend: {len(r['invocadas'])} rutas invocadas; fuera de openapi.yaml: {len(r['fuera'])}; sin registrar en routes.go: {len(r['sin_reg'])} (routes.go: {len(r['regs'])}; openapi: {len(r['ops'])} operaciones)")
    for k in r['fuera']: print('  FUERA', k)
    for k in r['sin_reg']: print('  SIN-REGISTRAR', k)
    pant = [i for i in r['ids'] if i.startswith('P-')]
    print(f"C3-02 pantallas U1: {len(r['nombres'])} (P-xx {sum(1 for i in r['nombres'] if '.' not in i)} + P-26.xx {sum(1 for i in r['nombres'] if '.' in i)}); filas: {len(r['filas'])} (P-: {len(pant)} + Global)")
    faltan_p = sorted(set(r['nombres']) - set(pant))
    print(f"C3-02 pantallas sin fila: {len(faltan_p)} {faltan_p}; filas sin operación ni Local: {len(r['informe']['sin_operacion'])} {r['informe']['sin_operacion']}")
    print(f"C3-02 operaciones: {len(r['ops'])}; sin pantalla: {len(r['exc']) + len(r['sin_justificar'])} (excepciones: {len(r['exc'])}, de ellas K53: {sum(1 for _, m in r['exc'] if 'K53' in m)}; sin justificar: {len(r['sin_justificar'])} {r['sin_justificar']})")
    print(f"C3-02 necesidades: U1 §3.29 {len(r['n_u1'])} N-xxx ({len(r['n_u1'] & LOCALES)} Local); C1 x-necesidades {len(r['n_c1'])}; de U1 sin operación ni Local: {len(r['n_sin_op'])} {r['n_sin_op']}; de C1 sin fila en U1: {len(r['n_sin_fila'])} {r['n_sin_fila']}")
    print(f"C3-04 Go: {len(r['go_ok'])} subpruebas TestContrato_Operaciones PASS · Playwright: {r['tot_pw'][0]} ejecutadas, {r['tot_pw'][1]} pasan, {r['tot_pw'][2]} omitidas")
    impl = sum(1 for f in r['filas'] if f.split('|')[4].strip() == 'sí')
    prob = sum(1 for f in r['filas'] if f.split('|')[6].strip() == 'sí')
    cap = sum(1 for f in r['filas'] if f.split('|')[8].strip() == 'sí')
    print(f"Estado: implementado sí {impl}/{len(r['filas'])}; probado sí {prob}/{len(r['filas'])}; capturado sí {cap}/{len(r['filas'])}")
    print(f"C3-05 capturas: {4 * sum(1 for i in r['ids'] if i != 'Global')} esperadas (4 × {sum(1 for i in r['ids'] if i != 'Global')} P-); filas con 4/4: {cap}; filas en «no»: {[f.split('|')[1].strip() for f in r['filas'] if f.split('|')[8].strip() != 'sí']}")
    for i, f in r['informe']['no_impl']: print('  NO-IMPLEMENTADO', i, f if isinstance(f, str) else ', '.join(f))
    for i in r['informe']['no_prob']: print('  NO-PROBADO', i)
    print(f"Operaciones con pantalla que el frontend no llama: {len(r['no_llamadas'])} {[o for o, _ in r['no_llamadas']]}; sin excepción justificada o sin cobertura: {len(r['no_llamadas_malas'])} {r['no_llamadas_malas']}")
    nueva = tabla(r)
    texto = leer(DOC) if DOC.exists() else ''
    actual = texto[texto.find(INI):texto.find(FIN) + len(FIN)] if INI in texto and FIN in texto else ''
    if a.escribir:
        if not actual:
            sys.exit(f'{DOC} no tiene los marcadores de la matriz')
        DOC.write_text(texto.replace(actual, nueva), encoding='utf-8')
        print('Matriz escrita en', DOC.relative_to(RAIZ))
    else:
        print('Matriz del documento = matriz generada:', 'sí' if actual == nueva else 'NO')
    malo = r['n_sin_op'] or r['n_sin_fila'] or r['fuera'] or r['sin_reg'] or faltan_p or r['informe']['sin_operacion'] or r['sin_justificar'] or r['no_llamadas_malas'] or (not a.escribir and actual != nueva)
    sys.exit(1 if malo else 0)


if __name__ == '__main__':
    main()
