# REVISION-U2-r4 · Sistema de diseño (DESIGN.md + C2) · WF-1 · 11-oct-2026

Revisor independiente, ronda 4 (extra; el plan preveía 3). Evalúo solo contra `docs/criterios/U2.md` 1.2 y `docs/criterios/C2.md` 1.2, decisiones §G, §I y §K (K26, K27, K41, K42) y plan §4 (§4.1 y §4.4). No corregí nada. Cada verificación la repetí con comandos y scripts propios en `datos/tmp/rev-u2r4/` (copia aislada de `docs/contratos` y `docs/design` en `datos/tmp/rev-u2r4/repo`, borrada al terminar). Los artefactos revisados tienen el mismo sha256 que en r3 salvo `contraste-componentes.md`: DESIGN.md `6eb3700c…`, C2-tokens.json `570ff553…`, tokens.css `2b9429bc…`, tailwind.theme.css `f66eb866…`, contraste-componentes.csv `df8cf761…`, C2-colores.json `48daf4d6…`.

## Veredicto

**APROBADO: 0 altos, 0 medios, 11 bajos.**

- **M-r3-1 (media) queda CERRADO.** Lo comprobé por mi cuenta. `contraste-componentes.md` línea 3 dice ahora «Fuente de los colores: `docs/contratos/C2-tokens.json` (única; `C2-colores.json`, §7.1 y las columnas de §10.5 como cruce secundario)». La línea 23 dice que el script «resuelve cada lado y su hex solo desde `C2-tokens.json`». `C2-colores` solo aparece en esas dos líneas, y las dos veces como cruce secundario. Concuerda con la cabecera del script (líneas 8–16), con DESIGN §10.5 (línea 1199) y con el README de tokens. Mi propia resolución de las 316 combinaciones desde C2-tokens coincide con el CSV.
- No encontré ningún hallazgo alto ni medio nuevo. Siguen abiertos los bajos B-1 a B-10 de r3; B-9 se amplía a `U2d-contraste.md`. Hay un bajo nuevo, B-11: `colorimetria.md` dice «tipo Notion/Linear», fuera del alcance de U2.

## 1. Tabla de criterios (evidencia propia)

