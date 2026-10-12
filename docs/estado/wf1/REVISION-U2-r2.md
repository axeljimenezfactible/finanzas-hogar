# REVISION-U2-r2 · Sistema de diseño (DESIGN.md + C2) · WF-1 · 11-oct-2026

Revisor independiente, ronda 2 de 3. Evalúo solo contra `docs/criterios/U2.md` 1.2 y `docs/criterios/C2.md`, decisiones §G, §I y §K (K26, K27, K41 y K42) y plan §4. No corregí nada. Todas las verificaciones se repitieron con mis propios comandos y scripts en `datos/tmp/rev-u2r2/`, en una copia aislada del repositorio (ya borrada).

## Veredicto

**NO APROBADO: 0 altos, 1 medio y 7 bajos.** Los 5 hallazgos altos y medios de la ronda 1 quedaron **cerrados**, y lo comprobé por mi cuenta. Queda abierto un hallazgo medio nuevo, N-1: la píldora de la barra inferior aparece en §7.10 como «32 × 56 px». El resto del documento escribe ancho × alto, y con esa convención la frase contradice los tokens `size.bottombar.indicator-width` = 56 px e `indicator-height` = 32 px. Se corrige en una línea y no cambia ninguna cifra ni ningún token.

## 1. Tabla de criterios (evidencia propia)

