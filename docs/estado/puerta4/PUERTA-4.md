# Puerta 4 · Finanzas Hogar · lo que Axel decide para pasar a operación

Redactor: agente de la puerta 4 (Opus) · 12-oct-2026, 21:20 · commit base **8ace01e** («[WF-3c ter] Comparador solo con líneas activas (A-1) y D-3 de Inicio Detalle espera a listPersons (M-1)…») · repositorio `/home/deepsec/finanzas-hogar`. Lo que pide la puerta (`docs/plan/plan.md` §3.8): «C3 firmado, capturas U5, informe de pruebas → aprobar paso a operación».

## 0. Cómo responder

**Contesta por id con la letra de la opción (p. ej. «Q4-31: B»); en las preguntas de varios puntos, punto y letra («Q4-42: 2b, 3a»). Puedes contestar por lotes: «todo lo recomendado salvo Q4-xx: B».**

- Hay **110 preguntas**, de Q4-01 a Q4-110, sin huecos. La §2 trae las 3 de aprobación y proceso; la §3, las 107 de bloqueos abiertos, en el orden de `docs/estado/bloqueos/INDICE.md` (WF-1, P2c, P2d, Pista 3, semilla, WF-3a, WF-3a bis, WF-3b, WF-3c y WF-3c bis) y, dentro de WF-3a, por tema. La hoja de la §5 las junta todas.
- Las preguntas de la puerta 2 que siguen abiertas van como lista corta con su id original (§3.1): se contestan con ese id y la letra de `docs/estado/puerta2/PUERTA-2.md` («Q2-50: B»).
- Las letras son las del archivo de origen. Cuando el archivo no pone letras, lo digo y las letras son de este documento. No agregué opciones: donde el archivo no las trae, dice «sin opciones escritas», y ahí contestas con tu valor o con «PROVISIONAL» para dejar el que corre.
- Etiquetas: **HECHO** (ya pasó o está escrito), **PREGUNTA** (decides tú), **PROPUESTA** (sugerencia de un agente, de un integrador o del orquestador; no vale hasta que la apruebes), **PROVISIONAL** (lo que corre hoy en el código mientras no contestas).
- «Recomendada» es la opción que propone un reporte, con su cita. «Sin recomendación» quiere decir que ningún reporte propuso una o que los agentes no coinciden.
- **PROPUESTA del redactor** (la misma regla de la puerta 2): si contestas «todo lo recomendado» y no tocas una pregunta sin recomendación, se queda abierta y sigue PROVISIONAL con la opción en uso; WF-4 trabaja con esa opción. No hace falta contestar todo hoy: para pasar a operación pesan sobre todo Q4-01 a Q4-03, los valores K31 (§3.6.1) y el tablero (§3.9 y §3.10).
- Las cifras salen de los archivos citados; este documento no agrega valores nuevos.

## 1. HECHO · lo que recibes en la puerta 4

### 1.1 C3 firmada (`docs/contratos/C3-integracion.md` §9)
- **C3 v1.4 FIRMADA** por `revisor-c3-cierre-r1` (independiente), ronda 1, 12-oct-2026. Veredicto: **aprobado**; se cumplen C3-01 a C3-11 con ejecuciones propias del revisor y quedan 3 observaciones bajas que no abren ningún criterio.
- **SHA-256 del contenido firmado:** `e341d182e08567d797ae503c3d3970b8d3542d650a632806559fc2df5507d163`. Lo recalculé hoy a las 21:20 con el comando del bloque de firma y da el mismo valor.
- Matriz: implementado 42 de 42, probado 42 de 42 y capturado 41 de 42. La fila Global sigue en «no» por `C3-capturas-global` (Q4-87). Hay 221 operaciones; las 5 sin fila tienen su excepción: las 4 del agente externo (K53) y `getOpenApiDocument` (`INTEGRACION-wf3c.md` §15.4).
- Evidencia propia del revisor en una copia limpia: Go 1151 PASS y 0 FAIL; Playwright 271 pasan y 4 omitidas; **8 de 8 flujos** contra el servidor real; C3-09: 120 cifras leídas en pantalla y 0 diferencias.
- Observaciones bajas de la firma: B-1, la prueba F-3A0 (`flujos.spec.ts:71`) sale en verde aunque su paso diga FALLA, porque le falta `expect(reg.ok)`; B-2, U1 y U3 siguen diciendo +$138.50 (Q4-81); B-3, la línea 3 de C3 todavía dice «sin firma», y corregirla pide una firma nueva porque es parte del contenido firmado.

### 1.2 Capturas y criterios U5 (`INTEGRACION-wf3c.md` §15.2, `docs/estado/wf3/u5-cierre.md`, `INTEGRACION-wf3b.md` §3)
- `u5-cierre` regeneró las **164 capturas** (41 pantallas × Android 412×915 y escritorio 1440×900 × claro y oscuro) con el `dist` de WF-3c bis, en `docs/design/capturas-u5/`.

| Criterio | Estado | Evidencia |
|---|---|---|
| U5-01 | Cumple | `verificar_capturas.py`: 164 de 164, sin faltantes ni medidas mal. El `ls \| wc -l` literal da 168 (Q4-88) |
| U5-02 | Cumple | Se ven las cifras de la semilla: $1,191.50, $10,000.00, $3,094.50, $6,293.44, $6,905.50 y $973.75; ninguna captura quedó cargando ni en error |
| U5-03 | Cumple | axe: 0 violaciones serias o críticas en las 4 combinaciones |
| U5-04 | Cumple | 25,818 pares evaluados y **0 bajo umbral**; la regla de «componente de UI» sigue PROVISIONAL (Q4-89) |
| U5-05 | Escalado | La revisión visual quedó «devuelta» en su ronda 3 (0 altos y 2 medios, M-7 y M-8) y no hay ronda 4: lo decides en Q4-82. M-7 (iconos de la barra lateral) se atendió en WF-3c con iconos PROVISIONALES (Q4-104); M-8 se atendió desde C3 v1.0 |
| U5-06 | Atendido, sin revisión visual nueva | El #10 (iconos) se atendió en WF-3c; el desborde de P-01 Detalle en Android se cerró en WF-3c bis: 412 de 412 px (`INTEGRACION-wf3c.md` §15.3) |
| U5-07 | Cumple | 0 fechas ISO, 0 palabras en inglés y 0 claves internas; la lista de nombres propios sigue PROVISIONAL (Q4-90) |
| U5-08 | Cumple | Un solo host, `finanzas.local:8443`, en las 164 capturas; 0 hosts externos |
| U5-09 | Cumple | 3 toques para registrar un gasto sin conexión (`cola.spec.ts`, WF-3b) |
| U5-10 | Cumple | Fondo del body igual al token de C2 en 164 de 164; `tema.spec.ts` pasa |
| U5-11 | Cumple | 41 filas P- de C3 con «capturado» y sus 4 rutas en disco |
| U5-12 | Cumple | Bloqueos `U5-*` escritos; nada se marcó aprobado con violaciones abiertas |
| U5-13 | Cumple | Cola sin conexión, reenvío idempotente y error 500 probados en un servidor aislado (WF-3b) |

### 1.3 Informe de pruebas
- **`make verify-wf3`: EXIT 0** (integrador de WF-3c bis, 12-oct-2026, de 20:05 a 20:25): Go 1.27.2 con `build`, `vet`, `gofmt` y `go mod tidy` limpios y 33 paquetes OK; cobertura K25 OK (dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8); contrato C1 1.5.1 con 221 operaciones en 176 rutas y 0 diferencias; interfaz con la API simulada: 271 pasan, 4 omitidas y 0 fallas; `verify-u5` como en la §1.2.
- **`make verify-all`: EXIT 0** («Verificación ampliada OK»); goose sube, baja a 0 y vuelve a subir las migraciones.
- `make verify-wf3-secretos`: 0 coincidencias. El registro del servidor no trae contraseñas ni tokens.
- **WF-3c ter (commit 8ace01e):** se corrigieron los dos hallazgos que WF-3c bis dejó sin aprobar. A-1 (alto): el comparador de escenarios ya solo ofrece líneas activas (`activeStates()`, prueba `TestOp_compareScenarios_LineaArchivada`). M-1 (medio): D-3 de P-01 Detalle espera a `listPersons` y ya no muestra «Hogar $973.75» sin la parte de Karen (3 pruebas nuevas). El commit registra `gofmt`, `vet`, `verify-wf3-go`, `verify-wf3-contrato` (221 de 221), `tsc`, `build` y `make verify-all` OK. Reportes: `docs/estado/wf3/correccion-a1-comparador.md` y `correccion-m1-detalle.md`.
- El sitio responde por **HTTP/2** en `https://finanzas.local:8443` y en `https://192.168.50.56:8443`. Escucha solo en `127.0.0.1:8443` y `192.168.50.56:8443`; por Tailscale no responde. Hoy corre como **servidor de desarrollo** (`ops/dev-run.sh`, con la semilla de desarrollo), no como servicio de systemd: eso es WF-4.

### 1.4 Servidor (S6) y avisos (S7) (`INTEGRACION-wf3a.md` §4, `INTEGRACION-wf3b.md` §8)
- **S6, 16 criterios:** 15 cumplen. **S6-11** cumple en la máquina (mkcert 1.4.4; el certificado lleva `finanzas.local`, `localhost` y las 2 IP; curl con la CA da 200), pero **Chrome en Android está NO VERIFICADO**: lo pruebas tú (§4). S6-09 cumple con hora y retención PROVISIONALES (Q4-29) y S6-10, con tipos PROVISIONALES (Q4-31).
- **S7, 14 criterios:** 13 cumplían desde WF-3a. **S7-07, el envío real por Telegram, lo hizo el orquestador hoy a las 21:15**: el servicio mandó el resumen diario y la alerta de pago al chat de Axel, con estado `sent` y HTTP 200. Falta que confirmes que los recibiste (§4). Ximena sigue «pendiente de acción»: su chat depende de Q4-43.

### 1.5 Lo que no está verificado o sigue abierto
- **S6-11 en tus dos Android** (manual; `ops/mkcert.md` §4).
- **Revisión independiente de WF-3c ter.** No hay `REVISION-*-r3` de esas dos correcciones; solo los reportes de autor y las verificaciones del commit.
- **OCR sin medir:** `datos/semilla/tickets/` está vacía; faltan tus 5 fotos (K24).
- **Hallazgos bajos y solicitudes sin cerrar** (`INTEGRACION-wf3c.md` §15.3 y §15.5): el `Makefile` dice «esquema 1.4» en 3 líneas (B-N4); `Estados.tsx` con `role="status"` (N-3); B-1 de la firma de C3; las 3 líneas de U1 y U3 con +$138.50; `deriveCredits` (`backend/internal/api/gastos-inicio/vistas.go:235`) sigue usando `mc.LiveDebt`, así que la vista de base de datos de créditos todavía cuenta las compras futuras; `GET /api/v1/debt/live` sin fecha respondió 2027-09-30, que coincide con la fecha por omisión PROVISIONAL de `getLiveDebt` (Q4-56, punto 3); `docs/criterios/C1.md` sigue en 1.5.0; `humo-inicio.spec.ts` espera las 6 tarjetas de U3 1.x y no entra en ningún `verify-*`; la insignia de `Calendario.tsx` (`U3-tablero-k4-personas`).
- **Avisos que hoy no salen:** en `backend/internal/reportes/avisos/fuente.go`, las fuentes de presupuesto (alertas de 80 % y 100 %), de precios viejos y del día de despensa responden «sin fuente», así que esas alertas se omiten y quedan en el registro; nunca inventan cifras (Q4-36). El resumen diario, el semanal y la alerta de pago sí tienen fuente.
- **Bloqueos en disco:** 154 archivos de bloqueo más `INDICE.md`. El comando del encargo (`grep -L "Respuesta de Axel"` más los que dicen «Estado: ABIERTO», sin `INDICE.md`) da **90**, la misma cifra que contó el integrador de WF-3c bis (`INTEGRACION-wf3c.md` §15.6.3). De esos 90, 2 dicen RESUELTO en su propio texto (`criterios-nxx.md` y `WF3-ajustes-registro.md`, §3.11 y Q4-64). Además, `C3-flujos.md` sigue abierto según el índice, aunque el comando no lo cuenta: trae la frase «Respuesta de Axel» porque cita K46. Los 91 están en la §3.

## 2. PREGUNTAS · aprobar el paso a operación y arrancar WF-4

- **Q4-01** · ¿Apruebas el paso a operación del sitio tal como está en 8ace01e (plan §3.8)?
  - Opciones (las plantea el orquestador en el encargo; letras de este documento): **A** aprobar · **B** aprobar con cambios (dices cuáles; van a WF-4 o a una ronda acotada) · **C** no aprobar (dices qué falta).
  - **Sin recomendación** escrita en un reporte. Lo que pesa (HECHO, §1): C3 firmada, U5 con 164 de 164 y 0 pares bajo umbral, `verify-wf3` y `verify-all` en verde, S7-07 enviado. En contra: S6-11 sin probar en Android, U5-05 escalado (Q4-82) y WF-3c ter sin revisión independiente.
  - Mientras tanto: el sitio sigue como servidor de desarrollo, sin systemd.
  - Impacto: con A o B arranca WF-4 (Q4-02); con C no hay WF-4.
