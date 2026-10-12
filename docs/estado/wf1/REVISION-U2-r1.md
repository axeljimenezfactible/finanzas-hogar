# REVISIÓN U2 · ronda 1 de 3 (revisor-U2-r1, WF-1, 11-oct-2026)

**Veredicto: NO APROBADO.** 2 hallazgos altos, 4 medios y 3 bajos. Ninguna cifra de color ni de contraste está mal: las 316 combinaciones de componentes y los 202 pares de la colorimetría dan 0 bajo umbral con mi cálculo propio. Lo que falla es de contrato y de coherencia. El script de contraste no lee `C2-tokens.json` (C2-05). DESIGN.md cita 4 nombres que en C2 no son tokens, y su §10 no está al día con C2 v1.0 (C2-04). Además hay una medida ambigua de la barra inferior, colores duplicados sin alias, una contradicción sobre cuál archivo es la fuente y no se dice cuál es la tabla oficial de contraste.

Base de evaluación: `docs/criterios/U2.md` 1.2 y `docs/criterios/C2.md` 1.2, decisiones §G, §I y §K (K26, K27, K41, K42) y plan §4. No corregí nada. Todas las verificaciones las volví a correr con scripts propios en `datos/tmp/rev-u2r1/` (borrados al terminar): conversión OKLab/OKLCH y contraste WCAG escritos por mí sin importar nada del repo, más los comandos de compilación.

## 1. Tabla de criterios (evidencia propia)

