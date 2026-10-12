# Reporte · revisor-bloqueos-r3 (11-oct-2026, noche)

Agente: `revisor-bloqueos-r3`, revisor independiente, ronda 3 de 3. No escribí nada de lo que revisé y no corregí nada. El informe completo está en `docs/estado/p2c/REVISION-bloqueos-r3.md`.

## Hecho
- **Veredicto: NO APROBADO.** Hay 1 hallazgo alto, 0 medios y 4 bajos.
- **Fidelidad sin fallas:**
  - Los 27 bloqueos tocados solo agregan texto al final y ninguno de los 42 resueltos cambió.
  - Las 53 letras con respuesta son la recomendada de PUERTA-2 §8.
  - Q2-11 = B y Q2-64 = A están escritas tal cual.
  - Las 8 preguntas «sin recomendación», más Q2-12 y la pregunta 4 de `S2-ventanas`, no llevan respuesta, quedan PROVISIONAL y están ABIERTAS.
  - En la muestra de 10 contra PUERTA-2, coinciden la letra y el sentido.
  - La tabla de K46 tiene 0 diferencias con §8 en las 98 filas.
  - En §K46 a §K52 no hay citas inventadas de Axel.
  - K50 trae los mismos valores que PUERTA-2 §4.
- **Cerrados:**
  - H-r2-02: el INDICE y K46 coinciden con el disco (75 bloqueos, 63 resueltos y 12 abiertos; las herramientas cuentan 67 y 8).
  - H-r2-03: la cita de Q2-12 va por sección y las filas existen.
- **Reabierto:** H-r2-01 pasa a ser H-r3-01 (alto), porque la verificación global sigue rota.
  - Hubo avance: el lado S1 ya pasa en el árbol.
  - Falta regenerar C1: `openapi.yaml` y `C1-datos.md` no han cambiado desde las 16:59.
  - `gen2` no conoce el tipo «hora HH:MM» y `verify_c1.py` falla K44 y C1-15.
  - No se arregla en los bloqueos.
- **Siguen abiertos:** los 4 bajos de la ronda 2 (ahora H-r3-02 a H-r3-05).

## Verificación
Todos los comandos los corrí yo. Usé tres clones en `datos/tmp/r3` (`a` = HEAD, `b` = HEAD más los bloqueos y `decisiones.md`, `t` = HEAD más todo el árbol) y mis scripts `v3.py` y `k3.py`. Al terminar borré todo `datos/tmp/r3`.

| Comando | Salida |
|---|---|
| `git diff --numstat -- docs/estado/bloqueos/` | 27 bloqueos con «+N 0» e `INDICE.md` con «117 5» |
| `python3 -I v3.py` | «HEAD bloqueos 69 · HEAD con «Respuesta de Axel» 41 · modificados 27 · sin prefijo HEAD [] · resueltos-HEAD modificados []»; «disco 75»; regla de las herramientas `{'RESUELTO': 67, 'ABIERTO': 8}`; «entradas 62 repetidas []»; «faltan Q2-11..71: []»; «sin rec con respuesta: []»; «letra != recomendada: []»; «Q2-11 [('S1-parametros-avisos.md', 'B')]»; «Q2-64 [('S5-terminos.md', 'A')]»; «literales no hallados en HEAD: [] · con literal: 53» |
| `python3 -I k3.py` | «K46 filas 98 dif []»; 15 sin respuesta; citas de Axel presentes en RESPUESTAS.md; «INDICE filas puerta 2: 62», con la única diferencia de forma en Q2-64 («RESUELTO (decisión); bloqueo ABIERTO»), que es la misma del archivo |
| `make -C a verify-c1` / `verify-s1` | «RESULTADO: TODO OK» / «RESUMEN TODO OK», rc=0 |
| `make -C b verify-c1` / `verify-s1` | «FALLA C1-15», rc=2 / «RESUMEN FALLAS: ['S1-16', 'K44']», rc=2 |
| `make -C t verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | rc=2: WF-0 OK y luego «gen2: tipo de parámetro desconocido: hora HH:MM» |
| `make -C t -i verify-c1` | «FALLA K44», «FALLA C1-15», «RESULTADO: 2 falla(s)» |
| `make -C t verify-s1` / `verify-u1` / `verify-criterios` | «RESUMEN TODO OK» / «RESULTADO: TODO OK» / rc=0 |

## Supuestos
- Para contar con la regla de la puerta 2, cuando un archivo no trae «Estado del bloqueo (puerta 2)» ni «Respuesta de Axel», tomé el estado que declara su encabezado. Pasa en `criterios-nxx` y en los 2 de P3.
- Juzgué la pregunta 4 de `S2-ventanas` (casos 1 y 2 abiertos) como fiel a PUERTA-2 §2.1. Q2-01 = A autoriza escribir la pregunta, pero no la contesta.

## Bloqueos
Ninguno nuevo.

## Pendientes
- H-r3-01: le toca al agente que regenera C1 (solicitud 1 de `docs/estado/p2c/S1-puerta2.md`) y al integrador. Se cierra cuando `make verify-all` pasa con todo el árbol.
- H-r3-02 y H-r3-03 (bajos): para una corrección de bloqueos, si el orquestador la lanza.
- H-r3-04 y H-r3-05 (bajos): H-r3-04 le toca al orquestador; H-r3-05 es una mejora de formato del INDICE.

## Propuestas para Axel
Ninguna nueva. Las preguntas abiertas que se pueden contestar por id son:
- Q2-12 (C o A);
- la pregunta 4 de `S2-ventanas`;
- Q2-48, Q2-50 a Q2-52 y Q2-56 a Q2-58;
- los 4 bloqueos nuevos de P2c;
- el siguiente paso de Q2-64 (B, URL legales, C o D).

## Solicitudes fuera de alcance
1. (Alta, H-r3-01) Al agente de C1 y al integrador: aplicar la solicitud 1 de `S1-puerta2.md`, que es regenerar C1 con `c1gen`:
   - agregar el tipo «hora HH:MM»;
   - exponer las claves de `app_setting`;
   - quitar las 6 marcas de `openapi.yaml` y las de `C1-datos.md`;
   - corregir las marcas de los criterios `U4.md` l. 6, `C1.md` (C1-22) y `S7.md`.

   No hacer el commit «[Puerta 2]» hasta que `make verify-all` pase.
2. (Baja, H-r3-04) Al orquestador: anotar en RESPUESTAS.md que Q2-12 queda abierta, o confirmarla con Axel.

## Consumo aproximado
32 acciones de herramienta · unos 120 KB leídos · unos 21 KB escritos (los dos informes y los scripts temporales).
