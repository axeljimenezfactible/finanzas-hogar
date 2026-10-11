# Revisor S1 · ronda 3 de 3 · reporte de agente (flujo P1)
Fecha: 2026-10-11 · Agente: revisor-S1-r3 · Informe: `docs/estado/p1/REVISION-S1-r3.md` · Veredicto: **APROBADO** (0 altos, 0 medios, 9 bajos).

## Hecho
- Revisé el esquema 1.2 de S1 solo contra `docs/criterios/S1.md` 1.1 y §K (K1–K3, K6, K8, K9, K11–K16, K18–K23). Revisé:
  - `00002_puerta1.sql`, entera;
  - los comentarios cambiados de `00001`;
  - `backend/db/tools/`;
  - el diccionario y las reglas, por partes;
  - el reporte `S1-esquema-1.2.md`.
- Reejecuté las verificaciones (a) a (j) con comandos y scripts propios. Además hice 4 mutaciones del verificador en una copia temporal y pruebas propias de K8, K4, K12, K18 y K23 sobre la base.
- R2-01 (alta) queda cerrado. Abrí un hallazgo bajo nuevo (R3-01). R2-02 a R2-09 siguen abiertos como bajos.
- No corregí nada ni escribí fuera de mi alcance.

## Archivos
- Nuevo: `docs/estado/p1/REVISION-S1-r3.md`.
- Nuevo: `docs/estado/p1/revisor-S1-r3.md` (este reporte; antes llevaba el progreso de la acción ~30).
- Temporales, ya borrados: `datos/tmp/rev.db`, `b1.db`, `b2.db`, `c.db`, `s1v.out`, `dic-regen.md`, `rev_s1_r3.py`, `rev_s1_r3b.py`, `rev_s1_r3c.py` y `datos/tmp/mut/`.

## Verificación
1. **goose v3.28.0 `up` en `datos/tmp/rev.db`:**
   - rc 0, versión 2;
   - `integrity_check` ok y `foreign_key_check` 0 filas;
   - 55 tablas STRICT (no lo es solo `goose_db_version`).
2. **Down de 00002:** `up-to 1` contra `up` + `down` da `sqlite_master` igual (232 objetos, en el mismo orden). Los datos son iguales sin contar las marcas de tiempo.
3. **goose up / down-to 0 / up:** códigos [0, 0, 0], versiones [2, 0, 2] y tablas [55, 0, 55].
4. **00001 contra HEAD, sin comentarios** (lexer propio): Up, Down, anotaciones goose, `sqlite_master` y semillas son iguales. Las 20 líneas `+`/`−` del diff son todas de comentario.
5. **`GOOSE=… python3 -I backend/db/tools/s1_verifica.py`:**
   - rc 0, «RESUMEN TODO OK»;
   - KUESKI 389500 y 239528; NEG 46/46; NEG12 26/26 y 12/12; S1-12 con 164 valores y 0 diferencias.
   - Mutaciones en `datos/tmp/mut`: quitar el rótulo de RN-06 o de RN-11 da FALLA, y 6,293.45 en RN-21 da FALLA. 5,226.40 dentro de la frase rotulada de RN-11 da OK (R3-01).
6. **`gen_dict.py datos/tmp/dic-regen.md` + `cmp`:** idéntico (174,605 bytes).
7. **K por grep con línea:** todas las de la lista están aplicadas (tabla en el informe).
8. **PROVISIONAL:**
   - 0 marcas citan un bloqueo resuelto;
   - las 10 marcas de las reglas citan 4 bloqueos abiertos;
   - 0 marcas en las migraciones; en el diccionario solo está la definición de la marca.
9. **Semillas:** la 00002 solo siembra valores de K21, K11 y K6. Las claves `optimizer.*` de K16 no tienen valor.
10. **Reglas** (escáner propio de 28 ejemplos):
    - 0 reglas sin ejemplo y 0 fechas fuera del oráculo sin rótulo;
    - única cifra fuera del oráculo sin rótulo: 1,814.50, la cuenta citada de R2-02;
    - 24 citas con ruta y 0 discrepancias.
11. **Criterios con pruebas propias:** S1-01 (31/31), S1-04 (55/55 y 5 comprobadas contra U1), S1-05, S1-07, S1-09 (18 filas, límites de K8, rechazo de la ventana 10→9 y equivalencia 12 de 12), S1-10 (17.27 y 25.8685), S1-13, S1-14, S1-15 y S1-17.

## Supuestos
- Mantengo como **bajos** R2-02 (cuenta citada exacta) y R2-09 (rótulo «ilustrativo» sin el paréntesis literal), igual que r2. Las cifras son exactas y nada se presenta como canon, así que se cumple el propósito de S1-03.
- La marca PROVISIONAL del título de RN-17 (R:36 y R:216) es del ajuste por objetivo de K23, no del IMC de K7. No la cuento contra S1-16.

## Bloqueos
Ninguno.

## Pendientes
Ninguno dentro de mi alcance. Los 9 bajos quedan para el integrador, con ruta y corrección en el informe.

## Solicitudes fuera de alcance
1. `backend/db/tools/s1_verifica.py` (bloque S1-03): comparar también las cifras con ruta dentro de las frases rotuladas (R3-01).
2. `docs/estado/p1/S1-esquema-1.2.md`: actualizar las cifras o enlazar `correccion-S1-r1.md` y `correccion-S1-r2.md` (R2-04).
3. Las correcciones de R2-02, R2-03 y R2-05 a R2-09 que propone el informe, si el integrador las toma antes del commit.

## Paquetes instalados
Ninguno. Usé goose v3.28.0 del scratchpad, sqlite3 3.53.4 y Python del sistema.

## Consumo aproximado
- 35 acciones de herramienta.
- ~210 KB leídos:
  - protocolo 6 y §K 8;
  - criterios S1 23;
  - 00002 16 y diff de 00001 10;
  - reglas ~40 y diccionario ~12;
  - reportes de r2 y de corrección ~25;
  - S1-esquema-1.2 11;
  - salidas de los scripts ~60.
- ~47 KB escritos: scripts temporales ~25, informe 16, este reporte ~5 y progreso 0.4.
