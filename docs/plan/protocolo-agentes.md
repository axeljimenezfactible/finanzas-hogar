# Protocolo para agentes de la Etapa B · Finanzas Hogar

Este archivo lo lee **todo agente** antes de empezar. Es parte del contrato de trabajo.

## 1. Lecturas obligatorias
1. `docs/plan/plan.md`: las secciones que indique tu prompt (siempre §1 HECHOS, §2 decisiones y la fila de tu fase en §3).
2. `docs/plan/decisiones.md`: decisiones cerradas de Axel (incluye matiz y pila).
3. Este protocolo.
4. Los criterios de tu fase en `docs/criterios/<fase>.md` cuando existan (en WF-0 se escriben en paralelo; el revisor sí los usa).

## 2. Regla suprema: no asumas
- Si te falta un dato, un criterio admite dos lecturas o algo no está en los archivos: **no inventes un valor**. Escribe un bloqueo en `docs/estado/bloqueos/<fase>-<tema>.md` con este formato y **sigue con todo lo que no depende de eso**:
  ```
  # Bloqueo <fase>-<tema>
  Contexto: qué estabas haciendo
  Qué falta: la pregunta exacta para Axel
  Opciones: 2-4 opciones con su consecuencia
  Impacto: qué queda sin hacer o marcado como provisional
  Cómo seguí: qué hiciste sin ese dato (nada provisional entra en cifras ni en esquema sin marcarse PROVISIONAL)
  ```
- Todo agente termina su reporte con una sección **SUPUESTOS** (puede ir vacía). Un supuesto que cambie cifras, esquema, endpoints o tokens es un bloqueo, no un supuesto.

## 3. Alcance disjunto
- Solo escribes en los archivos y carpetas que tu prompt te asigna. Si necesitas cambiar algo fuera de tu alcance, lo escribes como **solicitud** en tu reporte (ruta, cambio, motivo) y el integrador decide.
- No haces `git commit`: el integrador del flujo hace un commit por flujo. Tú solo dejas archivos.

## 4. Tope de 300K tokens por agente (reparto fino desde la puerta 1)
En WF-0, 13 de 14 agentes superaron 300K de contexto en su pico por leer y escribir archivos enteros muy grandes. Desde la puerta 1 rigen estos límites **por agente**:
- **Lectura ≤ 400 KB en total y escritura ≤ 150 KB en total** (suma de todo lo que lees y escribes con herramientas). Lleva la cuenta aproximada en tu archivo de progreso.
- **≤ 60 acciones de herramienta.** A las ~30 actualiza `docs/estado/<flujo>/<agente>.md` con el progreso (qué está hecho, qué falta). Si vas a pasar de 60 o lo restante no cabe, **cierra limpio**: deja todo consistente, escribe en el reporte la lista exacta de lo pendiente con rutas, marca `cerrado_por_tope: true` y termina. No comprimas ni resumas contexto para «que quepa».
- **Nunca leas enteros** `docs/contratos/openapi.yaml` (597 KB), `C1-datos.md` (113 KB), `diccionario-datos.md` (162 KB), `U1-arquitectura.md` (123 KB) ni la migración (137 KB). Usa los **extractos por módulo** de `docs/contratos/extractos/` (los genera `backend/internal/api/c1gen/build.py --extractos`) y `grep`/`sed -n` para lo demás.
- Escribe archivos grandes **por partes** (varias escrituras cortas o un script que los genere), nunca un solo volcado de más de 40 KB.
- Si una tarea no cabe en estos límites, no es tu problema resolverlo apretando: repórtalo como **Pendiente con rutas** y el orquestador la reparte.

## 5. Reporte final (obligatorio, es lo que devuelves)
Ruta: `docs/estado/<flujo>/<agente>.md`. Secciones: **Hecho** (lista), **Archivos** (rutas creadas o modificadas), **Verificación** (comandos ejecutados y su salida resumida: pruebas, compilaciones, cálculos), **Supuestos**, **Bloqueos** (rutas a los archivos de bloqueo), **Pendientes**, **Solicitudes fuera de alcance**, **Paquetes instalados** (si instalaste algo del sistema), **Consumo aproximado** (número de acciones de herramienta y tamaño de lo leído y escrito).

## 6. Cifras y dinero
- Dinero en centavos exactos (Decimal / enteros de centavos), redondeo ROUND_HALF_UP; el último pago de una serie ajusta centavos. MXN, símbolo `$`, dos decimales. Fechas en pantalla `10-oct-2026`; en archivos y API, ISO `2026-10-10`. Zona horaria `America/Mexico_City`.
- Toda cifra que exista en `oraculo/model.json` debe reproducirse exactamente. Si tu cálculo difiere, es un hallazgo: documenta ambos valores, no «ajustes» el tuyo.

## 7. Red, credenciales y sistema
- Internet solo para las fuentes listadas en el plan (PROFECO, USDA FoodData Central, Open Food Facts, documentación oficial de las herramientas). Nada de scraping de tiendas en WF-0 ni WF-1; la lectura asistida de sitios queda desactivada hasta que `docs/estado/S5-terminos.md` tenga veredicto.
- Nunca escribas, imprimas ni copies credenciales. El token de Telegram se lee en runtime desde `~/.openclaw/openclaw.json` y jamás se imprime; los **agentes no mandan Telegram**, lo hace el orquestador.
- Puedes instalar paquetes del sistema (Axel lo autorizó sin lista previa) solo si tu tarea lo exige; anótalo en **Paquetes instalados**.
- No toques nada fuera de `~/finanzas-hogar/` salvo instalaciones de paquetes y lecturas de `~/reporte-compra-latex/`.

## 8. Idioma y estilo
- Documentos, interfaz, mensajes al usuario y comentarios de negocio: **español de México**, sin jerga innecesaria.
- Identificadores de código (Go, TypeScript, SQL): inglés corto y consistente; los nombres de dominio del plan se traducen una sola vez en `docs/contratos/diccionario-datos.md` (columna `nombre_codigo`).
- Nada de dependencias que operen con servicios de terceros en runtime (fuentes, iconos y librerías autoalojados).

## 9. Lo que devuelves como salida final
Tu texto final **es** el valor que recibe el orquestador: devuelve el reporte en el formato estructurado que te pida el prompt (normalmente un JSON con rutas, supuestos y bloqueos). Nada de narrativa para humanos.
