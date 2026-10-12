# U5 · Revisión independiente contra DESIGN.md, U3 y U4

Registro de rondas que pide U5-05 (máximo 3; si la 3 no aprueba, se escala a Axel). Cada ronda: pantallas y criterios evaluados por ID, fallas (pantalla, combinación, regla), correcciones aplicadas desde la ronda anterior y veredicto. Detalle de cada hallazgo y corrección sugerida: `docs/estado/wf3/REVISION-u5-r<n>.md`.

Combinaciones: Ac = Android 412×915 claro · Ao = Android oscuro · Ec = escritorio 1440×900 claro · Eo = escritorio oscuro.

## Ronda 1

Revisor: `revisor-u5-r1` (no capturó ni implementó nada) · 12-oct-2026 · HEAD `6d691d0` con el árbol de trabajo de WF-3b; binario `backend/bin/finanzas` y `frontend/dist` (09:02) en marcha en `finanzas.local:8443`, base `datos/finanzas-dev.db` sembrada. Capturas evaluadas: `docs/design/capturas-u5/` (164 PNG de `u5-capturas`).

### Comandos reejecutados por el revisor (U5-05)
- **U5-01**: script propio en Python (lee los `P-xx` y `P-26.xx` de U1, abre cada PNG con PIL): T = 41, 4T = 164 → **162 de 164 correctos**; mal: `P-26.13-android-claro.png` y `P-26.13-android-oscuro.png` (415×1732, no 412). Mi propia captura de P-26.13 en Android también sale de 415 px (la página mide 415 de ancho). `ls capturas-u5/ | wc -l` literal = 168 (164 PNG + `axe/`, `contraste/`, `evidencia/`, `hosts.json`; bloqueo `U5-conteo-carpeta`).
- **U5-03**: comando literal sobre los JSON del capturista = **20**. Mi corrida propia de axe-core 4.13 (WCAG 2.0/2.1/2.2 A y AA) en las 164 cargas da lo mismo: 5 serias/críticas por combinación (Ac 5 · Ao 5 · Ec 5 · Eo 5), 0 moderadas/menores; incompletos 15/15/43/43.
- **U5-04**: implementación propia en el DOM (colores computados → sRGB 8 bits, fondo compuesto por ancestros, opacidad acumulada, texto grande ≥ 24 px o ≥ 18.66 px en negrita). Pares únicos por pantalla y bajo umbral: Ac 2,558 · 113; Ao 2,560 · 113; Ec 3,459 · 131; Eo 3,460 · 131. **Texto 0 e íconos 0 bajo umbral**; los bajo 3:1 son todos de componentes de UI (botón 94/94/112/112 y selección 19 por combinación). 20 pares elegidos al azar (semilla 12102026) de los JSON del capturista: **20 de 20 coinciden, Δ máx. 0.000**.

### Pantallas evaluadas
- Lista visual U5-06 (Read de los PNG, 4 combinaciones): P-01 (Resumen del dashboard, U3), P-26 y P-26.01 a P-26.14 (Ajustes, U4), y una por entrada de la barra: P-02 Gastos, P-06 Presupuesto, P-10 Comidas, P-13 Compras, P-17 Créditos, P-21 Salud, P-23 Reportes, P-26 Ajustes, P-03 (el «+» de la barra inferior). Además Inicio ▸ Tendencias y Plan ▸ Tendencias (gráficas) en Ao y Ec.
- Mediciones en el DOM (fuente, escala, `tnum`, iconos, radios, espaciado, sombras, toque, foco, saturación, desborde, hosts, fondo): las 41 pantallas × 4.

### Criterios evaluados
| Criterio | Resultado | Evidencia |
|---|---|---|
| U5-01 | **Falla** | 162 de 164; P-26.13 Ac/Ao a 415 px |
| U5-02 | Pasa | 8 cifras visibles = `model.json` ($10,000.00, $6,905.50, $973.75, $3,094.50, $6,293.44, $598.82, $374.93, $1,191.50) |
| U5-03 | **Falla** | 20 (5 reglas × 4) |
| U5-04 | **Falla** | texto/íconos 0; UI 113–131 por combinación bajo 3:1; bloqueo `U5-contraste-ui-controles` abierto |
| U5-06 | **Falla** | tabla de abajo |
| U5-07 | **Falla** | script: 1,404 importes en formato, 0 fechas ISO, 0 palabras de la lista; pero en las capturas se lee inglés en P-26.11 («Choose File», «No file chosen») y P-26.04 («groceries», «eating_out») |
| U5-08 | Pasa | 1 host: `finanzas.local:8443` |
| U5-09 | Pasa | 3 toques («+», Despensa, Guardar); aparece en Inicio y en la tabla de P-02 |
| U5-10 | Pasa | fondo = C2 en 164 de 164; interruptor: 7 estados correctos |
| U5-13 | **Falla** | paso 2 (H-5) reproducido; pasos 1 y 3–7 pasan |
| C3-09 | Pasa | `c3-09.spec` reejecutado: 111 cifras, 0 diferencias |

