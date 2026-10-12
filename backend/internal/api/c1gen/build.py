#!/usr/bin/env python3
"""Generador del contrato C1 · Finanzas Hogar.

Regenera, desde el esquema de S1 y las necesidades de U1, los tres archivos de C1 que no se editan a mano:
  docs/contratos/openapi.yaml       OpenAPI 3.1.1 (contrato C1)
  backend/internal/api/routes.go    rutas del enrutador, idénticas a las del contrato (S6-03)
  backend/internal/api/types_gen.go tipos Go de components/schemas y metadatos por operación (gen8_types.py, WF-3)
  docs/contratos/C1-datos.md        documento de C1, armado desde c1_plantilla.md (incluye la matriz §10)
  docs/contratos/extractos/         un extracto por etiqueta, componentes-comunes.yaml e INDICE.md (protocolo-agentes §4)
El esquema sale de TODAS las migraciones de backend/db/migrations/ aplicadas en orden (C1 v1.1).

Uso (solo biblioteca estándar de Python 3.11 o posterior; da igual la carpeta desde la que se llame):
  python3 -I backend/internal/api/c1gen/build.py            escribe los tres archivos en el repositorio
  python3 -I backend/internal/api/c1gen/build.py --check    regenera en una carpeta temporal y compara byte a byte (incluidos los extractos); no escribe
  python3 -I backend/internal/api/c1gen/build.py --extractos  regenera en una carpeta temporal y solo escribe docs/contratos/extractos/
  python3 -I backend/internal/api/c1gen/build.py --out DIR  escribe los tres archivos bajo DIR (árbol espejo)
Código de salida: 0 si todo salió bien (con --check, 0 = todos los archivos son idénticos); 1 si hay diferencias o falla una etapa.
Los intermedios (schema.json y openapi.json) van a una carpeta temporal bajo datos/tmp/ que se borra al terminar.
"""
import argparse, contextlib, difflib, glob, hashlib, io, os, runpy, shutil, sys, tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.abspath(os.path.join(HERE, os.pardir, os.pardir, os.pardir, os.pardir))
EXT_DIR = 'docs/contratos/extractos'


def inputs():
    migs = sorted(os.path.relpath(p, ROOT) for p in glob.glob(os.path.join(ROOT, 'backend/db/migrations/[0-9]*.sql')))
    return migs + ['docs/design/U1-arquitectura.md']


def ext_files(root):
    d = os.path.join(root, EXT_DIR)
    return sorted(os.path.join(EXT_DIR, f) for f in os.listdir(d) if f.endswith('.yaml') or f == 'INDICE.md') if os.path.isdir(d) else []


SOURCES = ['c1paths.py', 'parse_schema.py', 'gen1_base.py', 'gen2_schemas.py', 'gen3_ops.py', 'gen4_build.py', 'gen5_go.py', 'gen6_c1.py',
           'gen7_extractos.py', 'gen8_types.py', 'c1_plantilla.md', 'c1_inconsistencias.md']
STAGES = ['parse_schema.py', 'gen4_build.py', 'gen5_go.py', 'gen6_c1.py', 'gen7_extractos.py', 'gen8_types.py']
OUTPUTS = ['docs/contratos/openapi.yaml', 'backend/internal/api/routes.go', 'docs/contratos/C1-datos.md', 'backend/internal/api/types_gen.go']


def sha12(path):
    return hashlib.sha256(open(path, 'rb').read()).hexdigest()[:12]


def run_stages(out_root, work):
    os.environ['C1GEN_ROOT'] = ROOT
    os.environ['C1GEN_OUT'] = out_root
    os.environ['C1GEN_WORK'] = work
    sys.dont_write_bytecode = True  # no deja __pycache__ en el repositorio
    sys.path.insert(0, HERE)
    for st in STAGES:
        sys.argv = [os.path.join(HERE, st)]
        buf = io.StringIO()
        with contextlib.redirect_stdout(buf):
            runpy.run_path(os.path.join(HERE, st), run_name='__main__')
        first = (buf.getvalue().strip().splitlines() or [''])[0]
        print(f'  {st}: {first.replace(out_root + os.sep, "")}')


