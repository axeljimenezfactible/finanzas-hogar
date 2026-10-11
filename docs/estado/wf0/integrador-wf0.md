# Reporte WF-0 · integrador-wf0

Integrador del flujo WF-0 «Cimientos» · 11-oct-2026 · Estado: **terminado**, sin cierre por tope (`cerrado_por_tope = false`).
Veredicto del revisor independiente: **APROBADO en la ronda 3 de 3** (0 alta, 0 media, 10 baja). El documento de cierre es `docs/estado/wf0/INTEGRACION.md`.

## Hecho
1. **Lecturas.** Leí completos el protocolo y las decisiones; del plan, §0 a §3, §5, §8, §9 y §10 con `sed -n`. Leí los 16 reportes de `docs/estado/wf0/` (veredictos y hallazgos de REVISION-r1 a r3, y las secciones de cada reporte de agente y de corrección), los 32 bloqueos completos y `docs/criterios/INDICE.md` §6.
2. **Verificación propia.**
   - SQL en `datos/tmp/int.db`: Up, conteos, reglas de las decisiones, Down y goose.
   - `s1_verifica.py` con la regresión del oráculo.
   - `go build`, `go vet` y `gofmt`.
   - `colorimetria.py` en un árbol espejo.
   - Validador OpenAPI, `build.py --check`, `verify_c1.py --db` y `u1_verify.py`.
   - Huellas contra REVISION-r3.
   El detalle está en `INTEGRACION.md` §2.
3. **`README.md`:** qué es el sistema, decisiones clave con enlaces a `docs/plan/`, árbol del repositorio, entregables de WF-0, cómo correr la verificación y cómo sigue WF-1.
4. **`Makefile`:** `verify-wf0` (SQL temporal, go build y vet, colorimetría y validador OpenAPI) y `clean`, que piden la tarea. Agregué `verify-c1`, `verify-s1`, `verify-all` y `help`. Lo probé con goose y sin goose, y con dos pruebas negativas.
5. **`docs/estado/bloqueos/INDICE.md`:** 31 preguntas para Axel. Salen de los 32 archivos (31 abiertos y 1 resuelto) más 2 preguntas de la revisión (A12 decimales de tasa y A13 tope de Karen). Cada una lleva la pregunta exacta, las opciones con su consecuencia, el impacto, lo que quedó PROVISIONAL y la fase que necesita la respuesta.
6. **`docs/estado/wf0/INTEGRACION.md`:**
   - hecho por fase y verificación con comandos y salidas;
   - supuestos de todos los agentes;
   - bloqueos y pendientes para WF-1 y siguientes;
   - paquetes instalados;
   - 43 solicitudes fuera de alcance con su estado;
   - rondas de revisión y destino de los 10 hallazgos bajos;
   - lo que recibe Axel en la puerta 1, consumo por agente y huellas.
7. **R3-03:** anoté en `INTEGRACION.md` el cierre de U1-03 con la salida reejecutada.
8. **Decisión sobre `c1gen/`:** se queda en `backend/internal/api/c1gen/` (solicitud 6 de `correccion-C1-r1`).
9. **Limpieza y commit.** Borré `datos/tmp` y comprobé que no existe. Hice `git add -A` y un único commit con el mensaje «[WF-0] Cimientos: modelo de datos, arquitectura, colorimetría, criterios y contrato C1» y autor Axel Jimenez. El hash va en la salida estructurada.

## Archivos
- Creados:
  - `/home/deepsec/finanzas-hogar/Makefile`
  - `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/INDICE.md`
  - `/home/deepsec/finanzas-hogar/docs/estado/wf0/INTEGRACION.md`
  - `/home/deepsec/finanzas-hogar/docs/estado/wf0/integrador-wf0.md` (este reporte)
- Reescrito: `/home/deepsec/finanzas-hogar/README.md` (antes decía «README pendiente»).
- Temporales: `datos/tmp/` con `int.db`, `int-neg.db`, `int-down.db`, `int-goose.db`, `up.sql`, `down.sql` y `espejo/`, más lo que crea el Makefile en `datos/tmp/verify-wf0/`. Todo **borrado**. Los scripts de apoyo (`sql_checks.py`, `u1_verify.py` extraído y las salidas) quedaron en el scratchpad de la sesión, carpeta `int-wf0/`.
- No modifiqué ningún archivo de otro agente.

