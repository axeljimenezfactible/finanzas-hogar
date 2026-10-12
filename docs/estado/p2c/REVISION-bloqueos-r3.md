# Revisión independiente · bloqueos de la puerta 2 · ronda 3 de 3 (11-oct-2026, noche)

Revisor: `revisor-bloqueos-r3`. No escribí nada de lo revisado y no corregí nada. Base: HEAD 1ccb468 más el árbol de trabajo a las 23:01 (27 bloqueos tocados, 6 bloqueos nuevos, `INDICE.md`, `decisiones.md` §K46 a §K52 y los reportes `bloqueos-respuestas.md`, `correccion-bloqueos-r1.md` y `correccion-bloqueos-r2.md`). Comparo contra `docs/estado/puerta2/RESPUESTAS.md`, `docs/estado/puerta2/PUERTA-2.md` (§2.1, §3, §4 y §8), el protocolo §2 y §3 y el formato de las puertas 1 y 1b en `INDICE.md`. Las verificaciones las corrí yo con mis scripts (`datos/tmp/r3/s/v3.py` y `k3.py`) y en tres clones aislados en `datos/tmp/r3`: `a` = HEAD; `b` = HEAD más la carpeta de bloqueos y `decisiones.md` del árbol; `t` = HEAD más todo el árbol de trabajo. Al terminar los borré.

## Veredicto
**NO APROBADO: 1 alto, 0 medios y 4 bajos.**

La fidelidad está bien. Las 53 letras con respuesta son la recomendada de §8. Q2-11 = **B** y Q2-64 = **A** están escritas tal cual. Ninguna pregunta «sin recomendación» lleva respuesta. La tabla de K46 coincide con §8 en las 98 filas, sin diferencias, y K50 trae los valores de PUERTA-2 §4. Los dos medios de la ronda 2 quedaron cerrados:
- H-r2-02: el INDICE y K46 ya coinciden con el disco (75 bloqueos, 63 resueltos y 12 abiertos).
- H-r2-03: la cita de Q2-12 ahora va por sección.

Queda el alto H-r2-01, que reabro como H-r3-01. La verificación global sigue rota. Hubo avance: el lado S1 ya pasa en el árbol. Falta el lado C1: `openapi.yaml`, `C1-datos.md` y el generador `c1gen` siguen sin aplicar Q2-11 = B. El arreglo no está en los bloqueos.

