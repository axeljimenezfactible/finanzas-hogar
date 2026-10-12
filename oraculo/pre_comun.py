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


_ISO = re.compile(r'[0-9]{4}-[0-9]{2}-[0-9]{2}')
# RFC 3339 §5.6 (S5-interfaz §15.1, `created_at` «instante RFC 3339»): full-date «T» partial-time [time-secfrac] time-offset,
# con time-offset = «Z» o ±time-hour:time-minute. Rangos de §5.6/§5.7: hora 00-23, minuto 00-59, segundo 00-59; los mismos en la zona.
# PROVISIONAL (bloqueo S5-precio-vigente-captura pregunta 3): solo «T» y «Z» mayúsculas, sin segundo 60 y sin el año 0000 (dt.date
# no lo admite y el motor también lo rechaza). DIGIT es %x30-39 (RFC 5234): [0-9] y no \d, que acepta dígitos Unicode; y fullmatch,
# porque $ deja pasar un salto de línea final (M-1 de REVISION-motor-S5-r3).
_RFC = re.compile(r'([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})(?:\.([0-9]+))?(?:(Z)|([+-])([0-9]{2}):([0-9]{2}))')


def fecha(v, ruta):
    if not isinstance(v, str) or not _ISO.fullmatch(v):
        raise ErrorMotor('invalid_date', ruta)
    try:
        return dt.date.fromisoformat(v)
    except ValueError:
        raise ErrorMotor('invalid_date', ruta)


def instante(v, ruta):
    """`created_at` (§15.1) como instante exacto: segundos (Fraction) desde 0001-01-01T00:00:00Z.

    La zona se resta para que dos instantes con zonas distintas se comparen en UTC (§6, supuesto S-8). La fracción se toma
    completa, sin truncar (PROVISIONAL, bloqueo S5-precio-vigente-captura pregunta 2)."""
    m = _RFC.fullmatch(v) if isinstance(v, str) else None
    if not m:
        raise ErrorMotor('invalid_date', ruta)
    a, me, d, h, mi, s = (int(m.group(i)) for i in range(1, 7))
    try:
        dia = dt.date(a, me, d)
    except ValueError:
        raise ErrorMotor('invalid_date', ruta)
    zona = 0
    if m.group(8) is None:
        zh, zm = int(m.group(10)), int(m.group(11))
        if zh > 23 or zm > 59:
            raise ErrorMotor('invalid_date', ruta)
        zona = (zh * 60 + zm) * 60 * (1 if m.group(9) == '+' else -1)
    if h > 23 or mi > 59 or s > 59:
        raise ErrorMotor('invalid_date', ruta)
    f = m.group(7)
    fraccion = Fraction(int(f), 10 ** len(f)) if f else Fraction(0)
    return Fraction(dia.toordinal() * 86400 + h * 3600 + mi * 60 + s - zona) + fraccion


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
