# contrato-S3 · reporte final (WF-1, 11-oct-2026)

## Hecho
- Contrato de interfaz del motor de créditos `docs/contratos/motores/S3-interfaz.md` (46.6 KB, sin código): principios (centavos, `big.Rat`, dos modos de redondeo de K14, pago de referencia desde el total exacto), paquetes `motor_cred` y `motor_cred/comparador`, tipos y firmas, los 4 tipos de crédito (línea tipo Kueski con tasa cotizada, tabla por plazo K12, promociones RN-10, partes por persona, composición; disponible A/B con reposición K15; tarjeta; préstamo francés; tienda con abonos), calendario unificado (S2 y S7), deuda y compras vivas (RN-21), cotización por monto y sensibilidad, comparador (restricciones K16, 4 objetivos K17, mín-máx con empate declarado = 1, pesos iguales, desempates, búsqueda exhaustiva, tiempo máximo 5 s configurable K43/K44 → 503), frontera JSON con S2 (§2.5 de S2-interfaz), CLI `finanzas s3 dump` (`--fixture`, `--canon`), catálogo de 19 fixtures con formato `{caso, fuente, provisional, entradas, esperado}` y `$model` (JSON Pointer) para el canon, oráculo `cred_*.py` independiente, pruebas con los nombres de `-run` de los criterios, tabla de las 437 hojas y tabla criterio → sección.
- Tres bloqueos nuevos con propuesta: fórmula de tarjeta y tienda (+ calendario de la tarjeta), convención de tasa del préstamo, medida de los objetivos 2–4 y ventana del colchón.

## Archivos
- `docs/contratos/motores/S3-interfaz.md` (nuevo)
- `docs/estado/bloqueos/S3-formula-tarjeta.md`, `docs/estado/bloqueos/S3-convencion-prestamo.md`, `docs/estado/bloqueos/S3-objetivos-medida.md` (nuevos; el protocolo §2 pide escribirlos)
- `docs/estado/wf1/contrato-S3.md` (este reporte)

## Verificación
- Cobertura de hojas (script temporal en `datos/tmp/`, ya borrado, que recorre `model.json` y casa cada hoja numérica con los patrones de §12): `renglones 42 · por rama parametros 9, despensas 144, credito 143, sensibilidad 45, disposicion_por_despensa 68, interes_por_monto 28 · total 437 · sin dueño 0 · dos dueños 0 · conteo declarado distinto []`. La primera corrida encontró 56 hojas sin dueño por un patrón sin `[*]` en las series A/B; se corrigió el patrón (no el conteo).
- Fórmulas de la interfaz contra `model.json` (script temporal con `fractions.Fraction`, borrado): series RN-07 (montos truncados, capital HALF_UP, partes externas HALF_UP y titular con el resto), intereses por persona, intereses y totales de las 4 despensas, sensibilidad (45 hojas) e `interes_por_monto` (28): **0 diferencias**; `float(80050/309450) == parametros.tasa` y `float(× 100) == tasa_pct`: verdadero. Hallazgo: el pago de referencia calculado con el total redondeado da 0.32 (contra 0.31) y 1,258.69 (contra 1,258.68): la interfaz exige el total exacto.
- `sha256sum oraculo/model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` antes y después (intacto).
- Tamaño: `wc -c S3-interfaz.md` = 46,597 bytes (≤ 50 KB).
- Criterios S3-01…S3-20 → sección: tabla §13 de la interfaz (todos con sección; S3-20 en §14 y este reporte).

