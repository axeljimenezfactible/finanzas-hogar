# Índice de criterios de aceptación · Finanzas Hogar (Etapa B)

Versión 1.1 (10-oct-2026; filas y totales actualizados el 11-oct-2026, puerta 1; filas C1, U1 y U5 y total corregidos en la ronda 1 de P1; conteo reconciliado por el integrador de P1 el 11-oct-2026, §10) · Escrito en WF-0 por el agente «criterios»; corregido en la ronda 1 de revisión (§9) · Plan: `docs/plan/plan.md` §3.1 principio 4 («Criterios antes del trabajo… los revisores evalúan solo esos criterios»). Cada criterio se traza al plan, a `docs/plan/decisiones.md` o al protocolo (la parte que no deriva de una frase de esas fuentes se marca «convención propia del directorio» en su sección 5); los huecos son bloqueos en `docs/estado/bloqueos/criterios-*.md`.

## 1. Cómo usar estos archivos
- **Un archivo por ítem** (18 + este índice). Cada uno trae: 1 Objetivo · 2 Criterios numerados `<ID>-nn` (enunciado medible, «Verifica» con comando o cálculo, «Evidencia», severidad) · 3 No aceptable · 4 Umbrales numéricos · 5 Trazabilidad.
- **Severidad:** *alta* bloquea la puerta; *media* el revisor la devuelve y debe quedar corregida o justificada por escrito antes de la ronda final; *baja* se registra en el informe del revisor. (Hoy los 257 criterios de la tabla son alta o media: 197 alta y 60 media; ninguno baja.)
- **Evidencia:** cada criterio deja una línea que **empieza con su ID** en `docs/estado/<flujo>/<agente>.md` (flujos `wf0`, `wf1`, `wf3`, `wf4`), con el comando ejecutado y su salida resumida; los archivos con nombre fijo del plan son `docs/estado/S3-escenarios.md`, `S5-terminos.md`, `S5-cobertura-profeco.md`, `U5-revision.md`, `P3-iter-N.md` e `informe-final.md`. El revisor **reejecuta** los comandos; no se aprueba por la sola palabra del agente.
- **Lo que depende de Axel** (instalar la CA en los Android, abrir la PWA en Android, reiniciar la máquina, que Ximena escriba al bot, iniciar sesión en claude.ai/design) **no se declara cumplido** hasta que Axel lo confirme: queda «pendiente de Axel» o «no verificado» en el informe final.
- **Rutas y comandos:** absolutos, desde `/home/deepsec/finanzas-hogar`. `jq` no está instalado en la máquina (se usa `python3 -I`). Los scripts auxiliares los entrega cada agente y se nombran en su reporte; si falta una herramienta (`goose`, `govulncheck`, `mkcert`, validador OpenAPI, `style-dictionary`, `jsonschema`), se instala (autorizado) y se anota en «Paquetes instalados».
- **Identificadores de U1 y C1:** pantalla `P-xx` (inventario de U1 §2.1; `P-26.xx` para las 14 secciones de Ajustes, §2.3); necesidad de datos `N-xxx` (U1 §3; tipo L lectura o E escritura); componente global `G-xx`; flujo clave `F-x` (U1 §0.1). En los criterios, «pantalla» es siempre `P-xx` y «necesidad» es siempre `N-xxx`; T es el número de pantallas `P-xx` más secciones `P-26.xx` que declara U1-02.
- **Cifras:** toda cifra de `oraculo/model.json` (sha256 `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`, 596 hojas numéricas) se reproduce exactamente; si difiere es hallazgo (protocolo §6).

