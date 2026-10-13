# revisor-backend-cierre-r2 · WF-3c bis · 12-oct-2026

Agente revisor independiente (Opus). HEAD `0197dad`, con el diff sin commit de `backend-cierre` y `correccion-backend-cierre-r1`. No corregí nada, no toqué git y no reinicié el servidor de desarrollo. El informe completo está en `docs/estado/wf3/REVISION-backend-cierre-r2.md`.

## Hecho
- **HECHO:** Reejecuté (a) a (d) con mis propios comandos. Las pruebas, `verify-wf3-go`, `verify-wf3-contrato`, `build.py --check` y `verify-all` salen OK. Las cifras de K-3, P-17, G-5 y el barómetro coinciden con `getLiveDebt` y con el oráculo, también con una línea archivada que tiene pagos programados.
- **HECHO:** Los medios M-1, M-2 y M-3 y los bajos B-1 a B-4 de la ronda 1 quedan cerrados.
- **HECHO:** Encontré un alto nuevo, A-1. Al quitar el filtro de activas en `states`, el comparador de escenarios ofrece y recomienda líneas archivadas. Lo confirmé con una prueba overlay y con un servidor privado sobre una copia de la base de desarrollo. Con el filtro anterior, el comparador no las ofrece.
- **Dictamen: NO APROBADO** (1 alto, 0 medios y 0 bajos).

## Archivos
- `docs/estado/wf3/REVISION-backend-cierre-r2.md` (informe)
- `docs/estado/wf3/revisor-backend-cierre-r2.md` (este reporte)

## Verificación
- **S6-03:** `make verify-wf3-contrato` EXIT 0, con 221 operaciones en 176 rutas y 0 diferencias. `build.py --check` EXIT 0: 22 archivos idénticos.
- **S6-04:** sin cambios en el manejo de errores ni en el formato. Dinero en centavos enteros.
- **S6-01:** `make verify-wf3-go` EXIT 0, con 33 paquetes OK y 5 sin pruebas.
- **Cobertura:** api/creditos 82.5 %, api/presupuesto 77.7 %, api/tablero 87.8 % e internal/api 88.8 %. 0 SKIP.
- **S2/S3:** en el servidor de desarrollo, 389500 = `getLiveDebt` hoy. G-5 da 292125 en octubre, con 12 de 12 cierres iguales a `getLiveDebt`. `month_pending_household_cents` es 59882, `scale_max_pct` es 60 e `incomes/month?period=2026-11` da 2450000, 0 y 0.
- **Línea archivada (prueba propia):** K-3 da 779000 = `getLiveDebt`, los pagos del mes coinciden con el barómetro (260658 y 254632, según la prueba) y G-5 coincide en los 12 cierres. **Comparador:** recomienda la línea archivada (A-1).
- **C1-15:** OK, `S6-deuda-viva-hoy` está citado en C1. **C1-13:** `verify_c1.py` OK; el cambio de `docs/criterios/C1.md` está pedido.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT 0.

## Supuestos
- Califico A-1 como alto por dos razones: cambia la salida del comparador S3 (escenarios evaluados, puntajes y recomendación) con los mismos datos, y contradice la definición de `credit.is_active` en `diccionario-datos.md`, que es uno de los documentos que mandan.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Hay que corregir A-1 y hacer una ronda 3 de revisión.
- El servidor de `:8443` corre el binario de antes de la corrección. Hay que reiniciarlo cuando lo decida el orquestador.

## Propuestas para Axel
- Ninguna nueva. La del punto 5 de `S6-tablero-lineas-archivadas` sigue en pie.

## Paquetes instalados
Ninguno.

## Consumo aproximado
42 acciones, unos 110 KB leídos y unos 14 KB escritos. Todos los temporales de `datos/tmp/rev-bc2/` están borrados.
