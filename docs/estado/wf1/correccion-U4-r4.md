# correccion-U4-r4 · Ajustes estilo Notion (partes a y b) · ronda 4

Agente `correccion-U4-r4` · WF-1 · 11-oct-2026. Corrige el único hallazgo alto o medio de `docs/estado/wf1/REVISION-U4-r4.md` (R4-1, media, U4-05). No toqué nada fuera de ese hallazgo ni reabrí decisiones de §A–§K. Los bajos (R4-2, R4-3, R3-2 a R3-6) siguen abiertos y no son de esta ronda.

## Hecho
- **R4-1 · `docs/design/U4-ajustes-a.md` 1.3 → 1.4:**
  - **§0.2.** La excepción de las listas con vista previa en línea abarca ahora los **interruptores** con vista previa: «activo» de un gasto fijo (P-26.02) y «Activa» de una meta de ahorro (P-26.13). Al tocarlos se muestra la vista previa y no se guarda todavía. Se guarda con `Enter` o al salir del campo, y `Esc` vuelve al estado guardado. La regla nueva: la edición múltiple con `Ctrl+/` de esos campos muestra, antes de aplicar, la misma vista previa sumada sobre las filas seleccionadas. `Enter` aplica a todas y `Esc` cancela sin cambiar ninguna.
  - **P-26.02, fila «Gasto fijo: activo».** Tipo: «Interruptor con vista previa en línea (§0.2: no guarda al tocar; vista previa 3, punto g)».
  - **P-26.02, «Edición en el lugar».** «activo» entra en la lista de los campos que cambian el dinero, que ahora son 9. El interruptor «activo» no guarda al tocar. Se describe la edición múltiple con `Ctrl+/` (punto h). La frase ambigua «Interruptor: guarda al tocar.» pasa a decir a cuál se refiere: el de la proyección de los meses 7 a 12, que conserva su comportamiento y su vista previa 2 (K21).
  - **P-26.02, vista previa 3.** Encabezado con «activo» y con la edición múltiple. Punto nuevo **(g)**, al tocar «activo»:
    - (c) «Gastos fijos de {mes}», con la partida sumada o quitada;
    - (d) las partes por persona, también sumadas o quitadas;
    - (f) el aviso de P-07;
    - (a) no cambia, y (b) dice «No cuenta mientras esté inactiva».

    Punto nuevo **(h)**, `Ctrl+/` sobre categoría o activo de {n} partidas: antes de aplicar, (c), (d) y (f) sumados sobre las seleccionadas. Al cambiar la categoría, (c) no cambia, solo cambian las partes de las partidas sin reparto propio, y el bloque lo dice. Se aplica con un `updatePlannedItem` por partida (C1 no tiene operación masiva de partidas: solo `create`, `list`, `update`, `importPlannedItems` y `previewPlannedItemsImport`). `Esc` cancela.
  - **P-26.02, «Atajos».** `Ctrl+/` «con la vista previa 3 (punto h) antes de aplicar».
- **R4-1 · `docs/design/U4-ajustes-b.md` 1.2 → 1.3:**
  - **P-26.13, fila «Activa» de metas.** «Interruptor con vista previa en línea (parte a §0.2: no guarda al tocar; vista previa 3 de metas)». Las filas «Activa» de P-26.09 y P-26.11 quedan sin cambio.
  - **P-26.13, vista previa de metas.** Punto nuevo **(3)**: «Total reservado para metas en {mes en curso}: {antes} → {después}», que es la misma suma de `planned_monthly_cents` de las metas activas de la vista previa (2), con esta meta sumada o quitada. Una meta sin aportación planificada no cambia el total, y el bloque lo dice. `Enter` o salir guarda, y `Esc` revierte. Con `Ctrl+/`, el total sumado sobre las metas seleccionadas.
  - **P-26.13, «Edición en el lugar».** «Activa» entra en los campos que guardan después de su vista previa.
  - **P-26.13, «Atajos».** `Ctrl+/` cambia «Activa» de varias metas, con la vista previa (3) sumada antes de aplicar. Las convenciones de la parte b y U1 §5.3 ya ofrecen `Ctrl+/` en toda tabla. Si no se dice aquí, el mismo camino quedaba abierto sin vista previa.
  - **Cobertura de criterios de la parte b.** La línea U4-05 menciona R4-1.