### Lista de verificación U5-06 (14 elementos)
| # | Elemento (DESIGN.md) | Resultado | Fallas: pantalla · combinación · regla |
|---|---|---|---|
| 1 | Tokens de color y rol (índigo 280°, error rojo, crítica opuesta) | Pasa | Acción crítica no visible en la muestra |
| 2 | Tipografía autoalojada y escala | **Falla** | Inter 100 % y tamaños en la escala; pero el título de página mide 24 px/600 y no 32 px/700: todas · Ac Ao Ec Eo · §3.2 y §7.10 (`font.size.title-1`) |
| 3 | Cifras tabulares en todo dinero | Pasa | 0 importes sin `tabular-nums` |
| 4 | Iconos Lucide, trazo 1.5, 20/24 px | **Falla** | iconos de 16 px con trazo 1.2: P-17, P-18, P-19, P-21, P-23, P-26, P-26.01, .02, .04–.09, .11 · 4 comb. · §5.2; 24 px con trazo 1.8 en P-26.06; iconos deformados (7×16 en P-26.08, 19×20 en P-26.10) |
| 5 | Espaciado de 4 pt | **Falla** | separación de 2 px en la lista de la barra lateral (38 pantallas); relleno de 40 px en estados vacíos (P-05, P-08, P-13, P-14, P-16, P-21, P-23, P-26.02, .05, .09, .11, .13, .14) · 4 comb. · §4.1 |
| 6 | Radios 6/10/14 | Pasa | 0 radios fuera de la escala |
| 7 | Sombras (claro) y bordes (oscuro) | Pasa | 0 sombras en oscuro; diálogo P-03 con `shadow.2` en claro |
| 8 | Saturación apagada | Pasa | croma máx. de superficies grandes 0.30 (aviso ámbar oscuro) |
| 9 | Densidad cómoda | **Falla** | P-17 · Ac Ao: el aviso «Lectura B activa…» deja el texto en una columna de 1–2 palabras junto al enlace · §7.14/§7.15 |
| 10 | Barra lateral colapsable | **Falla** | ancho 256 px y no 240 (`layout.sidebar.width`); «Axel» se empalma con el botón «Ocultar barra lateral»: todas · Ec Eo · §7.9 y §4.6 |
| 11 | ≤ 6 tarjetas en el Resumen | Pasa | P-01: 6 |
| 12 | Objetivos táctiles ≥ 44 px | Pasa | los controles < 44 px tienen etiqueta ≥ 44 o área ampliada; barra inferior 82×55 (1 px menos que 56, baja) |
| 13 | Foco visible | Pasa | anillo 2 px de marca con separación de 2 px en 164 de 164 |
| 14 | Paleta de gráficas en orden fijo y dinero en texto | **Falla** | orden fijo y resumen en dinero correctos; ejes con seis marcas «$0» en lugar del estado vacío (P-10 y P-01 Tendencias · Ao Ec · §7.7 «sin ejes falsos») y etiquetas encimadas (P-01 Ec «Devengado/Efectivo acumulado», P-10 Ec leyenda de kcal, P-01 Ao «$1,197.64» cruzada por la barra) |

Fuera de la lista (U5-04 y U5-07, vistos en las capturas): la vista elegida (Tabla/Galería…) y las pestañas Gasto/Ingreso/Peso/Ticket de P-03 se marcan solo con relleno de 1.075:1 (claro) y 1.29:1 (oscuro), casi invisible en oscuro (todas las pantallas con barra de vista · 4 comb. · §7.5.1, que además pide `aria-pressed`).

