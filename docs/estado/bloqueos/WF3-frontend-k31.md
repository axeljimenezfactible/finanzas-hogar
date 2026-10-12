# Bloqueo WF3-frontend-k31 · Valores operativos del frontend (puerta 4)

- Fase: WF-3a, agente `frontend-base` (12-oct-2026). Estado: **PROVISIONAL** hasta la puerta 4 (K31, K44: todo valor configurable).
- **Valores propuestos (PROPUESTA) y justificación:**
  1. Sondeo de `/health` mientras no hay servidor: **cada 15 s** (`src/api/conexion.ts`). Detecta el regreso a la red de casa aunque el teléfono siga «en línea» con datos móviles; costo despreciable en LAN.
  2. Reintento automático de la cola: **cada 30 s** mientras haya pendientes, además de al volver la red y al abrir la app (`src/cola/cola.ts`).
  3. Copia de la última carga de cada lectura (IndexedDB `lecturas`): se conserva **hasta cerrar sesión** (U1 §2.4 «lo ya cargado sigue visible»; la página nunca cargada lo dice). Alternativa: caducar a los 7 días.
  4. La cola **no** se borra al cerrar sesión (U1 §2.4: nunca se pierde ni se borra sin confirmación); se envía al volver a entrar.
- **Impacto:** cambiar cualquiera es una constante (o, tras la puerta 4, un valor de Ajustes › Datos y respaldos).
