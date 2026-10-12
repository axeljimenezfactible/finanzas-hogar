# Reporte · revisor-datos-r1 (Pista 3, fase 1, ronda 1 de 3) · 11-oct-2026, noche

## Hecho
- Revisé `docs/design/maquetas/datos-sinteticos.json`, `LEEME.md`, `docs/design/tools/ds/verificar_sinteticos.py` y `docs/estado/wfp3/datos-sinteticos.md` contra P3-12, P3-08, K45.1, U1 §2.1, §2.3, §3.3, §3.6, §3.7 y §3.28, U3 §2 y §4, y U4 P-26.09.
- Informe: `docs/estado/wfp3/REVISION-datos-r1.md`. **No aprobado**: 0 altos, 4 medios (M-1 a M-4) y 4 bajos (B-1 a B-4).
- No corregí nada. No toqué nada fuera de `docs/estado/wfp3/` y `datos/tmp/rev-datos-r1/` (este último ya borrado).

## Verificación (comandos y salidas)
- `python3 -I docs/design/tools/ds/verificar_sinteticos.py` → 39 OK, `RESULTADO: OK`, salida 0.
- `grep -rnoE "<patrón de P3-12>" docs/design/maquetas/datos-sinteticos.json | wc -l` → 0; el mismo grep sobre `docs/design/maquetas/` → 0. Credenciales (patrón de P3-12) → 0; `@` → 0; IP, dominios y `http` → 0.
- `copia.py` (propio): 19 cifras canon más 8 del reporte y de U3, en centavos, flotante y texto con y sin coma → 0 de dinero. Nombres de productos iguales a los de `data.json` → 0; 4 de 105 precios iguales a un número de `data.json`, todos de producto distinto y redondos. 49 `*_cents` iguales a un número de model/data, todos redondos (B-4).
- `sumas.py` (propio): 16 de 16 OK (plan, real desde gastos, % y estados, ingreso y sin asignar, partes, flujo, próximo pago, ahorro, canasta, lista, recetas y plan, IMC, deuda viva, Kueski, 539 `_fmt`, separador de miles).
- Comando de fechas de P3-12 → 15 fechas, todas en día 15, 30 o 28-feb. `fechas.py` (propio): 0 en día 9 o 10 desde el 1-oct, 0 días de la semana mal puestos, 0 fechas coincidentes con monto canon o de Kueski. La única coincidencia, $1,000.00 de ahorro planificado, es otro concepto.
- `claves.py` (propio): llaves ausentes `confianza`, `renglon`, `notificaciones`, `semilla` (M-1 y M-2). `openapi.yaml` contra `inicio.tarjetas` → M-3. Rango de `calendario_hogar` 2026-10-15 a 2027-04-17 → M-4.

## Supuestos
1. «Cifra del reporte o de Kueski» incluye, además del grep de P3-12, las cifras reales que U3 cita ($6,905.50, $10,000.00, $3,233.00, $138.50 y $1,197.64) y las de plan §1.2 ($5,000.00, $4,000.00 y $199.00). No apareció ninguna como dinero.
2. Una coincidencia con un número redondo ($1,000.00, $2,000.00, $3,000.00, $9,000.00) o con un precio chico de otro producto no cuenta como copia. La de $9,000.00 va como baja (B-4) porque es el mismo concepto.
3. Los horarios 21:30, dom 19:00, día 1 08:00 y avisos 09:00 son FIRMES por K50; no son una decisión inventada.
4. La escritura fuera de las rutas permitidas y el exceso de tamaño del autor quedan como observaciones, fuera del veredicto, porque no forman parte de los criterios que me tocan.

## Bloqueos
- Ninguno.

## Pendientes
- Que el autor (o el agente de corrección) resuelva M-1 a M-4 y, si quiere, B-1 a B-4. Después, ronda 2.

## Consumo aproximado
- Acciones: unas 30. Leído: unos 130 KB (LEEME, reporte, P3-12 y P3-08, U1 §2.1, §2.3 y §3, U3 §2 y §4, U4 P-26.09 por grep, y extractos del JSON por script). Escrito: unos 20 KB (informe y reporte; los scripts temporales, unos 12 KB, ya borrados).