| Criterio | Estado | Evidencia del revisor |
|---|---|---|
| U2-01 (alta) | Cumple | 11 encabezados `## `: los 10 de plan §4.1 en orden y luego «Supuestos». Leí §1–§9, §10.1–§10.4, §10.6 y Supuestos, y ninguna sección está vacía. |
| U2-02 (alta) | Cumple | Matiz medido en hex (700 claro / 400 oscuro): marca 279.9/280.2, ingreso 135.1/134.7, gasto 44.7/45.0, ahorro 200.0/199.8, deuda 310.2/309.9, advertencia 75.4/75.7, error 19.5/20.1, éxito 150.0/150.0, info 239.8/239.6 y crítica 100.1. Pasos 50/950: marca `#F9FAFF`/`#1E203A` y neutro `#FAFAFB`/`#222325`, con croma de neutro máx. 0.0121. El ingreso en 135° se documenta con los dos valores y el motivo (§2.4). Karen está en 333° y Bodega Aurrera en 356°, con su refuerzo. |
| U2-03 (alta) | Cumple | Las 11 escalas (marca, neutro, crítica y 8 roles) tienen \|ΔL\| máx. 0.002 frente a la L fija por paso. Croma máx.: marca 0.1103, neutro 0.0121, roles ≤ 0.1408 y tokens oscuros (400) ≤ 0.120. El exceso de 0.0003–0.0008 viene del redondeo a hex, como explican colorimetría §5 y las revisiones previas. Las categóricas y las series llegan a 0.131 como máximo, así que la saturación sigue apagada. |
| U2-04 (alta) | Cumple en cifras | Recalculé **las 316 filas** de `contraste-componentes.csv` resolviendo cada par desde `C2-tokens.json`: 316 resueltas, 0 diferencias de token ni de hex, 0 diferencias de razón > 0.01 y **0 bajo umbral** (12 decorativas exentas). Mínimos: texto 4.51 (C-138 claro) e interfaz 3.06 (C-057 oscuro). Las 20 filas de muestra (una cada 16) coinciden. También recalculé los 202 pares de `contraste.csv`: 0 bajo umbral y diferencia 0. El script de U2d regenera el CSV idéntico (`cmp`) y `--check` sale con 0. Lo que falta del origen va en C2-05. |
| U2-05 (alta) | Cumple | 900/50 = 11.533 (el plan dice 11.4; la diferencia de 0.13 está documentada en §2.8), 700/50 = 5.397, 600/100 = 3.466, blanco/marca 600 = 3.985, 100/950 = 14.068, 300/900 = 7.975, marca 300/950 = 10.375, 200/50 = 1.196 y 800/950 = 1.895 (ambos declarados decorativos). Borde de campo: 600/50 = 3.715 (documentado) y 400/900 = 6.090. |
| U2-06 (alta) | Cumple | Rol claro (700) sobre fondo, tarjeta e insignia: el mínimo es 4.73 (éxito sobre tarjeta). Rol oscuro (400) sobre 950, 900 e insignia: el mínimo es 5.81. Con el paso 600 sobre fondo 50 todos dan de 3.49 a 3.98, y no se usa para texto. Los signos «+» y «−» (U+2212) y el icono son obligatorios (§2.4 y §3.3). |
| U2-07 (alta) | Cumple | C-001: 5.77 en claro y 10.47 en oscuro. C-002: 5.53 y 10.37. Hay 0 filas con `color.base.white` sobre `color.brand.600`. |
| U2-08 (alta) | Cumple | Mi barrido de 8 roles, 8 categóricas, marca (700/400) y crítica da 6 pares cercanos en claro y 5 en oscuro, **idénticos** a la tabla de §2.6, y cada uno trae su refuerzo. La regla K41 está escrita en §2.6, §2.10.3, §7.4 y §8.7. `color.store.walmart.dark` = `#7AA3EE`, L 0.717 medida (la colorimetría dice 0.716, que entra en el redondeo). Hay 0 descripciones con PROVISIONAL en Karen, Bodega Aurrera o Walmart. |
| U2-09 (alta) | Cumple | ΔL OKLCH entre series consecutivas: mín. 0.174 en claro y 0.218 en oscuro (≥ 0.10). Los 28 pares van solo como información (12 están bajo 0.10 por la alternancia A/B). Contraste mínimo sobre fondo y tarjeta: 3.145 y 3.108. La tabla de hex simulados está en colorimetría §8.4 (gris, protanopia, deuteranopia y tritanopia). |
| U2-10 (alta) | Cumple | Inter 4.1: los `woff2` tienen firma `wOF2` y los TTF traen `tnum`, con licencia OFL 1.1. Lucide 1.55.0: `sprite.svg` con 1,877 símbolos y 0 `stroke-width` fijos, licencia ISC más la MIT de Feather. El grep de terceros sobre docs/design y docs/contratos (md, css, json) da 0. |
| U2-11 (media) | Cumple (ver M-1) | Coinciden: espaciado 4…48, radios 6/10/14, tipografía 12…32, 150/200/250 ms con las curvas estándar y enfática, quiebres 360…1440, `size.touch.min` = 44px y «44 px» aparece 23 veces. `prefers-reduced-motion` está en DESIGN §6.3 y en tokens.css, con la excepción del indicador de carga. |
| U2-12 (media) | Cumple | Están los 12 grupos: §7.2, §7.3, §7.5.1–7.5.5, §7.6, §7.7, §7.8, §7.9, §7.11, §7.13, §7.14 y §7.15, cada uno con los 6 estados y sus tokens. |
| U2-13 (media) | Cumple | §8: barras, líneas, área apilada, dona ≤ 5 y dinero siempre en texto. §9.6: tema según el sistema con interruptor. Densidad cómoda (P6). El grep de `linear\|stripe\|material\|apple\|ios\|figma` en DESIGN.md da 0. Notion es la única referencia (§1.1). |
| U2-14 (media) | Cumple | Texto en es-MX, tokens en inglés y sección Supuestos. Bloqueos `U2-gasto-error`, `U2a-rosas-ximena-bodega` y `U2a-walmart-oscuro`: todos RESUELTOS. |
| C2-01 (alta) | Cumple | Con `rm -rf node_modules && npm ci` seguido de `npm run build`: código 0, `style-dictionary --version` = 5.6.1, `log.warnings: 'error'` y `brokenReferences: 'throw'`. |
| C2-02 (alta) | Cumple | 445 tokens: 0 sin `$value`, 0 sin tipo efectivo y 0 referencias rotas. Por tipo: color 342, dimension 79, fontWeight 4, shadow 3, number 3, duration 3, cubicBezier 2, fontFamily 1 y strokeStyle 8. |
| C2-03 (alta) | Cumple | Hay escalas de 11 pasos, 8 roles más la acción crítica, 8 categóricas, 8 series, tipografía, espaciado, radios, sombras 0/1/2, movimiento y quiebres. 78 `--ui-*` en claro, 78 en oscuro por media query y 78 con `data-theme`. |
| C2-04 (alta) | **No cumple (H-2)** | 190 de 190 colores de `C2-colores.json` están en C2-tokens con el mismo hex. De 162 hex en DESIGN (85 distintos), 0 quedan fuera de C2. 158 pares token·hex sin diferencias y las 41 filas de §7.1 sin diferencias. Pero 4 nombres citados en DESIGN no son tokens en C2, y §10 contradice a C2 v1.0. |
| C2-05 (alta) | **No cumple (H-1)** | Las cifras están bien (ver U2-04), pero el script no lee `C2-tokens.json`. |
| C2-06 (alta) | Cumple | Los conjuntos `--ui-*` son iguales en los 3 bloques (78/78/78). El oscuro va en `@media (prefers-color-scheme: dark) { :root:not([data-theme="light"]) }` y en `:root[data-theme="dark"]`. `prefers-reduced-motion` aparece 1 vez. |
| C2-07 (alta) | Cumple | `npm run test:tailwind` (4.3.3): código 0. En la salida están `#fafafb`, `#5d609f`, `#457723`, `var(--ui-role-income)`, `border-radius: 10px` y `14px`, `var(--ui-shadow-1)`, `cubic-bezier` y `48rem`. La salida se borró. |
| C2-08 (alta) | Cumple | 0 `http(s)://` en tokens.css, tailwind.theme.css y C2-tokens.json. Las 2 rutas de `@font-face` existen. |
| C2-09 (media) | Cumple | Tras `npm ci` limpio y `build`, `cmp` da tokens.css y tailwind.theme.css idénticos. `check:sync` reproduce los sha256 `4956…1c88` y `f66e…c88a`. `npm run generate` regenera C2-tokens.json idéntico (`cmp`). |
| C2-10 (media) | **No cumple (M-2)** | Hay 36 grupos de primitivos de color con el mismo hex y sin alias. |
| C2-11 (media) | Cumple en lo propio | El contrato está identificado como C2 v1.0 con su sección de cambios (README y `$extensions`). El consumo queda pendiente porque `docs/design/entrada/` está vacía y las `.tex` de S7 no existen todavía. |
| C2-12 (media) | Cumple | Hay sección Supuestos (U2c, 13 supuestos), no hay `docs/estado/bloqueos/C2-*` y los desvíos de plan §4.3 están documentados. |

