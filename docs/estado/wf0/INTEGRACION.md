# Integración WF-0 «Cimientos» · Finanzas Hogar

Integrador `integrador-wf0` · 11-oct-2026 · Etapa B, flujo WF-0 · Lo siguiente es la puerta 1 (C1 + U1 + criterios).

**Veredicto del revisor independiente: APROBADO en la ronda 3 de 3**, con 0 hallazgos de severidad alta, 0 de media y 10 de baja. Ninguno de los bajos bloquea la puerta 1 (`REVISION-r3.md`). El integrador repitió por su cuenta todas las verificaciones (§2) y los entregables siguen idénticos a los que aprobó el revisor (§11).

El flujo cierra con un solo commit: «[WF-0] Cimientos: modelo de datos, arquitectura, colorimetría, criterios y contrato C1». Su hash está en `git log` y en la salida del integrador, porque un archivo no puede citar el hash del commit que lo contiene.

## 1. Hecho por fase

### S1 · Dominio y modelo de datos
Agentes: `S1-modelo-datos` y `correccion-S1-r1`. Reportes: `S1-modelo-datos.md` y `correccion-S1-r1.md`.
- **Migración goose** `backend/db/migrations/00001_esquema.sql` (esquema 1.1). Tiene 54 tablas STRICT: las 31 entidades del plan §3.4.1, 14 de detalle y 9 complementarias. Suma 147 índices (110 explícitos) y 68 disparadores: 15 de reglas de las decisiones (Karen sin cuenta, ingresos solo del titular, salud y comidas sin externos, bitácora inmutable) y 53 de `updated_at`.
- **Convenciones.** El dinero va en enteros de centavos y las fechas en ISO con CHECK. Las tasas se guardan como fracción exacta (2 q = 53450/309450, 4 q = 80050/309450). Las semillas solo llevan valores decididos: 7 `app_setting`, 4 `meal_slot` y 7 `notification`.
- **Corrección r1.** Lista de compras (`shopping_list`, `shopping_list_item`, PROVISIONAL), descripción del gasto opcional, `recipe_meal_slot`, `inventory_movement`, `budget` con vigencia, `second_half_amount_cents`, `idempotency_key`, `category.kind`, límite de dos cuentas y disparador de participación.
- `docs/contratos/diccionario-datos.md` 1.1, generado desde la migración. `docs/contratos/reglas-negocio.md` trae 22 reglas, cada una con ejemplo del oráculo o marcada «ilustrativo», y 14 convenciones; su encabezado sigue en «Versión 1.0» (R2-04).
- **Hallazgos del oráculo**, documentados sin ajustar ninguna cifra:
  - conviven dos redondeos de serie: montos con base truncada y capital con base HALF_UP;
  - el «pago de referencia» no coincide con la serie real;
  - con la tasa redondeada a 25.8685 % no salen los mismos intereses: por eso se guarda la fracción;
  - P1–P6 dan exactamente los meses calendario de oct-2026 a mar-2027;
  - en las series de disponible, A − B = 3,094.50.
- 7 bloqueos.

### C1 · Contrato de datos y esqueleto del backend
Agentes: `S6a-contrato-c1`, `correccion-C1-r1` y `correccion-C1-r2`. Reportes: `S6a-contrato-c1.md`, `correccion-C1-r1.md` y `correccion-C1-r2.md`.
- **`docs/contratos/openapi.yaml`** (OpenAPI 3.1.1, C1 v1.0, `info.version` 1.0.0): 158 rutas, 197 operaciones y 237 esquemas.
  - Un solo esquema de error, paginación por cursor, 6 filtros globales y dinero `MoneyCents`.
  - `Idempotency-Key` en 41 operaciones de captura y cookie `fh_session`.
  - Solo login y salud son públicos.
- **`docs/contratos/C1-datos.md`**: convenciones, autenticación, errores, 17 métricas del dashboard con fórmula y ejemplo del oráculo, y la matriz pantalla → necesidad → operación (220 filas, 41 pantallas). Las cifras de ejemplo salen de `model.json` (16) o de `data.json` (2).
- **Backend Go que compila**: `backend/go.mod` (go 1.26), `cmd/finanzas/main.go` (`/health` y `/openapi.yaml` reales; 404 y 405 con el error único) y `internal/api/routes.go` (197 operaciones que responden 501).
- **Generador y verificador en el repositorio** (`backend/internal/api/c1gen/`): `build.py`, `build.py --check` y `verify_c1.py --db`, más `generate.go`.
- **Regeneración r2 contra el esquema 1.1.** Agregó 5 operaciones de la lista de compras (PROVISIONALES), presupuesto con vigencia (409 `budget_overlap`), `meal_slot_ids`, `InventoryMovement`, el monto del 15 y del 30, y la idempotencia en `x-tablas`.
- 6 bloqueos: 5 de C1 y `S6-version-go`.

### U1 · Arquitectura de información
Agente: `U1-arquitectura`. Reporte: `U1-arquitectura.md`.
- `docs/design/U1-arquitectura.md`:
  - barra lateral de 9 entradas;
  - 27 pantallas (P-00 a P-26) y 14 secciones de Ajustes (P-26.01 a P-26.14), con estados vacío, carga y error;
  - matriz de vistas para las 31 entidades;
  - 199 necesidades de datos (N-001 a N-315: 136 de lectura y 63 de escritura) y matriz pantalla ↔ necesidad (220 pares);
  - 8 flujos con pasos y conteo de toques (registrar gasto: 3 toques);
  - menú «/» con 27 comandos y 15 atajos (los de Notion verificados en su ayuda oficial);
  - reglas móvil primero (360, 412, 768, 1024 y 1440 px).
