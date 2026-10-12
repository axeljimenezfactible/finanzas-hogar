"""Oráculo S4 · piezas comunes: errores, redondeo HALF_UP, fechas y validación estricta de entradas.

Implementa S4-interfaz.md §12.1 (campos desconocidos o ausentes) y los redondeos de §1 (C-02) sin leer el código Go.
Solo biblioteca estándar. Dinero en centavos, cantidades en milésimas, todo exacto con Fraction.
"""
import datetime as dt
import re
from fractions import Fraction


class ErrorMotor(Exception):
    """Error del paquete (S4-interfaz §14 y §12.2): `code` en snake_case sin «Err»."""

    def __init__(self, code, detalle=''):
        super().__init__(f'{code}: {detalle}')
        self.code = code
        self.detalle = detalle


# ---- redondeo (C-02): ROUND_HALF_UP, la mitad se aleja de cero; una sola vez al final ----
def half_up(x):
    x = Fraction(x)
    n, d = x.numerator, x.denominator
    q = (2 * abs(n) + d) // (2 * d)
    return q if n >= 0 else -q


def techo(x):
    x = Fraction(x)
    return -((-x.numerator) // x.denominator)


def piso(x):
    x = Fraction(x)
    return x.numerator // x.denominator


def frac_json(x):
    x = Fraction(x)
    return {'num': x.numerator, 'den': x.denominator}


# ---- fechas ----
_ISO = re.compile(r'^\d{4}-\d{2}-\d{2}$')
_SEMANA = re.compile(r'^(\d{4})-W(\d{2})$')


def fecha(v, ruta):
    if not isinstance(v, str) or not _ISO.match(v):
        raise ErrorMotor('invalid_date', ruta)
    try:
        return dt.date.fromisoformat(v)
    except ValueError:
        raise ErrorMotor('invalid_date', ruta)


def semana_a_lunes(v, ruta):
    """'AAAA-Www' (semana ISO, lunes a domingo, K8) -> fecha del lunes."""
    m = _SEMANA.match(v) if isinstance(v, str) else None
    if not m:
        raise ErrorMotor('invalid_argument', ruta)
    try:
        return dt.date.fromisocalendar(int(m.group(1)), int(m.group(2)), 1)
    except ValueError:
        raise ErrorMotor('invalid_argument', ruta)


def lunes_a_semana(lunes):
    y, w, _ = lunes.isocalendar()
    return f'{y}-W{w:02d}'


# ---- validación estricta de entradas ----
def objeto(v, ruta, requeridos=(), opcionales=()):
    """Objeto con claves exactas: una desconocida es unknown_field; una obligatoria ausente, missing_field.
    Las opcionales ausentes equivalen a null (supuesto: «puede ser null» admite ausencia)."""
    if not isinstance(v, dict):
        raise ErrorMotor('invalid_argument', ruta)
    for k in v:
        if k not in requeridos and k not in opcionales:
            raise ErrorMotor('unknown_field', f'{ruta}.{k}')
    for k in requeridos:
        if k not in v:
            raise ErrorMotor('missing_field', f'{ruta}.{k}')
    return v


def lista(v, ruta):
    if not isinstance(v, list):
        raise ErrorMotor('invalid_argument', ruta)
    return v


def entero(v, ruta, lo=None, hi=None, nulo=False):
    if v is None and nulo:
        return None
    if not isinstance(v, int) or isinstance(v, bool) or (lo is not None and v < lo) or (hi is not None and v > hi):
        raise ErrorMotor('invalid_argument', ruta)
    return v


def texto(v, ruta, nulo=False):
    if v is None and nulo:
        return None
    if not isinstance(v, str) or not v:
        raise ErrorMotor('invalid_argument', ruta)
    return v


def booleano(v, ruta):
    if not isinstance(v, bool):
        raise ErrorMotor('invalid_argument', ruta)
    return v


def enum(v, ruta, valores, nulo=False):
    if v is None and nulo:
        return None
    if v not in valores:
        raise ErrorMotor('invalid_argument', ruta)
    return v


def a_json(x):
    if isinstance(x, dt.date):
        return x.isoformat()
    if isinstance(x, Fraction):
        return frac_json(x)
    if isinstance(x, dict):
        return {k: a_json(v) for k, v in x.items()}
    if isinstance(x, (list, tuple)):
        return [a_json(v) for v in x]
    return x
