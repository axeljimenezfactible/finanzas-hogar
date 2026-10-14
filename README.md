# Finanzas Hogar · rama ligera para leer desde la web (gemini-ligero)

Esta rama es una copia **solo de texto** del repositorio, pensada para importarla en Gemini (navegador) sin que el ZIP se trunque: no trae capturas PNG, PDF, fuentes, bases de datos ni dependencias. **No es la rama de trabajo**: el proyecto completo vive en la rama `main`.

Por dónde empezar:
1. `docs/design/entrada-gemini/BRIEF-WEB.md` (cómo trabajar y cómo entregar sin terminal).
2. `docs/design/entrada-gemini/CONTEXTO-TABLERO.md`: un solo archivo con todo lo necesario para el tablero de Inicio; si la importación vuelve a cortarse, adjunta solo ese archivo.
3. `docs/design/entrada-gemini/PROMPT-WEB.txt`: el texto que se pega en el chat.

Qué trae: especificación del tablero (`docs/design/U3-dashboard.md`), diseño (`DESIGN.md`, `U1-arquitectura.md`), tokens y contrato (`docs/contratos/`), el código del frontend (`frontend/`), las pruebas, los bloqueos vigentes del tablero y el backend del tablero (solo lectura).
