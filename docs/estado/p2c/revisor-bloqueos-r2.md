# P2c · revisor-bloqueos-r2 (11-oct-2026, noche)

Agente: `revisor-bloqueos-r2`, revisor independiente, ronda 2 de 3. Alcance de escritura: `docs/estado/p2c/REVISION-bloqueos-r2.md` y este reporte. No corregí nada de lo revisado.

## Hecho
- HECHO: revisé los 27 bloqueos tocados, `INDICE.md`, `decisiones.md` §K46 a §K52 y los reportes `bloqueos-respuestas.md` y `correccion-bloqueos-r1.md`. Medí contra RESPUESTAS.md, PUERTA-2.md (§2.1, §3, §4 y §8), el protocolo §2 y §3 y el formato de las puertas 1 y 1b del INDICE.
- HECHO: veredicto **NO APROBADO, con 1 alto, 2 medios y 4 bajos** (detalle en `docs/estado/p2c/REVISION-bloqueos-r2.md`).
  - Alto: H-r2-01. La verificación global sigue rota mientras no se aplique la solicitud 1.
  - Medios: H-r2-02 (el INDICE y K46 no registran 3 bloqueos que `contratos-puerta2` abrió después de la corrección) y H-r2-03 (la cita por línea de Q2-12 cambió de lugar con la reescritura de `S2-interfaz.md`).
  - Bajos: H-r2-04 a H-r2-07.
  - Cerrados de la r1: H-r1-02 y H-r1-04.

## Archivos
- `docs/estado/p2c/REVISION-bloqueos-r2.md` (informe)
- `docs/estado/p2c/revisor-bloqueos-r2.md` (este reporte)

## Verificación
Scripts propios (`v.py`, `m.py` y `k.py`) y dos clones locales de solo lectura (`git clone` de HEAD; en `b` copié los 27 bloqueos, `INDICE.md` y `decisiones.md` del árbol). Todo estuvo en `datos/tmp/rev-bloq-r2/` y lo borré al final.
- `v.py`:
  - 69 archivos en HEAD, 27 modificados, y los 27 solo agregan texto al final. 42 resueltos en HEAD y 0 tocados.
  - 62 entradas Q2 y 0 faltantes. 53 letras iguales a §8: Q2-11 «B» y Q2-64 «A». La única diferencia es de forma, en Q2-32 («A (propuesta del agente, letras de PUERTA-2)»).
  - Los 53 textos literales están en HEAD. Los 9 «sin recomendación» no traen respuesta.
  - Tabla del INDICE: 0 diferencias.
  - Carpeta real: 72 archivos, 63 RESUELTOS y 9 ABIERTOS; con la regla de los scripts, 67 y 5.
- `m.py`: muestra de 10 (Q2-11, 13, 14, 28, 32, 33, 42, 64, 70, 71). Letra y texto coinciden con PUERTA-2 y con el archivo.
- `k.py`: K46 tiene 98 filas y 0 diferencias con §8, con 15 sin respuesta. Las citas de Axel en §K son 3 y están todas en RESPUESTAS.md. K50 es igual a PUERTA-2 §4.
- `make verify-c1` y `make verify-s1`:
  - Clon `a`: «RESULTADO: TODO OK» y «RESUMEN TODO OK».
  - Clon `b`: «FALLA C1-15» y `RESUMEN FALLAS: ['S1-16', 'K44']`.
  - Árbol: «parse_schema: clave @param repetida» y `['S1-16', 'NEG12', 'K44', 'PROG13']`.
- `make verify-u1` y `make verify-criterios` en `a` y en `b`: rc=0. `verify-criterios` en el árbol: rc=0.
- `sed -n 250p;447p` de `S2-interfaz.md`: en HEAD son `variable_estimated_cents` y la tabla de bloqueos. En el árbol son `credit_payments.household_cents` y `TestInvarianteDeudaViva`, y el valor se movió a la l. 258 y la l. 460.

## Supuestos
- El árbol de trabajo a las 22:51 es lo que irá al commit. Por eso mido el INDICE contra los 72 archivos que hay en el disco y la cita de Q2-12 contra el `S2-interfaz.md` del árbol, aunque las dos diferencias las causó un agente que trabaja al mismo tiempo.
- Las cifras sin signo de pesos que se copian de una opción original cuentan como baja, igual que en la r1.

## Bloqueos
- Ninguno nuevo.

## Solicitudes fuera de alcance
1. Solicitud 1 de `correccion-bloqueos-r1.md`: quitar las 18 marcas PROVISIONAL que citan `S1-parametros-avisos`. La hace el agente que aplica Q2-11 = B. Sin esto no se cierra H-r2-01.
2. `contratos-puerta2`: pedir en su reporte que el INDICE registre `S2-ahorro-quincenal`, `S2-metas-vencida` y `S4-presentacion-respaldo`. Para H-r2-02.
3. Siguiente corrección de bloqueos o integrador: aplicar H-r2-02 a H-r2-05 y H-r2-07.

## Pendientes
- Ronda 3: comprobar H-r2-01 con `make verify-all` sobre el árbol completo, cuando terminen los agentes de esquema y contratos.

## Propuestas para Axel
- Ninguna nueva. Siguen abiertas Q2-12 (C o A), la pregunta 4 de `S2-ventanas` (casos 1 y 2), la variante de Q2-07 (`S4-presentacion-respaldo`), `S2-ahorro-quincenal`, `S2-metas-vencida` y qué sigue en `S5-terminos`.

## Consumo aproximado
Acciones de herramienta: ~31 · KB leídos: ~120 · KB escritos: ~19.
