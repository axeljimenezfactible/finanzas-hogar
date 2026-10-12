# correccion-U2-r3 · Sistema de diseño (DESIGN.md + C2) · WF-1 · 11-oct-2026

Corrige el único hallazgo alto o medio de `docs/estado/wf1/REVISION-U2-r3.md` (M-r3-1, media). No se tocan los bajos B-1 a B-10, que no estaban en la tarea, ni ninguna decisión de §A–§K. Ninguna cifra ni ningún token cambian.

## Hecho
- **M-r3-1 (media) · C2-05 / U2-04**: en `docs/design/contraste-componentes.md` cambié las dos líneas que describían una fuente de colores que ya no se usa. Usé el texto que sugirió el revisor, palabra por palabra.
  - Línea 3. Antes decía «Fuente de los colores: `docs/contratos/C2-colores.json`». Ahora dice «Fuente de los colores: `docs/contratos/C2-tokens.json` (única; `C2-colores.json`, §7.1 y las columnas de §10.5 como cruce secundario)».
  - Línea 23. Antes decía que el script toma «la resolución semántica de §7.1, y los hex **solo** de `C2-colores.json`». Ahora dice «El script lee la lista de pares de `DESIGN.md` §10.5 y resuelve cada lado y su hex solo desde `C2-tokens.json`; §7.1, las columnas de §10.5 y `C2-colores.json` son cruces secundarios que salen con código 1 si no coinciden». La frase final «No importa `colorimetria.py`.» se conserva.
- Antes de escribir, comprobé que el texto nuevo describe lo que hace el código:
  - `contraste_componentes.py`, líneas 8–16: «C2-tokens.json: UNICA fuente de la resolucion y de los hex (C2-05)».
  - Línea 42: `TOKENS`. Línea 43: `COLORES`, que solo se usa en el cruce secundario 3 (línea 301).
  - Línea 253: cruce §7.1. Línea 280: columnas de §10.5.
  - `main` devuelve código 1 si `bajo`, `difieren` o `avisos` no están vacíos.
- No hubo que regenerar nada: el script no escribe este `.md` (`grep -c 'contraste-componentes.md' contraste_componentes.py` = 0).
- Volví a correr toda la verificación del artefacto. Los 5 archivos generados quedaron byte a byte iguales a los anteriores.

## Archivos
- Modificado: `/home/deepsec/finanzas-hogar/docs/design/contraste-componentes.md` (solo las líneas 3 y 23). sha256 anterior `997a278b…1c`; nuevo `5c078362…7f`.
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-U2-r3.md` (este reporte).
- `npm run check` regeneró sin cambios `docs/contratos/C2-tokens.json`, `docs/contratos/tokens.css` y `docs/contratos/tailwind.theme.css`, con el mismo sha256 antes y después.

## Verificación
- **M-r3-1**:
  - `sed -n '3p;23p'` muestra el texto nuevo.
  - `grep -c 'C2-colores.json\`. No importa\|hex \*\*solo\*\*'` da 0: la redacción vieja ya no está.
  - `C2-colores` aparece solo en las líneas 3 y 23, y las dos veces como cruce secundario.
- **U2-04 / C2-05**: `python3 -I docs/design/tools/contraste_componentes.py --check` termina con código 0.
  - «CSV en disco IDENTICO al regenerado».
  - 158 pares y 316 combinaciones: 304 con umbral y 12 decorativas exentas. **0 bajo umbral**.
  - Cruce con contraste.csv: igual 266, nuevo 50 y difiere 0. Columna «en contraste.csv» que no coincide: 0.
  - Cruces secundarios (§7.1, columnas de §10.5 y C2-colores.json) que no coinciden con C2-tokens.json: 0.
- `--independiente`: 10 razones con segundo cálculo, todas coinciden, código 0. Ejemplos: C-061 claro 3.3827, C-101 claro 5.1269 y C-141 claro 5.0785.
- **U2-01 / U2-05 / U2-10 / C2-04**: `python3 -I docs/design/fuentes/scripts/verificar.py .` termina con código 0 y **RESULTADO: OK**.
  - 11 encabezados `##` en orden: los 10 de plan §4.1 y Supuestos.
  - 0 referencias a terceros.
  - Los 19 contrastes de §2.8 coinciden con el recálculo; por ejemplo, neutro 600/50 = 3.71 y marca 300/neutro 950 = 10.37.
