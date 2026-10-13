# Bloqueo U3-tablero-k4-personas · K-4 «Pagos próximos» cuando `listPersons` aún no responde o falla (U3 §2 K-4, C3-09)

Fecha: 2026-10-12 · Agente: `correccion-frontend-cierre-r1` (Opus, WF-3a cierre) · Origen: hallazgo M-2 de `docs/estado/wf3/REVISION-frontend-cierre-r1.md`. Ninguna cifra cambia.

## Qué falta (HECHO)
- K-4 (U3 §2) pinta por fila «Hogar [parte] · Karen [parte]» y el estado del pago de U3 §1.4 («Falta la parte de Karen» / «Falta la parte del hogar»). Las dos cosas necesitan saber quién es Karen: la persona con `participation` externa de `listPersons`.
- Antes, mientras `listPersons` no respondía, la fila decía un instante «Hogar $973.75» (todo al hogar), una cifra distinta del oráculo (R-13: Hogar $598.82 · Karen $374.93).
- U3 dice qué hace K-4 si falla la lectura de N (`getDashboardSettings`, §1.2: estado de error del bloque) pero no dice qué hace si falla `listPersons`.

## Qué se hizo (PROVISIONAL)
- Mientras `listPersons` no llega: esqueleto de filas (no se pinta ninguna fila).
- Si `listPersons` falla o nunca se cargó sin conexión: el estado de error de K-4 («No se pudieron cargar los pagos. Revisa que estés en la red de casa.» · «Reintentar»), igual que cuando falla N; «Reintentar» vuelve a pedir lo que falló.
- Código: `frontend/src/pantallas/inicio/Tablero.tsx` (`K4`); prueba: `frontend/pruebas/inicio/tablero.escritorio.spec.ts`.

## PREGUNTA
¿Qué muestra K-4 si falla `listPersons`?

| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. PROPUESTA (aplicada): estado de error del bloque con «Reintentar» | Nunca una cifra distinta del oráculo; misma conducta que con N | Ninguno | Bajo: se pierde la lista mientras falla una lectura secundaria |
| B. Las filas sin la línea «Hogar · Karen» y con el estado de pago sin las variantes «Falta la parte de…» | Se ven fecha, línea y monto | Bajo | Medio: un pago con una parte ya pagada diría «Pendiente» en vez de «Falta la parte de Karen» |

Impacto: solo el caso de falla de `listPersons`; con la lectura bien, K-4 no cambia. La misma conducta falta en D-3 (`Detalle.tsx`) y en el Calendario (`Calendario.tsx`), que usan `karen(c)` igual; quedó como solicitud para el agente de Inicio. Decide Axel en la puerta 4 o el orquestador.
