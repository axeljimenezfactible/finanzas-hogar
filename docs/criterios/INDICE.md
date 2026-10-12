# Índice de criterios de aceptación · Finanzas Hogar (Etapa B)

Versión 1.2 (11-oct-2026, flujo P1c, puerta 1b: respuestas K35 a K44 aplicadas, conteo propio con `tools/conteo.py`, cambios en §11 y correcciones de la revisión r1 de criterios (A-01 a A-03, M-01, M-02) y de la revisión r2 (hallazgos 1 a 4, fila de §11) la ronda 4 de P1d (hallazgos A a G de la revisión r3, I-02, B-01 de r2 de C1 y U3-01; filas de §11) y la ronda 5 de P1e (hallazgos H-1 a H-6 de la revisión r4; filas de §11); versión 1.1 del 10-oct-2026 con filas y totales de la puerta 1 y conteo reconciliado por el integrador de P1, §10) · Escrito en WF-0 por el agente «criterios»; corregido en la ronda 1 de revisión (§9) · Plan: `docs/plan/plan.md` §3.1 principio 4 («Criterios antes del trabajo… los revisores evalúan solo esos criterios»). Cada criterio se traza al plan, a `docs/plan/decisiones.md` o al protocolo (la parte que no deriva de una frase de esas fuentes se marca «convención propia del directorio» en su sección 5); los huecos son bloqueos en `docs/estado/bloqueos/criterios-*.md`.

## 1. Cómo usar estos archivos
- **Un archivo por ítem** (18 + este índice). Cada uno trae: 1 Objetivo · 2 Criterios numerados `<ID>-nn` (enunciado medible, «Verifica» con comando o cálculo, «Evidencia», severidad) · 3 No aceptable · 4 Umbrales numéricos · 5 Trazabilidad.
- **Severidad:** *alta* bloquea la puerta; *media* el revisor la devuelve y debe quedar corregida o justificada por escrito antes de la ronda final; *baja* se registra en el informe del revisor. (Hoy los 261 criterios de la tabla son alta o media: 201 alta y 60 media; ninguno baja.)
- **Evidencia:** cada criterio deja una línea que **empieza con su ID** en `docs/estado/<flujo>/<agente>.md` (flujos `wf0`, `wf1`, `wf3`, `wf4`), con el comando ejecutado y su salida resumida; los archivos con nombre fijo del plan son `docs/estado/S3-escenarios.md`, `S5-terminos.md`, `S5-cobertura-profeco.md`, `U5-revision.md`, `P3-iter-N.md` e `informe-final.md`. El revisor **reejecuta** los comandos; no se aprueba por la sola palabra del agente.
- **Lo que depende de Axel** (instalar la CA en los Android, abrir la PWA en Android, reiniciar la máquina, que Ximena escriba al bot, aprobar la referencia visual final en la URL privada del lienzo: la Pista 3 la ejecuta el orquestador con Artifacts, K45, que sustituye a K35) **no se declara cumplido** hasta que Axel lo confirme: queda «pendiente de Axel» o «no verificado» en el informe final.
- **Rutas y comandos:** absolutos, desde `/home/deepsec/finanzas-hogar`. `jq` no está instalado en la máquina (se usa `python3 -I`). Los scripts auxiliares los entrega cada agente y se nombran en su reporte; si falta una herramienta (`goose`, `govulncheck`, `mkcert`, validador OpenAPI, `style-dictionary`, `jsonschema`), se instala (autorizado) y se anota en «Paquetes instalados».
- **Identificadores de U1 y C1:** pantalla `P-xx` (inventario de U1 §2.1; `P-26.xx` para las 14 secciones de Ajustes, §2.3); necesidad de datos `N-xxx` (U1 §3; tipo L lectura o E escritura); componente global `G-xx`; flujo clave `F-x` (U1 §0.1). En los criterios, «pantalla» es siempre `P-xx` y «necesidad» es siempre `N-xxx`; T es el número de pantallas `P-xx` más secciones `P-26.xx` que declara U1-02.
- **Cifras:** toda cifra de `oraculo/model.json` (sha256 `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`, 596 hojas numéricas) se reproduce exactamente; si difiere es hallazgo (protocolo §6).

