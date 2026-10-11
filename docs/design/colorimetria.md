# Colorimetría · candidato C «Índigo 280°»

> Documento generado por `docs/design/tools/colorimetria.py` (Python 3, solo biblioteca estándar). No se edita a mano: se cambia la regla en el script y se regenera. Fase U2a del flujo WF-0, 2026-10-10. Es la base de color de `DESIGN.md` (U2, WF-1) y la parte de color del contrato C2 (`docs/contratos/C2-colores.json`); el resto de los tokens de C2 los hace U2b en WF-1. Queda sujeto a la puerta 2 de aprobación. Ronda 1 de revisión (2026-10-11): corrige el hallazgo H-03 (la serie 6 de gráficas usaba el matiz reservado del eje opuesto; §3 y §8.8). **Versión 1.1 (puerta 1, 2026-10-11):** aplica las respuestas K26 y K27 de Axel y el hallazgo R3-04; lo que cambia, con valores anteriores y nuevos, está en «Cambios 1.1».

## Índice

0. Cambios 1.1 (puerta 1): gasto, Karen, Bodega Aurrera, marca en pares cercanos y paleta de gráficas
1. Qué se genera y cómo regenerarlo
2. Método
3. Matiz base, eje opuesto y armonía complementaria dividida
4. Regla 60-30-10
5. Escalas de 11 pasos
6. Roles semánticos
7. Paleta categórica: personas y tiendas
8. Paleta de gráficas (8 series)
9. Reglas de uso por componente
10. Contraste (resumen de `contraste.csv`)
11. Tokens (`C2-colores.json`)
12. Hallazgos, decisiones derivadas del cálculo y bloqueos

Anexo A. Tokens: ruta → hex (lista que verifica el script)

## Cambios 1.1 (puerta 1, 2026-10-11)

Aplica las respuestas de Axel en la puerta 1 (`docs/plan/decisiones.md` §K): **K26** (gasto pasa a terracota 40–45° y error se queda en 20°; Karen y Bodega Aurrera, que estaban en 50°, se mueven en la paleta categórica; la marca entra al barrido de pares cercanos, hallazgo R3-04 de la revisión 3) y **K27** (paleta de gráficas: ΔL ≥ 0.10 entre series consecutivas con forma, patrón y etiqueta de refuerzo; todos los pares se reportan como información). Quedan resueltos los bloqueos `U2-gasto-error` y `criterios-deltal`. No cambian el método (L por paso, perfil de croma, umbrales), la marca, el neutro, los demás roles ni las 8 series de gráficas. Todo lo regenera el script y su verificación final lo comprueba.

### Parámetros que cambian

| Elemento | 1.0: matiz · croma pedido · L claro / oscuro | 1.1 | Origen |
|---|---|---|---|
| Gasto (`color.role.expense`) | matiz 25° | matiz 45° | K26 |
| Karen (`color.person.karen`) | 50° · croma 0.12 · L 0.640 / 0.725 | 333° · croma 0.12 · L 0.642 / 0.725 | K26; **PROVISIONAL** frente a Ximena, bloqueo `U2a-rosas-ximena-bodega` |
| Bodega Aurrera (`color.store.bodega-aurrera`) | 50° · croma 0.12 · L 0.410 / 0.830 | 356° · croma 0.11 · L 0.410 / 0.829 | K26; **PROVISIONAL** frente a Ximena, bloqueo `U2a-rosas-ximena-bodega` |
| Walmart (`color.store.walmart`) | 262° · croma 0.12 · L 0.515 / 0.720 | 262° · croma 0.12 · L 0.515 / 0.716 | **PROVISIONAL**, bloqueo `U2a-walmart-oscuro` |

### Tokens que cambian (medidos en el hex)

| Token | 1.0: hex · matiz · L · croma | 1.1: hex · matiz · L · croma |
|---|---|---|
| `color.role.expense.light` | `#A44541` · 24.9° · L 0.515 · C 0.127 | `#A04B20` · 44.7° · L 0.514 · C 0.126 |
| `color.role.expense.dark` | `#FC9A92` · 25.1° · L 0.785 · C 0.119 | `#F8A078` · 45.0° · L 0.785 · C 0.118 |
| `color.role.expense.surface-light` | `#FFEEEC` · 25.6° · L 0.962 · C 0.019 | `#FFEFE8` · 45.3° · L 0.963 · C 0.020 |
| `color.role.expense.surface-dark` | `#5A2220` · 24.7° · L 0.334 · C 0.083 | `#59260B` · 45.2° · L 0.336 · C 0.084 |
| `color.role.expense.50` | `#FFF8F8` · 17.3° · L 0.985 · C 0.007 | `#FFF9F6` · 48.7° · L 0.986 · C 0.008 |
| `color.role.expense.100` | `#FFEEEC` · 25.6° · L 0.962 · C 0.019 | `#FFEFE8` · 45.3° · L 0.963 · C 0.020 |
| `color.role.expense.200` | `#FFDDDA` · 23.8° · L 0.925 · C 0.038 | `#FFDED0` · 44.4° · L 0.924 · C 0.041 |
| `color.role.expense.300` | `#FFC1BB` · 25.0° · L 0.865 · C 0.072 | `#FFC3A9` · 44.5° · L 0.865 · C 0.078 |
| `color.role.expense.400` | `#FC9A92` · 25.1° · L 0.785 · C 0.119 | `#F8A078` · 45.0° · L 0.785 · C 0.118 |
| `color.role.expense.500` | `#E77770` · 25.0° · L 0.696 · C 0.140 | `#E37E4E` · 45.1° · L 0.696 · C 0.140 |
| `color.role.expense.600` | `#C85B56` · 24.8° · L 0.605 · C 0.140 | `#C46231` · 45.1° · L 0.605 · C 0.140 |
| `color.role.expense.700` | `#A44541` · 24.9° · L 0.515 · C 0.127 | `#A04B20` · 44.7° · L 0.514 · C 0.126 |
| `color.role.expense.800` | `#7E3330` · 24.9° · L 0.425 · C 0.105 | `#7C3815` · 44.9° · L 0.426 · C 0.105 |
| `color.role.expense.900` | `#5A2220` · 24.7° · L 0.334 · C 0.083 | `#59260B` · 45.2° · L 0.336 · C 0.084 |
| `color.role.expense.950` | `#3C1413` · 24.1° · L 0.255 · C 0.064 | `#3B1705` · 45.2° · L 0.256 · C 0.064 |
| `color.person.karen.light` | `#C67544` · 50.4° · L 0.641 · C 0.120 | `#B870AC` · 333.2° · L 0.643 · C 0.120 |
| `color.person.karen.dark` | `#E28F5E` · 50.5° · L 0.725 · C 0.120 | `#D389C6` · 333.3° · L 0.725 · C 0.120 |
| `color.store.bodega-aurrera.light` | `#773400` · 49.6° · L 0.410 · C 0.109 | `#762C4C` · 356.1° · L 0.410 · C 0.110 |
| `color.store.bodega-aurrera.dark` | `#FFB389` · 49.7° · L 0.829 · C 0.105 | `#FFAAC9` · 356.2° · L 0.829 · C 0.106 |
| `color.store.walmart.dark` | `#7BA4F0` · 262.1° · L 0.720 · C 0.120 | `#7AA3EE` · 261.9° · L 0.717 · C 0.119 |

Los otros 170 tokens no cambian de valor. Tablas regeneradas: §5.3 (escala de gasto), §6.1, §6.5, §6.6, §7.1 a §7.5, §7.7, §10, §12 y el Anexo A, además de `contraste.csv`, `C2-colores.json`, `paleta.html` y `paleta.png`.

### 1. Gasto: el grado de 40° a 45° más separado (K26)

Para cada grado entero del rango el script construye la escala de gasto con el mismo método (L por paso, perfil de croma, croma máx. 0.14), elige el paso de texto con la misma regla de 4.5:1 y mide en el hex el Δ matiz a error (20°) y a advertencia (75°), el ΔE normal y el ΔE mínimo con las tres simulaciones de daltonismo frente a esos dos roles, y la holgura mínima frente a las 8 categóricas (máx(Δ matiz / 20°, ΔL / 0.15); cumple con ≥ 1). Gana el mayor Δ matiz mínimo; a igualdad, el mayor ΔE normal y después el mayor ΔE con daltonismo.

| Matiz | Claro: paso · hex · matiz medido | Oscuro: paso · hex · matiz medido | Δ matiz a error (mín.) | Δ matiz a advertencia (mín.) | ΔE normal mín. | ΔE daltonismo mín. | Holgura frente a categóricas |
|---|---|---|---|---|---|---|---|
| 40° | 700 · `#A24A29` · 40.3° | 400 · `#FA9E7E` · 39.9° | 19.8° | 35.1° | 4.1 | 0.7 | 1.14 (Mercado) |
| 41° | 700 · `#A24A27` · 41.1° | 400 · `#FA9E7D` · 40.4° | 20.4° | 34.3° | 4.2 | 0.8 | 1.14 (Mercado) |
| 42° | 700 · `#A14A25` · 42.1° | 400 · `#F99F7B` · 42.5° | 22.4° | 33.3° | 4.6 | 0.9 | 1.14 (Mercado) |
| 43° | 700 · `#A14B24` · 43.1° | 400 · `#F99F7A` · 43.0° | 22.9° | 32.3° | 4.7 | 0.9 | 1.14 (Mercado) |
| 44° | 700 · `#A14B22` · 43.8° | 400 · `#F99F79` · 43.5° | 23.4° | 31.6° | 4.8 | 0.9 | 1.13 (Mercado) |
| **45°** | 700 · `#A04B20` · 44.7° | 400 · `#F8A078` · 45.0° | 24.9° | 30.7° | 5.1 | 1.0 | 1.14 (Mercado) |

**Se elige 45°.** Dentro del rango, subir el matiz aleja a gasto de error y lo acerca a advertencia, pero advertencia queda siempre más lejos que error, así que el Δ matiz mínimo crece con el matiz y el mejor grado es el extremo. En 45° gasto queda a 24.9° de error y a 30.7° de advertencia, con ΔE normal mínimo 5.1 y ΔE mínimo con daltonismo 1.0, los mayores del rango. Frente a las categóricas la holgura va de 1.13 a 1.14 en todo el rango y la fija Mercado, así que no decide. El paso de texto sigue siendo 700 en claro y 400 en oscuro, con 4.5:1 (§6.1).

### 2. Karen y Bodega Aurrera (K26)

Con gasto en 45° el rango cálido queda cerrado: entre error 20°, gasto 45° y advertencia 75° no hay ningún matiz a ≥ 20° de los tres, y de 80° a 120° está la zona reservada de acción crítica. Karen ocupa la L más alta de las personas en claro (la última que da 3:1) y Bodega Aurrera la L entre Mercado y Walmart; con esas L no llegan a ΔL 0.15 frente al texto de los roles (L 0.515 en claro y 0.785 en oscuro) en ningún tema, así que tienen que separarse por matiz de todos los roles, de acción crítica y de la marca. El script recorre los 360 matices enteros con la L y el croma finales de cada una y conserva los que cumplen todas las reglas de §7 (≥ 20° o ΔL ≥ 0.15 frente a roles, acción crítica y marca; contraste ≥ 3:1 sobre fondo y tarjeta; croma ≥ 0.10; fuera de la zona reservada; ΔL ≥ 0.10 en L OKLab y en gris frente a su grupo):

| Color | L claro / oscuro · croma pedido | Matices que cumplen todo | Con Walmart en L 0.720 (1.0) |
|---|---|---|---|
| Karen | 0.642 / 0.725 · 0.12 | 331°–359° | igual (otro grupo) |
| Bodega Aurrera | 0.410 / 0.829 · 0.11 | 331°–335°, 338°–341°, 344°–348°, 352°, 355°–359° | ninguna; con L 0.830 y croma 0.12 de la 1.0: ninguna |

Todas las ventanas quedan entre deuda (310°) y error (20° = 380°), y en oscuro Karen y Bodega Aurrera tienen ΔL 0.104 (< 0.15), así que también necesitan ≥ 20° de matiz entre sí. Se reparten los 70° que hay de deuda a error en tres huecos con el mayor hueco mínimo; a igualdad, Karen lo más lejos posible de Ximena (350°, su mismo grupo) y después la mayor holgura frente a roles, acción crítica y marca. Resultado: **Karen 333° y Bodega Aurrera 356°** (huecos de 23°, 23° y 24°). Los cinco mejores repartos:

| Karen | Bodega Aurrera | Huecos deuda → · → · → error | Hueco mínimo | Δ Karen–Ximena | Holgura mín. frente a roles y marca |
|---|---|---|---|---|---|
| 333° | 356° | 23° · 23° · 24° | 23° | 17° | 1.15 |
| 333° | 357° | 23° · 24° · 23° | 23° | 17° | 1.11 |
| 334° | 357° | 24° · 23° · 23° | 23° | 16° | 1.11 |
| 357° | 334° | 24° · 23° · 23° | 23° | 7° | 1.12 |
| 357° | 333° | 23° · 24° · 23° | 23° | 7° | 1.12 |

Karen queda a 17° de Ximena: dentro de las personas se distinguen por ΔL ≥ 0.10 en ambos temas (§7.3), por el orden fijo de la leyenda y por la etiqueta directa. Bodega Aurrera queda cerca de Ximena, pero son de grupos distintos y nunca comparten gráfica ni leyenda (§7.7).

**Frente a Ximena (ronda 1 de revisión, H-01; PROVISIONAL).** En oscuro, Bodega Aurrera (`#FFAAC9`) y Ximena (`#FFACD3`) quedan a 6.2°, ΔL 0.006 y **ΔE 1.3**: a simple vista son el mismo color. En claro quedan a 5.8°, ΔL 0.130 y ΔE 13.1. Karen y Ximena quedan a 17.1° y 16.7°, con ΔL 0.103 y 0.111 y ΔE 10.9 y 11.6 (claro y oscuro). Bodega Aurrera–Ximena es una colisión nueva de la 1.1: en la 1.0, el par persona–tienda más parecido era Karen–Bodega Aurrera en oscuro, con ΔE 10.5. Con la L y el croma de cada una, ningún matiz de sus ventanas (tabla de arriba) separa a Karen ni a Bodega Aurrera de Ximena por ≥ 20° o ΔL ≥ 0.15 en ambos temas; el script lo comprueba. §7.7 declara el refuerzo, y U2-08 lo admite. Pero U2-02 y la respuesta de Axel en `U2-gasto-error` piden moverlas «para conservar separación ≥ 20° o ΔL ≥ 0.15» frente a las demás entradas. Como el criterio admite dos lecturas, el bloqueo `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md` lleva las opciones para Axel. Mientras decide, van **PROVISIONAL** los 4 tokens que quedan bajo la regla frente a Ximena (`color.person.karen.dark`, `color.person.karen.light`, `color.store.bodega-aurrera.dark`, `color.store.bodega-aurrera.light`), con la nota en su descripción en `C2-colores.json`.

