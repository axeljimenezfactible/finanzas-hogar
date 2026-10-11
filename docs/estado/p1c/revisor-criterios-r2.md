# revisor-criterios-r2 · Flujo P1c (puerta 1b) · 11-oct-2026

## Hecho
- Revisé `docs/criterios/` 1.2 (diff contra HEAD `082235c`, más la corrección r1) solo contra INDICE §1 y decisiones §K (K35–K44). No cambié ningún criterio.
- Hice las verificaciones (a)–(g) con un script propio (`datos/tmp/rev2c/audit.py`, ya borrado), recalculé las cifras con Decimal sobre `oraculo/model.json` y corrí `verif_final.sh` como contraste.
- Comprobé los 5 hallazgos de r1: los 5 están cerrados.
- Veredicto: **NO APROBADO**, con 0 hallazgos altos, 4 medios (M-1 a M-4) y 6 bajos (B-1 a B-6). Detalle en `docs/estado/p1c/REVISION-criterios-r2.md`.

## Archivos
- `docs/estado/p1c/REVISION-criterios-r2.md` (informe)
- `docs/estado/p1c/revisor-criterios-r2.md` (este reporte)

## Verificación
- `python3 -I datos/tmp/rev2c/audit.py`:
  - 18 archivos con las secciones 1–5;
  - 261 criterios (201 / 60 / 0), 0 sin Verifica ni Evidencia, 0 IDs repetidos, 0 IDs citados que no existan;
  - INDICE §1, §2 y §10 coinciden con el conteo, incluidas las versiones;
  - 43 bloqueos, 2 sin «Respuesta de Axel»: `S1-parametros-avisos` (abierto) y `criterios-nxx` (resuelto con «Resolución»);
  - las 10 líneas que nombran un bloqueo resuelto junto a PROVISIONAL son avisos de que la marca se quitó o comprobaciones de que hay 0.
- Cifras: Decimal sobre `oraculo/model.json` (sha256 `1a5706…bfd4`). 973.75/31×7 = 219.88 y 2517.38/28×7 = 629.35, con un solo HALF_UP; con el diario redondeado salen 219.87 y 629.37. Pesos: 0.9836, 0.4918 y 0; el día 91 no compite.
- `grep -n` de K35–K44 en `docs/criterios/*.md`: cada K tiene línea en el archivo que le toca (tabla (b) del informe).
- `bash docs/criterios/tools/verif_final.sh` (con GOOSE): exit 0. 305 citas, 0 sin fuente; canon: 178 cifras, 0 diferencias; S1-09: 12 cifras, 0 diferencias; `conteo.py`: 0 diferencias. El paso 9 sigue mostrando N-320…N-327, el B-01 de r1, en C1-datos.
- Contraste con el esquema: `grep` de las claves de `app_setting` en 00001, 00002 y 00003. Cinco claves de K44 se siembran en 00001 o 00002 (de ahí M-3). `send_time` y `advance_days` están como `@param_sin_valor` en 00003, y el bloqueo `C1-valores-k31` (punto 3) dice que «los valores por defecto los propone S7» (de ahí M-2).

## Supuestos
- `criterios-nxx` cuenta como resuelto por su sección «Resolución», aunque no tenga «Respuesta de Axel»; ningún criterio lo cita con PROVISIONAL.
- S7-08 (4), «al menos una casilla dispara la alerta», lo tomé como lectura razonable y no como valor inventado (bajo B-6).

## Bloqueos
Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`.

## Pendientes
Para el agente de corrección de criterios:
- M-1: `docs/criterios/P3-claude-design.md:28`
- M-2: `docs/criterios/U4.md:67-68`
- M-3: `docs/criterios/C1.md:115-116` y `docs/criterios/U4.md:68`
- M-4: `docs/criterios/P3-claude-design.md:22,31,67`
- Bajos B-1 a B-6, opcionales.

## Solicitudes fuera de alcance
- `backend/db/migrations/00003_puerta1b.sql:45-46`: el comentario dice «lo propone S7 (WF-1)», pero INDICE §2 pone S7 en WF-3. Hay que alinear el flujo en el comentario (lo decide el integrador o S1).

## Propuestas para Axel
Ninguna.

## Consumo aproximado
39 acciones de herramienta, unos 160 KB leídos (por partes) y unos 17 KB escritos. Temporales en `datos/tmp/rev2c` borrados.
