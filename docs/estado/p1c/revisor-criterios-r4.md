# revisor-criterios-r4 · P1d, ronda 4 de 4 de criterios 1.2 (11-oct-2026)

## Hecho
- Revisión independiente de la ronda 4 de los criterios contra `INDICE.md` §1 y `decisiones.md` §K (K35 a K44). Veredicto: **NO APROBADO** (1 alta, 2 medias, 3 bajas). Informe: `docs/estado/p1c/REVISION-criterios-r4.md`.
- Comprobé el estado de los hallazgos A a G de r3, de I-02, de B-01 de C1-r2 y de U3/K38: A parcial (H-1, H-2); B, C, E, F e I-02 cerrados; D y G atendidos; B-01 parcial (H-3).
- No corregí nada y no hice commit.

## Archivos
- Creados: `docs/estado/p1c/REVISION-criterios-r4.md` y este reporte.
- Temporales: `datos/tmp/rev4c/` (scripts propios, base de goose y logs), borrada al terminar.

## Verificación
- `bash docs/criterios/tools/verif_final.sh`: rc = 0; 261 (201/60/0), traza y refs completas, citas 305/0, canon 178/0, S1-09 OK, matriz con 210 operaciones, INDICE contra conteo: 0 diferencias.
- `GOOSE=/tmp/claude-1000/…/gobin/goose make verify-criterios`: rc = 0, «las 11 comprobaciones dan lo esperado».
- `goose … sqlite3 datos/tmp/rev4c/t.db up`: 00001 a 00003 OK (versión 3); `alerts.*` = 80 / 100 / [1, 0]; `params_json`: `weekday` 7, `day_of_month` 1, sin `thresholds_pct` ni `days_before`.
- `c122.py` (Verifica de C1-22, literal): 0 fallas. Parte ejecutable de U4-12: 0 fallas (U4-ajustes.md todavía no existe).
- `verify_c1.py --db` sobre la base de goose: 29 OK, TODO OK.
- `c112.py` (Verifica literal de C1-12): públicas `getHealth` y `login`; 1 propiedad sensible por nombre, `SecuritySettings.auth.password_min_length` → H-3.
- `chk.py`: estructura y conteo propios = INDICE §2 y §10; Decimal 219.88 / 629.35 / 0.9836 / 0.4918.
- `prov.py`: 0 marcas PROVISIONAL activas que citen alguno de los 42 bloqueos con «Respuesta de Axel».
- `grep -c operationId docs/contratos/openapi.yaml` = 210 = `extractos/INDICE.md:31`.
- Regex de P3 con textos de muestra: 6 advertencias dan 0; credenciales 5/6 (P3-01/12) y 6/6 (P3-03).

## Supuestos
- U4-12 se evaluó en su parte ejecutable (cruce con la base y el inventario `@param`), porque `docs/design/U4-ajustes.md` se escribe en WF-1.
- Tomé la parte de C1-12 como parte de B-01 de C1-r2, según `revisor-C1-r2.md:39` e `INTEGRACION.md` §7, aunque el encargo resumía B-01 como C1-13 y C1-21.
- Mis corridas sobre el contrato son una foto: `correccion-C1-r3` trabajaba en paralelo en `docs/contratos/`.

## Bloqueos
Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`.

## Pendientes
Corrección de H-1, H-2 y H-3 (y de las bajas H-4 a H-6, si se quiere) en `docs/criterios/U4.md`, `S7.md`, `C1.md`, `INDICE.md` §6, `P3-claude-design.md` y `U3.md`, como dice el informe. Esta era la ronda 4 de 4: lo que sigue lo decide el orquestador.

## Solicitudes fuera de alcance
- Ninguna. Observación sin pedido: en los extractos, `weekday` y `day_of_month` solo aparecen en el texto de `params_json` de `ajustes-d.yaml` (línea 1180), no como propiedades declaradas; es de C1 y lo atiende el agente en paralelo si lo considera.

## Propuestas para Axel
- Ninguna. Los hallazgos son mecánicos. Sigue pendiente su respuesta a `S1-parametros-avisos` (opciones A, B o C).

## Consumo aproximado
~43 acciones de herramienta; lectura ≈ 185 KB; escritura ≈ 20 KB.
