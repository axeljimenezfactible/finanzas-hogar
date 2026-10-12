# integrador-wf3c · Reporte · WF-3c · 12-oct-2026

Agente `integrador-wf3c` (Opus): integrador final de WF-3c, «Inicio como tablero de control (K54)». La integración completa está en `docs/estado/wf3/INTEGRACION-wf3c.md`.

No mandé Telegram, no usé sudo, no instalé paquetes y no toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`.

## Hecho

- **HECHO:** corrí `make verify-wf3`. Da **EXIT 2**, y la única causa es U5-04: 2 pares bajo umbral. Comprobé que es un falso positivo: es el `<svg>` raíz `aria-hidden` del barómetro, que no tiene `fill` ni `stroke`, y `contraste-dom.js:91-96` toma su `fill` inicial negro. El resto de `verify-wf3` pasa:
  - Go, con 33 paquetes;
  - cobertura K25;
  - contrato: 221 operaciones en 176 rutas, con 0 diferencias;
  - frontend: 259 passed y 4 skipped;
  - U5-01, U5-03, U5-08 y U5-11.
  - No lo puse en verde: el arreglo está fuera de mi alcance y no metí una excepción en el Makefile, porque escondería la medida. Va como PREGUNTA en `WF3c-ronda3-escalamiento`, punto 4.
- **HECHO:** `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` da **EXIT 0**.
- **HECHO:** `make verify-wf3-secretos verify-wf3-respaldo verify-wf3-servidor` da EXIT 0.
- **HECHO:** el servidor de desarrollo corre con el binario actual (16:36; ningún `.go` es más nuevo) y con el `dist` nuevo (`index-CBAvKt4u.js`, el mismo que sirve). Responde en **HTTP/2**:
  - `curl -w %{http_version}` = `2` por loopback;
  - `2` por la LAN, 192.168.50.56;
  - `1.1` con `--http1.1`.
- **HECHO:** revisé la captura `P-01-escritorio-claro.png` (1440 × 3029): el tablero está completo y las cifras cuadran con lo que compararon `u5-tablero` y C3-09. También comprobé que `P-01-detalle-android-claro.png` mide 475 px de ancho, que es el desborde de U5-06.
- **HECHO:** escribí `INTEGRACION-wf3c.md` con:
  - K54, U3 2.0, DESIGN 1.1, esquema 1.5 y C1 1.5 con sus conteos;
  - API, pantalla y capturas de P-01;
  - C3 v1.3, **no firmada**;
  - los hallazgos abiertos, consolidados con archivo y cambio exacto;
  - los bloqueos, las propuestas para la puerta 4 y el consumo de los 23 agentes.
- **HECHO:** abrí el bloqueo `docs/estado/bloqueos/WF3c-ronda3-escalamiento.md`. Junta lo no aprobado (M-R3-1, M-2 y C3 H-1 a H-5), el rojo de `verify-wf3` y la página de 475 px, con opciones y una PROPUESTA de ronda acotada.
- **HECHO:** a `docs/estado/bloqueos/INDICE.md` le agregué la sección «Bloqueos nuevos de WF-3c»: 14 archivos y 1 actualizado. Con eso se cierran B-3 (backend r3) y B-r1-4 (especificación r3).
- **HECHO:** en `README.md` actualicé el estado al 12-oct-2026, agregué una nota de versiones de WF-3c, la sección «Qué entregó WF-3c» y el punto 8 de «Cómo sigue».
- **HECHO:** hice el commit con la ruta que pidió el orquestador, sin `.db`, certificados, `node_modules`, `dist` ni `datos/tmp`. Antes lo comprobé con `--dry-run`: 352 rutas.

## Archivos

- `docs/estado/wf3/INTEGRACION-wf3c.md` (nuevo)
- `docs/estado/wf3/integrador-wf3c.md` (nuevo)
- `docs/estado/bloqueos/WF3c-ronda3-escalamiento.md` (nuevo)
- `docs/estado/bloqueos/INDICE.md` (sección agregada al final)
- `README.md` (estado, versiones, «Qué entregó WF-3c» y «Cómo sigue»)
- El `Makefile` no lo toqué, porque `verify-wf3` no necesita operaciones nuevas.

## Verificación

| Criterio | Comando | Salida resumida |
|---|---|---|
| S6-01 | `make verify-wf3` (verify-wf3-go) | go1.27.2; build, vet, gofmt y tidy limpios; 33 paquetes OK y 5 sin pruebas |
| S6-02 | verify-wf3-cobertura | K25 OK: 87.0, 89.2, 94.0, 93.2 y 88.8 |
| S6-03 | verify-wf3-contrato | 22 archivos idénticos; «221 operaciones en 176 rutas, 0 diferencias» |
| C3-03 y U5 (interfaz) | verify-wf3-frontend | build con 80 entradas precacheadas; 212 rutas invocadas, 0 fuera y 0 sin registrar; Playwright: 259 passed y 4 skipped |
| U5-01, U5-03, U5-08 y U5-11 | verify-u5 | 164 de 164; axe 0; hosts solo `finanzas.local:8443`; 41 filas de C3 |
| **U5-04** | verify-u5 | **«24466 pares evaluados, bajo umbral = 2»: EXIT 2** (falso positivo comprobado) |
| Global | `make verify-all GOOSE=…/goose` | **EXIT 0**: «Verificación ampliada OK» |
| S6-13 | `make verify-wf3-secretos` | 0 coincidencias |
| S6-09 | `make verify-wf3-respaldo` | `integrity_check` ok; 56 tablas, 2907 filas y 0 diferencias |
| S6-11, S6-12 y S6-14 | `make verify-wf3-servidor` | curl 200 con el SAN; manifest e íconos; `sw.js` 200; escucha solo en 127.0.0.1 y 192.168.50.56; Tailscale no responde |
| HTTP/2 | `curl --cacert … -w %{http_version}` | 2 (loopback), 2 (LAN) y 1.1 (forzado) |

## Supuestos

1. «C3 1.1 firmada o no» se refiere a los criterios C3 1.1 aplicados a la matriz v1.3. El revisor la devolvió y la matriz **no** tiene firma.
2. Hice el commit aunque `verify-wf3` da EXIT 2, porque:
   - la causa única es un falso positivo de medición documentado;
   - el commit es la tarea del orquestador;
   - WF-3b ya se integró con escalamientos abiertos.
   
   El orquestador decide si hace falta una ronda acotada antes de la puerta 4.
3. El `dist` que sirve el servidor lo generó `verify-wf3-frontend` (`npm run verify`) a las 17:25. El servidor lee `frontend/dist` del disco, así que no hizo falta reiniciarlo.

## Bloqueos

- Nuevo: `docs/estado/bloqueos/WF3c-ronda3-escalamiento.md`.
- Siguen abiertos los 13 de WF-3c y `C3-flujos`.

## Pendientes

- Poner `make verify-wf3` en verde. La opción A es `fill="none"` en el `<svg>` raíz de `frontend/src/componentes/Barometro.tsx:61`, y luego regenerar las capturas de P-01 y correr `make verify-u5`.
- Las salidas de M-R3-1, M-2 y C3 H-1 a H-5, más la página de 475 px de P-01 Detalle en Android (`WF3c-ronda3-escalamiento`).
- B-N4: el `Makefile` (líneas 8, 66 y 149) y `s1_verifica.py:190` todavía dicen «esquema 1.4». No lo cambié porque no es una necesidad de `verify-wf3`.
- La tabla «Versiones vigentes» del README sigue con el corte de la puerta 2. Puse una nota con los cambios de WF-3c, y conviene rehacerla en la puerta 4.

## Propuestas para Axel

- **PROPUESTA:** una ronda acotada «WF-3c bis» antes de la puerta 4 (detalle en `WF3c-ronda3-escalamiento`). Saldrían `verify-wf3` en verde y C3 firmable, con 3 agentes cortos. El riesgo es bajo.
- **PREGUNTA (puerta 4):** ¿la «deuda viva total» de K-3 son solo las compras hechas a la fecha? Así hoy serían $3,895.00, igual que el oráculo; con las compras programadas son $18,999.22.
- **PREGUNTA (puerta 4):** ¿dónde se editan los 5 parámetros del tablero? La propuesta es Ajustes › Presupuesto.

## Paquetes instalados

Ninguno.

## Consumo aproximado

Unas 32 acciones. Unos 85 KB de texto leídos, por secciones y con `grep`, más 1 captura vista como imagen. Unos 45 KB escritos. Borré los temporales de `datos/tmp/integ-wf3c/` y `datos/tmp/verify-wf3/`.
