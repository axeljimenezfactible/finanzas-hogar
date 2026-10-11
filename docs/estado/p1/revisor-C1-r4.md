# Reporte de agente · revisor-C1-r4 (flujo P1b, 11-oct-2026)

Revisor independiente de la ronda 4 de C1 v1.1. No corrigió nada y no hizo commit. Informe completo: `docs/estado/p1/REVISION-C1-r4.md`.

## Hecho
- Leí el protocolo (completo), `decisiones.md` §K, `REVISION-C1-r3.md`, `correccion-C1-r4.md`, `git diff HEAD --stat` y el diff de la fuente (c1gen, migraciones, U1, criterios U1). De `C1-datos.md`, `openapi.yaml`, el diccionario y los extractos solo leí lo que salió de grep, sed -n y mis scripts.
- Reejecuté las verificaciones (a) a (n) con comandos propios y tres scripts en `datos/tmp/rev-c1-r4/` (`neg.py`, `chk.py` y `chk2.py`), con prueba negativa de las dos reglas nuevas de `verify_c1.py`.
- Veredicto: **APROBADO**, con 0 hallazgos altos y 0 medios. H-r3-01 a H-r3-04 cerrados. Los 5 bajos de r3 siguen abiertos (para WF-1). Hay 3 observaciones bajas nuevas (N-r4-01 a N-r4-03).

## Archivos
- `docs/estado/p1/REVISION-C1-r4.md` (nuevo)
- `docs/estado/p1/revisor-C1-r4.md` (este reporte)
- `datos/tmp/rev-c1-r4/` (temporal, borrado al terminar)

## Verificación
- `build.py --check`: IDÉNTICO en los 21 archivos.
- goose v3.28.0 up (00001 y 00002, versión 2), luego `verify_c1.py --db`: 26 OK, 0 FALLA, «RESULTADO: TODO OK».
- Prueba negativa (`neg.py` y `verify_c1.py --root`), con código 1 en cada caso y como única falla:
  - sin `ListWeek` en `getShoppingList`: FALLA C1-06;
  - sin `ListWeek` en `getBasketComparison`: FALLA C1-06;
  - sin `x-provisional` en `getMealPlanWeek`: FALLA C1-15 (regla del domingo).
- `openapi-spec-validator==0.9.0`: openapi.yaml y los 17 extractos OK.
- `verificar_u1.py`: TODO OK.
- `go build`, `go vet` y `gofmt -l`: OK, OK y vacío.
- Extractos: 17 archivos, el mayor de 76,202 B. Suma de operationId 201 = `grep -c` 201. 0 repetidas, 0 faltantes y 0 sobrantes.
- `gen_dict.py` a un archivo temporal: idéntico al diccionario (`cmp`). `diff` de las migraciones sin comentarios contra HEAD: vacío.
- `s1_verifica.py` con `GOOSE=…`: RESUMEN TODO OK.
- `make verify-all GOOSE=…`: código 0, «Verificación ampliada OK».
- Marcas PROVISIONAL: 7 líneas y 5 `x-provisional` en openapi, todas con un bloqueo abierto.
- Cifras nuevas en el diff: 0 inventadas.

## Supuestos
- Ninguno que cambie cifras, esquema o endpoints. Clasifiqué N-r4-01 como baja y no como media: §3 asigna 422 a la «validación de campos» y `default` cubre el código con `Error`, así que C1 no contradice a U1-16; solo le falta precisión.

## Bloqueos
- Ninguno nuevo. Se cita el bloqueo existente y abierto `docs/estado/bloqueos/U1-lista-vigente-domingo.md`.

## Pendientes
- Para WF-1 (no bloquean): H-r3-05, H-r3-06, H6, H7 y H8 de `REVISION-C1-r3.md`, y N-r4-01, N-r4-02 y N-r4-03 de `REVISION-C1-r4.md`, con rutas y corrección sugerida en ese informe.

## Solicitudes fuera de alcance
- Ninguna. El revisor no corrige. Las correcciones sugeridas están en la tabla de hallazgos del informe.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Unas 40 acciones de herramienta.
- Leído: unos 190 KB (diff de la fuente, reportes y salidas de grep, sed y scripts; ningún archivo grande se leyó entero).
- Escrito: unos 30 KB (informe de 17 KB, este reporte y los scripts temporales).
