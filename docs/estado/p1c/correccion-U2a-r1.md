# Corrección U2a · ronda 1 (flujo P1c, 11-oct-2026)

Agente `correccion-U2a-r1`. Corrige el único hallazgo asignado de `docs/estado/p1c/REVISION-U2a-r1.md`: **M-01 (media)**. Los bajos B-01 y B-02 no venían en el encargo y no se tocaron.

## Hecho
- **M-01 corregido.** El párrafo «Propuestas sujetas a la puerta 2» de §12 ya no incluye a Karen ni a Bodega Aurrera. Ahora dice, después de la lista: «**Ya no se proponen:** Karen en 333° y Bodega Aurrera en 356°, aceptadas por K41 en la puerta 1b; son FIRMES («Cambios 1.2»)». Los matices salen del modelo (`M.cat[...]["h"]`), no van escritos a mano. En la lista se quedan ingreso 135°, la resolución de colisiones de §7.5 (con gasto 45° de la 1.1) y la paleta de gráficas de la ronda 1. Esas tres no son tokens de `FIRMES_12` y el hallazgo no las menciona.
- **Candado agregado (la parte opcional de la corrección sugerida).** En `colorimetria.py` agregué la constante `MARCA_YA_DECIDIDAS = "**Ya no se proponen:**"` y la función `nombres_firmes_12()`, que da Karen, Bodega Aurrera y Walmart a partir de los tokens de `FIRMES_12` y de `PERSONAS`/`TIENDAS`. `verificar()` falla en dos casos: si el párrafo de §12 no aparece exactamente una vez con la marca, o si una persona o tienda con un token en `FIRMES_12` sale antes de la marca, es decir, como propuesta.
- Regeneré los artefactos con `colorimetria.py`. Solo cambia la línea 720 de `colorimetria.md`. `contraste.csv`, `C2-colores.json`, `paleta.html` y `paleta.png` salen idénticos byte a byte a como estaban antes de la corrección.

## Archivos
- `docs/design/tools/colorimetria.py`: modificado por parche con `str.replace` en 3 puntos: constante y ayudante tras `FIRMES_12` (líneas 374–385), párrafo de §12 (línea 1722) y candado en `verificar()` (líneas 2086–2096). Son 26 líneas de diferencia.
- `docs/design/colorimetria.md`: regenerado. Solo cambia la línea 720.
- `docs/design/contraste.csv`, `docs/contratos/C2-colores.json`, `docs/design/paleta.html` y `docs/design/paleta.png`: regenerados, sin cambio en su contenido.
- `docs/estado/p1c/correccion-U2a-r1.md`: este reporte.

## Verificación
- `python3 -I docs/design/tools/colorimetria.py`: salida 0, «VERIFICACIÓN OK». Escribe los 4 archivos de texto y `paleta.png` (en el intento 1). Resultados: contraste.csv con 202 filas y 0 con cumple=false; 190 tokens que coinciden con el Anexo A; 22/22 del plan §4.3; en claro y en oscuro, 0 colisiones con roles en personas y tiendas; tiendas en oscuro con ΔL mínimo de 0.101 en OKLab y en gris; «versión 1.2: 5 tokens FIRMES por K41 y K42, 0 con otro hex; marcas PROVISIONAL en el JSON: 0 (en el script: 0)». Línea nueva: «§12: personas o tiendas FIRMES por K41/K42 entre las propuestas de la puerta 2: 0 de 3 revisadas».
- `--solo-verificar`: salida 0, «VERIFICACIÓN OK».
- `make verify-color`: salida 0. Regenera en un árbol espejo y da «idéntico byte a byte» en `colorimetria.md`, `contraste.csv`, `C2-colores.json` y `paleta.html`. Da «contraste.csv del repositorio: 202 filas, 0 con cumple=false».
- Comparé contra lo que había antes de la corrección. Corrí el script previo, respaldado, en un espejo con `--sin-png`. `contraste.csv`, `C2-colores.json` y `paleta.html` salen iguales. En `colorimetria.md`, `diff` da solo `720c720`, el párrafo de §12.
- Contra HEAD: según `git diff --quiet HEAD`, `paleta.png`, `contraste.csv` y `paleta.html` son idénticos. Comparé los dos JSON aplanados con un script propio: 190 tokens en HEAD y 190 ahora, con los mismos nombres. `$value`, `$type` y `$extensions` tienen 0 diferencias. `$description` cambia solo en los 5 tokens de K41 y K42, y ese cambio viene de U2a-1.2, no de esta corrección.
- `grep -c PROVISIONAL` da 0 en `colorimetria.md`, `contraste.csv`, `paleta.html`, `paleta.png` (con `-a`) y `C2-colores.json`. En el script hay 19 coincidencias, todas del mecanismo; la corrección no agregó ninguna.
- Cálculo propio en OKLab con las matrices de Ottosson, a partir del hex del JSON. Karen `#B870AC` da 333.2°, L 0.6429 y C 0.1199; `#D389C6` da 333.3°, L 0.7251 y C 0.1196. Bodega Aurrera `#762C4C` da 356.1°, L 0.4096 y C 0.1095; `#FFAAC9` da 356.2°, L 0.8292 y C 0.1063. Walmart en oscuro `#7AA3EE` da 261.9°, L 0.7165 y C 0.1190. Todo coincide con la revisión y con «Cambios 1.2».
- Pruebas negativas del candado, cada una en un espejo en `datos/tmp`:
  - Con «Karen y Bodega Aurrera en sus matices nuevos» de vuelta en la lista: salida 1, «§12: ['Bodega Aurrera', 'Karen'] aparecen como propuesta de la puerta 2 pero sus tokens son FIRMES (FIRMES_12)».
  - Con el Markdown editado a mano para quitar la marca: salida 1, «se esperaba un párrafo «Propuestas sujetas a la puerta 2» con la marca **Ya no se proponen:**».
  - Con «Walmart en oscuro» metido en la lista: salida 1, «§12: ['Walmart'] aparecen como propuesta…».

## Supuestos
- Ninguno que cambie cifras, esquema, endpoints o tokens. Dejé «gasto 45°» en la lista de propuestas: K26 fija un rango de 40–45°, el valor exacto dentro del rango no está aceptado por separado y el hallazgo no lo pide.

## Bloqueos
- Ninguno.

## Pendientes
- Ninguno del encargo. Siguen abiertos los bajos de la revisión que no venían en el encargo: B-01, el encabezado de «Marcas que se quitan» sin «medidos en el hex», y B-02, el texto del refuerzo en las filas de K41 de §7.7 y en el punto 1 de la regla para `DESIGN.md`. Los dos se arreglan en `docs/design/tools/colorimetria.py` (`md_cambios_12()` y las cadenas `ref` del barrido de §7.7) si el orquestador los asigna.

## Solicitudes fuera de alcance
- Ninguna.

## Propuestas para Axel
- Ninguna.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones de herramienta: unas 19.
- Leído: unos 75 KB (protocolo 6, decisiones §K 15, revisión 11, fragmentos del script unos 25, reporte U2a-1.2 unos 4, salidas de comandos unos 14).
- Escrito a mano: unos 12 KB (parche unos 4, este reporte unos 6, pruebas negativas unos 2). Los archivos regenerados los escribe el script, no se editaron a mano.