## Verificaciones (evidencia propia)
| # | Qué | Comando o script | Resultado |
|---|---|---|---|
| a1 | Archivos tocados y resueltos intactos | `git diff --numstat -- docs/estado/bloqueos/`; `v3.py` compara cada archivo con `git show HEAD:` | En HEAD hay 69 bloqueos. Hay 27 modificados y todos solo agregan texto: el texto de HEAD es prefijo del nuevo, con 0 líneas borradas. Los 27 son exactamente los de HEAD sin «Respuesta de Axel», sin contar `criterios-nxx`. Ninguno de los 42 resueltos (41 con «Respuesta de Axel» más `criterios-nxx`) cambió. `decisiones.md` tiene 37 líneas nuevas y 0 borradas. **OK** |
| a2 | Bloqueos nuevos en disco | `v3.py` | Hay 6 nuevos. P2c abrió `S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo` y `S5-carritos-granel`; la Pista 3 abrió `P3-contraste-chip-hover` y `P3-huecos-componentes`. Los 4 de P2c traen Contexto, Qué falta, Opciones, Impacto y Cómo seguí (protocolo §2), están en «Estado: ABIERTO» y no traen recomendación. **OK** |
| b1 | Letra de cada entrada contra §8 | `v3.py`: 62 entradas `### Q2-NN` en los 27 archivos | No faltan entradas de Q2-11 a Q2-71 ni hay repetidas. Ninguna letra difiere de la recomendada. Ninguna «sin recomendación» trae «Respuesta:». En `S1-parametros-avisos` está Q2-11 = `B` y en `S5-terminos` está Q2-64 = `A`. **OK** |
| b2 | Texto literal contra el archivo de HEAD | `v3.py`: busca cada «Texto de la opción (literal del archivo)» en `git show HEAD:<archivo>` | Hay 53 entradas con texto literal y no falta ninguno. **OK** |
| b3 | Muestra de 10 contra PUERTA-2: Q2-12, 14, 15, 38, 47, 60, 62 y 70, sacadas con la semilla 20261012, más Q2-11 y Q2-64 forzadas | `grep '**Q2-NN**' PUERTA-2.md` y `awk` sobre la entrada del archivo | Las 10 coinciden en letra y en sentido. **Q2-11:** la B del archivo («En `app_setting` (claves `alerts.*`…)») es la B de PUERTA-2 §3.1 («En claves de `app_setting`…»). Va como «respuesta explícita de Axel» y cuenta la aclaración de la «C · Mixta», que está en RESPUESTAS.md. **Q2-64:** la A del archivo es la A de §3.21. El archivo anota que la decisión es RESUELTO y que el bloqueo sigue ABIERTO, porque los 3 sitios respondieron con un desafío anti-robot. **Q2-12:** no tiene respuesta, va como PROVISIONAL `null` y Estado ABIERTO, y explica que la recomendada de §8 es «C o A». **Q2-62:** B, que es la recomendada. En Q2-14, 15, 38, 47, 60 y 70, la letra y el texto son los de §3. **OK** |
| b4 | Las «sin recomendación» no llevan respuesta | `v3.py` y lectura de `S5-carritos-detalle` y `S2-ventanas` | Q2-12, 48, 50, 51, 52, 56, 57 y 58, y la pregunta 4 de `S2-ventanas` (Q2-01, casos 1 y 2), traen «PROVISIONAL: sigue en uso…», «Pasa a WF-3» y **Estado: ABIERTO**. **OK** |
| c1 | Conteos del INDICE contra el disco | `v3.py`. Regla de la puerta 2: cuenta la última línea «Estado del bloqueo (puerta 2)»; si falta, vale «Respuesta de Axel»; si no, el «Estado» del archivo | Hay **75 archivos: 63 RESUELTOS y 12 ABIERTOS**. Las 3 líneas que mi regex no leyó las revisé a mano: `criterios-nxx` está RESUELTO y los 2 de P3 dicen «Estado: ABIERTO» en el encabezado. La tabla del Resumen (34+8+1+26+4+2) y la l. 16 dan 75, 63 y 12. **OK** |
| c2 | Regla de las herramientas (`verify_c1.py` l. 322-323) | `v3.py` | Hay **67 resueltos y 8 abiertos**: `P3-contraste-chip-hover`, `P3-huecos-componentes`, `S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo`, `S5-carritos-detalle`, `S5-carritos-granel` y `S5-profeco-filtro`. Es lo que dice la l. 30 del INDICE, que ya no dice «por crear». **OK** |
| c3 | Tabla «Respuestas de la puerta 2» del INDICE contra los archivos | `k3.py`: 62 filas con bloqueo | Letra y estado: 0 diferencias. Q2-64 dice «RESUELTO (decisión); bloqueo ABIERTO», igual que el archivo. La fila Q2-01 dice «casos 1 y 2». **OK** |
| c4 | Sección «Bloqueos nuevos de P2c» contra los 4 archivos | lectura de `S2-ahorro-quincenal` y `S4-presentacion-respaldo`, y `grep` de los otros dos | Pregunta, opciones (P1–P3 A/B; i, ii y iii; A, B y C) y lo PROVISIONAL coinciden con los archivos. Ninguno trae recomendación. **OK** |
| d1 | Tabla de K46 contra §8 | `k3.py` | Tiene 98 filas y 0 diferencias. Q2-11 = B y Q2-64 = A. Hay 15 «—»: Q2-12, 48, 50 a 52, 56 a 58, 85, 87 a 91 y 94. La suma da 83 + 15 = 98, como dice K46. **OK** |
| d2 | Citas de Axel en §K46 a §K52 | `k3.py` busca el texto entre comillas y lo compara con RESPUESTAS.md y PUERTA-2 | Hay 12 textos entre comillas. Los 3 que se atribuyen a Axel («lo recomendado + excepciones», «C · Mixta» y «nada de sitios de tiendas») están en RESPUESTAS.md; el primero aparece con mayúscula inicial. Los demás son términos de PUERTA-2 («rebasado», «Atajos de una tecla», «Guardado», «falta por comprar», «C o A») o del reporte de S5-terminos («no verificado»). No hay citas inventadas. **OK** |
| d3 | K50 contra PUERTA-2 §4 (l. 291 a 302) | lectura lado a lado | Reportes: diario 21:30, semanal domingo 19:00 y mensual día 1 a las 08:00. Avisos: 1 día y 09:00. Comparador: 5 s, con límites de 1 a 30 s. Ventana del optimizador: 7 días. OCR: 30 s. Frescura: S ≥ 1, Z ≥ S y (Z − días) ÷ (Z − S + 1). Son iguales. **OK** |
| e1 | Cita de Q2-12 (H-r2-03) | `grep -n` en `S2-interfaz.md` | El archivo está en la versión 1.3 (l. 3). §3.11 ocupa las l. 236 a 270 y su fila `variable_estimated_cents` = `null` PROVISIONAL está en la l. 258. §8 empieza en la l. 457 y su fila «Variable estimado de M-10» está en la l. 460. La entrada cita por sección. **OK** |
| e2 | Verificación con la carpeta de bloqueos (clones `a` y `b`) | `make -C a\|b verify-c1` y `verify-s1` con `GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **`a`:** verify-c1 da «RESULTADO: TODO OK» y verify-s1 «RESUMEN TODO OK», con rc=0. **`b`:** verify-c1 da «FALLA C1-15 … (sin «Respuesta de Axel»)», con 6 marcas en `Notification/properties/params_json…`, y rc=2. verify-s1 da «S1-16 FALLA … líneas con PROVISIONAL: 19» y `RESUMEN FALLAS: ['S1-16', 'K44']`, con rc=2. → **H-r3-01** |
| e3 | Verificación global en el árbol (clon `t`) | `make -C t verify-all GOOSE=…` | rc=2. Pasa WF-0 completo (SQL, Go, colorimetría y OpenAPI) y se detiene en verify-c1 con «gen2: tipo de parámetro desconocido: hora HH:MM». Con `make -i verify-c1`, `verify_c1.py` da «FALLA K44» y «FALLA C1-15» (6 marcas), con «RESULTADO: 2 falla(s)». verify-s1 da «RESUMEN TODO OK», verify-u1 «RESULTADO: TODO OK» y verify-criterios rc=0. → **H-r3-01** |
| e4 | Marcas que todavía citan `S1-parametros-avisos` | `grep -c 'PROVISIONAL.*S1-parametros-avisos'` y `stat` | Siguen 6 en `openapi.yaml`, desde la l. 17911. En `C1-datos.md` siguen las l. 572, 576, 772 a 775 y 779. Esos dos archivos tienen fecha de las 16:59, anterior a P2c. Siguen 5 en `00003_puerta1b.sql`, que `s1_verifica.py` ya trata como historia por el `-- @resuelve` de la 00004. El diccionario y RN-29 ya no tienen marcas. En los criterios: `U4.md` l. 6 dice «la excepción del bloqueo abierto `S1-parametros-avisos`», y siguen las de `C1.md` (C1-22) y `S7.md`. |

## Hallazgos

### H-r3-01 · ALTA · la verificación global sigue rota (reabre H-r2-01)
- **Ruta:** `docs/estado/bloqueos/S1-parametros-avisos.md` trae «Respuesta de Axel», Q2-11 = B, y está bien. Choca con:
  - `docs/contratos/openapi.yaml` (6 marcas, l. 17911 a 17960);
  - `docs/contratos/C1-datos.md` (P-26.08 en la l. 572, P-26.09 en la l. 576, filas de `params_json` en las l. 772 a 775 y la l. 779);
  - `backend/internal/api/c1gen/` (`gen2_schemas.py` `_pv`, que no conoce «hora HH:MM»);
  - los criterios `docs/criterios/U4.md` l. 6, `C1.md` (C1-22) y `S7.md`.
  Las reglas que lo detectan están en `backend/internal/api/c1gen/verify_c1.py` l. 322-323 (C1-15) y en la comprobación K44.
- **Qué pasa:**
  - En el clon de HEAD más los bloqueos (e2), verify-c1 falla C1-15 y verify-s1 falla S1-16 y K44. HEAD solo pasa.
  - En el árbol completo (e3), `make verify-all` se detiene en verify-c1 por «gen2: tipo de parámetro desconocido: hora HH:MM». Si se salta ese error, quedan FALLA K44 y FALLA C1-15.
  - Avance desde la ronda 2: el lado S1 ya pasa en el árbol. La `00004_puerta2.sql` trae `-- @resuelve S1-parametros-avisos`, `s1_verifica.py` ya trata las marcas de la 00003 como historia, y el diccionario y RN-29 ya no tienen marcas.
  - Lo que falta es C1: `openapi.yaml` y `C1-datos.md` no han cambiado y el generador no conoce el tipo nuevo.
  - El criterio de cierre de la r1 («`make verify-all` pasa con los dos cambios juntos») no se cumple. Lo declaran `correccion-bloqueos-r2.md` (solicitud 1) y `S1-puerta2.md` (solicitud 1, l. 102 en adelante).
- **Corrección sugerida:** en los bloqueos no hay nada que cambiar; quitar «Respuesta de Axel» escondería la decisión de Axel. Hay que aplicar la solicitud 1 de `docs/estado/p2c/S1-puerta2.md`, que es la regeneración de C1 con `c1gen`:
  - agregar el tipo «hora HH:MM» a `gen2_schemas.py` `_pv`;
  - exponer las claves nuevas de `app_setting` por las operaciones de settings;
  - quitar las 6 marcas de `openapi.yaml` y las de `C1-datos.md` (P-26.08, P-26.09 y §11.7);
  - corregir K44 en `verify_c1.py` o en el contrato, según lo que pida;
  - quitar la mención «bloqueo abierto» de `U4.md` l. 6 y las marcas de C1-22 y S7.

  El integrador no hace el commit «[Puerta 2]» de la carpeta de bloqueos sin ese cambio. Se cierra cuando `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` pasa con todo el árbol junto.

### H-r3-02 · BAJA · en `S2-ventanas` la pregunta 4 se describe distinto según el lugar (sigue H-r2-04)
- **Ruta:** `docs/estado/bloqueos/S2-ventanas.md` l. 83 y l. 85; `docs/plan/decisiones.md` K46, fila «01 | A (caso 2: —)».
- **Qué pasa:**
  - El encabezado de la entrada (l. 83) solo nombra «caso 2», y la l. 85 («Qué falta … No se escriben opciones aquí: las redacta quien aplique Q2-01») sigue obsoleta. Pero la l. 88, la sección «Pregunta 4» (l. 90 en adelante), la fila Q2-01 del INDICE y la prosa de K46 dicen bien que quedan abiertos los casos 1 y 2.
  - La celda de la tabla de K46 copia §8 tal cual, así que no está mal; solo no menciona el caso 1.
  - Nadie tiene que adivinar nada, porque la l. 88 y la sección de abajo lo dicen.
- **Corrección sugerida:**
  - l. 83: «### Q2-01 (casos 1 y 2, pregunta 4 del bloqueo)».
  - l. 85: «Qué falta: la pregunta 4 de abajo (casos 1 y 2, sin recomendación)».
  - Opcional, en la celda de K46: «A (pregunta 4, casos 1 y 2: —)».

### H-r3-03 · BAJA · cifras sin signo de pesos (sigue H-r2-05)
- **Ruta:** `docs/estado/bloqueos/U3-t2-acumulados.md` l. 24, entrada Q2-65: «nov-2026: 2,921.25 y 3,895.00; mar-2027: 13,103.78 …».
- **Corrección sugerida:** en el texto nuevo, escribir «$2,921.25 y $3,895.00», y lo mismo en las demás cifras de esa línea.

### H-r3-04 · BAJA · RESPUESTAS.md no lista Q2-12 como abierta (sigue H-r2-06)
- **Ruta:** `docs/estado/puerta2/RESPUESTAS.md`, «Cómo se aplica», viñeta de Q2-12 a Q2-71. Choca con el INDICE, K46 y `S2-flujo-quincenal`, que la dejan abierta.
- **Corrección sugerida:** el orquestador la confirma con Axel o anota en RESPUESTAS.md que queda abierta. No es del alcance de bloqueos.

### H-r3-05 · BAJA · la tabla de la puerta 2 del INDICE no sigue el formato de las puertas 1 y 1b (sigue H-r2-07)
- **Ruta:** `docs/estado/bloqueos/INDICE.md`, «Respuestas de la puerta 2» (l. 420 en adelante).
- **Qué pasa:** la puerta 1 da la letra y un resumen de lo decidido, y la 1b da la decisión y su K. La puerta 2 da solo la letra y el origen.
- **Corrección sugerida:** agregar las columnas «Decisión (resumen)» y «K» (K46, K47 o K48).

## Estado de los anteriores
| Hallazgo r2 | Estado | Evidencia |
|---|---|---|
| H-r2-01 (alta) | **Reabierto como H-r3-01**, con avance parcial | e2 y e3: el clon `b` falla C1-15, S1-16 y K44; el árbol se detiene en verify-c1 por `gen2` y luego falla K44 y C1-15. verify-s1 ya pasa en el árbol. `openapi.yaml` y `C1-datos.md` siguen sin cambios (16:59) |
| H-r2-02 (media) | **Cerrado** | c1 a c4: el disco tiene 75 archivos, 63 RESUELTOS y 12 ABIERTOS, igual que el INDICE (l. 12, l. 16, tabla y l. 30). Las herramientas cuentan 67 y 8, como dice la l. 30. Ya existen las secciones «Bloqueos nuevos de P2c» y «de la Pista 3». K46 ya no dice «por crear» y nombra los casos 1 y 2 |
| H-r2-03 (media) | **Cerrado** | e1: la cita va por sección y las dos filas existen en `S2-interfaz.md` 1.3 (l. 258 y l. 460) |
| H-r2-04 (baja) | Sigue abierto (H-r3-02) | `S2-ventanas.md` l. 85 sin cambio |
| H-r2-05 (baja) | Sigue abierto (H-r3-03) | `U3-t2-acumulados.md` l. 24 sin cambio |
| H-r2-06 (baja) | Sigue abierto (H-r3-04) | Le toca al orquestador |
| H-r2-07 (baja) | Sigue abierto (H-r3-05) | La tabla no cambió |
