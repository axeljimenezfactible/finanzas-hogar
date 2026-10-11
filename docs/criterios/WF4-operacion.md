# WF-4 · Operación, semilla, seguridad e informe — Criterios de aceptación

Versión 1.1 (11-oct-2026, puerta 1) · Escrito en WF-0 por el agente «criterios» (ronda 1 de WF-0: `INDICE.md` §9); actualizado en el flujo P1 con las respuestas de la puerta 1 (`docs/plan/decisiones.md` §K) y los hallazgos bajos de la revisión r3 (agente `criterios-S-1.1`) · El revisor final evalúa **solo** contra este archivo y `S8.md`. Convención de evidencia, severidades y comandos: ver `INDICE.md`.

**Cambios de la versión 1.1 (11-oct-2026):**
- La etiqueta «1.1» que este archivo llevaba desde la ronda 1 de WF-0 (10-oct-2026) se reemplaza por la de la puerta 1; el contenido de esa ronda se conserva.
- Toda ruta temporal pasa de `/tmp` a `/home/deepsec/finanzas-hogar/datos/tmp/` (R2-07; los temporales se borran al terminar).
- WF4-01: el usuario del servicio lee el token por ACL sobre `~/.openclaw/openclaw.json`; el token nunca se copia (K30); se cierra el bloqueo `criterios-token-servicio`.
- WF4-03: el certificado cubre solo `finanzas.local`; Tailscale no se abre (K29).
- WF4-04 y WF4-08: retención, hora y política de permisos propuestas por el agente y aprobadas en la puerta 4 (K31).
- WF4-05: el procedimiento del chat de Ximena pasa a «confirmar en Ajustes tras escribirle al bot» (K4).
- §5: trazabilidad con K4, K29, K30 y K31.

## 1. Objetivo
Dejar el procedimiento verificable de los 5 agentes de WF-4 (ops: systemd, mkcert, respaldos y runbook; semilla; seguridad; informe final; revisor final) y la verificación extremo a extremo de plan.md §9, de modo que el cierre de la Etapa B se apruebe o devuelva solo con evidencia.

## 2. Criterios

**WF4-01 (alta) · Unidad systemd propia.** `ops/finanzas-hogar.service` define un usuario propio (`finanzas`), `WantedBy=multi-user.target` (arranca con la máquina) y ejecuta el binario de `backend/`; está instalada y activa. El usuario del servicio puede leer la base, los adjuntos y la carpeta de respaldos y, por **ACL de lectura** (K30), el archivo `~/.openclaw/openclaw.json` donde está el token del bot; no puede leer nada más fuera de lo necesario. El runbook trae el comando `setfacl` (con la ruta, nunca con el token). El token no se copia al repositorio, a la unidad systemd ni a ningún archivo.
- Verifica: `systemd-analyze verify /home/deepsec/finanzas-hogar/ops/finanzas-hogar.service`; `systemctl show finanzas-hogar -p User,ActiveState,UnitFileState`; `systemctl status finanzas-hogar --no-pager`; `getfacl ~/.openclaw/openclaw.json` muestra una entrada de lectura para `finanzas` (no se imprime el contenido); `sudo -u finanzas test -r ~/.openclaw/openclaw.json && echo legible` (sin imprimir el contenido); `namei -l ~/.openclaw/openclaw.json` confirma que `finanzas` puede recorrer cada carpeta del camino; `grep -rn "botToken" /home/deepsec/finanzas-hogar/ops` muestra a lo sumo el nombre de la clave, ningún valor. Esperado: `User=finanzas`, `UnitFileState=enabled`, `ActiveState=active`.
- Evidencia: `docs/estado/wf4/<ops>.md` línea `WF4-01`.

