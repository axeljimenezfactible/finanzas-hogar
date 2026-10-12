# Reporte · integrador-wf1 (WF-1, Opus) · 11-oct-2026

Cerrado por tope: no. Documento de integración: `docs/estado/wf1/INTEGRACION.md`.

## Hecho
- (1) `gofmt -l .` vacío, `go vet ./...` y `go build ./...` sin errores, `go test ./... -count=1 -cover` todo `ok`; cobertura por paquete en INTEGRACION §2.1 (K25 ≥ 85 % en motor_fin 89.2, motor_cred 93.7/96.0/90.8/90.0, recetas 93.2, precios 96.1/100/99.3/95.3; `dominio` aún no existe).
- (2) Oráculos: los 4 `*_run_fixtures.py` regeneran `oraculo/salidas/` idénticas; motor desde la raíz + `compara.py` s2 s3 s4 s5 con 0 diferencias; `fin_verify`, `pre_verify`, `model.py` y `verify3.py` reproducen; sha256 de `model.json` y `verify3.json` intactos.
- (3) `recetario.json` = unión validada de a + b (generador del oráculo `--verifica`, JSON Schema, 60 ids y nombres únicos, pruebas Go y fixture `s4-escalado-recetario`); parciales conservados; no hizo falta reescribirlo.
- (4) Tokens: `generar_c2_tokens.py` idéntico, `verificar_c2.py` OK, `sincronia.mjs` reproduce `tokens.css` y `tailwind.theme.css`; `contraste_componentes.py --check` IDENTICO, 0 bajo umbral.
- (5) `docs/design/U4-ajustes.md` unido (encabezado + índice + parte a + parte b sin cambios; 14 secciones P-26); parciales conservados.
- (5b) `paquete_entrada.sh` regenerado: solo cambió `U4-ajustes.pdf` (ahora del archivo unido) y su fila de `MANIFIESTO.md`; `sha256sum -c` 0 diferencias; sin código en `docs/design/entrada/`.
- (6) Makefile: `verify-wf1` (7 subobjetivos) y `verify-all` lo incluye; `make verify-all GOOSE=…/gobin/goose` rc 0.
- (7) `docs/estado/wf1/INTEGRACION.md` (secciones 0 a 11, incluidas «Interrupción, continuación y rondas extra», escalados a Axel y «Paquete para Claude Design»); `docs/estado/bloqueos/INDICE.md` con la sección «Bloqueos nuevos de WF-1» y las propuestas K31.
- (8) `README.md`: estado al cierre de WF-1, árbol, versiones, «Qué entregó WF-1» y «Cómo sigue»; fila de criterios corregida a aprobada en P1d.
- (9) `.gitignore`: `datos/semilla/*/raw/*` (conserva sus `.gitignore`), `profeco/db/*.sqlite*`, `**/node_modules/`, `.obsidian/` y la excepción `!datos/semilla/nutricion/tabla.csv`.
- (10) Commit único de WF-1 (hash y tamaño en la salida del integrador).

## Archivos
- Creados: `docs/design/U4-ajustes.md`, `docs/estado/wf1/INTEGRACION.md`, `docs/estado/wf1/integrador-wf1.md`.
- Modificados: `Makefile`, `README.md`, `.gitignore`, `docs/estado/bloqueos/INDICE.md`, `docs/design/entrada/U4-ajustes.pdf` y `docs/design/entrada/MANIFIESTO.md` (regenerados por el script).
- Sin cambios: `backend/cmd/motores/main.go` (los 4 subcomandos ya estaban registrados), `datos/semilla/recetario.json`, `docs/criterios/INDICE.md` (ningún agente cambió criterios) y todos los artefactos revisados.

## Verificación
- Ver INTEGRACION §2 (comandos y salidas). Resumen: go test todo ok; compara s2 2320/2320, s3 906/906, s4 8451/8451, s5 747/747, diferencias 0; `model.json` 596/596; recetario 0 errores; tokens y contraste idénticos; U4 unido True; paquete `sha256sum -c` OK; `make verify-all` → «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».
- Criterios: el integrador no tiene fase propia; los criterios de S2 a S5, U2 a U4 y P3 los evaluaron los revisores (veredictos en INTEGRACION §1 y §4) y el integrador volvió a correr su evidencia mecánica (oráculos, cobertura K25, tokens, contraste, paquete).

## Supuestos
- La lectura obligatoria (3) «criterios de tu fase» no aplica al integrador; no leí completos los 18 archivos de criterios.
- `recetario.json` ya era la unión que escribe `oraculo/rec_genera_recetario.py` y coincide con lo regenerado; no lo reescribí para no tocar contenido revisado.
- `.obsidian/` es configuración local de Obsidian y no es parte del proyecto.
- `datos/semilla/nutricion/tabla.csv` debe versionarse (la leen el recetario y las fixtures de S4).
- El commit lo hace el integrador por el protocolo §3 (un commit por flujo), en el repositorio local sin remoto.

## Bloqueos
- Ninguno nuevo del integrador. Abiertos: los 26 de WF-1 y `S1-parametros-avisos` (INTEGRACION §6).

## Pendientes
- Puerta 2: 6 artefactos sin aprobar (INTEGRACION §4), 27 bloqueos y valores K31 (§7). Resto en INTEGRACION §8.

## Solicitudes fuera de alcance
- `backend/go.mod`: `go mod tidy` para quitar `// indirect` de `modernc.org/sqlite` (I-2).
- `backend/cmd/motores/main.go` o `oraculo/fixtures/FORMATO.md` §4.1: resolver o documentar que las rutas de las fixtures son relativas a la raíz (I-1).
- `oraculo/fixtures/FORMATO.md` §6: cita `oraculo/cred_compara.py`, que no existe (I-4).
- `docs/design/entrada/LEEME.md` paso 2 (H-1 de `REVISION-paquete-r3.md`): se escala, no se corrigió.

## Paquetes instalados
- Ninguno. `jsonschema` efímero con `uv run --offline --with`.

## Propuestas para Axel
- Consolidadas en INTEGRACION §7 y en `docs/estado/bloqueos/INDICE.md`.

## Consumo aproximado
- ≈ 52 acciones; ≈ 175 KB leídos; ≈ 60 KB escritos con herramientas (más 143 KB de `U4-ajustes.md` generado por concatenación y los PDF del script).