## 2. Ítems, flujo, puerta y conteo
| Archivo | Ítem | Se construye en | Lo evalúa | Criterios (alta / media) |
|---|---|---|---|---|
| `S1.md` | Dominio y datos | WF-0 | Revisor WF-0 → puerta 1 | 17 (12 / 5) |
| `C1.md` | Contrato de datos y API (OpenAPI) | WF-0 | Revisor WF-0 → puerta 1 | 23 (21 / 2) |
| `U1.md` | Arquitectura de información | WF-0 | Revisor WF-0 → puerta 1 | 16 (10 / 6) |
| `S2.md` | Motor financiero (85 cifras canon) | WF-1 | Revisor motores → puerta 2 | 16 (13 / 3) |
| `S3.md` | Motor de créditos y optimizador (85 cifras canon) | WF-1 | Revisor motores → puerta 2 | 20 (17 / 3) |
| `S4.md` | Recetas y comidas | WF-1 | Revisor recetas+precios → puerta 2 | 15 (11 / 4) |
| `S5.md` | Precios y tiendas | WF-1 | Revisor recetas+precios → puerta 2 | 13 (10 / 3) |
| `U2.md` | Sistema de diseño y DESIGN.md | WF-0 (colorimetría) y WF-1 | Revisor diseño → puerta 2 | 14 (10 / 4) |
| `C2.md` | Contrato de tokens | WF-1 | Revisor diseño → puerta 2 | 12 (8 / 4) |
| `U3.md` | Dashboard | WF-1 | Revisor diseño → puerta 2 (y U5) | 12 (8 / 4) |
| `U4.md` | Ajustes estilo Notion | WF-1 | Revisor diseño → puerta 2 (y U5) | 13 (10 / 3) |
| `P3-claude-design.md` | Pista 3 · Claude Design | Pista 3, después de la puerta 2; termina en la puerta 3 (la ejecuta el orquestador con Artifacts, K45; antes K35) | Revisor visual Sonnet; puerta 4 | 13 (10 / 3) |
| `S6.md` | Backend local | WF-3 | Revisores backend ×2 → puerta 4 | 16 (12 / 4) |
| `S7.md` | Reportes | WF-3 | Revisor backend → puerta 4 | 14 (12 / 2) |
| `U5.md` | Prototipo y verificación visual | WF-3 | Revisor visual → puerta 4 | 13 (10 / 3) |
| `C3.md` | Contrato de integración | WF-3 | Revisor independiente → puerta 4 | 11 (9 / 2) |
| `S8.md` | Calidad y operación | WF-4 | Revisor final | 11 (8 / 3) |
| `WF4-operacion.md` | Operación, semilla, seguridad, informe | WF-4 | Revisor final | 12 (10 / 2) |
| **Total** | | | | **261 (201 / 60)** |

Notas: S4 y S5 también dependen de WF-1 y S6/S7 de WF-3 (plan §3.3). F1–F5 del frontend (plan §5) no tienen archivo propio: se evalúan por U3, U4, U5 y C3.

## 3. Puertas de aprobación y criterios que las gobiernan
| Puerta (decisiones J, plan §3.8) | Qué recibe Axel | Criterios que deben estar en verde |
|---|---|---|
| 1 · tras WF-0 | C1, U1, criterios de todas las fases | `C1.md`, `U1.md`, `S1.md`; este directorio completo |
| 2 · tras WF-1 | DESIGN.md, C2, regresión de motores, bloqueos acumulados | `S2.md`, `S3.md`, `S4.md`, `S5.md`, `U2.md`, `C2.md`, `U3.md`, `U4.md` |
| 3 · al final de la Pista 3 (K45; antes K35) | la URL privada del lienzo con la referencia visual final: Axel la aprueba o pide una vuelta más dentro del máximo de 3 (la Pista 3 la ejecuta el orquestador con Artifacts) | P3-01, P3-06 y P3-12 de `P3-claude-design.md` (PROPUESTA de `aplicar-k45`; la confirma el orquestador) |
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
8. **Colisión error–gasto (resuelta por K26 y K41):** con el gasto en terracota 40–45° y el error en 20° la colisión de 5° desaparece (U2-08); lo que queda es el refuerzo de Karen (333°) y Bodega Aurrera (356°), aceptado por K41 con orden fijo, forma y etiqueta directa, sin distinguir una persona y una tienda solo por color en una misma vista. El ingreso y el éxito a 0° siguen distinguiéndose por icono y contexto (o 135°, que decisiones I prevé).
9. **ΔL ≥ 0.10 entre series** no es alcanzable entre los 28 pares de 8 series con contraste ≥ 3:1 sobre fondo claro; se usa la lectura de series consecutivas (U2-09). **Resuelto por K27** (bloqueo `criterios-deltal`).
10. **Token de Telegram vs usuario del servicio (resuelto por K30):** el usuario del servicio recibe permiso de lectura (ACL) sobre `~/.openclaw/openclaw.json` y lee `bots.multix.botToken` en ejecución; el token no se copia ni se imprime (S7-06, WF4-01).
11. **Go 1.26.5 instalado frente a Go 1.27 aprobado:** resuelto por K28 (Go 1.27.2 instalado el 11-oct-2026; S6-01).
12. **Ventanas del oráculo (10 al 9) frente al mes calendario:** las decisiones B fijan el mes calendario como periodo base y P1–P6 de `model.json` van del día 10 al 9. Las 12 cifras de devengado y efectivo coinciden con oct-2026…mar-2027 solo porque la semilla no tiene eventos entre los días 1 y 9 de un mes. S1-09 pide la prueba de equivalencia con el oráculo y S2-09 la nota de que depende de la semilla (hallazgo H-07 de la ronda 1).