**WF4-02 (alta) · Vuelve solo tras reinicio (real).** Tras reiniciar la máquina, el servicio está activo sin intervención y la dirección `https://finanzas.local:8443` responde. El reinicio lo ejecuta o autoriza Axel/el orquestador; sin él, el criterio queda «no verificado».
- Verifica: `systemctl is-enabled finanzas-hogar` y, tras el arranque, `systemctl is-active finanzas-hogar` + `curl --cacert "$(mkcert -CAROOT)/rootCA.pem" -sS -o /dev/null -w "%{http_code}\n" https://finanzas.local:8443/`.
- Evidencia: línea `WF4-02` con hora de arranque y salidas; o «no verificado: reinicio pendiente de Axel».

**WF4-03 (alta) · mkcert y certificado.** mkcert está instalado (1.4.4 del plan §1.8; se anota en «Paquetes instalados»), el certificado cubre solo `finanzas.local` (K29: solo LAN; Tailscale no se abre) y se documenta su fecha de vencimiento y la renovación; `ops/mkcert.md` trae los pasos para instalar la CA en los dos Android. La resolución mDNS de `finanzas.local` funciona en la LAN (avahi activo) y el runbook sugiere reservar la IP en el router (riesgo R2).
- Verifica: `mkcert -version`; `openssl x509 -in <cert> -noout -enddate -ext subjectAltName`; `getent hosts finanzas.local` o `avahi-resolve -n finanzas.local`.
- Evidencia: línea `WF4-03`; la instalación de la CA en los Android es acción de Axel (plan §8.2) y se marca «pendiente de Axel» hasta su confirmación.

**WF4-04 (alta) · Respaldo diario programado y restaurado.** `ops/respaldo.sh` hace el respaldo con `VACUUM INTO` en una ubicación distinta de la base (otra carpeta o disco), versionado por fecha, ejecutado a diario por un temporizador systemd; verifica con `PRAGMA integrity_check`; se **restaura** en una ruta temporal con `integrity_check` = `ok` y mismos conteos de filas. La retención (cuántas versiones) y la hora no están en el plan: el agente las propone con su justificación en el runbook y quedan PROVISIONAL, con la propuesta escrita, hasta que Axel las apruebe en la puerta 4 (K31).
- Verifica: `systemctl list-timers --all | grep -i respaldo`; `bash /home/deepsec/finanzas-hogar/ops/respaldo.sh`; `ls -l <carpeta de respaldos>`; `sqlite3 <copia restaurada> "PRAGMA integrity_check;"`.
- Evidencia: línea `WF4-04` con salidas.

**WF4-05 (alta) · Runbook utilizable.** `ops/runbook.md` cubre como mínimo los 6 procedimientos que el plan menciona: arranque, parada y estado del servicio; restauración de respaldo paso a paso (R1, decisiones H); instalación de la CA de mkcert en los dos Android (R3, `ops/mkcert.md`); confirmación en Ajustes del chat de Ximena después de que ella le escriba al bot MultiX (R15, K4); reserva de IP en el router y nombre mDNS (R2); cambio del interruptor A/B cuando Kueski confirme (plan §8.2). Cada procedimiento fue ejecutado una vez por el agente (o marcado «no ejecutado» con motivo).
- Verifica: lista de los 6 procedimientos con su estado «ejecutado / no ejecutado» y la salida de cada uno; revisión por el revisor de 3 procedimientos al azar.
- Evidencia: línea `WF4-05`.

**WF4-06 (alta) · Semilla con importador listo.** La carga de la semilla toma `datos/semilla/data.json` y `oraculo/model.json`. El importador de `datos/semilla/gastos_fijos_ingresos.csv` (columnas `tipo, concepto, monto, periodicidad, dia_o_fecha, categoria, persona, porcentaje_conjunto, notas`; valores de `tipo` ingreso/gasto_fijo/gasto_venidero; de `periodicidad` quincenal/mensual/anual/unica; de `persona` axel/ximena/hogar/karen/conjunto, plan §8.2) existe y funciona: **sin el archivo** la carga termina bien y deja el importador listo (decisiones H); **con el archivo** valida cada fila, rechaza las inválidas indicando número de línea y no inventa filas. Las filas de prueba del importador viven en fixtures de prueba, nunca en la semilla.
- Verifica: ejecutar la carga con y sin el CSV (si Axel no lo ha dejado, con un CSV de fixture en `/home/deepsec/finanzas-hogar/datos/tmp/`): sin CSV → código 0 y mensaje claro; CSV con 1 fila inválida → esa línea se rechaza y las demás entran; `sqlite3 <base> "SELECT COUNT(*) FROM <tabla_ingreso>"` coherente.
- Evidencia: línea `WF4-06`.

