# artboards-creditos · Pista 3, fase 2 (lienzo) · 12-oct-2026

Agente: artboards-creditos (Sonnet). Área: `docs/design/maquetas/lienzo/project/creditos/`. Fragmento: `docs/design/maquetas/lienzo/fragmentos/creditos.json`.

## Hecho

8 artboards (móvil 412 px y escritorio PAGE fluida con `expand: fill`), todos con la forma obligatoria (`.dc.html`, `support.js`, `tokens.css`, `x-dc`, `helmet`, script `DCLogic` con `$preview`).

| Pantalla | Móvil | Escritorio | Qué muestra |
|---|---|---|---|
| P-17 Todas las líneas | `creditos-p17-lineas-m` | `creditos-p17-lineas-e` | Pestañas Resumen · Detalle · Tendencias · Calendario (estado real). Deuda viva total ($44,606.33, hogar/Karen), pagos de la quincena y del mes, intereses del mes, aviso A/B «Reposición de capital pendiente de confirmar con Kueski», galería de 4 líneas (disponible A/B en Kueski), calendario unificado de hoy a 60 días (hogar, Karen, estado por parte, fecha de corte), promociones como franjas, tabla de detalle, deuda viva por mes y línea (12 meses), calendario mensual de octubre. Escritorio: tarjeta «Estados» (vacío, carga, error). |
| P-18 Crédito | `creditos-p18-credito-m` | `creditos-p18-credito-e` | Selector de línea (Kueski Pay revolvente, Banorte Oro tarjeta, Préstamo personal, Crédito Coppel de tienda) y 4 pestañas. Revolvente: disponible lectura B activa y A en segundo plano, compras de la línea, tasas por plazo, serie de disponible, gráfica A/B. Tarjeta: ciclo (corte, límite de pago, saldo al corte, pago sin intereses, pago mínimo, saldo actual) y meses sin intereses. Préstamo y tienda: amortización y abonos con estado. Acciones: Marcar pago, Registrar compra, Simular compra, Editar condiciones. Escritorio: Estados con «Falta capturar · Capturar» (K18). |
| P-19 Compra a crédito | `creditos-p19-compra-m` | `creditos-p19-compra-e` | Alterna Detalle de compra / Nueva compra. Formulario (línea, comercio, fecha, capital, intereses cotizados con IVA, plazo, promoción, capital por persona, nota), vista previa del cálculo (total, partes Hogar y Karen con pago por quincena, promoción no aplicable), pagos con estado por parte, gasto vinculado. Escritorio: Estados (calculando, error, dato inválido). |
| P-20 Comparador | `creditos-p20-comparador-m` | `creditos-p20-comparador-e` | Datos de la compra, restricciones duras (solo esta comparación), tope por persona, pesos de los 4 objetivos, recomendación con puntaje, 4 óptimos lado a lado (apilados en móvil), tabla de los 6 escenarios con descartados y su motivo (plegada en móvil). Escritorio: Estados (ningún escenario, faltan restricciones, evaluando, no disponible). |

Todos los datos salen de `datos-sinteticos.json` (creditos, promociones, comparador, deuda_viva, meta). Enlaces entre artboards con rutas relativas al mismo ancho (barra lateral, botones principales, «Usar escenario recomendado» abre P-19). Fragmento JSON con las 8 rutas y `order`.

## Verificación

Script propio `/tmp/cr/verify.py` (python3 -I + node): `RESULTADO: OK · comprobaciones 433`, 0 fallas.

- Líneas obligatorias en los 8 archivos; JSON de `data-props` parsea; `$preview` 412 × alto ≥ 915 (móvil) y 1440 (escritorio).
- `html.parser`: 0 etiquetas sin cerrar ni mal anidadas; 1 `h1` por archivo; ids únicos; todo `label for` apunta a un control.
- 0 hex literales y 0 rgb/hsl fuera de `data-props`; 0 cifras prohibidas de P3-12 (con y sin coma) y 0 valores en centavos canon.
- Todos los `var(--x)` existen en `tokens.css`/`tokens.json`; todas las clases usadas están definidas en el `helmet` o en `tokens.css`.
- Sin `innerHTML`, `iframe`, red, `lorem`, `<img>` ni emoji; negativos sin guion corto; `tabular-nums` presente; iconos con trazo 1.5.
- Cada hueco `{{ … }}` es una ruta con puntos y resuelve contra `renderVals()` ejecutado con node en 4 a 7 estados por pantalla (todas las pestañas y las 4 líneas de P-18, modo alta de P-19).
- Fragmento: `boards` y `order` listan exactamente mis 8 archivos; `expand: fill` solo en escritorio; stems únicos.

