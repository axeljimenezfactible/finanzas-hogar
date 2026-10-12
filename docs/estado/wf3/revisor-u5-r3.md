# WF-3b · revisor-u5-r3 · Reporte

Agente `revisor-u5-r3` (revisor visual independiente, ronda 3 de 3) · 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo con `correccion-u5-r2`. No capturé ni implementé nada y no corregí nada. No hice `git add` ni commit.

## Hecho
1. Comprobé que el binario (11:15:41) y `frontend/dist` (11:15:25) en marcha desde las 11:16 son posteriores a todos los fuentes. Volví a correr completo el comando del capturista: `cd frontend && npx playwright test -c playwright.u5.config.ts` → **174 passed (8.8 min)**, 164 capturas nuevas de 11:34 a 11:43.
2. Recorrido propio (Playwright 1.64 + axe-core 4.13, scripts míos en `datos/tmp/rev-u5-r3/`, ya borrados) de las 41 pantallas × 4 combinaciones, con página nueva por pantalla: axe, contraste propio, hosts, texto visible, tema, h1, escala tipográfica, `tabular-nums`, iconos, radios, sombras, toque con prueba de impacto de 44 × 44, foco. Espaciado de 4 pt y tarjetas en la muestra de U5-06.
3. 20 pares de contraste del capturista recalculados con mi implementación.
4. F-1 (U5-09) y captura sin conexión (U5-13) con una prueba mía contra un servidor aislado (copia de la base en `datos/tmp/rev-u5-r3/srv`, `127.0.0.1:8445`, HOME aislado). La base de desarrollo no se tocó.
5. `c3-09.spec` reejecutado (solo lectura).
6. Lista visual U5-06 con Read de las capturas (hojas reducidas y recortes): P-00, P-01, Ajustes completos, una por entrada de la barra, P-03, P-22 y Tendencias, en las 4 combinaciones.
7. Escribí la sección «## Ronda 3» de `docs/estado/U5-revision.md`, `REVISION-u5-r3.md` y el bloqueo de escalamiento.

## Archivos
- `docs/estado/U5-revision.md` (sección Ronda 3 agregada)
- `docs/estado/wf3/REVISION-u5-r3.md` (nuevo)
- `docs/estado/wf3/revisor-u5-r3.md` (este)
- `docs/estado/bloqueos/U5-ronda3-escalamiento.md` (nuevo)
- Regeneradas por el comando del capturista: `docs/design/capturas-u5/` (164 PNG, `axe/`, `contraste/`, `hosts.json`, `evidencia/`) y `datos/tmp/u5-flujos/c3-09.json`.

