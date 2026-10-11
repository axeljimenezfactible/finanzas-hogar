# Reporte WF-0 · agente «criterios» (10-oct-2026)

## Hecho
- Leí completos `docs/plan/protocolo-agentes.md` y `docs/plan/decisiones.md`; del plan, §1.2 (cifras canon) y §3 completo (3.1 a 3.8), §4.2–§4.4, §5 (solo agentes), §6, §7 y §9 con `grep -n` y `sed -n`.
- Escribí los criterios de aceptación de las 18 fases/ítems y el índice: 248 criterios numerados (189 de severidad alta, 59 media), cada uno con enunciado medible, comando o cálculo de verificación, evidencia esperada y severidad; más las secciones «No aceptable», «Umbrales numéricos» y «Trazabilidad» con decisiones.md.
- S2 y S3 incluyen la lista completa de cifras canon de plan.md §1.2 como 85 casos de regresión (R-01…R-85), generada con un script que **afirma** cada valor contra `oraculo/model.json` y `datos/semilla/data.json` (no se transcribió a mano); S5 incluye K-01…K-08 de la comparativa de tiendas. Reparto de las 596 hojas numéricas de `model.json`: 159 de S2 + 437 de S3.
- C1 incluye explícitamente: toda `N-xxx` (necesidad de datos) de U1 mapeada a una operación y toda pantalla `P-xx` con operación o declaración Local (C1-02; redacción corregida en la ronda 1, ver la sección final), todo endpoint respaldado por tablas de S1 (C1-03), convención única de errores, paginación, filtros, fechas y dinero (C1-04 a C1-08) y OpenAPI 3.1 válido con validador (C1-01).
- Escribí 16 bloqueos de huecos del plan en `docs/estado/bloqueos/criterios-*.md` y listé 11 hallazgos en `INDICE.md` §5.
- Calculé y comprobé por script: la identidad exacta devengado/efectivo/deuda viva en los 6 cierres (el invariante literal del plan no se cumple en P1, P3, P5 ni el 10-abr-2027), los contrastes estructurales de plan.md §4.3 con los hex de C, y el origen de cada cifra del reporte de tiendas.

## Archivos
Directorio `/home/deepsec/finanzas-hogar/docs/criterios/` (19 archivos):
- `INDICE.md`, `S1.md`, `S2.md`, `S3.md`, `S4.md`, `S5.md`, `S6.md`, `S7.md`, `S8.md`, `U1.md`, `U2.md`, `U3.md`, `U4.md`, `U5.md`, `P3-claude-design.md`, `C1.md`, `C2.md`, `C3.md`, `WF4-operacion.md`.

Bloqueos en `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/` (16): `criterios-invariante-deuda.md`, `criterios-promedios.md`, `criterios-proyeccion-12m.md`, `criterios-credito-formulas.md`, `criterios-optimizador.md`, `criterios-frescura-peso.md`, `criterios-ocr-tickets.md`, `criterios-meta-kcal.md`, `criterios-login-bloqueo.md`, `criterios-bind.md`, `criterios-token-servicio.md`, `criterios-toques.md`, `criterios-nxx.md`, `criterios-deltal.md`, `criterios-cobertura.md`, `criterios-umbrales-ausentes.md`.

Este reporte: `/home/deepsec/finanzas-hogar/docs/estado/wf0/criterios.md`.

## Verificación
- `ls /home/deepsec/finanzas-hogar/docs/criterios | wc -l` → 19 (18 ítems + INDICE).
- Script `python3 -I` de estructura: cada uno de los 18 archivos tiene las 5 secciones (Objetivo, Criterios, No aceptable, Umbrales numéricos, Trazabilidad); en cada archivo el número de líneas «Verifica» y «Evidencia» es igual al número de criterios (248 = 248 = 248).
- Script de referencias cruzadas entre IDs (`S1-12`, `U1-04`, `C3-05`, etc.): 248 definidos, 0 referencias a IDs inexistentes (corregí dos referencias erróneas antes de cerrar: U3-02→C1-09 y U5-11→C3-05).
- `python3 -I gen_canon.py` (en el scratchpad): 85 filas R + 8 filas K generadas con `assert` contra `model.json`/`data.json` (incluye pagos 1,258.68 ×3 + 1,258.70, 883.75/883.77, 2027-02-28, 14 valores de cada serie, 33 renglones idénticos, +138.50 y +100.50) → sin errores.
- Conteo de hojas numéricas de `model.json` por rama: periodos 122, valor_estable 12, horizonte 13, papel 12 (=159 S2); parametros 9, despensas 144, credito 143, sensibilidad 45, disposicion_por_despensa 68, interes_por_monto 28 (=437 S3); total 596.
- Identidad deuda viva = (devengado acum − efectivo acum) + saldo no devengado: cierra a 0.00 en P1–P6 (2,921.25; 973.75; 3,776.06; 1,258.70; 3,776.06; 1,258.70) y el 10-abr-2027 (6,293.44).
- Contrastes con hex de C (WCAG 2.x calculado en Python): blanco/marca600 3.98; marca700/fondo50 5.54; texto900/50 11.53; texto700/50 5.40; texto600/100 3.47; oscuro 100/950 14.07; 300/900 7.97; marca300/950 10.38; neutro600/50 (borde) 3.71 (el plan dice ≥ 4.2: hallazgo).
- Prueba de la consulta de S1-05 en una base temporal (detecta columnas REAL/FLOAT): funciona.
- `sha256sum oraculo/model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`, byte a byte igual al de `~/reporte-compra-latex/model/model.json` (`cmp`).
- **Efecto secundario que debo declarar:** durante la lectura ejecuté una vez `python3 -I oraculo/model.py`, que **regenera `oraculo/model.json`** como efecto de importar el módulo; verifiqué con `cmp` y sha256 que el archivo quedó idéntico byte a byte al original (solo cambió la fecha de modificación). No modifiqué `model.py`, `verify3.py` ni los JSON.

