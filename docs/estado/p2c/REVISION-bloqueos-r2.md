# Revisión independiente · bloqueos de la puerta 2 · ronda 2 de 3 (11-oct-2026, noche)

Revisor: `revisor-bloqueos-r2`. No escribí nada de lo revisado y no corregí nada. Base: HEAD 1ccb468 más el árbol de trabajo a las 22:51 (27 bloqueos tocados, `INDICE.md`, `decisiones.md` §K46 a §K52, reportes `docs/estado/p2c/bloqueos-respuestas.md` y `docs/estado/p2c/correccion-bloqueos-r1.md`). Mido contra `docs/estado/puerta2/RESPUESTAS.md`, `docs/estado/puerta2/PUERTA-2.md` (§2.1, §3, §4 y §8), el formato de bloqueo del protocolo (`docs/plan/protocolo-agentes.md` §2; alcance §3) y el formato de las respuestas de las puertas 1 y 1b en `INDICE.md`. Mis scripts y dos clones locales del repositorio (`a` = HEAD; `b` = HEAD más los 27 bloqueos, `INDICE.md` y `decisiones.md` del árbol) estuvieron en `datos/tmp/rev-bloq-r2/` y los borré al final.

## Veredicto
**NO APROBADO: 1 alto, 2 medios y 4 bajos.**
La fidelidad sigue bien: las 53 letras con respuesta son la recomendada de la hoja §8, Q2-11 = B y Q2-64 = A coinciden al pie de la letra, ningún «sin recomendación» lleva respuesta, la tabla de K46 tiene 0 discrepancias con §8 y K50 trae los valores de PUERTA-2 §4. Se cerraron H-r1-02 y H-r1-04. Lo que falla:
- La verificación global sigue rota con los bloqueos del árbol. Ya está declarada y pedida, pero falta aplicar la solicitud 1 (H-r2-01).
- `INDICE.md` y K46 ya no coinciden con el disco, porque otro agente abrió 3 bloqueos después de la corrección (H-r2-02).
- La cita de Q2-12 volvió a quedar mal porque `S2-interfaz.md` se reescribió en el árbol (H-r2-03).