- 2 bloqueos.

### U2a · Colorimetría del matiz C (Índigo 280°)
Agentes: `U2a-colorimetria` y `correccion-U2a-r1`. Reportes: `U2a-colorimetria.md` y `correccion-U2a-r1.md`.
- `docs/design/tools/colorimetria.py` es determinista y usa solo la biblioteca estándar (OKLCH, WCAG 2.x y Machado 2009 en RGB lineal). Regenera:
  - `colorimetria.md`;
  - `contraste.csv` (202 pares, 0 bajo umbral);
  - `docs/contratos/C2-colores.json` (190 tokens DTCG);
  - `paleta.html` y `paleta.png`.
- Reproduce 22 de 22 valores de partida del plan §4.3.
- Ingreso propuesto en 135° y paleta categórica de personas y tiendas sin colisiones con los roles.
- **Corrección r1:** la serie 6 pasa de 100° (reservado a acción crítica) a 225°, con nuevo orden de series. La zona reservada 80°–120° queda verificada por script.
- 1 bloqueo.

### Criterios de aceptación
Agentes: `criterios` y `correccion-criterios-r1`. Reportes: `criterios.md` y `correccion-criterios-r1.md`.
- `docs/criterios/`: INDICE y 18 archivos (versión 1.1) con 248 criterios (189 de severidad alta y 59 de media), cada uno con verificación y evidencia.
- Las cifras canon R-01…R-85 y K-01…K-08 se generaron con asserts contra `model.json` y `data.json`.
- **Corrección r1:**
  - pantalla = `P-xx` y necesidad = `N-xxx`;
  - S1-09 pasa a meses calendario más la prueba de equivalencia con el oráculo;
  - 248 de 248 criterios trazados.
- 16 bloqueos; `criterios-nxx` quedó resuelto con (b).

### Revisión independiente
Hubo 3 rondas (detalle en §8). Las correcciones se asignaron por hallazgo y el revisor reejecutó todo con scripts propios en cada ronda.

### Integración (este agente)
- Verificación propia de todo lo anterior (§2).
- `README.md` y `Makefile` (`verify-wf0`, `verify-c1`, `verify-s1`, `verify-all`, `clean`).
- `docs/estado/bloqueos/INDICE.md`, que reúne 31 preguntas para Axel.
- Este documento, el reporte `integrador-wf0.md` y el commit único.

## 2. Verificación (reejecutada por el integrador el 11-oct-2026)

Todo se corrió desde `/home/deepsec/finanzas-hogar`. Las bases temporales vivieron en `datos/tmp/` (`int.db`, `int-neg.db`, `int-down.db` e `int-goose.db`) y se borraron al final; `datos/tmp` no existe al hacer el commit.

| # | Qué | Comando o método | Resultado |
|---|---|---|---|
| 1 | SQL en base temporal | Sección Up extraída con awk → `sqlite3 -bail datos/tmp/int.db`; consultas a `sqlite_master`, `pragma_table_list` y `pragma_foreign_key_check` | `UP_OK`: 54 tablas, las 54 STRICT; 147 índices (110 explícitos); 68 disparadores; `foreign_key_check` 0 filas; `integrity_check` ok. Semillas: 7 `app_setting`, 4 `meal_slot` y 7 `notification` |
| 2 | Reglas de las decisiones en la base | Script propio sobre una copia (`int-neg.db`) | 5 de 5 positivas aceptadas (titular, co-titular, Karen externa, feb-2027 y su quincena 2 que paga el 2027-02-28). 8 de 8 negativas rechazadas: cuenta para Karen, hash que no es argon2id, segunda co-titular, 12.5 en centavos, fecha `10-oct-2026`, ventana del 10 al 9 como mes, `pay_date` 2027-02-30 y Ximena con cuenta que pasa a externa. `foreign_key_check` 0 |
| 3 | Down y goose | Down sobre una copia; `goose v3.28.0 -dir backend/db/migrations sqlite3 datos/tmp/int-goose.db up/down/up` | Down deja 0 objetos y Up otra vez da 54 tablas con `integrity_check` ok. goose: «OK 00001_esquema.sql»; tras el down quedan solo `goose_db_version` y `sqlite_sequence`; tras el segundo up hay 54 tablas, `fk` 0 e `integrity_check` ok |
| 4 | Regresión de S1 contra el oráculo | `s1_verifica.py` extraído del Anexo B de `S1-modelo-datos.md` (`GOOSE=…/gobin/goose`) | «RESUMEN TODO OK»: S1-01 a S1-17 en OK; 46 de 46 negativas y 8 de 8 positivas; Down, Up y goose OK. Ida y vuelta: 164 valores, 0 diferencias. **Suma de pagos de la despensa 1 = 389500 y parte de Axel = 239528**, iguales a `model.json` (`despensas[0].total` 3895.0 y `axel.total` 2395.28, comprobados aparte) |
| 5 | Go | `cd backend && go build ./... && go vet ./...`; `gofmt -l .`; `go list ./...` | BUILD_OK y VET_OK con go1.26.5; `gofmt` sin salida; 2 paquetes; `go.mod` = `go 1.26`; no se escribió ningún binario |
| 6 | Colorimetría | `colorimetria.py --sin-png` corrido en un árbol espejo (`datos/tmp/espejo/`), para no reescribir el repositorio; luego `cmp` y conteo del CSV | «VERIFICACIÓN OK». `colorimetria.md`, `contraste.csv`, `C2-colores.json` y `paleta.html` **idénticos byte a byte** al repositorio. `contraste.csv`: 202 filas, **0 `cumple=false`**, 0 bajo umbral (mínimo de texto 4.73, de UI 3.05). 190 tokens = 190 filas del Anexo A; plan §4.3 22/22; zona 80°–120° sin tokens ajenos. `paleta.png` sin cambios (sha256 `7746da788ff19cbb`) |
| 7 | Validador OpenAPI | `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml`; carga con PyYAML | `docs/contratos/openapi.yaml: OK`, código 0. 3.1.1 · 1.0.0 · 158 rutas · 197 operaciones · 237 esquemas |
| 8 | Generador de C1 | `python3 -I backend/internal/api/c1gen/build.py --check` | IDÉNTICO en `openapi.yaml`, `routes.go` y `C1-datos.md`, código 0. Sin endpoint solo N-003 y N-301 (locales); sin necesidad solo `getOpenApiDocument`; 0 tablas sin operación |
| 9 | Verificador de C1 | `uv run --offline --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/int-goose.db` | 22 OK, 0 FALLA, «RESULTADO: TODO OK» (C1-01 a C1-15 y S6-03: 197 rutas = 197) |
| 10 | U1-03 con el script de U1 | `u1_verify.py` extraído del anexo de `U1-arquitectura.md` (sha256 `71e0d6e84453cf60`) | `OK U1-03 N-xxx sin endpoint en C1: []` y `OK U1-04 … = 3`. Su única «FALLA» es `GET /openapi.yaml`, la excepción «estáticos» que admiten C1-02 y U1-03 (C1 §11.2). **U1-03 cumple** (ver §8, R3-03) |
| 11 | Makefile | `make verify-wf0 GOOSE=…/gobin/goose`, `make verify-all` (sin goose en el PATH) y dos pruebas negativas | Las dos corridas terminan con código 0 y los mismos resultados de las filas 1, 3 a 9. Pruebas negativas: una migración con una tabla no STRICT da «FALLA: el esquema no cumple STRICT…» con código 2; un script de colorimetría con el título cambiado da «FALLA, difiere del regenerado: docs/design/colorimetria.md» con código 2 |
| 12 | Huellas | `sha256sum` | Los 12 artefactos que lista `REVISION-r3.md` tienen la misma huella (§11): nada cambió después de la aprobación |
| 13 | Limpieza | `make clean`, `rm -rf datos/tmp`, `find … __pycache__` | `datos/` queda solo con `semilla/`; no hay `__pycache__` ni `backend/bin` |

