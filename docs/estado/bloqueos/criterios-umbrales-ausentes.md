# Bloqueo criterios-umbrales-ausentes
Contexto: Varios criterios (S6, S7, WF-4, S8, U5) necesitan valores que el plan no fija.
Qué falta: Lista de lo que falta: (1) hora exacta y contenido del resumen diario (noche) y del semanal (domingo) y del PDF mensual (día 1) [S7-04, S7-10]; (2) anticipación de la alerta «precios viejos antes del día de despensa» y definición de «plan de comidas sin armar» (¿cualquier celda vacía de la semana siguiente?) [S7-08]; (3) hora y retención (cuántas versiones) del respaldo diario [WF4-04, S6-09]; (4) lista exacta de cabeceras de seguridad [S8-02]; (5) límites de tamaño y tipo de adjuntos [S6-10]; (6) tamaño máximo del PDF (más allá del límite de la API de Telegram) [S7-10]; (7) tiempo máximo aceptable de la búsqueda exhaustiva del optimizador [S3-12]; (8) permisos de archivos de la base y respaldos [WF4-08].
Opciones:
- (a) Que cada agente proponga el valor en su reporte con justificación y Axel los apruebe en la puerta correspondiente.
- (b) Axel fija los 8 valores ahora.
- (c) Valores por defecto estándar documentados (cabeceras OWASP, 90 días de retención, etc.) marcados PROVISIONAL.
Impacto: Esos valores quedan PROVISIONAL en cada criterio afectado; ninguno bloquea las partes no dependientes.
Cómo seguí: Los criterios piden que el valor esté escrito y marcado PROVISIONAL hasta aprobación; nadie los inventa en silencio.
