# Corrección U1 1.1 · ronda 1

Agente: `correccion-U1-r1` · Flujo P1 · 11-oct-2026. Corrige los 3 hallazgos de severidad media (M1, M2 y M3) de `docs/estado/p1/REVISION-U1-r1.md`. No hay hallazgos de severidad alta. Las bajas (B1 a B4) no estaban en la lista y no se tocaron, salvo la parte de B2 que M3 pide: el script ahora falla si la regla del domingo aparece fuera de §0.2 y §7. No se reabrió ninguna decisión de §A a §K y no cambió ninguna cifra en $.

## Hecho
1. **M1 · N-090 con la regla de K21 en los meses 7 a 12** (U1-12, U1-14). `docs/design/U1-arquitectura.md:504`. N-090 ahora tiene dos tramos:
   - **Meses 1 a 6:** la composición de antes (ingresos, pagos de crédito, gastos fijos y venideros, gasto variable estimado y aportación planificada a ahorro). El mes en curso es el 1, igual que en N-091.
   - **Meses 7 a 12:** rige N-318 (K21, RN-28), igual que en N-091. Con el interruptor apagado (el valor por defecto) solo entran los compromisos contratados (pagos programados de compras ya hechas), los gastos fijos y los ingresos del CSV, así que no entran el gasto variable estimado, la aportación a ahorro ni los gastos venideros. Con el interruptor encendido se suman la cadencia bimestral de despensa a crédito y el papel alterno. La pantalla dice qué incluye la proyección (N-093).
   - **En qué me baso:** el «solo» de K21 es una lista cerrada. Así lo leen también `docs/contratos/reglas-negocio.md:273-274` (RN-28 «Proyección de 12 meses · FIRME»: «se proyectan solo los compromisos contratados (pagos programados de compras ya hechas), los gastos fijos y los ingresos del CSV de Axel») y el criterio S2-09 de `docs/criterios/S2.md:52`. Como esto sí se puede afirmar, **no abrí** `U1-proyeccion-quincenas-7-12.md`. El supuesto 2 de `docs/estado/p1/U1-1.1.md:72` queda sustituido: en esos meses ya no queda nada abierto que S2 tenga que decidir. La fila 12 de la lista de cambios (línea 33) dice ahora que la regla de N-318 rige N-090 (por quincena) y N-091 (por mes).
2. **M2 · F-7 con K16** (U1-12, U1-05 b).
   - **Paso 3 (línea 953):** abre el comparador con las restricciones globales y los pesos de P-26.05 y el tope por persona de P-26.12 (N-225). Si faltan restricciones, muestra «Faltan restricciones · Capturar» (abre P-26.05) y no recomienda (K16).
   - **Paso 5 (línea 955):** solo hay recomendación cuando las restricciones están capturadas.
   - **Nota nueva «Sin restricciones capturadas» (línea 959):** mientras Axel no las capture no hay recomendación ni «Usar escenario recomendado». Así se verá el primer uso real.
   - **Atajo desde una canasta (línea 961):** tiene la misma rama.
   - **Lista de cambios:** la fila 11 (línea 32) suma F-7.
3. **M3 · la regla del domingo, solo en §0.2** (U1-14, U1-16, R3-01). La línea 353 (vista de plan_comida) dice ahora «Semana por defecto (§0.2)». La 876 (F-2, paso 1) dice «Abre el plan en Detalle en la semana por defecto (§0.2; N-120, N-121)». La fila 2 de la lista de cambios (línea 23) suma §2.5 (plan_comida) y F-2. La regla PROVISIONAL solo queda escrita en §0.2 (línea 65) y en §7 (línea 1150).
4. **`docs/design/tools/verificar_u1.py`:** 4 comprobaciones nuevas (líneas 141 a 162) y el docstring actualizado:
   - (M3) falla si la regla del domingo aparece fuera de §0.2 y §7, con la expresión `domingos?[\s,;:()]+(la )?(semana )?siguiente`, que también detecta «los domingos, la semana siguiente» y «domingos (la siguiente)»;
   - (M3) la vista de plan_comida y el paso 1 de F-2 remiten a §0.2;
   - (M1) N-090 y N-091 citan N-318 y K21 y tienen «Meses 7 a 12», y N-090, N-091 y N-318 dicen «solo compromisos contratados»;
   - (M2) F-7 tiene «Faltan restricciones» y K16, no dice «pesos por defecto», y la fila 11 de la lista de cambios incluye F-7.

## Archivos
- `docs/design/U1-arquitectura.md` (9 parches con `str.replace`, cada uno verificado con una sola coincidencia; el archivo no se reescribió a mano)
- `docs/design/tools/verificar_u1.py` (2 parches con 4 comprobaciones nuevas y un ajuste de rango)
- `docs/estado/p1/correccion-U1-r1.md` (este reporte)

