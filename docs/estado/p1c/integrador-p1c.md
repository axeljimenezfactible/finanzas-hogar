# integrador-p1c · Integración del flujo P1c (puerta 1b) · 11-oct-2026

## Hecho
- Leí el protocolo, decisiones §K (K35 a K44 completas, K1 a K34 por renglón), los veredictos y hallazgos de las revisiones finales, los reportes de los 26 agentes (por secciones) y `docs/estado/p1/INTEGRACION.md` §7.
- **Makefile.**
  - Etiquetas «esquema 1.3» y «U1 1.2».
  - `verify-sql` ya tomaba las migraciones por comodín. Ahora muestra el orden del Down: «Down en orden inverso (00003,00002,00001) deja 0 objetos».
  - `verify-s1` exige que `s1_verifica.py` aplique cada migración del comodín (probado con una migración falsa: falla).
  - `verify-criterios` exige «INDICE contra conteo propio: 0 diferencias» (solicitud de `criterios-1.2.md`) y rotula el subtotal del paso 1 (134 = S1 a S8 y WF4).
  - `verify-all` queda completo: WF-0, C1, S1, U1 y criterios.
- `make verify-all` dos veces con goose v3.28.0, después del último cambio al Makefile: las dos pasan y dan salidas idénticas.
- **README.md.** Estado al 11-oct-2026 (puerta 1b aplicada: esquema 1.3, C1 v1.2, U1 1.2, colorimetría 1.2 y criterios 1.2 no aprobados; la Pista 3 la ejecuta Axel con `docs/design/entrada/`; sigue WF-1). También la tabla de versiones, el árbol con la 00003, la salida esperada y «Cómo sigue».
- **Conteo de PROVISIONAL** en migraciones, diccionario, reglas, U1, `openapi.yaml`, `C1-datos.md`, `colorimetria.md`, `C2-colores.json` y criterios: 119 líneas. Toda marca activa cita solo `S1-parametros-avisos`, el único bloqueo abierto; ninguna cita un bloqueo con «Respuesta de Axel».
- **`docs/estado/bloqueos/INDICE.md`.** Conteo en «Resumen» (43 archivos: 42 resueltos y 1 abierto) y sección «Bloqueos nuevos de P1c».
- **`docs/criterios/INDICE.md`.** Conteos ya reconciliados (`conteo.py`: 261 = 201 / 60 / 0, 0 diferencias); no hizo falta tocarlo.
- `docs/estado/p1c/INTEGRACION.md` y el commit.

## Archivos
- Modificados: `Makefile`, `README.md`, `docs/estado/bloqueos/INDICE.md`.
- Creados: `docs/estado/p1c/INTEGRACION.md` y este reporte.
- El commit incluye todo el árbol de P1c (`git add -A`).

## Verificación
- `make verify-all GOOSE=…/scratchpad/gobin/goose` dos veces: exit 0 y 0, 10.2 s y 10.1 s, `diff` de las dos salidas vacío.
  - **SQL:** 55 tablas STRICT, 152 índices, 75 disparadores, FK 0, integridad ok. Down en orden inverso deja 0 objetos y Up otra vez da 55 tablas. goose up, down-to 0 y up sin errores.
  - **Go:** sin errores.
  - **Colorimetría y OpenAPI:** 202 pares con 0 que fallen; `openapi.yaml: OK`.
  - **C1:** 21 archivos idénticos al generador; `verify_c1` 29 OK.
  - **S1:** «RESUMEN TODO OK», con las 3 migraciones.
  - **U1:** 63 OK y 0 PENDIENTE.
  - **Criterios:** 261, 305 citas, 178 cifras canon con 0 diferencias, 12 de 12; las 11 comprobaciones dan lo esperado.
- `make verify-s1 MIGRATIONS="…00001_esquema.sql …00009_falsa.sql"`: «FALLA: s1_verifica.py no aplicó 00009_falsa.sql», exit 2. La prueba negativa funciona.
- `python3 -I docs/criterios/tools/conteo.py`: 261 (201 / 60 / 0) y «INDICE contra conteo propio: 0 diferencias».
- Bloqueos con «Respuesta de Axel» (o «Resolución»): 42 de 43; abierto: `S1-parametros-avisos`.
- Conteo de PROVISIONAL por archivo: tabla en `INTEGRACION.md` §4.

## Supuestos
- Un bloqueo cuenta como resuelto si trae «Respuesta de Axel». `criterios-nxx` trae «Resolución» de WF-0 y el índice ya lo daba por resuelto.
- La ronda 4 de C1 v1.1 atendió el punto 1 del §7 de P1 (commit `3f2533b`). Lo tomo de su mensaje; no repetí esa revisión.

## Bloqueos
- Ninguno nuevo del integrador. El único abierto es `docs/estado/bloqueos/S1-parametros-avisos.md`, de S1.

## Pendientes
- Están en `docs/estado/p1c/INTEGRACION.md` §7.
- Lo principal: los criterios 1.2 no quedaron aprobados (hallazgo A alto; B y C medios), con las rutas en §7.1. El orquestador decide.

## Solicitudes fuera de alcance
- `docs/criterios/C1.md`, `U4.md`, `S7.md` e `INDICE.md` (§6, líneas 68, 96 y 174): corregir el hallazgo A de `REVISION-criterios-r3.md`, porque los criterios contradicen K44 y el esquema 1.3. `P3-claude-design.md`: P3-06 y P3-04 (B y C).
- `backend/internal/api/c1gen/c1_plantilla.md`: cabecera de `C1-datos.md` con «puerta 1b» y «esquema 1.3» (I-01); regenerar con `build.py`.
- `docs/criterios/S4.md` S4-07: decir que `S4-meta-ajustes` es un bloqueo que S4 abriría, no uno existente (I-02).

## Propuestas para Axel
- Están consolidadas en `INTEGRACION.md` §8: el bloqueo `S1-parametros-avisos` (S1 recomienda B; criterios r1 recomendó A), los mínimos y máximos de los 15 parámetros con su justificación, «restablecer», la sección de cada parámetro, `advance_days`, el MB del adjunto y la propuesta K35 de WF-3.

## Consumo aproximado
- Acciones de herramienta: ~52 de 60.
- KB leídos: ~175 de 400.
- KB escritos: ~40 de 150.
- Sin paquetes instalados. Temporales en el scratchpad, borrados; `datos/tmp` no queda.
