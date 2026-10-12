# Revisión del contrato S3 · ronda 5 (agente revisor-contrato-S3-r5)

11-oct-2026 · WF-1 · Artefacto: `docs/contratos/motores/S3-interfaz.md` v1.4 (923 líneas, 133,628 bytes) · Evaluado solo contra `docs/criterios/S3.md` 1.2, decisiones §C y §K. No escribí nada de lo revisado y no corregí nada.

## Veredicto
**APROBADO.** 0 hallazgos altos y 0 medios. El hallazgo medio de la ronda 4 (P1) queda **cerrado**: lo comprobé con mi propio script, que arma el insumo con la receta nueva de §10.0 leyendo solo `model.json` y reproduce las 73 hojas de `sensibilidad` e `interes_por_monto` y R-37…R-43 (7 de 7). Quedan 3 bajos nuevos (R1–R3) y siguen abiertos los bajos anteriores (Q1–Q4, B1–B9). Ninguno impide implementar.

## Verificaciones propias
Script `datos/tmp/rev-s3-r5/verifica.py` (`python3 -I`, `fractions.Fraction`, solo biblioteca estándar; lee `model.json`, `data.json` y el texto de la interfaz; ya borrado).

| Qué | Método | Resultado |
|---|---|---|
| (a) Hojas por rama | recorrido recursivo de `model.json` (números, sin booleanos) | 596 en total; S3: `parametros` 9, `despensas` 144, `credito` 143, `sensibilidad` 45, `disposicion_por_despensa` 68, `interes_por_monto` 28 = **437**; S2: `periodos` 122, `valor_estable` 12, `papel` 12, `horizonte` 13 = 159 |
| (a) Cruce con la tabla de §12 | los 42 patrones se leen del texto del contrato (no los transcribí) y se pasan a expresión regular (`[*]`, `{a,b}`) | **0 sin dueño, 0 con dos dueños**, 0 renglones con conteo declarado distinto, 0 hojas de S2 que casen con un patrón de S3 |
| (c) 4 ramas desde la semilla | motor propio con las reglas de §4, §5.1 y §5.2 sobre la receta «Datos compartidos» (§10.0): `sched`, `SplitTrunc`, `SplitHalfUp`, intereses por persona, composición, línea de tiempo A/B con `From` = ancla | 364 de 364 hojas iguales (`parametros`, `despensas`, `credito` con sus 17 eventos y las 4 series, `disposicion_por_despensa`); 16 fechas de pago iguales, incluido el 28-feb-2027 |
| (c) P1: receta `sensitivityFromModel` | insumo armado solo desde hojas: montos [300000…500000], tasa 80050/309450, n 4, externa 119150, titular 190300, `accrual_months` 2, `cash_monthly_cents` 9950; reglas de §7 | **73 de 73** hojas (45 + 28); `axel_menor_que_actual` igual en las 5 filas |
| (c) R-37…R-43 (§11.1) | fila de 400000: `MonthlyTotal`, `HolderMonthly`, `ExternalMonthly`, su mitad HALF_UP y `MonthlyPurchase` | **7 de 7** (261687, 176751, 84936, 130844, 88376, 42468, 251737) |
| (c) R-05 y R-06 | d1 a 2 q con la tasa de la tabla 53450/309450 y las partes de §4 contra `data.json` | 9 de 9 (17.27 con HALF_UP a 2 decimales) |
| Cruces S2 (§6.1, §7) | efectivo por ventana del 10 al 9 y deuda viva | P1…P6 6 de 6; horizonte 1310378 = 1310378; deuda viva 629344 = 629344 |
| (b) Criterio → sección | lectura de §13 | S3-01…S3-20, los 20 con sección existente |
| Residuos de P1 | búsqueda de `kueski-semilla`, `sensitivityFromModel`, `s3-interes-por-monto`, `400000` | l. 602, 620, 774 y 803 coherentes entre sí, salvo R1 y R2 (bajos) |
| Bloqueos citados | `ls docs/estado/bloqueos/` | existen `S3-convencion-prestamo`, `S3-formula-tarjeta`, `S3-objetivos-medida`, `S3-origen-cargo-normal`, `S3-serie-sin-intereses` y `U4-iva-credito` |
| K43/K44 en el esquema | `grep` en `00003_puerta1b.sql` y `00001` | `comparator.timeout_seconds` = '5' (l. 86); `optimizer.*` de restricciones sin valor (l. 54–55); `optimizer.weights` sembrado 1,1,1,1 (00001 l. 1396) |
| S3-16 | `sha256sum oraculo/model.json` al inicio y al final | `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` las dos veces |

