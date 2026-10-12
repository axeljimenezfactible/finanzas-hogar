# revisor-contrato-S3-r2 · reporte de agente (WF-1, 11-oct-2026)

## Hecho
- Revisé `docs/contratos/motores/S3-interfaz.md` v1.1 completo, por partes, contra `docs/criterios/S3.md` 1.2 y decisiones §C y §K. Usé FORMATO 1.1, RN-11, RN-21, C1 `creditos-b.yaml` y el esquema solo para comprobar que los artefactos son coherentes entre sí.
- Reejecuté por mi cuenta, con scripts propios, las verificaciones numéricas: conteo de las 437 hojas por rama y cruce con §12 y §10, y las reglas del contrato contra `model.json` y `data.json` (296 comprobaciones).
- Comprobé el estado de los 8 hallazgos de la ronda 1: 7 cerrados y H1 reabierto en parte.
- Informe con veredicto, tabla de criterios, hallazgos y estado de los anteriores: `docs/estado/wf1/REVISION-contrato-S3-r2.md`. **Veredicto: NO APROBADO** (0 altos, 4 medios, 7 bajos).

## Archivos
- `docs/estado/wf1/REVISION-contrato-S3-r2.md` (nuevo)
- `docs/estado/wf1/revisor-contrato-S3-r2.md` (este reporte)
- Temporales `datos/tmp/rev-s3-r2/hojas.py` y `calc.py`: borrados al terminar.

## Verificación
- `python3 -I datos/tmp/rev-s3-r2/hojas.py`:
  - sha256 de `model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`;
  - ramas: `parametros` 9, `despensas` 144, `credito` 143, `sensibilidad` 45, `disposicion_por_despensa` 68, `interes_por_monto` 28 → s3 437, s2 159, total 596;
  - `filas 42 cubiertas 437 sin dueno 0 dos 0 filas con conteo distinto []`;
  - `credito.eventos` = 17, con tipos inicio 1, compra 3, pago 11, papel 2.
- `python3 -I datos/tmp/rev-s3-r2/calc.py` → `comprobaciones 296 fallas 0`. Cubre:
  - las 4 despensas con sus 16 pagos y las partes por persona;
  - los 17 eventos, las series A y B y los mínimos;
  - la propiedad de §5.2 en 16 fechas y los ejemplos 909450/600000 y 809450/500000;
  - `interes_por_monto` y `sensibilidad`;
  - las 30 filas S2 de §11.1;
  - R-05/R-06 y `tasa_pct` binary64;
  - los ejemplos de §3, promoción y MSI, y la fecha 2028-02-29;
  - la desviación estándar entera: 0/100/200 → 81650 y 5,000 series aleatorias contra la definición exacta;
  - la normalización 100/200/400.
  
  La primera corrida dio 31 fallas por un error de signo en mi propio cálculo del inicio de B. Corregí mi script, no el contrato.
- Criterios de la fase y su evidencia: tabla de `REVISION-contrato-S3-r2.md`. S3-05, S3-10 y S3-20 quedan parciales (M2, M3, M4); S3-12 está cubierto con hallazgos (M3, M4); los demás están cubiertos.

## Supuestos
- Evalué la coherencia de las fixtures y la CLI contra `oraculo/fixtures/FORMATO.md` 1.1, la versión vigente en disco, que dice regir para S3 sin excepciones. El contrato cita la 1.0.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Corregir M1–M4 (y, si conviene, B1–B7) en `docs/contratos/motores/S3-interfaz.md` y luego la ronda 3 de revisión.

## Solicitudes fuera de alcance
1. `docs/estado/bloqueos/S3-formula-tarjeta.md` y `S3-convencion-prestamo.md`: actualizar las referencias a §5.3/§5.4/§5.5, `CardPeriodInterest`, `CardMinimumPayment` y `Amortize`. Motivo: B1.
2. Integrador: decidir una sola clave (`decimales` o `decimals`) para las razones con decimales capturados en todos los motores, en `FORMATO.md` §3. Motivo: M1 (b).
3. Integrador: aplicar las solicitudes ya escritas por la corrección r1 a RN-11 y C1 (`Scenario.score_bp`: «entre los escenarios generados, los que respetan las restricciones duras») y a S3.md S3-12 (la parte HTTP queda en WF-2). Motivo: B7 y H8.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna nueva. Para la puerta 2 conviene que Axel ratifique la lectura «escenarios generados = los que respetan las restricciones duras» de §8.5 (B7) y vea la propuesta `optimizer.date_window_days` = 7 del reporte del contrato.

## Consumo aproximado
- Acciones de herramienta: ~29.
- Leído: ~255 KB (protocolo 5.6; decisiones §C y §K ~12; criterios S3 35; contrato 91; FORMATO 15; reportes ~10; bloqueos ~9; RN-11/RN-21 ~6; C1 y esquema ~12; salidas de scripts ~15).
- Escrito: ~31 KB (informe de revisión ~19, este reporte ~4, scripts temporales ~8).
