# Bloqueo S5-agente-externo

Abierto por `k53-aplicar` (flujo P2d, 12-oct-2026, madrugada) al registrar K53 (`docs/plan/decisiones.md` §K). Se pregunta en la puerta 4.

Contexto (HECHO): en la puerta 2 bis (11-oct-2026, 22:55) Axel contestó Q2-64 con estas palabras (literal de `docs/estado/puerta2/RESPUESTAS.md`): «Ignoralo, todo el tema de la investigacion de los precios lo haremos aparte, a que me refiero con esto, crearemos el sistema para que tenga yo una sesion de claude code que pueda mediante claude chrome extension sacar los precios y los productos, si bien estos datos los podemos insertar manualmente esto tambien impacta a las recetas tambien pueden ser agregadas de forma autonoma por eso la regla no funciona ya que usaremos a un agente externo». K53 lo registra: la aplicación no lee sitios de tiendas, se retira la «lectura asistida» y un agente externo (una sesión de Claude Code de Axel con la extensión Claude in Chrome) obtiene precios y productos, puede agregar recetas y los mete por la API con una credencial de servicio propia. Axel no decidió cómo pesan esos datos ni qué más puede escribir el agente.

## Qué falta (PREGUNTAS para la puerta 4)

Las tres que nombra la puerta 2 bis:
- **(a) Peso de los precios del agente externo en la comparativa.** ¿Compiten como `manual` (K37: mismo peso que el ticket y la misma regla de frescura) o son solo referencia, como PROFECO (se muestran y nunca ganan)?
- **(b) Qué puede escribir el agente.** ¿Solo productos, presentaciones, precios y recetas, o también inventario (existencias, ajustes)?
- **(c) Recetas del agente.** ¿Entran como borrador que Axel (o Ximena) confirma antes de usarlas en el plan, o entran directas al recetario?

Las que salieron al aplicar K53 contra el esquema 1.4 (HECHO de `backend/db/migrations/00001_esquema.sql`):
- **(d) Valor de `price.source`.** Es un CHECK cerrado: `ticket_manual`, `ticket_ocr`, `site_assisted`, `profeco`, `manual`. No existe `external_agent`. Opciones: (d1) S1 agrega `external_agent` con una migración y retira `site_assisted`; (d2) se reutiliza `site_assisted` y solo la documentación lo llama «agente externo»; (d3) se guarda como `manual` y la referencia del agente va en `source_ref`.
- **(e) Credencial de servicio.** El esquema 1.4 no tiene ninguna: la única credencial es la sesión de persona (`device_session.token_hash`). Además `idempotency_key.user_account_id` es NOT NULL y `audit_log.user_account_id` admite NULL («proceso del sistema»). Opciones de dónde vive la credencial: (e1) un archivo fuera del repositorio con permiso de lectura para el usuario del servicio, que la API compara en ejecución, como el token del bot en `criterios-token-servicio` (b); (e2) una tabla nueva con el SHA-256 de la credencial (migración de S1); (e3) sin credencial propia: el agente inicia sesión con la cuenta de Axel (cookie). Y a nombre de quién quedan la clave de idempotencia y la bitácora: (e4) a nombre de la cuenta de Axel; (e5) de una cuenta de servicio nueva (S1).
- **(f) Borrador y origen de una receta.** `recipe.origin` es un CHECK cerrado (`claude_seed`, `household`) y `recipe` no tiene estado de borrador. Opciones: (f1) S1 agrega `external_agent` a `origin` y un estado de borrador; (f2) borrador = `is_active = 0` y origen `household`, con la marca del agente en `notes`; (f3) si (c) = directa, solo cambia `origin`.