- **Q4-02** · ¿Arranca WF-4 en cuanto contestes esta puerta? WF-4 es «Ops systemd/mkcert/respaldos/runbook · semilla · seguridad · informe final · revisor final», con 5 agentes (`docs/plan/plan.md` l. 208 y l. 342).
  - Opciones (letras de este documento): **A** sí, al contestar · **B** después, cuando lo digas.
  - **Sin recomendación** escrita en un reporte; el plan pone WF-4 justo después de la puerta 4 (§3.8).
  - Mientras tanto: nada corre como servicio.
  - Impacto: WF-4 toma de esta puerta los valores K31 (§3.6.1), la red (Q4-47), la ruta del token (Q4-44), LaTeX (Q4-45), el candado de la semilla de desarrollo (Q4-79) y la semilla final (Q4-03, Q4-78 a Q4-81). Lo que no contestes entra con su opción PROVISIONAL.
- **Q4-03** · ¿WF-4 carga `datos/semilla/gastos_fijos_ingresos.csv` si existe?
  - Opciones (letras de este documento): **A** sí; y si no existe, carga sin esa parte y deja el importador listo, como dice el plan · **B** WF-4 espera a que lo entregues.
  - **Recomendada: A**, que es la regla de `docs/plan/plan.md` l. 186. HECHO: hoy el archivo no existe.
  - Mientras tanto: la base no trae tus ingresos ni gastos fijos reales.
  - Impacto: con A, WF-4 no se detiene; la correspondencia de columnas es Q4-55 y la plantilla está en `plan.md` §8.2 (l. 408).

## 3. PREGUNTAS · los bloqueos abiertos

Mientras no contestes, el código y los documentos corren con la opción PROVISIONAL de cada archivo, que cita su bloqueo. Ninguno cambia cifras canon de los motores ni sus oráculos (`INDICE.md`, secciones de WF-3a a WF-3c bis).

### 3.1 WF-1 y puerta 2 · lo que quedó abierto (lista corta, con su id Q2 original)

Axel dejó PROVISIONALES las preguntas sin recomendación de la puerta 2 «para preguntarse en la puerta 4» (`docs/estado/puerta2/RESPUESTAS.md`, puerta 2 bis). Las opciones completas están en `PUERTA-2.md` (§2, §3 y §5) y en cada archivo; se contestan con el id Q2. Siguen abiertos por ellas `S2-ventanas`, `S2-flujo-quincenal`, `S5-carrito-comparativa`, `S5-carritos-detalle` y `S5-profeco-filtro` (`INDICE.md`, «Bloqueos nuevos de WF-1», y `INTEGRACION-wf3c.md` §15.6.3).

| Q2 | Archivo | Tema | En uso (PROVISIONAL) |
|---|---|---|---|
| Q2-01 | `S2-ventanas.md` (pregunta 4) | Devengo con ventanas: caso 1, hueco que cruza un devengo; caso 2, ventanas contiguas que no duran un mes | caso 1: 1A; caso 2: 2A, sin recomendación |
| Q2-07 (variante) | `S4-presentacion-respaldo.md` | Presentación sin precio | se pregunta como Q4-06 |
| Q2-12 | `S2-flujo-quincenal.md` (1) | Gasto variable estimado de una quincena futura | sin respuesta; el motor lo da en null y el flujo lo pone en 0 (`S6-presupuesto-motor` punto 4) |
| Q2-48 | `S5-carrito-comparativa.md` (2) | Fuente de Bodega | A |
| Q2-50 | `S5-carritos-detalle.md` (1) | Ganador del sustituto | A |
| Q2-51 | `S5-carritos-detalle.md` (2) | `missing_count` | A |
| Q2-52 | `S5-carritos-detalle.md` (3) | `stale_count` | A |
| Q2-56 | `S5-profeco-filtro.md` (1) | Municipios | A |
| Q2-57 | `S5-profeco-filtro.md` (2) | Comparación de nombres | A |
| Q2-58 | `S5-profeco-filtro.md` (3) | Sucursal del precio | A |
| Q2-85 | — | U3 · supuesto 5 (ya retirado) | como está |
| Q2-87 | — | Recetas · etiquetas | como está |
| Q2-88 | — | Recetas · kcal de cortes con hueso | como está |
| Q2-89 | — | Recetas · recetas parecidas | como está |
| Q2-90 | `S4-empates-sin-confirmar.md` | Recetas · 40 empates de `revisar.tsv` | se pregunta también como Q4-13 |
| Q2-91 | — | PROFECO · cadenas del catálogo (A o B) | como está |
| Q2-94 | — | Lista · subtotal «falta por comprar» | como está |

HECHO sobre el conteo: `RESPUESTAS.md` habla de «las 28 sin recomendación» (98 preguntas menos las 70 recomendadas). En la hoja de `PUERTA-2.md` §8 y en `INDICE.md` encuentro 19 ids sin recomendación única: estos 17, más Q2-11 y Q2-64, que contestaste una por una. No pude reconstruir la diferencia con los archivos. **PREGUNTA para el orquestador**, no para Axel.

### 3.2 P2c · motores y contratos (13 bloqueos, `INDICE.md` «Bloqueos nuevos de P2c»)

Ninguno trae recomendación de la puerta; donde un agente dejó una PROPUESTA, la anoto.

- **Q4-04** · `S2-ahorro-quincenal.md` · El reparto mitad y mitad del ahorro mensual planificado ya es FIRME (Q2-13 = B). Falta: (1) ¿se reparte meta por meta o sobre la suma del mes? (2) ¿qué metas cuentan en una quincena futura? (3) ¿el ahorro planificado resta del saldo o solo se informa?
  - Opciones: (1) A sobre la suma del mes · B meta por meta. (2) A todas las metas activas con `planned_monthly_cents` · B además solo del mes de `start_date` al de `target_date` y sin las cumplidas. (3) A solo se informa · B resta del saldo como egreso comprometido. **Sin recomendación.**
  - PROVISIONAL: `savings_planned_cents` sigue en null y no entra a ningún saldo (`S2-interfaz.md` 1.3 §3.11). Impacto: motor S2 (M-02 y M-10).
- **Q4-05** · `S2-metas-vencida.md` · ¿Qué aportación muestra Metas para una meta vencida que no se cumplió (`months_left` 0 y falta dinero)? El motor y el oráculo difieren.
  - Opciones: A todo lo que falta, en el mes de hoy (el motor) · B sin aportación calculada: `needed_monthly_cents` null y `plan` vacío (el oráculo) · C `needed_monthly_cents` = lo que falta y `plan` vacío. **Sin recomendación.**
  - PROVISIONAL: A; el oráculo se alinea a ella. Impacto: pantalla de Metas.
- **Q4-06** · `S4-presentacion-respaldo.md` (variante de Q2-07) · Sin precio en S5, ¿qué presentación rige un artículo de la lista cuando (a) viene de varias recetas con presentación distinta, (b) es solo de reposición por ritmo o (c) es «calculado» con necesidad 0 y sin ingrediente esa semana?
  - Opciones: (i) la presentación típica del catálogo · (ii) para (a) la de la primera receta en orden de fecha y tiempo, y para (b) y (c) además i o iii · (iii) null en presentación, paquetes y sobrante, con el producto en `missing_prices[]`. **Sin recomendación.**
  - PROVISIONAL: la regla del motor: ii para (a), luego i, luego la primera presentación del catálogo (`S4-interfaz.md` 1.5 §6.2 paso 5). Falta alinear el oráculo y una fixture con los tres casos.
- **Q4-07** · `S5-carritos-granel.md` · Con cantidades a granel, ¿el importe de cada renglón se redondea como en el ticket o se suma exacto?
  - Opciones: A exacto: un solo redondeo HALF_UP de la suma (el motor) · B como el ticket: cada línea HALF_UP y el total es la suma de las redondeadas (el oráculo) · C otra regla por línea, documentada con un ticket real. **Sin recomendación.**
  - PROVISIONAL: A (`S5-interfaz.md` 1.5 §7.2 paso 8). Falta alinear el oráculo y una fixture a granel.
- **Q4-08** · `C1-hora-reportes-job.md` · ¿Qué manda en la hora y el día de un reporte programado: las claves `reports.*` de `app_setting` o `job.schedule`?
  - Opciones: A mandan las claves y la API reescribe `schedule` · B manda `job.schedule` y las claves solo informan · C quitar `schedule` de los jobs de reportes (migración nueva). **Recomendada: A**, PROPUESTA del agente `C1-puerta2` (el archivo).
  - PROVISIONAL: `updateJob` acepta `schedule` en todos los jobs. Se cruza con Q4-38 (ocultar `schedule` de los jobs de avisos).
- **Q4-09** · `C1-reparto-alcance.md` · ¿Qué gastos recalcula `CostShareRuleUpdate.scope = all`?
  - Opciones: A solo los que nombran la regla · B todos los que la usan, también por categoría o por omisión · C B limitado a un rango de fechas. **Recomendada: B**, PROPUESTA del agente `C1-puerta2` (el archivo).
  - PROVISIONAL: el alcance de `all` en C1-datos §7.2 y §13. Impacto: motor de reparto y U4 P-26.12.
- **Q4-10** · `S1-reparto-edicion-gasto.md` · Tras cambiar una regla con «solo los nuevos», si editas un gasto anterior, ¿con qué partes se recalcula?
  - Opciones: A la proporción de las partes congeladas · B guardar cómo se calculó cada parte (migración) · C con la regla vigente. **Sin recomendación** («esta ronda no las recomienda»).
  - PROVISIONAL: C, que es lo que hace `updateExpense`.
- **Q4-11** · `S3-formula-tarjeta-pdf.md` · Falta que entregues 2 PDF de estado de cuenta (una tarjeta y un crédito de tienda) para escribir la fórmula de intereses de cada producto. La opción (a) ya es respuesta tuya.
  - **Sin opciones nuevas**: solo falta la entrega (también en §4). **Sin recomendación.**
  - PROVISIONAL: la fórmula general de la opción (b) de `S3-formula-tarjeta` en `S3-interfaz.md` 1.5.
- **Q4-12** · `S3-iva-lectura.md` · (1) Un préstamo sin el dato de IVA, ¿qué hace? (2) Un crédito «no incluye IVA» sin porcentaje, ¿qué hace?
  - Opciones: (1) A no suma IVA, porque lee «la tasa ya incluye IVA» (lo que hace el motor) · B responde «falta capturar», y U4 y C1 tienen que capturarlo y avisarlo · C la API no deja guardar un préstamo sin ese dato. (2) A «falta capturar», y C1 lista `vat` en `missing_terms` (cambia RN-26 1.4) · B la API no deja guardar 0 sin porcentaje (cambia el esquema o C1). **Sin recomendación.**
  - PROVISIONAL: (1) A; (2) el motor responde «falta capturar» y no está decidido si C1 lo avisa.
- **Q4-13** · `S4-empates-sin-confirmar.md` (Q2-90) · ¿Confirmas o corriges los empates de nutrición y, cuando exista, la lista de ligas de PROFECO?
  - Opciones (letras de Q2-90): A confirmar tal cual · B corregir (dices cuáles). **Sin recomendación**: faltan datos que solo tú tienes.
  - PROVISIONAL: los renglones «propuesto» (`S4-interfaz.md` §9.1) y las fixtures `s4-kcal-*`. Se cruza con Q4-77 (kcal en 0).
- **Q4-14** · `S5-motor-oraculo-reglas.md` · Reglas en que el motor S5 y su oráculo difieren. El archivo tiene 8 puntos: (1) qué es un «renglón» para contar piezas; (2) qué precio toma lo emparejado si la canasta trae la misma presentación a precios distintos; (3) si las referencias del precio vigente llevan una por fuente o todas; (4) formatos de fecha de PROFECO; (5) liga a una presentación desconocida; (6) tope de rechazos; (7) gramática del precio; (8) liga repetida.
  - Opciones: en cada punto, **A** lo que hace el motor (en uso) y **B** lo que hace el oráculo; los puntos 1 y 2 tienen además **C** otra regla. Los textos completos están en el archivo. **Sin recomendación.**
  - PROVISIONAL: A en todos; el código de los dos lados cita el bloqueo.
- **Q4-15** · `S5-precio-vigente-captura.md` · Cuatro reglas que `S5-interfaz.md` 1.5 no fija: (1) desempate con el mismo `created_at`; (2) precisión de la fracción de segundo; (3) perfil de fecha y hora de `created_at` (con el año 0000, también en `observed_date`); (4) claves obligatorias de `prices[]`.
  - Opciones: por punto, en el archivo (en el 1, **A** es «el primero de la entrada», lo que hace el motor). **Sin recomendación.**
  - PROVISIONAL: lo que hace el motor. La revisión r3 pidió sumar el año 0000 a la pregunta 3.
- **Q4-16** · `U4-reparto-preseleccion.md` · En el diálogo de alcance de P-26.12, ¿qué opción aparece preseleccionada?
  - Opciones: A «Solo a los gastos nuevos» · B «También a los ya capturados» · C ninguna. **Recomendada: A**, PROPUESTA de `U4-puerta2` (el archivo).
  - PROVISIONAL: A (U4 1.5, P-26.12 y supuesto 6).

