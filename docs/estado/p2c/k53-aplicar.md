# P2d · k53-aplicar · reporte (12-oct-2026, madrugada)

Agente `k53-aplicar` (Opus), flujo «P2d · K53». Aplica la respuesta de Axel a Q2-64 en la puerta 2 bis (`docs/estado/puerta2/RESPUESTAS.md`) y dos remanentes de P2c (I-2 e I-6). No hice `git add` ni commit. No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`.

## Hecho

### K53 en `docs/plan/decisiones.md` §K
- Lleva las palabras literales de Axel y estas consecuencias:
  - La aplicación NO lee sitios de tiendas.
  - Se retiran la lectura asistida, su interruptor y el veredicto de términos. S5-07 pasa a «no aplica».
  - El agente externo es una sesión de Claude Code de Axel con Claude in Chrome. Mete precios, productos y recetas por la API con `serviceToken`.
  - El bloqueo `S5-terminos` queda RESUELTO.
  - Preguntas (a) a (f) para la puerta 4.
  - Lo PROVISIONAL va marcado como PROPUESTA del orquestador.
- K37 y K48 llevan al final la marca «Sustituida en parte por K53» y «Sustituida por K53». No borré texto.
- No edité §E: `docs/criterios/S5.md` la cita literal. K53 dice que la sustituye en lo de la lectura asistida.

### Bloqueos
- `S5-terminos.md`: sección «Respuesta de Axel (puerta 2 bis…)» con el literal, la decisión K53 y sus consecuencias. Cierra con **«Estado del bloqueo (puerta 2 bis): RESUELTO»** por K53.
- `S5-agente-externo.md` (NUEVO, ABIERTO):
  - Las preguntas de la puerta 2 bis: (a) peso de los precios, (b) qué escribe el agente, (c) receta en borrador o directa.
  - Tres técnicas que salen del esquema 1.4: (d) `price.source` es un CHECK cerrado sin `external_agent`; (e) no existe ninguna credencial de servicio, y `idempotency_key.user_account_id` es NOT NULL; (f) `recipe.origin` es un CHECK cerrado y `recipe` no tiene estado de borrador.
  - Opciones, lo PROVISIONAL e impacto.
- `INDICE.md`:
  - `p2c-remanentes` ya había registrado el conteo antes que yo, con la carpeta en disco: 85 archivos, 64 resueltos y 21 abiertos. También la fila P2d, la fila Q2-64 y la sección «Bloqueos nuevos de P2d».
  - Comprobé que coincide: hay 85 archivos más el índice.
  - No lo dupliqué. Solo marqué «RESUELTO por K53» en la fila de `S5-terminos.md` de la tabla de WF-1.

### `docs/estado/S5-terminos.md`
- Nota de cierre «No aplica desde K53» bajo el título. El registro del 11-oct-2026 se conserva sin cambios.

### `docs/contratos/motores/S5-interfaz.md` → 1.6
- Encabezado y «Cambios de la versión 1.6».
- §1: `asistida` se retira en WF-3, sin cliente HTTP hacia tiendas.
- §3, fila `site_assisted` = «agente externo»:
  - Es PROVISIONAL (d2) porque el CHECK es cerrado. No inventé un valor nuevo.
  - Peso PROVISIONAL «compite como `manual`».
  - HECHO: el motor y la fixture `s5-fuentes-compiten` todavía lo tratan como referencia.
- §13 reescrita como «Agente externo», con la lectura asistida «no aplica».
- Fila S5-07 de §18.
- No cambia ninguna cifra canon ni fixture.

### `backend/internal/precios/asistida/`
- Solo era el interruptor (`CanEnable`, `SiteFor`; 2 887 B) y su prueba (2 463 B).
- Lo dejé intacto. **Se retira en WF-3**, junto con las fixtures `s5-lectura-asistida*`.

### C1 v1.4.0, solo por el generador
Cambios en `backend/internal/api/c1gen/`: `gen2_schemas.py`, `gen3_ops.py`, `gen4_build.py`, `gen6_c1.py`, `c1_plantilla.md`, `verify_c1.py` y `README.md`. Después regeneré con `build.py`: `openapi.yaml`, `routes.go`, `C1-datos.md` y los extractos.

| | Antes (1.3.0) | Después (1.4.0) |
|---|---|---|
| Operaciones | 212 | 216 |
| Rutas | 168 | 172 |
| Esquemas | 255 | 262 |

- **Operaciones nuevas** (etiqueta «Compras y tiendas», marca `x-agente-externo`, sin N-xxx, `security: [{serviceToken: []}]`):
  - `getExternalAgentReference`: GET `/api/v1/external-agent/reference`.
  - `importExternalCatalog`: POST `.../catalog`.
  - `importExternalPrices`: POST `.../prices`, con `external_ref` → `price.source_ref`.
  - `importExternalRecipes`: POST `.../recipes`, en borrador.
  - Los tres POST son idempotentes por llave natural y llevan `Idempotency-Key`.
  - Todo lo que depende de la puerta 4 es `x-provisional` y cita `S5-agente-externo`.
- **Esquemas nuevos:** `ExternalAgentReference`, `ExternalCatalogItem`, `ExternalCatalogImport`, `ExternalPriceItem`, `ExternalPriceImport`, `ExternalRecipeImport` y `ExternalImportResult`.
- **Seguridad:** `serviceToken` (http bearer, PROVISIONAL).
- **Lectura asistida retirada:**
  - `createStore`, `updateStore`, `getPriceSettings` y `updatePriceSettings`: con `true`, responden siempre 409 `assisted_reading_not_allowed`.
  - `Store`, `PriceSources`, `PriceSettings` y `PriceSettingsUpdate` dicen «no aplica desde K53».
- **I-2:** `SavingsGoal.status` = `in_progress`, `achieved` y `overdue`, con la regla de S2 1.3 §3.13 (Q2-18 A).
- **I-6:** `getCashFlowFortnights`:
  - `planned_savings_cents` sale de `required` y es PROVISIONAL (`S2-ahorro-quincenal`).
  - `balance_cents` = ingresos − pagos de crédito − fijos y venideros. No resta el ahorro planificado ni el variable estimado.
- **`c1_plantilla.md`:** fila de Cambios y detalle «C1 v1.4 (puerta 2 bis)», viñeta de §2 y dos filas en §13. §11.2 da el motivo de cada excepción.
- **`verify_c1.py`:**
  - Exige 1.4.0 y «C1 v1.4».
  - C1-02 acepta sin N-xxx solo `GET /openapi.yaml` y las 4 del agente externo.
  - Esas 4 llevan solo `serviceToken` y la etiqueta «Compras y tiendas». Ninguna otra operación acepta el token.

### Criterios
- `S5.md` 1.4: S5-07 conserva id y severidad (alta) y su texto pasa a «No aplica desde K53».
  - Verifica, sin dominios de tiendas en el código: el grep de dominios en `.go` (sin `_test.go`) da 0 líneas.
  - Verifica que el interruptor nace en false y que ninguna migración lo enciende.
  - Verifica las 4 operaciones y la versión 1.4.0 en C1.
  - Al día: «Estado», §4 (fallas), §5 (cifras) y fuentes (K53).
- `C1.md` 1.4: C1-13 y C1-21 citan 1.4.0, 216 operaciones, 172 rutas y 262 esquemas. Lleva su bloque de cambios.
- `INDICE.md` de criterios 1.4:
  - Versiones de S5 y C1.
  - Fila K53, nota en K37, fila Q2-64 y PROPUESTAS (1) y (2) al día.
  - Siguen 261 criterios: 201 alta y 60 media.

### `plan.md` y `README.md`
- `plan.md` §3.4, fila S5: «sustituido por K53» junto al veredicto de términos.
- `plan.md` §7, R8: «Sustituido por K53».
- `README.md`: una frase de estado P2d (K53).

## Verificación (comandos y salidas)

| Comando | Salida |
|---|---|
| Conteo con pyyaml sobre `openapi.yaml`, antes | `ANTES ops 212 rutas 168 esquemas 255 version 1.3.0` |
| `python3 -I backend/internal/api/c1gen/build.py` | regenera `openapi.yaml`, `routes.go`, `C1-datos.md` y 20 extractos más su `INDICE.md` |
| Conteo con pyyaml, después | `DESPUES ops 216 rutas 172 esquemas 262 version 1.4.0`; `['in_progress', 'achieved', 'overdue']`; `required` de la fila sin `planned_savings_cents` |
| `make verify-c1 verify-openapi` | `RESULTADO: los 22 archivos coinciden byte a byte`, `comprobaciones OK: 33`, `TODO OK`; `openapi.yaml: OK` |
| `verify_c1.py`, líneas C1-02, C1-13 y C1-15 | `operaciones sin N-xxx: 5 [... 4 external-agent ..., 'GET /openapi.yaml']`; `OK C1-02 ... las 4 del agente externo`; `OK C1-13 ... C1 v1.4`; `C1-15 marcas PROVISIONAL en openapi.yaml: 15 · bloqueos abiertos citados: [..., 'S2-ahorro-quincenal', ..., 'S5-agente-externo'] · sin bloqueo abierto: []` |
| `cd backend && go build ./... && go vet ./internal/api/` | `GO-OK`; `grep -c external-agent internal/api/routes.go` = 4 |
| Verifica de S5-07, el grep de dominios en `.go` sin `_test.go` | 0 líneas |
| `grep -c "price.assisted_reading_enabled', 'false'" …/00001_esquema.sql` | 1 |
| `grep -rln "assisted_reading_enabled.*true" backend/db/migrations/` | 0 archivos |
| `grep -c 'operationId: "importExternal\|operationId: "getExternalAgentReference"' docs/contratos/openapi.yaml` | 4 |
| `grep -n -m1 'version: "1.4.0"' docs/contratos/openapi.yaml` | `7:  version: "1.4.0"` |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **exit 0**; `TOTAL 261 (alta 201 / media 60 / baja 0)`; `TOTAL comprobadas: 178 \| con diferencias: 0`; `RESUMEN TODO OK` (S1); `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` |

## Supuestos
1. **Etiqueta «Compras y tiendas», no una nueva.** `verify_c1.py` C1-13 exige exactamente las 9 etiquetas de la barra lateral (TAGS9).
2. **Excepción justificada a C1-02.** Las 4 operaciones van sin N-xxx porque no tienen pantalla. C1-02 admite «solo las excepciones justificadas», y `verify_c1.py` las comprueba una por una.
3. **`external_ref` va a `price.source_ref`.** Es la columna que ya existe («Referencia de la fuente … URL»). No agregué columnas. Catálogo y recetas no llevan `external_ref`.
4. **«Esquema ServiceToken mínimo».** Lo interpreté como esquema de seguridad `serviceToken` (http bearer) y no como objeto de datos: C1-12 prohíbe propiedades con «token» en las respuestas.
5. **`getExternalAgentReference` no lo pedía el encargo.** Es de solo lectura y lo agregué para que el agente obtenga ids de tiendas, grupos y tiempos sin cookie de sesión.
6. **«`true` → 409 siempre» sale de la regla vigente.** Antes: sin veredicto, 409. Ya no hay veredicto, así que no es un comportamiento nuevo.
7. **Reglas técnicas de rechazo por renglón que escribí yo:** unidad base distinta, código de barras ajeno y receta con nombre existente. La de recetas va dentro del texto PROVISIONAL.

## Bloqueos
- **`S5-agente-externo` (nuevo, ABIERTO, puerta 4):** preguntas (a) a (f).
- **Conflicto del encargo, documentado y no resuelto por mí.** El encargo pide tres cosas a la vez:
  - PROVISIONAL «compite como manual»;
  - `site_assisted` como valor;
  - no cambiar fixtures.

  Las tres no caben juntas: `s5-fuentes-compiten` trae 2 precios `site_assisted` más baratos que hoy nunca ganan.

  Quedó así: el contrato S5 1.6 dice PROVISIONAL «compite». Lo marqué como HECHO que el motor y la fixture siguen con la regla anterior hasta WF-3. Va también en el impacto del bloqueo.

## Pendientes
- **WF-3:**
  - Si (a) = compite, alinear el motor S5 y la salida esperada de `s5-fuentes-compiten`.
  - Retirar `precios/asistida` y las fixtures `s5-lectura-asistida*`.
  - Implementar las 4 operaciones y el middleware de `serviceToken`, según lo que decida la puerta 4 en (e).
- **S1, según la puerta 4:** valor de `price.source` (d), borrador y origen de la receta (f), credencial (e). Retirar `store.assisted_reading_enabled` y la clave `price.assisted_reading_enabled`.
- **I-6 bis (observación):** `variable_estimate_cents` sigue obligatorio en C1, aunque S2 1.3 §3.11 lo entrega en null. Su bloqueo `S2-flujo-quincenal` trae «Respuesta de Axel» y C1-15 no lo deja citar. El encargo solo pedía `planned_savings_cents`.
- **Fuera de mi alcance, para el integrador:**
  - `gen2_schemas.py`, «hora HH:MM» vieja (l. 364 y 374, punto 4 de INTEGRACION P2c).
  - Menciones de la lectura asistida en `reglas-negocio.md`, `diccionario-datos.md` (sale del esquema), `U4` y `protocolo-agentes.md`.
  - `plan.md` l. 141 (requisito 26).
  - El criterio C1 propio para la importación, PROPUESTA (2), no se añadió.
- **Concurrencia:** `p2c-remanentes` modificó en paralelo `backend/internal/precios/*.go`, `recetas/*.go`, varios bloqueos y el INDICE. No toqué nada de eso. `verify-all` pasó con su árbol de ese momento.

## Propuestas para Axel
- **(a) Peso.** El orquestador propone que los precios del agente compitan como `manual`. Yo agregaría que la pantalla muestre el origen «agente externo» en cada precio.
- **(d) Valor de la fuente.** Un valor propio `external_agent` (migración de S1 en WF-3) en lugar de reutilizar `site_assisted`: así no se mezcla con la regla vieja ni con la fixture.
- **(e) Credencial.** Archivo fuera del repositorio con ACL, igual que el token del bot (`criterios-token-servicio` (b)). La bitácora y la idempotencia, a nombre de la cuenta de Axel. No cambia el esquema.
- **(c) y (f) Recetas.** En borrador con `is_active = 0` hasta confirmar, si se quiere evitar la migración; un estado propio si se quiere distinguirlo en Ajustes.

## Consumo aproximado
- Acciones de herramienta: unas 55.
- KB leídos: unos 140. Todo por partes con `grep` y `sed -n`. No leí `openapi.yaml`, `C1-datos.md`, `S5-interfaz.md` ni `INDICE.md` enteros.
- KB escritos: unos 75, contando el texto de los scripts de parcheo en `/tmp/k53`: generador 22, documentación 30 y este reporte 9. Las salidas regeneradas por `build.py` no cuentan.