## 2. Ítems, flujo, puerta y conteo
| Archivo | Ítem | Se construye en | Lo evalúa | Criterios (alta / media) |
|---|---|---|---|---|
| `S1.md` | Dominio y datos | WF-0 | Revisor WF-0 → puerta 1 | 17 (12 / 5) |
| `C1.md` | Contrato de datos y API (OpenAPI) | WF-0 | Revisor WF-0 → puerta 1 | 21 (19 / 2) |
| `U1.md` | Arquitectura de información | WF-0 | Revisor WF-0 → puerta 1 | 16 (10 / 6) |
| `S2.md` | Motor financiero (85 cifras canon) | WF-1 | Revisor motores → puerta 2 | 16 (13 / 3) |
| `S3.md` | Motor de créditos y optimizador (85 cifras canon) | WF-1 | Revisor motores → puerta 2 | 20 (17 / 3) |
| `S4.md` | Recetas y comidas | WF-1 | Revisor recetas+precios → puerta 2 | 15 (11 / 4) |
| `S5.md` | Precios y tiendas | WF-1 | Revisor recetas+precios → puerta 2 | 13 (10 / 3) |
| `U2.md` | Sistema de diseño y DESIGN.md | WF-0 (colorimetría) y WF-1 | Revisor diseño → puerta 2 | 14 (10 / 4) |
| `C2.md` | Contrato de tokens | WF-1 | Revisor diseño → puerta 2 | 12 (8 / 4) |
| `U3.md` | Dashboard | WF-1 | Revisor diseño → puerta 2 (y U5) | 12 (8 / 4) |
| `U4.md` | Ajustes estilo Notion | WF-1 | Revisor diseño → puerta 2 (y U5) | 11 (8 / 3) |
| `P3-claude-design.md` | Pista 3 · Claude Design | Entre puertas 3 y 4 (orquestador) | Revisor visual Sonnet; puerta 4 | 13 (10 / 3) |
| `S6.md` | Backend local | WF-3 | Revisores backend ×2 → puerta 4 | 16 (12 / 4) |
| `S7.md` | Reportes | WF-3 | Revisor backend → puerta 4 | 14 (12 / 2) |
| `U5.md` | Prototipo y verificación visual | WF-3 | Revisor visual → puerta 4 | 13 (10 / 3) |
| `C3.md` | Contrato de integración | WF-3 | Revisor independiente → puerta 4 | 11 (9 / 2) |
| `S8.md` | Calidad y operación | WF-4 | Revisor final | 11 (8 / 3) |
| `WF4-operacion.md` | Operación, semilla, seguridad, informe | WF-4 | Revisor final | 12 (10 / 2) |
| **Total** | | | | **257 (197 / 60)** |

Notas: S4 y S5 también dependen de WF-1 y S6/S7 de WF-3 (plan §3.3). F1–F5 del frontend (plan §5) no tienen archivo propio: se evalúan por U3, U4, U5 y C3.

## 3. Puertas de aprobación y criterios que las gobiernan
| Puerta (decisiones J, plan §3.8) | Qué recibe Axel | Criterios que deben estar en verde |
|---|---|---|
| 1 · tras WF-0 | C1, U1, criterios de todas las fases | `C1.md`, `U1.md`, `S1.md`; este directorio completo |
| 2 · tras WF-1 | DESIGN.md, C2, regresión de motores, bloqueos acumulados | `S2.md`, `S3.md`, `S4.md`, `S5.md`, `U2.md`, `C2.md`, `U3.md`, `U4.md` |
| 3 · antes de Pista 3 | paquete de entrada | P3-01 y P3-02 de `P3-claude-design.md` |
| 4 · tras WF-3 | C3 firmado, capturas U5, informe de pruebas | `S6.md`, `S7.md`, `U5.md`, `C3.md`, `P3-claude-design.md` |
| Cierre · WF-4 | informe final | `S8.md`, `WF4-operacion.md` |

## 4. Cifras canon (regresión obligatoria)
- `S2.md` §6 y `S3.md` §6 llevan **la misma tabla de 85 casos R-01…R-85** (cada fila con su ruta en `oraculo/model.json` y el motor que la calcula: 30 de S2, 55 de S3); entre ambos cubren las 596 hojas numéricas (159 de S2 + 437 de S3; S2-02 y S3-02).
- `S5.md` §6 lleva la comparativa de tiendas K-01…K-08 (de `datos/semilla/data.json`).
- Las tablas se generaron con un script que **afirma** cada valor contra `model.json` y `data.json` (no se transcribió a mano).

