# DESIGN.md · Finanzas Hogar

Versión 1.0 en construcción (11-oct-2026, WF-1). Secciones 1 a 6: agente U2b-1-design. Secciones 7 a 10: agente U2b-2. Fuentes de verdad: `docs/plan/plan.md` §4 (esqueleto y método), `docs/plan/decisiones.md` §G, §I y §K, criterios `docs/criterios/U2.md` 1.2, colorimetría 1.2 (`docs/design/colorimetria.md`) y sus tokens `docs/contratos/C2-colores.json` (190 tokens). Este documento **no define ningún color nuevo**: todo hex citado sale de `C2-colores.json` y lo comprueba `docs/design/fuentes/scripts/verificar.py`. Los nombres de token están en inglés corto (identificadores de código); todo lo demás, en español de México.

Matiz aprobado: **C · Índigo 280°** (decisiones §I). Referencia visual única: **Notion** (decisiones §G).

## 1. Principios y tono

### 1.1 Referencia visual

Notion es la única referencia visual (decisiones §G). De Notion se toma: barra lateral jerárquica y colapsable, páginas hechas de bloques con propiedades arriba, vistas de base de datos (tabla, tablero, calendario, galería), menú `/`, edición en el lugar, superficies casi blancas en claro y casi negras en oscuro, tipografía neutra y color reservado para significar. No se cita ni se imita ningún otro producto.

### 1.2 Principios

| # | Principio | Qué significa en la interfaz | Cómo se comprueba |
|---|---|---|---|
| P1 | **Completo pero minimalista** | Todo dato existe, pero cada pantalla muestra solo lo que pide su decisión; lo demás queda a un paso (*progressive disclosure*): «Mostrar N más», menú «⋯», vista previa del registro (U1 G-06). El Resumen del dashboard lleva ≤ 6 tarjetas (decisiones §G). | Revisión por pantalla en U5: ningún dato secundario compite con la cifra principal. |
| P2 | **Una acción primaria por pantalla** | Un solo botón con relleno de marca visible a la vez (§2.9); las demás acciones son secundarias, de texto o van en «⋯». En móvil la primaria nunca se esconde (U1 §1.3). | Conteo de botones primarios por pantalla = 1 (las hojas y diálogos cuentan como pantalla propia). |
| P3 | **Móvil primero** | Se diseña a 360 px y se amplía (U1 §6.1). Registrar un gasto = 3 toques (K10). Todo control táctil mide ≥ 44 × 44 px (§4.6). PWA instalable en Android (decisiones §G). | U5 a 360 y 412 px; conteo de toques de F-1. |
| P4 | **El color significa, no decora** | Regla 60-30-10 y saturación apagada (§2.2); el color de un rol, persona, tienda o serie **nunca va solo**: lo acompaña signo, icono, forma o texto (§2.10). | Vista en escala de grises sin pérdida de significado. |
| P5 | **Cifras exactas y legibles** | Dinero `$1,258.68` (U1 §0.4) con cifras tabulares (§3.4); nunca se trunca ni se parte (U1 §6.1 regla 5); nunca se anima (§6.4). | `$99,999.99` completo en la tarjeta más estrecha a 360 px. |
| P6 | **Calma** | Densidad cómoda, sombras leves solo en claro (§4.5), movimiento corto (§6), sin rojos ni destellos salvo error y acción crítica. | Revisión de U5 contra §2 y §6. |
| P7 | **Claro y oscuro de primera** | Tema por defecto según el sistema, con interruptor (decisiones §G); cada token de color tiene variante clara y oscura. | Toda captura de U5 existe en los dos temas. |
| P8 | **Accesible por cálculo** | WCAG 2.2 AA mínimo, calculado: texto ≥ 4.5:1; texto grande, interfaz y gráficas ≥ 3:1 (decisiones §G). | Dos tablas oficiales, 0 bajo umbral: `contraste.csv` (paleta: 202 pares, §2.8) y `contraste-componentes.csv` (componentes: 316 combinaciones, §10.5). |

### 1.3 Tono de voz

Texto de interfaz en español de México, claro y breve, para una familia (Axel, Ximena y Karen), no para un banco.

- **Trato de tú** y frases cortas. Botones en infinitivo: «Guardar», «Registrar gasto», «Cambiar semana».
- **Sin jerga** financiera ni técnica: «Lo que debes hoy», no «saldo insoluto»; los términos del dominio son los del diccionario de datos.
- **Cifras con su formato**: `$1,258.68`, `10-oct-2026`, `80 %`, tasas como se capturaron (`17.27 %`; K12).
- **Mensajes que dicen qué pasó y qué hacer**: «Falta capturar la tasa de esta tarjeta», no «Error de validación».
- **Sin signos de exclamación ni emojis** en mensajes del sistema; el éxito es discreto: «Guardado · Deshacer» (U1 G-07).
- Textos ya fijados en U1 y decisiones que se respetan tal cual: «Sin conexión con el servidor de casa», «faltan restricciones» (K16), «falta capturar» (K18).

| Bien | Evitar | Por qué |
|---|---|---|
| «Guardado · Deshacer» | «¡Gasto guardado exitosamente!» | Breve, sin exclamación, con salida. |
| «Te quedan $1,320.50 de despensa este mes» | «Disponible: 1320.5» | Dice qué es y usa el formato de dinero. |
| «Este pago vence mañana» | «Pago próximo a vencer (D-1)» | Sin códigos internos. |

## 2. Colorimetría

Resumen fiel de la colorimetría 1.2 (`docs/design/colorimetria.md`; la 1.2 no cambia ningún valor de la 1.1, solo quita marcas provisionales con K41 y K42). El detalle, las tablas completas y la regeneración viven en ese documento y en `docs/design/tools/colorimetria.py`; aquí van las reglas que usa quien diseña o programa.

### 2.1 Método

| Elemento | Valor | Justificación |
|---|---|---|
| Espacio de color | OKLCH, conversión a sRGB de 8 bits | Luminosidad perceptual uniforme entre matices (plan §4.2). |
| Escalas | 11 pasos, 50 a 950 | Pasos suficientes para superficie, borde, texto y relleno en ambos temas. |
| L fija por paso | 0.985, 0.962, 0.925, 0.865, 0.785, 0.695, 0.605, 0.515, 0.425, 0.335, 0.255 | Mismo paso = misma luminosidad percibida en cualquier matiz. |
| Croma | Perfil 0.10 → 1.0 → 0.45 del máximo; máx. 0.11 marca, 0.14 roles, 0.012 neutro | Saturación apagada, tipo Notion (decisiones §G). |
| Gamut | Recorte de croma al sRGB por bisección | Ningún hex sale del gamut; se recorta croma, nunca L. |
| Contraste | WCAG 2.x, `(L1 + 0.05) / (L2 + 0.05)` con luminancia relativa sRGB | Umbrales: texto ≥ 4.5:1; texto grande, interfaz y gráficas ≥ 3:1. |
| Daltonismo | Protanopia, deuteranopia y tritanopia simuladas, más escala de grises | Ningún significado depende de distinguir matices (colorimetría §2.6). |

### 2.2 Matiz base, armonía y proporción 60-30-10

- **Base 280° (índigo)**: marca, superficies tintadas (neutro con croma 0.012) y acento principal.
- **Eje opuesto 100° (oliva amarillo)**: reservado **solo** para «acción crítica» de confirmación (decisiones §I). **Zona reservada 80°–120°**: ningún otro token con croma visible usa un matiz a menos de 20° del opuesto (colorimetría §3).
- **Complementario dividido**: acentos a ±30° del opuesto → advertencia 75° e ingreso 135°. Error y gasto (rojos) y éxito (verde) son convenciones semánticas fuera de la armonía, por eso siempre llevan icono y texto.

| Proporción | Papel | Claro | Oscuro |
|---|---|---|---|
| 60 % | Superficies y texto neutro | `color.neutral.50` fondo, `color.neutral.100` tarjeta, `color.neutral.200` divisor, `color.neutral.700` y `color.neutral.900` texto | `color.neutral.950` fondo, `color.neutral.900` tarjeta, `color.neutral.800` divisor, `color.neutral.300` y `color.neutral.100` texto |
| 30 % | Acento de marca en pasos de superficie: elemento activo, fila seleccionada, chips de filtro | `color.brand.50` a `color.brand.200` | `color.brand.900` y `color.brand.950` |
| 10 % | Énfasis: la acción primaria, enlaces, foco, roles, categóricas, series | `color.brand.700`, roles paso 700, categóricas y series `light` | `color.brand.300`, roles paso 400, categóricas y series `dark` |

### 2.3 Escalas de marca y neutro (280°)

Coinciden paso a paso con los valores de partida del plan §4.3 (marca `#F9FAFF … #1E203A`, neutro `#FAFAFB … #222325`; colorimetría §5.1–5.2).

| Paso | L | Marca (croma máx. 0.11) | Neutro tintado (croma máx. 0.012) |
|---|---|---|---|
| 50 | 0.985 | `color.brand.50` · `#F9FAFF` | `color.neutral.50` · `#FAFAFB` |
| 100 | 0.962 | `color.brand.100` · `#F0F1FF` | `color.neutral.100` · `#F2F2F4` |
| 200 | 0.925 | `color.brand.200` · `#E1E4FF` | `color.neutral.200` · `#E5E6E9` |
| 300 | 0.865 | `color.brand.300` · `#C9CFFF` | `color.neutral.300` · `#D1D2D8` |
| 400 | 0.785 | `color.brand.400` · `#ADB3F4` | `color.neutral.400` · `#B7B8C0` |
| 500 | 0.695 | `color.brand.500` · `#9095E0` | `color.neutral.500` · `#9B9CA4` |
| 600 | 0.605 | `color.brand.600` · `#7579C2` | `color.neutral.600` · `#808189` |
| 700 | 0.515 | `color.brand.700` · `#5D609F` | `color.neutral.700` · `#66676E` |
| 800 | 0.425 | `color.brand.800` · `#46487B` | `color.neutral.800` · `#4D4E53` |
| 900 | 0.335 | `color.brand.900` · `#303258` | `color.neutral.900` · `#36363A` |
| 950 | 0.255 | `color.brand.950` · `#1E203A` | `color.neutral.950` · `#222325` |

Blanco puro: `color.base.white` · `#FFFFFF`, solo como texto sobre rellenos 700 en claro.

### 2.4 Roles semánticos

Texto e icono del rol: paso 700 en claro y 400 en oscuro (tokens `light` y `dark`); insignia: paso 100 en claro y 900 en oscuro (`surface-light` y `surface-dark`). El paso 600 del plan no llega a 4.5:1 en claro con ningún matiz, por eso el texto usa 700 (colorimetría §6.2). La escala completa 50–950 de cada rol está en `color.role.<rol>.*` para los estados vecinos (hover: 800 en claro, 200 en oscuro).

| Rol | Matiz | Claro (texto) | Oscuro (texto) | Refuerzo obligatorio (nunca solo color) |
|---|---|---|---|---|
| Ingreso | 135° | `color.role.income.light` · `#457723` | `color.role.income.dark` · `#98CA7B` | Signo «+» y flecha de entrada (`arrow-down-left`) en cifras. |
| Gasto | 45° (terracota) | `color.role.expense.light` · `#A04B20` | `color.role.expense.dark` · `#F8A078` | Signo «−» y flecha de salida (`arrow-up-right`) en cifras. |
| Ahorro | 200° | `color.role.savings.light` · `#00767B` | `color.role.savings.dark` · `#3FD0D7` | Icono `piggy-bank` y etiqueta. |
| Deuda / crédito | 310° | `color.role.debt.light` · `#7E509D` | `color.role.debt.dark` · `#D0A3F2` | Icono `credit-card` y etiqueta. |
| Advertencia | 75° | `color.role.warning.light` · `#8B5D00` | `color.role.warning.dark` · `#E5AE5B` | Icono `triangle-alert` y texto. |
| Error | 20° (rojo) | `color.role.error.light` · `#A34449` | `color.role.error.dark` · `#FC9999` | Icono `circle-alert` y mensaje; solo en validación y avisos. |
| Éxito | 150° | `color.role.success.light` · `#207B3E` | `color.role.success.dark` · `#7ECF90` | Palomita (`circle-check`) y texto («Guardado», «Pagado»); nunca en cifras. |
| Información | 240° | `color.role.info.light` · `#006FA4` | `color.role.info.dark` · `#6BC3FD` | Icono `info` y texto. |

Decisiones de matiz que este resumen respeta:

- **Gasto en terracota 45°** (K26: rango 40–45°, error se queda en 20°). De los seis grados del rango, 45° es el más separado de error (24.9°) y de advertencia (30.7°) medidos en el hex (colorimetría «Cambios 1.1», punto 1).
- **Ingreso**: decisiones §I lo deja en **150°** «o U2 propone 135°». Valor vigente en los tokens: **135°** (propuesta de U2a ya prevista por decisiones §I). Motivo: separa ingreso de éxito 15° en vez de 0° y forma con advertencia el par de acentos de la armonía. Aun así, lado a lado ingreso y éxito apenas se distinguen (ΔE ≈ 3), por eso ingreso vive en cifras con «+» y éxito solo en avisos con palomita (colorimetría §6.4).
- **Gasto frente a error**: ya no colisionan (≥ 24.9° de matiz), pero con daltonismo rojo-verde siguen cerca: el gasto solo aparece como cifra con «−», el error solo como mensaje con icono; nunca comparten componente (colorimetría §6.5).

### 2.5 Acción crítica (eje opuesto 100°)

| Uso | Claro | Oscuro |
|---|---|---|
| Relleno del botón que confirma algo irreversible (borrar datos, restaurar un respaldo) | `color.critical.700` · `#766800` con texto `color.base.white` (5.00:1 sobre tarjeta) | `color.critical.300` · `#E0D58E` con texto `color.critical.950` |
| Texto o icono de acción crítica fuera del botón | `color.critical.light` · `#766800` | `color.critical.dark` · `#CBBB5A` |

Solo aparece dentro del diálogo de confirmación, una vez, y siempre con verbo explícito («Borrar 12 gastos»). No se usa para error (que es rojo 20°) ni para advertir. «Acción crítica» no es uno de los 8 roles: es el énfasis del eje opuesto.

### 2.6 Paleta categórica: personas y tiendas

Cada entrada tiene el **mismo matiz** en claro y oscuro y distinta L (la identidad se conserva al cambiar de tema); dentro de cada grupo la L se escalona con ΔL ≥ 0.10, así que se distinguen también en gris (colorimetría §7.3). **Orden fijo** de leyenda, apilado y columnas: personas Axel, Ximena, Karen, Conjunto; tiendas Walmart, Bodega Aurrera, Sam's Club, Mercado.

| Grupo | Nombre | Matiz | Claro | Oscuro |
|---|---|---|---|---|
| Persona | Axel | 256° | `color.person.axel.light` · `#002E64` | `color.person.axel.dark` · `#4F83C8` |
| Persona | Ximena | 350° | `color.person.ximena.light` · `#A14E78` | `color.person.ximena.dark` · `#FFACD3` |
| Persona | Karen | 333° | `color.person.karen.light` · `#B870AC` | `color.person.karen.dark` · `#D389C6` |
| Persona | Conjunto | 280° (croma 0.045) | `color.person.shared.light` · `#484B66` | `color.person.shared.dark` · `#E7EAFF` |
| Tienda | Walmart | 262° | `color.store.walmart.light` · `#4065AC` | `color.store.walmart.dark` · `#7AA3EE` |
| Tienda | Bodega Aurrera | 356° | `color.store.bodega-aurrera.light` · `#762C4C` | `color.store.bodega-aurrera.dark` · `#FFAAC9` |
| Tienda | Sam's Club | 175° | `color.store.sams-club.light` · `#009D82` | `color.store.sams-club.dark` · `#99FFE4` |
| Tienda | Mercado | 28° | `color.store.market.light` · `#5C0204` | `color.store.market.dark` · `#C2665B` |

- **Karen 333° y Bodega Aurrera 356°, en rosa con refuerzo declarado (K41, FIRME).** Con gasto en 45° el rango cálido quedó cerrado (entre error 20°, gasto 45° y advertencia 75° no hay matiz a ≥ 20° de los tres) y salieron de 50° hacia el hueco entre deuda 310° y error 20° (colorimetría «Cambios 1.1», punto 2). Frente a Ximena quedan bajo «≥ 20° o ΔL ≥ 0.15» y se distinguen por **orden fijo, forma y etiqueta directa**; entre categóricas basta ese refuerzo, también entre grupos distintos (U2-02, U2-08). Frente a los roles y a acción crítica la separación sí es obligatoria y se cumple (colorimetría §7.4).
- **Walmart en oscuro en L 0.716 (K42, FIRME)**: `color.store.walmart.dark` · `#7AA3EE`, para que las cuatro tiendas se escalonen con ΔL ≥ 0.10 en oscuro.
- **Conjunto** lleva el matiz de la marca con croma bajo: se lee como «de los dos», no como otra persona. Es la única categórica bajo el piso de croma 0.10, a propósito.

**Regla K41 (obligatoria): una persona y una tienda nunca se distinguen solo por color en una misma vista.** Si una vista muestra a la vez colores de personas y de tiendas (por ejemplo, gasto por persona y por tienda), cada color va con su nombre en texto y su forma (punto, anillo de avatar o marca de serie); personas y tiendas **nunca comparten gráfica ni leyenda**, y en tablas van en columnas distintas. La regla cubre el par Bodega Aurrera–Ximena en oscuro (ΔE 1.3: a simple vista, el mismo color).

Reglas de uso de las categóricas (colorimetría §7.6):

1. Nunca llevan texto encima: el nombre va en texto neutro junto a un punto, anillo de avatar o marca de serie del color.
2. En gráficas por persona o por tienda el color sigue a la entidad (no la paleta de series), con leyenda siempre y etiqueta directa cuando hay ≤ 4 elementos.
3. Los mercados adicionales usan el color de Mercado; si una gráfica necesita distinguir varios mercados, usa la paleta de series en su orden (§2.7).
4. La marca solo vive en controles (botón primario, enlace subrayado, anillo de foco, elemento activo) y nunca marca a una persona o tienda.

Pares cercanos (Δ matiz < 20° **y** ΔL < 0.15) y su refuerzo declarado, medidos en el hex (colorimetría §7.7; ninguno es categórica ni marca frente a rol):

| Tema | Par | Δ matiz | ΔL | Refuerzo |
|---|---|---|---|---|
| claro | Ingreso – Éxito (roles) | 14.9° | 0.000 | «+» y flecha de entrada en cifras; éxito solo con palomita y texto en avisos. |
| claro | Ximena – Karen (personas) | 17.1° | 0.103 | Mismo grupo: orden fijo, etiqueta directa y ΔL ≥ 0.10 (K41). |
| claro | Ximena – Bodega Aurrera | 5.8° | 0.130 | Persona y tienda: nunca comparten gráfica ni leyenda; columnas distintas; nombre siempre (K41). |
| claro | Conjunto – Walmart | 17.1° | 0.093 | Persona y tienda: ídem (K41). |
| claro | Conjunto – Marca 700 | 0.9° | 0.094 | La marca solo en controles; la categórica como punto o anillo junto a su nombre. |
| claro | Walmart – Marca 700 | 18.0° | 0.000 | Ídem. |
| oscuro | Ingreso – Éxito (roles) | 15.3° | 0.001 | Igual que en claro. |
| oscuro | Axel – Walmart | 5.9° | 0.112 | Persona y tienda (K41). |
| oscuro | Ximena – Karen (personas) | 16.7° | 0.111 | Mismo grupo: orden fijo, etiqueta directa y ΔL ≥ 0.10 (K41). |
| oscuro | Ximena – Bodega Aurrera | 6.2° | 0.006 | Persona y tienda (K41). |
| oscuro | Walmart – Marca 400 | 18.0° | 0.069 | La marca solo en controles. |

Colisiones que nombran las decisiones, resueltas: Karen 80° y Bodega Aurrera 60° frente a advertencia 75° (se movieron a 333° y 356°); error 20° frente a gasto (con 45° quedan a 24.9°–25.3°); ingreso frente a éxito (135° y refuerzo); Sam's Club 175° frente a éxito 150° y ahorro 200° (24.7° medidos en el hex); información 240° frente a marca 280° (40°).

### 2.7 Paleta de gráficas (8 series)

Series en **orden fijo de uso**, nunca recicladas; una novena serie se agrupa en «Otros» (`color.neutral.600` en claro, `color.neutral.400` en oscuro). Alternan dos clases de luminosidad (impares B baja, pares A alta) para separarse en gris: ΔL mínimo entre series consecutivas **0.162 en claro y 0.198 en oscuro** (≥ 0.10, K27); los 28 pares se reportan solo como información (colorimetría §8.3–8.5). Contraste ≥ 3:1 sobre fondo y tarjeta en ambos temas (mínimo 3.14 en claro y 3.10 en oscuro).

| Serie | Nombre | Matiz | Claro | Oscuro |
|---|---|---|---|---|
| 1 | índigo | 280° | `color.chart.1.light` · `#474698` | `color.chart.1.dark` · `#797DCD` |
| 2 | verde azulado | 180° | `color.chart.2.light` · `#009A87` | `color.chart.2.dark` · `#61E4CD` |
| 3 | óxido | 45° | `color.chart.3.light` · `#883500` | `color.chart.3.dark` · `#C16D45` |
| 4 | naranja | 55° | `color.chart.4.light` · `#BF6D2C` | `color.chart.4.dark` · `#FFB989` |
| 5 | magenta | 345° | `color.chart.5.light` · `#812E63` | `color.chart.5.dark` · `#B86797` |
| 6 | azul cielo | 225° | `color.chart.6.light` · `#0092B9` | `color.chart.6.dark` · `#76D9FF` |
| 7 | violeta | 310° | `color.chart.7.light` · `#683987` | `color.chart.7.dark` · `#9C70BC` |
| 8 | verde | 140° | `color.chart.8.light` · `#579749` | `color.chart.8.dark` · `#A0DE93` |

