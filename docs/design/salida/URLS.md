# URL de la Pista 3 (Claude Design por Artifacts)

Versión 1.0 · 12-oct-2026 · criterio P3-11 (`docs/criterios/P3-claude-design.md` 1.3) · lo escribe el agente `integracion-wfp3` con los datos que dio el orquestador; las URL las publicó el orquestador con la herramienta Artifact (K45).

**Estas dos URL son privadas y son solo para Axel.** Nadie más las abre, nadie las pega fuera del proyecto, no se comparten, no se exportan a Canva ni a otra herramienta sin que Axel lo confirme (P3-09). Los agentes y el frontend no dependen de ellas: el proyecto sigue con `DESIGN.md` y `C2-tokens.json` aunque las URL dejen de abrir (decisiones 42, 43 y K45).

## Las dos URL

| Qué es | URL privada | Versión | Fecha |
|---|---|---|---|
| Sistema de diseño «Finanzas Hogar» (tipo Design System) | https://claude.ai/artifact/D6r3woEKMBvpWiK6gzr3Ra | 3 | 12-oct-2026 |
| Lienzo «Finanzas Hogar · Pantallas» (tipo Design) | https://claude.ai/artifact/LGrqytmCEAdyokHbRBaHNz | 3 | 12-oct-2026 |

La versión del sistema de diseño que el lienzo tiene instalada está en su `canvas.json` (llave `designSystems`): namespace `finanzas-hogar`, versión `1791789757-a409`, copiada el 12-oct-2026.

## Qué contiene cada una

**Sistema de diseño (versión 3).** La base visual que los artboards usan; se reimplementa en la app con los tokens de C2.
- `README` con los 8 principios, el tono y las reglas de texto en español de México.
- `tokens.json` y `tokens.css`: los mismos tokens que `docs/contratos/C2-tokens.json` (los hex coinciden uno por uno; se comprobó con `comm`), con tema claro y oscuro, matiz índigo 280°, y la familia `fontWeight`.
- 20 componentes con vista previa y README (`components/<Componente>/preview.html`): AvisoBreve, BarraInferior, BarraLateral, BarraSuperior, BarraVista, Boton, CalendarioPagos, CampoTexto, ChipPersonaTienda, EstadosVista, Galeria, Grafica, HojaDialogo, InterruptorOpciones, MenuBarra, PanelAjustes, TablaBD, Tablero, TarjetaMetrica y VistaPrevia, más la portada (`Cover`).
- Fuente Inter variable (`fonts/`).

**Lienzo (versión 3).** Todas las pantallas de U1 en móvil 412 y escritorio 1440, con cifras sintéticas.
- 90 artboards (45 pantallas × 2 anchos) en 9 páginas: Shell (entrar, barra lateral y menú «/», registrar gasto y peso), Inicio, Gastos, Presupuesto, Comidas, Compras y tiendas, Créditos, Salud y Reportes, y Ajustes (que se parte en «Ajustes 1 a 7» y «Ajustes 8 a 14»).
- Las 10 áreas de revisión son esas 9 páginas con Ajustes dividido en las partes a y b.
- Todos los datos salen de `docs/design/maquetas/datos-sinteticos.json` (nada del reporte real ni de Kueski; los nombres Axel, Ximena y Karen sí aparecen, K45.1 y Q2-98).
- Claro y oscuro se cambian con el menú de tema del lienzo.

## Dónde están las fuentes en disco

| Qué | Ruta |
|---|---|
| Sistema de diseño (fuente de lo publicado) | `/home/deepsec/finanzas-hogar/docs/design/sistema/project/` |
| Lienzo: `canvas.json` y 90 artboards `.dc.html` por página (fuente de lo publicado) | `/home/deepsec/finanzas-hogar/docs/design/maquetas/lienzo/project/` |
| Fragmentos JSON por área que armaron el `canvas.json` | `/home/deepsec/finanzas-hogar/docs/design/maquetas/lienzo/fragmentos/` |
| Datos sintéticos de todas las pantallas | `/home/deepsec/finanzas-hogar/docs/design/maquetas/datos-sinteticos.json` (y su `LEEME.md`) |
| Brief de referencia (22 elementos) | `/home/deepsec/finanzas-hogar/docs/design/entrada/` |
| Herramientas que generan y verifican (tokens, datos, lienzo) | `/home/deepsec/finanzas-hogar/docs/design/tools/ds/` |

Lo que está en `project/` es el estado vigente (el `canvas.json` y el `design-system.json` declaran `v` = 3). Si una copia por vuelta y `project/` difieren, manda `project/`.

## Copias por vuelta

Cada vuelta deja en `/home/deepsec/finanzas-hogar/docs/design/salida/iter-N/<área>/` los artboards que cambiaron en esa vuelta, con el fragmento `<área>.json` de la área. No es el lienzo completo en cada carpeta.

| Carpeta | Artboards copiados | Áreas |
|---|---|---|
| `iter-1/` | 90 | las 10: shell, inicio, gastos, presupuesto, comidas, compras, creditos, salud-reportes, ajustes-a, ajustes-b |
| `iter-2/` | 76 | ajustes-a, ajustes-b, creditos, inicio, presupuesto, salud-reportes |
| `iter-3/` | 26 | ajustes-b, inicio, presupuesto |
| `iter-4/` | 10 | inicio (8 de Inicio y los 2 de Flujo y proyección de Presupuesto) |

Notas de las copias (para que nadie las tome por completas):
- No hay `iter-N/project/` ni un `canvas.json` por vuelta; el `canvas.json` vigente está en `docs/design/maquetas/lienzo/project/canvas.json`.
- Las correcciones de Ajustes b de la vuelta 4 y de Inicio y Flujo de la vuelta 5 están en `lienzo/project/` pero no tienen copia en `iter-4/` ni `iter-5/`.
- Las copias de los informes de revisión (`P3-iter-N-<área>.md`) todavía no están en `iter-N/`; los originales, con las líneas de evidencia, están en `/home/deepsec/finanzas-hogar/docs/estado/wfp3/`. Queda como pendiente del orquestador (ver `docs/estado/wfp3/INTEGRACION.md`).

## Qué falta: exportación y aprobación

- **Exportación PDF/HTML desde la página (menú Share › Export del lienzo): pendiente de Axel.** Se hace desde la página del lienzo, que solo Axel abre; él la hace al aprobar y la deja en `/home/deepsec/finanzas-hogar/docs/design/salida/final/` (hoy solo tiene `.gitkeep`). Es una puerta 3 nueva (P3-09): aprobar la referencia visual y exportar.
- **Aprobación de Axel (puerta 3, P3-01): pendiente.** Axel abre solo la URL del lienzo, revisa la versión 3 y decide: aceptar, o pedir una vuelta más dentro del máximo de 3. Las vueltas anteriores las revisaron agentes Sonnet (K45.3).
- Hasta que Axel responda, la Pista 3 queda «pendiente de Axel» y no se declara cumplida.
