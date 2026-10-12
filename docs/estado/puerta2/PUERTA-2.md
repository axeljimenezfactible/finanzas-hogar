# Puerta 2 · Finanzas Hogar · lo que Axel decide

Redactor: `redactor-puerta2` (Opus) · 11-oct-2026, noche · commit base **0926802** («[WF-1] Motores y sistema de diseño…») · repositorio `/home/deepsec/finanzas-hogar`.

## 0. Cómo responder

**Contesta por id con la letra de la opción (p. ej. «Q2-05: B»); «recomendada» es la opción que propusieron los agentes, con su cita; puedes contestar por lotes: «todo lo recomendado salvo Q2-05: B y Q2-48: A».**

- Hay **98 preguntas**, de Q2-01 a Q2-98, sin huecos. La hoja de respuestas de la §8 las junta todas.
- Las letras son las del archivo de origen. Cuando el archivo no pone letras, lo digo y las letras son de este documento.
- Etiquetas: **HECHO** (ya pasó o está escrito), **PREGUNTA** (decides tú), **PROPUESTA** (sugerencia de un agente o del orquestador, sin valor hasta que la apruebes).
- Donde dice «sin recomendación», ningún reporte propuso una opción, o los agentes no se pusieron de acuerdo. **PROPUESTA del redactor:** si contestas «todo lo recomendado» y no tocas una de esas, se queda abierta, sigue PROVISIONAL con la opción que ya está en uso y pasa a WF-3. No hace falta contestarlas todas hoy.
- Las cifras salen de los archivos citados. Este documento no agrega valores nuevos.

## 1. HECHO

**Lo que cerró WF-1** (`docs/estado/wf1/INTEGRACION.md` §0 a §2):
- Los 4 motores (S2 finanzas, S3 créditos, S4 recetas y nutrición, S5 precios) y sus oráculos reproducen sus fixtures y `oraculo/model.json` con **0 diferencias** (596 de 596 cifras canon). `model.json` y `verify3.json` no cambiaron.
- `make verify-all` termina en «Verificación ampliada OK». La cobertura de pruebas es de 85 % o más (K25) en `motor_fin`, `motor_cred`, `recetas` y `precios`.
- Quedaron aprobados los contratos S2 y S3, el motor S3, el motor S4 con su recetario y nutrición, el catálogo de 160 productos, el sistema de diseño (DESIGN.md, tokens C2 y contraste, revisión U2 en su ronda 4) y U3 Inicio.
- Quedan **6 artefactos sin aprobar** (sección 2) y **27 bloqueos abiertos**: 26 de WF-1 y `S1-parametros-avisos`, que viene de P1c (sección 3).
- La base local de PROFECO tiene 37,548,200 registros. El OCR está listo pero **sin medir**: faltan tus 5 fotos de tickets (K24).

**Lo que ya decidiste hoy: K45, que sustituye a K35** (`docs/estado/puerta2/k45-borrador.md`):
- La Pista 3 (Claude Design) la hace el orquestador con Artifacts de tipo **Design System** y **Design**. Nadie usa navegador ni sube archivos, y el paquete `docs/design/entrada/` pasa a ser solo un **brief de referencia**.
- Las maquetas llevan **cifras sintéticas** y cubren **todas las pantallas de U1** (unas 25) en móvil 412 y escritorio 1440, en claro y oscuro, más la barra lateral, el menú `/` y los 6 flujos clave.
- **Tú revisas solo la versión final.** Los agentes revisan las vueltas, con un máximo de 3, y la puerta 3 consiste en abrir la URL del lienzo y aprobar.
- HECHO, visto al redactar: el árbol de trabajo ya trae K45 aplicada **sin commit**. Están modificados `docs/plan/decisiones.md` (K45, y K35 marcada como sustituida), `docs/plan/plan.md`, `docs/criterios/P3-claude-design.md`, `docs/design/entrada/LEEME.md` y `MANIFIESTO.md`. No los toqué.

## 2. PREGUNTAS · los 6 artefactos sin aprobar

Fuente: `INTEGRACION.md` §4. Ningún hallazgo cambia una cifra canon. Son reglas que el contrato no escribe o textos incompletos. Todas las preguntas de esta sección son: **«¿Autorizas aplicar esa corrección en la ronda de aplicación de la puerta 2 (P2c)?»**, con estas opciones: **A** aplicar la corrección sugerida · **B** otra (dices cuál) · **C** dejarlo PROVISIONAL para WF-3. En todas, la recomendada es A, que es la corrección del revisor, salvo en Q2-10.

### 2.1 Motor S2 (`motor_fin`) · `docs/estado/wf1/REVISION-S2-r3.md`

**M1 (media) · S2-16 / S2-03.** Puedes pedir el resumen por periodos que tú defines («ventanas»). Con ventanas, un gasto que se reparte en n meses (devengo) se reparte en «n ventanas consecutivas» de la lista, no en n meses (l. 39). El revisor encontró **dos casos** con esa misma causa. Motor y oráculo coinciden en las cifras, pero la regla no está escrita en ningún lado (`REVISION-S2-r3.md` l. 38–45; el informe da las cifras en centavos).
- **Caso 1 · hueco que cruza un devengo.** Si las ventanas dejan un hueco sin movimientos, la parte del hueco pasa a la ventana siguiente. En modo meses, esa misma parte cae en el mes del hueco. Ejemplo del revisor: un gasto de contado de $1,000.00 (100000) el 2026-10-10, devengado en 2 meses. Con las ventanas A (2026-10-10 a 11-09) y C (2026-12-10 a 2027-01-09), C da devengado $500.00 (50000). Con A, B y C contiguas, C da 0. En modo meses, diciembre da 0 (l. 40–43).
- **Caso 2 · ventanas contiguas que no duran un mes.** La diferencia acumulada al mismo día de corte depende de cómo se corten las ventanas. Con la ventana P1 sola da $973.75 (97375) al 2026-11-09. Con dos ventanas, d = 2026-10-10 y r = 2026-10-11 a 11-09, da $2,921.25 (292125) al mismo 2026-11-09, porque la segunda mitad de la despensa 1 cae en r. La interfaz (§3.2) define esa diferencia «desde el primer evento… hasta `end`» (l. 45).
- Corrección sugerida: (1) agregar a `S2-ventanas` una pregunta 4, que el revisor redacta así (l. 48): «con ventanas, ¿el devengo "mitad por mes" sigue a la ventana o al mes, y qué pasa con un hueco sin eventos que cruza un devengo y con ventanas que no duran un mes? Poner opciones con sus cifras»; (2) mientras tanto, aplicar en motor y oráculo una regla PROVISIONAL: rechazar con `invalid_argument` toda lista con hueco, o al menos todo hueco que cruce un devengo; (3) agregar la consulta de las ventanas A + C a la fixture `s2-ventanas-sin-cubrir`; (4) que §2.12 diga «contiguas, ascendentes y disjuntas (PROVISIONAL, S2-ventanas)». Si se prefiere conservar el comportamiento actual, hay que escribirlo como regla PROVISIONAL con sus cifras.
- **Q2-01** · ¿Autorizas aplicarla en P2c? A / B / C. **Recomendada: A**, la corrección del revisor: escribir la pregunta 4 completa, con los dos casos, sus opciones y sus cifras, y poner la regla PROVISIONAL. Las opciones de la pregunta 4 todavía no existen; las escribe quien aplique la corrección (1).
  - HECHO: lo que `INTEGRACION.md` §7.2 recomienda para la pregunta 4, «exigir ventanas contiguas, ascendentes y disjuntas» (`correccion-S2-r1`, `correccion-S2-r2`, `revisor-S2-r3`), **cubre solo el caso 1**. No cierra el caso 2: d y r son contiguas, ascendentes y disjuntas, y siguen dando $2,921.25 contra $973.75 (`REVISION-S2-r3.md` l. 45). `revisor-S2-r3.md` l. 51 lo acota: la regla hace coincidir ventanas y meses «cuando las ventanas son las del oráculo». Las reglas de ejemplo de la corrección (2) tampoco lo tocan, porque d y r no dejan hueco.
  - **Caso 2: sin recomendación.** Ningún reporte propone opciones para él. Contestar A no lo cierra: queda abierto dentro de la pregunta 4 del bloqueo hasta que contestes esa pregunta. Mientras tanto, motor y oráculo siguen con el comportamiento actual, que no está escrito.

### 2.2 Motor S5 (`precios`, asistida, OCR, PROFECO) · `docs/estado/wf1/REVISION-S5-r3.md`

**H-r3-00 (alta) · S5-07.** Falta `docs/estado/S5-terminos.md`, que es el veredicto sobre si los términos de Walmart, Bodega Aurrera y Sam's Club permiten la lectura asistida. Ningún agente puede escribirlo, y la lectura asistida sigue apagada. **Esta pregunta es la del bloqueo `S5-terminos`: Q2-64.** Corrección del revisor: con la opción A o B del bloqueo se escriben los 3 veredictos con URL y fecha; con la D, el integrador reescribe S5-07 y el plan §3.4.

**H-r3-01 (media).** El contrato `S5-interfaz.md` 1.4 está desfasado. Describe archivos y una base que ya cambiaron (§11), no define varios conteos y órdenes de la comparativa por carritos (§7.2 y §15.2) y no menciona los tickets sintéticos (§12). Con 3 carritos, una misma presentación que está en 2 de las 3 canastas sale como «sustituto».
- Corrección sugerida: el integrador aplica las solicitudes 1 a 6 de `docs/estado/wf1/correccion-S5-r2.md` (contrato 1.5) y agrega la regla para 3 canastas o más: «idéntico = la misma presentación en todas las canastas; lo demás es sustituto o solo en una tienda». Si se quiere otra regla, la pregunta se abre en `S5-carritos-detalle`.
- **Q2-02** · ¿Autorizas aplicarla en P2c? A / B / C. **Recomendada: A.**