El color de una serie nunca va solo: cada serie lleva **forma de marcador, patrón de trazo o relleno y etiqueta directa** (K27). La simulación de daltonismo con los hex simulados de cada serie está en colorimetría §8.4. Las gráficas por persona o por tienda no usan esta paleta sino la categórica (§2.6). Los patrones de cada tipo de gráfica los fija la sección 8.

### 2.8 Contrastes estructurales

Recalculados por `verificar.py` con la fórmula de §2.1 (truncados a 2 decimales, como `contraste.csv`). La tabla completa de la colorimetría tiene **202 pares, 202 cumplen, 0 bajo umbral** (101 por tema; texto más ajustado 4.73:1, interfaz más ajustada 3.05:1). La segunda tabla oficial es la de componentes: `contraste-componentes.csv`, 316 combinaciones, 0 bajo umbral (§10.5).

| Par | Primer plano | Fondo | Plan §4.3 | Medido | Uso |
|---|---|---|---|---|---|
| Texto / fondo (claro) | `color.neutral.900` | `color.neutral.50` | 11.4 | **11.53** | texto normal; difiere 0.13 del plan, rige el medido |
| Texto secundario / fondo (claro) | `color.neutral.700` | `color.neutral.50` | 5.3 | 5.39 | texto secundario y marcador de posición |
| Neutro 600 / tarjeta (claro) | `color.neutral.600` | `color.neutral.100` | 3.45 | 3.46 | solo texto grande o iconos; no se usa para texto |
| Blanco / marca 600 | `color.base.white` | `color.brand.600` | 3.98 | 3.98 | **prohibido para texto** (tampoco en botones) |
| Marca 700 / fondo (claro) | `color.brand.700` | `color.neutral.50` | 5.1–5.5 | 5.53 | botón primario y enlaces sobre el fondo |
| Blanco / marca 700 | `color.base.white` | `color.brand.700` | 5.8 | 5.77 | texto del botón primario en claro |
| Texto / fondo (oscuro) | `color.neutral.100` | `color.neutral.950` | 14.1 | 14.06 | texto normal en oscuro |
| Texto 300 / tarjeta (oscuro) | `color.neutral.300` | `color.neutral.900` | 7.9 | 7.97 | texto secundario en oscuro |
| Marca 300 / fondo (oscuro) | `color.brand.300` | `color.neutral.950` | 10.4–10.7 | 10.37 | botón primario y enlaces en oscuro |
| Marca 950 / marca 300 | `color.brand.950` | `color.brand.300` | 10.4 | 10.47 | texto del botón primario en oscuro |
| Divisor (claro) | `color.neutral.200` | `color.neutral.50` | 1.2 | 1.19 | **decorativo** (exento) |
| Divisor (oscuro) | `color.neutral.800` | `color.neutral.950` | 1.9 | 1.89 | **decorativo** (exento) |
| Borde de campo (claro) | `color.neutral.600` | `color.neutral.50` | ≥ 4.2 | **3.71** | borde de campo y botón secundario; difiere del plan, cumple 3:1 |
| Borde de campo (oscuro) | `color.neutral.400` | `color.neutral.900` | — | 6.08 | borde de campo y botón secundario |
| Anillo de foco (claro) | `color.brand.600` | `color.neutral.100` | — | 3.56 | anillo de 2 px con separación de 2 px |
| Anillo de foco (oscuro) | `color.brand.400` | `color.neutral.900` | — | 6.03 | ídem |
| Texto más ajustado (claro) | `color.role.success.light` | `color.neutral.100` | — | 4.73 | cifra o texto de rol sobre tarjeta |
| Interfaz más ajustada (claro) | `color.store.sams-club.light` | `color.neutral.100` | — | 3.05 | punto o serie de Sam's Club sobre tarjeta |
| Interfaz más ajustada (oscuro) | `color.store.market.dark` | `color.neutral.900` | — | 3.06 | punto o serie de Mercado sobre tarjeta |

### 2.9 Uso por componente

| Elemento | Claro | Oscuro |
|---|---|---|
| Fondo de página · tarjeta | `color.neutral.50` · `color.neutral.100` | `color.neutral.950` · `color.neutral.900` |
| Texto normal · secundario | `color.neutral.900` · `color.neutral.700` | `color.neutral.100` · `color.neutral.300` |
| Texto deshabilitado | `color.neutral.700` sobre relleno `color.neutral.100`, borde punteado | `color.neutral.400` sobre relleno `color.neutral.900`, borde punteado |
| Enlace (siempre subrayado) · hover | `color.brand.700` · `color.brand.800` | `color.brand.300` · `color.brand.200` |
| **Botón primario** · hover | relleno `color.brand.700` · `color.brand.800`, texto `color.base.white` | relleno `color.brand.300` · `color.brand.200`, texto `color.brand.950` |
| Botón secundario · hover | borde `color.neutral.600`, texto `color.neutral.900` · relleno `color.neutral.200` | borde `color.neutral.400`, texto `color.neutral.100` · relleno `color.neutral.800` |
| Botón destructivo · hover | relleno `color.role.error.700` · `color.role.error.800`, texto `color.base.white` | relleno `color.role.error.300` · `color.role.error.200`, texto `color.role.error.950` |
| Botón de acción crítica · hover | relleno `color.critical.700` · `color.critical.800`, texto `color.base.white` | relleno `color.critical.300` · `color.critical.200`, texto `color.critical.950` |
| Elemento activo, fila seleccionada | relleno `color.brand.100`, texto `color.neutral.900`, icono `color.brand.700` | relleno `color.brand.900`, texto `color.neutral.100`, icono `color.brand.300` |
| Insignia de rol | texto `color.role.<rol>.light` sobre `color.role.<rol>.surface-light` | texto `color.role.<rol>.dark` sobre `color.role.<rol>.surface-dark` |
| Cifra de rol (con signo e icono) | `color.role.<rol>.light` | `color.role.<rol>.dark` |

El blanco sobre marca 600 (3.98:1) no aparece en ningún texto (U2-07). El estado de cada componente (reposo, hover, foco, activo, deshabilitado, error) lo detalla la sección 7.

### 2.10 Reglas de uso del color

1. **El color nunca es el único portador de significado** (WCAG 1.4.1): cifras con signo, estados con icono y texto, series con forma y etiqueta, categóricas con nombre.
2. Un rol se usa solo para su significado: el terracota de gasto no decora, el verde de éxito no marca ingresos, el oliva de acción crítica solo confirma.
3. Una persona y una tienda nunca se distinguen solo por color en una misma vista (K41, §2.6).
4. Ningún color nuevo fuera de `C2-colores.json`; transparencias y sombras se derivan de tokens existentes (§4.5), nunca de un hex suelto.
5. Todo par nuevo de primer plano y fondo entra, antes de usarse, en la tabla oficial que le toca: un par de paleta (rol, categórica o serie sobre una superficie) se agrega a `docs/design/tools/colorimetria.py` y sale en `docs/design/contraste.csv`; un par de componente se agrega a la lista de §10.5 y `docs/design/tools/contraste_componentes.py` lo calcula desde `C2-tokens.json` en `docs/design/contraste-componentes.csv`.

## 3. Tipografía

### 3.1 Familia y archivos

**Inter 4.1, variable, licencia SIL Open Font License 1.1, autoalojada** (propuesta del plan §4.1; decisiones §G pide tipografía autoalojada). Se eligió sobre las fuentes del sistema porque da la misma métrica en Android y Linux, trae cifras tabulares (`tnum`) y tamaño óptico, y su tono neutro encaja con la referencia de Notion. Archivos en `docs/design/fuentes/inter/` (procedencia y sha256 en `docs/design/fuentes/FUENTES.md`):

| Archivo | Para qué | Ejes medidos en el archivo |
|---|---|---|
| `web/InterVariable.woff2` (352,240 B) | interfaz web, redonda | `wght` 100–900 · `opsz` 14–32 |
| `web/InterVariable-Italic.woff2` (387,976 B) | cursiva en notas y citas | `wght` 100–900 · `opsz` 14–32 |
| `ttf/InterVariable.ttf` y `ttf/InterVariable-Italic.ttf` | documentos PDF generados en local y paquete de diseño | los mismos |
| `LICENSE.txt` | licencia OFL 1.1 | — |

Carga sin servicios de terceros (WF-3 copia los `woff2` a los recursos estáticos de la PWA; la ruta `/fonts/` es la propuesta):

```css
@font-face {
  font-family: "Inter";
  src: url("/fonts/InterVariable.woff2") format("woff2");
  font-weight: 100 900;
  font-style: normal;
  font-display: swap;
}
@font-face {
  font-family: "Inter";
  src: url("/fonts/InterVariable-Italic.woff2") format("woff2");
  font-weight: 100 900;
  font-style: italic;
  font-display: swap;
}
:root {
  --font-sans: "Inter", ui-sans-serif, system-ui, "Noto Sans", Roboto, sans-serif, "Noto Color Emoji";
  font-optical-sizing: auto;
}
```

- **Respaldo**: `ui-sans-serif, system-ui, "Noto Sans", Roboto` son las fuentes que ya traen Android y Linux, de métrica cercana; `"Noto Color Emoji"` para los iconos de página tipo Notion.
- **`font-display: swap`**: el texto se lee de inmediato con la fuente de respaldo; la PWA precarga los dos `woff2`, así que el cambio solo se ve en la primera visita.
- **Tamaño óptico automático** (`opsz` 14–32): el navegador ajusta el dibujo al tamaño; por eso el espaciado entre letras queda en 0 en toda la escala.

### 3.2 Escala

Seis tamaños (12/14/16/20/24/32), la misma escala en todos los anchos; base 16 px = 1 rem. Interlíneas múltiplos de 4 para caer en la retícula de 4 pt (§4.1).

| Token | px | rem | Interlínea | Peso | Uso |
|---|---|---|---|---|---|
| `font.size.caption` | 12 | 0.75 | 16 px | 400 (500 en etiquetas) | metadatos, etiquetas de eje, pies de tabla, insignias, fechas secundarias |
| `font.size.small` | 14 | 0.875 | 20 px | 400 | interfaz densa: barra lateral, tablas, propiedades, menús, botones en escritorio |
| `font.size.body` | 16 | 1 | 24 px | 400 | texto de páginas y notas, campos de formulario, botones en móvil |
| `font.size.title-3` | 20 | 1.25 | 28 px | 600 | títulos de sección y de tarjeta, cifra secundaria |
| `font.size.title-2` | 24 | 1.5 | 32 px | 600 | cifra principal de una tarjeta de métrica, título de hoja o diálogo |
| `font.size.title-1` | 32 | 2 | 40 px | 700 | título de página |

| Token | Valor | Uso |
|---|---|---|
| `font.weight.regular` | 400 | texto corrido y celdas |
| `font.weight.medium` | 500 | etiquetas, elemento activo, botones |
| `font.weight.semibold` | 600 | títulos 20 y 24, cifras principales |
| `font.weight.bold` | 700 | solo título de página |
| `font.letter-spacing` | 0 | el tamaño óptico ya corrige el espaciado |

Decisiones:

- **Nombres por papel, no por tamaño de Tailwind**: `caption … title-1` evita confundirlos con la escala por defecto de Tailwind (donde `text-lg` es 18 px) y deja la escala cerrada en seis valores.
- **Campos de formulario a 16 px** también en móvil: el navegador no amplía la página al enfocar y el monto de F-1 se lee de un vistazo.
- **Tablas y barra lateral a 14 px**: es la densidad de Notion en vistas de base de datos y cabe en 360 px con la cifra completa.
- **Máximo tres tamaños por pantalla además de `caption`**: la jerarquía la dan el peso y el color neutro (900 frente a 700), no muchos tamaños.

### 3.3 Cifras y dinero

| Regla | Valor | Justificación |
|---|---|---|
| Cifras tabulares en **todo** dinero (también porcentajes, cantidades y fechas en tablas) | `font-variant-numeric: tabular-nums` (rasgo `tnum`, presente en los dos archivos de Inter) | Columnas alineadas y cifras que no «bailan» al cambiar. |
| Formato | `$1,258.68`: MXN, símbolo `$`, coma de miles, punto y dos decimales (U1 §0.4) | Formato ya decidido. |
| Alineación en tablas | a la derecha; con dos decimales fijos y `tnum` los puntos decimales quedan en columna | Comparar montos de un vistazo. |
| Ingreso | `+$1,258.68`, `color.role.income.*` e icono `arrow-down-left` | Signo, color e icono: nunca solo color. |
| Gasto o negativo | `−$1,258.68`: signo menos tipográfico U+2212 (presente en Inter) antes de `$`, sin espacio, color del rol e icono | El guion corto (U+002D) se confunde con un guion de texto. |
| Negativo sin rol propio (p. ej., restante rebasado) | `−$230.00` con `color.role.warning.*`, icono `triangle-alert` y la palabra «rebasado» | Advertencia es el estado de un límite superado; error queda para validación (§2.4). |
| Montos neutros (totales, saldos sin rol) | sin signo, `color.neutral.900` (claro) o `color.neutral.100` (oscuro) | El color se reserva para significado (P4). |
| Sin truncar | la cifra nunca lleva puntos suspensivos ni se parte en dos líneas; se quita primero lo secundario (U1 §6.1 regla 5) | `$99,999.99` completo en la tarjeta más estrecha a 360 px; lo verifica U5. |
| Sin animación | las cifras no cuentan hacia arriba ni se deslizan (§6.4) | Un monto siempre es exacto en pantalla. |

Otros formatos con la misma tipografía: fechas `10-oct-2026` y `Vie 30 oct 2026` (U1 §0.4); porcentajes de avance enteros `80 %`; tasas como se capturaron, `17.27 %` y `25.8685 %` (K12).

## 4. Espaciado, rejilla, radios, sombras y bordes

### 4.1 Espaciado (retícula de 4 pt)

Siete valores y el cero; ningún margen, relleno o separación usa otro valor (los grosores de borde de §4.5 y las medidas de §4.3 y §4.6 son medidas de componente, no espaciado). En Tailwind 4, `--spacing: 0.25rem` hace que `p-1 … p-12` coincidan con estos tokens.

| Token | px | rem | Uso principal |
|---|---|---|---|
| `space.0` | 0 | 0 | sin separación |
| `space.1` | 4 | 0.25 | icono y texto dentro de un control; entre insignias |
| `space.2` | 8 | 0.5 | entre controles relacionados; etiqueta y campo; separación mínima entre objetivos táctiles |
| `space.3` | 12 | 0.75 | relleno horizontal de botones y campos; entre elementos de una lista |
| `space.4` | 16 | 1 | margen lateral y medianil en móvil; relleno de tarjeta en móvil |
| `space.6` | 24 | 1.5 | relleno de tarjeta desde 768; medianil desde 1024; entre bloques de una página |
| `space.8` | 32 | 2 | margen lateral a 1024; entre secciones de una página |
| `space.12` | 48 | 3 | margen lateral a 1440; aire sobre el título de página |

### 4.2 Rejilla

Columnas de 4, 8 y 12 por rango de ancho (rangos de U1 §6.1; los puntos de quiebre se detallan en la sección 9).

| Rango | Columnas | Margen lateral | Medianil | Área de contenido |
|---|---|---|---|---|
| 360 (hasta 411 px) | 4 | 16 px (`space.4`) | 16 px (`space.4`) | ancho completo |
| 412 (412–767 px) | 4 | 16 px | 16 px | ancho completo |
| 768 (768–1023 px) | 8 | 24 px (`space.6`) | 16 px | ancho completo (la barra lateral se superpone) |
| 1024 (1024–1439 px) | 12 | 32 px (`space.8`) | 24 px (`space.6`) | ancho − barra lateral de 240 px |
| 1440 (desde 1440 px) | 12 | 48 px (`space.12`) | 24 px | ancho − 240 px; páginas de texto con máximo de lectura de 720 px |

Decisiones: 4 columnas en teléfono porque las tarjetas van a una columna y los chips de F-1 en 2 o 3; 8 a 768 para las 2 columnas de tarjetas; 12 en escritorio para 3 tarjetas por fila o tabla con vista previa (U1 §6.2). Las vistas de datos (tabla, tablero, calendario) usan todo el ancho del área; solo las páginas de texto y notas se limitan a 720 px, como el ancho de página de Notion.

### 4.3 Medidas de estructura

| Token | Valor | Uso y justificación |
|---|---|---|
| `layout.sidebar.width` | 240 px | barra lateral fija desde 1024 (U1 G-01): cabe el nombre de cada sección a 14 px con icono de 20 px; colapsada se oculta por completo (Ctrl+\, U1 §6.2) |
| `layout.sidebar.overlay-width` | 280 px | barra lateral superpuesta a 768 (U1 §6.2 deja «ancho fijo (U2)»): un poco más ancha porque se usa con el dedo |
| `layout.sidebar.mobile-width` | 85 % | panel desde ☰ a 360 y 412 (fijado en U1 §6.2) |
| `layout.topbar.height.mobile` · `layout.topbar.height.desktop` | 56 px (< 768) · 48 px (≥ 768) | barra superior de página (U1 G-02): en móvil aloja botones de 44 px con 6 px de aire |
| `layout.bottombar.height` | 56 px + `env(safe-area-inset-bottom)` | barra inferior de 5 destinos (U1 G-05); respeta la barra de gestos de Android |
| `layout.reading.max-width` | 720 px | páginas de texto y notas a partir de 1440 (U1 §6.2) |
| `layout.preview.width.desktop` · `layout.preview.width.wide` | 400 px (1024) · 480 px (1440) | vista previa lateral de registro (U1 G-06, «más ancho» a 1440) |
| `layout.dialog.width.short` · `layout.dialog.width.long` | 480 px · 640 px | diálogo corto (confirmación) y largo (formularios), desde 768; en < 768 son hojas inferiores (U1 §1.3) |
| `layout.command.width` | 640 px | panel centrado del menú `/` y búsqueda desde 768 (U1 G-03, G-04) |

### 4.4 Radios

| Token | px | Uso |
|---|---|---|
| `radius.sm` | 6 | botones, campos, chips, elementos de menú, insignias |
| `radius.md` | 10 | tarjetas, menús emergentes, avisos breves, vista previa |
| `radius.lg` | 14 | diálogos, hojas inferiores (solo esquinas superiores), paneles grandes |
| `radius.full` | 9999 | forma circular (avatares, puntos de categórica, interruptores); no es un paso de la escala |

Regla: un radio interior nunca es mayor que el del contenedor (un botón de `radius.sm` dentro de una tarjeta de `radius.md`). Radios moderados porque Notion es casi recto; 14 px solo en lo que flota sobre todo.

### 4.5 Sombras y bordes

**Claro: dos niveles de sombra leve. Oscuro: bordes en vez de sombras** (plan §4.1), porque una sombra no se ve sobre `color.neutral.950` y ensucia la superficie. Las transparencias se derivan de `color.neutral.950` con `color-mix`; no hay ningún color nuevo.

| Token | Claro | Oscuro | Uso |
|---|---|---|---|
| `shadow.0` | ninguna | ninguna | tarjetas en reposo: planas, se separan por superficie (`color.neutral.100` sobre `color.neutral.50`; `color.neutral.900` sobre `color.neutral.950`) |
| `shadow.1` | `0 1px 2px 0` al 8 % + `0 1px 3px 0` al 6 % de `color.neutral.950` | sin sombra; borde 1 px `color.neutral.800` | menús emergentes, tarjeta mientras se arrastra, vista previa lateral |
| `shadow.2` | `0 4px 12px 0` al 10 % + `0 12px 32px 0` al 8 % de `color.neutral.950` | sin sombra; borde 1 px `color.neutral.800` | diálogos, hojas inferiores, menú `/`, avisos breves |
| `overlay.scrim` | `color.neutral.950` al 40 % | `color.neutral.950` al 60 % | fondo atenuado bajo diálogos, hojas y el panel ☰ |

```css
--shadow-1: 0 1px 2px 0 color-mix(in srgb, var(--color-neutral-950) 8%, transparent),
            0 1px 3px 0 color-mix(in srgb, var(--color-neutral-950) 6%, transparent);
--shadow-2: 0 4px 12px 0 color-mix(in srgb, var(--color-neutral-950) 10%, transparent),
            0 12px 32px 0 color-mix(in srgb, var(--color-neutral-950) 8%, transparent);
```

| Token | Valor | Uso |
|---|---|---|
| `border.width.default` | 1 px | divisores (`color.neutral.200` / `color.neutral.800`, decorativos), bordes de tarjeta y de campo (`color.neutral.600` / `color.neutral.400`, 3:1) |
| `border.width.focus` | 2 px | anillo de foco `color.brand.600` (claro) / `color.brand.400` (oscuro) |
| `focus.offset` | 2 px | separación entre el control y el anillo de foco |

El foco siempre es visible con `:focus-visible` y nunca se quita; con su separación de 2 px se ve sobre cualquier relleno (§2.8).

### 4.6 Objetivos táctiles y alturas de control

| Token | Valor | Cuándo |
|---|---|---|
| `size.touch.min` | **44 px** (44 × 44 px) | área de toque mínima de **todo** control táctil, aunque el dibujo sea menor (un icono de 20 px vive en un área de 44 × 44 px) |
| `size.control.coarse` | 44 px | altura de botones, campos, filas de menú y chips con puntero grueso, `@media (any-pointer: coarse)`: Android y cualquier pantalla táctil |
| `size.control.primary-sheet` | 48 px | botón primario de hojas inferiores en móvil (p. ej., «Guardar» de F-1), al alcance del pulgar |
| `size.control.fine` | 36 px | altura visual en escritorio solo con puntero fino, `@media (pointer: fine) and (not (any-pointer: coarse))`; no es un objetivo táctil |
| `size.bottombar.item.height` · `size.bottombar.item.min-width` | 56 px (alias de `layout.bottombar.height`) · ≥ 64 px | cada destino de la barra inferior: **alto 56 px**, el mismo de la barra sin la zona segura, y **ancho mínimo 64 px** (5 destinos en 360 px = 72 px cada uno) |

