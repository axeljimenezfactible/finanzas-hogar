# revisor-S3-r3 · WF-1 · 11-oct-2026

## Hecho
- Revisé S3 de forma independiente (ronda 3 de 3) contra `docs/criterios/S3.md` 1.2, decisiones §C y §K y `S3-interfaz.md` (§8.1 a §8.7 y partes de §10). No corregí nada.
- Volví a correr todas las verificaciones con mis propios comandos y escribí tres scripts de cálculo propios en `datos/tmp/rev-s3r3/` (ya borrados): regresión Kueski (117 de 117), comparador (0 errores en restricciones, puntaje, óptimos y recomendación de los 3 casos) y tabla francesa (3 de 3).
- Comprobé que M1 y M2 de r2 quedaron cerrados, incluida una prueba del oráculo con el núcleo compartido por dos procesos.
- Veredicto: **APROBADO**, con 0 hallazgos altos, 0 medios y 6 bajos (B1 a B4 siguen abiertos; B5 y B6 son nuevos).

## Archivos
- Nuevo: `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-S3-r3.md` (informe).
- Nuevo: `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-S3-r3.md` (este reporte).

## Verificación
- `gofmt -l` vacío; `go vet` limpio. `go test -cover`: 93.7 / 90.0 / 96.0 / 90.8 %, total 94.0 %. `-v`: 133 PASS, 0 FAIL, 0 SKIP; `TestCanon85` 85 de 85.
- `compara.py s3`: 437 de 437, 906/906 en las tres direcciones, diferencias 0, código 0, con la salida guardada por defecto (M1 cerrado).
- La salida guardada tiene el mismo sha256 que mi regeneración (`b84c0db7…`).
- Oráculo con el equipo cargado: la carga tardó 7,828 ms y dio `ok` con 3,870 evaluados; la salida tuvo el mismo sha256 (M2 cerrado).
- Kueski propio contra el dump y `model.json`: 117 de 117 (800.50, 973.75, 598.82/374.93, 5,034.74, 1,258.68/.70, 883.75/.77, series A y B con mínimos 8,094.50 y 5,000.00, 6,293.44, 2027-02-28).
- Comparador propio: 0 válidos con violación; mín-máx y `score_bp` iguales en los 46 válidos; 12 óptimos iguales y 4 distintos en cada caso; motor = oráculo. «Faltan restricciones» sin valores. Motor: 125 ms en la carga. La propuesta de K31 está escrita como propuesta.
- Tarjeta: 0.00 con el pago total antes de la fecha o el mismo día; un día después, 65.98, igual a mi cálculo. Préstamo francés: cuadra al centavo, con ajuste en el último pago.
- `model.json`: sha256 `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`, igual al empezar y al terminar. Independencia motor/oráculo comprobada con `grep`.
- Criterios S3-01 a S3-20: todos cumplen. La evidencia está en la tabla de `REVISION-S3-r3.md`.

## Supuestos
- B5 es baja y no media: el contrato no contradice el oráculo actual (nunca le pidió medir reloj) y la regla está escrita en `S3-escenarios.md` §8 y en el docstring de `cred_escenarios.py`. Así, el siguiente implementador no tiene que adivinarla.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ninguno de esta revisión. Los bajos B1 a B6 quedan para el integrador. Las rutas están en `REVISION-S3-r3.md` §Hallazgos.

## Solicitudes fuera de alcance
- Ninguna nueva. Las correcciones de B4 y B5 en `docs/contratos/motores/S3-interfaz.md` (§2, §8.2, §8.7, §10.0 y §11) le tocan al integrador del contrato.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna nueva. Confirmo la medición del motor: 125 ms en la carga de 3,870 candidatos, el 2.5 % de los 5 s de K43.

## Consumo aproximado
- Acciones de herramienta: unas 42.
- Leído: unos 165 KB (protocolo, §C, §K, criterios S3, extractos de REVISION-S3-r2 y correccion-S3-r2, §8 del contrato, partes del código, de las fixtures y de las salidas de comandos).
- Escrito: unos 22 KB (informe, este reporte y tres scripts temporales ya borrados).
- Temporales: borré `datos/tmp/rev-s3r3/`.
