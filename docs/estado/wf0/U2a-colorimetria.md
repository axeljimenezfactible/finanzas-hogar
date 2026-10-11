# U2a-colorimetria · Colorimetría del candidato C (Índigo 280°)

Flujo WF-0 «Cimientos», fase U2 parte a. Agente U2a-colorimetria. 2026-10-10.

## Nota de corrección · ronda 1 (2026-10-11)
Este reporte describe la entrega original de U2a y se conserva como registro. El agente `correccion-U2a-r1` corrigió el hallazgo H-03 de `docs/estado/wf0/REVISION-r1.md`; la evidencia completa está en `docs/estado/wf0/correccion-U2a-r1.md`. Lo que cambia frente a lo escrito abajo:
- Serie 6 de gráficas: «mostaza» 100° (`#978500` claro / `#DDCC6B` oscuro), dentro del matiz que decisiones §I reserva para acción crítica, pasa a «azul cielo» 225° (`#0092B9` / `#76D9FF`), el primer candidato del barrido de `colorimetria.md` §8.8.
- Con ese matiz, la búsqueda de orden de `colorimetria.md` §8.2 pone óxido en la serie 3 y magenta en la 5. Orden vigente: índigo, verde azulado, óxido, naranja, magenta, azul cielo, violeta, verde (antes: índigo, verde azulado, magenta, naranja, óxido, mostaza, violeta, verde).
- Cifras de gráficas que reemplazan a las de «Hecho» y «Verificación»: ΔL consecutivo mínimo 0.162 en claro (antes 0.170) y 0.198 en oscuro (igual); ΔE de vecinas con protanopia o deuteranopia 17.3 y 21.7 (antes 17.1 y 21.4) y en visión normal 17.8 y 22.1 (igual); variante culori 0.108 y 0.178 (igual); contraste mínimo 3.14 y 3.10 (truncado). Las tres primeras series pasan ahora la prueba de todos los pares de la skill dataviz (ΔE normal 22.7 y 21.3; con el orden anterior 13.9 y 12.8, bajo el piso 15).
- El script verifica además la zona reservada 80°–120° en todos los tokens, la variante culori y la prueba de todos los pares de las tres primeras series. Siguen iguales: 190 tokens (cambian 6: `color.chart.3.*`, `color.chart.5.*` y `color.chart.6.*`), 202 pares de contraste (cambian 12 filas, 0 bajo umbral) y el resto de la paleta.
- Para la puerta 2 se suma a lo pendiente la paleta de gráficas corregida (serie 6 en 225° y el intercambio de las series 3 y 5).

## Hecho
- Leí completos el protocolo y decisiones.md (secciones G e I aplicadas); del plan, §4 completo (4.1–4.4) y la fila U2 de §3.5; al final, por existir ya, los criterios `docs/criterios/U2.md` y `C2.md` y el bloqueo `criterios-deltal` (solo líneas de color, con grep).
- Escribí `docs/design/tools/colorimetria.py` (Python 3, solo biblioteca estándar, salida determinista). Implementa OKLCH → OKLab → sRGB lineal → sRGB (matrices de Ottosson, IEC 61966-2-1), recorte al gamut por bisección sobre el croma (60 iteraciones), cuantización ROUND_HALF_UP, luminancia relativa y contraste WCAG 2.x, simulación de protanopia, deuteranopia y tritanopia con las matrices de severidad 1.0 de Machado, Oliveira y Fernandes (2009), aplicadas en RGB lineal, y escala de grises por luminancia (L_gris = Y^(1/3)). Al ejecutarse regenera los 5 artefactos y luego los verifica leyéndolos del disco (código 1 si algo falla). Opciones: `--sin-png` y `--solo-verificar`.
- Método validado: deduje la tabla explícita del perfil de croma (0.10, 0.25, 0.40, 0.65, 0.85, 1.00, 1.00, 0.90, 0.75, 0.60, 0.45 del croma máximo) que reproduce **22/22** valores de partida del plan §4.3 (marca y neutro); documenté los intervalos que también los reproducen.
- Matiz base 280°, opuesto 100° reservado para acción crítica, armonía complementaria dividida (vecinos 70°/130°, ocupados por advertencia 75° e ingreso 135°) y regla 60-30-10 aplicada a superficies, acento y énfasis.
- 8 roles + acción crítica con escala completa de 11 pasos (croma máx. 0.14). Texto/icono en paso 700 en claro (el 600 pedido da 3.26–3.70:1 y no llega a 4.5; ajuste documentado en §6.2) y paso 400 en oscuro; fondo de insignia en paso 100/900.
- Ingreso vs éxito: propuse **135°** (la opción que prevé decisiones §I), con refuerzo obligatorio de signo «+» y flecha (ingreso) frente a palomita y texto (éxito); distancias calculadas en §6.4.
- Paleta categórica con variante clara y oscura por color y mismo matiz en ambos temas: personas Axel 256°, Ximena 350°, Karen 50°, Conjunto 280° (croma bajo a propósito); tiendas Walmart 262°, Bodega Aurrera 50°, Sam's Club 175°, Mercado 28°. Resueltas las 2 colisiones de decisiones §I (Karen y Bodega frente a advertencia) y 3 más que salieron al aplicar la regla (Axel frente a información, Mercado frente a gasto/error, Walmart en el límite exacto de 20°). ΔL ≥ 0.10 entre todos los pares de cada grupo (en OKLab y en gris) y separación ≥ 20° o ΔL ≥ 0.15 frente a cada rol (y frente a acción crítica).
- Paleta de 8 series en orden fijo (índigo, verde azulado, magenta, naranja, óxido, mostaza, violeta, verde) con dos clases de luminosidad alternadas; ΔL consecutivo mínimo 0.170 (claro) y 0.198 (oscuro) en gris y en las tres simulaciones.
- 202 pares de contraste en `contraste.csv` (texto normal, secundario, deshabilitado, placeholder, enlaces, botones primario/secundario/destructivo/acción crítica en reposo y hover, bordes de campo, foco, elemento activo, insignias y cifras de los 8 roles y de acción crítica, series, personas y tiendas sobre fondo y tarjeta, en claro y oscuro): 0 bajo umbral.
- 190 tokens DTCG en `C2-colores.json` (base, brand, neutral, role.<rol>.{light, dark, surface-light, surface-dark, 50–950}, critical, person, store, chart), cada uno con `$type`, `$value` hex y `$extensions.oklch` (l, c con 6 decimales, h).
- Muestrario `paleta.html` (sin CDN, fuentes ni JavaScript) y `paleta.png` 1400×3000 (chromium, primer intento), revisado a ojo.
- Bloqueo `U2-gasto-error` (gasto 25° y error 20°, ambos decididos, son indistinguibles: ΔE 1.2/1.0).

