# correccion-criterios-r1 · P1 · 11-oct-2026
Corrige los hallazgos H-01 a H-06 (media) de `docs/estado/p1/REVISION-criterios-r1.md`. No reabre decisiones de §A–§K; no hay valores nuevos sin decisión, así que no hay bloqueos.

## Hecho
- **H-01 (U2-06)**: se quita «4.80 (gasto 25°)» de los valores de partida; el texto remite a K26 (gasto en terracota 40–45°, contraste medido por el script de colorimetría 1.1). §5: la línea de K26 ahora traza también a U2-06; viñeta nueva en la lista de cambios.
- **H-02 (U4-03, U4-04)**: U4-03 suma por sección Presupuesto (alcance, K2), Notificaciones (candidatos y confirmación del chat, K4, sustituye al «registro del chat_id»), Créditos y condiciones (tasas como se capturaron K12; restricciones sin valor por defecto con «faltan restricciones» K16; parámetros de tarjeta y crédito de tienda con «falta capturar» K18; interruptor de meses 7–12 K21), Personas y reparto (tope por persona K13) y Salud (objetivo y meta sugerida/confirmada K23); la lista pasa de ≈ 45 a ≈ 53 elementos. U4-04 suma los defectos (alcance «de este mes en adelante», interruptor apagado, restricciones/tarjeta/tope/objetivo sin valor). U4-07 muestra tasas como se capturaron. §3 y §4 ajustadas; 8 trazas nuevas (K2, K4, K12, K13, K16, K18, K21, K23) en la §5; la lista de cambios deja de decir «sin cambios».
- **H-03 (U1-09)**: se agregan K2, K3, K4, K13, K16, K18, K21 y K23, cada una con `P-xx` y `N-xxx`; el umbral de la lista pasa de ≥ 20 a ≥ 28 capacidades; 8 trazas nuevas en la §5; viñeta en la lista de cambios.
- **H-04 (U5, nuevo U5-13, alta)**: Playwright con la red cortada captura un gasto y un peso, muestra «Sin conexión» y el contador, sobrevive a recargar, vuelve la red y deja 1 registro por captura con `Idempotency-Key` y `captured_at`; reintento sin duplicar y fallo de envío que conserva la cola. U5-07 incluye el texto «Sin conexión» y su contador. §3, §4 y §5 (traza a K9) actualizadas.
- **H-05 (S1-09, U3-05)**: S1-09 escribe los límites de K8 (quincena 1 del 1 al 15, quincena 2 del 16 al último día, el 31 en la segunda; semana lunes a domingo) separados del día de pago; la Verifica (a) comprueba inicio y fin de las 12 quincenas (`2026-10-16`…`2026-10-31`, `2027-02-16`…`2027-02-28`) y la semana. `tools/s109_sql.sh` ahora imprime y comprueba esos límites y la semana. U3-05 separa quincena (K8) y día de pago. Trazas de K8 en S1 y U3.
- **H-06 (S8-06)**: el punto (1) queda «cobertura ≥ 85 % en cada uno de los 5 paquetes de K25 (S6-02); los demás se miden sin umbral; oráculo 100 %»; traza de K25, §4 y lista de cambios.
- `INDICE.md`: filas C1 21 (19/2), U1 16 (10/6) (corrección de la revisión) y U5 13 (10/3) (U5-13); total 257 (197/60).

## Archivos
/home/deepsec/finanzas-hogar/docs/criterios/{U2,U4,U1,U5,S1,U3,S8,INDICE}.md
/home/deepsec/finanzas-hogar/docs/criterios/tools/s109_sql.sh
/home/deepsec/finanzas-hogar/docs/estado/p1/correccion-criterios-r1.md

## Verificación
`bash docs/criterios/tools/verif_final.sh` (reejecutado completo tras los cambios):
- cuenta.py: S1–S8 y WF4 = 134 (105/29/0), «incompletos: ninguno».
- traza.py: 18 archivos, 257 criterios, secciones 1–5 completas, Verifica = Evidencia = criterios, 0 sin línea en la §5.
- refs.py: 257 definidos, 257 citados, 0 inexistentes.
- citas.py: 275 citas, 5 sin fuente única, todas anteriores (mezclan dos fuentes: C1, C2 ×2, U2, U5; hallazgo H-07 baja, fuera de esta tarea); ninguna cita nueva falla.
- N-xx residuos: 0. canon.py: 178 cifras, 0 diferencias, sha256 `1a5706…bfd4`. s109_sql.sh: 18 filas, límites K8 (31-oct-2026 en la quincena 2; feb-2027 16→28; 0 fuera de límites) y semana `2026-10-12` a `2026-10-18` OK. s109_equiv.py: 12 de 12, 0 diferencias. matriz_c102.py: igual que antes (no se tocó C1). Rutas `/tmp`: 0.
- verificar_u1.py, build.py, go build/vet, make verify-color: no aplican (no se tocó ninguno de esos artefactos; solo criterios).

## Supuestos
- El umbral «≥ 28 capacidades» de U1-09 es 20 + 8 capacidades nuevas; no cambia ninguna cifra del oráculo.
- Dónde vive cada campo nuevo en U4 (K21 en Créditos y condiciones; K13 en Personas y reparto) sigue la corrección sugerida por el revisor y los nombres de K.
- Las filas C1 y U1 del INDICE se corrigieron con el conteo real de la revisión (confirmado con traza.py).

## Bloqueos
Ninguno nuevo.

## Pendientes
- Hallazgos bajos H-07 a H-13 de la revisión (citas con dos fuentes, bloqueo de frontera S5-04, convención de 80 KB, marca del valor por omisión de C1-16, línea de K1 en la §5 de C1): fuera de esta tarea.

## Solicitudes fuera de alcance
- `docs/design/U1-arquitectura.md` / `U1-1.1`: U1-09 ahora exige las necesidades K2, K3, K4, K13, K16, K18, K21 y K23; el agente de U1 debe confirmar que cada una tiene `P-xx` y `N-xxx` (ruta: docs/design/U1-arquitectura.md; motivo: coherencia C1-17/C1-18).

## Paquetes instalados
Ninguno.

## Consumo aproximado
Acciones de herramienta: 27. Leído ≈ 85 KB. Escrito ≈ 35 KB.
