"""Oráculo S5 · cargador PROFECO (§11) y parseo de TSV de Tesseract (§12) sobre archivos de oraculo/fixtures/s5/_datos/.

Sin red. Importes y coordenadas como texto/Decimal, nunca flotantes. Los motivos (`reason`) son texto libre y no se comparan.
"""
import csv
import datetime as dt
import io
import re
from decimal import Decimal, InvalidOperation
from fractions import Fraction

from pre_comun import ErrorMotor, half_up
from pre_datos import DATOS_S5, ruta_repo

COLUMNAS = ('producto', 'presentacion', 'marca', 'categoria', 'catalogo', 'precio', 'fecha_registro', 'cadena_comercial', 'giro',
            'nombre_comercial', 'direccion', 'estado', 'municipio', 'latitud', 'longitud')
_FECHAS = (re.compile(r'^(\d{4})-(\d{2})-(\d{2})(?:[ T]\d{2}:\d{2}(?::\d{2}(?:\.\d+)?)?)?$'),
           re.compile(r'^(\d{4})/(\d{2})/(\d{2})(?: \d{2}:\d{2}(?::\d{2})?)?$'),
           re.compile(r'^(\d{2})/(\d{2})/(\d{4})(?: \d{2}:\d{2}(?::\d{2})?)?$'))
_PRECIO = re.compile(r'^\d+(?:\.\d{1,2})?$')


def _lee(ruta):
    b = ruta.read_bytes()
    if b.startswith(b'\xef\xbb\xbf'):
        b = b[3:]
    try:
        return b.decode('utf-8')
    except UnicodeDecodeError:
        return b.decode('cp1252', errors='replace')


def _registros(ruta):
    rd = csv.reader(io.StringIO(_lee(ruta), newline=''))
    filas = [r for r in rd if r]
    return filas


def _fecha(txt):
    t = txt.strip()
    for i, rx in enumerate(_FECHAS):
        m = rx.match(t)
        if m:
            a, b, c = (int(x) for x in m.groups())
            try:
                return dt.date(c, b, a) if i == 2 else dt.date(a, b, c)
            except ValueError:
                return None
    return None


_SIN_ACENTO = str.maketrans('áàäâéèëêíìïîóòöôúùüûñ', 'aaaaeeeeiiiioooouuuun')
_SE_QUITAN = set("-'’.`")


def _normaliza(txt):
    """Bloqueo S5-profeco-filtro, opción A (PROVISIONAL): minúsculas, sin acentos, sin guion, apóstrofo, punto ni acento
    grave; cualquier otro signo o espacio separa palabras, con un solo espacio entre ellas. «Wal-mart» = «WALMART»."""
    if not txt:
        return ''
    palabras, actual = [], []
    for ch in txt.lower().translate(_SIN_ACENTO):
        if ch in _SE_QUITAN:
            continue
        if ch.isalnum():
            actual.append(ch)
        elif actual:
            palabras.append(''.join(actual))
            actual = []
    if actual:
        palabras.append(''.join(actual))
    return ' '.join(palabras)


