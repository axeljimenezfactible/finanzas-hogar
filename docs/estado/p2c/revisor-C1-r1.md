# P2c · revisor-C1-r1 · reporte

Revisor independiente de C1 v1.3 (puerta 2), ronda 1 de 3. Flujo P2c, 11-oct-2026, noche. El informe completo está en `docs/estado/p2c/REVISION-C1-r1.md`.

## Hecho
- **Veredicto: NO APROBADO**, con 0 altos, 3 medios y 5 bajos.
- **Medios:**
  - H-1: `visit_date_` y `cheapest_store_id_` salieron con un guion bajo de más (`gen2_schemas.py` l. 407 y 420).
  - H-2: M-11 de C1-datos no nombra los acumulados de Q2-65 A.
  - H-3: no está escrito qué responde la API a `vat_included` = true con porcentaje.
- **Bajos:**
  - B-1: marca de más en `vat_included`.
  - B-2: vocabulario 1/0/NULL en lugar de true/false/ausente.
  - B-3: `S3-iva-lectura` y RN-26 no registran lo que tomó C1.
  - B-4: `docs/criterios/C1.md` sigue en 1.2.
  - B-5: `bloqueos/INDICE.md` no lista los dos bloqueos nuevos de C1.
- Reejecuté todas las verificaciones con mis propios comandos y scripts en `datos/tmp/rc1r1/`, y borré esa carpeta al final.
- No toqué ningún archivo fuera de estos dos reportes. No hice `git add` ni `commit`.

## Verificación
- `make verify-c1` → rc 0, 21 archivos «IDÉNTICO», «comprobaciones OK: 29 · RESULTADO: TODO OK».
- `make verify-openapi` → «docs/contratos/openapi.yaml: OK».
- `goose … sqlite3 datos/tmp/rc1r1/c1.db up` → migra a la versión 4. Después, `verify_c1.py --db` sobre esa base → rc 0, 29 OK, 0 FALLA.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP_ROOT=/home/deepsec/finanzas-hogar/datos/tmp/rc1r1/tr` → rc 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».
- Scripts propios con PyYAML (`chk.py` a `chk9.py`):
  - conteos de HEAD contra el árbol de trabajo: de 167, 210 y 252 a 168, 212 y 255;
  - las 10 claves con su rango;
  - `params_json` sin propiedades;
  - campos de la puerta 2;
  - las 6 marcas;
  - extractos 212/212.
- Consultas `sqlite3` a la base que migré: valores sembrados y `notification.params_json`.
- Lo encontré así:
  - `grep -n 'visit_date_:\|cheapest_store_id_:'` → `openapi.yaml` l. 13517 y 13718, `compras-tiendas-a.yaml` l. 797 y 998, y `reportes.yaml` l. 802;
  - `grep -n 'C1_VERSION *='` → `gen4_build.py:6: C1_VERSION = '1.3.0'`.

## Supuestos
1. Para Q2-72 a Q2-77 usé PUERTA-2 §4 como respuesta de Axel. La orden nombra §2 y §3, pero §3.1 remite a esas preguntas y `RESPUESTAS.md` las da por aprobadas.
2. Leí `openapi.yaml` solo con scripts que imprimen resúmenes y `grep` de líneas sueltas. Nunca entero.
3. Califiqué H-3 como medio porque cae en «algo que el implementador siguiente tendría que adivinar». Es una regla técnica, no una decisión de Axel.

## Bloqueos
- No abrí ninguno.
- Los que cita C1 están abiertos (sin «Respuesta de Axel»): `C1-reparto-alcance`, `C1-hora-reportes-job`, `S3-iva-lectura` y `S1-reparto-edicion-gasto`.

## Pendientes
- **Ronda 2:** comprobar la corrección de H-1, H-2 y H-3 y, si se atienden, de B-1 a B-3.
- **Integrador:** B-4 (criterios C1 1.3) y B-5 (`INDICE.md` de bloqueos).
- **P2d (K53):** C1 todavía tiene la lectura asistida y le falta la importación del agente externo.

## Propuestas para Axel
- **Ninguna de fondo:** los tres medios son técnicos y no piden decisión.
- **PROPUESTA de proceso:**
  - Problema: `make verify-all` usa una carpeta temporal compartida, `datos/tmp/verify-wf0`. Dos agentes que lo corren a la vez chocan: mi primer intento dio «disk I/O error» porque el revisor `rev-s5-p2c-r1` lo estaba corriendo al mismo tiempo, y su corrida terminó bien.
  - Propuesta: cada agente lo corre con `TMP_ROOT=<su carpeta en datos/tmp>`, o el Makefile toma una carpeta por proceso.

## Solicitudes fuera de alcance
1. **Implementador de C1 (ronda de corrección):** H-1, H-2 y H-3 en `backend/internal/api/c1gen/`. Después, regenerar con `build.py`.
2. **Integrador:**
   - `docs/criterios/C1.md` a 1.3 (B-4);
   - `docs/estado/bloqueos/INDICE.md` con los dos bloqueos de C1 (B-5);
   - una línea en «Cómo seguí» de `S3-iva-lectura` y en RN-26 sobre la lectura de C1 v1.3 (B-3).

## Consumo aproximado
- Acciones de herramienta: unas 40.
- Leído: unos 200 KB, por partes. Incluye:
  - criterios C1 (~38 KB);
  - PUERTA-2 §0, §2, §3 y §4 (~40 KB);
  - reportes C1, S1 y contratos (~25 KB);
  - fragmentos de C1-datos y de los bloqueos (~35 KB);
  - salidas de scripts y logs (~60 KB).
- Escrito: unos 22 KB entre los dos reportes y los scripts temporales, que ya borré.