### Fallas que devuelven la ronda
- **Altas (5)**: U5-03 (5 reglas de axe × 4), U5-01 (P-26.13 Android a 415 px), U5-04 (pares de UI bajo 3:1, incluido el estado elegido), U5-07 (inglés visible en P-26.11 y P-26.04), U5-13 paso 2 (H-5).
- **Medias (6)**: título de página 24 px, barra lateral 256 px con empalme, iconos 16 px, espaciado fuera de 4 pt, aviso comprimido en P-17 Android, gráficas con ejes «$0» y etiquetas encimadas.
- **Bajas (6)**: ver `REVISION-u5-r1.md`.

### Correcciones aplicadas desde la ronda anterior
Ninguna: es la primera ronda.

Veredicto: devuelto

## Ronda 2

Revisor: `revisor-u5-r2` (no capturé ni implementé nada; no corrijo) · 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo con `correccion-u5-r1`. Binario `backend/bin/finanzas` y `frontend/dist` de las 10:33 (ningún fuente más nuevo que `dist`), en marcha en `finanzas.local:8443` con `datos/finanzas-dev.db`. Antes de opinar volví a correr completo el comando del capturista (`npx playwright test -c playwright.u5.config.ts` → **170 passed, 7.9 min**): las 164 capturas son de este binario. Las de las 10:25 eran 136 anteriores a él.

### Comandos reejecutados por el revisor (U5-05)
- **U5-01**: `python3 -I frontend/pruebas-u5/verificar_capturas.py` → T = 41, 4T = 164, **164 de 164 correctos**, 0 faltantes. `ls capturas-u5/ | wc -l` = 168 (164 PNG + `axe/`, `contraste/`, `evidencia/`, `hosts.json`; bloqueo `U5-conteo-carpeta`). Medición propia en el DOM: `scrollWidth = innerWidth` en las 164 cargas (P-26.13 Android = 412).
- **U5-03**: comando literal de U5-03 = **0**. Corrida propia de axe-core (wcag2a/2aa/21a/21aa/22a/22aa) en las 164 cargas, dos veces: **0 serias o críticas y 0 de otro impacto** en Ac, Ao, Ec y Eo. Revisé uno por uno los incompletos que no son de contraste. `aria-hidden-focus` en P-03 y P-22 (3 nodos) son las guardas de foco de Radix y la raíz oculta tras la hoja: en 30 Tab el foco nunca sale de la hoja, así que no es defecto. `aria-prohibited-attr` en P-26.10 es un `div` con `aria-label` y sin rol (B-8).
- **U5-04**: implementación propia (colores computados compuestos en canvas sRGB, fondo efectivo por ancestros, opacidad acumulada, texto grande ≥ 24 px o ≥ 18.66 px con 700, controles por borde, anillo o relleno, estado elegido por borde, relleno o pseudo-elemento). En reposo y con página nueva por pantalla: **0 pares bajo umbral** en texto, íconos y UI. Los 45 estados elegidos por combinación (`aria-pressed`, `aria-selected`, `aria-current`) dan todos ≥ 3:1. 20 pares al azar de los JSON del capturista (semilla 20261012): **20 de 20 coinciden, Δ máx. 0.000**, y los 20 aparecen en mi recorrido del DOM con la misma razón. Resumen del capturista: 4,812 / 4,826 / 6,022 / 6,057 elementos evaluados y 0 bajo umbral. Mi primera pasada marcó 1 par en P-19 (Ac, Ao) y 7 en P-04 a P-10 (Ec, Eo). Los 8 eran el estado *hover* que dejó el ratón de mi propio script. Comprobado: en reposo esos botones no tienen borde ni relleno.

### Pantallas evaluadas
- Lista visual U5-06 (Read de los PNG regenerados): P-01 Resumen en Ac, Ao, Ec y Eo. P-26 y P-26.01 a P-26.14 en Ec y Eo completas; en Android, P-26.09 Ac y medición en el DOM. Una por entrada de la barra: P-02, P-06, P-10, P-13, P-17, P-21, P-23 y P-26 en Ao y Ec; P-02 y P-17 en Ac. P-03 en Ao y Ec. Tendencias de P-01 (Ao y Ec) y de P-10 (Ec), con capturas mías.
- Mediciones en el DOM de las 41 pantallas × 4: h1, barra lateral, iconos, espaciado, radios, sombras, `tabular-nums`, escala tipográfica, toque, foco, fondo, hosts y texto visible.

