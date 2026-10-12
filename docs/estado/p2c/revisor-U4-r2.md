# revisor-U4-r2 · P2c · reporte del revisor independiente (ronda 2 de 3)

Agente `revisor-U4-r2` (Opus) · flujo P2c · 11-oct-2026, noche · árbol de trabajo sobre HEAD 4079cee. El informe completo está en `docs/estado/p2c/REVISION-U4-r2.md`.

## Hecho
- Revisé U4 a 1.6, b 1.5, el unido, U3 1.4 y el diff de DESIGN.md. Usé `docs/criterios/U4.md`, `docs/criterios/U3.md`, PUERTA-2 §2.4, §3 (Q2-11 y Q2-66 a Q2-70) y §5, el literal de los bloqueos que cita cada cambio, C1 1.3 (extractos) y los contratos de S2, S3 y S4.
- Volví a ejecutar con mis propios scripts las verificaciones (a) a (e), la forma, las cifras, la lista de PROVISIONAL contra el estado de cada bloqueo y `make verify-all`.
- **Veredicto: NO APROBADO.** 0 altos, 2 medios y 3 bajos:
  - M-1: D-2 de U3 da por FIRME el ahorro planificado, pero `S2-ahorro-quincenal` está abierto;
  - M-2: D-6 de U3 deja el avance como entero y lo da por resuelto con `S2-metas`, aunque Q2-17 C solo deja entero en la tarjeta;
  - B-1 y B-2 siguen de la ronda 1;
  - B-3: la sección de `ocr.timeout_seconds` es una PROPUESTA sin pregunta registrada.
- Los 5 medios de la ronda 1 están cerrados. En M-4 se cerró lo señalado, pero de su corrección salen M-1 y M-2.

## Verificación
Scripts en `datos/tmp/rev-u4-r2/` (`verif.py` y `prov.py`, con `python3 -I`), borrados al terminar.
- `verif.py`:
  - `openapi version 1.3.0 | ops openapi 212 | ops extractos 212 | iguales True` · `Ajustes ops 75`;
  - `U == head+a+"\n"+b: True` (encabezado de 38 líneas) · `secciones P-26 en U: 14 14 a: 7 b: 7`;
  - `candidatas sin existir: []` · `U4-02 Ajustes sin cita: []`;
  - las 8 claves `reports.*` y `alerts.*`, más `optimizer.date_window_days` y `ocr.timeout_seconds`, están en el Up de la 00004 con 21:30, 7, 19:00, 1, 08:00, 1, 7, 09:00, 7 y 30;
  - `claves citadas en U4 que no aparecen en ninguna migración: []`;
  - forma: a, b y U con 0 hex y 0 ISO.
- `prov.py`: cada PROVISIONAL de U4 cita un bloqueo abierto o una solicitud. En U3, la l. 186 no cita `S2-ahorro-quincenal`, que tiene «Estado: ABIERTO».
- DESIGN: `git diff -U0` da 4 líneas. Los números son iguales en las líneas − y + (solo se agrega el «2» de «puerta 2»), con 0 hex. `git diff 1ccb468 HEAD -- docs/design/DESIGN.md` está vacío.
- Cifras, con `fractions` y HALF_UP entero:
  - avance: 500 décimas;
  - aportación: 104,891, la última 104,890, y 157,336;
  - inventario: 1,486,571 milésimas, 32 días y 315,333 milésimas por semana;
  - tarjeta: 31-oct-2026, 20-nov-2026 y 30-nov-2026.
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: **rc=0** en 21.8 s:
  - U4: «U4-ajustes.md termina con parte a + parte b sin cambios: True · secciones P-26: 14 de 14»;
  - C1: «los 21 archivos coinciden byte a byte con el generador»;
  - S1: «RESUMEN TODO OK»;
  - oráculos S2 a S5 con 0 diferencias;
  - «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».

## Supuestos
1. **Q2-17 C («un decimal en Metas y entero en la tarjeta de Inicio»).** El entero vale solo para la tarjeta 6, como lo escriben C1 1.3 (`progress_tenths_pct`) y `S2-interfaz.md` 1.3 l. 295. Por eso D-6 es M-2. Como un agente podría leer que D-6 queda fuera de «Metas», no lo califiqué como alto, y la corrección admite marcarlo PROVISIONAL con un bloqueo.
2. **Ejemplo ilustrativo de fechas (a l. 295).** Corte 31 y pago 20 llevan la marca «no canon» y no son cifras de dinero. No lo cuento contra U4-11, igual que los ejemplos ilustrativos de RN-17 que aceptaron las rondas anteriores.
3. **Línea 6 de U3.** Es el registro histórico «Cambios de la versión 1.2», no una marca vigente.

## Bloqueos
No abrí ninguno: un revisor no corrige. Para M-2, si el corrector no aplica el decimal, el bloqueo sugerido es `docs/estado/bloqueos/U3-avance-detalle.md`.

## Pendientes
- Ronda 3 de 3: corregir M-1 y M-2 en `U3-dashboard.md`. Los bajos son opcionales.
- Al integrador:
  - registrar en `docs/estado/bloqueos/INDICE.md` `U4-reparto-preseleccion`, `C1-hora-reportes-job` y `C1-reparto-alcance`;
  - pedir la línea «Estado: ABIERTO» en los dos bloqueos de C1;
  - llevar al revisor de C1 el enum de `SavingsGoal.status`: C1 1.3 sigue con `on_track`, `behind` y `achieved`, contra Q2-18 A.

## Propuestas para Axel
Ninguna nueva. La pregunta de la puerta 4 sigue abierta en `U4-reparto-preseleccion` (A, B o C).

## Consumo aproximado
- Unas 38 acciones de herramienta.
- Unos 210 KB leídos: diffs por palabra de a, b y U3 (unos 86 KB), diff de DESIGN, reportes, criterios por secciones, PUERTA-2 por secciones, bloqueos y extractos con `grep`, `sed -n` y Python.
- Unos 19 KB escritos: este reporte, el informe y unos 5 KB de scripts temporales, ya borrados.
