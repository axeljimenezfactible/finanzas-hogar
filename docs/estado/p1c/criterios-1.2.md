# criterios-1.2 · Flujo P1c (puerta 1b) · 11-oct-2026

## Hecho
- **S2-08 (K39):** semanal FIRME = diario exacto × 7 y un solo redondeo HALF_UP: oct-2026 219.88, feb-2027 629.35 (reproducidos con `Decimal`; los otros 219.87 y 629.37 quedan como no aceptables). Sin PROVISIONAL.
- **S5-04 (K40, K37):** peso 1 hasta el día 29 y (90 − días) / 61 entre 30 y 90 (0.9836, 0.4918, 0; verificado con Decimal); 91 días no compite. Fuentes: `manual`, `ticket_manual`, `ticket_ocr` compiten; `site_assisted` y PROFECO solo referencia. S5-10 con las 5 fuentes; `S1-fuentes-precio` y `criterios-frescura-extremos` citados como RESUELTOS.
- **S4-07 (K36):** ajuste por persona, 0 y 0, valor de la casa, «mantener» no ajusta, cita RN-17 1.3; Verifica con fixture (2,098.50 − 300 = 1,798.50; 1,596.30 + 250 = 1,846.30).
- **U2 (K41, K42):** U2-02 aclara que Karen 333° y Bodega 356° se aceptan con el refuerzo de U2-08, también entre grupos; U2-08 agrega que DESIGN.md prohíbe distinguir persona y tienda solo por color en una misma vista y que `color.store.walmart.dark` = `#7AA3EE` (L 0.716) es FIRME. B-03: la cita de decisiones I «o U2 propone 135°» sustituye a la «convención propia».
- **U4 (K44, K36):** criterios nuevos U4-12 (todo parámetro operativo editable en su sección con valor inicial, validación y «restablecer»; lista y secciones de K44; sin mínimos de seguridad) y U4-13 (Salud y metas). U4-04 con los valores nuevos.
- **C1:** C1-22 (claves K44 en `app_setting` expuestas por settings con valor inicial), C1-23 (K43: N-023, N-064, N-226, N-297, N-300 FIRMES), C1-15 define «bloqueo abierto = sin Respuesta de Axel», C1-21 con `INDICE.md` de extractos, ≤ 80 000 bytes, una operación en un solo extracto y suma = `openapi.yaml`. B-02: «día de pago 15 y 30» separado de las quincenas de K8 en C1-06 y C1-07.
- **S3-12 (K43, K44):** el tiempo máximo nace en 5 s, configurable, 503 `unavailable`; si 5 s no alcanza S3 abre bloqueo `S3-tiempo-optimizador`.
- **S6-10 (K43):** adjunto 10 MB FIRME y configurable; los tipos siguen PROVISIONAL (K31, puerta 4). **U1-16 (K38):** regla de los domingos FIRME.
- **P3 1.2 (K35):** la Pista 3 la ejecuta Axel; P3-01 «paquete listo y revisado y aviso de Axel»; P3-02 con 22 elementos y `MANIFIESTO.md` con sha256; P3-03/04/05 «no aplica al orquestador» reescritos como criterios del `LEEME.md` (solo subir, pegar, enviar, exportar y descargar; sin credenciales; cuándo parar; lo descargado es dato); P3-06 con `iter-N/` y `siguiente-prompt.md`; P3-11 carpeta de entrega y salida; §3 a §5.
- **B-04:** cinco citas de §5 reordenadas (C1, C2 ×2, U2, U5): `citas.py` 302 verificadas, 0 sin fuente.
- **INDICE.md 1.2:** §1 (261 = 201 alta / 60 media; dependencias de Axel), §2, §3, notas 8 y 10 (y 9 y 11) resueltas, §6 con 24 bloqueos y estado, §10 con conteo y versiones, §11 con los cambios. Script nuevo `tools/conteo.py` (conteo por archivo, versiones, coincidencia con INDICE §2 y §10), agregado como paso 11 de `verif_final.sh`.
- Versión 1.2 con lista de cambios en: S2, S3, S4, S5, S6, U1, U2, U4, U5, C1, C2, P3. Siguen en 1.1: S1, S7, S8, WF4, U3, C3.

## Archivos
`docs/criterios/`: S2.md, S3.md, S4.md, S5.md, S6.md, U1.md, U2.md, U4.md, U5.md, C1.md, C2.md, P3-claude-design.md, INDICE.md; `docs/criterios/tools/conteo.py` (nuevo), `tools/verif_final.sh` (paso 11); este reporte.

