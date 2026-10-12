# Reporte · `revisor-c3-r2` (WF-3b, revisión de C3, ronda 2 de 3)

Fecha: 2026-10-12 · Agente: `revisor-c3-r2` (revisor independiente; no redactó la matriz ni implementó nada) · Objeto: `docs/contratos/C3-integracion.md` v1.1 · Criterios: `docs/criterios/C3.md` 1.1 · Informe: `docs/estado/wf3/REVISION-c3-r2.md`.

## Hecho
- Reejecuté todas las verificaciones de C3: C3-02 y C3-03 con un script propio (`python3 -I`, independiente de `frontend/scripts/c3-integracion.py`), `npm run build`, `go test ./... -count=1 -v`, `npx playwright test` (JSON) con cruce de los nombres de prueba, la existencia de las 4N capturas, el `grep` de C3-06, los flujos de C3-08 en copias aisladas nuevas, la prueba de cifras de C3-09 con lectura de la API, y C3-10 y C3-11.
- Comprobé los 3 hallazgos altos de la ronda 1: **H-2 y H-3 cerrados** (F-7 y F-8 tienen prueba y pasan) y **H-1 reabierto** (F-3 sigue fallando con la semilla).
- **Veredicto: devuelto, sin firma** (1 alto, 0 medios, 6 bajos). No modifiqué `C3-integracion.md`: su SHA-256 es `c9e622b1599d33b87fe76f61d19db1dbf96c2938789020979066bac617232044` antes y después.

## Archivos
- Escritos: `docs/estado/wf3/REVISION-c3-r2.md` y `docs/estado/wf3/revisor-c3-r2.md`.
- Temporales en `datos/tmp/c3r2/` (registros, scripts y respaldo de capturas) y la copia de prueba `datos/tmp/u5-flujos/`: los borré al terminar.
- Respaldé `frontend/pruebas-u5/flujos/capturas/` y `resultados/` antes de las corridas y los restauré después, así que la evidencia citada por la matriz quedó igual.

## Verificación
- **C3-01** · cumple · `grep -c "^| P-"` da 41 = T (27 + 14), más la fila Global. §1 trae las 3 definiciones.
- **C3-02** · cumple · Mi script da 0 pantallas sin fila, 0 filas sin pantalla, 211 = 211 `N-xxx` (3 Local), 0 sin operación ni Local, 0 de C1 fuera de U1 y 0 diferencias en la columna de operaciones. Hay 5 operaciones sin pantalla, todas con excepción. `openapi.yaml` y `routes.go` coinciden (216 = 216).
- **C3-03** · cumple · Hay 208 pares (método, ruta) y 166 literales, con 0 fuera del contrato y 0 sin registrar. En el cierre de importaciones solo falta N-054 en P-01, que U1 declara «vía P-03». `npm run build` sale con 0 y el árbol de git no cambia.
- **C3-04** · cumple · Go: 32 `ok`, 1126 PASS, 0 FAIL, 216/216 `TestContrato_Operaciones`. Playwright: 245 ejecutadas, 241 pasan, 4 omitidas, 0 fallan. Los 60 nombres visibles de la matriz pasaron.
- **C3-05** · cumple · 164 de 164 imágenes existen, son PNG válidos con el ancho correcto y las citas coinciden. La carpeta tiene 164 PNG.
- **C3-06** · cumple · `grep -c "| no |"` da 1: la fila Global, con el bloqueo `C3-capturas-global.md`.
- **C3-07** · no se firma · el veredicto es devuelto.
- **C3-08** · **no cumple** · `flujos.spec.ts` pasa 6 de 7 (falla `:122`, F-3), `f7-f8.spec.ts` 2 de 2 y `pendientes.spec.ts` 4 de 4, cada uno en una copia nueva. Pasan 7 de 8 flujos y F-1 se registra en 3 toques.
- **C3-09** · cumple · `c3-09.spec.ts`: 111 cifras con 0 diferencias. `humo-inicio.spec.ts`: 1 de 1. En la API, el próximo pago es 97375 el 2026-10-30, la lectura B es 690550 y `getLiveDebt` al 2027-04-10 da 629344, iguales a `model.json`.
- **C3-10** · cumple · Leí 13 filas y §3 a §8: es-MX y fechas ISO. El anglicismo de r1 ya no está.
- **C3-11** · cumple · `grep -n "v1"` encuentra v1.1 en la línea 3. Hay 3 bloqueos `C3-*`, y §Supuestos y §Bloqueos están en `c3-matriz.md` y `correccion-c3-r1.md`.

## Supuestos
1. Tomo N = 8 flujos para C3-08 (F-1 a F-8 de U1 §4), igual que la matriz (supuesto 4) y la ronda 1, porque U1 cubre los obligatorios (b) y (d) de U1-05 con F-7 y F-8.
2. Para C3-08 solo cuento como tope el de registrar gasto. Los toques de F-7 (7 hasta P-19, U1: 6) y de F-8 (15, U1: 5) no cambian el resultado; ya están anotados en `C3-flujos` puntos 6 y 7.
3. Para leer la API inicié sesión como `axel` en el servidor de desarrollo, solo con lecturas, y cerré la sesión al terminar. La contraseña se leyó del archivo fuera del repositorio y no se imprimió ni se guardó.

## Bloqueos
- No abro bloqueos nuevos. El que impide la firma ya existe: `docs/estado/bloqueos/C3-flujos.md`, puntos 2 y 3 (F-3), y lo decide Axel.

## Pendientes
- Ronda 3 (la última): reejecutar F-3 en una copia aislada cuando Axel decida los puntos 2 y 3 de `C3-flujos`, y firmar si todo lo demás sigue igual.
- Bajos para el corrector: R2-2 (citar `:141`), R2-3 (`getLiveDebt` y cerrar el supuesto 6), R2-4 (lista completa de pruebas), R2-5 (`expect(reg.ok)` en F-3A0).

## Propuestas para Axel
| Decisión | Resultado | Costo | Riesgo |
|---|---|---|---|
| Que la semilla traiga los precios de los ingredientes de las recetas (`C3-flujos` punto 2 b) | F-3 ruta A pasa con la semilla, como ya pasa con 18 precios de prueba (`pendientes.spec.ts:117`) | Bajo: datos de semilla o agente externo K53 | Bajo |
| Decidir qué cifra de Walmart vale: la de U1 ($3,233.00) o la del motor S5 (+$264.50) (`C3-flujos` punto 3) | F-3 ruta B pasa al corregir U1 o el motor | Bajo si se corrige U1; medio si se corrige el motor | Medio: es la cifra que ve el hogar |
| Recomendación | Hacer las dos cosas: es el único bloqueo que queda para firmar C3 en la ronda 3 | | |

## Consumo aproximado
- Leídos: unos 110 KB (C3.md, la matriz por partes, REVISION r1, `C3-flujos.md`, extractos de U1 y de `openapi.yaml`, y salidas de las ejecuciones).
- Escritos: unos 24 KB (2 reportes y scripts temporales ya borrados).
- Acciones: unas 36.
