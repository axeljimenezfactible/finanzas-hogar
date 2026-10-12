# WF-3b · revisor-u5-r2 · Reporte

Revisor visual independiente, ronda 2 de 3. No capturé ni implementé nada y no corregí código. 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo.

## Hecho
1. Volví a correr completo el comando del capturista (`npx playwright test -c playwright.u5.config.ts`: 170 passed, 7.9 min), porque 136 de las 164 capturas eran anteriores al binario en marcha. Corrí además los comandos literales de U5-01 y U5-03 y `resumen_u5.py`.
2. Medición propia en `datos/tmp/r2/`, independiente de los scripts del capturista, en las 41 pantallas × 4 combinaciones:
   - axe-core propio;
   - contraste propio: texto, íconos, controles y estado elegido;
   - DOM: h1, barra lateral, iconos, espaciado, radios, sombras, `tabular-nums`, tamaños, toque, foco, fondo, hosts y texto visible para U5-07.
3. Spec propia de U5-09 y U5-13 contra el servidor de prueba (:8444, copia de la base), con la secuencia de H-5. También corrí `captura-sin-conexion` (5 passed) y `c3-09` (111 cifras, 0 diferencias).
4. Revisión visual (Read de hojas de contacto armadas desde los PNG): P-01 × 4, Ajustes completos en Ec y Eo, entradas de la barra en Ao y Ec, P-02, P-03, P-17 y P-26.09 en Android, y Tendencias de P-01 (Ao, Ec) y P-10 (Ec).
5. Escribí `docs/estado/U5-revision.md` §Ronda 2 y `docs/estado/wf3/REVISION-u5-r2.md`. **Veredicto: devuelto.**

## Archivos
- `docs/estado/U5-revision.md` (sección «## Ronda 2» agregada).
- `docs/estado/wf3/REVISION-u5-r2.md` (nuevo) y este reporte.
- Regenerados por el comando del capturista: `docs/design/capturas-u5/` (164 PNG, `axe/`, `contraste/`, `evidencia/`, `hosts.json`) y la evidencia de `frontend/pruebas-u5/flujos/`.

## Verificación
- U5-01: `verificar_capturas.py` → 164 de 164, 0 faltantes. `ls | wc -l` = 168 (bloqueo `U5-conteo-carpeta`). DOM propio: sin desborde en 164. **Pasa**, con la salvedad de fidelidad Android (N-2).
- U5-02: corrida completa reejecutada, 170 passed. Cifras visibles = `model.json` ($10,000.00, $6,905.50, $973.75, $3,094.50, $6,293.44, $1,191.50). **Pasa**.
- U5-03: comando literal = 0. Mi axe: 0/0/0/0 serias o críticas y 0 de otro impacto. Incompletos revisados (Radix y P-26.10 → B-8). **Pasa**.
- U5-04: capturista 0 bajo umbral (4,812 / 4,826 / 6,022 / 6,057 evaluados). Mi implementación en reposo da 0, con 45 estados elegidos ≥ 3:1 por combinación. 20 pares al azar: 20 de 20, Δ máx. 0.000. **Pasa**.
- U5-05: ronda 2 registrada; reejecuté U5-01, U5-03 y U5-04. **Devuelto**.
- U5-06: 14 elementos. Fallan #2 (P-00 h1 600), #7 (sombra en Ajustes, claro), #12 (capturas Android sin puntero grueso) y #14 (etiquetas encimadas en P-01 Tendencias). **Falla**.
- U5-07: 541 importes por combinación en formato (el «$1.5 mil» es etiqueta de eje, que §7.7 permite); 0 fechas ISO; 0 palabras de la lista; «Sin conexión con el servidor de casa · n en cola» presente. Pero en P-26.13 se lee «Fixture s2-metas.json, savings_goals[2]». **Falla** (N-1).
- U5-08: 1 host, `finanzas.local:8443`. **Pasa**.
- U5-09: 3 toques medidos sin red. **Pasa**. La aparición en Inicio y en la tabla no la volví a comprobar.
- U5-10: fondo C2 en 164 de 164; `tema.spec` passed. **Pasa**.
- U5-11: **no evaluado** (sin acciones disponibles). Pendiente para la ronda 3.
- U5-12: supuestos y bloqueos presentes. **Pasa**.
- U5-13: 7 de 7 pasos con mi spec y 5 passed con la del capturista; H-5 cerrado. **Pasa**.
- C3-09: 111 cifras, 0 diferencias. **Pasa**.

## Supuestos
1. Una etiqueta de eje compacta como «$1.5 mil» no es un importe para U5-07, porque DESIGN §7.7 la permite en ejes.
2. «Stock por unidad» (P-16) no cuenta como inglés: es el término de `decisiones.md` y U1 (P-16, P-26.14).
3. El estado *hover* no entra en U5-04, que mide en reposo. Los 8 pares que vi con *hover* los provocó mi propio script.
4. N-2 es media y no alta: la app cumple con el puntero grueso; lo que falla es la fidelidad de las capturas Android a DESIGN §4.6.
5. N-1 es alta aunque el texto venga de la semilla de desarrollo: U5-07 pide que todo el texto visible de las capturas esté en español.

## Bloqueos
Ninguno nuevo. Siguen abiertos `U5-contraste-ui-controles` (con la lectura más estricta hoy da 0), `U5-nombres-propios`, `U5-tipografia-tres-tamanos` y `U5-conteo-carpeta`.

## Pendientes
- Ronda 3, la última: N-1, N-2, N-3, M-6b y M-1b. Evaluar también U5-11 y B-4.
- `git add` y commit los hace el integrador.

## Propuestas para Axel (PROPUESTA, no decidido)
| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| Corregir los 5 hallazgos y volver a capturar las 164 con puntero grueso comprobado | Ronda 3 con evidencia fiel de Android | Una tanda de corrección corta (texto de semilla, una clase de sombra, una clase en P-00, gráficas y script de captura) más 8 min de captura | Bajo; es la última ronda, así que conviene que el corrector vuelva a correr también mi comprobación de puntero grueso |
| Aceptar N-2 como nota de método sin volver a capturar | Ahorra la recaptura | Ninguno | La evidencia de Android seguiría mostrando controles de 36 px |

Recomendación: la primera opción.

## Consumo aproximado
≈ 60 acciones (en el límite). ≈ 150 KB de texto leído, más 8 imágenes de hojas de contacto. ≈ 30 KB escritos (reportes ≈ 20 KB; scripts temporales ≈ 10 KB, ya borrados). Las corridas generaron ≈ 25 MB de PNG y JSON en temporales.