| Criterio | Estado | Evidencia del revisor |
|---|---|---|
| U2-01 (alta) | Cumple | 11 encabezados `##` en orden: las secciones 1 a 10 y Supuestos. Leí cada sección y ninguna está vacía. |
| U2-02 (alta) | Cumple | Matiz OKLCH medido con conversión propia, paso 700 / paso 400. Ingreso 135.1 / 134.7: es la propuesta 135° de §I, documentada en §2.4 con 150° y el motivo. Gasto 44.7 / 45.0 (K26: 40–45). Error 19.5 / 20.1. Éxito 150.0. Advertencia 75.4. Información 239.8. Ahorro 200.0. Deuda 310.2. Acción crítica 100.1. Karen 333.2 y Bodega Aurrera 356.1 (K41). Marca 50 / 950: `#F9FAFF` / `#1E203A`. Neutro: `#FAFAFB` / `#222325`, croma máximo 0.012. |
| U2-03 (alta) | Cumple | Desviación de L por paso: como máximo 0.001 en marca y neutro y 0.002 en roles. Croma máximo: marca 0.110. Roles: 0.141 en ingreso y deuda, 0.001 sobre el tope por redondeo a hex (colorimetría 1.2 ya aprobada). Paso 400 en oscuro: como máximo 0.120. |
| U2-04 (alta) | Cumple | Regeneré el CSV y salió idéntico (`cmp`). `--check` da código 0: 316 combinaciones, 304 con umbral, 12 decorativas y 0 bajo umbral. Con resolución y fórmula propias recalculé **las 316**: 0 diferencias de token, de hex o de razón (±0.01) y 0 bajo umbral. La holgura menor es C-138 en claro: 4.51 frente a 4.5. Muestra de 10 (semilla 20261011): C-126 oscuro 7.4233 / 7.42 · C-089 oscuro 7.9750 / 7.97 · C-127 oscuro 10.9022 / 10.90 · C-087 oscuro 1.8951 / 1.89 (decorativo) · C-072 claro 4.8139 / 4.81 · C-015 claro 5.3970 / 5.39 · C-011 claro 3.4660 / 3.46 · C-092 oscuro 10.4268 / 10.42 · C-123 oscuro 10.7599 / 10.75 · C-010 oscuro 7.9623 / 7.96. |
| U2-05 (alta) | Cumple | Recalculé: 900/50 = 11.53 (difiere 0.13 del plan y §2.8 lo documenta). 700/50 = 5.40. 600/100 = 3.47. Blanco/marca 600 = 3.98. 100/950 = 14.07. 300/900 = 7.97. Marca 300/950 = 10.38 (10.47 sobre marca 950). 200/50 = 1.20 y 800/950 = 1.90, ambos decorativos. Neutro 600/50 = 3.71, documentado. Todo coincide con §2.8: allí se trunca y yo redondeo. |
| U2-06 (alta) | Cumple | 58 filas de texto de rol en el CSV, mínimo 4.73. El signo, la flecha y el icono son obligatorios (§2.4, §3.3, U+2212). |
| U2-07 (alta) | Cumple | C-001: 5.77 en claro y 10.47 en oscuro. C-002: 5.53 y 10.37. Hay 0 filas de blanco sobre `#7579C2`. `color.brand.600` solo se usa en `--ui-border-focus`. |
| U2-08 (alta) | Cumple | Mi barrido (8 roles, 8 categóricas y la marca 700/400) encuentra 6 pares en claro y 5 en oscuro. Son exactamente los 11 renglones de §2.6, cada uno con su refuerzo. La regla K41 está en §2.6, §2.10 regla 3, §7.4 y §8.7. `color.store.walmart.dark` = `#7AA3EE`, con L 0.7165, que corresponde a 0.716 (K42). Ninguna descripción de Karen, Bodega Aurrera ni Walmart dice PROVISIONAL. |
| U2-09 (alta) | Cumple | ΔL mínimo entre series consecutivas: OKLCH L 0.174 en claro y 0.218 en oscuro; gris equivalente 0.179 y 0.223. Contando también las simulaciones salen los 0.162 y 0.198 de §2.7 (colorimetría línea 144), todos ≥ 0.10. Contraste mínimo contra `ui.surface.card`: 3.15 y 3.11. Forma, trazo y patrón en §8.6. Simulación en colorimetría §8.4. Ver B-2. |
| U2-10 (alta) | Cumple | El grep del criterio contra terceros da 0 archivos. Los 8 sha256 coinciden con FUENTES.md. Licencias OFL 1.1 (Inter) e ISC + MIT (Lucide). El sprite tiene 1,877 `<symbol>`. `tabular-nums` está en §3.3. |
| U2-11 (media) | Cumple | 4 pt, rejilla 4/8/12, radios 6/10/14, 2 sombras en claro y bordes en oscuro, 150–250 ms y quiebres 360…1440: coinciden con C2. «44 px» aparece 23 veces. `prefers-reduced-motion` está en DESIGN y en tokens.css. La barra inferior quedó en alto 56 y ancho ≥ 64 (§4.6, §9.3, C2). |
| U2-12 (media) | Cumple con N-1 | 12 de 12 grupos (matriz §7.16), cada uno con estados y tokens. N-1 es una contradicción de medida en §7.10. |
| U2-13 (media) | Cumple | §8.1–§8.8 y tema según el sistema con interruptor (§9.6). Notion es la única referencia (§1.1). `grep -i 'linear\|stripe\|material'` = 0. |
| U2-14 (media) | Cumple | Español de México, tokens en inglés y 10 Supuestos. `U2-gasto-error`, `U2a-rosas-ximena-bodega` y `U2a-walmart-oscuro` están RESUELTOS. No hay ningún `U2-*` abierto. |
| C2-01 (alta) | Cumple | En la copia aislada: `npm ci` con código 0, Style Dictionary 5.6.1, build con código 0 (`warnings: 'error'`, `brokenReferences: 'throw'`). |
| C2-02 (alta) | Cumple | 445 tokens (444 hojas más `icon.stroke.$root`), 0 sin tipo y 0 referencias rotas. Por tipo: color 342, dimension 79, strokeStyle 8, fontWeight 4, shadow 3, number 3, duration 3, cubicBezier 2 y fontFamily 1. |
| C2-03 (alta) | Cumple | 76 grupos `ui`/`overlay` y todos tienen claro y oscuro. Escalas, roles, acción crítica, 16 categóricas, 16 series, tipografía, espaciado, radios, sombras, movimiento y quiebres. |
| C2-04 (alta) | Cumple | Los 190 colores de C2-colores están en C2-tokens con el mismo hex (siguiendo alias). Los 162 hex de DESIGN están en C2-colores. De 242 tokens citados, 0 no existen (los 2 que mi script no encuentra son los patrones `color.chart.N.light` y `viz.series.N.dash`). 0 grupos citados como token. Los valores no de color citados junto a su token coinciden con C2: revisé a mano las 12 alertas de mi heurística y las 12 eran falsos positivos. |
| C2-05 (alta) | Cumple | El script resuelve solo con C2-tokens. Mutación propia: `ui.text.default.light` → neutral 800 y `color.brand.300` → `#C0C6FF`. Cambian 42 filas del CSV con los valores mutados y los cruces secundarios avisan con código 1. |
| C2-06 (alta) | Cumple | `--ui-*` 78/78/78 y 0 diferencias. `prefers-color-scheme` 1, `data-theme` 2, `prefers-reduced-motion` 1. |
| C2-07 (alta) | Cumple | Tailwind 4.3.3 «Done». Aparecen `#fafafb`, `#5d609f`, `#457723`, `var(--ui-role-income)`, `border-radius: 10px` y `var(--ui-shadow-1)`. |
| C2-08 (alta) | Cumple | 0 URL. Las 2 `@font-face` apuntan a archivos que existen y la licencia está presente. |
| C2-09 (media) | Cumple | `generate` dos veces da un C2-tokens idéntico al versionado. El build da tokens.css `2b9429bc…` y tailwind.theme.css `f66eb866…`, idénticos a los versionados, y dos compilaciones salen idénticas. Hay una sola cadena, escrita en el README, en §10.1 y en el generador. |
| C2-10 (media) | Cumple | 0 hex literales repetidos sin alias (36 alias de paso) y 0 semánticos literales. Hay 18 medidas no de color iguales: §10.2 las declara como información por su regla explícita. Los 6 `font.line-height.*` se citan como grupo en §10.2 y como valores en §3.2. |
| C2-11 (media) | Sin evaluar todavía | C2 v1.0 con historial (README y `$extensions`). `docs/design/entrada/prompt.md` y las plantillas `.tex` de S7 aún no existen: no es hallazgo de U2. |
| C2-12 (media) | Cumple | Hay Supuestos en DESIGN y en los reportes U2b-1, U2b-2, U2c y U2d. 0 bloqueos `C2-*`. Los desvíos del plan §4.3 se documentan con los dos valores (§2.8). |