## 6. Bloqueos escritos (`docs/estado/bloqueos/`)
**Estado** = lo que dice el propio archivo: «RESUELTO» si trae la sección «Respuesta de Axel» (puerta 1 del 11-oct-2026 o puerta 1b); sin esa sección el bloqueo está abierto y la marca PROVISIONAL se conserva. Hoy, de los 25 de esta tabla, 24 están resueltos o respondidos y 1 está **abierto**: `S1-parametros-avisos.md` (C1-22, U4-12, S7-04, S7-08 y S7-10 lo citan con la marca PROVISIONAL sobre el lugar de `send_time`, `advance_days`, `weekday` y `day_of_month`: los CUATRO figuran en la excepción de U4-12 y C1-22 (columna «clave» `notification.params_json`, PROVISIONAL solo en el lugar, citando el bloqueo; `weekday` = 7 y `day_of_month` = 1 con valor inicial FIRME de la decisión F y «restablecer»; `send_time` y `advance_days` sin valor inicial, lo propone S7, K31, puerta 4), S7-04 y S7-08 leen el día de `weekday` y S7-10 el de `day_of_month`; los valores iniciales no cambian en ninguna opción, y los umbrales 80 y 100 % y los días antes del pago ya no son parte del bloqueo: son claves `alerts.*` de `app_setting` sembradas en 00003, decisión F y K44). Cuando traiga «Respuesta de Axel», se quitan esas marcas y la excepción de C1-22 y U4-12.