## Verificación
- `python3 -I docs/design/tools/verificar_u1.py`: sale con 0 y da «RESULTADO: TODO OK».
  - 9 entradas en la barra lateral, 27 pantallas (P-00 a P-26) y 14 secciones de Ajustes.
  - 203 N-xxx (138 L y 65 E); la última es N-319; las nuevas, N-316 a N-319. Las Local son N-003, N-301 y N-317.
  - F-1 = 3 toques. Hay 8 flujos con los 4 obligatorios; F-7 tiene 7 pasos, todos con una pantalla que existe.
  - Marcas PROVISIONAL en las líneas 65, 419, 1149 y 1150. Son las mismas de antes; las dos últimas bajaron 2 líneas por la nota nueva de F-7. Todas están justificadas.
  - 30 cifras en $ (21 distintas), todas en `model.json` o `data.json`, sin cambio. Fechas válidas. Ninguna palabra en inglés.
  - `openapi.yaml` 1.0.0: 197 operaciones, 196 con `x-necesidades`.
  - `PENDIENTE C1 v1.1: ['N-316', 'N-318', 'N-319']`: igual que antes y esperado.
- Las 4 comprobaciones nuevas dan OK:
  - «la regla del domingo solo aparece en §0.2 (líneas 53-68) y §7 (desde 1142); fuera: []»;
  - «la vista de plan_comida (§2.5) y el paso 1 de F-2 remiten a la semana por defecto de §0.2»;
  - «U1-12 K21: N-090 (por quincena) y N-091 (por mes) usan en los meses 7 a 12 la regla de N-318»;
  - «U1-05 F-7 con la rama «Faltan restricciones» (K16), sin restricciones por defecto, y citado en el cambio 11».
- **Prueba negativa:** corrí el mismo script contra una raíz temporal con la U1 anterior a esta corrección (con enlaces simbólicos al oráculo, la semilla, C1 y los bloqueos). Sale con 1 y 4 fallas: «domingo… fuera: [353, 876]», la remisión de plan_comida y F-2, K21 en N-090 y F-7. Las comprobaciones detectan exactamente M1, M2 y M3.
- Búsqueda de «la siguiente»: solo aparece en la línea 1150 (§7); en la 65 (§0.2) dice «la semana siguiente».
- No se agregó ni cambió ninguna cifra en $. Los meses 7 a 12 no tienen cifra canon en el oráculo.

## Supuestos
1. La acción «Capturar» del aviso «Faltan restricciones» de P-20 abre P-26.05 (F-7, paso 3, nota y atajo). Lo propone la corrección del revisor, y P-26.05 es donde se capturan las restricciones globales (N-288). Es navegación de la interfaz: no toca cifras, esquema ni operaciones.

## Bloqueos
- Ninguno nuevo. M1 se resolvió con K21, RN-28 (FIRME) y S2-09, sin abrir `U1-proyeccion-quincenas-7-12.md`.
- Sigue abierto, sin cambios, `docs/estado/bloqueos/U1-lista-vigente-domingo.md`.

## Pendientes
- Nada pendiente dentro del alcance.
- Fuera de la lista pedida quedan los hallazgos de severidad baja B1 (N-316 y la confirmación por la propia persona), B2 (resto: comprobar operación por operación N-075 y N-299 después de C1 v1.1), B3 (`period` en N-151 y N-153) y B4 (rotular la semana en las tarjetas de Inicio). No se tocaron, como pide la instrucción.

## Solicitudes fuera de alcance
1. **`docs/estado/p1/U1-1.1.md:72` (supuesto 2).** Anotar que esta corrección lo sustituye: N-090 aplica en los meses 7 a 12 la regla de N-318 (K21, RN-28) y no deja nada abierto para S2. Motivo: el reporte del agente U1 conserva un supuesto que la revisión clasificó como bloqueo.
2. **`docs/criterios/S2.md` (S2-10 «Flujo por quincena»).** Que el flujo por quincena de los meses 7 a 12 se verifique con la misma composición de RN-28 y S2-09, en los dos modos del interruptor. Con el interruptor apagado, sin gasto variable estimado, aportación a ahorro ni gastos venideros. Motivo: S2-09 solo prueba la proyección mensual y S2-10 no menciona K21. Así U1 (N-090), RN-28 y las pruebas de S2 quedan con una sola definición.

## Paquetes instalados
Ninguno.

## Consumo aproximado
- Acciones de herramienta: unas 33, contando este reporte y la salida final.
- Leído: unos 125 KB, todo por partes. U1 solo en rangos (unos 35 KB); `verificar_u1.py` entero (16 KB); la revisión (18 KB); el protocolo (6 KB); decisiones §K (8 KB); el resto en salidas de comandos.
- Escrito: unos 16 KB (parches de U1, unos 4.5 KB; del script, unos 3.5 KB; este reporte, unos 8 KB).
- Los temporales de `datos/tmp/correccion-u1-r1/` (copia de respaldo, raíz de prueba y salida) ya están borrados.