### 3. Walmart en oscuro (PROVISIONAL)

En oscuro las cuatro tiendas se escalonan con ΔL ≥ 0.10 en L OKLab y en gris entre Mercado (L 0.615, contraste 3.06, el mínimo del grupo) y Sam's Club (L 0.930, croma 0.103, cerca del piso de 0.10). El rosa pierde luminancia en gris (en Bodega Aurrera, gris = L − 0.014), y con Walmart en L 0.720 ningún matiz entero cumple para Bodega Aurrera (última columna de la tabla del punto 2). Walmart baja a L 0.716 (`#7BA4F0` → `#7AA3EE`, ΔE 0.4, imperceptible) y Bodega Aurrera usa croma 0.11 y L 0.829 en oscuro; el escalonamiento de las tiendas en oscuro queda en ΔL mínimo 0.101 (OKLab) y 0.101 (gris). Como es un token que K26 no nombra, va **PROVISIONAL** con el bloqueo `docs/estado/bloqueos/U2a-walmart-oscuro.md` (opciones para Axel) y la nota en la descripción del token en `C2-colores.json`.

### 4. La marca en el barrido de pares cercanos (R3-04)

§7.7 incluye ahora la marca en su paso de texto (700 en claro, 400 en oscuro). Ningún rol ni acción crítica queda cerca de la marca (el script lo verifica); quedan cerca estas categóricas, que se refuerzan con forma y etiqueta: la marca solo vive en controles (botón primario, enlace subrayado, anillo de foco, elemento activo) y la categórica siempre va como punto, anillo de avatar o marca de serie junto a su nombre en texto neutro.

| Tema | Color A | Color B | Δ matiz | ΔL |
|---|---|---|---|---|
| claro | Conjunto (persona) | Marca 700 (marca) | 0.9° | 0.094 |
| claro | Walmart (tienda) | Marca 700 (marca) | 18.0° | 0.000 |
| oscuro | Walmart (tienda) | Marca 400 (marca) | 18.0° | 0.069 |

### 5. Paleta de gráficas (K27)

No cambia ningún valor de las 8 series. K27 fija la lectura del criterio: el umbral duro es ΔL ≥ 0.10 entre series consecutivas en el orden de uso (mínimo en gris y en las tres simulaciones: 0.162 en claro y 0.198 en oscuro, §8.3) y los 28 pares se reportan como información (ΔL mínimo en gris 0.002 en claro y 0.001 en oscuro, §8.5). El color nunca va solo: cada serie lleva forma de marcador, patrón de trazo o relleno y etiqueta directa. Las tablas de simulación de daltonismo de las series (§8.4) no cambian porque las series no cambian; los ΔE con daltonismo de las categóricas (§7.3) se recalculan con Karen y Bodega Aurrera nuevas.

## 1. Qué se genera y cómo regenerarlo

| Archivo | Contenido |
|---|---|
| `docs/design/colorimetria.md` | este documento |
| `docs/design/contraste.csv` | 202 pares de contraste (102 de texto, 100 de interfaz y gráficas), claro y oscuro |
| `docs/contratos/C2-colores.json` | 190 tokens DTCG de color |
| `docs/design/paleta.html` | muestrario estático (sin CDN, sin fuentes externas, sin JavaScript) |
| `docs/design/paleta.png` | captura 1400×3000 con `chromium --headless=new` |

Comando: `python3 docs/design/tools/colorimetria.py` (agrega `--sin-png` para omitir la captura y `--solo-verificar` para revisar los archivos existentes sin regenerarlos). La salida es determinista: dos ejecuciones dan archivos idénticos. Al final el script verifica, leyendo los archivos ya escritos, que ningún par del CSV quede bajo su umbral, que los hex del JSON coincidan con el Anexo A de este documento, que la paleta de gráficas mantenga ΔL ≥ 0.10 entre series consecutivas en gris y en las tres simulaciones (con las dos variantes de §8.6) y que sus tres primeras series pasen la prueba de todos los pares (§8.5), que las categóricas cumplan su escalonamiento y su separación de los roles, que ningún token fuera de `color.critical` use la zona reservada del eje opuesto (§3), que marca y neutro reproduzcan los valores de partida del plan y, desde la versión 1.1, que gasto esté en 40–45° a ≥ 20° de error, que la marca no quede cerca de ningún rol y que Karen y Bodega Aurrera cumplan la regla de separación entre sí; si algo falla, termina con código 1.

## 2. Método

Resume y concreta el plan §4.2. Cada color se define en OKLCH (luminosidad L, croma C, matiz h en grados) y se convierte a sRGB de 8 bits.

### 2.1 Conversión

OKLCH → OKLab (a = C·cos h, b = C·sin h) → LMS no lineal (inversa de M2) → LMS (cubo) → sRGB lineal (inversa de M1) → sRGB codificado (función de transferencia de IEC 61966-2-1: 12.92·c si c ≤ 0.0031308; 1.055·c^(1/2.4) − 0.055 en otro caso) → hex con redondeo ROUND_HALF_UP por canal. Las matrices M1 y M2 son las publicadas por Björn Ottosson en «A perceptual color space for image processing» (bottosson.github.io/posts/oklab). La conversión inversa (hex → OKLab) se usa para medir.

### 2.2 Luminosidad fija por paso

| Paso | 50 | 100 | 200 | 300 | 400 | 500 | 600 | 700 | 800 | 900 | 950 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| L | 0.985 | 0.962 | 0.925 | 0.865 | 0.785 | 0.695 | 0.605 | 0.515 | 0.425 | 0.335 | 0.255 |

Todas las escalas (marca, neutro, roles y acción crítica) usan la misma L en cada paso, así que dos colores del mismo paso tienen la misma luminosidad percibida aunque cambie el matiz.

### 2.3 Perfil de croma por paso

| Paso | 50 | 100 | 200 | 300 | 400 | 500 | 600 | 700 | 800 | 900 | 950 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| Fracción del croma máximo | 0.10 | 0.25 | 0.40 | 0.65 | 0.85 | 1.00 | 1.00 | 0.90 | 0.75 | 0.60 | 0.45 |
| Intervalo que reproduce el plan | 0.064–0.180 | 0.160–0.262 | 0.326–0.416 | 0.624–0.746 | 0.846–0.852 | 0.994–1.000 | 0.992–1.004 | 0.890–0.906 | 0.750–0.764 | 0.594–0.610 | 0.438–0.456 |

Croma máximo por familia: marca 0.110, neutro tintado 0.012, roles y acción crítica 0.14 (plan §4.2: «saturación apagada»). El plan define el perfil solo por sus extremos y su pico (0.10 → 1.0 → 0.45); la tabla explícita de arriba reproduce los **22 de 22** valores de partida del plan §4.3 (marca y neutro del candidato C). La segunda fila muestra, por paso, el intervalo de fracciones que también los reproduce: en los pasos 50–300 la marca queda recortada por el gamut y el neutro se pierde en la cuantización, así que esos pasos solo quedan acotados por intervalo y se eligieron 0.10, 0.25, 0.40 y 0.65 dentro de él. Esto importa en las escalas de rol, cuyos pasos claros no siempre se recortan.

### 2.4 Recorte al gamut sRGB

Si el color pedido cae fuera de sRGB, se busca por bisección (60 iteraciones) el mayor croma que cabe, manteniendo L y h; un color está dentro si sus tres canales lineales quedan en [0, 1] con tolerancia 1e-9. Por eso las tablas muestran «croma pedido» y «croma final». El matiz y la luminosidad nunca se alteran para encajar.

### 2.5 Contraste (WCAG 2.x)

Luminancia relativa Y = 0.2126·R + 0.7152·G + 0.0722·B con R, G, B lineales (umbral 0.04045 de IEC 61966-2-1, el que cita la nota de WCAG 2.2; con valores de 8 bits da lo mismo que 0.03928). Razón = (Y₁ + 0.05) / (Y₂ + 0.05) con Y₁ ≥ Y₂. Umbrales: texto ≥ 4.5:1 (criterio 1.4.3, incluye texto deshabilitado por regla del proyecto) e interfaz, bordes de campo, foco, series y marcadores categóricos ≥ 3:1 (criterio 1.4.11). Las razones se truncan a dos decimales (nunca se redondean hacia arriba) y el cumplimiento se evalúa con el valor sin truncar.

### 2.6 Simulación de daltonismo

Modelo de Machado, Oliveira y Fernandes (2009), «A Physiologically-based Model for Simulation of Color Vision Deficiency», IEEE TVCG 15(6):1291-1298, doi:10.1109/TVCG.2009.113. Se usan las matrices de severidad 1.0 (dicromacia) de la Tabla 1 publicada por los autores en inf.ufrgs.br/~oliveira/pubs_files/CVD_Simulation/CVD_Simulation.html (consultada el 2026-10-10; coinciden con las que trae culori 4.0.2):

| Tipo | Fila R | Fila G | Fila B |
|---|---|---|---|
| Protanopia | 0.152286 1.052583 -0.204868 | 0.114503 0.786281 0.099216 | -0.003882 -0.048116 1.051998 |
| Deuteranopia | 0.367322 0.860646 -0.227968 | 0.280085 0.672501 0.047413 | -0.011820 0.042940 0.968881 |
| Tritanopia | 1.255528 -0.076749 -0.178779 | -0.078411 0.930809 0.147602 | 0.004733 0.691367 0.303900 |

Se aplican sobre sRGB **lineal** y el resultado se recorta a [0, 1]. El artículo asume RGB lineal (p. 1294) y así lo hace el paquete R colorspace desde la versión 2.1-0 (`simulate_cvd(linear = TRUE)`). Ojo: `filterDeficiencyProt/Deuter/Trit` de culori 4.0.2 aplican la matriz sobre sRGB codificado (`src/deficiency.js`); por eso el script también mide la paleta de gráficas con esa variante, como prueba de robustez (§8.6).

### 2.7 Escala de grises y distancias

El gris equivalente de un color es el gris de igual luminancia relativa Y; su luminosidad OKLab es exactamente L_gris = Y^(1/3) (para un gris, l = m = s = Y). La distancia entre colores es ΔE = 100·‖ΔOKLab‖ (euclidiana). Los umbrales de distinción entre series vecinas son los de la skill dataviz que el plan §4.1 adopta: ΔE ≥ 8 bajo protanopia y deuteranopia (piso 6 solo con codificación secundaria) y ΔE ≥ 15 en visión normal; además ΔL ≥ 0.10 entre series consecutivas en gris y en cada simulación (criterio U2 del plan §3.5; el encargo pide ≥ 0.08).

## 3. Matiz base, eje opuesto y armonía complementaria dividida

- **Matiz base 280° (índigo)**, aprobado por Axel (decisiones §I). Lleva la marca, las superficies tintadas (neutro con croma 0.012, apenas perceptible) y el acento principal. Lectura sobria y moderna, tipo Notion/Linear.
- **Eje opuesto 100° (oliva amarillo)** = 280° − 180°. Es el matiz de mayor contraste frente a la marca, por eso se **reserva para «acción crítica» de confirmación** (decisiones §I): un diálogo que pide confirmar algo irreversible es lo único que debe interrumpir. No se usa para error, que conserva la convención del rojo (20°).
- **Zona reservada 80°–120°.** Ningún otro token con croma visible (≥ 0.02; el neutro llega a 0.012) usa un matiz a menos de 20° del opuesto: ni roles, ni personas y tiendas, ni series de gráficas. El ancho es la separación mínima entre matices distintos que decisiones §I fija para roles y categóricas, y no puede ser mayor: advertencia (75°, también decidida) queda a 25° del opuesto (24.3° medidos en el hex). El script lo verifica en cada token, con el matiz de diseño y con el medido en el hex. Colores de identidad más cercanos al opuesto: advertencia (rol) a 24.3°; ingreso (rol) a 34.7°; verde (serie 8) a 40.1°; naranja (serie 4) a 44.7°; éxito (rol) a 50.0°. Hasta la ronda 1 de revisión la serie 6 de gráficas usaba 100° (hallazgo H-03); ahora usa 225° (§8.8).
- **Complementario dividido:** los dos vecinos del opuesto a ±30°, 70° y 130°, son los acentos que acompañan a la marca sin chocar con ella. Se asignan a los dos estados que más se ven en una app de finanzas además del gasto: **advertencia 75°** (decidida; 5° por encima del vecino exacto 70°) e **ingreso 135°** (la opción que decisiones §I prevé para U2; 5° por encima del vecino exacto 130°, el mismo desplazamiento que advertencia, §6.4). Así la armonía queda base 280° + acentos 75°/135° + opuesto 100° en reserva, y se evita el choque directo rojo/verde porque gasto y error (rojos) y éxito (verde) no forman parte de la armonía: son convenciones semánticas, con su propio matiz y siempre con icono y texto.

| Familia | Matiz | Distancia a la base | Distancia al opuesto | Origen |
|---|---|---|---|---|
| Marca | 280° | 0° | 180° | decisiones §I |
| Ingreso | 135° | 145° | 35° | Propuesta U2a con la opción que prevé decisiones §I (135°): a 15° de éxito y a 5° del acento dividido 130° |
| Gasto | 45° | 125° | 55° | K26 (puerta 1): terracota 40–45°; 45° es el grado del rango más separado de error y advertencia (Cambios 1.1) |
| Ahorro | 200° | 80° | 100° | decisiones §I |
| Deuda / crédito | 310° | 30° | 150° | decisiones §I |
| Advertencia | 75° | 155° | 25° | decisiones §I |
| Error | 20° | 100° | 80° | decisiones §I (rojo) |
| Éxito | 150° | 130° | 50° | decisiones §I (reasignado desde 280°) |
| Información | 240° | 40° | 140° | decisiones §I |
| Acción crítica | 100° | 180° | 0° | decisiones §I: eje opuesto reservado para confirmar acciones críticas |

## 4. Regla 60-30-10

Proporción orientativa del área visible de cada pantalla (no se mide en píxeles; guía la composición):