| Archivo | Pregunta para Axel | Criterios afectados | Estado |
|---|---|---|---|
| `criterios-invariante-deuda.md` | qué es «deuda viva» en el invariante | S2-07 | RESUELTO (K19) |
| `criterios-promedios.md` | fórmula de promedio diario y semanal | S2-08 | RESUELTO (K20) |
| `criterios-promedios-redondeo.md` | redondeo del promedio semanal | S2-08 | RESUELTO (K39: 219.88 y 629.35) |
| `criterios-proyeccion-12m.md` | supuestos de los meses 7–12 | S2-09 | RESUELTO (K21) |
| `criterios-credito-formulas.md` | mínimo e intereses de tarjeta; crédito de tienda | S3-08, S3-10 | RESUELTO (K18) |
| `criterios-optimizador.md` | normalización del puntaje y medidas de objetivos 2 y 4 | S3-13, S3-14 | RESUELTO (K17) |
| `criterios-frescura-peso.md` | factor de peso de precios viejos y frontera de 30 días | S5-04 | RESUELTO (K22) |
| `criterios-frescura-extremos.md` | peso entre los días 30 y 90 y después del 90 | S5-04 | RESUELTO (K40: (90 − días) / 61) |
| `criterios-ocr-tickets.md` | origen de los 5 tickets y qué es renglón correcto | S5-06 | RESUELTO (K24) |
| `criterios-meta-kcal.md` | cómo se obtiene la meta de calorías | S4-07 | RESUELTO (K23) |
| `criterios-login-bloqueo.md` | duración del bloqueo y de la sesión larga | S6-06, S6-07 | RESUELTO (K11) |
| `criterios-bind.md` | ¿Tailscale además de la LAN? | S6-14 | RESUELTO (K29: solo LAN) |
| `criterios-token-servicio.md` | cómo llega el token al servicio con usuario propio | S7-06, S7-07, WF4-01 | RESUELTO (K30) |
| `criterios-toques.md` | qué cuenta como «toque» | U1-04, U5-09, C3-08 | RESUELTO (K10) |
| `criterios-nxx.md` | «N-xx» = necesidad de datos `N-xxx`; la pantalla es `P-xx` | U1-02, U1-03, C1-02, C3-02, U5-01 | RESUELTO en la ronda 1 de correcciones con la opción (b) (sin pregunta para Axel) |
| `criterios-deltal.md` | ΔL ≥ 0.10 entre series consecutivas o todos los pares | U2-09 | RESUELTO (K27) |
| `criterios-cobertura.md` | si el 85 % aplica a recetas y precios | S4-15, S5-13, S6-02 | RESUELTO (K25) |
| `criterios-umbrales-ausentes.md` | 8 valores no fijados (horas, retención, cabeceras, adjuntos…) | S3-12, S6-09, S6-10, S7-04, S7-08, S7-10, S8-02, WF4-04, WF4-08 | RESPONDIDO (K31: cada agente propone y Axel aprueba en su puerta); adjunto 10 MB y comparador 5 s cerrados por K43 y K44 |
| `S1-ajuste-objetivo-kcal.md` | cuánto ajustan «bajar» y «subir» en la meta de calorías | S4-07, U4-13, C1-22 | RESUELTO (K36: 0 y 0, por persona) |
| `S1-fuentes-precio.md` | qué fuentes compiten en la comparativa | S5-04, S5-10 | RESUELTO (K37: opción C) |
| `U1-lista-vigente-domingo.md` | semana por defecto los domingos | U1-16 | RESUELTO (K38: la semana siguiente) |
| `U2a-rosas-ximena-bodega.md` | Karen 333° y Bodega Aurrera 356° con refuerzo | U2-02, U2-08 | RESUELTO (K41) |
| `U2a-walmart-oscuro.md` | Walmart oscuro en L 0.716 | U2-08 | RESUELTO (K42) |
| `C1-valores-k31.md` | cuatro valores de C1 (adjunto, comparador, hora y anticipación, casillas sin armar) | C1-23, S3-12, S6-10, U4-12 | RESUELTO (K43 y K44) |
| `S1-parametros-avisos.md` | dónde viven los parámetros de programación de avisos: `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1), en `notification.params_json` (opción A, K43 al pie de la letra) o en claves de `app_setting` (K44; opciones B y C, que cambiarían C1 v1.2, U4 y S7); los umbrales 80 y 100 % y los días antes del pago ya no están aquí: son claves `alerts.*` de `app_setting` (decisión F, K44) | C1-22, U4-12, S7-04, S7-08, S7-10 | **ABIERTO** (11-oct-2026, sin «Respuesta de Axel»; la marca PROVISIONAL sobre el lugar de esos cuatro parámetros se conserva) |

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

## 10. Conteo reconciliado (11-oct-2026; puerta 1 por el agente `integrador-p1`, actualizado en P1c por `criterios-1.2`)
Conteo propio por script sobre los 18 archivos (no incluye este índice): `python3 -I docs/criterios/tools/conteo.py` cuenta las líneas que empiezan con `**<ID>-nn` seguidas de la severidad entre paréntesis, comprueba que no haya IDs repetidos, que cada archivo traiga la lista «Cambios de la versión» de su versión y que las tablas de §2 y de esta sección coincidan fila por fila y en el total. Coincide con `docs/criterios/tools/traza.py` (261 criterios definidos) y `refs.py` (261 IDs citados, 0 inexistentes). La versión de cada archivo es la de su línea 3: pasan a 1.2 los archivos tocados en P1c; los demás siguen en 1.1.

| Archivo | Criterios (alta / media / baja) | Versión |
|---|---|---|
| `S1.md` | 17 (12 / 5 / 0) | 1.1 |
| `S2.md` | 16 (13 / 3 / 0) | 1.2 |
| `S3.md` | 20 (17 / 3 / 0) | 1.2 |
| `S4.md` | 15 (11 / 4 / 0) | 1.2 |
| `S5.md` | 13 (10 / 3 / 0) | 1.2 |
| `S6.md` | 16 (12 / 4 / 0) | 1.2 |
| `S7.md` | 14 (12 / 2 / 0) | 1.2 |
| `S8.md` | 11 (8 / 3 / 0) | 1.1 |
| `WF4-operacion.md` | 12 (10 / 2 / 0) | 1.1 |
| `U1.md` | 16 (10 / 6 / 0) | 1.2 |
| `U2.md` | 14 (10 / 4 / 0) | 1.2 |
| `U3.md` | 12 (8 / 4 / 0) | 1.2 |
| `U4.md` | 13 (10 / 3 / 0) | 1.2 |
| `U5.md` | 13 (10 / 3 / 0) | 1.2 |
| `C1.md` | 23 (21 / 2 / 0) | 1.2 |
| `C2.md` | 12 (8 / 4 / 0) | 1.2 |
| `C3.md` | 11 (9 / 2 / 0) | 1.1 |
| `P3-claude-design.md` | 13 (10 / 3 / 0) | 1.3 |
| **Total (18 archivos)** | **261 (201 / 60 / 0)** | |

Reconciliación con los cambios de P1c (el total pasa de 257 a 261): C1 suma C1-22 y C1-23 (21 alta / 2 media) y U4 suma U4-12 y U4-13 (10 alta / 3 media); los demás archivos conservan su número de criterios (P3 sigue en 13: P3-03, P3-04, P3-05 y P3-11 cambian de contenido, no de cantidad). Los 257 de la versión 1.1 fueron: `S1–S8` y `WF4` = 134 (105 / 29), sin cambio; las diferencias con los reportes de P1 (`criterios-S-1.1.md` daba 248 y `criterios-UC-1.1.md` 256) eran U5-13, agregado en la ronda 1 de correcciones, y los cambios de U1 y C1.
- Comando: `make verify-criterios` (corre `docs/criterios/tools/verif_final.sh`, que ahora incluye `conteo.py`).

## 11. Cambios de la versión 1.2 (11-oct-2026, flujo P1c, agente `criterios-1.2`)
Respuestas de la puerta 1b aplicadas (decisiones K35 a K44) y hallazgos bajos B-01 a B-04 de `docs/estado/p1/REVISION-criterios-r3.md`. Cada archivo tocado trae su lista «Cambios de la versión 1.2».

| Origen | Qué cambió | Archivos y criterios |
|---|---|---|
| K35 | La Pista 3 la ejecuta Axel; la puerta 3 se reduce a «paquete listo y revisado»; lista nueva del paquete con `LEEME.md` y `MANIFIESTO.md`; `iter-N/` y `siguiente-prompt.md` | `P3-claude-design.md` (P3-01 a P3-06, P3-11, §3 a §5), este índice (§1, §2, §3) |
| K45 (11-oct-2026, cierre de WF-1; `P3-claude-design.md` pasa a 1.3) | Pista 3 por Artifacts, sustituye a K35: la ejecuta el orquestador sin navegador ni subida de archivos; cifras sintéticas; todas las pantallas de U1; la puerta 3 pasa a «Axel abre la URL del lienzo y aprueba la referencia visual» | `P3-claude-design.md` 1.3 (P3-01 a P3-13, §1, §3 a §5), este índice (§1, §2, §3, §10) |
| K36 | Ajuste de calorías por persona, 0 y 0, «mantener» no ajusta | `S4.md` (S4-07), `U4.md` (U4-13), `C1.md` (C1-22) |
| K37 | `manual` compite; `site_assisted` y PROFECO solo referencia | `S5.md` (S5-04, S5-10) |
| K38 | Domingos: semana siguiente por defecto | `U1.md` (U1-16) |
| K39 | Promedio semanal: diario exacto × 7, un solo redondeo (219.88 y 629.35) | `S2.md` (S2-08) |
| K40 | Peso (90 − días) / 61 entre los días 30 y 90; 91 días no compite | `S5.md` (S5-04) |
| K41 y K42 | Karen 333° y Bodega Aurrera 356° con refuerzo; Walmart oscuro `#7AA3EE` | `U2.md` (U2-02, U2-08) |
| K43 | Adjunto 10 MB, comparador 5 s con 503, `send_time` y `advance_days`, casillas sin armar ≤ 21 | `C1.md` (C1-23), `S3.md` (S3-12), `S6.md` (S6-10) |
| K44 | Todo parámetro operativo existe y es configurable, con valor inicial y «restablecer» | `U4.md` (U4-12, U4-04), `C1.md` (C1-22), `S3.md` (S3-12) |
| Pendiente 4 de `INTEGRACION.md` §7 | C1-15 define «bloqueo abierto = sin Respuesta de Axel»; C1-21 pide `INDICE.md` en los extractos y suma igual a `openapi.yaml` | `C1.md` (C1-15, C1-21) |
| B-01 | §1 con 261; notas 8 y 10 (y 9 y 11) como resueltas; §6 con los 24 bloqueos y su estado; conteo propio y versiones | este índice (§1, §5, §6, §10), `tools/conteo.py` |
| B-02 | «Día de pago 15 y 30» separado de las quincenas de K8; C1-07 comprueba la declaración 1–15 y 16–último día | `C1.md` (C1-06, C1-07) |
| B-03 | U2-02 cita decisiones I «o U2 propone 135°»; se quita la marca de convención propia | `U2.md` (§5) |
| B-04 | Cinco citas de §5 reordenadas según su fuente (`citas.py` da 0 sin fuente) | `C1.md`, `C2.md`, `U2.md`, `U5.md` (§5) |
| Revisión r1, A-01 y A-02 (K43, K44, protocolo §2); acotada en la ronda 4 | Los parámetros de programación de avisos (`send_time`, `advance_days`, `weekday`, `day_of_month`) tienen su lugar PROVISIONAL (bloqueo abierto `S1-parametros-avisos`) y se excluyen de la comprobación «clave de `app_setting` en 00003»; los umbrales 80/100 % y los días antes del pago dejaron de ser excepción (hallazgo A de r3) | `C1.md` (C1-22), `U4.md` (U4-12) |
| Revisión r1, A-03 (K43, K44) | «Casillas sin armar» es regla FIRME (K43), sin marca; solo el valor inicial de la anticipación queda PROVISIONAL; hora, anticipación, umbrales y días antes del pago se leen de su parámetro, con prueba de cambio; S7 pasa a 1.2 | `S7.md` (S7-04, S7-08, S7-10, §3, §5) |
| Revisión r1, M-01 (§6) | §6 lista `S1-parametros-avisos` como ABIERTO (25 filas: 24 resueltas y 1 abierta) | este índice (§6) |
| Revisión r1, M-02 (K35) | P3-01 es solo «paquete listo y revisado»; el aviso de Axel «iter-N lista» pasa a P3-06 y la Pista 3 queda «pendiente de Axel» hasta que avise | `P3-claude-design.md` (P3-01, P3-06) |
| Revisión r2, hallazgo 1 (K35) | P3-02 pide su evidencia en el reporte del agente `paquete-claude-design` y de su revisor en `docs/estado/wf1/` (`ls -l` y `sha256sum -c`), no en `P3-iter-1.md`, que no existe al evaluar la puerta 3 | `P3-claude-design.md` (P3-02) |
| Revisión r2, hallazgo 2 (K43, K31, protocolo §2) | `send_time` y `advance_days` no tienen valor inicial decidido (lo propone S7, K31): la tabla de U4-12 los marca «sin valor inicial: lo propone S7, PROVISIONAL» y esa celda no cuenta como vacía; la lista mínima son 14 parámetros enumerados (16 desde la ronda 5, con `weekday` y `day_of_month`) | `U4.md` (U4-12), `C1.md` (C1-22) |
| Revisión r2, hallazgo 3 (K44) | Los cruces de C1-22 y U4-12 leen las claves sembradas en 00001 a 00003 (o el inventario `@param` de 00003) y aceptan un comentario que cita una decisión de §A–§K (letra o K-n) | `C1.md` (C1-22), `U4.md` (U4-12) |
| Revisión r2, hallazgo 4 (K35, §1) | Las Verifica de P3-01, P3-03 y P3-12 buscan patrones de credencial (`password *[:=]`, `bot[_-]?token`, `api[_-]?key`, formato de token de bot) y no la palabra suelta «token», «clave» o «contraseña» | `P3-claude-design.md` (P3-01, P3-03, P3-12) |
| Revisión r3, hallazgo A (K44, decisión F; esquema 1.3, corrección S1-r1 H-01) | Umbrales 80/100 % y días antes del pago [1, 0] = claves `alerts.budget_warning_pct`, `alerts.budget_limit_pct` y `alerts.payment_due_days_before` de `app_setting` (00003, «Decisión F, K44», `getAlertSettings` y `updateAlertSettings` con valor actual e inicial); se quita «no se siembran claves `alerts.*`» y la comprobación contra `params_json`; la excepción PROVISIONAL queda para `send_time`, `advance_days`, `weekday` y `day_of_month`; §6 reducido | `C1.md` (C1-22), `U4.md` (U4-12), `S7.md` (S7-04, S7-08, S7-10, §5), este índice (§6, §11) |
| Revisión r3, hallazgos B, C, D y F (K35) | P3-06 cuenta solo lo descargado (sin `siguiente-prompt.md`; archivo o página → artboard en `P3-iter-N.md`; `P-xx` es sugerencia); P3-04 comprueba 0 acciones de navegador con ruta absoluta del registro; patrones de credencial con valor de 6 o más caracteres; P3-11 comprueba también `salida/final/` y entra a la puerta 3 (§3) | `P3-claude-design.md` (P3-01, P3-03, P3-04, P3-06, P3-11, P3-12), este índice (§3) |
| Revisión r3, hallazgo E (K44) | `send_time` y `advance_days` llevan la sección que U4 proponga, con justificación en SUPUESTOS; sin valor inicial, su «restablecer» se activa cuando S7 proponga el valor | `U4.md` (U4-12) |
| Revisión r3, hallazgo G (K43) | S6-10 dice que el factor de bytes por MB lo fija C1 (propuesta de S1 1 MB = 1,048,576 bytes, a confirmar en la puerta 2); no se inventa | `S6.md` (S6-10) |
| I-02 de `docs/estado/p1c/INTEGRACION.md` | `S4-meta-ajustes` es un bloqueo que S4 abriría si hiciera falta, no uno existente | `S4.md` (S4-07) |
| B-01 de `REVISION-C1-r2` | C1-13 dice C1 v1.2 / `info.version` 1.2.0; C1-21 dice 210 operaciones (total de `extractos/INDICE.md`) | `C1.md` (C1-13, C1-21) |
| K38, U1 §0.2 | La tarjeta (4) de U3-01 abre por defecto la semana que sigue K38 | `U3.md` (U3-01) |
| Revisión r4, H-1 (K44, decisión F, protocolo §2) | La excepción del bloqueo abierto `S1-parametros-avisos` en U4-12 cubre los cuatro parámetros `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1): columna «clave» `notification.params_json` con PROVISIONAL solo en el lugar; `weekday` y `day_of_month` con valor inicial FIRME (decisión F) y «restablecer»; `send_time` y `advance_days` sin valor inicial (S7, K31); sección propuesta por U4 con justificación en SUPUESTOS; cruce aparte contra el `@param_sin_valor` de 00003 y el `params_json` de 00001, fuera del cruce «clave de `app_setting`» hasta la «Respuesta de Axel»; la lista mínima pasa de 14 a 16 filas; §6 lo refleja | `U4.md` (U4-12), este índice (§6) |
| Revisión r4, H-2 (K44, decisión F) | S7-04 y S7-08 (alerta 4) leen el día de `weekday` (7 = domingo, FIRME; lugar PROVISIONAL según `S1-parametros-avisos`) y `weekday` entra a las pruebas de configuración (cambio y restablecer) | `S7.md` (S7-04, S7-08) |
| Revisión r4, H-3 (B-01 de `REVISION-C1-r2`) | La Verifica de C1-12 excluye las propiedades con `x-app-setting` y nombra `auth.password_min_length` como excepción | `C1.md` (C1-12) |
| Revisión r4, H-4 (protocolo §2) | P3-04 cuenta llamadas `mcp__claude-in-chrome__*` en el registro de herramientas del orquestador, no texto de un reporte; P3-03 vuelve al patrón `(password\|contraseña\|clave) *[:=] *[^ *]{6,}` con justificación escrita del revisor para una coincidencia | `P3-claude-design.md` (P3-03, P3-04) |
| Revisión r4, H-5 (decisión F, K31) | S7-10: el valor inicial de la **hora** lo propone el agente (PROVISIONAL, K31); `day_of_month` = 1 es FIRME | `S7.md` (S7-10) |
| Revisión r4, H-6 (K38, K44) | U3 §5 suma «decisiones K38 → U3-01»; K38 aplica al destino de la tarjeta (4) (P-10), su contenido sigue U1 N-023 «Semana en curso»; U3 no necesita K44 (no fija ni lee umbrales ni días) | `U3.md` (U3-01, §5) |
| Conteo (rondas 4 y 5) | Sin criterios nuevos ni quitados: total **261 (201 / 60)** sin cambio; §2 y §10 no cambian | este índice |