## Supuestos
- `oraculo/fixtures/FORMATO.md` no existía al escribir: se usó `{caso, fuente, entradas, esperado}` (+ `provisional`) con `$model` para el canon; se alinea cuando contrato-S2 lo publique.
- La CLI de S2 aún no estaba definida en `S2-interfaz.md`: se fijó `go run ./cmd/finanzas s3 dump --fixture|--canon [--out]`, salida `{caso, salida}`, códigos 0/1/2; dinero en centavos y persona por `code` como S2 §1.1.
- Paquete `motor_cred` (igual al directorio).
- Promoción sin `store_id` en ambos lados: comercio comparado sin espacios a los lados y sin distinguir mayúsculas (como `COLLATE NOCASE`).
- Compra viva en el intervalo [fecha, último pago] inclusivo (lectura literal de RN-11 y de su ejemplo del 16-dic-2026).
- Empate en la normalización (máx = mín): el objetivo vale 1 en todos (S3-14 delega la declaración al agente; no cambia el orden).
- `violated_constraint` = la primera violada en el orden del enumerado de C1; fechas anteriores a `Today` no se generan; la desviación estándar se compara en milésimas de centavo con raíz entera exacta.
- Pagos con vencimiento ≥ `Today` sin registro se proyectan completos en su vencimiento; vencidos sin completar no reponen (K15).
- Línea revolvente semanal o mensual sin primera fecha capturada: `ErrUnsupported` (RN-06 solo define la quincenal).
- La parte de la co-titular la suma la API al titular antes de llamar (C-09).

## Bloqueos
- `docs/estado/bloqueos/S3-formula-tarjeta.md` (intereses y mínimo de tarjeta y tienda; calendario de la tarjeta)
- `docs/estado/bloqueos/S3-convencion-prestamo.md` (tasa del periodo e IVA del préstamo)
- `docs/estado/bloqueos/S3-objetivos-medida.md` (canasta y tienda del objetivo 4, serie del objetivo 2, ventana del colchón, créditos sin límite)
Ninguna cifra canon depende de ellos.

## Pendientes
- Implementación (fuera de este contrato): `backend/internal/motor_cred/*.go` (núcleo Opus), `backend/internal/motor_cred/comparador/*.go` (optimizador Opus), `oraculo/cred_lib.py`, `cred_run_fixtures.py`, `cred_compara.py`, `cred_escenarios.py`, `oraculo/fixtures/s3-*.json` (regresión Sonnet), subcomando `s3 dump` en `backend/cmd/finanzas`, `docs/estado/S3-escenarios.md` con los 3 casos fijados antes de correr y los tiempos medidos.

## Solicitudes fuera de alcance
1. `docs/contratos/openapi.yaml` (C1) `Scenario`: agregar `store_id`, `fortnight_payment_stddev_cents`, `over_cap` (personas cuyo pico pasa su tope) y `provisional`. Motivo: K17 mide la desviación estándar (hoy solo hay máx − mín), S3-12 pide marcar el tope y el objetivo 4 compara tiendas (bloqueo S3-objetivos-medida).
2. `docs/contratos/reglas-negocio.md` RN-11: escribir el empate de la normalización (vale 1), los desempates (puntaje y llave canónica) y el intervalo inclusivo de compra viva. Motivo: S3-14 pide la fórmula con el empate en RN. RN-26 y RN-06: la fórmula y el calendario de la tarjeta cuando se resuelva S3-formula-tarjeta.
3. Integrador WF-1: unificar `Date`, `Cents`, nombre de paquete, banderas de `dump` y `FORMATO.md` entre `motor_fin` y `motor_cred`.
4. `oraculo/verify3.py` línea 89 escribe en `/home/deepsec/reporte-compra-latex/model/` (S3-16 pide corregirlo en WF-1): asignarlo al agente de regresión u oráculo.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- `optimizer.date_window_days` (sin valor sembrado; 00003 dice «lo propone S3 con Axel»; puerta 2): **7 días**. Cubre una semana completa (cualquier día de compra o de mercado) y basta para que el primer pago de Kueski pase a la quincena anterior o siguiente (cada quincena dura unos 15 días); el peor caso estimado son unos 7,600 escenarios de ≤ 12 pagos, muy por debajo de los 5 s de `comparator.timeout_seconds`. Mientras Axel no lo capture, el comparador responde «faltan restricciones» (C1).
- Respuestas propuestas en los bloqueos: S3-objetivos-medida (1)(b) monto fijo y tienda como dimensión, (2) totales por quincena del hogar con el tope como marca, (3) vida de la compra del escenario, (4) solo créditos con límite; S3-formula-tarjeta (a) un estado de cuenta real; S3-convencion-prestamo (a) tasa nominal / 12 con IVA en la tasa.

## Consumo aproximado
45 acciones de herramienta; ~185 KB leídos; ~66 KB escritos.