def profeco_import(d, files, links_file, mode):
    if mode not in ('production', 'review'):
        raise ErrorMotor('invalid_argument', 'mode')
    if not isinstance(files, list) or not files or not all(isinstance(f, str) for f in files):
        raise ErrorMotor('invalid_argument', 'files')
    enlaces = {}
    filas = _registros(ruta_repo(links_file, 'links_file', DATOS_S5))
    cab = [c.strip().lower() for c in filas[0]] if filas else []
    for c in ('producto', 'presentacion', 'marca', 'presentacion_id', 'estado', 'confirmo', 'fecha_confirmacion'):
        if c not in cab:
            raise ErrorMotor('profeco_header', links_file)
    for f in filas[1:]:
        x = dict(zip(cab, (v.strip() for v in f)))
        if x['estado'] not in ('propuesto', 'confirmado'):
            raise ErrorMotor('invalid_argument', 'estado')
        if (x['estado'] == 'confirmado') != bool(x['confirmo'] and x['fecha_confirmacion']) or (x['estado'] == 'propuesto' and (x['confirmo'] or x['fecha_confirmacion'])):
            raise ErrorMotor('invalid_argument', 'confirmacion')
        enlaces.setdefault((x['producto'], x['presentacion'], x['marca']), x)
    # Filtro (PROVISIONAL, bloqueo S5-profeco-filtro, opción A): cadena = profeco_name de las tiendas chain;
    # municipio = municipality solo de las sucursales (branch); las dos comparaciones con _normaliza.
    cadenas = {_normaliza(s['profeco_name']): s for s in d.stores.values() if s['kind'] == 'chain' and _normaliza(s['profeco_name'])}
    munis = {_normaliza(s['municipality']) for s in d.stores.values() if s['kind'] == 'branch' and _normaliza(s['municipality'])}
    res = {'read': 0, 'rejected': 0, 'filtered_out': 0, 'skipped': 0, 'loaded': 0, 'duplicates': 0}
    cobertura, vistos, rechazos = {}, set(), []
    for fn in files:
        ruta = ruta_repo(fn, 'files', DATOS_S5)
        filas = _registros(ruta)
        cab = [c.strip().lower() for c in filas[0]] if filas else []
        if any(c not in cab for c in COLUMNAS):
            raise ErrorMotor('profeco_header', fn)
        for n, f in enumerate(filas[1:], 1):
            res['read'] += 1
            if len(f) != len(cab):
                res['rejected'] += 1
                rechazos.append({'file': fn, 'row': n, 'reason': 'columnas'})
                continue
            x = dict(zip(cab, (v.strip() for v in f)))
            if not _PRECIO.match(x['precio']):
                res['rejected'] += 1
                rechazos.append({'file': fn, 'row': n, 'reason': 'precio ilegible'})
                continue
            fch = _fecha(x['fecha_registro'])
            if fch is None:
                res['rejected'] += 1
                rechazos.append({'file': fn, 'row': n, 'reason': 'fecha ilegible'})
                continue
            cadena = cadenas.get(_normaliza(x['cadena_comercial']))
            if cadena is None or _normaliza(x['municipio']) not in munis:
                res['filtered_out'] += 1
                continue
            lg = enlaces.get((x['producto'], x['presentacion'], x['marca']))
            if lg is None or lg['presentacion_id'] not in d.pres or (mode == 'production' and lg['estado'] != 'confirmado'):
                res['skipped'] += 1
                continue
            # tienda del precio: la sucursal de esa misma cadena cuyo profeco_name = nombre_comercial; si no hay, la cadena
            nc = _normaliza(x['nombre_comercial'])
            tienda = next((s for s in d.stores.values() if s['kind'] == 'branch' and s['parent'] == cadena['code']
                           and nc and _normaliza(s['profeco_name']) == nc), None) or cadena
            clave = (lg['presentacion_id'], tienda['code'], fch)
            if clave in vistos:
                res['duplicates'] += 1
                continue
            vistos.add(clave)
            res['loaded'] += 1
            k = (x['cadena_comercial'], x['municipio'])
            cobertura[k] = cobertura.get(k, 0) + 1
    return {**res, 'coverage': [{'chain': c, 'municipality': m, 'price_count': v} for (c, m), v in sorted(cobertura.items())],
            'rejections': rechazos}


# ---------- OCR ----------
_IMPORTE = re.compile(r'^\$?(\d{1,3}(?:,\d{3})+|\d+)\.(\d{2})$')
_CANT = re.compile(r'^(\d+(?:[.,]\d+)?)$')
_CANT_UNI = re.compile(r'^(\d+(?:[.,]\d+)?)(kg|g|gr|l|lt|ml|pz|pza|pzas|pzs)$', re.I)
_UNIDADES = {'kg', 'g', 'gr', 'l', 'lt', 'ml', 'pz', 'pza', 'pzas', 'pzs'}


def _milli(txt):
    return half_up(Fraction(Decimal(txt.replace(',', '.')) * 1000))


def ocr_parse_tsv(tsv_file):
    ruta = ruta_repo(tsv_file, 'tsv_file', DATOS_S5)
    rd = csv.reader(io.StringIO(_lee(ruta), newline=''), delimiter='\t', quoting=csv.QUOTE_NONE)
    filas = list(rd)
    if not filas:
        return {'lines': []}
    cab = filas[0]
    for c in ('level', 'page_num', 'block_num', 'par_num', 'line_num', 'conf', 'text'):
        if c not in cab:
            raise ErrorMotor('invalid_argument', 'tsv_file')
    grupos, orden = {}, []
    for f in filas[1:]:
        x = dict(zip(cab, f))
        txt = x.get('text', '').strip()
        if x['level'].strip() != '5' or not txt:
            continue
        k = tuple(x[c] for c in ('page_num', 'block_num', 'par_num', 'line_num'))
        if k not in grupos:
            grupos[k] = []
            orden.append(k)
        try:
            c = Decimal(x['conf'])
        except InvalidOperation:
            c = Decimal(-1)
        grupos[k].append((txt, c))
    out = []
    for k in orden:
        pal = grupos[k]
        toks = [t for t, _ in pal]
        confs = [Fraction(c) for _, c in pal if c >= 0]
        conf = half_up(sum(confs) / len(confs)) if confs else 0
        imp, qty = None, None
        idx = [i for i, t in enumerate(toks) if _IMPORTE.match(t)]
        if idx:
            m = _IMPORTE.match(toks[idx[-1]])
            imp = int(m.group(1).replace(',', '')) * 100 + int(m.group(2))
            del toks[idx[-1]]
        if toks:
            m = _CANT.match(toks[0])
            mu = _CANT_UNI.match(toks[0])
            if m:
                qty = _milli(m.group(1))
                del toks[0]
                if toks and toks[0].lower() in _UNIDADES:
                    del toks[0]
            elif mu:
                qty = _milli(mu.group(1))
                del toks[0]
        out.append({'description': ' '.join(toks), 'quantity_milli': qty, 'amount_cents': imp, 'confidence': conf})
    return {'lines': out}
