# P2c · revisor-S1-r3 · reporte

Agente `revisor-S1-r3`, revisor independiente, ronda 3 de 3. Flujo P2c «Aplicación de la puerta 2», 11-oct-2026, noche. No corregí nada y no hice git add ni commit.
- Informe: `docs/estado/p2c/REVISION-S1-r3.md`.

## Hecho
- **Veredicto: NO APROBADO, con 1 alto, 0 medios y 7 bajos.**
  - El alto es H-01, que sigue abierto. `make verify-c1`, y con él `make verify-all`, fallan con `gen2: tipo de parámetro desconocido: hora HH:MM`.
  - Si se salta ese error, `verify_c1.py` da FALLA K44 (10 claves) y FALLA C1-15 (6 marcas en `openapi.yaml` y líneas de `C1-datos.md` que citan `S1-parametros-avisos`, ya resuelto).
  - El arreglo está en `backend/internal/api/c1gen/`, fuera del ALCANCE de S1, y nadie lo ha tomado.
- **Cerrados:**
  - N-01: la fila 5 del mapa de Ajustes ya trae `.vat_included`.
  - N-02: el par (0, NULL) dice lo mismo en las 5 fuentes de S1, como PROVISIONAL con la cita del bloqueo abierto `S3-iva-lectura`, pregunta 2, que no tiene recomendación.
- **Cumplen:**
  - las verificaciones (a) a (f) del encargo;
  - las 10 claves, con su valor igual a PUERTA-2 §4 y a la decisión F;
  - el traspaso y el Down de `params_json`;
  - `vat_included` según la opción 1 de U4-iva-credito, sin tocar `vat_on_interest_bp`;
  - 00001 a 00003 intactas;
  - el diccionario, idéntico al que genera `gen_dict.py`.
- **Bajos:**
  - siguen B-01 a B-06; B-06 se amplía porque el INDICE tampoco lista `S3-iva-lectura`;
  - B-07 es nuevo: huecos de `s1_verifica.py` en S1-16, K47 y S1-13, que hoy cubren mis scripts.

## Verificación
Detalle completo, con salidas, en la tabla de `REVISION-S1-r3.md`. Resumen:

| Comando | Resultado |
|---|---|
| `make verify-sql GOOSE=…` | rc=0; 55 tablas, 152 índices y 79 disparadores; integrity ok; el Down deja 0 objetos; goose up, down-to 0 y up |
| `make verify-s1 GOOSE=…` y `s1_verifica.py` directo | rc=0, `RESUMEN TODO OK` |
| goose con datos (`up-to 3`, 2 tarjetas, `up`, `down-to 3`, `up`) | rc 0/0/0/0; (0, 1600) y (1, NULL) con su `updated_at`; Down de vuelta a 1.3 |
| `python3 -I datos/tmp/rev-s1-r3/chk.py <repo>` | 15 de 15 OK: B-param, B-otras, B-desc, B-rangos (70 casos), C-up, C-semilla, C-down, C-down-valores, C-up-valores, D-traspaso, D-cols, D-ddl, D-check (26 casos) y D-down. `RESUMEN TODO OK` |
| `python3 -I datos/tmp/rev-s1-r3/chk2.py <repo>` | las líneas PROVISIONAL de 1.4 citan bloqueos ABIERTOS; las 10 claves están en su fila del mapa. Un falso positivo de mi lector (`attachment.ocr_draft_*`), revisado a mano |
| `git diff HEAD -- 00001… 00002… 00003… \| wc -l` | 0 |
| `gen_dict.py` a `datos/tmp` (dos veces) y `cmp` | idéntico al repositorio e idempotente |
| `make verify-c1 GOOSE=…` | rc=2, `gen2: tipo de parámetro desconocido: hora HH:MM` |
| `uv run … verify_c1.py --db` (base con las 4 Up) | rc=1, FALLA K44 y FALLA C1-15 |
| `make -k verify-all GOOSE=…` | rc=2; solo falla verify-c1 |

- En mis dos primeras corridas, C-up y C-down fallaron por error mío: comparaba `updated_at` entre dos bases creadas en instantes distintos.
  - Lo corregí para comparar antes y después sobre la misma base, y así pasan.
  - Mi primera versión de C-semilla también confundía la lectura de `params_json` en el INSERT con una escritura.
- Borré `datos/tmp/rev-s1-r3/` al terminar.

## Supuestos
1. H-01 se queda en alta aunque su arreglo no es de S1.
   - El encargo pone «verificación global rota» como alta, y Q2-97 A y K52 piden cerrar con `make verify-all`.
   - Por ser la ronda 3 de 3, lo escalo al orquestador en lugar de pedir otra ronda de S1.
2. N-02 está bien resuelto como PROVISIONAL.
   - `S3-iva-lectura` pregunta 2 no tiene recomendación. K46 deja esas preguntas con la opción en uso, y la opción en uso es la A: la usa S3 1.5 y el esquema admite el par.
   - Escribirlo FIRME por K18, como sugería la r2 en primer lugar, habría decidido por Axel.
3. Que el INDICE no liste dos bloqueos (B-06) lo dejo en baja, igual que r2. S1-16 pide citar un archivo abierto de `docs/estado/bloqueos/`, y eso se cumple.

## Bloqueos
- No abrí ninguno.

## Pendientes
- H-01: hace falta un agente de C1. Ver la corrección sugerida en el informe y la solicitud 1.
- Después, `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` con rc=0, antes del commit «[Puerta 2]».
- Bajos B-01 a B-07.

## Solicitudes fuera de alcance
1. **Orquestador → agente de C1:**
   - tipo «hora HH:MM» en `_pv`, como string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`;
   - las 10 claves con lectura, cambio y restablecer;
   - quitar las marcas que citan `S1-parametros-avisos`;
   - agregar `credit.vat_included` y `vat` en `Credit.missing_terms`, con (0, NULL) PROVISIONAL;
   - regenerar.
2. **Integrador o agente de bloqueos:** dar de alta `S1-reparto-edicion-gasto` y `S3-iva-lectura` en INDICE y en K46.
3. **Dueño de `s1_verifica.py`:** solicitud 3 de `correccion-S1-r2.md`.
4. **Integrador:** ratificar B-04.

## Propuestas para Axel
- Ninguna nueva. Sigue abierta la pregunta 2 de `S3-iva-lectura`: si marcas «no incluye IVA» y no escribes el porcentaje, ¿la tarjeta dice «falta capturar» (A) o la app no te deja guardar (B)? Mientras tanto se usa A, como PROVISIONAL.

## Consumo aproximado
- **Acciones de herramienta:** ~40.
- **KB leídos:** ~150.
  - Revisión r2 y correcciones: ~31.
  - 00004 completa: ~24.
  - Fragmentos del diccionario, de las reglas, de PUERTA-2, de decisiones, de criterios, de los bloqueos, de S3 y del Makefile: ~45.
  - Salidas de verificación: ~50.
- **KB escritos:** ~37.
  - `REVISION-S1-r3.md`: ~17.
  - Este reporte: ~5.
  - Scripts temporales, ya borrados: ~15.
