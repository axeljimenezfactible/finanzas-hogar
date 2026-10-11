# correccion-U2a-r1 · Corrección de U2a (colorimetría), ronda 1

Flujo WF-0 «Cimientos», ronda 1 de correcciones. Agente `correccion-U2a-r1`. 2026-10-11.
Hallazgo atendido: **H-03** (media, criterio U2-02) de `docs/estado/wf0/REVISION-r1.md`. La serie 6 de gráficas («mostaza», `#978500` claro / `#DDCC6B` oscuro) usaba el matiz 100°, que decisiones §I y U2-02 reservan solo para «acción crítica», y el documento no declaraba la excepción.

## Hecho
- Lecturas completas: protocolo, decisiones, `docs/criterios/U2.md` y el reporte de U2a. De `REVISION-r1.md` leí el veredicto, la tabla U2, la cobertura I2 y H-03. Del plan leí la fila U2 de §3.5 y §4.1–§4.4 (con `sed -n`). De `colorimetria.md` y `colorimetria.py` leí las secciones que tocan la corrección. También leí la skill dataviz (SKILL y `color-formula.md`).
- Respaldé los 7 archivos de U2a en el scratchpad. Confirmé la línea base: el script original regenera md, csv, json y html idénticos byte a byte.
- **Elegí la opción A de la sugerencia del revisor: mover la serie 6.** No la conservé con una excepción. La opción B habría pedido a Axel matizar una decisión cerrada.
  - Nuevo valor: «azul cielo» **225°**, `#0092B9` en claro (C 0.117) y `#76D9FF` en oscuro (C 0.107).
  - Cómo se eligió: la nueva función `barrido_serie` del script recorre los múltiplos de 5° fuera de la zona reservada. Son 65 candidatos (otros 7 caen dentro de la zona). Para cada uno repite todos los controles de §8 en los dos temas.
  - Pasan 35. De los 30 descartados, 28 caen por croma < 0.10 y 2 porque su matiz medido cae dentro de la zona.
  - Entre los que pasan gana el que más se aleja en matiz de las otras siete series (control 1 de la skill dataviz: ocho familias de matiz). 225° queda a 45.3° de la serie más cercana. 230° quedaría a 49.9°, pero su croma en oscuro es 0.099, bajo el piso de 0.10. La tabla está en `colorimetria.md` §8.8.
- **Repetí la búsqueda de orden documentada en §8.2 con el matiz nuevo.** Ahora la hace el script (`ranking_orden`): recorre los 12 órdenes que alternan clases, con índigo en la serie 1 y violeta y verde al final.
  - El orden de la ronda 0, ya con 225°, baja al puesto 4: mínimo ΔL consecutivo 0.154. Además, sus tres primeras series no pasan la prueba de todos los pares de la skill dataviz: ΔE normal 12.8, bajo el piso de 15.
  - Ese defecto ya existía en la paleta original (13.9 en claro y 12.8 en oscuro, entre índigo y magenta, las dos de clase B), aunque §8.5 reserva las tres primeras series para dispersión, mapas y pequeños múltiplos.
  - **Orden vigente (puesto 1): índigo, verde azulado, óxido, naranja, magenta, azul cielo, violeta, verde.** Es decir, óxido pasa a la serie 3 y magenta a la 5. Mínimo ΔL consecutivo 0.162; tres primeras con ΔE normal 22.7 (claro) y 21.3 (oscuro), y ΔE protan/deutan 19.1 y 20.1.
  - Empata en los dos criterios de la ronda 0 con otro orden (azul cielo en la serie 2). Ese orden pierde en el tercer criterio, que agregué solo para desempatar: el ΔE entre las tres primeras series (12.8 frente a 21.3). Es el requisito que §8.5 ya imponía, no un criterio nuevo de gusto.
