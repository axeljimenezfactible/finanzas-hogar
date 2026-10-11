# correccion-S1-r1 · Corrección de S1 tras la revisión WF-0 (ronda 1)

Fecha: 10-oct-2026 · Hallazgos atendidos: H-01 (alta), H-02 y H-04 (media) de `docs/estado/wf0/REVISION-r1.md`; además H-11 (baja), que también es de S1 · Sin `git commit`.

## Hecho
- **H-01:** tablas `shopping_list` y `shopping_list_item` en `00001_esquema.sql` (producto o texto libre con unidad, cantidad calculada y elegida, tienda, comprado, origen). Cubren N-151, N-152, N-153, N-165, N-127 y N-073. Quedan **PROVISIONALES** (opción A de `C1-lista-compras`, mismo tema en `S1-lista-compras`) hasta la puerta 1.
- **H-02:** `expense.description` admite NULL; `recipe_meal_slot` sustituye a `recipe.default_meal_slot_id`; `inventory_movement` con tipo y motivo; `budget` con vigencia `valid_from_date`/`valid_to_date` en lugar de `period_id`; `planned_item.second_half_amount_cents`; tabla `idempotency_key`. Además se corrigió un hueco que salió al armar la tabla de S1-13: `category.kind` (fijo o variable, N-281). Las tres preguntas de negocio quedaron en `S1-huecos-negocio`.
- **H-04:** 17 líneas `S1-xx` en `S1-modelo-datos.md`. Prueba reejecutable `s1_verifica.py` (Anexo B del reporte de S1, con comando de un renglón), con la ida y vuelta completa: 164 valores, 0 diferencias. Sección «Ajustes: 14 secciones → campos → tabla.columna» en el diccionario. Apariencia queda local por dispositivo, alineada con U1 (N-301) y C1, y se quitó `appearance.*`.
- **H-11:** `ux_person_single_co_holder` (dos cuentas como máximo) y disparador `trg_person_external_without_access`; el tope de Karen ya no se interpreta como 374.93.
- Diccionario 1.1 regenerado con un generador sin plantillas aparte (Anexo A): 54 tablas, 147 índices y 68 disparadores. Reglas: C-13, C-15, C-16, RN-05, RN-14, RN-15 y RN-22.

## Archivos
- `/home/deepsec/finanzas-hogar/backend/db/migrations/00001_esquema.sql` (modificado)
- `/home/deepsec/finanzas-hogar/docs/contratos/diccionario-datos.md` (regenerado, versión 1.1)
- `/home/deepsec/finanzas-hogar/docs/contratos/reglas-negocio.md` (modificado)
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/S1-modelo-datos.md` (evidencia, corrección, Anexos A y B)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-lista-compras.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-huecos-negocio.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/correccion-S1-r1.md` (este reporte)

## Verificación
- `s1_verifica.py` extraído del reporte de S1 y corrido desde el repositorio con goose v3.28.0: **RESUMEN TODO OK**. S1-01 a S1-17 en OK; 46 de 46 negativas rechazadas, 8 de 8 positivas aceptadas; Down deja 0 objetos; goose up/down/up con códigos [0, 0, 0] y 54 tablas. Las líneas completas están en `S1-modelo-datos.md` § Verificación.
- Datos clave: 54 tablas STRICT, 147 índices y 68 disparadores; `integrity_check` ok y `foreign_key_check` vacío. Ida y vuelta: 4 despensas, 16 pagos y 17 eventos; series A/B de 14 valores; mínimos 809450/500000; P1–P6 por persona; suma de pagos de la despensa 1 = 389500 y parte de Axel = 239528, con 164 valores comparados y 0 diferencias. S1-09: 18 periodos recuperados, la quincena 2 de feb-2027 paga el 2027-02-28, la ventana del 10 al 9 se rechaza, 12 de 12 coincidencias y ningún evento cae entre los días 1 y 9. S1-03: 22 reglas con ejemplo y todas las citas coinciden. S1-13: 14 de 14 secciones y 45 referencias tabla.columna, sin ninguna inexistente.
- Generador: antes de regenerar, reprodujo byte a byte el diccionario 1.0 desde la migración original (plantillas y versión sin plantillas). Al regenerar dio «columnas sin descripción: []».
- `datos/tmp` se borró; `datos/` tiene solo `semilla/`. No aplican go build ni go vet, colorimetría ni el validador OpenAPI, porque S1 no tiene código Go, colores ni OpenAPI. La matriz U1 ↔ C1 se cierra cuando C1 regenere sus archivos con este esquema (ver solicitudes).

