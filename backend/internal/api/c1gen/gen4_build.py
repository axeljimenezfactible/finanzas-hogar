# Generador C1 · parte 4: arma el documento OpenAPI y lo escribe
from gen3_ops import *
from c1paths import OUT, WORK
VIEW_ENT = ['person', 'user_account', 'device_session', 'income', 'category', 'expense', 'savings_goal', 'savings_movement', 'credit', 'credit_term_rate', 'credit_purchase', 'scheduled_payment', 'promotion', 'store', 'product', 'price', 'recipe', 'meal_plan_entry', 'inventory_item', 'weight_entry', 'notification_delivery', 'job', 'backup']
TAGS = [(I, 'Inicio: tarjetas del resumen, calendario del hogar, búsqueda y vistas de base de datos.'), (G, 'Gastos, adjuntos y captura de tickets con OCR.'), (P_, 'Presupuesto, flujo y proyección, ingresos y ahorro.'), (CO, 'Plan de comidas y recetario.'), (CT, 'Lista de compras, comparar canasta, productos, precios e inventario.'), (CR, 'Líneas, compras a crédito, pagos, promociones y comparador de escenarios.'), (SA, 'Perfil de salud, peso y calorías.'), (RE, 'Reportes, estado de cuenta de Karen, avisos enviados y jobs.'), (AJ, 'Las 14 secciones de Ajustes, sesiones y sistema.')]
def arr_ids(name, desc): return {'name': name, 'in': 'query', 'required': False, 'description': desc, 'style': 'form', 'explode': False, 'schema': {'type': 'array', 'items': {'type': 'integer', 'minimum': 1}}}
PARAMS = {
 'PersonId': arr_ids('person_id', 'Filtro global «persona»: ids de persona separados por coma (person.id). En gastos filtra expense.person_id; en partes de dinero, la persona de la parte.'),
 'Family': {'name': 'family', 'in': 'query', 'required': False, 'description': 'Filtro global «familia» (U1 §0.2, K5), tres valores incluyentes: household = Hogar, registros con reparto Hogar o Conjunto (incluye lo conjunto); external = Karen, registros con reparto Karen o Conjunto (incluye lo conjunto); joint = Conjunto, solo los registros con reparto Conjunto. Sin valor: todos.', 'schema': {'type': 'string', 'enum': ['household', 'external', 'joint']}},
 'CategoryId': arr_ids('category_id', 'Filtro global «categoría»: ids separados por coma (category.id).'),
 'StoreId': arr_ids('store_id', 'Filtro global «tienda»: ids separados por coma (store.id); en promociones, el comercio ligado.'),
 'Period': {'name': 'period', 'in': 'query', 'required': False, 'description': 'Filtro global «periodo», una sola sintaxis: `2026-10` mes calendario (periodo base); `2026-10-q1` / `2026-10-q2` quincena (q1 del 1 al 15, q2 del 16 al último día; días de pago 15 y 30, y en febrero el último día, p. ej. 2027-02-28); `2026-W42` semana ISO de lunes a domingo; `2026` año calendario; `2026-10-10..2026-11-09` rango inclusivo. Quincenas y semanas según K8 (el 31 cae en q2). Sin valor: el periodo por defecto de cada operación.', 'schema': {'type': 'string', 'pattern': '^([0-9]{4}|[0-9]{4}-(0[1-9]|1[0-2])|[0-9]{4}-(0[1-9]|1[0-2])-q[12]|[0-9]{4}-W(0[1-9]|[1-4][0-9]|5[0-3])|[0-9]{4}-[0-9]{2}-[0-9]{2}\\.\\.[0-9]{4}-[0-9]{2}-[0-9]{2})$', 'examples': ['2026-10']}},
 'CreditId': arr_ids('credit_id', 'Filtro global «crédito»: ids de línea separados por coma (credit.id). En gastos, la forma de pago a crédito con esa línea.'),
 'Cursor': {'name': 'cursor', 'in': 'query', 'required': False, 'description': 'Cursor opaco devuelto en `next_cursor` de la página anterior. Sin valor: primera página.', 'schema': {'type': 'string', 'maxLength': 512}},
 'Limit': {'name': 'limit', 'in': 'query', 'required': False, 'description': 'Registros por página: 50 por defecto (U1); máximo 200 (K6).', 'schema': {'type': 'integer', 'minimum': 1, 'maximum': 200, 'default': 50}},
 'Sort': {'name': 'sort', 'in': 'query', 'required': False, 'description': 'Orden: propiedades separadas por coma; prefijo `-` = descendente (p. ej. `-expense_date,-created_at`). El orden por defecto de cada lista está en su descripción. El id cierra siempre el orden para que el cursor sea estable.', 'schema': {'type': 'string', 'pattern': '^-?[a-z][a-z0-9_]*(,-?[a-z][a-z0-9_]*)*$'}},
 'BudgetScope': {'name': 'scope', 'in': 'query', 'required': False, 'description': 'Alcance de un cambio de monto (K2): this_month = «solo este mes»; onward = «de este mes en adelante» (por omisión; preseleccionado en pantalla).', 'schema': {'type': 'string', 'enum': ['this_month', 'onward'], 'default': 'onward'}},
 'BudgetMonth': {'name': 'month', 'in': 'query', 'required': False, 'description': 'Mes al que aplica el cambio de monto (AAAA-MM): el que muestra la pantalla. Sin valor, el mes en curso.', 'schema': {'$ref': '#/components/schemas/Month'}},
 'ListWeek': {'name': 'week', 'in': 'query', 'required': True, 'description': 'Llave de la lista de compras semanal (R3-01): dice qué lista se lee o se escribe. No es un filtro ni sustituye al filtro global `period` (C1-datos.md §6, «Llaves que no son filtros»). Semana ISO `2026-W42`, de lunes a domingo (K8). Obligatoria: la pantalla manda la semana que muestra; el contrato no tiene semana por omisión (la regla de la lista vigente vive en la PWA, U1 §0.2).', 'schema': {'type': 'string', 'pattern': '^[0-9]{4}-W(0[1-9]|[1-4][0-9]|5[0-3])$', 'examples': ['2026-W42']}},
 'Q': {'name': 'q', 'in': 'query', 'required': False, 'description': 'Texto libre, sin distinguir mayúsculas ni acentos.', 'schema': {'type': 'string', 'maxLength': 200}},
 'IdempotencyKey': {'name': 'Idempotency-Key', 'in': 'header', 'required': True, 'description': 'Obligatorio en todo POST de captura. Clave única por intento (UUID v4 recomendado). Repetir la misma clave con el mismo cuerpo devuelve la respuesta original con `Idempotent-Replayed: true`; con otro cuerpo, 409 `idempotency_key_reused`. La API la guarda por cuenta en la tabla idempotency_key de S1 (huella SHA-256 del cuerpo y respuesta); retención de 24 h al sembrar (K6; valor vigente de app_setting idempotency.retention_hours, que se cambia con updateDataSettings, K44). Gastos y pesos guardan además la clave en su propio renglón, sin purga, y responden 200 con el mismo recurso al reintento (K9).', 'schema': {'type': 'string', 'minLength': 16, 'maxLength': 128, 'pattern': '^[A-Za-z0-9_-]+$'}},
}
ERRS = {'400': ('BadRequest', 'Petición mal formada: JSON inválido o parámetro con formato incorrecto (`bad_request`).'), '401': ('Unauthorized', 'Sin sesión, sesión vencida o cerrada desde otro dispositivo (`unauthorized`).'), '403': ('Forbidden', 'Acción no permitida sobre ese registro (`forbidden`, `read_only_source`).'), '404': ('NotFound', 'El registro no existe (`not_found`).'), '409': ('Conflict', 'Conflicto con el estado actual (`conflict`, `already_exists`, `in_use`, `budget_overlap`, `idempotency_key_reused`, …).'), '413': ('PayloadTooLarge', 'Archivo demasiado grande (`payload_too_large`).'), '503': ('ServiceUnavailable', 'El cálculo pasó su tiempo máximo (`unavailable`); no se guardó nada.'), '415': ('UnsupportedMediaType', 'Tipo de archivo no admitido (`unsupported_media_type`).'), '422': ('UnprocessableEntity', 'Validación de negocio o de campos (`validation_failed`, con `fields`).'), '429': ('TooManyRequests', 'Demasiados intentos de inicio de sesión (`too_many_attempts`, con `retry_at`).')}
RESPONSES = {name: {'description': d, 'headers': {'X-Request-Id': ref('RequestId', 'headers')}, 'content': {'application/json': {'schema': ref('Error')}}} for code, (name, d) in ERRS.items()}
RESPONSES['Error'] = {'description': 'Cualquier otro error, incluidos 500 `internal_error`, 501 `not_implemented` y 503 `unavailable`, con el esquema único.', 'headers': {'X-Request-Id': ref('RequestId', 'headers')}, 'content': {'application/json': {'schema': ref('Error')}}}
HEADERS = {'RequestId': {'description': 'Identificador de la petición (también en Error.request_id y audit_log.request_id).', 'schema': {'type': 'string'}}, 'IdempotentReplayed': {'description': 'true cuando la respuesta es la repetición de una captura ya hecha con la misma Idempotency-Key.', 'schema': {'type': 'boolean'}}, 'SetCookie': {'description': 'Cookie de sesión `fh_session`: HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age de la sesión larga (90 días renovables con el uso, K11).', 'schema': {'type': 'string'}}}
MULTI = {'attachment': ('Adjunto de un gasto.', {'kind': {'type': 'string', 'enum': ['ticket_photo', 'invoice', 'bank_statement', 'other'], 'description': 'Tipo (attachment.kind).'}}, ['image/jpeg', 'image/png', 'image/webp', 'application/pdf']),
         'ticket': ('Foto del ticket.', {}, ['image/jpeg', 'image/png', 'image/webp']),
         'csv': ('Archivo gastos_fijos_ingresos.csv (plantilla de plan.md §8.2).', {}, ['text/csv']),
         'profeco': ('Corte de datos abiertos PROFECO (RAR o CSV).', {}, ['application/vnd.rar', 'text/csv'])}
