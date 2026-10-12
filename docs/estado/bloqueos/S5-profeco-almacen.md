# Bloqueo S5-profeco-almacen
Contexto: al especificar el cargador PROFECO (`docs/contratos/motores/S5-interfaz.md` §11, WF-1, agente `contrato-S4S5`, 11-oct-2026). S5-05 pide comparar `wc -l` de cada CSV contra `SELECT COUNT(*) FROM <tabla_profeco>`, pero el esquema 1.3 no tiene tabla de registros PROFECO: `importProfeco` (C1) escribe en `price`, `presentation`, `product` y `store`, y `price` exige una presentación del catálogo y es única por presentación, tienda, día y fuente.
Qué falta (tres preguntas para Axel):
1. ¿Dónde viven los registros crudos de PROFECO (millones de renglones, la mayoría sin producto del catálogo)?
2. ¿Dónde se guarda la liga de un producto PROFECO (producto, presentación, marca) con una presentación del catálogo?
3. Si hay varios registros de la misma presentación, tienda y día, ¿cuál se guarda en `price`?
Opciones:
- 1. A (PROVISIONAL): base local aparte `datos/profeco/qqp.sqlite` (tabla `profeco_record`, ignorada por git, fuera del esquema de la app); solo lo filtrado y ligado pasa a `price` como `profeco`. B: tabla `profeco_record` en una migración 00004 de la app (crece la base de producción y los respaldos). C: sin tabla: solo conteos en el reporte de cobertura (S5-05 se verifica contra el reporte, no contra SQL).
- 2. A (PROVISIONAL): archivo `datos/semilla/profeco/vinculos.csv` confirmado a mano. B: columna nueva en `presentation` (p. ej. `profeco_key`).
- 3. A (PROVISIONAL): el primero del archivo; los demás se cuentan como duplicados. B: el de menor precio. C: el promedio redondeado HALF_UP.
Impacto: dónde cuenta S5-05 y cuántos precios de referencia aparecen. PROFECO nunca compite (K22): ninguna cifra de la comparativa ni del canon depende de esto.
Cómo seguí: el contrato S5 especifica la opción A de las tres preguntas marcada PROVISIONAL con cita a este bloqueo; el formato de columnas, la codificación, el filtro y el reporte de cobertura no dependen de esto.
