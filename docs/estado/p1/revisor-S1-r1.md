# Reporte de agente · revisor-S1-r1 (flujo P1, ronda 1 de 3)
Fecha: 11-oct-2026 · Revisión independiente del esquema 1.2 de S1 contra `docs/criterios/S1.md` 1.1 y `docs/plan/decisiones.md` §K. No corregí nada: solo leí, reejecuté y escribí el informe.

## Hecho
- Leí `protocolo-agentes.md` completo, `decisiones.md` §K, `criterios/S1.md` 1.1, la migración `00002_puerta1.sql` completa, el diff de comentarios de `00001_esquema.sql`, lo agregado en el diff de `reglas-negocio.md` y, por grep y sed, las partes del diccionario 1.2 que tocan §K. También leí la respuesta de Axel de 9 bloqueos (los de S1 y C1 que se aplican aquí, más `S1-frescura-peso` y la opción D de `S1-meta-kcal`) y el formato de los 2 bloqueos nuevos de S1.
- Reejecuté las verificaciones (a) a (j) con scripts míos (`rev1.sh`, `rev2.py` y `rev3.sh` en `datos/tmp/rev-s1r1/`, ya borrados), corrí `s1_verifica.py` y `gen_dict.py`, y comparé con `git show HEAD`.
- Escribí el informe `docs/estado/p1/REVISION-S1-r1.md`: veredicto **NO APROBADO** con 1 hallazgo alto, 2 medios y 5 bajos, más el estado de los 6 hallazgos heredados de WF-0.

## Archivos
- Creados: `docs/estado/p1/REVISION-S1-r1.md` y `docs/estado/p1/revisor-S1-r1.md`.
- Temporales: `datos/tmp/rev-s1r1/` (bases, volcados y scripts), borrados al terminar. No toqué los temporales de otros agentes.

## Verificación
- (a) `sqlite3 rev.db < up 00001 && < up 00002`: integrity_check = ok · foreign_key_check = 0 filas · 55 de 55 tablas STRICT (`pragma_table_list`) · 69 disparadores.
- (b) Up+Down de 00002 contra 00001 sola: `sqlite_master` idéntico (1263 líneas, 111,282 bytes); `.dump` con datos idéntico al normalizar las marcas de tiempo.
- (c) `goose -dir backend/db/migrations sqlite3 g.db` up / version / down-to 0 / version / up / version: rc 0 en todos; versiones 2/0/2; 0 objetos tras down-to 0; esquema y datos iguales a los de sqlite3.
- (d) `git show HEAD:…/00001_esquema.sql` contra el archivo actual, sin comentarios: con `sed` y con un tokenizador que respeta las cadenas, 1302 = 1302, iguales; datos sembrados idénticos.
- (e) `GOOSE=… python3 -I backend/db/tools/s1_verifica.py`: rc 0, «RESUMEN TODO OK», KUESKI 389500 y 239528, NEG12 24/24 y 10/10. Mi suma directa de `model.json` da 389500, 239528 y 149972. Corrí 32 pruebas propias de los CHECK, la FK y los índices únicos nuevos: 32 de 32 como se espera.
- (f) `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev-s1r1/dic.md && cmp …`: idéntico (173,798 bytes).
- (g) grep de K1–K3, K6, K8, K9, K11–K16 y K18–K23 en los 4 artefactos: 18 de 18 presentes; K8 incompleta en `shopping_list.start_date` (H-02).
- (h) 11 líneas PROVISIONAL en las reglas, que citan 4 bloqueos, todos sin «Respuesta de Axel»; 0 en las migraciones; en el diccionario, 1 (la definición).
- (i) La 00002 agrega solo 3 filas a `app_setting`, todas con valores de §K (K21, K11, K6); `optimizer.*` sin las restricciones de K16.
- (j) 28 de 28 reglas con ejemplo; 62 citas resueltas y 17 de 17 «= valor» coinciden; RN-12 y RN-27 sin rótulo de ilustrativo y con cifras fuera del oráculo (H-01).

## Supuestos
- Cuento como rótulo aceptable «**Ejemplo ilustrativo** (motivo)» aunque no sea la leyenda literal del criterio; la diferencia queda como hallazgo bajo (H-08). Ninguno de mis supuestos cambia cifras, esquema ni endpoints.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ninguno de la revisión. Las correcciones de H-01 a H-08 le tocan al agente de corrección de S1; la ronda 2 debe reejecutar (a) a (j).

## Solicitudes fuera de alcance
1. `docs/contratos/reglas-negocio.md` L183, L270 y L173; `backend/db/tools/s1_verifica.py` (S1-03): rótulos y una comprobación más fuerte (H-01). Motivo: S1-03.
2. `backend/db/migrations/00001_esquema.sql` L1111 y `docs/contratos/diccionario-datos.md` L1773 (vía gen_dict): K8 en `shopping_list.start_date` (H-02). Motivo: K8.
3. `backend/db/migrations/00002_puerta1.sql` L31 (y L43): la necesidad de U1 correcta (H-03). Motivo: S1-04.

## Paquetes instalados
- Ninguno. Usé goose v3.28.0 del scratchpad del orquestador, sqlite3 3.53.4 y python3 del sistema.

## Consumo aproximado
- Acciones de herramienta: 32.
- Leído: unos 195 KB (protocolo 6, §K 8, criterios S1 22, migración 00002 15, reporte S1 11, diff de 00001 9, diff de reglas 29, partes del diccionario 20, bloqueos 8, salidas de verificación unos 60, U1 y otros greps 7).
- Escrito: unos 37 KB (scripts temporales 16, informe de revisión 15, este reporte 6).

cerrado_por_tope: false
