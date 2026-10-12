# U2b-2-design · reporte final (WF-1)

## Hecho
- `docs/design/DESIGN.md` §7 Componentes: 7.0 convenciones de estados; 7.1 tabla de 41 tokens semánticos `ui.*` (alias de primitivos de C2 en claro y oscuro, con hex); 7.2 botones y enlaces; 7.3 campos y edición en el lugar; 7.4 chips de persona y tienda con refuerzo de forma (K41: círculo/píldora contra cuadrado/rectángulo), avatares, chips de filtro e insignias; 7.5 vistas de base de datos (barra de la vista, tabla, tablero, calendario, galería); 7.6 tarjetas de métrica y barra de avance (umbrales `alerts.*` configurables, K44); 7.7 gráficas; 7.8 calendario de pagos (estados con icono y texto); 7.9 barra lateral (G-01); 7.10 barra superior, pestañas y barra inferior (G-02, G-05); 7.11 menú «/» y búsqueda (G-03, G-04); 7.12 vista previa, hojas y diálogos (G-06, diálogo de acción crítica); 7.13 panel de Ajustes (parámetros con «Restablecer», avisos PROVISIONAL citando `docs/estado/bloqueos/S1-parametros-avisos.md` con las dos variantes de operación); 7.14 estados vacío/carga/error y globales (G-07); 7.15 avisos breves (G-07); 7.16 matriz de cobertura. Cada componente con anatomía, variantes, estados (reposo, hover, foco, activo, deshabilitado, error), tokens, tamaños, móvil/escritorio y accesibilidad.
- §8 Patrones de visualización: marcas, elección de forma, área apilada devengado contra efectivo (capas efectivo + devengado sin pagar; ejemplo K19), dona ≤ 5 partes, dinero siempre en texto, 8 series en orden fijo con marcador, trazo y patrón (K27), personas y tiendas nunca en la misma gráfica (K41), lo que no se hace.
- §9 Accesibilidad y puntos de quiebre: tabla WCAG 2.2 AA criterio → cómo se cumple; foco visible; objetivo táctil ≥ 44 px; tokens `breakpoint.*` 360/412/768/1024/1440 y reajuste a 320 px; pruebas; preferencias del sistema.
- §10 Tokens: archivos y flujo; estructura DTCG de `C2-tokens.json` (grupos `color`, `ui`, `font`, `space`, `radius`, `border`, `focus`, `shadow`, `overlay`, `size`, `layout`, `icon`, `motion`, `breakpoint`, `viz`); `tokens.css` (primitivos estáticos `--color-*`, temáticas `--ui-*`, selectores de tema e interruptor); tema Tailwind 4 (`@theme inline`, solo semánticos como utilidades, variante `coarse`); **lista explícita de 158 pares × 2 temas (C-001…C-158)** para U2d; medidas de componente propuestas.
- Supuestos 4–9 agregados al final de la sección Supuestos (U2b-1 lo previó: «U2b-2 agrega los suyos»), sin tocar sus líneas.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/design/DESIGN.md` (secciones 7–10 y supuestos 4–9; 1–6 intactas, comprobado byte a byte).
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/U2b-2-design.md` (este reporte).
- Temporales en `datos/tmp/u2b2/` (generador de tablas y verificador) borrados.

## Verificación
- Generador propio (temporal) de §7.1 y §10.5 con control previo de contraste WCAG: «semánticos 41 · pares 158 · filas tema 316 · nuevos frente a contraste.csv 50 · mínimos texto 4.51, interfaz 3.06, texto grande 10.76 · bajo umbral: 0». La primera versión dio 3 pares bajo 4.5 (cifras de rol sobre `ui.surface.hover`: gasto oscuro 4.07, ingreso 4.30/4.37); se corrigió con `ui.surface.hover-subtle` y hover de solo borde en tarjetas (supuesto 6).
- `python3 -I docs/design/fuentes/scripts/verificar.py .` → código 0: «hex citados 162 · fuera de C2 [] · rutas de token citadas 958 · inexistentes [] · pares token→hex 158 · no coinciden []».
- Verificador propio de §7–10 → RESULTADO OK: hex en 7–10: 35, fuera de C2: 0; rutas `color.*`: 101, inexistentes 0; `ui.*` citados todos definidos en §7.1 y los 41 usados fuera de la tabla; secciones 1–6 y supuestos de U2b-1 intactos.
- **U2-01**: 11 encabezados `## `: los 10 del índice en orden + Supuestos; 7–10 con tablas y párrafos.
- **U2-12**: 15 de 15 grupos (los 12 del criterio más chips, barra superior/inferior y diálogos); los 15 subapartados de componente tienen los 6 estados (reposo, hover, foco, activo, deshabilitado, error) con «no aplica» y motivo donde corresponde; tokens semánticos con su valor claro y oscuro en §7.1. Cruce U1 §1.4: G-01…G-07 presentes en §7 (0 ausentes; matriz §7.16).
- **U2-13**: §8 contiene barras horizontales, líneas, área apilada, dona ≤ 5, «siempre en texto», K27 y orden fijo; tema por defecto según el sistema con interruptor (§9.6); `grep -i "linear|stripe|material"` en DESIGN.md = 0; Notion sigue como única referencia.
- U2-08/K41: «solo por color» aparece 3 veces; U2-11: «44 px» 20 veces en 7–10, ningún mínimo táctil menor; U2-10: 0 URL de terceros; 42 iconos citados en 7–10, todos existen en `sprite.svg`; 158 filas de pares, ningún id citado fuera de rango.