## Tabla de criterios (evaluación del contrato; el código todavía no existe)
| Criterio | Estado | Evidencia propia |
|---|---|---|
| S3-01 (alta) | Cumple | §11 `TestCanon85` y §11.1: las filas S3 salen de `seedFromModel` y `sensitivityFromModel`, y las S2 de funciones S3, sin archivos de otro autor. Recalculé 364 + 73 hojas, R-05/R-06, R-37…R-43, R-45…R-52 y R-53: todo exacto |
| S3-02 (alta) | Cumple | 437 hojas contadas por rama y cruzadas con los 42 renglones de §12: 0 sin dueño, 0 con dos; origen de las 6 ramas definido en §11 (4 desde el `Input`, 2 desde el insumo) |
| S3-03 (alta) | Cumple | §5.2: ancla, inicio en `From`, B por defecto y `PendingConfirmation`; recalculé las series A y B (28 de 28) y los mínimos |
| S3-04 (alta) | Cumple | §5.1 `sched`: las 16 fechas y el 28-feb-2027 salen con mi implementación; `s3-feb-2027` |
| S3-05 (alta) | Cumple | §5.2 K15 (abono que completa repone K e I) y ejemplo 1 de §10.3 con su `esperado` completo; Σ capital repuesto = capital en las 4 compras |
| S3-06 (alta) | Cumple | §5.1 orden de intereses, `ErrNoRate` sin interpolar, decimales K12; 400000 → 103474 recalculado |
| S3-07 (alta) | Cumple | §5.1 promoción aplicable (4 fronteras en `s3-promocion-0`); serie de capital con intereses 0 PROVISIONAL con bloqueo (Q1 abierto) |
| S3-08 (alta) | Cumple con PROVISIONAL | §5.3 y §5.8; fórmula y calendario PROVISIONALES con `S3-formula-tarjeta` y `U4-iva-credito`; sin datos → `missing_capture` (B1 abierto) |
| S3-09 (alta) | Cumple con PROVISIONAL | §5.4: propiedades FIRMES, tasa del periodo PROVISIONAL (`S3-convencion-prestamo`) |
| S3-10 (media) | Cumple con PROVISIONAL | §5.5 y §5.8: tasa 0 = abonos sin interés; sin datos `missing_capture`; con datos K18 `unsupported` PROVISIONAL |
| S3-11 (alta) | Cumple | §6: orden, llave sin duplicados, Σ por quincena, corte que sustituye (PROVISIONAL) (B8 abierto) |
| S3-12 (alta) | Cumple | §8.1–§8.3 y §8.7: sin valores por defecto, «faltan restricciones», 5 s de `app_setting` y 503 como dato; los 3 casos con restricciones de fixture (R3 bajo; B4, B5 abiertos) |
| S3-13 (alta) | Cumple | §8.4, §8.6, §10.0 (secuencia: tabla a mano antes de correr) y §10.1 (intención de cada óptimo y 4 distintos verificados a mano) (Q2, B9 abiertos) |
| S3-14 (alta) | Cumple | §8.5: mín-máx entre los generados (`valid`, decisión C), empate = 1 declarado, pesos 1,1,1,1 normalizados; 100/200/400 → 1, 2/3, 0 |
| S3-15 (alta) | Cumple | §1.5 y §2: sin `net/http`, sin solver |
| S3-16 (alta) | Cumple | §10.0 y §10.2: el oráculo no lee al motor; `model.json` intacto (sha256 verificado) |
| S3-17 (alta) | Cumple | §1.1–§1.3: `Cents`, `big.Rat`; `float64` solo para comparar `tasa` y `tasa_pct` en prueba |
| S3-18 (alta) | Cumple (en el contrato) | §11: cobertura ≥ 85 % por paquete, `gofmt`, `go vet` |
| S3-19 (media) | Cumple | §4: reparto proporcional al capital, externa HALF_UP y titular el resto; 492.28 / 308.22 y pagos 598.82 / 374.93 recalculados |
| S3-20 (media) | Cumple | §14 y §14.1; 6 bloqueos citados y existentes (Q4 abierto: el reporte del contrato describe la v1.0) |