## Verificaciones (evidencia propia)
| # | Qué | Comando o script | Resultado |
|---|---|---|---|
| a1 | Archivos tocados y resueltos intactos | `git status`, `git diff --stat HEAD`; `v.py` compara con `git show HEAD:<archivo>` | 69 archivos en HEAD: 27 modificados, y los 27 solo agregan texto al final (el texto de HEAD es su prefijo). Los 27 son exactamente los abiertos de HEAD. Hay 42 resueltos en HEAD (41 con «Respuesta de Axel» más `criterios-nxx`) y ninguno está modificado. Sin seguimiento: `S2-ahorro-quincenal.md`, `S2-metas-vencida.md` y `S4-presentacion-respaldo.md`, los 3 de `contratos-puerta2` (22:49 y 22:51), ver H-r2-02. **OK** |
| b1 | Letra de cada entrada contra §8 | `v.py`: 62 entradas (Q2-11 a Q2-71 y el caso 2 de Q2-01) | 0 faltantes y 0 repetidas. Hay 53 con letra, todas iguales a la recomendada. La única diferencia de forma es Q2-32, «A (propuesta del agente, letras de PUERTA-2)», que es la A. Q2-11 = **B** y Q2-64 = **A**, las dos como «respuesta explícita de Axel». **OK** |
| b2 | Texto literal contra el archivo de HEAD | `v.py`: busca cada «Texto de la opción (literal del archivo)» en `git show HEAD:` | Aparecen 53 de 53. **OK** |
| b3 | Muestra de 10 contra el texto de PUERTA-2 (Q2-11, 13, 14, 28, 32, 33, 42, 64, 70, 71; semilla 20261011, más Q2-11 y Q2-64 forzadas) | `m.py`: imprime la línea de PUERTA-2 y la entrada | Letra y sentido coinciden en las 10. En Q2-11, la B del archivo («En `app_setting`…») es la B de PUERTA-2 §3.1. En Q2-64, la A del archivo es la A de PUERTA-2 §3.21. Q2-13 anota la diferencia de redacción con PUERTA-2. **OK** |
| b4 | Las «sin recomendación» no llevan respuesta | `v.py` y `grep` de la línea PROVISIONAL | Q2-12, 48, 50, 51, 52, 56, 57, 58 y el caso 2 de Q2-01 no tienen «Respuesta:». Cada una trae «PROVISIONAL: sigue en uso A…» (Q2-12: `null`), «Pasa a WF-3» y **Estado: ABIERTO**. **OK** |
| c1 | Conteos del INDICE contra el disco | `v.py` (última línea «Estado del bloqueo (puerta 2)»; si falta, «Respuesta de Axel» o `criterios-nxx`; si no, «**Estado:**») | Los 69 de HEAD: 63 RESUELTOS y 6 ABIERTOS, igual que el INDICE. **Carpeta real: 72 archivos, 63 RESUELTOS y 9 ABIERTOS.** Con la regla de los scripts salen 67 y 5, y el INDICE l. 28 dice 67 y 2 → **H-r2-02** |
| c2 | Tabla «Respuestas de la puerta 2» contra los archivos | `v.py` | Hay 68 filas (Q2-01 caso 2, Q2-11 a Q2-71, Q2-72 a Q2-77). Letra y estado: 0 diferencias. **OK** |
| c3 | Fila «Flujo WF-1» y regla l. 26 (H-r1-02) | `grep -l 'Respuesta de Axel (puerta 2'` y `comm` | La sección está en 25 bloqueos: 1 de P1c y 24 de WF-1 (el grep da 26 porque cuenta también `INDICE.md`). De los 24, 20 están RESUELTOS y 4 ABIERTOS, como dice la fila. La l. 26 ya dice que es la regla de las puertas 1 y 1b. **OK** |
| d1 | Tabla de K46 contra §8 | `k.py` | 98 filas y 0 diferencias. Hay 15 sin respuesta (Q2-12, 48, 50–52, 56–58, 85, 87–91 y 94), y Q2-01 y Q2-07 traen su parte «—». **OK** |
| d2 | Citas de Axel en §K46–§K52 | `k.py` (texto entre comillas junto a «Axel») y lectura | Hay 3 citas: «lo recomendado + excepciones», «C · Mixta» y «nada de sitios de tiendas». Las 3 están en RESPUESTAS.md. No hay citas inventadas. **OK** |
| d3 | K50 contra PUERTA-2 §4 (l. 291–302) | lectura | Reportes: 21:30, domingo 19:00 y día 1 a las 08:00. Avisos: 1 día y 09:00. Comparador: 5 s con límites de 1 a 30 s. Ventana del optimizador: 7 días. OCR: 30 s. Frescura: S ≥ 1, Z ≥ S y (Z − días) ÷ (Z − S + 1). Son iguales. **OK** |
| e1 | Verificación que lee `docs/estado/bloqueos/` | `make verify-c1` y `make verify-s1` en el clon `a`, en el clon `b` y en el árbol | `a`: verify-c1 «RESULTADO: TODO OK» y verify-s1 «RESUMEN TODO OK», rc=0. `b`: «FALLA C1-15 … (sin «Respuesta de Axel»)», «S1-16 FALLA … líneas con PROVISIONAL: 19 (migración 6, diccionario 7, reglas 6)» y «K44 FALLA», con `RESUMEN FALLAS: ['S1-16', 'K44']`, rc=2. Árbol: verify-c1 se detiene en «parse_schema: clave @param repetida» (la `00004_puerta2.sql` de otro agente, sin seguimiento); verify-s1 da `RESUMEN FALLAS: ['S1-16', 'NEG12', 'K44', 'PROG13']`. Hay 18 líneas PROVISIONAL que citan `S1-parametros-avisos` en openapi, C1-datos, migraciones, diccionario y reglas, y `openapi.yaml` y `C1-datos.md` no están modificados → **H-r2-01** |
| e2 | Otras verificaciones globales | `make verify-u1` y `make verify-criterios` en `a` y en `b`; `verify-criterios` también en el árbol | rc=0 en todas. **OK** |
| e3 | Citas con número de línea del texto nuevo | `grep 'l\. [0-9]'` sobre el diff | `REVISION-S2-r3.md` l. 38-48 y l. 45: correctas. `S2-interfaz.md` l. 250 y l. 447: correctas en HEAD pero no en el árbol → **H-r2-03** |
| e4 | Marcas que citan los 4 ABIERTOS con «Respuesta de Axel» | `grep PROVISIONAL` en openapi, C1-datos, migraciones, diccionario y reglas | Hay 0 marcas que citen `S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa` o `S5-terminos`. El riesgo de H-r1-02 sigue latente hasta que se aplique la solicitud 2, pero hoy no rompe nada. |

## Hallazgos

