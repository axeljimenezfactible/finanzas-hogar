# Consumo real por flujo (tokens contados en los transcripts de los agentes; salida = output_tokens; entrada nueva = input + cache_creation; el pico de contexto por agente respeta el tope de 300K salvo donde se indica)

| Flujo | Run | Agentes | Duración | Salida | Entrada nueva | >300K ctx |
|---|---|---|---|---|---|---|
| S1-modelo-datos, S6a-contrato-c1, U1-arq… | wf_9684e5be-d78 | 14 | 3 h 48 min | 2.63 M | 12.68 M | 13 |
| C1-v1.1, S1-esquema-1.2, U1-1.1… | wf_e65bd01e-fc3 | 30 | 3 h 03 min | 2.66 M | 12.21 M | 1 |
| correccion-C1-r4, revisor-C1-r4… | wf_05a92cde-ec9 | 2 | 0 h 18 min | 0.10 M | 0.73 M | 0 |
| C1-v1.2, S1-esquema-1.3, U1-1.2… | wf_1df19586-380 | 27 | 1 h 42 min | 1.34 M | 8.26 M | 0 |
| correccion-C1-r3, correccion-criterios-r… | wf_81e482dd-440 | 4 | 0 h 14 min | 0.17 M | 1.28 M | 0 |
| correccion-criterios-r5, revisor-criteri… | wf_86fdbb95-0e3 | 2 | 0 h 11 min | 0.07 M | 0.62 M | 0 |
| U2b-1-design, U2b-2-design, U2c-tokens… | wf_98b77133-bf1 | 68 | 1 h 52 min | 4.06 M | 23.54 M | 1 |
| correccion-S2-r1, correccion-S2-r2, corr… | wf_b42c6402-0c2 | 48 | 2 h 14 min | 2.79 M | 16.97 M | 0 |
| aplicar-k45, aplicar-k45-fix, redactor-p… | wf_c5708e8e-192 | 8 | 0 h 35 min | 0.41 M | 2.57 M | 0 |
| S1-puerta2, S5-terminos, bloqueos-respue… | wf_cbf48e23-39e | 5 | 0 h 06 min | 0.15 M | 1.36 M | 0 |
| datos-sinteticos, ds-brandbook, ds-token… | wf_6b8245eb-3e3 | 3 | 0 h 06 min | 0.06 M | 0.78 M | 0 |
| correccion-datos-r1, correccion-datos-r2… | wf_79112df4-8e7 | 12 | 1 h 16 min | 0.92 M | 4.91 M | 0 |
| artboards-ajustes-a, artboards-ajustes-b… | wf_73bdc7c4-bb8 | 54 | 1 h 25 min | 3.39 M | 18.21 M | 2 |
| correccion-ajustes-b-r3, correccion-dato… | wf_0baea063-9bb | 5 | 0 h 16 min | 0.18 M | 1.12 M | 0 |
| C1-puerta2, S1-puerta2, U4-puerta2… | wf_1893738a-674 | 47 | 3 h 07 min | 2.72 M | 15.62 M | 1 |
| correccion-datos-r5, correccion-inicio-r… | wf_962d3959-638 | 3 | 0 h 13 min | 0.09 M | 0.61 M | 0 |
| correccion-k53-r1, correccion-k53-r2, k5… | wf_96c77e3b-13d | 8 | 0 h 46 min | 0.44 M | 2.60 M | 0 |
| api-ajustes, api-comidas-salud, api-comp… | wf_6af84047-b37 | 48 | 3 h 31 min | 4.38 M | 19.20 M | 1 |
| c3-regen, correccion-frontend-areas-r4, … | wf_aa499ddf-e8e | 15 | 2 h 24 min | 0.98 M | 6.21 M | 0 |
| c3-matriz, correccion-c3-r1, correccion-… | wf_ac0a0518-622 | 17 | 4 h 56 min | 1.19 M | 7.54 M | 0 |
| api-tablero, c3-tablero, correccion-espe… | wf_f337ec50-76d | 23 | 4 h 06 min | 1.35 M | 9.56 M | 1 |
| backend-cierre, c3-cierre, correccion-ba… | wf_a12c0227-086 | 13 | 2 h 46 min | 0.65 M | 6.27 M | 0 |
| **Total** | | **456** | | **30.73 M** | **172.86 M** | |

Plan §5 estimaba ≈ 8.2 M (rango 5.4–11.1 M) para toda la Etapa B contando solo salida y lectura ponderada; el conteo real incluye toda la entrada nueva de cada turno de cada agente y las rondas de revisión y corrección, por eso es mucho mayor. Fuente: `consumo.py` sobre cada run.
