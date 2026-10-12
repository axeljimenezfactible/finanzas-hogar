# revisor-contrato-S3-r3 · reporte final (WF-1, 11-oct-2026)

## Hecho
- Revisión independiente, ronda 3 de 3, de `docs/contratos/motores/S3-interfaz.md` v1.2 contra `docs/criterios/S3.md` 1.2 y decisiones §C y §K, con contrastes de consistencia (FORMATO 1.2, esquema 00001–00003, C1 por extractos, bloqueos abiertos). No corregí nada.
- Veredicto: **NO APROBADO**: 0 altos, 4 medios (N1–N4), 9 bajos. Hallazgos r2 M1–M4: **cerrados** los cuatro.
- Informe: `docs/estado/wf1/REVISION-contrato-S3-r3.md` (veredicto, verificaciones propias, tabla de criterios, estado de r2, hallazgos con ruta y corrección).

## Archivos
- `docs/estado/wf1/REVISION-contrato-S3-r3.md` (nuevo)
- `docs/estado/wf1/revisor-contrato-S3-r3.md` (este reporte; antes, nota de progreso)
- Temporales `datos/tmp/rev_s3_r3/hojas.py` y `cifras.py`: creados y borrados.

## Verificación
- `python3 -I datos/tmp/rev_s3_r3/hojas.py oraculo/model.json` (script propio): 596 hojas numéricas; rama S3 437 (parametros 9, despensas 144, credito 143, sensibilidad 45, disposicion_por_despensa 68, interes_por_monto 28); S2 159; 42 renglones de §12 con conteo igual al declarado; 0 sin dueño; 0 con dos dueños.
- `python3 -I datos/tmp/rev_s3_r3/cifras.py` (Fraction): 64 pagos canon = model.json; R-05/R-06: 362900, 53450, 17.27, 32870/223170/111585, 20580/139730/69865 = data.json (9 hojas); ejemplo 1 de §10.3 (4 eventos, series, mínimos) idéntico al `esperado`; series A/B canon = R-65/R-66, mínimos 809450/500000; inicios desde 2026-12-16 (909450/600000) y 2027-02-10 (809450/500000); desviación estándar: 81650 y 0 diferencias en 20,000 series aleatorias contra HALF_UP(1000 σ) exacto; pago de referencia desde total redondeado 32 y 125869 (contra 31 y 125868); binary64(8005000/309450) = 25.868476328970754; promoción 0 % de 309450: intereses por pago −1, −1, −1, +3.
- `sha256sum oraculo/model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (intacto).
- `grep` en `00001` l. 808 (`interest_cents CHECK >= 0`), l. 1396 (`optimizer.weights` 1,1,1,1); `00003` l. 86 (`comparator.timeout_seconds` = 5); `creditos-b.yaml` (orden de `missing_constraints` = §8.1); `00002` l. 123 (`credit_purchase.rate_decimals`); `U4-iva-credito.md` abierto (17:41) y sin cita en el contrato (0 coincidencias); `FORMATO.md` 1.2 (18:23) anterior a la v1.2 del contrato (18:33).
- Criterios: cubiertos S3-01, -02, -03, -04, -05, -09, -11, -12, -13, -14, -15, -16, -17, -18, -19; parciales S3-06 (N2), S3-07 (N1), S3-08 (N3), S3-10 (N1, N3), S3-20 (N1, N3). Evidencia por criterio en el informe.

## Supuestos
- Evalué la forma de las fixtures contra `FORMATO.md` vigente (1.2), porque rige a todos los motores sin excepciones (su §6) y ya existía cuando se cerró la v1.2 del contrato.
- Acepté como lectura válida que «escenarios generados» (K17) = los que respetan las restricciones (decisión C), como declara §8.5.

## Bloqueos
- Ninguno nuevo de mi parte (soy revisor). Propongo al corrector abrir `docs/estado/bloqueos/S3-serie-sin-intereses.md` (N1).

## Pendientes
- Corrección de N1–N4 en `docs/contratos/motores/S3-interfaz.md` (y bloqueo de N1); bajos B1–B9 opcionales.

## Solicitudes fuera de alcance
1. `docs/estado/bloqueos/S3-formula-tarjeta.md` y `S3-convencion-prestamo.md`: actualizar referencias a secciones y nombres de la v1.2 del contrato (B2).
2. `docs/estado/bloqueos/U4-iva-credito.md`: anotar en Impacto que también afecta a `S3-interfaz.md` §5.3 y §5.8 (`ErrMissingCapture`) (N3).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Ninguna de valores K31 (no es tarea del revisor). Para N1 la opción con menos cambios es «con intereses 0, capital = monto en todos los tipos»: no cambia montos ni cifras canon y respeta el CHECK del esquema; la decide Axel en la puerta 2.

## Consumo aproximado
~47 acciones de herramienta; ~225 KB leídos; ~33 KB escritos (informe 17 KB, reporte 5 KB, scripts temporales 9 KB, progreso 1 KB).
