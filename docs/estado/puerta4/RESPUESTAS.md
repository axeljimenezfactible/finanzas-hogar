# Puerta 4 · Respuestas de Axel

Registro del orquestador (Fable 5.1) · 13-oct-2026, 23:40 · sobre `PUERTA-4.md` (commit e365a3f).

## Modo elegido

Axel contestó por AskUserQuestion: **«Lo recomendado + excepciones»** (la misma regla de la puerta 2): toda pregunta con recomendación queda contestada con la opción recomendada; toda pregunta **sin recomendación** sigue **PROVISIONAL** con la opción en uso hasta que Axel la conteste aquí o en la puerta 5. Las excepciones y respuestas directas van abajo.

## Respuestas directas (13-oct-2026, 23:35)

| Id | Respuesta de Axel (opción elegida) | Nota |
|---|---|---|
| Q4-01 | **A · Aprobar y arrancar WF-4** | Texto de la opción: «WF-4 instala el servicio systemd, respaldos, nombre y runbook esta noche. Lo pendiente (alertas, tus respuestas) va a una ronda acotada P4c después y se redespliega con ops/instalar.sh» |
| Q4-02 | **A · sí, al contestar** | Incluido en la misma respuesta |
| Q4-80 | **B · Con los 36 precios** | «El comparador y "Mejor canasta" funcionan desde el primer día; se marcan viejos el 9-nov-2026 y tú o el agente externo los renuevan» |
| (nuevo) Alertas sin fuente | **Sí, en P4c después de WF-4** | Hallazgo del redactor de la puerta 4 (§1.5): presupuesto 80 %/100 %, precios viejos y día de despensa responden «sin fuente» en `backend/internal/reportes/avisos/fuente.go`. Ronda acotada P4c (2–3 agentes) tras WF-4; se redespliega con `ops/instalar.sh` |

## Hoja aplicada (regla «lo recomendado + excepciones»)