## Supuestos

1. «Hoy» es `meta.hoy` = miércoles 21-oct-2026 (el JSON manda, no la fecha de la orden).
2. Las alturas de `$preview` y de `h` del fragmento son estimadas sin navegador (no se podía medir); el orquestador debe medirlas. Móvil: 2200, 1800, 2250, 2400; escritorio: 2600, 1300, 1700, 1850.
3. Puntaje ponderado del comparador: el JSON no trae puntaje; lo derivé con normalización mínimo-máximo sobre los escenarios válidos y los pesos 40/30/30, con 0 % para «Deuda en cero» (100 menos la suma). Da E1 = 76, E6 = 70, E2 = 50, E5 = 43 y coincide con la recomendación `e1`. Es ilustrativo; manda el motor S3.
4. P-18 es UN artboard por ancho con selector de línea (estado interactivo) para mostrar los 4 tipos; los enlaces a «cada línea» (barra lateral, tarjetas de P-17) apuntan a ese artboard, que abre en Kueski Pay.
5. Intereses por pagar de la línea = deuda viva − capital pendiente (suma de pagos pendientes menos capital repuesto). No se muestran para la tarjeta porque el JSON no los trae de forma coherente; ahí van pago sin intereses y pago mínimo.
6. «Intereses pagados en el año» (N-191) no existe en el JSON: se muestra solo el mes.
7. Enlaces a Ajustes › Créditos y condiciones van a «#» porque el slug de `ajustes-p26-s05-…` no está en mi lista.
8. Gasto vinculado de P-19: el JSON no vincula la compra a un gasto; muestra el nombre y la fecha de la compra y enlaza a P-04.
9. Alta de P-19 reutiliza los mismos valores de ejemplo que el detalle; en móvil el bloque de pagos va abierto (U1 lo pliega) para que se vean los datos.
10. Iconos Lucide como símbolos `<symbol>`/`<use>` dentro del propio archivo (inline); `circle-dashed` es una aproximación con trazo discontinuo. La barra lateral y la inferior de mis artboards son versiones reducidas (sin iconos de sección); las canónicas son las del área shell.
11. `./support.js` y `../ds/finanzas-hogar/tokens.css` los aporta el orquestador (no existen en mi carpeta).
12. Las tablas usan `div` con roles ARIA (`table`, `row`, `columnheader`, `cell`) en vez de `<table>` porque `sc-for` dentro de `<tbody>` se reubica al parsear.

## Bloqueos

Ninguno duro. Cifras de tamaño: ver Consumo.

## Pendientes

- Medir las alturas reales de los 8 artboards y corregir `$preview`/fragmento.
- P-18 móvil › Detalle no incluye la tabla completa de la serie de disponible (sí mínimos A/B, tasas por plazo y gráfica); en escritorio muestra 10 de 21 eventos con «Mostrar 11 más».
- P-17 escritorio › Calendario: solo octubre de 2026 (sin navegación de mes).
- Filtros de P-17 y P-18 (línea, familia, estado, vivas/liquidadas) dibujados sin comportamiento.
- Enlaces a Ajustes › Créditos pendientes de slug; una vista de «pago vencido» (`estados_alternos.pago_vencido`) no está dibujada en P-17.

## Consumo aproximado

- Leído: ~190 KB (U1 por secciones, DESIGN §2.4/§3.3/§7/§9.4/§10.5, format.md, craft.md, LEEME, `datos-sinteticos.json` solo llaves de créditos, promociones, comparador, deuda_viva, tokens.css).
- Escrito: los 8 artboards suman ~218 KB en disco (generados por scripts), por encima del tope nominal de 150 KB; los scripts generadores (~108 KB, en `/tmp/cr`, fuera del repositorio) son lo que se tecleó. Más fragmento (~3 KB) y este reporte (~7 KB). Los artboards pesan más porque cada archivo lleva su `helmet` (~3 a 5 KB), su sprite de iconos y los datos de sus 4 líneas o estados.
- Acciones de herramienta: ~34.