**Hallazgo menor del integrador**, solo de presentación y sin cifra canon: R2-06 dice que con 4 decimales la tasa de 2 quincenas «mostraría 17.2725 %». Ese valor es el truncado. Con ROUND_HALF_UP (protocolo §6), 53450/309450 = 17.272580…% da **17.2726 %**. La pregunta A12 del INDICE de bloqueos usa el valor correcto.

## 3. Supuestos de todos los agentes

Según sus reportes y las tres revisiones, ninguno cambia cifras, esquema, operaciones ni tokens. El único que la revisión marcó como no reflejado en otro artefacto es el supuesto 4 de U1 (R3-01).

- **S1-modelo-datos:**
  - las convenciones técnicas están en `reglas-negocio.md` §Convenciones;
  - lecturas derivadas con su fuente: la regla por defecto es «Hogar» (decisión A); el presupuesto mide la parte del hogar (decisiones A y B, plan §1.3); la fecha ancla de la lectura A es el 2026-10-10 (`model.json: credito.eventos[0]`); no hay corrimiento por fin de semana (el oráculo conserva el 15-nov-2026, que es domingo).
- **correccion-S1-r1:**
  - la clave de idempotencia usa el formato de C1 §1 punto 8;
  - «tipo (fijo o variable)» de la categoría sale de U1 N-281, y los tipos y motivos de inventario de U1 N-181 y N-183;
  - la semana de la lista de compras dura 7 días y su inicio lo fija U1 (bloqueo A8).
- **S6a-contrato-c1:**
  - prefijo `/api/v1`, cookie `fh_session`, códigos de error en inglés y `operationId` como identificadores del contrato;
  - «colección» = lista paginada; los agregados no se paginan;
  - el estado de bloqueo del login viaja en la respuesta 429, para no revelar usuarios;
  - `GET /openapi.yaml` exige sesión en el contrato, pero el esqueleto la sirve sin sesión hasta S6;
  - el esqueleto escucha solo en 127.0.0.1:8080 hasta tener TLS y la respuesta D2.
- **correccion-C1-r1:**
  - el generador va en `backend/internal/api/c1gen/`;
  - C1 quedó con el esquema de r1 en esa ronda (se resolvió en r2);
  - C1-02 sigue el texto vigente del criterio;
  - los intermedios van a `datos/tmp/c1gen-*`.
- **correccion-C1-r2:**
  - el contrato sigue llamándose C1 v1.0, con una fila «corrección r2»;
  - forma de la API de la lista de compras: semana ISO por `period`, lista vacía si no existe, agregar un artículo crea la lista, asignar tienda responde 404 sin lista y no hay operación para quitar un artículo (cubierto por A1);
  - código 409 `budget_overlap`;
  - la existencia inicial del inventario se registra como movimiento `adjustment` con motivo `count`;
  - `idempotency_key` aparece en `x-tablas`.
- **U1-arquitectura:**
  1. Conteo de toques con la regla (a) (A10).
  2. Preferencias de vista, filtros, tema y barra lateral por dispositivo (PROVISIONAL).
  3. Barra inferior móvil con 5 destinos.
  4. Los domingos, el plan de comidas y la lista de compras abren la semana siguiente; C1 no lo refleja todavía (R3-01).
  5. La cuenta sugerida por renglón del OCR es la de la última compra confirmada.
