# revisor-tablero-backend-r3 · reporte · WF-3c · 12-oct-2026

Agente `revisor-tablero-backend-r3`: revisor independiente, ronda 3 de 3. No corrige nada. El informe completo está en `docs/estado/wf3/REVISION-tablero-backend-r3.md`.

## Hecho
- **HECHO · Dictamen: NO APROBADO.** Hay 0 altos, 1 medio (M-2, sigue abierto) y 6 bajos (B-1 a B-5 siguen; B-6 es nuevo).
- **HECHO · M-2 sigue abierto.**
  - `PaymentsByMonth` (`backend/internal/api/creditos/tablero.go`, sin cambios desde las 14:45) cuenta solo las líneas activas y devuelve `true` fijo.
  - RN-30 y la descripción de `payments_cents` en C1 dicen «de todas las líneas».
  - El subtest `opcion_a_PROVISIONAL` sigue en SKIP.
  - Con overlay: el código actual sin SKIP falla, y con la solicitud 1 de `correccion-tablero-backend-r2.md` pasa.
  - Lo que sí se atendió: el bloqueo ya dice la verdad sobre el código y la parte común de la prueba corre siempre.
- **HECHO · Reejecuté todas las verificaciones de (a) a (e)** con comandos propios y contra el servidor de desarrollo por HTTP/2. Las cifras coinciden con el oráculo en 6 grupos de valores y en la serie de 12 meses.

## Archivos
- Nuevos:
  - `docs/estado/wf3/REVISION-tablero-backend-r3.md` (informe);
  - este reporte.
- Temporales en `datos/tmp/revisor-tablero-backend-r3/` (script de curl, overlays, bitácoras de make y el archivo de cookies): borrados.
- No modifiqué ningún archivo del código, del contrato ni de los bloqueos.

## Verificación
- **S6-01:** `make verify-wf3-go` EXIT 0:
  - go1.27.2; build y vet sin errores; gofmt vacío; tidy sin cambios;
  - `go test ./...`: 33 paquetes OK y 5 sin pruebas.
- **S6-02:**
  - `make verify-wf3-cobertura` EXIT 0. Con umbral K25: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8. Sin umbral: api 82.0, auth 83.6, jobs 93.7, reportes 77.3 y telegram 80.0;
  - `go test -cover`: api/tablero 87.8 % y api/creditos 80.4 %.
- **S6-03:** `make verify-wf3-contrato` EXIT 0:
  - «22 archivos coinciden»;
  - «221 operaciones en 176 rutas, 0 diferencias»;
  - `TestContrato_Operaciones`: «221 de 221 operaciones cumplen».
- **S6-04:** convención de C1 en las respuestas reales. Observaciones: B-1, B-2, B-4 y B-6.
- **S6-14:**
  - `ss` muestra solo `127.0.0.1:8443` y `192.168.50.56:8443`;
  - Tailscale (100.115.28.35): curl con código 7;
  - pasan las 3 pruebas `TestEscucha_*`;
  - HTTP/2 por loopback y por la LAN; ALPN h2; `--http1.1` responde en 1.1.
- **S6-15:** identificadores en inglés y mensajes en español de México. El log tiene 0 coincidencias de password, contraseña, cookie o token.
- **S2 y S3:**
  - deuda viva 1,899,922;
  - octubre 59,882 · 97,375 · 1, con 2 % y 4 % y margen 920,118;
  - noviembre y diciembre, 5 % y 6 %;
  - serie de 12 meses, 12 de 12;
  - devengado 119,764, efectivo 59,882 y plan 200,000 (60 %).
  - Todo es igual a lo que recalculé desde `oraculo/model.json` con ROUND_HALF_UP.
- **(d) Umbrales:**
  - con 1 y 2: octubre `warning` y noviembre `exceeded`;
  - con amarillo 45 contra un tope de 2: 422 en español;
  - `reset_parameters`: vuelve a `normal` con 30/40 y los 5 parámetros en `is_initial` true.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`:** EXIT 0.
- Sesión de axel con la contraseña de desarrollo leída del archivo, sin imprimirla. Logout 204 y luego 401.

## Supuestos
- Mantengo la lectura de la ronda 2: M-2 se cierra solo si el código y los artefactos dicen lo mismo. Puede ser con la salida (i) o con la (ii). Que el bloqueo describa con honestidad la diferencia no basta, porque el código sigue contradiciendo RN-30 y C1.

## Bloqueos
- No abrí bloqueos nuevos.
- `S6-tablero-lineas-archivadas` sigue abierto, con dos preguntas:
  - al orquestador: ¿la salida (i) o la (ii)?
  - a Axel, en la puerta 4: ¿la opción (a) o la (b)?

## Pendientes
- **Escalado (ronda 3 de 3):** el orquestador decide la salida de M-2 (ver §3 del informe).
- Los bajos B-1 a B-6 van al integrador o a quien tenga `c1gen` y backend-base. B-3 le toca a `integrador-wf3c`.

## Propuestas para Axel
- **PROPUESTA:** la opción (a) de `S6-tablero-lineas-archivadas`. Si archivas una tarjeta que ya pagaste, sus pagos de meses pasados siguen en el barómetro y en las gráficas, igual que en el efectivo.
- **PROPUESTA (al orquestador):** la salida (i). Es un archivo de unas 10 líneas, probado con overlay en tres rondas, y no mueve las cifras del servidor de desarrollo.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 36 acciones, unos 105 KB leídos y unos 22 KB escritos (incluido el informe).
