# revisor-ajustes-b-r4 · vuelta 4 · 12-oct-2026

Agente `revisor-ajustes-b-r4`: revisor independiente de ajustes-b (P-26.08 a P-26.14, 14 artboards y el fragmento `ajustes-b.json`). No corrige.
Informe (criterio P3-07): `docs/estado/wfp3/P3-iter-4-ajustes-b.md`.

## Hecho
- Comprobé que los dos hallazgos de la vuelta 3 están cerrados:
  - **R3-A1 (alta).** Los 14 artboards tienen la regla `::placeholder{color:var(--ui-text-secondary);opacity:1}`. Los 49 marcadores de posición quedan en 5.40 en claro y 10.43 en oscuro, el par C-035.
  - **H-11 (media).** El dato regenerado da `umbral_reposicion` 12 pza, 2,000 ml y 1,000 g. s14-m y s14-e pintan esos niveles, y «Reponer» ahora cumple la regla K7 de la misma pantalla.
- Diff por tokens contra `docs/design/salida/iter-3/ajustes-b/`: solo cambió eso y no hay regresiones. El fragmento es idéntico.
- Pasé de nuevo la batería completa, de (a) a (h), con scripts propios: forma, tokens, contraste por cascada en los dos temas, P3-12, fidelidad a U1/U4, accesibilidad y objetivos táctiles, diseño móvil y escritorio, enlaces y fragmento.
- **Veredicto: APROBADO.** 0 altas, 0 medias y 21 bajas: las 18 abiertas de antes y 3 nuevas. Las nuevas son B-22 (fechas con cero a la izquierda), B-23 (días restantes enteros en ajustes y con decimal en compras-p16) y B-24 (falta enlace del prototipo a P-25).

## Verificación
- Scripts en `datos/tmp/rev-ajb-r4/` (`python3 -I`), borrados al terminar.
- Resultados:
  - (a) 14/14 de forma y `$preview`; 0 etiquetas sin cerrar y 0 cierres cruzados; stems únicos en los 90 artboards.
  - (b) 0 hex y todos los `var()` existen.
  - (c) 1,583 nodos por tema y 0 pares bajo 4.5:1; UI con 0 bajo 3:1.
  - (d) 0 cifras canon (21 patrones y 19 en centavos) y 0 credenciales; la muestra de 10 valores coincide con el JSON.
  - (e) 0 elementos de U4 faltantes en los escritorios; estados en los 6 escritorios que los llevan.
  - (f) 144 objetivos móviles con 0 bajo 44 px; 65 de 65 montos tabulares; 0 emoji y 0 fechas ISO.
  - (g) 7/7 móviles con raíz fija y 7/7 escritorios fluidos con `expand:fill`.
  - (h) 212 enlaces válidos, 30 a `#` y 0 rotos.

## Supuestos
- `support.js` y `ds/finanzas-hogar/tokens.css` no existen en el disco local del lienzo. Los pone el tipo «Design» al publicar, por eso la forma se juzga por las líneas y los tokens se resuelven con `docs/design/sistema/project/tokens.json` y `tokens.css`, que tienen los mismos 356 nombres.
- Las alturas del fragmento no se volvieron a medir (este revisor no usa navegador). Los cambios de la ronda no alteran el flujo: es una regla CSS y tres textos de largo parecido. Además, `correccion-ajustes-b-r3` midió s14-m en 2102 px.
- B-23 (redondeo de los días restantes) y B-16 (avance entero) quedan como decisiones del orquestador, porque U1 y U4 no fijan el redondeo de los días y el JSON trae el avance entero.

## Bloqueos
- Ninguno. No registré nada en `docs/estado/bloqueos/`.

## Pendientes
- Las 21 bajas no bloquean. Hay que copiar este informe a `docs/design/salida/iter-4/` (P3-06), lo hace el orquestador.

## Consumo aproximado
~150 KB leídos (U4 parte b por secciones, informe r3, `format.md` parcial, salidas de scripts), ~22 KB escritos (dos reportes y scripts temporales), ~52 acciones.
