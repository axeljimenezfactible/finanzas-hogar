# correccion-U4-r2 · Reporte final

Flujo WF-1 · 11-oct-2026 · Corrige los hallazgos R2-1 (alta), R2-2 y R2-3 (media) de `docs/estado/wf1/REVISION-U4-r2.md`. No toca los bajos (R2-4 a R2-11) ni reabre decisiones de §A–§K. Nada provisional entra en cifras: todo lo que depende de un bloqueo abierto lleva PROVISIONAL con la cita del bloqueo en la misma línea.

## Hecho
- **R2-1 · U4-05 (alta). Vista previa en los 4 casos que faltaban.**
  1. **P-26.13, perfil de salud** (`U4-ajustes-b.md`). La vista previa en línea (parte a §0.2) ahora cubre cualquier dato del perfil que entra a una fórmula, y se calcula en el navegador con RN-17:
     - Estatura: cambian el IMC con su clase OMS (K7: `weight_g` × 1000 ÷ `height_mm`², un decimal HALF_UP) y la meta sugerida.
     - Fecha de nacimiento, sexo, actividad, objetivo y los ajustes de la persona: cambia la meta sugerida, con desglose y un solo redondeo al final.
     - El IMC solo depende del peso y la estatura (RN-17), así que no aparece cuando cambia la fecha de nacimiento o el sexo. Es una precisión sobre la redacción del hallazgo.
     - Si falta un insumo, el lado sin cálculo dice «Falta {dato}».
     - Ejemplos tomados de RN-17: IMC 23.7 (normal), 2,555.5625 → 2,556 kcal y 2,556 → 2,056 kcal con un ajuste de 500. Se conserva el ejemplo de S4-07.
  2. **P-26.05, vista previa 9 (nueva): días de corte y de pago de la tarjeta.** Muestra el próximo corte y el pago de ese corte, antes → después, calculados en el navegador con el calendario del motor S3 (`S3-interfaz.md` §5.2). Lleva PROVISIONAL citando `S3-formula-tarjeta`, pregunta 4. La fecha que toca las dudas de esa pregunta dice «Vista previa pendiente»: cuando hoy es día de corte, o cuando el día no existe en el mes.
  3. **P-26.05, vista previa 10 (nueva): tasa anual del préstamo.** Muestra la tasa antes → después como se capturó (K12). Para el pago por periodo y los intereses dice «Vista previa pendiente» y lleva PROVISIONAL citando `S3-convencion-prestamo`.
  4. **P-26.05, vista previa 4 (reescrita).** Ahora rige tanto al completar las condiciones de tarjeta o tienda como al cambiar una ya capturada:
     - «Falta capturar» antes → después, con la regla de `missing_terms` (K18).
     - La tasa como se capturó.
     - El pago mínimo y los intereses dependen de R2-2.
  - El párrafo «Edición en el lugar» de P-26.05 y el de P-26.13 nombran los campos que guardan después de su vista previa.
- **R2-2 · U4-05 y U4-11 (media).** El pago mínimo y los intereses de la vista previa 4 llevan PROVISIONAL citando `S3-formula-tarjeta`, preguntas 1 a 3.
  - Mientras ese bloqueo no tenga respuesta, el bloque no muestra cifra: dice «Pago mínimo e intereses: vista previa pendiente».
  - Con la respuesta, se usarán `CardMinimumPayment` y `CardPeriodInterest` (S3), calculados en el navegador o con SC-U4-4.
  - El bloque sí afirma lo que ese bloqueo deja FIRME: si se paga el total a tiempo, los intereses son $0.00, y los meses sin intereses no cambian.
  - Se quitó la frase «según la fórmula general (K18)» como cifra firme.