| Proporción | Papel | Qué cubre | Claro | Oscuro |
|---|---|---|---|---|
| 60 % | Superficies | fondo de página, tarjetas, divisores, texto neutro | neutral 50 (fondo), 100 (tarjeta), 200 (divisor), 700/900 (texto) | neutral 950 (fondo), 900 (tarjeta), 800 (divisor), 300/100 (texto) |
| 30 % | Acento | marca en pasos de superficie: barra lateral, elemento activo, fila seleccionada, encabezados de vista, chips de filtro | brand 50–200 | brand 900–950 |
| 10 % | Énfasis | una acción primaria por pantalla, enlaces, foco, roles semánticos, categóricas, series de gráficas, acción crítica | brand 700/800, roles paso 700, categóricas y series claras | brand 300/200, roles paso 400, categóricas y series oscuras |

El neutro está tintado con el matiz base (croma ≤ 0.012), así que el 90 % de la pantalla ya pertenece a la familia índigo y el 10 % de énfasis es lo único con saturación visible.

## 5. Escalas de 11 pasos

### 5.1 Marca (280°, croma máx. 0.11)

| Paso | L | Croma pedido | Croma final | Matiz | Hex | Plan §4.3 | Coincide |
|---|---|---|---|---|---|---|---|
| 50 | 0.985 | 0.0110 | 0.0071 (recortado) | 280° | `#F9FAFF` | `#F9FAFF` | sí |
| 100 | 0.962 | 0.0275 | 0.0182 (recortado) | 280° | `#F0F1FF` | `#F0F1FF` | sí |
| 200 | 0.925 | 0.0440 | 0.0366 (recortado) | 280° | `#E1E4FF` | `#E1E4FF` | sí |
| 300 | 0.865 | 0.0715 | 0.0674 (recortado) | 280° | `#C9CFFF` | `#C9CFFF` | sí |
| 400 | 0.785 | 0.0935 | 0.0935 | 280° | `#ADB3F4` | `#ADB3F4` | sí |
| 500 | 0.695 | 0.1100 | 0.1100 | 280° | `#9095E0` | `#9095E0` | sí |
| 600 | 0.605 | 0.1100 | 0.1100 | 280° | `#7579C2` | `#7579C2` | sí |
| 700 | 0.515 | 0.0990 | 0.0990 | 280° | `#5D609F` | `#5D609F` | sí |
| 800 | 0.425 | 0.0825 | 0.0825 | 280° | `#46487B` | `#46487B` | sí |
| 900 | 0.335 | 0.0660 | 0.0660 | 280° | `#303258` | `#303258` | sí |
| 950 | 0.255 | 0.0495 | 0.0495 | 280° | `#1E203A` | `#1E203A` | sí |

### 5.2 Neutro tintado (280°, croma máx. 0.012)

| Paso | L | Croma pedido | Croma final | Matiz | Hex | Plan §4.3 | Coincide |
|---|---|---|---|---|---|---|---|
| 50 | 0.985 | 0.0012 | 0.0012 | 280° | `#FAFAFB` | `#FAFAFB` | sí |
| 100 | 0.962 | 0.0030 | 0.0030 | 280° | `#F2F2F4` | `#F2F2F4` | sí |
| 200 | 0.925 | 0.0048 | 0.0048 | 280° | `#E5E6E9` | `#E5E6E9` | sí |
| 300 | 0.865 | 0.0078 | 0.0078 | 280° | `#D1D2D8` | `#D1D2D8` | sí |
| 400 | 0.785 | 0.0102 | 0.0102 | 280° | `#B7B8C0` | `#B7B8C0` | sí |
| 500 | 0.695 | 0.0120 | 0.0120 | 280° | `#9B9CA4` | `#9B9CA4` | sí |
| 600 | 0.605 | 0.0120 | 0.0120 | 280° | `#808189` | `#808189` | sí |
| 700 | 0.515 | 0.0108 | 0.0108 | 280° | `#66676E` | `#66676E` | sí |
| 800 | 0.425 | 0.0090 | 0.0090 | 280° | `#4D4E53` | `#4D4E53` | sí |
| 900 | 0.335 | 0.0072 | 0.0072 | 280° | `#36363A` | `#36363A` | sí |
| 950 | 0.255 | 0.0054 | 0.0054 | 280° | `#222325` | `#222325` | sí |

### 5.3 Escalas de rol y de acción crítica (croma máx. 0.14)

| Familia | 50 | 100 | 200 | 300 | 400 | 500 | 600 | 700 | 800 | 900 | 950 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| Ingreso (135°) | `#F6FCF3` | `#E9F8E1` | `#D7EFCA` | `#B9E1A4` | `#98CA7B` | `#75B052` | `#5B9335` | `#457723` | `#335A18` | `#22400D` | `#142907` |
| Gasto (45°) | `#FFF9F6` | `#FFEFE8` | `#FFDED0` | `#FFC3A9` | `#F8A078` | `#E37E4E` | `#C46231` | `#A04B20` | `#7C3815` | `#59260B` | `#3B1705` |
| Ahorro (200°) | `#F0FDFE` | `#D8FAFC` | `#BBF2F5` | `#85E5EA` | `#3FD0D7` | `#00B2BA` | `#00949A` | `#00767B` | `#005A5E` | `#003F42` | `#00292B` |
| Deuda / crédito (310°) | `#FCF8FF` | `#F7EEFF` | `#F0DEFF` | `#E5C3FF` | `#D0A3F2` | `#B782DD` | `#9B67C0` | `#7E509D` | `#603C79` | `#442857` | `#2C1939` |
| Advertencia (75°) | `#FFF9F1` | `#FFF0DC` | `#FDE2BD` | `#F6CB8F` | `#E5AE5B` | `#CF8E1B` | `#AD7500` | `#8B5D00` | `#6B4600` | `#4C3000` | `#321E00` |
| Error (20°) | `#FFF8F8` | `#FFEEED` | `#FFDDDC` | `#FFC0BF` | `#FC9999` | `#E77678` | `#C75A5E` | `#A34449` | `#7E3236` | `#5A2124` | `#3C1416` |
| Éxito (150°) | `#F4FDF5` | `#E3FAE6` | `#CDF1D3` | `#A8E4B3` | `#7ECF90` | `#53B56D` | `#339852` | `#207B3E` | `#155E2D` | `#0A421E` | `#052B11` |
| Información (240°) | `#F6FBFF` | `#E7F5FF` | `#CFEBFF` | `#A6DAFF` | `#6BC3FD` | `#35A7EA` | `#008ACB` | `#006FA4` | `#00547E` | `#003B5A` | `#00263C` |
| Acción crítica (100°) | `#FCFBF0` | `#F8F4D9` | `#EEE8BD` | `#E0D58E` | `#CBBB5A` | `#B29E15` | `#948200` | `#766800` | `#5A4F00` | `#3F3700` | `#292300` |

Cada rol tiene su escala completa en el JSON (`color.role.<rol>.50` … `.950`, y `color.critical.*`), porque los estados hover usan pasos vecinos (800 en claro, 200 en oscuro).

## 6. Roles semánticos

### 6.1 Valores finales

Texto e icono del rol: paso 700 en claro y paso 400 en oscuro. Fondo de insignia: paso 100 en claro y paso 900 en oscuro. Contraste mínimo exigido: 4.5:1, porque los roles se usan como texto (cifras con signo, mensajes, insignias).

| Rol | Token | Matiz | Claro: paso · hex | Contraste fondo 50 / tarjeta 100 / insignia | Oscuro: paso · hex | Contraste fondo 950 / tarjeta 900 / insignia |
|---|---|---|---|---|---|---|
| Ingreso | `color.role.income` | 135° | 700 · `#457723` | 5.14 / 4.79 / 4.84 | 400 · `#98CA7B` | 8.29 / 6.34 / 6.12 |
| Gasto | `color.role.expense` | 45° | 700 · `#A04B20` | 5.71 / 5.33 / 5.32 | 400 · `#F8A078` | 7.70 / 5.89 / 6.02 |
| Ahorro | `color.role.savings` | 200° | 700 · `#00767B` | 5.19 / 4.84 / 4.90 | 400 · `#3FD0D7` | 8.39 / 6.42 / 6.26 |
| Deuda / crédito | `color.role.debt` | 310° | 700 · `#7E509D` | 5.71 / 5.33 / 5.29 | 400 · `#D0A3F2` | 7.64 / 5.84 / 6.04 |
| Advertencia | `color.role.warning` | 75° | 700 · `#8B5D00` | 5.49 / 5.12 / 5.11 | 400 · `#E5AE5B` | 7.89 / 6.03 / 6.09 |
| Error | `color.role.error` | 20° | 700 · `#A34449` | 5.78 / 5.39 / 5.37 | 400 · `#FC9999` | 7.60 / 5.81 / 6.05 |
| Éxito | `color.role.success` | 150° | 700 · `#207B3E` | 5.07 / 4.73 / 4.81 | 400 · `#7ECF90` | 8.39 / 6.42 / 6.19 |
| Información | `color.role.info` | 240° | 700 · `#006FA4` | 5.28 / 4.93 / 4.96 | 400 · `#6BC3FD` | 8.12 / 6.21 / 6.11 |
| Acción crítica | `color.critical` | 100° | 700 · `#766800` | 5.36 / 5.00 / 5.04 | 400 · `#CBBB5A` | 8.08 / 6.18 / 6.14 |

«Acción crítica» no es uno de los 8 roles: es el énfasis del eje opuesto (§3) y vive en `color.critical`.

### 6.2 Ajuste de paso: 600 → 700 en claro

El encargo pide el paso 600 sobre fondo 50 en claro. Con la L fija del paso 600 (0.605) ningún matiz llega a 4.5:1, así que se baja un paso (700, L 0.515) y se documenta aquí. En oscuro el paso 400 cumple sin cambios.

| Rol | Paso 600 | Contraste mínimo (fondo, tarjeta, insignia) | Paso elegido | Contraste mínimo | Cambio |
|---|---|---|---|---|---|
| Ingreso | `#5B9335` | 3.31 | `#457723` | 4.79 | 600 → 700 |
| Gasto | `#C46231` | 3.65 | `#A04B20` | 5.32 | 600 → 700 |
| Ahorro | `#00949A` | 3.29 | `#00767B` | 4.84 | 600 → 700 |
| Deuda / crédito | `#9B67C0` | 3.65 | `#7E509D` | 5.29 | 600 → 700 |
| Advertencia | `#AD7500` | 3.52 | `#8B5D00` | 5.11 | 600 → 700 |
| Error | `#C75A5E` | 3.70 | `#A34449` | 5.37 | 600 → 700 |
| Éxito | `#339852` | 3.26 | `#207B3E` | 4.73 | 600 → 700 |
| Información | `#008ACB` | 3.41 | `#006FA4` | 4.93 | 600 → 700 |
| Acción crítica | `#948200` | 3.44 | `#766800` | 5.00 | 600 → 700 |

### 6.3 Comparación con los valores del plan §4.3

El plan calculó los roles a L 0.56 en claro y L 0.74 en oscuro (sus contrastes son contra marca 50 `#F9FAFF`). Recalculados contra las superficies definitivas, varios quedan bajo 4.5:1 para texto, lo que confirma el ajuste anterior (7 de 7 con valor quedan bajo 4.5:1 en al menos una superficie; rango 3.92–4.47:1):

| Rol | Plan claro | Contraste fondo 50 / tarjeta 100 | Plan oscuro | Contraste fondo 950 / tarjeta 900 | Nota |
|---|---|---|---|---|---|
| Ingreso | `#218A45` | 4.20 / 3.92 | `#6FC082` | 7.14 / 5.46 | 150° en el plan |
| Gasto | `#B84D49` | 4.79 / 4.47 | `#ED8C84` | 6.50 / 4.97 | — |
| Ahorro | `#00858A` | 4.26 / 3.97 | `#24C1C9` | 7.15 / 5.46 | — |
| Deuda / crédito | `#8D5AB1` | 4.75 / 4.44 | `#C295E3` | 6.51 / 4.98 | — |
| Advertencia | `#9C6900` | 4.54 / 4.23 | `#D79F4C` | 6.70 / 5.12 | — |
| Error | `#857500` | 4.43 / 4.13 | `#BDAC4A` | 6.86 / 5.24 | 100° en el plan |
| Éxito | — | — | — | — | sin valor: 280° se reasignó |
| Información | `#007CB7` | 4.40 / 4.11 | `#5CB4EF` | 6.91 / 5.28 | — |

### 6.4 Ingreso frente a éxito (propuesta)

Decisiones §I deja ingreso en 150° (misma familia que éxito, distinguido por icono y contexto) «o U2 propone 135°». **Propuesta U2a: ingreso 135°**, la opción ya prevista, así que no se reabre nada. Razones: (1) separa ingreso de éxito 15° en vez de 0°; (2) queda a 5° del acento dividido 130° (100° + 30°), el mismo desplazamiento que advertencia tiene respecto a 70°, de modo que ingreso y advertencia forman el par de acentos de la armonía (§3); (3) queda a 35° del opuesto 100° (acción crítica) y a 60° de advertencia. Distancias calculadas entre ingreso y éxito (texto de rol): claro ΔE normal 3.3, protanopia 2.8, deuteranopia 2.8, tritanopia 3.8; oscuro 3.2, 2.6, 2.6, 4.1. Lado a lado apenas se distinguen en visión normal (ΔE ≈ 3, cerca del umbral de percepción) y no se distinguen con daltonismo, por eso el color nunca va solo: **ingreso** siempre lleva signo «+» y flecha de entrada en cifras y listas; **éxito** siempre lleva palomita y texto («Guardado», «Pagado»), y vive en avisos y estados, no en cifras.

### 6.5 Gasto frente a error

Versión 1.1 (K26; bloqueo `U2-gasto-error` resuelto con la opción 2): gasto pasa de 25° a 45° (terracota) y error se queda en 20°. Medidos en el hex quedan a 25.3° de matiz en claro y 24.9° en oscuro, con ΔE normal 5.5 y 5.1 (en la 1.0, con gasto en 25°, eran 1.2 y 1.0); el grado sale del barrido de «Cambios 1.1». Con daltonismo rojo-verde siguen cerca (ΔE protanopia 5.9 / 5.0, deuteranopia 4.4 / 4.1), así que se conserva la regla de uso: el **gasto** aparece como cifra con signo «−» en listas y tablas; el **error** solo como mensaje de validación o aviso con icono de alerta y texto. Nunca comparten componente.

