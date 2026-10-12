# revisor-creditos-r1 · Pista 3, fase 2 (lienzo) · 12-oct-2026

Agente: `revisor-creditos-r1` (revisor independiente, vuelta 1 de 3; no corrige). Área: `docs/design/maquetas/lienzo/project/creditos/` (8 artboards) y `docs/design/maquetas/lienzo/fragmentos/creditos.json`. Informe de revisión (P3-07): `docs/estado/wfp3/P3-iter-1-creditos.md`.

## Hecho

- Revisé los 8 artboards (P-17 Todas las líneas, P-18 Crédito, P-19 Compra a crédito y P-20 Comparador, en móvil y escritorio) y el fragmento contra P3 1.3 (P3-07, P3-08, P3-12), U5-04, U5-06 y U5-07, DESIGN.md, U1 (§0.4, §1.3, §2.1, §2.2, §2.4, §3.19 a §3.22, §6.4, F-4) y `format.md`.
- **Veredicto: no aprobado.** 0 altos, 9 medios y 10 bajos.
- Medios: M-01 partes por persona y filtro en la agenda móvil; M-02 serie de disponible en P-18 móvil; M-03 gráfica ilegible a 412 px; M-04 «Marcar pago» por parte (F-4); M-05 promociones como franjas en el mes; M-06 salida incompleta del comparador; M-07 intereses del año y compras de tarjeta y tienda, que dependen del JSON; M-08 barra lateral no colapsable; M-09 objetivo táctil de 20 px.

## Verificación

Scripts propios en `datos/tmp/rev-cred/` (`python3 -I` y un script de node para `renderVals()`), borrados al terminar:

- (a) Forma: 8/8 con las líneas obligatorias, JSON de `data-props` válido y `$preview` igual al fragmento. 0 etiquetas sin cerrar, `lang="es-MX"`, título propio, nombres válidos y stems únicos. Huecos `{{ }}` resueltos en todos los estados (0 fallas).
- (b) Tokens: 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`; todas las `var(--x)` existen en `tokens.json`.
- (c) Contraste: 26 pares del markup + 16 extra, en claro y oscuro, con los hex de `tokens.json`: 0 bajo 4.5:1 (texto) o 3:1 (interfaz).
- (d) P3-12: 0 cifras canon (21 patrones, con y sin coma), 0 credenciales; 10 valores muestreados idénticos a `datos-sinteticos.json`. Los montos derivados se comprobaron aritméticamente. Las 15 fechas que coinciden con `oraculo/model.json` llevan montos sintéticos.
- (e) Fidelidad: tabla pantalla → elementos en el informe; Estados de U1 §2.2 en los 4 artboards de escritorio.
- (f) Accesibilidad: `button`, `a` y `label for` reales; botones de solo icono con `aria-label`; 1 objetivo móvil < 44 px; sin emoji, sin barra de estado ni teclado falsos, sin inglés; dinero tabular.
- (g) Móvil con raíz fija de 412 px (4/4); escritorio fluido con `expand: fill` (4/4).
- (h) Enlaces: todos a archivos existentes del lienzo, del mismo ancho, o a `#`.
- Ruta del registro de herramientas: este agente no tiene journal propio en disco. Sus llamadas fueron Bash y lecturas, con 0 de navegador y 0 de Artifact (P3-04).

## Supuestos

1. «Pantalla principal» del área = P-17. Sus faltantes (intereses del año, franjas de promoción, partes en móvil) los califiqué como medios porque son secundarios: deuda viva, disponible A/B y calendario unificado sí están.
2. Lo que falta porque `datos-sinteticos.json` no lo trae (M-07, B-08) cuenta como hallazgo del área. La corrección pasa primero por el dueño de los datos.
3. Para la tasa de Coppel (B-06) interpreté `tasa_anual_pct` con `base_tasa: mensual` como tasa anual con base de cálculo mensual.
4. El contraste se calculó sobre colores declarados (clases del `helmet`, estilo en línea y herencia), sin navegador. Las alturas de `$preview` no se pueden medir sin navegador.
5. `is_interactive: true` en P-17, P-18 y P-19 no cuenta como defecto (lo pide format.md para artboards con manejadores); lo anoto en B-10 para que decida el orquestador.

## Bloqueos

Ninguno.

## Pendientes

- Vuelta 2: que el autor corrija M-01 a M-09 y, si cabe, los bajos.
- El dueño de `datos-sinteticos.json` debe agregar intereses del año (N-191), compras de la tarjeta (N-202), rendimiento por peso (N-226), puntaje del comparador y 4 pesos iguales (N-225).
- El orquestador debe medir alturas, copiar este informe a `docs/design/salida/iter-1/` (P3-06) y decidir sobre `is_interactive`.

## Consumo aproximado

- Leído: unos 190 KB (criterios P3 y U5 por secciones, U1 §0.4/§1.3/§2.2/§2.4/§3.19–3.22/§6.4/F-4, DESIGN §5.2/§7.6–7.10/§8.5/§9.3, format.md entero, craft.md en 3 secciones, reporte del autor, fragmento, partes de los 8 artboards y llaves de créditos del JSON).
- Escrito: unos 25 KB (informe y este reporte); scripts temporales de unos 15 KB, ya borrados.
- Acciones de herramienta: unas 50.
