# REVISION-U2-r3 · Sistema de diseño (DESIGN.md + C2) · WF-1 · 11-oct-2026

Revisor independiente, ronda 3 de 3. Evalúo solo contra `docs/criterios/U2.md` 1.2 y `docs/criterios/C2.md`, decisiones §G, §I y §K (K26, K27, K41, K42) y plan §4. No corregí nada. Todas las verificaciones las repetí con comandos y scripts propios en `datos/tmp/rev-u2r3/` (copia aislada del repositorio en `datos/tmp/rev-u2r3/repo`, borrada al terminar).

## Veredicto

**NO APROBADO: 0 altos, 1 medio y 10 bajos.**

- **N-1 (media, r2) queda CERRADO.** Lo comprobé por mi cuenta. §7.10 (línea 893) dice ahora «píldora de 56 × 32 px (`size.bottombar.indicator-width` × `size.bottombar.indicator-height`)». En C2 esos tokens valen 56px y 32px. La frase sigue la convención ancho × alto del resto del documento (interruptor 36 × 20, línea 675; asa 32 × 4, línea 923). «32 × 56» ya no aparece ninguna vez.
- **Hay un medio nuevo, M-r3-1.** El resumen oficial de la tabla de componentes, `docs/design/contraste-componentes.md`, sigue diciendo que el script toma los hex «solo de `C2-colores.json`» y la resolución de §7.1. Desde la corrección H-1 de r1 eso ya no es así: el script resuelve todo desde `C2-tokens.json` (C2-05). Es una inconsistencia entre artefactos sobre la fuente de verdad de un criterio alto.
- La corrección de M-r3-1 son dos líneas de texto. **Ninguna cifra ni ningún token cambian.** Todo lo demás cumple con mi propia evidencia.

## 1. Tabla de criterios (evidencia propia)

