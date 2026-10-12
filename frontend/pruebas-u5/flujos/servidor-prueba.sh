#!/usr/bin/env bash
# Servidor de PRUEBA aislado para u5-flujos (WF-3b): copia de la base de desarrollo en datos/tmp/u5-flujos/ y un
# «finanzas serve» en 127.0.0.1:8444 (mismo binario, mismo frontend compilado, HOME aislado → Telegram apagado).
# Uso: bash servidor-prueba.sh arrancar | detener | contar
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
[ -f "$ROOT/datos/finanzas-dev.db" ] || { echo "no existe la base de desarrollo"; exit 1; }
T="$ROOT/datos/tmp/u5-flujos"; PORT=8444
case "${1:-}" in
  arrancar)
    [ -f "$T/pid" ] && kill "$(cat "$T/pid")" 2>/dev/null && sleep 1 || true
    rm -rf "$T"; mkdir -p "$T/app" "$T/respaldos" "$T/home"; chmod 700 "$T"
    sqlite3 "$ROOT/datos/finanzas-dev.db" ".backup '$T/finanzas.db'"
    cd "$ROOT"
    HOME="$T/home" FINANZAS_DB="$T/finanzas.db" FINANZAS_ADDR="127.0.0.1:$PORT" \
    FINANZAS_CERT="$ROOT/datos/certs/finanzas.local.pem" FINANZAS_KEY="$ROOT/datos/certs/finanzas.local-key.pem" FINANZAS_CA="$(mkcert -CAROOT)/rootCA.pem" \
    FINANZAS_FRONTEND="$ROOT/frontend/dist" FINANZAS_OPENAPI="$ROOT/docs/contratos/openapi.yaml" \
    FINANZAS_DATOS="$T/app" FINANZAS_RESPALDOS="$T/respaldos" \
      setsid nohup "$ROOT/backend/bin/finanzas" serve > "$T/serve.log" 2>&1 < /dev/null &
    echo $! > "$T/pid"
    for _ in $(seq 1 40); do curl -s -m 2 --cacert "$(mkcert -CAROOT)/rootCA.pem" "https://127.0.0.1:$PORT/health" && break; sleep 0.5; done; echo ;;
  detener) [ -f "$T/pid" ] && kill "$(cat "$T/pid")" 2>/dev/null || true; sleep 1 ;;
  contar) sqlite3 "${2:-$T/finanzas.db}" "select 'expense',count(*) from expense union all select 'weight_entry',count(*) from weight_entry;" ;;
  *) echo "uso: arrancar|detener|contar"; exit 1 ;;
esac
