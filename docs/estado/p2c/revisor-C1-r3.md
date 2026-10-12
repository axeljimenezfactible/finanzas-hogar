# Reporte · revisor-C1-r3 (flujo P2c, ronda 3 de 3)

Agente: `revisor-C1-r3`, revisor independiente. 11-oct-2026, noche. No escribí nada de lo que reviso y no corregí nada. Sin `git add` ni `git commit`.

## Hecho
- Revisé C1 v1.3 después de `correccion-C1-r2`, contra `docs/criterios/C1.md` 1.2, `PUERTA-2.md` §2 y §3, las claves de la 00004 (S1-puerta2) y las solicitudes a C1 de `contratos-puerta2.md`.
- Informe: `docs/estado/p2c/REVISION-C1-r3.md`.
- **Veredicto: APROBADO** (0 altos, 0 medios, 6 bajos).
- **M-1 de la ronda 2: CERRADO.**
  - Lo comprobé con una base propia (Up de 00001 a 00004): el alta sin `vat_included` y con 1600, y `{"vat_included": null}` sobre (false, 1600), dan «CHECK constraint failed».
  - Los 8 textos (`CreditCreate`, `CreditUpdate`, `createCredit` y `updateCredit`, en `openapi.yaml` y en `ajustes-c.yaml`) dicen la regla completa: 422 `validation_failed`, `fields[].field = vat_on_interest_bp`, no guarda nada y no deduce `false`.
- Bajos nuevos:
  - **N-1:** `updateExpense` en `openapi.yaml` no lleva la marca PROVISIONAL de `S1-reparto-edicion-gasto` que sí trae `C1-datos.md` §13.
  - **N-2:** las descripciones de `vat_included` usan 1, 0 y NULL en un campo booleano, y `vat_on_interest_bp` conserva la frase de 1.3.
- Los bajos B-1 a B-5 de la ronda 2 siguen abiertos, sin cambio.

## Verificación (comandos y salidas)
Todas con `TMP_ROOT` propio bajo `datos/tmp/rev-c1-r3/`.

| Comando | Salida |
|---|---|
| `make verify-c1 GOOSE=…/gobin/goose TMP_ROOT=…/rev-c1-r3/mk` | rc 0. «los 21 archivos coinciden byte a byte», «comprobaciones OK: 33», «RESULTADO: TODO OK» |
| `uv run … python3 -I backend/internal/api/c1gen/verify_c1.py` | rc 0. 33 OK y 0 FALLA. C1-04 IVA: alta con 1600 {true: rechaza, null: rechaza, false: acepta, ausente: rechaza}; Merge Patch {true: rechaza, null: rechaza, true + null: acepta}; 0 faltantes |
| `make verify-openapi` | rc 0. «docs/contratos/openapi.yaml: OK» |
| `python3 -I backend/internal/api/c1gen/build.py --check` | rc 0. 21 «IDÉNTICO». Generador: 168 rutas, 212 operaciones y 255 esquemas; `routes.go` con 212 rutas; mayor extracto 79,094 bytes |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP_ROOT=…/rev-c1-r3/va` | rc 0. «== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» |
| `sqlite3` sobre mi base: `INSERT`/`UPDATE` de `credit` en `BEGIN … ROLLBACK` | (1, 1600), (NULL, 1600) y la columna ausente con 1600: rechazadas. (0, 1600): aceptada. Merge Patch sobre (0, 1600): `= 1` y `= NULL`, rechazados |
| scripts propios con `yaml.safe_load` (claves, `NotificationUpdate`, `scope`, `assigned_store_id`, `AccrualVsCash`, extractos y `x-provisional`) | ver V-6 a V-19 del informe. 13 claves en las 3 posiciones de settings, 0 `send_time`/`advance_days` en `NotificationUpdate`, 212 de 212 operaciones y 304 esquemas idénticos en los extractos, y 6 marcas con bloqueo abierto |
| `diff /tmp/*.bak-cc1r2` contra las fuentes | la corrección r2 solo tocó las 2 constantes del IVA, la comprobación nueva y 1 renglón del README |

## Supuestos
- **S-1.** Uso como fuente del CHECK la 00004 vigente (sha256 c2aa03512be5…), no el texto copiado en `S1-puerta2.md` l. 29, que es anterior.
- **S-2.** Clasifico N-1 como bajo porque C1-15 admite la marca en `C1-datos.md` («`x-provisional: true` o la marca PROVISIONAL en `C1-datos.md`»), y los textos no se contradicen.

## Bloqueos
Ninguno nuevo. Siguen abiertos `C1-reparto-alcance`, `C1-hora-reportes-job`, `S3-iva-lectura` y `S1-reparto-edicion-gasto`, sin «Respuesta de Axel».

## Pendientes
Los 6 bajos de `REVISION-C1-r3.md` (N-1, N-2 y B-1 a B-5 de la ronda 2). Ninguno bloquea la aprobación.

## Solicitudes fuera de alcance
1. **Integrador o agente de S1.** `S1-puerta2.md` l. 29 cita un CHECK de `credit.vat_included` que ya no es el de la 00004 (`vat_included = 0` en lugar de `vat_included IS 0`). Con NULL, ese texto aceptaría (NULL, 1600). Conviene anotar que lo cambió una corrección.
2. **Integrador.** `bloqueos/INDICE.md`: agregar `C1-reparto-alcance`, `C1-hora-reportes-job` y `S3-iva-lectura`. En `C1-puerta2.md` l. 84, «comprobaciones OK: 29» ahora es 33. `Makefile` l. 30: un TMP por corrida.
3. **Agente de criterios.** `docs/criterios/C1.md` a 1.3: C1-13 con «C1 v1.3», y C1-04 con la comprobación del CHECK del IVA contra la base.

## Propuestas para Axel
Ninguna. Esta revisión no abre preguntas nuevas. La PREGUNTA de `correccion-C1-r2` (¿la API completa `vat_included` = false cuando llega un porcentaje sin ese dato?) sigue en su reporte. Hoy C1 no lo hace, y así está documentado.

## Consumo aproximado
- Unas 38 acciones de herramienta.
- Unos 115 KB leídos:
  - la revisión r2 (15 KB) y la corrección r2 (12 KB);
  - fragmentos de `C1-datos.md` (unos 12 KB), de `verify_c1.py` (5 KB), de la 00004, de `PUERTA-2.md`, de `S3-iva-lectura.md`, de `S3-interfaz.md` y de `C1.md`;
  - las salidas de mis scripts.

  `openapi.yaml` solo se leyó con `python3 -I` y `grep`.
- Unos 22 KB escritos: el informe (14 KB), este reporte (5 KB) y los scripts temporales (3 KB). Los temporales ya están borrados.