## 2. Hallazgos

### H-1 (alta) · C2-05: el script de contraste no lee `C2-tokens.json`
- **Ruta:** `docs/design/tools/contraste_componentes.py` (líneas 9–10 y 35–36: `C2 = RAIZ / "docs/contratos/C2-colores.json"`, «única fuente de los hex»), y `docs/design/DESIGN.md` §10.5 (línea 1198), que afirma «U2d la calcula desde `C2-tokens.json` como única fuente (C2-05)».
- **Qué pasa:** C2-05 pide que el script de U2-04 lea `C2-tokens.json` como única fuente. El script toma la correspondencia `ui.*` → primitivo de DESIGN §7.1 y los hex de `C2-colores.json`. Si `C2-tokens.json` cambia (por ejemplo, un `ui.*` editado) y DESIGN no, la tabla no se entera. Lo que dice §10.5 es falso.
- **Impacto hoy:** ninguno en las cifras. Resolví las 316 combinaciones desde `C2-tokens.json` y dan 0 diferencias.
- **Corrección sugerida:** que el script resuelva cada par por las referencias de `C2-tokens.json`: `ui.<x>.light/.dark`; «rol X» → `ui.role.X`; «superficie de rol X» → `ui.role.X-surface`; «persona/tienda/serie» → `ui.person|store|chart.*`; «acción crítica» → `ui.critical`. DESIGN §10.5 quedaría solo como lista de pares, y §7.1 y `C2-colores.json` como cruce secundario. Después se vuelve a correr y se confirma el mismo CSV.

