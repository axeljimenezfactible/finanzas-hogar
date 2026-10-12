# Bloqueo WF3-integrador-red · IP fija, nombre finanzas.local y CA en el sistema (PROVISIONAL)

12-oct-2026 · agente **integrador-backend** (WF-3a) · S6-11, S6-14, K29 · lo cierra WF-4 o Axel. No detiene WF-3b: el servidor de desarrollo ya corre en 192.168.50.56:8443 y 127.0.0.1:8443.

## Qué falta (PREGUNTA)

1. **IP de la LAN fija.** Hoy `wlp4s0` tiene 192.168.50.56 por DHCP. El certificado de mkcert y `FINANZAS_ADDR` la llevan escrita: si el módem la cambia, el servidor no arranca en la IP vieja y el certificado deja de valer por IP.
   - Opciones: (A) reserva DHCP en el módem para la MAC de `wlp4s0`; (B) IP estática en la máquina; (C) dejarla dinámica y reemitir el certificado cuando cambie.
   - PROPUESTA: **A** (no toca la máquina y sobrevive reinstalaciones).
2. **Nombre `finanzas.local` permanente.** Hoy lo publica un proceso `avahi-publish -a -R` que arranca `ops/dev-run.sh` (se pierde al reiniciar la máquina). En Android, Chrome no siempre resuelve `.local` (ya preguntado en `S6-sistema`).
   - Opciones: (A) alias fijo de avahi en `/etc/avahi/hosts` (requiere sudo, fuera de lo permitido en WF-3); (B) unidad systemd que corra `avahi-publish`; (C) entrada DNS en el módem `finanzas.local → 192.168.50.56` (sirve también a los teléfonos); (D) `/etc/hosts` (solo esta máquina, requiere sudo).
   - PROPUESTA: **C** si el módem lo permite (resuelve los dos Android) y **A** para la propia máquina; si el módem no lo permite, los teléfonos abren `https://192.168.50.56:8443/` (la IP está en el certificado).
3. **CA de mkcert en el almacén del sistema de esta máquina.** Se instaló solo en el almacén del usuario (`TRUST_STORES=nss mkcert -install`) porque el `mkcert -install` completo usa `sudo trust anchor`, y en WF-3 sudo solo está permitido para instalar mkcert (decisión 64).
   - Opciones: (A) dejarla así (las verificaciones usan `curl --cacert`); (B) instalarla también en el sistema (`mkcert -install`, pide sudo).
   - PROPUESTA: **A**; B solo si Axel quiere que programas de esta máquina confíen sin `--cacert`.

## Impacto

- 1 y 2 afectan a S6-11 (abrir `https://finanzas.local:8443/` sin advertencia) y a S6-14 (la dirección declarada en `ops/runbook.md`). Mientras no se respondan, 192.168.50.56 y el anuncio por `avahi-publish` quedan **PROVISIONAL** en `ops/runbook.md` §1–§2 y `ops/mkcert.md` §2.
- 3 no afecta a ningún criterio.
