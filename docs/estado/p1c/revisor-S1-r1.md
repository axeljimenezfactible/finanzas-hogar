# revisor-S1-r1 · reporte final (flujo P1c, 11-oct-2026)

## Hecho
- Revisión independiente del esquema 1.3 contra `docs/criterios/S1.md` 1.1 y decisiones K36, K37, K39, K40, K43 y K44. Informe: `docs/estado/p1c/REVISION-S1-r1.md`. Veredicto: **NO APROBADO** (1 alta, 1 media, 5 bajas).
- Verificaciones (a)–(j) reejecutadas con scripts propios en `datos/tmp/rev-s1-r1/` (`rev.py`, `rev2.py`, `prov.py`) y con goose v3.28.0. No se corrigió ningún archivo revisado.

## Archivos
- Creados: `docs/estado/p1c/REVISION-S1-r1.md`, `docs/estado/p1c/revisor-S1-r1.md`.
- Temporales: `datos/tmp/rev-s1-r1/` (borrado al terminar; no se tocaron las carpetas de otros revisores).

## Verificación
- (a) `python3 -I datos/tmp/rev-s1-r1/rev.py`: `integrity_check` ok, `foreign_key_check` 0, 55/55 STRICT.
- (b) `sqlite_master` de 1.2 contra 1.3 con Down: idéntico (278 objetos). Datos: 55 tablas iguales sin contar marcas de tiempo; con goose en una misma base, `cmp` de `sqlite_master` (1295 líneas) y de `.dump --data-only` idénticos.
- (c) `${GOOSE} -dir backend/db/migrations sqlite3 … up | down-to 0 | up` → versión 3, 0 tablas, 55 tablas, ok.
- (d) `git diff --quiet HEAD` y `git show HEAD:… | cmp -` de 00001 y 00002: idénticos.
- (e) `GOOSE=… python3 -I backend/db/tools/s1_verifica.py` → exit 0, RESUMEN TODO OK, KUESKI 389500 / 239528.
- (f) `gen_dict.py datos/tmp/rev-s1-r1/dic.md` + `cmp` → idéntico (182,480 bytes).
- (g) 10 de 10 valores de K44 que viven en `app_setting` iguales a §K. 80 y 100 no están en `app_setting` (H-01).
- (h) 16 líneas con PROVISIONAL cruzadas contra 44 bloqueos: solo citan `S1-parametros-avisos`, que está abierto.
- (i) Decimal HALF_UP: 219.88, 629.35, 439.76; 0.9836, 0.4918, 0; 2,556 y 2,056.
- (j) Disparadores: solo tipo y signo; 0 y 100000 se aceptan, −5 y 0.5 se rechazan; sin mínimos razonables en el esquema.

## Supuestos
- Ninguno que cambie cifras, esquema ni endpoints. Califiqué H-01 como alta porque K44 y la respuesta de `C1-valores-k31` mandan las claves a `app_setting` y la tarea lista 80 y 100 en (g). Si el orquestador acepta la opción A del bloqueo como excepción aprobada por Axel, H-01 se cierra por decisión.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ronda 2: verificar H-01 y H-02 y, si se atienden, las bajas H-03 a H-07.

## Solicitudes fuera de alcance
- Orquestador: decidir H-01, es decir, aplicar K44 (opción B del bloqueo `S1-parametros-avisos`) o registrar la opción A como excepción explícita.

## Propuestas para Axel
- Ninguna nueva. Las del agente S1 sobre mínimos y máximos están bien ubicadas: en propuestas y no en el esquema.

## Consumo aproximado
- Acciones de herramienta: ~28. Leído: ~125 KB. Escrito: ~25 KB (informe 13.4 KB, este reporte ~3 KB, scripts temporales ~8 KB).
