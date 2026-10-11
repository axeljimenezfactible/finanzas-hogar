# Bloqueo S1-restricciones-optimizador
Contexto: Al formalizar el comparador de escenarios (reglas-negocio.md RN-11). La decisión C fija tres restricciones duras (colchón mínimo de disponible, máximo de compras vivas a la vez, deuda en cero antes de una fecha), pero no sus valores ni a qué aplican. Es complementario de `criterios-optimizador` (cómo se miden y normalizan los 4 objetivos) y de `criterios-umbrales-ausentes` punto 7 (tiempo máximo de la búsqueda).
Qué falta: (1) El valor de cada restricción y el «± días» de los escenarios. (2) ¿Son del hogar (una sola para todo), por crédito, o se fijan en cada compra planeada? (3) ¿El colchón se mide en la línea que se usa o en la suma del disponible de todas las líneas? (4) ¿Las compras vivas se cuentan en todos los créditos o solo en el que se usa? (5) ¿La deuda en cero es de todos los créditos o del que se usa?
Opciones:
A. Globales, con un valor por defecto en Ajustes que se puede cambiar en cada comparación. Colchón en la línea usada; compras vivas y deuda en cero en el crédito usado. Consecuencia: no cambia el esquema (claves `app_setting optimizer.*`).
B. Por crédito: columnas nuevas en `credit` (migración 00002). Consecuencia: más control y más campos en Ajustes «Créditos y condiciones».
C. Por compra planeada: columnas nuevas en `planned_item` (migración 00002). Consecuencia: cada plan lleva sus propias restricciones.
D. Globales sobre todo el hogar: colchón sobre la suma de las líneas, compras vivas y deuda en todos los créditos. Consecuencia: es lo más estricto y no cambia el esquema.
Impacto: S3 (optimizador) y U4 (Ajustes «Créditos y condiciones»). Las claves `optimizer.min_available_buffer_cents`, `optimizer.max_live_purchases`, `optimizer.zero_debt_by_date` y `optimizer.date_window_days` quedan sin valor.
Cómo seguí: Las definiciones de las tres restricciones quedaron FIRMES en RN-11, con ejemplos sobre fechas del oráculo. Los valores y el alcance quedaron PROVISIONALES. No agregué columnas por crédito ni por compra para no adelantar la decisión.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción A: restricciones globales en Ajustes › Créditos y condiciones, ajustables en cada comparación; el colchón se mide en la línea que se usa; compras vivas y deuda en cero se cuentan en el crédito que se usa. Los valores los captura Axel (sin valor por defecto inventado: mientras no estén capturados, el comparador avisa «faltan restricciones» y no recomienda).
- **Estado:** RESUELTO. Se aplica en: S1 1.2 (claves `optimizer.*` sin valor), S3, U4.
