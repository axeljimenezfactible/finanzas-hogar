# WF-1 · contrato-S2 · Reporte final
Agente: contrato-S2 (Opus) · 11-oct-2026 · Alcance: `docs/contratos/motores/S2-interfaz.md`, `oraculo/fixtures/FORMATO.md`, este reporte (más los bloqueos de S2 que permite el protocolo §2).

## Hecho
- Interfaz del motor financiero (sin código): convenciones de valores, calendario y redondeo (§1); 12 tipos de entrada con nombre de negocio, nombre de código y origen (§2), incluido el calendario de pagos de S3 como entrada; 15 funciones con firma Go, función Python, nombre JSON y salida (§3), con las vistas del oráculo que emiten las claves literales de `model.json`; tabla de las 159 hojas de S2 → vista/función, fixture y prueba, y destino de las 437 de S3 (§4); 24 fixtures canónicas (§5); cruce motor ↔ oráculo (§6); 14 pruebas Go con los patrones de `docs/criterios/S2.md` (§7); lo que no decide (§8); criterio → sección (§9).
- Formato común de fixtures para todos los motores: carpeta y nombres, forma JSON `{formato, motor, caso, descripcion, fuente, provisional, entradas{datos, consultas}, esperado}`, valores y hojas, conversión de `model.json` a centavos, salida común del motor y del oráculo, y `compara.py <motor>` con sus líneas y códigos de salida.
- Cuatro bloqueos con opciones (ninguno toca cifras canon).
- Verificado con `Decimal` desde `model.json`: las 7 filas de la identidad de la deuda viva (K19), los promedios de K39, el valor estable y su cruce con `sensibilidad`.

## Archivos
- `docs/contratos/motores/S2-interfaz.md` (nuevo, 44.4 KB)
- `oraculo/fixtures/FORMATO.md` (nuevo, 7.1 KB)
- `docs/estado/bloqueos/S2-proyeccion-interruptor.md`, `S2-flujo-quincenal.md`, `S2-metas.md`, `S2-presupuesto-total.md` (nuevos)
- `docs/estado/wf1/contrato-S2.md` (este reporte)