### 3.3 P2d · agente externo de precios y recetas (K53) · `S5-agente-externo.md`

HECHO: K53 dice que la aplicación no lee sitios de tiendas; un agente externo (tu sesión de Claude Code con la extensión de Chrome) mete productos, precios y recetas por la API con una credencial propia. C1 1.4 ya trae `getExternalAgentReference`, `importExternalCatalog`, `importExternalPrices` e `importExternalRecipes`, marcadas PROVISIONALES. Lo PROVISIONAL es PROPUESTA del orquestador para (a) a (d) y de `k53-aplicar` para (e) (`INTEGRACION-wf3c.md` §15.6.8).

- **Q4-17** · (a) ¿Los precios del agente compiten en la comparativa como `manual` (mismo peso y frescura que el ticket) o son solo referencia, como PROFECO (se muestran y nunca ganan)?
  - Opciones (letras de este documento): **A** compiten como `manual` · **B** solo referencia. **Recomendada: A**, PROPUESTA del orquestador (el archivo, «Mientras tanto»).
  - Impacto: con A, el motor S5 y la fixture `s5-fuentes-compiten` cambian su regla: hoy tratan `site_assisted` como referencia, y los 2 precios más baratos de esa fuente pasarían a competir.
- **Q4-18** · (b) ¿Qué puede escribir el agente?
  - Opciones (letras de este documento): **A** solo productos, presentaciones, precios y recetas · **B** también inventario (existencias y ajustes). **Recomendada: A**, PROPUESTA del orquestador.
- **Q4-19** · (c) ¿Las recetas del agente entran como borrador que tú o Ximena confirman, o directas al recetario?
  - Opciones (letras de este documento): **A** borrador · **B** directas. **Recomendada: A**, PROPUESTA del orquestador.
  - Impacto: WF-3 no pudo construir `importExternalRecipes` sin (c) y (f).
- **Q4-20** · (d) ¿Qué valor guarda `price.source` para un precio del agente? El CHECK solo admite `ticket_manual`, `ticket_ocr`, `site_assisted`, `profeco` y `manual`.
  - Opciones: **d1** S1 agrega `external_agent` con una migración y retira `site_assisted` · **d2** se reutiliza `site_assisted` y solo la documentación lo llama «agente externo» · **d3** se guarda como `manual` y la referencia del agente va en `source_ref`. **Recomendada: d2**, PROPUESTA del orquestador («no inventes un valor nuevo»).
- **Q4-21** · (e) ¿Dónde vive la credencial de servicio del agente y a nombre de quién quedan su clave de idempotencia y la bitácora?
  - Opciones: **e1** un archivo fuera del repositorio que la API compara al ejecutarse (como el token del bot) · **e2** una tabla nueva con el SHA-256 de la credencial (migración) · **e3** sin credencial propia: el agente entra con la cuenta de Axel. Y a nombre de: **e4** la cuenta de Axel · **e5** una cuenta de servicio nueva. **Recomendada: e1 + e4**, PROPUESTA de `k53-aplicar`.
  - PROVISIONAL: C1 1.4 solo declara el transporte `serviceToken` (`Authorization: Bearer`); sin verificador, las 4 operaciones responden 401 (`S6-valores-operativos-base`).
- **Q4-22** · (f) ¿Cómo se guardan el borrador y el origen de una receta del agente? `recipe.origin` solo admite `claude_seed` y `household`, y `recipe` no tiene estado de borrador.
  - Opciones: **f1** S1 agrega `external_agent` a `origin` y un estado de borrador · **f2** borrador = `is_active = 0` con origen `household` y la marca del agente en `notes` · **f3** si (c) = directa, solo cambia `origin`. **Sin recomendación.**
  - HECHO: `S6-compras-b-nutricion` punto 5 ya implementó f2 de forma PROVISIONAL, sin marca del agente en `notes`.
- **Q4-23** · (g) ¿Tope de renglones por lote?
  - Opciones: **g1** 500 productos, 500 precios y 100 recetas por lote · **g2** otros topes que des tú · **g3** sin tope; el lote solo pide al menos un renglón. **Sin recomendación**: g1 es un supuesto de `k53-aplicar`, no una PROPUESTA.
  - PROVISIONAL: g1, marcado en C1 1.4.

### 3.4 Pista 3 · sistema de diseño (2 bloqueos)

- **Q4-24** · `P3-contraste-chip-hover.md` · Cuatro marcas de color de categoría quedan bajo 3:1 sobre el fondo de hover del chip seleccionable. ¿Qué se hace?
  - Opciones: A el hover cambia solo el borde, como las tarjetas · B se deja y la marca se declara decorativa · C otra. **Recomendada: A**, del agente `ds-brandbook` (`INDICE.md`, Pista 3).
  - PROVISIONAL: DESIGN §7.4 tal cual, con excepción declarada en el verificador.
- **Q4-25** · `P3-huecos-componentes.md` · Huecos de DESIGN.md resueltos como supuestos con tokens existentes (iconos de 4 secciones de la barra lateral, bloque «Antes → Después», perilla del interruptor, círculo de «hoy», comandos del menú «/», pesos y anillo de 2 px, y más). Preguntas del archivo: ¿se confirman los iconos? ¿Se agrega una foto sintética de ticket a la maqueta? ¿Basta `CalendarioPagos` o hace falta un `Calendario` propio?
  - **Sin opciones con letra**; **sin recomendación**. Los iconos se deciden con Q4-104.
  - PROVISIONAL: todos los supuestos del archivo.

### 3.5 Semilla de WF-3a (3 bloqueos)

- **Q4-26** · `S6-semilla-despensas-futuras.md` · ¿Las despensas 2, 3 y 4 (10-dic-2026, 10-feb-2027 y 10-abr-2027) se siembran como compras a crédito activas, o las crea la proyección?
  - Opciones: A las 4 compras con sus 16 pagos de `motor_cred` · B solo la despensa 1 · C A con `status = 'cancelled'`. **Sin recomendación.**
  - PROVISIONAL: A (la deuda viva al 10-abr-2027 da $6,293.44, igual que el oráculo). Se cruza con Q4-108 (deuda viva de hoy).
- **Q4-27** · `S6-semilla-sucursales-municipio.md` · ¿En qué alcaldía o municipio y estado están Bodega Aurrera Vallejo y Walmart Acueducto de Guadalupe?
  - Opciones: A las capturas tú en Ajustes · B las dictas y van a la semilla · C suponerlas (descartada). **Sin recomendación.**
  - PROVISIONAL: sucursales sin municipio, así que 0 precios de la muestra PROFECO cargados.
- **Q4-28** · `S6-semilla-catalogo-alias.md` · ¿Dónde viven los sinónimos y los nombres de ticket del catálogo?
  - Opciones: A se leen de `productos.json` · B migración nueva con tablas de alias · C no usarlos. **Sin recomendación.**
  - PROVISIONAL: no se guardan en la base. Impacto: OCR y búsqueda.

### 3.6 WF-3a · construcción (40 bloqueos; `INDICE.md` «Bloqueos nuevos de WF-3a»)

#### 3.6.1 Valores operativos K31 PROVISIONALES

La lista consolidada está en `INTEGRACION-wf3b.md` §9 y `INTEGRACION-wf3c.md` §15.6.1; WF-3b, WF-3c y WF-3c bis no agregaron valores K31. Cambiar cualquiera es una constante o una opción, salvo donde se dice.

- **Q4-29** · `S6-respaldos.md` punto 1 y `S6-valores-operativos-base.md` punto 1 · ¿A qué hora se hace el respaldo diario y cuántos se guardan?
  - Opciones: el archivo no pone letras al valor; `S6-valores-operativos-base` da para todo su paquete (puntos 1 a 7): **(a)** aprobar la propuesta · **(b)** cambiar valores puntuales · **(c)** mover a `app_setting` la hora del respaldo, los límites de cuerpo y la purga (pide una migración nueva y C1). **Recomendada: (a)**, la PROPUESTA de `auth-sistema` y `backend-base`: diario a las **03:15** (America/Mexico_City), después del resumen de las 21:30 y antes del reporte de las 08:00, y se guardan **30** diarios.
  - PROVISIONAL: 03:15 y 30. Los respaldos manuales y los previos a una restauración (`pre_migration`) no se borran solos. Impacto: WF-4 lo copia a `ops/runbook.md`.
- **Q4-30** · `S6-respaldos.md` punto 2 · Al restaurar un respaldo, ¿qué se conserva del estado previo?
  - Opciones: **a)** todo se reemplaza salvo `backup` (el catálogo de respaldos) y `device_session` (las sesiones) · **b)** restaurar todo tal cual · **c)** conservar también `user_account` y `login_attempt` (si no, una contraseña cambiada después del respaldo vuelve a la anterior). **Sin recomendación** (el archivo lo deja como PREGUNTA).
  - PROVISIONAL: a). Además, respaldar y restaurar son síncronos (la base pesa ~1 MB) y un respaldo en curso da 409 a otro manual.
- **Q4-31** · `S6-adjuntos.md` punto 2 · ¿Qué tipos de archivo se admiten como adjunto?
  - Opciones: **a)** JPG, PNG, WebP y PDF · **b)** agregar HEIC. **Sin recomendación** (en uso a), porque es lo que declara C1).
  - PROVISIONAL: a), detectado por el contenido y no por el nombre; otro tipo responde 415. Límite de 10 MB.
- **Q4-32** · `S6-adjuntos.md` puntos 1, 3 y 4 · (1) ¿Cuántos bytes es «10 MB»? (3) ¿Qué `kind` lleva una subida que no lo dice? (4) ¿Cuándo da 409 borrar un adjunto?
  - Opciones: (1) a) 1,048,576 bytes · b) 1,000,000 bytes. (3) a) `other` · b) `ticket_photo` · c) obligatorio (422). (4) a) se borra y deja en null la referencia del precio o del estado de cuenta · b) 409 `in_use` si alguien lo usa. **Sin recomendación.**
  - PROVISIONAL: a) en los tres. Con (1) b) el límite baja 4.6 %.
- **Q4-33** · Cabeceras de seguridad y CSP · `S6-valores-operativos-base.md` punto 3, `WF3-frontend-servidor.md` punto 4 y `S6-sistema.md` punto 8.
  - Valor PROPUESTO: `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, `X-Frame-Options: DENY` y `Cache-Control: no-store` en la API; en la PWA, CSP `default-src 'self'; img-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; connect-src 'self'; font-src 'self'; worker-src 'self'; manifest-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'` ('unsafe-inline' en estilos por Radix, dnd-kit y ECharts); descargas con `Content-Disposition: attachment` y `Cache-Control: private, no-store`; **sin HSTS** mientras el certificado sea de mkcert.
  - Opciones: (a), (b) y (c) de Q4-29. **Recomendada: (a)**, PROPUESTA de `backend-base`.
  - HECHO: U5-08 vio 0 hosts externos en 164 cargas. WF-4 tiene en su alcance la seguridad (plan l. 342).
- **Q4-34** · `S6-valores-operativos-base.md` puntos 2, 4, 5, 6 y 7 · argon2id t=3, m=64 MiB, p=4 (~0.1 a 0.2 s por inicio de sesión); 8 MiB por cuerpo JSON (413) y hasta 32 MiB en memoria con Idempotency-Key; purga de idempotencia cada hora al minuto 5; rutas por omisión de `finanzas serve` (las fija WF-4 en `ops/runbook.md`; el binario rechaza 0.0.0.0 y ::); un campo desconocido en un JSON da 400.
  - Opciones: (a), (b) y (c) de Q4-29. **Recomendada: (a)**, PROPUESTA de `backend-base`.
  - PROVISIONAL: esos valores.
- **Q4-35** · Contenido de los resúmenes de Telegram (`INTEGRACION-wf3b.md` §9, bloqueo `S7-fuentes-avisos`) · ¿Qué lleva cada resumen?
  - **Sin opciones escritas; sin recomendación.** PROVISIONAL: el **diario** lleva fecha, gastos de hoy, gastos del mes y pagos de los próximos 7 días; el **semanal**, semana ISO, gastos de la semana y del mes, pagos de 7 días y casillas de comidas sin armar. Lo que llegó a tu chat a las 21:15 es el diario.
  - Impacto: cambiar el contenido toca `backend/internal/reportes`.
- **Q4-36** · `S7-fuentes-avisos.md` §1 · ¿Qué es el «día de despensa» del que se avisa con anticipación de precios viejos?
  - Opciones: **A** el lunes que abre cada lista de compras (semana K8) · **B** los días de pago 15 y 30 (fin de mes en febrero) · **C** una clave nueva `alerts.grocery_weekday` en `app_setting` (migración nueva, K44). **Recomendada: C**, de `s7-avisos` (el archivo).
  - PROVISIONAL: la fuente responde «sin fuente», así que la alerta de precios viejos **no sale**; queda anotada en el registro. HECHO, visto hoy en `backend/internal/reportes/avisos/fuente.go`: también la alerta de presupuesto (80 % y 100 %) y la de precios viejos esperan funciones de otras áreas y hoy responden «sin fuente» (§2 del archivo, solicitudes fuera de alcance). Ninguna inventa cifras.