| Q4 | Tema | Recomendada | Respuesta aplicada |
|---|---|---|---|
| Q4-01 | Aprobar el paso a operación tal cual (8ace01e) | sin recomendación | A · aprobar el paso a operación tal cual (8ace01e) |
| Q4-02 | Arrancar WF-4 al contestar | sin recomendación | A · WF-4 arranca al contestar |
| Q4-03 | WF-4 carga `gastos_fijos_ingresos.csv` si existe | A (plan l. 186) | la recomendada |
| Q4-04 | S2-ahorro-quincenal · reparto, metas y saldo | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-05 | S2-metas-vencida · aportación de una meta vencida | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-06 | S4-presentacion-respaldo · presentación sin precio (Q2-07) | sin recomendación (en uso: regla del motor) | PROVISIONAL (opción en uso) |
| Q4-07 | S5-carritos-granel · redondeo a granel | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-08 | C1-hora-reportes-job · qué manda en la hora del reporte | A | la recomendada |
| Q4-09 | C1-reparto-alcance · alcance de «todos» | B | la recomendada |
| Q4-10 | S1-reparto-edicion-gasto · editar un gasto viejo | sin recomendación (en uso C) | PROVISIONAL (opción en uso) |
| Q4-11 | S3-formula-tarjeta-pdf · entregar 2 PDF | sin opciones (entrega) | la recomendada |
| Q4-12 | S3-iva-lectura · préstamo y crédito sin dato de IVA | sin recomendación (en uso 1A) | PROVISIONAL (opción en uso) |
| Q4-13 | S4-empates-sin-confirmar · empates de nutrición (Q2-90) | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-14 | S5-motor-oraculo-reglas · 8 reglas motor contra oráculo | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-15 | S5-precio-vigente-captura · 4 reglas del precio vigente | sin recomendación (en uso: el motor) | PROVISIONAL (opción en uso) |
| Q4-16 | U4-reparto-preseleccion · opción preseleccionada | A | la recomendada |
| Q4-17 | S5-agente-externo (a) · peso de sus precios | A (compiten como manual) | la recomendada |
| Q4-18 | S5-agente-externo (b) · qué escribe | A | la recomendada |
| Q4-19 | S5-agente-externo (c) · recetas como borrador | A | la recomendada |
| Q4-20 | S5-agente-externo (d) · valor de `price.source` | d2 | la recomendada |
| Q4-21 | S5-agente-externo (e) · credencial y a nombre de quién | e1 + e4 | la recomendada |
| Q4-22 | S5-agente-externo (f) · borrador y origen de receta | sin recomendación (en uso f2) | PROVISIONAL (opción en uso) |
| Q4-23 | S5-agente-externo (g) · tope por lote | sin recomendación (en uso g1) | PROVISIONAL (opción en uso) |
| Q4-24 | P3-contraste-chip-hover | A | la recomendada |
| Q4-25 | P3-huecos-componentes · supuestos de DESIGN | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-26 | S6-semilla-despensas-futuras | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-27 | S6-semilla-sucursales-municipio | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-28 | S6-semilla-catalogo-alias | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-29 | K31 · respaldo: 03:15 y 30 diarios | (a) | la recomendada |
| Q4-30 | K31 · qué conserva una restauración | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-31 | K31 · tipos de adjunto | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-32 | Adjuntos · MB, `kind` y 409 | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-33 | K31 · cabeceras de seguridad y CSP | (a) | la recomendada |
| Q4-34 | K31 · argon2id, límites, purga, rutas, campo desconocido | (a) | la recomendada |
| Q4-35 | K31 · contenido de los resúmenes de Telegram | sin opciones escritas | la recomendada |
| Q4-36 | S7 · día de despensa | C | la recomendada |
| Q4-37 | S7 · semana que revisa la alerta 4 | A (la siguiente) | la recomendada |
| Q4-38 | K31 · cadencia y reintentos de avisos | (1) PROPUESTA; (2) A; (3) sin recomendación | la recomendada |
| Q4-39 | K31 · reporte mensual: secciones y negativos | (1) A; resto sin recomendación | la recomendada |
| Q4-40 | K31 · PWA sin conexión (15 s, 30 s, copia, cola) | los 4 valores propuestos | la recomendada |
| Q4-41 | K31 · cómo se sirve la PWA | A en producción, B en desarrollo | la recomendada |
| Q4-42 | S6-sistema · 8 puntos de autenticación y sistema | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-43 | S7-convivencia-multix · chat de Ximena | A | la recomendada |
| Q4-44 | S7-ruta-token | A | la recomendada |
| Q4-45 | S7-lualatex | A | la recomendada |
| Q4-46 | S7-estado-cuenta-fechas | A | la recomendada |
| Q4-47 | WF3-integrador-red · IP, nombre y CA | 1A, 2C (+A), 3A | la recomendada |
| Q4-48 | S6-ajustes-api · 6 puntos | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-49 | S6-comidas-calculos | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-50 | S6-compras-a-motores | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-51 | S6-compras-b-fuentes | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-52 | S6-compras-b-nutricion | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-53 | S6-creditos | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-54 | S6-gastos-inicio | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-55 | S6-importar-partidas-csv · columnas del CSV | sin opciones escritas | la recomendada |
| Q4-56 | S6-presupuesto-motor | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-57 | S6-presupuesto-periodos | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-58 | S6-presupuesto-quincenas | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-59 | C3-matriz · ratificar 1, 3c y 4a | ratificar | la recomendada |
| Q4-60 | WF3-ajustes-a · despensa programada, membresía, P-26 | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-61 | WF3-ajustes-compra-ejemplo | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-62 | WF3-ajustes-muestra-tema | A | la recomendada |
| Q4-63 | WF3-ajustes-precios-conteo | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-64 | WF3-ajustes-registro | A | la recomendada |
| Q4-65 | WF3-comidas-compras-interfaz · 7 puntos | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-66 | WF3-creditos-salud-reportes · 7 puntos | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-67 | WF3-frontend-rutas · direcciones en español | A | la recomendada |
| Q4-68 | WF3-gastos-contrato · 4 huecos de C1 | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-69 | WF3-gastos-ticket-credito | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-70 | WF3-inicio-tiendas · color de tienda | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-71 | WF3-inicio-vacios · 6 reglas de vacío | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-72 | WF3a-bis-ajustes-a · edición en grupo | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-73 | WF3a-bis-comidas-compras · 4 puntos | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-74 | WF3a-bis-gastos-presupuesto · 7 puntos | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-75 | Semilla de desarrollo · renta (punto 1) | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-76 | Semilla de desarrollo · presupuesto de otras categorías (punto 3) | sin opciones escritas | la recomendada |
| Q4-77 | Semilla de desarrollo · costo y kcal en 0 (punto 5) | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-78 | Semilla real · gastos proyectados 2 a 6 (punto 4) | (a) | la recomendada |
| Q4-79 | Candado de la semilla de desarrollo (punto 2) | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-80 | Q2-49 contra 36 precios de Walmart en la semilla | sin recomendación (en uso B) | B · se quedan los 36 precios de Walmart del 10-oct-2026 en la semilla de producción |
| Q4-81 | Canasta del 10-oct: +$264.50 o +$138.50 | A (+$264.50, K46) | la recomendada |
| Q4-82 | U5-ronda3-escalamiento | A | la recomendada |
| Q4-83 | C3-ronda3-escalamiento | A | la recomendada |
| Q4-84 | C3-flujos punto 2 · F-3 ruta A | (c) | la recomendada |
| Q4-85 | C3-flujos punto 4 · F-5 sep-2026 | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-86 | C3-flujos puntos 5, 6 y 7 | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-87 | C3-capturas-global | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-88 | U5-conteo-carpeta | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-89 | U5-contraste-ui-controles | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-90 | U5-nombres-propios | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-91 | U5-tipografia-tres-tamanos | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-92 | Barómetro · ingreso mensual neto | (a) | la recomendada |
| Q4-93 | Barómetro · pagos que cuentan | (a) | la recomendada |
| Q4-94 | Barómetro · redondeo, escala, mes y §8.9 | (a) en 3 y 4; 5 y 6 sin recomendación | la recomendada |
| Q4-95 | Horizonte de las gráficas | (a) | la recomendada |
| Q4-96 | «Más del hogar» · tarjetas de U3 1.x | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-97 | Pagos próximos · N, rango y forma | (a) en 1 y 2 | la recomendada |
| Q4-98 | Sugerencia de presupuesto · regla | (a) en 1, 2 y 3 | la recomendada |
| Q4-99 | C1-tablero · 5 puntos | (a) en 1 y 2 | la recomendada |
| Q4-100 | S6-tablero-sin-datos | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-101 | S6-tablero-senal-periodo | sin recomendación (en uso a) | PROVISIONAL (opción en uso) |
| Q4-102 | Líneas archivadas en el tablero (y punto 5) | (a) y (a) | la recomendada |
| Q4-103 | C3-ajustes-tablero · Parámetros del tablero | A | la recomendada |
| Q4-104 | WF3c-barra-iconos-seccion · 4 iconos | sin recomendación | PROVISIONAL (opción en uso) |
| Q4-105 | WF3c-tablero-interfaz · 6 puntos | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-106 | WF3c-barometro-sin-dato | sin recomendación (en uso A) | PROVISIONAL (opción en uso) |
| Q4-107 | WF3c-ronda3-escalamiento · ratificar | ratificar | la recomendada |
| Q4-108 | S6-deuda-viva-hoy · deuda de hoy, G-5 y falta pagar | (a) en los tres | la recomendada |
| Q4-109 | C1-unidad-parametros-tablero | A | la recomendada |
| Q4-110 | U3-tablero-k4-personas | A | la recomendada |

