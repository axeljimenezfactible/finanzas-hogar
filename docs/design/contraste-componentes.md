# Contraste de componentes · resumen de U2d

Fecha: 11-oct-2026 · Agente U2d-contraste · Fuente de la lista: `DESIGN.md` §10.5 · Fuente de los colores: `docs/contratos/C2-tokens.json` (única; `C2-colores.json`, §7.1 y las columnas de §10.5 como cruce secundario) · Criterios: U2-04 a U2-07 y U5-04 (lado de diseño).

## 1. Resultado

**0 pares bajo umbral.** Los 158 pares de §10.5 se calcularon en claro y en oscuro (316 combinaciones): 304 con umbral (texto 4.5:1; texto grande e interfaz 3:1) y 12 decorativas exentas (6 pares × 2 temas, marcadas `exento` en el CSV con su razón calculada de todos modos). No hubo que abrir ningún bloqueo de contraste y **no se cambió ningún token**.

| Tipo | Umbral | Claro | Oscuro | Razón mínima (claro · oscuro) |
|---|---|---|---|---|
| Texto | 4.5 | 86 | 86 | 4.51 · 5.03 |
| Interfaz (iconos, bordes de campo, foco, relleno de botones, series, marcas de forma) | 3.0 | 65 | 65 | 3.14 · 3.06 |
| Texto grande (cifra de 24 px de la tarjeta de métrica) | 3.0 | 1 | 1 | 10.75 · 10.75 (C-096) |
| Decorativo (exento) | — | 6 | 6 | 1.07 · 1.30 (mínimos) |
| **Total** | | **158** | **158** | **316** |

Archivos: `docs/design/contraste-componentes.csv` (316 filas), `docs/design/tools/contraste_componentes.py` (lo regenera idéntico).

## 2. Cómo se calculó

- Fórmula WCAG 2.x: `(L1 + 0.05) / (L2 + 0.05)` con L1 la luminancia mayor; luminancia relativa sRGB `0.2126 R + 0.7152 G + 0.0722 B` con canales linealizados (`c/12.92` si c ≤ 0.04045; si no, `((c + 0.055)/1.055)^2.4`).
- Aritmética `Decimal` de 60 dígitos. La columna `razon` va **truncada a 2 decimales** (igual que `contraste.csv`, §10.5); la columna `cumple` se decide con el valor exacto, nunca con el truncado.
- El script lee la lista de pares de `DESIGN.md` §10.5 y resuelve cada lado y su hex solo desde `C2-tokens.json`; §7.1, las columnas de §10.5 y `C2-colores.json` son cruces secundarios que salen con código 1 si no coinciden. No importa `colorimetria.py`.
- Comprobaciones cruzadas que corre el script (todas dieron 0 diferencias):
  1. Cada token `ui.*`, «rol», «persona», «tienda», «serie», «acción crítica» y «superficie de rol» del campo «Par» se resuelve y debe coincidir con los tokens de las columnas Claro y Oscuro de §10.5 (632 lados).
  2. Los hex impresos en §7.1 coinciden con C2.
  3. Cruce con los 202 pares de `contraste.csv`: de las 316 combinaciones, 266 ya estaban y salen **iguales al centésimo**; 50 son nuevas (las que §10.5 marca «no»), como declara el propio §10.5. Ninguna difiere.
  4. Forma esperada: 158 / 316 / 304 / 12; el script sale con código 1 si cambia.

## 3. Pares que no cumplen

Ninguno.

## 4. Pares con poco margen (los 15 más justos)

No son fallas. Si U2 mueve un token, estos son los que se rompen primero.