- **Q4-37** · `S7-fuentes-avisos.md` §3 · El día de revisar el plan de comidas (`alerts.meal_plan_missing_weekday`), ¿qué semana revisa la alerta 4?
  - Opciones (letras de este documento): **A** la que empieza el lunes siguiente · **B** la semana en curso. **Recomendada: A**, PROPUESTA de `s7-avisos`.
- **Q4-38** · `S7-fuentes-avisos.md` §4 · Cadencia y reintentos de los avisos.
  - Puntos: (1) los jobs del resumen diario, el semanal y las alertas revisan cada minuto y comparan con su hora en Ajustes; PROPUESTA: ocultar `schedule` de esos 3 jobs en Ajustes o dejarlo de solo lectura. (2) Un envío fallido queda `failed` y no se reintenta solo; la alternativa escrita es un reintento a los 10 minutos. (3) PREGUNTA: ¿se omite el aviso de un pago que ya liquidaste antes de su fecha? (sin opciones escritas; hoy la alerta no cruza los pagos hechos).
  - Opciones (letras de este documento): (2) **A** sin reintento · **B** un reintento a los 10 minutos. **Recomendada:** (1) la PROPUESTA de `s7-avisos`; (2) A, PROPUESTA del mismo agente; (3) sin recomendación.
- **Q4-39** · `S7-reporte-mensual-secciones.md` · (1) ¿C1 trae el devengado contra efectivo dentro del reporte mensual? (2) ¿Cómo comparte el reporte los cálculos de otras áreas? (3) ¿Valen los valores de esta versión? (4) ¿Cómo se escribe un negativo?
  - Opciones: (1) **A** C1 agrega `accrual_vs_cash` a `MonthlyReport` · **B** la vista web llama además a `getAccrualVsCash` (no sirve para meses cerrados). (2) **A** cada área exporta su constructor · **B** un `api.Registry.Invoke` interno. (3) sin opciones escritas: ingreso del mes = suma de ingresos del mes; la quincena de créditos en un mes cerrado es la segunda; intereses pagados = los de pagos ya cubiertos; despensa devengada = la de la persona en el mes. (4) letras de este documento: **A** `-$1,197.64` (en uso) · **B** `$-1,197.64` (lo que escribe S7-12). **Recomendada:** (1) A, de `s7-pdf`; (2), (3) y (4) sin recomendación.
  - PROVISIONAL: mientras no haya (2), el PDF mensual **no muestra plan de comidas ni canasta**. Las cifras que sí muestra cruzan con la API con 0 diferencias (S7-02, S7-11).
- **Q4-40** · `WF3-frontend-k31.md` · Valores de la PWA sin conexión: (1) sondeo de `/health` cada 15 s; (2) reintento de la cola cada 30 s; (3) la copia de lo ya cargado se guarda hasta cerrar sesión; (4) la cola no se borra al cerrar sesión.
  - **Sin opciones con letra**; la única alternativa escrita es para (3): caducar a los 7 días. **Recomendada:** los 4 valores, PROPUESTA de `frontend-base`.
  - HECHO: U5-13 y `cola.spec.ts` ejercieron la cola y el sondeo.
- **Q4-41** · `WF3-frontend-servidor.md` · ¿Cómo entrega el servidor la PWA en producción?
  - Opciones: **(A)** `dist/` incrustado en el binario · **(B)** la carpeta `frontend/dist` leída al ejecutarse. **Recomendada: A para producción y B en desarrollo**, de `frontend-base`. Los puntos 1 a 3 del archivo (mismo origen, respaldo de SPA, `Cache-Control` y tipo del manifiesto) son PROPUESTA.
  - HECHO: hoy corre B (servidor de desarrollo). En WF-3a el manifiesto salía como `text/plain` (hallazgo menor, `INTEGRACION-wf3a.md` S6-12). Impacto: decide cómo empaqueta WF-4.
- **Q4-42** · `S6-sistema.md` puntos 1 a 7 y 9 · Decisiones de autenticación y sistema.
  - Puntos (PROVISIONAL → alternativa escrita): (1) C1 deja `auth.password_min_length` y `auth.max_failed_attempts` en 0, y K6 pide 12 o más → subir el mínimo de C1 a 12 y 1. (2) Cambiar la contraseña no cierra las otras sesiones → cerrar todas salvo la actual. (3) El usuario distingue mayúsculas → compararlo sin mayúsculas. (4) Estado de la semilla por la bitácora → tabla nueva. (5) Ruta de la CA por opción o variable → campo en la configuración. (6) IA (Ollama y Claude) instalada y **apagada**; OCR con IA y recetas con IA responden «no configurado» → decidir modelo y llave antes de encenderla. (7) Orden de `listJobs` por `code` → por `next_run_at`. (9) `finanzas.local` en Android → entrada en el DNS del router e IP fija (se decide en Q4-47).
  - Opciones (letras de este documento): en cada punto **A** lo PROVISIONAL · **B** la alternativa escrita. **Sin recomendación.**
  - Impacto: (1) es seguridad: hoy el sistema aceptaría una contraseña corta si se cambia el parámetro.

#### 3.6.2 Telegram, PDF y red

- **Q4-43** · `S7-convivencia-multix.md` · ¿Quién le pasa al servicio el mensaje de Ximena, si MultiX ya consume las actualizaciones del bot (la API de bots entrega cada actualización a un solo lector)?
  - Opciones: **A** MultiX reenvía el `Update` de un chat desconocido (o que empiece con `/finanzas`) a una carpeta local o a un endpoint local del servicio · **B** código de enlace de 6 dígitos que Ximena manda al bot, reenviado como en A (pide C1 nuevo) · **C** el servicio lee la bitácora de MultiX (frágil) · **D** un segundo bot solo para Finanzas (contradice K30). **Recomendada: A**, de `s7-avisos` (el archivo).
  - PROVISIONAL: nadie entra solo a la lista de candidatos; tu chat se registró a mano con tu OK (`ops/telegram-disparo.md`). Impacto: Ximena no recibe avisos hasta decidir; es la acción 5 de la §4.
- **Q4-44** · `S7-ruta-token.md` · ¿Cuál es la ruta real del token en `~/.openclaw/openclaw.json`: `bots.multix.botToken` (K30) o `channels.telegram.bots.multix.botToken`?
  - Opciones: **A** dejar las dos · **B** dejar solo la que confirmes. **Recomendada: A** hasta WF4-01, de `s7-avisos`.
  - HECHO: el envío de las 21:15 salió con estado `sent` y HTTP 200, así que una de las dos rutas funciona.
- **Q4-45** · `S7-lualatex.md` · `lualatex` no carga fuentes OpenType porque falta el paquete `texlive-luatex`.
  - Opciones: **A** instalas `sudo pacman -S texlive-luatex` (unos 15 MB) y el preámbulo pasa a Inter OpenType · **B** quedarse con Inter Type 1, que ya funciona. **Recomendada: A**, de `s7-pdf`.
  - PROVISIONAL: B. Impacto: solo calidad tipográfica; S7-01, S7-12 y S7-13 pasan con B.
- **Q4-46** · `S7-estado-cuenta-fechas.md` · El estado de cuenta de Karen necesita la fecha de cada pago y C1 no la trae en `ExternalStatement.payment_shares`.
  - Opciones: **A** C1 1.5 agrega `scheduled_payment_id` y `due_date` · **B** seguir con la cita cruzada a `getMonthlyReport`. **Recomendada: A**, de `s7-pdf`.
  - PROVISIONAL: B; el cruce de S7-03 da 0 diferencias en 20 de 20 cifras y fechas.
- **Q4-47** · `WF3-integrador-red.md` · (1) IP fija de la máquina (hoy 192.168.50.56 por DHCP); (2) nombre `finanzas.local` permanente (hoy lo publica un `avahi-publish` que se pierde al reiniciar); (3) CA de mkcert en el almacén del sistema.
  - Opciones: (1) **A** reserva DHCP en el módem · **B** IP estática en la máquina · **C** dinámica y reemitir el certificado al cambiar. (2) **A** alias fijo de avahi en `/etc/avahi/hosts` (sudo) · **B** unidad systemd con `avahi-publish` · **C** entrada DNS en el módem · **D** `/etc/hosts` (sudo, solo esta máquina). (3) **A** dejarla en el almacén del usuario · **B** instalarla también en el sistema (sudo). **Recomendada: 1A, 2C (y A para la propia máquina; si el módem no deja, los teléfonos usan `https://192.168.50.56:8443/`), 3A**, PROPUESTA de `integrador-backend`.
  - PROVISIONAL: 192.168.50.56 y `avahi-publish` en `ops/runbook.md` §1–§2 y `ops/mkcert.md` §2. Impacto: S6-11 y S6-14; lo cierra WF-4 o tú.

#### 3.6.3 API · reglas que C1 no fija

- **Q4-48** · `S6-ajustes-api.md` · Seis puntos de Ajustes.
  - Opciones: (1) `PlannedItem.next_date`: **A** exportar `NextOccurrence` en `motor_fin` sin cambiar cifras · **B** dejarlo siempre ausente · **C** calcularlo en la API (descartada). (2) `updateCostShareRule` con `scope = all`: **A** exportar `SplitRule` en `motor_fin` y recalcular con el alcance de Q4-09 · **B** responder 409 hasta entonces. (3) `reorderCategories`: **A** todas las categorías · **B** solo las activas. (4) Orden por omisión (sin opciones escritas). (5) Metas desactivadas: **A** sin cifras · **B** pedir al motor también las inactivas. (6) PATCH vacío o fuera del enum = 400 (sin opciones escritas). **Sin recomendación.**
  - PROVISIONAL: 1B (P-26.02 y P-26.03 no muestran «próxima»), **2B: hoy un cambio de reparto no se puede aplicar a gastos ya capturados** (responde «Elige «Solo los nuevos»»), 3B, 5A.
- **Q4-49** · `S6-comidas-calculos.md` · Cálculos de comidas y salud: (1) C1 pide enteros en costo y kcal y el motor S4 los deja ausentes; (2) ventana de «últimas 8 semanas»; (3) semanas y promedio de `getMealPlanTrends`; (4) estado del día en el mes; (5) `copyPreviousWeek` sin periodo.
  - Opciones: (2) **A** de hoy − 55 días a hoy · **B** las 8 semanas ISO completas anteriores · **C** incluir lo planeado a futuro. (3) promedio **A** suma de kcal ÷ 7 · **B** ÷ días con plan. (5) **A** semana por omisión (K38) · **B** 400. Los puntos 1 y 4 traen su regla en el archivo. **Sin recomendación.**
  - PROVISIONAL: 2A, 5A y lo que el archivo marca en 1, 3 y 4.
- **Q4-50** · `S6-compras-a-motores.md` · (1) el motor S4 en Go no trae la cotización por elección de S4 1.5 (M-r5-1); (2) quién tiene plan de comidas; (3) a qué «vertical» va un artículo de texto libre; (4) ¿compite una tienda dada de baja o una presentación inactiva?
  - Opciones: (1) **a** actualizar el motor S4 a 1.5 · **b** dejarlo (ver archivo). (2) **a** columna nueva · **b** regla por participación. (3) **a** siempre cadena · **b** según la tienda elegida · **c** grupo propio. (4) sin opciones escritas. **Sin recomendación.**
  - PROVISIONAL: (2) Axel y Ximena sí, Karen no; (3) mercado si la tienda elegida es mercado, si no cadena; (4) tiendas inactivas fuera y presentaciones inactivas dentro.
- **Q4-51** · `S6-compras-b-fuentes.md` · Fuentes de precio tras K53: (1) veredicto de cada cadena con sitio; (2) fecha de la última importación de PROFECO; (3) RAR en la importación de PROFECO; (4) dónde vive `vinculos.csv`; (5) etiqueta del corte; (6) bitácora de la importación.
  - Opciones: por punto, (a), (b) y a veces (c) en el archivo; por ejemplo (3) **(a)** RAR responde 415 y se sube el CSV · **(b)** descomprimir RAR (pide un paquete); (4) **(a)** `<DataDir>/profeco/vinculos.csv` · **(b)** `datos/semilla/profeco/vinculos.csv` · **(c)** tabla nueva. **Sin recomendación.**
- **Q4-52** · `S6-compras-b-nutricion.md` · (1) qué archivo y llave usa la recarga de nutrición; (2) dato preferido nuevo cuando ya hay uno; (3) pasar un artículo a ritmo de consumo sin fecha de conteo; (4) pasar a existencias con ritmo guardado; (5) recetas del agente externo (= Q4-19 y Q4-22); (6) orden permitido en colecciones.
  - Opciones: por punto en el archivo; por ejemplo (2) **(a)** 409 `already_exists` · **(b)** el nuevo pasa a preferido; (4) **(a)** 422 si no se mandan los campos del ritmo en null · **(b)** la API los borra sola. **Sin recomendación.**
