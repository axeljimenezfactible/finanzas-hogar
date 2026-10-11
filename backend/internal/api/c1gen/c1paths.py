# Rutas del generador C1 (sin rutas absolutas fijas). Lo importan las etapas parse_schema.py y gen1…gen6.
#   C1GEN_ROOT  raíz del repositorio de donde se leen las entradas (por defecto, cuatro niveles arriba de esta carpeta)
#   C1GEN_OUT   raíz donde se escriben las salidas (por defecto, la misma raíz: el repositorio)
#   C1GEN_WORK  carpeta de trabajo para los intermedios schema.json y openapi.json (la crea y la borra build.py)
import os
HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.abspath(os.environ.get('C1GEN_ROOT') or os.path.join(HERE, os.pardir, os.pardir, os.pardir, os.pardir))
OUT = os.path.abspath(os.environ.get('C1GEN_OUT') or ROOT)
WORK = os.environ.get('C1GEN_WORK', '')
if not WORK or not os.path.isdir(WORK):
    raise SystemExit('c1gen: falta la carpeta de trabajo C1GEN_WORK; ejecuta el generador con build.py')
