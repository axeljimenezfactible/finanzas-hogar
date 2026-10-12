# revisor-U4-r3 · P2c · reporte del revisor independiente, ronda 3 de 3

Agente `revisor-U4-r3` (Opus) · flujo P2c · 11-oct-2026, noche · árbol de trabajo sobre HEAD 4079cee. Informe completo: `docs/estado/p2c/REVISION-U4-r3.md`.

## Hecho
- Revisé U4 a 1.6, U4 b 1.5, el unido, U3 1.5 y el diff de DESIGN contra HEAD. Lo medí contra `docs/criterios/U4.md`, `docs/criterios/U3.md`, PUERTA-2 §2.4, §3 (Q2-11 y Q2-66 a Q2-70) y §5, con RESPUESTAS.md como regla: lo recomendado más excepciones, y no hay excepciones.
- **Veredicto: APROBADO, con 0 altos, 0 medios y 5 bajos** (B-1 a B-5).
- **Ronda 2:** M-1 (D-2, ahorro planificado) y M-2 (D-6, un decimal) quedan **cerrados**. B-1, B-2 y B-3 siguen abiertos.
- **Bajos nuevos:**
  - B-4: U3 l. 18 declara «C1 v1.2», pero usa C1 1.3.
  - B-5: el «—» de D-2 no lleva texto de motivo, y la regla de U3 l. 22 lo pide.
- No escribí ni corregí nada de lo revisado. No hice `git add` ni `git commit`. No toqué las carpetas de la Pista 3.

## Verificación
Comandos y scripts propios en `datos/tmp/rev-u4-r3/`, todos con `python3 -I`. Ya los borré: `ls datos/tmp | grep -c rev-u4-r3` → 0. Las salidas completas están en la tabla del informe; aquí va lo principal.

- **(a) R5-1.** La regla general está en a §0.2, puntos 1 a 3. Cada sección la aplica:
  - P-26.04: el ejemplo de $199.00, al soltar una tarjeta y con `Ctrl+/`.
  - P-26.14: estado antes → después en cada fila y el aviso «{k} pasan a Reponer…».
  - P-26.13 y P-26.02 también la aplican.
- **(b) Claves de la 00004.** Las 8 claves `reports.*` y `alerts.*` están en el Up de la 00004 y en U4 b con los mismos valores: "21:30", "19:00", "08:00", "09:00", 1, 7, 1 y 7. Ninguna clave de `app_setting` citada en U4 falta en las migraciones. C1 1.3 las pide así:
  - `ReportSettings`: las 5 claves `reports.*`.
  - `AlertSettings`: las 3 claves nuevas de `alerts.*`.
- **(c) Q2-66 a Q2-70.** Están como dicen las respuestas. Recalculé con `fractions`:
  - existencia: 1,486,571 milésimas;
  - días: piso(31,217,991 ÷ 946,000) = 32;
  - consumo por semana: 315,333 milésimas;
  - avance: 500 décimas, y 50 en entero en la tarjeta.
- **Q2-11 B.** Hay una sola variante. Las citas de `S1-parametros-avisos` en U4 dicen RESUELTO.
- **(d) DESIGN.** Cambian 4 líneas. Del lado + solo aparecen los números de «puerta 2», K51 y los ids Q2. Hex: 0. Lo que se quitó son marcas «PROPUESTA» y la frase condicional de Q2-84. `git diff --quiet 1ccb468 HEAD -- docs/design/DESIGN.md` no muestra diferencias.
- **(e) Unido.** `U == U[:2461] + A + "\n" + B` da True. Secciones: 7 + 7 = 14. Contra 212 operaciones (75 de Ajustes):
  - citadas que no existen: 0;
  - operaciones de Ajustes sin cita: 0.
- **Diff de U3 1.4 → 1.5** contra el respaldo `/tmp/correccion-U4-r2/U3-dashboard.md.bak`: solo cambian las líneas 3, 186, 190, 475 y 476 (nuevas), 478 y 479.
- **PROVISIONAL.** En U3 y en U4, cada «PROVISIONAL» cita un bloqueo abierto, una solicitud o K31 punto 3. Hay 0 casos en que el texto diga «RESUELTO» o «abierto» y el archivo diga otra cosa.
- **Global.** `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` da **rc=0**. Salidas:
  - «U4-ajustes.md termina con parte a + parte b sin cambios: True · secciones P-26: 14 de 14»
  - «RESUMEN TODO OK»
  - 21 archivos de C1 iguales byte a byte
  - oráculos S2 a S5 con 0 diferencias
  - «WF-1: verificación completa OK…» y «Verificación ampliada OK…»

## Supuestos
1. Cuento como cerrado M-2 porque D-6 cae dentro de «Metas»: cita la métrica M-16, y `S2-metas.md` l. 22 dice literalmente «un decimal en Metas (`progress_tenths_pct`) y entero en la tarjeta de Inicio (`progress_pct`)». No es una decisión mía: es la lectura de Q2-17 C que también usan C1 1.3 y S2-interfaz 1.3.
2. Los bajos B-1 a B-3 no tienen texto nuevo, porque los md5 de a, b y DESIGN no cambiaron desde la r2. Mantengo su severidad de la r2.

## Bloqueos
- No abrí bloqueos nuevos. Siguen abiertos los que citan U3 y U4: `S2-ahorro-quincenal`, `S2-flujo-quincenal` (Q2-12), `S2-metas-vencida`, `S3-iva-lectura`, `S3-formula-tarjeta-pdf`, `U4-reparto-preseleccion`, `C1-hora-reportes-job` y `C1-reparto-alcance`.

## Pendientes
Son para el integrador o el revisor de C1; no bloquean la aprobación de U4.
- **C1 1.3, `getCashFlowFortnights`.** `planned_savings_cents` está en `required` y `balance_cents` está descrito como «income − egresos − ahorro». Las dos cosas contradicen el bloqueo abierto `S2-ahorro-quincenal`. U3 §15 l. 476 ya lo pide.
- **C1 1.3, enum de `SavingsGoal.status`.** Todavía es `on_track` y `behind`; según Q2-18 A tiene que ser `in_progress`, `achieved` y `overdue`.
- **C1-datos M-16.** Dice «% de avance entero» y Q2-17 C pide un decimal.
- **`docs/estado/bloqueos/INDICE.md`.** Faltan `U4-reparto-preseleccion`, `C1-hora-reportes-job` y `C1-reparto-alcance`. A los dos de C1 les falta la línea «Estado: ABIERTO».
- **Bajos B-1 a B-5.** Si se corrige B-1 en la parte a, hay que regenerar el unido.

## Propuestas para Axel
- Contestar `S2-ahorro-quincenal` por id. Son tres preguntas sin recomendación (A o B en cada una). Mientras no las contestes, D-2 muestra «—» en el ahorro planificado.
- Llevar a la puerta 4 la sección de `ocr.timeout_seconds` (B-3): ¿«Datos y respaldos» o «Tiendas y precios»?

## Consumo aproximado
- 39 acciones de herramienta.
- Unos 95 KB leídos, por partes: REVISION r2, correccion r2, PUERTA-2 §2.4, §3.1, §3.3, §3.4, §3.22 a §3.27, §4 y §5, RESPUESTAS, criterios U3 y U4 (encabezados), U3 (l. 1 a 8, 100 a 125, 186, 190 y 468 a 485), U4 (secciones citadas), bloqueos, extractos de C1, C1-datos (M-02, M-10 y M-16), diff de DESIGN y U4-puerta2.
- Unos 21 KB escritos: este reporte, el informe y unos 6 KB de scripts temporales ya borrados.