**H-r3-02 (media).** Al comparar carritos con cantidades a granel (por ejemplo 0.5 kg), motor y oráculo difieren en 1 o 2 centavos: un total sale $1.01 en el motor y $1.02 en el oráculo, otro $1.03 contra $1.04. El oráculo redondea cada renglón y el motor redondea una sola vez. Hoy no cambia el canon, porque sus cantidades son enteras, pero podría voltear al ganador de un sustituto.
- Corrección sugerida: fijar en §7.2 cómo se redondea el importe de un renglón. Si no es una regla técnica, abrir el bloqueo `S5-carritos-granel` y marcarla PROVISIONAL. El precio por unidad base siempre sale de la fracción exacta. Después, alinear el oráculo y agregar una fixture a granel.
- **Q2-03** · ¿Autorizas aplicarla en P2c? A / B / C. **Recomendada: A.**

**H-r3-03 (media).** Al importar PROFECO, motor y oráculo leen distinto algunos casos raros: un precio de $0.00, un precio escrito «1,021.00», una cadena vacía, «ATUN» contra «ATÚN» y los textos de los motivos de rechazo. Los cortes reales no se afectan (0 rechazos en 37,548,200 registros).
- Corrección sugerida: escribir en §11 cómo se lee un precio (positivo, ¿con separador de miles?), la lista cerrada de motivos con su texto, qué pasa con la cadena vacía y si la liga se compara normalizada. Si no es técnico, va con PROVISIONAL y bloqueo. Después, alinear oráculo o motor y agregar una fixture con estos casos.
- **Q2-04** · ¿Autorizas aplicarla en P2c? A / B / C. **Recomendada: A.** Se relaciona con `S5-profeco-filtro` (Q2-56 a Q2-58).

**H-r3-04 (media).** «Cuánto más cara que la más barata» (`difference_vs_cheapest_cents`) se mide distinto según el modo. Al comparar carritos se mide contra el menor total aunque a esa tienda le falten productos; en la canasta se mide contra la primera tienda completa.
- Corrección sugerida: agregar una pregunta 4 al bloqueo `S5-carritos-detalle` (o decidirlo en §7.2), usar la misma referencia en los dos modos o documentar la diferencia, y fijarlo en la fixture.
- **Q2-05** · ¿Autorizas aplicarla en P2c? A / B / C. **Recomendada: A.** Se toca con M-r5-3 (Q2-08), que propone diferencia `null` para tiendas incompletas.

Además: S5-06 (OCR) queda «sin medir» porque faltan las 5 fotos de K24. Es una acción tuya y está en la sección 7.

### 2.3 Contrato S4/S5 (tras 5 rondas) · `docs/estado/wf1/REVISION-contrato-S4S5-r5.md`

**M-r5-1 (media) · S4-03, S4-06.** Un mismo producto puede tener varias elecciones en la lista (la del artículo calculado y la de uno o más artículos manuales). Las cotizaciones no dicen a qué elección responde cada una, así que con una entrada válida el motor puede dar error. Ejemplo del revisor: arroz con Walmart 1 kg a $28.00 y 750 g a $22.50, y Bodega 500 g a $13.50; debería regir Walmart 1 kg ($28.00) y sale `invalid_argument`.
- Corrección sugerida: que cada cotización diga a qué elección responde (`choice_store` y `choice_presentation`, `null` y `null` en la recomendada) y que el paso 1 tome la de exactamente esa elección. Otra forma: que `list_choices[]` traiga su propia cotización. Ajustar §12.1, el caso `s4-costo-cotizacion-elegida` y `TestCostoCotizacionQueRige`.
- **Q2-06** · ¿Autorizas aplicarla en P2c? A / B / C. **Recomendada: A.**

**M-r5-2 (media) · S4-03.** Cuando un producto aún no tiene precio (lo normal al arrancar), la regla «rige la presentación del ingrediente» no da una sola presentación en 3 casos: (a) el producto viene de varias recetas con presentación distinta; (b) el artículo solo está por reposición; (c) un artículo calculado con necesidad 0. Ejemplo: papel higiénico por ritmo de consumo sin precio. Quedan sin definir la presentación, los paquetes y el sobrante.
- Corrección sugerida: una regla única de respaldo. El revisor da 3 variantes sin marcar preferida: (i) la presentación típica del catálogo (`typical_presentation` en `catalog.products[]`); (ii) para el caso (a), la primera receta por fecha y tiempo; (iii) `null` en los tres campos, como el artículo manual. Agregar el caso a una fixture.
- **Q2-07** · ¿Autorizas aplicarla en P2c? A (di la variante i, ii o iii) / B / C. **Recomendada: A**, pero la variante queda **sin recomendación**.

**M-r5-3 (media) · S4-03, S5-03 y oráculos.** Los contratos no fijan el **orden** de varias listas de salida (artículos de la lista, faltantes de precio, casillas del plan, renglones de la comparativa, tiendas, etc.). Como las fixtures se comparan por posición, motor y oráculo pueden «diferir» aunque los dos calculen bien. Tampoco dice dónde va una tienda incompleta.
- Corrección sugerida: una línea de orden por lista. Por ejemplo: `items[]` por grupo, luego calculado antes que manual, luego código de producto; `missing_prices[]` por código y sin repetir; `cells[]` por fecha, tiempo y persona; `rows[]` con idénticos, luego sustitutos y luego «solo en una». En `stores[]`, primero las completas por total y después las incompletas, con diferencia `null` y `cheapest_store` `null` si no hay ninguna completa (PROVISIONAL, `S5-reglas-comparativa` pregunta 2).
- **Q2-08** · ¿Autorizas aplicarla en P2c? A / B / C. **Recomendada: A.**

### 2.4 U4 Ajustes (tras 5 rondas) · `docs/estado/wf1/REVISION-U4-r5.md`

**R5-1 (media) · U4-05.** La regla general de Ajustes dice que interruptores y listas de opciones «guardan al elegir». Pero hay campos de lista con vista previa que debería verse antes de guardar: el modo de reposición (P-26.14), el reparto de una categoría (P-26.04) y sexo, actividad y objetivo (P-26.13). Además, arrastrar una tarjeta en el tablero o editar varias filas con `Ctrl+/` no dice si muestra la vista previa. Quien programe tendría que adivinar.
- Corrección sugerida: volver general la regla de §0.2: (1) toda lista o interruptor con vista previa la muestra al elegir y guarda con `Enter` o al salir (`Esc` vuelve); (2) `Ctrl+/` y el arrastre entre columnas de un campo con vista previa la muestran antes de aplicar; (3) decir qué muestra cada caso (P-26.14: estado antes y después por fila; P-26.04: el ejemplo de $199.00 con el reparto nuevo). Otra opción: quitar el arrastre y `Ctrl+/` de esos campos.
- **Q2-09** · ¿Autorizas aplicarla en P2c? A / B / C. **Recomendada: A.** Siguen abiertos como bajos R5-2, R5-3, R4-2, R4-3 y R3-2 a R3-6, sin pregunta.

### 2.5 Paquete para Claude Design · `docs/estado/wf1/REVISION-paquete-r3.md`

**H-1 (media) · P3-03.** El paso 2 del LEEME («Crea un proyecto nuevo…») no decía qué elegir si la herramienta pedía tipo de proyecto o sistema de diseño.
- Corrección sugerida por el revisor (antes de K45): agregar una frase al paso 2 y regenerar `MANIFIESTO.md`.
- HECHO: con K45 el paquete ya no se sube. El LEEME del árbol de trabajo ya se reescribió como brief de referencia: el texto «Crea un proyecto nuevo» ya no aparece, y `MANIFIESTO.md` también cambió. Todo eso está sin commit.
- **Q2-10** · ¿Qué hacemos en P2c? A: aplicar la frase del revisor · B: dar H-1 por superado con el LEEME y el MANIFIESTO reescritos por K45, que P2c revisa y mete al commit · C: PROVISIONAL. **Recomendada: B.** Fuente: `docs/estado/puerta2/k45-borrador.md`, «Archivos del repo a actualizar», donde dice que los pasos 1–9 se sustituyen. La A quedó superada por K45.

## 3. PREGUNTAS · los 27 bloqueos abiertos

Van en el orden de `docs/estado/bloqueos/INDICE.md`: primero el de P1c y luego los 26 de WF-1. Abiertos son los archivos sin «Respuesta de Axel», sin contar `criterios-nxx`, que está resuelto. Las recomendaciones salen de `INTEGRACION.md` §7.2 y de los reportes que se citan. En cada bloqueo pongo qué sigue PROVISIONAL y qué pasa si no se contesta hoy.

