# revisor-U4-r1 · reporte del agente

Flujo P2c · ronda 1 de 3 · 11-oct-2026, noche (cierre pasada la medianoche) · informe completo en `docs/estado/p2c/REVISION-U4-r1.md`.

## Hecho
- Revisé `U4-ajustes-a.md` 1.5, `U4-ajustes-b.md` 1.4 (las líneas cambiadas), el `U4-ajustes.md` unido, el diff de `U3-dashboard.md` 1.3 y el diff de `DESIGN.md`, contra U4.md, U3.md, PUERTA-2 §2.4, §3 (Q2-11 y Q2-66 a Q2-70) y §5, con RESPUESTAS.md como mandato.
- **Veredicto: NO APROBADO**, con 0 altos, 5 medios (M-1 a M-5) y 2 bajos.
  - Lo que se asignó a `U4-puerta2` está bien aplicado.
  - M-1 (preselección del alcance sin decisión ni PROVISIONAL) y M-2 (dato falso sobre C1 v1.2) son de su alcance.
  - M-3 a M-5 son marcas o citas que la propia puerta 2 y el C1 v1.3 en curso dejaron viejas. Hay que asignarlos.
- No corregí nada. No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`. Sin git add ni commit, sin Telegram y sin navegador.

## Verificación (mis comandos)
- **Esquema:** `goose -dir backend/db/migrations sqlite3 datos/tmp/rev-u4-r1/t.db up` → versión 4.
  - `app_setting` tiene 30 claves; `params_json` de los reportes es `{}`.
  - `interval_count` lleva `CHECK BETWEEN 1 AND 24`.
- **`python3 -I verif.py`:**
  - unido == encabezado (38) + a + "\n" + b: True;
  - P-26: 7 + 7 = 14;
  - 210 `operationId` en los extractos 1.2, 73 de Ajustes;
  - U4 cita 109, con 0 inexistentes y 0 de Ajustes sin cita;
  - 28 claves citadas, 0 inexistentes.
- **`python3 -I u3.py`:** 6 tarjetas, 0 ISO y 0 hex. Solo $138.50 y $2,921.25 no aparecen literales en `model.json`/`data.json`; los dos son de antes de P2c y se derivan con su fórmula.
- **Q2-70 en milésimas enteras:** 1,486,571; 32 (cociente 32.99999…); 315,333.
- **Diff de DESIGN por script:** 4 líneas, 0 valores, hex o tokens cambiados. Los tokens C2 no tienen cambios.
- **`make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`:** rc 2.
  - Pasan WF-0, S1 («RESUMEN TODO OK»), U1, Criterios, Go con K25, Oráculos S2 a S5 (0 diferencias), Recetario, Tokens, **U4** y Paquete, con «WF-1: verificación completa OK».
  - **Falla solo `verify-c1`:** «21 archivo(s) difieren del generador». El generador de C1 v1.3 está modificado y sus salidas no se regeneraron; no viene de U4.
- **`build.py --out datos/tmp/…`** (C1 v1.3 candidato, sin escribir en el repositorio): 212 operaciones, 75 de Ajustes. U4 no cita `getReportSettings` ni `updateReportSettings` (M-5).
- Borré `datos/tmp/rev-u4-r1/` al terminar.

## Supuestos
1. «Marcas PROVISIONAL de Q2-78 a Q2-84» en DESIGN = las «PROPUESTA… lo decide Axel», como las entiende el agente. La l. 933 queda fuera del diff (B-2).
2. Para M-5 tomé como referencia el C1 v1.3 que genera el árbol de trabajo. Si `C1-puerta2` cambia la operación de las claves `reports.*`, M-5 se ajusta a lo que quede.
3. «PROVISIONAL sin bloqueo abierto» (media, según la escala de la tarea) incluye las marcas que citan un bloqueo que la puerta 2 dejó RESUELTO.

## Bloqueos
No abrí ninguno. M-1 propone abrir `U4-reparto-preseleccion` o sumar la pregunta a `C1-reparto-alcance`; eso le toca al corrector.

## Pendientes
- Ronda 2: comprobar M-1 a M-5 y B-1 a B-2, y repetir el cruce de `operationId` contra el `openapi.yaml` 1.3 cuando se regenere.
- Siguen abiertos los bajos R3-2 a R3-6 de WF-1 (nadie los pidió).

## Solicitudes fuera de alcance
1. **Orquestador:** asignar M-3 (U4: `S2-metas`, `S3-formula-tarjeta` y `S3-convencion-prestamo`) y M-4 (U3: `S2-presupuesto-total`, `S2-flujo-quincenal` Q2-14, `S2-metas` y `S2-proyeccion-interruptor`). Las respuestas son de las Q2-13 a Q2-32, que no estaban en el alcance de `U4-puerta2`.
2. **C1-puerta2 / integrador:** regenerar C1 v1.3; hoy `verify-c1` falla. Avisar a U4 la operación final de las claves `reports.*` (M-5).
3. **Dueño de DESIGN o integrador:** l. 933 (B-2).

## Propuestas para Axel
- Elegir la preselección del diálogo de alcance del reparto (P-26.12): «Solo a los gastos nuevos», «También a los ya capturados» o ninguna (M-1).

## Consumo aproximado
Unas 47 acciones de herramienta, unos 210 KB leídos (por partes, con `sed -n`, `grep` y scripts) y unos 18 KB escritos (este reporte y el informe; scripts temporales ya borrados).
