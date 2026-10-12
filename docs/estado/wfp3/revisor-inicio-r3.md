# revisor-inicio-r3 (vuelta 3 de 3, Pista 3, área Inicio)

Fecha: 12-oct-2026. Revisor independiente; no corregí artboards ni datos.

## Hecho
- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/` y `fragmentos/inicio.json` contra P3 1.3 (P3-06, P3-07, P3-08, P3-12), U5-04/06/07, DESIGN.md, U1 §2.1/§2.2, U3 y `referencia-tipo/format.md`.
- Informe P3-07: `docs/estado/wfp3/P3-iter-3-inicio.md`.
- Veredicto: **NO APROBADO**, con 1 alto, 0 medios y 10 bajos.
  - A-2 (vuelta 2, decimal en D-6 Metas): **cerrado**.
  - A-3 (nuevo, alta): D-2 «Quincena en curso» muestra «Ahorro planificado $1,000.00» y «Gasto variable estimado $0.00» y resta el ahorro del saldo («Saldo de la quincena $5,045.50»). U3 D-2 y S2 1.3 §3.11 piden «—» mientras llegue `null`, y que no entre a ningún saldo. El bloqueo `S2-ahorro-quincenal` sigue abierto. La tarjeta 2 dice «$6,045.50 al cierre» para la misma quincena. El origen está en `datos-sinteticos.json` `proyeccion.flujo_por_quincena`.
  - B-12 (nuevo, baja, datos): las décimas de Metas no están en el JSON, y `LEEME.md` regla 4 sigue con «enteros».

## Verificación
- Scripts propios en `datos/tmp/rev-inicio-r3/`, todos borrados al final: `rv.js` (node), y `rev.py`, `c10.py` y `txt.py` (python3 -I).
- Forma 8/8 OK: líneas obligatorias, `data-props` válido, `$preview` = fragmento, `html.parser` sin etiquetas abiertas, `lang="es-MX"`, `<title>`, stems únicos.
- Tokens: 0 hex y todos los `var()` en `tokens.json`.
- Contraste: 22 pares de texto por tema, 0 bajo 4.5 (mínimo 4.80); 43 pares de UI por tema, 0 bajo 3 salvo los decorativos exceptuados. Recalculé 13 pares a mano.
- P3-12: grep canon = 0 y credenciales = 0. 12 montos trazados al JSON. Fechas coincidentes revisadas: ninguna va con un monto real.
- Enlaces: todos válidos. Táctil ≥ 44 px en móvil. Dinero tabular con U+2212. Sin emoji, sin fechas ISO visibles y sin barra de estado falsa.
- `diff` contra `salida/iter-2/inicio/`: solo cambiaron las 8 líneas de A-2.

## Supuestos
- El móvil y el escritorio de una misma pantalla comparten `<title>`: el encargo pide «nombre de la pantalla». No lo cuento como falla.
- Las sumas de T-2 ($242,520.65 y $257,740.56) son sumas de datos a la vista (U3 §4) y no cifras inventadas.
- A-3 va en alta porque contradice el contrato S2 1.3 y U3 en la pantalla principal, y deja dos saldos distintos para la misma quincena. Se decide así una pregunta de un bloqueo abierto.

## Bloqueos
- Ninguno para la revisión. A-3 necesita primero la corrección del orquestador en `datos-sinteticos.json`.

## Pendientes
- Orquestador: A-3 (datos y aviso a presupuesto P-07), B-11, B-12, B-7, y la copia de esta vuelta e informe en `docs/design/salida/iter-3/` (P3-06), que está vacía.
- Autor: A-3 en detalle-m l. 100-103 y detalle-e l. 139-142; bajas opcionales.

## Consumo aproximado
~190 KB leídos, ~25 KB escritos (informe, reporte y scripts temporales), ~45 acciones de herramienta. cerrado_por_tope = false.