Entre dos objetivos táctiles contiguos hay al menos `space.2` (8 px). Los 24 × 24 px de WCAG 2.2 (criterio 2.5.8) son solo el piso que nada puede bajar, nunca el objetivo de diseño de un control táctil (U1 §6.1 regla 4).

## 5. Iconografía

### 5.1 Juego de iconos y archivos

**Lucide 1.55.0 (paquete estático `lucide-static`), licencia ISC, autoalojado** (plan §4.1; decisiones §G). Los iconos derivados de Feather conservan además su licencia MIT, incluida en el mismo archivo. Archivos en `docs/design/fuentes/lucide/`: `sprite.svg` (1,877 símbolos, todos con `viewBox="0 0 24 24"` y sin trazo fijo), `LICENSE` (ISC + MIT de Feather) y la versión, el sha512 del paquete y la procedencia en `docs/design/fuentes/FUENTES.md`. Se eligió Lucide por su trazo lineal uniforme, cercano a los iconos de Notion, y porque ya es el juego que usan los componentes de shadcn/ui de la pila A (decisiones §I).

### 5.2 Tamaños y trazo

| Token | Valor | Uso |
|---|---|---|
| `size.icon.sm` | 20 px | junto a texto de 14 y 16 px: barra lateral, filas de tabla, botones, campos, insignias |
| `size.icon.md` | 24 px | junto a títulos de 20 px o más y solos en barras: barra superior, barra inferior, botones de solo icono |
| `icon.stroke` | **1.5 px** dibujados | trazo constante en los dos tamaños |
| `icon.stroke.user-units.md` · `icon.stroke.user-units.sm` | 1.5 a 24 px · 1.8 a 20 px | valor de `stroke-width` en unidades del `viewBox` de 24 que da 1.5 px dibujados (1.5 × 24 / 20 = 1.8) |

- **Trazo de 1.5 px y no el 2 de Lucide por defecto**: más ligero, acorde al peso 400 de Inter a 14–16 px y a la calma de Notion.
- **Trazo constante**: un icono de 20 px no se ve más delgado que uno de 24 px. En la PWA se logra con el ancho de trazo absoluto del componente de iconos (WF-3); con el sprite, con la regla de abajo.
- **Solo dos tamaños**: un icono de 16 px a 1.5 px pierde detalle; uno de 32 px compite con las cifras.
- Los iconos de solo icono viven en un área táctil de 44 × 44 px (§4.6) y siempre llevan nombre accesible (`aria-label`) y, en escritorio, descripción emergente.

Uso del sprite (documentos, Storybook, páginas estáticas; los símbolos heredan el trazo del `<svg>` que los contiene):

```html
<svg class="icon icon-20" aria-hidden="true"><use href="/icons/sprite.svg#receipt"></use></svg>
```

```css
.icon { fill: none; stroke: currentColor; stroke-linecap: round; stroke-linejoin: round; }
.icon-20 { width: 20px; height: 20px; stroke-width: 1.8; }
.icon-24 { width: 24px; height: 24px; stroke-width: 1.5; }
```

### 5.3 Color de los iconos

El icono toma el color del texto que acompaña (`currentColor`): `color.neutral.700` / `color.neutral.300` en reposo, `color.neutral.900` / `color.neutral.100` en el elemento activo, el color de rol cuando refuerza un rol (§2.4) y `color.brand.700` / `color.brand.300` solo en el icono del elemento activo de la barra lateral (§2.9). Un icono que es el único contenido de un control cumple 3:1 contra su fondo (interfaz); todos los pasos citados ya pasan 4.5:1 como texto sobre fondo y tarjeta (§2.8, §2.9), así que cumplen con holgura. `color.neutral.600` sobre tarjeta (3.46:1) solo se admite en iconos decorativos.

### 5.4 Mapa de iconos

Nombres de Lucide 1.55.0 (`verificar.py` comprueba que cada nombre de la segunda columna existe en el sprite). Los iconos de cada pantalla los completan las secciones 7 y U3/U4; este mapa fija los significados que se repiten.

| Significado | Icono | Dónde |
|---|---|---|
| Ingreso | `arrow-down-left` | cifras con «+» (§3.3) |
| Gasto | `arrow-up-right` | cifras con «−» (§3.3) |
| Ahorro · metas de ahorro | `piggy-bank` | insignias, tarjetas |
| Deuda / crédito | `credit-card` | insignias, tarjetas, Créditos |
| Advertencia | `triangle-alert` | avisos, límite rebasado |
| Error | `circle-alert` | validación de campos, avisos de error |
| Éxito | `circle-check` | «Guardado», «Pagado» |
| Información | `info` | avisos informativos |
| Acción crítica | `shield-alert` | encabezado del diálogo de confirmación |
| Inicio | `house` | barra inferior y lateral |
| Comidas | `utensils` | barra inferior y lateral |
| Registrar | `plus` | «+» de la barra inferior, «+ Registrar» |
| Compras | `shopping-cart` | barra inferior y lateral |
| Buscar | `search` | barra inferior, búsqueda (G-04) |
| Menú · más acciones | `menu` · `ellipsis` | ☰ en móvil · «⋯» del encabezado |
| Colapsar barra lateral | `panel-left` | barra lateral en escritorio |
| Desplegar · navegar | `chevron-right` · `chevron-down` | árbol de la barra lateral, migas |
| Calendario | `calendar` · `calendar-clock` | vista calendario · pagos próximos |
| Ticket con OCR · foto | `scan-line` · `camera` | captura de ticket |
| Adjunto | `paperclip` | comprobantes |
| Filtrar · ordenar | `list-filter` · `arrow-up-down` | barra de la vista |
| Deshacer · cerrar | `undo-2` · `x` | aviso «Guardado · Deshacer», hojas y diálogos |
| Borrar · editar | `trash` · `pencil` | acciones de registro |
| Arrastrar | `grip-vertical` | filas y tarjetas reordenables |
| Cargando | `loader-circle` | indicador de carga |
| Sin conexión | `wifi-off` | franja «Sin conexión con el servidor de casa» (G-07) |
| Vistas de base de datos | `table-2` · `kanban` · `calendar` · `layout-grid` | tabla, tablero, calendario, galería |
| Ajustes | `settings` | entrada a Ajustes |
| Ajustes › Cuenta y roles | `user-cog` | sección 1 de 14 (decisiones §G) |
| Ajustes › Presupuesto | `wallet` | sección 2 |
| Ajustes › Ingresos | `banknote` | sección 3 |
| Ajustes › Categorías | `tags` | sección 4 |
| Ajustes › Créditos y condiciones | `credit-card` | sección 5 |
| Ajustes › Tiendas y precios | `store` | sección 6 |
| Ajustes › Recetas | `chef-hat` | sección 7 |
| Ajustes › Reportes y frecuencias | `file-chart-column` | sección 8 |
| Ajustes › Notificaciones | `bell` | sección 9 |
| Ajustes › Apariencia | `palette` · `sun` · `moon` · `monitor` | sección 10; tema claro, oscuro y del sistema |
| Ajustes › Datos y respaldos | `database` | sección 11 |
| Ajustes › Personas y reparto | `users` | sección 12 |
| Ajustes › Salud y metas de ahorro | `heart-pulse` | sección 13 |
| Ajustes › Inventario | `package` | sección 14 |

## 6. Movimiento

### 6.1 Duraciones y curvas

| Token | Valor | Uso |
|---|---|---|
| `motion.duration.fast` | 150 ms | cambios de estado: hover, presionado, color, casillas, interruptores; salidas de menús y avisos |
| `motion.duration.base` | 200 ms | menús emergentes, aviso breve que entra, indicador de pestaña, desplegar un nodo de la barra lateral |
| `motion.duration.slow` | 250 ms | diálogos, hojas inferiores, panel ☰ y barra lateral superpuesta, vista previa lateral |
| `motion.easing.standard` | `cubic-bezier(0.2, 0, 0, 1)` | cambios dentro de la pantalla y salidas |
| `motion.easing.emphasized` | `cubic-bezier(0.05, 0.7, 0.1, 1)` | entradas de lo que aparece sobre la pantalla (diálogo, hoja, panel, aviso) |

Decisiones:

- **150–250 ms** (plan §4.1): por debajo de 150 ms el cambio no se percibe como transición; por encima de 250 ms la interfaz se siente lenta en el teléfono.
- **Curva enfática solo para entrar**: arranca rápido y frena suave, así lo que aparece llega pronto al lugar donde se lee; salir usa la estándar con `fast` porque lo que se va no debe hacer esperar.
- **Animar solo `opacity` y `transform`**: no cambian el tamaño de la página ni provocan saltos de diseño en teléfonos modestos.
- El aviso «Guardado · Deshacer» permanece 5 s (U1 G-07); la duración de permanencia no es movimiento y no cambia con las preferencias.

### 6.2 Qué se mueve

| Elemento | Entrada | Salida |
|---|---|---|
| Hoja inferior (< 768) | desliza desde abajo, `slow`, enfática | baja, `fast`, estándar |
| Diálogo (≥ 768) | aparece con opacidad y escala 0.98 → 1, `slow`, enfática | opacidad, `fast`, estándar |
| Panel ☰ y barra lateral superpuesta | desliza desde la izquierda, `slow`, enfática, con `overlay.scrim` en opacidad | regresa, `fast` |
| Menú emergente, menú `/` | opacidad y desplazamiento de 4 px, `base`, enfática | opacidad, `fast` |
| Aviso breve | sube 8 px con opacidad, `base`, enfática | opacidad, `fast` |
| Hover, foco, presionado | color de fondo y borde, `fast`, estándar; el anillo de foco aparece sin transición | — |
| Nodo de la barra lateral | el chevron gira 90°, `base`; los hijos aparecen sin animar su altura | ídem |
| Carga | esqueleto estático en `color.neutral.100` / `color.neutral.900`; `loader-circle` gira (1 vuelta por segundo) solo en acciones de más de 400 ms | — |

### 6.3 Movimiento reducido

Con `prefers-reduced-motion: reduce` se desactivan todas las transiciones no esenciales: nada se desliza, escala ni gira por decoración; los cambios de estado y la entrada de diálogos, hojas y paneles son instantáneos. Lo único que se conserva es el indicador de carga `loader-circle`, porque es información (la app está trabajando), y aun así gira a la mitad de velocidad.

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    scroll-behavior: auto !important;
  }
  .icon-spinner { animation: spin 2s cubic-bezier(0, 0, 1, 1) infinite !important; }
}
```

### 6.4 Lo que nunca se anima

- **Cifras de dinero**: sin conteo, sin desplazamiento de dígitos y sin parpadeo al actualizar; el valor nuevo reemplaza al anterior de inmediato.
- Gráficas: aparecen ya dibujadas (sin crecimiento de barras ni trazado de líneas), así la primera lectura ya es la exacta.
- El desplazamiento de la página (`scroll-behavior: auto`).
- Nada parpadea más de 3 veces por segundo (WCAG 2.3.1).

## 7. Componentes

Secciones 7 a 10: agente U2b-2-design (11-oct-2026). Esta sección especifica cada componente de U1 §1.4 y del plan §4.1 punto 7 con su anatomía, variantes, estados, tokens, tamaños, comportamiento en móvil y escritorio y accesibilidad. No define colores nuevos: los tokens `ui.*` son alias de los primitivos de `C2-colores.json` (§7.1) y U2c los lleva a `C2-tokens.json` (§10). Cada par de primer plano y fondo que aparece aquí tiene su fila en la lista de pares de §10.5 (ids `C-xxx`), que U2d calcula.

### 7.0 Convenciones de todos los componentes

| Estado | Cuándo | Regla común |
|---|---|---|
| Reposo | sin interacción | tokens de la tabla del componente |
| Hover | solo con `@media (hover: hover)`; en Android no existe | cambio de fondo o borde con `motion.duration.fast` y `motion.easing.standard` (§6.2) |
| Foco | `:focus-visible` (teclado o lector) | anillo de 2 px (`border.width.focus`) en `ui.border.focus`, separado 2 px (`focus.offset`); en filas y celdas, anillo interior; aparece sin transición y nunca se quita (§4.5) |
| Activo | presionado (`:active`) o seleccionado (`aria-pressed`, `aria-selected`, `aria-current`) | presionado: un paso más de fondo (`*-pressed`); seleccionado: `ui.surface.selected` |
| Deshabilitado | la acción no se puede usar ahora | relleno `ui.surface.disabled`, texto `ui.text.disabled`, borde punteado de 1 px en `ui.border.control`; se prefiere `aria-disabled="true"` con la razón en texto («Requiere conexión con el servidor de casa») a ocultar el control |
| Error | validación o falla al guardar | borde `ui.border.error`, icono `circle-alert` y mensaje en `color.role.error.light` (claro) o `color.role.error.dark` (oscuro) que dice qué pasó y qué hacer; `aria-invalid="true"` y `aria-describedby` al mensaje |

- Cuando un estado no aplica a un componente, su tabla lo dice con el motivo.
- Medidas por tipo de puntero (§4.6): con puntero grueso todo control mide `size.control.coarse` (44 px) y su área táctil `size.touch.min` (44 × 44 px); con puntero fino solo el dibujo baja a `size.control.fine` (36 px).
- Superficies, sombras y radios: lo que flota usa `ui.surface.raised`, borde `ui.border.divider`, `shadow.1` (menús) o `shadow.2` (diálogos, hojas, avisos) en claro y solo el borde en oscuro (§4.5).
- Roles, categóricas y series usan su par `light` / `dark` de C2 directamente (`color.role.<rol>.light` / `.dark`, `color.person.<nombre>.light` / `.dark`, `color.chart.N.light` / `.dark`); en CSS cada par es una sola variable que cambia con el tema (§10.3).

### 7.1 Tokens semánticos de los componentes

41 tokens, cada uno alias de un primitivo de C2 en cada tema (ningún valor nuevo). Los que repiten el valor de otro llevan su alias declarado para que U2c los relacione en `C2-tokens.json` (criterio C2-10). Los componentes usan estos nombres y no combinan primitivos a mano.

| Token semántico | Claro | Oscuro | Alias de | Uso |
|---|---|---|---|---|
| `ui.surface.page` | `color.neutral.50` · `#FAFAFB` | `color.neutral.950` · `#222325` | — | fondo de página, barra superior, tablas y tableros |
| `ui.surface.card` | `color.neutral.100` · `#F2F2F4` | `color.neutral.900` · `#36363A` | — | tarjetas, barra lateral, fondo de gráfica |
| `ui.surface.raised` | `color.neutral.50` · `#FAFAFB` | `color.neutral.900` · `#36363A` | — | lo que flota: menús, menú «/», diálogos, hojas, vista previa, avisos breves; fondo de chips |
| `ui.surface.field` | `color.neutral.50` · `#FAFAFB` | `color.neutral.950` · `#222325` | — | interior de campos |
| `ui.surface.hover` | `color.neutral.200` · `#E5E6E9` | `color.neutral.800` · `#4D4E53` | — | hover de filas, elementos de menú, botones secundario, de texto y de icono |
| `ui.surface.hover-subtle` | `color.neutral.100` · `#F2F2F4` | `color.neutral.900` · `#36363A` | `ui.surface.card` | hover de filas de tabla y de la edición en el lugar sobre la página (las cifras de rol conservan 4.5:1) |
| `ui.surface.pressed` | `color.neutral.300` · `#D1D2D8` | `color.neutral.700` · `#66676E` | — | presionado de secundario, de texto y de icono |
| `ui.surface.selected` | `color.brand.100` · `#F0F1FF` | `color.brand.900` · `#303258` | — | elemento activo, fila seleccionada, chip de filtro activo |
| `ui.surface.disabled` | `color.neutral.100` · `#F2F2F4` | `color.neutral.900` · `#36363A` | — | relleno de control deshabilitado |
| `ui.skeleton.on-page` | `color.neutral.100` · `#F2F2F4` | `color.neutral.900` · `#36363A` | — | esqueleto de carga sobre la página (decorativo) |
| `ui.skeleton.on-card` | `color.neutral.200` · `#E5E6E9` | `color.neutral.800` · `#4D4E53` | — | esqueleto de carga sobre tarjeta (decorativo) |
| `ui.text.default` | `color.neutral.900` · `#36363A` | `color.neutral.100` · `#F2F2F4` | — | texto normal, cifras neutras |
| `ui.text.secondary` | `color.neutral.700` · `#66676E` | `color.neutral.300` · `#D1D2D8` | — | texto secundario, marcador de posición, etiquetas de eje |
| `ui.text.disabled` | `color.neutral.700` · `#66676E` | `color.neutral.400` · `#B7B8C0` | — | texto de control deshabilitado (con borde punteado) |
| `ui.text.link` | `color.brand.700` · `#5D609F` | `color.brand.300` · `#C9CFFF` | — | enlace (siempre subrayado) |
| `ui.text.link-hover` | `color.brand.800` · `#46487B` | `color.brand.200` · `#E1E4FF` | — | enlace en hover |
| `ui.text.on-primary` | `color.base.white` · `#FFFFFF` | `color.brand.950` · `#1E203A` | — | texto e icono sobre el botón primario |
| `ui.text.on-destructive` | `color.base.white` · `#FFFFFF` | `color.role.error.950` · `#3C1416` | — | texto sobre el botón destructivo |
| `ui.text.on-critical` | `color.base.white` · `#FFFFFF` | `color.critical.950` · `#292300` | — | texto sobre el botón de acción crítica |
| `ui.icon.default` | `color.neutral.700` · `#66676E` | `color.neutral.300` · `#D1D2D8` | `ui.text.secondary` | icono en reposo |
| `ui.icon.selected` | `color.brand.700` · `#5D609F` | `color.brand.300` · `#C9CFFF` | `ui.action.primary.bg` | icono del elemento activo (barra lateral, barra inferior) |
| `ui.border.divider` | `color.neutral.200` · `#E5E6E9` | `color.neutral.800` · `#4D4E53` | — | divisores y borde de lo que flota (decorativo) |
| `ui.border.control` | `color.neutral.600` · `#808189` | `color.neutral.400` · `#B7B8C0` | — | borde de campo, casilla, botón secundario (≥ 3:1) |
| `ui.border.control-hover` | `color.neutral.700` · `#66676E` | `color.neutral.300` · `#D1D2D8` | `ui.text.secondary` | borde de campo y botón secundario en hover |
| `ui.border.focus` | `color.brand.600` · `#7579C2` | `color.brand.400` · `#ADB3F4` | — | anillo de foco de 2 px con separación de 2 px |
| `ui.border.selected` | `color.brand.700` · `#5D609F` | `color.brand.300` · `#C9CFFF` | `ui.action.primary.bg` | indicador de pestaña activa, destino de arrastre, tarjeta elegida |
| `ui.border.error` | `color.role.error.light` · `#A34449` | `color.role.error.dark` · `#FC9999` | — | borde de campo o celda con error |
| `ui.action.primary.bg` | `color.brand.700` · `#5D609F` | `color.brand.300` · `#C9CFFF` | — | relleno del botón primario, casilla marcada, interruptor encendido, «hoy» |
| `ui.action.primary.bg-hover` | `color.brand.800` · `#46487B` | `color.brand.200` · `#E1E4FF` | — | botón primario en hover |
| `ui.action.primary.bg-pressed` | `color.brand.900` · `#303258` | `color.brand.100` · `#F0F1FF` | — | botón primario presionado |
| `ui.action.destructive.bg` | `color.role.error.700` · `#A34449` | `color.role.error.300` · `#FFC0BF` | — | botón destructivo |
| `ui.action.destructive.bg-hover` | `color.role.error.800` · `#7E3236` | `color.role.error.200` · `#FFDDDC` | — | botón destructivo en hover |
| `ui.action.destructive.bg-pressed` | `color.role.error.900` · `#5A2124` | `color.role.error.100` · `#FFEEED` | — | botón destructivo presionado |
| `ui.action.critical.bg` | `color.critical.700` · `#766800` | `color.critical.300` · `#E0D58E` | — | botón de acción crítica (solo en su diálogo) |
| `ui.action.critical.bg-hover` | `color.critical.800` · `#5A4F00` | `color.critical.200` · `#EEE8BD` | — | acción crítica en hover |
| `ui.action.critical.bg-pressed` | `color.critical.900` · `#3F3700` | `color.critical.100` · `#F8F4D9` | — | acción crítica presionada |
| `ui.viz.grid` | `color.neutral.200` · `#E5E6E9` | `color.neutral.800` · `#4D4E53` | `ui.border.divider` | retícula de gráficas, 1 px (decorativa) |
| `ui.viz.baseline` | `color.neutral.600` · `#808189` | `color.neutral.400` · `#B7B8C0` | `ui.border.control` | línea base (cero) de gráficas |
| `ui.viz.other` | `color.neutral.600` · `#808189` | `color.neutral.400` · `#B7B8C0` | `ui.border.control` | serie «Otros» (novena en adelante) |
| `ui.viz.track` | `color.neutral.200` · `#E5E6E9` | `color.neutral.800` · `#4D4E53` | `ui.border.divider` | pista sin llenar de barras de avance (decorativa) |
| `ui.viz.meter` | `color.neutral.700` · `#66676E` | `color.neutral.300` · `#D1D2D8` | `ui.text.secondary` | llenado de barra de avance bajo el primer umbral |

### 7.2 Botones y enlaces

**Anatomía.** Contenedor con `radius.sm`; icono opcional a la izquierda (`size.icon.sm`, a `space.1` de la etiqueta); etiqueta en infinitivo (§1.3) con `font.weight.medium`, `font.size.small` en escritorio y `font.size.body` en móvil (§3.2). El atajo, si existe, va en la descripción emergente, nunca dentro del botón.

**Variantes.**

