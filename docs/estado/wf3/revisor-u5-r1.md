# WF-3b · revisor-u5-r1 · Reporte

Agente `revisor-u5-r1` (revisor visual independiente, ronda 1 de 3) · 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo. No capturé ni implementé nada y no corregí la aplicación. **Veredicto: devuelto** (5 altos, 6 medios, 6 bajos).

## Hecho
- Reejecuté U5-01, U5-03 y U5-04 con comandos propios antes de opinar (U5-05); recalculé 20 pares de contraste con mi implementación (20 de 20, Δ 0.000).
- Corrí por mi cuenta, en un servidor de prueba aislado (copia de la base en `datos/tmp/`, `127.0.0.1:8445`, sin Telegram), F-1 (U5-09) y la captura sin conexión (U5-13), y repetí H-5 con la secuencia del capturista.
- Medí en el DOM las 41 pantallas × 4 combinaciones (axe, contraste, texto visible, hosts, tema y 11 reglas de DESIGN.md) y probé el interruptor de tema (U5-10).
- Lista visual U5-06 (14 elementos) leyendo las capturas: Resumen (P-01), Ajustes (P-26 y P-26.01 a P-26.14) y una pantalla por entrada de la barra (P-02, P-03, P-06, P-10, P-13, P-17, P-21, P-23, P-26), en las 4 combinaciones; además las gráficas de Tendencias.
- Reejecuté `c3-09.spec.ts` (C3-09).