### H-2 (alta) · C2-04: DESIGN.md cita nombres que en C2 no son tokens, y §10 no refleja C2 v1.0
- **Ruta:** `docs/design/DESIGN.md` §4.3 (líneas 358, 361 y 362), §4.6 (línea 410), §7.10 (línea 886), §7.12 (líneas 922 y 924), §9.3, §10.1 (línea 1142), §10.2 (línea 1151), §10.4 (línea 1190) y §10.6 (línea 1378).
- **Qué pasa:** `layout.topbar.height`, `layout.preview.width`, `layout.dialog.width` y `size.bottombar.item` aparecen como tokens con valor, pero en `C2-tokens.json` son grupos sin `$value`. Los tokens reales son `.mobile/.desktop`, `.desktop/.wide`, `.short/.long` y `.min-width/.min-height` (Supuesto 5 de U2c), y en CSS no existe `--layout-dialog-width`. Además: §10.1 llama al tema `tailwind-theme.css`, pero el archivo es `tailwind.theme.css`. §10.2 dice «41 semánticos» `ui.*`, y C2 tiene 75 pares (41 más 34 de rol, crítica, persona, tienda y serie); también lista `shadow.1` y `shadow.2`, y C2 trae además `shadow.0`. §10.4 dice que el tema hace `@import "tailwindcss"`, cuando en realidad solo importa `tokens.css` y es la hoja del frontend la que importa tailwindcss (Supuesto 9 de U2c). Es la solicitud 1 de U2c, que nadie aplicó.
- **Corrección sugerida:** en DESIGN, citar los 4 tokens con sus nombres desdoblados y su valor, renombrar el archivo del tema, poner el conteo real de `ui.*` (75 pares, más `overlay.scrim`), agregar `shadow.0` y describir la importación real (la del README de tokens). Luego correr `npm run check` y confirmar que `generar_c2_tokens.py` sigue leyendo §7.1.

### M-1 (media) · Barra inferior: «56 × ≥ 64 px» se leyó al revés y choca con el alto de la barra
- **Ruta:** `docs/design/DESIGN.md` §4.6 (línea 410) y §9.3; `docs/contratos/C2-tokens.json` `size.bottombar.item.min-width` = 56px y `min-height` = 64px; `layout.bottombar.height` = 56px.
- **Qué pasa:** un destino con alto mínimo de 64 px no cabe en una barra de 56 px. Lo que dice la propia DESIGN («5 destinos en 360 px = 72 px cada uno») apunta a 56 de alto × ≥ 64 de ancho. Quien implemente tendría que adivinar.
- **Corrección sugerida:** que DESIGN lo escriba sin ambigüedad («alto 56 px = `layout.bottombar.height`, ancho mínimo 64 px», si es lo que se quiso) y que C2 se regenere con los nombres y valores que salgan de ahí.

### M-2 (media) · C2-10: 36 grupos de colores con el mismo valor y sin alias
- **Ruta:** `docs/contratos/C2-tokens.json` (grupos `color.role.<rol>.*` y `color.critical.*`), `docs/design/tools/tokens/generar_c2_tokens.py`, `docs/design/tools/tokens/verificar_c2.py` (solo revisa los tokens que no son de color) y DESIGN §10.2 (afirma que «dos tokens con el mismo valor siempre quedan relacionados»).
- **Qué pasa:** en los 8 roles y en la acción crítica, `light` = `700`, `dark` = `400`, `surface-light` = `100` y `surface-dark` = `900`, y todos van como hex literal. Por ejemplo, `color.role.income.light` = `color.role.income.700` = `#457723`. Es la misma idea de diseño (DESIGN §2.4: «texto paso 700 en claro y 400 en oscuro») y no tiene alias.
- **Corrección sugerida:** que el generador emita `light`, `dark` y `surface-*` como referencias (`{color.role.income.700}`…). El hex no cambia y `C2-colores.json` no se toca. Extender `verificar_c2.py` a los colores. La otra salida es que el integrador apruebe la excepción por escrito.