## Archivos
- `docs/design/U4-ajustes-a.md`: modificado, 8 reemplazos exactos. Sigue con 405 líneas. sha256 `1bbf3ffc…6e97`.
- `docs/design/U4-ajustes-b.md`: modificado, 6 reemplazos exactos. Sigue con 388 líneas. sha256 `85c25673…c6889`.
- `docs/estado/wf1/correccion-U4-r4.md`: este reporte.
- Fuera del repositorio, en el scratchpad (`u4r4/`): los respaldos `U4-ajustes-a.antes.md` y `U4-ajustes-b.antes.md`, el editor `edita.py` (cada reemplazo exige 1 coincidencia) y la verificación `v4.py`, que es `v3.py` de la ronda 3 más 18 comprobaciones de R4-1.

## Verificación
Comando: `python3 -I <scratchpad>/u4r4/v4.py /home/deepsec/finanzas-hogar`. Además, un conteo de celdas vacías. Salida resumida:
- **U4-01:** 14 secciones, P-26.01 a P-26.14 en orden, sin bloques faltantes. 321 filas de tabla (191 en la parte a y 130 en la b), 0 con celdas vacías.
- **U4-02:** 210 `operationId` en los extractos, igual que en `openapi.yaml`. Se citan 100 en la parte a y 45 en la b, 0 inexistentes. Las 73 de 73 operaciones de Ajustes están cubiertas. `updatePlannedItem`, `updateSavingsGoal` y `listSavingsGoals` existen.
- **U4-05 (R4-1, todo en True):**
  - En §0.2, la excepción incluye «activo» y «Activa», y está la regla de `Ctrl+/` con la vista previa sumada.
  - P-26.02: la fila «activo» dice «no guarda al tocar». «activo» está en la lista de «Edición en el lugar». La frase vieja «Interruptor: guarda al tocar.» aparece 0 veces. La proyección sigue guardando al tocar (vista previa 2). La vista previa 3 trae (g) y (h) sin elementos faltantes. «Atajos» lleva `Ctrl+/` con el punto h.
  - P-26.13: la fila «Activa» con vista previa. La vista previa (3) trae el total reservado, «Sin aportación planificada…» y `Esc`. «Edición en el lugar» y «Atajos» incluyen «Activa».
  - Las filas «Activa» de P-26.09 y P-26.11 siguen intactas (1 y 1).
  - Ejemplo de (g), RN-03 y K14: apagar un gasto fijo de $199.00 con Conjunto 50/50 y una fecha en el mes resta 19,900 a «Gastos fijos de {mes}», 9,950 a Karen y 9,950 a Hogar.
- **U4-05 (regresión):**
  - Caso A/B: las series A y B de 14 valores aparecen literales y son iguales a `calculado_*`. Están los mínimos $8,094.50 y $5,000.00 y el disponible de hoy. A − B es constante.
  - Despensa 1: $3,094.50, $800.50, $3,895.00 y $973.75. La cotización da 800.50.
  - R3-1 sigue cerrado: 8 campos de la ronda 3 nombrados, los elementos (a) a (f), el reparto 9,950 y 9,950, y la regla de fechas igual a `despensas[*]` del oráculo (incluidos febrero y el recorte a fin de mes).
  - P-26.05 conserva las vistas previas 1 a 10, y R2-1, R2-2 y R2-3 siguen cerrados.