## Supuestos
- Ninguno cambia cifras. Estas lecturas salen de las fuentes citadas: el formato de la clave de idempotencia es el de C1 §1 punto 8; «tipo (fijo o variable)» de la categoría viene de U1 N-281; los tipos y motivos de inventario son los de U1 N-181 y N-183; la semana de la lista de compras dura 7 días y su día de inicio lo fija U1 (bloqueo `U1-periodos-semana-quincena`).

## Bloqueos
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-lista-compras.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S1-huecos-negocio.md`

## Pendientes
- Respuestas de Axel a los dos bloqueos en la puerta 1. Con la opción B o C de la lista, `shopping_list` y `shopping_list_item` salen del esquema.
- C1 debe regenerar sus archivos con el esquema 1.1 para que la matriz quede en 0 N-xxx sin endpoint (solicitud 1).

## Solicitudes fuera de alcance
1. **C1** (`docs/contratos/openapi.yaml`, `backend/internal/api/routes.go`, `docs/contratos/C1-datos.md` §10, §11.3, §11.4, §13 y §14). Pedido: operaciones para N-151, N-152, N-153 y N-165 sobre `shopping_list_item`, y persistencia de N-127 y N-073 (PROVISIONAL). Además: `expense.description` opcional, sin rellenarla con el nombre de la categoría; `recipe_meal_slot` en lugar de `default_meal_slot_id`; `inventory_movement` (N-181, y N-183 con `reason`); vigencia de `budget`; `planned_item.second_half_amount_cents` (N-279, N-280); `category.kind` (N-281); y las claves de `Idempotency-Key` guardadas en `idempotency_key`, con solo la retención PROVISIONAL. Motivo: H-01 y H-02.
2. **C1:** actualizar «Cómo seguí» en `docs/estado/bloqueos/C1-lista-compras.md` y `C1-huecos-u1-s1.md`, que remiten a `S1-lista-compras` y `S1-huecos-negocio`. Motivo: hoy dicen que S1 no tiene tablas.
3. **Integrador (H-05):** copiar los Anexos A y B de `S1-modelo-datos.md` a `backend/db/tools/gen_dict.py` y `backend/db/tools/s1_verifica.py`, o a la ruta que decida. Motivo: que se puedan reejecutar sin extraerlos del reporte.
4. **S6:** que la purga de `idempotency_key` use la retención que se decida. Que la API escriba `inventory_movement` en la misma transacción que la existencia. Que la API valide los invariantes 7 (vigencias de `budget` sin traslape) y 13 (regenerar la lista conserva lo elegido) del diccionario.

## Paquetes instalados
- Ninguno. Usé el binario goose v3.28.0 que S1 compiló en el scratchpad (`gobin/goose`).

## Consumo aproximado
- Unas 50 acciones de herramienta, sin subagentes. Leí el protocolo, las decisiones, el plan §1 a §3, la revisión, los criterios de S1, partes de U1 y C1 y la migración por tramos. Escribí unos 25 KB de SQL y texto, más el diccionario regenerado (unos 150 KB) y dos scripts (unos 30 KB).
- El contador de la sesión pasó de 300K tokens hacia el final, contando el contexto inicial. Al pasarlo ya no empecé trabajo nuevo: solo terminé la documentación para dejar todo consistente.
