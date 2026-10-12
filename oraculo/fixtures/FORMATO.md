# Formato común de fixtures de los motores

Versión 1.2 · 11-oct-2026 · WF-1 (1.0: agente `contrato-S2`; 1.1: agente `correccion-contrato-S2-r1`; 1.2: agente `correccion-contrato-S2-r2`). Rige para todos los motores (S2 `motor_fin`, S3 `motor_cred`, S4, S5 y los que vengan), sin excepciones (§6). Cada motor declara en su interfaz (`docs/contratos/motores/<motor>-interfaz.md`) sus funciones, sus tipos de entrada y sus casos; este archivo fija solo la forma común, la salida del motor y del oráculo, la comparación y quién escribe las piezas compartidas.

**Cambios de la 1.1** (hallazgos M1 y M2 de `docs/estado/wf1/REVISION-contrato-S2-r1.md`): §2.1 nuevo, regla única para ligar una consulta con su hoja de `model.json` o `data.json` y vistas de cruce con claves literales (M1); §3, razones con `decimales`; §4, el recorte deja las listas completas y solo recorta claves de objetos (M1); §4.1 nuevo, dueño y API de registro del binario común `backend/cmd/motores` (M2); §5, dueño de `oraculo/compara.py`, ramas de todos los motores y línea de `data.json` (M2); §6 nuevo, regla para S3 y para los motores que lleguen (M2).

**Cambios de la 1.2** (hallazgos M1, M2 y M5 de `docs/estado/wf1/REVISION-contrato-S2-r2.md`; ninguna regla de comparación cambia): §1, cada interfaz declara el autor de cada archivo de su carpeta y qué hace el otro agente mientras falten (M5); §2, `argumentos` lleva las claves que la interfaz declara para cada función (M1); §4, regla 5: los mapas `{code: cents}` siguen la regla de contenido de la interfaz (M2).

## 1. Carpetas y nombres
- Un archivo por caso: `oraculo/fixtures/<motor>/<caso>.json`, con `<motor>` en minúsculas (`s2`, `s3`, …) y `<caso>` = `<motor>-<tema>[-<variante>]` en minúsculas con guiones (`s2-periodos-p1`, `s2-fin-de-mes-feb-2027`). El nombre del archivo es el valor de `caso`.
- Datos compartidos por varios casos: `oraculo/fixtures/<motor>/_datos/<nombre>.json` (lo que empieza con `_` no es un caso).
- **Autoría:** la interfaz de cada motor declara quién escribe cada archivo de `oraculo/fixtures/<motor>/` (casos y `_datos/`) y qué hace el otro agente mientras no existan (S2: `S2-interfaz.md` §5.0). Un archivo sin autor declarado es un hueco de la interfaz, no de quien lo necesita.
- Salida del oráculo: `oraculo/salidas/<motor>.json`. Salida del motor Go: `datos/tmp/motor-<motor>.json` (temporal; se borra al terminar).
- UTF-8, JSON estricto (sin comentarios ni comas finales), sangría de 2 espacios, claves en el orden en que se escriben, salto de línea al final.

## 2. Forma de una fixture
```json
{
  "formato": "fixture/1",
  "motor": "s2",
  "caso": "s2-periodos-p1",
  "descripcion": "Octubre de 2026 (P1 del oráculo) con la semilla",
  "fuente": [{"tipo": "model.json", "ruta": "periodos[0]"}],
  "provisional": null,
  "entradas": {
    "datos": {"$ref": "_datos/semilla-oraculo.json"},
    "consultas": [
      {"id": "periodos[0]", "funcion": "vista_periodo", "argumentos": {"month": "2026-10"}}
    ]
  },
  "esperado": {
    "periodos[0]": {"despensa_que_cubre": 1, "devengado": {"total": 194750}}
  }
}
```
| Clave | Oblig. | Contenido |
|---|---|---|
| `formato` | sí | `"fixture/1"` |
| `motor`, `caso` | sí | Como en §1 |
| `descripcion` | sí | Una línea en español |
| `fuente` | sí | Lista de `{tipo, ruta}`: `tipo` ∈ `model.json` (ruta con la sintaxis de §3), `data.json` (`datos/semilla/data.json`), `decision` (`K21`), `criterio` (`S2-06`), `regla` (`RN-22`), `manual` (cálculo a mano explicado en `descripcion`). Toda consulta cuyo `id` sea una ruta (§2.1) exige en `fuente` esa ruta o una que la contenga |
| `provisional` | sí | `null` o `{"bloqueo": "docs/estado/bloqueos/<archivo>.md", "que": "<qué parte>"}`; un caso provisional no entra en cifras canon |
| `entradas.datos` | sí | El `Input` del motor, en línea o `{"$ref": "_datos/<archivo>.json"}` (ruta relativa a la carpeta del caso); se permite `{"$ref": …, "parche": {…}}`, que aplica un JSON Merge Patch (RFC 7396) sobre el archivo referido |
| `entradas.consultas` | sí | Lista de `{id, funcion, argumentos}`; `id` único en el caso, con la regla de §2.1; `argumentos` es un objeto con las claves que la interfaz del motor declara para esa `funcion` (`{}` si no lleva; S2: línea «Argumentos» de cada función, `S2-interfaz.md` §3) |
| `esperado` | sí | Objeto `{id: salida}`: cada salida tiene la forma de la salida de la función o un subconjunto de las claves de sus objetos (las listas van completas, §3); o `{"error": "<código>"}` |

