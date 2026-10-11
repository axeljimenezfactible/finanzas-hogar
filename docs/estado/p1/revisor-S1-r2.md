# Reporte del agente revisor-S1-r2 · flujo P1 «Aplicar respuestas de la puerta 1»

Revisor independiente del esquema 1.2 (ronda 2 de 3) · 11-oct-2026 · No escribí nada de lo que reviso y no corregí nada. Informe: `docs/estado/p1/REVISION-S1-r2.md`.

## Hecho
- Leí, siempre por partes:
  - el protocolo completo y decisiones §K;
  - `docs/criterios/S1.md` 1.1 completo y los reportes `S1-esquema-1.2.md` y `correccion-S1-r1.md`;
  - `00002_puerta1.sql` entera y el `git diff` de comentarios de 00001;
  - en `reglas-negocio.md`, la cabecera, todos los ejemplos y las definiciones de RN-02, RN-04, RN-06, RN-07, RN-11, RN-12, RN-15, RN-17 y RN-21 a RN-28;
  - del diccionario, por `grep` y `sed -n`: semillas, claves sin valor, Ajustes, notas, `login_attempt` y filas nuevas;
  - de U1, solo las filas K3, K4, K12 y K23 y los nodos citados;
  - la respuesta de Axel en los bloqueos `S1-meta-kcal` y `S1-redondeo-series`;
  - el estado (con o sin «Respuesta de Axel») de los 42 archivos de bloqueo, y el contenido de `U1-lista-vigente-domingo`.
- Reejecuté las verificaciones (a) a (j) con scripts propios. Pasan de la (a) a la (i); la (j) no pasa (R2-01).
- Evalué S1-01 a S1-17: 16 cumplen y S1-03 no cumple. S1-04 y S1-15 los comprobé con muestras que elegí yo.
- Comprobé los 3 hallazgos de r1: S1-09/K8 y S1-04 quedan cerrados; S1-03 queda reabierto en parte. También revisé las bajas H-04 a H-08: cuatro siguen abiertas y H-08 está abierta en parte.
- Escribí `REVISION-S1-r2.md` con veredicto, tabla de criterios, hallazgos con corrección sugerida, estado de r1 y trazabilidad K con líneas.

## Archivos
- Creados: `docs/estado/p1/REVISION-S1-r2.md` y `docs/estado/p1/revisor-S1-r2.md` (este reporte; a las 29 acciones sirvió de archivo de progreso).
- Temporales borrados al terminar: `datos/tmp/rev.db` y `datos/tmp/rev-s1-r2/`. Ahí estaban `abcd.sh`, `strip.py`, `ej.py` y `neg.py`, los volcados, `g.db` de goose, `dic.md` y la salida del verificador.

## Verificación
- (a) `sqlite3 -bail datos/tmp/rev.db` con el Up de 00001 y 00002: rc 0 y 0 · integrity ok · foreign_key_check 0 · 55/55 STRICT · 152 índices · 71 disparadores.
- (b) `sqlite_master` y `.dump` completo antes de 00002 contra después de su Down: `cmp` idénticos (1263 líneas; 102,696 bytes). Down de 00001: 0 objetos.
- (c) goose v3.28.0 sobre `datos/tmp/rev-s1-r2/g.db`:
  - up, version, down-to 0, version, up, version y status, todos con rc 0;
  - versiones 2, 0 y 2;
  - el esquema queda igual al de sqlite3, quitando las tablas propias de goose.
- (d) `strip.py` (quita comentarios respetando literales): tokens de HEAD 8599 = actual 8599.
- (e) `GOOSE=… python3 -I backend/db/tools/s1_verifica.py` → rc 0, RESUMEN TODO OK. KUESKI 389500 y 239528. NEG 46/46 y NEG12 26/26 rechazadas, con 12/12 positivas. S1-12 0 diferencias en 164 valores.
- (f) `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev-s1-r2/dic.md` + `cmp` → idéntico (174,605 bytes).
- (g) `grep -n -E "K<n>([^0-9]|$)"` en los 4 artefactos: tabla en el informe.
- (h) Las 11 líneas PROVISIONAL de reglas citan 4 bloqueos abiertos; 0 en las migraciones; 1 en el diccionario (la definición).
- (i) `app_setting`: 3 claves nuevas con valores de K21, K11 y K6; las 3 claves de K16 sin valor.
- (j) `ej.py`: 0 citas «ruta = valor» distintas de `model.json`/`data.json`. Todas las cifras fuera del oráculo están rotuladas o son pasos intermedios, salvo 1,814.50 (cuenta citada, R2-02). Un `grep` de fechas en el oráculo encontró 0 apariciones de 2026-10-05, 2026-10-11, 2026-10-31 y 2026-12-16 (R2-01).
- Pruebas propias (`neg.py`):
  - `telegram_chat_candidate`: rechaza confirmado sin persona, chat_id con letras, chat_id repetido y confirmación sin «Z»;
  - `shopping_list`: rechaza martes, domingo y cambio a miércoles; acepta lunes con fin en domingo.
- S1-01 propio: 31/31 entidades con tabla y sección del diccionario.
- S1-14: 0 tokens con el patrón de Telegram en `backend/db`.

## Supuestos
- Las fechas de un ejemplo cuentan como «cifras» para S1-03. Me baso en el propio criterio, que llama «ejemplo ilustrativo, sin dato en el oráculo» al ejemplo de los periodos de K8. No cambia cifras ni esquema.

## Bloqueos
- Ninguno.

## Pendientes
- Ninguno de esta revisión. Los hallazgos para la corrección están en `REVISION-S1-r2.md`, con rutas: R2-01 (alta) y R2-02 a R2-09 (bajas).

## Solicitudes fuera de alcance
- Ninguna nueva. Sigue pendiente la solicitud 1 de `correccion-S1-r1.md`: actualizar la línea S1-03 de `S1-esquema-1.2.md` (R2-04).

## Paquetes instalados
- Ninguno. Usé goose v3.28.0 del scratchpad, sqlite3 3.53.4 y Python del sistema.

## Consumo aproximado
- Acciones de herramienta: 37.
- Leído: unos 210 KB (criterios 23, reportes 22, 00002 16, reglas ~40, diccionario ~20, salidas de verificación ~25, lo demás en `grep`).
- Escrito: unos 31 KB (informe 14, este reporte 5, scripts temporales 11, progreso 1).
