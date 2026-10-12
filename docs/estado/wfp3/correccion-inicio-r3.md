# correccion-inicio-r3 (vuelta 4, Pista 3, área Inicio)

Fecha: 12-oct-2026. Alcance: artboards de `docs/design/maquetas/lienzo/project/inicio/`, fragmento `inicio.json`, P-07 (`presupuesto-p07-flujo-m/-e`, porque muestran la quincena con el ahorro restado) y este reporte.

## Hecho
- Copia previa (P3-06) en `docs/design/salida/iter-3/inicio/` (8 `.dc.html` + `inicio.json`) y `docs/design/salida/iter-3/presupuesto/` (los 2 P-07).
- Hallazgo A-3 (alta, U3 §3 D-2 / S2-interfaz 1.3 §3.11 / U1 §2.2 / P3-12), con los datos corregidos de `correccion-datos-r4`:
  - `inicio-p01-detalle-m` y `-e`, panel D-2 «Quincena en curso»: «Gasto variable estimado —» y «Ahorro planificado —», «Saldo de la quincena $6,045.50» (igual a la tarjeta 2 del Resumen) y «Saldo acumulado $12,827.78». Se agrega un `caption` con el motivo: «Gasto variable estimado y ahorro planificado: pendientes de definir. No se restan de ningún saldo.»
  - Alto del fragmento y `$preview` de detalle ajustados por el renglón nuevo: móvil 4250 -> 4298 (también `min-height` de la raíz), escritorio 3000 -> 3032.
  - `presupuesto-p07-flujo-m` y `-e`: la serie `F.q` embebida se reemplazó por `proyeccion.flujo_por_quincena` actual (24 quincenas, mismas llaves). Con `null` en la 2.ª quincena de oct-2026 se muestra «—» (móvil: «— (pendiente de definir)»; escritorio: «—» en las columnas Variable est. y Ahorro plan.), con una nota al pie con el motivo. Saldo +$6,045.50 y acumulado $12,827.78; KPI «al cierre de oct-2026» $12,827.78 y «acumulado más bajo» $4,171.54 (dic-2026, 1.ª quincena), todo calculado desde los datos.

## Verificación (python3 -I)
- Inicio, 8/8 artboards (`/tmp/fh-inicio/verificar.py`): líneas obligatorias, `data-props` parsea y `$preview` = w x h del fragmento, `html.parser` sin etiquetas sin cerrar, 0 hex fuera de `data-props`, todos los `var(--x)` existen en tokens.css, 0 cifras prohibidas P3-12 (con y sin coma), dinero bien formado, raíces móvil/escritorio, enlaces. Fragmento válido, lista exactamente los 8 archivos. RESULTADO: OK.
- P-07, 2/2 (`/tmp/fh-r3/chk.py`): `data-props` parsea, parser sin etiquetas abiertas, 0 hex, 0 cifras prohibidas, `var()` existentes, sin «11,827.78» ni 1182778, serie con `null` en la quincena en curso. OK.
- `renderVals()` ejecutado con node en los 4 artboards: sin error; P-07 produce el renglón de la 2.ª quincena de oct-2026 con «—», +$6,045.50 y $12,827.78.
- Sin restos de $5,045.50, $11,827.78, 504550 ni 1182778 en ningún artboard del lienzo (grep).

## Supuestos
- El $12,827.78 es el acumulado correcto con el alcance pedido (el $13,827.78 del encargo es un desliz, según el aviso del orquestador).
- Los altos nuevos son estimados; se miden al publicar.

## Bloqueos
- Ninguno.

## Pendientes
- Orquestador: el alto de P-07 en `presupuesto.json` (móvil) crece unos 20-30 px por la nota nueva; el fragmento de presupuesto no es mío y no lo toqué.
- Los hallazgos bajos B-1 a B-12 de la vuelta 3 quedan fuera de esta ronda (solo mecánicos alta y media).

## Consumo aproximado
~75 KB leídos, ~10 KB de artboards modificados más este reporte, 17 acciones de herramienta.