## (c) Suficiencia para implementar
Núcleo, comparador y oráculo se implementan sin preguntas en lo que P1 dejaba abierto: el origen del `SensitivityInput` y de los montos de `interes_por_monto` está escrito hoja por hoja y lo reproduje. Preguntas menores que quedan (todas bajas): de qué hoja sale el 400000 que elige la fila de §11.1 (R2); si `TestCanon85` exige igualdad también con `s3-interes-por-monto` (R1); las ya listadas en B1, B4, B7, B8, Q2 y Q3.

## (d) Coherencia con K12–K18 y K31
K12 (fracción con decimales capturados; −1 sin texto), K13 (tope opcional, solo marca y desempata), K14 (dos modos, externas HALF_UP con el titular como resto, pago de referencia desde el total exacto), K15 (§5.2), K16 (globales y ajustables por campo, sin valores por defecto, colchón en la línea usada y vivas/deuda en el crédito usado), K17 (mín-máx entre generados, desviación estándar exacta con `isqrt`, unidades base por peso PROVISIONAL), K18 (datos por producto, `ErrMissingCapture`, lista de 4 PROVISIONAL `U4-iva-credito`), K31 (`optimizer.date_window_days` propuesto en 7 días, sin valor sembrado), K43 y K44 (5 s configurables desde `app_setting`): **coherente**. La v1.4 no tocó ninguna regla de cálculo.

## (e) Los 3 casos del comparador
§10.1 exige en cada caso los 4 valores de restricción en la fixture (capturas de fixture, no de Axel), ≥ 2 créditos con límite, ≥ 2 plazos con tasa, una promoción 0 % en la tienda de peor U/C, ≥ 2 tiendas, compras existentes desiguales, un descartado por restricción, la intención de cada óptimo y 4 óptimos distintos verificados a mano antes de correr; §10.0 fija la secuencia (tabla a mano → transcripción sin cambiar números → corrida). El ejemplo 2 trae las restricciones explícitas. Cumple en el contrato; las cifras las fija `docs/estado/S3-escenarios.md`, que aún no existe (lo escribe el optimizador).

## Estado de los hallazgos de la ronda 4
| Hallazgo | Estado | Evidencia |
|---|---|---|
| P1 (media; S3-01, S3-02) | **Cerrado** | §10.0 l. 602 (segundo ayudante `sensitivityFromModel`), receta l. 616–618, l. 620; §11 l. 773–774 (4 ramas desde el `Input`, 2 desde el insumo); §11.1 l. 803 (ya no lee archivos del agente de regresión). Receta reproducida: 73 de 73 hojas y R-37…R-43 7 de 7 |
| Q1–Q4 (bajos) | Abiertos | sin cambios en la v1.4 (§5.1, §10.1, §5.2 `TimelineInput.State`, reporte `contrato-S3.md`) |
| B1–B9 (bajos, ronda 3) | Abiertos | sin cambios en la v1.4 |

## Hallazgos nuevos (todos bajos)
- **R1 (baja; S3-01).** §11.1 l. 803 exige igualdad «si existen `_datos/kueski-semilla.json` y `s3-sensibilidad`», pero §10.0 l. 602 y §11 l. 774 nombran también `s3-interes-por-monto`, y R-79…R-85 de `TestCanon85` salen de ese insumo. Ruta: `docs/contratos/motores/S3-interfaz.md` §11.1. Corrección: «… `s3-sensibilidad` y `s3-interes-por-monto` …».
- **R2 (baja; S3-01).** §11.1 l. 806 elige la «fila `Amount` = 400000» de `Sensitivity` para R-37…R-43, pero no dice de qué hoja sale ese 400000, y las pruebas no pueden llevar cifras canon escritas (§11). Las dos lecturas obvias (`despensas[1].capital` o `sensibilidad[2].N`) dan la misma fila, así que no cambia el resultado. Corrección: «la fila cuyo `Amount` = `despensas[1].capital` (despensa base), leída de `model.json`».
- **R3 (baja; S3-12).** El título de §11 dice que los nombres casan con `-run` de los criterios, pero `TestTiempoClaveSembrada` (la prueba que lee el valor sembrado de `comparator.timeout_seconds`) no casa con `-run Restricciones`, que es el comando de S3-12; con ese comando no se ejecuta la parte «el tiempo máximo se lee de la clave de `app_setting`». Sí se ejecuta con la corrida completa de S3-18. Corrección: renombrarla `TestRestriccionesClaveSembrada` (o decir en §11 que la evidencia de S3-12 incluye `-run TiempoClaveSembrada`).
