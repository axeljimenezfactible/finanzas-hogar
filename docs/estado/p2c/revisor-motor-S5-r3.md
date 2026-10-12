# P2c · revisor-motor-S5-r3 · 11-oct-2026, noche

Agente `revisor-motor-S5-r3`, revisor independiente, ronda 3 de 3 del flujo P2c «Aplicación de la puerta 2». El informe completo está en `docs/estado/p2c/REVISION-motor-S5-r3.md`.

## Hecho
- **Veredicto: NO APROBADO (0 altos, 1 medio, 6 bajos).**
- M-1 de la ronda 2 (`created_at` con zona) quedó **CERRADO** en lo que pedía:
  - `-06:00` da 2100 `ticket_manual` y `+00:00` da 2000 `manual`, igual en el motor y en el oráculo.
  - Mis 7 sondas válidas coinciden en los dos lados y con mi cálculo.
  - Las fixtures nuevas detectan el oráculo y el motor de `HEAD`.
- **M-1 nuevo [media]:** quedan 6 textos de `created_at` en los que un lado acepta y el otro da `invalid_date`. El bloqueo `S5-precio-vigente-captura` dice que no hay divergencia.
  - El motor acepta la hora de un dígito (`T7:00:00-06:00` decide S-8 con 2100) y la coma decimal.
  - El oráculo acepta un salto de línea final y dígitos no ASCII.
  - El año 0000 es RFC 3339 válido: el motor lo acepta y el oráculo no.
- Bajas:
  - B-1 a B-5 siguen de rondas anteriores, y la B-3 creció con las 4 fixtures nuevas fuera de la tabla de §15.
  - B-6 es nueva: `types.go:4` cita el contrato 1.4, y el vigente es 1.5.
- No escribí código, fixtures ni bloqueos. No hice `git add` ni `commit`. No toqué `docs/design/`, `docs/estado/wfp3/` ni `data.json`.

## Verificación (comandos propios desde `/home/deepsec/finanzas-hogar`)
| Qué | Salida |
|---|---|
| `gofmt -l internal/precios cmd/motores`; `go vet ./internal/precios/... ./cmd/motores` | vacío; limpio |
| `go test -count=1 -coverprofile=… ./internal/precios/...` | 96.0 %, 100.0 %, 99.3 % y 95.3 %; total 88.8 % |
| `pre_genera_fixtures.py` en una raíz aislada y `diff -r` | 36 fixtures, idénticas |
| `pre_run_fixtures.py --out datos/tmp/rev-s5-r3/oraculo-s5.json` y `cmp` | 36/36, 1359 hojas, 0 diferencias; idéntico a `salidas/s5.json` |
| `motores s5 dump` (mi binario) y `compara.py s5 --motor …` | 958/958 en las tres comparaciones; provisionales 315/315; diferencias 0 |
| Mis 26 fixtures `rv3-*` (granel, diferencia, PROFECO de 17 registros y 22 sondas de `created_at`) | 20 iguales. Las 6 distintas son todas de `created_at` (M-1) |
| Granel (0.5, 1.01 y 1.03) | 4332 contra 4336 (+4), igual en los dos lados y en mi cálculo |
| `profeco_import`, 17 registros | read 17 = 7 + 1 + 2 + 6 + 0 + 1 (producción) y 7 + 1 + 1 + 6 + 1 + 1 (revisión); iguales |
| Canon | +13850 y +10050; 61/60/55; 22/3; `ties` 8 |
| `git diff --quiet HEAD -- datos/semilla/data.json` | sin cambios (sha256 `2fe7afae…1031d31b`) |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | EXIT=0, «WF-1: verificación completa OK», «Verificación ampliada OK»; los archivos s5 quedaron con el mismo sha256 |

## Supuestos
- S-a: llamo RFC 3339 a la gramática de §5.6, con `DIGIT` = %x30-39 (RFC 5234). Así lo leyó también la corrección r2 para `+24:00` y `+05:60`.
- S-b: clasifico M-1 como medio porque es la misma clase y el mismo efecto que M-1 de r2, y porque el bloqueo afirma lo contrario. No mueve ninguna cifra canon.

## Bloqueos
- No abrí ninguno. El año 0000 debe sumarse a la pregunta 3 de `S5-precio-vigente-captura`, y eso le toca a quien corrija M-1.

## Pendientes
- M-1: las gramáticas estrictas en ASCII en los dos lados, cuatro fixtures `s5-ingesta-*` y corregir el texto del bloqueo.
- B-1 a B-6: van al integrador, a contratos o al motor.

## Propuestas para Axel
- Ninguna que decidir. M-1 es técnico: §15.1 ya dice RFC 3339.

## Solicitudes fuera de alcance
- Al integrador: esta es la ronda 3 de 3 y queda 1 medio. Decide si abre una corrección r3 acotada a M-1 (es poco código en cada lado) o si registra M-1 como pendiente con su bloqueo.

## Consumo aproximado
- 33 acciones de herramienta.
- Unos 175 KB leídos.
- Unos 21 KB escritos: los dos informes y mis scripts y sondas temporales en `datos/tmp/rev-s5-r3/`, ya borrados.