**WF4-07 (alta) · Semilla = oráculo.** Tras la carga, `GET` a la API reproduce el 100 % de las 85 cifras canon y las 596 hojas de `model.json` (S8-04) y la comparativa de tiendas de `data.json` (K-01…K-08). Sin datos históricos más allá de la semilla.
- Verifica: ver S8-04; más `sqlite3 <base> "SELECT COUNT(*) FROM <tabla_compra_credito>"` = 4 despensas y `<tabla_pago_programado>` = 16 pagos de las despensas.
- Evidencia: línea `WF4-07`.

**WF4-08 (alta) · Informe de seguridad.** `docs/estado/wf4/<agente-seguridad>.md` incluye: resultado de `govulncheck` y `npm audit` (S8-01); tabla de cabeceras (S8-02); prueba de autenticación (login, 5 intentos en 15 min, cierre de sesión, Karen sin cuenta); barrido de secretos del árbol de trabajo (token de Telegram, llaves, contraseñas); permisos de archivos de la base, los adjuntos y los respaldos (propietario `finanzas`, sin lectura para otros usuarios; la política exacta la propone el agente en el runbook y queda PROVISIONAL hasta que Axel la apruebe en la puerta 4, K31); lista de puertos escuchando.
- Verifica: `grep -rn "[0-9]\{8,10\}:[A-Za-z0-9_-]\{30,\}" /home/deepsec/finanzas-hogar -I | wc -l` = 0; `ls -l` de la base, adjuntos y respaldos; `ss -ltnp | grep -E "8443|finanzas"`.
- Evidencia: línea `WF4-08` con cada salida.

**WF4-09 (alta) · Informe final completo.** `docs/estado/informe-final.md` cumple S8-05: hecho, verificado con evidencia, consumo real de tokens por fase contra el presupuesto de plan.md §5 (WF-0 850K, WF-1 3,000K, WF-3 3,480K, Pista 3 ≈ 200K + 300K de sesión principal, WF-4 700K, total ≈ 8.2M) o marcado «no medido» donde no haya registro, paquetes instalados, bloqueos y estado, lista de lo no verificado y tabla de extremo a extremo (S8-06). El texto está en español de México y las cifras salen de `docs/estado/` o del oráculo.
- Verifica: lectura del revisor con comprobación cruzada de 10 afirmaciones.
- Evidencia: nota del revisor.

**WF4-10 (alta) · Revisor final independiente.** El revisor final es un agente distinto de los que construyeron; evalúa **solo** contra `docs/criterios/`; vuelve a **ejecutar** al menos los comandos de S8-01, S8-02, S8-04 y WF4-04 en lugar de confiar en los reportes; aprueba o devuelve con lista de criterios fallidos por ID; máximo **3 rondas**, y si no aprueba escala a Axel.
- Verifica: reporte del revisor con el número de ronda, los IDs evaluados (todos los de `S8.md` y `WF4-operacion.md`), las salidas que reejecutó y el veredicto.
- Evidencia: `docs/estado/wf4/<revisor-final>.md`.

**WF4-11 (media) · Aviso final por Telegram del orquestador.** Al cerrar WF-4 y la Etapa B el **orquestador** (no un agente) envía a Axel un mensaje de texto plano con prefijo `[Finanzas-Hogar]` y formato Qué / Archivos / Verificación / Estado / Pendientes; HTTP 200 registrado.
- Verifica: respuesta de la API de Telegram con `ok: true` registrada por el orquestador (sin imprimir el token).
- Evidencia: `docs/estado/wf4/<orquestador>.md` línea `WF4-11`.

