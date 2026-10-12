# Bloqueo U5-tipografia-tres-tamanos

Contexto: WF-3b, `revisor-u5-r1` (12-oct-2026). DESIGN.md §3.2 dice «Máximo tres tamaños por pantalla además de `caption`», pero los componentes del mismo DESIGN ya piden más en una pantalla con barra lateral y tarjetas: barra lateral y tablas 14 px, cuerpo y campos 16 px, título de tarjeta 20 px, cifra de tarjeta 24 px y título de página 32 px. Hoy 68 de 96 pantallas de la muestra usan 4 tamaños (14/16/20/24) y, al corregir el título de página a 32 px (hallazgo M-1), usarán 5.

Qué falta: saber qué cuenta en la regla.
- A: la regla cuenta solo el área de contenido (sin barra lateral, barra superior ni barra inferior) y excluye las cifras de tarjeta → se mide así en la ronda 2.
- B: la regla se reescribe como «máximo cinco tamaños por pantalla (toda la escala)».
- C: la regla vale literal y hay que rediseñar tarjetas y títulos.

Consecuencia: mientras no se decida, la ronda 1 lo registra como bajo (B-3) y no cambia el veredicto. Lo decide Axel.
