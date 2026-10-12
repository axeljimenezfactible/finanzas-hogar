# P2c · revisor-contratos-r3 · 11-oct-2026, noche

Agente revisor independiente, ronda 3 de 3, de la cadena de contratos (S2 a S5). No escribí nada de lo revisado y no corregí nada. El informe completo está en `docs/estado/p2c/REVISION-contratos-r3.md`.

## Hecho
- **Veredicto: NO APROBADO.** Queda 1 alto, 1 medio y 5 bajos.
  - **A-1 (alta).** `verify-c1` falla con «gen2: tipo de parámetro desconocido: hora HH:MM». Está fuera del alcance de los contratos.
  - **M-1 (media).** `INDICE.md` y K46 todavía no registran los 4 bloqueos que abrieron las correcciones de P2c. La corrección está preparada en `correccion-contratos-r2.md` y no se aplicó.
- **Cerrado: M-2 de la ronda 2** (lista para WF-3). Lo que queda de él es un bajo: `correccion-contratos-r1.md` no remite a r2.
- **Comprobé con mis scripts (a) a (e):**
  - Q2-01 a Q2-08 están aplicadas.
  - Ninguna pregunta «sin recomendación» quedó decidida.
  - Toda marca PROVISIONAL activa cita un bloqueo ABIERTO.
  - Ninguna cifra canon desaparece y los oráculos dan 0 diferencias.
  - Las listas de cambios cubren los 101 bloques del diff y los 106 renglones de código con bloqueos resueltos.
- Los 9 archivos de la cadena tienen el mismo sha256 que al cerrar la ronda 2 (`sha256sum -c` contra `datos/tmp/correccion-contratos-r2/sha-antes.txt`: 9 OK).

## Verificación (comandos y salidas)
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP=…/datos/tmp/revisor-contratos-r3/vw` → **rc=2**.
  - Falla: `gen2: tipo de parámetro desconocido: hora HH:MM` y `make: *** [Makefile:130: verify-c1] Error 1`.
  - Pasan:
    - WF-0: «verificación completa OK».
    - S1: «RESUMEN TODO OK».
    - U1: 63 OK, 0 pendientes.
    - Criterios: «TOTAL comprobadas: 178 | con diferencias: 0» e «INDICE contra conteo propio: 0 diferencias».
    - K25: motor_cred 93.7 %, motor_fin 89.2 %, recetas 93.2 % y precios 96.1 %.
    - Oráculos: S2 30/30 (2320 hojas); S3 906/906, canon 437/437; S4 26/26 (8451); S5 29/29 (747).
    - `model.json` sha256 `1a57066689dd…` intacto.
    - WF-1: «verificación completa OK».
- `r3.py c`, PROVISIONAL contra el estado de cada bloqueo:
  - Hay 79 archivos: 62 RESUELTO, 14 ABIERTO y 3 sin línea de estado (`criterios-nxx` y los 2 de la Pista 3).
  - Ninguna marca activa cita solo bloqueos resueltos. Los renglones que salen son definiciones, encabezados o bitácora: S2 l. 7 y 458; S3 l. 612, 658, 884, 893 y 896; S4 l. 7, 20, 26 y 31; S5 l. 7, 22, 28 y 34.
- `r3.py d`, cifras de `HEAD` contra el árbol de trabajo: no desaparece ninguna. Las nuevas son de ejemplos: S2 100000, 199.00 y 4,000.00; S4 1350, 2250 y 2300; S5 0.00, 1,021.00 y 102100.
- `hunks.py` y `hk.py`, diff por bloque y por palabra: 101 bloques (S2 21, S3 34, S4 23, S5 23), cada uno con su renglón en las listas de los 3 reportes de la cadena.
- `code3.py`: 52 renglones con PROVISIONAL y un bloqueo RESUELTO a ±2. Los 40 de motores y oráculos están en las listas; los 12 de `api/c1gen` quedan fuera del objeto.
- Recetas de `correccion-contratos-r2.md` §E: la 1 da **40** y la 2 da **106**. Los 106 están en las listas: 0 sin listar.
- Imprimí los 7 sitios de M-2, `card.go` l. 89 y las pruebas `engine_test.go` l. 171 y 183, `comparador_test.go` l. 394 y `timeline_test.go` l. 460–461. Están en esos renglones, con ese texto.
- `grep` de los 4 bloqueos nuevos en `INDICE.md` y en `decisiones.md`: 0 menciones cada uno. `grep Remisión` en los 3 bloqueos resueltos: 0.
- Q2-07: `recetas/catalog.go` l. 66–77 (`fallbackPresentation`) aplica la típica y luego la primera del catálogo. Es la regla que S4 l. 185 escribe como «la que aplica hoy el motor».

## Supuestos
1. Uso como fuente del estado de un bloqueo su última línea «Estado del bloqueo (puerta 2)» o «- **Estado:**», que es la regla de conteo del INDICE desde la puerta 2.
2. Una marca PROVISIONAL en la bitácora de versiones o de hallazgos no es una marca activa. Si contradice el estado actual sin decir que quedó superada, la cuento como baja (B-5).
3. Lo que cité de `api/c1gen` (C1 y S1-parametros-avisos) queda fuera del objeto de esta revisión. Va como observación y no cuenta para el veredicto.
4. `datos/tmp/r3/` no es mío: lo creó otro agente a las 23:43–23:45. No lo toqué.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- A-1: el agente de C1 o el integrador aplica la solicitud 1 de `S1-puerta2.md` en `gen2_schemas.py`.
- M-1: el integrador aplica la solicitud 2 de `correccion-contratos-r2.md`, con el script o con el Anexo.
- B-1 a B-5: los bajos, en la integración o en WF-3.
- Era la ronda 3 de 3. Según el protocolo, lo que queda abierto pasa al integrador de P2c.

## Propuestas para Axel
- Ninguna nueva. Siguen abiertas las de la cadena:
  - contestar por id `S3-iva-lectura`, Q2-90 (`S4-empates-sin-confirmar`), `S2-ventanas` pregunta 4, `S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo` y `S5-carritos-granel`;
  - entregar los PDF de un estado de cuenta de tarjeta y uno de crédito de tienda (`S3-formula-tarjeta-pdf`).

## Solicitudes fuera de alcance
1. **Integrador o agente de C1:** A-1. Al mismo tiempo, alinear `gen2_schemas.py` l. 364 y l. 374 con Q2-11 = B (todo en `app_setting`), porque hoy citan como PROVISIONAL `S1-parametros-avisos`, que está RESUELTO.
2. **Integrador:** M-1, y la remisión de r1 a r2 (B-1).

## Consumo aproximado
44 acciones de herramienta · unos 150 KB leídos · unos 25 KB escritos: 12 KB del informe, 5 KB de este reporte y unos 8 KB de scripts temporales, ya borrados.
