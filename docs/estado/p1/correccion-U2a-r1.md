# Corrección U2a · ronda 1 (flujo P1) · `correccion-U2a-r1`

Fecha: 11-oct-2026. Corrige el único hallazgo de severidad alta o media de `docs/estado/p1/REVISION-U2a-r1.md`: **H-01** (media; U2-02 y K26). Los hallazgos bajos H-02, H-03 y H-04 no estaban en el encargo y no se tocaron. No se reabre ninguna decisión de §A–§K.

## Hecho
- Nuevo bloqueo `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md` con el formato del protocolo §2. Trae 4 opciones, cada una con su consecuencia calculada:
  1. Aceptar con el refuerzo declarado.
  2. Mover a Ximena.
  3. Reordenar la L de las tiendas en oscuro.
  4. Otra salida, con cambios fuera de K26.
- El cálculo de las opciones encontró algo que no estaba escrito: con las reglas de §7, ninguna de las opciones 2 a 4 cumple «≥ 20° o ΔL ≥ 0.15» frente a Ximena si solo se mueve a Karen, a Bodega Aurrera o a Ximena (cifras en el bloqueo).
- Cuatro tokens quedan PROVISIONAL: `color.person.karen.light`, `color.person.karen.dark`, `color.store.bodega-aurrera.light` y `color.store.bodega-aurrera.dark`. Son los valores de la 1.1 que quedan bajo la regla frente a Ximena. La nota que cita el bloqueo va en su `$description`. Ningún `$value` cambia.
- Parches a `docs/design/tools/colorimetria.py` (11 reemplazos, cada uno verificado con conteo exacto):
  - Constantes `BLOQUEO_ROSAS` y `ROSAS_XIMENA` y cuatro entradas nuevas en `PROVISIONAL`. La variante `light` de las categóricas ahora también lee `PROVISIONAL`.
  - Funciones nuevas `rosas_ximena`, `md_rosas_ximena` y `hallazgo_rosas_ximena`. La guarda aborta en dos casos: si las marcas no coinciden con los tokens que quedan bajo la regla, o si algún matiz de las ventanas los separa de Ximena.
  - Nota PROVISIONAL en el motivo de Karen (`PERSONAS`) y de Bodega Aurrera (`TIENDAS`), en la columna «Origen» de «Cambios 1.1», en la introducción de §7 y en los 4 renglones de §7.7.
  - Hallazgo 14 en §12 y el bloqueo nuevo en su línea **Bloqueos**.
- `docs/design/colorimetria.md` regenerado por el script. En «Cambios 1.1» §2 entra un párrafo con las cifras:
  - Bodega Aurrera–Ximena: ΔE 1.3 en oscuro y 13.1 en claro.
  - Karen–Ximena: ΔE 10.9 y 11.6.
  - Referencia de la 1.0: ΔE 10.5.
  §7.1, §7.2 y §7.7 marcan los tokens PROVISIONAL, y §12 trae el hallazgo y el bloqueo.
- `docs/contratos/C2-colores.json` regenerado: cambian solo 4 `$description`.
- `contraste.csv`, `paleta.html` y `paleta.png` quedan sin cambios. El HTML no muestra descripciones y sale idéntico byte a byte, así que el PNG no cambia.

## Archivos
- `docs/design/tools/colorimetria.py`: parches.
- `docs/design/colorimetria.md`: regenerado.
- `docs/contratos/C2-colores.json`: regenerado.
- `docs/design/contraste.csv`, `docs/design/paleta.html` y `docs/design/paleta.png`: regenerados e idénticos a los anteriores.
- `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md`: nuevo.
- `docs/estado/p1/correccion-U2a-r1.md`: este reporte.

## Verificación
- `python3 -I docs/design/tools/colorimetria.py --sin-png`:
  - Escribe md, csv, json y html.
  - `contraste.csv`: 202 filas, 0 con cumple=false.
  - 190 tokens, igual a las 190 filas del Anexo A.
  - Plan §4.3: 22/22.
  - Personas y tiendas: 0 colisiones con roles en ambos temas. Zona reservada: 0 tokens dentro.
  - Resultado: **VERIFICACIÓN OK**.