### H-r2-01 · ALTA · la verificación global sigue rota con los bloqueos del árbol (reabre H-r1-01)
- **Ruta:** `docs/estado/bloqueos/S1-parametros-avisos.md` (con «Respuesta de Axel», Q2-11 = B), que choca con `docs/contratos/openapi.yaml` (6 marcas), `docs/contratos/C1-datos.md`, `backend/db/migrations/00003_puerta1b.sql`, `docs/contratos/diccionario-datos.md` y `docs/contratos/reglas-negocio.md`. Las reglas que deciden están en `backend/internal/api/c1gen/verify_c1.py` l. 322-323 y en `backend/db/tools/s1_verifica.py` l. 242 y 575.
- **Qué pasa:** la corrección r1 ya declaró la falla y pidió el arreglo (solicitud 1 de `correccion-bloqueos-r1.md`). Eso cierra la parte «no declarada». La parte «rota» sigue (e1): HEAD pasa y HEAD más los bloqueos falla C1-15, S1-16 y K44. En el árbol, `openapi.yaml` y `C1-datos.md` siguen sin cambios. La aplicación de Q2-11 = B en el esquema va en curso (`00004_puerta2.sql`), pero hoy rompe además `parse_schema`, NEG12 y PROG13. El criterio de cierre de la r1 era «`make verify-all` pasa con los dos cambios juntos», y no se cumple. El reporte original (`bloqueos-respuestas.md`) todavía dice «No se corrió `make verify-all`: solo documentos…». Lo corrige el reporte de la corrección, no el original.
- **Corrección sugerida:** en los bloqueos no hay nada que cambiar; quitar «Respuesta de Axel» escondería una decisión de Axel. Hay que aplicar la solicitud 1: el agente que aplica Q2-11 = B quita las 18 marcas que citan `S1-parametros-avisos` (6 en `openapi.yaml` por `c1gen`, más `C1-datos.md` P-26.08/P-26.09, migración, diccionario y RN-29), y criterios C1-22, U4-12 y S7. El integrador no hace commit de la carpeta de bloqueos sin ese cambio. Se cierra cuando `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` pasa con los dos cambios juntos.

### H-r2-02 · MEDIA · `INDICE.md` y K46 no coinciden con el disco (nuevo)
- **Ruta:** `docs/estado/bloqueos/INDICE.md`: l. 12 (estado de la puerta 2), l. 16 (conteo al cierre de P2c), tabla del resumen, l. 28 («Está por crear el bloqueo `S4-presentacion-respaldo` … 69 pasa a 70 y los abiertos a 7»; «cuentan 67 resueltos y 2 abiertos») y sección «Respuestas de la puerta 2». También `docs/plan/decisiones.md` K46 («`S4-presentacion-respaldo`, **por crear** en P2c, que aún no existe y no está en el conteo de 69»).
- **Qué pasa:** la carpeta tiene **72** bloqueos: **63 RESUELTOS y 9 ABIERTOS**. Los 3 nuevos (`S2-ahorro-quincenal`, `S2-metas-vencida` y `S4-presentacion-respaldo`) los abrió `contratos-puerta2` a las 22:49 y a las 22:51, después de la corrección r1 (22:50). El INDICE no los nombra en ninguna parte y sigue con 69, 63 y 6. Con la regla de los scripts salen 67 y 5, no 67 y 2. K46 dice que `S4-presentacion-respaldo` no existe. Además, el «Estado del bloqueo» de `S2-ventanas` dice ahora que en la pregunta 4 quedan abiertos los casos 1 y 2 (1A y 2A como PROVISIONAL). El INDICE (l. 12 y la fila de Q2-01) y K46 solo nombran el caso 2. Quien lea el INDICE en WF-3 no ve 3 de las 9 preguntas abiertas.
- **Corrección sugerida:** hacerla cuando `contratos-puerta2` cierre, en la siguiente corrección de bloqueos (el INDICE está en su alcance) o por el integrador.
  - Conteo 72, 63 y 9, con una fila «Flujo P2c | 3 | 0 | 3».
  - Una sección «Bloqueos nuevos de P2c» con la pregunta, las opciones y lo que queda PROVISIONAL de cada uno, como «Bloqueos nuevos de WF-1».
  - Quitar «por crear» de l. 28 y de K46, y poner «67 y 5» para las herramientas.
  - En l. 12, la fila de Q2-01 y K46: la pregunta 4 de `S2-ventanas` queda abierta con sus dos casos.
  - `contratos-puerta2` debe dejar esto como solicitud en su reporte.

