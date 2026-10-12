"""Oráculo financiero S2 · piezas comunes (S2-interfaz.md §1, §2, §3.15).

Errores, redondeo HALF_UP exacto con Fraction, fechas y lectura validada de `Input`.
Independiente del motor Go (S2-03): no lee nada de backend/. Dinero en centavos enteros.
"""
import calendar
import datetime as dt
import re
from fractions import Fraction

MESES = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']
ZONA = 'America/Mexico_City'


class ErrorMotor(Exception):
    """Error de la interfaz (§3.15): `code` es el código; `detalle`, la ruta o el motivo."""

    def __init__(self, code, detalle=''):
        super().__init__(f'{code}: {detalle}')
        self.code = code
        self.detalle = detalle


def hu(x):
    """HALF_UP a entero desde una fracción exacta; la mitad se aleja del cero (RN-01)."""
    x = Fraction(x)
    a = abs(x)
    n, r = divmod(a.numerator, a.denominator)
    if 2 * r >= a.denominator:
        n += 1
    return -n if x < 0 else n


# ---- fechas ---------------------------------------------------------------
def fecha(v, ruta):
    if not isinstance(v, str) or not re.fullmatch(r'\d{4}-\d{2}-\d{2}', v):
        raise ErrorMotor('invalid_date', ruta)
    try:
        return dt.date.fromisoformat(v)
    except ValueError:
        raise ErrorMotor('invalid_date', ruta) from None


def fecha_o_nulo(v, ruta):
    return None if v is None else fecha(v, ruta)


