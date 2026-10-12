# Corrección del contrato S3 · ronda 4 (agente correccion-contrato-S3-r4)

11-oct-2026 · WF-1 · Artefacto: `docs/contratos/motores/S3-interfaz.md` v1.3 → **v1.4** (916 → 923 líneas; 129,503 → 133,628 bytes) · Informe atendido: `docs/estado/wf1/REVISION-contrato-S3-r4.md` (1 hallazgo medio, P1; 0 altos).

## Hecho
- **P1 (media; S3-01, S3-02; residuo de r3-N4; FORMATO 1.2 §1): cerrado.**
  - §10.0, bullet «Núcleo S3»: además de `seedFromModel` (el `Input`), el núcleo arma con un **segundo ayudante de prueba `sensitivityFromModel`** las entradas de `sensibilidad` e `interes_por_monto`, que no están en el `Input`; con los dos corre `TestCanon85`, `TestCanonHojas437` y las demás pruebas con cifras canon (incluidas `TestTasaPorMonto` y `TestTasaSensibilidad`), **sin leer archivos de otro autor**. La igualdad con fixtures ajenas solo se exige si existen.
  - §10.0, receta nueva «Insumo de sensibilidad y de cotización por monto» (entre «Datos compartidos» y la tabla de casos), la que sugirió el revisor, con cada origen: `amounts_cents` = `sensibilidad[*].N`; `rate` = tasa de 4 q de la semilla (`despensas[0].intereses` / `despensas[0].capital`, `decimales` 4 de K12, igual que la semilla); `term_count` = plazo de d1 (número de `despensas[0].pagos`); `external` = karen con `principal_cents` = `parametros.karen_fija.capital` (Joint y PaperIn en null: `Sensitivity` solo usa el capital, §7); `holder_current_principal_cents` = `despensas[0].axel.capital`; `accrual_months` = 2 (de §7); `cash_monthly_cents` = HALF_UP(`parametros.papel` / 2); montos de `vista_interes_por_monto` = `interes_por_monto[*].monto`, con la misma tasa y el mismo plazo. Las cifras entre paréntesis son solo la comprobación y no van en las pruebas.
  - Tabla de casos de §10.0: `s3-sensibilidad` y `s3-interes-por-monto` dicen que sus argumentos son el insumo de `sensitivityFromModel`, así el agente de regresión usa la misma receta y la igualdad queda bien definida.
  - §11: `TestCanon85` (las filas S3 salen con `seedFromModel` y `sensitivityFromModel`) y `TestCanonHojas437` (4 ramas sobre el `Input`; `sensibilidad` e `interes_por_monto` sobre el insumo; igualdad con `_datos/kueski-semilla.json`, `s3-sensibilidad` y `s3-interes-por-monto` solo si existen).
  - §11.1, l. 797: reemplazada por la frase del revisor: «con el `Input` de `seedFromModel` y el insumo de sensibilidad de §10.0 (`sensitivityFromModel`), sin leer archivos de otro autor; si existen `_datos/kueski-semilla.json` y `s3-sensibilidad`, exige además igualdad, como `TestCanonHojas437`».
  - Encabezado (versión 1.4 y autor de r4), títulos de §10, §10.0 y §11.1 con «r4-P1», y renglón «Cambios de la versión 1.4» en §14.
- No toqué los hallazgos bajos (Q1–Q4, B1–B9) ni otras secciones: el encargo pedía solo altos y medios.

## Archivos
- Modificado: `/home/deepsec/finanzas-hogar/docs/contratos/motores/S3-interfaz.md` (encabezado, §10, §10.0, §11, §11.1, §14).
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-contrato-S3-r4.md` (este reporte).
- Temporales en `datos/tmp/corr-s3-r4/` (`verifica.py`, `edita.py`, copia de la v1.3, `s3.cov`): borrados al terminar.

## Verificación
- **Receta de P1 contra el canon** (`python3 -I datos/tmp/corr-s3-r4/verifica.py`, `fractions.Fraction`, solo lee `model.json`): el insumo sale `amounts` [300000, 350000, 400000, 450000, 500000], `rate` 80050/309450, `n` 4, externa karen 119150, titular actual 190300, `accrual_months` 2, `cash_monthly` 9950; y para `interes_por_monto` [100, 10000, 100000, 200000, 309450, 400000, 500000]. Es igual a los argumentos de las filas `s3-sensibilidad` y `s3-interes-por-monto` de §10.0. Con las reglas de §7: **73 de 73 hojas coinciden** (45 de `sensibilidad` y 28 de `interes_por_monto`) y `axel_menor_que_actual` coincide en las 5 filas. **R-37…R-43 con las aserciones de §11.1 sobre la fila 400000: 7 de 7** (261687, 176751, 84936, 130844, 88376, 42468, 251737).
- **Integridad del contrato** (script de contraste v1.3 contra v1.4): cambiaron solo el encabezado, §10, §11, §11.1 y §14; §1–§9, §12 (las 437 hojas), §13 y §14.1 quedaron iguales byte a byte. Ningún número de la v1.3 desaparece salvo el «1.3» del encabezado (pasa a 1.4). En todo el archivo, 0 filas de tabla tienen un número de columnas distinto del de su tabla. `sensitivityFromModel` aparece en 10 lugares que se citan entre sí (§10.0 bullet, receta, tabla de casos, §11, §11.1, §14).
- `sha256sum oraculo/model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (intacto, S3-16).
- Código existente (no es mío; no lo edité): `cd backend && go build ./internal/motor_cred/... && go vet ./internal/motor_cred/...` → limpio; `gofmt -l internal/motor_cred` → vacío; `go test -count=1 -coverprofile=… ./internal/motor_cred/...` → `coverage: 0.0% of statements`, porque todavía no hay `*_test.go`, `comparador/` ni `evaluar/`. Sin pruebas Go, no se pueden correr todavía S3-01, S3-02, S3-03…S3-14, S3-18 ni S3-19 sobre código.
- Oráculo (no es mío): `oraculo/fixtures/s3/` sin fixtures y `_datos/` vacío. `python3 -I oraculo/compara.py s3 --sin-motor` → «falta oráculo: …/oraculo/salidas/s3.json», código 2. `python3 -I oraculo/cred_run_fixtures.py` → «0 fixtures · 0 consultas»; creó `oraculo/salidas/s3.json` con `casos` vacío, y lo borré porque no existía antes de mi corrida. `python3 -I oraculo/cred_escenarios.py` → «escenarios aceptados con violación = 0», código 0, aunque faltan las fixtures del comparador (ver Solicitud 1).
- Criterios de la fase que toca P1: **S3-01** (las entradas de `TestCanon85` ya no dependen de archivos de otro autor; R-37…R-43 y R-74…R-85 se arman desde `model.json`: 7 de 7 y 73 de 73 hojas) y **S3-02** (las 6 ramas de `TestCanonHojas437` tienen origen definido: 364 hojas desde `seedFromModel` y 73 desde `sensitivityFromModel`, 437 en total, y las 73 reproducidas). S3-20 sigue igual: ningún bloqueo nuevo, porque todo origen sale de hojas de `model.json` o de constantes que ya estaban en §7 y §10.0.

