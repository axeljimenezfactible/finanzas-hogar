# P2d · revisor-remanentes-r1 · reporte (12-oct-2026, madrugada)

Agente `revisor-remanentes-r1` (Opus), flujo «P2d · K53». Revisión independiente del trabajo de `p2c-remanentes`; no corregí nada. El informe es `docs/estado/p2c/REVISION-remanentes-r1.md`. Sin `git add` ni commit, sin Telegram ni navegador.

## Hecho
- Revisé `p2c-remanentes.md` contra:
  - M-1 de `REVISION-contratos-r3`;
  - M-1, B-2, B-5 y B-6 de `REVISION-motor-S5-r3`;
  - I-1, I-3, I-4 e I-5 de `integrador-p2c` (INTEGRACION §5);
  - los criterios S5 y S1 (S5-10, S5-13, K25).
- **Veredicto: APROBADO, con 0 altos, 0 medios y 2 bajos.**
  - B-1: marcas viejas en `oraculo/rec_*.py` y `oraculo/README.md`. Están declaradas y en la lista para WF-3.
  - B-2: una frase de alcance inexacta en el reporte.
- Escribí 32 sondas propias de `created_at`, `observed_date` y `today`. Incluyen los 6 textos de M-1 y controles válidos. Motor y oráculo coinciden en las 32.

## Verificación (comandos y salidas)
| Comando | Salida |
|---|---|
| `python3 -I /tmp/rev-rem-r1/estados.py docs/estado/bloqueos` | `85 {'RESUELTO': 64, 'ABIERTO': 19, 'SIN-LINEA(sin)': 2}`; los 2 son de P3, así que hay 21 abiertos. El Total del INDICE da 85, 64 y 21 |
| Bucle «Respuesta de Axel» | `total=85 resp=67`, es decir, 18 sin ella (regla de las herramientas) |
| `grep` de los 13 nombres en INDICE y en `decisiones.md` l. 120 | todos al menos 1 y 1 |
| `sondas.py`, `pre_run_fixtures.py --fixtures datos/tmp/rev-rem-r1/fx` y `motores s5 dump` | `iguales 32 de 32 json idéntico: True` |
| Motor y oráculo de `476f4a9` (`git archive`) con las 7 fixtures nuevas | oráculo viejo: «7 fixtures, 5 coinciden … 2 diferencias»; `compara.py`: 7 diferencias en total |
| `pre_genera_fixtures.py` en una raíz aislada y `diff -r` | «43 fixtures»; vacío |
| `pre_run_fixtures.py --out datos/tmp/rev-rem-r1/s5-or.json` y `cmp` | «43 fixtures, 43 coinciden, 1366 hojas de esperado, 0 diferencias»; idéntico |
| `compara.py s5 --oraculo … --motor …` | 963 de 963; provisionales 317 de 317; diferencias: 0 |
| `go test -cover ./internal/precios/... ./internal/recetas/...` | 96.0, 100.0, 99.3, 95.3 y 93.2 %; total 90.4 %; las 3 pruebas de fechas PASS |
| `bash /tmp/rev-rem-r1/paq.sh` | 5 sellos iguales a la fuente actual; 21 OK; `/api/v1` = 0 y nombres de operación = 0 en los 4 PDF |
| `ventana.py` sobre los archivos de `p2c-remanentes` | 56 renglones PROVISIONAL; solo DESIGN l. 933, que dice FIRME y cita el bloqueo abierto en lo PROVISIONAL |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | `EXIT=0`; RESUMEN TODO OK; 178 y 0; `git status` sin cambios por la corrida |
| `git status --short --untracked-files=all \| grep -E 'design/(sistema\|maquetas)\|wfp3'` | 0 |

## Supuestos
1. **5 rutas son de `k53-aplicar`.** `README.md`, `docs/criterios/{C1,INDICE,S5}.md` y `docs/plan/plan.md` se las atribuyo por su mtime, posterior al reporte de `p2c-remanentes`, y porque el reporte dice que no las tocó. No las evalué.
2. **I-5 se mide en los sitios que lista.** Lo medí contra `integrador-p2c` e INTEGRACION §5. Lo que queda en `oraculo/rec_*.py` lo dejé como bajo: I-5 lo mandaba a WF-3, y está en la lista para WF-3.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- B-1 y B-2 (ver el informe), para el dueño de S4 o WF-3, y para `p2c-remanentes` si hay ronda.
- Que el integrador vuelva a medir el INDICE cuando cierre `k53-aplicar`.

## Propuestas para Axel
- Ninguna de negocio.

## Consumo aproximado
- Unas 45 acciones de herramienta.
- Unos 120 KB leídos, por partes con `grep`, `sed` y `awk`. No leí `openapi.yaml`, `DESIGN.md` ni el INDICE completos.
- Unos 40 KB escritos: el informe, 13 KB; este reporte, 4 KB; los scripts de `/tmp/rev-rem-r1`, unos 12 KB; las sondas de `datos/tmp/rev-rem-r1`, unos 12 KB.
