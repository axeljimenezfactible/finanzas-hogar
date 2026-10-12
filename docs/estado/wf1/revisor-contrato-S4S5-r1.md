# revisor-contrato-S4S5-r1 · reporte final (WF-1, 11-oct-2026)

## Hecho
- Revisé de forma independiente `docs/contratos/motores/S4-interfaz.md`, `S5-interfaz.md` y `productos-formato.md` contra `docs/criterios/S4.md` 1.2, `docs/criterios/S5.md` 1.2 y decisiones §D, §E y §K. No corregí nada.
- Reejecuté con mis propios scripts (`python3 -I`, `Decimal` y `Fraction`) el canon de la comparativa (+138.50 y +100.50, K-01 a K-08, 22/8/3), la unidad base de los sustitutos, la frescura de K40, Mifflin-St Jeor e IMC, RN-15, RN-16, el ejemplo de costo §8.3 y la cobertura del catálogo (61 piezas, 39 renglones, 3,094.50). Todo coincide.
- Cotejé con C1 (extractos por grep), el esquema 1.3 (grep), RN-12 a RN-20, `FORMATO.md` y los 3 bloqueos nuevos.
- Informe: `docs/estado/wf1/REVISION-contrato-S4S5-r1.md`. Veredicto **NO APROBADO**: 0 altos, 7 medios (M1 a M7) y 8 bajos.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-contrato-S4S5-r1.md` (creado)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-contrato-S4S5-r1.md` (este reporte)

## Verificación
- `python3 -I verif.py`: K-01 309450 = 309450; K-02 323300; K-03 13850 = 10050 + 5900 − 2100; K-06 10050; piezas 61/60/55; renglones 33/3; 0 renglones con unitario × cantidad ≠ total; ganadores 22/8/3; sustitutos 2933/2800, 12692/16406, 3375/2833; §8.3 4425/2213/8850; frescura 1, 60/61, 59/61, 30/61, 0 y «no compite»; `[0.9836, 0.4918, 0.0]`; MSJ 2099, 1596, 1799, 1846, 2556, 2056; IMC 185, 250, 300, 237 y 185 (18.4975); RN-15 1486571; membresía 8 y 83100.
- `python3 -I cat.py`: 36 productos, 39 presentaciones, 0 nombres de `data.json` faltantes o de más, 39 renglones de §5.2 iguales a `cuentas` en orden, texto, piezas e importe, 61 piezas, 3,094.50, 0 errores de precio unitario, 0 diferencias de piezas por presentación.
- RN-16: 2026-10-10 es ISO 6 y el siguiente martes es 2026-10-13. `tesseract 5.5.3`, `osd spa`; `unrar`, `bsdtar` y `sqlite3` presentes; `datos/semilla/tickets/` no existe.
- grep: S4-01…S4-15 y S5-01…S5-13 citados; PROVISIONAL con bloqueo (dos excepciones justificadas); bloques JSON válidos (2); columnas, enumerados y claves de `app_setting` existen en 00001 y 00003; campos de C1 existen salvo `identical_pieces`, `target_status`, `missing_inputs` y `missing_nutrition`, que son salidas propias del motor.
- Criterios: la tabla completa con estado y evidencia está en el informe. Parciales: S4-08, S4-09, S4-11, S4-13 y S5-06.

## Supuestos
- Evalúo el contrato por suficiencia: un criterio está «cubierto» si el contrato fija todo lo que se necesita para cumplirlo. Las pruebas Go aún no existen.
- M1 se basa en `data.json: supuestos` («Precios del carrito…») y en que `cuentas` y el reparto solo usan la compra de Bodega. Si existe un ticket de Walmart, M1 se reduce a la incompatibilidad con C1 (una sola `expense_id`).

## Bloqueos
Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`; no toca S4 ni S5.

## Pendientes
- Para el autor (ronda 2): M1 a M7 y B1 a B8 en `docs/contratos/motores/S4-interfaz.md`, `S5-interfaz.md`, `productos-formato.md` y `docs/estado/wf1/contrato-S4S5.md`.

## Solicitudes fuera de alcance
1. C1 `getBasketComparison`: soportar «compra contra carrito» o dos compras (M1), o confirmar el modelo de una sola compra.
2. C1 `BasketComparison.stores[].identical_count`: describirlo como «renglones idénticos» y evaluar `piece_count` (M6).
3. `docs/contratos/reglas-negocio.md`: RN-13 equilibrio en «compras al año» (K7) (B3) y prioridad «agotado sobre reponer» en RN-15 (B4).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna de valores K31. Las decisiones que M1 y M5 podrían llevar a la puerta 2: de dónde salen los precios de Walmart del 10-oct (carrito o ticket) y quién confirma los empates de nutrición y los vínculos PROFECO.

## Consumo aproximado
39 acciones de herramienta; ~215 KB leídos; ~27 KB escritos (informe 18 KB, este reporte, progreso y scripts temporales ya borrados).
