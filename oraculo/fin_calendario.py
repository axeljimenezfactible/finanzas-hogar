"""Oráculo S2 · calendario (S2-interfaz.md §1.2, §3.1; S2-06, S2-14)."""
import datetime as dt

from fin_comun import (ErrorMotor, como_datos, etiqueta, fecha, hoy_local, instante, lim_mes, lim_quincena, mes_de,
                       mes_txt, q_txt, quincena_de, semana_txt, ultimo)


def hoy(now, zona):
    return hoy_local(now, zona)


def limites_mes(m):
    return lim_mes(m)


def limites_quincena(q):
    return lim_quincena(q)


def limites_semana(d):
    lunes = d - dt.timedelta(days=d.weekday())
    return lunes, lunes + dt.timedelta(days=6)


def semana_de(d):
    return semana_txt(d)


def calendario(datos, date=None, now=None):
    D = como_datos(datos)
    if (date is None) == (now is None):
        raise ErrorMotor('invalid_argument', 'argumentos: exactamente uno de date/now')
    d = fecha(date, 'argumentos.date') if date is not None else hoy(instante(now, 'argumentos.now'), D.zone)
    m = mes_de(d)
    a, b = lim_mes(m)
    q = quincena_de(d)
    qa, qb, pay = lim_quincena(q)
    wa, wb = limites_semana(d)
    return dict(today=d, month=mes_txt(m), month_start=a, month_end=b, month_label=etiqueta(a, b), days_in_month=ultimo(*m),
                fortnight=q_txt(q), fortnight_start=qa, fortnight_end=qb, pay_date=pay, week=semana_de(d),
                week_start=wa, week_end=wb)