### 2.1 El `id` de una consulta (una sola regla)
- **Ruta de `model.json`.** El `id` empieza con una de las 10 raíces de `model.json` (`parametros`, `periodos`, `despensas`, `credito`, `valor_estable`, `sensibilidad`, `papel`, `disposicion_por_despensa`, `horizonte`, `interes_por_monto`) seguida de fin, `.` o `[`, con la sintaxis de rutas de §3 (`periodos[0]`, `horizonte`, `sensibilidad[2]`). La consulta reproduce ese nodo: su `funcion` es una vista que emite las **claves literales** de `model.json` (dinero en centavos, sin sufijo) y la hoja de `model.json` de cada hoja de la salida es `id` + su ruta relativa (`periodos[0]` + `.efectivo.detalle_pagos[1].monto`). No hay otra correspondencia: ni mapas ni renombres.
- **Ruta de `data.json`.** `data.json:` + ruta con la misma sintaxis (`data.json:financiamiento.dos_quincenas`); igual que la anterior, contra `datos/semilla/data.json`.
- **Nombre libre.** Cualquier otro `id`, que debe casar con `^[a-z0-9][a-z0-9_-]*$` (sin puntos, corchetes ni dos puntos, para que nunca parezca ruta). Sus hojas se comparan solo entre esperado, oráculo y motor, aunque `fuente` cite `model.json` para trazabilidad (p. ej. un flujo cuyos valores salen de `despensas[k].pagos`).
- **Cruces.** Una consulta cuyo `id` es ruta de una rama de otro motor (por ejemplo `sensibilidad[0]` en una fixture de `s2`) usa una vista de claves literales propia del motor que la consulta (S2: `vista_sensibilidad`, `vista_interes_por_monto`, `S2-interfaz.md` §3.14); su `esperado` puede traer solo las claves que ese motor reproduce. Se cuentan en su línea (§5) y no cambian de dueño.
- Un `id` con forma de ruta que no existe en `model.json` o `data.json` es formato inválido (código 2 de `compara.py`).

## 3. Valores y hojas
- **Dinero:** entero de centavos (`97375`). Nunca decimales ni texto.
- **Fechas:** `"AAAA-MM-DD"`; instantes RFC 3339 con «Z» (solo en entradas); mes `"AAAA-MM"`; quincena `"AAAA-MM-q1"`; semana `"AAAA-Www"`.
- **Conteos y ordinales:** enteros. **Razones** (tasas): `{"num": 80050, "den": 309450}` exacto, nunca un flotante. **Porcentajes:** el entero que pida la interfaz del motor (`_pct`, `_tenths_pct`, `_bp`).
- **Claves:** las salidas genéricas usan claves en inglés (con `_cents` para dinero); las vistas que reproducen `model.json` usan sus claves literales sin sufijo.
- **Razones capturadas con decimales** (K12): `{"num": 53450, "den": 309450, "decimales": 2}`; se comparan como «HALF_UP(num/den, ×100 si la clave termina en `_pct`) a `decimales` decimales es igual al decimal guardado» (17.27 en `data.json`).
- **Hoja:** todo valor que no es objeto ni lista (número, texto, booleano o `null`). Ruta de una hoja: `id` de la consulta + `.clave` y `[índice]` (`periodos[0].efectivo.detalle_pagos[1].monto`). Las listas se comparan por índice y su largo también es hoja (`…detalle_pagos.#`); por eso `esperado` escribe cada lista **completa** (todos sus elementos, aunque cada uno sea un subconjunto de claves o `{}`).
- **De `model.json` (y `data.json`) a centavos:** se lee con `json.load(…, parse_float=decimal.Decimal)`; según la última clave de la ruta: conteos (`n`, `n_pagos`, `despensa`, `despensa_que_cubre`, `paquete`) quedan enteros; razones (`tasa`, `tasa_pct`) se comparan como «el flotante de `model.json` es igual al binary64 correctamente redondeado de num/den (×100 en `tasa_pct`)» (si no, es hallazgo con ambos valores, protocolo §6); todo lo demás es dinero: `int(valor × 100)`, que debe ser exacto (si queda fracción, es error de formato). Los booleanos y textos de `model.json` no se comparan contra él.

