# P2c · integrador-p2c · reporte (11-oct-2026, noche)

Agente `integrador-p2c` (Opus), flujo P2c «Aplicación de la puerta 2». El documento de integración es `docs/estado/p2c/INTEGRACION.md`. Aquí van solo mi trabajo y mis comprobaciones.

## Hecho
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: pasó completo antes y después de mis cambios (exit 0 las dos veces).
- `bash docs/design/tools/paquete_entrada.sh`: **falla** con `ERROR: U4-ajustes.pdf trae 2 líneas con /api/v1`.
  - Causa: U4 1.6, l. 490 y 848 de `U4-ajustes.md`, cita `/api/v1/settings/reports` sin método HTTP, y el filtro del script no quita esa forma.
  - El script dejó `docs/design/entrada/` a medias: 4 archivos nuevos y el `MANIFIESTO.md` viejo.
  - Restauré `docs/design/entrada/` y `docs/design/tools/paquete/.cache/` con `git checkout`. No edité nada a mano.
  - Resultado: `sha256sum -c`, 21 OK.
  - No lo arreglé: U4 está aprobado y el script no es de mi ALCANCE.
- Makefile: 4 rótulos de `verify-s1` pasan de «esquema 1.3» a «esquema 1.4», y uno de «hoy 00001 a 00003» a «00001 a 00004». La verificación no cambió, porque el comodín ya tomaba la 00004. Es la solicitud 7 de S1, cosmética.
- `docs/estado/bloqueos/INDICE.md`:
  - Conteo final: 84 archivos, 63 resueltos y 21 abiertos.
  - Fila P2c: 13 abiertos. Total: 84, 63 y 21.
  - Subsección «Dados de alta al cierre de P2c» con los 9 bloqueos que faltaban. Cierra B-06 de S1, B-4 de motor S5 y M-1 de contratos en lo que toca al INDICE.
  - No cambié el contenido de ningún bloqueo ni decidí ninguno.
- `docs/estado/p2c/INTEGRACION.md` (§0 a §11) y README con «Estado al 11-oct-2026, puerta 2», la tabla de versiones vigentes y el punto 4 de «Cómo sigue».
- Registré la excepción de independencia de §15.3 del motor S5 (B-1) y ratifiqué la edición de herramientas de S1 (B-04), como pedían las revisiones.
- Commit único, con las exclusiones de §10 de INTEGRACION.

## Verificación (comandos y salidas)
| Comando | Salida |
|---|---|
| `make verify-all GOOSE=…` (1.ª, antes de mis cambios) | exit 0, 21.7 s; `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` |
| `bash docs/design/tools/paquete_entrada.sh` | exit 1; `ERROR: U4-ajustes.pdf trae 2 líneas con /api/v1` |
| `pdftotext docs/design/entrada/U4-ajustes.pdf - \| grep -n /api/v1` (PDF a medias) | 2 líneas: `( /api/v1/settings/reports , N-296 y N-297)` y `( /api/v1/settings/reports );` |
| `git checkout -- docs/design/entrada docs/design/tools/paquete/.cache` y `sha256sum -c` | 21 OK; `git status` de esas rutas, vacío |
| script propio de estados de bloqueo (`/tmp/p2c-integ/estados.py`) | `total 84 Counter({'RESUELTO': 63, 'ABIERTO': 16, 'ABIERTO?': 5})`; los 5 «ABIERTO?» son los 2 de la Pista 3, `C1-hora-reportes-job`, `C1-reparto-alcance` y `S5-precio-vigente-captura`, abiertos por su texto |
| `grep -n -A2 on_track docs/contratos/extractos/presupuesto-b.yaml` | `SavingsGoal.status` enum `on_track`, `behind`, `achieved` (I-2) |
| `make verify-all GOOSE=…` (2.ª, después de Makefile, INDICE, README e INTEGRACION) | exit 0, 21.7 s; `esquema 1.4: todas las migraciones`; `RESUMEN TODO OK`; `TOTAL comprobadas: 178 \| con diferencias: 0`; `sha256sum -c: 21 archivos sin diferencias`; `Verificación ampliada OK` |
| `git show --name-only --format= 0271f3b` filtrado por rutas fuera de la Pista 3 | vacío: [P3c] no tocó nada de P2c |
| `git diff --cached --name-only \| wc -l` antes del `git add` | 0 |

## Supuestos
1. **docs/design/salida y tools/ds, fuera del commit.** Excluí también `docs/design/salida` y `docs/design/tools/ds`, que no estaban en la lista de exclusión del encargo. Son de la Pista 3: lo confirma su propio commit [P3c], y siguen cambiando en el árbol. Meterlos en «[Puerta 2]» atribuiría a P2c trabajo a medias.
2. **El alto común de las r3 está cerrado de hecho.** El alto de bloqueos, S1 y contratos (`verify-c1` roto) lo doy por cerrado porque `verify-all` pasa completo. No hubo ronda 4 que lo confirme, así que esos artefactos siguen «no aprobados» y se escalan.
3. **Restaurar el paquete no es editarlo.** Restaurar `docs/design/entrada/` con `git checkout`, después de que el script fallara, solo deshace mi propia corrida. No lo cuento como edición a mano.
4. **S5-terminos sigue abierto.** Queda ABIERTO en el INDICE. K53 se registra en P2d; no me toca a mí cerrarlo.

## Bloqueos
- Ninguno nuevo. I-1 (paquete) e I-2 (C1 `SavingsGoal.status`) son defectos técnicos con dueño, no decisiones de Axel. Van en INTEGRACION §5 y §9.

## Pendientes
- I-1: regenerar el paquete cuando el dueño de U4 o del script arregle las 2 rutas. Después, `bash docs/design/tools/paquete_entrada.sh` y `make verify-all`.
- I-2, I-3, I-4, I-5 e I-6 (INTEGRACION §5). Motor S5, M-1. Registrar en `decisiones.md` §K46 los 13 bloqueos de P2c. K53 en P2d.
- I-7: ¿de dónde sale «28 sin recomendación» en `RESPUESTAS.md`? La hoja da 16 filas más 3 parciales, y K46 cuenta 15.

## Propuestas para Axel
- Ninguna propia. Las de los agentes están consolidadas en INTEGRACION §7.

## Consumo aproximado
- Acciones de herramienta: unas 52.
- KB leídos: unos 205. Todo por partes con `sed -n`, `grep` y `awk`. No leí `openapi.yaml`, `DESIGN.md` ni `PUERTA-2.md` enteros: de la hoja §8 solo saqué las filas con `awk`.
- KB escritos: unos 65. INTEGRACION 29, INDICE 13, este reporte 7, README 5, Makefile menos de 1, y scripts temporales en `/tmp/p2c-integ`, unos 10.
