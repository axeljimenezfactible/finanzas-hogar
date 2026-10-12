# Bloqueo S5-motor-oraculo-reglas

Lo abrió el agente `correccion-motor-S5-r1` el 11-oct-2026 en la ronda P2c, al corregir M-1 y M-2 de `docs/estado/p2c/REVISION-motor-S5-r1.md`. Se decide en WF-3, o antes si contratos o Axel contestan por número.

Hay 8 reglas que `S5-interfaz.md` 1.5 no fija y en las que el motor (`backend/internal/precios`) y el oráculo (`oraculo/pre_*.py`) hacen cosas distintas, o en las que el oráculo copió al motor. Ninguna fixture las prueba, así que hoy `compara.py s5` da 0 diferencias. Ninguna mueve una cifra canon: +$138.50, +$100.50, 61/60/55 piezas, 33/3 renglones y 22/8/3 ganadores.

**Regla mientras no se conteste:** PROVISIONAL, **en uso: el motor**, que es lo que corre en la app. El código de los dos lados cita este bloqueo con el número de la pregunta. No se eligió qué lado tiene razón. Cuando se fije una regla, se alinea el lado que no la cumpla y se agrega una fixture con el caso (las cifras de abajo sirven de base).

Las cifras de las preguntas 1 a 3 las medí con una fixture borrador, sin `esperado`, en `datos/tmp/corr-s5-r1/` (fuera de git), con el binario `motores s5 dump` y con `pre_run_fixtures.py`. Las de la pregunta 4 salen de una copia literal de `dateLayouts` y de `_fecha` del oráculo.

## Preguntas sobre la comparativa (§7.2) y el precio vigente (§6): de negocio o de pantalla

Si contratos las juzga técnicas, las fija en el contrato. Si no, son para Axel.

### 1. ¿Qué es un «renglón» en el supuesto S-11 de `piece_count`? (§7.2 paso 6; M-1)
S-11 dice: «un renglón con cantidad no entera (a granel) cuenta como 1 pieza». No dice qué renglón.
- **A) La presentación agrupada del paso 1.** Es lo que hace el motor, y es lo que está en uso.
- **B) Cada línea de `rows[]`,** con la parte emparejada y el excedente por separado. Es lo que hace el oráculo.
- **C) Cada línea capturada** (la del ticket o del carrito).

Cifras (`piece_count` de Bodega y de Walmart):

| Caso | A | B | C |
|---|---|---|---|
| `piezas-1`: Bodega 0.5 de arroz y 0.5 + 0.25 de atún; Walmart 0.5 y 0.5 | 2 y 2 | 3 y 2 | 3 y 2 |
| `piezas-2`: Bodega 0.25 + 0.25 de atún y 0.5 de arroz; Walmart 0.5 de atún y 1 bolsa de arroz (sobra 0.5) | 2 y 2 | 2 y 3 | 3 y 2 |

`identical_pieces` da 2 en los dos lados.

**Impacto:** solo el conteo de piezas (K-07) con cantidades no enteras. No toca dinero, ganadores ni el canon (todas sus cantidades son enteras). Sin recomendación.

### 2. ¿Qué precio toma la parte emparejada si una canasta trae la misma presentación a precios distintos? (§7.2 pasos 1 y 2; M-2a)
- **A) El orden de captura:** lo emparejado sale de las primeras líneas y el excedente de las últimas. Es el motor, en uso.
- **B) El precio promedio** de esa presentación en la canasta (importe exacto ÷ cantidad), tanto para lo emparejado como para el excedente. Es el oráculo.
- **C) Otra regla,** por ejemplo emparejar primero la línea más barata.

Cifras: Bodega trae 1 atún a $20.00 y otro a $22.00; Walmart, 1 a $21.00.

| Salida | A (motor) | B (oráculo) |
|---|---|---|
| `identical_total_cents` de Bodega | 2000 | 2100 |
| `winner_store` del renglón idéntico | Bodega | sin ganador |
| `ties` | 0 | 1 |
| `winners` de Bodega | 1 | 0 |
| Línea «solo en una tienda» de Bodega | 2200 | 2100 |
| Precio por unidad base, línea idéntica de Bodega | 14286 | 15000 |
| `total_cents` de Bodega | 4200 | 4200 |

**Impacto:** cambia el desglose de idénticos, los ganadores y los empates que se ven en pantalla. No cambia el total de cada tienda ni el canon. Sin recomendación.

### 3. ¿`current_price.references[]` lleva una referencia por fuente o todas? (§6; M-2b)
- **A) La más reciente de cada fuente** (`profeco`, `site_assisted`). Es el motor, en uso.
- **B) Todas** las de fecha hasta hoy. Es el oráculo.

