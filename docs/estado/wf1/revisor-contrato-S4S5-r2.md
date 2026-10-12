# revisor-contrato-S4S5-r2 · reporte final (WF-1, 11-oct-2026)

## Hecho
- Revisé `S4-interfaz.md` 1.1, `S5-interfaz.md` 1.1 y `productos-formato.md` 1.0 contra los criterios S4 y S5 1.2 y contra decisiones §D, §E y §K.
- Recalculé el canon y crucé el catálogo con mis propios scripts.
- Comprobé los 7 hallazgos medios de la ronda 1: los 7 están cerrados.
- Veredicto **NO APROBADO**: 0 altos y 4 medios nuevos (M-r2-1 a M-r2-4), más 6 bajos nuevos. Siguen abiertos B1, B2, B3, B4, B6, B7 y B8 de r1.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-contrato-S4S5-r2.md` (informe: veredicto, tabla de criterios, hallazgos y estado de r1)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-contrato-S4S5-r2.md` (este reporte)

## Verificación
- `python3 -I datos/tmp/rev-s4s5-r2/canon.py` sobre `data.json`:
  - K-01 = 309450, K-02 = 323300 y K-03 = 13850 (= 10050 + 5900 − 2100).
  - K-04 / K-05 = 277350 / 287400 y K-06 = 10050.
  - K-07 = 61 / 60 / 55 y K-08 = 33 / 3; ganadores 22 / 8 / 3.
  - Unidad base: 2933 / 2800, 12692 / 16406 y 3375 / 2833.
  - S4 §8.3: 4425 / 2213 / 8850.
  - Frescura: 60/61, 59/61, 30/61 y 0.
  - Metas 2099 / 1596 / 1799 / 1846 / 2556 / 2056; IMC 185 / 250 / 300 / 237 / 185; RN-15 1486571. Todo coincide.
- `python3 -I datos/tmp/rev-s4s5-r2/catalogo.py`:
  - 36 productos y 39 presentaciones; 39 de 39 nombres de `data.json` existen tal cual.
  - §5.2 tiene 39 renglones, 61 piezas y 309450, idéntico a `cuentas`.
  - Precios unitarios y piezas por presentación cuadran.
- Revisé además:
  - Las claves de `app_setting` en 00001 y 00003, RN-12, RN-16 y RN-17.
  - El extracto C1 `BasketComparison`.
  - FORMATO.md §4.1 (firma del evaluador).
  - El formato de los 5 bloqueos.
- Criterios: S4-01 a S4-15 y S5-01 a S5-13 tienen sección. S4-12 queda incumplido (M-r2-4). S4-06, S4-08 y S5-05 quedan con observación (M-r2-2, M-r2-3 y M-r2-1). S5-06 queda «sin medir», porque faltan las fotos.
- S4-09 «S4 verificación manual»: no aplica a esta revisión de contrato, porque `recetario.json` y `tabla.csv` no existen. Lo hará el revisor de la implementación con 5 recetas que él elija.

## Supuestos
- Evalué el contrato como documento sin código: no corrí pruebas Go ni cobertura, porque no hay código que medir.

## Bloqueos
- Ninguno nuevo. Sigue abierto `S1-parametros-avisos.md`, que no toca S4 ni S5.

## Pendientes
- Para la corrección r2 del contrato: M-r2-1 (`S5-interfaz.md` §15.1 y §15.2), M-r2-2 (`S4-interfaz.md` §12.1, §8.1 y §6.2), M-r2-3 (`S4-interfaz.md` §3.2 y §9.1) y M-r2-4 (`S4-interfaz.md` §6.5).
- Bajos opcionales: B-r2-1 a B-r2-6, y B1 a B4, B6 a B8 de r1.

## Solicitudes fuera de alcance
- `oraculo/fixtures/FORMATO.md` §4.1: solo si la corrección de M-r2-1 elige pasar la carpeta de fixtures al evaluador en lugar de usar rutas desde la raíz.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna. Esta revisión no fija valores de K31.

## Consumo aproximado
39 acciones de herramienta; unos 265 KB leídos; unos 30 KB escritos (informe 17 KB, este reporte 3 KB, scripts temporales 7 KB, ya borrados).
