# Revisión independiente · bloqueos de la puerta 2 · ronda 1 de 3 (11-oct-2026, noche)

Revisor: `revisor-bloqueos-r1`. No escribí nada de lo revisado y no corregí nada. Base: HEAD 1ccb468 + árbol de trabajo (27 bloqueos, `INDICE.md`, `decisiones.md` §K46–§K52, reporte `docs/estado/p2c/bloqueos-respuestas.md`). Fuentes contra las que se mide: `docs/estado/puerta2/RESPUESTAS.md`, `docs/estado/puerta2/PUERTA-2.md` (§0, §2.1, §3, §4, §5, §6, §8), protocolo de agentes §2, §3 y §5, y el formato de la puerta 1 en `INDICE.md`. Scripts propios en `datos/tmp/rev-bloq-r1/` (borrados al final).

## Veredicto
**NO APROBADO: 1 alto, 3 medios y 2 bajos.**
La fidelidad está bien: las 53 letras con respuesta coinciden con la hoja §8, Q2-11 = B y Q2-64 = A coinciden literalmente, ningún «sin recomendación» tiene respuesta, la tabla de K46 tiene 0 discrepancias y los valores de K50 son iguales a PUERTA-2 §4. Lo que falla: el cambio rompe `verify-c1` y `verify-s1` y el reporte no lo declara; la nueva regla de conteo no la siguen las herramientas; hay una cita a un archivo equivocado; y faltan las solicitudes fuera de alcance.

## Verificaciones (evidencia propia)
| # | Qué | Comando o script | Resultado |
|---|---|---|---|
| a1 | Archivos modificados en `docs/estado/bloqueos/` | `git diff --stat HEAD -- docs/estado/bloqueos` | 28 archivos (27 bloqueos + INDICE), 665 inserciones, 4 borrados (todos en INDICE). 0 sin seguimiento en la carpeta. |
| a2 | Los 42 resueltos no se tocaron y solo se agregó texto | `chk.py`: compara el texto de HEAD con el actual por prefijo; `cnt.py` | Ningún archivo cambió su prefijo de HEAD. 42 resueltos en HEAD (41 con «Respuesta de Axel» + `criterios-nxx`); ninguno en el diff. **OK** |
| b1 | Letra de cada Q2 contra la hoja §8 | `chk.py` (62 entradas: Q2-11 a Q2-71 + caso 2 de Q2-01) | 53 con respuesta: todas iguales a la recomendada; Q2-11 = **B** y Q2-64 = **A** («respuesta explícita de Axel»); Q2-32 = «A (propuesta del agente, letras de PUERTA-2)». Faltantes en Q2-11 a Q2-71: ninguno. **OK** |
| b2 | Texto literal de la opción contra el archivo en HEAD | `chk.py` (cada «Texto de la opción (literal del archivo)» buscado en `git show HEAD:<archivo>`) | 53 de 53 aparecen literales. **OK** |
| b3 | Muestra de 15 contra PUERTA-2 §3 (Q2-11, 13, 17, 21, 26, 33, 34, 62, 64, 65, 66, 67, 68, 70, 71) | lectura de PUERTA-2 l. 97–290 y de las entradas | Texto y sentido coinciden. Q2-13: el archivo dice «HALF_UP en la 1 y el resto en la 2» y PUERTA-2 «el centavo de redondeo en la 1»; son equivalentes y la entrada lo anota. **OK** |
| b4 | «Sin recomendación» sin respuesta | `chk.py` | Q2-48, 50, 51, 52, 56, 57, 58 y el caso 2 de Q2-01 llevan «PROVISIONAL … Pasa a WF-3» y **Estado: ABIERTO**, sin «Respuesta:». Q2-12 («C o A») también abierta, ver H-r1-06. **OK** |
| c1 | Conteo del INDICE (63 resueltos, 6 abiertos) | `cnt.py` (línea «Estado del bloqueo (puerta 2)» en los tocados; «Respuesta de Axel» o `criterios-nxx` en los demás) | 69 · 63 · 6; los 6 abiertos son los que nombra el INDICE. Filas 34+8+1+26 = 69 y 34+8+1+20 = 63. **OK** con la regla nueva; con la regla anterior del mismo INDICE (l. 26) salen 67 y 2 → H-r1-02 |
| c2 | Tabla «Respuestas de la puerta 2» del INDICE | lectura del diff | 63 filas, letras iguales a las de los archivos. **OK** |
| d1 | Tabla de K46 contra la hoja §8 | `k46.py` | 98 entradas, 0 discrepancias; 15 sin respuesta (Q2-12, 48, 50–52, 56–58, 85, 87–91, 94). 98 − 15 = 83. **OK** |
| d2 | K50 contra PUERTA-2 §4 | lectura | 21:30 / domingo 19:00 / día 1 08:00; 1 día y 09:00; 5 s con 1 a 30 s; 7 días; 30 s; S ≥ 1, Z ≥ S, (Z − días) ÷ (Z − S + 1). Iguales. **OK** |
| d3 | Citas de Axel en §K46–§K52 | lectura | Solo las de RESPUESTAS.md (modo, aclaración «C · Mixta», autorización de Q2-64). K51 y K52 iguales a PUERTA-2 §5 y §6. **OK** |
| e1 | Verificación que lee `docs/estado/bloqueos/` | `make verify-c1` y `make verify-s1` en el árbol; las mismas en una copia `git archive HEAD`, y en esa copia con solo los 28 archivos de bloqueos del árbol | Árbol: **FALLA C1-15**, **S1-16 FALLA**, **K44 FALLA**. Copia de HEAD: verify-c1 «RESULTADO: TODO OK», verify-s1 rc=0 (S1-16 OK, K44 OK). Copia + bloqueos del árbol: vuelven las 3 fallas → H-r1-01 |
| e2 | Citas con número de línea en el texto nuevo | `cites.py` | 2 citas; una apunta a otro archivo → H-r1-03 |

