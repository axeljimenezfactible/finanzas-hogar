# Disparo manual de avisos por Telegram (S7-07)

Agente s7-avisos · 12-oct-2026. **Lo ejecuta el orquestador** (protocolo §7: ningún agente manda Telegram). El servicio lee el token en ejecución de `~/.openclaw/openclaw.json` (K30) con su ACL de lectura (WF4-01); el token no aparece en ningún comando de esta página ni se imprime.

## Qué hace

`sendTestNotification` (`POST /api/v1/notifications/{notification_id}/test`, cuerpo `{"person_id": <id>}`) arma **el texto real** del aviso con las cifras de hoy (motor/API), lo manda como texto plano sin `parse_mode` y responde **202** con el `NotificationDelivery`: `status` (`sent` o `failed`), `http_status` (el código HTTP de Telegram: S7-07 pide 200) y `error` resumido sin token. Si sale bien, marca `verified_at` del destinatario (Ajustes lo muestra como «registrado»).

## Requisitos

1. Servicio en marcha (`finanzas serve`) en `https://finanzas.local:8443`, con `avisos.Register` en `backend/cmd/finanzas/areas.go`.
2. ACL de lectura del usuario del servicio sobre `~/.openclaw/openclaw.json` (WF4-01). Si falta, el envío queda `failed` con «sin permiso de lectura … (K30)».
3. El chat de Axel registrado como destinatario. **PROVISIONAL**: la vía de detección del chat está en `docs/estado/bloqueos/S7-convivencia-multix.md`; hasta que Axel decida, el único camino es que el orquestador, con el OK de Axel, anote el chat de Axel como candidato y lo confirme por la API (paso 2).

## Pasos (desde la máquina de casa)

```sh
BASE=https://finanzas.local:8443
CA="$(mkcert -CAROOT)/rootCA.pem"
JAR=/home/deepsec/finanzas-hogar/datos/tmp/s7-07.cookies
mkdir -p /home/deepsec/finanzas-hogar/datos/tmp

# 1. Sesión de Axel: la contraseña se teclea, no se escribe en el comando ni en el historial.
printf 'Contraseña de axel: '; read -rs PW; echo
printf '{"username":"axel","password":"%s","device_label":"orquestador S7-07"}' "$PW" \
  | curl -sS --cacert "$CA" -c "$JAR" -H 'Content-Type: application/json' -d @- "$BASE/api/v1/auth/login" -o /dev/null -w 'login HTTP %{http_code}\n'
unset PW

# 2. (Solo si Axel aún no tiene destinatario y él lo aprueba) chat pendiente → confirmar.
curl -sS --cacert "$CA" -b "$JAR" "$BASE/api/v1/telegram-chat-candidates"            # id y chat enmascarado
curl -sS --cacert "$CA" -b "$JAR" -H 'Content-Type: application/json' \
  -H "Idempotency-Key: $(cat /proc/sys/kernel/random/uuid)" -d '{}' \
  "$BASE/api/v1/telegram-chat-candidates/<candidate_id>/confirm"

# 3. Ids de los avisos y de la persona de Axel.
curl -sS --cacert "$CA" -b "$JAR" "$BASE/api/v1/notifications" | jq -r '.items[] | "\(.id) \(.code)"'
AXEL_PERSON=$(curl -sS --cacert "$CA" -b "$JAR" "$BASE/api/v1/auth/me" | jq -r '.person_id')   # o el id de listPersons

# 4. Disparo: resumen diario (daily_summary) y una alerta de pago (payment_due).
for N in <id daily_summary> <id payment_due>; do
  curl -sS --cacert "$CA" -b "$JAR" -H 'Content-Type: application/json' \
    -d "{\"person_id\": $AXEL_PERSON}" "$BASE/api/v1/notifications/$N/test" \
    | jq '{notification_code, status, http_status, error}'
done
rm -f "$JAR"
```

Evidencia para `docs/estado/wf3/<orquestador>.md`, línea `S7-07`: por cada mensaje `status: "sent"` y `http_status: 200`, más la confirmación visual de Axel. Ximena queda **pendiente de acción** hasta que su chat se detecte y ella lo confirme (K4, S7-09).

La alerta de pago toma el siguiente pago programado de los próximos 31 días (`scheduled_payment`); si no hay, manda «No hay pagos programados en los próximos 31 días.» (sigue sirviendo para comprobar HTTP 200). Cada mensaje de prueba termina con «Mensaje de prueba enviado desde Ajustes.».

## Prueba contra el servidor falso (hecha por el agente)

`cd /home/deepsec/finanzas-hogar/backend && go test ./internal/api/avisos/ -run TestOp_sendTestNotification -count=1 -v`: los dos disparos (diario y pago) responden 202 con `status: sent` y `http_status: 200` contra un `httptest` con un `openclaw.json` falso; el cuerpo no lleva `parse_mode` y el token falso no aparece en el registro.

Nota: el paso 3 usa `getCurrentSession` (`GET /api/v1/auth/me`); si no trae `person_id`, usar el de `listPersons` (código `axel`).
