# REVISION-u5-r2 · Revisión visual independiente U5, ronda 2

Revisor `revisor-u5-r2` · 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo. Registro por ronda en `docs/estado/U5-revision.md` §Ronda 2. **Veredicto: devuelto** (1 alta, 4 medias). Combinaciones: Ac Android claro · Ao Android oscuro · Ec escritorio claro · Eo escritorio oscuro.

## Hallazgos de la ronda 1
| ID r1 | Estado | Evidencia r2 |
|---|---|---|
| A-1 U5-03 axe | Cerrado | Comando literal 0; corrida propia 0 serias/críticas y 0 de otro impacto × 4 |
| A-2 U5-01 P-26.13 | Cerrado | 164 de 164; `scrollWidth = innerWidth = 412` |
| A-3 U5-04 UI | Cerrado | 0 bajo umbral en reposo; 45 estados elegidos ≥ 3:1 por combinación; barra de vista con `aria-pressed` |
| A-4 U5-07 inglés | Cerrado | P-26.11 «Elegir archivo CSV · Ningún archivo elegido»; P-26.04 «Despensa», «Comida fuera» |
| A-5 U5-13 H-5 | Cerrado | Spec propia: tras recargar sin red sale la franja «· 2 en cola» y se sigue capturando |
| M-1 h1 | **Parcial** | 32/700 en 40 pantallas; P-00 sigue a 600 (M-1b) |
| M-2 barra lateral | Cerrado | 240 px; «‹‹» en x 199–231, ningún texto debajo |
| M-3 iconos | Cerrado | 0 fuera de 20/24 o con trazo ≠ 1.5 en 164 |
| M-4 espaciado | Cerrado | Solo quedan la estructura (56 px) y B-7 |
| M-5 franja P-17 | Cerrado | 314 px, 2 líneas |
| M-6 gráficas | **Parcial** | Estados vacíos correctos; siguen etiquetas encimadas (M-6b) |

## Altos
### N-1 · U5-07 · texto en inglés y claves internas visibles en P-26.13
- Pantalla y combinación: P-26.13 «Salud y metas de ahorro» (tabla «Metas de ahorro y fondo») · Ac Ao Ec Eo.
- Qué se ve: «Dato de desarrollo (finanzas semilla --desarrollo): no existe en producción. Fixture s2-metas.json, savings_goals[2].» Incluye «Fixture», un nombre de archivo y la clave `savings_goals[2]`. El script de innerText del capturista no lo detecta porque esas palabras no están en su lista. Mi búsqueda de `[a-z]+_[a-z_]+` sí lo encuentra (1 coincidencia por combinación, solo en esa pantalla).
- Archivo probable: `backend/internal/semilla/desarrollo.go` (nota de la meta sembrada por `finanzas semilla --desarrollo`).
- Corrección sugerida: dejar la nota en español y sin nombres internos, p. ej. «Dato de desarrollo: no existe en producción». Después, volver a sembrar la base de desarrollo y regenerar las 4 capturas de P-26.13 (y P-09 si muestra la nota). Agregar `\b[a-z]+_[a-z_]+\b` y «Fixture» al script de U5-07.

## Medios
### N-2 · U5-06 #12 y fidelidad de U5-01 · capturas Android sin puntero grueso
- Pantallas y combinaciones: capturas Android (comprobado en P-02 y P-17; probablemente todas después de P-00) · Ac Ao · DESIGN §4.6 (`size.touch.min` 44, `size.control.coarse` 44).
- Evidencia: en `P-02-android-claro.png`, el botón primario «Registrar gasto» mide **36 px** de alto (y 204–239). En un contexto Android nuevo, `matchMedia('(any-pointer: coarse)')` da `true` y el mismo botón mide **44 px**: `Boton` usa `h-9 … coarse:h-11` y la regla está en `dist`. Reproduje el efecto: al reutilizar la página de sesión después de abrir y cerrar la página anónima de P-00, salieron 1,392 controles < 44 px. Con una página nueva por pantalla, el puntero grueso queda activo y solo hay 241 casillas, todas dentro de una etiqueta de 44 px (pasa), más B-9. La app cumple; las capturas Android no muestran la app como se ve en Android.
- Archivo probable: `frontend/pruebas-u5/capturas.spec.ts` y `apoyo.ts` (`abrirAnonima` y sesión única por proyecto).
- Corrección sugerida: capturar P-00 en otro navegador o al final, o abrir una página nueva por pantalla. Antes de cada captura Android, comprobar `matchMedia('(any-pointer: coarse)').matches === true` y guardarlo en `manifiesto.json`. Regenerar las 82 capturas Android y repetir axe (`target-size`) y contraste.