def main():
    ap = argparse.ArgumentParser(description='Generador del contrato C1 (openapi.yaml, routes.go, C1-datos.md y extractos).')
    grp = ap.add_mutually_exclusive_group()
    grp.add_argument('--check', action='store_true', help='regenera en una carpeta temporal y compara con el repositorio; no escribe')
    grp.add_argument('--extractos', action='store_true', help='regenera en una carpeta temporal y solo escribe docs/contratos/extractos/')
    grp.add_argument('--out', metavar='DIR', help='escribe las salidas bajo DIR en lugar del repositorio')
    args = ap.parse_args()
    INPUTS = inputs()
    missing = [rel for rel in INPUTS if not os.path.isfile(os.path.join(ROOT, rel))] + [s for s in SOURCES if not os.path.isfile(os.path.join(HERE, s))]
    if missing:
        print('c1gen: faltan archivos:', ', '.join(missing))
        return 1
    print(f'c1gen: raíz {ROOT}')
    for rel in INPUTS:
        print(f'  entrada {rel} (sha256 {sha12(os.path.join(ROOT, rel))}…)')
    tmp_parent = os.path.join(ROOT, 'datos', 'tmp')
    made_parent = not os.path.isdir(tmp_parent)
    os.makedirs(tmp_parent, exist_ok=True)
    work = tempfile.mkdtemp(prefix='c1gen-', dir=tmp_parent)
    try:
        out_root = os.path.join(work, 'salida') if args.check or args.extractos else os.path.abspath(args.out) if args.out else ROOT
        for rel in OUTPUTS:
            os.makedirs(os.path.dirname(os.path.join(out_root, rel)), exist_ok=True)
        run_stages(out_root, work)
        if args.extractos:
            new, old = ext_files(out_root), ext_files(ROOT)
            os.makedirs(os.path.join(ROOT, EXT_DIR), exist_ok=True)
            for rel in old:
                if rel not in new:
                    os.remove(os.path.join(ROOT, rel))
            for rel in new:
                shutil.copyfile(os.path.join(out_root, rel), os.path.join(ROOT, rel))
                print(f'  escrito {os.path.join(ROOT, rel)} (sha256 {sha12(os.path.join(ROOT, rel))}…)')
            return 0
        if not args.check:
            for rel in OUTPUTS + ext_files(out_root):
                print(f'  escrito {os.path.join(out_root, rel)} (sha256 {sha12(os.path.join(out_root, rel))}…)')
            return 0
        diffs = 0
        allrel = OUTPUTS + sorted(set(ext_files(out_root)) | set(ext_files(ROOT)))
        for rel in allrel:
            new_path = os.path.join(out_root, rel)
            new = open(new_path, 'rb').read() if os.path.isfile(new_path) else b''
            cur_path = os.path.join(ROOT, rel)
            cur = open(cur_path, 'rb').read() if os.path.isfile(cur_path) else b''
            if new == cur:
                print(f'IDÉNTICO {rel} ({len(new)} bytes, sha256 {hashlib.sha256(new).hexdigest()[:12]}…)')
                continue
            diffs += 1
            print(f'DIFIERE  {rel} (repositorio {len(cur)} bytes, generador {len(new)} bytes)')
            udiff = difflib.unified_diff(cur.decode('utf-8', 'replace').splitlines(), new.decode('utf-8', 'replace').splitlines(),
                                         'repositorio/' + rel, 'generador/' + rel, lineterm='', n=1)
            for i, line in enumerate(udiff):
                if i >= 40:
                    print('    …')
                    break
                print('    ' + line[:200])
        print('RESULTADO:', f'los {len(allrel)} archivos coinciden byte a byte con el generador' if not diffs else f'{diffs} archivo(s) difieren del generador')
        return 1 if diffs else 0
    finally:
        shutil.rmtree(work, ignore_errors=True)
        if made_parent:
            with contextlib.suppress(OSError):
                os.rmdir(tmp_parent)


if __name__ == '__main__':
    sys.exit(main())
