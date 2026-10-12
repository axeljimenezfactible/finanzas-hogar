# Runbook · Finanzas Hogar (BORRADOR)

Versión 0.1 · 12-oct-2026 · agente **integrador-backend** (WF-3a). **Borrador:** WF-4 lo completa (unidad systemd, temporizador, usuario del servicio, rutas de producción, permisos del token de Telegram, procedimiento de restauración desde Ajustes). Todo lo marcado **PROVISIONAL** espera a Axel en la puerta 4 (K31) o a WF-4.

Convención: HECHO = comprobado en esta máquina el 12-oct-2026 · PREGUNTA = falta decisión de Axel · PROPUESTA = valor del agente con su justificación.

## 1. Dirección (S6-14)

| Qué | Valor | Estado |
|---|---|---|
| Interfaz y dirección de la LAN | `wlp4s0` · **192.168.50.56**/24 (red 192.168.50.0/24) | HECHO hoy; **PROVISIONAL** hasta WF-4 |
| Puerto | **8443** (HTTPS, sin HTTP plano) | FIRME (S6-11) |
| Escucha | solo `192.168.50.56:8443` y `127.0.0.1:8443` (`FINANZAS_ADDR`) | HECHO (`ss -ltnp`) |
| Tailscale | `tailscale0` (100.115.28.35) **no** se usa: el servicio no escucha ahí | FIRME (K29); HECHO |
| Comodines | `serve` rechaza `0.0.0.0`, `::` y `:8443` | HECHO (código de `cmd/finanzas`) |

- PREGUNTA (WF-4/Axel): la IP llega por DHCP del módem. **PROPUESTA:** reservar 192.168.50.56 para esta máquina en el módem (reserva DHCP por MAC de `wlp4s0`), porque el certificado y `FINANZAS_ADDR` la llevan escrita; si cambia, hay que reemitir el certificado (`ops/mkcert.md` §2) y cambiar `FINANZAS_LAN_IP`.
- `enp5s0` (cable) está abajo hoy. Si la máquina pasa a cable, cambia la interfaz y quizá la IP: mismo procedimiento.

## 2. Nombre `finanzas.local`

| Dónde | Cómo | Estado |
|---|---|---|
| Esta máquina | mDNS con avahi: `avahi-publish -a -R finanzas.local 192.168.50.56` en segundo plano (lo arranca `ops/dev-run.sh` si el nombre no resuelve; lo detiene `ops/dev-stop.sh`). `avahi-daemon` activo y `/etc/nsswitch.conf` con `mdns_minimal` | HECHO: `getent hosts finanzas.local` → 192.168.50.56 |
| Otras computadoras de la LAN | el mismo anuncio mDNS (Linux con nss-mdns, macOS, Windows 10+) | sin probar |
| Android (Axel y Ximena) | Chrome no siempre resuelve `.local` | PREGUNTA (bloqueo `S6-sistema`) |

- No se usó `/etc/hosts` (habría requerido sudo, que solo está permitido para instalar mkcert).
- PROPUESTA para producción (WF-4): en lugar del proceso `avahi-publish`, un alias fijo de avahi (`/etc/avahi/hosts` con `192.168.50.56 finanzas.local`, requiere sudo y reiniciar `avahi-daemon`) o una unidad systemd del usuario que corra `avahi-publish`. Para los teléfonos, la opción robusta es una entrada DNS en el módem (`finanzas.local → 192.168.50.56`); si el módem no lo permite, abrir por IP: `https://192.168.50.56:8443/` (la IP está en el certificado).

## 3. Cómo correr

### 3.1 Desarrollo (WF-3) · HECHO

```bash
cd /home/deepsec/finanzas-hogar
bash ops/dev-run.sh             # compila backend/bin/finanzas; si no existe datos/finanzas-dev.db: migra y siembra; arranca en segundo plano
bash ops/dev-run.sh --sembrar   # además vuelve a correr migrate + semilla (idempotente)
bash ops/dev-stop.sh            # SIGTERM (a los 10 s SIGKILL) al servidor y a avahi-publish
bash ops/dev-stop.sh --conservar-nombre   # detiene solo el servidor
make verify-wf3-servidor        # S6-11, S6-12 y S6-14 contra el servidor corriendo
```

| Qué | Dónde (todo ignorado por git) |
|---|---|
| Base sembrada | `datos/finanzas-dev.db` (+ `-wal`, `-shm`) |
| Adjuntos y reportes | `datos/dev/app/` |
| Respaldos del job diario | `datos/dev/respaldos/` |
| PID, log del servidor, log de la semilla, PID de avahi-publish | `datos/dev/run/finanzas.pid`, `finanzas.log`, `semilla.log`, `avahi-publish.pid` |
| HOME aislado del proceso | `datos/dev/home/` (vacío a propósito: ver Telegram abajo) |
| Certificado y llave | `datos/certs/finanzas.local.pem`, `finanzas.local-key.pem` (`ops/mkcert.md`) |
| Contraseñas iniciales de la base de DESARROLLO | `~/.config/finanzas-hogar/dev-semilla.txt` (fuera del repositorio, modo 600, `axel=…` y `ximena=…`, aleatorias; el script nunca las imprime) |

