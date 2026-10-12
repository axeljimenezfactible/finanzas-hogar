# Bloqueo WF3-frontend-rutas · Direcciones (URL) de las pantallas

- Fase: WF-3a, agente `frontend-base` (12-oct-2026). Estado: **PROVISIONAL** (no detiene el trabajo).
- **Qué falta (PREGUNTA):** U1 §2.1 define pantallas P-00..P-26, subsecciones y vistas, pero no la dirección que ve la persona en la barra del navegador ni cómo se codifican subsección, vista y hojas.
- **Lo que quedó (PROPUESTA, en `frontend/src/app/rutas.ts`):** direcciones en español sin acentos que siguen el árbol de U1 §1.1: `/entrar` (P-00), `/inicio`, `/gastos`, `/gastos/:id`, `/gastos/ticket`, `/presupuesto`, `/presupuesto/flujo|ingresos|ahorro`, `/comidas`, `/comidas/recetas[/:id]`, `/compras`, `/compras/comparar|productos|inventario`, `/creditos`, `/creditos/:id`, `/creditos/:id/compras/:compra`, `/creditos/comparador`, `/salud`, `/reportes`, `/reportes/mensual|karen`, `/ajustes/:seccion` (01–14). Subsección `?sub=detalle` (sin parámetro = la de por defecto), vista `?vista=tablero`, hoja `?hoja=ingreso`. P-03 y P-22 son hojas sin dirección (`?hoja=gasto|peso` las abre).
- **Opciones:** (A) quedarse con la propuesta; (B) direcciones en inglés (`/expenses`…), coherentes con los identificadores (§8) pero visibles para la persona; (C) direcciones con el id de pantalla (`/p/02`).
- **Recomendación:** A (lo que se ve en pantalla va en es-MX, U5-07).
- **Impacto:** solo `rutas.ts` (un archivo) y los enlaces que usan `rutaDe()`; ninguna área escribe direcciones a mano. Accesos directos del manifest (`/inicio?hoja=gasto`).