## 4. Salida del motor y del oráculo
El motor (`backend/cmd/motores <motor> dump --fixtures <carpeta> --out <archivo>`) y el oráculo (`oraculo/fin_run_fixtures.py` para S2; cada motor nombra el suyo) escriben el mismo formato:
```json
{"formato": "salida/1", "motor": "s2", "casos": {"s2-periodos-p1": {"periodos[0]": {"despensa_que_cubre": 1}}}}
```
- Por cada caso y consulta, la salida calculada **recortada con `esperado` como plantilla**, recursivamente:
  1. *Objeto* en `esperado`: se conservan solo sus claves, en ese orden; una clave pedida que la función no produce se escribe `"__falta__"`.
  2. *Lista* en `esperado`: se escriben **todos** los elementos calculados, en su orden; el elemento i se recorta con el elemento i de `esperado`, y un elemento calculado más allá del largo de `esperado` se escribe `"__sobra__"`. El largo escrito es siempre el calculado, así que la hoja `.#` sí detecta listas de más o de menos.
  3. *Hoja* en `esperado`: se escribe el valor calculado. Si el tipo no casa (objeto o lista donde `esperado` tiene hoja, o al revés) se escribe `"__tipo__"`.
  4. `{"error": "<código>"}` si la función falló; si `esperado` pedía error y la función no falló, la regla 1 da `{"error": "__falta__"}`.
  5. *Mapas* `{code: cents}` (por persona o por categoría): son objetos y siguen la regla 1. Qué claves trae cada uno lo fija la interfaz del motor (S2: `S2-interfaz.md` §1.1: titular y personas externas siempre, con 0; categorías solo con importe distinto de 0), así que `esperado` solo pide claves que esa regla garantiza; una que falta en una implementación sale `__falta__` y es diferencia.
- Determinista: sin fechas de generación ni rutas absolutas; dos corridas dan bytes iguales.
- `dump`: `--fixtures` y `--out` obligatorios, `--caso <nombre>` opcional; rutas relativas al directorio actual; `--out` debe quedar bajo `datos/tmp/`. Código de salida 0 si escribió el archivo (los errores de consulta van dentro), 1 si no pudo leer o escribir, 2 por uso incorrecto.

### 4.1 Piezas compartidas: dueño y registro (alcance disjunto, protocolo §3)
- **`backend/cmd/motores/main.go` y `main_test.go` los escribe el agente Go de S2** (`motor_fin`), antes que su propio motor: argumentos, lista de fixtures (`*.json` de `--fixtures` que no empiezan con `_`, por nombre), `$ref` y `parche` (RFC 7396), validación de `formato` y de `motor` (= subcomando), una llamada al evaluador por consulta, recorte (§4) y escritura determinista con los códigos de §4. Nadie más edita `main.go`; si otro motor necesita un cambio, lo pide como solicitud al integrador.
- **Cada motor agrega solo `backend/cmd/motores/<motor>.go`** (y, si quiere, `<motor>_test.go`) con `func init() { register("<motor>", <evaluador>) }`. API que fija `main.go` (paquete `main`):
```go
type Query struct {
	ID   string          `json:"id"`
	Func string          `json:"funcion"`
	Args json.RawMessage `json:"argumentos"`
}

// Evaluator calcula una consulta. input es entradas.datos ya resuelto ($ref y parche aplicados).
// Devuelve la salida (un valor que encoding/json serializa, con dinero en enteros de centavos)
// o, si la función falla, code != "" con el código de error que declara la interfaz del motor.
type Evaluator func(input json.RawMessage, q Query) (out any, code string)

// register da de alta un motor; un nombre repetido es pánico al arrancar.
func register(motor string, ev Evaluator)
```
- Una `funcion` desconocida la resuelve el evaluador con el código que declare su interfaz (S2: `invalid_argument`). Un subcomando sin motor registrado es código de salida 2.
- Si `main.go` aún no existe cuando un motor termina, su agente no lo crea ni escribe un sustituto: deja `<motor>.go` contra esta API, prueba su evaluador desde las pruebas de su propio paquete y lo anota como pendiente; el integrador compila el binario al juntar.
- **Ejecutores del oráculo:** cada motor escribe el suyo (`oraculo/fin_run_fixtures.py` en S2, `oraculo/cred_run_fixtures.py` en S3, el que nombre cada interfaz), con esta misma lectura de fixtures, recorte y salida `salida/1` en `oraculo/salidas/<motor>.json`. No hay módulo Python compartido entre motores: cada agente trabaja sin esperar a otro.