### 3.1 `S1-parametros-avisos.md` · fase P1c (S1 esquema 1.3) · lo usan C1, U4 y S7
Choque entre K43, que pone la hora y la anticipación de los avisos en `notification.params_json`, y K44, que pone todo parámetro en `app_setting`. Los valores de `weekday` (7) y `day_of_month` (1) son FIRMES por la decisión F; lo que está en duda es **dónde** se guardan.
- **Q2-11** · ¿Dónde viven `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1)?
  - A. En `notification.params_json`, como dice K43 al pie de la letra. Los parámetros quedan en dos lugares y hay que registrar la excepción a K44.
  - B. En claves de `app_setting`, como los umbrales de 80 y 100 % y los días antes del pago. Todo queda en una tabla. Una migración siembra las 3 claves con los valores de F, las quita de `params_json` y su Down las repone.
  - C. Mixta: la hora y el día de cada aviso en su `params_json`, y la anticipación (`advance_days`) en `app_setting`.
  - **Sin recomendación única: los agentes no coinciden.** S1 recomienda **B**, por tener una sola fuente por valor (`docs/estado/p1c/correccion-S1-r2.md` l. 59). La corrección de criterios recomendó **A**, porque no cambia ni esquema ni contrato (`docs/estado/p1c/correccion-criterios-r1.md`, citado en INDICE). HECHO: el hallazgo A, alto, de `docs/estado/p1c/REVISION-criterios-r3.md` (los criterios describían los umbrales como si estuvieran en `params_json`) ya está cerrado. `REVISION-criterios-r4.md` l. 17 comprueba los umbrales como claves `alerts.*` y deja dos pendientes que la ronda 5 cierra (`REVISION-criterios-r5.md` l. 8 y 39); la r5 aprueba los criterios 1.2 con 0 altos y 0 medios (l. 6; commit f10d3bb). C1-22 (`docs/criterios/C1.md` l. 124) deja fuera los cuatro parámetros de esta pregunta y los marca PROVISIONAL mientras siga abierto este bloqueo.
- PROVISIONAL: solo el lugar de los cuatro. Hay marcas en `00003_puerta1b.sql`, `diccionario-datos.md`, RN-29, `openapi.yaml`, `C1-datos.md` y en los criterios C1-22, U4-12 y S7.
- Si no se contesta: U4 («Notificaciones» y «Reportes y frecuencias») sigue con dos variantes, S7 (WF-3) no sabe dónde leer la hora, y Q2-72 y Q2-73 no tienen dónde sembrarse.

### 3.2 `S2-flujo-quincenal.md` · fase S2
Ningún documento fija cómo se estima el gasto variable ni el ahorro de una quincena futura, ni cuándo una partida planificada ya se capturó.
- **Q2-12** · (1) ¿Cómo se estima el gasto variable de una quincena futura? A) presupuesto de las categorías variables del mes, en proporción a los días de la quincena · B) promedio diario del gasto variable real de los últimos 3 meses × días de la quincena · C) 0 hasta que Axel decida. **Recomendada: C, o A.** Por qué: A usa solo datos que ya capturas (`docs/estado/wf1/contrato-S2.md`, «Propuestas»; §7.2: «(1) C con 0 o A»).
- **Q2-13** · (2) ¿Cómo se reparte el ahorro mensual planificado entre las dos quincenas? A) todo en la quincena 1 · B) mitad y mitad, con el centavo de redondeo en la 1. **Recomendada: B.** Es simétrico con el sueldo quincenal (`contrato-S2.md`).
- **Q2-14** · (3) ¿Cuándo una partida planificada (sueldo, renta) se da por cumplida, para no contarla dos veces? A) hay un gasto o ingreso con la misma partida (`planned_item_id`) en la misma quincena, o en el mismo mes si es mensual o de mayor plazo · B) solo con la misma partida y la fecha exacta. **Recomendada: A** (`contrato-S2.md`; §7.2).
- PROVISIONAL: el gasto variable estimado, el ahorro planificado y la regla de cumplimiento. La interfaz usa 3A como lectura de trabajo. Las cifras canon (R-13 a R-31) no dependen de esto.
- Si no se contesta: el saldo de las quincenas futuras (M-02, M-10) se construye en WF-3 con marca PROVISIONAL.

### 3.3 `S2-metas-resumen.md` · fase S2
- **Q2-15** · ¿El resumen del mes (tarjeta «Ahorro») cuenta solo metas y fondos activos, o también los inactivos con saldo? A) solo activas, el mismo conjunto que la lista de metas · B) todas para «ahorrado» y «aportado del mes», y solo activas para «planificado». **Recomendada: A.** Así la tarjeta de Inicio y la pantalla de Metas suman lo mismo (`docs/estado/wf1/correccion-contrato-S2-r2.md` l. 63; §7.2).
- PROVISIONAL: `saved_cents` y `contributed_this_month_cents`. Las dos opciones dan lo mismo si no hay metas inactivas con movimientos.
- Si no se contesta: la tarjeta «Ahorro» queda PROVISIONAL en WF-3.

### 3.4 `S2-metas.md` · fase S2
- **Q2-16** · (1) ¿Cómo se cuentan los «meses que faltan» de una meta? A) del mes de hoy al mes de la fecha meta, los dos incluidos · B) del mes siguiente al de hoy hasta el de la fecha meta · C) meses completos entre hoy y la fecha meta. **Recomendada: A.** La aportación de este mes cuenta (`contrato-S2.md`; también `correccion-U4-r1`, según `correccion-U4-r2.md` l. 122). Contesta también la antigua pregunta 2 de `U4-vista-previa-formulas`.
- **Q2-17** · (2) ¿Qué formato tiene el porcentaje de avance? A) un decimal en todas partes (RN-18, «50.0 %») · B) entero en todas partes (C1) · C) un decimal en Metas y entero en la tarjeta de Inicio. **Recomendada: C** (`contrato-S2.md`).
- **Q2-18** · (3) ¿Qué estados puede tener una meta? A) en curso, cumplida, vencida · B) además «atrasada», cuando lo ahorrado va por debajo del plan acumulado. **Recomendada: A**, tres estados simples (`contrato-S2.md`).
- PROVISIONAL: meses que faltan, aportación necesaria, plan y estado, la fixture `s2-metas` y P-26.13 de U4. Pendiente también alinear RN-18 con M-07/M-16 cuando decidas.
- Si no se contesta: la pantalla de Metas sale con esas 4 cifras PROVISIONALES.

### 3.5 `S2-presupuesto-total.md` · fase S2
- **Q2-19** · (1) Si existe un renglón de «presupuesto total del mes», ¿el plan total es ese monto o la suma de las categorías? A) manda el renglón del total y la suma se muestra aparte · B) manda la suma de categorías (M-01) y el renglón solo se informa · C) el mayor de los dos. **Recomendada: B**, como dice M-01 (`contrato-S2.md`).
- **Q2-20** · (2) Con plan $0.00 en una categoría, ¿qué porcentaje y banderas se muestran? A) sin porcentaje ni banderas (estado «sin plan») · B) cualquier gasto mayor que 0 enciende las dos banderas. **Recomendada: A** (`contrato-S2.md`).
- PROVISIONAL: el plan total cuando existe el renglón del total, y el caso de plan 0. Las banderas de 80 y 100 % con plan mayor que 0 no dependen de esto.
- Si no se contesta: la interfaz sigue con B y A como PROVISIONAL.

### 3.6 `S2-proyeccion-interruptor.md` · fase S2
Con el interruptor de Ajustes encendido (K21), los meses 7 a 12 suman una despensa a crédito cada dos meses y el papel de $199.00 en los meses alternos. RN-28 lo da solo como ejemplo. **Las opciones A, B y C del archivo contestan las 4 preguntas a la vez:** A = repetir el patrón del oráculo; B = que salga de partidas planificadas que tú capturas (sin partidas no hay cadencia); C = A como valor inicial sembrado y B para cambiarlo.
- **Q2-21** · (1) Monto de cada despensa proyectada: A) $4,000.00 de capital, como las despensas 2 a 4 · B) el de la partida planificada que captures · C) A al inicio, B para cambiarlo. **Recomendada: C.**
- **Q2-22** · (2) Fecha: A) día 10 cada dos meses desde la última despensa (10-jun-2027, 10-ago-2027) · B) la fecha de inicio de la partida planificada · C) A al inicio, B para cambiarlo. **Recomendada: C.**
- **Q2-23** · (3) Reparto: A) las mismas partes que las despensas 2 a 4 (la parte fija de Karen del oráculo: $1,191.50 de capital y $1,499.72 con intereses) · B) lo que diga la partida planificada (el archivo no detalla el reparto en B; la pregunta menciona como alternativa la regla de la categoría Despensa) · C) A al inicio, B para cambiarlo. **Recomendada: C.**
- **Q2-24** · (4) Papel alterno: A) $199.00 de contado, pagado por Karen, el día 10 de mayo, julio y septiembre de 2027 · B) solo si se captura como partida · C) A al inicio, B para cambiarlo. **Recomendada: C.**
- Por qué C: conserva el ejemplo de RN-28 y cumple K44, que pide todo configurable (`contrato-S2.md` l. 59; §7.2).
- PROVISIONAL: con el interruptor encendido, los meses 7 a 12 y la fixture `s2-proyeccion-7-12-cadencia`. Apagado, que es el valor por omisión, no depende de esto.
- Si no se contesta: la proyección con interruptor encendido sale PROVISIONAL. La proyección por omisión no cambia.

### 3.7 `S2-ventanas.md` · fase S2
Ventanas son los periodos que defines tú, en lugar de meses calendario. La pregunta 4 (devengo con ventanas: hueco que cruza un devengo y ventanas que no duran un mes) es la Q2-01; ver §2.1, donde el caso de las ventanas que no duran un mes va sin recomendación.
- **Q2-25** · (1) ¿Qué pasa con un costo o pago anterior a la primera ventana pedida? A) rechazar con `invalid_argument`: quien pide tiene que declarar ventanas contiguas desde el primer movimiento · B) devengar lo anterior con ventanas implícitas hacia atrás (con P2 sola daría R-46) · C) contar solo desde la primera ventana (cambia el texto de §3.2, y P2 sola ya no reproduce R-46). **Recomendada: A.** La diferencia acumulada siempre queda completa y cuadra con el canon (§7.2; `correccion-S2-r1`, `correccion-S2-r2`, `revisor-S2-r3`).
- **Q2-26** · (2) ¿Y un costo o pago fechado en un hueco entre ventanas? H1) rechazar · H2) devengar el costo del hueco en n ventanas a partir de la siguiente · H3) que cuente en la diferencia acumulada, pero en ningún periodo. **Recomendada: H1** (§7.2).
- **Q2-27** · (3) ¿Y ventanas fuera de orden, traslapadas o una dentro de otra? O1) rechazar toda lista en la que una ventana no empiece después del fin de la anterior · O2) ordenarlas por inicio para calcular y devolver los resultados en el orden recibido (los traslapes se siguen rechazando). **Recomendada: O1** (§7.2).
- PROVISIONAL: A + H1 + O1 están aplicadas en `ledger.go` y `fin_devengo.py`, con la fixture `s2-ventanas-sin-cubrir`. No cambia ninguna cifra canon.
- Si no se contesta: siguen A + H1 + O1 como PROVISIONAL, y la interfaz no escribe la regla en §2.12 y §3.2.

### 3.8 `S3-convencion-prestamo.md` · fase S3
- **Q2-28** · ¿Cómo se calcula la tasa de cada periodo de un préstamo a plazo fijo y cómo entra el IVA? (a) tasa nominal: anual ÷ 12, ÷ 24 o ÷ 52; el IVA multiplica la tasa del periodo (r × 1.16) y el pago fijo se calcula con ella · (b) base diaria: días reales × anual ÷ 360, con el IVA aparte; los pagos dejan de ser idénticos · (c) capturas del contrato el pago fijo y el motor solo reparte capital e intereses. **Recomendada: (a)**, que el archivo describe como la convención usual de las tablas francesas en México (`docs/estado/wf1/contrato-S3.md` l. 51).
- PROVISIONAL: `Loan.Amortize` y la fixture `s3-prestamo-frances`, con datos de fixture y no de un préstamo real.
- Si no se contesta: la tabla del préstamo usa (a) como PROVISIONAL.

### 3.9 `S3-formula-tarjeta.md` · fase S3
K18 dice qué datos capturas de cada tarjeta, pero la fórmula que los usa no está escrita ni tiene fuente. **Las opciones (a), (b) y (c) son comunes a las preguntas 1 a 3:** (a) entregas un estado de cuenta de una tarjeta y uno de un crédito de tienda (PDF en `datos/semilla/estados-cuenta/`, fuera de git) y S3 escribe la fórmula de ese producto con un ejemplo comprobable; (b) se adopta una fórmula general PROVISIONAL: saldo promedio diario × tasa anual ÷ 360 × días en base diaria, o saldo al corte × tasa anual ÷ 12 en base mensual, con IVA si aplica y mínimo = el mayor entre % × saldo y el mínimo fijo; (c) el motor solo calcula lo ya decidido y para lo demás responde «falta capturar».
- **Q2-29** · (1) Intereses cuando no pagas el total: ¿sobre qué saldo y con qué días? (a) / (b) / (c). **Recomendada: (a).**
- **Q2-30** · (2) Pago mínimo: ¿el mayor entre % × saldo y el fijo, o % × saldo + intereses + IVA con piso en el fijo? (a) / (b) / (c). **Recomendada: (a).**
- **Q2-31** · (3) Crédito de tienda con intereses o comisión: ¿la misma fórmula que la tarjeta? (a) / (b) / (c). **Recomendada: (a).**
- Por qué (a): da una fórmula con fuente y cierra S3-08 y S3-10 (`contrato-S3.md` l. 51). Si eliges (a), tienes que entregar los PDF (sección 7).
- **Q2-32** · (4) Calendario de la tarjeta (letras de este documento): A) propuesta del agente: el corte es la primera fecha igual o posterior a la compra cuyo día sea el día de corte (recortado al último día del mes); el pago, la primera fecha posterior al corte cuyo día sea el día de pago; los meses sin intereses pagan en las fechas de pago de cortes sucesivos · B) la compra hecha el día de corte entra en el corte siguiente. **Recomendada: A**, la propuesta del agente en el propio archivo.
- PROVISIONAL: en uso está (b) en `Card.InterestForPeriod`, `Card.MinimumPayment` y `StoreCredit` con intereses, y en las fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos` (variante con interés). Son FIRMES: pago total a tiempo ($0.00), meses sin intereses y abonos sin interés.
- Si no se contesta: intereses y mínimos de tarjeta y tienda salen marcados PROVISIONAL, y S3-08 y S3-10 quedan parciales.