def instante(v, ruta):
    if not isinstance(v, str) or not re.fullmatch(r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z', v):
        raise ErrorMotor('invalid_date', ruta)
    try:
        return dt.datetime.fromisoformat(v[:-1] + '+00:00')
    except ValueError:
        raise ErrorMotor('invalid_date', ruta) from None


def mes(v, ruta):
    m = re.fullmatch(r'(\d{4})-(\d{2})', v) if isinstance(v, str) else None
    if not m or not 1 <= int(m[2]) <= 12:
        raise ErrorMotor('invalid_date', ruta)
    return (int(m[1]), int(m[2]))


def quincena(v, ruta):
    m = re.fullmatch(r'(\d{4})-(\d{2})-q([12])', v) if isinstance(v, str) else None
    if not m or not 1 <= int(m[2]) <= 12:
        raise ErrorMotor('invalid_date', ruta)
    return (int(m[1]), int(m[2]), int(m[3]))


def ultimo(y, m):
    return calendar.monthrange(y, m)[1]


def mes_de(d):
    return (d.year, d.month)


def mes_n(m):
    return m[0] * 12 + m[1] - 1


def mes_mas(m, k):
    t = mes_n(m) + k
    return (t // 12, t % 12 + 1)


def mes_txt(m):
    return f'{m[0]:04d}-{m[1]:02d}'


def lim_mes(m):
    return dt.date(m[0], m[1], 1), dt.date(m[0], m[1], ultimo(*m))


def lab(d):
    return f'{d.day}-{MESES[d.month - 1]}-{d.year}'


def etiqueta(a, b):
    return f'{lab(a)} a {lab(b)}'


def dias(a, b):
    d = a
    while d <= b:
        yield d
        d += dt.timedelta(days=1)


def quincena_de(d):
    return (d.year, d.month, 1 if d.day <= 15 else 2)


def lim_quincena(q):
    """(inicio, fin, día de pago) con K8 y RN-06: q1 = 1..15 paga el 15; q2 = 16..último, paga el 30 o el último."""
    y, m, h = q
    u = ultimo(y, m)
    if h == 1:
        return dt.date(y, m, 1), dt.date(y, m, 15), dt.date(y, m, 15)
    return dt.date(y, m, 16), dt.date(y, m, u), dt.date(y, m, min(30, u))


def q_txt(q):
    return f'{q[0]:04d}-{q[1]:02d}-q{q[2]}'


def q_mas(q, k):
    t = (q[0] * 12 + q[1] - 1) * 2 + q[2] - 1 + k
    return (t // 24, (t // 2) % 12 + 1, t % 2 + 1)


def semana_txt(d):
    y, w, _ = d.isocalendar()
    return f'{y:04d}-W{w:02d}'


def a_json(x):
    """Convierte fechas a ISO y tuplas a listas, recursivamente."""
    if isinstance(x, dt.date):
        return x.isoformat()
    if isinstance(x, dict):
        return {k: a_json(v) for k, v in x.items()}
    if isinstance(x, (list, tuple)):
        return [a_json(v) for v in x]
    return x


# ---- validación de entradas (§2) ------------------------------------------
CAMPOS = {
    'Input': ('params', 'persons', 'categories', 'expenses', 'credit_purchases', 'incomes', 'planned_items',
              'budgets', 'savings_goals', 'savings_movements', 'cadence', 'stable_items'),
    'Params': ('today', 'now', 'zone', 'holder', 'groceries_category', 'paper_category', 'repeat_cadence_enabled',
               'budget_warning_pct', 'budget_limit_pct'),
    'Person': ('code', 'display_name', 'participation', 'pays_share_directly'),
    'Category': ('code', 'name', 'kind', 'is_active'),
    'Expense': ('id', 'kind', 'category', 'amount_cents', 'expense_date', 'accrual_months', 'payment_method', 'paid_by',
                'shares', 'credit_purchase_id', 'planned_item_id', 'projected'),
    'Share': ('person', 'amount_cents'),
    'Purchase': ('id', 'credit_id', 'expense_id', 'category', 'purchase_date', 'principal_cents', 'interest_cents',
                 'total_cents', 'accrual_months', 'status', 'shares', 'embedded', 'payments', 'projected'),
    'PurchaseShare': ('person', 'principal_cents', 'interest_cents', 'total_cents'),
    'Embedded': ('category', 'person', 'amount_cents'),
    'Payment': ('seq', 'due_date', 'amount_cents', 'principal_cents', 'interest_cents', 'shares', 'paid'),
    'Paid': ('person', 'amount_cents', 'paid_date'),
    'Income': ('id', 'person', 'kind', 'amount_cents', 'income_date', 'planned_item_id'),
    'Planned': ('id', 'kind', 'description', 'amount_cents', 'second_half_amount_cents', 'frequency', 'interval_count',
                'start_date', 'end_date', 'category', 'payment_method', 'accrual_months', 'person', 'share_rule',
                'credit_id', 'is_active'),
    'Rule': ('parts', 'rest_person'),
    'RulePart': ('person', 'pct_bp', 'fixed_cents'),
    'Budget': ('id', 'category', 'amount_cents', 'valid_from', 'valid_to'),
    'Goal': ('id', 'name', 'kind', 'target_cents', 'target_date', 'planned_monthly_cents', 'start_date', 'is_active'),
    'Movement': ('goal_id', 'movement_date', 'kind', 'amount_cents'),
    'Cadence': ('purchases', 'expenses'),
    'Stable': ('label', 'cycle_months', 'total_cents', 'shares'),
    'StableShare': ('person', 'total_cents'),
    'Window': ('label', 'start', 'end'),
}


def obj(o, tipo, ruta):
    if not isinstance(o, dict):
        raise ErrorMotor('invalid_argument', ruta)
    for k in o:
        if k not in CAMPOS[tipo]:
            raise ErrorMotor('unknown_field', f'{ruta}.{k}')
    return o


def ent(v, ruta, minimo=None, maximo=None, nulo=False):
    if v is None and nulo:
        return None
    if not isinstance(v, int) or isinstance(v, bool):
        raise ErrorMotor('invalid_argument', ruta)
    if (minimo is not None and v < minimo) or (maximo is not None and v > maximo):
        raise ErrorMotor('invalid_argument', ruta)
    return v


def booleano(v, ruta, defecto):
    if v is None:
        return defecto
    if not isinstance(v, bool):
        raise ErrorMotor('invalid_argument', ruta)
    return v


def lista(o, k, ruta):
    v = o.get(k)
    if v is None:
        return []
    if not isinstance(v, list):
        raise ErrorMotor('invalid_argument', f'{ruta}.{k}')
    return v


def opcion(v, validas, ruta):
    if v not in validas:
        raise ErrorMotor('invalid_argument', ruta)
    return v


class Datos:
    """`Input` validado y normalizado (fechas como date, partes como dict persona -> centavos)."""

    def __init__(self, raw):
        obj(raw, 'Input', 'datos')
        p = obj(raw.get('params'), 'Params', 'params')
        if ('today' in p) == ('now' in p):
            raise ErrorMotor('invalid_argument', 'params: exactamente uno de today/now')
        self.zone = p.get('zone') or ZONA
        if self.zone != ZONA:
            raise ErrorMotor('invalid_zone', 'params.zone')
        if 'today' in p:
            self.today = fecha(p['today'], 'params.today')
        else:
            self.today = hoy_local(instante(p['now'], 'params.now'), self.zone)
        for k in ('holder', 'groceries_category', 'paper_category'):
            if not isinstance(p.get(k), str):
                raise ErrorMotor('invalid_argument', f'params.{k}')
        self.holder = p['holder']
        self.groceries = p['groceries_category']
        self.paper = p['paper_category']
        self.repeat = booleano(p.get('repeat_cadence_enabled'), 'params.repeat_cadence_enabled', False)
        self.warn = ent(p.get('budget_warning_pct', 80), 'params.budget_warning_pct', 0)
        self.limit = ent(p.get('budget_limit_pct', 100), 'params.budget_limit_pct', 0)
        self.persons = []
        for i, x in enumerate(lista(raw, 'persons', 'datos')):
            obj(x, 'Person', f'persons[{i}]')
            opcion(x.get('participation'), ('holder', 'co_holder', 'external'), f'persons[{i}].participation')
            if not isinstance(x.get('code'), str):
                raise ErrorMotor('invalid_argument', f'persons[{i}].code')
            self.persons.append(x)
        if [x['code'] for x in self.persons if x['participation'] == 'holder'] != [self.holder]:
            raise ErrorMotor('invalid_argument', 'persons: un solo titular igual a params.holder')
        self.claves = [x['code'] for x in self.persons if x['participation'] in ('holder', 'external')]
        self.externas = [c for c in self.claves if c != self.holder]
        self.nombre = {x['code']: x.get('display_name') or x['code'] for x in self.persons}
        self.directo = {x['code'] for x in self.persons if x.get('pays_share_directly')}
        self.codigos = {x['code'] for x in self.persons}
        self.categories = {}
        for i, x in enumerate(lista(raw, 'categories', 'datos')):
            obj(x, 'Category', f'categories[{i}]')
            self.categories[x.get('code')] = dict(kind=x.get('kind'), active=booleano(x.get('is_active'), f'categories[{i}].is_active', True), name=x.get('name'))
        self.expenses = [self._gasto(x, f'expenses[{i}]', False) for i, x in enumerate(lista(raw, 'expenses', 'datos'))]
        self.purchases = {}
        for i, x in enumerate(lista(raw, 'credit_purchases', 'datos')):
            c = self._compra(x, f'credit_purchases[{i}]', False)
            if c['id'] in self.purchases:
                raise ErrorMotor('invalid_argument', f'credit_purchases[{i}].id repetido')
            self.purchases[c['id']] = c
        cad = raw.get('cadence') or {}
        obj(cad, 'Cadence', 'cadence')
        self.cad_purchases = {}
        for i, x in enumerate(lista(cad, 'purchases', 'cadence')):
            c = self._compra(x, f'cadence.purchases[{i}]', True)
            if c['id'] in self.purchases or c['id'] in self.cad_purchases:
                raise ErrorMotor('invalid_argument', f'cadence.purchases[{i}].id repetido')
            self.cad_purchases[c['id']] = c
        self.cad_expenses = [self._gasto(x, f'cadence.expenses[{i}]', True) for i, x in enumerate(lista(cad, 'expenses', 'cadence'))]
        self._liga_credito(self.expenses, self.purchases, 'expenses')
        self._liga_credito(self.cad_expenses, {**self.purchases, **self.cad_purchases}, 'cadence.expenses')
        self.incomes = [self._ingreso(x, f'incomes[{i}]') for i, x in enumerate(lista(raw, 'incomes', 'datos'))]
        self.planned = [self._partida(x, f'planned_items[{i}]') for i, x in enumerate(lista(raw, 'planned_items', 'datos'))]
        self.budgets = renglones(lista(raw, 'budgets', 'datos'))
        self.goals = []
        for i, x in enumerate(lista(raw, 'savings_goals', 'datos')):
            r = f'savings_goals[{i}]'
            obj(x, 'Goal', r)
            g = dict(id=ent(x.get('id'), r + '.id'), name=x.get('name'), kind=opcion(x.get('kind'), ('fund', 'goal'), r + '.kind'),
                     target=ent(x.get('target_cents'), r + '.target_cents', 1, nulo=True),
                     target_date=fecha_o_nulo(x.get('target_date'), r + '.target_date'),
                     planned=ent(x.get('planned_monthly_cents'), r + '.planned_monthly_cents', 0, nulo=True),
                     start=fecha_o_nulo(x.get('start_date'), r + '.start_date'),
                     active=booleano(x.get('is_active'), r + '.is_active', True))
            if g['kind'] == 'goal' and (g['target'] is None or g['target_date'] is None):
                raise ErrorMotor('invalid_argument', r + ': una meta lleva monto y fecha')
            self.goals.append(g)
        self.movements = []
        for i, x in enumerate(lista(raw, 'savings_movements', 'datos')):
            r = f'savings_movements[{i}]'
            obj(x, 'Movement', r)
            self.movements.append(dict(goal=ent(x.get('goal_id'), r + '.goal_id'), date=fecha(x.get('movement_date'), r + '.movement_date'),
                                       kind=opcion(x.get('kind'), ('deposit', 'withdrawal'), r + '.kind'),
                                       amount=ent(x.get('amount_cents'), r + '.amount_cents', 1)))
        self.stable_items = partidas_estables(lista(raw, 'stable_items', 'datos'), self, 'stable_items')

    # -- personas de dinero (§1.4)
    def persona(self, v, ruta):
        if v != self.holder and v not in self.externas:
            raise ErrorMotor('invalid_argument', ruta)
        return v

    def _gasto(self, x, r, proyectado):
        obj(x, 'Expense', r)
        g = dict(id=ent(x.get('id'), r + '.id'), kind=opcion(x.get('kind'), ('fixed', 'variable', 'upcoming'), r + '.kind'),
                 cat=x.get('category'), amount=ent(x.get('amount_cents'), r + '.amount_cents'),
                 date=fecha(x.get('expense_date'), r + '.expense_date'),
                 n=ent(x.get('accrual_months', 1) if x.get('accrual_months') is not None else 1, r + '.accrual_months', 1, 24),
                 method=opcion(x.get('payment_method'), ('cash', 'debit', 'transfer', 'credit'), r + '.payment_method'),
                 paid_by=self.persona(x['paid_by'], r + '.paid_by') if x.get('paid_by') is not None else self.holder,
                 cpid=x.get('credit_purchase_id'), pid=x.get('planned_item_id'))
        if booleano(x.get('projected'), r + '.projected', False) != proyectado:
            raise ErrorMotor('invalid_argument', r + '.projected')
        g['shares'] = {}
        for j, s in enumerate(lista(x, 'shares', r)):
            obj(s, 'Share', f'{r}.shares[{j}]')
            per = self.persona(s.get('person'), f'{r}.shares[{j}].person')
            g['shares'][per] = g['shares'].get(per, 0) + ent(s.get('amount_cents'), f'{r}.shares[{j}].amount_cents')
        if sum(g['shares'].values()) != g['amount']:
            raise ErrorMotor('shares_mismatch', r + '.shares')
        if g['method'] == 'credit' and g['cpid'] is None:
            raise ErrorMotor('missing_purchase', r + '.credit_purchase_id')
        return g

    def _compra(self, x, r, proyectado):
        obj(x, 'Purchase', r)
        c = dict(id=ent(x.get('id'), r + '.id'), credit_id=x.get('credit_id'), expense_id=x.get('expense_id'),
                 cat=x.get('category'), date=fecha(x.get('purchase_date'), r + '.purchase_date'),
                 principal=ent(x.get('principal_cents'), r + '.principal_cents', 0),
                 interest=ent(x.get('interest_cents'), r + '.interest_cents', 0),
                 total=ent(x.get('total_cents'), r + '.total_cents', 0),
                 n=ent(x.get('accrual_months') if x.get('accrual_months') is not None else 1, r + '.accrual_months', 1, 24),
                 status=opcion(x.get('status', 'active'), ('active', 'paid_off', 'cancelled'), r + '.status'))
        if booleano(x.get('projected'), r + '.projected', False) != proyectado:
            raise ErrorMotor('invalid_argument', r + '.projected')
        if c['principal'] + c['interest'] != c['total']:
            raise ErrorMotor('sum_mismatch', r + ': capital + intereses != total')
        c['share_tot'] = {}
        for j, s in enumerate(lista(x, 'shares', r)):
            obj(s, 'PurchaseShare', f'{r}.shares[{j}]')
            per = self.persona(s.get('person'), f'{r}.shares[{j}].person')
            c['share_tot'][per] = c['share_tot'].get(per, 0) + ent(s.get('total_cents'), f'{r}.shares[{j}].total_cents')
        if sum(c['share_tot'].values()) != c['total']:
            raise ErrorMotor('sum_mismatch', r + '.shares')
        c['embedded'] = []
        for j, e in enumerate(lista(x, 'embedded', r)):
            obj(e, 'Embedded', f'{r}.embedded[{j}]')
            c['embedded'].append(dict(cat=e.get('category'), person=self.persona(e.get('person'), f'{r}.embedded[{j}].person'),
                                      amount=ent(e.get('amount_cents'), f'{r}.embedded[{j}].amount_cents', 0)))
        c['payments'] = []
        for j, p in enumerate(lista(x, 'payments', r)):
            rp = f'{r}.payments[{j}]'
            obj(p, 'Payment', rp)
            pg = dict(seq=ent(p.get('seq'), rp + '.seq'), due=fecha(p.get('due_date'), rp + '.due_date'),
                      amount=ent(p.get('amount_cents'), rp + '.amount_cents'), principal=ent(p.get('principal_cents'), rp + '.principal_cents'),
                      interest=ent(p.get('interest_cents'), rp + '.interest_cents'), shares={}, paid=[])
            if pg['principal'] + pg['interest'] != pg['amount']:
                raise ErrorMotor('sum_mismatch', rp + ': monto != capital + intereses')
            for k, s in enumerate(lista(p, 'shares', rp)):
                obj(s, 'Share', f'{rp}.shares[{k}]')
                per = self.persona(s.get('person'), f'{rp}.shares[{k}].person')
                pg['shares'][per] = pg['shares'].get(per, 0) + ent(s.get('amount_cents'), f'{rp}.shares[{k}].amount_cents')
            if sum(pg['shares'].values()) != pg['amount']:
                raise ErrorMotor('sum_mismatch', rp + '.shares')
            for k, a in enumerate(lista(p, 'paid', rp)):
                obj(a, 'Paid', f'{rp}.paid[{k}]')
                pg['paid'].append(dict(person=self.persona(a.get('person'), f'{rp}.paid[{k}].person'),
                                       amount=ent(a.get('amount_cents'), f'{rp}.paid[{k}].amount_cents', 0),
                                       date=fecha(a.get('paid_date'), f'{rp}.paid[{k}].paid_date')))
            c['payments'].append(pg)
        if sum(p['amount'] for p in c['payments']) != c['total']:
            raise ErrorMotor('sum_mismatch', r + ': suma de montos != total')
        if sum(p['principal'] for p in c['payments']) != c['principal']:
            raise ErrorMotor('sum_mismatch', r + ': suma de capital != capital')
        return c

    def _liga_credito(self, gastos, compras, r):
        for i, g in enumerate(gastos):
            if g['method'] != 'credit':
                continue
            c = compras.get(g['cpid'])
            if c is None:
                raise ErrorMotor('missing_purchase', f'{r}[{i}].credit_purchase_id')
            if g['amount'] != c['principal']:
                raise ErrorMotor('sum_mismatch', f'{r}[{i}].amount_cents != principal_cents')

    def _ingreso(self, x, r):
        obj(x, 'Income', r)
        if x.get('person') != self.holder:
            raise ErrorMotor('income_not_holder', r + '.person')
        return dict(id=x.get('id'), person=x['person'], kind=opcion(x.get('kind'), ('fixed', 'extra'), r + '.kind'),
                    amount=ent(x.get('amount_cents'), r + '.amount_cents', 1), date=fecha(x.get('income_date'), r + '.income_date'),
                    pid=x.get('planned_item_id'))

    def _partida(self, x, r):
        obj(x, 'Planned', r)
        it = dict(id=ent(x.get('id'), r + '.id'), kind=opcion(x.get('kind'), ('income', 'fixed_expense', 'upcoming_expense'), r + '.kind'),
                  amount=ent(x.get('amount_cents'), r + '.amount_cents', 0),
                  second=ent(x.get('second_half_amount_cents'), r + '.second_half_amount_cents', 0, nulo=True),
                  freq=opcion(x.get('frequency'), ('weekly', 'fortnightly', 'monthly', 'yearly', 'once'), r + '.frequency'),
                  interval=ent(x.get('interval_count') if x.get('interval_count') is not None else 1, r + '.interval_count', 1, 24),
                  start=fecha(x.get('start_date'), r + '.start_date'), end=fecha_o_nulo(x.get('end_date'), r + '.end_date'),
                  cat=x.get('category'), n=ent(x.get('accrual_months') if x.get('accrual_months') is not None else 1, r + '.accrual_months', 1, 24),
                  person=x.get('person'), credit_id=x.get('credit_id'),
                  active=booleano(x.get('is_active'), r + '.is_active', True), rule=None)
        if it['person'] is not None and it['person'] not in self.codigos:
            raise ErrorMotor('invalid_argument', r + '.person')
        if it['kind'] == 'income':
            if it['person'] not in (None, self.holder):
                raise ErrorMotor('income_not_holder', r + '.person')
            return it
        regla = x.get('share_rule')
        if regla is None:
            raise ErrorMotor('invalid_argument', r + '.share_rule')
        obj(regla, 'Rule', r + '.share_rule')
        partes = []
        for j, pt in enumerate(lista(regla, 'parts', r + '.share_rule')):
            rr = f'{r}.share_rule.parts[{j}]'
            obj(pt, 'RulePart', rr)
            if ('pct_bp' in pt) == ('fixed_cents' in pt):
                raise ErrorMotor('invalid_argument', rr)
            per = self.persona(pt.get('person'), rr + '.person')
            if 'pct_bp' in pt:
                partes.append((per, 'pct', ent(pt['pct_bp'], rr + '.pct_bp', 0, 10000)))
            else:
                partes.append((per, 'fixed', ent(pt['fixed_cents'], rr + '.fixed_cents', 0)))
        it['rule'] = dict(parts=partes, rest=self.persona(regla.get('rest_person'), r + '.share_rule.rest_person'))
        reparto(it['rule'], it['amount'], r)
        if it['second'] is not None:
            reparto(it['rule'], it['second'], r)
        return it


def reparto(regla, monto, ruta='share_rule'):
    """RN-03 sobre una ocurrencia: porcentaje = HALF_UP(monto·pb/10000); fija = fixed_cents; el resto absorbe."""
    partes = {}
    usado = 0
    for per, tipo, v in regla['parts']:
        x = hu(Fraction(monto * v, 10000)) if tipo == 'pct' else v
        if x > monto:
            raise ErrorMotor('invalid_argument', ruta + ': parte mayor que el monto')
        partes[per] = partes.get(per, 0) + x
        usado += x
    if usado > monto:
        raise ErrorMotor('invalid_argument', ruta + ': partes mayores que el monto')
    partes[regla['rest']] = partes.get(regla['rest'], 0) + monto - usado
    return partes


def renglones(raw):
    """BudgetRow validados; sin traslapes por categoría (budget_overlap)."""
    out = []
    for i, x in enumerate(raw):
        r = f'budgets[{i}]'
        obj(x, 'Budget', r)
        b = dict(id=ent(x.get('id'), r + '.id'), category=x.get('category'), amount=ent(x.get('amount_cents'), r + '.amount_cents', 0),
                 valid_from=fecha(x.get('valid_from'), r + '.valid_from'), valid_to=fecha_o_nulo(x.get('valid_to'), r + '.valid_to'))
        if b['valid_from'].day != 1 or (b['valid_to'] is not None and (b['valid_to'] != lim_mes(mes_de(b['valid_to']))[1] or b['valid_to'] < b['valid_from'])):
            raise ErrorMotor('invalid_argument', r + ': vigencia')
        out.append(b)
    for i, a in enumerate(out):
        for b in out[i + 1:]:
            if a['category'] == b['category']:
                fa = a['valid_to'] or dt.date.max
                fb = b['valid_to'] or dt.date.max
                if a['valid_from'] <= fb and b['valid_from'] <= fa:
                    raise ErrorMotor('budget_overlap', f'renglones {a["id"]} y {b["id"]}')
    return out


def partidas_estables(raw, D, ruta):
    out = []
    for i, x in enumerate(raw):
        r = f'{ruta}[{i}]'
        obj(x, 'Stable', r)
        it = dict(label=x.get('label'), cycle=ent(x.get('cycle_months'), r + '.cycle_months', 1),
                  total=ent(x.get('total_cents'), r + '.total_cents'), shares={})
        for j, s in enumerate(lista(x, 'shares', r)):
            obj(s, 'StableShare', f'{r}.shares[{j}]')
            per = D.persona(s.get('person'), f'{r}.shares[{j}].person')
            it['shares'][per] = it['shares'].get(per, 0) + ent(s.get('total_cents'), f'{r}.shares[{j}].total_cents')
        out.append(it)
    return out


def hoy_local(t, zona=ZONA):
    """Fecha civil de un instante en la zona (C-05, S2-14), con la base de zonas del sistema."""
    from zoneinfo import ZoneInfo
    if zona != ZONA:
        raise ErrorMotor('invalid_zone', zona)
    return t.astimezone(ZoneInfo(zona)).date()


def como_datos(datos):
    return datos if isinstance(datos, Datos) else Datos(datos)


def mapa(D, acc):
    """Mapa {code: cents}: titular y externas siempre, en el orden de persons (§1.1)."""
    return {c: acc.get(c, 0) for c in D.claves}


def por_categoria(acc):
    return {k: acc[k] for k in sorted(acc, key=lambda s: s.encode()) if acc[k] != 0}
