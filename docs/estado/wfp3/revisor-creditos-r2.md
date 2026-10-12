# revisor-creditos-r2 · vuelta 2 de 3 · Créditos (P-17 a P-20)

Agente: `revisor-creditos-r2` (revisor independiente; no corrige) · Fecha: 12-oct-2026 · Informe P3-07: `docs/estado/wfp3/P3-iter-2-creditos.md`.

## Hecho

- Comprobé los 9 hallazgos medios de la vuelta 1. M-01, M-02, M-03, M-04, M-05, M-06, M-08 y M-09 están cerrados en los archivos señalados. M-07 está parcial y su residuo de datos pasa a bajo (B-03), porque lo bloquea `datos-sinteticos.json`.
- Encontré 2 hallazgos nuevos que impiden aprobar:
  - **A-01 (alto)**: el relleno `--ui-chart-1` sobre la pista `--ui-viz-track` da 2.22:1 en oscuro en P-17 m, P-17 e, P-18 m y P-18 e. Se corrige cambiando el relleno a `--ui-viz-meter` (C-102).
  - **M-01 (medio)**: los pagos de la agenda de P-17 m no abren la vista previa ni «Marcar pago» (DESIGN §7.8).
- Registré 11 bajos (B-01 a B-11). Varios vienen de la vuelta 1: iconos de 14 y 16 px, ejes, barra inferior, input de fecha y enlaces a `#`.
- **Veredicto: no aprobado (1 alto, 1 medio).**

## Verificación

- Usé scripts propios en `datos/tmp/rev-cred2/` (`python3 -I` más un arnés en node que ejecuta `renderVals()` y explora estados llamando a los manejadores). Los borré al terminar.
- Resultado de lo que se puede medir en el markup:
  - Forma: 0 fallas en 8 artboards.
  - `html.parser`: 0 etiquetas sin cerrar.
  - Plantilla: 0 huecos sin resolver.
  - Tokens: 0 hex y 0 `var()` inexistentes.
  - P3-12: 0 cifras prohibidas y 0 credenciales.
  - Objetivos táctiles: 0 por debajo de 44 px en móvil.
  - Enlaces: 0 rotos.
  - Accesibilidad: 0 controles sin nombre o sin `label`.
  - Contraste: 0 pares de texto bajo umbral y 1 par de interfaz bajo umbral (A-01).
- El detalle está en el informe, en «Verificación» y en «Contraste recalculado».

## Supuestos

- «Hoy» es el 21-oct-2026 (`meta.hoy` del JSON).
- El contraste se calculó con los hex de `docs/design/sistema/project/tokens.css`, que se genera de `tokens.json`.
- El fondo efectivo es el del ancestro más cercano con fondo. La pista de las barras y la retícula de las gráficas son decorativas (C-105, §8); el relleno de la barra no lo es.
- A-01 se califica alto por la regla de severidad de la orden («contraste bajo umbral») y por U5-04, que solo exime lo que DESIGN declara decorativo. Atenuante: la barra es `aria-hidden` y el monto va en texto al lado.
- M-01 se trata como hallazgo nuevo. La lista de archivos de M-04 en la vuelta 1 no incluía P-17 m, aunque su texto decía «en P-17 y P-18».
- Los 28 montos que no están literalmente en el JSON son derivados por suma o resta de cifras del JSON. Verifiqué 4 a mano: saldo Kueski de oct $4,327.47, nov $906.98, dic $2,797.02 y total de MSI $14,099.04. No los cuento como cifras inventadas.
- `is_interactive: true` en P-17, P-18 y P-19, que tienen manejadores, sigue a `format.md`. P-20 queda en `false` porque es estático. No es hallazgo.

## Bloqueos

- B-03 depende del dueño de `datos-sinteticos.json`. Faltan `creditos.totales.intereses_anio_*`, `creditos.tarjeta.compras[]` completas, el rendimiento por peso de cada escenario y los pesos de los 4 objetivos. Mientras no estén, el área de créditos no puede cerrarlo sin inventar cifras.
- Sin navegador (K45) no se pueden medir las alturas reales (B-11).

## Pendientes

- **Autor de créditos, vuelta 3:**
  - A-01: un token en 4 archivos.
  - M-01: botón en las filas de la agenda de P-17 m que abra la hoja de P-18 m.
  - Si cabe, los bajos baratos: B-04 (4 `href`), B-05 (clase `tn`), B-06 (iconos a 20 px), B-09 (fecha como texto) y B-10 (filtros de compras en P-18 m).
- **Orquestador:**
  - Medir `$preview`/`h` reales.
  - Decidir la copia de este informe en `docs/design/salida/iter-2/` (P3-06). Yo no la hice porque mi alcance de escritura son los dos reportes.

## Consumo aproximado

- Lectura: unos 190 KB. Incluye `format.md` y `craft.md` por secciones, el informe y la corrección de la vuelta 1, U1 y DESIGN por secciones, `tokens.css`, el JSON por consultas y el texto renderizado de unos 25 estados.
- Escritura: unos 45 KB (scripts temporales, ya borrados, y los dos reportes).
- Acciones de herramienta: unas 57.