### 3.10 `S3-objetivos-medida.md` · fase S3
El archivo marca con * la propuesta del agente. Para las preguntas 2 a 4 no trae letras: A es la propuesta * y B la alternativa del archivo.
- **Q2-33** · (1) ¿La tienda es una dimensión del escenario del comparador y qué es «la canasta»? (a) la lista de compras de la semana con precios vigentes de cada tienda, y lo financiado es lo que cuesta en esa tienda · (b)* el monto de la petición es fijo; la tienda es dimensión solo si la petición no fija tienda; rendimiento = (monto × unidades ÷ costo de la canasta) ÷ (monto + intereses); `Scenario` necesita `store_id` (solicitud a C1) · (c) la tienda no es dimensión y solo se distinguen promociones. **Recomendada: (b).**
- **Q2-34** · (2) ¿Sobre qué serie se mide la desviación de pagos y qué peso tiene el pico contra el tope? A)* totales por quincena de todos los créditos del hogar en la ventana del escenario; el pico solo se marca y desempata · B) sumar a la medida el exceso sobre el tope. (La pregunta también menciona la serie de la compra sola; el archivo no la desarrolla como opción.) **Recomendada: A.**
- **Q2-35** · (3) Horizonte del colchón y del disponible: A)* la vida de la compra del escenario, de su fecha a su último pago · B) hasta el último pago de cualquier compra viva de la línea. **Recomendada: A.**
- **Q2-36** · (4) Créditos sin límite (préstamo, tienda sin límite): A)* solo entran créditos con límite; los demás quedan fuera con el motivo «sin límite» · B) entran y el colchón no aplica. **Recomendada: A.**
- Fuente de las cuatro recomendaciones: el propio archivo y `contrato-S3.md` l. 51.
- PROVISIONAL: los objetivos 2, 3 y 4, el colchón, `docs/estado/S3-escenarios.md` y los casos `s3-comparador-caso-1..3`.
- Si no se contesta: el comparador elige con las propuestas *. Si después cambian, cambian también los óptimos de los 3 casos.

### 3.11 `S3-origen-cargo-normal.md` · fase S3
- **Q2-37** · ¿Qué origen de intereses (`interest_source`) lleva un cargo normal de tarjeta? (a) `quoted`, leído como «lo fija el banco en su estado de cuenta»; sin cambio de esquema · (b) un valor nuevo (`none` o `statement`), con migración · (c) que admita vacío en el cargo normal, con migración. **Recomendada: (a)**, porque no toca el esquema (`correccion-contrato-S3-r2` y `-r3`; §7.2).
- PROVISIONAL: el origen del cargo normal y las fixtures `s3-unificado` y `s3-unificado-sin-corte`. Los intereses $0.00 son FIRMES.
- Si no se contesta: sale `quoted` con marca PROVISIONAL y ningún monto cambia.

### 3.12 `S3-serie-sin-intereses.md` · fase S3
Con intereses 0, la regla RN-07 deja intereses de −$0.01 por pago. Ejemplo: $3,094.50 a 4 quincenas da pagos de $773.62 × 3 + $773.64 y capital de $773.63 × 3 + $773.61. El esquema no admite intereses negativos.
- **Q2-38** · ¿Qué capital lleva cada pago cuando los intereses son 0? (a) capital = monto en los tres tipos (línea tipo Kueski, meses sin intereses y tienda); cabe en el esquema · (b) RN-07 literal, y el esquema pasa a admitir intereses negativos por pago (migración) · (c) RN-07 literal, ajustando en el primer pago; tampoco evita negativos y también exige migración. **Recomendada: (a)**: es la única serie con intereses de 0 o más en cada pago (§7.2; `correccion-contrato-S3-r2` y `-r3`).
- PROVISIONAL: la serie de capital de toda compra con intereses 0 y 2 pagos o más. Hoy la línea usa (b) y los meses sin intereses y la tienda usan (a).
- Si no se contesta: **WF-2 no puede guardar la serie de la promoción al 0 % de la línea con el esquema actual** (lo dice el archivo).

