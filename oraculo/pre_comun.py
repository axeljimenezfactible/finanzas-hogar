"""Oráculo S5 · piezas comunes: errores, redondeo HALF_UP, fechas y validación estricta de entradas.

Implementa S5-interfaz.md §1 (aritmética exacta), §15.1 y §15.2 (unknown_field, missing_field, invalid_date)
sin leer el código Go. Solo biblioteca estándar. Dinero en centavos, cantidades en milésimas, todo exacto con Fraction.
"""
import datetime as dt
import re
from fractions import Fraction


class ErrorMotor(Exception):
    """Error del paquete (S5-interfaz §17 y §15.2): `code` en snake_case sin «Err»."""

    def __init__(self, code, detalle=''):
        super().__init__(f'{code}: {detalle}')
        self.code = code
        self.detalle = detalle


def half_up(x):
    """ROUND_HALF_UP (la mitad se aleja de cero), una sola vez al publicar."""
    x = Fraction(x)
    n, d = x.numerator, x.denominator
    q = (2 * abs(n) + d) // (2 * d)
    return q if n >= 0 else -q


def techo(x):
    x = Fraction(x)
    return -((-x.numerator) // x.denominator)


def frac(x):
    x = Fraction(x)
    return {'num': x.numerator, 'den': x.denominator}


_ISO = re.compile(r'^\d{4}-\d{2}-\d{2}$')
_RFC = re.compile(r'^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z$')


def fecha(v, ruta):
    if not isinstance(v, str) or not _ISO.match(v):
        raise ErrorMotor('invalid_date', ruta)
    try:
        return dt.date.fromisoformat(v)
    except ValueError:
        raise ErrorMotor('invalid_date', ruta)


def instante(v, ruta):
    if not isinstance(v, str) or not _RFC.match(v):
        raise ErrorMotor('invalid_date', ruta)
    try:
        return dt.datetime.fromisoformat(v[:-1] + '+00:00')
    except ValueError:
        raise ErrorMotor('invalid_date', ruta)


def objeto(x, ruta, oblig, opc=()):
    if not isinstance(x, dict):
        raise ErrorMotor('invalid_argument', ruta)
    for k in x:
        if k not in oblig and k not in opc:
            raise ErrorMotor('unknown_field', f'{ruta}.{k}')
    for k in oblig:
        if k not in x:
            raise ErrorMotor('missing_field', f'{ruta}.{k}')
    return x


def lista(x, ruta):
    if not isinstance(x, list):
        raise ErrorMotor('invalid_argument', ruta)
    return x


def texto(x, ruta, nulo=False):
    if x is None and nulo:
        return None
    if not isinstance(x, str) or x == '':
        raise ErrorMotor('invalid_argument', ruta)
    return x


def entero(x, ruta, minimo=None, nulo=False):
    if x is None and nulo:
        return None
    if isinstance(x, bool) or not isinstance(x, int) or (minimo is not None and x < minimo):
        raise ErrorMotor('invalid_argument', ruta)
    return x


def booleano(x, ruta):
    if not isinstance(x, bool):
        raise ErrorMotor('invalid_argument', ruta)
    return x


def enum(x, ruta, valores):
    if x not in valores:
        raise ErrorMotor('invalid_argument', ruta)
    return x
