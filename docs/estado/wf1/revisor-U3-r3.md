# Revisor U3 · ronda 3 de 3 (WF-1) · reporte final

Agente `revisor-U3-r3` · 11-oct-2026 · Revisor independiente de `docs/design/U3-dashboard.md` v1.2 contra `docs/criterios/U3.md` v1.2, decisiones §G y §K. Veredicto: **APROBADO** (0 altos, 0 medios). Informe: `docs/estado/wf1/REVISION-U3-r3.md`.

## Hecho
- Leí el protocolo completo, decisiones §B, §G y §K, los criterios U3 completos y U3-dashboard.md completo, por partes. También `correccion-U3-r2.md`, el reporte `U3-dashboard.md`, el bloqueo `U3-t2-acumulados.md`, DESIGN §4.2, §7.6 a §7.8 y §8, y los extractos (`inicio`, `presupuesto-a`, `presupuesto-b`, `creditos-a`, `creditos-b`, `ajustes-e` y `componentes-comunes`), por grep y awk.
- Rehíce cada verificación con scripts propios en `datos/tmp/rev-u3-r3/` (`t2.py`, `check.py`, `fields.py`; `python3 -I`, solo biblioteca estándar). Ya los borré.
- Comprobé cerrados H-1 (alto), H-2 y H-3 (medios) de la ronda 2.
- Escribí la tabla de los 12 criterios con mi evidencia, el estado de L-1 a L-6 (siguen abiertos, bajos) y 3 bajos nuevos (N-1 a N-3).

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-U3-r3.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-U3-r3.md` (este reporte)

## Verificación
- `t2.py` (Decimal sobre `oraculo/model.json`):
  - Acumulados desde P1: P1 973.75 / 1,947.50; P2 2,921.25 / 3,895.00; P6 13,103.78 / 14,362.48.
  - En los 6 periodos, devengado acumulado − efectivo acumulado = `diferencia_acumulada`.
  - Ventana desde P2 (variante b): 12,130.03 / 12,414.98.
  - La regla de la v1.1 daba 2,921.25 y 13,388.73, y 13,388.73 ya no aparece en U3 (grep = 0).
  - Todas las cifras de A-11 coinciden.
- `check.py`:
  - 27 de 27 `operationId` existen; 13 de 13 rutas GET existen.
  - M-01 a M-11 y M-13 a M-17 existen en C1-datos.
  - 0 fechas ISO; 0 hex o `rgb(`; 30 fechas con día de la semana, 0 erróneas.
  - 57 montos (20 distintos): 18 literales, $138.50 = K-03 y $2,921.25 = R-45 + R-46 (también en `S2.md` S2-07).
  - 58 tokens, 0 sin definir.
- `fields.py`: 147 identificadores en minúsculas. Los únicos ausentes de los extractos son iconos, `tnum`, los 2 campos pedidos a C1 (variante (a), PROVISIONAL) y `weekday`/`day_of_month`, que U3 cita solo para decir que no los usa. Los enums coinciden con `inicio.yaml`.
- Parámetros reales de cada operación contra §6: cada «S» corresponde a un parámetro existente. `credit_id` es un arreglo. `Period` admite quincena. `Family` tiene los valores de K5.
- Criterios: U3-01 a U3-08 (altos) cumplen. U3-09 a U3-12 (medios) cumplen, con observaciones bajas. Evidencia por criterio en la tabla de `REVISION-U3-r3.md`.

## Supuestos
- Mantengo L-2 (T-3 repinta series al filtrar) en severidad baja, como en r2: DESIGN §8.6 fija la regla y U3 admite la lectura «posición en `listCredits`», así que quien implemente no tiene que adivinar.

## Bloqueos
- Ninguno nuevo. U3 cita, con su marca PROVISIONAL, `docs/estado/bloqueos/U3-t2-acumulados.md` y los de S2 (`S2-presupuesto-total`, `S2-flujo-quincenal`, `S2-metas`, `S2-proyeccion-interruptor`), todos abiertos y sin «Respuesta de Axel».

## Pendientes
- Ninguno dentro del alcance. Los bajos L-1 a L-6 y N-1 a N-3 quedan para quien corrija U3 si el orquestador lo decide; no bloquean.

## Solicitudes fuera de alcance
1. `docs/estado/bloqueos/INDICE.md`: actualizar el conteo, porque ya no hay un solo bloqueo abierto (`U3-t2-acumulados` y los de S2 de WF-1). Es la misma solicitud 4 de `correccion-U3-r2.md`.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna de K31: U3 no tiene parámetros operativos. Coincido con que la opción A de `U3-t2-acumulados` es la que conserva DESIGN §8.3 sin desviaciones (lo decide la puerta 2).

## Consumo aproximado
- Acciones de herramienta: 46.
- Leído: unos 240 KB (U3 79, criterios 13, protocolo y decisiones 20, DESIGN 22, reportes y bloqueos 35, extractos por partes 20, salidas de scripts 25, otros 26).
- Escrito: unos 27 KB (scripts 7, informe 12, este reporte y el progreso 8).
