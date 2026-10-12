# Bloqueo S5-profeco-filtro
Fecha: 11-oct-2026 · Abierto por `correccion-S5-r2` (WF-1) · Origen: `docs/estado/wf1/REVISION-S5-r2.md` R2-02 (antes `REVISION-S5-r1.md` H-01). Se decide en la puerta 2; el integrador puede cerrarlo antes como convención del contrato si lo juzga técnico.

Contexto: el filtro de `profeco_import` (`docs/contratos/motores/S5-interfaz.md` 1.4 §11) dice «`cadena_comercial` ∈ nombres PROFECO de las cadenas (`store.profeco_name`) y `municipio` ∈ municipios de las sucursales (`store.municipality`)» y «tienda del precio: la sucursal cuyo `profeco_name` = `nombre_comercial`; si no hay, la cadena». El motor (`backend/internal/precios/profeco/filter.go`) y el oráculo (`oraculo/pre_profeco.py`) lo leían distinto y daban cifras distintas: con una cadena con municipio y sin sucursal, el motor cargaba 1 y filtraba 3; el oráculo cargaba 2 y filtraba 2. Con `profeco_name` «Wal-mart» y el CSV en «WALMART», el motor cargaba 2 y el oráculo 1.

Qué falta: tres preguntas para Axel (o para el dueño del contrato):
1. ¿El `municipality` de una tienda que **no** es sucursal (una cadena o un mercado) entra al filtro de municipios?
2. ¿Cómo se comparan `cadena_comercial`, `municipio` y `nombre_comercial` de PROFECO con `profeco_name` y `municipality` del catálogo: igualdad exacta o comparación normalizada?
3. Para elegir la sucursal del precio por `nombre_comercial`, ¿se buscan solo las sucursales de la cadena que casó o cualquier sucursal del catálogo?

Opciones:
- **A (PROVISIONAL vigente):** (1) solo los municipios de las sucursales (`kind = branch`), como dice el texto de §11; una cadena o un mercado con municipio no amplía el filtro. (2) Comparación normalizada en las tres columnas: minúsculas; sin acentos (á à ä â → a, é è ë ê → e, í ì ï î → i, ó ò ö ô → o, ú ù ü û → u, ñ → n); se quitan guion, apóstrofo (' y ’), punto y acento grave; cualquier otro signo o espacio separa palabras y se deja un solo espacio entre ellas, sin espacios al inicio ni al final. Así «Wal-mart» = «WALMART», «Sam's Club» = «SAMS CLUB» y «Gustavo A. Madero» = «GUSTAVO A. MADERO». (3) Solo las sucursales de la cadena que casó; si ninguna, el precio va a la cadena. `coverage[]` agrupa por el texto de `cadena_comercial` y `municipio` tal como llega en el CSV (sin normalizar). Consecuencia: Axel captura nombres y municipios como quiera (mayúsculas, acentos); para que entre un municipio, la sucursal de ese municipio debe existir en el catálogo; dos cadenas que solo difieran en signos se tomarían como la misma (no se conoce ningún caso en los 3 cortes).
- **B:** (1) el municipio de cualquier tienda con municipio entra; (2) igualdad exacta tras quitar espacios; (3) cualquier sucursal. Consecuencia: `profeco_name` y `municipality` se capturan con la ortografía exacta de PROFECO (mayúsculas, guiones, acentos como vengan), o el filtro no deja pasar nada; un mercado con municipio deja pasar los registros de todas las cadenas en ese municipio.
- **C:** mezcla: (1) de A y (2) exacto de B. Consecuencia: el filtro depende de que Axel teclee el municipio en mayúsculas igual que PROFECO.

Impacto: cambian `loaded`, `filtered_out` y `coverage[]` de `profeco_import` y los precios `source = profeco` que pasan a `price` (solo referencia; PROFECO nunca compite, K22). No cambia ninguna cifra canon ni el conteo de S5-05 (37,548,200 registros en `qqp.sqlite`, que no se filtran).

Cómo seguí: motor y oráculo implementan la opción A (el motor ya lo hacía; se alineó el oráculo `oraculo/pre_profeco.py`). La fixture nueva `oraculo/fixtures/s5/s5-profeco-filtro.json` (generada por `oraculo/pre_genera_fixtures.py`) prueba los tres puntos y va marcada PROVISIONAL con este bloqueo. Solicitud al dueño del contrato: escribir la opción A en §11 con la marca «PROVISIONAL (bloqueo `S5-profeco-filtro`)». Esta opción sustituye la convención 11 de `docs/estado/wf1/oraculo-S5.md` («filtros por igualdad exacta; municipio = de cualquier tienda»).


## Puerta 2 (11-oct-2026): sin respuesta
Fuente: `docs/estado/puerta2/RESPUESTAS.md` («Las preguntas sin recomendación quedan sin respuesta en esta puerta: PROVISIONAL con la opción que ya está en uso, y pasan a WF-3») y hoja de respuestas de `PUERTA-2.md` §8. Axel puede mandar una excepción por id cuando lea el documento.

### Q2-56 · (1) ¿El municipio de una tienda que no es sucursal entra al filtro?
- PROVISIONAL: sigue en uso A (cita de este archivo: «**A (PROVISIONAL vigente):**» y «motor y oráculo implementan la opción A»; `docs/estado/wf1/correccion-S5-r2.md`). Pasa a WF-3.
- **Estado:** ABIERTO

### Q2-57 · (2) ¿Cómo se comparan los nombres y municipios de PROFECO con los del catálogo?
- PROVISIONAL: sigue en uso A, comparación normalizada (misma cita). Con B, Axel capturaría nombres y municipios con la ortografía exacta de PROFECO. Pasa a WF-3.
- **Estado:** ABIERTO

### Q2-58 · (3) Para asignar la sucursal del precio, ¿solo las de la cadena que casó o cualquiera?
- PROVISIONAL: sigue en uso A (misma cita). Pasa a WF-3.
- **Estado:** ABIERTO

**Estado del bloqueo (puerta 2): ABIERTO** — Q2-56, Q2-57 y Q2-58 sin respuesta; pasan a WF-3 con A como PROVISIONAL.
