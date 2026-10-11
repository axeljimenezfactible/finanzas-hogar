# correccion-criterios-r2 · Corrección de los hallazgos de REVISION-criterios-r2 (11-oct-2026)

## Hecho
- **Hallazgo 1 (media, P3-02):** la evidencia pasa a «línea `P3-02` en el reporte del agente `paquete-claude-design` y en el de su revisor, en `docs/estado/wf1/`, con `ls -l` y `sha256sum -c`», igual que P3-01; se aclara que `P3-iter-1.md` no es evidencia (lo escribe el agente Sonnet después del aviso de Axel, K35).
- **Hallazgo 2 (media, U4-12):** la lista mínima queda enumerada (14 parámetros: umbral 80 %, umbral 100 %, días antes del pago [1, 0], `send_time`, `advance_days`, precio viejo 30, peso cero 90, sesión 90, 5 intentos, ventana 15 min, comparador 5 s, adjunto 10 MB, ajuste kcal bajar 0 y subir 0). Para `send_time` y `advance_days` la celda «valor inicial» dice «sin valor inicial: lo propone S7 (K31), PROVISIONAL» y no cuenta como vacía; no se inventa valor. Ajuste análogo en la Verifica de avisos de C1-22 (su «restablecer» no exige valor sembrado para esas dos claves).
- **Hallazgo 3 (media, C1-22 y U4-12):** los cruces leen las claves sembradas en 00001 a 00003 y el inventario `@param` de 00003; el comentario debe citar una decisión de §A–§K (`decisión [A-Z]` o `K[0-9]+`); se quita «(migración 00003)» de C1-22 y se nombran las 5 claves sembradas antes (3 en 00001, 2 en 00002).
- **Hallazgo 4 (media, P3-01, P3-03, P3-12):** las tres Verifica buscan patrones de credencial (`password *[:=]`, `contraseña *[:=]`, `clave *[:=]` solo en P3-03, `bot[_-]?token`, `api[_-]?key`, formato `[0-9]{8,}:[A-Za-z0-9_-]{30,}`) y no la palabra suelta; se aclara que nombres de archivo (`C2-tokens.json`) y advertencias son correctos.
- Cada archivo documenta el cambio en su lista de cambios; `INDICE.md` suma 4 filas en §11 y cita la revisión r2 en el encabezado. Conteo sin cambios: 261 criterios (201 alta / 60 media).

## Archivos
`docs/criterios/P3-claude-design.md`, `docs/criterios/U4.md`, `docs/criterios/C1.md`, `docs/criterios/INDICE.md` y este reporte. No se tocó `docs/criterios/tools/`.

## Verificación
- `bash docs/criterios/tools/verif_final.sh`: conteo 261 (201/60/0); traza 18 archivos, 0 criterios sin línea en §5; refs 261 definidos y 261 citados, 0 inexistentes; citas 305 verificadas, 0 sin fuente; 0 residuos N-xx; canon 178 comprobadas, 0 diferencias; S1-09 SQL y equivalencia con el oráculo 12 cifras, 0 diferencias; 0 rutas /tmp; INDICE contra conteo propio 0 diferencias.
- Paso 9 (`matriz_c102.py`) sigue reportando N-320…N-327 sin operación: hallazgo bajo B-01 de r1 (C1-datos.md, fuera de alcance), sin cambio.
- Prueba de las regex nuevas con texto de muestra: un LEEME que nombra `C2-tokens.json` y advierte «no escribas tu contraseña» da 0 coincidencias; `password: abc` y `bot_token=12` dan 2.
- `grep -n "P3-iter-1" docs/criterios/*.md`: solo la aclaración y el registro del cambio (ya no es evidencia de ningún criterio).

## Supuestos
- La Verifica de P3 sigue apuntando a `docs/design/entrada/`, que se crea en WF-1; no se pudo ejecutar contra archivos reales.
- Los «14 parámetros» cuentan cada umbral, cada ajuste de kcal y cada clave por separado; los días antes del pago [1, 0] son una fila.

## Bloqueos
Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md` (sin «Respuesta de Axel»).

## Pendientes
- Al responder Axel `S1-parametros-avisos`: quitar la excepción de C1-22 y U4-12.
- Cuando S7 proponga `send_time` y `advance_days` (K31, puerta 4): quitar la marca «sin valor inicial» de U4-12.
- Hallazgos bajos de las revisiones r1 y r2 no pedidos.

## Solicitudes fuera de alcance
- Ninguna.

## Propuestas para Axel
- K43 dice que `advance_days` «nace» con valor, pero ningún texto de §A–§K da ese valor (00003 lo declara `@param_sin_valor`). Confirmar en la puerta 2 si el valor lo propone S7 (como asumen los criterios) o si Axel lo da antes.
- Mínimos razonables de seguridad: siguen sin entrar en criterios ni esquema (K44).

## Consumo aproximado
Acciones de herramienta: 14; leído ≈ 75 KB (por partes); escrito ≈ 12 KB.