- **Declaré la zona reservada en §3: 80°–120°.** El ancho es la separación mínima de 20° que fija decisiones §I. No puede ser mayor: advertencia (75°, también decidida) queda a 25° del opuesto (24.3° medidos en el hex). §3 también lista los colores de identidad más cercanos al opuesto.
- Agregué tres verificaciones a `verificar()`. Las dos últimas comprueban afirmaciones que el documento ya hacía pero que el script no revisaba:
  - zona reservada en cada token fuera de `color.critical` con croma ≥ 0.02, con el matiz de diseño y con el medido;
  - ΔL consecutivo ≥ 0.10 también con la variante de culori (§8.6);
  - prueba de todos los pares de las tres primeras series (§8.5).
  - Además, la consola ahora muestra los contrastes truncados, como dice §2.5 (antes los redondeaba).
- Partes del documento que reescribí (todo lo genera el script; no hay ediciones a mano):
  - encabezado y §1;
  - viñeta de la zona en §3;
  - §8.2, con texto calculado y tabla de los 4 mejores de 12 órdenes;
  - §8.5, con las cifras de las tres primeras series;
  - §8.8 nueva, sobre el matiz reservado y la elección de la serie 6;
  - §12: hallazgos 9 y 10, y la paleta corregida en las propuestas para la puerta 2.
- Regeneré los cinco artefactos: md, csv, json, html y png.
- Agregué al principio de `U2a-colorimetria.md` una nota de corrección que lista qué cifras quedan reemplazadas. El resto del reporte original no cambió.

Tokens que cambian (de 190; todas las rutas siguen iguales):

| Token | Ronda 0 | Ronda 1 |
|---|---|---|
| `color.chart.3.light` / `.dark` | `#812E63` / `#B86797` (magenta) | `#883500` / `#C16D45` (óxido) |
| `color.chart.5.light` / `.dark` | `#883500` / `#C16D45` (óxido) | `#812E63` / `#B86797` (magenta) |
| `color.chart.6.light` / `.dark` | `#978500` / `#DDCC6B` (mostaza 100°) | `#0092B9` / `#76D9FF` (azul cielo 225°) |

## Archivos
- /home/deepsec/finanzas-hogar/docs/design/tools/colorimetria.py (modificado: serie 6, orden, zona reservada, `barrido_serie`, `ranking_orden`, verificaciones nuevas, textos de §1, §3, §8.2, §8.5, §8.8 y §12)
- /home/deepsec/finanzas-hogar/docs/design/colorimetria.md (regenerado)
- /home/deepsec/finanzas-hogar/docs/design/contraste.csv (regenerado; cambian 12 filas, las de las series 3, 5 y 6)
- /home/deepsec/finanzas-hogar/docs/contratos/C2-colores.json (regenerado; cambian 6 tokens)
- /home/deepsec/finanzas-hogar/docs/design/paleta.html (regenerado)
- /home/deepsec/finanzas-hogar/docs/design/paleta.png (regenerada con chromium, 1400×3000)
- /home/deepsec/finanzas-hogar/docs/estado/wf0/U2a-colorimetria.md (nota de corrección al inicio)
- /home/deepsec/finanzas-hogar/docs/estado/wf0/correccion-U2a-r1.md (este reporte)

## Verificación
Líneas de evidencia por criterio (convención de `INDICE.md` §1):
- U2-02 · `python3 -I docs/design/tools/colorimetria.py` → «zona reservada 80°–120° (acción crítica): 149 tokens cromáticos fuera de color.critical, 0 dentro; el más cercano color.role.warning.400 a 24.3°». La serie 6 mide 225.3° en los dos temas (125° del opuesto). Marca y neutro: 22/22 del plan.
- U2-02 (prueba negativa) · el script corregido con `--solo-verificar` sobre los archivos originales respaldados → «2 dentro … el más cercano color.chart.6.dark a 0.0°», «VERIFICACIÓN FALLIDA: zona reservada del eje opuesto usada fuera de color.critical: ['color.chart.6.light', 'color.chart.6.dark']», código 1. La verificación nueva detecta el defecto original.
- U2-04 · `contraste.csv`: 202 filas, 0 con `cumple=false`. Recálculo independiente de las 202 (implementación propia, sin importar el script): 0 difieren más de 0.01, 0 bajo umbral, y los hex del CSV coinciden con el JSON. Contraste de la serie 6: claro 3.46 (fondo) y 3.23 (tarjeta); oscuro 9.84 y 7.52.
- U2-09 · Claro: ΔL consecutivo mínimo gris 0.179, protanopia 0.174, deuteranopia 0.162, tritanopia 0.177; variante culori 0.108. ΔE de vecinas: protan/deutan 17.3, normal 17.8. Contraste mínimo 3.14; croma mínimo 0.112.
- U2-09 · Oscuro: ΔL 0.223, 0.224, 0.205 y 0.198; variante culori 0.178. ΔE de vecinas: 21.7 y 22.1. Contraste 3.10; croma 0.103.
- U2-09 · Tres primeras series, todos los pares: ΔE normal 22.7 (claro) y 21.3 (oscuro); protan/deutan 19.1 y 20.1. Pares no consecutivos con ΔL gris < 0.10: 12 de 28 en cada tema (informativo; bloqueo `criterios-deltal`).
- U2-03, U2-05, U2-06, U2-07 y U2-08 · sin cambios: la corrección no toca escalas, roles, botones ni categóricas. El script vuelve a dar 22/22 del plan, 0 colisiones de categóricas con roles y la tabla §7.7 igual.
- U2-14 · Documento en español de México; los tokens conservan nombres en inglés. Supuestos y bloqueos en este reporte.