| Variante | Para qué | Regla |
|---|---|---|
| Primario | la acción principal de la pantalla, hoja o diálogo | uno solo visible (P2); nunca blanco sobre marca 600 (U2-07) |
| Secundario | acciones de apoyo junto a la primaria | borde de 1 px |
| De texto | acciones terciarias, «Cancelar», «Mostrar N más» | sin borde ni relleno en reposo |
| De icono | «⋯», cerrar, colapsar, filtrar | solo icono, con `aria-label` y descripción emergente |
| Destructivo | borrar un registro o un renglón | solo en el diálogo que confirma o en el menú «⋯»; verbo explícito («Borrar gasto») |
| Acción crítica | lo irreversible: restaurar un respaldo, borrar datos | solo dentro de su diálogo, una vez (§2.5); verbo con cantidad («Borrar 12 gastos») |
| Enlace | navegar a otra página | siempre subrayado; `ui.text.link`, hover `ui.text.link-hover` |

**Estados** (los tokens valen para claro y oscuro; su valor por tema está en §7.1):

| Estado | Primario | Secundario | De texto y de icono | Destructivo | Acción crítica |
|---|---|---|---|---|---|
| Reposo | relleno `ui.action.primary.bg`, texto `ui.text.on-primary` | borde `ui.border.control`, texto `ui.text.default`, sin relleno | texto `ui.text.secondary`, icono `ui.icon.default` | relleno `ui.action.destructive.bg`, texto `ui.text.on-destructive` | relleno `ui.action.critical.bg`, texto `ui.text.on-critical` |
| Hover | `ui.action.primary.bg-hover` | relleno `ui.surface.hover`, borde `ui.border.control-hover` | relleno `ui.surface.hover`, texto `ui.text.default` | `ui.action.destructive.bg-hover` | `ui.action.critical.bg-hover` |
| Foco | anillo común (§7.0) | ídem | ídem | ídem | ídem |
| Activo (presionado) | `ui.action.primary.bg-pressed` | `ui.surface.pressed` | `ui.surface.pressed`, texto `ui.text.default` | `ui.action.destructive.bg-pressed` | `ui.action.critical.bg-pressed` |
| Activo (alternado, `aria-pressed`) | no aplica | relleno `ui.surface.selected`, borde `ui.border.selected` | relleno `ui.surface.selected`, icono `ui.icon.selected` | no aplica | no aplica |
| Deshabilitado | común (§7.0): pierde el relleno de marca | común | texto `ui.text.disabled` | común | común |
| Error | no aplica: el error vive en el campo o en un aviso (§7.14) | — | — | — | — |
| Ocupado | `loader-circle` en lugar del icono, la etiqueta se queda, el ancho no cambia, `aria-busy="true"`; no admite un segundo toque (evita duplicados, K9) | ídem | ídem | ídem | ídem |

**Tamaños.** Alto `size.control.coarse` (44 px) con puntero grueso y `size.control.fine` (36 px) con puntero fino; `size.control.primary-sheet` (48 px) para el primario de una hoja inferior. Relleno horizontal `space.3`; `space.4` en el de hoja. De icono: cuadrado de 44 px (grueso) o 36 px (fino) con icono `size.icon.sm`, o `size.icon.md` en barras. Entre botones contiguos, `space.2`.

**Móvil y escritorio.** Móvil (< 768): la primaria siempre visible (U1 §1.3); en hojas, a todo el ancho y pegada abajo, sobre el teclado (U1 §6.1 regla 6); secundarias en «⋯». Escritorio: primaria a la derecha de la barra superior (G-02) y secundarias a su izquierda.

**Accesibilidad.** `<button>` nativo (o `<a>` para el enlace); nombre visible o `aria-label`; el destructivo y el crítico no dependen del color: verbo, cantidad e icono (`trash`, `shield-alert` en el encabezado del diálogo). Pares: C-001 a C-033.

### 7.3 Campos de formulario y edición en el lugar

**Anatomía.** Etiqueta arriba (`font.size.small`, `font.weight.medium`, `ui.text.default`) a `space.2` del control; control con `radius.sm`, borde `border.width.default` en `ui.border.control`, fondo `ui.surface.field` y valor en `font.size.body` (16 px también en móvil, §3.2); prefijo o sufijo (`$`, `%`, `kg`, `kcal`) en `ui.text.secondary`; ayuda o error debajo en `font.size.caption`. Entre campos, `space.4`.

**Variantes.**

| Variante | Regla |
|---|---|
| Texto, área de texto | notas largas con el editor de bloques (U1 §5.1 punto 2) |
| Dinero | prefijo `$`, `inputmode="decimal"`, `tabular-nums`, dos decimales al salir del campo; a la derecha en tablas |
| Monto grande | primer campo de «Registrar gasto» (P-03): `font.size.title-2`, `font.weight.semibold`, enfocado al abrir |
| Porcentaje y tasa | se muestra con los decimales capturados (`17.27 %`, `25.8685 %`; K12) |
| Número con unidad | entero o decimal con sufijo; flechas solo en escritorio |
| Fecha | muestra `10-oct-2026`; el selector empieza en lunes (K8) y marca «hoy» como el calendario (§7.5.4) |
| Hora | `HH:MM` de 24 h, hora de Ciudad de México |
| Selección y combobox | lista corta: menú emergente; lista larga (categoría, producto, tienda): combobox con búsqueda sin acentos |
| Chips de opción rápida | categorías frecuentes de F-1 en 2 o 3 columnas, cada chip de 44 px (§7.4) |
| Casilla | cuadro de 20 px dentro de un área de 44 px; marcada: relleno `ui.action.primary.bg` con `check` en `ui.text.on-primary`; sin marcar: borde `ui.border.control` |
| Interruptor | pista de 36 × 20 px dentro de 44 px; encendido: pista `ui.action.primary.bg`, perilla `ui.text.on-primary`; apagado: pista `ui.border.control`, perilla `ui.surface.page`; la posición y el texto de la fila dicen el estado, no solo el color |
| Opción segmentada | 2 a 4 opciones excluyentes (tema Sistema · Claro · Oscuro; alcance «solo este mes» · «de este mes en adelante», K2): la elegida en `ui.surface.selected` con borde `ui.border.selected` |
| Calculado (solo lectura) | sin borde ni fondo de campo, valor en `ui.text.default`; descripción emergente «Se calcula» (U1 §2.6) |

**Estados.**

| Estado | Tokens |
|---|---|
| Reposo | borde `ui.border.control`, fondo `ui.surface.field`, valor `ui.text.default`, marcador `ui.text.secondary` |
| Hover | borde `ui.border.control-hover` |
| Foco (escribiendo) | anillo común (§7.0) fuera del borde; el borde se conserva |
| Activo | casilla, interruptor y segmentada: ver variantes; un campo de texto activo es el foco |
| Deshabilitado | común (§7.0) |
| Error | borde `ui.border.error`, icono `circle-alert` y mensaje en `color.role.error.light` / `color.role.error.dark` debajo; el valor capturado no se borra (U1 §2.4) |

**Edición en el lugar** (U1 §2.6, como en Notion): el valor se ve como texto; hover: fondo `ui.surface.hover-subtle` sobre la página o solo borde `ui.border.control-hover` sobre una tarjeta (así las cifras de rol conservan 4.5:1; pares C-046 a C-049); tocar o Enter lo vuelve campo con foco; Enter o salir guarda solo y muestra «Guardado · Deshacer»; Esc revierte; si falla, regresa al valor anterior con el aviso de error. Lo calculado no muestra hover.

**Tamaños.** Alto 44 px (grueso) o 36 px (fino); relleno `space.3`; menú emergente de selección con `radius.md` y `shadow.1`.

**Móvil y escritorio.** Móvil: los formularios viven en hojas; `inputmode` abre el teclado correcto; selección y combobox abren una hoja con búsqueda. Escritorio: menú emergente anclado al campo; `Ctrl+Enter` guarda la hoja (U1 §5.3).

**Accesibilidad.** `<label for>` siempre (el marcador no sustituye la etiqueta); `aria-describedby` a ayuda y error; `aria-required`; el combobox sigue el patrón ARIA *combobox* con *listbox*; el error se anuncia sin mover el foco. Pares: C-034 a C-049.

### 7.4 Chips de persona y tienda, avatares, chips de filtro e insignias

**Regla K41 aplicada a la forma.** Persona y tienda se distinguen siempre por **forma, nombre y posición**, no solo por color:

| Grupo | Marca de forma (10 px) | Contenedor del chip | Orden fijo |
|---|---|---|---|
| Persona | círculo relleno (`radius.full`); Conjunto: anillo hueco con trazo de 2 px | píldora (`radius.full`) | Axel, Ximena, Karen, Conjunto |
| Tienda | cuadrado relleno de esquinas rectas; los mercados adicionales usan el de Mercado | rectángulo con `radius.sm` | Walmart, Bodega Aurrera, Sam's Club, Mercado |

**Anatomía del chip.** Alto visual 24 px; fondo `ui.surface.raised` (así la marca se mide siempre sobre un fondo verificado, también en filas seleccionadas); borde `ui.border.divider` (decorativo); marca de forma en `color.person.<nombre>.light` / `.dark` o `color.store.<nombre>.light` / `.dark`, a `space.1` del nombre; nombre en `font.size.small` y `ui.text.default`. Nunca hay texto encima del color (§2.6 regla 1).

**Avatar.** Círculo de 24 px con fondo `ui.surface.hover`, inicial en `font.size.caption` y `font.weight.medium` (`ui.text.default`) y anillo de 2 px del color de la persona; Conjunto muestra el icono `users` en lugar de inicial. Va en el menú de usuario de la barra lateral y en la columna «Persona».

**Variantes.** Estático (celdas, tarjetas); seleccionable (filtros y chips de F-1, `aria-pressed`); quitable (filtro activo con `x`); insignia de rol (texto `color.role.<rol>.light` sobre `color.role.<rol>.surface-light`, en oscuro `.dark` sobre `.surface-dark`, con el icono del rol de §2.4, `radius.sm`, `font.size.caption`, alto 20 px).

| Estado | Chip estático | Chip seleccionable o de filtro |
|---|---|---|
| Reposo | tokens de la anatomía | ídem |
| Hover | no aplica (no es interactivo) | fondo `ui.surface.hover` |
| Foco | no aplica | anillo común |
| Activo | no aplica | seleccionado: fondo `ui.surface.selected`, borde `ui.border.selected`, icono `check`; filtro activo con `x` en `ui.icon.default` |
| Deshabilitado | no aplica | común (§7.0) |
| Error | no aplica | no aplica |

**Tamaños.** Chip de 24 px visuales; el seleccionable mide 44 px de alto con puntero grueso (área táctil). Separación entre chips `space.2`.

**Móvil y escritorio.** Móvil: los chips de filtro activos van bajo el título (U1 §1.3) con desplazamiento horizontal propio. Escritorio: en la barra de la vista.

**Accesibilidad.** La marca de forma lleva `aria-hidden="true"`: el nombre en texto es el portador; seleccionable = `<button aria-pressed>`; quitable con `aria-label` («Quitar filtro Walmart»). En una vista que mezcla personas y tiendas, van en columnas distintas y nunca comparten leyenda ni gráfica (§2.6). Pares: C-050 a C-074.

### 7.5 Vistas de base de datos

Las cuatro vistas de Notion (decisiones §G) sobre las entidades de U1 §2.5; cada entidad ofrece solo las vistas que marca U1 §2.5.a, y la vista, filtros, orden, agrupación y columnas se recuerdan por dispositivo.

#### 7.5.1 Barra de la vista

Selector de vista (pestañas con icono: `table-2` tabla, `kanban` tablero, `calendar` calendario, `layout-grid` galería), «Filtrar» (`list-filter`), «Ordenar» (`arrow-up-down`), «Agrupar» y búsqueda por texto donde U1 §2.5.b lo indica. Son botones de texto (§7.2); la vista elegida va como alternado (`aria-pressed`). Móvil: un solo botón «Filtrar» abre una hoja con filtros, orden y agrupación, y los filtros activos quedan como chips quitables bajo el título (§7.4). Escritorio: una fila sobre la vista, a todo el ancho del área.

#### 7.5.2 Tabla

**Anatomía.** Encabezado fijo al desplazar (icono del tipo de propiedad en `size.icon.sm` y nombre en `font.size.small`, `font.weight.medium`, `ui.text.secondary`); filas a `font.size.small` con divisor `ui.border.divider`; primera columna (título del registro) fija al desplazar en horizontal; columna de selección con casilla y `grip-vertical` para reordenar donde el orden es manual (categorías); dinero a la derecha con `tabular-nums` (§3.3); encabezado de grupo plegable con conteo y suma; pie con totales y «Cargar 200 más» (200 registros por página, K6).

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | fondo `ui.surface.page`, texto `ui.text.default`, cifras de rol con signo e icono (`color.role.<rol>.light` / `.dark`) |
| Hover | fila en `ui.surface.hover-subtle`; aparecen casilla y `grip-vertical`; la celda editable muestra el cursor de texto |
| Foco | fila: anillo interior de 2 px; celda en edición: el campo con su anillo (§7.3) |
| Activo | fila seleccionada (`aria-selected="true"`): `ui.surface.selected`; las cifras de rol conservan 4.5:1 (C-083, C-084) |
| Deshabilitado | no hay filas deshabilitadas: lo calculado simplemente no se edita (U1 §2.6); un gasto en la cola del teléfono lleva `clock` y «En cola» |
| Error | celda: borde `ui.border.error`, regresa al valor anterior y avisa (U1 §2.4); fila que no se pudo enviar: `circle-alert` y «No se pudo enviar · Revisar» |

**Tamaños.** Alto de fila 36 px con puntero fino y 44 px con grueso; relleno de celda `space.2` horizontal; columna de título de 240 px mínimo en escritorio.

**Móvil y escritorio** (U1 §6.2). 360: filas compactas de 2 líneas (título y cifra; una propiedad debajo), el resto en la vista previa; 412: una propiedad más; 768: hasta 5 columnas; 1024: completa con columnas configurables; 1440: completa. Tocar una fila en móvil abre la vista previa a página completa (G-06).

**Accesibilidad.** `<table>` con `<th scope="col">` y `aria-sort`; la tabla editable sigue el patrón ARIA *grid*: ↑ ↓ entre filas, ← → entre celdas, Enter edita, Esc sale, `Ctrl+/` edita las filas seleccionadas (U1 §5.3); «Guardado» se anuncia con `aria-live="polite"`. Pares: C-075 a C-087.

#### 7.5.3 Tablero

**Anatomía.** Columnas por grupo (U1 §2.5.a) de 272 px con encabezado (nombre del grupo, conteo y suma de dinero en `tabular-nums`); tarjetas con fondo `ui.surface.card`, `radius.md`, relleno `space.3`, título en `font.size.small` y `font.weight.medium` (`ui.text.default`), hasta 3 propiedades en `ui.text.secondary` y la cifra a la derecha; «+ Nuevo» al pie de cada columna.

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | tarjeta plana sobre `ui.surface.page` (`shadow.0`) |
| Hover | borde de 1 px `ui.border.control-hover` (sin cambiar el fondo, para no bajar el contraste de las cifras de rol); aparece `grip-vertical` |
| Foco | anillo común en la tarjeta |
| Activo | seleccionada: `ui.surface.selected`; arrastrando: `shadow.1` en claro o borde `ui.border.divider` en oscuro; columna destino con borde punteado de 2 px `ui.border.selected` y fondo `ui.surface.selected` |
| Deshabilitado | tarjetas de solo lectura (P-17, P-18) sin `grip-vertical` ni arrastre |
| Error | si el cambio no se guarda, la tarjeta vuelve a su columna y sale el aviso de error |

**Móvil y escritorio** (U1 §6.2). 360 y 412: una columna visible, se cambia deslizando o con el selector de grupo; 768: 2; 1024: 3 a 4; 1440: todas las que quepan. El tablero es uno de los tres componentes con desplazamiento horizontal propio (U1 §6.1 regla 3).

**Accesibilidad.** Arrastrar nunca es la única forma (WCAG 2.5.7): el menú «⋯» de cada tarjeta ofrece «Mover a…» con la lista de grupos; con teclado, Espacio levanta, ← → cambian de columna, Espacio suelta y Esc cancela, con anuncio en `aria-live`. Pares: C-088 a C-090.

#### 7.5.4 Calendario

**Anatomía.** Mes en rejilla de 7 columnas **de lunes a domingo** (K8), encabezados «Lun … Dom» en `font.size.caption` y `ui.text.secondary`; en cada celda, número del día (`ui.text.secondary`; días de otro mes en `ui.text.disabled`) y registros como renglones de una línea (título y cifra) con fondo `ui.surface.card` y `radius.sm`; «+N más» como botón de texto; navegación «‹ Hoy ›» (← → con teclado, U1 §5.3). **Hoy**: el número va en un círculo `ui.action.primary.bg` con `ui.text.on-primary` (C-091).

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | celdas sobre `ui.surface.page` con divisores `ui.border.divider` |
| Hover | celda en `ui.surface.hover-subtle`; aparece «+» para crear en esa fecha |
| Foco | anillo interior en la celda o en el renglón |
| Activo | día seleccionado en `ui.surface.selected` |
| Deshabilitado | no aplica (los días fuera de rango solo cambian de color de número, con texto) |
| Error | aviso en la vista con Reintentar (§7.14) |

**Móvil y escritorio** (U1 §6.2). 360: agenda por día (lista agrupada por fecha con encabezado fijo `Vie 30 oct 2026`); 412: agenda y mes compacto opcional; 768 y 1024: mes; 1440: mes con el detalle del día seleccionado a la derecha.

**Accesibilidad.** Patrón ARIA *grid*; cada día con nombre completo («viernes 30 de octubre de 2026, 2 pagos») y `aria-current="date"` en hoy. Pares: C-091 a C-093.

#### 7.5.5 Galería

**Anatomía.** Tarjetas en rejilla con portada 4:3 (foto del ticket; sin imagen, el icono de la entidad en `ui.icon.default` sobre `ui.surface.hover`), título, 2 o 3 propiedades y la cifra; fondo `ui.surface.card`, `radius.md`, relleno `space.3`. Estados: reposo plano (`shadow.0`); hover con borde `ui.border.control-hover`; foco con anillo común; activo (seleccionada) con `ui.surface.selected` y borde `ui.border.selected`; deshabilitado no aplica; error de imagen: icono de la entidad y texto «Foto no disponible». Columnas: 1 a 360 (tickets en 2), 2 a 412, 3 a 768, 3 a 4 a 1024, 4 a 5 a 1440 (U1 §6.2). Imágenes con carga diferida y texto alternativo con comercio, fecha y total («Ticket de Walmart, 10-oct-2026, $1,258.68»). Par: C-094.

### 7.6 Tarjetas de métrica

**Anatomía.** Etiqueta (`font.size.small`, `ui.text.secondary`, sin dos puntos); cifra principal (`font.size.title-2`, `font.weight.semibold`, `tabular-nums`; neutra en `ui.text.default` o de rol con signo e icono, §3.3); dato secundario opcional (uno a 360, uno más a 412; comparación contra un periodo con nombre, «12 % más que septiembre», con signo, icono y color de rol); y opcionalmente una **barra de avance** o una mini serie (12 puntos, serie 1, sin ejes, desde 768). El Resumen lleva ≤ 6 tarjetas en orden fijo (decisiones §G).

**Barra de avance** (presupuesto, metas): pista `ui.viz.track` de 8 px con `radius.full`; llenado `ui.viz.meter` bajo el primer umbral; desde el primer umbral de alerta (`alerts.*` de `app_setting`, 80 % inicial y configurable, K44) el llenado pasa a `color.role.warning.light` / `.dark` con `triangle-alert` y el texto «80 % usado»; desde el segundo (100 % inicial) se agrega «rebasado» y el excedente `−$230.00` (§3.3). El llenado nunca pasa del ancho de la pista: el excedente se dice en texto. Los umbrales se leen de Ajustes, nunca se fijan en el componente.

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | fondo `ui.surface.card`, `shadow.0`, `radius.md` |
| Hover | borde de 1 px `ui.border.control-hover` (la tarjeta entera es enlace al Detalle) |
| Foco | anillo común en la tarjeta |
| Activo | presionado: borde `ui.border.selected`, sin cambiar el fondo |
| Deshabilitado | no aplica |
| Error | «No se pudo cargar · Reintentar» dentro de la tarjeta, sin cifra (U1 §2.4) |
| Carga · vacío · falta capturar | esqueleto de etiqueta y cifra en `ui.skeleton.on-card`; «Sin datos este mes» con su acción; «Falta capturar · Capturar» (K18) |

**Tamaños.** Relleno `space.4` en móvil y `space.6` desde 768; todas las tarjetas de una fila con el mismo alto; `$99,999.99` cabe completo a 360 px (P5). Columnas: 1 (360, 412), 2 (768), 3 (1024, 1440).

**Accesibilidad.** La tarjeta es un solo enlace con nombre completo («Gasto del mes, $8,240.10, 12 % más que septiembre»); la comparación usa palabras («más», «menos») además de flecha y color; la barra de avance es `role="meter"` con `aria-valuetext` («82 % del plan»). Pares: C-095 a C-105.

### 7.7 Gráficas

**Anatomía.** Título con el estilo de la etiqueta de tarjeta y periodo en `ui.text.secondary`; resumen de una frase debajo del título («El gasto subió 12 % respecto a septiembre»); área de trazado con retícula horizontal de 1 px sólida en `ui.viz.grid` (decorativa) y línea base en `ui.viz.baseline`; ejes en `font.size.caption`, `ui.text.secondary` y `tabular-nums` (dinero del eje compacto: `$0`, `$5 mil`, `$10 mil`); marcas según §8; leyenda (con 2 o más series) con marca de forma y nombre en `ui.text.default`; etiquetas directas en `ui.text.default`; descripción emergente con fondo `ui.surface.raised`, `radius.md` y `shadow.1` (fecha o categoría y una fila por serie: marca, nombre y valor exacto); botón de texto «Ver como tabla» (`table-2`) que alterna una tabla con los mismos datos.

**Variantes.** Barras horizontales, columnas, líneas, área apilada, dona y mini serie (reglas en §8); barra de avance (§7.6).

