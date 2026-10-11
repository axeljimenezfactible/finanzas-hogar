# revisor-criterios-r2 · P1 · 11-oct-2026
Estado: **terminado** · cerrado_por_tope: no · Veredicto: **NO APROBADO**. Hay 0 hallazgos altos, 1 medio (R2-M01) y 6 bajos nuevos. Informe: `docs/estado/p1/REVISION-criterios-r2.md`.

## Hecho
- Leí el protocolo completo y §K completa. Del INDICE leí completo el archivo (incluye la §1). De `docs/estado/wf0/REVISION-r3.md` leí las líneas de R2-07. Leí los reportes `criterios-S-1.1.md`, `criterios-UC-1.1.md` y `correccion-criterios-r1.md`, y el índice de hallazgos de `REVISION-criterios-r1.md`. De `docs/criterios/tools/` leí `verif_final.sh`, `cuenta.py` y `traza.py`.
- Escribí mis propios scripts en `datos/tmp/rev-crit-r2/`:
  - `estructura.py`: revisa por criterio, no por total, que haya exactamente 1 Verifica y 1 Evidencia. Revisa también las secciones 1–5 en orden, la §5 con rangos, los huecos y duplicados de ID, la versión 1.1, la tabla del INDICE y las rutas `/tmp`.
  - `kmat.py`: matriz K → archivo:línea(criterio).
  - `kw.py`: 64 comprobaciones de palabras clave por decisión K.
  - `residuos.py`: busca términos viejos de §A–§J que §K reemplazó.
  - `crit.py`: imprime criterios completos.
- Reejecuté `bash docs/criterios/tools/verif_final.sh` completo.
- Leí completos 40 criterios: los de los 6 hallazgos anteriores y los de K8, K9, K10, K11, K16 a K24, K26, K28, K30 y K31.
- Recalculé desde `oraculo/model.json` las cifras nuevas de S2-07, S2-08, S2-09, S4-07 y S5-09.
- Comprobé que S2-06 no cambió en esta ronda (`git diff HEAD -- docs/criterios/S2.md`: 0 líneas de S2-06) y que C1.md no cambió desde la r1 (fecha de modificación 11:31, antes de la r1 de las 11:53).

## Archivos
- Creados: `/home/deepsec/finanzas-hogar/docs/estado/p1/REVISION-criterios-r2.md` y `/home/deepsec/finanzas-hogar/docs/estado/p1/revisor-criterios-r2.md`.
- Temporales en `/home/deepsec/finanzas-hogar/datos/tmp/rev-crit-r2/`, ya borrados. No modifiqué ningún criterio.

## Verificación
- `python3 -I estructura.py`:
  - 18/18 archivos con las 5 secciones en orden.
  - 257 criterios (197/60/0), cada uno con 1 Verifica y 1 Evidencia. 0 sin línea en su §5, 0 huecos, 0 duplicados.
  - Las 18 filas del INDICE y el total 257 (197/60) coinciden. S y WF4: 134 (105/29). U, C y P3: 123 (92/31).
  - Rutas `/tmp` en comandos: 0. INDICE.md:7 todavía dice «248».
- `bash docs/criterios/tools/verif_final.sh`:
  - `cuenta.py` 134 (105/29), incompletos ninguno. `traza.py` 257, 0. `refs.py` 257/257.
  - `citas.py` 275 citas, 5 sin fuente (C1, C2 ×2, U2, U5).
  - `canon.py` 178, 0 diferencias.
  - `s109_sql.sh`: UP OK, 54 tablas, 18 filas; límites de K8 bien, 0 quincenas fuera; semana del 12 al 18 de octubre; CHECK rechaza la ventana 10→9.
  - `s109_equiv.py` 12/12. Rutas `/tmp` en S y WF4: 0.
  - `matriz_c102.py`: N-316, N-318 y N-319 sin operación, porque C1 v1.1 todavía no se regenera. Eso es del artefacto, no de los criterios.
- `python3 -I kmat.py` y `kw.py`:
  - 32 de las 33 decisiones K que piden criterio están aplicadas con línea (tabla §3.2 del informe). K8 está aplicada en S1-09, C1-06, C1-07, U3-05 y U1-16, pero **no en S2-06**.
  - R2-07 está aplicado en sus 6 partes. Las paráfrasis siguen abiertas como H-07.
- `python3 -I residuos.py`: «(15 y 30, o último día)» queda solo como día de pago, salvo en S2-06 y S1-03, más textos viejos en el INDICE. Sin «25°» de gasto vigente, sin Tailscale abierto, sin `SameSite` distinto de Strict y sin Go 1.26 en los criterios.
- Cifras (Python sobre `model.json`): S2-07 da 2,921.25 / 973.75 / 3,776.06 / 1,258.70 / 3,776.06 / 1,258.70 / 6,293.44 (= `horizonte.deuda_viva_10_abr`). S2-08 da 31.41 / 64.92 / 89.91, con 219.88 contra 219.87 y 629.35 contra 629.37. S2-09: pagos de la despensa 4 del `2027-04-30` al `2027-06-15`. Todo coincide.

## Supuestos
- R2-M01 es media, la misma severidad que la r1 dio a H-05: es el mismo texto ambiguo frente a K8, y aquí con la frase más fuerte «la quincena de fin de mes es el 30, no el 31».
- En la comprobación (b), «la lista» son las decisiones K que piden un criterio (K1 a K32 y K34; K33 es de proceso) más las 6 partes de R2-07.
- Las bajas R2-B04 (S1 sin los campos nuevos del esquema 1.2) y R2-B05 (umbral de U1-09) no bloquean, porque otra verificación ya cubre cada K.

## Bloqueos
Ninguno.

## Pendientes
Ninguno dentro de mi alcance. Para la ronda 3: comprobar R2-M01 y las bajas R2-B01 a R2-B06, H-07 a H-11 y H-13.

## Solicitudes fuera de alcance
- `docs/criterios/S2.md:39-40` (S2-06) y `docs/criterios/S1.md:29` (S1-03), para el agente de criterios S. Aplicar K8 como dice R2-M01: separar el periodo del día de pago y agregar un caso con un gasto del 31-oct-2026 en la quincena 2, con traza K8 y viñeta en la lista de cambios. Motivo: §K K8 FIRME y coherencia con S1-09.
- `docs/criterios/INDICE.md`, para el integrador. Reconciliar §1 (257), §5 notas 8 a 11, §6 (resueltos y los 2 bloqueos nuevos), §7 (fila K), §8 («≥ 2 meses») y §9 (flujo P1). Motivo: R2-B01.
- `docs/criterios/U3.md:6`, `U2.md:100`, `U1.md:51`, `C1.md:26` y `:30`, y `S1.md:61` y `:69`: correcciones bajas R2-B02 a R2-B06, con el detalle en el informe.

## Paquetes instalados
Ninguno.

## Consumo aproximado
- Acciones de herramienta: 29, con el cierre.
- Leído: unos 190 KB (protocolo, §K, INDICE, reportes, tools, salidas de mis scripts y 40 criterios). No leí ningún archivo grande entero.
- Escrito: unos 37 KB (scripts temporales, ya borrados, unos 13 KB; informe unos 17 KB; este reporte unos 7 KB).
