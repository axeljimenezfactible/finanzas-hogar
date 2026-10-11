# Bloqueo criterios-bind
Contexto: S6-14: decisiones H dice «solo red local». La máquina de casa tiene también interfaz Tailscale (100.115.28.35, plan §1.4).
Qué falta: ¿El servicio debe escuchar solo en la LAN (192.168.50.x) o también por Tailscale (acceso remoto de Axel y Ximena fuera de casa)?
Opciones:
- (a) Solo LAN (interpretación literal de «solo red local»).
- (b) LAN + Tailscale (más cómodo; amplía la superficie).
- (c) LAN ahora y Tailscale como decisión posterior.
Impacto: La dirección de escucha en S6, WF-4 y el certificado de mkcert (¿incluye el nombre de Tailscale?) quedan sin fijar.
Cómo seguí: S6-14 exige declarar la dirección y no abrir Tailscale mientras no haya decisión.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción (a): el servicio escucha solo en la LAN (192.168.50.x) con el nombre `finanzas.local`. Tailscale no se abre.
- **Estado:** RESUELTO. Se aplica en: criterios S6-14, S6, WF-4 (mkcert solo para finanzas.local).