- **U2a-colorimetria:**
  - «fondo 50 / tarjeta 100» y «950 / 900» son pasos del neutro tintado;
  - el «texto 950» del botón primario en oscuro es brand 950;
  - ΔL ≥ 0.10 se aplicó a todos los pares de cada grupo (la lectura más estricta) y también frente a acción crítica;
  - los nombres de los tokens van en inglés; las personas usan `person.code` y `store` no tiene código estable.
- **correccion-U2a-r1:**
  - la zona reservada mide 20° a cada lado de 100° (80°–120°), y no puede ser más ancha porque advertencia está en 75°;
  - el matiz de la serie 6 y el orden de las series son trabajo de diseño sujeto a la puerta 2.
- **criterios:**
  - el formato de evidencia y las severidades media y baja son convenciones del directorio;
  - cada agente fija los nombres de sus pruebas y scripts;
  - medidas propias: tolerancia de ±2° al medir matiz (U2-02) y «≥ 2 meses» para «sesión de meses» (S6-07);
  - el supuesto «N-xx = pantalla» se corrigió en r1.
- **correccion-criterios-r1:**
  - agregó la resolución a `criterios-nxx`;
  - T = pantallas P-xx más las 14 secciones;
  - P-26 se cubre con sus secciones;
  - en fixtures con eventos entre los días 1 y 9, el esperado sale del oráculo con las ventanas de la fixture;
  - versión 1.1 y etiqueta «convención propia del directorio»;
  - usó un subdirectorio propio en `datos/tmp`.
- **Revisores r1, r2 y r3:**
  - r1 corrió la colorimetría en una copia, no en el repositorio, y fijó las severidades de H-01 y H-03 con su justificación;
  - r2 evaluó C1 tal como estaba publicado y usó el goose del scratchpad;
  - r3 calificó como baja todo hallazgo que no cambia cifras, esquema, rutas ni contratos, y consideró que las propiedades anidadas de las entradas cumplen C1-03.
- **integrador-wf0:**
  1. Verificó la colorimetría en un árbol espejo, igual que los revisores.
  2. No editó archivos fuera de su alcance: el cierre de U1-03 se anota aquí y no en el reporte de U1.
  3. No cierra `criterios-toques` sin la confirmación de Axel (A10).
  4. `make clean` solo borra lo que crea el Makefile.
  5. La línea de coautoría del commit nombra al modelo que hizo la integración (detalle en su reporte).

## 4. Bloqueos

El detalle está en [`docs/estado/bloqueos/INDICE.md`](../bloqueos/INDICE.md).
- **Archivos:** 32 en `docs/estado/bloqueos/`, todos con el formato del protocolo: 16 de criterios (1 resuelto), 7 de S1, 5 de C1, 2 de U1, 1 de U2 y 1 de S6.
- **Preguntas:** con las 2 que dejó la revisión (A12 decimales de tasa y A13 tope de Karen) son 31 preguntas para Axel:
  - 13 fijan C1, U1 y el esquema (puerta 1);
  - 12 las necesitan los motores de WF-1;
  - 2 son de diseño (puerta 2);
  - 4 son de construcción y operación (WF-3 y WF-4).
- **Una respuesta cierra dos archivos** en 4 casos: lista de compras, huecos de negocio, frescura de precios y meta de calorías.
- Este integrador no creó bloqueos nuevos.

## 5. Pendientes para WF-1 (y los flujos siguientes)

### Antes de arrancar WF-1 (orquestador, en la puerta 1)
1. Llevar a Axel las 31 preguntas del INDICE de bloqueos y anotar sus respuestas en cada archivo.
2. Con las respuestas, regenerar C1 v1.1 con `backend/internal/api/c1gen/` (`build.py`, `build.py --check`, `verify_c1.py --db`). En esa regeneración entran:
   - R3-01: la semana por omisión de la lista de compras, alineada con U1 §0.2;
   - R3-02: los cuatro límites técnicos, según A6;
   - R2-06: el «Z» solo en entradas y los decimales de tasa, según A12.
   Si alguna respuesta cambia el esquema (A1 con B o C, A8 con D, A9 con 2 o 3, B3 con B o C, B10 con D, A3 con B), primero va la migración 00002 de S1.
3. Crear `docs/estado/wf1/` (solicitud 5 de `criterios`).
4. Recalcular el presupuesto de tokens de WF-1 con el consumo real de WF-0 (§10; plan §5 pide calibrar con WF-0). Cuatro agentes reportaron más de 300K tokens de contexto acumulado sin comprimir.

### WF-1 · S2 y S3 (motores financiero y de créditos)
- **Insumos obligatorios de S1:**
  - tasas como fracción `rate_num`/`rate_den`;
  - series guardadas ya calculadas (`scheduled_payment`, `scheduled_payment_share`);
  - montos con base truncada y capital con base HALF_UP para reproducir el oráculo (RN-07);
  - la regla por defecto «Hogar» y la fecha ancla de la lectura A (2026-10-10).
- **Higiene del oráculo antes de extenderlo** (solicitudes 1 y 2 de `criterios`; S2-04 y S3-16): `oraculo/verify3.py` línea 89 escribe en `~/reporte-compra-latex/model/verify3.json`, fuera del repositorio, y `oraculo/model.py` regenera `model.json` al importarse. Mientras no se corrija, nadie debe ejecutar esos dos archivos.
- **Fixtures de S2:** declarar las ventanas de cada fixture y no reutilizar R-45…R-50 como esperado cuando haya eventos entre los días 1 y 9 (solicitud 4 de `correccion-criterios-r1`).
- Respuestas que necesitan: B1 a B8 y el punto 7 de D4.

