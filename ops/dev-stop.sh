#!/usr/bin/env bash
# ops/dev-stop.sh · Finanzas Hogar · detiene el servidor de DESARROLLO que arrancó ops/dev-run.sh
# (WF-3a, agente integrador-backend, 12-oct-2026).
#
# Uso:  bash ops/dev-stop.sh              detiene «finanzas serve» (SIGTERM; a los 10 s, SIGKILL) y avahi-publish
#       bash ops/dev-stop.sh --conservar-nombre   detiene solo el servidor; finanzas.local sigue publicado
# No borra la base (datos/finanzas-dev.db) ni datos/dev/: para empezar de cero, bórralos a mano después.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN="$ROOT/datos/dev/run"
PIDF="$RUN/finanzas.pid"
AVAHI_PIDF="$RUN/avahi-publish.pid"
PORT="${FINANZAS_PORT:-8443}"

conservar=0
case "${1:-}" in
  "") ;;
  --conservar-nombre) conservar=1 ;;
  *) echo "dev-stop: opción desconocida: $1" >&2; exit 1 ;;
esac

detener() { # detener <archivo-pid> <nombre>
  local f="$1" nombre="$2" pid
  if [ ! -f "$f" ]; then echo "dev-stop: $nombre no tiene archivo de PID ($f)"; return 0; fi
  pid="$(cat "$f")"
  if ! kill -0 "$pid" 2>/dev/null; then echo "dev-stop: $nombre (PID $pid) ya no corría"; rm -f "$f"; return 0; fi
  kill -TERM "$pid"
  for _ in $(seq 1 40); do kill -0 "$pid" 2>/dev/null || break; sleep 0.25; done
  if kill -0 "$pid" 2>/dev/null; then kill -KILL "$pid"; echo "dev-stop: $nombre (PID $pid) no respondió a SIGTERM: SIGKILL"; fi
  rm -f "$f"
  echo "dev-stop: $nombre detenido (PID $pid)"
}

detener "$PIDF" "finanzas serve"
[ "$conservar" = 1 ] || detener "$AVAHI_PIDF" "avahi-publish (finanzas.local)"

if ss -ltn "sport = :$PORT" | grep -q LISTEN; then
  echo "dev-stop: AVISO: algo sigue escuchando en :$PORT:"; ss -ltnp "sport = :$PORT"; exit 1
fi
echo "dev-stop: nada escucha en :$PORT"