def path_param(n):
    if n == 'entity': s = {'type': 'string', 'enum': VIEW_ENT}; d = 'Entidad (nombre_codigo de la tabla) de la vista.'
    elif n == 'provider': s = {'type': 'string', 'enum': ['local_model', 'claude_api']}; d = 'Proveedor de IA.'
    else: s = {'type': 'integer', 'minimum': 1}; d = f'Id del registro ({n}).'
    return {'name': n, 'in': 'path', 'required': True, 'description': d, 'schema': s}
def build():
    schemas = dict(S); paths = {}
    for o in OPS:
        if o['list_of']:
            it = o['list_of']; pg = it + 'Page' if it != 'ViewItem' else 'ViewPage'
            if pg not in schemas:
                ext = {'type': 'object', 'properties': {'items': {'type': 'array', 'items': ref(it)}}}
                if it == 'ViewItem':
                    ext['properties']['mode'] = {'type': 'string', 'enum': ['table', 'board', 'calendar', 'gallery'], 'description': 'Modo pedido.'}
                    ext['properties']['groups'] = {'type': 'array', 'items': ref('ViewGroup'), 'description': 'Columnas del tablero o días del calendario (vacío en tabla y galería).'}
                    ext['required'] = ['mode', 'groups']
                schemas[pg] = {'allOf': [ref('Page'), ext], 'description': f'Página de {it} con la envoltura única (Page).'}
            o['resp'] = pg
        oper = {'operationId': o['oid'], 'tags': [o['tag']], 'summary': o['summary']}
        desc = o['desc'] or o['summary'] + '.'
        oper['description'] = desc
        if 'PROVISIONAL' in desc: oper['x-provisional'] = True
        oper['x-necesidades'] = o['needs']
        # Los POST de captura leen y escriben su clave en idempotency_key (S1 1.1, C-15; corrección r2).
        oper['x-tablas'] = sorted(set(o['tables']) | ({'idempotency_key'} if o['idem'] else set()))
        prm = [path_param(n) for n in re.findall(r'{(\w+)}', o['path'])]
        names = list(o['params'])
        if o['list_of']: names += ['Cursor', 'Limit', 'Sort']
        if o['idem']: names.append('IdempotencyKey')
        seen = set()
        for p in names:
            if isinstance(p, str):
                if p in seen: continue
                seen.add(p); prm.append(ref(p, 'parameters'))
            else: prm.append(p)
        if prm: oper['parameters'] = prm
        if o['public']: oper['security'] = []
        if o['req']:
            oper['requestBody'] = {'required': True, 'content': {'application/json': {'schema': ref(o['req'])}}}
        if o['multipart']:
            d, extra, mts = MULTI[o['multipart']]
            props = {'file': {'type': 'string', 'contentMediaType': 'application/octet-stream', 'description': d + ' Tipos admitidos: ' + ', '.join(mts) + '.'}}; props.update(extra)
            oper['requestBody'] = {'required': True, 'content': {'multipart/form-data': {'schema': {'type': 'object', 'required': ['file'], 'properties': props}, 'encoding': {'file': {'contentType': ', '.join(mts)}}}}}
        code = o['code'] or '200'
        hdr = {'X-Request-Id': ref('RequestId', 'headers')}
        if o['idem']: hdr['Idempotent-Replayed'] = ref('IdempotentReplayed', 'headers')
        if o['oid'] == 'login': hdr['Set-Cookie'] = ref('SetCookie', 'headers')
        okd = {'200': 'Correcto.', '201': 'Creado.', '202': 'Aceptado: el trabajo quedó registrado o en curso.', '204': 'Hecho, sin contenido.'}[code]
        r = {'description': okd, 'headers': hdr}
        if o['binary']:
            mt, bd = o['binary']; r['content'] = {mt: {'schema': {'type': 'string', 'contentMediaType': mt, 'description': bd}}}
        elif code != '204': r['content'] = {'application/json': {'schema': ref(o['resp'])}}
        resps = {}
        if o.get('replay'):  # K9: el reintento de una captura de la cola responde 200 con el mismo recurso
            resps['200'] = {'description': 'Reintento de una captura ya registrada con la misma Idempotency-Key (K9): devuelve el mismo recurso, sin crear otro, con Idempotent-Replayed: true.', 'headers': hdr, 'content': {'application/json': {'schema': ref(o['resp'])}}}
        resps[code] = r
        errs = ['400']
        if not o['public']: errs.append('401')
        if prm and any(isinstance(p, dict) and p.get('in') == 'path' for p in prm): errs.append('404')
        if o['m'] in ('post', 'patch', 'put', 'delete') and o['oid'] not in ('login', 'logout', 'previewCreditPurchase', 'compareScenarios', 'previewPlannedItemsImport'): errs.append('409')
        if o['req'] or o['multipart'] or o['m'] == 'get' and o['params']: errs.append('422')
        errs += o['extra_err']
        for e in sorted(set(errs)): resps[e] = ref(ERRS[e][0], 'responses')
        resps['default'] = ref('Error', 'responses')
        oper['responses'] = resps
        paths.setdefault(o['path'], {})[o['m']] = oper
    doc = {'openapi': '3.1.1',
      'info': {'title': 'Finanzas Hogar · API (C1)', 'version': '1.2.0', 'summary': 'Contrato C1 v1.2 del backend Go local de Finanzas Hogar.',
               'description': 'Contrato C1 v1.2 (11-oct-2026, respuestas de las puertas 1 y 1b). Convenciones únicas (detalle en docs/contratos/C1-datos.md): dinero entero en centavos MXN con sufijo `_cents` (esquema MoneyCents, ROUND_HALF_UP, el último pago de una serie ajusta centavos); fechas ISO `AAAA-MM-DD` (Date) e instantes con desplazamiento de America/Mexico_City (DateTime; «Z» solo en entradas); un solo esquema de error (Error); paginación por cursor con la envoltura Page; seis filtros globales (person_id, family, category_id, store_id, period, credit_id); `Idempotency-Key` obligatoria en los POST de captura; parámetros operativos (K44) por clave de app_setting, con valor vigente, valor inicial y «restablecer», en las operaciones /api/v1/settings/*. Solo `POST /api/v1/auth/login` y `GET /health` son públicos; todo lo demás exige la cookie de sesión. Axel y Ximena tienen los mismos permisos; Karen no tiene cuenta. Extensiones: `x-tablas` (tablas de S1 que lee o escribe la operación), `x-necesidades` (necesidades N-xxx de U1), `x-calculado` (campo calculado, no es columna), `x-entrada` (dato de acción que no se guarda), `x-columna` (columna a la que corresponde una entrada en lista), `x-app-setting` (clave de app_setting), `x-tabla` (tabla de un esquema de entidad) y `x-provisional` (depende de un bloqueo abierto).'},
      'servers': [{'url': 'https://finanzas.local:8443', 'description': 'Máquina de casa, solo red local (decisión H).'}],
      'security': [{'cookieAuth': []}],
      'tags': [{'name': n, 'description': d} for n, d in TAGS],
      'paths': paths,
      'components': {'securitySchemes': {'cookieAuth': {'type': 'apiKey', 'in': 'cookie', 'name': 'fh_session', 'description': 'Sesión larga por dispositivo. Cookie HttpOnly, Secure y SameSite=Strict; el token en claro solo viaja en la cookie (la base guarda su SHA-256 en device_session.token_hash). Duración 90 días renovable con el uso (K11, app_setting session.max_age_days).'}},
                     'parameters': PARAMS, 'headers': HEADERS, 'responses': RESPONSES, 'schemas': schemas}}
    return doc