- `make verify-color`: **VERIFICACIÓN OK**. md, csv, json y html son idénticos byte a byte a la regeneración en el árbol espejo, y el CSV tiene 202 filas, 0 con cumple=false.
- Recálculo propio (`verif_propia.py`). Usa mi propia implementación de OKLab de Ottosson y de la conversión OKLCH→sRGB, sin importar `colorimetria.py`:
  - Tokens: 190 en la versión nueva, en la previa y en HEAD, con las mismas rutas. Frente a la previa no cambia ningún `$value`, `$type` ni `$extensions`. La `$description` cambia solo en los 4 tokens, que son exactamente los marcados con el bloqueo. En total hay 5 PROVISIONAL: esos 4 más `color.store.walmart.dark`.
  - Ximena–Karen: 17.1°, ΔL 0.103 y ΔE 10.93 en claro; 16.7°, 0.111 y 11.59 en oscuro.
  - Ximena–Bodega Aurrera: 5.8°, 0.130 y 13.10 en claro; 6.2°, 0.006 y 1.33 en oscuro.
  - Los 4 pares quedan bajo la regla, y las cifras coinciden con el md.
  - En la 1.0 (`git show HEAD`), los pares persona–tienda bajo la regla eran Karen–Bodega Aurrera en oscuro (ΔE 10.51), Axel–Walmart en oscuro (11.59) y Conjunto–Walmart en claro (12.19). El mínimo, 10.5, coincide con la revisión.
  - Réplica del barrido: entre 331° y 359°, ningún matiz separa de Ximena en ambos temas a Bodega Aurrera (L 0.410/0.829, croma 0.11) ni a Karen (L 0.642/0.725, croma 0.12). La réplica reproduce los hex `#FFAAC9` y `#D389C6`.
- Prueba de la guarda: al quitar la marca de `color.store.bodega-aurrera.dark`, `rosas_ximena()` aborta con «Bloqueo U2a-rosas-ximena-bodega: bajo la regla frente a Ximena …».
- Diff de `colorimetria.py` frente al estado previo: 11 bloques y 92 líneas, todas del parche.
- Cifras de las opciones del bloqueo (`analisis.py`, `analisis2.py` y `opcion3.py`, con el modelo del script):
  - Ximena con matiz y L libres: se separa de Bodega Aurrera solo en 331°–336°. De Bodega Aurrera y Karen a la vez, en ninguno.
  - Bodega Aurrera en oscuro, con L de 0.55 a 1.00 y croma de 0.10 a 0.12: ninguna L la separa de Ximena.
  - Las 24 permutaciones de L de las tiendas en oscuro: las 6 que la separan de Ximena hacen fallar a Mercado frente a los roles.
  - Bodega Aurrera con matiz, L y croma libres: solo cabe en 331°–359°, y en ningún caso con la lectura estricta. Con Sam's Club en L 0.829, no cabe en ningún matiz.
  - Intercambio de Walmart y Bodega Aurrera: ΔE 11.9. Walmart queda con croma 0.085, a 18.2° y ΔL 0.044 de la marca 400.

## Supuestos
- Ninguno que cambie valores. El criterio para marcar PROVISIONAL sigue el precedente de `U2a-walmart-oscuro`: se marca el valor entregado que está en duda, no cada token que una opción tocaría. Por eso entran los tokens de la 1.1 (Karen y Bodega Aurrera) que quedan bajo la regla frente a Ximena, y no los de Ximena, que conserva su valor de la 1.0. Marcar también a Ximena (opción 2) sería una entrada más en `PROVISIONAL`. Walmart en oscuro (opción 3) ya está marcado.

## Bloqueos
- `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md` (nuevo): 4 tokens PROVISIONAL.
- `docs/estado/bloqueos/U2a-walmart-oscuro.md` sigue abierto, sin cambios.

## Pendientes
- Ninguno dentro del alcance. Los hallazgos bajos H-02, H-03 y H-04 de la revisión no estaban en el encargo.

## Solicitudes fuera de alcance
- `docs/estado/bloqueos/INDICE.md`: registrar `U2a-rosas-ximena-bodega` (PROVISIONAL, 4 tokens) para la puerta que toque. Motivo: el índice lo lleva el orquestador.
- `docs/criterios/U2.md` (U2-02 y U2-08): cuando Axel decida, aclarar si «las demás entradas» admite el refuerzo de U2-08 entre categóricas, tanto de grupos distintos (persona–tienda) como del mismo grupo (Karen–Ximena). Motivo: esa doble lectura originó H-01; es la «otra vía» que propuso el revisor.
- `docs/estado/p1/U2a-1.1.md` (Supuestos): anotar que la «Observación para Axel» sobre Bodega Aurrera–Ximena pasó al bloqueo `U2a-rosas-ximena-bodega`, con ΔE 1.3. Motivo: es el reporte de otro agente.
- U2 `DESIGN.md` (WF-1): llevar lo que decida el bloqueo. Con la opción 1, una persona y una tienda nunca se distinguen solo por color en la misma vista.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones de herramienta: 32 de 60.
- Leído: ≈ 125 KB de 400. Ningún archivo grande se leyó entero; solo se usaron grep y sed -n.
- Escrito: ≈ 50 KB de 150, entre scripts temporales, parche, bloqueo y este reporte. Además, el script regeneró el md y el json, y el parche reescribió `colorimetria.py`. Esas escrituras las hizo el programa y no pasaron por el contexto.
- Los temporales de `datos/tmp/corr-u2a-r1` se borraron al terminar.