| Criterio | Estado | Evidencia del revisor |
|---|---|---|
| U2-01 (alta) | Cumple | Hay 11 encabezados `##` en orden: los 10 de plan §4.1 (lo leí) y Supuestos. Leí por partes §1–§6, §7.0, §7.1 (cabecera), §7.12–§7.16, §8, §9, §10.1, §10.3–§10.6 y Supuestos; §7.2–§7.11 y §10.5 los revisé por script. Ninguna sección está vacía. Las preguntas que le quedarían al implementador son menores y ya están en los bajos B-1, B-3, B-4 y B-6. |
| U2-02 (alta) | Cumple | Matiz OKLCH medido con mi conversión, paso 700 / paso 400. Marca: 280.2 / 279.9. Ingreso: 135.1 / 134.7 (la propuesta 135° de §I, con 150° y el motivo en §2.4). Gasto: 44.7 / 45.0 (K26). Error: 19.5 / 20.1. Éxito: 150.0 / 150.0. Advertencia: 75.4 / 75.7. Información: 239.8 / 239.6. Ahorro: 200.0 / 199.8. Deuda: 310.2 / 309.9. Acción crítica: 100.1 / 100.1. Karen 333.2° y Bodega Aurrera 356.1° (K41). Marca 50/950: `#F9FAFF` / `#1E203A`. Neutro 50/950: `#FAFAFB` / `#222325`, croma máximo 0.0121. |
| U2-03 (alta) | Cumple (ver B-10) | Desviación máxima de L por paso: marca 0.0010, neutro 0.0009 y roles hasta 0.0016 (ingreso). Croma máximo: marca 0.1103 y neutro 0.0121. Roles en el paso 400 (oscuro): como máximo 0.1200. En claro el máximo es 0.1408 (ingreso) y 0.1406 (deuda), sobre el 0.14 declarado por el redondeo a hex de 8 bits (B-10). |
| U2-04 (alta) | Cumple | Regeneré el CSV en la copia aislada: código 0, `cmp` IDÉNTICO y `--check` con código 0. Salen 158 pares y 316 combinaciones (304 con umbral y 12 decorativas), con 0 bajo umbral. Recalculé con mi propia fórmula las 316: 0 diferencias de hex (el hex de cada fila es el de su token en C2-tokens), 0 diferencias de razón truncada y 0 bajo umbral. Muestra listada de 30 pares (C-001 5.7705, C-022 12.5168, C-071 5.3744, C-099 4.8441, C-113 3.1684, C-138 4.5110, C-057o 3.0622 y otros). Mínimos: texto 4.511 / 5.035, interfaz 3.145 / 3.062 y texto grande 10.760. También recalculé las 202 filas de `contraste.csv`: 0 diferencias y 0 bajo umbral (mínimos 4.738 y 3.057). |
| U2-05 (alta) | Cumple | Valores recalculados: 900/50 = 11.533 (§2.8 documenta la diferencia de 0.13 con el plan). 700/50 = 5.397. 600/100 = 3.466. Blanco/marca 600 = 3.985. 100/950 = 14.068. 300/900 = 7.975. Marca 300/950 = 10.375. 200/50 = 1.196 y 800/950 = 1.895, los dos decorativos. Neutro 600/50 = 3.715, documentado como 3.71. Neutro 400/900 = 6.090. |
| U2-06 (alta) | Cumple | Ninguna fila de texto del CSV de componentes queda bajo 4.5 (mínimo 4.511). En `contraste.csv` el texto más bajo es 4.738. El signo, la flecha y el icono son obligatorios (§2.4, y §3.3 con U+2212). Un negativo sin rol lleva advertencia, icono y la palabra «rebasado». |
| U2-07 (alta) | Cumple | Por mi cuenta: blanco/marca 700 = 5.770 y marca 700/fondo = 5.532; en oscuro, marca 950/marca 300 = 10.471 y marca 300/fondo = 10.375. `ui.action.primary.bg` = `#5D609F` / `#C9CFFF`. No hay ninguna fila con fondo `#7579C2` y 0 de blanco sobre marca 600. |
| U2-08 (alta) | Cumple | Mi barrido cubre los 8 roles, las 8 categóricas y la marca 700/400. Encuentra 6 pares en claro y 5 en oscuro, los mismos 11 renglones de §2.6, cada uno con su refuerzo y con Δ matiz y ΔL iguales a mi medición (p. ej. Walmart–marca 700: 18.0° / 0.000; Conjunto–marca: 0.9° / 0.094). La regla K41 está en §2.6, §2.10 regla 3 y §8.7. `color.store.walmart.dark` = `#7AA3EE` con L 0.7165 (K42). No hay ninguna descripción con PROVISIONAL, ni en C2-colores ni en C2-tokens. |
| U2-09 (alta) | Cumple | Series en OKLCH L. Claro: 0.440 / 0.615 / 0.439 / 0.616 / 0.439 / 0.614 / 0.439 / 0.615. Oscuro: 0.621 / 0.840 / … El ΔL mínimo entre consecutivas es 0.174 en claro y 0.218 en oscuro (≥ 0.10, K27). Contraste mínimo: 3.145 / 3.108 sobre tarjeta y 3.371 / 4.063 sobre página. Forma, trazo y patrón están en §8.6. Los hex simulados de protanopia, deuteranopia, tritanopia y gris están en colorimetría §8.4 (línea 527). |
| U2-10 (alta) | Cumple | El grep del criterio contra terceros no da ningún archivo, con o sin `node_modules`. Los 7 sha256 (2 woff2, 2 ttf, LICENSE.txt, sprite y LICENSE de Lucide) coinciden con FUENTES.md. Licencias: OFL 1.1 en Inter e ISC en Lucide. El sprite tiene 1,877 `<symbol>` sin `stroke-width` fijo. `tabular-nums` está en §3.3, y `--font-numeric-tabular` en tokens.css. |
| U2-11 (media) | Cumple | En C2: espaciado 0/4/8/12/16/24/32/48, radios 6/10/14, duraciones 150/200/250 ms con 2 curvas, quiebres 360/412/768/1024/1440, `size.touch.min` = 44px, control grueso 44 (alias) y fino 36. «44 px» aparece 23 veces. No hay ningún mínimo táctil menor de 44: los 36 px son de puntero fino y los 24 × 24 se citan solo como piso de WCAG. `prefers-reduced-motion` está en §6.3 y en tokens.css. |
| U2-12 (media) | Cumple | La matriz §7.16 cubre los 12 grupos del criterio, y además barra superior, barra inferior, vista previa y chips. Por script, §7.2–§7.15 tienen tabla de estados o la regla equivalente: §7.12 y §7.14 dan los estados en prosa, con «los de sus controles». Usan tokens `ui.*`, que cambian con el tema, así que no hace falta repetir claro y oscuro. |
| U2-13 (media) | Cumple | §8.1–§8.8 están completos: barras, líneas, área apilada, dona ≤ 5 partes y etiquetas de dinero en texto. El tema sigue al sistema y tiene interruptor (§9.6). Notion es la única referencia (§1.1: «No se cita ni se imita ningún otro producto»). El grep `linear\|stripe\|material\|figma\|apple\|airbnb\|fluent\|carbon` en DESIGN da 0. La saturación apagada está en §2.1 y §2.2. Fuera de U2, ver B-11. |
| U2-14 (media) | Cumple | Español de México, tokens en inglés y 10 Supuestos. `U2-gasto-error`, `U2a-rosas-ximena-bodega` y `U2a-walmart-oscuro` están RESUELTOS. No hay ningún `U2-*` abierto. El bloqueo S1 se cita como PROVISIONAL en §7.13. |
| C2-01 (alta) | Cumple | En la copia aislada: `npm ci --offline` y Style Dictionary 5.6.1. `generate` reescribe C2-tokens.json (100073 B), IDÉNTICO al versionado. Compilé dos veces con código 0 y ninguna línea de advertencia o error (`warnings: 'error'`, `brokenReferences: 'throw'`). |
| C2-02 (alta) | Cumple | Hay 445 `$value` (444 hojas más `icon.stroke.$root`), 0 sin tipo efectivo y 0 referencias rotas. Por tipo: color 342, dimension 79, strokeStyle 8, fontWeight 4, shadow 3, number 3, duration 3, cubicBezier 2 y fontFamily 1. |
| C2-03 (alta) | Cumple | 101 grupos con `light` y `dark`, ninguno incompleto: 75 `ui.*`, 8 roles, 8 categóricas, 8 series, acción crítica y velo. Hay escalas de 11 pasos, tipografía de 6 tamaños, espaciado, radios, sombras 0/1/2 (`none` en oscuro), movimiento y 5 quiebres. |
| C2-04 (alta) | Cumple | Los 190 colores de C2-colores están en C2-tokens con el mismo hex, siguiendo los alias. DESIGN cita 162 hex (85 distintos), todos en C2-colores. Los 158 pares «token · hex» de DESIGN no tienen diferencias. De las 242 rutas citadas, las únicas 2 «inexistentes» son plantillas (`color.chart.N.light`, `viz.series.N.dash`). |
| C2-05 (alta) | Cumple | El script resuelve todo solo desde C2-tokens (cabecera y código). Mi resolución propia de las 194 filas «`ui.x` sobre `ui.y`» coincide, igual que el hex de las 316 filas. Hay 0 bajo umbral y existen tokens semánticos de par (`ui.text.on-primary`, etc.). |
| C2-06 (alta) | Cumple | `prefers-color-scheme: dark` 1, `data-theme` 2 y `prefers-reduced-motion` 1. `--ui-*`: 78 en claro, 78 en `@media` y 78 en `[data-theme="dark"]`, con 0 diferencias. Hay 190 `--color-*` estáticos. |
| C2-07 (alta) | Cumple | `tailwindcss` 4.3.3 da «Done» con código 0. En `c2-prueba.css` aparecen `#fafafb`, `#5d609f`, `#457723`, `border-radius: 10px`, `ui-shadow-1` y `ui-role-income`. |
| C2-08 (alta) | Cumple | Ninguna URL en tokens.css, C2-tokens.json ni tailwind.theme.css. Las dos `@font-face` apuntan a woff2 locales que existen, y la licencia está presente. |
| C2-09 (media) | Cumple | Dos compilaciones dan los mismos sha256 (`2b9429bc41836fc5…`, `f66eb866449f391a…`), y `cmp` sale IDÉNTICO contra lo versionado. |
| C2-10 (media) | Cumple | `verificar_c2.py` en la copia: 0 hex repetidos sin alias (36 alias de paso) y 0 semánticos con valor literal. Los 18 grupos de medidas iguales sin alias (p. ej. 12px = `font.size.caption` = `space.3`) son ideas de diseño distintas, como justifica U2c en su Supuesto 12; lo acepto. Los 6 `font.line-height.*` se citan como `font.line-height.*` en §10.2, con sus valores en §3.2. |
| C2-11 (media) | Cumple en lo evaluable | C2 v1.0 tiene historial en `$extensions.finanzas.cambios` y en el README. `docs/design/entrada/prompt.md` tiene 90 hex distintos, todos en C2. `entrada/C2-tokens.json` y `entrada/tokens.css` son IDÉNTICOS a los de `docs/contratos`, y el MANIFIESTO cita sus sha256. Las plantillas `.tex` de S7 aún no existen (fase S7): no es hallazgo de U2. |
| C2-12 (media) | Cumple | Hay Supuestos en DESIGN (10) y en los 4 reportes. No hay ningún bloqueo `C2-*`. §2.8 documenta los desvíos del plan §4.3 con los dos valores. |

