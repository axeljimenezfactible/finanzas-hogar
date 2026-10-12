# Corrección U3 · ronda 2 (WF-1)

Agente `correccion-U3-r2` · 11-oct-2026 · Artefacto: `docs/design/U3-dashboard.md`, de la v1.1 a la v1.2 (479 líneas, 79.3 KB) · Hallazgos atendidos: H-1 (alto), H-2 y H-3 (medios) de `docs/estado/wf1/REVISION-U3-r2.md`. No toqué los bajos L-1 a L-6 ni ninguna decisión de §A a §K.

## Hecho
- **H-1 (alto) · T-2 con dos orígenes.** Reproduje el hallazgo con `oraculo/model.json` (Decimal). Con la ventana desde nov-2026, la regla de la v1.1 da cima 2,921.25 y en mar-2027 13,388.73. El oráculo da 3,895.00 y 14,362.48 (R-51). La causa está en U3, no en el oráculo: sumaba `cash_cents` desde el primer mes visible y encima apilaba `cumulative_difference_cents`, que C1 acumula desde P1. Corrección:
  - T-2 exige que las dos capas tengan el origen de `cumulative_difference_cents`: capa 1 = efectivo acumulado, cima = devengado acumulado, capa 2 = `cumulative_difference_cents`. Prohíbe mezclar orígenes.
  - C1 v1.2 no trae esos dos acumulados. Por eso las capas, sus etiquetas finales y su resumen quedan **PROVISIONAL**, con dos variantes escritas en T-2: (a) C1 agrega `cumulative_accrued_cents` y `cumulative_cash_cents`, que es lo que U3 pide en §15; (b) todo relativo a la ventana, con rótulos «en estos 12 meses» y `cumulative_difference_cents` en la tabla como «Diferencia acumulada desde el inicio».
  - Abrí el bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` con 3 opciones; recomiendo A.
  - A-11 suma la corrida con reloj en nov-2026: fila nov-2026 (R-46); capas (a) 2,921.25 y 3,895.00 en nov-2026, y 13,103.78 y 14,362.48 en mar-2027 (R-52, R-51); la fórmula de (b).
- **H-2 (medio) · mes calendario contra periodo.** Retiré la frase «no son meses calendario: T-2 no los muestra en ninguna fila». La regla de ejemplos, al inicio del documento, ahora dice que oct-2026 a mar-2027 son P1 a P6. Cita C1-datos M-01 («P1 = oct-2026»), `is_canonical` («oct-2026 a mar-2027») y `docs/criterios/S2.md` S2-09 (nota de equivalencia: con la semilla dan las mismas cifras; una prueba con eventos entre los días 1 y 9 usa sus propias ventanas). La tarjeta 1, D-7 y A-06 citan «P1 = oct-2026». La fila «Total» de T-2 aclara que R-51 y R-52 son totales de seis filas (P1 a P6), no de doce. No hace falta bloqueo: C1 v1.2 y S2-09 ya lo deciden.
- **H-3 (medio) · resumen de DESIGN §7.7.** §4 tiene una regla nueva: `<figure>` con `<figcaption>` (título y resumen), el resumen en `ui.text.secondary` debajo del título, solo con cifras de C1 y sin variaciones ni porcentajes calculados. Si falta una cifra, «—». En carga, vacío o error no hay resumen. Cada gráfica tiene su texto:
  - T-1: real y plan del último mes.
  - T-2: devengado y efectivo acumulados al último mes, sin rotular deuda; PROVISIONAL con las capas.
  - T-3: deuda viva al cierre del primer y del último mes, o por línea si hay filtro.
  - T-4a y T-4b: costo y kcal contra meta de la última semana.
  La desviación frente al ejemplo de §7.7 («subió 12 %») queda declarada y pedida a DESIGN (solicitud 5 de §15). Supuesto 9 nuevo (solo texto). A-11 prueba `<figcaption>` y que T-2 no use la palabra «deuda».
- Encabezado v1.2 con la lista de cambios; §11.1, §13 (U3-05, U3-08, U3-11 y U3-12), §14 y §15 actualizados.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/design/U3-dashboard.md` (modificado, v1.2; respaldo de la v1.1 fuera del repositorio, en el scratchpad de la sesión).
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U3-t2-acumulados.md` (nuevo).
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-U3-r2.md` (este reporte).

