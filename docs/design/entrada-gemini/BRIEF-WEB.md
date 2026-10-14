# BRIEF-WEB · Gemini desde el navegador (sin terminal) · Tablero de Inicio · Finanzas Hogar

Versión 1.0 · 14-oct-2026. **Este archivo manda sobre `BRIEF-TABLERO.md` en todo lo que dice cómo trabajar.** De `BRIEF-TABLERO.md` sigues vigentes las secciones 1 a 5: contexto, lecturas, **reglas duras**, mapa de rutas y API y la **lista de lo que falta**. Ignora su sección 0 (clonar, rama, push) y los comandos de su sección 6: en la web no puedes correrlos; los corre Axel.

## 1. Por qué existe este archivo
La primera vez, el repositorio completo llegó a la web como un ZIP que se truncó entre las capturas PNG de `docs/design/capturas-u5/`, antes de llegar a `frontend/` y a este brief. Dos arreglos ya hechos:
- **Rama ligera `gemini-ligero`** del mismo repositorio: solo texto (sin PNG, PDF, fuentes ni bases de datos). URL: `https://github.com/axeljimenezfactible/finanzas-hogar/tree/gemini-ligero`. El repositorio completo sigue en la rama `main`.
- **Archivo único `docs/design/entrada-gemini/CONTEXTO-TABLERO.md`** (también en la rama ligera): trae, uno tras otro y separados por `## ARCHIVO: <ruta>`, todos los archivos que necesitas leer. Si la importación del repo vuelve a truncarse, adjunta solo este archivo.
- Las capturas actuales del Inicio (`P-01-android-claro.png`, `P-01-android-oscuro.png`, `P-01-escritorio-claro.png`, `P-01-escritorio-oscuro.png`) **no** están en la rama ligera; Axel puede adjuntarlas aparte si las pides.

## 2. Lo que no puedes hacer y qué haces en su lugar
No puedes ejecutar `git`, `npm`, Playwright ni compilar. Por eso:
1. **Tu entrega son archivos o parches**, que Axel o Claude aplican y verifican. No la calificas de «terminada» hasta que ellos corran las pruebas.
2. **Nunca declares que algo pasa si no lo corriste.** En tu reporte, cada verificación va como «NO EJECUTADO (revisión estática)» con lo que sí revisaste a ojo: tipos, imports, nombres de props y funciones.
3. **No inventes símbolos.** Usa solo funciones, props, rutas, claves de textos y tokens que aparezcan en los archivos que tienes. Si necesitas leer un archivo que no está en el contexto, no adivines su contenido: pídelo en un bloqueo (`docs/estado/bloqueos/GEMINI-<tema>.md`, formato de `BRIEF-TABLERO.md` regla 9) y sigue con lo que no depende de eso.
4. Si el contexto llega incompleto otra vez (archivo cortado, ZIP truncado), **detente y dilo en la primera línea de tu respuesta**: qué archivo es el último que alcanzaste a leer. No rellenes lo que falta.

## 3. Formato de entrega (obligatorio)
Entrega **por tandas** y espera «continúa» entre una y otra. Cada mensaje trae como máximo 5 archivos y termina con una línea «Tanda k de N · Pendiente: …».
- **Archivo pequeño o nuevo (hasta unas 300 líneas):** un encabezado `### ARCHIVO: ruta/relativa/desde/la/raíz` y debajo el archivo **completo** en un bloque de código con su lenguaje.
- **Archivo grande (por ejemplo `Tablero.tsx`, 422 líneas):** un diff unificado (`diff --git a/ruta b/ruta`, con 3 líneas de contexto) que se pueda aplicar con `git apply`. Dilo en el encabezado: `### PARCHE: ruta`.
- Rutas siempre exactas y relativas a la raíz del repositorio. No cambies nombres de archivos existentes.
- Tandas sugeridas: **T1** punto 1 de la lista (cadena P-04 → P-15 › Tendencias) con su prueba; **T2** puntos 2 a 5; **T3** puntos 6 a 8 y sus pruebas; **T4** la tabla de criterios A-01 a A-22 (punto 9), el reporte `docs/estado/gemini/tablero.md`, `docs/estado/gemini/solicitudes.md` y los bloqueos; el punto 10 (texto) solo si sobra espacio.
- El reporte final (`docs/estado/gemini/tablero.md`) lleva: Hecho, Archivos (rutas), Verificación (todo «NO EJECUTADO» salvo lo que sí revisaste, y cómo), Supuestos, Bloqueos y Pendientes, más el bloque de comandos de abajo.

## 4. Lo que corre Axel después (ponlo tal cual al final de tu reporte)
```
git checkout main && git pull
git checkout -b gemini/tablero
# aplicar cada ARCHIVO (copiar el contenido completo a su ruta) y cada PARCHE:
git apply <archivo.patch>
cd frontend && npm ci && npm run typecheck && npm run build && npm run test:e2e
```
Si algo falla, Axel te pega la salida y tú corriges solo lo que falló, en una tanda nueva.

## 5. Recordatorios que más se olvidan
- La interfaz **no calcula** porcentajes, diferencias ni estados; todo llega de la API.
- Texto de interfaz solo en `frontend/src/pantallas/inicio/textos.ts` o en `frontend/src/textos/es-MX.ts`; español de México, tuteo, sin emojis ni signos de exclamación.
- Colores y tamaños solo con tokens (`docs/contratos/tokens.css`); el color nunca va solo.
- Una tarjeta es un solo enlace: nunca un botón dentro de una tarjeta-enlace.
- No edites `docs/contratos/`, `docs/design/capturas-u5/` ni el backend. Lo que encuentres ahí va a `docs/estado/gemini/solicitudes.md` (ruta, cambio, motivo).
