# P2c · S5-terminos (11-oct-2026)

## Hecho
- HECHO: Q2-64 = A aplicada. Se intentó la portada de los 3 sitios con WebFetch (3 de 12 solicitudes). Los 3 devolvieron un desafío anti-robot; no se leyó ninguna página legal.
- HECHO: `docs/estado/S5-terminos.md` escrito con 3 veredictos «no verificado», evidencia del bloqueo y propuesta.
- HECHO: sección «Resultado» añadida a `docs/estado/bloqueos/S5-terminos.md` (el bloqueo sigue abierto).

## URL solicitadas, en orden (22:28 CST)
| # | URL | Resultado |
|---|---|---|
| 1 | https://www.walmart.com.mx/ | Código HTTP no expuesto por WebFetch; contenido = «Verifica tu identidad / mantén presionado el botón» (anti-robot). Sin enlaces. |
| 2 | https://www.bodegaaurrera.com.mx/ | Igual: desafío anti-robot, sin enlaces. |
| 3 | https://www.sams.com.mx/ | Igual: desafío anti-robot, sin enlaces. |
No se reintentó ningún bloqueo, no se adivinaron URL legales, no se usó curl ni navegador.

## Verificación
- `grep -c "Veredicto" docs/estado/S5-terminos.md` = 3 (mecánico). OJO: S5-07 NO se cumple en el fondo: los 3 son «no verificado».
- El código no lee ese archivo (grep en *.go sin referencias); el interruptor sigue apagado.
- `make verify-all` no se corrió: solo se tocaron documentos de estado, sin código.

## Supuestos
- El mensaje «Verifica tu identidad» equivale a un bloqueo (captcha) según la regla; no se intentó una página legal distinta porque la regla dice registrar y pasar al siguiente sitio y los tres fallaron igual.

## Bloqueos
- `docs/estado/bloqueos/S5-terminos.md` sigue abierto: los sitios exigen verificación humana.

## Pendientes
- Veredicto real por sitio, tras decisión de Axel.

## Propuestas para Axel
- Opción 1 (recomendada, equivale a B): abre en tu navegador los términos de cada sitio y pega URL y cláusulas sobre acceso automatizado; un agente redacta el veredicto.
- Opción 2: darme URL legales exactas para un intento directo (probable mismo desafío).
- Opción 3/4: dejar pendiente (C) o descartar la lectura asistida (D).
- Mientras tanto: lectura asistida apagada; `site_assisted` solo referencia (K37).

## Consumo aproximado
Acciones: 8 · KB leídos: ~12 · KB escritos: ~7
