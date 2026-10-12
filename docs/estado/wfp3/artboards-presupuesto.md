# artboards-presupuesto · Pista 3, fase 2 (lienzo)

Fecha: 12-oct-2026 · Agente: artboards-presupuesto (Sonnet) · Alcance: P-06 a P-09, móvil y escritorio (8 artboards).

## Hecho
- 8 artboards en `docs/design/maquetas/lienzo/project/presupuesto/`:
  - `presupuesto-p06-mes-{m,e}.dc.html` (Presupuesto del mes)
  - `presupuesto-p07-flujo-{m,e}.dc.html` (Flujo y proyección)
  - `presupuesto-p08-ingresos-{m,e}.dc.html` (Ingresos)
  - `presupuesto-p09-ahorro-{m,e}.dc.html` (Ahorro y metas)
- Fragmento `docs/design/maquetas/lienzo/fragmentos/presupuesto.json` (8 boards + order, sin x ni y; `expand: "fill"` solo en escritorio; `is_interactive: false`).
- Contenido por pantalla (nombres de U1 §3.8 a §3.11):
  - P-06: selector de mes; 6 tarjetas N-080 (Ingreso del mes, Plan de gasto, Ahorro planificado, Sin asignar, Real, Restante) y barra real contra plan; alertas N-084; quincenas N-082; despensa devengada N-083; 13 categorías con plan, real, diferencia, barra con marcas de 80 y 100 %, estado y número de gastos, ordenadas por % usado; venideros aparte; diálogo «Ajustar presupuesto» con alcance K2 («De este mes en adelante» preseleccionado, «Solo este mes») y vista previa del efecto (U4 P-26.02).
  - P-07: saldo inicial, cierre de octubre y saldo acumulado más bajo; 24 quincenas de efectivo (ingresos esperados marcados, base de cada quincena); gráfica devengado contra efectivo de 12 meses con patrón en la segunda serie, leyenda y «Ver como tabla»; deuda viva a una fecha con selector de fecha, partes del hogar y de Karen y desglose por línea; supuestos (N-093), interruptor «Repetir la despensa bimestral y el papel» (apagado) y lectura B de Kueski; acción primaria «Simular compra» enlazada a P-20.
  - P-08: totales (capturado, esperado pendiente, extras, total del mes); tabla de quincenas y extras por fecha con estado capturada/esperada/extra y quién/cuándo; calendario de ingresos de octubre (hoy = 21); formulario «Capturar quincena» prellenado con el esperado.
  - P-09: ahorrado total, fondo, aportado contra planificado del mes; galería (fondo, 2 metas con avance, fecha meta, aportación planificada y necesaria, estado Atrasada / En ruta, tarjeta «Nueva meta»); movimientos con chips de destino; formulario «Registrar aportación». En móvil, movimientos al abrir cada tarjeta.
  - Escritorio: fila «Estados» (Vacío, Carga, Error) con los textos de U1 §2.2.
- Enlaces entre artboards con `../<area>/<archivo>.dc.html` del mismo ancho (barra lateral, barra inferior, acciones primarias); tabs y destinos sin artboard van a `#`.

## Verificación
Script propio `python3 -I` (en `/tmp/p3pres/verify.py`, fuera del repositorio), resultado: 310 comprobaciones, 0 fallas. Cubre, por archivo: cabecera, `support.js`, `tokens.css`, `<x-dc><helmet><style>`, `data-props` con `$preview` (JSON válido, 412 u 1440 de ancho), clase `Component extends DCLogic`, etiquetas balanceadas con `html.parser` (incluye `sc-if`/`sc-for`), 0 hex y 0 rgb/hsl fuera de `data-props`, todos los `var(--x)` existen en `tokens.css` (356 variables), 0 cifras prohibidas de P3-12 (con y sin coma) y 0 valores canon en centavos, sin `innerHTML`, `iframe`, red ni emoji, un solo `h1`, ids únicos, botones de icono con `aria-label`, enlaces dentro de la lista de pantallas, y que cada hueco `{{...}}` resuelve contra lo que devuelve `renderVals()` (se ejecutó con node). Fragmento: válido, lista exactamente los 8 archivos, w y h iguales a `$preview`.
Cotejo de datos: los 234 pares `*_fmt` / `*_cents` de `categorias`, `presupuesto`, `ingresos`, `ahorro`, `proyeccion` y `deuda_viva` coinciden con el formato `$1,234.56` / `−$12.30` que genera el artboard; las sumas de categorías cuadran con el resumen (plan $27,586.00, real $20,525.78).