### WF-1 · S4 y S5 (recetas, comidas y salud; precios y tiendas)
- Respuestas B9 a B12 y A7.
- **S5:** la lectura asistida de sitios sigue desactivada hasta que exista `docs/estado/S5-terminos.md` con veredicto (protocolo §7).
- **S1:** campo `store.code` (a-z, 0-9 y guion bajo), para mapear cada tienda a `color.store.<código>` sin depender del nombre (solicitud de U2a).

### WF-1 · U2 (DESIGN.md y tokens C2), U3 y U4
- Integrar `docs/contratos/C2-colores.json` sin cambiar valores. Ojo: `color.chart.3`, `.5` y `.6` cambiaron en la corrección r1.
- **R3-04:** sumar la marca al barrido de pares cercanos de `colorimetria.md` §7.7 (Walmart 262° queda a 18° de la marca 700 con ΔL 0.000) y declarar el refuerzo, o separar Walmart.
- Simular daltonismo con culori convirtiendo antes a `lrgb`, y mantener `.light` y `.dark` por tema también en series, personas y tiendas.
- No copiar del plan §4.3 «bordes de campos neutro 600 ≥ 4.2:1» (se mide 3.71:1, que sí cumple 3:1) ni los roles a L 0.56.
- **Puerta 2:** C1 y C2 del INDICE, ingreso en 135°, colisiones categóricas resueltas y la paleta de gráficas corregida.
- **U4:** valores por defecto que no están decididos: hora del resumen diario, anticipación de alertas, holgura del comparador y precisión del peso (U1, pendiente 3). Tienen relación con D4.

### Criterios (dueño de `docs/criterios/` en WF-1)
- **R2-07:**
  - U2-02: admitir 135° como propuesta;
  - U1-05: pedir 6 flujos o más, con los 4 obligatorios;
  - C1-01: poner primero el comando con `uvx`;
  - S1-03: permitir ejemplos «ilustrativos» marcados;
  - cambiar `/tmp` por `datos/tmp` en los «Verifica»;
  - C1-02, C1-03 y S6-03: citar `verify_c1.py --db` y `build.py --check`;
  - volver literales las citas parafraseadas.
- **U2-02 y U2-09** (solicitud de `correccion-U2a-r1`): fijar 20° como ancho de la zona reservada, medido en el hex, y exigir la prueba de todos los pares en las tres primeras series.
- Actualizar los criterios afectados por cada respuesta (columna «Criterios afectados» de `docs/criterios/INDICE.md` §6).

### S1 (la próxima vez que se toque el esquema)
- **R2-03:** redacción neutral de `person.fortnight_payment_cap_cents`; la pregunta va como A13.
- **R2-04:** `reglas-negocio.md` pasa a versión 1.1 con su lista de cambios.
- Al cambiar la migración, regenerar el diccionario con `gen_dict.py` (Anexo A de `S1-modelo-datos.md`) y correr `make verify-s1`.

### U1
- **R2-05:** 44 px como mínimo táctil; 24 px solo como piso de WCAG 2.5.8.
- **Solicitud 3 de S6a:**
  - alinear «familia» de §0.2 con N-040 según A5;
  - remitir la columna «Endpoint C1» de §3.29 a `C1-datos.md` §10;
  - N-011 se cubre con la respuesta 429 del login.
- Opcional, según A1: la necesidad «Quitar artículo» de la lista. C1 agregaría `DELETE /api/v1/shopping-list/items/{shopping_list_item_id}`.

### Reproducibilidad (R2-08; dueños de cada artefacto o integrador de WF-1)
- Llevar a rutas del repositorio los scripts que hoy viven solo en anexos de reportes o en el scratchpad de la sesión, que se borra:
  - `gen_dict.py` y `s1_verifica.py` → `backend/db/tools/`;
  - `u1_verify.py` → `docs/design/tools/verificar_u1.py`, aceptando la excepción de C1 §11.2;
  - los scripts de criterios (`traza.py`, `refs.py`, `citas.py`, `canon.py`, `s109_sql.sh`, `s109_equiv.py`, `matriz_c102.py`, `verif_final.sh`) → `docs/criterios/tools/`;
  - opcional: `verif_indep.py` de U2a → `docs/design/tools/`.
- Mientras tanto, `make verify-s1` ejecuta `s1_verifica.py` desde su anexo.
- `gen_canon.py` de criterios solo existió en el scratchpad. Sus 85 filas R y 8 filas K están en los criterios, y `canon.py` (anexo de `correccion-criterios-r1.md`) las reverifica.
- Agregar `datos/tmp/` y `__pycache__/` al `.gitignore`.

### WF-3 (S6, S7, F1 a F4, U5)
- **S6:**
  - DSN de modernc con `_pragma=foreign_keys(1)`, `journal_mode(WAL)` y `busy_timeout`;
  - `audit_log` e `inventory_movement` en la misma transacción que el cambio;
  - purga de `idempotency_key` según A6;
  - validar los invariantes 7 (vigencias de `budget` sin traslape) y 13 (regenerar la lista conserva lo elegido);
  - exigir sesión en `GET /openapi.yaml`;
  - implementar los 197 manejadores con una prueba por operación;
  - sesiones visibles y cerrables por ambos co-titulares (N-271);
  - certificado raíz de la CA local descargable (N-274);
  - bloqueo de login según A11.
- **S7:** registro del chat (A4) y claves de `notification.params_json` (D4).
- **U5:** comprobar en Android a 412 × 915 que el teclado numérico abre con el primer toque de F-1 y que la hoja cabe sin desplazar a 360 px.
- Respuestas D1 a D4 y A9.