**WF4-12 (media) · Instalaciones declaradas.** Todo paquete instalado en WF-4 (p. ej. `mkcert`, `govulncheck`) figura en «Paquetes instalados» del reporte del agente con nombre y versión.
- Verifica: comparar `pacman -Qqe` o el historial de instalaciones del día contra los reportes; 0 paquetes sin declarar.
- Evidencia: línea `WF4-12`.

## 3. No aceptable
- Reiniciar la máquina de casa sin autorización; modificar el router, la CA de los Android o cualquier cosa fuera de `~/finanzas-hogar/` salvo paquetes y el servicio.
- Declarar verificado lo que depende de Axel (instalar la CA, Android, Ximena en el bot, reinicio) sin su confirmación.
- Respaldo sin restauración probada; respaldo en la misma carpeta que la base.
- Semilla con filas inventadas para «probar» o con datos históricos; importador que falla sin el CSV de Axel.
- Revisor que reaprueba sin reejecutar comandos o que evalúa contra criterios que no están en este directorio.
- Token de Telegram impreso en cualquier salida; mensajes de aviso con `parse_mode`.
- Informe con consumos de tokens estimados presentados como medidos.

## 4. Umbrales numéricos (del plan)
- Rondas del revisor: ≤ 3 (WF4-10).
- 0 vulnerabilidades críticas (WF4-08).
- 4 despensas y 16 pagos de la semilla; 596 hojas, 85 cifras, 0 diferencias (WF4-07).
- 6 procedimientos mínimos del runbook (WF4-05).
- Puerto 8443; `finanzas.local`; mkcert 1.4.4.
- Presupuesto WF-4: 700K (5 agentes); tope por agente 300K.

## 5. Trazabilidad
- plan.md §3.4 fila S8 y §5 fila WF-4 (ops, semilla, seguridad, informe, revisor final) → WF4-01…WF4-12.
- plan.md §9 puntos 3 a 7 → WF4-02, WF4-03, WF4-04, WF4-09.
- plan.md §3.1 principio 5 (revisor independiente, máximo 3 rondas) y principio 7 (Telegram al cerrar cada flujo) → WF4-10, WF4-11.
- plan.md §3.1 principio 9 y §8.2 (semilla CSV de Axel; importador listo) → WF4-06.
- decisiones K «Token del bot: el usuario del servicio recibe permiso de lectura (ACL) sobre ~/.openclaw/openclaw.json y lee bots.multix.botToken en tiempo de ejecución; el token jamás se copia, imprime ni entra al repositorio» → WF4-01.
- decisiones K «Red: solo LAN; Tailscale no se abre» → WF4-03.
- decisiones K «cada agente propone el valor con justificación; Axel aprueba en la puerta que toque (tiempo máximo del optimizador en la puerta 2; reportes, respaldos, cabeceras, adjuntos, PDF y permisos en la puerta 4). Hasta entonces, PROVISIONAL con la propuesta escrita» → WF4-04, WF4-08.
- decisiones K «Chat de Telegram de Ximena: el bot lo detecta cuando ella escribe y ella confirma en Ajustes» → WF4-05.
- decisiones H «Go + SQLite (un archivo) + servicio systemd nativo (usuario propio, arranca con la máquina). Respaldos solo locales… restauración probada» → WF4-01, WF4-02, WF4-04.
- decisiones H «Instalar paquetes del sistema: autorizado… se reporta lo instalado» → WF4-12.
- decisiones F «Avisos de desarrollo del orquestador con prefijo `[Finanzas-Hogar]`» → WF4-11.
- plan.md §7 riesgos R1, R2, R3, R15, R17 → WF4-04, WF4-03, WF4-03, WF4-05, WF4-06.