### 6.6 Pares de roles cercanos (ΔE normal < 10)

Todos los roles comparten L y croma, así que la distancia depende solo del ángulo de matiz. Se listan los pares con ΔE normal < 10 (los que pueden confundirse a simple vista):

| Tema | Par | Δ matiz | ΔE normal | ΔE protanopia | ΔE deuteranopia | ΔE tritanopia |
|---|---|---|---|---|---|---|
| claro | Ingreso – Éxito | 15° | 3.3 | 2.8 | 2.8 | 3.8 |
| claro | Advertencia – Acción crítica | 25° | 4.6 | 1.9 | 0.4 | 5.2 |
| claro | Gasto – Error | 25° | 5.5 | 5.9 | 4.4 | 1.0 |
| claro | Gasto – Advertencia | 30° | 6.4 | 3.4 | 1.1 | 5.9 |
| claro | Ingreso – Acción crítica | 35° | 7.3 | 2.9 | 2.1 | 9.5 |
| claro | Ahorro – Información | 40° | 7.6 | 7.3 | 7.4 | 1.6 |
| claro | Ahorro – Éxito | 50° | 9.7 | 9.4 | 9.3 | 1.7 |
| claro | Éxito – Acción crítica | 50° | 10.0 | 4.8 | 4.9 | 13.3 |
| oscuro | Ingreso – Éxito | 15° | 3.2 | 2.6 | 2.6 | 4.1 |
| oscuro | Advertencia – Acción crítica | 25° | 5.0 | 2.2 | 0.2 | 5.7 |
| oscuro | Gasto – Error | 25° | 5.1 | 5.0 | 4.1 | 1.1 |
| oscuro | Gasto – Advertencia | 30° | 6.3 | 4.6 | 2.8 | 3.4 |
| oscuro | Ingreso – Acción crítica | 35° | 7.1 | 3.9 | 3.8 | 9.6 |
| oscuro | Ahorro – Información | 40° | 8.1 | 6.7 | 6.1 | 2.0 |

Todos estos pares viven en contextos distintos (cifra, aviso, diálogo de confirmación) y todos llevan icono y texto; ningún significado depende solo del color (criterio 1.4.1).

## 7. Paleta categórica: personas y tiendas

Reglas (plan §4.2–4.3 y encargo): luminosidad escalonada con ΔL ≥ 0.10 entre todos los pares de cada grupo (medida en L OKLab y en L del gris equivalente), separación ≥ 20° de matiz **o** ΔL ≥ 0.15 respecto a cada rol del mismo tema (también se exige frente a acción crítica), y contraste ≥ 3:1 sobre fondo y tarjeta. Un solo valor por color no puede cumplir 3:1 en los dos temas y a la vez escalonarse: la ventana común es L_gris ∈ [0.596, 0.641] (ancho 0.045 < 0.10). Por eso cada color tiene variante clara y oscura con el **mismo matiz** (la identidad se conserva) y distinta L.

### 7.1 Personas

| Nombre | Token | Matiz | Croma pedido | Claro: hex · L · L gris · contraste mín. | Oscuro: hex · L · L gris · contraste mín. | Motivo |
|---|---|---|---|---|---|---|
| Axel | `color.person.axel` | 256° | 0.120 | `#002E64` · L 0.310 · gris 0.306 · 11.92 | `#4F83C8` · L 0.605 · gris 0.604 · 3.10 | Azul a 256° (6° más violeta que el 250° del plan para conservar croma ≥ 0.10 en L 0.31); como queda a 16° de información, se separa por luminosidad (ΔL ≥ 0.15 frente al rol en ambos temas). |
| Ximena | `color.person.ximena` | 350° | 0.120 | `#A14E78` · L 0.540 · gris 0.524 · 4.84 | `#FFACD3` · L 0.835 · gris 0.822 · 6.92 | Conserva 350° del plan (30° de error, 40° de deuda). |
| Karen | `color.person.karen` | 333° | 0.120 | `#B870AC` · L 0.642 · gris 0.628 · 3.15 | `#D389C6` · L 0.725 · gris 0.710 · 4.67 | Versión 1.1 (K26): el 50° de la 1.0 quedaba a 5° de gasto 45°. Con su L (la más alta del grupo en claro, la última que da 3:1) se separa de los roles por matiz, y sus ventanas factibles quedan entre deuda 310° y error 20°; 333° queda a 23° de deuda y de Bodega Aurrera, lo más lejos de Ximena (350°) que permite el reparto, y de ella se separa por ΔL ≥ 0.10. L en claro 0.642 (antes 0.640) para dar margen al escalonamiento. PROVISIONAL (bloqueo `U2a-rosas-ximena-bodega`): frente a Ximena queda bajo ≥ 20° o ΔL ≥ 0.15 en ambos temas y se distingue por el refuerzo de §7.7 (Cambios 1.1, punto 2). |
| Conjunto | `color.person.shared` | 280° | 0.045 | `#484B66` · L 0.420 · gris 0.419 · 7.59 | `#E7EAFF` · L 0.940 · gris 0.940 · 10.08 | Lo compartido lleva el matiz de la marca con croma bajo (casi neutro), como el croma 0.03 del plan; 280° queda a 30° de deuda y 40° de información. Es la única excepción deliberada al piso de croma 0.10: se lee como «de ambos», no como una persona más. |

### 7.2 Tiendas

| Nombre | Token | Matiz | Croma pedido | Claro: hex · L · L gris · contraste mín. | Oscuro: hex · L · L gris · contraste mín. | Motivo |
|---|---|---|---|---|---|---|
| Walmart | `color.store.walmart` | 262° | 0.120 | `#4065AC` · L 0.515 · gris 0.511 · 5.11 | `#7AA3EE` · L 0.716 · gris 0.715 · 4.75 | 220° del plan quedaba a 20.0° exactos de ahorro e información (frágil al cuantizar a 8 bits); se mueve al azul de 262° (22° de información). PROVISIONAL en oscuro (bloqueo `U2a-walmart-oscuro`): L 0.716 en vez de 0.720 para que Bodega Aurrera quepa en el escalonamiento de las tiendas (Cambios 1.1). |
| Bodega Aurrera | `color.store.bodega-aurrera` | 356° | 0.110 | `#762C4C` · L 0.410 · gris 0.395 · 8.40 | `#FFAAC9` · L 0.829 · gris 0.815 · 6.78 | Versión 1.1 (K26): el 50° de la 1.0 quedaba a 5° de gasto 45°. Con su L (entre Mercado y Walmart) se separa de los roles por matiz en ambos temas: 356° queda a 24° de error y a 23° de Karen, de la que en oscuro no se separa por luminosidad. Croma 0.11 y L 0.829 en oscuro: el rosa pierde luminancia en gris y con 0.12 no cabe en el escalonamiento de las tiendas (Cambios 1.1). PROVISIONAL (bloqueo `U2a-rosas-ximena-bodega`): frente a Ximena queda bajo ≥ 20° o ΔL ≥ 0.15 en ambos temas (en oscuro, a simple vista el mismo color) y se distingue solo por el refuerzo de §7.7 (Cambios 1.1, punto 2). |
| Sam's Club | `color.store.sams-club` | 175° | 0.120 | `#009D82` · L 0.620 · gris 0.636 · 3.05 | `#99FFE4` · L 0.930 · gris 0.943 · 10.18 | Conserva 175° del plan (25° de éxito y de ahorro); va en la L más alta de cada tema, donde el verde azulado conserva croma ≥ 0.10. |
| Mercado | `color.store.market` | 28° | 0.120 | `#5C0204` · L 0.300 · gris 0.286 · 12.81 | `#C2665B` · L 0.615 · gris 0.601 · 3.06 | Terracota a 28° (2° menos que el plan, para quedar a más de 20° de Bodega Aurrera); como queda a 3° de gasto y 8° de error, se separa por luminosidad (ΔL ≥ 0.15 frente a los roles en ambos temas). |

### 7.3 Escalonamiento y distinción dentro de cada grupo

| Grupo | Tema | ΔL mín. (OKLab, todos los pares) | ΔL mín. (gris, todos los pares) | ΔE vecinos normal | ΔE vecinos protan/deutan | ΔE vecinos tritan | Contraste mín. | Croma mín. (sin Conjunto) |
|---|---|---|---|---|---|---|---|---|
| Personas | claro | 0.103 | 0.104 | 10.9 | 10.5 | 10.7 | 3.15 | 0.106 |
| Personas | oscuro | 0.105 | 0.106 | 11.6 | 11.1 | 9.9 | 3.10 | 0.108 |
| Tiendas | claro | 0.105 | 0.110 | 19.8 | 15.5 | 23.2 | 3.05 | 0.110 |
| Tiendas | oscuro | 0.101 | 0.101 | 20.0 | 9.5 | 22.9 | 3.06 | 0.103 |

Orden fijo de leyenda y apilado: personas Axel, Ximena, Karen, Conjunto; tiendas Walmart, Bodega Aurrera, Sam's Club, Mercado (los «vecinos» de la tabla son los contiguos en ese orden). Con cuatro niveles separados ≥ 0.10 no cabe la banda de L de la skill dataviz (§8.7): el escalonamiento ocupa casi todo el rango que permite el contraste, así que algunos colores quedan muy oscuros en claro (L 0.30–0.31) o muy claros en oscuro (L 0.93–0.94). El piso de croma 0.10 se cumple en todos salvo Conjunto, que es bajo a propósito.

### 7.4 Separación respecto a los roles

Para cada color, el rol más exigente (menor holgura) y la regla con la que cumple:

| Color | Tema | Rol más cercano | Δ matiz | ΔL | Cumple por | Cumple con todos |
|---|---|---|---|---|---|---|
| Axel | claro | Información | 16.4° | 0.207 | ΔL ≥ 0.15 | sí |
| Axel | oscuro | Información | 16.5° | 0.180 | ΔL ≥ 0.15 | sí |
| Ximena | claro | Error | 29.2° | 0.026 | Δ matiz ≥ 20° | sí |
| Ximena | oscuro | Error | 30.1° | 0.051 | Δ matiz ≥ 20° | sí |
| Karen | claro | Deuda / crédito | 22.9° | 0.127 | Δ matiz ≥ 20° | sí |
| Karen | oscuro | Deuda / crédito | 23.4° | 0.059 | Δ matiz ≥ 20° | sí |
| Conjunto | claro | Deuda / crédito | 30.9° | 0.094 | Δ matiz ≥ 20° | sí |
| Conjunto | oscuro | Deuda / crédito | 30.5° | 0.157 | Δ matiz ≥ 20° | sí |
| Walmart | claro | Información | 22.4° | 0.002 | Δ matiz ≥ 20° | sí |
| Walmart | oscuro | Información | 22.3° | 0.069 | Δ matiz ≥ 20° | sí |
| Bodega Aurrera | claro | Error | 23.4° | 0.104 | Δ matiz ≥ 20° | sí |
| Bodega Aurrera | oscuro | Error | 23.9° | 0.044 | Δ matiz ≥ 20° | sí |
| Sam's Club | claro | Éxito | 24.7° | 0.106 | Δ matiz ≥ 20° | sí |
| Sam's Club | oscuro | Ahorro | 24.7° | 0.144 | Δ matiz ≥ 20° | sí |
| Mercado | claro | Gasto | 16.9° | 0.214 | ΔL ≥ 0.15 | sí |
| Mercado | oscuro | Error | 8.0° | 0.170 | ΔL ≥ 0.15 | sí |

### 7.5 De la propuesta del plan a los valores finales

| Nombre | Plan (valor único) | Contraste del plan claro / oscuro | Final: matiz · claro / oscuro |
|---|---|---|---|
| Axel | `#0F74C5` 250° L 0.55 | 4.34 / 2.47 | 256° · `#002E64` / `#4F83C8` |
| Ximena | `#D375A4` 350° L 0.68 | 2.74 / 3.92 | 350° · `#A14E78` / `#FFACD3` |
| Karen | `#E0AE57` 80° L 0.78 | 1.81 / 5.93 | 333° · `#B870AC` / `#D389C6` |
| Conjunto | `#495766` 250° L 0.45 | 6.61 / 1.62 | 280° · `#484B66` / `#E7EAFF` |
| Walmart | `#2FB5D8` 220° L 0.72 | 2.15 / 4.99 | 262° · `#4065AC` / `#7AA3EE` |
| Bodega Aurrera | `#FFA658` 60° L 0.80 | 1.73 / 6.21 | 356° · `#762C4C` / `#FFAAC9` |
| Sam's Club | `#11957C` 175° L 0.60 | 3.34 / 3.21 | 175° · `#009D82` / `#99FFE4` |
| Mercado | `#9C4438` 30° L 0.50 | 5.69 / 1.88 | 28° · `#5C0204` / `#C2665B` |

Colisiones resueltas: las dos que lista decisiones §I (Karen 80° y Bodega 60° frente a advertencia 75°) y tres más que aparecen al aplicar la regla a todos los roles: Axel 250° frente a información 240°, Mercado 30° frente a gasto 25° y error 20°, y Walmart 220° justo en el límite de ahorro 200° e información 240°. Además, 4 valores del plan no llegaban a 3:1 sobre las superficies claras y 3 sobre las oscuras. En la versión 1.1 (K26) gasto pasa a terracota, Karen y Bodega Aurrera salen de 50° y Walmart ajusta su L en oscuro (PROVISIONAL); frente a Ximena, Karen y Bodega Aurrera quedan PROVISIONALES (bloqueo `U2a-rosas-ximena-bodega`). Detalle en «Cambios 1.1».

### 7.6 Reglas de uso de las categóricas

- Nunca llevan texto encima: el nombre va en texto neutro junto a un punto, anillo de avatar o marca de serie del color.
- La misma persona o tienda conserva su matiz en ambos temas; solo cambia la luminosidad.
- En gráficas por persona o por tienda se usan estos colores (el color sigue a la entidad) y no la paleta de series; con leyenda siempre y etiqueta directa cuando hay ≤ 4 elementos.
- Mercados adicionales que Axel dé de alta usan el color de Mercado; si una gráfica necesita distinguir varios mercados, se usan las series de gráficas en su orden.

### 7.7 Pares cercanos en todo el sistema y su refuerzo

