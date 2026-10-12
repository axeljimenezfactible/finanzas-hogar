# revisor-U4-r5 · reporte de agente

Agente `revisor-U4-r5` · WF-1 · 11-oct-2026. Revisión independiente de `docs/design/U4-ajustes-a.md` 1.4 y `docs/design/U4-ajustes-b.md` 1.3 contra `docs/criterios/U4.md` 1.2 y `docs/plan/decisiones.md` §A–§K. No corregí nada de lo revisado.

## Hecho
- Leí el protocolo completo, las decisiones §A–§K completas (incluida la sección K) y los criterios U4 completos. También leí las dos partes de U4 completas, por partes con `sed -n`, y U1 §2.3, §2.5, §2.6, §5.2 y §5.3 con `sed -n`. Además: `REVISION-U4-r4.md`, `correccion-U4-r4.md`, las secciones de solicitudes y propuestas de `U4a-ajustes.md` y `U4b-ajustes.md`, el inventario `@param` de 00003 y las siembras de `app_setting` de 00001 a 00003 con `grep`.
- Escribí y corrí mi propio script `datos/tmp/revU4r5/rev.py` (`python3 -I`, solo la biblioteca estándar). Revisa secciones, celdas vacías, `operationId` contra los extractos, parámetros contra el esquema, atajos contra U1, hex, ISO, tokens, bloqueos, la muestra de valores por defecto y las cifras del oráculo. Lo borré al terminar.
- Comprobé que R4-1 está cerrado.
- Abrí R5-1 (media) y dos bajos nuevos (R5-2 y R5-3). Siguen abiertos los bajos R4-2, R4-3 y R3-2 a R3-6.
- Veredicto: NO APROBADO. Informe en `docs/estado/wf1/REVISION-U4-r5.md`.

## Archivos
- `docs/estado/wf1/REVISION-U4-r5.md`: creado.
- `docs/estado/wf1/revisor-U4-r5.md`: creado (este reporte).
- `datos/tmp/revU4r5/rev.py`: temporal, creado y borrado.

## Verificación
| Criterio | Estado | Evidencia (resumen; el detalle está en REVISION-U4-r5) |
|---|---|---|
| U4-01 | Cumple | 14 secciones con los nombres de la decisión G. 285 filas de tabla, 0 celdas vacías. |
| U4-02 | Cumple | 210 = 210 operaciones. 0 inexistentes. Las 73 operaciones de Ajustes están citadas. |
| U4-03 | Cumple | ≈ 53 elementos presentes en el documento (K21 está en P-26.02: R3-5, baja). |
| U4-04 | Cumple | Muestra al azar de 15 valores: los 15 tienen origen en §A–§K o en el esquema. 0 valores de K31 propuestos dentro del diseño. |
| U4-05 | Cumple, con R5-1 (media) | Caso A/B con las series del oráculo y sus mínimos. 800.50 y 973.75. K40. RN-17. S4-07. RN-15. RN-18. Hueco nuevo: R5-1 (`Ctrl+/` y arrastre). |
| U4-06 | Cumple | Atajos y comandos «/» dentro de U1 §5.2 y §5.3, sin duplicados. Vistas según U1 §2.5, salvo R3-4. |
| U4-07 | Cumple, con R4-2 | 15 mensajes en español de México. Tasas con 4 decimales (K12). |
| U4-08 | Cumple | 8 acciones críticas con `color.critical.*`. |
| U4-09 | Cumple | 0 hex. 21 tokens, todos existentes. 412 px y 1440 px en las 14 secciones. |
| U4-10 | Cumple | Estados en las 14 y regla única en §0.3. |
| U4-11 | Cumple, con R4-3 | 0 fechas ISO. SUPUESTOS presentes. 8 bloqueos, todos existentes. |
| U4-12 | Cumple | 15 `@param` en su sección con valor inicial y «Vuelve a…». 20 claves sembradas citadas. Los 4 de `params_json` van con lugar PROVISIONAL y cita de `S1-parametros-avisos`. |
| U4-13 | Cumple | Ajuste de la casa 0 y 0 y ajuste por persona. «Mantener» sin campo. La vista previa da 1,799. |

- Comando: `python3 -I datos/tmp/revU4r5/rev.py /home/deepsec/finanzas-hogar`. La salida está resumida en REVISION-U4-r5 §Verificaciones. Los falsos positivos del script (`address`, las claves con «= false» y los nombres largos de sección) los comprobé aparte con `grep` y un segundo script corto.

## Supuestos
- Tomé la lista de excepciones de §0.2 como cerrada, porque enumera los casos y la regla de `Ctrl+/` nombra «(P-26.02 y P-26.13)». Si se lee como ilustrativa, la parte de R5-1 sobre «guardan al elegir» se atenúa, pero quedan sin definir el arrastre y `Ctrl+/` de P-26.14 y P-26.04.

## Bloqueos
Ninguno nuevo.

## Pendientes
- Que el siguiente corrector atienda R5-1 (media), en `docs/design/U4-ajustes-a.md` §0.2 y P-26.04 y en `docs/design/U4-ajustes-b.md` P-26.14 y P-26.13. Los bajos R5-2, R5-3, R4-2, R4-3 y R3-2 a R3-6 son opcionales.

## Solicitudes fuera de alcance
Ninguna.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna. La revisión no pide valores operativos de K31.

## Consumo aproximado
- 27 acciones de herramienta.
- ≈ 290 KB leídos.
- ≈ 33 KB escritos: script ≈ 11 KB, informe ≈ 15.5 KB y este reporte ≈ 5 KB.