### 3.13 `S4-asignacion-mercado.md` · fase S4
RN-16 dice que un artículo de mercado va «a un mercado que surte su grupo», pero no dice a cuál ni qué pasa si no hay ninguno. Al arrancar no hay ninguno dado de alta.
- **Q2-39** · (1) Si varios mercados surten el grupo, ¿a cuál va? A) al de visita más próxima; en empate, al primero dado de alta · B) a un mercado «preferido» por grupo que marcas en Ajustes (cambia esquema y C1) · C) al más barato según S5; si ninguno tiene precio, A.
- **Q2-40** · (2) Si ningún mercado surte el grupo, ¿dónde aparece? A) en «mercado» sin tienda, con la leyenda «sin mercado dado de alta» y el total con el precio de S5 (exige cambiar C1) · B) pasa a la cadena que recomienda S5 y vuelve al mercado cuando exista uno.
- **Q2-41** · (3) Si una cadena es más barata para un producto de mercado: A) se queda en el mercado; la cadena se ve en la comparativa y puedes elegirla (exige cambiar C1) · B) pasa a la cadena si el ahorro supera un umbral (valor nuevo de K44) · C) va siempre a la tienda que recomienda S5.
- **Q2-42** · (4) Fecha de la próxima visita, orden de los grupos y total con artículos sin precio: A) la visita se cuenta desde la fecha más tardía entre hoy y el lunes de la semana de la lista; primero los grupos de mercado, luego los de cadena, y los que no tienen tienda al final; el total suma solo lo que tiene precio y se avisa cuáles faltan · B) la visita se cuenta desde hoy y el total queda vacío si falta un precio (cambia C1).
- **Recomendada: A** en las cuatro. `INTEGRACION.md` §7.2 dice «A, A y A», tres valores para cuatro preguntas; el reporte `correccion-contrato-S4S5-r2.md` l. 82 aplica A como PROVISIONAL en las 4. Nota del archivo: con 2B y 3C, C1 no cambia.
- PROVISIONAL: solo tienda, fecha de visita, orden y total de los grupos de la lista (P-13). Ningún precio cambia.
- Si no se contesta: la lista sigue con A, y la solicitud de cambio a C1 (`assigned_store_id`) queda pendiente.

### 3.14 `S4-confirmacion-empates.md` · fase S4
- **Q2-43** · ¿Quién confirma los empates producto ↔ dato nutricional (`nutricion/tabla.csv`) y las ligas de PROFECO (`profeco/vinculos.csv`), y cuándo? A) tú, en la puerta 2, sobre la lista de renglones «propuesto» · B) tú o Ximena en la app, después del despliegue (agrega una columna al esquema y operaciones a C1) · C) tú, antes de cerrar WF-1. HECHO: WF-1 ya cerró, así que C ya no es posible tal como está escrita. Un agente no puede confirmar: R9 y S4-11 piden una persona. **Recomendada: A.** No toca esquema ni C1, y nada entra sin confirmar (`correccion-contrato-S4S5-r1.md` l. 73). La lista concreta de nutrición es la Q2-90.
- PROVISIONAL: quién y cuándo. Ya está fijado que en modo producción no se carga ningún renglón «propuesto».
- Si no se contesta: la app arranca sin las kcal ni las ligas que no estén confirmadas.

### 3.15 `S4-reglas-lista.md` · fase S4
- **Q2-44** · (1) Sobrante previsto: A) se cocina lo planeado; el plan no genera sobrantes nuevos · B) se cocina la receta completa en múltiplos de sus porciones base y la diferencia es sobrante (la lista sube) · C) un campo nuevo «porciones a cocinar» por casilla (cambia esquema y C1).
- **Q2-45** · (2) Reposición por ritmo de consumo: A) si la existencia estimada al domingo es igual o menor que el nivel de reposición, se agrega el consumo de una semana; sin nivel no se repone · B) reponer hasta cubrir la semana más el nivel · C) agregar una presentación típica completa.
- **Q2-46** · (3) IMC: A) se clasifica con el valor que se muestra a un decimal · B) con el valor exacto (18.4975 se muestra como 18.5, pero sería «bajo peso»).
- **Recomendada: A** en las tres (`docs/estado/wf1/contrato-S4S5.md`, «Propuestas»; §7.2).
- PROVISIONAL: sobrantes, cantidad de los artículos por ritmo y clase de IMC en la franja de ±0.05 de cada corte.
- Si no se contesta: siguen las tres A como PROVISIONAL.

### 3.16 `S5-carrito-comparativa.md` · fase S5
El canon compara carrito contra carrito: Bodega Aurrera $3,094.50 contra Walmart $3,233.00, +$138.50. Solo la compra de Bodega se registró.
- **Q2-47** · (1) Al comparar una compra registrada, ¿contra qué se compara en las otras tiendas? A) cada tienda se cotiza con sus precios vigentes para los productos y cantidades de la compra base; no cambia C1; la API no reproduce +$138.50, pero el motor sí · B) C1 agrega un segundo parámetro (otra compra o un carrito guardado) y la API reproduce el canon; cambia C1 · C) las dos: A por omisión y B con el segundo parámetro. **Recomendada: A** (`correccion-contrato-S4S5-r1`; §7.2).
- **Q2-48** · (2) ¿Con qué fuente se siembran los precios de la compra de Bodega del 10-oct-2026? A) `manual`, captura sin ticket · B) `ticket_manual`, hay ticket impreso. El canon da lo mismo con las dos. **Sin recomendación en §7.2.** En uso está A, «porque no consta un ticket» (el archivo). Solo tú sabes si hubo ticket.
- **Q2-49** · (3) ¿Los precios del carrito de Walmart del 10-oct-2026 se siembran en la app? A) no, solo como entrada del motor y de la fixture; la app arranca sin precios de Walmart · B) sí, en `price` con fuente `manual`; se marcarían viejos el 9-nov-2026 (K22). **Recomendada: A** (§7.2).
- PROVISIONAL: la parte de la API, la fuente de la semilla de Bodega y si se siembra Walmart. El motor ya está fijado con el canon exacto.
- Si no se contesta: siguen las tres A como PROVISIONAL.

### 3.17 `S5-carritos-detalle.md` · fase S5
**Las opciones A, B y C del archivo contestan las 3 preguntas a la vez:** A = (1) lo que pagó por unidad, (2) faltantes de la unión, (3) por presentación agrupada; B = (1) mejor presentación, (2) siempre 0, (3) por línea; C = (1) de B y (2) y (3) de A.
- **Q2-50** · (1) Ganador de un sustituto cuando una tienda tiene varias presentaciones: ¿precio por unidad base de lo que esa tienda cobró (importe total ÷ contenido total) o el de su mejor presentación? A / B / C.
- **Q2-51** · (2) `missing_count` en carritos: ¿cuenta los productos de las otras canastas que faltan en esta tienda (A), o siempre 0 (B)? A / B / C.
- **Q2-52** · (3) `stale_count` en carritos: ¿cuenta cada presentación agrupada una vez, si tiene 30 días o más (A), o cada línea de `rows[]` (B)? A / B / C.
- **Sin recomendación explícita en §7.2.** En uso está A, que motor y oráculo ya implementan (`correccion-S5-r2.md`). El revisor sugiere sumar una pregunta 4 sobre la diferencia contra la más barata, que es la Q2-05.
- PROVISIONAL: ganador de sustitutos con varias presentaciones, `missing_count` y `stale_count`. No cambia ninguna cifra canon.
- Si no se contesta: sigue A, y el contrato S5 1.5 la escribe como PROVISIONAL (Q2-02).

### 3.18 `S5-profeco-almacen.md` · fase S5
- **Q2-53** · (1) ¿Dónde viven los registros crudos de PROFECO? A) en una base local aparte, fuera del esquema de la app e ignorada por git; solo lo filtrado y ligado pasa a `price` · B) en una tabla nueva de la app (migración 00004); crecen la base de producción y los respaldos · C) sin tabla, solo conteos en el reporte de cobertura. HECHO: hoy la base está en `datos/semilla/profeco/db/qqp.sqlite` y pesa 17 GB (`INTEGRACION.md` §8.6).
- **Q2-54** · (2) ¿Dónde se guarda la liga entre un producto PROFECO y una presentación del catálogo? A) en `datos/semilla/profeco/vinculos.csv`, confirmado a mano · B) en una columna nueva de `presentation`.
- **Q2-55** · (3) Con varios registros de la misma presentación, tienda y día, ¿cuál se guarda? A) el primero del archivo; los demás se cuentan como duplicados · B) el de menor precio · C) el promedio redondeado.
- **Recomendada: A** en las tres, porque no infla la base de producción (`contrato-S4S5.md`; §7.2). PROFECO nunca compite (K22).
- Si no se contesta: no se decide dónde vive la base de 17 GB en producción (pendiente 6 de `INTEGRACION.md` §8).

### 3.19 `S5-profeco-filtro.md` · fase S5
**Las opciones A, B y C del archivo contestan las 3 preguntas a la vez:** A = (1) solo municipios de sucursales, (2) comparación normalizada (minúsculas, sin acentos ni signos: «Wal-mart» = «WALMART»), (3) solo sucursales de la cadena que casó; B = (1) cualquier tienda con municipio, (2) igualdad exacta, (3) cualquier sucursal; C = (1) de A y (2) exacta, de B.
- **Q2-56** · (1) ¿El municipio de una tienda que no es sucursal (una cadena o un mercado) entra al filtro? A / B / C.
- **Q2-57** · (2) ¿Cómo se comparan los nombres y municipios de PROFECO con los del catálogo? A / B / C.
- **Q2-58** · (3) Para asignar la sucursal del precio, ¿solo se buscan las de la cadena que casó o cualquiera? A / B / C.
- **Sin recomendación explícita en §7.2.** En uso está A: el motor ya lo hacía y el oráculo se alineó (`correccion-S5-r2.md`). Con B, capturarías nombres y municipios con la ortografía exacta de PROFECO.
- PROVISIONAL: cargados, filtrados y cobertura de la importación (solo precios de referencia).
- Si no se contesta: sigue A como PROVISIONAL.