- **R2-3 · U4-05 y U4-11 (media).**
  1. P-26.13 cita `S2-metas` en cinco lugares:
     - Pregunta 1, los meses que faltan, y pregunta 3, el estado, en la aportación necesaria.
     - Pregunta 2, el número de decimales del avance.
     - Pregunta 3 también en el tablero por estado y en la columna calculada.
     - Se reconoce que C1-datos M-16 da la forma «(monto meta − ahorrado) ÷ meses que faltan».
  2. El avance de la vista previa usa el formato de la columna (`progress_pct`, entero HALF_UP de C1) para que «Antes» sea igual a lo que ya se ve. El número de decimales lleva PROVISIONAL (`S2-metas` (2)). Ejemplo de RN-18: 50 % (50.0 % con un decimal).
  3. El tablero por estado lleva PROVISIONAL (`S2-metas` (3)) y muestra las dos opciones del bloqueo frente a las columnas de U1 §2.5.
  4. Se retiró la pregunta 2 de `U4-vista-previa-formulas.md`, que remite a `S2-metas` (1):
     - Su opción (a) equivale a la A de S2-metas y su (b) a la B.
     - Su (c), días ÷ 30.4375, se retira.
     - El redondeo al centavo es la regla general de dinero (protocolo §6).
     - Las preguntas 1 y 3 conservan su número.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/design/U4-ajustes-a.md`, versión 1.2: encabezado y P-26.05 (vistas previas 4, 9 y 10, y «Edición en el lugar»).
- `/home/deepsec/finanzas-hogar/docs/design/U4-ajustes-b.md`, versión 1.2: encabezado y P-26.13 (vista de base de datos, columna calculada de metas, vista previa y edición en el lugar). En «Cobertura de criterios», cambian las líneas U4-05 y U4-11.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U4-vista-previa-formulas.md`: se retira la pregunta 2. Lo pide el hallazgo R2-3 y es un bloqueo `U4-*` del propio artefacto (U4-11).
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-U4-r2.md`: este reporte.

## Verificación
Corrí `python3 -I <scratchpad>/u4r2/v.py /home/deepsec/finanzas-hogar`. El script vive fuera del repositorio y lee las dos partes, los extractos, `openapi.yaml` (solo la expresión regular), `C2-colores.json`, `DESIGN.md`, las migraciones 00001 a 00003 y `oraculo/model.json`. Salida:
- **U4-01:** 14 secciones, P-26.01 a P-26.14 en orden. A ninguna le falta «Vista de base de datos», «Edición en el lugar», «Vista previa del efecto», «Estados», «Atajos», «412 px» ni «1440 px».
- **U4-02:** 210 `operationId` en los extractos, los mismos 210 que en `openapi.yaml`. Se citan 100 en a y 45 en b, con 0 inexistentes. Cubre 73 de 73 operaciones de Ajustes.
- **U4-04 y U4-12:**
  - Las 15 claves `@param` de 00003 están en U4.
  - Las 19 claves de `app_setting` sembradas están citadas.
  - `send_time`, `advance_days`, `weekday` y `day_of_month` están presentes, con `S1-parametros-avisos`.
  - «Mínimo de seguridad» fuera de SUPUESTOS: 0.
- **U4-05:**
  - P-26.05 tiene las vistas previas 1 a 10.
  - R2-2: «según la fórmula general (K18)» aparece 0 veces. La vista previa 4 cita `S3-formula-tarjeta` (preguntas 1 a 3) con «vista previa pendiente».
  - La vista previa 9 cita la pregunta 4 de ese bloqueo y la 10 cita `S3-convencion-prestamo`.
  - P-26.13 cubre estatura (IMC), fecha de nacimiento y sexo.
  - Caso A/B: las series A y B de 14 valores aparecen literales y son iguales a `calculado_*`. Mínimos, disponible de hoy y despensa 1 están presentes; A − B es constante.
- **Cifras, con Decimal y ROUND_HALF_UP, 0 diferencias:**
  - IMC 23.7.
  - Meta: 1,648.75 × 1.55 = 2,555.5625 → 2,556; menos 500 → 2,056.
  - S4-07: 2,098.50 − 300 → 1,799.
  - RN-18: 50 % entero y 50.0 % con un decimal.
  - Cotización 3,094.50 × 25.8685 % = 800.50.
- **U4-09 y U4-11:**
  - 0 hex, 0 `rgb(` y 0 fechas ISO.
  - 21 tokens citados, 0 inexistentes.
  - Bloqueos citados que no existen: 0.
- **R2-3:**
  - `S2-metas` aparece 5 veces en P-26.13.
  - «con un decimal (RN-18» como valor firme: 0.
  - Citas de «`U4-vista-previa-formulas.md`, pregunta 2»: 0.
  - El tablero por estado lleva PROVISIONAL con `S2-metas` pregunta 3.
  - En el bloqueo U4, la pregunta 2 queda retirada y remite a S2-metas.
- **PROVISIONAL:**
  - Parte a: 9 líneas, todas con cita.
  - Parte b: 27 líneas. Las 5 sin cita en la misma línea son las mismas del revisor (26, 104, 209, 361 y 362, de R2-10/H-7, baja y fuera de esta ronda). Las líneas nuevas citan su bloqueo.
- **No aplican** pruebas Go, `compara.py` ni la compilación de tokens: U4 es un documento de diseño sin código. No se crearon temporales en `datos/tmp`.

## Supuestos
- Vista previa 9: «hoy» hace el papel de la fecha de compra del calendario propuesto en `S3-interfaz.md` §5.2. Va bajo PROVISIONAL, con la cita de `S3-formula-tarjeta` (4), y no cambia cifras ni contratos.
- Mientras `S2-metas` (2) siga abierto, el avance de la vista previa usa el formato entero de C1. Lleva PROVISIONAL con cita y no es una decisión.

## Bloqueos
- Ninguno nuevo. Esta corrección cita los bloqueos abiertos `S3-formula-tarjeta` (preguntas 1 a 4), `S3-convencion-prestamo` y `S2-metas` (1 a 3).
- `U4-vista-previa-formulas` queda con 2 preguntas vigentes, la 1 y la 3.
- Siguen abiertos `S1-parametros-avisos`, `U4-iva-credito`, `U4-dia-mes-reporte` y `U4-reparto-retroactivo`.

## Pendientes
- Hallazgos bajos de `REVISION-U4-r2.md` que no se asignaron a esta ronda:
  - R2-4: vista previa 7 de P-26.05 con `first_payment_date`.
  - R2-5: restricciones alimentarias editables en P-26.07 y P-26.13.
  - R2-6: línea 3 de la parte b, «§7 aún no existe». Su segundo punto, U4-11 sin `U4-vista-previa-formulas`, quedó resuelto de paso al reescribir esa línea.
  - R2-7: P-26.01, «Tus otras sesiones siguen abiertas».
  - R2-8: vistas que se apartan de U1 §2.5.a.
  - R2-9: `/promocion` y `/reporte`.
  - R2-10: tabla consolidada de la parte b.
  - R2-11: formatos de a y b.
- Al responderse `S3-formula-tarjeta`, `S3-convencion-prestamo` y `S2-metas`, hay que quitar «Vista previa pendiente» y la marca de P-26.05 (vistas previas 4, 9 y 10) y de P-26.13.

## Solicitudes fuera de alcance
- **SC-U4-4 · C1, condicional** (`docs/contratos/openapi.yaml`, vía `c1gen`). Aplica cuando se respondan `S3-formula-tarjeta` y `S3-convencion-prestamo`, y solo si el pago mínimo, los intereses o la tabla del préstamo no se pueden calcular en el navegador porque necesitan el saldo diario o el saldo del corte. Pide una vista previa sin guardar que reciba las condiciones tecleadas y devuelva:
  - Condiciones de entrada: tasa, base, IVA, porcentaje y mínimo fijo de pago mínimo, día de corte y día de pago.
  - Respuesta: el pago mínimo, los intereses, el próximo corte y su pago, y el pago por periodo del préstamo.
  - Motivo: U4-05, vistas previas 4, 9 y 10 de P-26.05.
- **SC-U4-2 (reiterada y acotada):**
  - Para la meta, queda condicionada a `S2-metas` (1) y (3).
  - Para el inventario, queda condicionada a `U4-vista-previa-formulas` (3).
- **`docs/estado/bloqueos/S2-metas.md` («Impacto»), al integrador:** agregar U4 P-26.13. Ahí dependen del bloqueo la vista previa de la aportación necesaria, el formato del avance y el tablero por estado.
- **`docs/design/U1-arquitectura.md`:** marcar PROVISIONAL citando `S2-metas` (2) y (3), o alinear cuando haya respuesta. Afecta:
  - Lugares: la fila `meta_ahorro` de la tabla de vistas (línea 353), la fila P-09 de la tabla de páginas (línea 224), N-025 y N-111.
  - Contenido: los estados «en ruta, atrasada o cumplida» y el «% de avance (entero)».
  - Motivo: es la misma inconsistencia de R2-3.
- **`docs/estado/bloqueos/S3-formula-tarjeta.md` y `S3-convencion-prestamo.md` («Impacto»), opcional:** agregar U4 P-26.05. Las vistas previas 4 y 9 dependen del primero y la 10 del segundo.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- No hay valores K31 en esta corrección.
- La recomendación de `correccion-U4-r1` para la antigua pregunta 2 de U4, opción (a), equivale a **`S2-metas` (1) opción A**: del mes de hoy al de `target_date`, los dos incluidos. Es simple y se calcula en el navegador sin una operación nueva.

## Consumo aproximado
34 acciones de herramienta, unos 170 KB leídos y unos 42 KB escritos. Lo escrito incluye los scripts de edición y de verificación del scratchpad, unos 26 KB, y el contenido nuevo de las dos partes, el bloqueo y este reporte, unos 16 KB.