Cifras: atún en Bodega con PROFECO del 01-oct a $19.00, PROFECO del 05-oct a $19.50 y `site_assisted` del 03-oct a $20.50.
- A: 2 referencias (05-oct PROFECO y 03-oct `site_assisted`).
- B: 3 referencias, con la del 01-oct al final.

El orden de §15.2 (fecha descendente y luego fuente) es el mismo en los dos.

**Impacto:** solo lo que se muestra como referencia. Las referencias nunca entran en ganadores, totales ni el costo de S4 (§6). Sin recomendación.

## Solicitudes técnicas a contratos (§11, cargador PROFECO; M-2c y M-2d)

### 4. Formatos de `fecha_registro` aceptados (§11 paso 4)
§11 dice que el cargador «acepta los formatos que encuentre, los declara en el reporte de cobertura y rechaza con motivo los que no reconozca», pero no da la lista. Medido:
- **Solo el motor acepta:** `2026/8/3` y `2026-8-3` (sin ceros), `3/8/2026`, `20260831` y `2026-08-31T10:20:30Z`.
- **Solo el oráculo acepta:** la hora sin segundos (`2026-08-31 10:20`, `2026-08-31T10:20` y `2026/08/31 10:20`) y `31/08/2026 10:20:30`.
- **Los dos aceptan:** `2026/08/31`, `2026-08-31`, `31/08/2026`, `2026-08-31 10:20:30` y la fracción de segundo (`2026-08-31 10:20:30.5`). El informe del motor decía que solo el oráculo aceptaba fracciones; no es así, porque `time.Parse` de Go las acepta después de los segundos.

**Opciones:** A) la lista del motor (en uso); B) la del oráculo; C) una lista cerrada en §11 con los formatos que de verdad traen los tres cortes (`S5-cobertura-profeco.md`). **Impacto:** el conteo de rechazados por «fecha ilegible». El corte 2026 trae `AAAA/MM/DD`, que aceptan los dos.

### 5. Liga a un `presentacion_id` desconocido
- **A)** `s5 dump` (`backend/cmd/motores/s5.go`) responde a la consulta con `unknown_presentation`. Es el motor, en uso.
- **B)** El oráculo cuenta el registro en `skipped`.

En producción las ligas llevan el id numérico. No medí qué pasa con un id que no está en el catálogo. **Impacto:** el resultado de `profeco_import` cuando `vinculos.csv` trae un código que no existe.

### 6. Tope de `rejections[]`
- **A)** El motor guarda como máximo 100 rechazos con detalle (`Options.MaxRejections`, 0 = 100); los conteos siguen completos. En uso.
- **B)** El oráculo no tiene tope.

Hoy ninguna fixture pasa de 100 rechazos. **Impacto:** la lista de `rejections[]` de un corte real (37.5 millones de registros).

### 7. Gramática del precio: supuestos S-1 a S-3 (§11 paso 3; M-2d)
- **S-1.** Un espacio en medio («1 021.00») da «precio ilegible». Alternativa: quitar todos los espacios.
- **S-2.** Más de 15 dígitos enteros da «precio ilegible»: es la protección del motor contra el desborde de int64.
- **S-3.** «Más de dos decimales» cuenta todo lo que sigue al primer punto, aunque no sean dígitos: «1.2.3» da «precio con más de dos decimales». Alternativa: «precio ilegible».

El oráculo **copia** al motor (lo decía su comentario: «supuesto del agente motor-S5-puerta2»). Por eso que coincidan no es evidencia independiente (§15.3).

### 8. Liga repetida después de normalizar
- **A)** Error de todo `vinculos.csv` (`invalid_argument`). Lo hacen el motor (`profeco/links.go`) y el oráculo, que lo copia. En uso.
- **B)** Se queda la primera y la repetida se reporta.

**Impacto:** con la liga normalizada (Q2-57, PROVISIONAL), «ATUN DOLORES» y «ATÚN DOLORES» en dos renglones hacen que falle toda la importación.

## Cómo seguí
Los dos lados quedan como están y llevan la marca «PROVISIONAL, bloqueo S5-motor-oraculo-reglas pregunta N». Lo que el oráculo copió del motor (preguntas 7 y 8) quedó marcado como copia. No hay fixtures nuevas: el revisor pidió agregarlas «con la regla fijada».

Solicitudes, en `docs/estado/p2c/correccion-motor-S5-r1.md`:
- A contratos: fijar las 8 en `S5-interfaz.md` (§7.2 paso 6, §7.2 pasos 1 y 2, §6 y §11) o pasarlas a Axel.
- Al integrador: agregar este archivo al índice y decidir si las preguntas 1 y 2 se suman a `S5-carritos-granel` y `S5-carritos-detalle`.

- **Estado:** ABIERTO