## Supuestos
- `external.joint_cents` y `external.paper_in_cents` van en null en el insumo de sensibilidad. Son composición informativa (§4, `ShareInput`) y `Sensitivity` no los usa (§7), así que no cambian ninguna cifra, esquema ni endpoint. Lo fijé solo para que la igualdad JSON con `s3-sensibilidad` quede bien definida.
- El «/ 2» de `cash_monthly_cents` es el del valor estable («199/2», R-44, `model.py` `PAP_M`), una constante de la interfaz como `accrual_months` = 2. No es hoja ni cifra canon escrita.

## Bloqueos
- Ninguno nuevo. Siguen abiertos los que el contrato ya cita: `S3-formula-tarjeta`, `S3-convencion-prestamo`, `S3-objetivos-medida`, `S3-origen-cargo-normal`, `S3-serie-sin-intereses` y `U4-iva-credito`.

## Pendientes
- Hallazgos bajos de la revisión r4 que no me tocaban (siguen abiertos): Q1 (disparador «algún interés por pago < 0» en §5.1, §14.1 y `S3-serie-sin-intereses.md`), Q2 (marca sobre el tope solo en el caso 3, §10.1), Q3 (si los cortes de la tarjeta mueven el disponible, `TimelineInput.State` en §5.2), Q4 (nota en `docs/estado/wf1/contrato-S3.md` y frase de `S3-formula-tarjeta.md`), y B1–B9 de la ronda 3.
- Verificación sobre código (pruebas Go con cobertura ≥ 85 %, `TestCanon85` 85 de 85, `compara.py s3` 437 de 437) cuando existan las pruebas del núcleo y las fixtures del agente de regresión.

## Solicitudes fuera de alcance
1. `oraculo/cred_escenarios.py` (agente de regresión S3): sin argumentos y sin las fixtures `s3-comparador-*`, escribe «escenarios aceptados con violación = 0» y termina con código 0. §10.0 pide que escriba «falta oraculo/fixtures/s3/<archivo> (autor: optimizador S3)» por cada caso que falte y termine con código 1. Motivo: con la salida actual parece que S3-12 pasa sin haber evaluado nada.
2. Núcleo S3 (`backend/internal/motor_cred/*_test.go`, todavía sin escribir): implementar `seedFromModel` y `sensitivityFromModel` con las recetas de §10.0 v1.4; `TestCanonHojas437` debe comparar con `_datos/kueski-semilla.json`, `s3-sensibilidad` y `s3-interes-por-monto` solo si existen.
3. Agente de regresión S3 (`oraculo/cred_genera_fixtures.py`, todavía no existe): escribir los argumentos de `s3-sensibilidad` y `s3-interes-por-monto` con la receta «Insumo de sensibilidad y de cotización por monto» de §10.0, incluidos `decimales` 4 en la tasa y `joint_cents` y `paper_in_cents` en null.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna nueva. Sigue en §14 la propuesta K31 de `optimizer.date_window_days` = 7 días, sin valor sembrado; se decide en la puerta 2.

## Consumo aproximado
- Acciones de herramienta: 29.
- Leído: ~140 KB (protocolo 5.6; revisión r4 15.5; decisiones §C y §K 17; criterios S3 completos 35; S3-interfaz por partes ~55: encabezado, §7, §9.3, §10, §11, §11.1, §12 y §14; reporte r3 §Verificación 3; salidas de scripts y `grep` ~9).
- Escrito: ~19 KB redactados (texto nuevo del contrato ~4.1 KB; scripts temporales ~8 KB; este reporte ~7 KB). El script reescribió el archivo completo del contrato (133.6 KB), pero solo cambió lo anterior.
