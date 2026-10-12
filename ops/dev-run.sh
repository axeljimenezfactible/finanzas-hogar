#!/usr/bin/env bash
# ops/dev-run.sh · Finanzas Hogar · servidor de DESARROLLO (WF-3a, agente integrador-backend, 12-oct-2026)
#
# Arranca «finanzas serve» en segundo plano con la base sembrada datos/finanzas-dev.db (ignorada por git),
# HTTPS con el certificado de mkcert (datos/certs/, ignorado por git) en 192.168.50.56:8443 y 127.0.0.1:8443.
# Lo detiene ops/dev-stop.sh. Producción (systemd, rutas, hora de respaldos) es de WF-4: ops/runbook.md.
#
# Uso:  bash ops/dev-run.sh             arranca (si la base no existe: migra y siembra)
#       bash ops/dev-run.sh --sembrar   además vuelve a correr «finanzas migrate» y «finanzas semilla» (idempotente)
#
# Seguridad:
# - Telegram APAGADO en desarrollo: el proceso corre con HOME=datos/dev/home, donde no existe .openclaw/openclaw.json;
#   los jobs de avisos (daily_summary, weekly_summary, alerts, reporte mensual) no encuentran token y no envían nada.
#   Los envíos reales los dispara solo el orquestador (S7-07).
# - Contraseñas iniciales de la base de desarrollo: archivo FUERA del repositorio, modo 600
#   (~/.config/finanzas-hogar/dev-semilla.txt, renglones axel=… y ximena=…). Si no existe se crea con valores
#   aleatorios; este script nunca las imprime. No son las contraseñas reales de la casa.
# - rootCA-key.pem de mkcert nunca se toca: el servidor solo recibe rootCA.pem (FINANZAS_CA).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LAN_IP="${FINANZAS_LAN_IP:-192.168.50.56}"   # PROVISIONAL hasta que WF-4 la fije en ops/runbook.md (S6-14)
PORT="${FINANZAS_PORT:-8443}"
DB="${FINANZAS_DEV_DB:-$ROOT/datos/finanzas-dev.db}"
DEV="$ROOT/datos/dev"                         # adjuntos, reportes, respaldos, HOME aislado y PID (ignorado por git)
RUN="$DEV/run"
PIDF="$RUN/finanzas.pid"
AVAHI_PIDF="$RUN/avahi-publish.pid"
LOG="$RUN/finanzas.log"
CERT="$ROOT/datos/certs/finanzas.local.pem"
KEY="$ROOT/datos/certs/finanzas.local-key.pem"
BIN="$ROOT/backend/bin/finanzas"
PASSFILE="${FINANZAS_SEED_PASSWORD_FILE:-$HOME/.config/finanzas-hogar/dev-semilla.txt}"

fail() { echo "dev-run: $*" >&2; exit 1; }

sembrar=0
for a in "$@"; do
  case "$a" in
    --sembrar) sembrar=1 ;;
    -h|--help) sed -n '2,19p' "$0"; exit 0 ;;
    *) fail "opción desconocida: $a (usa --sembrar o --help)" ;;
  esac
done

if [ -f "$PIDF" ] && kill -0 "$(cat "$PIDF")" 2>/dev/null; then
  echo "dev-run: el servidor ya está corriendo (PID $(cat "$PIDF")); para detenerlo: bash ops/dev-stop.sh"
  exit 0
fi

command -v mkcert >/dev/null || fail "falta mkcert (ops/mkcert.md §1)"
CA="$(mkcert -CAROOT)/rootCA.pem"
[ -f "$CA" ] || fail "no existe $CA: corre «TRUST_STORES=nss mkcert -install» (ops/mkcert.md §1)"
[ -f "$CERT" ] && [ -f "$KEY" ] || fail "faltan $CERT o $KEY (ops/mkcert.md §2)"
ip -br addr | grep -qw "$LAN_IP/[0-9]*" || fail "la IP $LAN_IP no está en ninguna interfaz (ip -br addr); exporta FINANZAS_LAN_IP"
if ss -ltn "sport = :$PORT" | grep -q LISTEN; then fail "el puerto $PORT ya está ocupado: $(ss -ltnp "sport = :$PORT" | tail -n +2)"; fi