### H-r2-03 · MEDIA · la fuente del valor en uso de Q2-12 vuelve a apuntar a otras líneas (reabre H-r1-03)
- **Ruta:** `docs/estado/bloqueos/S2-flujo-quincenal.md`, entrada Q2-12 («`docs/contratos/motores/S2-interfaz.md` l. 250 y l. 447»).
- **Qué pasa:** la cita era correcta contra HEAD. Pero `contratos-puerta2` reescribió `S2-interfaz.md` en el árbol (versión 1.3, 22:49), y ese archivo entra en el mismo commit. Ahora la l. 250 es `| credit_payments.household_cents | by_person[holder] |` y la l. 447 es la fila de `TestInvarianteDeudaViva`. El valor en uso está en la l. 258 (§3.11, `variable_estimated_cents` = `null`, PROVISIONAL, Q2-12) y en la l. 460 (§8, «Variable estimado de M-10»).
- **Corrección sugerida:** citar por sección y no por línea, porque el archivo sigue cambiando: «`docs/contratos/motores/S2-interfaz.md` 1.3 §3.11 (fila `variable_estimated_cents`) y §8 (fila «Variable estimado de M-10»)». Otra opción es actualizar a l. 258 y l. 460 cuando `S2-interfaz.md` quede final.

### H-r2-04 · BAJA · «Qué falta» del caso 2 en `S2-ventanas` quedó obsoleto
- **Ruta:** `docs/estado/bloqueos/S2-ventanas.md` l. 85.
- **Qué pasa:** dice «No se escriben opciones aquí: las redacta quien aplique Q2-01 … cuando exista, esta línea se sustituye por la pregunta completa». La pregunta 4 ya está escrita abajo (l. 90 en adelante: opciones 1A–1D y 2A–2C, impacto y «Cómo seguí») y la línea no se sustituyó. La línea PROVISIONAL de arriba sí la remite.
- **Corrección sugerida:** «Qué falta: la pregunta 4 de abajo (casos 1 y 2, sin recomendación)».

### H-r2-05 · BAJA · cifras sin signo de pesos (sigue H-r1-05)
- **Ruta:** `docs/estado/bloqueos/U3-t2-acumulados.md` l. 24 y l. 25 (entrada Q2-65): «2,921.25 y 3,895.00». Vienen de la opción A del original (l. 10), que tampoco lo trae.
- **Corrección sugerida:** en el texto nuevo, «$2,921.25 y $3,895.00».

### H-r2-06 · BAJA · Q2-12: RESPUESTAS.md no la lista como abierta (sigue H-r1-06)
- **Ruta:** `docs/estado/puerta2/RESPUESTAS.md` l. 18 frente a INDICE, K46 y `S2-flujo-quincenal`.
- **Corrección sugerida:** el orquestador la confirma con Axel o anota en RESPUESTAS.md que queda abierta. No es del alcance de bloqueos.

### H-r2-07 · BAJA · la tabla de la puerta 2 del INDICE no sigue el formato de las puertas 1 y 1b
- **Ruta:** `docs/estado/bloqueos/INDICE.md`, «Respuestas de la puerta 2» (l. 396 en adelante).
- **Qué pasa:** la puerta 1 da la letra y un resumen de lo decidido («A: lista compartida en el servidor»), y la 1b da la decisión y su K. La de la puerta 2 da solo la letra («C», «(a)», «1»). Para saber qué se decidió hay que abrir cada archivo, y la fila no dice qué K la consolida (K46, K47 o K48).
- **Corrección sugerida:** agregar una columna «Decisión (resumen)» de una línea y otra «K».

## Estado de los anteriores
| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H-r1-01 (alta) | **Reabierto como H-r2-01.** La parte «no declarada» quedó cerrada; la verificación sigue rota | e1: el clon `b` falla C1-15, S1-16 y K44, y el clon `a` pasa. La solicitud 1 está en `correccion-bloqueos-r1.md` y no se ha aplicado |
| H-r1-02 (media) | **Cerrado** | c3: la l. 26 está enmendada, la l. 28 explica las dos reglas y la fila dice «24 con la sección, 20 RESUELTOS», cifra que comprobé. La solicitud 2 a los scripts está pedida. Los números de la l. 28 cambiaron por los bloqueos nuevos (H-r2-02) |
| H-r1-03 (media) | **Reabierto como H-r2-03** | Correcto contra HEAD, mal contra el `S2-interfaz.md` 1.3 del árbol |
| H-r1-04 (media) | **Cerrado** | (i) La línea «Qué falta» está, y además la pregunta 4 ya existe completa (obsoleta en lo de detalle: H-r2-04). (ii) K46 cita `S4-presentacion-respaldo` (que ya existe: H-r2-02). (iii) `correccion-bloqueos-r1.md` trae «Archivos» y «Solicitudes fuera de alcance» (1 a 6) |
| H-r1-05 (baja) | Sigue abierto (H-r2-05) | El reporte de la corrección lo deja como pendiente |
| H-r1-06 (baja) | Sigue abierto (H-r2-06) | Es del orquestador |
