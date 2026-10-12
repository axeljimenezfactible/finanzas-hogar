# Bloqueo S6-presupuesto-periodos · cuándo se llena la tabla period

Agente: api-presupuesto, 12-oct-2026. PROVISIONAL. Ninguna migración inserta renglones en `period`. `listPeriods` crea a pedido (con `audit_log`) los meses y quincenas que falten del rango pedido; sin `period`, del mes en curso − 11 al mes en curso + 12; tope 36 meses (422). Opciones: a) así; b) que `migrate`/`semilla` llenen un rango fijo; c) un job mensual. Impacto: solo el selector de periodo.