## 5. Hallazgos de este agente (no son bloqueos; cada uno tiene su criterio)
1. **Invariante del plan no se cumple literalmente** (S2-07): «devengado acumulado − efectivo acumulado = deuda viva» solo se cumple al cierre de P2, P4 y P6; al 10-abr-2027 la deuda viva es 6,293.44 y la diferencia 1,258.70. Se cumple siempre la identidad con el saldo no devengado. Bloqueo `criterios-invariante-deuda`.
2. **`oraculo/verify3.py` línea 89 escribe fuera del repositorio** (`/home/deepsec/reporte-compra-latex/model/verify3.json`), contra protocolo §7. S2-04 y S3-16 lo exigen corregido en WF-1. Además `oraculo/model.py` **regenera `model.json` al ejecutarse** (en la lectura de este agente corrió una vez; el archivo quedó **idéntico byte a byte** al original de `~/reporte-compra-latex/model/`, mismo sha256).
3. **17.27 % a 2 quincenas es un redondeo**: 534.50 / 3,094.50 = 17.2725 %; 17.27 % × 3,094.50 daría 534.42 y `data.json` guarda 534.50. Las pruebas comparan contra el valor guardado (nota de la tabla en S2/S3).
4. **Dos cifras parecidas pero distintas:** valor estable por quincena de Axel **883.76** (R-41) y pago por quincena de la despensa 2 **883.75** (último 883.77).
5. **«33 idénticos»** son 33 renglones (55 piezas); los totales completos incluyen 3 sustitutos y el sexto atún de Bodega (S5, nota de §6).
6. **Contraste de bordes de campo:** con los hex de partida de C, neutro 600 sobre neutro 50 da **3.71:1**, no «≥ 4.2:1» como dice plan.md §4.3; cumple el umbral de UI (3:1). Documentarlo (U2-05).
7. **Roles como texto:** los pasos 600 de ingreso (4.21), ahorro (4.26) e información (4.41) no llegan a 4.5:1 sobre fondo 50; para texto debe usarse un paso más oscuro (U2-06).
8. **Colisión no listada en decisiones I:** error (20°) y gasto (25°) quedan a 5° (y ingreso/éxito a 0°, ya aceptado). Se resuelven con refuerzo de icono o forma (U2-08).
9. **ΔL ≥ 0.10 entre series** no es alcanzable entre los 28 pares de 8 series con contraste ≥ 3:1 sobre fondo claro; se usa la lectura de series consecutivas (U2-09, bloqueo `criterios-deltal`).
10. **Token de Telegram vs usuario del servicio:** el archivo está en el home de `deepsec` y el servicio correría como `finanzas` (bloqueo `criterios-token-servicio`).
11. **Go 1.26.5 instalado frente a Go 1.27 aprobado:** el plan ya pide actualizar (S6-01).
12. **Ventanas del oráculo (10 al 9) frente al mes calendario:** las decisiones B fijan el mes calendario como periodo base y P1–P6 de `model.json` van del día 10 al 9. Las 12 cifras de devengado y efectivo coinciden con oct-2026…mar-2027 solo porque la semilla no tiene eventos entre los días 1 y 9 de un mes. S1-09 pide la prueba de equivalencia con el oráculo y S2-09 la nota de que depende de la semilla (hallazgo H-07 de la ronda 1).