| Estado | Comportamiento |
|---|---|
| Reposo | dibujada completa, sin animación de entrada (§6.4) |
| Hover | puntero fino: guía vertical de 1 px en `ui.viz.baseline` y descripción emergente del punto o barra más cercano; las demás series no se atenúan |
| Foco | el contenedor recibe el foco; ← → recorren los puntos, ↑ ↓ cambian de serie; el punto enfocado lleva anillo `ui.border.focus` y muestra su descripción |
| Activo | tocar en móvil fija la descripción; tocar fuera la cierra |
| Deshabilitado | no aplica |
| Error | «No se pudo cargar la gráfica · Reintentar» en el alto reservado |
| Carga · vacío | esqueleto del alto final en `ui.skeleton.on-card`, sin ejes falsos; «Aún no hay datos para este periodo» con su acción |

**Tamaños.** Alto 200 px en < 768 y 280 px desde 768; marcas según §8.1; área de toque de cada punto o barra ≥ 44 px aunque la marca sea menor.

**Móvil y escritorio** (U1 §6.2). < 768: máximo 6 series, leyenda debajo, valor al tocar. ≥ 768: hasta 8 series, leyenda a la derecha, valor al pasar el puntero.

**Accesibilidad.** `<figure>` con `<figcaption>` (título y resumen); cada punto enfocable con nombre («Octubre, $8,240.10»); la tabla alterna siempre disponible; con `forced-colors: active` se conservan forma, trazo y patrón. Pares: C-106 a C-120.

### 7.8 Calendario de pagos

Subsección Calendario de Créditos (P-17, P-18; `pago_programado`). Parte del calendario de §7.5.4 y le agrega el estado de cada pago.

**Anatomía.** Arriba, resumen del mes en cifra neutra («Este mes pagas $3,512.40 en 4 pagos»). Agenda (< 768): una fila de 2 líneas por pago con icono de estado, nombre de la línea, «Pago 3 de 6», monto a la derecha en `tabular-nums` y, si hay parte de otra persona, su chip con el monto. Mes (≥ 768): hasta 2 pagos por día y «+N más», con el total del día en `font.size.caption`. Tocar un pago abre su vista previa (G-06) con «Marcar pago» como primaria (N-210; nunca en el lugar, U1 §2.6).

**Estado de cada pago** (siempre icono y texto):

| Estado | Cuándo | Icono | Color | Texto |
|---|---|---|---|---|
| Pendiente | aún no vence | `circle-dashed` | icono `ui.icon.default`, monto `ui.text.default` | «Pendiente» |
| Por vencer | dentro de los días de aviso (`alerts.*`: 1 y 0 días antes, configurables, K44) | `calendar-clock` | `color.role.warning.light` / `.dark` | «Vence mañana», «Vence hoy» |
| Vencido | la fecha pasó sin pago completo | `calendar-x` | `color.role.warning.light` / `.dark` | «Vencido» |
| Parcial | alguna parte pagada y otra no (el capital se repone al completarse, K15) | `circle-dot` | `ui.icon.default` | «Falta la parte de Karen» |
| Pagado | todas las partes pagadas | `circle-check` | `color.role.success.light` / `.dark` | «Pagado» |

| Estado del componente | Comportamiento |
|---|---|
| Reposo · hover | celda o fila sobre `ui.surface.page`; hover en `ui.surface.hover-subtle` |
| Foco | anillo interior en la celda o fila |
| Activo | día seleccionado en `ui.surface.selected` (C-124) |
| Deshabilitado | sin conexión, «Marcar pago» dice «Requiere conexión con el servidor de casa» (U1 §2.4) |
| Error | si marcar el pago falla, el estado no cambia y sale el aviso; tarjeta o crédito de tienda sin sus datos: «Falta capturar · Capturar» en lugar de pago mínimo (K18) |

**Tamaños.** Fila de agenda de 56 px; celda de mes de 96 px de alto a 1024 y 120 px a 1440. **Accesibilidad.** Cada pago se lee completo («jueves 15 de octubre, Kueski, pago 3 de 6, $1,258.68, pendiente»); vencido y por vencer no dependen del color. Pares: C-121 a C-124.

### 7.9 Barra lateral colapsable (G-01)

**Anatomía** (U1 §1.2). Fondo `ui.surface.card`, ancho `layout.sidebar.width`; encabezado con «Finanzas Hogar» y el avatar del usuario (menú: Cuenta, Apariencia, Cerrar sesión); dos acciones fijas («Buscar», «Registrar gasto»); las 9 secciones en orden fijo con icono `size.icon.sm`, nombre en `font.size.small` y `chevron-right` (gira 90° a `chevron-down` al desplegar, §6.2); subpáginas con sangría de `space.4`; pie con el interruptor de tema; botón `panel-left` («‹‹») para colapsar.

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | elemento en `ui.text.secondary`, icono `ui.icon.default` |
| Hover | fondo `ui.surface.hover`, texto `ui.text.default` |
| Foco | anillo interior en el elemento |
| Activo | página abierta (`aria-current="page"`): fondo `ui.surface.selected`, texto `ui.text.default` con `font.weight.medium`, icono `ui.icon.selected` |
| Deshabilitado | no aplica: las secciones no se ocultan ni se deshabilitan (U1 §1.2 punto 2) |
| Error | no aplica: la falta de conexión la dice la franja global (§7.14) |

**Tamaños.** Elemento de 32 px de alto con puntero fino y 44 px con grueso; relleno horizontal `space.2`; ancho fijo 240 px (`layout.sidebar.width`), superpuesta 280 px (`layout.sidebar.overlay-width`), panel móvil 85 % (`layout.sidebar.mobile-width`).

