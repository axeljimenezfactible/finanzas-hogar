# Bloqueo S6-compras-b-fuentes

Abierto por `api-compras-b` (WF-3a, 12-oct-2026). Se pregunta en la puerta 4. Código: `backend/internal/api/comprasb/fuentes.go`, `comprasb.go`.

## Qué falta (PREGUNTAS)

1. **Veredicto de cada cadena con sitio (`getPriceSources.sites[].verdict`).** La base 1.4 no guarda veredictos; `docs/estado/S5-terminos.md` es registro histórico y desde K53 la app no lee sitios. Opciones: (a) `missing` para todas mientras no haya columna; (b) `denied` para todas por K53; (c) S1 agrega una columna con el veredicto histórico.
2. **`profeco_last_import_at`.** C1 dice «(job)», pero no hay un job de importación PROFECO registrado. Opciones: (a) el alta más reciente de un precio `profeco` (`max(price.created_at)`); (b) registrar un job `profeco_import` y que `importProfeco` actualice su `last_run_at`.
3. **RAR en `importProfeco`.** C1 admite RAR o CSV, pero el motor `profeco` solo lee CSV. Opciones: (a) RAR responde 415 y se sube el CSV; (b) descomprimir RAR (necesita `unrar` o `bsdtar`: paquete y decisión nuevos); (c) C1 deja solo CSV.
4. **Dónde vive `vinculos.csv`** (ligas confirmadas registro → presentación). Opciones: (a) `<DataDir>/profeco/vinculos.csv`; (b) `datos/semilla/profeco/vinculos.csv` (hoy no existe); (c) tabla nueva en S1.
5. **Etiqueta del corte** (`Options.Corte`, va en `price.source_ref`). Opciones: (a) AAAA-MM del día de la importación; (b) la manda la persona en el formulario (C1 no tiene ese campo).
6. **Bitácora de la importación PROFECO.** El motor (`LoadIntoDB`) escribe los precios sin `audit_log`. Opciones: (a) sin bitácora por renglón (puede ser mucho volumen); (b) un renglón de resumen (el CHECK de `audit_log` pide una tabla y un `row_id`); (c) por renglón.
7. **Discrepancia motor ↔ C1 1.4 (pregunta a de S5-agente-externo).** C1 dice PROVISIONAL que los precios del agente externo (`site_assisted`) «compiten como manual», pero el motor aprobado `precios.Competes("site_assisted")` = false. `getPriceSummary` usa el motor sin cambios, así que un precio del agente no cuenta como fresco ni da «precio más reciente» por tienda. Resolverlo es del dueño del motor S5 (no de WF-3).

## Mientras tanto (PROVISIONAL)

1 (a) · 2 (a) · 3 (a): 415 con mensaje «descomprime el RAR y sube el CSV» · 4 (a): sin archivo ningún registro se liga (`skipped_count`) · 5 (a) · 6 (a) · 7: se respeta el motor.

## Impacto

Solo `getPriceSources`, `importProfeco` y `getPriceSummary`. Ninguna cifra canon cambia; si 7 cambia en el motor, `getPriceSummary` lo sigue sin tocar la API.

**Estado: ABIERTO.**