Revisión de todos los pares entre los 8 roles, acción crítica, las 4 personas, las 4 tiendas y la marca (paso 700 en claro y 400 en oscuro; hallazgo R3-04, versión 1.1), en cada tema, con matiz y L medidos en el hex final. Se listan los que quedan con Δ matiz < 20° **y** ΔL < 0.15, con el refuerzo que evita confundirlos (ninguno es categórica ni marca frente a rol):

| Tema | Color A | Color B | Δ matiz | ΔL | Refuerzo |
|---|---|---|---|---|---|
| claro | Ingreso (rol) | Éxito (rol) | 14.9° | 0.000 | Ingreso siempre con «+» y flecha de entrada en cifras; éxito siempre con palomita y texto en avisos. |
| claro | Ximena (persona) | Karen (persona) | 17.1° | 0.103 | Mismo grupo: orden fijo, etiqueta directa y ΔL ≥ 0.10 dentro del grupo. **PROVISIONAL**: bloqueo `U2a-rosas-ximena-bodega` (¿basta este refuerzo frente a «≥ 20° o ΔL ≥ 0.15» de K26?). |
| claro | Ximena (persona) | Bodega Aurrera (tienda) | 5.8° | 0.130 | Grupos distintos (persona y tienda): nunca comparten gráfica ni leyenda; en tablas van en columnas distintas y con su nombre. **PROVISIONAL**: bloqueo `U2a-rosas-ximena-bodega` (¿basta este refuerzo frente a «≥ 20° o ΔL ≥ 0.15» de K26?). |
| claro | Conjunto (persona) | Walmart (tienda) | 17.1° | 0.093 | Grupos distintos (persona y tienda): nunca comparten gráfica ni leyenda; en tablas van en columnas distintas y con su nombre. |
| claro | Conjunto (persona) | Marca 700 (marca) | 0.9° | 0.094 | Marca frente a categórica (R3-04): la marca solo vive en controles (botón primario, enlace subrayado, anillo de foco, elemento activo) y nunca marca a una persona o tienda; la categórica siempre va como punto, anillo de avatar o marca de serie junto a su nombre en texto neutro (forma y etiqueta). |
| claro | Walmart (tienda) | Marca 700 (marca) | 18.0° | 0.000 | Marca frente a categórica (R3-04): la marca solo vive en controles (botón primario, enlace subrayado, anillo de foco, elemento activo) y nunca marca a una persona o tienda; la categórica siempre va como punto, anillo de avatar o marca de serie junto a su nombre en texto neutro (forma y etiqueta). |
| oscuro | Ingreso (rol) | Éxito (rol) | 15.3° | 0.001 | Ingreso siempre con «+» y flecha de entrada en cifras; éxito siempre con palomita y texto en avisos. |
| oscuro | Axel (persona) | Walmart (tienda) | 5.9° | 0.112 | Grupos distintos (persona y tienda): nunca comparten gráfica ni leyenda; en tablas van en columnas distintas y con su nombre. |
| oscuro | Ximena (persona) | Karen (persona) | 16.7° | 0.111 | Mismo grupo: orden fijo, etiqueta directa y ΔL ≥ 0.10 dentro del grupo. **PROVISIONAL**: bloqueo `U2a-rosas-ximena-bodega` (¿basta este refuerzo frente a «≥ 20° o ΔL ≥ 0.15» de K26?). |
| oscuro | Ximena (persona) | Bodega Aurrera (tienda) | 6.2° | 0.006 | Grupos distintos (persona y tienda): nunca comparten gráfica ni leyenda; en tablas van en columnas distintas y con su nombre. **PROVISIONAL**: bloqueo `U2a-rosas-ximena-bodega` (¿basta este refuerzo frente a «≥ 20° o ΔL ≥ 0.15» de K26?). |
| oscuro | Walmart (tienda) | Marca 400 (marca) | 18.0° | 0.069 | Marca frente a categórica (R3-04): la marca solo vive en controles (botón primario, enlace subrayado, anillo de foco, elemento activo) y nunca marca a una persona o tienda; la categórica siempre va como punto, anillo de avatar o marca de serie junto a su nombre en texto neutro (forma y etiqueta). |

## 8. Paleta de gráficas (8 series)

### 8.1 Orden de uso

| Serie | Nombre | Matiz | Clase | Claro: hex · L · C · contraste mín. | Oscuro: hex · L · C · contraste mín. |
|---|---|---|---|---|---|
| 1 | índigo | 280° | B | `#474698` · L 0.440 · C 0.130 · 7.23 | `#797DCD` · L 0.620 · C 0.120 · 3.22 |
| 2 | verde azulado | 180° | A | `#009A87` · L 0.615 · C 0.112 · 3.14 | `#61E4CD` · L 0.840 · C 0.120 · 7.72 |
| 3 | óxido | 45° | B | `#883500` · L 0.440 · C 0.126 · 7.34 | `#C16D45` · L 0.620 · C 0.120 · 3.17 |
| 4 | naranja | 55° | A | `#BF6D2C` · L 0.615 · C 0.130 · 3.45 | `#FFB989` · L 0.840 · C 0.103 · 7.19 |
| 5 | magenta | 345° | B | `#812E63` · L 0.440 · C 0.130 · 7.49 | `#B86797` · L 0.620 · C 0.120 · 3.10 |
| 6 | azul cielo | 225° | A | `#0092B9` · L 0.615 · C 0.117 · 3.23 | `#76D9FF` · L 0.840 · C 0.107 · 7.52 |
| 7 | violeta | 310° | B | `#683987` · L 0.440 · C 0.130 · 7.44 | `#9C70BC` · L 0.620 · C 0.120 · 3.12 |
| 8 | verde | 140° | A | `#579749` · L 0.615 · C 0.130 · 3.16 | `#A0DE93` · L 0.840 · C 0.120 · 7.65 |

Tokens `color.chart.<n>.light` y `color.chart.<n>.dark`. Las series se asignan en este orden fijo, nunca se reciclan: una novena serie se agrupa en «Otros» (neutral 600 en claro, neutral 400 en oscuro, que cumplen 3:1 como los bordes de campo).

### 8.2 Regla de luminosidad

Las series alternan dos clases de luminosidad para que las vecinas se separen en gris y con daltonismo: clase A (alta) L 0.615 en claro y 0.84 en oscuro; clase B (baja) L 0.44 en claro y 0.62 en oscuro. Las impares son B y las pares A. Los matices de clase B (índigo, óxido, magenta, violeta) conservan croma en luminosidad baja y los de clase A (verde azulado, naranja, azul cielo, verde) en luminosidad alta, por eso la clase de cada matiz es la misma en los dos temas. Los límites salen del contraste: en claro la clase A queda bajo el máximo que da 3:1 sobre la tarjeta 100 y en oscuro la clase B sobre el mínimo que da 3:1 sobre la tarjeta 900. Croma pedido 0.13 en claro y 0.12 en oscuro; todas las series quedan con croma final ≥ 0.10 (piso de la skill dataviz para que el matiz siga identificando). El orden se eligió en la ronda 0 por búsqueda exhaustiva entre las permutaciones que alternan clases, con índigo fijo en la serie 1 (gráficas de una sola serie) y violeta y verde al final (para que las primeras cinco tengan matices bien distintos y el verde no sugiera «positivo» en gráficas comunes), maximizando el mínimo ΔL consecutivo en las cuatro vistas y, después, el ΔE en visión normal entre vecinas. En la ronda 1 de revisión (hallazgo H-03) cambió el matiz de la serie 6 (§8.8) y el script repite la búsqueda con los matices vigentes (12 órdenes posibles; tabla siguiente). Como tercer criterio, solo para desempatar, se agrega el ΔE en visión normal entre todos los pares de las tres primeras series, las únicas que se usan en dispersión, mapas y pequeños múltiplos (§8.5). El orden vigente queda en el puesto 1. Respecto a la ronda 0 pone óxido en la serie 3 y magenta en la 5: con el orden anterior y el matiz nuevo, el mínimo ΔL consecutivo bajaría a 0.154 (puesto 4) y las tres primeras series quedarían con ΔE normal mínimo 12.8 entre todos sus pares, bajo el piso 15 de la skill dataviz; con el orden vigente son 0.162 y 21.3. Empata en los dos primeros criterios con el orden índigo, azul cielo, magenta, naranja, óxido, verde azulado, violeta, verde, que pierde en el tercero (12.8 frente a 21.3).

| Puesto | Orden (series 1 a 8) | ΔL consecutivo mín. | ΔE normal mín. entre vecinas | ΔE normal mín. entre las 3 primeras | ΔE protan/deutan mín. entre las 3 primeras | Nota |
|---|---|---|---|---|---|---|
| 1 | índigo, verde azulado, óxido, naranja, magenta, azul cielo, violeta, verde | 0.162 | 17.8 | 21.3 | 19.1 | vigente |
| 2 | índigo, azul cielo, magenta, naranja, óxido, verde azulado, violeta, verde | 0.162 | 17.8 | 12.8 | 8.4 |  |
| 3 | índigo, azul cielo, óxido, naranja, magenta, verde azulado, violeta, verde | 0.154 | 17.8 | 20.8 | 18.5 |  |
| 4 | índigo, verde azulado, magenta, naranja, óxido, azul cielo, violeta, verde | 0.154 | 17.8 | 12.8 | 8.4 | ronda 0 con la serie 6 nueva |

Los cuatro primeros de 12 órdenes; cada valor es el mínimo de los dos temas.

### 8.3 Distancias entre series consecutivas