## Supuestos
1. Área apilada: capas efectivo acumulado y devengado − efectivo (cima = devengado); adelanto con capa en cero y texto. No cambia cifras ni contrato.
2. Pago vencido en advertencia, no en error (error solo validación y fallas, §2.4).
3. Capa `ui.*` (41 alias sin valores nuevos) y `ui.surface.hover-subtle`; prefijo `ui` para no chocar con la comprobación `color.*` de `verificar.py`.
4. Medidas de §10.6, nombre del archivo del tema Tailwind (`docs/contratos/tailwind-theme.css`) y nombres `breakpoint.*`: propuestas para U2c.
5. Patrón siempre presente desde la serie 2 (K27) y `tabular-nums` en todo dinero (§3.3) por encima de la guía de visualización.
6. Estado «Parcial» del calendario de pagos se deriva de las partes pagadas (K15).

## Bloqueos
- Ninguno nuevo. Se cita el abierto `docs/estado/bloqueos/S1-parametros-avisos.md` (§7.13 PROVISIONAL con las dos variantes: `updateNotification` o la operación de settings).

## Pendientes
- U2c: escribir `docs/contratos/C2-tokens.json`, `docs/contratos/tokens.css` y el tema Tailwind 4 según §10.2–10.4, incluidos los 41 `ui.*` de §7.1 y las medidas de §10.6.
- U2d: calcular la unión de §10.5 (316 combinaciones) con `docs/design/contraste.csv` (202) desde `C2-tokens.json`; se espera 0 bajo umbral.

## Solicitudes fuera de alcance
1. `docs/plan/protocolo-agentes.md` §4: agregar `docs/design/DESIGN.md` (155 KB) a la lista de «nunca leer entero»; leer por secciones (`grep -n '^##'`). Motivo: tope de lectura.
2. `docs/design/U1-arquitectura.md` §5.4 y U4 (Ajustes › Apariencia): interruptor local «Atajos de una tecla» para WCAG 2.1.4, si Axel lo aprueba en la puerta 2.
3. U3 / C1 (`docs/contratos/extractos/` de créditos): confirmar que la lectura del calendario de pagos trae el estado pagado por parte (persona) para mostrar «Parcial»; si no, abrir bloqueo en C1.
4. U2c y su revisor: comprobar que todo `ui.*` citado en DESIGN.md existe en `C2-tokens.json` (C2-04); `verificar.py` solo comprueba `color.*`.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Interruptor «Atajos de una tecla» (encendido por omisión, guardado en el dispositivo) para cumplir WCAG 2.1.4 con `/`, `?` y flechas.
- Indicador «Guardado» de 2 s junto al control en Ajustes: suficiente para leerlo sin estorbar; el aviso breve sigue en 5 s (U1 G-07).
- Constantes de interfaz (5 s, 2 s, 400 ms, 10 s) como diseño y no como parámetros operativos de K44; si las quiere configurables, irían en Apariencia.
- Pago vencido mostrado como advertencia con «Vencido» y `calendar-x`, no como error.

## Consumo aproximado
- 30 acciones de herramienta.
- Leído ≈ 165 KB (protocolo, decisiones §G y §K, criterios U2 y partes de C2, DESIGN §1–6, plan §4.1–4.4, U1 §1.2–1.4, §2.4–2.6, §5 y §6, nombres de C2-colores.json, guía de visualización y salidas de scripts).
- Escrito ≈ 135 KB a mano (DESIGN §7–10 y supuestos ≈ 107 KB, de ellos 40 KB generados por script; scripts temporales ≈ 19 KB; reporte y progreso ≈ 7 KB). El generador reescribió además sus salidas temporales (≈ 40 KB) en tres corridas, ya borradas.