if __name__ == '__main__':
    import sys
    doc = build()
    out = sys.argv[1] if len(sys.argv) > 1 else os.path.join(OUT, 'docs/contratos/openapi.yaml')
    head = ('# Finanzas Hogar · Contrato C1 v1.1 (OpenAPI 3.1). Archivo generado: no se edita a mano.\n'
            '# Generador: backend/internal/api/c1gen/ (agente S6a-contrato-c1, WF-0, 2026-10-10; v1.1: agente C1-v1.1, puerta 1, 2026-10-11). Regenerar: python3 -I backend/internal/api/c1gen/build.py (con --check compara sin escribir).\n'
            '# Fuente: esquema de S1 (backend/db/migrations/*.sql, aplicadas en orden) + necesidades de U1 (docs/design/U1-arquitectura.md). Ver docs/contratos/C1-datos.md.\n')
    open(out, 'w', encoding='utf-8').write(head + '\n'.join(emit(doc)) + '\n')
    json.dump(doc, open(os.path.join(WORK, 'openapi.json'), 'w', encoding='utf-8'), ensure_ascii=False)
    print('escrito', out, len(doc['paths']), 'rutas', sum(len(v) for v in doc['paths'].values()), 'operaciones', len(doc['components']['schemas']), 'esquemas')
