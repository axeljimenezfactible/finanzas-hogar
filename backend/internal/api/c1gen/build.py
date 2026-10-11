#!/usr/bin/env python3
"""Generador del contrato C1 · Finanzas Hogar.

Regenera, desde el esquema de S1 y las necesidades de U1, los tres archivos de C1 que no se editan a mano:
  docs/contratos/openapi.yaml       OpenAPI 3.1.1 (contrato C1)
  backend/internal/api/routes.go    rutas del enrutador, idénticas a las del contrato (S6-03)
  docs/contratos/C1-datos.md        documento de C1, armado desde c1_plantilla.md (incluye la matriz §10)

Uso (solo biblioteca estándar de Python 3.11 o posterior; da igual la carpeta desde la que se llame):
  python3 -I backend/internal/api/c1gen/build.py            escribe los tres archivos en el repositorio
  python3 -I backend/internal/api/c1gen/build.py --check    regenera en una carpeta temporal y compara byte a byte; no escribe
  python3 -I backend/internal/api/c1gen/build.py --out DIR  escribe los tres archivos bajo DIR (árbol espejo)
Código de salida: 0 si todo salió bien (con --check, 0 = los tres archivos son idénticos); 1 si hay diferencias o falla una etapa.
Los intermedios (schema.json y openapi.json) van a una carpeta temporal bajo datos/tmp/ que se borra al terminar.
"""
import argparse, contextlib, difflib, hashlib, io, os, runpy, shutil, sys, tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.abspath(os.path.join(HERE, os.pardir, os.pardir, os.pardir, os.pardir))
INPUTS = ['backend/db/migrations/00001_esquema.sql', 'docs/design/U1-arquitectura.md']
SOURCES = ['c1paths.py', 'parse_schema.py', 'gen1_base.py', 'gen2_schemas.py', 'gen3_ops.py', 'gen4_build.py', 'gen5_go.py', 'gen6_c1.py',
           'c1_plantilla.md', 'c1_inconsistencias.md']
STAGES = ['parse_schema.py', 'gen4_build.py', 'gen5_go.py', 'gen6_c1.py']
OUTPUTS = ['docs/contratos/openapi.yaml', 'backend/internal/api/routes.go', 'docs/contratos/C1-datos.md']


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
    ap = argparse.ArgumentParser(description='Generador del contrato C1 (openapi.yaml, routes.go y C1-datos.md).')
    grp = ap.add_mutually_exclusive_group()
    grp.add_argument('--check', action='store_true', help='regenera en una carpeta temporal y compara con el repositorio; no escribe')
    grp.add_argument('--out', metavar='DIR', help='escribe las salidas bajo DIR en lugar del repositorio')
    args = ap.parse_args()
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
        out_root = os.path.join(work, 'salida') if args.check else os.path.abspath(args.out) if args.out else ROOT
        for rel in OUTPUTS:
            os.makedirs(os.path.dirname(os.path.join(out_root, rel)), exist_ok=True)
        run_stages(out_root, work)
        if not args.check:
            for rel in OUTPUTS:
                print(f'  escrito {os.path.join(out_root, rel)} (sha256 {sha12(os.path.join(out_root, rel))}…)')
            return 0
        diffs = 0
        for rel in OUTPUTS:
            new = open(os.path.join(out_root, rel), 'rb').read()
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
        print('RESULTADO:', 'los 3 archivos coinciden byte a byte con el generador' if not diffs else f'{diffs} archivo(s) difieren del generador')
        return 1 if diffs else 0
    finally:
        shutil.rmtree(work, ignore_errors=True)
        if made_parent:
            with contextlib.suppress(OSError):
                os.rmdir(tmp_parent)


if __name__ == '__main__':
    sys.exit(main())
