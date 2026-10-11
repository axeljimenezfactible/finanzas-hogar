# Bloqueo U2-gasto-error
Contexto: Al calcular los 8 roles del candidato C (U2a, colorimetría) con los matices de decisiones §I, gasto (25°) y error (20°) quedan a 5° de matiz con la misma luminosidad y croma. Claro: gasto `#A44541`, error `#A34449` (ΔE OKLab×100 = 1.2). Oscuro: gasto `#FC9A92`, error `#FC9999` (ΔE = 1.0). Con ΔE ≈ 1 la diferencia no se percibe: a efectos prácticos son el mismo rojo. Decisiones §I fija ambos matices y no menciona esta colisión (el plan §4.3 sí la señaló para el candidato B: «gasto se separa a 35°»).
Qué falta: ¿Aceptas que gasto y error compartan prácticamente el mismo rojo, distinguidos solo por contexto, signo e icono, o prefieres separarlos con alguna de las opciones 2 a 4?
Opciones:
1. Mantener gasto 25° y error 20° (es lo entregado). Gasto solo aparece como cifra con signo «−» en listas y tablas; error solo como mensaje de validación o aviso con icono de alerta y texto; nunca comparten componente. Consecuencia: ningún cambio; el color no ayuda a distinguirlos.
2. Mover gasto a 40–45° (terracota o naranja quemado): ΔE frente a error 4.6–5.5 en claro y 4.1–5.1 en oscuro. Consecuencia: gasto queda a 5–10° de Karen y de Bodega Aurrera (50°), así que hay que rehacer esa parte de la paleta categórica (el script lo recalcula en una iteración).
3. Gasto sin color de rol: cifra en texto neutro (neutral 900 en claro, 100 en oscuro) con signo «−»; el rojo queda solo para error. Consecuencia: `color.role.expense` pasa a alias de neutro o se elimina; las gráficas de gasto ya usan series o categóricas.
4. Gasto 25° con croma bajo (máximo 0.06–0.08, «terracota apagado»): ΔE frente a error 5.6–7.3. Consecuencia: rompe la regla de croma uniforme de los roles solo para gasto y su cifra se ve más apagada que ingreso.
Impacto: solo los tokens `color.role.expense.*` (y, con la opción 2, `color.person.karen.*` y `color.store.bodega-aurrera.*`). No hay nada marcado PROVISIONAL: lo entregado son los valores decididos.
Cómo seguí: Entregué gasto 25° y error 20° como están en decisiones §I, documenté las distancias en `docs/design/colorimetria.md` §6.5 y §6.6 y fijé la regla de uso de la opción 1. Cambiar de opción es editar el matiz o croma de gasto en `docs/design/tools/colorimetria.py` y regenerar; el script vuelve a verificar contraste, colisiones y tokens.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción 2: gasto pasa a terracota, 40–45°; error se queda en 20°. Karen y Bodega Aurrera (hoy 50°) se mueven en la paleta categórica para conservar separación ≥ 20° o ΔL ≥ 0.15; el script regenera y verifica contraste y colisiones. Además, el barrido de pares cercanos incluye la marca (R3-04).
- **Estado:** RESUELTO. Se aplica en: U2a 1.1 (`colorimetria.py`, C2-colores.json), U2 DESIGN.md en WF-1.
