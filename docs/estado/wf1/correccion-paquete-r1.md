# correccion-paquete-r1 · Corrección del paquete para Claude Design (ronda 1)

Fecha: 11-oct-2026. Contra `docs/estado/wf1/REVISION-paquete-r1.md` (6 medios: M-1 a M-6). Sin tocar bajos ni reabrir decisiones §A–§K.

## Hecho
- **M-1 (canal del aviso).** `entrada/LEEME.md` abre con «Dónde me avisas»: todos los avisos («iter-N lista», preguntas de Claude Design, fallas) se escriben en la sesión de Claude Code del proyecto Finanzas Hogar (la terminal de `/home/deepsec/finanzas-hogar`); quien los recibe es Claude Code, el orquestador. «Yo» queda definido en la línea 3. El paso 1 de «Vueltas» y «Si algo falla» repiten el canal.
- **M-2 (cuántos archivos y entrega parcial).** Paso 7: se esperan hasta 84 imágenes. Paso 9 nuevo: si entrega solo una parte, guarda lo que haya, avisa «iter-1 lista» con la lista de pantallas que faltan y no pidas el resto por tu cuenta; el orquestador decide en la revisión si hace falta otra vuelta. Igual en `salida/LEEME.md`.
- **M-3 (preguntas o plan en lugar de maquetas).** Paso 6 y «Si algo falla»: no contestes, copia lo que te pida o proponga y avisa en la sesión de Claude Code; el orquestador entrega el mensaje que hay que pegar.
- **M-4 (una sola convención de nombre).** `prompt.md` pide el artboard como `NN-pantalla-tema-dispositivo` (NN = número de la lista; minúsculas, sin acentos ni espacios), con ejemplos y la regla de Ajustes (`NN-ajustes-<sección>`). La pantalla 19 pasa a **dos estados con nombre propio** (`19-barra-lateral-expandida-...` y `19-barra-lateral-colapsada-...`), por eso el conteo es 20 pantallas, 21 con los dos estados de la 19, × 4 combinaciones = **84 artboards** (antes 80; es el conteo derivado de `prompt.md` que usa P3-06, y la lista de 20 pantallas no cambia). El ejemplo ambiguo `menu-barra` desaparece: la pantalla 20 es `20-menu-comandos-...` y `salida/LEEME.md` explica que es el menú «/». `entrada/LEEME.md` y `salida/LEEME.md` dicen «No renombres nada» y usan la misma convención. Se conservan la sugerencia `P-xx` y la regla de PDF de N páginas = N.
- **M-5 (sección «Comprobar»).** Se quitó el comando `sha256sum -c` y la mención de regenerar con el script del LEEME de Axel; ahora dice que no tiene que comprobar nada (lo hace el orquestador con `MANIFIESTO.md`).
- **M-6 (rutas y nombres de operación en U3/U4).** `docs/design/tools/paquete_entrada.sh`: nueva función `filtrar_api` (perl) que genera en `datos/tmp/paquete` una copia filtrada de U3 y de U4 antes de pandoc (las fuentes `docs/design/U3-dashboard.md`, `U4-ajustes-a/b.md` no se tocan): quita `(GET `/api/v1/...`)` (conservando el «campo `...`» cuando va en el mismo paréntesis), sustituye los nombres de operación de lectura por «consulta de C1» y los de escritura y demás por «operación de C1», colapsa «X / X», concuerda «de la consulta de C1». Nuevas comprobaciones en el paso 4/5 que hacen fallar el script: `sin_api` (`pdftotext | grep -c '/api/v1'` = 0 en U3, U4, C1-resumen y DESIGN) y `sin_operaciones` (0 líneas con `get/list/update/create/delete` + Mayúscula en U3, U4 y C1-resumen). Se añadió pdftotext y perl a los requisitos del script.
- Regeneración: `paquete_entrada.sh` completo; MANIFIESTO regenerado.

## Archivos
- `docs/design/entrada/LEEME.md` (reescrito en partes), `docs/design/entrada/prompt.md` (líneas 12 y 98: nombres y conteo; 12 284 bytes, antes 11 958), `docs/design/salida/LEEME.md` (convención de nombres).
- `docs/design/tools/paquete_entrada.sh` (filtro y comprobaciones).
- Regenerados por el script: `docs/design/entrada/U3-dashboard.pdf` (19 págs.), `U4-ajustes.pdf` (37 págs.; antes 31: «consulta de C1» ocupa más que el nombre de la operación), `MANIFIESTO.md`. DESIGN.pdf y C1-resumen.pdf se reconstruyeron (el sello del caché incluye el hash del script) con el mismo texto (md5 de `pdftotext -layout` idéntico) y hash de archivo distinto por la fecha de Chromium; DESIGN.docx, `paleta.png`, `C2-tokens.json`, `tokens.css` y los 11 PNG, idénticos a los de antes.
- Este reporte.

