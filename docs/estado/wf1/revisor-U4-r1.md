# revisor-U4-r1 · reporte de agente (WF-1, 11-oct-2026)

Revisor independiente de U4, ronda 1 de 3. No escribí nada de lo revisado y no corregí nada. Informe completo: `docs/estado/wf1/REVISION-U4-r1.md`.

## Hecho
- Leí completos `docs/design/U4-ajustes-a.md` y `docs/design/U4-ajustes-b.md`, `docs/criterios/U4.md` 1.2, `docs/plan/decisiones.md` §A–§K y el protocolo; U1 §2.3–§2.6 y §5.2–§5.4 por `sed -n`; índice de extractos y esquemas puntuales de `ajustes-a/c/d/e.yaml` por `grep`; inventario `@param` y siembras de 00001–00003 por `grep`/`sed -n`; reportes U4a y U4b; bloqueos `U4-*`.
- Reejecuté las verificaciones con scripts propios (`ops.py`, `tok.py`, `chk.py` en `datos/tmp/revisor-U4-r1/`, borrados al final) y recalculé con `Decimal` todas las cifras de ejemplo.
- Veredicto: **NO APROBADO** (1 alta, 2 medias, 5 bajas).

## Archivos
- `docs/estado/wf1/REVISION-U4-r1.md` (nuevo)
- `docs/estado/wf1/revisor-U4-r1.md` (este reporte)

## Verificación
- `python3 -I datos/tmp/revisor-U4-r1/ops.py .` → 210 `operationId` en `openapi.yaml` (solo la línea `operationId:`); parte a 99 citados, 0 inexistentes; parte b 45, 0 inexistentes; 73 operaciones de Ajustes, 0 sin campo o acción. **U4-02: cumple.**
- `python3 -I datos/tmp/revisor-U4-r1/tok.py .` → 0 hex, 0 `rgb(`, 0 fechas ISO; todos los tokens existen en C2 o DESIGN. **U4-09 (parte medible) y U4-11: cumplen.**
- `python3 -I datos/tmp/revisor-U4-r1/chk.py .` → 14 secciones en orden con nombres de decisión G (**U4-01**); 15/15 `@param` en su sección con valor inicial y «restablecer», 20 claves de `app_setting` citadas; los 4 de `params_json` en `@param_sin_valor` y sembrados 7/7/1 en 00001; 0 mínimos de seguridad en tablas (**U4-12**); atajos 0 fuera de U1 §5.3 y comandos 0 fuera de §5.2 (**U4-06**); series A/B de 14 valores, mínimos 8,094.50 y 5,000.00, disponible de hoy 10,000.00 y 6,905.50 = oráculo; 800.50, 25.8685 %, 0.9836/0.4918/0, 9-nov-2026 y 8-ene-2027, $99.50 y 1,748.75 × 1.2 − 300 = 1,798.50 → 1,799 = S4-07 (**U4-05 caso A/B, U4-13**). 0 diferencias numéricas.
- Criterios: U4-01, U4-02, U4-03, U4-04, U4-06, U4-08, U4-11, U4-12, U4-13 cumplen; U4-07 cumple con hallazgo medio H-3; U4-05 no cumple (H-1, alta); U4-10 no cumple y U4-09 parcial (H-2, media). Detalle y evidencia por criterio en la tabla del informe.

## Supuestos
- Conté como «ajuste con efecto numérico» (U4-05) todo campo de Ajustes cuyo cambio mueve una cifra que la app muestra (disponible, calendario, aportación necesaria, días restantes). No conté el tiempo máximo del comparador ni las contraseñas (no mueven cifras).
- Acepté saturación y densidad como valores fijos (no configurables) porque U1 §2.3 aprobado dice «un solo control» y U4b lo declara en SUPUESTOS.

## Bloqueos
Ninguno nuevo. Abierto y respetado por U4: `docs/estado/bloqueos/S1-parametros-avisos.md`.

## Pendientes
Para la ronda 2 (autores U4a/U4b o integrador): H-1 a H-8 de `docs/estado/wf1/REVISION-U4-r1.md`.

## Solicitudes fuera de alcance
1. `docs/design/U1-arquitectura.md` §2.3 y §2.5.a: alinear las vistas de inventario (tablero por modo o por estado; galería), de líneas de crédito en P-26.05 (galería o tabla) y de `regla_reparto` (con o sin lista en P-26.12). Motivo: U1 se contradice y U4 elige sin poder citar una sola fuente (H-5).
2. `docs/contratos/openapi.yaml` (vía `build.py`), `changeOwnPassword`: fijar si cambiar la contraseña cierra las demás sesiones. Motivo: H-4.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (el revisor no propone valores operativos).

## Consumo aproximado
- Acciones de herramienta: 33.
- Leído: ~265 KB (protocolo 6, decisiones 26, criterios U4 27, U4-a 55, U4-b 59, U1 por partes ~34, extractos e índices ~12, migraciones ~8, reportes U4a/U4b 22, bloqueos 5, salidas de scripts ~11).
- Escrito: ~30 KB (informe 17.5, este reporte ~4.5, scripts temporales ~8, ya borrados).