## Supuestos
1. Alturas (`$preview` y `h` del fragmento) ESTIMADAS a mano porque no se usó navegador: móvil 4400 / 3400 / 1850 / 2800 y escritorio 2650 / 2300 / 1450 / 1850 (P-06 a P-09). El orquestador debe medirlas al publicar y corregir `h` (el contenido define la altura real; el mínimo móvil de 915 se cumple).
2. Tablas hechas con `div` y roles ARIA (`role="table"`, `row`, `columnheader`, `cell`) y no con `<table>`: el navegador mueve las etiquetas `sc-for` fuera de un `<tbody>` (foster parenting) y el formato no documenta tablas con bucles. Misma razón para la gráfica: columnas en HTML y CSS, no SVG con bucles.
3. P-07: el «interruptor A/B» del encargo no existe en U1 para P-07; se dibujó el interruptor de N-318 (apagado) y, solo lectura, la lectura B activa de Kueski con enlace a Ajustes. La gráfica usa columnas por mes (no área apilada §8.3) porque con estos datos el efectivo acumulado supera siempre al devengado (la capa «devengado sin pagar» valdría $0.00); el subtítulo dice «pagado por adelantado» con la diferencia al cierre de dic-2026. Se grafican los totales (hogar + Karen) para que cuadren con la diferencia acumulada de N-091; la tabla de apoyo también.
4. P-07, deuda viva: el selector abre en 31-dic-2026 (existe en `deuda_viva.a_fecha`); el valor por defecto de U1 (fin del horizonte, 30-sep-2027) no tiene entrada en `a_fecha` y valdría $0.00. Se compara contra «hoy».
5. P-06: el diálogo «Ajustar presupuesto» se dibuja abierto sobre Mascotas con el monto nuevo igual al real ($1,098.60); es un cálculo con datos del JSON (plan total antes $27,586.00, después $27,684.60), no una cifra inventada.
6. P-08 abre en la subsección Detalle (default de U1), P-06, P-07 y P-09 en Resumen. Pestañas Tendencias y Calendario: solo la fila de pestañas (enlace `#`).
7. Enlaces a Ajustes › secciones apuntan a `ajustes-p26-indice-{m,e}` porque el slug de cada sección lo define otro autor. «Nueva meta» (P-26.13) y «Ir a Ajustes › Ingresos» (P-26.03) usan el índice.
8. Colores de ahorro: aportaciones y retiros llevan solo signo y texto (+$ / −$), sin color de rol; ahorro «en ruta» usa insignia de éxito y «atrasada» de advertencia con icono y texto.

## Bloqueos
Ninguno.

## Pendientes
- Medir alturas reales y corregir `h` en el fragmento y en `$preview`.
- Revisar en el lienzo (sin navegador no se vio): que los huecos dentro de `style="{{...}}"` y `class="... {{...}}"` se interpreten como describe format.md, y las barras de la gráfica (`height:%`) y el calendario de P-08.
- Subsecciones Detalle, Tendencias y Calendario de estas pantallas no tienen artboard propio (no están en la lista del lienzo).
- Las fechas meta del calendario de ahorro (N-115) y la tendencia de ahorro (N-114) quedan para esas pestañas.

## Consumo aproximado
- Leído: ~230 KB (fuentes y datos por secciones).
- Escrito a mano: ~104 KB (plantillas 80 KB + CSS, barra lateral, ayudantes y script de inyección 24 KB, en `/tmp/p3pres/`).
- Entregables en disco: 237 KB (los artboards llevan el CSS, la barra lateral, los iconos y los datos incrustados; superan los 150 KB porque cada archivo es autocontenido).
- Acciones de herramienta: ~62.