- **Telegram apagado en desarrollo:** el servidor corre con `HOME=datos/dev/home`, donde no existe `.openclaw/openclaw.json`; los jobs `daily_summary`, `weekly_summary`, `alerts` y el reporte mensual corren pero no tienen token y no envían nada. Los envíos reales los dispara solo el orquestador (S7-07, `ops/telegram-disparo.md`).
- Variables que acepta `ops/dev-run.sh`: `FINANZAS_LAN_IP` (192.168.50.56), `FINANZAS_PORT` (8443), `FINANZAS_DEV_DB`, `FINANZAS_SEED_PASSWORD_FILE`.
- Para entrar a la base de desarrollo (WF-3b): usuario `axel` o `ximena`, contraseña en el archivo de arriba (leerlo, nunca copiarlo a reportes ni capturas).
- Para empezar de cero: `bash ops/dev-stop.sh && rm -rf datos/finanzas-dev.db* datos/dev && bash ops/dev-run.sh`.

### 3.2 Producción · PENDIENTE (WF-4)

PROPUESTA de partida para WF-4: el mismo `backend/bin/finanzas serve` con las variables de `ops/mkcert.md` §3, como unidad systemd con `Restart=on-failure`, usuario dedicado sin shell, `FINANZAS_DB`, `FINANZAS_DATOS` y `FINANZAS_RESPALDOS` en una carpeta fuera del repositorio, y la PWA embebida en el binario (bloqueo `WF3-frontend-servidor`, opción A). El orden de la primera carga es: migrar → dejar el archivo de contraseñas iniciales fuera del repositorio (modo 600) → `finanzas semilla` → entrar y cambiar las contraseñas en Ajustes › Cuenta y roles (`docs/estado/wf3/semilla.md`).

## 4. Cómo respaldar (S6-09)

| Qué | Valor | Estado |
|---|---|---|
| Método | `VACUUM INTO` + `PRAGMA integrity_check` + restauración de prueba en carpeta temporal y conteos por tabla | HECHO (`dominio.Backup`) |
| Dónde | solo local: `FINANZAS_RESPALDOS` (desarrollo: `datos/dev/respaldos/`) | FIRME: nunca fuera del equipo (decisión H) |
| Hora del respaldo diario | **03:15** de `America/Mexico_City` | **PROVISIONAL** (K31, puerta 4) |
| Retención | **30** respaldos diarios; los manuales y los previos a restaurar no se borran solos | **PROVISIONAL** (K31, puerta 4) |

- **Justificación de la PROPUESTA** (bloqueo `docs/estado/bloqueos/S6-respaldos.md`, constantes en `backend/cmd/finanzas/main.go`): a las 03:15 la casa duerme y nadie captura; queda después del resumen diario (21:30) y antes del reporte mensual (08:00 del día 1), así el respaldo nunca compite con ellos. 30 diarios = un mes para volver atrás; la base sembrada pesa ~1.8 MB, así que 30 copias ocupan ~55 MB.
- PROPUESTA (WF-4/Axel): una segunda copia en **otro disco** de la misma máquina (S6-09 pide «otra carpeta o disco»); hoy la carpeta está en el mismo disco que la base.
- Respaldo manual:

```bash
backend/bin/finanzas respaldo -db datos/finanzas-dev.db -respaldos datos/dev/respaldos -tmp datos/tmp -tipo manual
# imprime: respaldo N: <archivo> (bytes, sha256), integrity_check: ok, restauración de prueba: ok y los conteos por tabla
make verify-wf3-respaldo   # lo mismo sobre una COPIA de la base de desarrollo en datos/tmp, con la comparación de conteos
```

- Comprobado el 12-oct-2026 sobre copia de `datos/finanzas-dev.db`: `integrity_check` = ok en origen y restaurada; 56 tablas, 2,476 filas, 0 diferencias de conteo.
- Restaurar: desde Ajustes › Datos y respaldos (`restoreBackup`). Qué tablas se conservan al restaurar es PREGUNTA abierta (`S6-respaldos.md` punto 2). WF-4 escribe aquí el procedimiento de restauración manual con el servicio detenido.

## 5. Revisión rápida de salud

```bash
curl --cacert "$(mkcert -CAROOT)/rootCA.pem" -sS https://finanzas.local:8443/health    # 200, JSON (getHealth, público)
ss -ltnp | grep 8443                                                                   # solo 192.168.50.56 y 127.0.0.1
tail -f datos/dev/run/finanzas.log                                                     # log del servidor (sin contraseñas, cookies ni tokens: S6-15)
```

## 6. Pendientes para WF-4

1. Unidad systemd del servidor y del respaldo (S6-09 dice que la unidad y el temporizador son de WF-4), usuario del servicio y rutas de producción.
2. Fijar la IP (reserva DHCP) y el nombre `finanzas.local` de forma permanente (avahi o DNS del módem); confirmar en los dos Android (`ops/mkcert.md` §4, manual de Axel).
3. ACL de lectura del token de Telegram para el usuario del servicio (K30) sin copiarlo.
4. Hora y retención de respaldos, tipos de adjunto y cabeceras: aprobación de Axel en la puerta 4 (K31).
5. Segunda copia de los respaldos en otro disco local.