Decisiones: K26 (gasto 44.7°/45.0°, error 19.5°/20.1°), K27 (ΔL consecutivo ≥ 0.10 con forma, trazo y patrón; 28 pares informativos), K41 (Karen y Bodega Aurrera sin mover, con refuerzo y la regla persona/tienda escrita) y K42 (`#7AA3EE`, L 0.7165, FIRME) quedan reflejadas. Nada en DESIGN contradice §G (Notion única referencia, saturación apagada, tema según el sistema, autoalojado) ni §I (pila y matiz 280°).

## 2. Hallazgos

No hay hallazgos altos ni medios.

### Bajos (no bloquean)
- **B-1 (de r1, abierta).** `DESIGN.md` línea 21, §1.2 P5, cita «§3.4», que no existe; debe ser §3.3.
- **B-2 (de r1, abierta).** `DESIGN.md` línea 169, §2.7: 0.162 / 0.198 son el mínimo en gris y en las tres simulaciones, y el texto no lo dice. En OKLCH L yo mido 0.174 / 0.218. Corrección: agregar «mínimo en gris y en las tres simulaciones».
- **B-3 (de r1, abierta).** `DESIGN.md` línea 63, §2.2: «±30° del opuesto». Los ángulos reales son −25° (75°) y +35° (135°).
- **B-4 (de r2, abierta).** `DESIGN.md` línea 116, §2.5: 5.00 es el contraste del relleno sobre la tarjeta, no el del texto. Corrección: «relleno 5.00:1 sobre tarjeta; texto blanco 5.59:1».
- **B-5 (de r2, abierta).** `docs/design/tools/tokens/verificar_c2.py` no compara con C2 los valores que no son de color de las tablas de DESIGN. Hoy coinciden: lo comprobé con mi script (espaciado, radios, duraciones, curvas, quiebres, tamaños de letra, iconos y táctil).
- **B-6 (de r2, abierta).** §7.0 (línea 572) y §7.2–§7.15 citan pares primitivos (`color.role.<rol>.light` / `.dark`), mientras §10.4 solo expone `--ui-*`. Corrección: citar `ui.role.<rol>`, `ui.person.<id>`, `ui.store.<id>` y `ui.chart.<n>`, o poner una tabla de equivalencias en §7.0.
- **B-7 (de r2, abierta).** `DESIGN.md` línea 3: «Versión 1.0 en construcción». Hay que quitar «en construcción» antes de la puerta 2.
- **B-8 (de r3, abierta y ampliada).** `docs/design/contraste-componentes.md` §4, «los 15 más justos»:
  - Faltan C-111 claro (+0.23) y C-056 claro (+0.28), y sobran C-072 y C-066.
  - El orden no sigue el margen.
  - Además, el párrafo «Lectura» nombra C-099, que no está en la tabla.
  - Mi lista ordenada por margen exacto: C-138c, C-057o, C-050o, C-110o, C-112o, C-107c, C-113c, C-108o, C-106o, C-111c, C-121c, C-056c, C-083c, C-048c, C-097c.
  - Corrección: regenerar la tabla y el párrafo.
