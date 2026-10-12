# revisor-shell-r1 · Pista 3, fase 2 (lienzo) · vuelta 1 de 3 · 12-oct-2026

Soy un revisor independiente: evalúo y no corrijo. El informe completo, con la tabla artboard × regla, el contraste, la muestra de datos y la fidelidad, está en `docs/estado/wfp3/P3-iter-1-shell.md`.

## Hecho
- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/shell/` y el fragmento `fragmentos/shell.json` contra:
  - P3-08 y P3-12 (`docs/criterios/P3-claude-design.md` 1.3);
  - U5-04, U5-06 y U5-07, en lo que se puede ver en el marcado;
  - DESIGN.md §4.6, §5.4, §6.3, §7.0, §7.2, §7.3, §7.9 a §7.15 y §10.5;
  - U1 §1.1 a §1.4, §2.1, §2.2, §2.4, §3.1, §3.2, §3.5, §3.24, F-1, F-6 y §5;
  - `referencia-tipo/format.md`.
- **Veredicto: no aprobado.** 0 altos, 2 medios y 7 bajos.
  - **S1-01 (media)**: en móvil, los `<input>` miden 32 a 40 px de alto propio dentro de contenedores de 44 a 64 px: `#monto`, `#peso`, `#usuario`, `#contrasena` y `#q-m`. El resto de los controles móviles cumple 44 px.
  - **S1-02 (media)**: en escritorio, «Ajustes» lleva a `ajustes-p26-indice-e`. Según U1 §1.2.3 debe abrir «01 Cuenta y roles», y ese artboard ya existe (`ajustes/ajustes-p26-s01-cuenta-roles-e.dc.html`). Aparece en 4 `href`: barra y tarjeta de Main, y las barras de fondo de P-03 e y P-22 e.
  - Bajos:
    - S1-03: el menú «/» muestra un subconjunto sin aviso de que hay más, fuera del orden de U1 §5.2.
    - S1-04: a los avisos les faltan `undo-2` y `x`.
    - S1-05: el diálogo de P-22 usa el ancho corto; los errores van en `small`; `#peso-e` hereda 16 px.
    - S1-06: «Buscar» y «Registrar gasto» de la barra de escritorio llevan los atajos dibujados dentro.
    - S1-07: «08:10» es un valor derivado y a «Tienda» le faltan 2 tiendas.
    - S1-08: la portada de Main no es una pantalla de U1.
    - S1-09: con movimiento reducido, `loader-circle` se detiene.

## Verificación
Todo lo que sigue lo comprobé con scripts propios (`python3 -I`, en `datos/tmp/rev-shell-r1/`, ya borrados).
- **Forma (8 de 8)**:
  - líneas obligatorias exactas, `lang="es-MX"` y un `<title>` propio por artboard;
  - el JSON de `data-props` parsea y `$preview` coincide con el fragmento;
  - `html.parser` no encuentra etiquetas sin cerrar;
  - stems únicos en todo el lienzo y nombres válidos;
  - el fragmento lista exactamente los 8 archivos.
- El verificador detecta errores: lo probé sobre copias alteradas y marcó los 10 tipos de error que inyecté.
- **Tokens**: 0 hex literales fuera de `data-props`, y 0 `rgb()` y 0 `hsl()`. Todos los `var(--x)` existen en `tokens.json` (365 nombres más `font-sans`) y en `tokens.css`.
- **Contraste**: 46 pares en claro y oscuro, recalculados con los hex de `tokens.json`: **0 bajo umbral**.
- **P3-12**:
  - 0 cifras canon (con y sin coma) y 0 credenciales;
  - 0 montos ≥ 1,000 sin separador;
  - las fechas del shell (19 y 21 de octubre de 2026) no coinciden con el calendario real;
  - 10 de 10 valores de la muestra salen de `datos-sinteticos.json`; «08:10» es derivado (S1-07).
- **Enlaces**: las 116 referencias existen y apuntan al mismo ancho, o a «#». La salvedad es S1-02.
- **Anchos**: en móvil la raíz es fija de 412 px con border-box. En escritorio es fluida, con `flex-wrap`, y el fragmento lleva `expand:"fill"`.

## Supuestos
1. Mido el objetivo táctil de los campos por el alto de la caja del `<input>`, calculado con interlineado más relleno porque K45 no permite renderizar, y no por el contenedor dibujado. Si el orquestador prefiere contar el contenedor, S1-01 baja a «baja».
2. S1-02 lo clasifico como media porque el prototipo lleva a un destino distinto del que fija U1, y ese destino ya existe en el lienzo. Lo trato como un enlace roto, aunque el archivo al que apunta sí existe.
3. Leo «Error» en las tarjetas «Estados» como palabra española y como nota del lienzo, no como texto en inglés (U5-07).
4. El fondo atenuado bajo el velo (`aria-hidden` o `inert`) es decorativo y no entra en el contraste.

## Bloqueos
- Ninguno.

## Pendientes
- Autor del shell, vuelta 2: corregir S1-01 y S1-02 (los dos son de pocas líneas) y, si cabe, los bajos S1-03 a S1-09.
- Orquestador:
  - decidir S1-07 (agregar la hora de desbloqueo al JSON sintético o aceptar el supuesto);
  - copiar este informe a `docs/design/salida/iter-1/` (P3-06);
  - al renderizar, confirmar las alturas de escritorio (P-00 e podría pasar de 1000).
- Dueño de los datos: `personas[axel].rol` («Titular») no coincide con `ajustes.usuarios[axel].rol` («Co-titular»).

## Consumo aproximado
- Leídos: unos 240 KB (8 artboards con 104 KB; P3 y las secciones U5-04, U5-06 y U5-07; format.md; secciones de U1 y DESIGN; extractos de tokens y del JSON sintético; reporte del autor).
- Escritos: unos 38 KB (este reporte, el informe P3-iter-1-shell y los scripts temporales).
- Acciones de herramienta: unas 40.
- Navegador y Artifact: 0.
