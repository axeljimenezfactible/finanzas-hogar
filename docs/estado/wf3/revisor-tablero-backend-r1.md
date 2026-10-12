# revisor-tablero-backend-r1 · reporte · WF-3c · 12-oct-2026

Revisor independiente, ronda 1 de 3, del trabajo de `api-tablero`. No corrige nada. El informe completo está en `docs/estado/wf3/REVISION-tablero-backend-r1.md`.

## Hecho
- **HECHO:** reejecuté todas las verificaciones con comandos propios, y con el servidor de desarrollo en `https://finanzas.local:8443` (HTTP/2):
  - `make verify-wf3-go`, `verify-wf3-cobertura`, `verify-wf3-contrato` y `verify-all`;
  - `go test -cover` de los paquetes tocados.
- **HECHO:** recalculé contra `oraculo/model.json` y contra los motores, y todo coincide:
  - pagos del hogar y totales de oct, nov y dic-2026;
  - devengado y efectivo de octubre;
  - ratio, margen y plan.
- **HECHO:** comprobé que los umbrales salen de `app_setting`: con 1/2 el color cambia a warning y exceeded, y al restablecer vuelve a normal y 30/40.
- **HECHO:** comprobé HTTP/2 en LAN y en loopback, HTTP/1.1 cuando se fuerza, y que no hay escucha en Tailscale.
- **HECHO:** con una prueba temporal propia (`go test -overlay`, sin tocar el árbol) confirmé el hallazgo M-2, las líneas archivadas.
- **Dictamen:** NO APROBADO, con 0 altos, 2 medios (M-1 y M-2) y 4 bajos (B-1 a B-4).

## Archivos
- `docs/estado/wf3/REVISION-tablero-backend-r1.md` (nuevo)
- `docs/estado/wf3/revisor-tablero-backend-r1.md` (nuevo, este reporte)

## Verificación
- **S6-01:** `make verify-wf3-go` EXIT 0, con go1.27.2 y go.mod 1.27.
  - build, vet, gofmt y `tidy -diff` limpios;
  - 33 paquetes OK y 5 sin pruebas.
- **S6-02:** `make verify-wf3-cobertura` EXIT 0.
  - Con umbral K25: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8.
  - Sin umbral: api 82.0, auth 83.6, jobs 93.0, reportes 77.3 y telegram 80.0.
  - api/tablero 87.6 % y api/gastos-inicio 82.0 %.
- **S6-03:** `make verify-wf3-contrato` EXIT 0.
  - «22 archivos coinciden byte a byte»;
  - «221 operaciones en 176 rutas, 0 diferencias»;
  - `TestAreas_SinOperacionesEn501` 221/221 y `TestContrato_Operaciones` 221/221.
- **S6-04:** cumple: `AAAA-MM`, centavos enteros, ErrorBody único, 400 de formato y 422 de regla, mensajes en es-MX. Tiene observaciones B-1, B-2 y B-4.
- **S6-14:** cumple.
  - Escucha solo en 127.0.0.1:8443 y 192.168.50.56:8443; por Tailscale curl sale con código 7.
  - HTTP/2 por ALPN (`h2`).
  - Las tres pruebas `TestEscucha_*` pasan.
- **S6-15:** cumple. Identificadores en inglés, mensajes en es-MX y 0 secretos en `datos/dev/run/finanzas.log`. Tiene la observación B-2.
- **S2/S3 (cifras):** 7 grupos de valores coinciden con el oráculo y con los motores (tabla en el §2 del informe).
- **C1:** RN-31 no se cumple con `period` futuro (M-1). RN-30, «de todas las líneas», no se cumple con líneas archivadas (M-2).
- **verify-all** (`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`): EXIT 0.

## Supuestos
- El alcance de la revisión de M-2 incluye `api/creditos/tablero.go`, porque es la fuente S3 de las dos operaciones del tablero que revisé. Si no le toca a api-tablero, su corrección va como solicitud fuera de alcance.
- B-1 se clasifica baja y no media: la misma frase de C1 ya existía en 8 operaciones de settings y la implementación sigue C1 §3.

## Bloqueos
Ninguno nuevo; el revisor no abre bloqueos. Para M-1 y M-2 propongo abrir bloqueos PROVISIONALES; las opciones están en el informe.

## Pendientes
- Ronda 2: cerrar M-1 y M-2.
- Base de desarrollo:
  - quedó con los valores iniciales del tablero (`is_initial` = true);
  - la bitácora conserva los PATCH de prueba;
  - la sesión de revisión se cerró (logout 204).

## Propuestas para Axel
- **PREGUNTA** (M-1): si se pide la sugerencia de un mes futuro, ¿qué prefieres? Que no salga ninguna sugerencia (los meses aún no cerrados no cuentan), o que la API rechace la petición.
- **PREGUNTA** (M-2): cuando archivas una tarjeta ya liquidada, ¿sus pagos de meses pasados deben seguir en las gráficas y en el barómetro de esos meses? Propongo que sí, como dice RN-30: «de todas las líneas».

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 47 acciones; unos 140 KB leídos y unos 19 KB escritos. Los temporales de `datos/tmp/revisor-tablero-backend-r1/` (log, cookie, overlay y prueba temporal) se borraron.
