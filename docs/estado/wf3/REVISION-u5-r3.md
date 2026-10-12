# WF-3b · REVISION-u5-r3 · Hallazgos de la revisión visual (ronda 3 de 3)

Revisor `revisor-u5-r3` · 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo con `correccion-u5-r2`. Se evalúa solo contra `docs/criterios/U5.md`, DESIGN.md, U3 y U4. Combinaciones: Ac (Android claro), Ao (Android oscuro), Ec (escritorio claro), Eo (escritorio oscuro). Registro de la ronda: `docs/estado/U5-revision.md` § Ronda 3.

**Veredicto: devuelto** (0 altos, 2 medios). Es la ronda 3: se escala a Axel (`docs/estado/bloqueos/U5-ronda3-escalamiento.md`).

## Hallazgos de r2: estado
| Id | Sev. | Estado | Comprobación propia |
|---|---|---|---|
| N-1 · U5-07 · nota de semilla en P-26.13 | alta | **Cerrado** | innerText de P-26.13 en las 4: «Dato de desarrollo: no existe en producción.»; 0 «Fixture», 0 `.json`, 0 `[n]`, 0 claves `a_b` en las 164 cargas |
| N-2 · U5-06 #12 · capturas Android sin puntero grueso | media | **Cerrado** | `manifiesto.json` 82 de 82 con `puntero_grueso`; en mis 82 cargas `any-pointer: coarse` = true y 0 controles con área efectiva < 44 × 44 |
| N-3 · U5-06 #7 · sombra en secciones de Ajustes | media | **Cerrado** | `box-shadow` de las secciones = `none` en P-26.01 a P-26.14 × 4; fondo `#F2F2F4` / `#36363A`, borde 1 px, radio 14 |
| M-6b · U5-06 #14 · etiquetas encimadas en Tendencias | media | **Cerrado** | `--grep Tendencias` 4 de 4 en mi corrida; a ojo (Ec, Ao): «$18,025.47» a la derecha de «$18 mil» / «$21 mil»; T2 sin texto sobre el área y leyenda abajo con círculo y cuadro |
| M-1b · U5-06 #2 · h1 de P-00 | media | **Cerrado** | 32px/700 en las 4 (y en las 164 cargas) |

## Hallazgos nuevos

### M-7 · media · U5-06 #10 · Barra lateral sin iconos de sección
- **Pantallas**: las 40 con barra lateral (todas menos P-00) · **Ec, Eo**.
- **Regla**: DESIGN §7.9 («las 9 secciones en orden fijo con icono `size.icon.sm`, nombre en `font.size.small` y `chevron-right`»; reposo «icono `ui.icon.default`», activo «icono `ui.icon.selected`») y §5.4 (`house`, `utensils`, `shopping-cart`… «barra inferior y lateral»).
- **Qué se ve**: cada sección muestra solo `chevron-right` y el nombre (recorte a tamaño real de P-01 Ec y P-02 Eo). En el código, `frontend/src/shell/BarraLateral.tsx` importa de Lucide solo `ChevronRight, ChevronsLeft, LogOut, Monitor, Moon, Plus, Search, Sun, UserRound`.
- **Archivo probable**: `frontend/src/shell/BarraLateral.tsx` (lista `SECCIONES.map`) y el mapa de iconos por sección (los mismos de la barra inferior en `shell/Barras.tsx`).
- **Corrección sugerida**: agregar el icono de §5.4 de cada sección (20 px, trazo 1.8, `ui.icon.default` y `ui.icon.selected` en la activa) entre el `chevron-right` y el nombre, sin cambiar el ancho de 240 px.

### M-8 · media · U5-11 · La matriz de C3 no marca «capturado»
- **Dónde**: `docs/contratos/C3-integracion.md` v0.95 §2, columnas «Capturado» y «Evidencia (capturado)»: «no · pendiente WF-3b (`docs/design/capturas-u5/`)» en las 41 filas (P-00 a P-26.14).
- **Regla**: U5-11 («Cada pantalla capturada marca el estado «capturado» en la matriz de C3 con la ruta de sus 4 imágenes»). El cruce literal da 0 pantallas implementadas sin captura (las 41 tienen sus 4 PNG), pero falta la marca.
- **Corrección sugerida** (documental, no de la app): el integrador o `c3-regen` pone «sí» y las 4 rutas `docs/design/capturas-u5/<P-xx>-{android,escritorio}-{claro,oscuro}.png` en cada fila.

### Bajas nuevas
- **B-10** · U5-07 · P-26.11 · 4 comb.: «Versión **dev**». Es el valor por omisión de `backend/internal/api/sistema/sistema.go:24` (`Version = "dev"`); ningún script de `ops/` lo fija con `-ldflags`. No es palabra de la lista del criterio; lo dejo bajo como identificador de compilación. Sugerencia: mostrar «desarrollo» cuando la versión es `dev`, o fijar la versión al compilar.
- **B-11** · §7.10 y §4.5 · P-03 y P-22 · Ac Ec: la pestaña elegida del selector de tipo (Gasto · Ingreso · Peso · Ticket) lleva relleno `#F2F2F4` y `shadow.1`, además del indicador inferior de 2 px. La sombra es para lo que flota. Archivo probable: componente de pestañas de la hoja de registro (`frontend/src/componentes/`). Sugerencia: quitar la sombra y dejar el indicador.
- **B-12** · §7.3 · P-07, P-26.02, .04, .06, .08, .09, .10, .11 · Ac Ec: la pista del interruptor mide 40 × 24 (DESIGN pide 36 × 20 dentro de 44) y la perilla lleva `shadow.1` en claro. Junto con B-7.
- **B-13** · P-00 · 4 comb.: en una visita sin sesión previa (contexto nuevo) aparece «Tu sesión terminó. Vuelve a entrar para seguir donde estabas.» El aviso debería salir solo si hubo una sesión que expiró. Archivo probable: `frontend/src/app/IniciarSesion.tsx` / `sesion.tsx`.
- **B-14** · §7.9 · todas · Ec Eo: el encabezado de la barra muestra «Finanzas H…» truncado por «Axel». Sugerencia: el nombre del usuario en el menú del avatar, o reducir su ancho.
- **B-15** · §7.12 · P-03, P-22 · Ac: la hoja inferior no tiene sombra con desenfoque en el DOM (DESIGN: `shadow.2`); bajo el velo apenas se nota.
- **Observación (por confirmar, sin severidad)**: P-06 muestra «Quincena 2 · 16-oct-2026 – 31-oct-2026 Real $1,197.64», el mismo real que la quincena 1, aunque el único gasto es del 10-oct. No está entre las cifras de C3-09 ni de U5-02; conviene que lo mire el dueño de Presupuesto.

### Bajas que siguen
B-1 (destinos de la barra inferior, no remedido), B-2 (bloqueo `U5-nombres-propios`), B-3 (bloqueo `U5-tipografia-tres-tamanos`), B-5, B-6, B-7, B-8 (axe sigue dando `aria-prohibited-attr` incompleto en P-26.10). **B-9 cerrado**: 0 controles con área efectiva < 44 px en Android.