## Hallazgos
### H-r1-01 · ALTA · la verificación global queda rota y el reporte dice que no hacía falta correrla
- **Ruta:** `docs/estado/bloqueos/S1-parametros-avisos.md` (sección nueva) junto con `backend/internal/api/c1gen/verify_c1.py` l. 321–340 y `backend/db/tools/s1_verifica.py` l. 242 y 575. Reporte: `docs/estado/p2c/bloqueos-respuestas.md` §Verificación.
- **Qué pasa:** las herramientas dan por resuelto un bloqueo si su archivo contiene «Respuesta de Axel». `S1-parametros-avisos` ya la trae, pero lo citan como PROVISIONAL 6 marcas de `openapi.yaml` (`Notification.params_json`, `send_time`, `advance_days`, `NotificationUpdate`…), 5 o más líneas de `C1-datos.md` (P-26.08, P-26.09, `params_json send_time/advance_days/weekday`) y 19 líneas del esquema (migración 6, diccionario 7, reglas 6). Resultado: falla C1-15 de verify-c1, y S1-16 y K44 de verify-s1. En HEAD pasan; con solo los archivos de bloqueos del árbol fallan (e1). El reporte dice «No se corrió `make verify-all`: solo documentos de estado y `decisiones.md`, sin código», pero esa carpeta es entrada de las dos verificaciones.
- **Corrección sugerida:** correr `make verify-c1` y `make verify-s1` y anotar la falla en Verificación. Abrir una solicitud fuera de alcance: el agente que aplica Q2-11 = B (00003 o su continuación, `diccionario-datos.md`, RN-29, `openapi.yaml` por `c1gen`, `C1-datos.md`, criterios C1-22, U4-12 y S7) quita esas marcas en la misma ronda P2c. El integrador no hace commit de los bloqueos sin ese cambio. Se cierra cuando `make verify-all` pasa con los dos cambios juntos.

### H-r1-02 · MEDIA · dos reglas de conteo; las herramientas siguen la anterior
- **Ruta:** `docs/estado/bloqueos/INDICE.md` l. 26 (regla anterior, sin enmendar) y l. 28 (regla nueva); `verify_c1.py` l. 322–323; `s1_verifica.py` l. 242 y 575; fila «Flujo WF-1» del resumen.
- **Qué pasa:** la regla nueva cuenta por la línea «Estado del bloqueo (puerta 2)». La regla de l. 26 («resuelto si su archivo trae la sección «Respuesta de Axel»») sigue en el INDICE sin decir que ya no aplica, y es la que usan los scripts y PUERTA-2 §3. Cuatro archivos ABIERTOS traen «Respuesta de Axel (puerta 2…)»: `S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa` y `S5-terminos`. Con la regla de los scripts salen 67 resueltos y 2 abiertos; el INDICE dice 63 y 6. Hoy ninguna marca revisada por los scripts cita esos 4, pero la primera marca PROVISIONAL de WF-3 en C1 o en el esquema que los cite (Q2-12, Q2-48, caso 2 de Q2-01) fallará C1-15 o S1-16 aunque el bloqueo siga abierto. Además, la fila dice «20 con «Respuesta de Axel (puerta 2)»», pero esa sección está en 24 archivos de WF-1; 20 de ellos están RESUELTOS.
- **Corrección sugerida:** solicitud fuera de alcance para que los dos scripts lean la última línea «Estado del bloqueo (puerta 2): …» cuando exista y usen «Respuesta de Axel» solo si falta. En el INDICE, l. 26 debe decir que esa regla es la de las puertas 1 y 1b, y la fila debe decir «24 con la sección de la puerta 2, 20 RESUELTOS».