### M-3 (media) · No queda claro cuál es la fuente ni cuál es la tabla oficial de contraste
- **Ruta:** `docs/design/tools/tokens/README.md` («lo demás se cambia en `C2-tokens.json` y se vuelve a compilar»; `npm run check` empieza con `generate`), `docs/design/DESIGN.md` §10.1 («única fuente» `C2-tokens.json`, pero «nombres y valores de los tokens no de color» en DESIGN §3–§9 y en el generador), §1.2 P8, §2.8 y §2.10 regla 5 (todo par nuevo «entra en `contraste.csv`»), y §10.5 (no da la ruta de `contraste-componentes.csv` ni de `contraste_componentes.py`; habla de una «tabla final» como unión que no existe como archivo).
- **Qué pasa:** si alguien cambia un valor no de color en `C2-tokens.json` como dice el README, `npm run check` lo pisa al regenerar desde DESIGN y el generador. Y quien agregue un par de componente no sabe si va a `colorimetria.py`/`contraste.csv` o a §10.5/`contraste_componentes.py`.
- **Corrección sugerida:** declarar una sola cadena, por ejemplo «DESIGN §3–§10.6 + `C2-colores.json` → `generar_c2_tokens.py` → `C2-tokens.json` (no se edita a mano) → CSS». También nombrar en §10.5 y §2.10.5 las dos tablas oficiales: `contraste.csv`, para los pares de paleta de colorimetría, y `contraste-componentes.csv` con su script, para los pares de componentes. Es la solicitud de U2d.

### B-1 (baja) · Referencia rota
`docs/design/DESIGN.md` §1.2 P5 (línea 21) cita «cifras tabulares (§3.4)», pero esa sección no existe: es §3.3.

### B-2 (baja) · ΔL de la paleta de gráficas sin decir cómo se mide
`docs/design/DESIGN.md` §2.7: «ΔL mínimo … 0.162 en claro y 0.198 en oscuro» es el mínimo de gris más las 3 simulaciones (colorimetría §8.3). En OKLCH normal da 0.174 y 0.218, y en gris 0.179 y 0.223. Conviene decir «mínimo en gris y en las simulaciones».

### B-3 (baja) · Armonía dividida descrita a ±30°
`docs/design/DESIGN.md` §2.2: «acentos a ±30° del opuesto → advertencia 75° e ingreso 135°». En realidad están a −25° y +35° del opuesto 100°. Mejor escribir «a unos ±30°» o dar los grados reales.

## 3. Estado de lo anterior
Es la primera revisión de U2 en WF-1, así que no hay hallazgos previos de este artefacto. Solicitudes de los agentes que tocan este artefacto:
- U2c, solicitación 1 (sincronizar DESIGN §10 y §4 con C2): **no aplicada**, pasa a H-2.
- U2d (ruta del CSV en §10.5): **no aplicada**, pasa a M-3.
- U2b-2, solicitud 4 (comprobar `ui.*` citados contra C2): hecha; los 41 coinciden.
- U2b-1, solicitudes 2 y 3 (colorimetría «Notion/Linear», U1 §6.2): fuera del alcance de U2 y sin efecto en DESIGN.md, que no cita a Linear.

## 4. Para la ronda 2
Basta con cerrar H-1, H-2 y M-1 a M-3. Las cifras de color y contraste no necesitan cambio. Después de corregir, hay que volver a correr `npm run check`, `python3 -I docs/design/tools/contraste_componentes.py --check` y `python3 -I docs/design/fuentes/scripts/verificar.py .`.