mkdir -p "$RUN" "$DEV/home" "$DEV/app" "$DEV/respaldos"
chmod 700 "$DEV" "$DEV/home"

echo "dev-run: compilando backend/bin/finanzas"
(cd "$ROOT/backend" && go build -o bin/finanzas ./cmd/finanzas)

if [ ! -f "$DB" ] || [ "$sembrar" = 1 ]; then
  if [ ! -f "$PASSFILE" ]; then
    mkdir -p "$(dirname "$PASSFILE")"; chmod 700 "$(dirname "$PASSFILE")"
    ( umask 077
      printf '# Contraseñas iniciales de la base de DESARROLLO (ops/dev-run.sh). No son las de la casa.\naxel=%s\nximena=%s\n' \
        "$(openssl rand -base64 18)" "$(openssl rand -base64 18)" > "$PASSFILE" )
    echo "dev-run: contraseñas iniciales de desarrollo creadas en $PASSFILE (modo 600; no se imprimen)"
  fi
  echo "dev-run: migrando y sembrando $DB"
  (cd "$ROOT" && "$BIN" migrate -db "$DB")
  (cd "$ROOT" && FINANZAS_SEED_PASSWORD_FILE="$PASSFILE" "$BIN" semilla -db "$DB" -dir "$ROOT/datos/semilla") > "$RUN/semilla.log" 2>&1 \
    || { tail -20 "$RUN/semilla.log" >&2; fail "falló «finanzas semilla» (detalle en $RUN/semilla.log)"; }
  echo "dev-run: semilla OK (resumen en $RUN/semilla.log)"
fi

# Nombre finanzas.local en esta máquina: mDNS con avahi (nsswitch: mdns_minimal). Si ya resuelve, no se publica otra vez.
if ! getent hosts finanzas.local | grep -qw "$LAN_IP"; then
  command -v avahi-publish >/dev/null || fail "finanzas.local no resuelve y falta avahi-publish (ops/runbook.md §2)"
  setsid nohup avahi-publish -a -R finanzas.local "$LAN_IP" > "$RUN/avahi-publish.log" 2>&1 < /dev/null &
  echo $! > "$AVAHI_PIDF"
  for _ in $(seq 1 20); do getent hosts finanzas.local | grep -qw "$LAN_IP" && break; sleep 0.25; done
  getent hosts finanzas.local | grep -qw "$LAN_IP" || fail "avahi-publish no logró publicar finanzas.local (ver $RUN/avahi-publish.log)"
  echo "dev-run: finanzas.local → $LAN_IP publicado con avahi-publish (PID $(cat "$AVAHI_PIDF"))"
fi

cd "$ROOT"
HOME="$DEV/home" \
FINANZAS_DB="$DB" \
FINANZAS_ADDR="$LAN_IP:$PORT,127.0.0.1:$PORT" \
FINANZAS_CERT="$CERT" FINANZAS_KEY="$KEY" FINANZAS_CA="$CA" \
FINANZAS_FRONTEND="$ROOT/frontend/dist" FINANZAS_OPENAPI="$ROOT/docs/contratos/openapi.yaml" \
FINANZAS_DATOS="$DEV/app" FINANZAS_RESPALDOS="$DEV/respaldos" \
  setsid nohup "$BIN" serve > "$LOG" 2>&1 < /dev/null &
echo $! > "$PIDF"

for _ in $(seq 1 60); do
  kill -0 "$(cat "$PIDF")" 2>/dev/null || { tail -20 "$LOG" >&2; rm -f "$PIDF"; fail "el servidor terminó al arrancar (log: $LOG)"; }
  n=$(ss -ltn "sport = :$PORT" | grep -c LISTEN || true)
  [ "$n" -ge 2 ] && break
  sleep 0.25
done
[ "${n:-0}" -ge 2 ] || fail "el servidor no abrió $LAN_IP:$PORT y 127.0.0.1:$PORT en 15 s (log: $LOG)"

echo "dev-run: servidor corriendo · PID $(cat "$PIDF") · log $LOG"
echo "         https://finanzas.local:$PORT/  https://$LAN_IP:$PORT/  https://127.0.0.1:$PORT/"
echo "         base $DB · Telegram apagado (HOME aislado) · detener: bash ops/dev-stop.sh"
