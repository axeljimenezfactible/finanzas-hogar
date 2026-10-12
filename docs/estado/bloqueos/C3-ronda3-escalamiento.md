# Bloqueo C3-ronda3-escalamiento

Contexto: WF-3b, `revisor-c3-r3` (12-oct-2026). El revisor independiente tiene como máximo 3 rondas (plan §3.1, principio 5; trazabilidad de C3-07). La ronda 3 de `docs/contratos/C3-integracion.md` v1.2 sale **devuelta y sin firma**, con 1 hallazgo alto, 0 medios y 6 bajos (`docs/estado/wf3/REVISION-c3-r3.md`):
- R3-1 (alta, C3-08): F-3 «Comparar canasta» falla con la semilla por la **ruta A**. La tarjeta «Mejor canasta» lleva a «Productos y precios» porque la lista del plan no tiene precios de ingredientes. Con la lista vacía, la tarjeta anuncia la compra del 10-oct-2026, pero P-14 abre la semana 42. La ruta B ya pasa con la cifra de K46 (Q2-47 = A). Pasan 7 de 8 flujos.

Los demás criterios altos (C3-01 a C3-06 y C3-09) se cumplen con comandos reejecutados por el revisor.

Qué falta: la decisión de Axel sobre el **punto 2 de `docs/estado/bloqueos/C3-flujos.md`** y cómo se firma C3 sin una ronda 4.

Opciones (PROPUESTA, sin decidir):
| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. Axel elige la opción (c) del punto 2: U1 aclara que la ruta A pide una lista con todos los precios en al menos una tienda. Se autoriza una verificación acotada de C3-08 (solo F-3) y la firma | C3 firmable sin otra ronda completa; la evidencia de la ruta A ya existe (`pendientes.spec.ts:117`, 3 toques) | Bajo: una aclaración en U1 y una verificación de unos 5 min | Bajo. Con la semilla, la tarjeta sigue anunciando una canasta y abriendo otra (cara 2b), salvo que también se elija (b) |
| B. Axel elige la opción (b) del punto 2: la semilla trae los precios de los ingredientes (dueño: semilla o agente externo K53). Después, verificación acotada y firma | F-3 pasa con la semilla tal cual | Medio: datos nuevos en la semilla y repetir las corridas | Medio: inventar precios sería asumir; necesitan una fuente |
| C. Aceptar C3 con F-3 ruta A como pendiente con fecha, firmado con la excepción citada | Se llega a la puerta 4 | Ninguno ahora | C3-08 pide N de N flujos; la firma tendría una excepción explícita |

Recomendación del revisor: A, y además la cara (b) de 2b (la tarjeta abre P-14 con la compra que anuncia) como mejora posterior.

Estado: ABIERTO (lo decide Axel).
