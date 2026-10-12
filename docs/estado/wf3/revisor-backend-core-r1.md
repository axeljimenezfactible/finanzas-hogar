# revisor-backend-core-r1 · reporte

12-oct-2026 · revisor independiente, ronda 1 de 3 · cerrado_por_tope = false. Informe completo: `docs/estado/wf3/REVISION-backend-core-r1.md`.

## Hecho
- Revisé contra S6.md 1.2 (todo) y C3-03, C3-04 y C3-09 (parte del servidor) el trabajo de backend-base, auth-sistema, semilla, api-gastos-inicio, api-presupuesto y api-comidas-salud.
- Volví a correr todas las verificaciones con mis propios comandos. Mis pruebas de las 4 áreas de API (5 o más operaciones cada una, validadas contra `openapi.yaml` con `testutil`) las metí con `go test -overlay`, sin escribir en `backend/`.
- Construí el binario y lo probé (migrar, sembrar, respaldar, servir). Restauré el respaldo y comparé los conteos de las 56 tablas. Crucé la semilla con `oraculo/model.json`.
- Veredicto: **NO APROBADO** (6 altos, 3 medios, 4 bajos).

## Archivos
- Escritos: `docs/estado/wf3/REVISION-backend-core-r1.md` y este reporte.
- Temporales en `datos/tmp/rev-bc-r1/`, `datos/tmp/s6-*.cov` y `datos/tmp/restaura.db*`: borrados al terminar.
- No toqué código, migraciones, contratos ni la Pista 3. No hice git add ni git commit.

## Verificación (resumen; detalle en el informe)
- **S6-01**: build y vet OK; `gofmt -l` lista 2 archivos de `internal/reportes`. NO CUMPLE.
- **S6-02**: `go test ./...` → FAIL en `internal/api` (contrato). Coberturas: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8 (todas ≥ 85); api 81.9, auth 83.6, jobs 92.3, reportes 77.2 y telegram 80.0. NO CUMPLE.
- **S6-03**: `build.py --check` = 0; `verify_c1.py --db` = TODO OK (216/216, 0 sobrantes y 0 faltantes). `TestOp_*` cubre 214 de 216. El binario responde 501 en `login`. NO CUMPLE.
- **S6-04**: la prueba de contrato da 191 de 216. K9 de gasto y de peso correcto (1 fila); límite de página 200/201 correcto; 20 mensajes, 0 en inglés. NO CUMPLE.
- **S6-05**: PARCIAL. Cookie con HttpOnly, Secure, SameSite=Strict y Max-Age 7776000; 11 caracteres → 422 y 12 → 204; Karen → 401; `$argon2id$` ×2.
- **S6-06**, **S6-07**, **S6-08**, **S6-13** y **S6-15**: CUMPLEN.
- **S6-09**: CUMPLE (integrity `ok` y 0 diferencias en 56 tablas); hora y retención PROVISIONAL.
- **S6-10**: CUMPLE; el factor de bytes es PROVISIONAL.
- **S6-11**: NO VERIFICADO (sin mkcert). **S6-12**: fuera de esta ronda. **S6-14**: parcial, pendiente de WF-4.
- **S6-16**: NO CUMPLE (los 3 bloqueos de la semilla no se llaman `S6-*`).
- **C3-03**: NO CUMPLE. **C3-04**: NO CUMPLE. **C3-09** (servidor): CUMPLE (A 10,000.00 / B 6,905.50 · 973.75 el 2026-10-30 · 6,293.44 al 2027-04-10 · 16 pagos = oráculo).

## Supuestos
- La falla de la prueba de contrato cuenta para los agentes revisados aunque la prueba sea de pruebas-contrato: las causas están en su código (H1, H2 y H3).
- Las fallas de s7-pdf y api-compras-a (H5 y H6) se reportan porque S6-01 y S6-02 son globales. Su corrección es de esos dueños.
- S6-11, S6-12 y S6-14 no estaban en la lista mínima de esta ronda; quedan como «no verificado» o «parcial» sin hallazgo.

## Bloqueos
Ninguno nuevo. PREGUNTA B4 (`auth.password_min_length` con `minimum: 0` en C1) para el dueño de C1 y para Axel.

## Pendientes
Correcciones H1 a H6 y M1 a M3 (rutas y dueños en el informe). Ronda 2 después de las correcciones.

## Propuestas para Axel
- Fijar en C1 el mínimo de `auth.password_min_length` en 12 (K6) para que no se pueda bajar desde Ajustes.
- Confirmar qué código interno usa el motor para las categorías que crea el usuario (H1), si eso cambia algo visible.

## Paquetes instalados
Ninguno.

## Consumo aproximado
~42 acciones · ~200 KB leídos · ~32 KB escritos (informe 16 KB, reporte 4 KB y pruebas temporales 7 KB, ya borradas).