### WF-4 y acciones de Axel
- `datos/semilla/gastos_fijos_ingresos.csv`, opcional; el plan arranca sin él.
- Confirmar con Kueski cómo repone el capital: es el interruptor A/B, y B2 opción D depende de eso.
- Instalar la CA de mkcert en los dos Android.
- Ximena escribe una vez al bot.

## 6. Paquetes instalados (consolidado)

- **Nada con pacman** ni en el Python del sistema. Go 1.27 no se instaló (bloqueo D1).
- **goose v3.28.0:** lo compiló S1 con `go install … github.com/pressly/goose/v3/cmd/goose@v3.28.0` en el scratchpad de la sesión, fuera del PATH y del repositorio. Como efecto lateral llenó `~/go/pkg/mod` (pressly/goose, modernc.org/sqlite, modernc.org/libc, modernc.org/memory, golang.org/x/sys, golang.org/x/sync y sus dependencias). Lo usaron después las correcciones, los revisores y el integrador, sin reinstalarlo.
- **openapi-spec-validator 0.9.0** con `uvx`, en la caché de uv. Lo usaron C1, las correcciones, los revisores y el integrador (este último con `--offline`).
- **pyyaml 6.0.3 y jsonschema 4.26.0** con `uv run --with`, en entornos efímeros de la caché de uv (`verify_c1.py`).
- **@redocly/cli 2.62.1** con `npx -y`, en la caché de npm (C1).
- **culori 4.0.2** (`npm pack`) y **style-dictionary 5.6.1** (`npm install --ignore-scripts`, 106 paquetes): solo en el scratchpad, para contraverificar U2a.
- **Red:** registro de npm (U2a) y documentación oficial leída con Obscura: ayuda de atajos de Notion (U1), tabla de matrices de Machado et al. y documentación de `simulate_cvd` del paquete R colorspace (U2a). Nadie usó credenciales ni mandó Telegram.
- **Integrador:** ninguno; todo salió de las cachés, sin red.

## 7. Solicitudes fuera de alcance: resueltas y pendientes

Estados: **Resuelta** (por quién), **Decidida** (decisión del integrador, sin escribir fuera de su alcance) y **Pendiente** (dueño y momento).