La que salió en la revisión `revisor-k53-r1` (hallazgo M-1 de `docs/estado/p2c/REVISION-k53-r1.md`; la abrió `correccion-k53-r1`):
- **(g) Tamaño máximo de un lote.** HECHO: C1 1.4 dice «De 1 a 500 productos» (`ExternalCatalogImport`), «De 1 a 500 precios» (`ExternalPriceImport`) y «De 1 a 100 recetas» (`ExternalRecipeImport`), y `importExternalCatalog` responde 422 a un lote vacío o de más de 500 renglones. Esos números no salen de RESPUESTAS.md, de §K ni de C1 1.3: los eligió `k53-aplicar` y no los declaró como supuesto. Opciones: (g1) dejar 500 productos, 500 precios y 100 recetas por lote; (g2) otros topes que dé Axel (por ejemplo, uno solo para los tres lotes); (g3) sin tope de renglones en C1: el lote solo pide al menos un renglón y el agente parte sus envíos como quiera.

## Mientras tanto (PROVISIONAL)

Es PROPUESTA del orquestador, escrita en el encargo de P2d. No es decisión de Axel:
- (a) Los precios del agente externo **compiten como `manual`**, con la misma frescura (K22, K40, K44).
- (c) Las recetas del agente **entran como borrador** con origen «external_agent».
- (b) El agente escribe solo lo que nombra K53: productos, presentaciones, precios y recetas. Inventario no.
- (d) Opción (d2), por instrucción del encargo («no inventes un valor nuevo»): el valor guardado es `site_assisted` y la documentación lo llama «agente externo».
- (e) C1 1.4 declara solo el transporte: el esquema de seguridad `serviceToken` (encabezado `Authorization: Bearer`), aceptado únicamente por las cuatro operaciones del agente externo. Dónde vive la credencial y a nombre de quién queda la bitácora quedan sin decidir. PROPUESTA de `k53-aplicar`: (e1) y (e4), porque no cambian el esquema y siguen el patrón ya aprobado para el token del bot.
- (f) Sin decidir. C1 1.4 describe el borrador. Con el esquema 1.4, el backend no puede guardar ni el estado de borrador ni el origen «external_agent».
- (g) Supuesto de `k53-aplicar`, declarado en la corrección `correccion-k53-r1`; no es decisión de Axel ni PROPUESTA del orquestador: la opción (g1). C1 1.4 lo marca PROVISIONAL, con esta pregunta, en los tres esquemas de lote y en `importExternalCatalog`.

## Impacto

- C1 v1.4: `getExternalAgentReference`, `importExternalCatalog`, `importExternalPrices` e `importExternalRecipes`, y el esquema de seguridad `serviceToken`, con `x-provisional`.
- `docs/contratos/motores/S5-interfaz.md` 1.6, §3 y §13.
- **Motor S5 y fixture `s5-fuentes-compiten`.** Hoy tratan `site_assisted` como referencia, según la regla anterior de K37. Si (a) = compiten y (d) = (d2), WF-3 cambia esa regla del motor y la salida esperada de la fixture: los 2 precios `site_assisted` más baratos, que hoy nunca ganan, pasarían a competir. P2d no toca motor ni fixtures, así que hasta WF-3 el motor va detrás del contrato en esa regla.
- WF-3 no puede construir `importExternalRecipes` sin decidir (c) y (f), ni la verificación de la credencial sin decidir (e).
- (g): si cambia el tope, cambian las descripciones de C1 (los tres esquemas de lote y `importExternalCatalog`), la validación 422 que construya WF-3 y sus pruebas. No toca el esquema de S1 ni cifras canon.
- Criterios: S5-07 («no aplica desde K53») y los futuros criterios de C1 para la importación.

## Cómo seguí

C1 1.4 y S5-interfaz 1.6 usan lo PROVISIONAL de arriba y citan este bloqueo. No cambió ninguna cifra canon ni fixture, ni el esquema de S1. El código de `backend/internal/precios/asistida/` queda igual: se retira en WF-3.

**Estado del bloqueo (P2d): ABIERTO**. Se pregunta en la puerta 4.