## 6. Bloqueos escritos (`docs/estado/bloqueos/`)
| Archivo | Pregunta para Axel | Criterios afectados |
|---|---|---|
| `criterios-invariante-deuda.md` | qué es «deuda viva» en el invariante | S2-07 |
| `criterios-promedios.md` | fórmula de promedio diario y semanal | S2-08 |
| `criterios-proyeccion-12m.md` | supuestos de los meses 7–12 | S2-09 |
| `criterios-credito-formulas.md` | mínimo e intereses de tarjeta; crédito de tienda | S3-08, S3-10 |
| `criterios-optimizador.md` | normalización del puntaje y medidas de objetivos 2 y 4 | S3-13, S3-14 |
| `criterios-frescura-peso.md` | factor de peso de precios viejos y frontera de 30 días | S5-04 |
| `criterios-ocr-tickets.md` | origen de los 5 tickets y qué es renglón correcto | S5-06 |
| `criterios-meta-kcal.md` | cómo se obtiene la meta de calorías | S4-07 |
| `criterios-login-bloqueo.md` | duración del bloqueo y de la sesión larga | S6-06, S6-07 |
| `criterios-bind.md` | ¿Tailscale además de la LAN? | S6-14 |
| `criterios-token-servicio.md` | cómo llega el token al servicio con usuario propio | S7-06, S7-07, WF4-01 |
| `criterios-toques.md` | qué cuenta como «toque» | U1-04, U5-09, C3-08 |
| `criterios-nxx.md` | **Resuelto** en la ronda 1 de correcciones con la opción (b): «N-xx» = necesidad de datos `N-xxx`; la pantalla es `P-xx` (no queda pregunta para Axel) | U1-02, U1-03, C1-02, C3-02, U5-01 |
| `criterios-deltal.md` | ΔL ≥ 0.10 entre series consecutivas o todos los pares | U2-09 |
| `criterios-cobertura.md` | si el 85 % aplica a recetas y precios | S4-15, S5-13, S6-02 |
| `criterios-umbrales-ausentes.md` | 8 valores no fijados (horas, retención, cabeceras, adjuntos…) | S3-12, S6-09, S6-10, S7-04, S7-08, S7-10, S8-02, WF4-04, WF4-08 |

## 7. Trazabilidad con `docs/plan/decisiones.md`
| Sección de decisiones | Frases cubiertas | Archivos |
|---|---|---|
| A · Personas y acceso | Axel y Ximena co-titulares; login, sesión larga, argon2id, 5 en 15 min; Karen sin cuenta, paga directo, estado de cuenta PDF; plan de comidas solo 2 adultos; avisos a Axel y Ximena | S1, S4, S6, S7, U1, U3, U4, C1 |
| B · Dinero, periodos y reparto | MXN, centavos, ROUND_HALF_UP, último pago ajusta; ISO/pantalla; mes calendario y quincena; despensa devengada mitad por mes; ingresos solo de Axel; reparto 50/50; ahorro y metas; comida fuera | S1, S2, S3, S4, S7, C1, U3, U4, U5 |
| C · Créditos | 4 tipos; tasa por compra y tabla por plazo (17.27 / 25.8685); promociones; interruptor A/B por defecto B; optimizador 4 objetivos y 3 restricciones; regresión Kueski | S1, S3, U3, U4, C1 |
| D · Comida, salud e inventario | recetario 60; tiempos; salud, IMC, kcal; fuentes nutricionales; inventario híbrido; verticales de compra | S1, S4, U4, C1 |
| E · Precios y tiendas | fuentes autorizadas, OCR con confirmación, lectura asistida desactivada hasta veredicto, PROFECO; frescura 30 d; Sam's membresía; unidad base | S1, S5, U4, C1 |
| F · Reportes y avisos | diario, semanal, mensual PDF, vista web, estado de Karen; 4 alertas; canal MultiX texto plano | S2 (80/100 %), S7, U4, WF4 |
| G · Interfaz | Notion, barra lateral, vistas, menú `/`, atajos; Android + escritorio, PWA, mkcert; ≤ 3 toques; tema del sistema; ≤ 6 tarjetas; 14 secciones; autoalojado; WCAG 2.2 AA | U1, U2, U3, U4, U5, C2, S6, P3 |
| H · Infraestructura | máquina de casa, solo LAN; Go + SQLite + systemd; respaldos locales; IA opcional apagada; semilla | S6, S8, WF4, S1, S4 |
| I · Pila y matiz | pila A; matiz C 280° y sus roles; botón primario; sin solver externo | U2, C2, S3, S6, P3 |
| J · Puertas | puertas 1 a 4 | este índice §3, C3, P3 |