### N-3 · U5-06 #7 · secciones de Ajustes con sombra en reposo (claro)
- Pantallas y combinaciones: P-26.01 a P-26.14 · Ac Ec · DESIGN §4.5 (`shadow.0` en tarjetas en reposo; `shadow.1` solo en menús, arrastre y vista previa) y §7.13 (Reposo: «filas sobre `ui.surface.page` con divisor»).
- Evidencia: cada `section` mide `box-shadow` = `shadow.1` (0 1px 2px 8 % + 0 1px 3px 6 % de neutral-950), con fondo #F2F2F4, borde 1 px #E5E6E9 y radio 14. En oscuro no tiene sombra (correcto).
- Archivo probable: el componente de sección de Ajustes (`frontend/src/pantallas/ajustes-a/comunes.tsx` y `ajustes-b/comun.tsx`).
- Corrección sugerida: quitar `shadow-1` (`shadow.0`) y dejar la separación por superficie o divisor que pide §7.13.

### M-6b · U5-06 #14 · etiquetas de gráficas encimadas (reabre M-6 en parte)
- Pantalla y combinación: P-01 ▸ Tendencias · Ec y Ao (vistas en capturas mías de esa pestaña) · DESIGN §7.7 y §8.5.
- Qué se ve:
  - «Deuda viva proyectada»: la etiqueta «$18,025.47» cae sobre la etiqueta de eje «$18 mil» (Ec) y «$21 mil» (Ao).
  - «Devengado contra efectivo»: «Efectivo acumulado $19,397.22» (Ec) y «Efectivo acumulado» (Ao) se dibujan encima del área.
  - En Ec, la leyenda de «Efectivo acumulado» no muestra marca.
  - P-10 Tendencias, en cambio, está bien: muestra los dos estados vacíos.
- Archivo probable: `frontend/src/componentes/Grafica.tsx` (`FiguraGrafica`: etiqueta del primer punto y de fin de línea; `labelLayout`) y `inicio/GraficaInicio.tsx`.
- Corrección sugerida: no poner etiqueta al primer punto cuando choca con el eje, o dejar margen izquierdo (`grid.left`) para ella. Llevar a la leyenda las etiquetas de fin de serie de las áreas apiladas, porque §8.3 ya las nombra. Revisar el símbolo de leyenda de la serie de área en claro.

### M-1b · U5-06 #2 · título de P-00 a 600
- Pantalla y combinación: P-00 «Finanzas Hogar» · Ac Ao Ec Eo · DESIGN §3.2 (`font.size.title-1` = 32 px / 700, «solo título de página»).
- Archivo probable: la pantalla de entrada en `frontend/src/app/` (P-00, fuera de `shell/Barras.tsx`).
- Corrección sugerida: `text-title-1 font-bold`, como el resto de los h1.

## Bajos
- **B-1** (sigue): destinos de la barra inferior de 82×55 (DESIGN pide 56 de alto) · Ac Ao.
- **B-2** (sigue, bloqueo `U5-nombres-propios`): Tesseract, OCR, CSV, API, etc.
- **B-3** (sigue, bloqueo `U5-tipografia-tres-tamanos`): 29 (Android) y 30 (escritorio) pantallas usan más de 3 tamaños además de `caption`.
- **B-4**: no lo volví a revisar.
- **B-5** (sigue): P-13 «Sin plan esta semana» mientras P-10 tiene plan; P-23 «Resumen diario … último resultado: correcto» frente a P-26.08 «Aún no se ha enviado».
- **B-6** (sigue): tablas angostas en Android. En P-02 Ac la fecha «10-oct-2026» se parte en 3 líneas y el monto queda cortado en el borde.
- **B-7** (nuevo): el interruptor de P-07 usa 2 px de relleno interno (4 nodos), fuera de §4.1. Convendría declararlo medida de componente o usar 4 px.
- **B-8** (nuevo): en P-26.10, `<div aria-label="Muestra de tarjeta">` sin rol. axe lo da como incompleto `aria-prohibited-attr` y el nombre no llega a lectores. Agregar `role="img"` o `role="group"`.
- **B-9** (por confirmar): en P-26.09 «Seleccionar todas» (`ajustes-b/varias.tsx:18`, etiqueta `coarse:min-h-11 coarse:min-w-11`), mi sonda del área efectiva falla en un borde aun con el puntero grueso activo. Puede ser un control vecino que se empalma.

## Notas de método (no son fallas de la app)
- Las capturas de las 10:25 eran anteriores al binario de las 10:33. Volví a correr el comando completo, así que la evidencia de `capturas-u5/` es del binario actual.
- En Android, la captura de página completa deja la barra inferior fija a media imagen; eso tapa parte de las gráficas de Tendencias en Ao.