### 3.20 `S5-reglas-comparativa.md` · fase S5
- **Q2-59** · (1) ¿Cómo «pesa menos» un precio viejo al elegir la tienda más barata? A) se ordena por precio por unidad base ÷ peso, y totales y pantallas usan el precio real · B) el peso solo desempata · C) un precio viejo solo gana si no hay uno fresco en otra tienda. **Recomendada: A.**
- **Q2-60** · (2) Empates, y si puede ganar una tienda a la que le faltan precios: A) empate = precios ajustados iguales; el renglón queda sin ganador; la tienda recomendada es la que ya tiene más renglones asignados, luego la de menos precios viejos y luego la de menor id; una tienda con renglones sin precio no puede ser la más barata (precisión r3: en un renglón se desempata por el menor precio del paquete) · B) gana la tienda preferida que marques en Ajustes (campo nuevo) · C) gana la más cercana al día de despensa. **Recomendada: A.**
- **Q2-61** · (3) ¿Cuánto dura vigente un precio de promoción? A) solo el día en que se observó · B) columna nueva `price.valid_to` (cambia esquema y C1) · C) ligarlo a la tabla de promociones de crédito. **Recomendada: A.**
- **Q2-62** · (4a) ¿Con cuántas compras al año se anualiza el ahorro de la membresía de Sam's? A) una clave nueva que capturas en Tiendas y precios · B) las compras de despensa en cadenas de los últimos 365 días · C) 52 · D) 6. **Recomendada: B**: usa la frecuencia del historial real, sin inventar.
- **Q2-63** · (4b) Si costo ÷ ahorro no da entero, ¿se redondea hacia arriba? A) sí (8.3 compras → 9) · B) se muestra con un decimal (cambia el tipo entero de C1). **Recomendada: A.**
- Fuente: `contrato-S4S5.md`, «Propuestas»; §7.2: «A en 1 a 3 y 4b, B en 4a».
- PROVISIONAL: ganador con precios viejos, tienda recomendada en empates, promociones y los campos de la membresía cuando la división no es exacta.
- Si no se contesta: siguen esas opciones como PROVISIONAL. Ninguna cifra canon cambia.

### 3.21 `S5-terminos.md` · fase S5 · es también el hallazgo alto H-r3-00
- **Q2-64** · ¿Cómo se obtiene el veredicto de términos de walmart.com.mx, bodegaaurrera.com.mx y sams.com.mx? (A) autorizas por escrito que un agente lea **solo** las páginas legales de los 3 sitios, sin buscar ni extraer precios, y escriba `docs/estado/S5-terminos.md` con veredicto, URL y fecha; es una excepción explícita a «nada de sitios de tiendas» · (B) tú lees los términos y dictas el veredicto (permitido, no permitido o dudoso) con URL y fecha, y un agente lo transcribe · (C) queda pendiente con la lectura asistida apagada y S5-07 en «NO CUMPLE» · (D) descartas la lectura asistida en esta etapa; el integrador ajusta S5-07 y el plan §3.4. **Sin recomendación única:** §7.2 pide decidir entre A y B (`revisor-S5-r1`).
- PROVISIONAL: nada en cifras ni en esquema. La lectura asistida sigue apagada y no se puede activar, y `site_assisted` no compite (K37).
- Si no se contesta: S5-07 sigue sin cumplirse, y por este hallazgo alto el motor S5 sigue sin aprobar (`REVISION-S5-r3.md`). Nada más depende de esto.

### 3.22 `U3-t2-acumulados.md` · fase U3
La gráfica T-2 de Inicio («Devengado contra efectivo») mezclaba dos orígenes. Con una ventana que empieza en nov-2026, la cima marcaba $2,921.25 como «Devengado acumulado», y el devengado acumulado real es $3,895.00.
- **Q2-65** · A) C1 agrega el devengado acumulado y el efectivo acumulado con el mismo origen que la diferencia acumulada; son dos campos más en C1 y en S2 · B) sin cambiar C1, la gráfica se dibuja relativa a los 12 meses visibles; deja de ser la de DESIGN §8.3 y puede salir negativa · C) C1 agrega solo los saldos de apertura y la interfaz suma encima. **Recomendada: A.** La cima es el devengado acumulado sea cual sea el mes. El archivo la marca «Recomendada por U3» (`correccion-U3-r2`; §7.2).
- PROVISIONAL: las capas de T-2, sus etiquetas, su resumen y las cifras de capas de A-11. La tabla mensual de T-2 no depende de esto.
- Si no se contesta: T-2 sigue con dos variantes y WF-3 no puede dibujar la versión final.

### 3.23 `U4-dia-mes-reporte.md` · fase U4
- **Q2-66** · Si cambias el día del reporte mensual a 29, 30 o 31, ¿qué pasa en los meses que no tienen ese día? A) solo se aceptan los días 1 a 28 · B) se aceptan del 1 al 31 y en un mes más corto sale el último día · C) se aceptan del 1 al 31 y en un mes más corto no sale. **Recomendada: B**, para que nunca falte el reporte (`docs/estado/wf1/U4b-ajustes.md` l. 63).
- PROVISIONAL: qué pasa en los meses cortos. Con el valor inicial (1) no hay diferencia.
- Si no se contesta: la pantalla valida del 1 al 31 y S7 no tiene regla para los meses cortos.

### 3.24 `U4-iva-credito.md` · fase U4
- **Q2-67** · ¿Cómo se distingue «no capturado» de «la tasa ya incluye IVA»? 1) columna nueva `credit.vat_included`, que admite vacío como «falta capturar» (migración 00004 y C1 v1.3) · 2) usar `vat_on_interest_bp = 0` para «incluye IVA» y vacío para «falta capturar»; cambia el significado y la semilla Kueski pasaría a 0 · 3) dejarlo así: «incluye IVA» por omisión y sin aviso, con K18 parcial. **Recomendada: 1.** Cumple K18 sin cambiar el significado de lo que ya hay (`docs/estado/wf1/U4a-ajustes.md` l. 126).
- PROVISIONAL: el campo «¿la tasa incluye IVA?» de P-26.05. Ninguna cifra cambia: Kueski ya incluye IVA.
- Si no se contesta: el aviso «falta capturar» cubre 4 de los 5 datos de K18.

### 3.25 `U4-reparto-retroactivo.md` · fase U4
- **Q2-68** · Si cambias el porcentaje de una regla de reparto reutilizable (por ejemplo «Conjunto 50/50»), ¿vale para los gastos ya capturados o solo para los nuevos? A) para todo, como hoy; incluye meses anteriores y estados de cuenta ya entregados a Karen · B) solo para los nuevos, con una regla nueva o una vigencia «desde»; cambian S1 y C1 · C) se pregunta el alcance en cada cambio, como el presupuesto en K2; cambian C1 y el motor de reparto. **Recomendada: C**, para no cambiar sin querer estados de cuenta ya entregados a Karen (`U4b-ajustes.md` l. 64).
- PROVISIONAL: la vista previa y el texto de confirmación de P-26.12.
- Si no se contesta: queda A, el comportamiento actual de C1, y un cambio al 50/50 reescribiría lo ya capturado.

### 3.26 `U4-vista-previa-formulas.md` · fase U4
La pregunta 2 se retiró y se contesta en `S2-metas` (Q2-16 a Q2-18).
- **Q2-69** · (pregunta 1) ¿Las comisiones de un crédito (anual, por pago tardío y por apertura) entran al costo, a la deuda o al disponible? (a) no; solo se guardan como dato, como hoy · (b) la de apertura se suma al costo de la primera compra y la anual se reparte por mes (cambian S3 y la proyección) · (c) se registran como cargos manuales cuando ocurren. **Recomendada: (a)** (`correccion-U4-r1`; §7.2).
- **Q2-70** · (pregunta 3) Días restantes y consumo por semana de un producto por ritmo: (a) días = piso(existencia × días del ritmo ÷ cantidad del ritmo) y consumo semanal = cantidad × 7 ÷ días, redondeado · (b) los días también redondeados en lugar de piso · (c) sin días restantes, solo el estado (K7). **Recomendada: (a)** (`correccion-U4-r1`; §7.2).
- PROVISIONAL: el «Antes → Después» de esas cifras en P-26.05 (vista previa 8) y P-26.14.
- Si no se contesta: esas vistas previas solo dicen que la cifra se recalcula al guardar.

### 3.27 `WF1-catalogo-huevo-por-kilo.md` · fase WF1 (catálogo)
- **Q2-71** · ¿Cómo se registra un precio por kilo del huevo, cuya unidad base es la pieza? 1) el huevo se deja en pieza y capturas por pieza o por cartón (si compras por kilo, divides) · 2) S5 convierte kilo a pieza con 50 g por pieza; requiere un campo nuevo y cambia el esquema · 3) pasar el huevo a gramos; las recetas van en gramos · 4) dos productos: huevo por pieza y huevo a granel. **Recomendada: 1**: cero cambios (`docs/estado/wf1/correccion-catalogo-r1.md` l. 47).
- PROVISIONAL: nada. El catálogo sigue válido.
- Si no se contesta: el precio por kilo del mercado no se captura tal cual. Nada se bloquea.

## 4. PROPUESTA · valores operativos K31 y K44

Fuente: `INTEGRACION.md` §7.1, recogida en INDICE. Son PROVISIONALES hasta que los apruebes. Opciones en cada uno: **A** aprobar · **B** cambiar a ___. Recomendada: A, que es la propuesta del agente.

