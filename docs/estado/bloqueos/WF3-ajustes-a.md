# Bloqueo WF3-ajustes-a · Ajustes, secciones 01 a 07 (front-ajustes-a, 12-oct-2026)

Estado: ABIERTO. Lo marcado PROVISIONAL en `frontend/src/pantallas/ajustes-a/` cita este archivo.

## 1. Despensa programada (P-26.02): tipo de partida y alta
- **Falta**: U4 y C1 no dicen con qué `planned_item.kind` se guarda la despensa programada ni cómo se distingue de un gasto fijo.
- **Hoy (PROVISIONAL)**: la pantalla toma como despensa la partida que no es ingreso y tiene `accrual_months = 2` (esquema 00001: «despensa bimestral = 2»). La saca de la tabla de gastos fijos y la muestra en su bloque. Si no existe, dice «Sin despensa programada» y no ofrece el alta.
- **Opciones**: (a) `kind = fixed_expense` con `accrual_months = 2`; (b) `kind = upcoming_expense`; (c) una marca propia en C1, por ejemplo `planned_item.is_grocery_plan`.
- **Impacto**: no se puede programar la despensa desde Ajustes mientras no se decida. La semilla y el CSV sí la cargan.

## 2. Punto de equilibrio de la membresía (P-26.06)
- **Falta**: `getBasketComparison` exige `week` (la lista de la semana) y U4 no dice qué semana usa Ajustes.
- **Hoy (PROVISIONAL)**: sin costo capturado dice «Captura el costo anual para evaluar». Con costo, dice «Vista previa pendiente».
- **Opciones**: (a) la semana en curso de la lista de compras; (b) la última semana con lista; (c) una operación propia en C1 que no dependa de la semana.
- **Impacto**: la membresía de Sam's Club no se evalúa en Ajustes. En Compras sí puede verse.

## 3. P-26 compartida con la parte b (secciones 08 a 14)
- **Falta**: CONVENCIONES-frontend.md no dice cómo dos áreas comparten P-26. `registro.ts` junta los `pantallas` con `Object.assign`, así que si otra área también exporta `'P-26'`, gana la última en orden alfabético.
- **Hoy (PROPUESTA)**: `ajustes-a` registra P-26 (la columna de las 14 secciones). Cada área aporta sus secciones en `src/pantallas/<area>/secciones-ajustes.ts`, con `export const secciones = { '08': lazy(...), … }`. P26Ajustes los descubre con `import.meta.glob`. Una sección sin componente muestra «Pantalla en construcción».
- **Pedido**: el área de las secciones 08 a 14 no exporta `'P-26'` en su `indice.ts`, y frontend-base agrega esta regla a CONVENCIONES-frontend.md §3.
