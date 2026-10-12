# Bloqueo S5-precio-vigente-captura

Lo abrió el agente `correccion-motor-S5-r2` el 11-oct-2026, en la ronda P2c. Lo abrió al corregir M-1 de `docs/estado/p2c/REVISION-motor-S5-r2.md` (`created_at` con zona horaria). Se decide en WF-3, o antes si contratos o Axel contestan por número.

**Contexto.**
- `S5-interfaz.md` 1.5, §6 (supuesto S-8): si dos precios que compiten tienen la misma fecha, rige el de `created_at` más reciente.
- §15.1: `created_at` es un «instante RFC 3339».
- M-1 ya quedó corregido. Los dos lados comparan el instante en UTC y rechazan con `invalid_date` lo que RFC 3339 §5.6 no admite: zona `+24:00` o `+05:60` y hora `24`.
- Lo fijan cuatro fixtures nuevas: `s5-precio-vigente-zona`, `s5-ingesta-zona-hora`, `s5-ingesta-zona-minuto` y `s5-ingesta-hora-24`.
- Al medir salieron cuatro puntos que el contrato no fija. Ninguna fixture los usa, así que `compara.py s5` da 0 diferencias. Ninguno mueve una cifra canon.

**Regla mientras no se conteste:** PROVISIONAL, **en uso: el motor**, igual que en `S5-motor-oraculo-reglas`.
- No elegí qué lado tiene razón.
- El código cita este bloqueo con el número de la pregunta:
  - `backend/internal/precios/pricing.go` (`later`, pregunta 1);
  - `backend/internal/precios/eval.go` (lectura de `created_at`, preguntas 2 y 3);
  - `oraculo/pre_calculo.py` (`vigente`, pregunta 1);
  - `oraculo/pre_comun.py` (`_RFC` e `instante`, preguntas 2 y 3).
- Cuando se fije una regla, se alinea el lado que no la cumpla y se agrega una fixture con el caso. Las cifras de abajo sirven de base.

**Cómo medí.**
- Sondas sin `esperado` en `datos/tmp/corr-s5-r2/fx`, `fx2` y `fx3` (fuera de git), corridas con `motores s5 dump` y con `pre_run_fixtures.py`.
- En cada sonda compiten dos precios de la misma presentación y la misma tienda, del 2026-10-10.
  - Salvo que se diga otra cosa, son `manual` a 2000 con `created_at` 12:30:00Z y `ticket_manual` a 2100 con el `created_at` que se indica.
  - La consulta es `current_price`.

## 1. Desempate con el mismo `created_at` (§6, supuesto S-8): de negocio o técnica

| Precios, en el orden de la entrada, todos con `created_at` 2026-10-10T12:30:00Z | Motor | Oráculo |
|---|---|---|
| `manual` 2000, `ticket_manual` 2100 | 2000 `manual` | 2100 `ticket_manual` |
| `ticket_manual` 2100, `manual` 2000 | 2100 `ticket_manual` | 2000 `manual` |
| `ticket_ocr` 2200, `manual` 2000 | 2200 `ticket_ocr` | 2000 `manual` |
| `manual` 2000, `ticket_ocr` 2200 | 2000 `manual` | 2200 `ticket_ocr` |

- Con el mismo instante, el motor deja el **primero** de la entrada y el oráculo el **último**.
- En el oráculo la misma comparación elige también `promotion` y ordena `references[]` antes de ordenarlas por fecha y fuente.

**Opciones:**
- **A.** El primero de la entrada. Es lo que hace hoy el motor. En la API dependería del orden en que lea `price`, y §15.1 no lo fija.
- **B.** El último de la entrada. Es lo que hace hoy el oráculo. Si la API lee por id ascendente, es la última fila insertada.
- **C.** Un criterio que no dependa del orden de la entrada, por ejemplo el menor precio o un orden fijo de fuentes.