| Criterio | Estado | Evidencia del revisor |
|---|---|---|
| U2-01 (alta) | Cumple | Hay 11 encabezados `##` en orden: los 10 de plan §4.1 y Supuestos. Leí DESIGN.md completo (§10.5 por script): ninguna sección está vacía. Para implementar no quedan preguntas abiertas de severidad media. Las ambigüedades menores son los bajos B-1, B-3, B-4 y B-6. |
| U2-02 (alta) | Cumple | Matiz OKLCH medido con mi conversión, paso 700 / paso 400. Ingreso: 135.1 / 134.7, la propuesta 135° de §I, documentada junto a 150° y su motivo en §2.4. Gasto: 44.7 / 45.0 (K26). Error: 19.5 / 20.1. Éxito: 150.0 / 150.0. Advertencia: 75.4 / 75.7. Información: 239.8 / 239.6. Ahorro: 200.0 / 199.8. Deuda: 310.2 / 309.9. Acción crítica: 100.1. Marca 700: 280.2. Karen 333° y Bodega Aurrera 356° (K41). Marca 50/950: `#F9FAFF` / `#1E203A`. Neutro 50/950: `#FAFAFB` / `#222325`, croma máximo 0.0121. |
| U2-03 (alta) | Cumple (ver B-10) | Desviación máxima de L por paso en las 11 escalas: 0.0016 (ingreso). En marca y neutro: 0.0010 y 0.0009. Croma máximo: marca 0.110 y neutro 0.012. Roles en el paso 400 (oscuro): como máximo 0.120. En los roles el máximo es 0.1408 (ingreso 500) y 0.1406 (deuda 600). El objetivo declarado es 0.14 y la diferencia viene del redondeo a hex de 8 bits. |
| U2-04 (alta) | Cumple | Regeneré el CSV en la copia aislada y salió idéntico (`cmp`). `--check` termina con código 0: 158 pares, 316 combinaciones (304 con umbral y 12 decorativas) y 0 bajo umbral. Recalculé **las 316** con mi propia resolución: 194 por el par semántico `ui.*` y 122 por la regla de §10.5 (rol, superficie de rol, persona, tienda, serie, acción crítica). Resultado: 0 diferencias de hex y 0 de razón (truncada a 2 decimales), y 0 bajo umbral. También recalculé las 202 filas de `contraste.csv`: 0 diferencias y 0 bajo umbral, con mínimos de 4.738 en texto y 3.057 en interfaz. Muestra de 12 (semilla 20261014), propio / CSV: C-140o 10.4268/10.42 · C-045o 7.9623/7.96 · C-039c 5.7845/5.78 · C-124c 10.7304/10.73 · C-152c 10.9320/10.93 · C-131c 5.3970/5.39 · C-085o 8.0805/8.08 · C-092o 10.4268/10.42 · C-082o 10.9022/10.90 · C-135o 10.7599/10.75 · C-091o 10.4714/10.47 · C-039o 7.6063/7.60. |
| U2-05 (alta) | Cumple | Valores recalculados: 900/50 = 11.533 (§2.8 documenta la diferencia de 0.13 con el plan). 700/50 = 5.397. 600/100 = 3.466. Blanco/marca 600 = 3.985. 100/950 = 14.068. 300/900 = 7.975. Marca 300/950 = 10.375 (10.471 sobre marca 950). 200/50 = 1.196 y 800/950 = 1.895, los dos decorativos. Neutro 600/50 = 3.715, documentado como 3.71. Neutro 400/900 = 6.090. |
| U2-06 (alta) | Cumple | El CSV tiene 53 filas de texto con color de rol; la mínima es 4.73. El signo, la flecha y el icono son obligatorios (§2.4 y §3.3, con U+2212). |
| U2-07 (alta) | Cumple | C-001 da 5.77 en claro y 10.47 en oscuro; C-002, 5.53 y 10.37. Hay 0 filas de blanco sobre `#7579C2`. |
| U2-08 (alta) | Cumple | Mi barrido cubre los 8 roles, las 8 categóricas y la marca 700/400. Encuentra 6 pares en claro y 5 en oscuro, que son exactamente los 11 renglones de §2.6, cada uno con su refuerzo. No hay ningún par cercano a la acción crítica. La regla K41 está escrita en §2.6, §2.10 regla 3, §7.4 y §8.7. `color.store.walmart.dark` = `#7AA3EE`, con L 0.7165 (K42). Hay 0 descripciones con PROVISIONAL en Karen, Bodega Aurrera y Walmart, tanto en C2-colores como en C2-tokens. ΔE de Ximena contra Bodega Aurrera en oscuro: 1.33, como dice §2.6. |
| U2-09 (alta) | Cumple | ΔL mínimo entre series consecutivas, en OKLCH L: 0.174 en claro y 0.218 en oscuro (≥ 0.10, K27). Los 28 pares son solo informativos; su mínimo es 0.000 porque las series alternan dos clases de L. Contraste mínimo de las series: 3.145 / 3.108 sobre tarjeta y 3.371 / 4.063 sobre página. Forma, trazo y patrón en §8.6. La simulación de daltonismo con los hex simulados está en colorimetría §8.4 (línea 527). |
| U2-10 (alta) | Cumple | El grep del criterio contra terceros (con `--include` entre comillas) no da ningún archivo. Los 7 sha256 (2 woff2, 2 ttf, LICENSE.txt, sprite y LICENSE de Lucide) aparecen en FUENTES.md. Licencias: OFL 1.1 en Inter; ISC y MIT de Feather en Lucide. El sprite tiene 1,877 `<symbol>`. `tabular-nums` está en §3.3. |
| U2-11 (media) | Cumple | En C2: espaciado 0/4/8/12/16/24/32/48, radios 6/10/14, duraciones 150/200/250 ms con 2 curvas, quiebres 360/412/768/1024/1440, `size.touch.min` = 44px, control grueso 44, fino 36. «44 px» aparece 23 veces. No hay ningún mínimo táctil menor de 44: los 24 × 24 citados son el piso de WCAG. `prefers-reduced-motion` está en DESIGN §6.3 y en tokens.css. |
| U2-12 (media) | Cumple | Los 12 grupos del criterio están, más chips, barras y diálogos (matriz §7.16), cada uno con sus 6 estados y sus tokens. N-1 está cerrado. |
| U2-13 (media) | Cumple | §8.1–§8.8 completos. Tema según el sistema con interruptor (§9.6). Notion es la única referencia (§1.1). `grep -i 'linear\|stripe\|material'` = 0. |
| U2-14 (media) | Cumple | Español de México, tokens en inglés y 10 Supuestos. `U2-gasto-error`, `U2a-rosas-ximena-bodega` y `U2a-walmart-oscuro` están RESUELTOS. No hay ningún `U2-*` abierto. |
| C2-01 (alta) | Cumple | En la copia aislada: `npm ci` con código 0 y Style Dictionary 5.6.1. Generé y compilé dos veces, con código 0 y 0 advertencias o errores en el registro (`warnings: 'error'`, `brokenReferences: 'throw'`). |
| C2-02 (alta) | Cumple | 445 tokens: 0 sin tipo efectivo y 0 referencias rotas. Por tipo: color 342, dimension 79, strokeStyle 8, fontWeight 4, shadow 3, number 3, duration 3, cubicBezier 2 y fontFamily 1. |
| C2-03 (alta) | Cumple | `--ui-*`: 78 en claro, 78 en `@media` oscuro y 78 en `data-theme`, con 0 diferencias. Están las escalas, los 8 roles con la acción crítica, las 16 categóricas, las 16 series, tipografía, espaciado, radios, sombras, movimiento y quiebres. |
| C2-04 (alta) | Cumple | Los 190 colores de C2-colores están en C2-tokens con el mismo hex, siguiendo los alias. Los 85 hex distintos de DESIGN están en C2-colores (0 fuera). Los 158 pares «token · hex» de DESIGN no tienen diferencias. De las 254 rutas citadas, 0 no existen. |
| C2-05 (alta) | Cumple en el script; ver M-r3-1 | La cabecera del script, líneas 8–16, y su código (`TOKENS`, línea 42) resuelven solo desde C2-tokens; C2-colores es un cruce secundario. Mi resolución propia de las 316 coincide. El resumen `.md` contradice esto. |
| C2-06 (alta) | Cumple | `prefers-color-scheme: dark` 1, `data-theme` 2, `prefers-reduced-motion` 1 y 0 URL. Los conjuntos de `--ui-*` son iguales en los tres bloques. |
| C2-07 (alta) | Cumple | `test:tailwind` (4.3.3) da «Done». En `c2-prueba.css` aparecen `#fafafb` y `#5d609f` (sin distinguir mayúsculas), además de `#457723`, `var(--ui-role-income)`, `border-radius: 10px` y `var(--ui-shadow-1)`. |
| C2-08 (alta) | Cumple | 0 URL. Las dos `@font-face` apuntan a woff2 locales que existen, y la licencia está presente. |
| C2-09 (media) | Cumple | Dos compilaciones dan el mismo sha256. `cmp` contra lo versionado sale IDÉNTICO en C2-tokens.json, tokens.css y tailwind.theme.css. `check:sync` da OK (`2b9429bc…`, `f66eb866…`). |
| C2-10 (media) | Cumple | `verify`: 0 hex repetidos sin alias (36 alias de paso) y 0 semánticos literales. De las 255 hojas que no son color, 0 quedan sin cita en DESIGN. |
| C2-11 (media) | Sin evaluar todavía | C2 v1.0 tiene historial (README). `docs/design/entrada/` está vacía y las plantillas `.tex` de S7 no existen: no es hallazgo de U2. |
| C2-12 (media) | Cumple | Hay Supuestos en DESIGN y en los reportes. 0 bloqueos `C2-*`. Los desvíos del plan §4.3 se documentan con los dos valores (§2.8). |