## Supuestos
1. ~~«N-xx» = identificador de pantalla o nodo del inventario de U1~~ (supuesto inicial; el plan no lo define). **Corregido en la ronda 1:** pantalla = `P-xx` (`P-26.xx` para las secciones de Ajustes) y necesidad de datos = `N-xxx`, como los fija U1 §0.1; el bloqueo `criterios-nxx` quedó resuelto con la opción (b).
2. El formato de evidencia (línea que empieza con el ID del criterio en `docs/estado/<flujo>/<agente>.md`) y las definiciones de severidad media/baja son convenciones de este directorio y no alteran cifras, esquema, endpoints ni tokens.
3. Los nombres de pruebas, scripts y rutas auxiliares los fija cada agente en su reporte (`<declarado por el agente>`).
4. Dos medidas operativas propias, señaladas en su criterio: tolerancia ±2° al medir matiz (U2-02) y «≥ 2 meses» para «sesión de meses» (S6-07).
5. Leí páginas del plan fuera de las indicadas solo para trazabilidad (§1.2 cifras canon, §6 pila, §7 riesgos, §8.2 plantilla CSV, §1.5/§1.6/§1.8); no cambian ninguna decisión.

## Bloqueos
Los 16 archivos listados en «Archivos» (`criterios-nxx.md` quedó **resuelto** en la ronda 1 de correcciones; los otros 15 siguen abiertos). Los de mayor impacto: `criterios-invariante-deuda.md` (el invariante del plan no se cumple literalmente), `criterios-token-servicio.md` (token de Telegram y usuario `finanzas`), `criterios-ocr-tickets.md` (no hay tickets de prueba) y `criterios-deltal.md` (ΔL ≥ 0.10 entre los 28 pares no es alcanzable).

## Pendientes
- Ninguno dentro de mi alcance. Pendiente del integrador/revisor WF-0: confirmar con Axel los 16 bloqueos en la puerta 1 y, según su respuesta, actualizar los criterios afectados (lista en `INDICE.md` §6).

## Solicitudes fuera de alcance
1. `/home/deepsec/finanzas-hogar/oraculo/verify3.py` línea 89: escribe en `/home/deepsec/reporte-compra-latex/model/verify3.json` (fuera del repositorio, contra protocolo §7). Cambiar a una ruta dentro de `oraculo/` en WF-1 (S2-04, S3-16). Motivo: los criterios exigen 0 escrituras fuera del repo.
2. `/home/deepsec/finanzas-hogar/oraculo/model.py`: la regeneración de `model.json` al importar el módulo es un riesgo para el canon; pasarla a ejecución explícita (`__main__`) o protegerla con comparación de hash. Motivo: S2-04 y S3-16.
3. `docs/plan/plan.md` §3.4 fila S2: precisar el invariante «devengado acumulado − efectivo acumulado = deuda viva» cuando Axel decida el bloqueo `criterios-invariante-deuda`.
4. `docs/plan/plan.md` §4.3: el contraste «neutro 600 ≥ 4.2:1» de los bordes de campo no se reproduce para C (3.71:1); corregir el texto o el paso de neutro usado (U2-05).
5. Crear los directorios `docs/estado/wf1`, `wf3` y `wf4` al arrancar esos flujos (los criterios los usan como destino de evidencia).

## Paquetes instalados
Ninguno.

## Consumo aproximado
- ≈ 70 acciones de herramienta (lecturas del protocolo, decisiones y plan, scripts de verificación, escritura de 19 archivos de criterios, 16 bloqueos y este reporte).
- Leído: ≈ 70 KB (protocolo 4.7 KB, decisiones 11 KB, partes del plan ≈ 40 KB, `model.json`/`data.json` por script).
- Escrito: ≈ 1,950 líneas (≈ 250 KB) de criterios, ≈ 16 bloqueos cortos y este reporte.
- Cerrado por tope: no.

## Correcciones de la ronda 1 de revisión (10-oct-2026, agente `correccion-criterios-r1`)
Origen: `docs/estado/wf0/REVISION-r1.md` §5, hallazgos H-06, H-07 y H-08. Detalle, verificación y consumo: `docs/estado/wf0/correccion-criterios-r1.md`. Los criterios pasan a la versión 1.1 (17 archivos) con el mismo conteo: 248 (189 alta, 59 media). Resumen en `docs/criterios/INDICE.md` §9.
- **H-06:** los criterios usaban «N-xx» como identificador de pantalla; ahora pantalla = `P-xx` (`P-26.xx` para Ajustes) y necesidad de datos = `N-xxx`. C1-02 cruza `N-xxx` ↔ `operationId` y exige que cada `P-xx` tenga operación o la declaración Local. Bloqueo `criterios-nxx` resuelto con la opción (b).
- **H-07:** S1-09 pide meses calendario oct-2026…mar-2027 más la prueba de equivalencia con P1–P6 del oráculo (12 de 12 cifras) en lugar de insertar las ventanas del 10 al 9; S2-09 y las notas de R-45…R-50 dicen que la equivalencia depende de la semilla.
- **H-08:** los 35 criterios sin línea en su sección 5 la tienen ahora (248 de 248).
- Sin cambios en las 85 cifras canon R-01…R-85 ni en K-01…K-08 (filas idénticas byte a byte; se reverificaron contra `model.json` y `data.json`).