## Archivos
- /home/deepsec/finanzas-hogar/docs/design/tools/colorimetria.py (nuevo)
- /home/deepsec/finanzas-hogar/docs/design/colorimetria.md (generado)
- /home/deepsec/finanzas-hogar/docs/design/contraste.csv (generado)
- /home/deepsec/finanzas-hogar/docs/contratos/C2-colores.json (generado)
- /home/deepsec/finanzas-hogar/docs/design/paleta.html (generado)
- /home/deepsec/finanzas-hogar/docs/design/paleta.png (generado)
- /home/deepsec/finanzas-hogar/docs/estado/bloqueos/U2-gasto-error.md (nuevo)
- /home/deepsec/finanzas-hogar/docs/estado/wf0/U2a-colorimetria.md (este reporte)

## Verificación
- `python3 -I docs/design/tools/colorimetria.py` → escribe los 5 artefactos; «paleta.png generada (intento 1)»; contraste.csv 202 filas, 0 con cumple=false; C2-colores.json 190 tokens = Anexo A de colorimetria.md 190 filas (coinciden); plan §4.3 reproducido 22/22; gráficas claro ΔL consecutivo mín. gris 0.179 · protan 0.174 · deutan 0.170 · tritan 0.174, ΔE vecinas CVD 17.1 y normal 17.8, contraste mín. 3.15, croma mín. 0.112; gráficas oscuro 0.223 · 0.224 · 0.214 · 0.198, ΔE 21.4 / 22.1, contraste 3.11, croma 0.103; personas claro ΔL 0.102 (OKLab) / 0.105 (gris), oscuro 0.105 / 0.106; tiendas claro 0.105 / 0.112, oscuro 0.101 / 0.103; 0 colisiones con roles; paleta.html con 0 referencias externas; «VERIFICACIÓN OK» (código 0).
- `python3 -I docs/design/tools/colorimetria.py --solo-verificar` → VERIFICACIÓN OK. `python3 -I -m py_compile` → compila.
- Determinismo: sha256 de md, csv, json y html idénticos entre dos ejecuciones.
- Contraverificación independiente con culori 4.0.2 (node, en el scratchpad): 190/190 hex reproducidos desde `$extensions.oklch`; `wcagContrast` difiere del CSV < 0.01 (solo truncado), 0 pares bajo umbral; con su simulación (matriz sobre sRGB codificado) ΔL consecutivo mín. 0.108 en claro y 0.189 en oscuro (también ≥ 0.10).
- Validador de la skill dataviz (`validate_palette.py`, Machado 2009 en RGB lineal) sobre la tarjeta de cada tema: gráficas claro, todo PASS (CVD 17.1, normal 17.8, contraste ≥ 3:1); gráficas oscuro, PASS salvo la banda de L (0.84 > 0.67, documentado en colorimetria.md §8.7); personas y tiendas, CVD, piso normal y contraste en PASS; banda de L fuera por construcción (cuatro niveles separados ≥ 0.10) y piso de croma fuera solo en Conjunto (bajo a propósito).
- Style Dictionary 5.6.1 (instalado en el scratchpad) compila C2-colores.json a CSS y JS: 190 tokens, 190 variables `--color-*`, sin advertencias.
- Matrices de Machado: coinciden con la Tabla 1 de la página de los autores (consultada con Obscura) y con culori 4.0.2; la página del paquete R colorspace confirma la aplicación en RGB lineal desde la versión 2.1-0.
- Umbral de ΔL: 0.10 entre series consecutivas (criterio U2 del plan §3.5 y U2-09; más estricto que el 0.08 del encargo); los 28 pares se reportan como informativos (§8.4–8.5), igual que pide el bloqueo `criterios-deltal`.