## 5. Comparación · `python3 -I oraculo/compara.py <motor>`
**Dueño: el agente Python de S2** (oráculo de `motor_fin`). `compara.py` es genérico: no tiene lógica por motor salvo la tabla de ramas de abajo y las rutas por omisión, y ningún otro agente lo edita (solicitud al integrador). Opciones: `--fixtures oraculo/fixtures/<motor>`, `--oraculo oraculo/salidas/<motor>.json`, `--motor datos/tmp/motor-<motor>.json`, `--modelo oraculo/model.json`, `--sin-motor` (compara solo oráculo y esperado mientras el motor no existe). Solo biblioteca estándar; no escribe nada.
- Para cada hoja de `esperado` compara oráculo, motor y esperado (`"__falta__"`, `"__sobra__"` o `"__tipo__"` en cualquiera es diferencia). Además, para cada consulta cuyo `id` es ruta de `model.json` (§2.1) convierte (§3) la hoja numérica de `model.json` en `id` + ruta relativa y exige que las cuatro coincidan; con `data.json:` igual contra `datos/semilla/data.json`. Los largos `.#` de esas consultas también se contrastan con la lista de `model.json` (una diferencia es `DIF`), pero no cuentan en las líneas de `model.json`.
- Ramas de `model.json` por motor (tabla fija en `compara.py`): **s2** = `periodos`, `valor_estable`, `horizonte`, `papel` (159 hojas numéricas); **s3** = `parametros`, `despensas`, `credito`, `sensibilidad`, `disposicion_por_despensa`, `interes_por_monto` (437); **s4** y **s5** = ninguna (sus líneas de `model.json` dicen `0 de 0`). Total 596. Las hojas de la rama de otro motor que una fixture reproduce (cruces, §2.1) van en su propia línea.
- Salida (en este orden):
```
compara s2 · 27 fixtures
oráculo contra esperado: N de N hojas coinciden
motor contra esperado: N de N hojas coinciden
oráculo contra motor: N de N hojas coinciden
model.json (ramas de s2): 159 de 159 hojas coinciden
model.json (cruces con otras ramas): 43 de 43 hojas coinciden
data.json: 0 de 0 hojas coinciden
diferencias: 0
```
  seguida de una línea por diferencia, `DIF <caso> <ruta> esperado=<v> oraculo=<v> motor=<v>`, y una por hoja de las ramas del motor que ningún caso cubre, `SIN-DUEÑO <ruta>`.
- Código de salida: 0 si todo coincide y no hay hojas sin dueño; 1 si hay diferencias u hojas sin dueño; 2 por archivo faltante o formato inválido (con el motivo).
- Las fixtures marcadas `provisional` se comparan igual, pero sus hojas se informan en una línea propia (`provisionales: N de N`) y no cuentan para la cobertura de `model.json`.

## 6. Motores que no nacieron con este formato (S3) y los que lleguen
- **Todos los motores adoptan este formato; no hay excepciones.** S3 escribió su contrato antes de que existiera (su `S3-interfaz.md` §10 lo declara SUPUESTO y dice que «se alinea cuando exista»), así que se alinea: fixtures en `oraculo/fixtures/s3/<caso>.json` con `fixture/1`; subcomando `s3 dump` en `backend/cmd/motores/s3.go` (no en `backend/cmd/finanzas`); salida `salida/1`; ejecutor `oraculo/cred_run_fixtures.py` → `oraculo/salidas/s3.json`; comparación oficial `python3 -I oraculo/compara.py s3`.
- Lo propio de S3 que no choca se conserva como verificación adicional, nunca como sustituto: `oraculo/cred_compara.py` (sus 437 hojas contra `model.json`) y la salida `--canon`.
- `{"$model": "<puntero>"}` no es parte de `fixture/1`: en una consulta con `id` ruta de `model.json`, `esperado` escribe el valor convertido y `compara.py` lo contrasta con `model.json` en cada corrida, así que una cifra mal copiada sale como `DIF`. R-05 y R-06 van con `id` `data.json:financiamiento.dos_quincenas` (§2.1) y la razón con `decimales` (§3).
- El cambio en `S3-interfaz.md` §2 y §9–§11 lo hace quien corrija ese contrato (solicitud del reporte `docs/estado/wf1/correccion-contrato-S2-r1.md`).