## 8. Supuestos de este agente
- Pantalla = `P-xx` (y `P-26.xx` para las secciones de Ajustes) y necesidad de datos = `N-xxx`, como los fija U1 §0.1 y los usan C1, C3 y U5. El bloqueo `criterios-nxx` quedó cerrado con la opción (b) en la ronda 1 de correcciones; la lectura inicial «N-xx = pantalla» se descartó.
- El formato de evidencia (`<ID>` al inicio de la línea en el reporte del agente) y las definiciones de severidad son convenciones de este directorio; no cambian cifras, esquema, endpoints ni tokens.
- Los nombres de pruebas, scripts y rutas auxiliares que aparecen como `<declarado por el agente>` los fija cada agente en su reporte; el criterio exige que existan y se reejecuten.
- La tolerancia de ±2° en la medición del matiz (U2-02) y la regla «≥ 2 meses» para «sesión de meses» (S6-07) son medidas operativas de este directorio, no datos del plan; ambas están señaladas en su criterio.

## 9. Cambios de la ronda 1 de correcciones (10-oct-2026, agente `correccion-criterios-r1`)
Origen: `docs/estado/wf0/REVISION-r1.md` §5, hallazgos H-06, H-07 y H-08 (severidad media). No cambia el número de criterios (248; 189 alta y 59 media) ni ningún umbral numérico. Los hallazgos H-09 a H-12 (baja) no se aplicaron en esta ronda.

| Hallazgo | Qué cambió | Archivos y criterios |
|---|---|---|
| H-06 · identificadores | Pantalla = `P-xx` (`P-26.xx` para las 14 secciones de Ajustes); necesidad de datos = `N-xxx`; componentes globales `G-xx` y flujos `F-x` (los fija U1 §0.1). C1-02 cruza `N-xxx` ↔ `operationId` y exige que cada `P-xx` tenga al menos una operación o la declaración Local; una `N-xxx` bloqueada o PROVISIONAL sin endpoint no cuenta como cubierta. El conteo `N` pasa a `T` (pantallas `P-xx` más secciones `P-26.xx`) para no chocar con `N-xxx`. Bloqueo `criterios-nxx` cerrado con la opción (b). | `U1.md` (objetivo, U1-02, U1-03, U1-04, U1-05, U1-07, U1-08, U1-09, §3, §4), `C1.md` (C1-02, C1-10, §3, §4), `C3.md` (C3-01, C3-02, C3-08, §4), `U5.md` (U5-01, U5-03, §4), `U4.md` (U4-01), `P3-claude-design.md` (P3-06) |
| H-07 · periodos | S1-09 ya no pide insertar P1–P6 (ventanas del 10 al 9): pide los meses calendario oct-2026…mar-2027 con sus quincenas, el rechazo de la ventana por el `CHECK` y la prueba de equivalencia con el oráculo (12 de 12 cifras). S2-09 y las notas de R-45…R-50 (`S2.md` §6 (e), `S3.md` §6 (d)) dicen que la equivalencia depende de la semilla; S2-06 usa una etiqueta de mes calendario; S2-03, S2-07 y el §4 de S2, C1-09, C3-09 y S7-02 (consumidores de R-45…R-50) aclaran que P1–P6 son oct-2026…mar-2027. | `S1.md` (S1-09, §3, §5), `S2.md` (S2-03, S2-06, S2-07, S2-09, §4, nota (e)), `S3.md` (nota (d)), `C1.md` (C1-09), `C3.md` (C3-09), `S7.md` (S7-02) |
| H-08 · trazabilidad | Los 35 criterios que no aparecían en la sección 5 de su archivo tienen ahora su línea (al plan, a decisiones o al protocolo; lo que no deriva de una frase se marca «convención propia del directorio»). Todo criterio de los 18 archivos aparece ya en su sección 5 (248 de 248). | C1-04; S1-15; U1-14; C2-09, C2-10, C2-11; C3-03, C3-10, C3-11; P3-13; S2-08, S2-09, S2-14, S2-16; S3-04, S3-11, S3-20; S4-06, S4-15; S5-13; S6-08, S6-10, S6-15, S6-16; S7-14; S8-08; U3-08, U3-10, U3-12; U4-07, U4-09, U4-10, U4-11; U5-06, U5-12 |