- **Q4-53** · `S6-creditos.md` · Decisiones PROVISIONALES de la API de créditos: (1) alta del gasto ligado dentro de la compra a crédito; (2) ventana de la tendencia de deuda («12 meses») e intereses pendientes; (3) «saldo» de una tarjeta y abonos a un corte; (4) el comparador sin partes ni canastas; (5) horizonte de la serie de disponible A/B; (6) compra con promoción.
  - Opciones: **(a)** lo PROVISIONAL de cada punto y **(b)** su alternativa, en el archivo (por ejemplo (2) (a) mes en curso y 11 hacia adelante; (5) (a) de hoy al último vencimiento · (b) una ventana fija). **Sin recomendación.**
- **Q4-54** · `S6-gastos-inicio.md` · Huecos de Gastos e Inicio (18 secciones). Las que piden tu respuesta: §1 reparto RN-03 fuera del motor (**A** exportar `motor_fin.SplitRule` · **B** dejarlo en dominio con un caso del oráculo); §2 confianza de los renglones del OCR (**A** columna nueva y umbrales en Ajustes · **B** umbrales fijos 85 y 60 · **C** dejarla baja); §3 promedio mensual (**A** así · **B** omitirlo en periodos de menos de un mes · **C** días ÷ 30); §5 código de la categoría de papel; §7 si filtrar por una categoría padre incluye sus hijas; §9 «compras particulares» de Karen; §12 tipo `page` en la búsqueda; §14 estado del plan de comidas en el calendario (**A** leer por día · **B** quitar el tipo de C1); §15.4 segunda tienda de «Mejor canasta»; §16 seis reglas de las vistas de base de datos (A, B y a veces C por regla); §17 OCR fallido (**A** así · **B** reintentar N veces); §18 a) inventario al confirmar un ticket (**A** así · **B** crear existencias · **C** tomarlo como conteo nuevo) y b) reparto de la compra a crédito por renglón o capturado aparte.
  - **Sin recomendación**; la PROVISIONAL de cada sección está en el archivo. HECHO: §11 (500 con categorías sin código) ya se cerró en `correccion-backend-core-r1`.
- **Q4-55** · `S6-importar-partidas-csv.md` · ¿Vale esta correspondencia de columnas de `gastos_fijos_ingresos.csv` (plan §8.2)? `tipo` y `periodicidad` a sus claves (sin «semanal»); `monto` en pesos con punto, acepta `$` y comas; `dia_o_fecha` como fecha o como día; `categoria` por nombre o código; `persona` axel, ximena, hogar, karen o conjunto; `porcentaje_conjunto` igual a la parte de Karen de la regla; `notas`.
  - **Sin opciones escritas; sin recomendación.** PROVISIONAL: esa correspondencia. Impacto: lo usa WF-4 si entregas el archivo (Q4-03).
- **Q4-56** · `S6-presupuesto-motor.md` · (1) códigos de las categorías de despensa y papel; (2) cadencia de los meses 7 a 12 de la proyección; (3) fecha por omisión de `getLiveDebt`; (4) gasto variable estimado (= Q2-12); (5) categorías tuyas sin código.
  - Opciones: (1) a) constantes y las alternativas del archivo; (3) a) el último día del mes en curso + 11 · b) la fecha canon 10-abr-2027. **Sin recomendación.**
  - PROVISIONAL: (3) a), por eso `/debt/live` sin fecha responde 2027-09-30; (2) con el interruptor encendido la respuesta lo dice, pero la cadencia todavía no se arma; (4) en 0.
- **Q4-57** · `S6-presupuesto-periodos.md` · ¿Cuándo se llena la tabla de periodos?
  - Opciones: **a)** a pedido, del mes en curso − 11 al + 12, tope 36 meses · **b)** que `migrate` o la semilla llenen un rango fijo · **c)** un job mensual. **Sin recomendación.** PROVISIONAL: a). Impacto: solo el selector de periodo.
- **Q4-58** · `S6-presupuesto-quincenas.md` · ¿Cómo se reparte el plan mensual entre las dos quincenas?
  - Opciones: **a)** por días: Q1 = plan × 15 ÷ días del mes, Q2 = el resto · **b)** mitad y mitad · **c)** una función nueva del motor S2. **Sin recomendación.** PROVISIONAL: a). Impacto: la gráfica de quincenas de P-06.

#### 3.6.4 Pantallas y matriz C3

- **Q4-59** · `C3-matriz.md` · ¿Ratificas las decisiones del orquestador de WF-3a bis que siguen PROVISIONALES: (1) la fila P-26 como «Local (contenedor de P-26.01 a P-26.14)»; (3c) excepción técnica para `queryView`, `listPeriods` y `listIncomes` cuando la pantalla ya cubre la necesidad con otra operación; (4a) «probado» = título de la prueba más 3 reglas escritas, con la API simulada?
  - Opciones: las (a), (b) y (c) de cada punto del archivo; las decisiones son 1 (a), 3 (b) y (c), 4 (a). **Recomendada:** ratificarlas, como pide `INTEGRACION-wf3b.md` §8.4.
- **Q4-60** · `WF3-ajustes-a.md` · Ajustes, secciones 01 a 07: (1) con qué tipo de partida se guarda la despensa programada en P-26.02 y cómo se distingue de un gasto fijo; (2) qué semana usa el punto de equilibrio de la membresía en P-26.06; (3) cómo comparten P-26 las dos áreas de Ajustes.
  - Opciones: (1) **(a)** gasto fijo con `accrual_months = 2` · **(b)** gasto venidero · **(c)** una marca propia en C1. (2) **(a)** la semana en curso de la lista · **(b)** la última semana con lista · **(c)** una operación propia en C1. (3) sin opciones escritas; PROPUESTA en uso: `ajustes-a` registra P-26 y descubre las secciones de cada área. **Sin recomendación.**
  - PROVISIONAL: (1) la despensa es la partida con `accrual_months = 2`; **mientras no decidas, no se puede programar la despensa desde Ajustes** (la semilla y el CSV sí la cargan). (2) Ajustes no evalúa la membresía de Sam's Club; Compras sí.
- **Q4-61** · `WF3-ajustes-compra-ejemplo.md` · Intereses de la compra de ejemplo en las vistas previas 6 y 7 de P-26.05.
  - Opciones: **A** la PROVISIONAL del archivo · **B** usar siempre la tabla de plazos y, sin tabla, «Sin tasa para calcular» · **C** en líneas no revolventes, solo fechas sin montos. **Sin recomendación.**
- **Q4-62** · `WF3-ajustes-muestra-tema.md` · P-26.10 pide en escritorio una muestra de tarjeta en claro y en oscuro a la vez, y los tokens de C2 solo cambian de tema en `:root`.
  - Opciones: **(A)** quedarse así · **(B)** que C2 publique los mismos valores bajo `[data-theme]` en cualquier elemento (cambio en `tokens.css`, dueño C2) · **(C)** quitar la muestra doble de U4. **Recomendada: A**, de `front-ajustes-b` («no toca C2»).
- **Q4-63** · `WF3-ajustes-precios-conteo.md` · Conteo de precios que pasan de frescos a viejos al cambiar el valor en P-26.06.
  - Opciones: **A** la PROVISIONAL del archivo · **B** C1 agrega un parámetro o una vista previa para que el servidor cuente · **C** contar cada precio por producto y tienda (pide también B). **Sin recomendación.**
- **Q4-64** · `WF3-ajustes-registro.md` · Dice «RESUELTO entre áreas en WF-3a» y deja la PROPUESTA A para la base: cómo se registra P-26 cuando la construyen dos agentes.
  - Opciones: **(A)** un cambio de una línea en `Layout.tsx` y las dos áreas exportan solo sus `P-26.nn` · **(B)** quedarse con el despacho de `ajustes-b` · **(C)** otra en el archivo. **Recomendada: A**, de `front-ajustes-b`.
- **Q4-65** · `WF3-comidas-compras-interfaz.md` · Siete puntos de P-10 a P-16: (1) cómo pasa P-13 la semana a P-05; (2) cómo pasa P-14 el total a P-20; (3) color de la tienda sin código en C1; (4) «Sugerir con IA» con la IA apagada; (5) unidad del sobrante con receta; (6) semana en la dirección; (7) historial de precios en fechas sin captura.
  - Opciones: **a** y **b** por punto en el archivo. **Sin recomendación.** PROVISIONAL: 1a, 2a, 3b, 4b, 5a, 6a y 7a.
- **Q4-66** · `WF3-creditos-salud-reportes.md` · Siete puntos de P-17 a P-25: (1) quién es «Hogar»; (2) desmarcar un pago; (3) destinatarios de envíos y avisos; (4) fecha de las partes de Karen (= Q4-46); (5) peso al inicio y al cierre en el reporte; (6) qué se edita de una compra a crédito; (7) rutas propias.
  - Opciones: (1) a) la persona externa es su nombre y titular y co-titular son «Hogar» · b) C1 marca quién es hogar; (2) a) C1 agrega los ids de pago a cada parte · b) una lista de pagos hechos; (3) a) C1 agrega destinatarios; (4) a (6) «C1 lo agrega o U1 lo quita o se ajusta». **Sin recomendación.** PROVISIONAL: lo que dice el archivo en cada punto (por ejemplo, en P-19 solo se editan comercio y nota).
- **Q4-67** · `WF3-frontend-rutas.md` · Direcciones de las pantallas en la barra del navegador.
  - Opciones: **(A)** la propuesta: en español sin acentos (`/entrar`, `/inicio`, `/gastos`…) · **(B)** en inglés (`/expenses`…) · **(C)** con el id de pantalla (`/p/02`). **Recomendada: A**, de `frontend-base` («lo que se ve en pantalla va en es-MX»).
- **Q4-68** · `WF3-gastos-contrato.md` · Huecos de C1 para P-02, P-04 y P-05: (1) id del adjunto en `Expense`; (2) vista previa del reparto al editar el monto; (3) filtros de tipo, forma de pago y «con ticket» en los totales; (4) búsqueda de productos desde el ticket.
  - Opciones: **(a)** y **(b)** por punto en el archivo (por ejemplo (1) (a) agregar `attachments` a `Expense`). **Sin recomendación.** Se cruza con Q4-74.
- **Q4-69** · `WF3-gastos-ticket-credito.md` · Pagar con crédito desde P-05 (ticket).
  - Opciones: **(a)** compra a plazos con capital = total y se piden línea, intereses y plazo · **(b)** lo mismo más frecuencia y primer pago · **(c)** mandar a P-19 después de confirmar el ticket. **Sin recomendación.**
- **Q4-70** · `WF3-inicio-tiendas.md` · El color de cada tienda: C1 no trae un código estable de tienda.
  - Opciones: **A** dejar la regla por nombre · **B** C1 agrega `Store.code` · **C** mapa por id en Ajustes › Tiendas. **Sin recomendación.** PROVISIONAL: A; el nombre siempre va en texto (K41).
- **Q4-71** · `WF3-inicio-vacios.md` · Reglas de «vacío» de Inicio y Detalle y destinos que C1 y U3 no fijan (6 puntos: quincena, metas, comidas, fondo, filtros de «Ver todos los gastos» y destino del calendario).
  - Opciones: **A** la regla PROVISIONAL de cada punto · **B** que C1 agregue un campo (por ejemplo `income_source`, `has_goals`, `planned_slots`) o la alternativa del punto. **Sin recomendación.**

### 3.7 WF-3a bis · corrección del frontend y semilla de desarrollo (4 bloqueos)

- **Q4-72** · `WF3a-bis-ajustes-a.md` · Edición en grupo con `Ctrl+/` en Ajustes: (1) avisos de texto sin cifra en el bloque sumado de P-26.02 y P-26.04; (2) qué campos de varias promociones se editan a la vez en P-26.05; (3) cómo se cambian a la vez los tiempos aptos (selección múltiple) en P-26.07.
  - Opciones: (a) PROVISIONAL y su alternativa en el archivo, por punto (por ejemplo (2) (a) Activa, Línea, Plazo, Compra mínima y Fin de vigencia). **Sin recomendación.**
- **Q4-73** · `WF3a-bis-comidas-compras.md` · (1) P-14 con la semilla de desarrollo no da las cifras canon de U1 §3.16 (Walmart $3,359.00 contra $3,233.00; se decide en Q4-81); (2) P-13 sin lista sembrada; (3) nombre visible de una sucursal («Vallejo» sin su cadena); (4) `Esc` dentro de una hoja cierra la hoja entera.
  - Opciones: por punto en el archivo; (4) (a) PROVISIONAL, ya aplicada. **Sin recomendación.**
- **Q4-74** · `WF3a-bis-gastos-presupuesto.md` · Siete puntos para el dueño de C1: filtros de `queryView` que faltan para N-040, id del adjunto en `Expense` (N-064), galería N-046, N-073, N-053 con préstamos a plazo y tarjetas, `listPeriods` como excepción técnica y detalles menores.
  - Opciones: **(a)** PROVISIONAL y su alternativa, por punto, en el archivo (por ejemplo (5) (a) toda línea activa lleva a P-19 prellenada). **Sin recomendación.** Se cruza con Q4-68.