**Móvil y escritorio** (U1 §1.3, §6.2). < 768: oculta; se abre desde ☰ (`menu`) como panel con `overlay.scrim`, nunca deslizando desde el borde (gesto «atrás» de Android); 768: superpuesta desde ☰; ≥ 1024: fija y colapsable con «‹‹», `Ctrl+\` o /barra; colapsada reaparece como panel flotante (`shadow.2`) al llevar el puntero al borde izquierdo. Plegado y desplegado se recuerdan por dispositivo.

**Accesibilidad.** `<nav aria-label="Secciones">` con lista anidada; el chevron es un botón propio con `aria-expanded` (tocar el nombre navega, tocar ▸ despliega); el panel móvil atrapa el foco y Esc lo cierra. Pares: C-125 a C-128.

### 7.10 Barra superior, pestañas de subsección y barra inferior (G-02, G-05)

**Barra superior (G-02).** Alto `layout.topbar.height.mobile` (56 px en < 768) o `layout.topbar.height.desktop` (48 px desde 768) sobre `ui.surface.page`; migas en `font.size.small` y `ui.text.secondary` («‹ Padre» en móvil, completas desde 768 recortadas al centro si no caben); título de página en `font.size.title-1` debajo, con sus propiedades clave; a la derecha, acción primaria y «⋯» (`ellipsis`) con las secundarias; desde 768, acceso a búsqueda y «+ Registrar». Estados: los de sus botones (§7.2); la migaja actual lleva `aria-current="page"`.

**Pestañas de subsección** (Resumen · Detalle · Tendencias · Calendario). Texto en `font.size.small`: reposo `ui.text.secondary`; hover `ui.text.default`; foco con anillo común; activa en `ui.text.default` con `font.weight.medium` y un indicador inferior de 2 px en `ui.border.selected`; deshabilitado y error no aplican. Fila con desplazamiento horizontal en < 768 (la activa siempre visible) y completa desde 768. Patrón ARIA *tabs* (`role="tablist"`, ← → entre pestañas).

**Barra inferior (G-05, solo < 768).** Alto `layout.bottombar.height` con la zona segura de Android; fondo `ui.surface.raised` y borde superior `ui.border.divider`; 5 destinos en orden fijo: Inicio (`house`), Comidas (`utensils`), «+» (`plus`), Compras (`shopping-cart`), Buscar (`search`); cada destino mide 56 px de alto (`size.bottombar.item.height`, el de la barra) y al menos 64 px de ancho (`size.bottombar.item.min-width`), con icono `size.icon.md` y nombre en `font.size.caption`.

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | icono y nombre en `ui.text.secondary` |
| Hover | no aplica (solo táctil) |
| Foco | anillo interior del destino |
| Activo | destino actual: píldora de 56 × 32 px (`size.bottombar.indicator-width` × `size.bottombar.indicator-height`) en `ui.surface.selected` detrás del icono `ui.icon.selected`, nombre en `ui.text.default` y `font.weight.medium`, `aria-current="page"` |
| «+» | círculo de 40 px en `ui.action.primary.bg` con `plus` en `ui.text.on-primary`; abre P-03 (primer toque de F-1) |
| Deshabilitado · error | no aplican (sin conexión, «+» sigue abriendo P-03: gastos y peso se capturan en la cola, K9) |

Pares: C-129 a C-134.

### 7.11 Menú «/» y búsqueda (G-03, G-04)

**Anatomía.** Un solo panel (U1 §5.1): campo de búsqueda arriba (icono `search`, `font.size.body`); resultados agrupados («En esta pantalla», «Registrar», «Ir a», «Vista», «Herramientas»; encabezados en `font.size.caption` y `ui.text.secondary`); cada fila con icono, nombre del comando en `ui.text.default`, alias en `ui.text.secondary` y, en escritorio, su atajo en una tecla dibujada (`radius.sm`, borde `ui.border.divider`, `font.size.caption`). Fondo `ui.surface.raised`, `radius.lg`, `shadow.2`.

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | filas sobre `ui.surface.raised` |
| Hover · resaltado | una sola fila resaltada a la vez (puntero o ↑ ↓) en `ui.surface.hover` |
| Foco | el foco se queda en el campo; la fila resaltada se indica con `aria-activedescendant` |
| Activo | Enter ejecuta el comando resaltado; nada se guarda desde el menú (U1 §5.1 punto 4) |
| Deshabilitado | comandos que requieren conexión: se muestran con `ui.text.disabled` y «Requiere conexión» |
| Error · vacío | «Sin resultados para "texto"» con sugerencia de /ir |

**Tamaños.** Panel de `layout.command.width` (640 px) centrado a 20 % del alto desde 768; filas de 36 px (fino) o 44 px (grueso); máximo 8 filas visibles con desplazamiento.

**Móvil y escritorio.** < 768: hoja a pantalla completa desde la pestaña «Buscar»; el «/» al inicio del texto cambia a comandos. Escritorio: «/» fuera de campos, `Ctrl+K` o `Ctrl+P`; Esc cierra y devuelve el foco a donde estaba.

**Accesibilidad.** Patrón ARIA *combobox* con *listbox* agrupado; búsqueda sin acentos ni mayúsculas; el número de resultados se anuncia con `aria-live="polite"`. Pares: C-135 a C-138.

### 7.12 Vista previa de registro, hojas y diálogos (G-06)

| Pieza | Anatomía y medidas | Móvil · escritorio |
|---|---|---|
| Vista previa (G-06) | encabezado con «‹» o `x`, título editable en el lugar, propiedades (4 + «Mostrar N más» en móvil) y cuerpo; fondo `ui.surface.raised`, `shadow.1`, ancho `layout.preview.width.desktop` (400 px, a 1024) o `layout.preview.width.wide` (480 px, a 1440) | página completa en < 768; panel central a 768; panel lateral derecho desde 1024 (más ancho a 1440) |
| Hoja inferior | asa de 32 × 4 px en `ui.border.control`, título en `font.size.title-2`, contenido con desplazamiento propio y primaria de 48 px fija abajo; `radius.lg` solo arriba, `shadow.2`, `overlay.scrim` detrás | todas las hojas y formularios en < 768 (U1 §1.3); se ajusta al teclado |
| Diálogo | título, cuerpo en `ui.text.secondary`, acciones a la derecha (primaria al final); `layout.dialog.width.short` 480 px (confirmación) o `layout.dialog.width.long` 640 px (formulario); `radius.lg`, `shadow.2` | desde 768; en < 768 se vuelve hoja |
| Diálogo de acción crítica | encabezado con `shield-alert` en `color.critical.light` / `color.critical.dark` (C-139), explica el efecto antes de confirmar (U1 §2.6: restaurar un respaldo, interruptor A/B), botón «Cancelar» de texto y botón de acción crítica con verbo y cantidad | ídem; nunca se confirma con Enter por omisión: el foco inicial va a «Cancelar» |

Estados: los de sus controles; al abrir, el foco va al primer campo (o a «Cancelar» en confirmaciones), queda atrapado dentro y Esc cierra devolviendo el foco al disparador; con cambios sin guardar, cerrar pregunta «¿Descartar cambios?». Error: el mensaje va junto al campo y la hoja no se cierra. `role="dialog"` y `aria-modal="true"` con `aria-labelledby` al título.

### 7.13 Panel de Ajustes (P-26)

**Anatomía.** Las 14 secciones de decisiones §G en su orden, con los iconos de §5.4. Escritorio: lista de secciones a la izquierda (240 px, mismos estados que la barra lateral pero sobre `ui.surface.page`) y contenido a la derecha con ancho de lectura `layout.reading.max-width`. Cada sección es una lista de filas: nombre (`font.size.body`, `ui.text.default`), descripción (`font.size.small`, `ui.text.secondary`) y control a la derecha (interruptor, campo, selección o segmentada, §7.3). **Parámetro operativo** (K44): la fila muestra además «Valor inicial: 80 %» en `ui.text.secondary` y el botón de texto «Restablecer al valor inicial» (`rotate-ccw`), visible solo cuando el valor difiere del inicial. Guardado automático: junto al control aparece `circle-check` en `color.role.success.light` / `.dark` con «Guardado» durante 2 s. La zona de acciones irreversibles va al final de «Datos y respaldos», separada por un divisor, y abre el diálogo de acción crítica (§7.12).

**Avisos y frecuencias.** En «Notificaciones» y «Reportes y frecuencias», cada aviso o reporte muestra hora de envío (`HH:MM`, §7.3), días de anticipación (entero ≥ 0), día de la semana (opción segmentada Lun … Dom) y día del mes (número de 1 a 31). Los umbrales de alerta (80 % y 100 %) y los días antes del pago (1 y 0) son claves `alerts.*` de `app_setting` y se muestran como parámetros operativos. **PROVISIONAL (bloqueo `docs/estado/bloqueos/S1-parametros-avisos.md`):** dónde se guardan hora, anticipación, día de la semana y día del mes, en `params_json` de cada aviso (operación `updateNotification`) o en claves de `app_setting` (operación de settings); el componente es el mismo en las dos variantes y solo cambia la operación que guarda.

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | filas sobre `ui.surface.page` con divisor `ui.border.divider` |
| Hover | fila sin cambio; hover del control (§7.3) |
| Foco | anillo del control |
| Activo | sección abierta en la lista: `ui.surface.selected`, `aria-current="page"` |
| Deshabilitado | sin conexión, todos los controles en solo lectura con «Requiere conexión con el servidor de casa»; Apariencia sigue funcionando (es local) |
| Error | el control vuelve al valor anterior con el mensaje de §7.0; un valor fuera de rango no se guarda y dice el rango |

**Móvil.** < 768: la lista de 14 secciones es la página; cada sección abre como página con «‹ Ajustes». **Accesibilidad.** Cada fila es un grupo con su etiqueta; «Restablecer» dice qué restablece («Restablecer umbral de alerta a 80 %»). Pares: C-140 y C-141 más los de sus controles.

### 7.14 Estados vacío, carga, error y estados globales (G-07)

| Estado | Anatomía y texto | Tokens |
|---|---|---|
| Carga | esqueleto con la forma final del contenido (U1 §2.4), estático (§6.2); `loader-circle` solo en acciones de más de 400 ms; a los 10 s: «Esto está tardando más de lo normal · Reintentar»; nunca cifras provisionales | `ui.skeleton.on-page`, `ui.skeleton.on-card` |
| Vacío | icono de la entidad en `size.icon.md`, una frase de qué falta y el botón que lo llena («Aún no hay gastos este mes · Registrar gasto»); nunca una pantalla en blanco | icono `ui.icon.default`, texto `ui.text.secondary`, botón §7.2 |
| Error | franja en la página con `circle-alert`, qué pasó y qué hacer, y «Reintentar» (`refresh-cw`); lo capturado no se pierde | fondo `color.role.error.surface-light` / `.surface-dark`, icono del rol, texto `ui.text.default` |
| Falta capturar (K18) | franja con `triangle-alert`: «Falta capturar · Capturar», abre P-26.05 en esa línea | superficie e icono de advertencia, texto `ui.text.default` |
| Sin conexión (G-07, K9) | franja fija bajo la barra superior: `wifi-off`, «Sin conexión con el servidor de casa» y «· n en cola» si hay capturas guardadas; lo cargado sigue visible con su hora de carga | superficie e icono de advertencia, texto `ui.text.default` |
| Sesión vencida | P-00 con aviso; al entrar vuelve a la página donde estaba | franja informativa |
| Registro inexistente o cambiado | «Ya no existe · Volver» (`file-x`); «Este registro cambió · Recargar» | franja informativa |

Franjas en página (avisos tipo *callout* de Notion): `radius.md`, relleno `space.3`, icono `size.icon.sm` del rol y texto en `ui.text.default` sobre la superficie del rol (pares C-145 a C-152: información, advertencia, error y éxito). Las franjas de error y sin conexión usan `role="alert"` la primera vez y luego quedan como texto; las demás, `role="status"`.

### 7.15 Avisos breves (G-07)

**Anatomía.** Fondo `ui.surface.raised`, `radius.md`, `shadow.2` en claro y borde `ui.border.divider` en oscuro; icono del rol (`circle-check` éxito, `circle-alert` error, `triangle-alert` advertencia, `info` información) en su color `light` / `dark`; texto en `font.size.small` y `ui.text.default`; acción opcional como enlace («Deshacer», `undo-2`) y `x` para cerrar. Texto: «Guardado · Deshacer»; sin conexión, «Gasto guardado en el teléfono · Se enviará al volver a la red de casa · Deshacer» (U1 §2.4).

**Comportamiento.** Uno visible a la vez y los demás en cola; el de éxito dura 5 s (U1 G-07) y se pausa con hover o foco; los de error se quedan hasta cerrarlos; entra y sale según §6.2. Posición: centrado abajo, a `space.2` sobre la barra inferior en < 768; abajo a la derecha a `space.6` de los bordes desde 768. Ancho máximo 480 px; en 360 ocupa el ancho menos `space.4` por lado.

| Estado | Comportamiento |
|---|---|
| Reposo | visible con su temporizador |
| Hover · foco | el temporizador se pausa; la acción y `x` tienen sus estados (§7.2) |
| Activo | «Deshacer» revierte y el aviso cambia a «Se deshizo» |
| Deshabilitado | «Deshacer» desaparece cuando ya no se puede deshacer |
| Error | variante de error con icono y sin cierre automático |

**Accesibilidad.** Éxito e información en `role="status"` (`aria-live="polite"`), error en `role="alert"`; el aviso no roba el foco; «Deshacer» se alcanza con Tab desde el control que lo provocó. Pares: C-153 a C-157.

### 7.16 Matriz de cobertura

| Componente | U1 §1.4 | Plan §4.1 punto 7 / U2-12 | Sección |
|---|---|---|---|
| Barra lateral | G-01 | barra lateral colapsable | §7.9 |
| Barra superior | G-02 | — | §7.10 |
| Menú «/» | G-03 | menú `/` | §7.11 |
| Búsqueda | G-04 | — | §7.11 |
| Barra inferior | G-05 | — | §7.10 |
| Vista previa de registro | G-06 | — | §7.12 |
| Avisos breves y estados globales | G-07 | estados vacío/carga/error y toasts | §7.14, §7.15 |
| Botones | — | botones | §7.2 |
| Campos | — | campos | §7.3 |
| Tabla, tablero, calendario, galería | — | tablas y vistas de base de datos | §7.5 |
| Tarjetas de métrica | — | tarjetas de métrica | §7.6 |
| Gráficas | — | gráficas | §7.7, §8 |
| Calendario de pagos | — | calendario de pagos | §7.8 |
| Panel de Ajustes | — | panel de Ajustes | §7.13 |
| Chips de persona y tienda | — | refuerzo de forma (K41) | §7.4 |

## 8. Patrones de visualización

Reglas para toda gráfica de la app, del PDF de reportes y de la Pista 3. Siguen el método de la guía de visualización de datos que cita el plan (§4.1 punto 8) con los valores de este sistema; donde la guía y las decisiones difieren, rigen las decisiones: cifras tabulares en todo dinero (§3.3, decisiones §G) y patrón de refuerzo siempre presente (K27), no solo como opción.

### 8.1 Marcas

| Marca | Especificación |
|---|---|
| Barra y columna | grosor ≤ 24 px (el resto de la banda queda como aire); extremo de datos redondeado 4 px, recto en la base; siempre desde cero |
| Línea | 2 px, uniones y extremos redondeados |
| Marcador | 10 px, relleno del color de la serie, con anillo de 2 px del color de la superficie (`ui.surface.card`) para leerse donde cruza otra línea |
| Área | relleno del color de la serie mezclado al 20 % con la superficie (derivado con `color-mix`, ningún color nuevo) y borde superior de 2 px en el color pleno |
| Separación | 2 px del color de la superficie entre segmentos apilados y entre barras contiguas; nunca un borde dibujado alrededor de una marca |
| Retícula y ejes | retícula horizontal de 1 px sólida en `ui.viz.grid`; línea base en `ui.viz.baseline`; sin retícula vertical salvo «Hoy» |
| Eje | **un solo eje Y**: nunca doble eje; dos medidas de escala distinta van en dos gráficas |
| Texto | etiquetas, valores y leyendas en `ui.text.default` o `ui.text.secondary`, **nunca en el color de la serie**; la identidad la da la marca junto al texto |

### 8.2 Qué forma usar

| Pregunta | Forma | Reglas |
|---|---|---|
| ¿Cuánto en cada categoría? (gasto por categoría, presupuesto contra real) | **barras horizontales** | orden descendente por monto (personas y tiendas en su orden fijo); una sola serie en `color.chart.1.light` / `.dark` sin patrón; el plan, como marca vertical de 2 px en `ui.viz.baseline` sobre la barra del real |
| ¿Cómo cambia en el tiempo? (gasto mensual, peso, kcal, precio) | **líneas** | eje X en el tiempo; marcadores en cada punto con ≤ 12 puntos y, con más, solo en el último y en el enfocado; el eje Y puede no empezar en cero (peso) y entonces se rotula |
| ¿Cuánto en cada mes? (pocas series) | **columnas** | máximo 3 series agrupadas; más series pasan a líneas |
| ¿Cuánto se ha devengado y cuánto se ha pagado? | **área apilada** (§8.3) | — |
| ¿De qué se compone un total? (≤ 5 partes) | **dona** (§8.4) | más de 5 partes o partes parecidas: barras horizontales |
| Una sola cifra | **tarjeta de métrica** (§7.6), no gráfica | — |

### 8.3 Área apilada: devengado contra efectivo

Muestra la diferencia entre lo que el crédito ya generó (devengado) y lo que ya se pagó (efectivo), base de la deuda viva (K19).

- Eje X: fechas del calendario de pagos; eje Y: dinero acumulado desde cero.
- **Capa 1 (abajo), efectivo acumulado**: serie 1, relleno liso.
- **Capa 2 (arriba), devengado sin pagar** = devengado acumulado − efectivo acumulado: serie 2 con su patrón; la cima de la pila es el devengado acumulado.
- Línea vertical «Hoy» de 1 px en `ui.viz.baseline` con su etiqueta.
- Si en un tramo lo pagado supera a lo devengado (pago adelantado), la capa 2 vale cero en ese tramo y el texto lo dice («pagado por adelantado»); nunca se dibuja una capa negativa.
- Etiquetas en texto al final: devengado acumulado y efectivo acumulado; el subtítulo da la deuda viva completa con su suma, por ejemplo al 10-abr-2027 «Deuda viva $6,293.44 = $1,258.70 devengado sin pagar + $5,034.74 por devengar» (cifras de K19). La descripción emergente da las dos series y su diferencia.

### 8.4 Dona (composición de ≤ 5 partes)

Solo para partes de un total en un mismo periodo, nunca para el tiempo. Con más de 5 partes se muestran las 4 mayores y «Otros» en `ui.viz.other`, o se usan barras. Grosor del anillo 40 % del radio; 2 px de separación entre partes; partes en orden descendente desde las 12 en sentido horario («Otros» al final); el total va al centro en `font.size.title-3` con `tabular-nums`; la leyenda es una lista con marca, nombre, monto y porcentaje en texto (`$1,258.68 · 23 %`). Colores: series en su orden (§8.6) o, si las partes son personas o tiendas, sus categóricas (§8.7).

### 8.5 Etiquetas de dinero siempre en texto

Ningún monto vive solo en la descripción emergente: barras con el valor al final de cada barra; columnas con el valor sobre la columna si cabe con aire, y si no, en la tabla; líneas con el último valor y los extremos; área con los valores finales; dona con la leyenda. Formato completo (`$1,258.68`) en etiquetas y leyendas; el compacto (`$5 mil`) solo en ejes. Una etiqueta que no cabe no se recorta: se mueve fuera de la marca o pasa a la tabla. «Ver como tabla» siempre está disponible (§7.7). Etiquetado selectivo: nunca un número en cada punto.

### 8.6 Series: orden fijo, forma y patrón (K27)

Las 8 series de §2.7 se asignan **en este orden y nunca se reciclan**; una novena se agrupa en «Otros». El color sigue a la entidad, no a su rango: un filtro que quita series no repinta las que quedan. Cada serie lleva color, forma de marcador, trazo y patrón de relleno; con 2 o más series hay leyenda siempre y, con ≤ 4, también etiqueta directa.

| Serie | Color | Marcador | Trazo (líneas; guion y hueco en px) | Patrón (barras, áreas, dona) |
|---|---|---|---|---|
| 1 | `color.chart.1.light` / `color.chart.1.dark` | círculo | continuo | liso |
| 2 | `color.chart.2.light` / `color.chart.2.dark` | cuadrado | 8 · 4 | rayas a 45° cada 6 px |
| 3 | `color.chart.3.light` / `color.chart.3.dark` | triángulo | 2 · 4 (puntos) | rayas a 135° cada 6 px |
| 4 | `color.chart.4.light` / `color.chart.4.dark` | rombo | 8 · 4 · 2 · 4 | puntos de 2 px cada 6 px |
| 5 | `color.chart.5.light` / `color.chart.5.dark` | triángulo invertido | 4 · 4 | cuadrícula diagonal cada 8 px |
| 6 | `color.chart.6.light` / `color.chart.6.dark` | cruz | 12 · 4 | rayas a 45° cada 10 px |
| 7 | `color.chart.7.light` / `color.chart.7.dark` | aspa | 8 · 3 · 2 · 3 · 2 · 3 | rayas a 135° cada 10 px |
| 8 | `color.chart.8.light` / `color.chart.8.dark` | anillo | 2 · 2 | puntos de 3 px cada 8 px |

El patrón se dibuja con trazo de 1.5 px en el color de la superficie sobre el relleno pleno (en áreas, con el color pleno de la serie sobre el relleno al 20 %); el color de la serie cubre al menos el 60 % de la marca, así su contraste contra la superficie es el de la tabla de pares (C-106 a C-113). Con `forced-colors: active`, marcador, trazo y patrón se conservan.

### 8.7 Personas y tiendas en gráficas

- El color sigue a la entidad con su categórica (§2.6), no con la paleta de series; el marcador y el trazo siguen la posición de la entidad en su orden fijo (la primera persona o tienda usa el marcador y el trazo de la serie 1, la segunda los de la serie 2, y así).
- **Personas y tiendas nunca comparten gráfica ni leyenda** (K41); en un tablero que muestra las dos, cada gráfica lleva título con el grupo («Gasto por persona», «Gasto por tienda») y etiquetas directas con el nombre.
- Ximena y Karen (17° de separación) y Axel y Walmart en oscuro se distinguen por marcador, trazo, orden y etiqueta, no por color (§2.6).

### 8.8 No se hace

Doble eje; 3D; arcoíris o colores generados fuera de C2; dona con más de 5 partes o para el tiempo; montos solo en la descripción emergente; animación de entrada o de cambio de cifras (§6.4); atenuar o repintar series al filtrar; texto en el color de la serie; color como único portador de significado.

## 9. Accesibilidad y puntos de quiebre

### 9.1 Objetivo y cómo se cumple

**WCAG 2.2 nivel AA como mínimo, calculado** (decisiones §G). Idioma del documento `lang="es-MX"`.

| Criterio WCAG 2.2 | Cómo lo cumple este sistema |
|---|---|
| 1.3.1 Información y relaciones · 4.1.2 Nombre, función y valor | HTML semántico y patrones ARIA declarados en cada componente (§7): *grid*, *tabs*, *combobox*, *dialog*, *meter* |
| 1.3.4 Orientación | ninguna pantalla se bloquea en vertical u horizontal; en horizontal rige el ancho resultante (U1 §6.1 regla 7) |
| 1.4.1 Uso del color | signo, icono, forma, patrón y texto acompañan siempre al color (§2.10, §7.4, §8.6) |
| 1.4.3 Contraste (mínimo) · 1.4.11 Contraste no textual | texto ≥ 4.5:1; texto grande, interfaz, foco y gráficas ≥ 3:1; lista de pares de §10.5 calculada por U2d, 0 bajo umbral |
| 1.4.4 Cambio de tamaño · 1.4.10 Reajuste | a 200 % de texto y a 320 px de ancho CSS nada se pierde ni pide desplazamiento horizontal de la página; solo tabla, tablero y fila de pestañas se desplazan en su propio contenedor (U1 §6.1 regla 3) |
| 1.4.12 Espaciado del texto | ningún contenedor con alto fijo recorta texto si se amplía interlínea o espaciado; las cifras no se truncan (§3.3) |
| 1.4.13 Contenido al pasar el puntero o al enfocar | las descripciones emergentes se cierran con Esc, se pueden recorrer con el puntero y no desaparecen solas |
| 2.1.1 Teclado · 2.1.2 Sin trampas | todo se hace con teclado (U1 §5.3); solo diálogos y hojas atrapan el foco, y Esc los cierra |
| 2.1.4 Atajos de una tecla | `/`, `?` y las flechas solo actúan fuera de campos de texto (U1 §5.4); PROPUESTA: interruptor «Atajos de una tecla» en Apariencia para desactivarlos (§9.6) |
| 2.3.1 Tres destellos | nada parpadea (§6.4) |
| 2.4.3 Orden del foco · 2.4.7 Foco visible · 2.4.11 Foco no oculto | el orden sigue la lectura; anillo de §9.2; `scroll-padding` igual a las barras fijas para que nada tape el elemento enfocado |
| 2.5.7 Movimientos de arrastre | arrastrar siempre tiene alternativa de un toque («Mover a…», §7.5.3) |
| 2.5.8 Tamaño del objetivo | todo control táctil ≥ 44 × 44 px (§9.3), por encima del piso de 24 px |
| 3.3.1 Identificación de errores · 3.3.3 Sugerencia | el mensaje dice qué pasó y qué hacer (§1.3, §7.0) |
| 3.3.7 Entrada redundante · 3.3.8 Autenticación accesible | los formularios proponen lo ya capturado (últimas categorías, tienda); el acceso admite gestores de contraseñas y pegar, sin pruebas cognitivas |

### 9.2 Foco visible

Anillo de 2 px (`border.width.focus`) en `ui.border.focus` (`color.brand.600` en claro, `color.brand.400` en oscuro) con separación de 2 px (`focus.offset`), por `:focus-visible`; en filas, celdas y elementos de lista va por dentro (`outline-offset: -2px`). Contraste medido contra página, tarjeta y fila seleccionada (C-006 a C-008). Nunca `outline: none` sin reemplazo; el anillo aparece sin transición. Con `forced-colors: active` el anillo usa el color del sistema `Highlight`.

### 9.3 Objetivo táctil ≥ 44 px

- Todo control táctil ocupa al menos `size.touch.min` (**44 × 44 px**) aunque su dibujo sea menor: casilla de 20 px, icono de 20 px, chip de 24 px, punto de una gráfica.
- Con puntero grueso (`any-pointer: coarse`) los controles miden 44 px de alto (`size.control.coarse`); el primario de hoja, 48 px; los destinos de la barra inferior, 56 px de alto (`size.bottombar.item.height`, igual a `layout.bottombar.height`) y al menos 64 px de ancho (`size.bottombar.item.min-width`) (§4.6).
- Entre objetivos contiguos hay `space.2` (8 px) como mínimo.
- Los 36 px de `size.control.fine` solo aplican con puntero fino y sin pantalla táctil.

### 9.4 Puntos de quiebre

Se diseña a 360 px y se amplía (móvil primero, U1 §6.1). Cinco anchos de referencia; las consultas son `min-width`.

| Token | Valor | Rango | Qué cambia (detalle en U1 §6.2 y §4.2) |
|---|---|---|---|
| `breakpoint.base` | 360 px | hasta 411 px | base sin consulta: 4 columnas, barra inferior, hojas a pantalla completa, tabla de 2 líneas, agenda |
| `breakpoint.phone` | 412 px | 412–767 px | un dato secundario más por tarjeta y fila; galería de 2 columnas |
| `breakpoint.tablet` | 768 px | 768–1023 px | 8 columnas; desaparece la barra inferior; tarjetas en 2 columnas; diálogos centrados; calendario mensual; barra lateral superpuesta |
| `breakpoint.desktop` | 1024 px | 1024–1439 px | 12 columnas; barra lateral fija y colapsable; tabla completa; vista previa lateral |
| `breakpoint.wide` | 1440 px | desde 1440 px | más aire; páginas de texto a 720 px; vista previa más ancha; propiedades en 2 columnas |

Además: a 320 px (200 % de zoom en 640 px o 400 % en 1280 px) todo se reajusta a una columna sin desplazamiento horizontal de la página (1.4.10). Referencias de prueba de U5: Android 412 × 915 y escritorio 1440 × 900 (U1 §6.1).

### 9.5 Pruebas que pide este sistema

U5 y WF-3 comprueban: navegación completa solo con teclado en 1440; TalkBack con Chrome en el Android de referencia; un lector de pantalla de escritorio en Linux; zoom de texto al 200 % y ancho de 320 px; vista en escala de grises (P4) y con las simulaciones de daltonismo de colorimetría §8.4; `forced-colors`; y la tabla de contraste de U2d con 0 pares bajo umbral.

### 9.6 Preferencias del sistema

| Preferencia | Comportamiento |
|---|---|
| `prefers-color-scheme` | tema por defecto según el sistema; el interruptor (Sistema · Claro · Oscuro, en Apariencia, en el pie de la barra lateral, `Ctrl+Shift+L` y /tema) fija `data-theme` en la raíz y se recuerda por dispositivo |
| `prefers-reduced-motion` | se anulan las transiciones no esenciales (§6.3) |
| `forced-colors` | colores del sistema; se conservan anillo de foco, bordes de controles, forma, trazo y patrón de las series |
| Atajos de una tecla | PROPUESTA de U2b-2 para WCAG 2.1.4: interruptor local en Apariencia, encendido por omisión; lo decide Axel en la puerta 2 y lo agregan U1 y U4 |

## 10. Tokens

### 10.1 Archivos y flujo

| Archivo | Quién | Qué es |
|---|---|---|
| `docs/contratos/C2-colores.json` | colorimetría 1.2 (`docs/design/tools/colorimetria.py`) | 190 primitivos de color; no se edita a mano |
| `docs/design/DESIGN.md` §3–§10.6 | U2b-1, U2b-2 | nombres y valores de los tokens no de color (§3–§6, §9.3–§9.4, §10.6), los 41 semánticos `ui.*` (§7.1) y la lista de pares de componente (§10.5) |
| `docs/design/tools/tokens/generar_c2_tokens.py` | U2c | arma `C2-tokens.json`: copia los primitivos sin cambiar un hex (con alias para los pasos de rol, §10.2), lee §7.1 y transcribe los valores no de color de este documento |
| `docs/contratos/C2-tokens.json` | U2c, **generado** | contrato C2 en DTCG y **única fuente de valores de diseño para quien lo consume** (frontend, PDF de reportes, Pista 3 y la tabla de contraste de componentes); no se edita a mano |
| `docs/contratos/tokens.css` | U2c, compilado con Style Dictionary 5 | variables CSS con tema claro, oscuro e interruptor |
| `docs/contratos/tailwind.theme.css` | U2c, compilado con Style Dictionary 5 | tema de Tailwind 4: bloque `@theme inline` que expone los semánticos como utilidades (§10.4) |

**Cadena única:** DESIGN.md §3–§10.6 + `C2-colores.json` → `generar_c2_tokens.py` → `C2-tokens.json` (no se edita a mano) → Style Dictionary → `tokens.css` y `tailwind.theme.css`. Un valor no de color se cambia en este documento y en el generador, que lo transcribe; un color nuevo entra por `colorimetria.py`, nunca directo en los tokens (§2.10 regla 4). Después se corre `npm run check` en `docs/design/tools/tokens/`, que regenera, compila, verifica, comprueba la sincronía y prueba Tailwind (README de esa carpeta). Nadie escribe un valor en `C2-tokens.json`, `tokens.css` ni el tema (C2-09).

### 10.2 Estructura de `C2-tokens.json` (DTCG)

| Grupo | `$type` | Contenido | Origen |
|---|---|---|---|
| `color` | `color` | los 190 primitivos con sus nombres actuales (`color.brand.50` … `color.chart.8.dark`); en los 8 roles y la acción crítica, `light`, `dark`, `surface-light` y `surface-dark` son alias de los pasos 700, 400, 100 y 900 (p. ej. `color.role.income.light` = `{color.role.income.700}`), con el mismo hex | `C2-colores.json`, mismos hex |
| `ui` | `color` | 75 semánticos: los 41 de §7.1 más 34 pares de rol (8 de texto e icono y 8 de insignia, p. ej. `ui.role.income` y `ui.role.income-surface`), acción crítica (`ui.critical` y `ui.critical-surface`), persona (4), tienda (4) y serie (8); el velo `overlay.scrim` va en su propio grupo. Cada uno es un grupo con hojas `light` y `dark` cuyo `$value` es una referencia, p. ej. `ui.surface.page.light` = `{color.neutral.50}`; los alias de §7.1 apuntan a la hoja del otro semántico (`ui.icon.default.light` = `{ui.text.secondary.light}`) | §7.1 (41); §2.4–§2.7 (34) |
| `font` | `fontFamily`, `dimension`, `fontWeight`, `number` | `font.family.sans`, `font.size.*` (6), `font.line-height.*` (6), `font.weight.*` (4), `font.letter-spacing` | §3 |
| `space` | `dimension` | `space.0` … `space.12` (8 valores) | §4.1 |
| `radius` | `dimension` | `radius.sm`, `radius.md`, `radius.lg`, `radius.full` | §4.4 |
| `border` · `focus` | `dimension` | `border.width.default`, `border.width.focus`, `focus.offset` | §4.5 |
| `shadow` | `shadow` | `shadow.0` (sin sombra: capa nula con alfa 0 que se emite como `none`), `shadow.1` y `shadow.2` (solo claro: capas con color `{color.neutral.950}` y su alfa, de 0.06 a 0.10, en `$extensions.finanzas.alpha`; en oscuro `--ui-shadow-1` y `--ui-shadow-2` valen `none`) | §4.5 |
| `overlay` | `color` | `overlay.scrim` con hojas `light` (alfa 0.40) y `dark` (alfa 0.60) sobre `color.neutral.950` | §4.5 |
| `size` · `layout` · `icon` | `dimension`, `number` | medidas de §4.3, §4.6, §5.2 y las de componente de §10.6; una medida con dos valores es un grupo sin `$value` y se cita por sus hojas: `layout.topbar.height.mobile` y `layout.topbar.height.desktop`, `layout.preview.width.desktop` y `layout.preview.width.wide`, `layout.dialog.width.short` y `layout.dialog.width.long`, `size.bottombar.item.height` y `size.bottombar.item.min-width`, `icon.stroke.user-units.md` y `icon.stroke.user-units.sm` | §4–§5, §10.6 |
| `motion` | `duration`, `cubicBezier` | `motion.duration.fast/base/slow`, `motion.easing.standard/emphasized` | §6 |
| `breakpoint` | `dimension` | `base` 360, `phone` 412, `tablet` 768, `desktop` 1024, `wide` 1440 | §9.4 |
| `viz` | `strokeStyle`, `string` en `$extensions` | `viz.series.N.dash` (trazo de §8.6) y, en `$extensions`, marcador y patrón de cada serie | §8.6 |

Reglas: cada token lleva `$description` en español de México que cita su sección; nombres en inglés corto; ningún semántico tiene valor propio (todos son referencias) y ningún primitivo repite el hex de otro sin ser su alias (los `light`, `dark` y `surface-*` de rol y de acción crítica apuntan a su paso), así dos colores con el mismo valor siempre quedan relacionados (C2-10; `verificar_c2.py` lo comprueba). Entre los no de color, dos medidas iguales con distinta idea de diseño (p. ej. 24 px de `font.size.title-2` y de `size.icon.md`) no llevan alias; `verificar_c2.py` las lista como información.

```json
{
  "ui": {
    "surface": {
      "page": {
        "$type": "color",
        "$description": "Fondo de página (DESIGN §7.1)",
        "light": { "$value": "{color.neutral.50}" },
        "dark": { "$value": "{color.neutral.950}" }
      }
    }
  }
}
```

### 10.3 `tokens.css`

- **Primitivos estáticos** en `:root`, iguales en los dos temas: `--color-neutral-50`, `--color-role-income-700`, `--color-role-income-light`, `--color-chart-1-dark`, etc. (ruta con guiones).
- **Variables temáticas `--ui-*`** (78 por tema), con el mismo conjunto de nombres en los dos temas (C2-06): los 41 semánticos (`--ui-surface-page`, `--ui-text-default`…), una por cada par `light` / `dark` de C2 sin el sufijo (`--ui-role-income`, `--ui-role-income-surface`, `--ui-critical`, `--ui-critical-surface`, `--ui-person-axel`, `--ui-store-walmart`, `--ui-chart-1`…), `--ui-overlay-scrim` y las sombras `--ui-shadow-1` y `--ui-shadow-2` (valor en claro, `none` en oscuro).
- **Temas**: claro en `:root`; oscuro en `@media (prefers-color-scheme: dark) { :root:not([data-theme="light"]) { … } }` y en `:root[data-theme="dark"] { … }`; `color-scheme` acompaña a cada tema.
- **No de color**: `--font-sans`, `--font-size-caption`, `--line-height-caption`, `--font-weight-medium`, `--space-4`, `--radius-md`, `--duration-fast`, `--ease-standard`, `--size-touch-min`, `--layout-sidebar-width`, etc.; los puntos de quiebre se emiten como referencia (no sirven dentro de `@media`).
- `@media (prefers-reduced-motion: reduce)` con la regla de §6.3, y las `@font-face` de §3.1 con rutas locales.

### 10.4 Tema de Tailwind 4

- El tema importa solo `tokens.css` (`@import "./tokens.css";`), no `tailwindcss`: la hoja de entrada del frontend importa primero `tailwindcss` y luego el tema (`@import "tailwindcss";` seguido de `@import ".../docs/contratos/tailwind.theme.css";`, como dice el README de `docs/design/tools/tokens/`). Luego `@theme inline { … }` vacía las escalas por defecto con `--color-*: initial;` y lo mismo en `--font-*`, `--text-*`, `--font-weight-*`, `--radius-*`, `--shadow-*`, `--ease-*` y `--breakpoint-*`; solo quedan `transparent` y `currentColor`. Solo se exponen como utilidades los semánticos y las variables temáticas, nunca los primitivos: `--color-surface-page: var(--ui-surface-page);` (`bg-surface-page`), `--color-text-default: var(--ui-text-default);` (`text-text-default`), `--color-role-income: var(--ui-role-income);`, `--color-chart-1: var(--ui-chart-1);` y así para cada `--ui-*`. Así un componente no puede combinar primitivos a mano.
- Tipografía: `--font-sans`, `--text-caption` con `--text-caption--line-height` … `--text-title-1`, `--font-weight-regular` … `--font-weight-bold`; `--spacing: 0.25rem` (solo se usan los pasos 0, 1, 2, 3, 4, 6, 8 y 12; lo revisa WF-3); `--radius-sm`, `--radius-md`, `--radius-lg`; `--shadow-1: var(--ui-shadow-1)` y `--shadow-2: var(--ui-shadow-2)`; `--ease-standard`, `--ease-emphasized` (duraciones con `duration-150`, `duration-200`, `duration-250`).
- Puntos de quiebre: `--breakpoint-*: initial;` y `--breakpoint-phone: 25.75rem`, `--breakpoint-tablet: 48rem`, `--breakpoint-desktop: 64rem`, `--breakpoint-wide: 90rem` (360 es la base sin prefijo).
- Variante `coarse` con `@custom-variant coarse (@media (any-pointer: coarse));` para las medidas de 44 px. Los componentes no usan la variante `dark:`: el tema lo resuelven las variables `--ui-*`.
- Los valores no de color del tema se generan literales desde `C2-tokens.json` (no con `var()` a `tokens.css`) para evitar referencias circulares cuando el nombre coincide.

### 10.5 Pares de contraste que calcula U2d

Lista explícita de los pares por componente y estado que usan las secciones 2.9 y 7. La calcula `docs/design/tools/contraste_componentes.py` (U2d), que lee `docs/contratos/C2-tokens.json` como única fuente (C2-05), con la fórmula de §2.1 truncada a 2 decimales como `contraste.csv`, y escribe `docs/design/contraste-componentes.csv` (con `--check` comprueba que el archivo en disco es idéntico). **Hay dos tablas oficiales de contraste:** `docs/design/contraste.csv` para los pares de paleta (colorimetría, 202 pares) y `docs/design/contraste-componentes.csv` para los pares de componente de esta lista; juntas cubren todos los pares usados. **158 pares × 2 temas = 316 combinaciones**: 304 con umbral (texto 4.5:1; texto grande e interfaz 3:1) y 12 decorativas, listadas y exentas (divisores, pista de avance, retícula, esqueletos, borde de lo que flota). 50 combinaciones aún no están en `contraste.csv`. Se espera 0 bajo umbral (U2-04); si U2d encuentra alguna, abre un bloqueo con el par y no cambia ningún token por su cuenta. La resolución por tema sale de las referencias de `C2-tokens.json`: un semántico `ui.<x>` usa sus hojas `light` y `dark`; «rol X» usa `ui.role.<rol>`; «superficie de rol X», `ui.role.<rol>-surface`; «persona», «tienda» y «serie», `ui.person.<id>`, `ui.store.<id>` y `ui.chart.<n>`; «acción crítica», `ui.critical`. Las columnas «Claro» y «Oscuro» de esta tabla y la tabla de §7.1 son cruce secundario: si no coinciden con C2, el script avisa y sale con código 1.

| Id | Componente | Estado · elemento | Par | Tipo y umbral | Claro: primer plano / fondo | Oscuro: primer plano / fondo | En `contraste.csv` (claro · oscuro) |
|---|---|---|---|---|---|---|---|
| C-001 | Botón primario | reposo · texto | `ui.text.on-primary` sobre `ui.action.primary.bg` | texto 4.5 | `color.base.white` / `color.brand.700` | `color.brand.950` / `color.brand.300` | sí · sí |
| C-002 | Botón primario | reposo · relleno sobre página | `ui.action.primary.bg` sobre `ui.surface.page` | interfaz 3 | `color.brand.700` / `color.neutral.50` | `color.brand.300` / `color.neutral.950` | sí · sí |
| C-003 | Botón primario | reposo · relleno sobre tarjeta | `ui.action.primary.bg` sobre `ui.surface.card` | interfaz 3 | `color.brand.700` / `color.neutral.100` | `color.brand.300` / `color.neutral.900` | sí · sí |
| C-004 | Botón primario | hover · texto | `ui.text.on-primary` sobre `ui.action.primary.bg-hover` | texto 4.5 | `color.base.white` / `color.brand.800` | `color.brand.950` / `color.brand.200` | sí · sí |
| C-005 | Botón primario | presionado · texto | `ui.text.on-primary` sobre `ui.action.primary.bg-pressed` | texto 4.5 | `color.base.white` / `color.brand.900` | `color.brand.950` / `color.brand.100` | no · no |
| C-006 | Todos los controles | foco · anillo sobre página | `ui.border.focus` sobre `ui.surface.page` | interfaz 3 | `color.brand.600` / `color.neutral.50` | `color.brand.400` / `color.neutral.950` | sí · sí |
| C-007 | Todos los controles | foco · anillo sobre tarjeta | `ui.border.focus` sobre `ui.surface.card` | interfaz 3 | `color.brand.600` / `color.neutral.100` | `color.brand.400` / `color.neutral.900` | sí · sí |
| C-008 | Todos los controles | foco · anillo sobre fila seleccionada | `ui.border.focus` sobre `ui.surface.selected` | interfaz 3 | `color.brand.600` / `color.brand.100` | `color.brand.400` / `color.brand.900` | no · no |
| C-009 | Botón secundario | reposo · texto | `ui.text.default` sobre `ui.surface.page` | texto 4.5 | `color.neutral.900` / `color.neutral.50` | `color.neutral.100` / `color.neutral.950` | sí · sí |
| C-010 | Botón secundario | reposo · borde sobre página | `ui.border.control` sobre `ui.surface.page` | interfaz 3 | `color.neutral.600` / `color.neutral.50` | `color.neutral.400` / `color.neutral.950` | sí · sí |
| C-011 | Botón secundario | reposo · borde sobre tarjeta | `ui.border.control` sobre `ui.surface.card` | interfaz 3 | `color.neutral.600` / `color.neutral.100` | `color.neutral.400` / `color.neutral.900` | sí · sí |
| C-012 | Botón secundario | hover · texto | `ui.text.default` sobre `ui.surface.hover` | texto 4.5 | `color.neutral.900` / `color.neutral.200` | `color.neutral.100` / `color.neutral.800` | sí · sí |
| C-013 | Botón secundario | hover · borde sobre página | `ui.border.control-hover` sobre `ui.surface.page` | interfaz 3 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-014 | Botón secundario | presionado · texto | `ui.text.default` sobre `ui.surface.pressed` | texto 4.5 | `color.neutral.900` / `color.neutral.300` | `color.neutral.100` / `color.neutral.700` | no · no |
| C-015 | Botón de texto | reposo · texto sobre página | `ui.text.secondary` sobre `ui.surface.page` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-016 | Botón de texto | reposo · texto sobre tarjeta | `ui.text.secondary` sobre `ui.surface.card` | texto 4.5 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-017 | Botón de icono | reposo · icono sobre página | `ui.icon.default` sobre `ui.surface.page` | interfaz 3 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-018 | Botón de icono | reposo · icono sobre tarjeta | `ui.icon.default` sobre `ui.surface.card` | interfaz 3 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-019 | Botón de icono | presionado · icono | `ui.text.default` sobre `ui.surface.pressed` | interfaz 3 | `color.neutral.900` / `color.neutral.300` | `color.neutral.100` / `color.neutral.700` | no · no |
| C-020 | Botón destructivo | reposo · texto | `ui.text.on-destructive` sobre `ui.action.destructive.bg` | texto 4.5 | `color.base.white` / `color.role.error.700` | `color.role.error.950` / `color.role.error.300` | sí · sí |
| C-021 | Botón destructivo | hover · texto | `ui.text.on-destructive` sobre `ui.action.destructive.bg-hover` | texto 4.5 | `color.base.white` / `color.role.error.800` | `color.role.error.950` / `color.role.error.200` | sí · sí |
| C-022 | Botón destructivo | presionado · texto | `ui.text.on-destructive` sobre `ui.action.destructive.bg-pressed` | texto 4.5 | `color.base.white` / `color.role.error.900` | `color.role.error.950` / `color.role.error.100` | no · no |
| C-023 | Botón destructivo | reposo · relleno sobre diálogo | `ui.action.destructive.bg` sobre `ui.surface.raised` | interfaz 3 | `color.role.error.700` / `color.neutral.50` | `color.role.error.300` / `color.neutral.900` | sí · sí |
| C-024 | Botón de acción crítica | reposo · texto | `ui.text.on-critical` sobre `ui.action.critical.bg` | texto 4.5 | `color.base.white` / `color.critical.700` | `color.critical.950` / `color.critical.300` | sí · sí |
| C-025 | Botón de acción crítica | hover · texto | `ui.text.on-critical` sobre `ui.action.critical.bg-hover` | texto 4.5 | `color.base.white` / `color.critical.800` | `color.critical.950` / `color.critical.200` | sí · sí |
| C-026 | Botón de acción crítica | presionado · texto | `ui.text.on-critical` sobre `ui.action.critical.bg-pressed` | texto 4.5 | `color.base.white` / `color.critical.900` | `color.critical.950` / `color.critical.100` | no · no |
| C-027 | Botón de acción crítica | reposo · relleno sobre diálogo | `ui.action.critical.bg` sobre `ui.surface.raised` | interfaz 3 | `color.critical.700` / `color.neutral.50` | `color.critical.300` / `color.neutral.900` | sí · sí |
| C-028 | Botones (todos) | deshabilitado · texto | `ui.text.disabled` sobre `ui.surface.disabled` | texto 4.5 | `color.neutral.700` / `color.neutral.100` | `color.neutral.400` / `color.neutral.900` | sí · sí |
| C-029 | Enlace | reposo · sobre página | `ui.text.link` sobre `ui.surface.page` | texto 4.5 | `color.brand.700` / `color.neutral.50` | `color.brand.300` / `color.neutral.950` | sí · sí |
| C-030 | Enlace | reposo · sobre tarjeta | `ui.text.link` sobre `ui.surface.card` | texto 4.5 | `color.brand.700` / `color.neutral.100` | `color.brand.300` / `color.neutral.900` | sí · sí |
| C-031 | Enlace | hover · sobre página | `ui.text.link-hover` sobre `ui.surface.page` | texto 4.5 | `color.brand.800` / `color.neutral.50` | `color.brand.200` / `color.neutral.950` | sí · sí |
| C-032 | Enlace | hover · sobre tarjeta | `ui.text.link-hover` sobre `ui.surface.card` | texto 4.5 | `color.brand.800` / `color.neutral.100` | `color.brand.200` / `color.neutral.900` | sí · sí |
| C-033 | Enlace | sobre aviso breve o diálogo | `ui.text.link` sobre `ui.surface.raised` | texto 4.5 | `color.brand.700` / `color.neutral.50` | `color.brand.300` / `color.neutral.900` | sí · sí |
| C-034 | Campo | reposo · valor | `ui.text.default` sobre `ui.surface.field` | texto 4.5 | `color.neutral.900` / `color.neutral.50` | `color.neutral.100` / `color.neutral.950` | sí · sí |
| C-035 | Campo | reposo · marcador de posición | `ui.text.secondary` sobre `ui.surface.field` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-036 | Campo | reposo · borde contra el interior | `ui.border.control` sobre `ui.surface.field` | interfaz 3 | `color.neutral.600` / `color.neutral.50` | `color.neutral.400` / `color.neutral.950` | sí · sí |
| C-037 | Campo | hover · borde sobre tarjeta | `ui.border.control-hover` sobre `ui.surface.card` | interfaz 3 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-038 | Campo | error · borde sobre tarjeta | `ui.border.error` sobre `ui.surface.card` | interfaz 3 | `color.role.error.light` / `color.neutral.100` | `color.role.error.dark` / `color.neutral.900` | sí · sí |
| C-039 | Campo | error · borde sobre página | `ui.border.error` sobre `ui.surface.page` | interfaz 3 | `color.role.error.light` / `color.neutral.50` | `color.role.error.dark` / `color.neutral.950` | sí · sí |
| C-040 | Campo | error · mensaje e icono sobre tarjeta | rol error sobre `ui.surface.card` | texto 4.5 | `color.role.error.light` / `color.neutral.100` | `color.role.error.dark` / `color.neutral.900` | sí · sí |
| C-041 | Campo | error · mensaje e icono sobre página | rol error sobre `ui.surface.page` | texto 4.5 | `color.role.error.light` / `color.neutral.50` | `color.role.error.dark` / `color.neutral.950` | sí · sí |
| C-042 | Campo | deshabilitado · valor | `ui.text.disabled` sobre `ui.surface.disabled` | texto 4.5 | `color.neutral.700` / `color.neutral.100` | `color.neutral.400` / `color.neutral.900` | sí · sí |
| C-043 | Casilla e interruptor | marcado · palomita o perilla | `ui.text.on-primary` sobre `ui.action.primary.bg` | interfaz 3 | `color.base.white` / `color.brand.700` | `color.brand.950` / `color.brand.300` | sí · sí |
| C-044 | Casilla e interruptor | marcado · relleno sobre tarjeta | `ui.action.primary.bg` sobre `ui.surface.card` | interfaz 3 | `color.brand.700` / `color.neutral.100` | `color.brand.300` / `color.neutral.900` | sí · sí |
| C-045 | Interruptor | apagado · perilla sobre pista | `ui.surface.page` sobre `ui.border.control` | interfaz 3 | `color.neutral.50` / `color.neutral.600` | `color.neutral.950` / `color.neutral.400` | no · no |
| C-046 | Edición en el lugar | hover · valor | `ui.text.default` sobre `ui.surface.hover-subtle` | texto 4.5 | `color.neutral.900` / `color.neutral.100` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-047 | Edición en el lugar | hover · cifra de rol (gasto) | rol expense sobre `ui.surface.hover-subtle` | texto 4.5 | `color.role.expense.light` / `color.neutral.100` | `color.role.expense.dark` / `color.neutral.900` | sí · sí |
| C-048 | Edición en el lugar | hover · cifra de rol (ingreso) | rol income sobre `ui.surface.hover-subtle` | texto 4.5 | `color.role.income.light` / `color.neutral.100` | `color.role.income.dark` / `color.neutral.900` | sí · sí |
| C-049 | Edición en el lugar | hover sobre tarjeta · borde | `ui.border.control-hover` sobre `ui.surface.card` | interfaz 3 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-050 | Chip de persona o tienda | reposo · marca de forma | persona axel sobre `ui.surface.raised` | interfaz 3 | `color.person.axel.light` / `color.neutral.50` | `color.person.axel.dark` / `color.neutral.900` | sí · sí |
| C-051 | Chip de persona o tienda | reposo · marca de forma | persona ximena sobre `ui.surface.raised` | interfaz 3 | `color.person.ximena.light` / `color.neutral.50` | `color.person.ximena.dark` / `color.neutral.900` | sí · sí |
| C-052 | Chip de persona o tienda | reposo · marca de forma | persona karen sobre `ui.surface.raised` | interfaz 3 | `color.person.karen.light` / `color.neutral.50` | `color.person.karen.dark` / `color.neutral.900` | sí · sí |
| C-053 | Chip de persona o tienda | reposo · marca de forma | persona shared sobre `ui.surface.raised` | interfaz 3 | `color.person.shared.light` / `color.neutral.50` | `color.person.shared.dark` / `color.neutral.900` | sí · sí |
| C-054 | Chip de persona o tienda | reposo · marca de forma | tienda walmart sobre `ui.surface.raised` | interfaz 3 | `color.store.walmart.light` / `color.neutral.50` | `color.store.walmart.dark` / `color.neutral.900` | sí · sí |
| C-055 | Chip de persona o tienda | reposo · marca de forma | tienda bodega-aurrera sobre `ui.surface.raised` | interfaz 3 | `color.store.bodega-aurrera.light` / `color.neutral.50` | `color.store.bodega-aurrera.dark` / `color.neutral.900` | sí · sí |
| C-056 | Chip de persona o tienda | reposo · marca de forma | tienda sams-club sobre `ui.surface.raised` | interfaz 3 | `color.store.sams-club.light` / `color.neutral.50` | `color.store.sams-club.dark` / `color.neutral.900` | sí · sí |
| C-057 | Chip de persona o tienda | reposo · marca de forma | tienda market sobre `ui.surface.raised` | interfaz 3 | `color.store.market.light` / `color.neutral.50` | `color.store.market.dark` / `color.neutral.900` | sí · sí |
| C-058 | Chip de persona o tienda | reposo · nombre | `ui.text.default` sobre `ui.surface.raised` | texto 4.5 | `color.neutral.900` / `color.neutral.50` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-059 | Avatar | reposo · anillo sobre página | persona axel sobre `ui.surface.page` | interfaz 3 | `color.person.axel.light` / `color.neutral.50` | `color.person.axel.dark` / `color.neutral.950` | sí · sí |
| C-060 | Avatar | reposo · anillo sobre página | persona ximena sobre `ui.surface.page` | interfaz 3 | `color.person.ximena.light` / `color.neutral.50` | `color.person.ximena.dark` / `color.neutral.950` | sí · sí |
| C-061 | Avatar | reposo · anillo sobre página | persona karen sobre `ui.surface.page` | interfaz 3 | `color.person.karen.light` / `color.neutral.50` | `color.person.karen.dark` / `color.neutral.950` | sí · sí |
| C-062 | Avatar | reposo · anillo sobre página | persona shared sobre `ui.surface.page` | interfaz 3 | `color.person.shared.light` / `color.neutral.50` | `color.person.shared.dark` / `color.neutral.950` | sí · sí |
| C-063 | Avatar | reposo · inicial | `ui.text.default` sobre `ui.surface.hover` | texto 4.5 | `color.neutral.900` / `color.neutral.200` | `color.neutral.100` / `color.neutral.800` | sí · sí |
| C-064 | Chip de filtro | activo · texto | `ui.text.default` sobre `ui.surface.selected` | texto 4.5 | `color.neutral.900` / `color.brand.100` | `color.neutral.100` / `color.brand.900` | sí · sí |
| C-065 | Chip de filtro | activo · icono quitar | `ui.icon.default` sobre `ui.surface.selected` | interfaz 3 | `color.neutral.700` / `color.brand.100` | `color.neutral.300` / `color.brand.900` | no · no |
| C-066 | Insignia de rol | reposo · income | rol income sobre superficie de rol income | texto 4.5 | `color.role.income.light` / `color.role.income.surface-light` | `color.role.income.dark` / `color.role.income.surface-dark` | sí · sí |
| C-067 | Insignia de rol | reposo · expense | rol expense sobre superficie de rol expense | texto 4.5 | `color.role.expense.light` / `color.role.expense.surface-light` | `color.role.expense.dark` / `color.role.expense.surface-dark` | sí · sí |
| C-068 | Insignia de rol | reposo · savings | rol savings sobre superficie de rol savings | texto 4.5 | `color.role.savings.light` / `color.role.savings.surface-light` | `color.role.savings.dark` / `color.role.savings.surface-dark` | sí · sí |
| C-069 | Insignia de rol | reposo · debt | rol debt sobre superficie de rol debt | texto 4.5 | `color.role.debt.light` / `color.role.debt.surface-light` | `color.role.debt.dark` / `color.role.debt.surface-dark` | sí · sí |
| C-070 | Insignia de rol | reposo · warning | rol warning sobre superficie de rol warning | texto 4.5 | `color.role.warning.light` / `color.role.warning.surface-light` | `color.role.warning.dark` / `color.role.warning.surface-dark` | sí · sí |
| C-071 | Insignia de rol | reposo · error | rol error sobre superficie de rol error | texto 4.5 | `color.role.error.light` / `color.role.error.surface-light` | `color.role.error.dark` / `color.role.error.surface-dark` | sí · sí |
| C-072 | Insignia de rol | reposo · success | rol success sobre superficie de rol success | texto 4.5 | `color.role.success.light` / `color.role.success.surface-light` | `color.role.success.dark` / `color.role.success.surface-dark` | sí · sí |
| C-073 | Insignia de rol | reposo · info | rol info sobre superficie de rol info | texto 4.5 | `color.role.info.light` / `color.role.info.surface-light` | `color.role.info.dark` / `color.role.info.surface-dark` | sí · sí |
| C-074 | Insignia de acción crítica | reposo · texto | acción crítica sobre `ui.surface.raised` | texto 4.5 | `color.critical.light` / `color.neutral.50` | `color.critical.dark` / `color.neutral.900` | sí · sí |
| C-075 | Tabla | encabezado · texto e icono | `ui.text.secondary` sobre `ui.surface.page` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-076 | Tabla | reposo · celda | `ui.text.default` sobre `ui.surface.page` | texto 4.5 | `color.neutral.900` / `color.neutral.50` | `color.neutral.100` / `color.neutral.950` | sí · sí |
| C-077 | Tabla | reposo · cifra de ingreso | rol income sobre `ui.surface.page` | texto 4.5 | `color.role.income.light` / `color.neutral.50` | `color.role.income.dark` / `color.neutral.950` | sí · sí |
| C-078 | Tabla | reposo · cifra de gasto | rol expense sobre `ui.surface.page` | texto 4.5 | `color.role.expense.light` / `color.neutral.50` | `color.role.expense.dark` / `color.neutral.950` | sí · sí |
| C-079 | Tabla | reposo · cifra rebasada | rol warning sobre `ui.surface.page` | texto 4.5 | `color.role.warning.light` / `color.neutral.50` | `color.role.warning.dark` / `color.neutral.950` | sí · sí |
| C-080 | Tabla | hover · celda | `ui.text.default` sobre `ui.surface.hover-subtle` | texto 4.5 | `color.neutral.900` / `color.neutral.100` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-081 | Tabla | hover · cifra rebasada | rol warning sobre `ui.surface.hover-subtle` | texto 4.5 | `color.role.warning.light` / `color.neutral.100` | `color.role.warning.dark` / `color.neutral.900` | sí · sí |
| C-082 | Tabla | seleccionada · celda | `ui.text.default` sobre `ui.surface.selected` | texto 4.5 | `color.neutral.900` / `color.brand.100` | `color.neutral.100` / `color.brand.900` | sí · sí |
| C-083 | Tabla | seleccionada · cifra de ingreso | rol income sobre `ui.surface.selected` | texto 4.5 | `color.role.income.light` / `color.brand.100` | `color.role.income.dark` / `color.brand.900` | no · no |
| C-084 | Tabla | seleccionada · cifra de gasto | rol expense sobre `ui.surface.selected` | texto 4.5 | `color.role.expense.light` / `color.brand.100` | `color.role.expense.dark` / `color.brand.900` | no · no |
| C-085 | Tabla | seleccionada · texto secundario | `ui.text.secondary` sobre `ui.surface.selected` | texto 4.5 | `color.neutral.700` / `color.brand.100` | `color.neutral.300` / `color.brand.900` | no · no |
| C-086 | Tabla | error · borde de celda | `ui.border.error` sobre `ui.surface.page` | interfaz 3 | `color.role.error.light` / `color.neutral.50` | `color.role.error.dark` / `color.neutral.950` | sí · sí |
| C-087 | Tabla | divisor de filas | `ui.border.divider` sobre `ui.surface.page` | decorativo | `color.neutral.200` / `color.neutral.50` | `color.neutral.800` / `color.neutral.950` | no · no |
| C-088 | Tablero | tarjeta · título | `ui.text.default` sobre `ui.surface.card` | texto 4.5 | `color.neutral.900` / `color.neutral.100` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-089 | Tablero | tarjeta · propiedad | `ui.text.secondary` sobre `ui.surface.card` | texto 4.5 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-090 | Tablero | arrastre · destino | `ui.border.selected` sobre `ui.surface.page` | interfaz 3 | `color.brand.700` / `color.neutral.50` | `color.brand.300` / `color.neutral.950` | sí · sí |
| C-091 | Calendario (vista) | día de hoy · número | `ui.text.on-primary` sobre `ui.action.primary.bg` | texto 4.5 | `color.base.white` / `color.brand.700` | `color.brand.950` / `color.brand.300` | sí · sí |
| C-092 | Calendario (vista) | día · número | `ui.text.secondary` sobre `ui.surface.page` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-093 | Calendario (vista) | día de otro mes · número | `ui.text.disabled` sobre `ui.surface.page` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.400` / `color.neutral.950` | sí · sí |
| C-094 | Galería | portada sin imagen · icono | `ui.icon.default` sobre `ui.surface.hover` | interfaz 3 | `color.neutral.700` / `color.neutral.200` | `color.neutral.300` / `color.neutral.800` | no · no |
| C-095 | Tarjeta de métrica | etiqueta | `ui.text.secondary` sobre `ui.surface.card` | texto 4.5 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-096 | Tarjeta de métrica | cifra principal (24 px) | `ui.text.default` sobre `ui.surface.card` | texto grande 3 | `color.neutral.900` / `color.neutral.100` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-097 | Tarjeta de métrica | dato secundario de rol · income | rol income sobre `ui.surface.card` | texto 4.5 | `color.role.income.light` / `color.neutral.100` | `color.role.income.dark` / `color.neutral.900` | sí · sí |
| C-098 | Tarjeta de métrica | dato secundario de rol · expense | rol expense sobre `ui.surface.card` | texto 4.5 | `color.role.expense.light` / `color.neutral.100` | `color.role.expense.dark` / `color.neutral.900` | sí · sí |
| C-099 | Tarjeta de métrica | dato secundario de rol · savings | rol savings sobre `ui.surface.card` | texto 4.5 | `color.role.savings.light` / `color.neutral.100` | `color.role.savings.dark` / `color.neutral.900` | sí · sí |
| C-100 | Tarjeta de métrica | dato secundario de rol · debt | rol debt sobre `ui.surface.card` | texto 4.5 | `color.role.debt.light` / `color.neutral.100` | `color.role.debt.dark` / `color.neutral.900` | sí · sí |
| C-101 | Tarjeta de métrica | dato secundario de rol · warning | rol warning sobre `ui.surface.card` | texto 4.5 | `color.role.warning.light` / `color.neutral.100` | `color.role.warning.dark` / `color.neutral.900` | sí · sí |
| C-102 | Barra de avance | bajo el umbral · llenado sobre pista | `ui.viz.meter` sobre `ui.viz.track` | interfaz 3 | `color.neutral.700` / `color.neutral.200` | `color.neutral.300` / `color.neutral.800` | no · no |
| C-103 | Barra de avance | bajo el umbral · llenado sobre tarjeta | `ui.viz.meter` sobre `ui.surface.card` | interfaz 3 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-104 | Barra de avance | umbral alcanzado · llenado sobre pista | rol warning sobre `ui.viz.track` | interfaz 3 | `color.role.warning.light` / `color.neutral.200` | `color.role.warning.dark` / `color.neutral.800` | no · no |
| C-105 | Barra de avance | pista sobre tarjeta | `ui.viz.track` sobre `ui.surface.card` | decorativo | `color.neutral.200` / `color.neutral.100` | `color.neutral.800` / `color.neutral.900` | no · no |
| C-106 | Gráfica | serie 1 sobre tarjeta | serie 1 sobre `ui.surface.card` | interfaz 3 | `color.chart.1.light` / `color.neutral.100` | `color.chart.1.dark` / `color.neutral.900` | sí · sí |
| C-107 | Gráfica | serie 2 sobre tarjeta | serie 2 sobre `ui.surface.card` | interfaz 3 | `color.chart.2.light` / `color.neutral.100` | `color.chart.2.dark` / `color.neutral.900` | sí · sí |
| C-108 | Gráfica | serie 3 sobre tarjeta | serie 3 sobre `ui.surface.card` | interfaz 3 | `color.chart.3.light` / `color.neutral.100` | `color.chart.3.dark` / `color.neutral.900` | sí · sí |
| C-109 | Gráfica | serie 4 sobre tarjeta | serie 4 sobre `ui.surface.card` | interfaz 3 | `color.chart.4.light` / `color.neutral.100` | `color.chart.4.dark` / `color.neutral.900` | sí · sí |
| C-110 | Gráfica | serie 5 sobre tarjeta | serie 5 sobre `ui.surface.card` | interfaz 3 | `color.chart.5.light` / `color.neutral.100` | `color.chart.5.dark` / `color.neutral.900` | sí · sí |
| C-111 | Gráfica | serie 6 sobre tarjeta | serie 6 sobre `ui.surface.card` | interfaz 3 | `color.chart.6.light` / `color.neutral.100` | `color.chart.6.dark` / `color.neutral.900` | sí · sí |
| C-112 | Gráfica | serie 7 sobre tarjeta | serie 7 sobre `ui.surface.card` | interfaz 3 | `color.chart.7.light` / `color.neutral.100` | `color.chart.7.dark` / `color.neutral.900` | sí · sí |
| C-113 | Gráfica | serie 8 sobre tarjeta | serie 8 sobre `ui.surface.card` | interfaz 3 | `color.chart.8.light` / `color.neutral.100` | `color.chart.8.dark` / `color.neutral.900` | sí · sí |
| C-114 | Gráfica | serie «Otros» sobre tarjeta | `ui.viz.other` sobre `ui.surface.card` | interfaz 3 | `color.neutral.600` / `color.neutral.100` | `color.neutral.400` / `color.neutral.900` | sí · sí |
| C-115 | Gráfica | línea base | `ui.viz.baseline` sobre `ui.surface.card` | interfaz 3 | `color.neutral.600` / `color.neutral.100` | `color.neutral.400` / `color.neutral.900` | sí · sí |
| C-116 | Gráfica | retícula | `ui.viz.grid` sobre `ui.surface.card` | decorativo | `color.neutral.200` / `color.neutral.100` | `color.neutral.800` / `color.neutral.900` | no · no |
| C-117 | Gráfica | etiqueta de eje | `ui.text.secondary` sobre `ui.surface.card` | texto 4.5 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-118 | Gráfica | etiqueta directa y leyenda | `ui.text.default` sobre `ui.surface.card` | texto 4.5 | `color.neutral.900` / `color.neutral.100` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-119 | Gráfica | descripción emergente · texto | `ui.text.default` sobre `ui.surface.raised` | texto 4.5 | `color.neutral.900` / `color.neutral.50` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-120 | Gráfica | descripción emergente · texto secundario | `ui.text.secondary` sobre `ui.surface.raised` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-121 | Calendario de pagos | pagado · icono y texto | rol success sobre `ui.surface.card` | texto 4.5 | `color.role.success.light` / `color.neutral.100` | `color.role.success.dark` / `color.neutral.900` | sí · sí |
| C-122 | Calendario de pagos | por vencer o vencido · icono y texto | rol warning sobre `ui.surface.card` | texto 4.5 | `color.role.warning.light` / `color.neutral.100` | `color.role.warning.dark` / `color.neutral.900` | sí · sí |
| C-123 | Calendario de pagos | pendiente · monto | `ui.text.default` sobre `ui.surface.card` | texto 4.5 | `color.neutral.900` / `color.neutral.100` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-124 | Calendario de pagos | día seleccionado · monto | `ui.text.default` sobre `ui.surface.selected` | texto 4.5 | `color.neutral.900` / `color.brand.100` | `color.neutral.100` / `color.brand.900` | sí · sí |
| C-125 | Barra lateral | reposo · elemento | `ui.text.secondary` sobre `ui.surface.card` | texto 4.5 | `color.neutral.700` / `color.neutral.100` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-126 | Barra lateral | hover · elemento | `ui.text.default` sobre `ui.surface.hover` | texto 4.5 | `color.neutral.900` / `color.neutral.200` | `color.neutral.100` / `color.neutral.800` | sí · sí |
| C-127 | Barra lateral | activo · elemento | `ui.text.default` sobre `ui.surface.selected` | texto 4.5 | `color.neutral.900` / `color.brand.100` | `color.neutral.100` / `color.brand.900` | sí · sí |
| C-128 | Barra lateral | activo · icono | `ui.icon.selected` sobre `ui.surface.selected` | interfaz 3 | `color.brand.700` / `color.brand.100` | `color.brand.300` / `color.brand.900` | sí · sí |
| C-129 | Barra superior | migas | `ui.text.secondary` sobre `ui.surface.page` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-130 | Pestañas | activa · indicador | `ui.border.selected` sobre `ui.surface.page` | interfaz 3 | `color.brand.700` / `color.neutral.50` | `color.brand.300` / `color.neutral.950` | sí · sí |
| C-131 | Barra inferior | reposo · icono y nombre | `ui.text.secondary` sobre `ui.surface.raised` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-132 | Barra inferior | activo · icono | `ui.icon.selected` sobre `ui.surface.selected` | interfaz 3 | `color.brand.700` / `color.brand.100` | `color.brand.300` / `color.brand.900` | sí · sí |
| C-133 | Barra inferior | activo · nombre | `ui.text.default` sobre `ui.surface.raised` | texto 4.5 | `color.neutral.900` / `color.neutral.50` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-134 | Barra inferior | «+» · relleno sobre barra | `ui.action.primary.bg` sobre `ui.surface.raised` | interfaz 3 | `color.brand.700` / `color.neutral.50` | `color.brand.300` / `color.neutral.900` | sí · sí |
| C-135 | Menú «/» y búsqueda | reposo · comando | `ui.text.default` sobre `ui.surface.raised` | texto 4.5 | `color.neutral.900` / `color.neutral.50` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-136 | Menú «/» y búsqueda | reposo · alias y grupo | `ui.text.secondary` sobre `ui.surface.raised` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.900` | sí · sí |
| C-137 | Menú «/» y búsqueda | resaltado · comando | `ui.text.default` sobre `ui.surface.hover` | texto 4.5 | `color.neutral.900` / `color.neutral.200` | `color.neutral.100` / `color.neutral.800` | sí · sí |
| C-138 | Menú «/» y búsqueda | resaltado · alias y atajo | `ui.text.secondary` sobre `ui.surface.hover` | texto 4.5 | `color.neutral.700` / `color.neutral.200` | `color.neutral.300` / `color.neutral.800` | no · no |
| C-139 | Diálogo de acción crítica | icono de encabezado | acción crítica sobre `ui.surface.raised` | interfaz 3 | `color.critical.light` / `color.neutral.50` | `color.critical.dark` / `color.neutral.900` | sí · sí |
| C-140 | Ajustes | fila · descripción y valor inicial | `ui.text.secondary` sobre `ui.surface.page` | texto 4.5 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-141 | Ajustes | guardado · palomita | rol success sobre `ui.surface.page` | interfaz 3 | `color.role.success.light` / `color.neutral.50` | `color.role.success.dark` / `color.neutral.950` | sí · sí |
| C-142 | Estado vacío | icono | `ui.icon.default` sobre `ui.surface.page` | interfaz 3 | `color.neutral.700` / `color.neutral.50` | `color.neutral.300` / `color.neutral.950` | sí · sí |
| C-143 | Esqueleto de carga | sobre página | `ui.skeleton.on-page` sobre `ui.surface.page` | decorativo | `color.neutral.100` / `color.neutral.50` | `color.neutral.900` / `color.neutral.950` | no · no |
| C-144 | Esqueleto de carga | sobre tarjeta | `ui.skeleton.on-card` sobre `ui.surface.card` | decorativo | `color.neutral.200` / `color.neutral.100` | `color.neutral.800` / `color.neutral.900` | no · no |
| C-145 | Aviso en página (franja) | info · icono | rol info sobre superficie de rol info | interfaz 3 | `color.role.info.light` / `color.role.info.surface-light` | `color.role.info.dark` / `color.role.info.surface-dark` | sí · sí |
| C-146 | Aviso en página (franja) | info · texto | `ui.text.default` sobre superficie de rol info | texto 4.5 | `color.neutral.900` / `color.role.info.surface-light` | `color.neutral.100` / `color.role.info.surface-dark` | no · no |
| C-147 | Aviso en página (franja) | warning · icono | rol warning sobre superficie de rol warning | interfaz 3 | `color.role.warning.light` / `color.role.warning.surface-light` | `color.role.warning.dark` / `color.role.warning.surface-dark` | sí · sí |
| C-148 | Aviso en página (franja) | warning · texto | `ui.text.default` sobre superficie de rol warning | texto 4.5 | `color.neutral.900` / `color.role.warning.surface-light` | `color.neutral.100` / `color.role.warning.surface-dark` | no · no |
| C-149 | Aviso en página (franja) | error · icono | rol error sobre superficie de rol error | interfaz 3 | `color.role.error.light` / `color.role.error.surface-light` | `color.role.error.dark` / `color.role.error.surface-dark` | sí · sí |
| C-150 | Aviso en página (franja) | error · texto | `ui.text.default` sobre superficie de rol error | texto 4.5 | `color.neutral.900` / `color.role.error.surface-light` | `color.neutral.100` / `color.role.error.surface-dark` | no · no |
| C-151 | Aviso en página (franja) | success · icono | rol success sobre superficie de rol success | interfaz 3 | `color.role.success.light` / `color.role.success.surface-light` | `color.role.success.dark` / `color.role.success.surface-dark` | sí · sí |
| C-152 | Aviso en página (franja) | success · texto | `ui.text.default` sobre superficie de rol success | texto 4.5 | `color.neutral.900` / `color.role.success.surface-light` | `color.neutral.100` / `color.role.success.surface-dark` | no · no |
| C-153 | Aviso breve | texto | `ui.text.default` sobre `ui.surface.raised` | texto 4.5 | `color.neutral.900` / `color.neutral.50` | `color.neutral.100` / `color.neutral.900` | sí · sí |
| C-154 | Aviso breve | icono · success | rol success sobre `ui.surface.raised` | interfaz 3 | `color.role.success.light` / `color.neutral.50` | `color.role.success.dark` / `color.neutral.900` | sí · sí |
| C-155 | Aviso breve | icono · error | rol error sobre `ui.surface.raised` | interfaz 3 | `color.role.error.light` / `color.neutral.50` | `color.role.error.dark` / `color.neutral.900` | sí · sí |
| C-156 | Aviso breve | icono · warning | rol warning sobre `ui.surface.raised` | interfaz 3 | `color.role.warning.light` / `color.neutral.50` | `color.role.warning.dark` / `color.neutral.900` | sí · sí |
| C-157 | Aviso breve | icono · info | rol info sobre `ui.surface.raised` | interfaz 3 | `color.role.info.light` / `color.neutral.50` | `color.role.info.dark` / `color.neutral.900` | sí · sí |
| C-158 | Lo que flota | borde | `ui.border.divider` sobre `ui.surface.raised` | decorativo | `color.neutral.200` / `color.neutral.50` | `color.neutral.800` / `color.neutral.900` | no · no |

### 10.6 Medidas de componente (propuestas de U2b-2 para `C2-tokens.json`)

| Token | Valor | Uso |
|---|---|---|
| `size.sidebar.item` | 32 px (fino; 44 px con puntero grueso) | elemento de la barra lateral y de la lista de Ajustes (§7.9, §7.13) |
| `size.chip.height` · `size.badge.height` | 24 px · 20 px | chip de persona, tienda o filtro; insignia de rol (§7.4) |
| `size.avatar` | 24 px | avatar con anillo de 2 px (§7.4) |
| `size.mark.category` | 10 px | marca de forma de persona o tienda (§7.4) |
| `size.checkbox` · `size.switch.width` · `size.switch.height` | 20 px · 36 px · 20 px | casilla e interruptor, dentro de 44 px (§7.3) |
| `size.meter.height` | 8 px | barra de avance (§7.6) |
| `size.bottombar.add` · `size.bottombar.indicator-width` · `size.bottombar.indicator-height` | 40 px · 56 px · 32 px | «+» y píldora del destino activo (§7.10) |
| `size.sheet.handle-width` · `size.sheet.handle-height` | 32 px · 4 px | asa de la hoja inferior (§7.12) |
| `layout.board.column-width` | 272 px | columna del tablero (§7.5.3) |
| `size.chart.height-compact` · `size.chart.height` | 200 px · 280 px | alto de gráfica en < 768 y desde 768 (§7.7) |
| `size.calendar.cell` · `size.calendar.cell-wide` | 96 px · 120 px | celda del mes a 1024 y a 1440 (§7.8) |
| `size.agenda.row` | 56 px | fila de agenda y de pagos (§7.8) |
| `size.viz.bar-max` · `size.viz.bar-radius` · `size.viz.line` · `size.viz.marker` · `size.viz.gap` | 24 · 4 · 2 · 10 · 2 px | marcas de §8.1 |
| `layout.toast.max-width` | 480 px (alias de `layout.dialog.width.short`) | aviso breve (§7.15) |

## Supuestos

Secciones 1 a 6 (U2b-1-design); U2b-2 agrega los suyos.

1. La ruta de los archivos estáticos en la PWA (`/fonts/`, `/icons/`) es una propuesta: la fija WF-3 al armar el frontend; los archivos fuente están en `docs/design/fuentes/`.
2. En la PWA los iconos se usan con el componente de iconos de Lucide para React que trae shadcn/ui, con la misma versión que el sprite (1.55.0) y ancho de trazo absoluto de 1.5; la dependencia exacta la decide WF-3. El sprite sirve a documentos, Storybook y páginas estáticas.
3. Los nombres de token que no son de color (`font.*`, `space.*`, `radius.*`, `shadow.*`, `size.*`, `layout.*`, `motion.*`, `icon.*`, `border.*`, `focus.*`, `overlay.*`) son propuestas de U2b-1 listas para que la sección 10 los lleve a `C2-tokens.json` y `tokens.css`; no cambian ningún valor de color.

Secciones 7 a 10 (U2b-2-design):

4. **Área apilada devengado contra efectivo** (§8.3): las capas son efectivo acumulado y devengado sin pagar (devengado − efectivo), de modo que la cima es el devengado acumulado; un pago adelantado deja la capa superior en cero con texto. Es la lectura de diseño de plan §4.1 punto 8; no cambia ninguna cifra ni el contrato de datos.
5. **Pago vencido en advertencia y no en error** (§7.8): error queda para validación y fallas (§2.4); un vencimiento es un límite superado, como el presupuesto rebasado (§3.3).
6. **Capa semántica `ui.*`** (§7.1, §10): 41 alias de primitivos de C2, sin valores nuevos; el prefijo `ui` evita confundirlos con los primitivos `color.*` que comprueba `verificar.py`. `ui.surface.hover-subtle` existe porque las cifras de rol sobre `ui.surface.hover` bajan de 4.5:1 (gasto en oscuro 4.07; ingreso 4.30 en claro y 4.37 en oscuro; cálculo previo de U2b-2): filas y edición en el lugar usan el hover sutil y las tarjetas solo cambian el borde.
7. Las medidas de componente de §10.6, el nombre del archivo del tema de Tailwind (quedó `docs/contratos/tailwind.theme.css`, §10.1) y los nombres `breakpoint.*` son propuestas listas para U2c; no cambian ningún valor de §3–§6.
8. Patrón de relleno siempre presente desde la serie 2 (K27) y cifras tabulares en todo dinero (§3.3), aunque la guía de visualización los plantea como opcional y como proporcionales en cifras grandes: rigen las decisiones.
9. Las constantes de interfaz (aviso breve 5 s de U1 G-07; «Guardado» 2 s en Ajustes; indicador de carga desde 400 ms; mensaje de lentitud a los 10 s de U1 §2.4) se tratan como diseño y no como parámetros operativos de K44; si Axel los quiere configurables, irían en Apariencia.
10. **Barra inferior** (§4.6, §9.3; corrección de la revisión U2 r1): «56 × ≥ 64 px» se lee como alto 56 px, el mismo de `layout.bottombar.height` (un destino no puede ser más alto que su barra), y ancho mínimo 64 px (5 destinos en 360 px dan 72 px cada uno). No se agrega ningún valor: son los dos números que ya traía este documento.