## 2. Hallazgos

### M-r3-1 (media) · C2-05 / U2-04 · El resumen de contraste de componentes describe una fuente que ya no se usa
- Ruta: `docs/design/contraste-componentes.md`, líneas 3 y 23.
- Defecto:
  - La línea 3 dice «Fuente de los colores: `docs/contratos/C2-colores.json`».
  - La línea 23 dice que el script toma «la resolución semántica de §7.1, y los hex **solo** de `C2-colores.json`».
  - Desde la corrección de H-1 (r1), `contraste_componentes.py` declara en las líneas 8–16 que «C2-tokens.json: ÚNICA fuente de la resolución y de los hex (C2-05)», y lo cumple. §7.1, las columnas de §10.5 y C2-colores quedan solo como cruces secundarios. Lo mismo dicen DESIGN §10.5 y el README de tokens.
  - El documento que resume la tabla oficial afirma lo contrario a lo que exige un criterio alto. Quien mantenga la tabla después lo tomaría como cierto.
- Corrección sugerida, sin cambiar ninguna cifra:
  - Línea 3: «Fuente de los colores: `docs/contratos/C2-tokens.json` (única; `C2-colores.json`, §7.1 y las columnas de §10.5 como cruce secundario)».
  - Línea 23: «El script lee la lista de pares de `DESIGN.md` §10.5 y resuelve cada lado y su hex solo desde `C2-tokens.json`; §7.1, las columnas de §10.5 y `C2-colores.json` son cruces secundarios que salen con código 1 si no coinciden».