- **Q4-75** · `WF3a-bis-semilla-dev.md` punto 1 · La renta de la fixture (8,000.00 mensuales desde el 20-oct) y las cifras P1 a P6 de `model.json` no caben juntas en la semilla de desarrollo.
  - Opciones: **(a)** por omisión sin renta; con `--con-gasto-fijo` se agrega · **(b)** siempre con renta, y el cotejo C3-09 de P1 a P6 pasa a «model.json + 8,000.00» · **(c)** renta ligada a un gasto capturado. **Sin recomendación.** PROVISIONAL: (a).
- **Q4-76** · `WF3a-bis-semilla-dev.md` punto 3 · ¿Qué presupuesto quieres para las demás categorías? La única fixture con presupuesto real es Despensa, $2,000.00 desde oct-2026; no se inventaron montos para Comida fuera ni Vivienda.
  - **Sin opciones escritas; sin recomendación.** PROVISIONAL: «Mes contra plan» compara contra $2,000.00 (60 % usado con la semilla).
- **Q4-77** · `WF3a-bis-semilla-dev.md` punto 5 · El costo y las kcal de la semana salen en 0: solo 36 de 160 productos tienen precio y la tabla de nutrición está vacía.
  - Opciones (sin letras en el archivo; letras de este documento): **A** que el agente externo (K53) cargue los precios que faltan · **B** que confirmes filas de `nutricion/tabla.csv` (Q4-13) · **C** aceptar el 0 en las capturas U5. **Sin recomendación.** No se inventó ningún precio ni caloría.
- **Q4-78** · `WF3a-bis-semilla-dev.md` punto 4 · ¿Ratificas que la semilla real agregue los gastos 2 a 6 del oráculo como gastos futuros marcados «(proyectado)» y los ligue a sus compras a crédito? ¿Y que la lista de Gastos muestre esas 5 filas con fecha futura (`semilla-c3.md` l. 39)?
  - Opciones: **(a)** la semilla real los agrega · **(b)** lo hace un overlay de desarrollo · **(c)** C3-09 compara solo P1 y P2. **Recomendada: (a)**, decisión del orquestador (12-oct-2026) aplicada por `semilla-c3`; con ella P1 a P6 cuadran con `model.json` y la deuda viva al 10-abr-2027 da $6,293.44. Las 5 filas: sin opciones escritas.
  - Impacto: el esquema no distingue «proyectado»; si quieres un campo propio, es un cambio de esquema.
- **Q4-79** · `WF3a-bis-semilla-dev.md` punto 2 · Candado para que la semilla de desarrollo nunca corra en producción.
  - Opciones (sin letras en el archivo; letras de este documento): **A** la variable `FINANZAS_ENTORNO=produccion` (en uso) · **B** una marca `entorno` en `app_setting` · **C** una etiqueta de compilación: el binario de producción no incluye la semilla de desarrollo. **Sin recomendación** («decide Axel/WF-4»).
  - Impacto: con A, WF-4 fija la variable en la unidad systemd y en `ops/runbook.md`.
- **Q4-80** · `C3-flujos.md` (observación) y `WF3a-bis-semilla-dev.md` punto 7 · Q2-49 = A dice «la app arranca sin precios de Walmart», pero la decisión PROVISIONAL del orquestador metió en la semilla real los 36 precios de Walmart del 10-oct-2026 (fuente `manual`). ¿Cuál vale?
  - Opciones (letras de este documento): **A** vale Q2-49 = A: la semilla de producción sin precios de Walmart; la ruta B de F-3 y «Mejor canasta» no tienen con qué comparar hasta que captures precios · **B** vale el punto 7: se quedan los 36 precios, que es la opción B de Q2-49 (se marcarían viejos el 9-nov-2026, K22). **Sin recomendación** («que Axel diga cuál de las dos vale»).
  - PROVISIONAL: B.
- **Q4-81** · Canasta del 10-oct-2026: **+$138.50 contra +$264.50** (`semilla-c3.md` l. 21 y 39; `backend/internal/api/comprasa/mejor_canasta_test.go` l. 77–105; `INTEGRACION-wf3c.md` §15.6.7).
  - HECHO: U1 §3.16, U1 §4 F-3 y el ejemplo de la tarjeta 5 de U3 dicen Walmart $3,233.00 (+$138.50), que viene de `data.json` (carrito contra carrito). Con K46 (Q2-47 = A, cada tienda se cotiza con los productos y cantidades de la compra base), la API y la pantalla dan Walmart $3,359.00 contra Bodega $3,094.50: **+$264.50**. La diferencia de $126.00 es de método: +$21.00 por la sexta lata de atún y +$105.00 porque para 780 g de Coffee Mate el motor compra dos de 640 g. En los artículos idénticos las dos dan +$100.50. La prueba `TestBestBasket_SemillaConPreciosDeWalmart` exige 309450 y 26450.
  - Opciones (letras de este documento): **A** +$264.50 (motor S5 y K46), y el dueño de U1 y U3 corrige las 3 líneas · **B** +$138.50 (el cálculo manual de `data.json`). **Recomendada: A**, de `INTEGRACION-wf3c.md` §15.6.10; `C3-flujos` punto 3 dice que K46 ya lo contesta. B contradice K46.

### 3.8 WF-3b · verificación visual y C3 (8 bloqueos y `C3-flujos`)

- **Q4-82** · `U5-ronda3-escalamiento.md` · U5 quedó devuelto en la ronda 3 (0 altos y 2 medios: M-7, iconos de la barra lateral; M-8, «capturado» en C3) y no hay ronda 4.
  - Opciones: **A** corrección acotada a M-7 y M-8 y verificación puntual · **B** aceptar U5 con M-7 y M-8 como pendientes con fecha · **C** otra ronda completa (contradice el límite de 3 rondas). **Recomendada: A**, del revisor `revisor-u5-r3` y del integrador (`INTEGRACION-wf3b.md` §11).
  - HECHO: M-8 se cerró desde C3 v1.0 y M-7 con los iconos de WF-3c (Q4-104); `verify-u5` da 0 pares bajo umbral. No hubo una verificación visual independiente después.
- **Q4-83** · `C3-ronda3-escalamiento.md` · C3 v1.2 quedó sin firma en la ronda 3 porque F-3 ruta A fallaba con la semilla.
  - Opciones: **A** eliges la (c) del punto 2 de `C3-flujos` (U1 aclara que la ruta A pide una lista con todos sus precios en al menos una tienda), verificación acotada y firma · **B** la (b): la semilla trae precios de los ingredientes · **C** firmar con F-3 ruta A como pendiente con fecha. **Recomendada: A**, del revisor `revisor-c3-r3`, que sugiere además la cara (b) de 2b como mejora posterior.
  - HECHO: así se hizo; C3 v1.4 está firmada con la (c) PROVISIONAL. Contestar A aquí equivale a la (c) de Q4-84.
- **Q4-84** · `C3-flujos.md` punto 2 · F-3 ruta A con la semilla.
  - Opciones: **(a)** registrarla como «falla con la semilla» · **(b)** que la semilla traiga precios de ingredientes (dueño: semilla o agente externo) · **(c)** U1 aclara que la ruta A solo aplica con una lista completa en precios. **Recomendada: (c)**, `INTEGRACION-wf3c.md` §15.6.10 («Confirmar la (c)»). La tarjeta que abre P-14 con la compra elegida (cara 2b) queda como PROPUESTA.
  - PROVISIONAL: (c); si falta un precio, la tarjeta «Mejor canasta» lleva a «Registrar precios».
- **Q4-85** · `C3-flujos.md` punto 4 · F-5 «Mensual · sep-2026» no existe con la semilla; la prueba usa oct-2026 con 6 toques (U1: 5).
  - Opciones: **(a)** «pasa con reserva» con oct-2026 · **(b)** la semilla genera el reporte de sep-2026 · **(c)** U1 cambia el ejemplo a «el último mes disponible». **Sin recomendación.** PROVISIONAL: (a).
- **Q4-86** · `C3-flujos.md` puntos 5, 6 y 7 · (5) en F-4, marcar la parte del hogar no cambia el disponible ni la deuda viva; (6) la «fecha de deuda en cero» de P-20: U1 la llama opcional y C1 y el servidor la exigen; (7) F-8 (ticket con OCR) pide al menos 8 toques y U1 dice 5.
  - Opciones: (5) **(a)** comportamiento previsto mientras la reposición esté «pendiente de confirmar con Kueski» · **(b)** confirmas o pides el recálculo. (6) **(a)** se sigue a C1 (obligatoria) y la etiqueta «(opcional)» se queda · **(b)** opcional: se corrigen C1 y el servidor · **(c)** obligatoria: se corrigen U1 N-225 y la etiqueta. (7) **(a)** «pasa con diferencia de toques» · **(b)** con tu orden, ajustar el frontend («Ticket» abre la cámara, se propone la categoría) y el lector S5 · **(c)** U1 actualiza el conteo. **Sin recomendación.** PROVISIONAL: (a) en los tres; en (6) la etiqueta de P-20 contradice lo que responde el servidor.
- **Q4-87** · `C3-capturas-global.md` · «Capturado» de la fila Global de C3 (G-01 a G-07), que no es una pantalla.
  - Opciones: **(a)** «no» · **(b)** aceptar P-01 y U5-13 como evidencia · **(c)** capturas propias de G-03, G-04 y G-06 · **(d)** sacarla del conteo. **Sin recomendación**; `INTEGRACION-wf3b.md` §9 la pone en la lista de PROVISIONALES a confirmar. PROVISIONAL: (a), por eso C3 dice 41 de 42.
- **Q4-88** · `U5-conteo-carpeta.md` · El comando literal de U5-01 da 168 (164 PNG más `axe/`, `contraste/`, `evidencia/` y `hosts.json`).
  - Opciones: **A** contar los PNG (164 de 164) · **B** mover los PNG a una subcarpeta · **C** reescribir el comando como `ls capturas-u5/*.png | wc -l`. **Sin recomendación.** PROVISIONAL: A.
- **Q4-89** · `U5-contraste-ui-controles.md` · ¿Qué cuenta como «componente de UI» en el 3:1 de U5-04?
  - Opciones: **A** contar los tres subtipos (campo, selección y botón) e informar cada uno · **B** solo campo y selección · **C** solo campo. **Sin recomendación.** PROVISIONAL: A; con A hoy da 0 bajo umbral.
- **Q4-90** · `U5-nombres-propios.md` · ¿Se aceptan en pantalla Claude, Ollama, Tesseract y las siglas API, CSV, PDF y OCR?
  - Opciones: **A** ampliar la lista con esos nombres y siglas; las claves internas no se muestran o van en español · **B** ampliar y aceptar también las claves internas · **C** no ampliar y cambiar los textos. **Sin recomendación.** HECHO: U5-07 da hoy 0 claves internas y 0 palabras en inglés.
- **Q4-91** · `U5-tipografia-tres-tamanos.md` · «Máximo tres tamaños por pantalla» (DESIGN §3.2) contra los 4 o 5 que piden sus propios componentes.
  - Opciones: **A** la regla cuenta solo el área de contenido y excluye las cifras de tarjeta · **B** «máximo cinco tamaños por pantalla» · **C** la regla literal y rediseñar tarjetas y títulos. **Sin recomendación.** PROVISIONAL: A.

### 3.9 WF-3c · Inicio como tablero de control (K54; 14 bloqueos)

Lo FIRME de K54 no se reabre: la fórmula «pagos mensuales de todos los créditos ÷ ingreso mensual neto», los umbrales 30 y 40 y que sean configurables. `INTEGRACION-wf3c.md` §15.6.10 recomienda «confirmarlos tal cual» para los valores de su §15.6.2; lo indico en cada pregunta.

- **Q4-92** · `U3-tablero-barometro.md` punto 1 · ¿Qué es el «ingreso mensual neto» del barómetro?
  - Opciones: **(a)** las quincenas del mes, capturadas más las esperadas, **sin extras** (un extra de un mes pintaría verde y escondería el riesgo; se ve aparte en K-2) · **(b)** quincenas más extras del mes · **(c)** promedio de las quincenas de los últimos 3 meses. **Recomendada: (a)** (§15.6.2 y §15.6.10).
- **Q4-93** · `U3-tablero-barometro.md` punto 2 · ¿Qué pagos cuentan? En Kueski, Karen paga directo su parte.
  - Opciones: **(a)** solo la parte del hogar de los pagos con fecha en el mes, pagados o no; el porcentaje con la parte de Karen solo se informa desde 768 px · **(b)** el monto total de cada pago · **(c)** solo lo que falta pagar del mes. **Recomendada: (a)** (§15.6.2 y §15.6.10).
- **Q4-94** · `U3-tablero-barometro.md` puntos 3 a 6 · (3) redondeo y fronteras; (4) escala del arco; (5) mes que se mide; (6) número de la sección de DESIGN.
  - Opciones: (3) **(a)** entero con ROUND_HALF_UP: ≤ 30 verde, 31 a 40 amarillo, ≥ 41 rojo · **(b)** un decimal y el estado con el valor exacto. (4) **(a)** de 0 al múltiplo de 10 igual o mayor que 1.5 × el tope (0 a 60 %) · **(b)** 0 a 100 % · **(c)** clave nueva. (5) **(a)** el mes en curso · **(b)** promedio de los próximos 3 meses. (6) **(a)** el barómetro entra como §8.9 sin renumerar · **(b)** renumerar. **Recomendada: (a) en 3 y 4** (§15.6.2 y §15.6.10); 5 y 6 sin recomendación (en uso (a)). Nota: la justificación del punto 4 trae una cuenta errónea (B-R3-1).
