# Bloqueo criterios-token-servicio
Contexto: S7-06: protocolo §7 y decisiones F: el token del bot MultiX se lee en ejecución desde `~/.openclaw/openclaw.json`. Plan §6 y decisiones H: el servicio corre con usuario propio (`finanzas`).
Qué falta: ¿Cómo llega el token al servicio si el archivo está en el home de otro usuario (`deepsec`) y el servicio corre como `finanzas`? Un usuario propio no puede leer ese archivo sin cambiar permisos.
Opciones:
- (a) El servicio corre como `deepsec` (sin usuario propio).
- (b) ACL de lectura del archivo para `finanzas` (sigue leyéndose del mismo archivo en ejecución).
- (c) Archivo de entorno de systemd (`EnvironmentFile`) fuera del repositorio con permisos 600, creado por Axel (copia del token fuera de `openclaw.json`).
Impacto: El envío real de Telegram desde el servicio (S7-07) y la unidad systemd (WF4-01) dependen de esta decisión.
Cómo seguí: S7-06 exige que el token jamás se imprima ni entre al repositorio y que las pruebas usen archivo falso.
