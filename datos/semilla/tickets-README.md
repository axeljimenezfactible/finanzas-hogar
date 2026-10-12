# Tickets para medir el OCR (S5-06, K24)

Este archivo vive **fuera** de `datos/semilla/tickets/` para que esa carpeta guarde solo las 5 fotos y `ls datos/semilla/tickets/ | wc -l` dé 5 (S5-interfaz §12, criterio S5-06; REVISION-S5-r1 H-08). La carpeta trae un `.gitkeep` oculto para existir en git; `ls` no lo cuenta.

**Estado (11-oct-2026): faltan las fotos de Axel. S5-06 queda «sin medir».**

## Qué necesito de ti, Axel
Fotografía **5 tickets reales** de compra (Bodega Aurrera y Walmart; si puedes, 3 de una tienda y 2 de la otra) y guárdalos en `datos/semilla/tickets/` con estos nombres (nada más dentro de esa carpeta):

`01.jpg`, `02.jpg`, `03.jpg`, `04.jpg`, `05.jpg` (también sirve `.jpeg` o `.png`; el número es lo que importa).

Cómo tomar la foto (el OCR rinde mucho mejor así):
- Ticket extendido y plano sobre una mesa; sin dobleces sobre los renglones.
- Luz pareja, sin sombra de tu mano ni del teléfono, sin flash directo (el papel térmico brilla).
- Toda la hoja en el cuadro, de arriba a abajo, derecha (el programa endereza un poco, pero no milagros).
- Si el ticket es largo y no cabe con letra legible, parte la foto en dos **solo** si lo avisas; por defecto, una foto por ticket.
- El ticket real trae datos tuyos (tarjeta, folio). Si no quieres que queden en el repositorio, tapa esa parte con un dedo o recórtala; no afecta la medición de los renglones.

## Qué se transcribe a mano (antes de medir)
Por cada foto, un archivo **fuera de la carpeta de fotos**: `datos/semilla/tickets-esperado/NN.csv` (mismo `NN`). Va fuera para que `datos/semilla/tickets/` cuente solo las 5 fotos.

```
descripcion,cantidad,importe
LECHE ENTERA 1 L,,28.50
ARROZ BLANCO 1 KG,2,58.00
JITOMATE,0.750,22.13
```

Reglas (congeladas, no se aflojan después de medir):
1. Una fila por producto, tal como está **impreso** en el ticket. Sin código de barras, sin la letra de IVA ni totales.
2. `cantidad`: solo si el ticket la imprime (por ejemplo `2` o `0.750`); vacía si no la imprime.
3. `importe`: el importe impreso del renglón, con punto decimal y sin `$` ni comas (`1234.50`).
4. Un renglón sale **correcto** si la descripción coincide (sin importar mayúsculas, acentos, espacios repetidos ni puntuación al final) **y** la cantidad **y** el importe son exactos.
5. Se necesita **80 %** o más de renglones correctos sumando los 5 tickets.

## Cómo medir cuando estén las fotos
```
cd /home/deepsec/finanzas-hogar/backend
go run ./internal/precios/ocr/cmd/ocrmedir
```
Imprime la tabla ticket × correctos/esperados y el porcentaje. Si no hay 5 fotos o falta una transcripción, imprime «sin medir» y no declara nada cumplido.

## Contenido de esta carpeta
- `datos/semilla/tickets-sinteticos/`: 2 imágenes **sintéticas** (generadas, no son fotos reales ni imitan el formato de una tienda) y su transcripción en `tickets-sinteticos/esperado/`. Solo prueban la herramienta (`integration_test.go`). **No cuentan** para S5-06 y viven fuera de `datos/semilla/tickets/`.
- Verificación literal de S5-06: `ls datos/semilla/tickets/ | wc -l` = 5 cuando estén las fotos (hoy 0).