Comandos y salidas (desde `/home/deepsec/finanzas-hogar`):
1. `python3 -I docs/design/tools/colorimetria.py` → escribe los 4 archivos de texto; «paleta.png generada (intento 1)»; 190 tokens = 190 filas del Anexo A; «VERIFICACIÓN OK», código 0.
2. `python3 -I docs/design/tools/colorimetria.py --solo-verificar` → «VERIFICACIÓN OK».
3. Reproducibilidad, en un árbol espejo del scratchpad:
   - el script original + el parche (`patch_r1.py` con 11 fragmentos) da el mismo `colorimetria.py` del repositorio, byte a byte;
   - regenerar desde cero da md, csv, json y html idénticos a los del repositorio;
   - dos ejecuciones seguidas dan los mismos sha256.
4. Contraverificación independiente con `verif_indep.py` (scratchpad, solo biblioteca estándar; OKLab de Ottosson, WCAG y Machado 2009 reescritos aparte):
   - sobre el repositorio → «RESULTADO: OK»;
   - sobre los archivos originales → detecta la serie 6 en la zona (0.0°) y las tres primeras series bajo el piso (12.8 / 8.4).
5. Validador de la skill dataviz, `validate_palette.py` sobre la tarjeta de cada tema (`#F2F2F4` / `#36363A`):
   - claro, 8 series adyacentes: todo PASS (CVD 17.3, normal 17.8);
   - claro, tres primeras con `--pairs all`: todo PASS (CVD 19.1, normal 22.7); con el orden de la ronda 0, FAIL del piso normal (13.9);
   - oscuro: PASS en croma, CVD 21.7, normal 22.1 y contraste; tres primeras PASS (CVD 20.1, normal 21.3). Sigue el FAIL de la banda de L (0.84 > 0.67), igual que en la ronda 0: el documento lo trata en §8.7.
6. Style Dictionary 5.6.1 (instalado en el scratchpad desde la caché local de npm) compila `C2-colores.json` a CSS y JS: 190 variables `--color-*`, 190 exportaciones, código 0, 0 advertencias. En CSS sale, por ejemplo, `--color-chart-6-light: #0092b9`.
7. `py_compile` → OK. El bytecode quedó en el scratchpad: no hay `__pycache__` en el repositorio.
8. Revisión visual de `paleta.png` (recorte de la zona de gráficas):
   - la leyenda muestra el orden nuevo en los dos temas;
   - el oliva-amarillo ya solo aparece en los botones «Confirmar» (acción crítica);
   - conteo de todos los píxeles de la imagen: 0 de `#978500` y 0 de `#DDCC6B`; 3,945 de `#0092B9` y 3,945 de `#76D9FF`.
9. No aplican a U2a, porque sus rutas no tienen SQL, Go, OpenAPI ni matriz: base temporal en `datos/tmp`, `go build` / `go vet`, validador OpenAPI y consistencia de la matriz. No creé `datos/tmp`.