## Supuestos
- «Fondo 50 / tarjeta 100» (claro) y «950 / 900» (oscuro) son pasos del neutro tintado (el 60 % de superficies); el plan medía los roles contra marca 50. No cambia valores: la diferencia de contraste entre ambas superficies es ≤ 0.01.
- «Texto 950» del botón primario en oscuro se tomó como brand 950 (el 950 de la misma escala del relleno); con neutral 950 también cumple (valor en colorimetria.md §9). No cambia tokens, solo la regla de uso.
- «ΔL ≥ 0.10 dentro de cada grupo» se aplicó a todos los pares del grupo, medido en L OKLab y en L del gris (la lectura más estricta). La separación de las categóricas se exigió también frente a acción crítica, no solo frente a los 8 roles.
- Los nombres de las hojas de token son propios (en inglés, protocolo §8): roles income, expense, savings, debt, warning, error, success, info; personas axel, ximena, karen, shared; tiendas walmart, bodega-aurrera, sams-club, market. Las personas coinciden con `person.code` del diccionario de S1; `store` no tiene código estable (ver solicitudes).

## Bloqueos
- /home/deepsec/finanzas-hogar/docs/estado/bloqueos/U2-gasto-error.md: gasto 25° y error 20° (decididos) son indistinguibles por color (ΔE 1.2 en claro, 1.0 en oscuro). Entregado con los valores decididos y regla de uso por contexto; 4 opciones con consecuencias para Axel. Nada marcado PROVISIONAL.
- Relacionado y ya existente: `criterios-deltal.md` (lectura de «ΔL entre series»); esta entrega cumple la lectura mínima (a) con umbral 0.10 y reporta los 28 pares.

## Pendientes
- Nada pendiente dentro de mi alcance (no cerré por tope).
- Puerta 2: Axel decide el bloqueo U2-gasto-error y aprueba ingreso 135° y la resolución de colisiones categóricas (colorimetria.md §6.4 y §7.5).
- WF-1 (U2b): integrar C2-colores.json sin cambiar valores en C2-tokens.json, tokens.css y el tema de Tailwind, con las reglas de uso de colorimetria.md §9.

## Solicitudes fuera de alcance
- `docs/criterios/U2.md` (U2-02 y U2-08): U2-02 lista ingreso 150°; decisiones §I permite que U2 proponga 135° y así se entregó. Conviene que el criterio acepte el valor propuesto y que U2-08 acepte el refuerzo declarado para los pares gasto–error e ingreso–éxito (tabla en colorimetria.md §7.7). Motivo: evitar un rechazo por el valor por defecto.
- `docs/contratos/diccionario-datos.md` (S1): la tabla `store` no tiene un código estable; un campo `code` (a-z, 0-9, guion bajo) permitiría mapear cada tienda a `color.store.<código>` sin depender del nombre visible; la regla «Conjunto 50/50» de `cost_share_rule` corresponde a `color.person.shared`. Motivo: que el color siga a la entidad.
- U2b (WF-1): si simula daltonismo con culori, convertir a `lrgb` antes de aplicar la matriz (culori 4.0.2 la aplica sobre sRGB codificado); mantener `.light/.dark` por tema también para series, personas y tiendas (un valor único no cumple 3:1 en ambos temas con escalonamiento: ventana L_gris [0.596, 0.641]).
- DESIGN.md (U2, WF-1): no copiar del plan §4.3 «bordes de campos neutro 600 ≥ 4.2:1» (medido 3.71:1, cumple 3:1) ni los roles a L 0.56 (3.92–4.47:1 sobre superficies claras, bajo 4.5:1).

## Paquetes instalados
- Ninguno del sistema ni dentro del repositorio.
- Descargas solo en el scratchpad de la sesión, fuera del repositorio, para verificar: culori 4.0.2 (`npm pack`) y style-dictionary 5.6.1 (`npm install`, 106 paquetes, `--ignore-scripts`).
- Red: registro de npm y dos páginas de documentación leídas con Obscura (tabla de matrices de Machado et al. y documentación de `simulate_cvd` del paquete R colorspace). Sin credenciales ni Telegram.

## Consumo aproximado
- Unas 95 acciones de herramienta.
- Leído: unos 60 KB (protocolo, decisiones, plan §4, guía y validador de la skill dataviz, fragmentos de culori, Machado y colorspace, líneas de criterios y diccionario).
- Escrito: script 102 KB; colorimetria.md 63 KB; contraste.csv 24 KB; C2-colores.json 53 KB; paleta.html 67 KB; paleta.png 864 KB; bloqueo 2.4 KB; este reporte unos 9 KB.