## Verificación
- **U5-01** · Pasa · `python3 -I frontend/pruebas-u5/verificar_capturas.py` → T = 41, 4T = 164, 164 de 164, 0 faltantes; `ls capturas-u5 | wc -l` = 168 (bloqueo `U5-conteo-carpeta`); `scrollWidth = innerWidth` en 164 de 164.
- **U5-02** · Pasa · corrida completa 174 passed; $10,000.00, $6,905.50, $973.75, $6,293.44, $3,094.50 y $1,191.50 visibles = `model.json` (`disponible_A`, `disponible_B`, `pagos`, `deuda_viva_10_abr`).
- **U5-03** · Pasa · comando literal = 0; axe propio con wcag22a: 0 serias o críticas y 0 de otro impacto en las 4; incompletos 15 · 15 · 6 · 6.
- **U5-04** · Pasa · mío: 4,327 · 4,333 · 5,567 · 5,597 pares, 0 bajo umbral; capturista: 4,588 · 4,623 · 5,833 · 5,875, 0; 20 de 20 pares con Δ máx. 0.0005 (15 también en mi DOM con la misma razón).
- **U5-05** · Ronda 3 registrada; reejecuté U5-01, U5-03 y U5-04 antes de opinar. **Veredicto: devuelto** → escala a Axel.
- **U5-06** · Falla · 13 de 14 elementos pasan; #10 falla: la barra lateral va sin iconos de sección (M-7, Ec y Eo, §7.9 y §5.4). Los 4 medios de r2 están cerrados.
- **U5-07** · Pasa · 351 importes por combinación, 0 fuera de formato, 0 fechas ISO, 0 palabras de la lista, 0 claves internas; N-1 cerrado; «Sin conexión con el servidor de casa · 2 en cola»; 10 pantallas leídas completas. Baja B-10 («Versión dev»).
- **U5-08** · Pasa · 1 host `finanzas.local:8443` en 164 cargas (capturista y mío).
- **U5-09** · Pasa · 3 toques («+», Despensa, Guardar); la tarjeta 1 de Inicio cambia; «$187.25» en la tabla de P-02.
- **U5-10** · Pasa · fondo del body = C2 (`#FAFAFB` / `#222325`) en 164 de 164, sin `data-theme`; `tema.spec` pasa.
- **U5-11** · Falla (media, M-8) · C3 v0.95 marca «Capturado: no» en las 41 filas; el cruce da 0 pantallas sin sus 4 capturas.
- **U5-12** · Pasa · bloqueos `U5-conteo-carpeta`, `U5-contraste-ui-controles`, `U5-nombres-propios`, `U5-tipografia-tres-tamanos` y el nuevo `U5-ronda3-escalamiento`; nada aprobado con violaciones abiertas.
- **U5-13** · Pasa · los 7 pasos: franja y contador 1 → 2; 2 tras recargar; claves de 36 caracteres y `captured_at` ISO con desplazamiento anterior al envío; base 6 → 9 gastos (uno por captura: 187.25, 43.21 y 55.55) y 0 → 1 peso; reenvío = 200, `Idempotent-Replayed: true`, mismo id; un 500 deja «No se pudo enviar · 1 en cola» y 0 registros.
- **C3-09** · Pasa · 111 cifras, 0 diferencias.

## Supuestos
1. Página nueva por pantalla en mi recorrido, igual que el capturista desde r2, para medir con el puntero grueso real en Android.
2. B-10 («Versión dev») va como baja: es el valor de compilación de `sistema.Version`, no texto de la interfaz ni palabra de la lista de U5-07.
3. La sombra de la pestaña elegida (B-11) y la de la perilla (B-12) van como bajas, porque son detalles de componente y no secciones en reposo.
4. U5-11 cuenta como medio aunque el arreglo sea documental, porque el criterio pide la marca y la ruta en la matriz.
5. Las mediciones de foco «sin foco» de mi pasada automática en 6 pantallas de escritorio salían de mi propio script. Las comprobé a mano con Tab: anillo de 2 px en cada parada.

## Bloqueos
- Nuevo: `docs/estado/bloqueos/U5-ronda3-escalamiento.md` (ronda 3 devuelta; decide Axel, sin ronda 4).
- Siguen: `U5-conteo-carpeta`, `U5-contraste-ui-controles`, `U5-nombres-propios`, `U5-tipografia-tres-tamanos`.

## Pendientes
- M-7 (iconos de la barra lateral) y M-8 (columna «Capturado» de C3), según lo que decida Axel.
- Bajas B-1, B-2, B-3, B-5 a B-8 y B-10 a B-15.
- Observación sin severidad: P-06 muestra en la quincena 2 el mismo real ($1,197.64) que en la quincena 1.

## Propuestas para Axel (PROPUESTA, no decidido)
| Qué | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. Corrección acotada de M-7 y M-8 con verificación puntual de esos 2 puntos | U5 cerrable sin otra ronda completa | Bajo | Bajo |
| B. Aceptar U5 con M-7 y M-8 como pendientes con fecha | Avanza a la puerta 4 | Ninguno ahora | Barra lateral distinta de DESIGN; C3 incompleto en la puerta |
| C. Otra ronda completa | Revisión completa de nuevo | Medio | Rompe el límite de 3 rondas de U5-05 |
Recomendación: A.

## Consumo aproximado
≈ 59 acciones de herramienta · ≈ 150 KB de texto leído, más ≈ 0.84 MB en 7 hojas JPEG reducidas de capturas · ≈ 30 KB escritos (scripts temporales ≈ 17 KB, ya borrados, y reportes ≈ 13 KB). Corridas: capturas 8.8 min, recorrido propio 6.6 min, muestra de espaciado ≈ 2.5 min.