Conteo: 50 contestadas con la recomendada, 57 PROVISIONALES, 3 directas.

## Respuesta del orquestador a la PREGUNTA del redactor (§3.1)

La cifra «28 sin recomendación» de `docs/estado/puerta2/RESPUESTAS.md` fue aritmética (98 preguntas − 70 contestadas con la recomendada), no una lista; incluía preguntas con recomendación parcial y las que Axel contestó una por una. La lista real de ids sin recomendación única es la de 19 que reconstruyó el redactor (17 abiertas + Q2-11 y Q2-64). Vale la lista de 17 de `PUERTA-4.md` §3.1.

## Lo que WF-4 toma de esta puerta

- Q4-29, Q4-33, Q4-34 = (a): respaldo diario 03:15 America/Mexico_City con 30 copias; cabeceras y CSP propuestas; argon2id y límites de cuerpo como están.
- Q4-41 = A en producción (dist incrustado) si el binario ya lo permite; si no, B con `FINANZAS_FRONTEND` y se anota como PROVISIONAL.
- Q4-44 = A (las dos rutas del token); Q4-47 = 1A, 2C, 3A recomendadas: WF-4 aplica lo que depende de esta máquina (unidad de avahi, /etc/hosts) y deja la reserva de IP y el DNS del módem como acción de Axel.
- Q4-80 = B: la semilla de producción lleva los 36 precios de Walmart.
- Q4-03 = A: carga `datos/semilla/gastos_fijos_ingresos.csv` si existe; hoy no existe.
- Lo sin recomendación (Q4-30, Q4-31, Q4-32, Q4-42, Q4-43 tiene recomendación A, etc.) corre con su opción PROVISIONAL.