- **C2-01 a C2-10**: `cd docs/design/tools/tokens && npm run check` termina con código 0.
  - generate: escribe 100073 B, idéntico al anterior.
  - build: Style Dictionary escribe tokens.css y tailwind.theme.css sin advertencias.
  - verify: **RESULTADO: OK**.
    - C2-02: 445 tokens; 0 sin tipo, 0 sin `$value` y 0 referencias rotas.
    - C2-03: 76 semánticos con claro y oscuro, 0 incompletos.
    - C2-04: 190 colores con 0 diferencias de hex; 158 pares de DESIGN con 0 diferencias.
    - C2-05: 170 combinaciones, 0 bajo umbral.
    - C2-06: `--ui-*` 78/78/78, 0 diferencias.
    - C2-08: 0 URL a terceros; las fuentes y la licencia existen.
    - C2-10: 0 hex repetidos sin alias (36 alias de paso) y 0 semánticos literales. Los 18 grupos de medida y la línea de 6 `font.line-height.*` sin cita literal son informativos y no cambiaron, porque no toqué DESIGN.md (sha256 `6eb3700c…` igual).
  - check:sync: tokens.css `2b9429bc…708d` y tailwind.theme.css `f66eb866…c88a`, iguales a los versionados.
  - test:tailwind: tailwindcss v4.3.3, «Done».
- `sha256sum -c` contra lo tomado antes de editar da OK en C2-tokens.json (`570ff553…`), tokens.css (`2b9429bc…`), tailwind.theme.css (`f66eb866…`), contraste-componentes.csv (`df8cf761…`) y DESIGN.md (`6eb3700c…`). Solo cambia contraste-componentes.md, que es el cambio buscado.
- `python3 -I -m py_compile` de los 4 scripts del alcance (`contraste_componentes.py`, `tokens/generar_c2_tokens.py`, `tokens/verificar_c2.py` y `fuentes/scripts/verificar.py`): OK.
- **U2-02, U2-03, U2-06 a U2-09 y U2-11 a U2-14**: no cambian. No se tocó ningún token, ni DESIGN.md, ni el CSV. La evidencia de REVISION-U2-r3 §1 sigue vigente.
- Temporales borrados: `datos/tmp/c2-prueba.css`, que crea `test:tailwind`, y los `__pycache__` de py_compile, que git ignora. `datos/tmp/c2-sync` no quedó. `datos/tmp/rev-s2` y `datos/tmp/motor-s2.json` son de otros agentes y no los toqué.

## Supuestos
- (ninguno) La corrección es la que sugirió el revisor, palabra por palabra. Ninguna cifra ni ningún token cambian.

## Bloqueos
- Ninguno nuevo. El único bloqueo abierto del proyecto, `docs/estado/bloqueos/S1-parametros-avisos.md`, no afecta a U2.

## Pendientes
- Quedan abiertos los bajos de REVISION-U2-r3, que están fuera de esta tarea porque solo se pidieron los altos y medios:
  - B-1: `DESIGN.md` §1.2 P5 cita §3.4, y debe ser §3.3.
  - B-2: §2.7 debe aclarar «mínimo en gris y en las tres simulaciones».
  - B-3: §2.2 dice ±30°; los ángulos reales son −25° y +35°.
  - B-4: §2.5, línea 116: decir «relleno 5.00:1 sobre tarjeta; texto blanco 5.59:1».
  - B-5: `docs/design/tools/tokens/verificar_c2.py` debe comparar los valores no de color de las tablas de DESIGN.
  - B-6: §7.2–§7.15 deben citar `ui.*` o dar una tabla de equivalencias en §7.0.
  - B-7: quitar «en construcción» del encabezado de DESIGN.md.
  - B-8: `docs/design/contraste-componentes.md` §4 debe regenerar «los 15 más justos» ordenados por margen. Faltan C-111c y C-056c, y sobran C-072 y C-066.
  - B-9: `docs/estado/wf1/U2c-tokens.md` necesita una nota «superado por correccion-U2-r1.md», porque su Supuesto 5 y su sha256 de tokens.css ya no son los vigentes. Esta ruta está fuera del alcance de esta corrección.
  - B-10: DESIGN §2.1 debe decir «medido en el hex: hasta 0.141 por redondeo de 8 bits».

  El revisor recomienda corregir B-1 a B-4, B-7, B-8 y B-10 antes del paquete de la Pista 3.

## Solicitudes fuera de alcance
- Ninguna.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna. Esta corrección no toca valores de K31.

## Consumo aproximado
- Acciones de herramienta: 16.
- Leído: unos 95 KB (protocolo, §K de decisiones, criterios U2, REVISION-U2-r3, correccion-U2-r2, cabecera y extractos de `contraste_componentes.py`, primeras 40 líneas de contraste-componentes.md y salidas de verificación).
- Escrito: unos 7 KB (2 líneas de contraste-componentes.md y este reporte). Los archivos que regeneró `npm run check` salieron idénticos.