- **Q4-95** · `U3-tablero-horizonte.md` · ¿Cuántos meses ve cada gráfica? (1) G-1 ingreso, gasto y pagos; (2) G-6 tendencia del barómetro; (3) una o varias peticiones.
  - Opciones: (1) **(a)** los 11 meses anteriores y el en curso · **(b)** 6 anteriores, el en curso y 5 siguientes con el plan. (2) **(a)** 6 anteriores, el en curso y 5 siguientes · **(b)** 12 anteriores · **(c)** el en curso y 11 siguientes. (3) **(a)** una sola de 17 meses · **(b)** una por gráfica. **Recomendada: (a) en los tres** (§15.6.2 y §15.6.10).
- **Q4-96** · `U3-tablero-mas-del-hogar.md` · (1) ¿Qué pasa con las tarjetas 2 (quincena), 4 (comidas), 5 (mejor canasta) y 6 (ahorro) de U3 1.x? (2) ¿Y con la 1 (mes contra plan) y la 3 (próximo pago y disponible)?
  - Opciones: (1) **(a)** siguen sin cambio al final, bajo «Más del hogar» · **(b)** se quitan de Inicio (su información sigue en Detalle y en sus pantallas; cambian F-3 y sus pruebas) · **(c)** solo queda la 2. (2) **(a)** la 1 pasa a K-1 y la 3 se reparte entre K-3 y K-4 · **(b)** conservar la 3 además de K-4. **Sin recomendación** (en uso (a)).
- **Q4-97** · `U3-tablero-recordatorio.md` · «Pagos próximos» (K-4): (1) cuántos pagos se ven; (2) rango; (3) pagos o créditos.
  - Opciones: (1) **(a)** 5, en `dashboard.payment_reminder_count` (de 1 a 20) · **(b)** 3 · **(c)** todos los del rango. (2) **(a)** primero los vencidos y luego de hoy a hoy + 30 días · **(b)** solo el mes en curso · **(c)** de hoy a hoy + 60. (3) **(a)** una fila por pago · **(b)** una fila por línea con su próximo pago. **Recomendada: (a) en 1 y 2** (§15.6.2 y §15.6.10); 3 sin recomendación (en uso (a)).
- **Q4-98** · `U3-tablero-senal-presupuesto.md` · Regla de «Sugerencia: subir o recortar presupuesto» (G-3): (1) cuántos meses; (2) umbral; (3) meses sin plan o sin gasto; (4) qué dice y dónde.
  - Opciones: (1) **(a)** 3 meses cerrados (`budget.signal_months`, de 2 a 12) · **(b)** 2 · **(c)** 6. (2) **(a)** un margen de 10 % para los dos lados (`budget.signal_margin_pct`) · **(b)** márgenes distintos para subir y para recortar · **(c)** promedio contra el plan, sin racha. (3) **(a)** un mes sin plan corta la racha y un mes con plan sin gasto cuenta como 0 % · **(b)** se saltan. (4) **(a)** insignia «Sugerencia: …» sin monto · **(b)** además un monto sugerido · **(c)** además un aviso de S7 el día 1. **Recomendada: (a) en 1, 2 y 3** (§15.6.2 y §15.6.10); 4 sin recomendación (en uso (a)).
- **Q4-99** · `C1-tablero.md` · Lo que C1 v1.5 fija sin decisión: (1) sección de Ajustes de las 5 claves; (2) rangos; (3) necesidades de U1 de las 5 operaciones nuevas; (4) nombre de la clave del recordatorio; (5) formato de `period`.
  - Opciones: (1) **(a)** las cinco en Presupuesto (P-26.02) · **(b)** repartidas entre Créditos, Presupuesto y Notificaciones · **(c)** las cinco en otra sección (ver archivo). (2) **(a)** los rangos de U3 2.0 §2.6, validados por disparadores · **(b)** solo enteros de 0 o más. (3) **(a)** colgarlas de necesidades vigentes · **(b)** U1 agrega necesidades propias. (4) **(a)** `dashboard.payment_reminder_count` · **(b)** el plural. (5) **(a)** solo `AAAA-MM`; una quincena da 422 · **(b)** aceptar quincena y tomar su mes. **Recomendada: (a) en 1** (es la A de `C3-ajustes-tablero`, §15.6.10) **y en 2** (§15.6.2); 3, 4 y 5 sin recomendación (en uso (a)).
- **Q4-100** · `S6-tablero-sin-datos.md` · ¿Cuándo un mes «no tiene datos» (se muestra «—» y no $0.00)?
  - Opciones: **(a)** cuando el devengado y el efectivo del hogar son 0 · **(b)** los meses anteriores al primer gasto o pago registrado · **(c)** nunca se omite. **Sin recomendación** (en uso (a); también decide el mes anterior de K-2).
- **Q4-101** · `S6-tablero-senal-periodo.md` · «Si pides la sugerencia de un mes que todavía no llega, ¿prefieres que no sugiera nada (a), que el sistema diga que ese mes no se puede pedir (b) o que use los últimos meses ya cerrados (c)?»
  - Opciones: las tres de la pregunta. **Sin recomendación** (en uso (a)). La PWA pide hoy el mes en curso, donde las tres dan lo mismo.
- **Q4-102** · `S6-tablero-lineas-archivadas.md` · (1) «Cuando archives una tarjeta que ya terminaste de pagar, ¿quieres que sus pagos de meses pasados sigan apareciendo en el barómetro y en las gráficas del tablero (a) o que desaparezcan (b)?» (5) «Si archivas una línea que todavía tiene pagos por hacer, ¿quieres que su deuda y sus pagos sigan sumando en «Deuda viva», en «Créditos del mes» y en la gráfica de deuda proyectada (a), o que salgan aunque falte pagarlos (b)?»
  - Opciones: las de cada pregunta. **Recomendada: (a) en las dos**, `INTEGRACION-wf3c.md` §15.6.10 («Confirmar la (a) y el punto 5»).
  - PROVISIONAL: (a) en el código. HECHO: desde WF-3c ter el comparador de escenarios ya no ofrece líneas archivadas para compras nuevas.
- **Q4-103** · `C3-ajustes-tablero.md` · `updateDashboardSettings` no tenía pantalla. ¿Qué se hace?
  - Opciones: **A** apartado «Parámetros del tablero» en P-26.02 con las 5 claves y «Restablecer» · **B** excepción técnica · **C** quitar la operación de C1 o dejarla para el agente externo. **Recomendada: A**, PROPUESTA de `c3-tablero` y `INTEGRACION-wf3c.md` §15.6.10.
  - HECHO: A ya está aplicada en P-26.02 con su vista previa y su validación; el texto del archivo (l. 21) todavía dice «No se tomó ninguna decisión» (B-2 de frontend).
- **Q4-104** · `WF3c-barra-iconos-seccion.md` · Iconos de sección de la barra lateral que DESIGN §5.4 no nombra: Presupuesto `wallet`, Salud `heart-pulse`, Reportes y Gastos (los cuatro, en el archivo).
  - Opciones: **A** confirmar los 4 y que DESIGN §5.4 los agregue · **B** Gastos con `arrow-up-right` y el resto igual · **C** otros. **Sin recomendación.** Contesta también la pregunta de iconos de Q4-25.
- **Q4-105** · `WF3c-tablero-interfaz.md` · Seis puntos de interfaz que U3 2.0 no cierra: (1) K-4 dice «Hay más pagos en Calendario» sin número; (2) K-3 sin «Reintentar» dentro de la tarjeta-enlace; (3) la insignia de K-1 es solo texto; (4) alto del barómetro; (5) separación de columnas de G-1 (8 %); (6) G-2 y G-3 como barras HTML.
  - Opciones: **A** así (lo hecho) y **B** su alternativa, por punto, en el archivo (por ejemplo (1) B: C1 agrega un conteo). **Sin recomendación.**
- **Q4-106** · `WF3c-barometro-sin-dato.md` · Rol accesible del barómetro sin dato.
  - Opciones: **A** `role="img"` con su texto (hecho) · **B** texto solo para lector de pantalla. **Sin recomendación** (en uso A). DESIGN §8.9 debe fijarlo.
- **Q4-107** · `WF3c-ronda3-escalamiento.md` · ¿Ratificas lo que WF-3c bis aplicó de este escalamiento (puntos 1, 2 (i), 3a a 3d, 4 y 5): `scale_max_pct` en C1 v1.5.1 (G-6 y el barómetro con la misma escala); los pagos de todas las líneas; los flujos y C3-09 al día con U3 2.0; el `<svg>` del barómetro con `fill="none"`; y P-01 Detalle a 412 px en Android?
  - Opciones: por punto en el archivo; 3b y 3c se contestan en Q4-103 y Q4-108. **Recomendada:** ratificar, porque cada pieza coincide con lo que recomienda §15.6.10 (Q4-102, Q4-103 y Q4-108); el archivo sigue abierto solo por esa ratificación (`INDICE.md`, WF-3c bis).

### 3.10 WF-3c bis · cierre del tablero (3 bloqueos)

- **Q4-108** · `S6-deuda-viva-hoy.md` · (1) La deuda viva a una fecha, ¿cuenta solo las compras con fecha hasta el corte o también las programadas a futuro? (2) La deuda viva proyectada (G-5), ¿aplica la misma regla en cada cierre? (3) ¿Qué es «falta pagar» del mes en K-3?
  - Opciones: (1) **(a)** solo las compras hasta el corte: hoy $3,895.00 en K-3, P-17, P-07 y `getLiveDebt`, igual que el oráculo · **(b)** todas las activas: K-3 y P-17 dirían $18,999.22. (2) **(a)** sí: G-5 arranca en $2,921.25 al cierre de oct-2026 · **(b)** todas las activas: arranca en $18,025.47. (3) **(a)** parte del hogar menos lo que ya abonaron Axel o Ximena, nunca negativa ($598.82 en oct-2026) · **(b)** la parte del hogar completa de los pagos pendientes o parciales. **Recomendada: (a) en los tres**, `INTEGRACION-wf3c.md` §15.6.10.
  - PROVISIONAL: (a). Pendiente técnico: `deriveCredits` de las vistas de base de datos todavía usa la regla anterior (§1.5).
- **Q4-109** · `C1-unidad-parametros-tablero.md` · El servidor manda como «unidad» de las 5 claves de K54 una descripción de 57 a 123 caracteres, y C1 promete una unidad corta.
  - Opciones: **A** la API da la unidad corta con una migración nueva de anotaciones, sin editar la 00005 · **B** P-26.02 pone sus unidades en `textos.ts` («%», «%», «pagos», «meses», «puntos») · **C** `Parametro` muestra la unidad solo si es corta. **Recomendada: A**, del archivo y de `INTEGRACION-wf3c.md` §15.6.10 («se retira el puente (B)»).
  - PROVISIONAL: B.
- **Q4-110** · `U3-tablero-k4-personas.md` · «¿Qué muestra K-4 si falla `listPersons`?»
  - Opciones: **A** el estado de error del bloque con «Reintentar»; nunca una cifra distinta del oráculo · **B** las filas sin la línea «Hogar · Karen» y sin las variantes «Falta la parte de…». **Recomendada: A**, PROPUESTA aplicada del archivo y `INTEGRACION-wf3c.md` §15.6.10.
  - HECHO: en WF-3c ter, D-3 de P-01 Detalle hace lo mismo; falta la insignia de `Calendario.tsx`.

### 3.11 Registrado como resuelto en su propio texto

- `criterios-nxx.md`: «RESUELTO el 10-oct-2026 con la opción (b)» (`INDICE.md`, sección «Resuelto»). No pide respuesta; aparece en el conteo del comando porque no trae la frase «Respuesta de Axel».
- `WF3-ajustes-registro.md`: «RESUELTO entre áreas en WF-3a»; su PROPUESTA pendiente va en Q4-64.

## 4. Acciones de Axel (sin id)

1. **Instalar la CA de mkcert en los dos Android** (el tuyo y el de Ximena) y abrir `https://finanzas.local:8443` en Chrome, sin advertencia, y luego «Instalar app» (`ops/mkcert.md` §4, pasos 1 a 5). El certificado se baja en **Ajustes › Seguridad › Descargar certificado de la casa** o por cable USB. Si `finanzas.local` no abre en el teléfono, usa `https://192.168.50.56:8443/`, que también vale para el certificado (Q4-47). Con esto se cierra S6-11.
2. **Confirmar que llegaron los 2 mensajes de Telegram de las 21:15** (el resumen diario y la alerta de pago) y que su contenido te sirve (Q4-35).
3. **Entregar `datos/semilla/gastos_fijos_ingresos.csv`** con las columnas de `plan.md` §8.2 (l. 408): `tipo`, `concepto`, `monto`, `periodicidad`, `dia_o_fecha`, `categoria`, `persona`… (cómo se lee cada una: Q4-55). Hoy no existe; con Q4-03 = A, WF-4 no lo espera.
4. **5 fotos de tickets reales** en `datos/semilla/tickets/` (K24); la carpeta está vacía y el OCR sigue sin medir.
5. **Que Ximena le escriba al bot**, cuando esté decidido Q4-43 (`S7-convivencia-multix`); después confirmas su chat en Ajustes (K4).
6. **Cambiar las contraseñas iniciales en Ajustes después de WF-4.** Hoy las de la base de desarrollo viven fuera del repositorio (`ops/runbook.md` l. 52); cómo se dan y se cambian, en `docs/estado/wf3/semilla.md` («Contraseñas iniciales»).

