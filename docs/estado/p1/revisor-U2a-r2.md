# Reporte de agente · `revisor-U2a-r2` (flujo P1, revisión U2a, ronda 2 de 3)
Fecha: 11-oct-2026. Rol: revisor independiente; no escribí nada de lo revisado y no corregí nada. `cerrado_por_tope: false`.

## Hecho
- Revisé la colorimetría 1.1 contra U2-02 a U2-09, K26, K27 y R3-04. El informe está en `docs/estado/p1/REVISION-U2a-r2.md`.
- **Veredicto: APROBADO.** No hay hallazgos altos ni medios. Hay 5 bajos: B-01 y B-02 son nuevos; H-02, H-03 y H-04 vienen de la ronda 1 y siguen abiertos.
- Volví a correr `make verify-color`.
- Escribí dos scripts propios que no importan `colorimetria.py`: OKLab de Ottosson, contraste WCAG 2.x, simulación de Machado 2009 y recorte de gamut por bisección. Con ellos verifiqué los puntos (a) a (f) del encargo.
- Comprobé que H-01 (media) quedó cerrado:
  - El bloqueo cumple el protocolo §2.
  - Los 4 tokens llevan la marca PROVISIONAL.
  - El ΔE 1.3 está en «Cambios 1.1» y el hallazgo 14 en §12.
- Repliqué las cifras del bloqueo: las ventanas de matiz, la búsqueda libre de Bodega Aurrera y la opción 3.
- Probé la guarda `rosas_ximena()` en un espejo: al quitar una marca PROVISIONAL, el script aborta.

## Archivos
- `docs/estado/p1/REVISION-U2a-r2.md` (nuevo): veredicto, verificaciones (a) a (f), tabla de criterios, hallazgos y estado de los hallazgos anteriores.
- `docs/estado/p1/revisor-U2a-r2.md` (nuevo): este reporte.
- Temporales en `datos/tmp/revisor-u2a-r2` (`verif.py`, `verif2.py` y `espejo/`): borrados al terminar.

## Verificación
- **`make verify-color`:** VERIFICACIÓN OK y código de salida 0. md, csv, json y html son idénticos byte a byte, y el CSV tiene 202 filas, 0 con cumple=false.
- **`verif.py`:**
  - JSON frente a HEAD: 190 rutas iguales y en el mismo orden. Cambian 20 `$value` y los mismos 20 `$extensions`, y 5 `$description` con la marca PROVISIONAL.
  - Tabla «Tokens que cambian»: 20 filas, 0 discrepancias.
  - CSV: 202 filas recalculadas. Ninguna difiere en más de 0.01, ninguna tiene hex desfasado y ninguna queda bajo umbral.
  - Gasto: 44.73° en claro y 45.03° en oscuro. Gasto–error: 25.26° y 24.93°.
  - Karen y Bodega Aurrera frente a roles, acción crítica y marca: holgura mínima de 1.146 a 1.197.
  - Karen–Bodega Aurrera: 22.95° y 22.90°.
  - Barrido: 6 pares en claro y 5 en oscuro, iguales a §7.7.
  - Gráficas, ΔL consecutivo mínimo: 0.174 y 0.218 en OKLab; 0.179 y 0.223 en gris; 0.162 y 0.198 con simulación en RGB lineal.
  - Contraste mínimo de las series: 3.11.
- **`verif2.py`:**
  - Ventanas de Karen y Bodega Aurrera: idénticas a las del md, y ninguna las separa de Ximena.
  - Búsqueda libre: Bodega Aurrera solo cabe en 331°–359°, sin separarse de Ximena.
  - Opción 3: `#B86787`, ΔL 0.221 y contraste 3.05.
  - Matiz de los roles: todos dentro de ±2°.
  - Marca y neutro, pasos 50 y 950: iguales al plan.
- **Prueba de la guarda:** salida 1, con el mensaje «Bloqueo U2a-rosas-ximena-bodega: bajo la regla frente a Ximena …».

## Supuestos
- Ninguno cambia cifras, esquema, endpoints ni tokens.
- Doy U2-02 y U2-08 por cumplidos aunque 4 tokens sigan PROVISIONAL. Cuando un criterio admite dos lecturas, el protocolo §2 pide un bloqueo y marcar los valores PROVISIONAL, y U2-08 admite el refuerzo declarado. Decidir el fondo le toca a Axel.
- Simulé el daltonismo con las matrices de Machado 2009 (severidad 1.0), en RGB lineal y en sRGB con gamma, sin `culori`. La variante lineal reproduce las cifras del md.
- «Gris» es Y^(1/3), como define `colorimetria.md` §2.7.
- El rango de gasto de 40° a 45° lo evalué con la tolerancia de ±2° de U2-02 (ver B-01).

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ninguno dentro de mi alcance.

## Solicitudes fuera de alcance
- `docs/estado/bloqueos/INDICE.md`: registrar `U2a-walmart-oscuro` y `U2a-rosas-ximena-bodega`, que suman 5 tokens PROVISIONAL. Le toca al orquestador; U2a ya lo pidió.
- `docs/estado/p1/U2a-1.1.md`: anotar que la «Observación para Axel» pasó al bloqueo `U2a-rosas-ximena-bodega`. Le toca al integrador.
- `docs/estado/wf0/REVISION-r3.md`: marcar R3-04 como atendido. Le toca al integrador.
- B-01, B-02 (1), H-02, H-03 y H-04 son bajos y no bloquean. Se corrigen en otra ronda de U2a si el orquestador lo decide.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones de herramienta: 29 de 60, contando la salida final.
- Leído: unos 130 KB de 400. Ningún archivo grande se leyó entero; usé grep, sed -n y scripts.
- Escrito: unos 40 KB de 150 (scripts 17 KB, informe 16 KB y este reporte 6 KB).
- Fuera de esa cuenta, porque no pasó por el contexto: la copia de `colorimetria.py` al espejo con `cp` y lo que regeneraron `make` y el script.
