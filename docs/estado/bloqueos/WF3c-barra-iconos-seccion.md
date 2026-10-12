# Bloqueo WF3c-barra-iconos-seccion · Iconos de sección de la barra lateral (G-01)

Abierto por `front-tablero` (WF-3c, 12-oct-2026) al cerrar el hallazgo M-7 · U5-06 #10 de `docs/estado/wf3/REVISION-u5-r3.md`.

**Qué falta.** DESIGN §7.9 pide un icono `size.icon.sm` por cada una de las 9 secciones y §5.4 solo nombra el de 5: Inicio `house`, Comidas `utensils`, Compras `shopping-cart` («barra inferior y lateral»), Créditos `credit-card` («Deuda / crédito … Créditos») y Ajustes `settings`. No nombra el de **Gastos**, **Presupuesto**, **Salud** ni **Reportes**.

**Qué se hizo (PROVISIONAL)** en `frontend/src/shell/BarraLateral.tsx` (`ICONO_SECCION`): el icono de la sección de Ajustes que corresponde, que ya está en §5.4: Presupuesto `wallet` (Ajustes › Presupuesto), Salud `heart-pulse` (Ajustes › Salud y metas de ahorro), Reportes `file-chart-column` (Ajustes › Reportes y frecuencias); Gastos `receipt` (el ejemplo de sprite de §5.2).

**Opciones.** A) Confirmar estos 4 y que DESIGN §5.4 los agregue con «barra lateral». B) Gastos `arrow-up-right` (el significado «Gasto» de §5.4); el resto igual. C) Otros que elija Axel o DESIGN.

**Impacto.** Solo el dibujo de 4 iconos (20 px, trazo 1.8, `ui.icon.default` y `ui.icon.selected`); el ancho de 240 px y el resto de la barra no cambian.