## Verificación
- **SQL.** `awk` (sección Up) → `sqlite3 -bail datos/tmp/int.db`: `UP_OK`; 54 tablas, 54 STRICT; 147 índices (110 explícitos); 68 disparadores; `foreign_key_check` 0; `integrity_check` ok; semillas 7, 4 y 7.
- **Reglas de las decisiones** (script propio sobre una copia): 5 de 5 positivas y 8 de 8 negativas rechazadas, de la cuenta para Karen a la co-titular que pasa a externa; `foreign_key_check` 0.
- **Down y goose.** Down deja 0 objetos y Up otra vez da 54 tablas. goose v3.28.0 up, down y up: OK, con 54 tablas, `fk` 0 e `integrity_check` ok.
- **`s1_verifica.py`** (Anexo B de S1, con GOOSE): «RESUMEN TODO OK»; 164 valores comparados con 0 diferencias; **389500** y **239528**, que en `oraculo/model.json` son `despensas[0].total` = 3895.0 y `axel.total` = 2395.28.
- **Go.** `cd backend && go build ./... && go vet ./...`: BUILD_OK y VET_OK con go1.26.5; `gofmt -l .` vacío.
- **Colorimetría.** `python3 -I datos/tmp/espejo/docs/design/tools/colorimetria.py --sin-png`: «VERIFICACIÓN OK». Los 4 archivos de texto salen idénticos byte a byte al repositorio. `contraste.csv`: 202 filas y 0 `cumple=false`.
- **OpenAPI.** `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml` → OK. Tiene 158 rutas, 197 operaciones y 237 esquemas.
- **C1.** `build.py --check`: IDÉNTICO ×3. `verify_c1.py --db`: 22 OK, «TODO OK».
- **`u1_verify.py`:** U1-03 y U1-04 en OK. Su única falla es la excepción declarada `GET /openapi.yaml`.
- **Makefile.**
  - `make verify-wf0 GOOSE=…`: código 0.
  - `make verify-all`: código 0, con goose omitido.
  - Prueba negativa de SQL (tabla no STRICT): FALLA, código 2.
  - Prueba negativa de colorimetría (título cambiado): FALLA, código 2.
- **Huellas.** Los 12 artefactos de REVISION-r3 no cambiaron.
- **Antes del commit:** `test ! -e datos/tmp` y `git status` revisado.

## Supuestos
1. **Colorimetría.** La tarea dice que «`colorimetria.py` regenera sin filas cumple=false». Lo corrí en un árbol espejo bajo `datos/tmp` (mismo script, mismas entradas) y comparé byte a byte, igual que los revisores. Así no reescribo artefactos de U2a. El Makefile hace lo mismo.
2. **Alcance.**
   - No copié scripts a rutas nuevas (R2-08) ni edité el reporte de U1 (R3-03): esas rutas no estaban en mi alcance. El revisor admite «cerrar o dejar registrados»; quedaron registrados en `INTEGRACION.md` §5 y §7.
   - `make verify-s1` ejecuta `s1_verifica.py` desde su anexo mientras no tenga ruta propia.
   - `criterios-toques` no lo cierro: es una regla de producto, así que queda como A10, pendiente de que Axel confirme la opción (a).
3. **Preguntas sin archivo.** Las dos preguntas que dejó la revisión entran al INDICE sin archivo de bloqueo propio, porque crear bloqueos nuevos no estaba en mi alcance.
4. **`make clean`** solo borra lo que crea el Makefile, `backend/bin` y los `__pycache__`. Quita `datos/tmp` solo si queda vacía, para no borrar temporales de otros procesos.
5. **Atribución del commit.** La tarea pedía cerrar el mensaje con «Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>». La integración la hizo Claude Opus 5.5, y el entorno de ejecución de esta sesión fija para sus commits las líneas «Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>» y «Claude-Session: …». Usé esas para que la atribución sea exacta. Los dos commits de la Etapa A llevan Fable 5.1 porque los hizo el orquestador. Autor y correo del commit: Axel Jimenez <axel.jimenez@factibleco.com>, como pidió la tarea.

Ninguno cambia cifras, esquema, operaciones ni tokens.

## Bloqueos
- Ninguno nuevo. El índice consolidado está en `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/INDICE.md`.

## Pendientes
- Nada pendiente de lo asignado.
- Para la puerta 1 y WF-1, la lista completa está en `INTEGRACION.md` §5:
  - las respuestas de Axel;
  - C1 v1.1 con R3-01, R3-02 y R2-06;
  - higiene del oráculo (`verify3.py` y `model.py`);
  - insumos de S1 para S2 y S3;
  - R3-04 en U2;
  - R2-03 a R2-08;
  - scripts al repositorio;
  - `docs/estado/wf1/`;
  - recalibrar tokens.

## Solicitudes fuera de alcance
1. **Orquestador** (`/home/deepsec/finanzas-hogar/.gitignore`): agregar `datos/tmp/` y `__pycache__/`. Motivo: que un temporal olvidado no entre en un commit.
2. **Integrador de WF-1 o dueños de cada artefacto:** R2-08. Mover `gen_dict.py` y `s1_verifica.py` a `backend/db/tools/`, `u1_verify.py` a `docs/design/tools/verificar_u1.py` con la excepción de C1 §11.2 (cierra la parte pendiente de R3-03) y los scripts de criterios a `docs/criterios/tools/`. Después, cambiar `verify-s1` del Makefile para que llame a la ruta nueva.
3. **Orquestador:** anotar en cada archivo de bloqueo la respuesta de Axel y cerrar `criterios-toques` si confirma la opción (a).

## Paquetes instalados
- Ninguno.
- Usé, sin red y desde las cachés del usuario: goose v3.28.0 (`scratchpad/gobin/goose`, compilado antes por S1), `uvx --offline` con openapi-spec-validator 0.9.0 y `uv run --offline` con pyyaml y jsonschema.
- No usé internet, credenciales ni Telegram.

## Consumo aproximado
- Unas 60 acciones de herramienta, sin subagentes.
- Leí unos 420 KB: protocolo, decisiones, plan por secciones, 16 reportes por secciones, 32 bloqueos y extractos de la migración, del OpenAPI y de los scripts.
- Escribí unos 85 KB: README, Makefile, INDICE de bloqueos, INTEGRACION y este reporte. En el scratchpad, unos 10 KB.