sha256 finales:
- `beadf8708b445a6071e404de08d4731d7d974c77ae2c7c3a3014a4cf85f37534` colorimetria.md
- `9f995320c6fed1f622b796fc78b902b1b1fc497af499fbfea99a384318924bd4` contraste.csv
- `1446c402b3226d8e13c3166459315ba97ca105fc5d7169cfa35aca92d6c1a350` C2-colores.json
- `60fa92bbccfd202a5fab21068282f79ad54dd33a7c021f078da97e8a7dc0edd2` paleta.html
- `7746da788ff19cbbc6ea19bd5a8240dfbe341f27be059fed4493633471d9cb68` paleta.png
- `c70e643bb4300282248f423dcf20aee7a05cdb69e4bcf1f07c78adacd1e76bd3` colorimetria.py

## Supuestos
- **Ancho de la zona reservada: 20° a cada lado de 100° (80°–120°).** Decisiones §I reserva «el opuesto 100°» sin dar un ancho. Tomé la separación mínima entre matices distintos que la misma decisión fija para roles y categóricas; no puede ser mayor, porque advertencia (75°, decidida) quedaría dentro. Ninguna cifra depende de este ancho: el matiz elegido está a 125° del opuesto.
- El matiz de la serie 6 y el orden de las series son trabajo de diseño de U2a, con criterios escritos y calculados en §8.2 y §8.8. Igual que el resto de la colorimetría, quedan sujetos a la puerta 2 y no son supuestos sobre datos faltantes.

## Bloqueos
- Ninguno nuevo. `docs/estado/bloqueos/U2-gasto-error.md` sigue igual: no lo toca esta corrección.

## Pendientes
- Nada pendiente dentro de mi alcance (no cerré por tope).
- Puerta 2: Axel aprueba la paleta de gráficas corregida (serie 6 en 225°, óxido en la serie 3 y magenta en la 5), junto con ingreso 135° y las colisiones categóricas que ya estaban pendientes.
- WF-1 (U2b): integrar `C2-colores.json` sin cambiar valores. Ojo: `color.chart.3`, `.5` y `.6` ya no tienen los valores de la ronda 0.

## Solicitudes fuera de alcance
- `docs/criterios/U2.md` (dueño de criterios). Motivo: que el revisor lo mida igual que el script.
  - U2-02: precisar que la reserva del opuesto 100° vale para todo token (roles, categóricas y series) y fijar un ancho medible: 20°, con el matiz de diseño y el medido en el hex.
  - U2-09: agregar la prueba de todos los pares para las tres primeras series (ΔE normal ≥ 15 y protan/deutan ≥ 8). Es lo que exige la skill dataviz, adoptada por plan §4.1 punto 8, para dispersión, mapas y pequeños múltiplos, y §8.5 de la colorimetría limita esas gráficas a las tres primeras series.
- Integrador (H-05), opcional. Para reproducir U2a basta `docs/design/tools/colorimetria.py`. Si quiere guardar en el repositorio la contraverificación independiente, puede copiarla a `docs/design/tools/`: está en `/tmp/claude-1000/-home-deepsec/ad9b2909-7e27-4a0a-9e87-9e6a5395adce/scratchpad/u2a-r1/tools/verif_indep.py` (uso: `python3 -I verif_indep.py /home/deepsec/finanzas-hogar`). En la misma carpeta están el parche (`patch_r1.py` + `snip/`) y `dataviz_val.sh`.

## Paquetes instalados
- Ninguno en el sistema ni en el repositorio.
- En el scratchpad: style-dictionary 5.6.1 con `npm install --offline --ignore-scripts` (106 paquetes). Salió de la caché local de npm, sin red.
- No usé internet. Sin credenciales ni Telegram.

## Consumo aproximado
- Unas 65 acciones de herramienta.
- Leído: unos 190 KB (protocolo, decisiones, revisión, criterios U2, plan §3.5 y §4, el script casi completo, secciones de colorimetria.md, reporte U2a y skill dataviz).
- Escrito: script 121.8 KB (+22 KB); colorimetria.md 70 KB; contraste.csv 24 KB; C2-colores.json 53 KB; paleta.html 67 KB; paleta.png 864 KB; nota en el reporte U2a unos 2 KB; este reporte unos 11 KB.
