# Bloqueo WF3-inicio-tiendas · Color de tienda sin código en C1

Fase: WF-3a · agente `front-inicio` · 12-oct-2026 · Estado: **ABIERTO**.

## Qué falta
C2 tiene un color por tienda (`color.store.walmart`, `bodega-aurrera`, `sams-club`, `market`) y DESIGN §2.6 pide punto de color + nombre, pero el esquema `Store` de C1 v1.4 no trae un código estable (solo `id`, `name`, `kind`). La tarjeta 5, D-5 y los filtros de tienda necesitan saber qué color le toca a cada tienda.

## Regla PROVISIONAL en el código (`src/pantallas/inicio/comun.tsx`, `codigoTienda`)
`kind` = `market` → `market`; si no, el nombre sin acentos ni apóstrofos, en minúsculas y con guiones, que empiece con `walmart`, `bodega-aurrera` o `sams-club`; cualquier otra tienda usa el punto de `market`. El nombre siempre va en texto, así que el color nunca es el único dato (K41).

## Opciones
A) Dejar la regla por nombre. B) C1 agrega `Store.code` (enum de C2) y la PWA lo usa tal cual. C) Mapa por `id` en Ajustes › Tiendas.

## Impacto
Solo el color del punto; ninguna cifra. Con B se cambia una función. Lo mismo necesitan las áreas de Compras (P-13 a P-16).