| # | Origen → destino | Pedido | Estado |
|---|---|---|---|
| 1 | S1 → integrador (H-05, R2-08) | Copiar `gen_dict.py` y `s1_verifica.py` a `backend/db/tools/` | **Pendiente** (WF-1). Paliativo: `make verify-s1` |
| 2 | S1 → C1 | Regenerar C1 con el esquema 1.1 | **Resuelta** por `correccion-C1-r2` |
| 3 | S1 → S6 | Pragmas en el DSN y `audit_log` en la misma transacción | **Pendiente** (WF-3, S6) |
| 4 | S1 → S2 y S3 | Tasas como fracción, series guardadas, redondeos del oráculo | **Pendiente** (WF-1; insumo obligatorio, §5) |
| 5 | correccion-S1-r1 → C1 (1 y 2) | Operaciones de la lista de compras y columnas nuevas; «Cómo seguí» de los bloqueos de C1 | **Resuelta** por `correccion-C1-r2` |
| 6 | correccion-S1-r1 → S6 (4) | Purga de `idempotency_key`, movimientos de inventario en la misma transacción, invariantes 7 y 13 | **Pendiente** (WF-3) |
| 7 | S6a → integrador (1) | Generador de C1 en el repositorio | **Resuelta** por `correccion-C1-r1` (`backend/internal/api/c1gen/`) |
| 8 | S6a → S1 (2) | Migración con los huecos de U1 | **Resuelta** en el esquema 1.1, salvo A2 a A4 |
| 9 | S6a → U1 (3) | Familia en §0.2, columna de §3.29, N-011 con 429 y excepciones de `u1_verify.py` | **Pendiente** (U1, tras A5) |
| 10 | S6a → criterios (4) | C1-01 con `uvx`; N-xxx en C1-02 y C1-10; nombres de los filtros en C1-06 | N-xxx **resuelta** (H-06); el resto **pendiente** (R2-07) |
| 11 | S6a → criterios (5) | S6-01 espera Go 1.27 | **Pendiente** (depende de D1) |
| 12 | S6a → S7 (6) | Registro del chat y claves de `notification.params_json` | **Pendiente** (A4 y D4; WF-3) |
| 13 | correccion-C1-r1 → integrador (1 a 4) | Scripts de S1, `u1_verify.py`, `gen_canon.py` y scripts del revisor al repositorio | **Pendiente** (R2-08, §5) |
| 14 | correccion-C1-r1 (5) y correccion-C1-r2 (4) → criterios | Citar `verify_c1.py --db` y `build.py --check` en C1-02, C1-03 y S6-03 | **Pendiente** (R2-07) |
| 15 | correccion-C1-r1 → integrador (6) | Dejar o mover `c1gen/` | **Decidida**: se queda en `backend/internal/api/c1gen/`, junto a `routes.go` y `go generate` |
| 16 | correccion-C1-r2 (1) y revisor r3 (1) → integrador o U1 | Anotar el cierre de U1-03 y ajustar `u1_verify.py` | **Decidida** en parte: el cierre queda anotado aquí (§2, fila 10; §8). El reporte de U1 se conserva como registro de su entrega. El ajuste del script queda **pendiente** con R2-08 |
| 17 | correccion-C1-r2 (2) → S1 | R2-03, redacción del tope por quincena | **Pendiente** (S1; pregunta A13) |
| 18 | correccion-C1-r2 (3) → U1 | Necesidad opcional «Quitar artículo» | **Pendiente**, opcional (según A1) |
| 19 | U1 → criterios (1) | (a) N-xx = necesidad; (b) U1-05 con 6 + 2 flujos | (a) **resuelta** (H-06); (b) **pendiente** (R2-07) |
| 20 | U1 (2) y correccion-criterios-r1 (2) → integrador | Cerrar `criterios-toques` con (a) | **Decidida**: no se cierra sin Axel; queda como A10, pendiente de su confirmación |
| 21 | U1 → S1 (3) | Necesidades sin entidad explícita | **Resuelta** en el esquema 1.1 (r3: C1-03 54 de 54, U1-12 cumple) |
| 22 | U1 → S1 y S2 (4) | Equivalencia P1–P6 = meses calendario | **Resuelta** en S1 (S1-09, 12 de 12); para S2 queda la nota de S2-09 (WF-1) |
| 23 | U1 → C1 (5) | Matriz con una fila por necesidad | **Resuelta** (C1 §10, 220 filas) |
| 24 | U1 → S6 (6) | Sesiones de ambos, certificado raíz (N-274), bloqueo de login | **Pendiente** (WF-3) |
| 25 | U1 → integrador (7) | Guardar `u1_verify.py` en el repositorio | **Pendiente** (R2-08) |
| 26 | U2a → criterios | U2-02 con 135° y U2-08 con el refuerzo declarado | **Pendiente** (R2-07) |
| 27 | U2a → S1 | `store.code` | **Pendiente** (S1, WF-1) |
| 28 | U2a → U2b | culori en `lrgb`; `.light` y `.dark` por tema | **Pendiente** (WF-1) |
| 29 | U2a → DESIGN.md | No copiar «neutro 600 ≥ 4.2:1» ni roles a L 0.56 | **Pendiente** (WF-1) |
| 30 | correccion-U2a-r1 → criterios | Ancho de 20° (U2-02) y prueba de todos los pares (U2-09) | **Pendiente** |
| 31 | correccion-U2a-r1 → integrador (opcional) | Copiar `verif_indep.py` a `docs/design/tools/` | **Pendiente**, opcional (R2-08) |
| 32 | criterios → oráculo (1 y 2) | `verify3.py` escribe fuera del repositorio; `model.py` regenera `model.json` al importarse | **Pendiente** (WF-1, antes de extender el oráculo) |
| 33 | criterios → plan (3 y 4) | Invariante de deuda viva (B6) y «neutro 600 ≥ 4.2:1» | **Pendiente** (orquestador tras las respuestas; `plan.md` lo aprobó Axel) |
| 34 | criterios → orquestador (5) | Crear `docs/estado/wf1`, `wf3` y `wf4` | **Pendiente** (al arrancar cada flujo) |
| 35 | correccion-criterios-r1 → integrador (1) | Scripts de criterios a `docs/criterios/tools/` | **Pendiente** (R2-08) |
| 36 | correccion-criterios-r1 → corrección de S1 (3) | S1-09 con el texto nuevo | **Resuelta** (S1-09 en OK: 18 filas, rechazo de la ventana del 10 al 9, 12 de 12) |
| 37 | correccion-criterios-r1 → WF-1 (4) | Fixtures de S2 con ventanas propias | **Pendiente** (WF-1, S2) |
| 38 | revisor r1 → integrador | Generadores y verificadores en el repositorio (H-05) | **Resuelta** para C1; el resto **pendiente** (R2-08) |
| 39 | revisor r2 → orquestador | Asignar a C1 la regeneración (R2-01, R2-02) | **Resuelta** (`correccion-C1-r2`; ambos cerrados en r3) |
| 40 | revisor r3 → C1 (2) | R3-01, R3-02 y R2-06 en la próxima regeneración | **Pendiente** (C1 v1.1 tras la puerta 1; dependen de A6 y A12) |
| 41 | revisor r3 → U2 (3) | R3-04 | **Pendiente** (WF-1) |
| 42 | revisor r3 → S1 y criterios (4) | R2-03, R2-04, R2-07 y R2-08 | **Pendiente** (§5) |
| 43 | integrador → orquestador | Agregar `datos/tmp/` y `__pycache__/` al `.gitignore` | **Pendiente**: `.gitignore` no está en el alcance de este integrador |

## 8. Revisión independiente: veredicto y rondas

| Ronda | Fecha | Veredicto | Alta | Media | Baja | Qué bloqueaba | Correcciones asignadas |
|---|---|---|---|---|---|---|---|
| r1 | 10-oct-2026 | NO APROBADO | 1 | 7 | 4 | H-01: la lista de compras sin tabla en S1 ni operación en C1 (N-151 a N-153, N-165) | `correccion-S1-r1` (H-01, H-02, H-04, H-11), `correccion-C1-r1` (H-05), `correccion-U2a-r1` (H-03), `correccion-criterios-r1` (H-06, H-07, H-08) |
| r2 | 11-oct-2026 | NO APROBADO | 2 | 0 | 6 | R2-01 y R2-02: C1 no se había regenerado contra el esquema 1.1 | `correccion-C1-r2` (R2-01, R2-02) |
| r3 | 11-oct-2026 | **APROBADO** | 0 | 0 | 10 | — | ninguna; los bajos se registran |

Los 10 hallazgos bajos que siguen abiertos no cambian cifras, esquema ni rutas. Su destino:

| Id | Artefacto | Qué es | Destino |
|---|---|---|---|
| R3-01 | C1 | La semana por omisión de la lista de compras no sigue la regla del domingo de U1 | C1 v1.1 (§5) |
| R3-02 | C1 y S1 | 4 límites técnicos sin marca PROVISIONAL | A6 del INDICE y C1 v1.1 |
| R3-03 | U1 | El reporte de U1 dice «U1-03 PENDIENTE» y su script marca como falla la excepción declarada | **Anotado aquí**: U1-03 cumple (reejecución del integrador: `OK U1-03 N-xxx sin endpoint en C1: []`; la única falla del script es `GET /openapi.yaml`, la excepción de C1 §11.2). El ajuste del script va con R2-08 |
| R3-04 | U2a | La marca no entra en el barrido de pares cercanos (Walmart a 18°) | WF-1, U2 |
| R2-03 | S1 | La descripción del tope de Karen presenta una interpretación como hecho | S1 y A13 |
| R2-04 | S1 | `reglas-negocio.md` sigue en versión 1.0 | S1 |
| R2-05 | U1 | Mínimo táctil de 24 px en lugar de 44 px | U1 |
| R2-06 | C1 | `DateTime` admite «Z»; decimales de `rate_pct` | A12 y C1 v1.1 |
| R2-07 | criterios | Incoherencias menores de redacción y de comandos | Criterios en WF-1 |
| R2-08 | S1, U1 y criterios | Scripts solo en anexos | Reproducibilidad (§5) |

## 9. Puerta 1: qué recibe Axel

- **Para aprobar o corregir:**
  - C1 v1.0 (`docs/contratos/openapi.yaml` y `docs/contratos/C1-datos.md`);
  - U1 (`docs/design/U1-arquitectura.md`);
  - los criterios de las 18 fases (`docs/criterios/`).
- **Informativo:**
  - modelo de datos S1 1.1 (migración, diccionario y reglas);
  - colorimetría U2a (se aprueba en la puerta 2);
  - paquetes usados (§6): ninguno instalado con pacman.
- **Para responder:** las 31 preguntas de [`docs/estado/bloqueos/INDICE.md`](../bloqueos/INDICE.md).
- **Para verificar por su cuenta:** `make verify-wf0` (o `make verify-all`) desde la raíz del repositorio.

## 10. Consumo aproximado por agente

Datos tomados de cada reporte. «Contexto» es el contador acumulado de la sesión que reportó el agente, que incluye las relecturas del contexto; ninguno comprimió contexto.

| Agente | Acciones de herramienta | Contexto reportado | Leído / escrito (aprox.) |
|---|---|---|---|
| S1-modelo-datos | ~38 | ~400K | ~60 KB / ~300 KB |
| S6a-contrato-c1 | ~80 | ~200K | ~260 KB / ~720 KB (570 KB generados) |
| U1-arquitectura | ~47 | — | ~75 KB / ~150 KB |
| U2a-colorimetria | ~95 | — | ~60 KB / ~1.2 MB (paleta.png 864 KB) |
| criterios | ~70 | — | ~70 KB / ~250 KB |
| revisor-wf0-r1 | ~75 | — | ~420 KB / ~69 KB |
| correccion-S1-r1 | ~50 | > 300K | — / ~205 KB |
| correccion-C1-r1 | ~70 | — | ~200 KB / ~236 KB (más ~700 KB regenerados) |
| correccion-U2a-r1 | ~65 | — | ~190 KB / ~1.2 MB regenerados |
| correccion-criterios-r1 | ~95 | ~400K | ~280 KB / +118 y −70 líneas |
| revisor-wf0-r2 | ~85 | — | ~330 KB / ~80 KB |
| correccion-C1-r2 | ~100 | ~370K | ~250 KB / ~30 KB editados (más ~735 KB regenerados) |
| revisor-wf0-r3 | ~75 | — | ~380 KB / ~45 KB |
| integrador-wf0 | ~60 | — | ~420 KB / ~85 KB |
| **Total** | **~1,005** | | |

El plan (§5) estimaba 850K tokens para WF-0 con 7 agentes. Hubo 14 agentes: 5 de trabajo, 3 rondas de revisión, 5 correcciones y el integrador. Varios pasaron los 300K de referencia en el contador de su sesión. Es dato para recalibrar WF-1 en la puerta 1.

## 11. Huellas de los entregables (sha256, primeros 16 caracteres)

Las 12 primeras coinciden con las que lista `REVISION-r3.md` y la de `paleta.png` con la de `correccion-U2a-r1.md`, así que no hubo cambios después de la aprobación.

| Artefacto | sha256 | Bytes |
|---|---|---|
| `backend/db/migrations/00001_esquema.sql` | `747fc8fd88e576f4` | 137,332 |
| `docs/contratos/diccionario-datos.md` | `edda526135849067` | 161,619 |
| `docs/contratos/reglas-negocio.md` | `b93aaf9bc094ccd9` | 37,505 |
| `docs/contratos/openapi.yaml` | `95c8866589d173bc` | 597,527 |
| `docs/contratos/C1-datos.md` | `e3ec3c882ee4e274` | 113,439 |
| `backend/internal/api/routes.go` | `b4b4d33296eb4940` | 24,507 |
| `docs/design/U1-arquitectura.md` | `45fcfb61da0ca09f` | 123,059 |
| `docs/design/colorimetria.md` | `beadf8708b445a60` | 70,117 |
| `docs/design/contraste.csv` | `9f995320c6fed1f6` | 24,386 |
| `docs/contratos/C2-colores.json` | `1446c402b3226d8e` | 52,561 |
| `docs/design/tools/colorimetria.py` | `c70e643bb4300282` | 121,782 |
| `oraculo/model.json` | `1a57066689ddeed3` | 21,017 |
| `docs/design/paleta.png` | `7746da788ff19cbb` | 863,848 |
