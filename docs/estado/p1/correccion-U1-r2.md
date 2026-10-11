# Corrección U1 · ronda 2 (`correccion-U1-r2`)

Flujo P1 · 11-oct-2026. Corrige el hallazgo A1 (severidad alta) de `docs/estado/p1/REVISION-U1-r2.md`. Esa revisión no tiene hallazgos de severidad media. Los bajos (B1 a B5) no se tocan, porque el encargo cubre solo altos y medios.

## Hecho
- **A1 · K18 en U1** (`docs/design/U1-arquitectura.md`, por parches): 10 cambios, 7 líneas reemplazadas y 3 nuevas.
  - **Lista de cambios, fila 17** (línea 38): K18 y K12, con todos los lugares donde se aplican.
  - **§2.2, P-18** (251): en la tarjeta bancaria y el crédito de tienda, el estado «Falta capturar · Capturar» (abre P-26.05 en esa línea) ocupa el lugar de los intereses y del pago mínimo.
  - **§2.3, P-26.05** (271): para cada tarjeta bancaria y crédito de tienda se capturan 5 datos: tasa anual tal como se capturó, base diaria o mensual, si incluye IVA, porcentaje de pago mínimo y mínimo fijo. Los captura Axel desde su estado de cuenta y no tienen valor por defecto. Mientras falte alguno, la línea dice «Falta capturar».
  - **§2.4** (289): estado común «Falta capturar» (K18). Ninguna pantalla muestra intereses ni pago mínimo de esa línea calculados con valores supuestos (RN-26).
  - **N-200** (629): la tasa se muestra como se capturó (K12). Los cinco datos se muestran como se capturaron, o «Falta capturar».
  - **N-201** (630): los intereses por pagar y el pago mínimo salen de la fórmula general de S3. Si falta un dato, las cifras que dependen de él dicen «Falta capturar · Capturar».
  - **N-205** (634): si falta un dato, el pago mínimo de la tarjeta dice «Falta capturar · Capturar» y no se calcula.
  - **N-283** (729): los cinco datos, sin valor por defecto:
    - tasa anual con sus decimales (K12);
    - base diaria o mensual;
    - si incluye IVA y, si no lo incluye, el IVA que se suma a los intereses;
    - porcentaje de pago mínimo;
    - mínimo fijo.
    También dice cuáles faltan.
  - **N-289** (735): edita los cinco datos, sin valor por defecto. Si se guarda la línea sin alguno, queda en «Falta capturar». La tasa anual se guarda con los decimales tecleados.
  - **§7** (1155): fila K18 «Decidido», con el bloqueo `criterios-credito-formulas` (resuelto) y RN-26.
- **Lo que no cambia:**
  - No hay N-xxx nuevas: siguen 203 (138 de lectura y 65 de escritura).
  - La matriz §3.29 queda igual. N-200, N-201 y N-205 ya estaban en la fila de P-18; N-283 y N-289, en la de P-26.05.
  - No hay cifras en $ nuevas.
- **`docs/design/tools/verificar_u1.py`** (líneas 15, 168 a 191 y 270 a 281): 7 comprobaciones nuevas de K18 y una línea PENDIENTE para C1.
  - **U1-09:** revisa P-26.05, N-283, N-289, P-18 (§2.2, N-200, N-201 y N-205), §2.4 y la fila de la lista de cambios.
  - **U1-12:** revisa las columnas de `credit` en S1 1.2: en 00002, `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents` y `annual_rate_decimals`; en 00001, `vat_on_interest_bp`.
  - **U1-03 (C1):** si `Credit`, `CreditCreate` o `CreditUpdate` de `openapi.yaml` no tienen los cuatro campos nuevos, imprime «PENDIENTE C1 v1.1», que no cuenta como falla. Cuando C1 v1.1 los agregue, pasa a OK.
- **Adenda «Para C1»:** va en la sección siguiente. El integrador la suma a `docs/estado/p1/U1-1.1.md` (solicitud 1).

## Para C1 (adenda K18 · corrección r2)
Necesidades que cambian el contrato:
- **N-200 y N-201** (`getCredit`) y **N-283** (`listCredits`):
  - En `bank_card` y `store_credit`, `Credit` expone `rate_base` (`daily` o `monthly`), `min_payment_pct_bp`, `min_payment_fixed_cents` y `annual_rate_decimals` (S1 1.2, 00002). Van junto a `annual_rate_num`/`annual_rate_den` y `vat_on_interest_bp`, que ya existen.
  - Cada campo puede venir vacío, que significa «falta capturar». Conviene que la respuesta diga cuáles faltan (K18, K12, RN-25 y RN-26).
- **N-289** (`createCredit` y `updateCredit`):
  - `CreditCreate` y `CreditUpdate` aceptan esos cuatro campos, opcionales y sin valor por defecto.
  - La base, el porcentaje y el mínimo fijo solo se aceptan en esos dos tipos (CHECK de 00002).
  - La tasa anual se guarda con los decimales tecleados.
- **N-201** (`getCredit`) y **N-205** (`getCardStatus`):
  - Si a una tarjeta bancaria o a un crédito de tienda le falta alguno de los cinco datos, los intereses por pagar y el pago mínimo llegan sin valor, nunca en cero, y con la indicación «falta capturar» (RN-26).
  - En `bank_card` y `store_credit`, N-201 también lleva el pago mínimo.
