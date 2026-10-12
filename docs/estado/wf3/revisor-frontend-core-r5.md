# revisor-frontend-core-r5 · WF-3a · 12-oct-2026

Revisor independiente, ronda 5. No corregí nada y no hice `git add` ni `git commit`. El informe completo está en `docs/estado/wf3/REVISION-frontend-core-r5.md`.

## Hecho

- **Veredicto: NO APROBADO.** Hay 0 altos y 1 medio (R5-01). Además hay 4 bajos nuevos: R5-02, R5-03, R5-04 y R5-05.
- **Hallazgos de la ronda 4.** R4-01 y R4-02 (los dos medios) quedan **cerrados**, con medición propia. R4-03, R4-04, R4-05, H-07, H-08 y H-09 (bajos) siguen abiertos; no se encargaron.
- **Verificaciones que reejecuté con mis comandos:**
  - build y tsc;
  - rutas contra el contrato y contra `routes.go`;
  - hosts externos en `dist` y en ejecución;
  - grep de hex y de medidas;
  - barrido de 36 vistas × 2 combinaciones: inglés, ISO, dinero tabular y formato, SVG de gráficas, fondo del tema y tarjetas;
  - fechas con formato híbrido;
  - objetivos táctiles, iconos y foco;
  - F-1 en 3 toques;
  - cola U5-13, pasos 1–3 y 6–7;
  - interruptor de tema;
  - estados de U1;
  - P-03 (N-102 y N-053);
  - la suite e2e del grupo;
  - `make verify-all`.

## Archivos

- `docs/estado/wf3/REVISION-frontend-core-r5.md`: veredicto, tabla, hallazgos y estado de los anteriores.
- `docs/estado/wf3/revisor-frontend-core-r5.md`: este reporte.
- Temporales en `datos/tmp/rev-fcore-r5/`: build, configuración de Playwright, pruebas propias y `rutas.py`. Ya están borrados.

## Verificación

- **C3-03.**
  - `npm run build -- --outDir ../datos/tmp/rev-fcore-r5/dist` → exit 0 (`tsc --noEmit` limpio, «precache 81 entries»).
  - `rutas.py`: 0 rutas fuera del contrato y 0 métodos no declarados en 208 llamadas.
  - `node scripts/rutas-api.mjs` → «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)».
- **U5-08.** Hosts en ejecución: solo `localhost:4895`. `hosts-externos.mjs` → «hosts externos que se pedirían = 0».
- **U5-06.**
  - 498 controles, 0 < 44 px sin una zona de 44.
  - 222 iconos, 0 fuera de norma (trazo dibujado de 1.5).
  - Foco: `outline 2px`.
  - Dinero sin tabular: 0, también en el SVG.
  - Eje compacto y etiquetas completas en P-01, P-02, P-06, P-07, P-08, P-09 y P-10.
  - 0 hex y 0 `dark:`.
  - 6 tarjetas en Inicio.
  - Los componentes de §7 se usan en las 3 áreas.
- **U5-07.**
  - 0 términos en inglés en unos 170 buscados, en 36 vistas × 2 combinaciones.
  - 0 fechas ISO.
  - **1 vista con fecha híbrida, P-10 Tendencias («Lun 5 oct-2026») → R5-01.**
- **U5-09.** 3 toques en línea y 3 sin conexión. «$97.40» aparece en P-02.
- **U5-10.** Fondo claro `rgb(250,250,251)` y oscuro `rgb(34,35,37)` (`ui.surface.page`). El sistema y el interruptor funcionan: claro, oscuro y «Según el sistema».
- **U5-13 (pasos 1–3 y 6–7, mocks).**
  - Se ven «Sin conexión con el servidor de casa» y «2 en cola», también tras recargar.
  - Con un 500, aparece el aviso y la captura sigue en la cola.
  - Con el reintento queda 1 gasto y 1 peso, y los 2 POST de gasto llevan 1 sola clave.
  - `captured_at` lleva desplazamiento y es la hora de captura.
- **U1 §2.2.** Funcionan los estados de carga, error en 13 vistas, vacío en 3 vistas y sin conexión.
- **U1 §3.5.** P-03 tiene N-102 y N-053 en la misma hoja.
- **e2e del grupo.** 111 passed, 3 skipped (`Ctrl+/` en móvil) y 1 failed intermitente: `correccion-bis.escritorio.spec.ts:79`. Al reejecutarlo pasó 3/3, y 57/57 con `--repeat-each=3` → R5-05.
- **`make verify-all GOOSE=…`** → exit 0, «Verificación ampliada OK».

## Supuestos

- **SUPUESTO 1.** El formato «Lun 5 oct-2026» cuenta como incumplimiento de U5-07 y de CONVENCIONES §5, porque no es ninguno de los dos formatos admitidos. Lo clasifiqué medio y no alto: el script de U5-07 (ISO = 0, inglés = 0) pasa, y la fecha se lee sin ambigüedad.
- **SUPUESTO 2.** El fallo intermitente se clasifica bajo: pasó 6 de 6 veces al reejecutarlo, y en ese momento corría a la par otra suite de Playwright de otro revisor.
- **SUPUESTO 3.** Los SVG de P-15 Tendencias no se pintan con los datos falsos del área. Lo cerré por el código: es la misma `Grafica` de dinero, `P15Productos.tsx:135`.

## Bloqueos

Ninguno nuevo de mi parte. Sugiero que el dueño de R5-03 abra uno (ver «Preguntas para Axel» en el informe).

## Pendientes

- **`front-comidas-compras`: R5-01.** En `P10Plan.tsx:245`, usar `fechaConDia()` o `fecha()`.
- **`frontend-base`: R5-02.** Es decisión de Axel.
- **`front-gastos-presupuesto`:**
  - R5-03: es PREGUNTA y lleva bloqueo;
  - R5-05: hacer estable la prueba de `Ctrl+/`.
- **`front-inicio`: R5-04.** Color de la leyenda.
- **Integrador.** Lo que dejó pendiente `correccion-frontend-core-r4`: regenerar la matriz C3 con `GET /api/v1/credits` y `POST /api/v1/incomes` en P-03, y recompilar el binario de `https://finanzas.local:8443` para ver los cambios en la LAN.

## Propuestas para Axel

- **PROPUESTA.** Que la suite fije los formatos de fecha de las gráficas con una aserción en cada área. Así R5-01 no habría pasado 4 rondas sin detectarse: el barrido de U5-07 solo busca ISO.
- **PROPUESTA (R5-02 y R5-03).** Decidir qué líneas ofrece N-053 y qué hace «Guardar» con «A crédito».

## Paquetes instalados

Ninguno.

## Consumo aproximado

- Unas 50 acciones.
- Unos 150 KB leídos: criterios, convenciones, reportes, código por partes y salidas.
- Unos 45 KB escritos: pruebas propias temporales, script, informe y este reporte.