- **U4-13 y RN-17:** IMC 23.7. Meta 2,555.5625, que da 2,556 y 2,056. S4-07: 2,098.50 − 300, que da 1,799. RN-18: 50 % (50.0 %).
- **U4-04 y U4-12:** las 15 claves `@param` de 00003 y las 19 claves sembradas en las migraciones 00001 a 00003 están citadas. `send_time`, `advance_days`, `weekday` y `day_of_month` llevan `S1-parametros-avisos`. «Mínimo de seguridad» fuera de SUPUESTOS: 0.
- **U4-06:** atajos de U4 fuera de U1 §5.3: 0. No agregué atajos ni comandos «/»; `Ctrl+/` ya estaba en U1 §5.3, en §0.7 y en las convenciones de la parte b.
- **U4-09 y U4-11:**
  - 0 hex, 0 `rgb(` y 0 fechas ISO en las dos partes.
  - 21 tokens citados, 0 inexistentes. 0 bloqueos citados que no existan.
  - PROVISIONAL en la parte a: 9 líneas, todas con cita. En la b: 27 líneas. Las 5 sin cita en la misma línea (26, 104, 209, 361 y 362) son las mismas de antes (R4-3, baja, fuera de esta ronda). Los números de línea no se movieron porque no agregué líneas.
- **No aplican** las pruebas Go, la compilación de tokens ni `compara.py`. `compara.py` exige el argumento `motor`: compara la salida de un motor contra el oráculo, y U4 es un documento de diseño sin código. Las cifras del oráculo que cita U4 se cotejan directo contra `oraculo/model.json` en `v4.py`. No creé temporales en `datos/tmp`.

## Supuestos
- Para `Ctrl+/` elegí la primera opción de la corrección sugerida: mostrar el bloque sumado antes de aplicar. No saqué categoría ni activo de la edición múltiple, porque U1 §5.3 la ofrece en tablas y tableros. No cambia cifras, esquema, endpoints ni tokens.
- En P-26.13 escribí `Ctrl+/` para «Activa» porque las convenciones de la parte b ya ofrecen `Ctrl+/` en toda tabla. Callarlo dejaba implícito un cambio múltiple sin vista previa. Es el mismo camino que señala R4-1.
- El interruptor de la proyección de los meses 7 a 12 (K21) sigue guardando al tocar con su vista previa 2: R4-1 no lo incluye. Solo nombré a cuál se refería la frase vieja.

## Bloqueos
Ninguno nuevo. Siguen citados los que ya existían (`S1-parametros-avisos`, `S2-metas`, `S3-formula-tarjeta`, `S3-convencion-prestamo`, `U4-iva-credito`, `U4-vista-previa-formulas` y los demás que verifica `v4.py`, 0 inexistentes).

## Pendientes
- Hallazgos bajos que siguen abiertos y no son de esta ronda: R4-2 (`U4-ajustes-a.md` P-26.02, «cada cuánto» con V-21 y sin tope 1 a 24), R4-3 (`U4-ajustes-b.md` líneas 104, 209, 361 y 362, PROVISIONAL sin cita en la misma línea) y R3-2 a R3-6 (ver `REVISION-U4-r3.md`).
- Observación para la próxima revisión. No la cambié porque no está en R4-1. En `U4-ajustes-a.md` P-26.04, cambiar el reparto de una categoría, también con `Ctrl+/` («tipo, reparto o activa de varias»), cambia las partes de los gastos fijos de esa categoría sin reparto propio: es la vista previa 3 (d) de P-26.02. La vista previa de P-26.04 solo da el ejemplo de $199.00 y el aviso «Aplica a los gastos que registres desde ahora». El revisor decide si es un hueco.
- El integrador une las dos partes en `docs/design/U4-ajustes.md`.

## Solicitudes fuera de alcance
Ninguna.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna. R4-1 no pide valores operativos (K31).

## Consumo aproximado
- 21 acciones de herramienta.
- ≈ 125 KB leídos: protocolo, §K, `REVISION-U4-r4.md`, `criterios/U4.md`, extractos con `sed`/`grep` de las dos partes y de U1, `v3.py` de la ronda 3 y las salidas.
- ≈ 30 KB escritos por mí: `edita.py` ≈ 9.5 KB, lo agregado a `v4.py` ≈ 4.5 KB, ≈ 4.3 KB de texto nuevo en las dos partes (el script reescribe los archivos completos) y este reporte ≈ 9 KB.
