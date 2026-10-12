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

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-53 · (1) ¿Dónde viven los registros crudos de PROFECO?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «base local aparte `datos/profeco/qqp.sqlite` (tabla `profeco_record`, ignorada por git, fuera del esquema de la app); solo lo filtrado y ligado pasa a `price` como `profeco`.»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: dónde cuenta S5-05 y cuántos precios de referencia aparecen. PROFECO nunca compite (K22): ninguna cifra de la comparativa ni del canon depende de esto.
- **Qué cambia:** `S5-interfaz.md` §11 quita la marca PROVISIONAL; S5-05 se verifica contra la base local. HECHO de PUERTA-2 §3.18: hoy la base está en `datos/semilla/profeco/db/qqp.sqlite` y pesa 17 GB (el archivo de bloqueo cita `datos/profeco/qqp.sqlite`; el integrador alinea la ruta). PREGUNTA que sigue abierta fuera de este bloqueo: dónde vive esa base en producción (pendiente 6 de `INTEGRACION.md` §8); esta respuesta solo fija que va aparte del esquema de la app.
- **Estado:** RESUELTO

### Q2-54 · (2) ¿Dónde se guarda la liga entre un producto PROFECO y una presentación del catálogo?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «archivo `datos/semilla/profeco/vinculos.csv` confirmado a mano.»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: dónde cuenta S5-05 y cuántos precios de referencia aparecen. PROFECO nunca compite (K22): ninguna cifra de la comparativa ni del canon depende de esto.
- **Qué cambia:** `S5-interfaz.md` §11 quita la marca PROVISIONAL; S5-05 se verifica contra la base local. La liga sigue en `vinculos.csv` (sin columna nueva en `presentation`).
- **Estado:** RESUELTO

### Q2-55 · (3) Con varios registros de la misma presentación, tienda y día, ¿cuál se guarda?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «el primero del archivo; los demás se cuentan como duplicados.»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: dónde cuenta S5-05 y cuántos precios de referencia aparecen. PROFECO nunca compite (K22): ninguna cifra de la comparativa ni del canon depende de esto.
- **Qué cambia:** `S5-interfaz.md` §11 quita la marca PROVISIONAL; S5-05 se verifica contra la base local. 
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