## 2. Hallazgos

### N-1 (media) · U2-12 / C2-04 · Medida de la píldora de la barra inferior invertida
- Ruta: `docs/design/DESIGN.md` línea 893 (§7.10, estado Activo).
- Defecto: dice «píldora de 32 × 56 px». DESIGN escribe ancho × alto en todo lo demás (asa «32 × 4 px», línea 923, = `size.sheet.handle-width`/`-height`; pista «36 × 20 px», línea 675, = `size.switch.width`/`.height`). Con esa convención la píldora mediría 32 de ancho y 56 de alto. Pero C2 y §10.6 dicen `size.bottombar.indicator-width` = 56 px e `indicator-height` = 32 px, y una píldora de 56 px de alto llenaría toda la barra de 56 px. Es el mismo tipo de ambigüedad que M-1 de la ronda 1, en el mismo componente, y llega a la Pista 3 a través de DESIGN.pdf.
- Corrección sugerida: «píldora de 56 × 32 px (`size.bottombar.indicator-width` × `size.bottombar.indicator-height`)». Ningún token cambia.

### B-1 (baja, de r1, abierta) · §1.2 P5 cita «§3.4», que no existe; la referencia correcta es §3.3.
### B-2 (baja, de r1, abierta) · §2.7 da ΔL 0.162 / 0.198 sin decir que es el mínimo en gris y en las tres simulaciones (colorimetría línea 144). Mido 0.174 / 0.218 en OKLCH L y 0.179 / 0.223 en gris. Corrección: agregar «mínimo en gris y en las tres simulaciones».
### B-3 (baja, de r1, abierta) · §2.2 dice «acentos a ±30° del opuesto», pero los ángulos reales son −25° (75°) y +35° (135°). Corrección: escribir los ángulos reales.
### B-4 (baja, nueva) · §2.5, línea 116: «texto `color.base.white` (5.00:1 sobre tarjeta)» se lee como el contraste del texto. En realidad 5.00 es relleno sobre tarjeta (`#766800`/`#F2F2F4` = 5.008). Blanco sobre `#766800` da 5.60. Corrección: «relleno 5.00:1 sobre tarjeta; texto blanco 5.59:1».
### B-5 (baja, nueva) · `docs/design/tools/tokens/generar_c2_tokens.py` líneas 10–11: los valores no de color se transcriben a mano y `verificar_c2.py` solo compara nombres. Hoy coinciden porque lo verifiqué yo, pero un cambio en DESIGN que no se copie al generador pasaría `npm run check`. Corrección: que `verificar_c2.py` compare las tablas «Token | Valor» de §3–§6, §9.4 y §10.6 con C2.
### B-6 (baja, nueva) · §7.2–§7.15 citan los pares primitivos (`color.role.error.light` / `.dark`, `color.person.<nombre>.light`…), aunque §7.1 dice que los componentes no combinan primitivos y §10.4 solo expone `--ui-*`. La correspondencia (`ui.role.error` → `--ui-role-error` → `text-role-error`) se deduce de §7.0 y §10.3. Corrección: citar `ui.role.<rol>`, `ui.person.<id>`, `ui.store.<id>` y `ui.chart.<n>`, o dar la tabla de equivalencias en §7.0.
### B-7 (baja, nueva) · El encabezado dice «Versión 1.0 en construcción». Corrección: quitar «en construcción» antes del paquete de la Pista 3.

