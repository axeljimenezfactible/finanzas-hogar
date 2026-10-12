# revisor-contrato-S4S5-r3 · reporte final (WF-1, 11-oct-2026)

## Hecho
- Revisé, como revisor independiente (ronda 3 de 3), `docs/contratos/motores/S4-interfaz.md` 1.2, `S5-interfaz.md` 1.2 y `productos-formato.md` 1.1. Evalué solo contra los criterios S4 1.2 y S5 1.2 y decisiones §D, §E y §K. No corregí nada.
- Rehíce todas las cifras con un script propio (canon de la comparativa, unidad base, catálogo de los 61 artículos, §8.3, frescura, salud, RN-15, RN-16) y revisé la estructura (criterio → sección, JSON, `-run`, bloqueos citados, claves de `app_setting`).
- Veredicto: **NO APROBADO**. Hay 0 altos, 2 medios nuevos (M-r3-1, M-r3-2) y 2 bajos nuevos. Los 4 medios de r2 están cerrados. Informe: `docs/estado/wf1/REVISION-contrato-S4S5-r3.md`.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-contrato-S4S5-r3.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-contrato-S4S5-r3.md` (este reporte)
- Temporal: `datos/tmp/rev-s4s5-r3/verif.py`, borrado al terminar.

## Verificación
- `python3 -I datos/tmp/rev-s4s5-r3/verif.py` → `OK 63 FALLAS 1`. La falla es un falso positivo del filtro: la línea S4-12 de §15 cita `S4-asignacion-mercado` sin la palabra «bloqueo».
- Canon: Bodega 309450, Walmart 323300, **+13850** = 10050 + 5900 − 2100. Idénticos 277350 / 287400, **+10050**. Piezas 61 / 60 / 55, renglones 33 / 3, ganadores 22 / 8 / 3. Todo igual a `totales` y `resultados_por_producto`.
- Unidad base: 2933/2800, 12692/16406 y 3375/2833.
- Catálogo: 36 productos y 39 presentaciones. Los 39 nombres de `data.json` existen tal cual. §5.2 es igual a `cuentas`, renglón por renglón: 39 renglones, 61 piezas y 3,094.50. Precio y unidad coherentes; piezas por presentación igual a `items_identicos` (atún 5 + 1).
- §8.3: 4425 / 2213 / 8850. Frescura: 1, 60/61, 59/61, 30/61 y 0; el día 91 no compite; 0.9836 / 0.4918 / 0.0000.
- Salud: MSJ 1748.75, 1330.25 y 1648.75; metas 2099, 1596, 1799, 1846, 2556 y 2056; factores 1.2 a 1.9 exactos. IMC 185, 250, 300, 237 y 185 (18.4975).
- RN-15: 1486571. RN-16: 2026-10-13.
- Estructura:
  - S4 §15 = S4-01 a S4-15 y S5 §18 = S5-01 a S5-13.
  - 2 bloques JSON válidos.
  - 17 expresiones `-run` con pruebas y 0 mezclas.
  - Los bloqueos citados existen (`S4-meta-ajustes` es hipotético, como dice S4-07).
  - Las claves `price.*` y `health.goal_adjust_*` están en 00001 y 00003.
- Comprobé contra el texto vigente que los 4 medios de r2 están cerrados. Comparé la lista con C1 (`compras-tiendas-a.yaml`: `ShoppingListItem`, `ShoppingList.groups`) y de ahí salen M-r3-1 y M-r3-2.
- Criterios cubiertos y su evidencia: tabla de criterios del informe (28 criterios; 6 con observación por M-r3-1 o M-r3-2).

## Supuestos
- Ninguno que cambie cifras. Califiqué la contradicción con C1 como media («inconsistencia entre artefactos») aunque la regla de S4 está marcada PROVISIONAL, porque ni el bloqueo ni las solicitudes la mencionan.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`; no lo toca este artefacto.

## Pendientes
- Para el corrector, si hay ronda posterior o integración: M-r3-1 en `S4-interfaz.md` §6.5 y §12.2, el bloqueo `S4-asignacion-mercado` y una solicitud a C1; M-r3-2 en `S5-interfaz.md` §14, §7.3, §10 y §15.2, y `S4-interfaz.md` §8.1. Los bajos B-r3-1 y B-r3-2 y los bajos abiertos de r1 y r2.

## Solicitudes fuera de alcance
1. C1 v1.2 (`ShoppingListItem`, `ShoppingList.groups`): alinear `recommended_store_id`, `groups[].store_id` y `visit_weekday` con la asignación de mercados de S4 §6.5 (M-r3-1). La decide el integrador junto con el bloqueo `S4-asignacion-mercado`.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Ninguna nueva. Para la puerta 2: decidir el bloqueo `S4-asignacion-mercado` junto con el campo de C1 donde vive la tienda asignada.

## Consumo aproximado
33 acciones de herramienta; unos 230 KB leídos; unos 36 KB escritos (script 9 KB, informe 16 KB, reporte y progreso 7 KB).