| Q2 | Valor | Propuesta | Por qué | Agente |
|---|---|---|---|---|
| **Q2-72** | Hora de los reportes (K31 punto 1) | diario 21:30; semanal domingo 19:00; mensual día 1 a las 08:00 | de noche ya se capturó el día; el domingo se planea la semana; el día 1 el mes ya cerró. Dónde se guarda depende de Q2-11 | `U4b-ajustes` |
| **Q2-73** | Avisos (K31 punto 2) | anticipación del aviso de precios viejos = 1 día; hora de avisos 09:00 | un día antes de la despensa da tiempo de actualizar precios | `U4b-ajustes` |
| **Q2-74** | Tiempo máximo del comparador (K31 punto 7) | mantener 5 s (K43); límites de seguridad de 1 s a 30 s | el peor caso medido es de 111 a 123 ms (3,870 candidatos), el 2.5 % del tope; los casos reales tardan 1 ms o menos | `motor-S3-comparador-go`, `oraculo-S3` |
| **Q2-75** | `optimizer.date_window_days` | 7 días | cubre una semana y mueve el primer pago a la quincena vecina; el peor caso es de unos 7,600 escenarios, muy por debajo de 5 s | `contrato-S3` |
| **Q2-76** | `ocr.timeout_seconds` | 30 s | 2 imágenes sintéticas tardan unos 2 s; una foto de 12 MP en un equipo lento, de 5 a 10 s; deja un margen de 3× | `ocr`, `contrato-S4S5` |
| **Q2-77** | Rango de frescura (K44) | `price.stale_after_days` de 1 o más y `price.zero_weight_after_days` igual o mayor que el anterior; fórmula del peso (Z − días) ÷ (Z − S + 1), que con 30 y 90 da K40 | con otros valores la recta no baja de 1 a 0; el motor los rechaza con `invalid_argument` | `motor-S5-go`, `U4a-ajustes` |

## 5. PROPUESTA · otras decisiones a ratificar

Fuente: `INTEGRACION.md` §7.3; el detalle está en cada reporte citado. Si no se dice otra cosa, las opciones son **A** ratificar · **B** cambiar (dices cómo).

**Diseño** (`docs/estado/wf1/U2b-1-design.md` y `U2b-2-design.md`, «Propuestas para Axel»):
- **Q2-78** · Negativos con el signo menos tipográfico (U+2212) antes de `$`: `−$1,258.68` (DESIGN.md §3.3). Recomendada: A.
- **Q2-79** · Un restante rebasado se muestra en color de advertencia, con icono y la palabra «rebasado» (§3.3). Recomendada: A.
- **Q2-80** · La interfaz te habla de tú (§1.3). Recomendada: A.
- **Q2-81** · El color del ingreso queda en 135°, que ya está en los tokens desde U2a. Decisiones §I dice 150° «o U2 propone 135°». A) 135° · B) 150°. Recomendada: A.
- **Q2-82** · Interruptor «Atajos de una tecla», encendido por omisión y guardado en el dispositivo, para cumplir WCAG 2.1.4. Recomendada: A.
- **Q2-83** · Indicador «Guardado» de 2 s junto al control en Ajustes; el aviso breve sigue en 5 s (U1 G-07). Recomendada: A.
- **Q2-84** · Las constantes de interfaz (5 s, 2 s, 400 ms, 10 s) son de diseño y no parámetros de K44. A) ratificar · B) que se puedan configurar, en Apariencia. Recomendada: A.

**U3 Inicio** (`docs/estado/wf1/U3-dashboard.md` l. 51, `correccion-U3-r1`):
- **Q2-85** · Supuesto 5 (cómo se dibuja T-2). HECHO: ya se **retiró** en U3 1.1 (`docs/design/U3-dashboard.md` l. 459), y cómo se dibuja T-2 lo decide ahora la Q2-65. A) darlo por cerrado · B) reabrirlo. **Sin recomendación** de los agentes: el reporte que lo pedía es anterior al retiro.
- **Q2-86** · Nombre de la segunda tienda en la tarjeta 5 («$138.50 menos que…»): hoy C1 no lo entrega y se ve en el detalle D-5. A) aprobar la solicitud a C1 para que lo entregue · B) dejarlo como hoy. Recomendada: A (U3 pide confirmar la solicitud).

**Recetas** (`correccion-S4-r1`, `revisor-S4-r2`, `nutricion`):
- **Q2-87** · ¿Las recetas llevan etiquetas («económica», «rápida», «sobrantes») para filtrar en el plan? (a) no; el recetario no cambia · (b) sí, como campo opcional: contrato S4 1.5, esquema, cargador y una columna o tabla nueva. Los borradores ya las traen para las 60 recetas. **Sin recomendación.**
- **Q2-88** · Kcal de recetas con cortes con hueso (B7): (a) aceptarlas como están, sobreestimadas · (b) escribir las cantidades en peso comestible · (c) agregar un factor de parte comestible por producto (cambia contrato y esquema). **Sin recomendación.**
- **Q2-89** · Recetas parecidas (B8): huevos con nopales / huevos a la mexicana, y chilaquiles verdes con huevo / con pollo. Letras de este documento: A) no cambiar ninguna · B) cambiar alguna (dices cuál). **Sin recomendación.**
- **Q2-90** · Confirmar o corregir los 40 empates de `datos/semilla/nutricion/revisar.tsv`. Lo que más importa: el % de grasa de la carne molida que compras, si la milanesa es empanizada, la jamaica como flor seca, el café molido (OFF repite 396.6 kcal en tres marcas) y el chipotle en adobo (49 o 73 kcal). Letras de este documento: A) confirmar tal cual · B) corregir (dices cuáles). **Sin recomendación**: necesita datos que solo tú tienes. El «cuándo» es la Q2-43.

**PROFECO y lista** (`profeco`, `contrato-S4S5`, `correccion-contrato-S4S5-r4`):
- **Q2-91** · Cadenas del catálogo: `Wal-mart` como nombre PROFECO de Walmart, `Bodega Aurrera`, y Sam's Club sin nombre PROFECO (no hay datos). Las dos variantes `Express` van como cadenas aparte si quieres sus precios. A) aprobar con Express aparte · B) aprobar sin Express · C) cambiar. **Sin recomendación** (A o B): el agente deja a tu elección si quieres precios de Express (`docs/estado/wf1/profeco.md`).
- **Q2-92** · Municipios por defecto en Ajustes: el de cada sucursal registrada (Gustavo A. Madero para tus tres sucursales, si se confirman). Recomendada: A.
- **Q2-93** · Los 9 grupos iniciales del catálogo (`productos-formato.md` §4), editables en Ajustes. Recomendada: A.
- **Q2-94** · Subtotal opcional «falta por comprar» en la lista, que no cuenta lo ya marcado como comprado (campo calculado nuevo en C1, sin efecto en cifras canon). A) agregarlo · B) no por ahora. **Sin recomendación**: el agente lo marca como opcional.

## 6. PROPUESTAS de proceso

- **Q2-95** · (a) **WF-3 no espera a la Pista 3.** Ya estaba anotada como propuesta para esta puerta: K35 dice «WF-3 no espera a la Pista 3 (plan R4: referencia visual, no dependencia)» (`docs/plan/decisiones.md` l. 105) y `k45-borrador.md` dice que «sigue en pie para la puerta 2». WF-3 arranca sobre DESIGN.md y los tokens C2, y lo que devuelva la Pista 3 en `docs/design/salida/` se incorpora cuando llegue. A) aprobar · B) WF-3 espera a la referencia visual final. Recomendada: A.
- **Q2-96** · (b) **La Pista 3 por Artifacts arranca en cuanto apruebes DESIGN.md y C2 como base del sistema de diseño** (`docs/estado/puerta2/pista3-plan-borrador.md` §1: el Design System se arma desde `C2-tokens.json` y DESIGN.md). HECHO: DESIGN.md dice en su encabezado «Versión 1.0 en construcción»; `C2-tokens.json` es la versión 1.0; los dos quedaron aprobados en la revisión U2 de la ronda 4, con 11 bajos. A) aprobar DESIGN.md 1.0 y C2 1.0 tal cual · B) aprobarlos con cambios (dices cuáles) · C) todavía no; la Pista 3 espera. Recomendada: A, porque la revisión U2 los aprobó en la ronda 4 (`INTEGRACION.md` §1, `REVISION-U2-r4.md`).
- **Q2-97** · (c) **Ronda de aplicación P2c** después de tus respuestas. El orquestador anota «Respuesta de Axel» en cada bloqueo y la consolida en `decisiones.md` §K, como en las puertas 1 y 1b (INDICE). Luego aplica las correcciones autorizadas en la sección 2, los valores de la 4 y lo ratificado en la 5; incluye los cambios de K45 que ya están en el árbol de trabajo (`k45-borrador.md`, «Archivos del repo a actualizar»); corre `make verify-all` con sus revisores, y cierra con un commit «[Puerta 2]». A) aprobar · B) cambiar (dices cómo). Recomendada: A (propuesta del orquestador).
- **Q2-98** · (d) **Nombres en las maquetas de la Pista 3.** Es un supuesto declarado del orquestador en K45.1: aparecen Axel, Ximena y Karen, porque son los usuarios y la paleta los nombra; las cifras sí son sintéticas. A) nombres reales · B) seudónimos (solo cambia el archivo de datos sintéticos). Recomendada: A (supuesto del orquestador, `k45-borrador.md`).

## 7. Acciones pendientes de Axel