**Impacto:** solo cuando dos precios que compiten tienen la misma presentación, tienda y fecha, y además el mismo instante de captura. Si la API guarda `created_at` al segundo, dos capturas en el mismo segundo empatan. No medí con qué resolución lo guarda la API.

## 2. Precisión de la fracción de segundo de `created_at` (§15.1): técnica, para contratos

| `created_at` del `ticket_manual` contra 12:30:00Z del `manual` | Motor | Oráculo |
|---|---|---|
| 2026-10-10T12:30:00.0000001Z (7 dígitos) | 2100 `ticket_manual` | 2100 `ticket_manual` |
| 2026-10-10T12:30:00.0000000001Z (10 dígitos) | 2000 `manual` | 2100 `ticket_manual` |

- El motor guarda hasta nanosegundos. Con 10 dígitos los dos instantes le quedan iguales y cae en la pregunta 1.
- El oráculo toma la fracción completa y exacta, así que el `ticket_manual` es posterior.

**Opciones:**
- **A.** Exacta, con cualquier número de dígitos. Es lo que hace hoy el oráculo.
- **B.** Hasta nanosegundos, cortando lo demás. Es lo que hace hoy el motor.
- **C.** §15.1 limita la fracción (por ejemplo a 9 dígitos, o a la que genere la API) y lo demás es `invalid_date`.

**Impacto:** ninguno mientras la API genere `created_at` con 9 dígitos de fracción o menos.

## 3. Perfil de RFC 3339 para `created_at` (§15.1): técnica, para contratos

RFC 3339 deja abiertas tres cosas:
- §5.6 permite «t» y «z» en minúscula, y la especificación que lo usa puede exigir mayúsculas.
- Una nota de §5.6 deja que una aplicación use un espacio en lugar de «T».
- §5.7 admite el segundo 60 (segundo intercalar).

Los dos lados rechazan las cuatro formas con `invalid_date` (sondas `p-t-minus`, `p-z-minus`, `p-espacio` y `p-seg60`). Aquí **no** hay divergencia; falta que el contrato lo diga.

**Opciones:**
- **A.** §15.1 fija el perfil en uso: «T» y «Z» mayúsculas, segundos 00-59 y zona «Z» o ±hh:mm con hora 00-23 y minuto 00-59.
- **B.** Aceptar todo lo que admite RFC 3339. En ese caso cambian los dos lados.

**Impacto:** ninguno mientras la API genere `created_at` con «T», «Z» u offset, y sin segundo 60.

## 4. Claves obligatorias de `prices[]` (§15.1): técnica, para contratos

| Precio (`manual` 2000 del 2026-10-10) | Motor | Oráculo |
|---|---|---|
| sin `created_at` | 2000 `manual` (lo acepta; `eval.go` deja el instante en cero) | `missing_field` |
| `created_at` `null` | 2000 `manual` | `invalid_date` |
| sin `is_promotional` | 2000 `manual` (lo acepta como regular) | `missing_field` |

- La fila `prices` de §15.1 solo dice «(o `null`)» de `source_ref`. No dice qué claves internas son obligatorias.
- Un precio sin `created_at` en el motor pierde siempre el desempate S-8.

**Opciones:**
- **A.** Son obligatorias las 7 claves salvo `source_ref`. Es lo que hace hoy el oráculo.
- **B.** `created_at` e `is_promotional` son opcionales, y el contrato fija su valor cuando faltan. Es lo que hace hoy el motor, con los valores que el contrato todavía no fija.

**Impacto:** ninguna fixture omite esas claves. Falta confirmar si la API siempre las manda.

## Cómo seguí
- Corregí M-1 en los dos lados, cada uno con su marca:
  - el oráculo, escrito desde §15.1 y RFC 3339 §5.6 antes de leer el Go;
  - el motor, para que cumpla los rangos de la zona.
- Agregué las cuatro fixtures de arriba. Ninguna depende de este bloqueo: no hay pares que empaten, la fracción es de 3 dígitos y las letras van en mayúscula.
- Las cuatro preguntas quedan abiertas.