| Id | Componente · estado | Tema | Frente / fondo | Razón | Umbral | Margen |
|---|---|---|---|---|---|---|
| C-138 | Menú «/» y búsqueda · resaltado · alias y atajo | claro | `#66676E` / `#E5E6E9` | 4.51 | 4.5 | +0.01 |
| C-057 | Chip de persona o tienda · reposo · marca de forma | oscuro | `#C2665B` / `#36363A` | 3.06 | 3.0 | +0.06 |
| C-050 | Chip de persona o tienda · reposo · marca de forma | oscuro | `#4F83C8` / `#36363A` | 3.10 | 3.0 | +0.10 |
| C-110 | Gráfica · serie 5 sobre tarjeta | oscuro | `#B86797` / `#36363A` | 3.10 | 3.0 | +0.10 |
| C-112 | Gráfica · serie 7 sobre tarjeta | oscuro | `#9C70BC` / `#36363A` | 3.12 | 3.0 | +0.12 |
| C-107 | Gráfica · serie 2 sobre tarjeta | claro | `#009A87` / `#F2F2F4` | 3.14 | 3.0 | +0.14 |
| C-121 | Calendario de pagos · pagado · icono y texto | claro | `#207B3E` / `#F2F2F4` | 4.73 | 4.5 | +0.23 |
| C-113 | Gráfica · serie 8 sobre tarjeta | claro | `#579749` / `#F2F2F4` | 3.16 | 3.0 | +0.16 |
| C-108 | Gráfica · serie 3 sobre tarjeta | oscuro | `#C16D45` / `#36363A` | 3.17 | 3.0 | +0.17 |
| C-083 | Tabla · seleccionada · cifra de ingreso | claro | `#457723` / `#F0F1FF` | 4.78 | 4.5 | +0.28 |
| C-048 | Edición en el lugar · hover · cifra de rol (ingreso) | claro | `#457723` / `#F2F2F4` | 4.79 | 4.5 | +0.29 |
| C-097 | Tarjeta de métrica · dato secundario de rol · income | claro | `#457723` / `#F2F2F4` | 4.79 | 4.5 | +0.29 |
| C-072 | Insignia de rol · reposo · success | claro | `#207B3E` / `#E3FAE6` | 4.81 | 4.5 | +0.31 |
| C-106 | Gráfica · serie 1 sobre tarjeta | oscuro | `#797DCD` / `#36363A` | 3.22 | 3.0 | +0.22 |
| C-066 | Insignia de rol · reposo · income | claro | `#457723` / `#E9F8E1` | 4.84 | 4.5 | +0.34 |

Lectura: los márgenes más chicos son (a) el texto secundario sobre el hover de los elementos del menú «/» (C-138, claro, 4.51 contra 4.5: **solo 0.01**), (b) en oscuro, las marcas de forma de persona/tienda sobre `ui.surface.raised` y las series 5 y 7 de la gráfica sobre tarjeta (3.06 a 3.12 contra 3.0) y (c) en claro, las series 2 y 8. Todos son de interfaz (3:1), salvo C-138, C-121, C-083, C-048, C-097, C-072, C-066 y C-099, que son texto (4.5:1) con margen de 0.2 a 0.34 salvo C-138.

## 5. Correcciones sugeridas (solo si U2 decide ganar holgura; ninguna es necesaria)

El sistema no tiene pares fallidos, así que no hay corrección obligatoria. Hallazgos de fragilidad para U2, sin tocar tokens aquí:

| Hallazgo | Par | Qué pasa | Corrección posible (decide U2) |
|---|---|---|---|
| F-1 | C-138 (claro) | `ui.text.secondary` (neutro 700, `#66676E`) sobre `ui.surface.hover` (neutro 200, `#E5E6E9`) = 4.51; un cambio de un paso de hex lo deja bajo 4.5 | En la fila resaltada del menú «/», pintar alias y atajo con `ui.text.default` (neutro 900); o conservar el diseño y dejar este CSV como prueba de regresión |
| F-2 | C-057, C-050 (oscuro) | Marcas de forma de Mercado y de Axel sobre `ui.surface.raised` en oscuro, 3.06 y 3.10 | Basta con conservarlas: son marcas con forma y etiqueta (K41); si se quiere holgura, subir un poco la L oscura de esas dos entradas |
| F-3 | C-110, C-112, C-108 (oscuro); C-107, C-113 (claro) | Series 5, 7, 3 (oscuro) y 2, 8 (claro) sobre tarjeta cerca de 3:1 | Conservar con el refuerzo de forma, patrón y etiqueta (K27); no apoyar nunca el significado solo en el color de esas series |

## 6. Pares exentos por decorativos

