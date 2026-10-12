# revisor-especificacion-r2 · Reporte del revisor (WF-3c, K54, ronda 2 de 3) · 12-oct-2026

## Hecho
- Revisé, de forma independiente, a `u3-v2` y `s1-c1-v15` después de `correccion-especificacion-r1`. Medí contra K54, U3.md, U2.md, C1.md y S1.md.
- Escribí el informe en `/home/deepsec/finanzas-hogar/docs/estado/wf3/REVISION-especificacion-r2.md`.
- **Resultado: NO APROBADO.**
  - 1 alto: A-1 = H-1 reabierto, `make verify-all` da EXIT 2 porque las 5 operaciones nuevas responden 501.
  - 1 medio nuevo: M-N1, la «M-01 ampliada» de K-1 no está en C1-datos §9.
  - 5 bajos nuevos.
- Cerrados: H-2, M-1, M-2, M-3 y M-4 de la ronda 1.
- No corregí nada. No hice `git add`, no envié Telegram, no usé sudo ni instalé paquetes.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/REVISION-especificacion-r2.md`, el informe.
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/revisor-especificacion-r2.md`, este reporte.

## Verificación
- **K54 (a):** las citas de Axel son idénticas a las de `k54-tablero.md` (script). Están las 4 decisiones y §G se conserva.
- **U3-01:** K = 4, B = 1, G = 6, Tarjeta = 4; §2.2, §2.3 y §2.4 están.
- **U3-02 y C1-09:**
  - 24 operaciones citadas, 0 faltan.
  - Campos de §2.6 iguales a los extractos.
  - M-N1 y B-N3 abiertos.
- **U3-03:** 25 bloques con sus 4 estados.
- **U3-08:** 23 montos, 21 directos en `model.json` y `data.json`. Los otros 2 ya tenían su origen documentado en 1.x. El 27 % lleva la marca SINTÉTICO.
- **U3-09 y U3-11:** 0 fechas ISO y 0 colores en hex o `rgb(`.
- **U2-04:**
  - C-159 a C-172 recalculados con mi propio código desde `C2-tokens.json`: coinciden exactamente en los dos temas.
  - 40 pares más del CSV: 0 diferencias y 0 bajo umbral.
  - `contraste_componentes.py --check` rc 0 e `--independiente` OK.
- **S1-02:** goose `up`, `down-to 0` y `up` con rc 0. El `down` de la 00005 deja 0 claves y 0 disparadores. Los disparadores rechazan valores fuera de rango.
- **S1-16 y diccionario:** `gen_dict.py` regenera el diccionario idéntico byte a byte.
- **C1-13:** `version: "1.5.0"` en la línea 7, rc 0.
- **C1-21:** 221 operaciones, 176 rutas y 267 esquemas, con 22 extractos idénticos. `verify_c1 --db` da TODO OK.
- **Criterios:** 261 (201 / 60 / 0), 0 ids o severidades cambiados contra HEAD. `make verify-criterios` rc 0.
- **Verificación global:** `make verify-all` da **EXIT 2** (`verify-wf1-cobertura`, 4 pruebas en 3 paquetes). Con `make -k`, es el único objetivo que falla.

## Supuestos
1. Numerar el barómetro como §8.9 y no como §8.5 se acepta porque va PROVISIONAL y con bloqueo. El bajo B-r1-2 de la ronda 1 sigue abierto.
2. Que A-1 le toque a `api-tablero` no baja su severidad. La regla dice que una verificación rota es alta.

## Bloqueos
No abrí ninguno.

## Pendientes
- **`api-tablero`:** A-1.
- **Corrector de la ronda 3:**
  - M-N1, en `c1_plantilla.md` y después con `build.py`.
  - Los bajos B-N1 a B-N5.
  - Los que siguen abiertos de la ronda 1: B-r1-1 a B-r1-4.

## Propuestas para Axel
Ninguna nueva. Las PREGUNTAS de `correccion-especificacion-r1` siguen en pie: el tope de §G en los Resúmenes que no son Inicio y la opción (a) o (b) de `C1-tablero`, punto 3.

## Paquetes instalados
Ninguno.

## Consumo aproximado
55 acciones, unos 190 KB leídos (por secciones; `openapi.yaml` solo con un grep y por extractos) y unos 17 KB escritos. Borré la carpeta temporal `datos/tmp/rev-esp-r2/`.
