# HTTPS de la casa con mkcert (S6-11)

Versión 1.1 · 12-oct-2026 · v1.0 del agente **auth-sistema**; v1.1 del **integrador-backend** (WF-3a), que instaló mkcert, creó la CA, emitió el certificado y corrió S6-11. La prueba en los dos Android es **manual de Axel** (S6-11: mientras no la haga, queda «no verificado»).

**Cambios de la v1.1:** carpeta `datos/certs/` (antes `datos/tls/`; ignorada por git); IP de la LAN `192.168.50.56` (PROVISIONAL hasta WF-4); `mkcert -install` solo en el almacén del usuario (`TRUST_STORES=nss`, sin sudo); el certificado incluye también `localhost`; el arranque de desarrollo es `bash ops/dev-run.sh`.

## 1. Instalar mkcert y crear la CA local (una vez, en la máquina de la casa) · HECHO 12-oct-2026

```bash
sudo pacman -S --noconfirm mkcert      # única orden con sudo permitida (decisión 64 y G) → mkcert 1.4.4-3 (extra)
mkcert -version                        # v1.4.4 (plan §1.8)
TRUST_STORES=nss mkcert -install       # crea la CA local del usuario deepsec y la registra solo en ~/.pki/nssdb (Chrome/Firefox del usuario)
mkcert -CAROOT                         # /home/deepsec/.local/share/mkcert: rootCA.pem (público) y rootCA-key.pem (SECRETO, modo 400)
```

- **Por qué `TRUST_STORES=nss`:** un `mkcert -install` a secas también instala la CA en el almacén del sistema de Arch con `sudo trust anchor`, y ese sudo **no** está permitido (decisión 64: sudo solo para instalar mkcert). Con `nss` la CA queda en el almacén del usuario y no hace falta sudo. Las verificaciones usan `curl --cacert "$(mkcert -CAROOT)/rootCA.pem"`, así que no dependen del almacén del sistema.
- PREGUNTA para Axel (no bloquea): ¿quiere la CA también en el almacén del sistema de esta máquina (`mkcert -install` sin `TRUST_STORES`, pide sudo)? Solo sirve para que `curl` y otros programas de la propia máquina confíen sin `--cacert`.
- `rootCA-key.pem` nunca se copia, no sale de la máquina y no entra al repositorio. El servidor solo ofrece archivos con bloques `CERTIFICATE` (si `FINANZAS_CA` apunta a una llave, responde 503).

## 2. Emitir el certificado del servidor · HECHO 12-oct-2026

La IP de la LAN es **192.168.50.56** (la de `wlp4s0` hoy); queda **PROVISIONAL** hasta que WF-4 la fije en `ops/runbook.md` (S6-14). Si cambia, se vuelve a emitir con la IP nueva.

```bash
cd /home/deepsec/finanzas-hogar
mkdir -p datos/certs && chmod 700 datos/certs          # datos/certs/ está en .gitignore
mkcert -cert-file datos/certs/finanzas.local.pem -key-file datos/certs/finanzas.local-key.pem \
       finanzas.local 192.168.50.56 127.0.0.1 localhost
chmod 600 datos/certs/finanzas.local-key.pem
openssl x509 -in datos/certs/finanzas.local.pem -noout -ext subjectAltName -enddate
#   DNS:finanzas.local, DNS:localhost, IP Address:192.168.50.56, IP Address:127.0.0.1 · vence el 12-ene-2029
```

## 3. Arrancar el servidor

Desarrollo (WF-3): `bash ops/dev-run.sh` (base `datos/finanzas-dev.db`, Telegram apagado) y `bash ops/dev-stop.sh`; detalle en `ops/runbook.md` §3. A mano, lo mismo que hace el script:

```bash
FINANZAS_ADDR=192.168.50.56:8443,127.0.0.1:8443 \
FINANZAS_CERT=/home/deepsec/finanzas-hogar/datos/certs/finanzas.local.pem \
FINANZAS_KEY=/home/deepsec/finanzas-hogar/datos/certs/finanzas.local-key.pem \
FINANZAS_CA="$(mkcert -CAROOT)/rootCA.pem" \
backend/bin/finanzas serve
```

- `serve` rechaza `0.0.0.0`, `::` y `:8443` (comodines); se le dan solo `127.0.0.1` y la IP de la LAN (S6-14, K29). La red de Tailscale (`tailscale0`, 100.x) no se usa.
- Sin `--cert/--key` no arranca; `--dev-http` (sin TLS) solo se permite en `127.0.0.1` para desarrollo.

## 4. Instalar la CA en los dos Android (Axel y Ximena) · PENDIENTE DE AXEL

1. Bajar el certificado: en la app, **Ajustes › Seguridad › Descargar certificado de la casa** (`GET /api/v1/system/ca-certificate`, archivo `finanzas-hogar-ca.pem`), o copiar `rootCA.pem` por cable USB. La primera vez, como el teléfono aún no confía, se baja por USB o aceptando la advertencia una sola vez.
2. Abrir **Configuración › Seguridad y privacidad › Más ajustes de seguridad › Cifrado y credenciales › Instalar un certificado › Certificado de CA** (los nombres cambian un poco según la marca).
3. Elegir `finanzas-hogar-ca.pem`, confirmar con el PIN del teléfono.
4. Comprobar en **Credenciales de confianza › Usuario** que aparece «mkcert …».
5. En Chrome abrir `https://finanzas.local:8443/`: debe cargar **sin advertencia**; luego «Instalar app» (PWA).

PREGUNTA (bloqueo `S6-sistema`): Android no siempre resuelve nombres `.local` (mDNS); si `finanzas.local` no abre en el teléfono, hace falta una entrada en el DNS del módem o router (finanzas.local → 192.168.50.56). Mientras tanto, `https://192.168.50.56:8443/` también es válido para el certificado (la IP está en el SAN).

## 5. Verificación (S6-11; la corrió el integrador el 12-oct-2026)

```bash
make verify-wf3-servidor     # S6-11, S6-12 y S6-14 contra el servidor corriendo
# o a mano, los comandos literales de S6-11 y S6-14:
curl --cacert "$(mkcert -CAROOT)/rootCA.pem" -sS -o /dev/null -w "%{http_code}\n" https://finanzas.local:8443/     # 200
openssl s_client -connect finanzas.local:8443 -servername finanzas.local </dev/null 2>/dev/null | openssl x509 -noout -ext subjectAltName
ss -ltnp | grep 8443        # solo 192.168.50.56:8443 y 127.0.0.1:8443
```

El certificado de mkcert vence a los 2 años y 3 meses (12-ene-2029): se vuelve a emitir con el paso 2 (la CA no cambia, los teléfonos no se tocan).
