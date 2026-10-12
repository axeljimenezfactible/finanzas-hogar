# Bloqueo S6-semilla-sucursales-municipio

Contexto: WF-3a, agente `semilla`. `data.json` solo trae el nombre de las dos sucursales de la compra del 10-oct-2026 («Bodega Aurrera (sucursal Vallejo)» y «Walmart (sucursal Acueducto de Guadalupe)»): sin dirección, municipio ni estado. El filtro de PROFECO (`internal/precios/profeco`, S5-interfaz §11) deja pasar un precio solo si la cadena y el municipio de alguna sucursal coinciden (bloqueo `S5-profeco-filtro`).

Qué falta: ¿en qué alcaldía o municipio y estado están esas dos sucursales (y qué nombre comercial usa PROFECO para cada una, si existe)?

Opciones:
- A: Axel los captura en Ajustes → Tiendas y precios (`updateStore`); entonces se vuelve a correr la importación de PROFECO. Consecuencia: la semilla no inventa nada; hoy 0 de 1,974 renglones de la muestra entran.
- B: Axel los dicta ahora y se agregan a la semilla. Consecuencia: la muestra podría cargar precios solo si además existe `profeco/vinculos.csv` confirmado (no existe).
- C: sembrar con el municipio que yo suponga por la colonia (Vallejo, Acueducto de Guadalupe). Descartada: es suponer (regla suprema).

Impacto: la muestra de PROFECO queda en 0 precios cargados (ver AVISO del reporte). `profeco_name` de las cadenas sí se sembró («Bodega Aurrera», «Wal-mart»), tomado de `cadena_comercial` de la propia muestra; PROVISIONAL hasta importar PROFECO completo.

Cómo seguí: opción A en la práctica; municipio, estado y dirección quedan NULL.

- **Estado:** ABIERTO