- **Regla general (§2.4):** cualquier otra respuesta que lleve intereses o pago mínimo calculados de esa línea sigue la misma regla: sin valor y con la indicación.
- **Si incluye IVA:** hoy `vat_on_interest_bp` NULL significa «la tasa ya trae IVA». Por eso la API no distingue «falta capturar» de «incluye IVA» (solicitud 2).

## Archivos
- `docs/design/U1-arquitectura.md` (modificado por parches; de 138,943 a 141,969 bytes; 1,161 líneas).
- `docs/design/tools/verificar_u1.py` (modificado; de 17,718 a 21,744 bytes).
- `docs/estado/p1/correccion-U1-r2.md` (este reporte, nuevo).

## Verificación
- **`python3 -I docs/design/tools/verificar_u1.py`:** sale con 0 y dice «RESULTADO: TODO OK». Da 52 OK (45 de antes y 7 de K18) y 0 FALLA. Lo demás no cambia:
  - 9 entradas de la barra lateral, 27 pantallas y 14 secciones de Ajustes;
  - 203 N-xxx (138 L y 65 E; la última es N-319);
  - F-1 = 3 toques;
  - 30 cifras en $ (21 distintas), todas del oráculo o de la semilla;
  - 77 nombres revisados, ninguno en inglés;
  - `openapi.yaml` 1.0.0 con 197 operaciones, 196 con `x-necesidades`.
  - Las cuatro marcas PROVISIONAL de antes ahora están en las líneas 66, 421, 1151 y 1152 (las inserciones las recorrieron 1 o 2 líneas). Cada una cita su bloqueo.
- **Líneas PENDIENTE**, esperadas hasta C1 v1.1:
  - `['N-316', 'N-318', 'N-319']`;
  - nueva: los cuatro campos de K18 en `Credit`, `CreditCreate` y `CreditUpdate`.
- **Prueba negativa:** en una raíz temporal (`datos/tmp/correccion-u1-r2/`, ya borrada) con el U1 de antes de esta corrección, salen 6 FALLA (las seis comprobaciones U1-09 de K18) y el script sale con 1. La comprobación de S1 da OK.
- **Prueba positiva:** en la misma raíz, con el U1 corregido y una copia de `openapi.yaml` que tiene los cuatro campos en los tres esquemas:
  - «U1-03 K18 … exponen …» sale OK;
  - la línea PENDIENTE de K18 desaparece;
  - el resultado es TODO OK.
- **`diff` contra la copia previa:** 10 bloques (7 líneas reemplazadas y 3 nuevas), solo en las líneas que lista «Hecho».

## Supuestos
- Ninguno cambia cifras, esquema, endpoints ni tokens. Todo sale de K18, K12, RN-25, RN-26 y de 00002 («NULL = falta capturar»).
- Que N-201 muestre el pago mínimo de la tarjeta y del crédito de tienda es una decisión de diseño de U1: P-18 no tenía otro lugar para el pago mínimo del crédito de tienda.

## Bloqueos
- Ninguno nuevo. Los cinco datos tienen lugar claro: se capturan en P-26.05 (N-283 y N-289) y se leen en P-18 (N-200, N-201 y N-205).
- Lo del IVA es un tema de representación técnica (solicitud 2), no una pregunta para Axel: K18 ya fijó el comportamiento.

## Pendientes
- **C1 v1.1:** aplicar la adenda «Para C1». Después, `verificar_u1.py` deja de imprimir el PENDIENTE de K18.
- **Hallazgos bajos B1 a B5 de la ronda 2:** no se tocaron, porque quedan fuera del encargo (solo altos y medios).

## Solicitudes fuera de alcance
1. **`docs/estado/p1/U1-1.1.md`, sección «Para C1»** (en «Cambiadas, con cambio de contrato»).
   - Cambio: sumar las entradas de «Para C1 (adenda K18 · corrección r2)» de este reporte.
   - Motivo: lo pide A1 de `REVISION-U1-r2.md`. Sin ellas, C1 v1.1 no tiene ninguna necesidad que exponga las columnas K18 de S1 1.2.
2. **`backend/db/migrations/00002_puerta1.sql` (S1) y `docs/contratos/reglas-negocio.md` (RN-26).**
   - Cambio: que «si incluye IVA» pueda quedar sin capturar en `bank_card` y `store_credit`. Hay dos opciones:
     - una columna nueva 0/1 en la que NULL = falta capturar;
     - en esos dos tipos, que `vat_on_interest_bp` NULL signifique falta capturar, 0 que la tasa ya trae IVA y 1600 que se suma 16 %.
   - Motivo: K18 y RN-26 piden «falta capturar» si falta cualquiera de los cinco datos. Con `vat_on_interest_bp` NULL = «ya trae IVA», ese dato nunca puede faltar.
   - Lo decide S1 con el integrador.

## Paquetes instalados
Ninguno.

## Consumo aproximado
- Acciones de herramienta: 26, contando este reporte y la salida final.
- Leído: unos 135 KB (U1 solo por rangos y grep; `openapi.yaml` lo leyeron los scripts).
- Escrito: unos 28 KB (parches a U1 de unos 3 KB, al script de unos 4 KB, scripts de parche en el scratchpad de unos 13 KB y este reporte de unos 8 KB).
