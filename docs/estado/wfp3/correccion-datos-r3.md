# correccion-datos-r3 · Corrección de M-1 (y dos bajas de texto) de REVISION-datos-r3

Fecha: 12-oct-2026. Alcance: `docs/design/maquetas/` (JSON y LEEME), `docs/design/tools/ds/generar_sinteticos.py` y `verificar_sinteticos.py`. Sin git add/commit.

## Hecho
- **M-1 (media): `seis_mayores` ahora sigue U1 N-044** («las 6 categorías con más gasto del rango»). `generar_sinteticos.py`: se calcula `GASTO_RANGO` (Σ `real_cents` de los 6 meses de `TPC_MESES` por categoría) y `SEIS_MAYORES` = las 6 de mayor gasto, desempate por el orden de `CAT` (P-26.04). `otras_*` se recalculan con el mismo criterio. Resultado: vivienda, super, transporte, despensa-comun, mercado, restaurantes (antes: vivienda, super, servicios, transporte, despensa-comun, mercado). El texto de `tendencia_por_categoria.nota` dice el criterio.
- **Verificador:** la comprobación 36 recalcula `seis_mayores` desde los datos (Σ real de los 6 meses, mismo desempate) y falla si no coincide.
- **LEEME:** la fila de `gastos_resumen` dice «las 6 con más gasto del hogar devengado de mayo a octubre, U1 N-044».
- **Bajas (solo texto del LEEME, sin tocar el JSON):** B-2 (regla de 12 meses: se cita U3 §4 y U1 N-044/N-085/N-103 y se marca como supuesto), B-1 (cómo armar el aportado y qué hacer con el planificado y con los intereses sin separar) y B-4 (el CSV pendiente de P-26.11 se lee como captura a mano). Nueva sección «Notas de lectura (ronda 3)».

## Verificación
- Regeneración con `generar_sinteticos.py`: determinista (`cmp` idéntico en dos corridas).
- Diff de llaves antes/después (todas las rutas de llaves, listas con `[]`): idéntico. Diferencias de valor: solo `seis_mayores[2..5]`, `nota`, y `otras_{real,plan}_{cents,fmt}` de los 6 meses.
- `python3 -I verificar_sinteticos.py`: 57 OK, `RESULTADO: OK`.
- Prueba negativa: copia con `seis_mayores[5] = "servicios"` -> `FALLA 36` y `RESULTADO: FALLA` (antes daba OK).
- Grep de P3-12 (21 cifras, con y sin coma) sobre `docs/design/maquetas` (incluye LEEME): 0.

## Supuestos
- Desempate por el orden de `categorias` (P-26.04), como pide la corrección sugerida; en los datos actuales no hay empates entre las 6 y la 7.ª.
- `seis_mayores` queda ordenado de mayor a menor gasto del rango (antes, por plan).
- No se cambió la llave `meta.version` (sigue 1.2) para no mover otra llave.

## Bloqueos
Ninguno.

## Pendientes
- B-1 (llaves nuevas `aportado_cents`/`planificado_cents` en `ahorro.tendencia` y separación de intereses): resueltos solo por regla en el LEEME; añadir las series exigiría llaves nuevas, fuera de «no cambies otras llaves».
- B-3 (bajas de la ronda 2: campos sin cifras, ids internos en 18 títulos del calendario, $9,000.00 coincidente, 105 flotantes sin `_fmt`, textos de estado): no tocadas.

## Consumo aproximado
~35 KB leídos, ~6 KB escritos, ~14 acciones.
