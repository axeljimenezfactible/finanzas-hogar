# revisor-backend-areas-r3 · 12-oct-2026

Revisor independiente (ronda 3 de 3) del backend de áreas: api-compras-a, api-compras-b, api-creditos, api-ajustes, s7-avisos, s7-pdf y pruebas-contrato, más la corrección `correccion-backend-areas-r2`. Lo comparé solo contra S6-03, S6-04, S6-15, S6-16 y S7-01 a S7-14 (S7-07 es del orquestador).

No corregí nada. No hice git add ni commit, no usé la red ni Telegram y no usé credenciales reales: el token y la credencial de servicio son falsos y se arman en ejecución. HEAD = `fae1890`.

## Hecho

- **Veredicto: APROBADO, con 0 altos, 0 medios y 8 bajos.** El detalle está en `docs/estado/wf3/REVISION-backend-areas-r3.md`.
- **H-1 (alta, ronda 2) está CERRADO.** `getDashboardBestBasket` ya responde 200. El contrato da 216 de 216 y 0 sin implementar, y `make verify-all` da RC=0.
- **Volví a correr cada verificación con comandos propios:**
  - generador y paridad de C1;
  - contrato completo y tabla operación → prueba (216 filas, 0 FALLA);
  - las 7 pruebas literales de S7 (`Cruce`, `Karen`, `Resumen`, `Alerta`, `Mensual`, `VistaMensual` y `telegram`);
  - compilación LaTeX desde cero;
  - `pdftotext` con las expresiones de S7-12;
  - `definecolor` contra C2 y el contraste;
  - `go vet`, la cobertura y el grep de URL;
  - el grep del patrón de token en todo el repositorio, que da 0.
- **Pruebas propias, con overlay y sin tocar el repositorio:**
  - **Inventario de operaciones.** `TestRevR3_AreasOps` cuenta 109 de 109 operaciones con `TestOp_` y elige 5 al azar por área.
  - **Pruebas por área.** Escribí `TestRevR3_ComprasA`, `ComprasB`, `Creditos`, `Ajustes`, `Avisos` y `Reportes`: 30 operaciones, más K53, el comparador con `app_setting`, la idempotencia, el servidor falso de Telegram y el archivo sin clave. Todas en PASS.
  - **Ajustes a mis pruebas.** En la primera corrida 3 de ellas fallaron por supuestos míos, no por el código:
    - la lista agrupa los renglones en `groups[].items`;
    - el pago va por la parte de quien paga;
    - el CSV exige que la parte venga como `text/csv`.

    Las ajusté y volví a correrlas.

## Archivos

- `docs/estado/wf3/REVISION-backend-areas-r3.md` (informe)
- `docs/estado/wf3/revisor-backend-areas-r3.md` (este reporte)
- Temporales en `datos/tmp/rev-ba-r3/`: el registro de `verify-all`, la base de `goose`, las pruebas de overlay, la tabla, los `.tex`, los PDF y los PNG. **Ya los borré.**

## Verificación

- **S6-03**:
  - `build.py --check` → 23 idénticos.
  - `verify_c1.py --db` → TODO OK.
  - `TestContrato_Rutas` → «servidor 216, sobrantes 0, faltantes 0, sin implementar 0».
  - `TestContrato_Operaciones` → «216 de 216».
  - Por área: comprasa 18/18, comprasb 23/23, creditos 29/29, ajustes 19/19, avisos 10/10 y reportes 10/10, o sea 109 de 109.
- **S6-04**:
  - `CapturasAtrasadas` y `PaginaMaxima` en PASS.
  - Mis pruebas: `limit=201` → 400; misma clave → `Idempotent-Replayed: true` y 1 renglón; otro cuerpo → 409; clave corta → 400; centavos con decimales → 400/422.
- **S6-15**: `-run Log` → ok. De 20 mensajes al azar, 0 en inglés. Hay identificadores técnicos entre paréntesis (B-8).
- **S6-16**: 7 de 7 reportes con Supuestos y Bloqueos. Hay 25 bloqueos `S6-*` y `S7-*`.
- **S7-01**: con `latexmk`, rc 0, 0 «!» y 0 «undefined» en los dos PDF; 2 + 1 páginas y 3 PNG.
- **S7-02**: M = 256 en 6 PDF, 0 diferencias, y R-45…R-50 = oráculo.
- **S7-03**: 20 de 20, 0 diferencias.
- **S7-04**: diario a las 21:30 y semanal el domingo a las 19:00. Un cambio de hora o de día se respeta y se restablece, y `params_json` queda `{}`.
- **S7-05**: sin `parse_mode`; 12,296 bytes → 3 trozos.
- **S7-06**: el grep del token da 0. `TestToken` en PASS. El archivo sin la clave deja la entrega en `failed` con un mensaje claro y sin el token.
- **S7-07**: no lo evalúo (es del orquestador).
- **S7-08**: la matriz de 11 fechas da 0 diferencias; el presupuesto avisa una vez en 80 y otra en 100; plan «1 de 21»; configuración y restablecer correctos.
- **S7-09**: `Destinatarios` en PASS. Karen da 422 en confirmar y en probar.
- **S7-10**: PDF v1 de 96,403 bytes, por debajo de 50 MB. Con la plantilla rota → error y 0 envíos. El día 31 cae en el 30-nov.
- **S7-11**: 34 de 34. El endpoint está en C1 y P-24 en U1.
- **S7-12**: 0 fechas ISO y 0 cifras fuera del patrón en el PDF y en los textos.
- **S7-13**: 7 de 7 colores en C2. Contraste mínimo 5.04.
- **S7-14**: `vet` limpio; cobertura de 77.3 %, 77.3 % y 80.0 %. La única URL es `api.telegram.org`. LaTeX corre con tiempo límite.
- **Global**: `make verify-all GOOSE=…` → RC=0.

## Supuestos

- Mis pruebas usan los ayudantes de `testutil`. En comprasb uso además el verificador falso `fakeService` del paquete, porque no hay otra forma de entrar como servicio sin una credencial real. Lo que compruebo lo escribí yo.
- La revisión de los 20 mensajes la hice sobre literales que extraje de las 6 áreas con grep. No incluye los mensajes que se arman en tiempo de ejecución, salvo `badEnum`.

## Bloqueos

Ninguno nuevo. Los de los agentes siguen abiertos y documentados: `S7-fuentes-avisos`, `S5-agente-externo` (e) y `S6-gastos-inicio` §15.

## Pendientes

Los 8 bajos (B-1 a B-8) quedan para sus dueños. Ninguno bloquea la aprobación.

## Propuestas para Axel

- **PROPUESTA**: conectar `reportes.SetSender` en el binario (B-5) antes de la puerta 4. Sin eso, el reporte mensual del día 1 se guarda, pero no sale por Telegram.
- **PREGUNTA (de s7-avisos, sigue abierta)**: ¿cuál es el «día de despensa» para el aviso de precios viejos? Mientras no se defina, ese aviso no sale con la base real (B-4).

## Paquetes instalados

Ninguno.

## Consumo aproximado

Unas 52 acciones, unos 190 KB leídos y unos 60 KB escritos (pruebas temporales e informes).
