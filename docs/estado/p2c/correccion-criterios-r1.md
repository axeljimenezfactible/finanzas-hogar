# Corrección de criterios · ronda 1 (P2c) · `correccion-criterios-r1`

Fecha: 11-oct-2026, noche. Alcance: `docs/criterios/*.md`, `INDICE.md` y este reporte. Origen: `docs/estado/p2c/REVISION-criterios-r1.md` (H-1 alta, M-1 y M-2 medias). Las bajas B-1 a B-4 no se pedían y no se tocaron.

## Hecho
- **H-1 (alta) · U4-05.** El cuerpo ya no dice que «Solo a los gastos nuevos» sea requisito ni lo atribuye a Q2-68 = C: dice que Q2-68 C no fija la preselección y que esta es PROVISIONAL (bloqueo `docs/estado/bloqueos/U4-reparto-preseleccion.md`, abierto, se pregunta en la puerta 4; en uso «Solo a los gastos nuevos»). La Verifica pide «las dos opciones de alcance y la preselección marcada PROVISIONAL con cita del bloqueo». La lista «Cambios de la versión 1.3» de U4 lo registra. No decidí A, B ni C.
- **M-1 (media) · C1-15.** La Verifica agrega `S1-reparto-edicion-gasto` y deja dicha la regla general (Estado ABIERTO y sin «Respuesta de Axel»; la lista no es cerrada).
- **M-2 (media) · S5-13.** Separa las marcas del contrato S5 1.5 (las cuatro de siempre) de las del motor y el oráculo, que además citan `S5-motor-oraculo-reglas` y `S5-precio-vigente-captura` (abiertos, P2c); anotado en «Cambios de la versión 1.3» de S5.
- **INDICE §6**: la lista de bloqueos que sostienen marcas agrega `S5-motor-oraculo-reglas`, `S5-precio-vigente-captura` (S5-13), `U4-reparto-preseleccion` (U4-05) y `S1-reparto-edicion-gasto` (C1-15). **INDICE §12** registra la corrección.
- Sin criterios nuevos ni quitados; ids y severidades iguales.

## Verificación
- Los cuatro bloqueos citados existen y no tienen «Respuesta de Axel» (`grep -c` = 0 en los cuatro); `U4-reparto-preseleccion`, `S1-reparto-edicion-gasto` y `S5-motor-oraculo-reglas` dicen «Estado: ABIERTO»; `S5-precio-vigente-captura` no trae línea de estado pero su texto dice que se decide en WF-3 (sin respuesta).
- Marcas del código comprobadas con `grep -n PROVISIONAL`: `pricing.go` l. 215 y `pre_comun.py` l. 41 y 58 citan `S5-precio-vigente-captura`; `pricing.go` l. 227 y `compare.go` l. 105 citan `S5-motor-oraculo-reglas`; `C1-datos.md` l. 1044 cita `S1-reparto-edicion-gasto`.
- `make verify-criterios`: «TOTAL 261 (alta 201 / media 60 / baja 0) | archivos: 18», «INDICE contra conteo propio: 0 diferencias», «residuos N-xx y rutas /tmp: 0 y 0; las 11 comprobaciones dan lo esperado». 321 citas, 0 sin fuente; 178 cifras canon, 0 diferencias.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: rc 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».

## Supuestos
- En S5-13 y C1-15 la nota «no son defectos» aplica solo a bloqueos con Estado ABIERTO y sin «Respuesta de Axel».
- `S5-precio-vigente-captura` se trata como abierto: no tiene sección de respuesta, como dice el revisor.

## Bloqueos
Ninguno nuevo.

## Pendientes
- Bajas B-1 a B-4 de `REVISION-criterios-r1` (B-1 y B-3 en criterios; B-2 para P2d al registrar K53; B-4 en U2, fuera de este alcance si se quiere ronda siguiente).

## Propuestas para Axel
Ninguna nueva; la preselección del diálogo de P-26.12 se pregunta en la puerta 4 (opciones A, B, C del bloqueo).

## Consumo aproximado
Acciones de herramienta: 12; KB leídos: ≈ 75; KB escritos: ≈ 5.