- **B-9 (de r3, abierta y ampliada).** Reportes de agente desactualizados que el agente del paquete podría tomar como vigentes:
  - `docs/estado/wf1/U2c-tokens.md`: el Supuesto 5 (`size.bottombar.item.{min-width,min-height}`) y el sha256 de tokens.css que cita en C2-09.
  - `docs/estado/wf1/U2d-contraste.md`, línea 4: dice que el script toma «los hex solo de `docs/contratos/C2-colores.json`», cuando ahora resuelve solo desde C2-tokens (corrección H-1 de r1 y M-r3-1).
  - Corrección: una nota «superado por `correccion-U2-r1.md` / `correccion-U2-r3.md`» en los dos.
- **B-10 (de r3, abierta).** DESIGN §2.1 dice «máx. 0.14 roles». Medido en el hex llega a 0.1408 (ingreso) y 0.1406 (deuda). Corrección: agregar «medido en el hex: hasta 0.141 por redondeo de 8 bits».
- **B-11 (nueva, fuera del alcance de U2).** `docs/design/colorimetria.md`, línea 209 (§3): «Lectura sobria y moderna, tipo Notion/Linear».
  - Decisiones §G dicen «Notion es la única referencia visual». U2b-1 ya lo había pedido como solicitud fuera de alcance y sigue sin aplicarse.
  - No afecta a DESIGN.md: §1.1 lo prohíbe explícitamente y el paquete `entrada/` no lo contiene. Por eso es bajo.
  - Corrección: quitar «/Linear» en la próxima versión de la colorimetría. Lo decide el integrador.

