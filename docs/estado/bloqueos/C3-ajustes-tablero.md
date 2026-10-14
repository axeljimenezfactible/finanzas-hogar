# Bloqueo C3-ajustes-tablero · `updateDashboardSettings` no tiene pantalla (C3-02, C3-03)

Fecha: 2026-10-12 · Agente: `c3-tablero` (Sonnet, WF-3c) · Contexto: `docs/contratos/C3-integracion.md` v1.3 (sin firma). Hay una pregunta abierta y nada de esto cambia cifras ni motores.

## Qué falta (HECHO)
- C1 v1.5 trae la pareja `getDashboardSettings` / `updateDashboardSettings` (GET y PATCH `/api/v1/settings/dashboard`; `x-necesidades`: N-278 la de escritura).
- U1 §3.3 (fila «Parámetros del tablero»: «las cinco claves de K54 en Ajustes: leer y cambiar», N-275 y N-278) y U4 la ponen en Ajustes › Presupuesto (P-26.02).
- En `frontend/src` solo Inicio llama `getDashboardSettings` (`pantallas/inicio/comun.tsx`, para N de «Pagos próximos»). **Ninguna pantalla llama `updateDashboardSettings`**; `S02Presupuesto.tsx` no tiene el apartado.
- `python3 -I frontend/scripts/c3-integracion.py` lo informa: «Operaciones con pantalla que el frontend no llama: 2 [listPeriods, updateDashboardSettings]; sin excepción justificada o sin cobertura: 1 [updateDashboardSettings]» y sale con 1.
- La fila P-26.02 sigue en «sí» en implementado porque la regla de C3 es por necesidad (N-278 la cubren `createBudget`, `updateBudget`, `createPlannedItem`, `updatePlannedItem`, `deleteBudget`); pero el apartado «Parámetros del tablero» no existe.

## PREGUNTA
¿Qué se hace con `updateDashboardSettings`?

| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. PROPUESTA (la recomendada): el agente de Ajustes (frontend, P-26.02) agrega el apartado «Parámetros del tablero» con las cinco claves de K54 (`debt.capacity_warning_pct`, `debt.capacity_limit_pct`, `dashboard.payment_reminder_count`, `budget.signal_months`, `budget.signal_margin_pct`), con «Restablecer», que llama PATCH con `parameters` y `reset_parameters` | C3-02 y C3-03 en verde sin excepción; cierra el ciclo K44 | Medio: una sección más con su prueba y sus capturas de P-26.02 | Bajo |
| B. Declararla excepción técnica PROVISIONAL (como `listPeriods`) | El script sale con 0 | Ninguno | Alto: queda el valor del barómetro (30 % y 40 %) sin forma de cambiarse desde la app, contra K44 y U1 |
| C. Quitar la operación de C1 o dejarla solo para el agente externo | Contrato más chico | Medio: cambia C1 y el generador | Contradice U1 y U3 §2.6 |

Impacto: sin decidir, C3 no se firma con 0 filas abiertas en C3-03 (el script falla por esta operación). **No se tomó ninguna decisión**: el script no la admite como excepción; el pendiente queda a la vista en la tabla al pie de la matriz. Decide Axel o el orquestador.