## Verificación
- `bash docs/design/tools/paquete_entrada.sh` → «pdftotext | grep -c '/api/v1' = 0 en los 4 PDF; sin nombres de operación en U3, U4 y C1-resumen», «sha256sum -c: 0 diferencias», «Listo: 22 archivos». Primera corrida con la comprobación de operaciones sobre DESIGN.pdf falló (1 línea, ver Solicitudes); se acotó a U3, U4 y C1-resumen.
- `pdftotext` por PDF (rutas /api/v1 y nombres de operación): DESIGN 0/1, U3 0/0 (antes 13/25), U4 0/0 (antes 0/89), C1-resumen 0/0. Texto de DESIGN.pdf y C1-resumen.pdf idéntico al anterior (md5 `32f6c533c3a1`, `f5048f129cc6`).
- Segunda corrida sin `--forzar`: los 6 productos «sin cambios en las fuentes»; MANIFIESTO idéntico (`cmp`).
- `grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c` → 0 líneas fuera de OK; `ls -A entrada | wc -l` = 22.
- Patrón de credencial de P3-01/P3-03 sobre `LEEME.md` = 0. Los 90 hex de `prompt.md` están todos en `C2-tokens.json` (script propio con `python3 -I`). Rutas absolutas del LEEME existentes (salvo `siguiente-prompt.md`, `P3-iter-1.md` e `iter-N/`, que se crean después o son plantilla).
- Revisión visual de una página de U4 (tabla de Ajustes 03 Ingresos): legible, sin rutas ni nombres de operación.
- Criterios cubiertos: P3-03 (M-1 a M-5: canal, cantidad, parcial, preguntas, convención, sin acciones fuera de la lista), P3-02 y P3-06 (conteo 84 derivado de `prompt.md`), P3-12 (sin credenciales) y verificación (e) del revisor. No hay pruebas Go, oráculo ni compara.py en este artefacto.

## Supuestos
1. El «canal» del aviso es la sesión de Claude Code del proyecto (lo sugiere el revisor; no hay otro canal en las decisiones). Si Axel prefiere Telegram u otro, cambiar solo el párrafo «Dónde me avisas».
2. La pantalla 19 cuenta como dos estados con artboard propio (84 en total), según la corrección sugerida; si se prefiere un solo artboard con ambos estados, hay que volver a 80.
3. Los nombres de operación se sustituyen por «consulta/operación de C1», sin más contexto; se pierde el nombre exacto de la operación en los PDF (siguen en `docs/design/U3-dashboard.md` y `U4-ajustes-*.md`).
4. `prompt.md` queda en 12 284 bytes, bajo 12 KiB (la observación del revisor «≤ 12 KB»); no hay tope en los criterios.

## Bloqueos
Ninguno nuevo. Se mantiene S1 (`docs/estado/bloqueos/S1-parametros-avisos.md`), que no afecta a este paquete.

## Pendientes
- Bajos B-1 a B-7 de `REVISION-paquete-r1.md` no se tocaron (fuera de esta corrección): B-1 («No contiene código»), B-2 (LEEME l. 52 y colores del diálogo crítico), B-3 y B-4 (`tools/paquete/C1-resumen.md`, fuera del alcance), B-5 (`SCRIPT_HASH` con ruta relativa, extensiones del filtro de código), B-6 (rutas absolutas restantes del LEEME; las nuevas ya van absolutas), B-7 (tarjeta 3 en `prompt.md`).
- Revisor de la ronda 2 del paquete.

## Solicitudes fuera de alcance
- `docs/design/DESIGN.md` (texto de DESIGN.pdf, línea ~3971 del texto extraído): menciona la operación `updateNotification` en la descripción PROVISIONAL del bloqueo S1. Cambio: decir «la operación de escritura de avisos» sin nombre, si se quiere que DESIGN.pdf también salga sin nombres de operación; se dejó porque no está en M-6 y es texto PROVISIONAL de U4.
- `docs/criterios/P3-claude-design.md` y `docs/design/entrada/MANIFIESTO.md`/reportes `paquete-claude-design.md`: dicen «80 artboards» (supuesto 3 y línea 7 del reporte); el integrador debe actualizarlos a 84 (el criterio P3-06 usa «número de artboards de `prompt.md`», que ahora es 84). `docs/design/tools/paquete/C1-resumen.md` no menciona el conteo.
- Opcional: versionar `docs/design/tools/paquete/.cache/` (ya solicitado antes) para que MANIFIESTO no cambie al regenerar en otro clon.

## Paquetes instalados
Ninguno (perl y pdftotext ya estaban).

## Propuestas para Axel
Ninguna (no se pidió ningún valor K31).

## Consumo aproximado
~32 acciones de herramienta; ~75 KB leídos; ~30 KB escritos (más los PDF regenerados por el script).