### Criterios evaluados
| Criterio | Resultado | Evidencia |
|---|---|---|
| U5-01 | Pasa (conteo) | 164 de 164; P-26.13 a 412. Fidelidad Android: ver N-2 |
| U5-02 | Pasa | Corrida completa reejecutada (170 passed). Cifras = `model.json`: $10,000.00, $6,905.50, $973.75, $3,094.50, $6,293.44, $1,191.50 |
| U5-03 | Pasa | 0 en las 4 combinaciones (literal, capturista y mía) |
| U5-04 | Pasa | 0 bajo umbral; 20 de 20 pares, Δ 0.000; estados elegidos ≥ 3:1 |
| U5-06 | **Falla** | Tabla de abajo (N-2, N-3, M-6 parcial, M-1 en P-00) |
| U5-07 | **Falla** | Script propio: 541 importes por combinación, 0 fuera de formato (salvo «$1.5 mil», etiqueta de eje permitida por §7.7), 0 fechas ISO, 0 palabras de la lista, 0 `input file` visibles; «Choose File» y «groceries/eating_out» cerrados. **Nuevo**: P-26.13 muestra «Fixture s2-metas.json, savings_goals[2]» (N-1) |
| U5-08 | Pasa | 1 host: `finanzas.local:8443` (capturista y mío) |
| U5-09 | Pasa | Conteo propio sin red: 3 toques («+», Despensa, Guardar). No volví a comprobar que el gasto aparezca en Inicio y en la tabla (r1 y `flujos.spec`) |
| U5-10 | Pasa | Fondo del body = C2 en 164 de 164 (#FAFAFB / #222325); `tema.spec` dentro de los 170 passed |
| U5-11 | No evaluado | Sin acciones disponibles en esta ronda; pendiente para la ronda 3 |
| U5-12 | Pasa | Supuestos y bloqueos U5-* presentes; nada aprobado con violaciones abiertas |
| U5-13 | Pasa | Spec propia (secuencia H-5): tras recargar sin red sale «Sin conexión con el servidor de casa · 2 en cola» y se puede seguir capturando (3 en cola). Clave de 36 caracteres; `captured_at` ISO con desplazamiento y anterior al envío; 6→8 gastos y 0→1 peso; el reenvío da 200 con el mismo id e `Idempotent-Replayed: true`; un 500 deja «No se pudo enviar · 1 en cola» y la base sin cambio. `captura-sin-conexion` del capturista: 5 passed |
| C3-09 | Pasa | `c3-09.spec`: 111 cifras, 0 diferencias |

### Lista de verificación U5-06 (14 elementos)
| # | Elemento (DESIGN.md) | Resultado | Fallas: pantalla · combinación · regla |
|---|---|---|---|
| 1 | Tokens de color y rol | Pasa | — |
| 2 | Tipografía autoalojada y escala | **Falla** | h1 a 32/700 en 40 pantallas; **P-00 «Finanzas Hogar» a 32/600** · 4 comb. · §3.2 (`font.size.title-1` = 700). Todos los tamaños en la escala; 29–30 pantallas con > 3 tamaños (B-3, bloqueo) |
| 3 | Cifras tabulares | Pasa | 0 de 497 importes sin `tabular-nums` |
| 4 | Iconos Lucide 20/24, trazo 1.5 | Pasa | 0 fuera en 164 |
| 5 | Espaciado de 4 pt | Pasa | Solo quedan el `sr-only` (−1, invisible), los 56 px de estructura (`layout.bottombar.height`) y 2 px dentro del interruptor de P-07 (B-7) |
| 6 | Radios 6/10/14 | Pasa | 0 fuera de la escala |
| 7 | Sombras (claro) y bordes (oscuro) | **Falla** | Secciones de Ajustes con `shadow.1` en reposo: P-26.01 a P-26.14 · Ac Ec · §4.5 y §7.13 (N-3). En oscuro, 0 sombras |
| 8 | Saturación apagada | Pasa | — |
| 9 | Densidad cómoda | Pasa | Franja de P-17/P-18 Android a 314 px en 2 líneas |
| 10 | Barra lateral colapsable | Pasa | 240 px; botón «‹‹» en x 199–231 sin texto debajo |
| 11 | ≤ 6 tarjetas en el Resumen | Pasa | P-01: 6 |
| 12 | Objetivos táctiles ≥ 44 px | **Falla** | En las **capturas** Android los controles salen a 36 px: el capturista las tomó sin puntero grueso. La app real con `any-pointer: coarse` da 44 · Ac Ao · §4.6 (N-2) |
| 13 | Foco visible | Pasa | Anillo de 2 px en 164 de 164 (interior −2 px donde §7.9/§7.10 lo piden) |
| 14 | Gráficas: orden fijo y dinero en texto | **Falla** | Vacíos cerrados («Aún no hay datos para este periodo», sin ejes). Siguen etiquetas encimadas: P-01 Tendencias «Deuda viva proyectada», «$18,025.47» sobre el eje «$18 mil» (Ec) o «$21 mil» (Ao); «Devengado contra efectivo», «Efectivo acumulado» sobre el área (Ec y Ao) y la leyenda de «Efectivo acumulado» sin marca visible (Ec) · §7.7 y §8.5 |

### Fallas que devuelven la ronda
- **Alta (1)**: N-1 · U5-07 · P-26.13 · 4 comb.: texto de semilla con inglés y claves internas visibles.
- **Medias (4)**: N-2 capturas Android sin puntero grueso; N-3 sombra en las secciones de Ajustes en claro; M-6 (reabierto en parte) etiquetas de gráficas encimadas; M-1 (reabierto en P-00) h1 a 600.
- **Bajas**: B-1, B-2, B-3, B-5, B-6 siguen; B-4 no lo volví a revisar. Nuevas: B-7 (interruptor de P-07 con 2 px), B-8 (`aria-label` sin rol en P-26.10), B-9 (P-26.09 «Seleccionar todas» por confirmar).

### Correcciones aplicadas desde la ronda anterior (verificadas)
- Cerrados: A-1 (axe 0), A-2 (P-26.13 a 412), A-3 (estados y controles ≥ 3:1 con `aria-pressed`), A-4 («Choose File» y claves en inglés de P-26.04), A-5 (H-5: la sesión sobrevive a recargar sin red), M-2 (barra de 240 sin empalme), M-3 (iconos), M-4 (espaciado), M-5 (franja de P-17).
- Cerrados en parte: M-1 (falta P-00) y M-6 (faltan las etiquetas encimadas).

Veredicto: devuelto

## Ronda 3

Revisor: `revisor-u5-r3` (no capturé ni implementé nada; no corrijo) · 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo con `correccion-u5-r2`. Binario `backend/bin/finanzas` (11:15:41) y `frontend/dist` (11:15:25) en marcha en `finanzas.local:8443` desde las 11:16; ningún fuente de `frontend/src` ni del backend es más nuevo que el binario. Base `datos/finanzas-dev.db` recreada por la corrección. Antes de opinar volví a correr completo el comando del capturista (`npx playwright test -c playwright.u5.config.ts` → **174 passed, 8.8 min**; las 164 capturas son de 11:34:46 a 11:43:15, todas de esta corrida).

### Comandos reejecutados por el revisor (U5-05)
- **U5-01**: `python3 -I frontend/pruebas-u5/verificar_capturas.py` → T = 41 (27 P-xx + 14 P-26.xx), 4T = 164, **164 de 164 correctos**, 0 faltantes. `ls capturas-u5/ | wc -l` = 168 (bloqueo `U5-conteo-carpeta`). Recorrido propio (página nueva por pantalla): `scrollWidth = innerWidth` en las 164 cargas; P-26.13 Android = 412.
- **U5-03**: comando literal = **0**. axe-core 4.13 propio (wcag2a/2aa/21a/21aa/22a/22aa) en las 164 cargas: **0 serias o críticas y 0 de otro impacto** en Ac, Ao, Ec y Eo. Incompletos 15 · 15 · 6 · 6: `color-contrast` 12/12/3/3 (cubiertos por mi cálculo de contraste), `aria-hidden-focus` 2 (guardas de foco de la hoja, ya revisadas en r2) y `aria-prohibited-attr` 1 (B-8).
- **U5-04**: implementación propia (color computado resuelto en canvas sRGB, fondo efectivo por capas de ancestros, opacidad acumulada, texto grande ≥ 24 px o ≥ 18.66 px con 700, iconos Lucide ≥ 3:1, límite de campos e interruptores ≥ 3:1): 4,327 · 4,333 · 5,567 · 5,597 pares de texto e icono, **0 bajo umbral**, 0 con imagen de fondo. Campos: 99 · 99 · 110 · 110; los 3 marcados por combinación eran el `button[role=switch]` envolvente; la pista medida aparte da 3.47:1 (`#808189` sobre `#F2F2F4`). Capturista: 4,588 · 4,623 · 5,833 · 5,875 y 0 bajo umbral. 20 pares al azar de sus JSON (semilla 20261012): **20 de 20 coinciden, Δ máx. 0.0005**; 15 aparecen en mi recorrido del DOM con la misma razón y los 5 de UI usan los tokens de C2 (`#808189`, `#C9CFFF`, `#B7B8C0`).

### Pantallas evaluadas
- Lista visual U5-06 (Read de los PNG regenerados, en hojas reducidas y recortes): P-00 y P-01 en Ac, Ao, Ec y Eo; P-26 y P-26.01 a P-26.14 en las 4 combinaciones; una por entrada de la barra (P-02, P-06, P-10, P-13, P-17, P-21, P-23 y P-26) en las 4; P-03 (Ac, Ec) y P-22 (Eo); Tendencias de P-01 (`evidencia/tendencias/`, Ec y Ao); barra lateral de P-01 Ec y P-02 Eo a tamaño real.
- Mediciones en el DOM de las 41 pantallas × 4: h1, tamaños de letra, `tabular-nums`, iconos, radios, sombras, toque con prueba de impacto de 44 × 44, puntero grueso, foco, fondo y `data-theme`, hosts y texto visible. Espaciado de 4 pt y tarjetas de P-01 en la muestra de U5-06 (23 pantallas × 4).

### Criterios evaluados
| Criterio | Resultado | Evidencia |
|---|---|---|
| U5-01 | Pasa | 164 de 164; 0 desbordes (P-26.13 a 412) |
| U5-02 | Pasa | Corrida completa reejecutada (174 passed). Cifras = `model.json` (`disponible_A` 10000.0, `disponible_B` 6905.5, `pagos` 973.75, `deuda_viva_10_abr` 6293.44): $10,000.00, $6,905.50, $973.75, $6,293.44, $3,094.50 y $1,191.50 visibles en las 4 combinaciones |
| U5-03 | Pasa | 0 en las 4 combinaciones (literal, capturista y mía) |
| U5-04 | Pasa | 0 bajo umbral (capturista y mía); 20 de 20 pares, Δ 0.0005 |
| U5-06 | **Falla** | Tabla de abajo: barra lateral sin iconos de sección (#10). Los 4 medios de r2 quedan cerrados |
| U5-07 | Pasa | Script propio: 351 importes por combinación, 0 fuera de formato, 0 fechas ISO, 0 palabras de la lista del criterio (y de mi lista ampliada, salvo «Versión dev» en P-26.11, B-10), 0 claves internas, 0 nombres de archivo. P-26.13: «Dato de desarrollo: no existe en producción.» Sin red: «Sin conexión con el servidor de casa · 2 en cola». Leí completo el texto de 10 pantallas (P-01, P-02, P-06, P-10, P-13, P-17, P-21, P-23, P-26.11, P-26.13) |
| U5-08 | Pasa | 1 host: `finanzas.local:8443` (capturista y mío, 164 cargas) |
| U5-09 | Pasa | Prueba propia en servidor aislado (copia en `datos/tmp/`, puerto 8445): 3 toques («+», Despensa, Guardar); la tarjeta «Mes contra plan» cambia y «$187.25» aparece en la tabla de P-02 |
| U5-10 | Pasa | Fondo del body = C2 (`#FAFAFB` / `#222325`) en 164 de 164, sin `data-theme`; `tema.spec` (interruptor) dentro de los 174 passed |
| U5-11 | **Falla** | `C3-integracion.md` v0.95 marca «Capturado: no · pendiente WF-3b» en las 41 filas, aunque las 41 pantallas tienen sus 4 capturas (el cruce literal da 0 pantallas sin captura) |
| U5-12 | Pasa | Supuestos en los reportes; bloqueos `U5-*` (4) presentes; nada aprobado con violaciones abiertas |
| U5-13 | Pasa | Prueba propia (servidor aislado): «Sin conexión con el servidor de casa», 1 → 2 en cola, 3 toques sin red; tras recargar sigue «2 en cola»; claves de 36 caracteres; `captured_at` ISO con desplazamiento y anterior al envío; base 6 → 9 gastos (1 por captura) y 0 → 1 peso; reenvío con la misma clave = 200, `Idempotent-Replayed: true`, mismo id; un 500 deja «No se pudo enviar · 1 en cola» y 0 registros, y al quitarlo llega 1 sola vez. Capturas del estado sin conexión en las 4 combinaciones en `frontend/pruebas-u5/flujos/capturas/` |
| C3-09 | Pasa | `c3-09.spec` reejecutado: 111 cifras, 0 diferencias |

### Lista de verificación U5-06 (14 elementos)
| # | Elemento (DESIGN.md) | Resultado | Fallas: pantalla · combinación · regla |
|---|---|---|---|
| 1 | Tokens de color y rol | Pasa | — |
| 2 | Tipografía autoalojada y escala | Pasa | h1 a 32/700 en 164 de 164 (P-00 incluida); 0 textos fuera de 12/14/16/20/24/32. B-3 sigue (bloqueo) |
| 3 | Cifras tabulares | Pasa | 0 importes sin `tabular-nums` |
| 4 | Iconos Lucide 20/24, trazo 1.5 | Pasa en tamaño y trazo | 0 fuera de 20 px/1.8 y 24 px/1.5 (§5.2). Los iconos que faltan en la barra lateral van en #10 |
| 5 | Espaciado de 4 pt | Pasa | 0 rellenos o separaciones fuera de la retícula en la muestra; solo márgenes `auto` y el −1 px oculto. B-7 no remedido |
| 6 | Radios 6/10/14 | Pasa | 0 fuera de la escala |
| 7 | Sombras (claro) y bordes (oscuro) | Pasa | Secciones de Ajustes sin sombra en las 4 (N-3 cerrado). En claro solo flotan diálogos; quedan bajas B-11 y B-12. Oscuro: 0 sombras |
| 8 | Saturación apagada | Pasa | — |
| 9 | Densidad cómoda | Pasa | — |
| 10 | Barra lateral colapsable | **Falla** | 240 px y «‹‹» en las 40 pantallas con barra, pero **las 9 secciones van sin icono** (solo `chevron-right` y nombre) · Ec Eo · §7.9 («las 9 secciones … con icono `size.icon.sm`») y §5.4 (`house`, `utensils`, `shopping-cart` «barra inferior y lateral») |
| 11 | ≤ 6 tarjetas en el Resumen | Pasa | P-01: 6 en las 4 |
| 12 | Objetivos táctiles ≥ 44 px | Pasa | Puntero grueso en 82 de 82 cargas Android (capturista y mías); 0 controles con área efectiva < 44 × 44 (prueba de impacto) |
| 13 | Foco visible | Pasa | Anillo `outline` sólido de 2 px en cada parada de Tab comprobada |
| 14 | Gráficas: orden fijo y dinero en texto | Pasa | «$18,025.47» a la derecha de los ejes; T2 sin texto sobre el área y leyenda abajo con marcas y valores finales; `--grep Tendencias` 4 de 4 en mi corrida |

### Fallas que devuelven la ronda
- **Medias (2)**: M-7 barra lateral sin iconos de sección (Ec, Eo, 40 pantallas; §7.9 y §5.4); M-8 U5-11: la matriz de C3 no marca «capturado» ni la ruta de las 4 imágenes (41 filas).
- **Bajas**: nuevas B-10 («Versión dev», P-26.11), B-11 (pestaña elegida de P-03/P-22 con `shadow.1` en claro), B-12 (interruptor: pista de 40 × 24 y perilla con sombra), B-13 (P-00 dice «Tu sesión terminó» sin sesión previa), B-14 («Finanzas H…» truncado en el encabezado de la barra), B-15 (hoja de P-03/P-22 en Android sin `shadow.2`). Siguen B-1, B-2, B-3, B-5, B-6, B-7 y B-8. B-9 cerrado (0 controles < 44 px).

### Correcciones aplicadas desde la ronda anterior (verificadas)
- Cerrados: **N-1** (alta, U5-07: P-26.13 sin «Fixture», nombre de archivo ni clave interna), **N-2** (capturas Android con puntero grueso, 82 de 82), **N-3** (secciones de Ajustes sin sombra), **M-6b** (etiquetas de Tendencias sin encimarse), **M-1b** (h1 de P-00 a 32/700).

Veredicto: devuelto

Es la ronda 3: según U5-05 no hay ronda 4 y se escala a Axel (`docs/estado/bloqueos/U5-ronda3-escalamiento.md`).