**Tema claro** (series sobre fondo #FAFAFB y tarjeta #F2F2F4)

| Par | ΔL gris | ΔL protanopia | ΔL deuteranopia | ΔL tritanopia | ΔE normal | ΔE protanopia | ΔE deuteranopia | ΔE tritanopia |
|---|---|---|---|---|---|---|---|---|
| 1-2 | 0.198 | 0.206 | 0.186 | 0.193 | 25.5 | 25.2 | 22.5 | 20.4 |
| 2-3 | 0.202 | 0.270 | 0.162 | 0.182 | 28.1 | 27.8 | 19.1 | 31.9 |
| 3-4 | 0.179 | 0.183 | 0.176 | 0.177 | 17.8 | 18.5 | 17.7 | 17.7 |
| 4-5 | 0.183 | 0.174 | 0.185 | 0.180 | 23.2 | 23.6 | 23.0 | 18.4 |
| 5-6 | 0.200 | 0.261 | 0.162 | 0.191 | 27.6 | 26.1 | 17.6 | 30.1 |
| 6-7 | 0.198 | 0.232 | 0.173 | 0.196 | 24.2 | 23.6 | 17.3 | 24.7 |
| 7-8 | 0.203 | 0.217 | 0.186 | 0.179 | 31.3 | 29.9 | 25.3 | 21.4 |

Mínimo ΔL consecutivo en gris y simulaciones: **0.162** (umbral 0.10; el encargo pide 0.08). ΔE mínimo entre vecinas: normal 17.8 (piso 15), protan/deutan 17.3 (objetivo 8), tritan 17.7.

**Tema oscuro** (series sobre fondo #222325 y tarjeta #36363A)

| Par | ΔL gris | ΔL protanopia | ΔL deuteranopia | ΔL tritanopia | ΔE normal | ΔE protanopia | ΔE deuteranopia | ΔE tritanopia |
|---|---|---|---|---|---|---|---|---|
| 1-2 | 0.241 | 0.252 | 0.227 | 0.223 | 28.6 | 28.6 | 25.3 | 24.5 |
| 2-3 | 0.245 | 0.310 | 0.205 | 0.212 | 31.1 | 31.6 | 22.8 | 36.1 |
| 3-4 | 0.223 | 0.232 | 0.218 | 0.198 | 22.1 | 23.2 | 21.8 | 20.4 |
| 4-5 | 0.228 | 0.224 | 0.226 | 0.202 | 25.5 | 25.8 | 25.3 | 20.4 |
| 5-6 | 0.242 | 0.296 | 0.208 | 0.211 | 29.3 | 29.6 | 21.7 | 33.4 |
| 6-7 | 0.240 | 0.272 | 0.217 | 0.217 | 26.8 | 27.5 | 21.7 | 28.2 |
| 7-8 | 0.246 | 0.259 | 0.230 | 0.223 | 32.6 | 32.1 | 27.9 | 24.8 |

Mínimo ΔL consecutivo en gris y simulaciones: **0.198** (umbral 0.10; el encargo pide 0.08). ΔE mínimo entre vecinas: normal 22.1 (piso 15), protan/deutan 21.7 (objetivo 8), tritan 20.4.

### 8.4 Hex simulados y luminosidad por serie

**Tema claro** (hex · L; en gris, L = Y^(1/3))

| Serie | Normal | Gris | Protanopia | Deuteranopia | Tritanopia |
|---|---|---|---|---|---|
| 1 | `#474698` · 0.440 | `#505050` · 0.431 | `#19539B` · 0.447 | `#0E4D96` · 0.426 | `#255868` · 0.434 |
| 2 | `#009A87` · 0.615 | `#898989` · 0.629 | `#939086` · 0.652 | `#828389` · 0.612 | `#009D94` · 0.627 |
| 3 | `#883500` · 0.439 | `#4F4F4F` · 0.427 | `#4D4300` · 0.382 | `#625600` · 0.450 | `#961F2D` · 0.444 |
| 4 | `#BF6D2C` · 0.616 | `#828282` · 0.605 | `#857623` · 0.565 | `#99892C` · 0.627 | `#D15B60` · 0.621 |
| 5 | `#812E63` · 0.439 | `#4E4E4E` · 0.422 | `#364565` · 0.391 | `#4E5361` · 0.442 | `#892D44` · 0.441 |
| 6 | `#0092B9` · 0.614 | `#878787` · 0.622 | `#7C8FBB` · 0.652 | `#6780B9` · 0.604 | `#009D9F` · 0.632 |
| 7 | `#683987` · 0.439 | `#4E4E4E` · 0.424 | `#204C8A` · 0.420 | `#314F85` · 0.432 | `#64475A` · 0.436 |
| 8 | `#579749` · 0.615 | `#888888` · 0.627 | `#9B8C42` · 0.637 | `#92864E` · 0.617 | `#529385` · 0.615 |

**Tema oscuro** (hex · L; en gris, L = Y^(1/3))

| Serie | Normal | Gris | Protanopia | Deuteranopia | Tritanopia |
|---|---|---|---|---|---|
| 1 | `#797DCD` · 0.621 | `#848484` · 0.614 | `#6187D0` · 0.628 | `#5C81CB` · 0.610 | `#5F8D9B` · 0.615 |
| 2 | `#61E4CD` · 0.840 | `#CFCFCF` · 0.855 | `#DBD8CC` · 0.881 | `#C7C9CF` · 0.836 | `#00E7DD` · 0.838 |
| 3 | `#C16D45` · 0.621 | `#838383` · 0.610 | `#847841` · 0.570 | `#988A44` · 0.631 | `#D25D64` · 0.626 |
| 4 | `#FFB989` · 0.840 | `#C8C8C8` · 0.833 | `#CDBF85` · 0.803 | `#DECE89` · 0.849 | `#FFACAD` · 0.824 |
| 5 | `#B86797` · 0.620 | `#818181` · 0.605 | `#6D7999` · 0.579 | `#828695` · 0.622 | `#C16779` · 0.622 |
| 6 | `#76D9FF` · 0.839 | `#CCCCCC` · 0.847 | `#C4D5FF` · 0.874 | `#B0C6FF` · 0.830 | `#00E4E5` · 0.833 |
| 7 | `#9C70BC` · 0.619 | `#828282` · 0.606 | `#6080BF` · 0.603 | `#6B83BA` · 0.613 | `#987B8D` · 0.617 |
| 8 | `#A0DE93` · 0.840 | `#CECECE` · 0.852 | `#E2D28E` · 0.861 | `#D9CC97` · 0.843 | `#9CD9CB` · 0.840 |

ΔL en gris entre los 28 pares (informativo por K27: el umbral duro es ΔL ≥ 0.10 entre series consecutivas, §8.3, con forma del marcador, patrón de trazo o relleno y etiqueta directa como refuerzo; bloqueo `criterios-deltal` resuelto):

**Tema claro**

| Serie | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 |
|---|---|---|---|---|---|---|---|---|
| 1 |  |  |  |  |  |  |  |  |
| 2 | 0.198 |  |  |  |  |  |  |  |
| 3 | 0.004 | 0.202 |  |  |  |  |  |  |
| 4 | 0.175 | 0.023 | 0.179 |  |  |  |  |  |
| 5 | 0.009 | 0.207 | 0.005 | 0.183 |  |  |  |  |
| 6 | 0.191 | 0.007 | 0.195 | 0.017 | 0.200 |  |  |  |
| 7 | 0.007 | 0.205 | 0.003 | 0.182 | 0.002 | 0.198 |  |  |
| 8 | 0.196 | 0.002 | 0.200 | 0.021 | 0.205 | 0.005 | 0.203 |  |

Pares con ΔL gris < 0.10: 12 de 28 (todos no consecutivos).

**Tema oscuro**

| Serie | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 |
|---|---|---|---|---|---|---|---|---|
| 1 |  |  |  |  |  |  |  |  |
| 2 | 0.241 |  |  |  |  |  |  |  |
| 3 | 0.004 | 0.245 |  |  |  |  |  |  |
| 4 | 0.219 | 0.022 | 0.223 |  |  |  |  |  |
| 5 | 0.009 | 0.250 | 0.005 | 0.228 |  |  |  |  |
| 6 | 0.233 | 0.008 | 0.237 | 0.014 | 0.242 |  |  |  |
| 7 | 0.008 | 0.248 | 0.004 | 0.227 | 0.001 | 0.240 |  |  |
| 8 | 0.238 | 0.003 | 0.242 | 0.019 | 0.247 | 0.005 | 0.246 |  |

Pares con ΔL gris < 0.10: 12 de 28 (todos no consecutivos).

### 8.5 Pares no consecutivos

El criterio U2 del plan §3.5 dice «ΔL ≥ 0.10 entre series». Entre **todas** las series es imposible con 8 series: harían falta 7 saltos de 0.10, un rango de 0.70 en L, y el contraste de 3:1 deja como mucho L_gris ≤ 0.641 en claro (rango [0, 0.641] aun aceptando negro) y L_gris ≥ 0.596 en oscuro (rango [0.596, 1], 0.404). Por eso el criterio se aplica, como pide el encargo, a series **consecutivas** del orden de uso, que son las que se tocan en barras apiladas, barras agrupadas y leyendas. Mínimos entre todos los pares (informativo): claro ΔE normal 6.7 y ΔL gris 0.002; oscuro ΔE normal 6.2 y ΔL gris 0.001. Consecuencia de uso: en dispersión, mapas o pequeños múltiplos (donde cualquier par puede quedar junto) se usan como máximo las 3 primeras series, y toda gráfica con 2 o más series lleva leyenda y, con ≤ 4, etiqueta directa; en impresión en gris o modo de alto contraste se activa la textura de trama.

Las tres primeras series pasan la prueba de todos los pares de la skill dataviz (`--pairs all`): ΔE mínimo entre ellas de 22.7 en visión normal y 19.1 con protanopia o deuteranopia en claro, y de 21.3 y 20.1 en oscuro (pisos 15 y 8). El script lo verifica.

### 8.6 Robustez frente a la implementación de la simulación

Con la variante que aplica las matrices sobre sRGB codificado (la de culori 4.0.2), el mínimo ΔL consecutivo es 0.108 en claro y 0.178 en oscuro: la paleta también pasa el umbral 0.10 si U2b o U5 simulan con culori sin linealizar. Recomendación para ellos: convertir a `lrgb` antes de aplicar la matriz, para coincidir con este documento.

### 8.7 Banda de luminosidad de la skill dataviz

La skill propone L 0.43–0.77 en claro y 0.48–0.67 en oscuro, calibrada para su superficie oscura `#1a1a19`. Con la tarjeta 900 de este sistema (L 0.335), 3:1 exige L ≥ 0.596, así que su banda oscura dejaría solo [0.60, 0.67], insuficiente para alternar con ΔL ≥ 0.10. Se adoptan sus demás controles (croma ≥ 0.10, ΔE de vecinas con daltonismo ≥ 8 y en visión normal ≥ 15, contraste ≥ 3:1) y se reemplaza la banda por los límites de contraste de §8.2.

### 8.8 Matiz reservado y elección de la serie 6

Decisiones §I reserva el eje opuesto 100° para «acción crítica» de confirmación, y la zona reservada de §3 (80°–120°) vale también para las series. Hasta la ronda 1 de revisión la serie 6 era «mostaza» 100° (`#978500` en claro, `#DDCC6B` en oscuro), dentro de esa zona y sin excepción declarada (hallazgo H-03 de la revisión independiente de WF-0). Para elegir su nuevo matiz, el script barre los múltiplos de 5° fuera de la zona (65 candidatos; 7 quedan dentro de ella) con la clase de la serie (A) y las otras siete series en su lugar, y repite los controles de esta sección en ambos temas: croma final ≥ 0.10, matiz medido en el hex también fuera de la zona, contraste ≥ 3:1 sobre fondo y tarjeta, ΔL ≥ 0.10 con las vecinas (series 5 y 7) en gris y en las tres simulaciones con las dos variantes de §8.6, y ΔE de vecinas ≥ 15 en visión normal y ≥ 8 con protanopia y deuteranopia. Pasan 35. Descartes por motivo (un matiz puede tener varios): croma final < 0.10: 28; matiz medido dentro de la zona reservada: 2.

Entre los que pasan se elige el que más se aleja en matiz de las otras siete series: una paleta categórica identifica cada serie por una familia de matiz distinta (control 1 de la skill dataviz, ocho familias en orden fijo), y las series de la misma clase solo se distinguen por matiz cuando no son vecinas (§8.5). El desempate es el mayor ΔL con las vecinas. Los cinco primeros:

| Puesto | Matiz | Claro: hex · C | Oscuro: hex · C | Δ matiz mín. a las otras series | ΔL mín. con vecinas | ΔE normal mín. con vecinas | Elección |
|---|---|---|---|---|---|---|---|
| 1 | 225° | `#0092B9` · C 0.117 | `#76D9FF` · C 0.107 | 45.3° | 0.162 | 24.2 | **elegido** |
| 2 | 220° | `#0093B3` · C 0.112 | `#65DCFF` · C 0.117 | 40.3° | 0.163 | 24.5 |  |
| 3 | 215° | `#0095AD` · C 0.109 | `#5BDEFB` · C 0.120 | 34.5° | 0.166 | 25.1 |  |
| 4 | 210° | `#0095A8` · C 0.106 | `#58E0F6` · C 0.120 | 30.2° | 0.165 | 25.2 |  |
| 5 | 205° | `#0096A3` · C 0.105 | `#55E1F0` · C 0.120 | 25.4° | 0.166 | 25.6 |  |

Se usa el puesto 1, 225° («azul cielo»): queda a 45.3° de la serie más cercana en matiz y a 125.0° del opuesto. El descartado con más separación de matiz es 230° (49.9°), que falla por croma final < 0.10. Como el resto de las series, la 6 no se separa de los roles: queda a 14° de información (rol); la serie 7 (violeta) comparte matiz con deuda y la 8 (verde) queda a 5° de ingreso. Los roles viven en cifras, insignias y avisos, siempre con signo, icono o texto; las series viven en gráficas con leyenda, y cuando una serie significa un estado (por ejemplo, presupuesto excedido) usa el token de rol y no el de serie (regla de la skill dataviz).

## 9. Reglas de uso por componente

| Elemento | Claro | Oscuro | Contraste mínimo claro / oscuro |
|---|---|---|---|
| Fondo de página | neutral 50 | neutral 950 | — |
| Tarjeta | neutral 100 | neutral 900 | — |
| Divisor decorativo | neutral 200 | neutral 800 | exento (decorativo) |
| Texto normal | neutral 900 | neutral 100 | 10.75 / 10.75 |
| Texto secundario | neutral 700 | neutral 300 | 5.03 / 7.97 |
| Texto deshabilitado | neutral 700 + relleno neutral 100 | neutral 400 + relleno neutral 900 | 5.03 / 6.08 |
| Marcador de posición | neutral 700 | neutral 300 | 5.39 / 10.42 |
| Enlace (subrayado) · hover | brand 700 · brand 800 | brand 300 · brand 200 | 5.16 / 7.93 |
| Botón primario · hover | relleno brand 700 · 800, texto blanco | relleno brand 300 · 200, texto brand 950 | 5.16 / 7.93 |
| Botón secundario · hover | borde neutral 600, texto neutral 900 · relleno neutral 200 | borde neutral 400, texto neutral 100 · relleno neutral 800 | 3.46 / 6.08 |
| Botón destructivo · hover | relleno error 700 · 800, texto blanco | relleno error 300 · 200, texto error 950 | 5.39 / 7.75 |
| Botón de acción crítica · hover | relleno critical 700 · 800, texto blanco | relleno critical 300 · 200, texto critical 950 | 5.00 / 8.06 |
| Borde de campo | neutral 600 | neutral 400 | 3.46 / 6.08 |
| Anillo de foco (2 px, separación 2 px) | brand 600 | brand 400 | 3.56 / 6.03 |
| Elemento activo, fila seleccionada | relleno brand 100, texto neutral 900, icono brand 700 | relleno brand 900, texto neutral 100, icono brand 300 | 5.14 / 8.04 |
| Insignia de rol | texto rol 700 sobre rol 100 | texto rol 400 sobre rol 900 | 4.81 / 6.02 |
| Cifra de rol (con signo e icono) | rol 700 | rol 400 | 4.73 / 5.81 |
| Series de gráficas | `color.chart.<n>.light` | `color.chart.<n>.dark` | 3.14 / 3.10 |
| Personas y tiendas | `color.person|store.*.light` | `color.person|store.*.dark` | 3.05 / 3.06 |

Notas: (1) «texto 950» del botón primario en oscuro se lee como brand 950, el 950 de la misma escala del relleno; con neutral 950 daría 10.37:1, también válido. (2) WCAG 1.4.3 exime al texto de controles deshabilitados, pero el encargo exige 4.5:1 a todo texto; en claro no hay paso entre 600 (3.71:1 sobre fondo, insuficiente) y 700, así que el deshabilitado usa el mismo neutral 700 que el secundario y se reconoce por el relleno neutral 100, el borde punteado y el cursor. (3) Los divisores y bordes de tarjeta (neutral 200 en claro, 800 en oscuro) son decorativos: los campos y botones tienen su propio borde de 3:1. (4) Los enlaces van subrayados para no depender del color frente al texto normal.

## 10. Contraste (resumen de `contraste.csv`)

Columnas: `primer_plano`, `fondo` (rutas de token), `tema`, `uso`, `ratio` (truncada a 2 decimales), `umbral`, `cumple`, y además `primer_plano_hex` y `fondo_hex` para volver a medir sin el JSON.

| Tema | Tipo | Pares | Cumplen | Razón mínima |
|---|---|---|---|---|
| claro | texto | 51 | 51 | 4.73 |
| claro | interfaz y gráficas | 50 | 50 | 3.05 |
| oscuro | texto | 51 | 51 | 5.81 |
| oscuro | interfaz y gráficas | 50 | 50 | 3.06 |

Total: **202 pares, 202 cumplen, 0 bajo umbral.** Par de texto más ajustado: Cifra o texto de éxito sobre tarjeta (claro, 4.73:1). Par de interfaz más ajustado: Tienda Sam's Club (punto, anillo o serie) contra tarjeta (claro, 3.05:1).

**Pares estructurales del plan §4.3, recalculados** (tolerancia ±0.1; si la diferencia es mayor se documentan ambos valores):

| Par | Hex | Plan | Medido | Tolerancia |
|---|---|---|---|---|
| Texto 900 / fondo 50 | `#36363A` / `#FAFAFB` | 11.4 | 11.53 | **difiere 0.13** |
| Texto secundario 700 / fondo 50 | `#66676E` / `#FAFAFB` | 5.3 | 5.39 | dentro |
| Neutro 600 / tarjeta 100 (solo texto grande o iconos; aquí no se usa para texto) | `#808189` / `#F2F2F4` | 3.45 | 3.46 | dentro |
| Blanco / marca 600 (no se usa para texto de botón) | `#FFFFFF` / `#7579C2` | 3.7–4.0 (3.98 en decisiones §I) | 3.98 | dentro |
| Marca 700 / fondo 50 | `#5D609F` / `#FAFAFB` | 5.1–5.5 | 5.53 | dentro |
| Blanco / marca 700 (texto del botón primario) | `#FFFFFF` / `#5D609F` | 5.8 (criterios U2-07) | 5.77 | dentro |
| Texto 100 / fondo 950 | `#F2F2F4` / `#222325` | 14.1 | 14.06 | dentro |
| Texto 300 / tarjeta 900 | `#D1D2D8` / `#36363A` | 7.9 | 7.97 | dentro |
| Marca 300 / fondo 950 | `#C9CFFF` / `#222325` | 10.4–10.7 | 10.37 | dentro |
| Marca 950 / marca 300 (texto del botón primario en oscuro) | `#1E203A` / `#C9CFFF` | 10.4 (criterios U2-07) | 10.47 | dentro |
| Divisor 200 / fondo 50 (decorativo) | `#E5E6E9` / `#FAFAFB` | 1.2 | 1.19 | dentro |
| Divisor 800 / fondo 950 (decorativo) | `#4D4E53` / `#222325` | 1.9 | 1.89 | dentro |
| Borde de campo 600 / fondo 50 | `#808189` / `#FAFAFB` | ≥ 4.2 | 3.71 | **difiere 0.49** |

Donde la diferencia supera 0.1 se documentan ambos valores y rige el medido. En particular, el borde de campo neutral 600 sobre fondo 50 da 3.71:1, no ≥ 4.2:1 como indicaba el plan; cumple el 3:1 que exige el borde de un componente (criterio 1.4.11).

## 11. Tokens (`C2-colores.json`)

Formato DTCG con `$value` en hex sRGB, como pide el encargo. Cada token lleva `$type: "color"`, `$value` y `$extensions.oklch` con `l`, `c` (6 decimales) y `h` de diseño, antes de cuantizar a 8 bits; `c` es el croma ya recortado al gamut. Con 6 decimales, convertir esos valores a hex reproduce exactamente `$value` (con 4 decimales un canal de `color.role.warning.50` cruza el medio punto de redondeo). Grupos:

| Grupo | Tokens | Contenido |
|---|---|---|
| `color.base.white` | 1 | blanco puro para texto sobre rellenos 700 en claro |
| `color.brand.50…950` | 11 | marca 280° |
| `color.neutral.50…950` | 11 | neutro tintado 280° |
| `color.role.<rol>.{light,dark,surface-light,surface-dark,50…950}` | 120 | 8 roles: income, expense, savings, debt, warning, error, success, info |
| `color.critical.{light,dark,surface-light,surface-dark,50…950}` | 15 | acción crítica, eje opuesto 100° |
| `color.person.<persona>.{light,dark}` | 8 | axel, ximena, karen, shared (conjunto) |
| `color.store.<tienda>.{light,dark}` | 8 | walmart, bodega-aurrera, sams-club, market (mercado) |
| `color.chart.<1…8>.{light,dark}` | 16 | series de gráficas en orden de uso |

`light`/`dark` de cada rol duplican el paso elegido (700/400) y `surface-*` los pasos 100/900, para que U2b mapee variables por tema sin conocer la regla. Los nombres de código de roles, personas y tiendas están en inglés (protocolo §8); si el diccionario de datos de S1 fija otros códigos, se agregan alias en WF-1 sin cambiar valores.

## 12. Hallazgos, decisiones derivadas del cálculo y bloqueos

**Hallazgos**

1. Los 22 valores de partida del plan §4.3 (marca y neutro) se reproducen exactamente con la tabla de perfil de §2.3; el plan no la explicitaba.
2. Los valores de rol del plan (L 0.56) quedan entre 3.92 y 4.47:1 sobre las superficies claras; 7 de 7 quedan bajo 4.5:1 en al menos una de las dos. Se usan pasos de escala (claro 700, oscuro 400); el paso 600 que pedía el encargo da 3.26–3.70:1 (§6.2).
3. Versión 1.1 (K26): gasto pasa de 25° a 45° y deja de confundirse con error 20° (ΔE normal 5.5 en claro y 5.1 en oscuro; en la 1.0 eran 1.2 y 1.0). Bloqueo `U2-gasto-error` resuelto (opción 2).
4. Además de las dos colisiones listadas en decisiones §I, la regla de separación detecta tres más en la propuesta categórica del plan: Axel–información, Mercado–gasto/error y Walmart en el límite exacto de 20°. Las cinco quedan resueltas (§7.5).
5. 4 colores categóricos del plan no llegaban a 3:1 sobre las superficies claras y 3 sobre las oscuras; no existe un valor único que cumpla 3:1 en ambos temas y además se escalone (ventana L_gris [0.596, 0.641]). Personas, tiendas y series tienen variante por tema (§7, §8).
6. «ΔL ≥ 0.10 entre series» (plan §3.5) solo es alcanzable entre series consecutivas; entre todos los pares es matemáticamente imposible con 8 series (§8.5). K27 fija esa lectura: umbral duro entre series consecutivas y los 28 pares como información (bloqueo `criterios-deltal` resuelto).
7. culori 4.0.2 aplica las matrices de Machado sobre sRGB codificado, no lineal; la paleta pasa con ambas variantes (§8.6).
8. Pares de roles con ΔE normal < 15 (ingreso–éxito, advertencia–acción crítica, entre otros, §6.6): ninguno comparte componente y todos llevan icono y texto.
9. Ronda 1 de revisión (H-03): la serie 6 de gráficas («mostaza», 100°) usaba el matiz que decisiones §I reserva para acción crítica y el documento no declaraba la excepción. Pasa a «azul cielo» 225° (§8.8); con ese matiz, la búsqueda de orden de §8.2 pone óxido en la serie 3 y magenta en la 5. El script verifica la zona reservada 80°–120° en todos los tokens (§3).
10. Encontrado al corregir H-03: con el orden de la ronda 0, las tres primeras series (índigo, verde azulado, magenta), que §8.5 reserva para dispersión, mapas y pequeños múltiplos, no llegaban al piso de ΔE 15 entre todos sus pares en visión normal (13.9 en claro y 12.8 en oscuro, entre índigo y magenta, ambas de clase B y con la misma luminosidad). Con el orden vigente dan 22.7 y 21.3, y el script lo verifica (§8.5).
11. Versión 1.1 (K26): con gasto en 45° no queda hueco cálido a ≥ 20° de error, gasto y advertencia, y la L de Karen y de Bodega Aurrera las obliga a separarse de los roles por matiz; sus ventanas factibles quedan entre deuda y error y van a 333° y 356° (Cambios 1.1).
12. Versión 1.1 (R3-04): con la marca en el barrido de §7.7 quedan cerca Conjunto (persona) y Marca 700 (marca) en claro (0.9°, ΔL 0.094); Walmart (tienda) y Marca 700 (marca) en claro (18.0°, ΔL 0.000); Walmart (tienda) y Marca 400 (marca) en oscuro (18.0°, ΔL 0.069); todos con refuerzo de forma y etiqueta, y ningún rol cerca de la marca.
13. Versión 1.1: en oscuro, Bodega Aurrera en rosa no cabe en el escalonamiento ΔL ≥ 0.10 de las tiendas con Walmart en L 0.720 (ningún matiz entero cumple); Walmart pasa a L 0.716, PROVISIONAL con el bloqueo `U2a-walmart-oscuro`, y Bodega Aurrera usa croma 0.11 (Cambios 1.1).
14. Ronda 1 de revisión U2a (H-01): en oscuro, Bodega Aurrera y Ximena quedan a 6.2°, ΔL 0.006 y ΔE 1.3 (a simple vista, el mismo color; en la 1.0, el par persona–tienda más parecido tenía ΔE 10.5), y en claro a 5.8°, ΔL 0.130 y ΔE 13.1. Karen y Ximena quedan a 17.1° y 16.7°, con ΔL 0.103 y 0.111. Con las L actuales, ningún matiz de sus ventanas las separa de Ximena por ≥ 20° o ΔL ≥ 0.15. §7.7 declara el refuerzo, pero U2-02 admite dos lecturas: bloqueo `U2a-rosas-ximena-bodega`, con 4 tokens PROVISIONALES (Cambios 1.1, punto 2).

**Decisiones derivadas del cálculo** (no son supuestos: las fuerza un requisito explícito)

- Paso 700 en claro para los roles (el 600 incumple 4.5:1).
- Variantes `light`/`dark` para personas, tiendas y series (un valor único no cumple 3:1 en ambos temas con escalonamiento).
- Criterio de ΔL aplicado a series consecutivas, con umbral 0.10 (más estricto que el 0.08 del encargo).
- Texto deshabilitado en claro con neutral 700 (el 600 incumple 4.5:1).
- Versión 1.1: Bodega Aurrera con croma 0.11 y Karen con L 0.642 en claro (márgenes del escalonamiento de su grupo; Cambios 1.1).

**Propuestas sujetas a la puerta 2:** ingreso 135° (§6.4, opción prevista en decisiones §I), la resolución de colisiones categóricas (§7.5; en la versión 1.1, gasto 45°, Karen y Bodega Aurrera en sus matices nuevos, «Cambios 1.1») y la paleta de gráficas corregida en la ronda 1 (serie 6 en 225°, con óxido en la serie 3 y magenta en la 5; §8.2 y §8.8).

**Bloqueos:** `docs/estado/bloqueos/U2a-walmart-oscuro.md` (versión 1.1: Walmart en oscuro con L 0.716 en vez de 0.720, PROVISIONAL; «Cambios 1.1») y `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md` (ronda 1 de revisión, H-01: Karen y Bodega Aurrera bajo ≥ 20° o ΔL ≥ 0.15 frente a Ximena; 4 tokens PROVISIONALES; «Cambios 1.1», punto 2). Resueltos en la puerta 1: `U2-gasto-error` (K26) y `criterios-deltal` (K27).

## Anexo A. Tokens: ruta → hex

Lista completa y verificable de `C2-colores.json` (el script compara esta tabla con el JSON en ambos sentidos).

| Token | Hex | L | C | h |
|---|---|---|---|---|
| `color.base.white` | `#FFFFFF` | 1.000000 | 0.000000 | 0.00 |
| `color.brand.50` | `#F9FAFF` | 0.985000 | 0.007136 | 280.00 |
| `color.brand.100` | `#F0F1FF` | 0.962000 | 0.018245 | 280.00 |
| `color.brand.200` | `#E1E4FF` | 0.925000 | 0.036551 | 280.00 |
| `color.brand.300` | `#C9CFFF` | 0.865000 | 0.067413 | 280.00 |
| `color.brand.400` | `#ADB3F4` | 0.785000 | 0.093500 | 280.00 |
| `color.brand.500` | `#9095E0` | 0.695000 | 0.110000 | 280.00 |
| `color.brand.600` | `#7579C2` | 0.605000 | 0.110000 | 280.00 |
| `color.brand.700` | `#5D609F` | 0.515000 | 0.099000 | 280.00 |
| `color.brand.800` | `#46487B` | 0.425000 | 0.082500 | 280.00 |
| `color.brand.900` | `#303258` | 0.335000 | 0.066000 | 280.00 |
| `color.brand.950` | `#1E203A` | 0.255000 | 0.049500 | 280.00 |
| `color.neutral.50` | `#FAFAFB` | 0.985000 | 0.001200 | 280.00 |
| `color.neutral.100` | `#F2F2F4` | 0.962000 | 0.003000 | 280.00 |
| `color.neutral.200` | `#E5E6E9` | 0.925000 | 0.004800 | 280.00 |
| `color.neutral.300` | `#D1D2D8` | 0.865000 | 0.007800 | 280.00 |
| `color.neutral.400` | `#B7B8C0` | 0.785000 | 0.010200 | 280.00 |
| `color.neutral.500` | `#9B9CA4` | 0.695000 | 0.012000 | 280.00 |
| `color.neutral.600` | `#808189` | 0.605000 | 0.012000 | 280.00 |
| `color.neutral.700` | `#66676E` | 0.515000 | 0.010800 | 280.00 |
| `color.neutral.800` | `#4D4E53` | 0.425000 | 0.009000 | 280.00 |
| `color.neutral.900` | `#36363A` | 0.335000 | 0.007200 | 280.00 |
| `color.neutral.950` | `#222325` | 0.255000 | 0.005400 | 280.00 |
| `color.role.income.light` | `#457723` | 0.515000 | 0.126000 | 135.00 |
| `color.role.income.dark` | `#98CA7B` | 0.785000 | 0.119000 | 135.00 |
| `color.role.income.surface-light` | `#E9F8E1` | 0.962000 | 0.035000 | 135.00 |
| `color.role.income.surface-dark` | `#22400D` | 0.335000 | 0.084000 | 135.00 |
| `color.role.income.50` | `#F6FCF3` | 0.985000 | 0.014000 | 135.00 |
| `color.role.income.100` | `#E9F8E1` | 0.962000 | 0.035000 | 135.00 |
| `color.role.income.200` | `#D7EFCA` | 0.925000 | 0.056000 | 135.00 |
| `color.role.income.300` | `#B9E1A4` | 0.865000 | 0.091000 | 135.00 |
| `color.role.income.400` | `#98CA7B` | 0.785000 | 0.119000 | 135.00 |
| `color.role.income.500` | `#75B052` | 0.695000 | 0.140000 | 135.00 |
| `color.role.income.600` | `#5B9335` | 0.605000 | 0.140000 | 135.00 |
| `color.role.income.700` | `#457723` | 0.515000 | 0.126000 | 135.00 |
| `color.role.income.800` | `#335A18` | 0.425000 | 0.105000 | 135.00 |
| `color.role.income.900` | `#22400D` | 0.335000 | 0.084000 | 135.00 |
| `color.role.income.950` | `#142907` | 0.255000 | 0.063000 | 135.00 |
| `color.role.expense.light` | `#A04B20` | 0.515000 | 0.126000 | 45.00 |
| `color.role.expense.dark` | `#F8A078` | 0.785000 | 0.119000 | 45.00 |
| `color.role.expense.surface-light` | `#FFEFE8` | 0.962000 | 0.020215 | 45.00 |
| `color.role.expense.surface-dark` | `#59260B` | 0.335000 | 0.084000 | 45.00 |
| `color.role.expense.50` | `#FFF9F6` | 0.985000 | 0.007839 | 45.00 |
| `color.role.expense.100` | `#FFEFE8` | 0.962000 | 0.020215 | 45.00 |
| `color.role.expense.200` | `#FFDED0` | 0.925000 | 0.041095 | 45.00 |
| `color.role.expense.300` | `#FFC3A9` | 0.865000 | 0.077809 | 45.00 |
| `color.role.expense.400` | `#F8A078` | 0.785000 | 0.119000 | 45.00 |
| `color.role.expense.500` | `#E37E4E` | 0.695000 | 0.140000 | 45.00 |
| `color.role.expense.600` | `#C46231` | 0.605000 | 0.140000 | 45.00 |
| `color.role.expense.700` | `#A04B20` | 0.515000 | 0.126000 | 45.00 |
| `color.role.expense.800` | `#7C3815` | 0.425000 | 0.105000 | 45.00 |
| `color.role.expense.900` | `#59260B` | 0.335000 | 0.084000 | 45.00 |
| `color.role.expense.950` | `#3B1705` | 0.255000 | 0.063000 | 45.00 |
| `color.role.savings.light` | `#00767B` | 0.515000 | 0.087553 | 200.00 |
| `color.role.savings.dark` | `#3FD0D7` | 0.785000 | 0.119000 | 200.00 |
| `color.role.savings.surface-light` | `#D8FAFC` | 0.962000 | 0.035000 | 200.00 |
| `color.role.savings.surface-dark` | `#003F42` | 0.335000 | 0.056952 | 200.00 |
| `color.role.savings.50` | `#F0FDFE` | 0.985000 | 0.014000 | 200.00 |
| `color.role.savings.100` | `#D8FAFC` | 0.962000 | 0.035000 | 200.00 |
| `color.role.savings.200` | `#BBF2F5` | 0.925000 | 0.056000 | 200.00 |
| `color.role.savings.300` | `#85E5EA` | 0.865000 | 0.091000 | 200.00 |
| `color.role.savings.400` | `#3FD0D7` | 0.785000 | 0.119000 | 200.00 |
| `color.role.savings.500` | `#00B2BA` | 0.695000 | 0.118154 | 200.00 |
| `color.role.savings.600` | `#00949A` | 0.605000 | 0.102853 | 200.00 |
| `color.role.savings.700` | `#00767B` | 0.515000 | 0.087553 | 200.00 |
| `color.role.savings.800` | `#005A5E` | 0.425000 | 0.072252 | 200.00 |
| `color.role.savings.900` | `#003F42` | 0.335000 | 0.056952 | 200.00 |
| `color.role.savings.950` | `#00292B` | 0.255000 | 0.043351 | 200.00 |
| `color.role.debt.light` | `#7E509D` | 0.515000 | 0.126000 | 310.00 |
| `color.role.debt.dark` | `#D0A3F2` | 0.785000 | 0.119000 | 310.00 |
| `color.role.debt.surface-light` | `#F7EEFF` | 0.962000 | 0.024076 | 310.00 |
| `color.role.debt.surface-dark` | `#442857` | 0.335000 | 0.084000 | 310.00 |
| `color.role.debt.50` | `#FCF8FF` | 0.985000 | 0.009411 | 310.00 |
| `color.role.debt.100` | `#F7EEFF` | 0.962000 | 0.024076 | 310.00 |
| `color.role.debt.200` | `#F0DEFF` | 0.925000 | 0.048281 | 310.00 |
| `color.role.debt.300` | `#E5C3FF` | 0.865000 | 0.089202 | 310.00 |
| `color.role.debt.400` | `#D0A3F2` | 0.785000 | 0.119000 | 310.00 |
| `color.role.debt.500` | `#B782DD` | 0.695000 | 0.140000 | 310.00 |
| `color.role.debt.600` | `#9B67C0` | 0.605000 | 0.140000 | 310.00 |
| `color.role.debt.700` | `#7E509D` | 0.515000 | 0.126000 | 310.00 |
| `color.role.debt.800` | `#603C79` | 0.425000 | 0.105000 | 310.00 |
| `color.role.debt.900` | `#442857` | 0.335000 | 0.084000 | 310.00 |
| `color.role.debt.950` | `#2C1939` | 0.255000 | 0.063000 | 310.00 |
| `color.role.warning.light` | `#8B5D00` | 0.515000 | 0.108645 | 75.00 |
| `color.role.warning.dark` | `#E5AE5B` | 0.785000 | 0.119000 | 75.00 |
| `color.role.warning.surface-light` | `#FFF0DC` | 0.962000 | 0.030581 | 75.00 |
| `color.role.warning.surface-dark` | `#4C3000` | 0.335000 | 0.070672 | 75.00 |
| `color.role.warning.50` | `#FFF9F1` | 0.985000 | 0.011918 | 75.00 |
| `color.role.warning.100` | `#FFF0DC` | 0.962000 | 0.030581 | 75.00 |
| `color.role.warning.200` | `#FDE2BD` | 0.925000 | 0.056000 | 75.00 |
| `color.role.warning.300` | `#F6CB8F` | 0.865000 | 0.091000 | 75.00 |
| `color.role.warning.400` | `#E5AE5B` | 0.785000 | 0.119000 | 75.00 |
| `color.role.warning.500` | `#CF8E1B` | 0.695000 | 0.140000 | 75.00 |
| `color.role.warning.600` | `#AD7500` | 0.605000 | 0.127631 | 75.00 |
| `color.role.warning.700` | `#8B5D00` | 0.515000 | 0.108645 | 75.00 |
| `color.role.warning.800` | `#6B4600` | 0.425000 | 0.089658 | 75.00 |
| `color.role.warning.900` | `#4C3000` | 0.335000 | 0.070672 | 75.00 |
| `color.role.warning.950` | `#321E00` | 0.255000 | 0.053795 | 75.00 |
| `color.role.error.light` | `#A34449` | 0.515000 | 0.126000 | 20.00 |
| `color.role.error.dark` | `#FC9999` | 0.785000 | 0.119000 | 20.00 |
| `color.role.error.surface-light` | `#FFEEED` | 0.962000 | 0.018693 | 20.00 |
| `color.role.error.surface-dark` | `#5A2124` | 0.335000 | 0.084000 | 20.00 |
| `color.role.error.50` | `#FFF8F8` | 0.985000 | 0.007238 | 20.00 |
| `color.role.error.100` | `#FFEEED` | 0.962000 | 0.018693 | 20.00 |
| `color.role.error.200` | `#FFDDDC` | 0.925000 | 0.038101 | 20.00 |
| `color.role.error.300` | `#FFC0BF` | 0.865000 | 0.072486 | 20.00 |
| `color.role.error.400` | `#FC9999` | 0.785000 | 0.119000 | 20.00 |
| `color.role.error.500` | `#E77678` | 0.695000 | 0.140000 | 20.00 |
| `color.role.error.600` | `#C75A5E` | 0.605000 | 0.140000 | 20.00 |
| `color.role.error.700` | `#A34449` | 0.515000 | 0.126000 | 20.00 |
| `color.role.error.800` | `#7E3236` | 0.425000 | 0.105000 | 20.00 |
| `color.role.error.900` | `#5A2124` | 0.335000 | 0.084000 | 20.00 |
| `color.role.error.950` | `#3C1416` | 0.255000 | 0.063000 | 20.00 |
| `color.role.success.light` | `#207B3E` | 0.515000 | 0.126000 | 150.00 |
| `color.role.success.dark` | `#7ECF90` | 0.785000 | 0.119000 | 150.00 |
| `color.role.success.surface-light` | `#E3FAE6` | 0.962000 | 0.035000 | 150.00 |
| `color.role.success.surface-dark` | `#0A421E` | 0.335000 | 0.084000 | 150.00 |
| `color.role.success.50` | `#F4FDF5` | 0.985000 | 0.014000 | 150.00 |
| `color.role.success.100` | `#E3FAE6` | 0.962000 | 0.035000 | 150.00 |
| `color.role.success.200` | `#CDF1D3` | 0.925000 | 0.056000 | 150.00 |
| `color.role.success.300` | `#A8E4B3` | 0.865000 | 0.091000 | 150.00 |
| `color.role.success.400` | `#7ECF90` | 0.785000 | 0.119000 | 150.00 |
| `color.role.success.500` | `#53B56D` | 0.695000 | 0.140000 | 150.00 |
| `color.role.success.600` | `#339852` | 0.605000 | 0.140000 | 150.00 |
| `color.role.success.700` | `#207B3E` | 0.515000 | 0.126000 | 150.00 |
| `color.role.success.800` | `#155E2D` | 0.425000 | 0.105000 | 150.00 |
| `color.role.success.900` | `#0A421E` | 0.335000 | 0.084000 | 150.00 |
| `color.role.success.950` | `#052B11` | 0.255000 | 0.063000 | 150.00 |
| `color.role.info.light` | `#006FA4` | 0.515000 | 0.117966 | 240.00 |
| `color.role.info.dark` | `#6BC3FD` | 0.785000 | 0.119000 | 240.00 |
| `color.role.info.surface-light` | `#E7F5FF` | 0.962000 | 0.020158 | 240.00 |
| `color.role.info.surface-dark` | `#003B5A` | 0.335000 | 0.076735 | 240.00 |
| `color.role.info.50` | `#F6FBFF` | 0.985000 | 0.007889 | 240.00 |
| `color.role.info.100` | `#E7F5FF` | 0.962000 | 0.020158 | 240.00 |
| `color.role.info.200` | `#CFEBFF` | 0.925000 | 0.040341 | 240.00 |
| `color.role.info.300` | `#A6DAFF` | 0.865000 | 0.074271 | 240.00 |
| `color.role.info.400` | `#6BC3FD` | 0.785000 | 0.119000 | 240.00 |
| `color.role.info.500` | `#35A7EA` | 0.695000 | 0.140000 | 240.00 |
| `color.role.info.600` | `#008ACB` | 0.605000 | 0.138582 | 240.00 |
| `color.role.info.700` | `#006FA4` | 0.515000 | 0.117966 | 240.00 |
| `color.role.info.800` | `#00547E` | 0.425000 | 0.097351 | 240.00 |
| `color.role.info.900` | `#003B5A` | 0.335000 | 0.076735 | 240.00 |
| `color.role.info.950` | `#00263C` | 0.255000 | 0.058410 | 240.00 |
| `color.critical.light` | `#766800` | 0.515000 | 0.107068 | 100.00 |
| `color.critical.dark` | `#CBBB5A` | 0.785000 | 0.119000 | 100.00 |
| `color.critical.surface-light` | `#F8F4D9` | 0.962000 | 0.035000 | 100.00 |
| `color.critical.surface-dark` | `#3F3700` | 0.335000 | 0.069646 | 100.00 |
| `color.critical.50` | `#FCFBF0` | 0.985000 | 0.014000 | 100.00 |
| `color.critical.100` | `#F8F4D9` | 0.962000 | 0.035000 | 100.00 |
| `color.critical.200` | `#EEE8BD` | 0.925000 | 0.056000 | 100.00 |
| `color.critical.300` | `#E0D58E` | 0.865000 | 0.091000 | 100.00 |
| `color.critical.400` | `#CBBB5A` | 0.785000 | 0.119000 | 100.00 |
| `color.critical.500` | `#B29E15` | 0.695000 | 0.140000 | 100.00 |
| `color.critical.600` | `#948200` | 0.605000 | 0.125779 | 100.00 |
| `color.critical.700` | `#766800` | 0.515000 | 0.107068 | 100.00 |
| `color.critical.800` | `#5A4F00` | 0.425000 | 0.088357 | 100.00 |
| `color.critical.900` | `#3F3700` | 0.335000 | 0.069646 | 100.00 |
| `color.critical.950` | `#292300` | 0.255000 | 0.053014 | 100.00 |
| `color.person.axel.light` | `#002E64` | 0.310000 | 0.105844 | 256.00 |
| `color.person.axel.dark` | `#4F83C8` | 0.605000 | 0.120000 | 256.00 |
| `color.person.ximena.light` | `#A14E78` | 0.540000 | 0.120000 | 350.00 |
| `color.person.ximena.dark` | `#FFACD3` | 0.835000 | 0.108508 | 350.00 |
| `color.person.karen.light` | `#B870AC` | 0.642000 | 0.120000 | 333.00 |
| `color.person.karen.dark` | `#D389C6` | 0.725000 | 0.120000 | 333.00 |
| `color.person.shared.light` | `#484B66` | 0.420000 | 0.045000 | 280.00 |
| `color.person.shared.dark` | `#E7EAFF` | 0.940000 | 0.029064 | 280.00 |
| `color.store.walmart.light` | `#4065AC` | 0.515000 | 0.120000 | 262.00 |
| `color.store.walmart.dark` | `#7AA3EE` | 0.716000 | 0.120000 | 262.00 |
| `color.store.bodega-aurrera.light` | `#762C4C` | 0.410000 | 0.110000 | 356.00 |
| `color.store.bodega-aurrera.dark` | `#FFAAC9` | 0.829000 | 0.106641 | 356.00 |
| `color.store.sams-club.light` | `#009D82` | 0.620000 | 0.116902 | 175.00 |
| `color.store.sams-club.dark` | `#99FFE4` | 0.930000 | 0.103566 | 175.00 |
| `color.store.market.light` | `#5C0204` | 0.300000 | 0.120000 | 28.00 |
| `color.store.market.dark` | `#C2665B` | 0.615000 | 0.120000 | 28.00 |
| `color.chart.1.light` | `#474698` | 0.440000 | 0.130000 | 280.00 |
| `color.chart.1.dark` | `#797DCD` | 0.620000 | 0.120000 | 280.00 |
| `color.chart.2.light` | `#009A87` | 0.615000 | 0.111579 | 180.00 |
| `color.chart.2.dark` | `#61E4CD` | 0.840000 | 0.120000 | 180.00 |
| `color.chart.3.light` | `#883500` | 0.440000 | 0.125796 | 45.00 |
| `color.chart.3.dark` | `#C16D45` | 0.620000 | 0.120000 | 45.00 |
| `color.chart.4.light` | `#BF6D2C` | 0.615000 | 0.130000 | 55.00 |
| `color.chart.4.dark` | `#FFB989` | 0.840000 | 0.102608 | 55.00 |
| `color.chart.5.light` | `#812E63` | 0.440000 | 0.130000 | 345.00 |
| `color.chart.5.dark` | `#B86797` | 0.620000 | 0.120000 | 345.00 |
| `color.chart.6.light` | `#0092B9` | 0.615000 | 0.116522 | 225.00 |
| `color.chart.6.dark` | `#76D9FF` | 0.840000 | 0.106925 | 225.00 |
| `color.chart.7.light` | `#683987` | 0.440000 | 0.130000 | 310.00 |
| `color.chart.7.dark` | `#9C70BC` | 0.620000 | 0.120000 | 310.00 |
| `color.chart.8.light` | `#579749` | 0.615000 | 0.130000 | 140.00 |
| `color.chart.8.dark` | `#A0DE93` | 0.840000 | 0.120000 | 140.00 |