Otras entregas que salen de las preguntas: los 2 PDF de estado de cuenta (Q4-11); la alcaldía o municipio de las dos sucursales (Q4-27); los empates de nutrición (Q4-13); instalar `texlive-luatex` si eliges Q4-45 = A; y, en el módem, la reserva de IP y la entrada DNS si eliges Q4-47 = 1A y 2C.

## 5. Hoja de respuestas

| Q4 | Tema | Recomendada | Respuesta de Axel |
|---|---|---|---|
| Q4-01 | Aprobar el paso a operación tal cual (8ace01e) | sin recomendación | |
| Q4-02 | Arrancar WF-4 al contestar | sin recomendación | |
| Q4-03 | WF-4 carga `gastos_fijos_ingresos.csv` si existe | A (plan l. 186) | |
| Q4-04 | S2-ahorro-quincenal · reparto, metas y saldo | sin recomendación | |
| Q4-05 | S2-metas-vencida · aportación de una meta vencida | sin recomendación (en uso A) | |
| Q4-06 | S4-presentacion-respaldo · presentación sin precio (Q2-07) | sin recomendación (en uso: regla del motor) | |
| Q4-07 | S5-carritos-granel · redondeo a granel | sin recomendación (en uso A) | |
| Q4-08 | C1-hora-reportes-job · qué manda en la hora del reporte | A | |
| Q4-09 | C1-reparto-alcance · alcance de «todos» | B | |
| Q4-10 | S1-reparto-edicion-gasto · editar un gasto viejo | sin recomendación (en uso C) | |
| Q4-11 | S3-formula-tarjeta-pdf · entregar 2 PDF | sin opciones (entrega) | |
| Q4-12 | S3-iva-lectura · préstamo y crédito sin dato de IVA | sin recomendación (en uso 1A) | |
| Q4-13 | S4-empates-sin-confirmar · empates de nutrición (Q2-90) | sin recomendación | |
| Q4-14 | S5-motor-oraculo-reglas · 8 reglas motor contra oráculo | sin recomendación (en uso A) | |
| Q4-15 | S5-precio-vigente-captura · 4 reglas del precio vigente | sin recomendación (en uso: el motor) | |
| Q4-16 | U4-reparto-preseleccion · opción preseleccionada | A | |
| Q4-17 | S5-agente-externo (a) · peso de sus precios | A (compiten como manual) | |
| Q4-18 | S5-agente-externo (b) · qué escribe | A | |
| Q4-19 | S5-agente-externo (c) · recetas como borrador | A | |
| Q4-20 | S5-agente-externo (d) · valor de `price.source` | d2 | |
| Q4-21 | S5-agente-externo (e) · credencial y a nombre de quién | e1 + e4 | |
| Q4-22 | S5-agente-externo (f) · borrador y origen de receta | sin recomendación (en uso f2) | |
| Q4-23 | S5-agente-externo (g) · tope por lote | sin recomendación (en uso g1) | |
| Q4-24 | P3-contraste-chip-hover | A | |
| Q4-25 | P3-huecos-componentes · supuestos de DESIGN | sin recomendación | |
| Q4-26 | S6-semilla-despensas-futuras | sin recomendación (en uso A) | |
| Q4-27 | S6-semilla-sucursales-municipio | sin recomendación | |
| Q4-28 | S6-semilla-catalogo-alias | sin recomendación | |
| Q4-29 | K31 · respaldo: 03:15 y 30 diarios | (a) | |
| Q4-30 | K31 · qué conserva una restauración | sin recomendación (en uso a) | |
| Q4-31 | K31 · tipos de adjunto | sin recomendación (en uso a) | |
| Q4-32 | Adjuntos · MB, `kind` y 409 | sin recomendación (en uso a) | |
| Q4-33 | K31 · cabeceras de seguridad y CSP | (a) | |
| Q4-34 | K31 · argon2id, límites, purga, rutas, campo desconocido | (a) | |
| Q4-35 | K31 · contenido de los resúmenes de Telegram | sin opciones escritas | |
| Q4-36 | S7 · día de despensa | C | |
| Q4-37 | S7 · semana que revisa la alerta 4 | A (la siguiente) | |
| Q4-38 | K31 · cadencia y reintentos de avisos | (1) PROPUESTA; (2) A; (3) sin recomendación | |
| Q4-39 | K31 · reporte mensual: secciones y negativos | (1) A; resto sin recomendación | |
| Q4-40 | K31 · PWA sin conexión (15 s, 30 s, copia, cola) | los 4 valores propuestos | |
| Q4-41 | K31 · cómo se sirve la PWA | A en producción, B en desarrollo | |
| Q4-42 | S6-sistema · 8 puntos de autenticación y sistema | sin recomendación | |
| Q4-43 | S7-convivencia-multix · chat de Ximena | A | |
| Q4-44 | S7-ruta-token | A | |
| Q4-45 | S7-lualatex | A | |
| Q4-46 | S7-estado-cuenta-fechas | A | |
| Q4-47 | WF3-integrador-red · IP, nombre y CA | 1A, 2C (+A), 3A | |
| Q4-48 | S6-ajustes-api · 6 puntos | sin recomendación | |
| Q4-49 | S6-comidas-calculos | sin recomendación | |
| Q4-50 | S6-compras-a-motores | sin recomendación | |
| Q4-51 | S6-compras-b-fuentes | sin recomendación | |
| Q4-52 | S6-compras-b-nutricion | sin recomendación | |
| Q4-53 | S6-creditos | sin recomendación | |
| Q4-54 | S6-gastos-inicio | sin recomendación | |
| Q4-55 | S6-importar-partidas-csv · columnas del CSV | sin opciones escritas | |
| Q4-56 | S6-presupuesto-motor | sin recomendación | |
| Q4-57 | S6-presupuesto-periodos | sin recomendación (en uso a) | |
| Q4-58 | S6-presupuesto-quincenas | sin recomendación (en uso a) | |
| Q4-59 | C3-matriz · ratificar 1, 3c y 4a | ratificar | |
| Q4-60 | WF3-ajustes-a · despensa programada, membresía, P-26 | sin recomendación | |
| Q4-61 | WF3-ajustes-compra-ejemplo | sin recomendación (en uso A) | |
| Q4-62 | WF3-ajustes-muestra-tema | A | |
| Q4-63 | WF3-ajustes-precios-conteo | sin recomendación (en uso A) | |
| Q4-64 | WF3-ajustes-registro | A | |
| Q4-65 | WF3-comidas-compras-interfaz · 7 puntos | sin recomendación | |
| Q4-66 | WF3-creditos-salud-reportes · 7 puntos | sin recomendación | |
| Q4-67 | WF3-frontend-rutas · direcciones en español | A | |
| Q4-68 | WF3-gastos-contrato · 4 huecos de C1 | sin recomendación | |
| Q4-69 | WF3-gastos-ticket-credito | sin recomendación | |
| Q4-70 | WF3-inicio-tiendas · color de tienda | sin recomendación (en uso A) | |
| Q4-71 | WF3-inicio-vacios · 6 reglas de vacío | sin recomendación (en uso A) | |
| Q4-72 | WF3a-bis-ajustes-a · edición en grupo | sin recomendación | |
| Q4-73 | WF3a-bis-comidas-compras · 4 puntos | sin recomendación | |
| Q4-74 | WF3a-bis-gastos-presupuesto · 7 puntos | sin recomendación | |
| Q4-75 | Semilla de desarrollo · renta (punto 1) | sin recomendación (en uso a) | |
| Q4-76 | Semilla de desarrollo · presupuesto de otras categorías (punto 3) | sin opciones escritas | |
| Q4-77 | Semilla de desarrollo · costo y kcal en 0 (punto 5) | sin recomendación | |
| Q4-78 | Semilla real · gastos proyectados 2 a 6 (punto 4) | (a) | |
| Q4-79 | Candado de la semilla de desarrollo (punto 2) | sin recomendación (en uso A) | |
| Q4-80 | Q2-49 contra 36 precios de Walmart en la semilla | sin recomendación (en uso B) | |
| Q4-81 | Canasta del 10-oct: +$264.50 o +$138.50 | A (+$264.50, K46) | |
| Q4-82 | U5-ronda3-escalamiento | A | |
| Q4-83 | C3-ronda3-escalamiento | A | |
| Q4-84 | C3-flujos punto 2 · F-3 ruta A | (c) | |
| Q4-85 | C3-flujos punto 4 · F-5 sep-2026 | sin recomendación (en uso a) | |
| Q4-86 | C3-flujos puntos 5, 6 y 7 | sin recomendación (en uso a) | |
| Q4-87 | C3-capturas-global | sin recomendación (en uso a) | |
| Q4-88 | U5-conteo-carpeta | sin recomendación (en uso A) | |
| Q4-89 | U5-contraste-ui-controles | sin recomendación (en uso A) | |
| Q4-90 | U5-nombres-propios | sin recomendación | |
| Q4-91 | U5-tipografia-tres-tamanos | sin recomendación (en uso A) | |
| Q4-92 | Barómetro · ingreso mensual neto | (a) | |
| Q4-93 | Barómetro · pagos que cuentan | (a) | |
| Q4-94 | Barómetro · redondeo, escala, mes y §8.9 | (a) en 3 y 4; 5 y 6 sin recomendación | |
| Q4-95 | Horizonte de las gráficas | (a) | |
| Q4-96 | «Más del hogar» · tarjetas de U3 1.x | sin recomendación (en uso a) | |
| Q4-97 | Pagos próximos · N, rango y forma | (a) en 1 y 2 | |
| Q4-98 | Sugerencia de presupuesto · regla | (a) en 1, 2 y 3 | |
| Q4-99 | C1-tablero · 5 puntos | (a) en 1 y 2 | |
| Q4-100 | S6-tablero-sin-datos | sin recomendación (en uso a) | |
| Q4-101 | S6-tablero-senal-periodo | sin recomendación (en uso a) | |
| Q4-102 | Líneas archivadas en el tablero (y punto 5) | (a) y (a) | |
| Q4-103 | C3-ajustes-tablero · Parámetros del tablero | A | |
| Q4-104 | WF3c-barra-iconos-seccion · 4 iconos | sin recomendación | |
| Q4-105 | WF3c-tablero-interfaz · 6 puntos | sin recomendación (en uso A) | |
| Q4-106 | WF3c-barometro-sin-dato | sin recomendación (en uso A) | |
| Q4-107 | WF3c-ronda3-escalamiento · ratificar | ratificar | |
| Q4-108 | S6-deuda-viva-hoy · deuda de hoy, G-5 y falta pagar | (a) en los tres | |
| Q4-109 | C1-unidad-parametros-tablero | A | |
| Q4-110 | U3-tablero-k4-personas | A | |

Preguntas de la puerta 2 que siguen abiertas (§3.1; la variante de Q2-07 es Q4-06):

| Q2 | Tema | Recomendada | Respuesta de Axel |
|---|---|---|---|
| Q2-01 | S2-ventanas · casos 1 y 2 de la pregunta 4 | sin recomendación en el caso 2 | |
| Q2-12 | S2-flujo-quincenal (1) gasto variable estimado | «C o A» en la puerta 2 | |
| Q2-48 | S5-carrito-comparativa (2) fuente de Bodega | sin recomendación (en uso A) | |
| Q2-50 | S5-carritos-detalle (1) ganador del sustituto | sin recomendación (en uso A) | |
| Q2-51 | S5-carritos-detalle (2) `missing_count` | sin recomendación (en uso A) | |
| Q2-52 | S5-carritos-detalle (3) `stale_count` | sin recomendación (en uso A) | |
| Q2-56 | S5-profeco-filtro (1) municipios | sin recomendación (en uso A) | |
| Q2-57 | S5-profeco-filtro (2) comparación de nombres | sin recomendación (en uso A) | |
| Q2-58 | S5-profeco-filtro (3) sucursal del precio | sin recomendación (en uso A) | |
| Q2-85 | U3 · supuesto 5 (ya retirado) | sin recomendación | |
| Q2-87 | Recetas · etiquetas | sin recomendación | |
| Q2-88 | Recetas · kcal de cortes con hueso | sin recomendación | |
| Q2-89 | Recetas · recetas parecidas | sin recomendación | |
| Q2-90 | Recetas · 40 empates de `revisar.tsv` (= Q4-13) | sin recomendación | |
| Q2-91 | PROFECO · cadenas del catálogo | sin recomendación (A o B) | |
| Q2-94 | Lista · subtotal «falta por comprar» | sin recomendación | |
