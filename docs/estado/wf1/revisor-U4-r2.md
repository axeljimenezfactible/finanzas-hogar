# revisor-U4-r2 · Reporte final

Flujo WF-1 · 11-oct-2026 · Revisor independiente de U4, ronda 2 de 3. No escribió nada de lo que revisa y no corrige nada.

## Hecho
- Leí las lecturas obligatorias en orden:
  - el protocolo completo;
  - decisiones §A–§J y §K completas;
  - los criterios U4 1.2 completos;
  - las dos partes de U4 completas, por rangos con `sed -n`;
  - U1 §2.3, §2.4, §2.5.a y §2.5.b, §2.6, §5.2 y §5.3;
  - DESIGN §7.13;
  - los extractos de C1 por `grep` de `operationId` y de esquemas puntuales;
  - el inventario `@param` y los `INSERT INTO app_setting` de 00001 a 00003;
  - los bloqueos `U4-vista-previa-formulas`, `S2-metas`, `S3-formula-tarjeta`, `S3-convencion-prestamo` y `S1-parametros-avisos` (este último solo por `grep` de «Respuesta de Axel»);
  - los reportes `correccion-U4-r1`, U4a (Supuestos a Propuestas) y U4b (Pendientes a Propuestas), y la sección de hallazgos de `REVISION-U4-r1`.
- Volví a correr cada verificación con scripts propios (`datos/tmp/revisor-U4-r2/v.py` y `o.py`, con `python3 -I`; ya borrados).
- Escribí el informe `docs/estado/wf1/REVISION-U4-r2.md`: veredicto, tabla de los 13 criterios, 11 hallazgos y estado de los anteriores.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-U4-r2.md` (nuevo).
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-U4-r2.md` (este reporte).

## Verificación
- `v.py` (tareas a, f, g):
  - 210 `operationId` en los extractos, 73 de ellos de Ajustes.
  - Citados: 100 en la parte a y 45 en la b, **0 inexistentes**. **Ajustes sin citar: []**.
  - 14 secciones de P-26.01 a P-26.14, sin repetidas. Las 14 tienen vista de BD, edición en el lugar, vista previa, estados, atajos, 412 px y 1440 px.
  - Hex 0, `rgb(` 0, ISO 0, «al guardar» 0. V-16: 1 en a y 2 en b. «Suman más de 100»: 0.
  - 21 tokens citados, 0 inexistentes.
  - PROVISIONAL sin cita en la misma línea: a = []; b = líneas 26, 104, 209, 361 y 362 (bajas).
- `o.py` (cifras):
  - Las series A y B del oráculo aparecen literales en la parte a.
  - `calculado` = `esperado`; mínimos 8094.5 y 5000.0; `eventos[0]` 10000.0 y 6905.5; A − B = {3094.5}.
  - Pagos de `despensas[0]`: 4 × 973.75 = 3,895.00; intereses 800.50 (tanto con la tasa 25.8685 % como con 80050/309450).
  - 199 ÷ 2 = 99.50. Meta de calorías: 2,098.500 → 1,799.
  - RN-15: 1351429 y 1486571. RN-18: 50.0 %. K40: 0.9836 y 0.4918. Fechas: 2026-11-09 y 2027-01-08.
  - **0 diferencias contra el oráculo.**
- Cruce K44 (tarea c): 15 de 15 claves `@param` en su sección y con su valor sembrado. Los 4 de `params_json` están en el `@param_sin_valor` (líneas 56–59) y PROVISIONAL citando `S1-parametros-avisos`, que no tiene respuesta.
- Vistas y edición en el lugar (tarea d): coherentes con U1 §2.5 y §2.6, salvo los bajos de H-5 (R2-8).
- Atajos (tarea e): todos están en U1 §5.3, sin nuevos ni duplicados. Quedan los bajos de menú «/» de H-6 (R2-9).
- Muestra de 15 valores por defecto (tarea b): todos con origen en §A–§K o en la semilla.
- Criterios:
  - Cumplen: U4-01, 02, 03, 04, 06, 07, 08, 09, 10, 12 y 13.
  - **No cumplen: U4-05** (R2-1 alta; R2-2 y R2-3 medias) **y U4-11** (R2-2 y R2-3).
- No aplican pruebas Go, `gofmt` ni `go vet`: U4 es un documento de diseño.

## Supuestos
- Las vistas previas marcadas PROVISIONAL con un bloqueo que lo justifica cuentan como «definidas» para U4-05. Las que muestran como firme una cifra cuya fórmula está bloqueada, no.
- Un ajuste que cambia una fecha calculada que se muestra (próximo pago) tiene «efecto numérico».

## Bloqueos
Ninguno nuevo.

## Pendientes
Corregir R2-1, R2-2 y R2-3 en `docs/design/U4-ajustes-a.md` (P-26.05) y `docs/design/U4-ajustes-b.md` (P-26.13), y en `docs/estado/bloqueos/U4-vista-previa-formulas.md` (pregunta 2). Después, la ronda 3.

## Solicitudes fuera de alcance
- Integrador o S2: cuando Axel responda `S2-metas`, alinear RN-18 (un decimal) con C1 `progress_pct` (entero) y con la vista previa de U4 P-26.13.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna. No hay valores de K31 en esta revisión.

## Consumo aproximado
45 acciones de herramienta · ~330 KB leídos · ~27 KB escritos (informe, este reporte, progreso y scripts temporales).
