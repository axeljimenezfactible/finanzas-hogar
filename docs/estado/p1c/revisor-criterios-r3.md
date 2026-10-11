# revisor-criterios-r3 · Revisión independiente de los criterios 1.2 (flujo P1c, ronda 3 de 3) · 11-oct-2026

## Hecho
- Revisé `docs/criterios/` (diff contra HEAD, `tools/`, `INDICE.md` y el reporte `criterios-1.2.md`) solo contra `INDICE.md` §1 y decisiones §K (K35–K44). No corregí nada.
- Reejecuté por mi cuenta las comprobaciones (a)–(g) y revisé los 4 hallazgos medios de la ronda 2.
- Veredicto: **NO APROBADO**, con 1 alto (A: parámetros de avisos contra K44 y el esquema 1.3), 2 medios (B: conteo de P3-06; C: Verifica de P3-04) y 4 bajos (D–G). Los 4 hallazgos de r2 quedan cerrados, dos con residuo bajo.
- Informe: `docs/estado/p1c/REVISION-criterios-r3.md`.

## Archivos
- Creados: `docs/estado/p1c/REVISION-criterios-r3.md` y este reporte.
- Temporales: `datos/tmp/rev3c/` (script `chk.py`, textos de muestra, base `t.db`), ya borrados.

## Verificación
- `bash docs/criterios/tools/verif_final.sh`: 261 criterios (201/60/0); traza, refs y citas (305) sin fallas; canon 178/0; S1-09 12/0; conteo contra INDICE 0 diferencias; paso 9 con N-320…N-327 (B-01 de r1, fuera de los criterios).
- `python3 -I datos/tmp/rev3c/chk.py` (propio):
  - Estructura: 18 archivos, 261 criterios; 0 sin Verifica, Evidencia o §5. S2, S3 y S5 traen un §6 anexo desde WF-0.
  - 0 PROVISIONAL activos que citen alguno de los 42 bloqueos resueltos.
  - Cifras: 219.88, 629.35, 0.9836, 0.4918, 0 y «no compite» a los 91 días. 973.75 y 2,517.38 están en `model.json`.
- `goose` v3.28.0 `up` (00001–00003) en una base temporal:
  - `app_setting` trae `alerts.budget_warning_pct` = 80, `alerts.budget_limit_pct` = 100 y `alerts.payment_due_days_before` = [1, 0].
  - En `params_json`, `payment_due` queda `{}` y `budget_threshold` queda `{"scopes": [...]}`.
  - Las 23 claves citan una decisión según la regex de C1-22.
  - Esto confirma el hallazgo A.
- `grep -n -i -E` con las regex de P3-01, P3-03 y P3-12 sobre dos textos de muestra:
  - Prosa correcta: 0 coincidencias.
  - Credenciales: 3 coincidencias.
  - Advertencia en formato de etiqueta: 1 o 2 coincidencias (hallazgo D).

## Supuestos
- Evalué el estado de los archivos a las 15:30 del 11-oct-2026. La 00003 cambió a las 15:23:44 (corrección S1-r1) y los criterios a las 15:24. Tomo el esquema actual como el vigente.
- Mantengo la calibración de r1 y r2: una Verifica que falla con el artefacto correcto es media. A es alta porque además reabre K44 y deja PROVISIONAL algo decidido.

## Bloqueos
Ninguno nuevo.

## Pendientes
- Corregir A en `docs/criterios/C1.md` (C1-22), `U4.md` (U4-12), `S7.md` (S7-04, S7-08, S7-10, lista de cambios y §5) e `INDICE.md` (§6, líneas 68 y 96; §11, línea 174).
- Corregir B y C en `docs/criterios/P3-claude-design.md` (P3-06 y P3-04).
- Los bajos D–G quedan a criterio del orquestador.

## Solicitudes fuera de alcance
- `backend/db/migrations/00003_puerta1b.sql:53-54`: los `@param_sin_valor` de `send_time` y `advance_days` dicen «lo propone S7 (WF-1)». Según `INDICE.md` §2, S7 es de WF-3, y K31 aprueba los valores de reportes en la puerta 4. Sugerencia: «lo propone S7 (WF-3, K31, puerta 4)».

## Propuestas para Axel
- Ninguna nueva. Sigue vigente la de `correccion-criterios-r2.md`: K43 dice que `advance_days` «nace» con valor, pero §A–§K no lo da.

## Consumo aproximado
- Acciones de herramienta: 33 de 60.
- KB leídos: ~150 de 400 (protocolo 6, §K 15, reportes 12, criterios por partes ~70, migraciones y bloqueos ~20, salidas de scripts ~25).
- KB escritos: ~18 de 150.
- Sin paquetes instalados.