Divisores, pista de avance, retícula, esqueletos y borde de lo que flota (§10.5). Se calculan, se listan y no cuentan contra el umbral; ninguno porta información que no esté repetida en texto, forma o posición.

| Id | Componente · estado | Claro | Oscuro |
|---|---|---|---|
| C-087 | Tabla · divisor de filas | 1.19 | 1.89 |
| C-105 | Barra de avance · pista sobre tarjeta | 1.11 | 1.44 |
| C-116 | Gráfica · retícula | 1.11 | 1.44 |
| C-143 | Esqueleto de carga · sobre página | 1.07 | 1.30 |
| C-144 | Esqueleto de carga · sobre tarjeta | 1.11 | 1.44 |
| C-158 | Lo que flota · borde | 1.19 | 1.44 |

## 7. Cobertura de criterios

- **U2-04**: 316 combinaciones calculadas por script (código de salida 0), 0 bajo umbral; el recálculo independiente del revisor se hace con la columna `hex_frente`/`hex_fondo` del CSV.
- **U2-07**: botón primario claro, blanco sobre marca 700 = 5.77 (C-001) y marca 700 sobre fondo = 5.53 (C-002); oscuro, texto marca 950 sobre marca 300 = 10.47 (C-001) y relleno sobre fondo = 10.37 (C-002). Blanco sobre marca 600 no aparece en ninguna fila de texto (el relleno usa marca 700).
- **U5-04** (en el DOM): este CSV es la tabla de referencia contra la que U5 puede comparar lo medido.

## 8. Segundo cálculo independiente (10 razones)

`python3 -I docs/design/tools/contraste_componentes.py --independiente` recalcula 10 razones con otra implementación escrita aparte dentro del script: coma flotante, umbral de linealización 0.03928 de WCAG 2.0, canales en lista y producto punto. Resultado: las 10 coinciden con el CSV (el CSV trunca, así que se exige `csv ≤ real < csv + 0.01`); la mayor diferencia es +0.0095. Las diez razones elegidas por regla: las 4 de menor margen, 2 decorativas y 4 repartidas por el archivo.

| Id | Tema | Frente / fondo | CSV | Independiente |
|---|---|---|---|---|
| C-138 | claro | `#66676E` / `#E5E6E9` | 4.51 | 4.5110 |
| C-057 | oscuro | `#C2665B` / `#36363A` | 3.06 | 3.0622 |
| C-050 | oscuro | `#4F83C8` / `#36363A` | 3.10 | 3.1009 |
| C-110 | oscuro | `#B86797` / `#36363A` | 3.10 | 3.1078 |
| C-087 | claro | `#E5E6E9` / `#FAFAFB` | 1.19 | 1.1964 |
| C-158 | oscuro | `#4D4E53` / `#36363A` | 1.44 | 1.4495 |
| C-021 | claro | `#FFFFFF` / `#7E3236` | 8.76 | 8.7623 |
| C-061 | claro | `#B870AC` / `#FAFAFB` | 3.38 | 3.3827 |
| C-101 | claro | `#8B5D00` / `#F2F2F4` | 5.12 | 5.1269 |
| C-141 | claro | `#207B3E` / `#FAFAFB` | 5.07 | 5.0785 |

Ejemplo a mano (C-138 claro): canales de `#66676E` linealizados = 0.13287, 0.13563, 0.15593, luminancia 0.13651; canales de `#E5E6E9` = 0.78354, 0.79130, 0.81485, luminancia 0.79135; razón = (0.79135 + 0.05) / (0.13651 + 0.05) = 0.84135 / 0.18651 = 4.511.

## 9. Columnas del CSV

`id` (C-001 a C-158 de §10.5), `componente`, `estado`, `tema` (claro/oscuro), `par` (token semántico frente sobre fondo), `token_frente`, `token_fondo` (primitivos de C2), `hex_frente`, `hex_fondo`, `razon` (truncada a 2 decimales), `tipo`, `umbral` (4.5, 3.0 o `exento`), `cumple` (`true`, `false` o `exento`), `cruce_contraste_csv` (`igual`, `nuevo` o `difiere (...)`).