## 3. Estado de los hallazgos anteriores

| Hallazgo | Estado | Comprobación propia |
|---|---|---|
| M-r3-1 (media, r3) C2-05 / U2-04 | **Cerrado** | Líneas 3 y 23 de `contraste-componentes.md` con el texto nuevo; `C2-colores` solo como cruce secundario. Coincide con el script (líneas 8–16), con DESIGN §10.5 y con el README. Las 316 filas se resuelven desde C2-tokens con 0 diferencias. |
| N-1 (media, r2) | Sigue cerrado | §10.6 da `indicator-width` 56 px e `indicator-height` 32 px; C2 coincide. |
| H-1, H-2, M-1, M-2, M-3 (r1) | Siguen cerrados | H-1: resolución desde C2-tokens. H-2: 0 rutas inexistentes, salvo plantillas. M-1: §4.6 y §9.3 con alto 56 y ancho mínimo 64. M-2: 36 alias de paso, 0 hex repetidos. M-3: una sola cadena (README, §10.1, generador). |
| B-1 a B-10 (bajas) | Abiertas | Las comprobé una por una (arriba). B-8 y B-9 se amplían. |

## 4. Para cerrar
Nada bloquea. Antes del paquete final de la Pista 3 y de la puerta 2 conviene aplicar las ediciones de texto B-1 a B-4, B-7, B-8 y B-10, que no cambian ninguna cifra ni ningún token, y la nota de B-9. B-11 va al dueño de la colorimetría.
