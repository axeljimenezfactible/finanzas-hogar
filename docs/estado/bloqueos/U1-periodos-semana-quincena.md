# Bloqueo U1-periodos-semana-quincena

Contexto: al definir las necesidades de datos de U1 (`docs/design/U1-arquitectura.md` §0.2 y §3) que agrupan por semana (plan de comidas N-120, tarjeta «Comidas de la semana» N-023, promedios semanales N-042, tendencias de 12 semanas N-129 y N-233) y por quincena (tarjeta «Flujo quincenal» N-021, flujo por quincena N-090, plan y real por quincena N-082). Las decisiones fijan el mes calendario como periodo base y la quincena «15 y 30, o último día del mes» como subperiodo, pero no dicen qué días abarca cada semana ni cada quincena.

Qué falta: dos preguntas para Axel.
1. ¿La semana va de lunes a domingo o de domingo a sábado?
2. ¿Qué días abarca cada quincena: del 1 al 15 y del 16 al último día del mes, o de fecha de pago a fecha de pago (del 15 al 29 y del 30 al 14 del mes siguiente)? En particular, ¿a qué quincena pertenece el día 31?

Opciones:
- Semana A, lunes a domingo (norma ISO 8601). Consecuencia: encaja con el resumen semanal del domingo (cierra la semana ese día) y con el aviso del domingo de «plan de comidas sin armar», que se referiría a la semana que empieza el lunes siguiente.
- Semana B, domingo a sábado. Consecuencia: el resumen del domingo cubriría la semana anterior (domingo a sábado) y el aviso de plan sin armar llegaría el mismo día en que empieza la semana.
- Quincena C, del 1 al 15 y del 16 al último día del mes (quincena calendario). Consecuencia: el ingreso y los pagos del 15 y del 30 quedan dentro de la quincena que cierran; el día 31 pertenece a la segunda quincena del mes; cada quincena es una mitad del mes, lo que facilita sumar al mes.
- Quincena D, de pago a pago: del 15 al 29 y del 30 (o último día) al 14. Consecuencia: el ingreso abre la quincena y paga lo que viene después; el día 31 cae en la quincena que empezó el 30; una quincena cruza dos meses, lo que complica comparar con el presupuesto mensual.

Impacto: cambian las cifras de promedios semanales, costo semanal de comidas, tendencias por semana, flujo por quincena y la tarjeta «Flujo quincenal». No cambian las pantallas, los flujos ni las necesidades de datos (solo sus fechas de corte).

Cómo seguí: en U1 usé Semana A y Quincena C marcadas **PROVISIONAL** (§0.2 del documento). Ninguna cifra se calculó en U1 con esas fechas. S1 y S2 deben usar la misma respuesta en `docs/contratos/reglas-negocio.md`.