## 3. Estado de los hallazgos anteriores (r1)

| Hallazgo | Estado | Comprobación propia |
|---|---|---|
| H-1 (alta) C2-05 | **Cerrado** | `contraste_componentes.py` resuelve solo con C2-tokens (`semantico_de`, `primitivo` y `hex_de`). La prueba de mutación propia sigue a C2-tokens (42 filas cambian). §7.1, §10.5 y C2-colores quedan como cruces secundarios con código 1. |
| H-2 (alta) C2-04 | **Cerrado** | Los 4 grupos se citan por sus hojas, con valores iguales a C2. `icon.stroke` es un token válido (`$root` = 1.5px, `--icon-stroke`). §10.1 dice `tailwind.theme.css`. §10.2: 75 semánticos y `shadow.0`. §10.4 dice que el tema importa solo `./tokens.css`, lo comprobé en la línea 9 del tema. |
| M-1 (media) U2-11 | **Cerrado** | `size.bottombar.item.height` = `{layout.bottombar.height}` (56px) y `min-width` = 64px. §4.6, §7.10, §9.3 y el Supuesto 10 coinciden. |
| M-2 (media) C2-10 | **Cerrado** | 36 alias a 700/400/100/900, 0 hex repetidos sin alias y 0 semánticos literales. La regla de §10.2 es verdadera. |
| M-3 (media) C2-09 / U2-04 | **Cerrado** | El README, §10.1 y el generador dan la misma cadena única. §1.2 P8, §2.8, §2.10 regla 5 y §10.5 nombran las dos tablas y el script. |
| B-1, B-2, B-3 (bajas) | Abiertas | Quedaron fuera de la corrección r1. Siguen como bajas. |

## 4. Para la ronda 3
Corregir N-1 (una línea de DESIGN §7.10). De paso conviene corregir B-1 a B-4 y B-7, que son ediciones de texto en DESIGN.md. Después hay que volver a correr `python3 -I docs/design/fuentes/scripts/verificar.py .`, `npm run check` (en `docs/design/tools/tokens/`) y `python3 -I docs/design/tools/contraste_componentes.py --check`. Ninguna cifra ni ningún token cambian.

Fuera de mi alcance, y no es hallazgo: el Supuesto 9 de DESIGN trata como diseño, y no como parámetros de K44, el aviso de 5 s, el «Guardado» de 2 s, la carga desde 400 ms y la lentitud a los 10 s. Lo marca como pendiente de que Axel decida.
