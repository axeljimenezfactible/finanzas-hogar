# Bloqueo U4-dia-mes-reporte
Contexto: agente U4b-ajustes (WF-1, 11-oct-2026), Ajustes › Reportes y frecuencias (P-26.08). Por K44 el día del mes del reporte mensual (`day_of_month`, valor inicial 1, decisión F) es configurable, y la validación de U4 solo puede ser de tipo y rango técnico (U4-12).
Qué falta: si Axel cambia el día del mes del reporte mensual a 29, 30 o 31, ¿qué pasa en los meses que no tienen ese día (febrero, abril, junio, septiembre, noviembre)?
Opciones:
- A. Solo se aceptan los días 1 a 28. Consecuencia: el reporte sale todos los meses el mismo día; la validación lo dice («Elige un día del 1 al 28»).
- B. Se aceptan 1 a 31 y, en un mes más corto, sale el último día de ese mes. Consecuencia: el reporte nunca se salta; S7 calcula el último día.
- C. Se aceptan 1 a 31 y, en un mes más corto, no sale ese mes. Consecuencia: puede haber meses sin reporte automático (sigue el «Generar ahora»).
Impacto: ninguna cifra cambia; con el valor inicial (1) no hay diferencia. Afecta la validación del campo «Día del mes» en P-26.08 y el programador de S7.
Cómo seguí: `docs/design/U4-ajustes-b.md` valida «entero de 1 a 31» (rango técnico de un día del mes) y marca PROVISIONAL, citando este bloqueo, lo que pasa en los meses sin ese día. Nada entra en cifras ni en el esquema.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-66 · Si el día del reporte mensual es 29, 30 o 31, ¿qué pasa en los meses que no tienen ese día?
- **Respuesta:** **B** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Se aceptan 1 a 31 y, en un mes más corto, sale el último día de ese mes.»
- **Consecuencia:** el reporte nunca se salta; S7 calcula el último día.
- **Qué cambia:** Validación del campo «Día del mes» en P-26.08 (`docs/design/U4-ajustes-b.md`: entero de 1 a 31 sin marca PROVISIONAL en meses cortos) y programador de S7 (último día del mes corto).
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