### H-r1-03 · MEDIA · la fuente de lo PROVISIONAL en uso de Q2-12 apunta a otro archivo
- **Ruta:** `docs/estado/bloqueos/S2-flujo-quincenal.md`, entrada «Q2-12» de «Puerta 2 (11-oct-2026): sin respuesta».
- **Qué pasa:** cita «`docs/estado/wf1/INTEGRACION.md` l. 250; l. 447». Ese archivo tiene 323 líneas, y su l. 250 es `| correccion-U4-r2 | 34 | 170 | 42 |`. El valor en uso (`variable_estimated_cents` = `null`, «no entran a ningún saldo») está en `docs/contratos/motores/S2-interfaz.md` l. 250, y su l. 447 lo liga al bloqueo. WF-3 llega a una fila que no tiene que ver.
- **Corrección sugerida:** cambiar la ruta a `docs/contratos/motores/S2-interfaz.md` l. 250 y l. 447.

### H-r1-04 · MEDIA · dependencias fuera de alcance sin solicitud; al reporte le faltan secciones del protocolo §5
- **Ruta:** `docs/estado/p2c/bloqueos-respuestas.md`; `docs/estado/bloqueos/S2-ventanas.md` (caso 2); `docs/plan/decisiones.md` K46 (Q2-07).
- **Qué pasa:** el reporte no tiene «Archivos» ni «Solicitudes fuera de alcance». Quedan cosas que otros tienen que hacer y no están pedidas: (i) la pregunta 4 de `S2-ventanas`. El caso 2 es un PROVISIONAL («comportamiento actual, que no está escrito») cuyo bloqueo no tiene Qué falta, Opciones ni Impacto (formato del protocolo §2); el archivo dice «la redacta quien aplique Q2-01», pero nadie la pidió. (ii) RESPUESTAS.md manda documentar la variante de Q2-07 como PROVISIONAL en `S4-presentacion-respaldo`, que no existe; K46 dice «Q2-07 … con una parte … que sigue abierta» sin citar ningún bloqueo, y el conteo de 69 del INDICE cambiará cuando se cree. (iii) y (iv): H-r1-01 y H-r1-02.
- **Corrección sugerida:** agregar «Archivos» y «Solicitudes fuera de alcance» (ruta, cambio, motivo, responsable) para (i) a (iv). En K46, citar `S4-presentacion-respaldo` como bloqueo por crear. En `S2-ventanas`, dejar la línea «Qué falta» del caso 2 con la cita `REVISION-S2-r3.md` l. 38–48 y PUERTA-2 §2.1, sin inventar opciones.

### H-r1-05 · BAJA · cifras sin signo de pesos
- **Ruta:** `docs/estado/bloqueos/U3-t2-acumulados.md`, entrada Q2-65, «Qué cambia»: «2,921.25 y 3,895.00».
- **Corrección sugerida:** «$2,921.25 y $3,895.00».

### H-r1-06 · BAJA · Q2-12: RESPUESTAS.md y la aplicación no listan lo mismo
- **Ruta:** `RESPUESTAS.md` («Los «sin recomendación» (Q2-48, Q2-50 a Q2-52, Q2-56 a Q2-58) quedan abiertos») frente a INDICE, K46 y `S2-flujo-quincenal`, que dejan abierta también Q2-12.
- **Qué pasa:** la hoja §8 da «C o A», dos opciones. Dejarla abierta cumple «no asumas», y el reporte se la pregunta a Axel (Propuestas). Pero RESPUESTAS.md da a entender que tiene respuesta.
- **Corrección sugerida:** el orquestador confirma Q2-12 con Axel o anota en RESPUESTAS.md que queda abierta. No hay nada que corregir en los bloqueos.

## Estado de los anteriores
Primera ronda: no hay hallazgos anteriores.