### Bajos (no bloquean)
- **B-1 (de r1, abierta).** `DESIGN.md` línea 21, §1.2 P5, cita «§3.4», que no existe; debe ser §3.3.
- **B-2 (de r1, abierta).** `DESIGN.md` línea 169, §2.7: los valores 0.162 / 0.198 son el mínimo en gris y en las tres simulaciones (colorimetría), pero el texto no lo dice. Yo mido 0.174 / 0.218 en OKLCH L. Corrección: agregar «mínimo en gris y en las tres simulaciones».
- **B-3 (de r1, abierta).** `DESIGN.md` línea 63, §2.2, dice «±30° del opuesto». Los ángulos reales son −25° (75°) y +35° (135°).
- **B-4 (de r2, abierta).** `DESIGN.md` línea 116, §2.5: 5.00 es el contraste del relleno sobre la tarjeta (5.008), no el del texto. El blanco sobre `#766800` da 5.599. Corrección: «relleno 5.00:1 sobre tarjeta; texto blanco 5.59:1».
- **B-5 (de r2, abierta).** `docs/design/tools/tokens/verificar_c2.py` no compara los valores que no son de color de las tablas «Token | Valor» de DESIGN con C2. Hoy coinciden; lo comprobé con mi script para U2-11.
- **B-6 (de r2, abierta).** §7.0 y §7.2–§7.15 citan pares primitivos (`color.role.error.light` / `.dark`, `color.person.<nombre>.*`), aunque §10.4 solo expone `--ui-*`. La equivalencia se deduce de §7.0 (último punto) y de §10.3. Corrección: citar `ui.role.<rol>`, `ui.person.<id>`, `ui.store.<id>` y `ui.chart.<n>`, o dar una tabla de equivalencias en §7.0.
- **B-7 (de r2, abierta).** `DESIGN.md` línea 3 dice «Versión 1.0 en construcción». Hay que quitar «en construcción» antes del paquete de la Pista 3.
- **B-8 (nueva).** `docs/design/contraste-componentes.md` §4, «los 15 más justos»:
  - Faltan C-111 claro (serie 6 sobre tarjeta, 3.23, +0.23) y C-056 claro (marca de forma, 3.27, +0.27).
  - En su lugar aparecen C-072 (+0.31) y C-066 (+0.34).
  - El orden no sigue el margen: C-121 va antes que C-113 y C-106 va al final.
  - Mi lista ordenada: C-138c, C-057o, C-050o, C-110o, C-112o, C-107c, C-113c, C-108o, C-106o, C-111c, C-121c, C-056c, C-083c, C-048c, C-097c.
  - Corrección: regenerar la tabla ordenada por margen.
- **B-9 (nueva).** `docs/estado/wf1/U2c-tokens.md` está desactualizado en dos puntos:
  - El Supuesto 5 todavía dice `size.bottombar.item.{min-width,min-height}`, «leído como ancho mínimo 56 y alto mínimo 64». La corrección M-1 de r1 lo superó: ahora es `item.height` = 56 (alias) y `min-width` = 64.
  - El sha256 de tokens.css que cita en C2-09 (`495664…`) ya no es el vigente.
  - Corrección: una nota «superado por `correccion-U2-r1.md`» para que el agente del paquete de la Pista 3 no lo lea como vigente.
- **B-10 (nueva).** U2-03: el croma de ingreso 500 mide 0.1408 y el de deuda 600 mide 0.1406, sobre el máximo de 0.14 por el redondeo a hex. DESIGN §2.1 dice «máx. 0.14 roles» sin mencionarlo. Corrección: agregar «medido en el hex: hasta 0.141 por redondeo de 8 bits».

## 3. Estado de los hallazgos anteriores

| Hallazgo | Estado | Comprobación propia |
|---|---|---|
| N-1 (media, r2) U2-12 / C2-04 | **Cerrado** | Línea 893: «56 × 32 px (`size.bottombar.indicator-width` × `size.bottombar.indicator-height`)». En C2-tokens, `indicator-width` = 56px e `indicator-height` = 32px. `grep -c '32 × 56'` = 0. La píldora mide 32 de alto, menos que `layout.bottombar.height` (56), y 56 de ancho, menos que `item.min-width` (64). |
| H-1, H-2, M-1, M-2, M-3 (r1) | Siguen cerrados | H-1: el script resuelve desde C2-tokens (pero ver M-r3-1 en su resumen). H-2 y C2-04: 0 rutas inexistentes. M-1: `item.height` → 56px y `min-width` 64px. M-2: 36 alias, 0 hex repetidos. M-3: una sola cadena (README, §10.1 y generador). |
| B-1 a B-7 (bajas) | Abiertas | Las comprobé una por una (líneas arriba). Ninguna se tocó en la corrección r2, que solo atendía N-1. |

## 4. Para cerrar
Hay que corregir M-r3-1 (dos líneas de `docs/design/contraste-componentes.md`). Ninguna cifra ni ningún token cambian y no hace falta regenerar nada: el `.md` no lo genera el script (`grep -c 'contraste-componentes.md' contraste_componentes.py` = 0). De paso conviene corregir B-1 a B-4, B-7, B-8 y B-10, que son ediciones de texto, antes del paquete de la Pista 3.