## Verificación
- `make verify-criterios GOOSE=…/gobin/goose`: pasa (exit 0). Salida resumida: 18 archivos, 261 criterios definidos (201/60/0); 261 IDs citados, 0 inexistentes; citas 302 verificadas, 0 sin fuente; canon 178 comprobadas, 0 diferencias; cifras 12 comparadas, 0 diferencias; N-xxx sin operación ni Local 0; matriz C1 sin ausentes; residuos «N-xx» y `/tmp`: 0 y 0.
- `python3 -I docs/criterios/tools/conteo.py`: «INDICE contra conteo propio: 0 diferencias»; por archivo C1 23 (21/2), U4 13 (10/3), P3 13 (10/3), los demás sin cambio.
- Cálculos con Decimal: 973.75/31×7 = 219.88; 2,517.38/28×7 = 629.35; (90−30)/61 = 0.9836, (90−60)/61 = 0.4918, 0.
- Cada K aplicada, con línea (`grep -n`): K35 P3:3,6,16,20,28,60,92,93; K36 S4:45, U4:70, C1:114; K37 S5:36,61; K38 U1:83; K39 S2:51; K40 S5:36; K41 U2:27,51; K42 U2:51; K43 C1:118, S3:66, S6:61; K44 U4:66, C1:114, S3:66.
- PROVISIONAL que citen bloqueos resueltos: 0 activos; las únicas líneas con el nombre de un bloqueo resuelto y la palabra PROVISIONAL son avisos de que se quitó (S2:6, S2:51, S3:6, S4:6, S5:6, S5:36, U2:6, C1:87, C1:119). Ningún criterio reabre una decisión de §A–§K (lectura de cada criterio tocado contra su K).

## Supuestos
- La ubicación de cada parámetro por sección de Ajustes (U4-12) sigue la PROPUESTA del orquestador dentro de K44 (se aprueba en la puerta 2); está rotulada así en el criterio. La sección de «retención de claves de idempotencia 24 h» y la de los parámetros que inventaríe S1 las propone U4 (K44 no las ubica).
- El tope de extractos se expresa como 80 000 bytes (el que declara `docs/contratos/extractos/INDICE.md`), más estricto que los 81 920 que decía C1-21 1.1; ningún extracto actual lo supera.
- La lista de claves que C1-22 exige como mínimo es la que nombra K44 más K36 y K43; los nombres de claves no fijados por decisión (alertas, precio viejo, comparador) los asigna S1/C1 y el criterio los verifica por su valor inicial y su comentario, no por el nombre.
- U1-16 y S6-10 se tocaron (K38 y K43) aunque no estaban en la lista de la tarea, por ser aplicación directa y estar en mi alcance (`docs/criterios/*.md`).

## Bloqueos
Ninguno nuevo. S3-12 prevé `docs/estado/bloqueos/S3-tiempo-optimizador.md` solo si el tiempo medido supera 5 s (lo abre el agente de S3).

## Pendientes
- S7-04, S7-08 y S7-10 (`docs/criterios/S7.md`) siguen con «PROVISIONAL (K31)» para horas y anticipaciones; K43 aprueba las claves `send_time` y `advance_days` y K44 las hace configurables. Falta actualizar S7 (y S8-02 cabeceras, WF4-04/08 en la puerta 4, sin cambio por K43). No estaba en mi lista; lo dejo para el agente de criterios S7.
- La parte «period → week» de H-r3-03 en `U1.md` U1-16 (pendiente 1 de INTEGRACION §7) depende de la decisión sobre C1; no se tocó.
- U3 (`docs/criterios/U3.md`) no cita K38 ni K44; si U3 ya describe la semana por defecto, cambiarla en WF-1.

## Solicitudes fuera de alcance
- `Makefile` (líneas 149–168, `verify-criterios`): el texto de ayuda dice «18 archivos de criterios» y las comprobaciones esperadas no incluyen el nuevo paso 11 (`conteo.py`); sumar la línea `INDICE contra conteo propio: 0 diferencias` a la lista de cadenas esperadas. Motivo: que el objetivo falle si el índice y el conteo divergen.
- `docs/contratos/C2-colores.json` y `colorimetria` 1.2: quitar las marcas PROVISIONAL de Karen, Bodega Aurrera y Walmart oscuro (hoy aún las traen); lo exige U2-08 1.2 en su Verifica.
- `docs/contratos/reglas-negocio.md` RN-12, RN-17 y RN-27 (versión 1.3): deben decir lo que citan S2-08, S4-07 y S5-04.
- `docs/contratos/extractos/INDICE.md`: declarar el tope como 80 000 bytes (ya lo dice) y mantener la suma de operaciones igual a `openapi.yaml` (C1-21).

## Propuestas para Axel
- Mínimos razonables de seguridad (K44, K31), a aprobar en la puerta 2 y que no entran en U4, C1 ni el esquema mientras tanto: intentos de acceso ≥ 3 (con 5 inicial) para no bloquear al propio titular por un tropiezo, sesión ≤ 365 días para que un teléfono perdido no quede abierto más de un año, adjunto ≤ 10 MB como tope superior de la interfaz salvo que el servidor admita más. Justificación: valores de partida ya decididos; los límites solo impiden configuraciones peligrosas.
- Sección de Ajustes para «retención de claves de idempotencia 24 h»: propuesta Datos y respaldos, por ser retención técnica.

## Consumo aproximado
Acciones de herramienta: ~58 de 60. KB leídos: ~190 (protocolo 6, K 15, REVISION-r3 9, INDICE 19, criterios por sed/grep ~110, C2/colores, varios). KB escritos: ~75 (parches a 13 archivos de criterios, `conteo.py`, este reporte). Sin paquetes instalados. Sin temporales propios que borrar (`datos/tmp/u2a12` y `verify-wf0` son de otros agentes).