No llevan id porque no son decisiones: son cosas que solo tú puedes hacer.
1. **Fotografiar 5 tickets reales** y guardarlos en `datos/semilla/tickets/` (K24). Con eso se mide el OCR (S5-06) con `go run ./internal/precios/ocr/cmd/ocrmedir` (`INTEGRACION.md` §8.2). Hoy la carpeta solo tiene `.gitkeep` (`REVISION-S5-r3.md`).
2. **Entregar `datos/semilla/gastos_fijos_ingresos.csv`**, con la plantilla de `plan.md` §8.2 (`docs/plan/decisiones.md` l. 58). Si no existe al llegar WF-4, la carga se hace sin esa parte y el importador queda listo (`plan.md` l. 186).
3. **Confirmar con Kueski** cómo repone el capital. Mientras tanto, el interruptor A/B de Ajustes está en B por omisión y se ve el aviso «pendiente confirmar» (`docs/plan/decisiones.md` l. 22; riesgo R6 de `plan.md`, que también incluye la tasa personalizada).
4. **Instalar la CA de mkcert en los dos Android** (WF-4). Sin HTTPS no se puede instalar la PWA (riesgo R3 de `plan.md`; los pasos van en `ops/mkcert.md`).
5. **Ximena escribe al bot** una vez para registrar su chat, y tú lo confirmas en Ajustes (K4; riesgo R15 de `plan.md`).
- Según lo que contestes: con (a) en Q2-29 a Q2-31, entregar el PDF de un estado de cuenta de tarjeta y uno de crédito de tienda; con (B) en Q2-64, leer los términos de los 3 sitios y dictar el veredicto; con A en Q2-43, confirmar las listas de empates y ligas (Q2-90); en Q2-48, decir si hubo ticket impreso de Bodega; y al final de la Pista 3, abrir la URL del lienzo y aprobar (puerta 3 nueva, K45).

## 8. Hoja de respuestas

| Q2 | Tema | Recomendada | Respuesta de Axel |
|---|---|---|---|
| Q2-01 | Motor S2 · M1 devengo con ventanas (hueco y ventanas que no duran un mes) | A (el caso 2 va sin recomendación) | |
| Q2-02 | Motor S5 · H-r3-01 contrato S5 1.5 | A | |
| Q2-03 | Motor S5 · H-r3-02 compras a granel | A | |
| Q2-04 | Motor S5 · H-r3-03 lectura de PROFECO | A | |
| Q2-05 | Motor S5 · H-r3-04 diferencia contra la más barata | A | |
| Q2-06 | Contrato S4/S5 · M-r5-1 cotización por elección | A | |
| Q2-07 | Contrato S4/S5 · M-r5-2 presentación sin precio | A (variante sin recomendación) | |
| Q2-08 | Contrato S4/S5 · M-r5-3 orden de las listas | A | |
| Q2-09 | U4 · R5-1 vista previa en listas, arrastre y Ctrl+/ | A | |
| Q2-10 | Paquete · H-1 paso 2 del LEEME | B (superado por K45) | |
| Q2-11 | S1-parametros-avisos · dónde viven los 4 parámetros | sin recomendación (B o A) | |
| Q2-12 | S2-flujo-quincenal (1) gasto variable estimado | C o A | |
| Q2-13 | S2-flujo-quincenal (2) ahorro entre quincenas | B | |
| Q2-14 | S2-flujo-quincenal (3) partida cumplida | A | |
| Q2-15 | S2-metas-resumen · activas o todas | A | |
| Q2-16 | S2-metas (1) meses que faltan | A | |
| Q2-17 | S2-metas (2) formato del avance | C | |
| Q2-18 | S2-metas (3) estados de la meta | A | |
| Q2-19 | S2-presupuesto-total (1) renglón total o suma | B | |
| Q2-20 | S2-presupuesto-total (2) plan $0.00 | A | |
| Q2-21 | S2-proyeccion-interruptor (1) monto | C | |
| Q2-22 | S2-proyeccion-interruptor (2) fecha | C | |
| Q2-23 | S2-proyeccion-interruptor (3) reparto | C | |
| Q2-24 | S2-proyeccion-interruptor (4) papel alterno | C | |
| Q2-25 | S2-ventanas (1) antes de la primera ventana | A | |
| Q2-26 | S2-ventanas (2) movimiento en un hueco | H1 | |
| Q2-27 | S2-ventanas (3) orden y traslape | O1 | |
| Q2-28 | S3-convencion-prestamo · tasa e IVA | (a) | |
| Q2-29 | S3-formula-tarjeta (1) intereses | (a) | |
| Q2-30 | S3-formula-tarjeta (2) pago mínimo | (a) | |
| Q2-31 | S3-formula-tarjeta (3) crédito de tienda | (a) | |
| Q2-32 | S3-formula-tarjeta (4) calendario de corte y pago | A | |
| Q2-33 | S3-objetivos-medida (1) canasta y tienda | (b) | |
| Q2-34 | S3-objetivos-medida (2) serie y pico | A | |
| Q2-35 | S3-objetivos-medida (3) horizonte del colchón | A | |
| Q2-36 | S3-objetivos-medida (4) créditos sin límite | A | |
| Q2-37 | S3-origen-cargo-normal | (a) | |
| Q2-38 | S3-serie-sin-intereses | (a) | |
| Q2-39 | S4-asignacion-mercado (1) varios mercados | A | |
| Q2-40 | S4-asignacion-mercado (2) ningún mercado | A | |
| Q2-41 | S4-asignacion-mercado (3) cadena más barata | A | |
| Q2-42 | S4-asignacion-mercado (4) visita, orden y total | A | |
| Q2-43 | S4-confirmacion-empates · quién y cuándo | A | |
| Q2-44 | S4-reglas-lista (1) sobrante | A | |
| Q2-45 | S4-reglas-lista (2) reposición por ritmo | A | |
| Q2-46 | S4-reglas-lista (3) IMC | A | |
| Q2-47 | S5-carrito-comparativa (1) contra qué | A | |
| Q2-48 | S5-carrito-comparativa (2) fuente de Bodega | sin recomendación (en uso A) | |
| Q2-49 | S5-carrito-comparativa (3) sembrar Walmart | A | |
| Q2-50 | S5-carritos-detalle (1) ganador del sustituto | sin recomendación (en uso A) | |
| Q2-51 | S5-carritos-detalle (2) missing_count | sin recomendación (en uso A) | |
| Q2-52 | S5-carritos-detalle (3) stale_count | sin recomendación (en uso A) | |
| Q2-53 | S5-profeco-almacen (1) registros crudos | A | |
| Q2-54 | S5-profeco-almacen (2) liga | A | |
| Q2-55 | S5-profeco-almacen (3) duplicados | A | |
| Q2-56 | S5-profeco-filtro (1) municipios | sin recomendación (en uso A) | |
| Q2-57 | S5-profeco-filtro (2) comparación de nombres | sin recomendación (en uso A) | |
| Q2-58 | S5-profeco-filtro (3) sucursal del precio | sin recomendación (en uso A) | |
| Q2-59 | S5-reglas-comparativa (1) peso del precio viejo | A | |
| Q2-60 | S5-reglas-comparativa (2) empates | A | |
| Q2-61 | S5-reglas-comparativa (3) vigencia de promoción | A | |
| Q2-62 | S5-reglas-comparativa (4a) compras al año (Sam's) | B | |
| Q2-63 | S5-reglas-comparativa (4b) redondeo hacia arriba | A | |
| Q2-64 | S5-terminos · veredicto de los 3 sitios | sin recomendación (A o B) | |
| Q2-65 | U3-t2-acumulados · gráfica T-2 | A | |
| Q2-66 | U4-dia-mes-reporte · días 29 a 31 | B | |
| Q2-67 | U4-iva-credito · «falta capturar» del IVA | 1 | |
| Q2-68 | U4-reparto-retroactivo · alcance del cambio | C | |
| Q2-69 | U4-vista-previa-formulas (1) comisiones | (a) | |
| Q2-70 | U4-vista-previa-formulas (3) días restantes | (a) | |
| Q2-71 | WF1-catalogo-huevo-por-kilo | 1 | |
| Q2-72 | K31 · hora de los reportes | A | |
| Q2-73 | K31 · avisos de precios viejos | A | |
| Q2-74 | K31 · tiempo máximo del comparador | A | |
| Q2-75 | K44 · ventana de fechas del optimizador (7 días) | A | |
| Q2-76 | K44 · tiempo máximo del OCR (30 s) | A | |
| Q2-77 | K44 · rango de frescura | A | |
| Q2-78 | Diseño · signo menos U+2212 | A | |
| Q2-79 | Diseño · «rebasado» en advertencia | A | |
| Q2-80 | Diseño · trato de tú | A | |
| Q2-81 | Diseño · ingreso en 135° | A | |
| Q2-82 | Diseño · interruptor «Atajos de una tecla» | A | |
| Q2-83 | Diseño · «Guardado» de 2 s | A | |
| Q2-84 | Diseño · constantes de interfaz fuera de K44 | A | |
| Q2-85 | U3 · supuesto 5 (ya retirado) | sin recomendación | |
| Q2-86 | U3 · segunda tienda en la tarjeta 5 | A | |
| Q2-87 | Recetas · etiquetas | sin recomendación | |
| Q2-88 | Recetas · kcal de cortes con hueso | sin recomendación | |
| Q2-89 | Recetas · recetas parecidas | sin recomendación | |
| Q2-90 | Recetas · 40 empates de revisar.tsv | sin recomendación | |
| Q2-91 | PROFECO · cadenas del catálogo | sin recomendación (A o B) | |
| Q2-92 | PROFECO · municipios por defecto | A | |
| Q2-93 | Catálogo · 9 grupos iniciales | A | |
| Q2-94 | Lista · subtotal «falta por comprar» | sin recomendación | |
| Q2-95 | Proceso · WF-3 no espera a la Pista 3 | A | |
| Q2-96 | Proceso · DESIGN.md 1.0 y C2 1.0 como base de la Pista 3 | A | |
| Q2-97 | Proceso · ronda P2c y commit «[Puerta 2]» | A | |
| Q2-98 | Pista 3 · nombres reales en las maquetas | A | |