## 10. Conteo reconciliado de la puerta 1 (11-oct-2026, agente `integrador-p1`)
Conteo propio por script sobre los 18 archivos (no incluye este índice): cuenta las líneas que empiezan con `**<ID>-nn` seguidas de la severidad entre paréntesis, y comprueba que no haya IDs repetidos y que cada archivo diga «Versión 1.1». Coincide fila por fila con la tabla de §2 y con los totales de `docs/criterios/tools/traza.py` (257 criterios definidos) y `refs.py` (257 IDs citados, 0 inexistentes).

| Archivo | Criterios (alta / media / baja) | Versión |
|---|---|---|
| `S1.md` | 17 (12 / 5 / 0) | 1.1 |
| `S2.md` | 16 (13 / 3 / 0) | 1.1 |
| `S3.md` | 20 (17 / 3 / 0) | 1.1 |
| `S4.md` | 15 (11 / 4 / 0) | 1.1 |
| `S5.md` | 13 (10 / 3 / 0) | 1.1 |
| `S6.md` | 16 (12 / 4 / 0) | 1.1 |
| `S7.md` | 14 (12 / 2 / 0) | 1.1 |
| `S8.md` | 11 (8 / 3 / 0) | 1.1 |
| `WF4-operacion.md` | 12 (10 / 2 / 0) | 1.1 |
| `U1.md` | 16 (10 / 6 / 0) | 1.1 |
| `U2.md` | 14 (10 / 4 / 0) | 1.1 |
| `U3.md` | 12 (8 / 4 / 0) | 1.1 |
| `U4.md` | 11 (8 / 3 / 0) | 1.1 |
| `U5.md` | 13 (10 / 3 / 0) | 1.1 |
| `C1.md` | 21 (19 / 2 / 0) | 1.1 |
| `C2.md` | 12 (8 / 4 / 0) | 1.1 |
| `C3.md` | 11 (9 / 2 / 0) | 1.1 |
| `P3-claude-design.md` | 13 (10 / 3 / 0) | 1.1 |
| **Total (18 archivos)** | **257 (197 / 60 / 0)** | |

Reconciliación con los reportes de P1:
- `docs/estado/p1/criterios-S-1.1.md`: S1–S8 y WF4 = 134 (105 / 29), sin cambio. **Coincide.** Su total de 248 era la tabla de antes de los cambios de U1 y C1; el mismo reporte anticipaba 256 con esos cambios.
- `docs/estado/p1/criterios-UC-1.1.md`: daba 256 (196 / 60) en total. La diferencia es 1 criterio alto: U5-13, que agregó la ronda 1 de correcciones de P1 (`docs/estado/p1/correccion-criterios-r1.md`; U5 pasa a 13 = 10 / 3), con lo que el total queda en 257 (197 / 60). El conteo propio y `traza.py` confirman 257.
- Comando: `make verify-criterios` (corre `docs/criterios/tools/verif_final.sh` y exige sus resultados esperados).

Siguen registrados (hallazgo bajo B-01 de `docs/estado/p1/REVISION-criterios-r3.md`; el integrador no cambia contenido): las notas 8 y 10 de §5 todavía describen como problemas la colisión error–gasto y el token del bot, ya resueltos por K26 y K30, y §6 no lista los bloqueos nuevos `criterios-promedios-redondeo` y `criterios-frescura-extremos` (están en `docs/estado/bloqueos/INDICE.md`, sección «Bloqueos nuevos de P1»).