## Verificación
- Hojas numéricas por rama (`python3 -I` sobre `oraculo/model.json`): parametros 9, periodos 122, despensas 144, credito 143, valor_estable 12, sensibilidad 45, papel 12, disposicion_por_despensa 68, horizonte 13, interes_por_monto 28 → **S2 = 159, S3 = 437, total 596**.
- Tabla §4 (expansión de llaves con un analizador por niveles): 20 renglones → 159 hojas, 159 distintas, 0 duplicadas, 0 «n» mal declarado; **hojas de S2 sin dueño: ninguna; rutas sobrantes: ninguna**.
- Criterio → sección: S2-01…S2-16 están en la tabla de §9 y citados en el cuerpo; todas las secciones citadas existen.
- Pruebas Go declaradas: TestCanon, TestOraculoHojas, TestOraculoFixtures, TestSumaPartes, TestQuincenaFinDeMes, TestInvarianteDeudaViva, TestPromedios, TestProyeccionDefecto, TestProyeccionCadencia, TestFlujoQuincenal, TestPresupuestoUmbrales, TestPresupuestoScope, TestMetasAhorro, TestZonaHoy; patrones de S2.md (Canon, Oraculo, Suma, Quincena, FinDeMes, Invariante, Promedio, Proyeccion, Flujo, Presupuesto, Meta, Zona) sin prueba: ninguno.
- Los 18 casos canónicos de la tarea están en §5 (más s2-sensibilidad, s2-deuda-viva-identidad, s2-semana, s2-zona-medianoche, s2-ventanas-dia-5, s2-flujo-quincenal: 24).
- Identidad K19 (deuda = diferencia acumulada + saldo no devengado, en pesos): 31-oct-2026 2,921.25 = 973.75 + 1,947.50; 30-nov 973.75 = 973.75 + 0; 31-dic 3,776.06 = 1,258.69 + 2,517.37; 31-ene-2027 1,258.70 = 1,258.70 + 0; 28-feb 3,776.06 = 1,258.69 + 2,517.37; 31-mar 1,258.70 = 1,258.70 + 0; 10-abr 6,293.44 = 1,258.70 + 5,034.74: las 7 iguales a S2-07, y deuda = Σ costos − Σ efectivo en los 7.
- Promedios (`Decimal`, HALF_UP): 31.41/219.88 (oct), 64.92/454.42 (nov), 89.91/629.35 (feb), 62.82/439.76 (devengado de despensa oct); mensuales 2,183.96 y 2,393.75 (R-56, R-57).
- Valor estable: 2,517.37 + 99.50 = 2,616.87; Axel 1,767.51, Karen 849.36; quincena 1,308.44 / 883.76 / 424.68, diferencia 0 (R-37…R-44). Cruce con `sensibilidad`: costo_mensual_total 1,987.53 · 2,302.20 · 2,616.87 · 2,931.54 · 3,246.21 y axel/karen_mensual reproducidos con la misma regla.
- `sha256sum oraculo/model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (sin cambios). Bloqueos citados en la interfaz: los 4 existen. Tamaño de la interfaz ≤ 45 KB.
- Cobertura de criterios por la interfaz: S2-01 §4/§5/§7; S2-02 §4/§6/FORMATO §5; S2-03 §0.4/§5/§6; S2-04 §0.5/§7; S2-05 §1.3/§2.5/§3.2; S2-06 §1.2/§3.1; S2-07 §3.4; S2-08 §1.3/§3.5; S2-09 §3.2/§3.12/§2.10/§2.12; S2-10 §3.11/§3.6; S2-11 §3.8/§3.9; S2-12 §3.13; S2-13 §0.3/§1.1/§7; S2-14 §1.2/§3.1; S2-15 §7; S2-16 §8.

## Supuestos
- Devengado acumulado de la identidad = meses cerrados (fin de mes ≤ corte): es la única lectura que reproduce las 7 filas de S2-07 (la del 10-abr-2027 en particular); no cambia cifras.
- Las categorías «despensa» y «papel» se identifican por parámetro (`groceries_category`, `paper_category`), porque las migraciones no siembran categorías; el ordinal «despensa n» = orden por (fecha, id).
- El valor estable recibe sus partidas como entrada (`stable_items`): ningún documento dice de dónde salen en producción y C1 v1.2 no lo consume; la semilla las toma de `despensas[1]` y `parametros.papel`.
- Las vistas del oráculo usan claves literales de `model.json` y exigen la forma de la semilla.
- Banderas de presupuesto con comparación exacta en enteros (lectura literal de S2-11 «frontera exacta»).

## Bloqueos
- `docs/estado/bloqueos/S2-proyeccion-interruptor.md` (previsto por S2-09)
- `docs/estado/bloqueos/S2-flujo-quincenal.md`
- `docs/estado/bloqueos/S2-metas.md`
- `docs/estado/bloqueos/S2-presupuesto-total.md`
- Abierto y ajeno: `S1-parametros-avisos` no afecta a S2 (umbrales 80/100 ya son `alerts.*`).

## Pendientes
- Agente Go: `backend/internal/motor_fin/*.go` y `backend/cmd/motores/` (subcomando `s2 dump`) según §3, §6 y §7.
- Agente Python: `oraculo/fin_calendario.py`, `fin_devengo.py`, `fin_presupuesto.py`, `fin_flujo.py`, `fin_metas.py`, `fin_vistas.py`, `fin_run_fixtures.py`, `oraculo/compara.py`, `oraculo/fixtures/s2/*.json` (24), `oraculo/fixtures/s2/_datos/semilla-oraculo.json` y `_datos/insumo-s3.json`.
- Axel, en la puerta 2: los 4 bloqueos de S2.

## Solicitudes fuera de alcance
1. `docs/criterios/S2.md` S2-03: la lista mínima de fixtures nombra «interés por monto», cuyas 28 hojas son de S3 (S2-02, S3-01); proponer «sensibilidad (cruce del valor estable)» y dejar interés por monto a S3. Motivo: evitar que el revisor exija a S2 una fixture que no calcula.
2. `docs/contratos/reglas-negocio.md` RN-28, ejemplo: con el interruptor apagado los meses 7–12 también traen el último pago de la despensa 3 (2027-04-15, R-26), no solo los de la despensa 4. Motivo: el ejemplo omite un pago contratado (S2-09 lo dice con «p. ej.»).
3. `docs/contratos/reglas-negocio.md` RN-18 contra `docs/contratos/C1-datos.md` M-07/M-16: porcentaje con un decimal contra entero; alinear cuando Axel decida (`S2-metas`).
4. `docs/contratos/C1-datos.md` M-01: aclarar el renglón del total (`S2-presupuesto-total`).
5. U1/U3 (`docs/design/`): el `% usado` entero puede mostrar 80 con estado normal entre 79.5 % y 79.99 %; decidir cómo se presenta.
6. `oraculo/verify3.py` línea 89 escribe en `~/reporte-compra-latex/` (S2-04): corregir en WF-1 (agente del oráculo).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- S2-proyeccion-interruptor: opción C (patrón del oráculo como valor inicial sembrado y partidas planificadas para cambiarlo); conserva el ejemplo de RN-28 y cumple K44 (todo configurable).
- S2-flujo-quincenal: (1) C, 0 hasta decidir, o A (presupuesto variable prorrateado por días) por usar solo datos que Axel ya captura; (2) B, mitad y mitad, simétrico con el sueldo quincenal; (3) A, misma quincena o mismo mes.
- S2-metas: (1) A, del mes de hoy al de la meta incluidos (la aportación de este mes cuenta); (2) C, decimal en Metas y entero en la tarjeta; (3) A, tres estados simples.
- S2-presupuesto-total: (1) B, la suma de categorías como dice M-01; (2) A, sin porcentaje ni banderas con plan 0.
- No hay valores operativos de K31 en el alcance de S2.

## Consumo aproximado
38 acciones de herramienta; ~180 KB leídos; ~67 KB escritos (interfaz 44.4, FORMATO 7.1, bloqueos 6.5, reporte y progreso ~9).
