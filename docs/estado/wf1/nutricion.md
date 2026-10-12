# WF-1 · agente nutricion · reporte final (11-oct-2026, reanudación)

## Hecho
- Reanudé el intento interrumpido a las 19:00. Comprobé la integridad del volcado de Open Food Facts con `gzip -t` (rc 0, completo; no hizo falta reanudar la descarga).
- `empates.tsv` no existía en ninguna carpeta; lo reescribí completo (160 productos, con confianza y nota).
- Leí las secciones vigentes de S4 §9.1 (ronda 4): `tabla.csv` lleva `estado`, `confirmo` y `fecha_confirmacion`; todo sale `propuesto`. La licencia va en `PROCEDENCIA.md` (nombre que fija el contrato); `fuentes.md` solo apunta a él.
- Filtré OFF en flujo (4,535,553 filas leídas, 16,851 de México) para los 8 empacados sin equivalente USDA y guardé solo esas 8 filas.
- `construir_tabla.py` regenera `tabla.csv`, `tabla_detalle.tsv` y `revisar.tsv`.

## Archivos (datos/semilla/nutricion/)
tabla.csv · empates.tsv · off_filas_usadas.tsv · construir_tabla.py · extraer_off.py · tabla_detalle.tsv · revisar.tsv · PROCEDENCIA.md · fuentes.md · raw/ (ignorado por git; su .gitignore no se tocó). También este reporte.

## Verificación
- `python3 -I construir_tabla.py` dos veces seguidas y `cmp tabla.csv` contra la copia: IDÉNTICA.
- Validación tipo cargador (script ad hoc): cabecera exacta de 12 columnas, 150 renglones, ids únicos y existentes en `productos.json`, decimales de 0 a 3 cifras, `id_fuente` no vacío (USDA dígitos; OFF 8 a 14 dígitos), `estado=propuesto`, `confirmo` y `fecha_confirmacion` vacíos, UTF-8 sin BOM, sin CR.
- Cobertura: 160 de 160 productos. 150 con renglón; 10 sin fuente justificada (no son alimentos: blanqueador, pasta dental, jabón de manos, acondicionador, desinfectante, detergente, jabón de trastes, papel higiénico, servilletas, toallas de papel).
- Conteo por fuente: USDA 142 (Foundation 43, SR Legacy 99) · OFF 8 · manual 0 · sin fuente 10. Por confianza (los 150 con renglón): alta 110, media 35, baja 5; las 40 media o baja van en `revisar.tsv` junto con los 10 «sin fuente» (50 líneas).
- 6 valores comprobados contra la API oficial de USDA (api.nal.usda.gov/fdc/v1/food/<fdc_id>, clave DEMO_KEY): 171413 aceite de oliva 884 kcal, 100 g grasa (coincide); 173418 queso crema 350 kcal, 6.15 P, 34.4 G, 5.52 C (tabla 34.44 de grasa en el CSV, la API redondea a 34.4); 175176 tilapia cruda 96 kcal, 20.1 P, 1.7 G (tabla 20.08 en CSV); 169736 pasta seca enriquecida 371 kcal, 13.0 P, 1.51 G, 74.7 C (tabla 13.04 y 74.67 en CSV); 2647441 queso Oaxaca (Foundation) Atwater específico 298 kcal, general 297 (tabla 298.114 = nutriente 2048); 2346396 avena en hojuelas (Foundation) Atwater específico 379 kcal (tabla 378.866). Los decimales extra del CSV son el valor exacto sin redondeo de la API.
- Revisión de coherencia: kcal contra 4P+9G+4C; solo difieren las especias y el polvo para hornear (fibra y sales no digeribles), esperado.

## Supuestos
- Las fuentes de Foundation que no publican energía (nutriente 1008) se leen por Atwater específico (2048) y, si falta, general (2047); coincide con lo que muestra la API.
- Para productos del catálogo sin marca (chipotle en adobo, chocolate de mesa, queso panela, café molido) escogí una marca de México de OFF y lo anoté; Axel decide en la puerta 2.
- Los 10 productos no alimenticios no llevan renglón en `tabla.csv` (el contrato exige fuente USDA, OFF o manual).

## Bloqueos
Ninguno nuevo. Sigue vigente el PROVISIONAL de S4 sobre quién confirma los `propuesto` (`docs/estado/bloqueos/S4-confirmacion-empates.md`).

## Pendientes
- Axel revisa `datos/semilla/nutricion/revisar.tsv` (40 empates: 35 de confianza media y 5 de baja) y pasa a `confirmado` los que acepte (puerta 2). Las 5 bajas: chile-habanero, milanesa-res, jamaica, chipotle-adobo, cafe-molido; (petit-suisse y yogur bebible quedaron en media pero con datos de OFF variables).
- El cargador Go (`backend/internal/recetas/nutricion/`, del motor) lee `tabla.csv`; nada más que hacer de mi lado.

## Solicitudes fuera de alcance
- Ninguna. (Nota: el catálogo cita 171562 para el cubo de caldo de pollo; usé 171563, «cubos», para la nutrición. Si se quiere un solo registro, la decisión es de catálogo/Axel.)

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Confirmar o corregir los 40 empates de `revisar.tsv`, sobre todo carne molida (¿qué % de grasa compra?), milanesa (¿empanizada?), jamaica (flor seca), café molido (OFF repite 396.6 kcal en tres marcas) y chipotle en adobo (OFF dice 49 o 73 kcal).
- No hay valores K31 en mi alcance.

## Consumo aproximado
Acciones de herramienta: ~27. Leído: ~190 KB (protocolo 7 KB, S4 §9 12 KB, productos.json ~20 KB resumido, listados USDA/OFF ~140 KB). Escrito: ~55 KB (empates.tsv 11 KB, PROCEDENCIA.md 6 KB, scripts y reporte; tablas generadas 40 KB).