## Archivos
- `docs/estado/U5-revision.md` (nuevo: «## Ronda 1»).
- `docs/estado/wf3/REVISION-u5-r1.md` (hallazgos con pantalla, combinación, archivo probable y corrección).
- `docs/estado/bloqueos/U5-tipografia-tres-tamanos.md` y `docs/estado/bloqueos/U5-nombres-propios.md` (nuevos).
- Este reporte.
- Efecto lateral: `frontend/pruebas-u5/flujos/resultados/c3-09.json` reescrito por la reejecución de su prueba (mismo resultado: 111 cifras, 0 diferencias). Mis inicios de sesión en `:8443` dejan sesiones nuevas en `datos/finanzas-dev.db` (solo lectura de datos).
- Temporales borrados: `datos/tmp/rev-u5/` (specs, servidor de prueba y su base, salidas, láminas) y `datos/tmp/rev-u5-py/`.

## Verificación
- **U5-01 · falla (alta)**: T = 41, 4T = 164; 162 de 164 correctos; P-26.13 Ac/Ao 415×1732 (mi captura: 415×1792; `innerWidth` = 415). `ls capturas-u5/ | wc -l` = 168 (bloqueo `U5-conteo-carpeta`).
- **U5-02 · pasa**: 8 cifras visibles = `model.json` (A $10,000.00; B $6,905.50; $973.75; $3,094.50; $6,293.44; $598.82; $374.93; $1,191.50).
- **U5-03 · falla (alta)**: comando literal = 20; mi corrida: Ac 5 · Ao 5 · Ec 5 · Eo 5 (P-06 y P-09 `aria-meter-name`, P-08 `button-name` crítica, P-26.05 y P-26.06 `nested-interactive`); 0 moderadas/menores.
- **U5-04 · falla (alta)**: pares únicos / bajo umbral: Ac 2,558 / 113; Ao 2,560 / 113; Ec 3,459 / 131; Eo 3,460 / 131; texto 0, íconos 0, todos los bajo 3:1 son de UI (botón y estado elegido 1.075:1); 20 pares al azar coinciden ±0.01 (20 de 20).
- **U5-05 · ronda 1 devuelta**: `grep -c "^## Ronda" docs/estado/U5-revision.md` = 1, «Veredicto: devuelto».
- **U5-06 · falla (media)**: 14 elementos; fallan tipografía (título 24 px), iconos (16 px), espaciado (2 y 40 px), densidad (P-17 Android), barra lateral (256 px y empalme) y gráficas (ejes «$0», etiquetas encimadas); pasan color, `tnum`, radios, sombras/bordes, saturación, ≤ 6 tarjetas (6), objetivos táctiles y foco.
- **U5-07 · falla (alta)**: script: 1,404 importes en formato, 0 fechas ISO, 0 palabras de la lista; capturas: «Choose File»/«No file chosen» (P-26.11) y «groceries»/«eating_out» (P-26.04) en las 4 combinaciones; «Sin conexión» y «2 en cola» en español.
- **U5-08 · pasa**: 1 host, `finanzas.local:8443`.
- **U5-09 · pasa**: 3 toques; el gasto aparece en la tarjeta de Inicio y en la tabla de P-02.
- **U5-10 · pasa**: fondo = C2 en 164 de 164 sin `data-theme`; interruptor 7 de 7.
- **U5-11 · no evaluado** (lo marca el integrador en C3-05).
- **U5-12 · pasa**: §Supuestos en los dos reportes; bloqueos `U5-*` en `docs/estado/bloqueos/`; nada aprobado con violaciones abiertas.
- **U5-13 · falla (alta)**: pasos 1 y 3–7 pasan (3 toques sin red, clave de 36 caracteres, `captured_at` con desplazamiento anterior al envío, contador a 0, 1 registro por captura, reenvío 200 con `Idempotent-Replayed: true` y mismo id, fallo 500 con aviso y cola intacta); paso 2 falla con la secuencia del capturista (pantalla de entrada tras recargar sin red), pasa si antes hubo navegaciones completas.
- **C3-09 · pasa**: `c3-09.spec.ts` → «111 cifras leídas en pantalla · diferencias != 0: 0».

## Supuestos
1. «Pantallas del Resumen» = P-01 Inicio en su subsección Resumen (el dashboard de U3); las otras subsecciones no están en las capturas de U5-01, salvo Tendencias, que revisé aparte por las gráficas.
2. «Una por entrada de la barra» = las 9 secciones de la barra lateral más el «+» de la barra inferior (P-03).
3. U5-04: mientras el bloqueo `U5-contraste-ui-controles` siga abierto, cuento los pares de UI como fallas (lectura literal; U5-12 impide aprobar con violaciones abiertas).
4. Un control de menos de 44 px cumple el objetivo táctil si su `<label>` mide ≥ 44 × 44 o tiene área ampliada con pseudoelemento.
5. «groceries» y «eating_out» son texto en inglés visible mientras Axel no los acepte como identificadores técnicos (bloqueo `U5-nombres-propios`).
6. El relleno de 64 px del contenedor puede ser la holgura de la barra inferior (medida de estructura); lo dejo dentro de M-4 para que se confirme.

## Bloqueos
- `docs/estado/bloqueos/U5-tipografia-tres-tamanos.md` (nuevo).
- `docs/estado/bloqueos/U5-nombres-propios.md` (nuevo).
- Siguen abiertos los del capturista: `U5-contraste-ui-controles`, `U5-conteo-carpeta`.

## Pendientes
- Corregir A-1 a A-5 y M-1 a M-6 y volver a capturar (`npx playwright test -c playwright.u5.config.ts`) y repetir `captura-sin-conexion`; ronda 2 de revisión.
- Decisiones de Axel sobre los 4 bloqueos `U5-*`.

## Propuestas para Axel (PROPUESTA, no decidido)
| Qué | Resultado | Costo | Riesgo |
|---|---|---|---|
| Corregir las 5 fallas de axe y el desborde de P-26.13 | Pasan U5-01 y U5-03 | Bajo | Ninguno |
| Quitar el selector de archivo nativo y la clave en inglés de Categorías | Pasa U5-07 | Bajo | Ninguno |
| Arrancar sin red con la sesión guardada y la franja «Sin conexión» (H-5) | Pasa U5-13 | Medio | Revisar que no muestre datos tras cerrar sesión |
| Indicador de 2 px (o borde ≥ 3:1) para la vista elegida y pestañas de P-03 | Se va la parte de U5-04 que es estado | Bajo | Cambia un poco el aspecto |
| Decidir la regla de U5-04 para tarjetas-botón y chips con texto | Cierra U5-04 | Ninguno (decisión) | Si se exige 3:1, sube bordes en C2 |
| Título 32 px, barra lateral 240 px, iconos 20 px, espaciado de 4 pt, aviso de P-17 apilado y gráficas | Pasa U5-06 | Medio | Ninguno |
Recomendación: corregir todo lo alto y medio antes de la ronda 2 y responder los 4 bloqueos en la misma pasada; así la ronda 2 puede aprobar.

## Consumo aproximado
≈ 47 acciones de herramienta · ≈ 115 KB de texto leído (más 13 láminas JPEG de capturas, ≈ 1.9 MB, para la lista visual) · ≈ 32 KB escritos en `docs/` (más ≈ 38 KB de scripts temporales en `datos/tmp/`, borrados).