## Verificación
Las ediciones se aplicaron con `apply.py` (27 reemplazos exactos; antes de cambiar nada comprueba que cada texto viejo aparezca las veces esperadas; resultado «reemplazos: 27 bytes: 79341»). Revisé el documento con `verify.py` (`python3 -I`). Los dos scripts están en el scratchpad de la sesión.
- **Oráculo:** una copia de `oraculo/model.py` corrida en el scratchpad reproduce `model.json` (igualdad JSON: True). Con la ventana desde P1, las capas (a) y (b) coinciden: P1 973.75 / 1,947.50 … P6 13,103.78 / 14,362.48. En cada periodo, efectivo acumulado − devengado acumulado da `diferencia_acumulada` (aserción cumplida). Con la ventana desde P2: (a) P2 2,921.25 / 3,895.00 y P6 13,103.78 / 14,362.48; (b) P2 1,947.50 / 1,947.50 y P6 12,130.03 / 12,414.98 (= R-52 − 973.75 y R-51 − 1,947.50); la regla de la v1.1 daba cima 2,921.25 y 13,388.73, lo mismo que encontró el revisor.
- **U3-01:** `^### Tarjeta` = 6; tarjeta 4 sin cambios (K38, U1 §0.2).
- **U3-02:** 28 identificadores camelCase: 27 son `operationId` de los extractos y el 28.º es la palabra `operationId`. Métricas citadas: M-01 a M-11 y M-13 a M-17, 0 faltantes en C1-datos. 148 nombres en minúsculas; los que no están en los extractos son 5 iconos (`banknote`, `house`, `plus`, `utensils`, `wallet`), `tnum` y los 2 campos que U3 **pide** a C1 (`cumulative_accrued_cents`, `cumulative_cash_cents`). Los dos aparecen solo en la variante (a), marcada PROVISIONAL con su bloqueo, y en la solicitud de §15.
- **U3-03:** §7 tiene 18 filas y 0 celdas vacías; el resumen no aparece en carga, vacío ni error (§4).
- **U3-04:** §8 sin cambios.
- **U3-05:** Detalle, Tendencias y Calendario cubren 15 de 15. T-2 tiene capas del mismo origen (PROVISIONAL `U3-t2-acumulados`). 0 donas.
- **U3-06:** las menciones nuevas de Karen solo están en la tabla de T-2 («partes del hogar y de Karen»), que ya estaba permitida por §9 y A-13.
- **U3-07:** tarjeta 3 sin cambios.
- **U3-08:** 57 montos, 20 distintos; 18 son literales de `model.json` o `data.json`. Los otros 2 son sumas explícitas de referencias: $138.50 = K-03 (3,233.00 − 3,094.50) y $2,921.25 = efectivo de R-45 + R-46 (973.75 + 1,947.50). Esa suma está escrita en el documento. 0 montos sin origen.
- **U3-09:** 0 fechas ISO; 30 fechas con día de la semana y 0 con el día equivocado.
- **U3-10:** §6 tiene 18 filas y 0 celdas vacías.
- **U3-11:** 0 hex y 0 `rgb(`; 54 tokens citados y 0 sin definir en C2 o DESIGN. 4 bloques «**Resumen** (DESIGN §7.7)» (T-1 a T-4); `figcaption` aparece en §4, A-11, §13 y el encabezado.
- **U3-12:** §11, §14 (supuesto 9 nuevo) y §15. `ls docs/estado/bloqueos/U3-*` → `U3-t2-acumulados.md`.
- **H-2:** «no son meses calendario» = 0 apariciones; «P1 = oct-2026» = 3; «U3-t2-acumulados» = 8.
- No hay pruebas Go ni tokens que correr: U3 es un documento, y no toqué código, C2 ni el contrato.

## Supuestos
- Supuesto 9 de U3: redacción de los resúmenes de una frase y de los rótulos de la variante (b). Es solo texto y no cambia cifras, esquema, operaciones ni tokens.

## Bloqueos
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U3-t2-acumulados.md`: origen de las capas de T-2 (A: C1 agrega los dos acumulados; B: relativo a la ventana; C: saldos de apertura). Se decide en la puerta 2. Mientras tanto, las capas, sus etiquetas y su resumen quedan PROVISIONAL.

## Pendientes
- Nada dentro del alcance. Los hallazgos bajos L-1 a L-6 de la revisión r2 siguen abiertos; esta ronda no los tenía en su alcance.

## Solicitudes fuera de alcance
1. **C1** (`docs/contratos/openapi.yaml`, generado con `backend/internal/api/c1gen/build.py`; esquema `AccrualVsCash.months[]`; `docs/contratos/C1-datos.md` M-11; extractos regenerados): agregar `cumulative_accrued_cents` y `cumulative_cash_cents` (`x-calculado`) con el mismo origen que `cumulative_difference_cents`. Motivo: H-1 y opción A del bloqueo. Solo si la puerta 2 elige A.
2. **S2** (`docs/contratos/motores/S2-interfaz.md`, salida de devengado contra efectivo): entregar esos dos acumulados, que ya calcula para la diferencia. Mismo motivo y misma condición.
3. **DESIGN** (`docs/design/DESIGN.md` §7.7): cambiar el ejemplo «El gasto subió 12 % respecto a septiembre» o aclarar que la variación llega calculada desde C1, porque la interfaz no recalcula (U1 §0.2). Es la solicitud 5 de U3 §15. Si la puerta 2 elige B, además declarar en §8.3 la desviación «cima = devengado en la ventana».
4. **Índice de bloqueos** (`docs/estado/bloqueos/INDICE.md`): sumar `U3-t2-acumulados`. Ya no es cierto que el único bloqueo abierto sea `S1-parametros-avisos`.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Puerta 2, bloqueo `U3-t2-acumulados`: **opción A**. Mantiene la regla de DESIGN §8.3 (la cima es el devengado acumulado), hace que la gráfica y la tabla den la misma diferencia acumulada y deja comprobar las cifras canon con cualquier reloj. Cuesta dos campos calculados que el motor S2 ya suma para sacar la diferencia. La opción B no toca C1, pero hace que la gráfica y la tabla muestren dos diferencias distintas y que «Pagado por adelantado» deje de servir cuando se paga deuda de antes de la ventana.
- No propongo valores de K31: U3 no tiene parámetros operativos.

## Consumo aproximado
- Acciones de herramienta: 31.
- Leído: unos 118 KB (protocolo 6, decisiones §K 15, criterios U3 13, revisión r2 13, partes de U3 unos 30, extractos y C1-datos unos 6, S2.md y DESIGN.md por partes unos 10, reporte r1 3, salidas de verificación unos 10, más otros fragmentos).
- Escrito por mí: unos 34 KB (scripts de edición 17, `verify.py` 5, bloqueo 4.4 con su ajuste, este reporte 8). El script regeneró `U3-dashboard.md` completo (79.3 KB; el cambio neto es de +8.5 KB).
