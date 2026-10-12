# MANIFIESTO · Paquete de entrada para Claude Design

Lo genera `docs/design/tools/paquete_entrada.sh`; no se edita a mano. Cambia solo si cambian los archivos de origen.

- Carpeta de entrada (lo que se sube): `/home/deepsec/finanzas-hogar/docs/design/entrada/`
- Carpeta de cada vuelta (lo que se baja): `/home/deepsec/finanzas-hogar/docs/design/salida/iter-N/`, con N = 1, 2 o 3
- Carpeta de los exportes finales: `/home/deepsec/finanzas-hogar/docs/design/salida/final/`
- Ningún archivo del paquete contiene código ni credenciales.

## Archivos (21; este MANIFIESTO no se lista a sí mismo)

| # | Archivo | Tamaño | sha256 | Para qué sirve |
|---|---|---|---|---|
| 1 | `LEEME.md` | 8212 bytes (8.0 KB) | `f5a78d091d4fd8e3961979dfcc0e27bd6194d2f45a442fb8addfc396a1ab3793` | Pasos para Axel: qué subir, qué pegar, dónde guardar lo que baje y cuándo parar. |
| 2 | `prompt.md` | 12284 bytes (12.0 KB) | `693e3d0530da8d327451994c8fee16ef77899916793e028acaeea94448b08a46` | Mensaje que se pega en Claude Design: pantallas a maquetar, tokens de color en texto, reglas de dinero y accesibilidad. |
| 3 | `DESIGN.pdf` | 2859638 bytes (2792.6 KB) | `670430d909b2232cff97b160039ac1e25847a02cfe053a4800f3d146f475f56c` | Sistema de diseño completo (principios, color, tipografía, componentes, accesibilidad). |
| 4 | `DESIGN.docx` | 82222 bytes (80.3 KB) | `a48a7693d8840639c8d72fac9ae2fb35ce5349ea6d8ece54eabe72fa3d8bcbff` | El mismo sistema de diseño en Word, por si el PDF no se puede leer. |
| 5 | `U3-dashboard.pdf` | 1163454 bytes (1136.2 KB) | `0bd9664cb765d56425ca5160f345457318b465e95aa97fe009486160afd1629f` | Qué lleva el Inicio: Resumen (6 tarjetas), Detalle, Tendencias y Calendario. |
| 6 | `U4-ajustes.pdf` | 1909203 bytes (1864.5 KB) | `2eb165334c330bd63bde7c0d01236cb8a2a2dccb7193fb554413a7e4bc6e3c07` | Las 14 secciones de Ajustes con sus campos y estados. |
| 7 | `C1-resumen.pdf` | 280159 bytes (273.6 KB) | `6c529d29e049ac5d247618e73905041e0242ad880b1929cc96df0ff713633dcf` | Resumen sin código del contrato de datos: módulos, pantallas, entidades, métricas y formatos. |
| 8 | `C2-tokens.json` | 100073 bytes (97.7 KB) | `570ff553cb58e03d096ab65c8b5f3486488633b31f8ec6088e9450386f49fe5d` | Todos los valores de diseño (colores, tipografía, espaciado, radios) en un solo archivo. |
| 9 | `tokens.css` | 26121 bytes (25.5 KB) | `2b9429bc41836fc5a17fa74040ab519033bfdfec47b559a713de57bbaa06708d` | Los mismos valores como variables de hojas de estilo, para consulta. |
| 10 | `paleta.png` | 865657 bytes (845.4 KB) | `8b6bb79f31d439896bcc8af0537b66c9588649e4c6eb0b0fbf02640b1cab70e4` | Imagen de la paleta de color en claro y oscuro. |
| 11 | `reporte-p01.png` | 193383 bytes (188.9 KB) | `2d8b5215a7a27708f3e5bd01545033ee89c9f70b9df4ebaabdb8da93cd477d97` | Página 1 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 12 | `reporte-p02.png` | 126390 bytes (123.4 KB) | `f23e15e788c724246e1271e3ee522dc20532fad222982f69ed41139dfd9327d2` | Página 2 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 13 | `reporte-p03.png` | 217623 bytes (212.5 KB) | `a16f3234097b00fe33acbdebff12f7917f09d380ffa49c99a0da958c0bb82ae8` | Página 3 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 14 | `reporte-p04.png` | 224544 bytes (219.3 KB) | `493a14a3d2f62d830b91fb56b3012f92b28831ce0674a5705b60d07623297101` | Página 4 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 15 | `reporte-p05.png` | 213182 bytes (208.2 KB) | `639cbe5b84e79730ffcc267d17eee9d7c51e6051203582170c74fefd5538e04e` | Página 5 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 16 | `reporte-p06.png` | 174292 bytes (170.2 KB) | `60f154879f501bdb1f807931478be8ba69817a9b060aac79f7ed32caa0d7bfcd` | Página 6 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 17 | `reporte-p07.png` | 247150 bytes (241.4 KB) | `c553dcc0cdc46dce0bf7539cbbf7d7f7062835fd367df0f0ee481b9b077067f3` | Página 7 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 18 | `reporte-p08.png` | 189586 bytes (185.1 KB) | `2743264a00438579fd9a7a3c2699dca67ac2c731ebee2509767f294d5f728412` | Página 8 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 19 | `reporte-p09.png` | 194116 bytes (189.6 KB) | `2bd286c7bc2be0abbb9df99a7e17cd6c5dde83a2a4a0e0cdbf2fd53feb19577a` | Página 9 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 20 | `reporte-p10.png` | 229803 bytes (224.4 KB) | `60ae8cb3a34ef620b7ca6033fca29a3cc3a1ac718a8e2eed0bbcf25fcd292edf` | Página 10 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |
| 21 | `reporte-p11.png` | 310313 bytes (303.0 KB) | `80126d7b9425a89a75b7fa1e918e708c56c86dad4f4fcc66a6159f65993c0021` | Página 11 de 11 del reporte previo de compras: referencia de cifras reales del hogar. |

## Cómo verificarlo

Desde la carpeta de entrada, este comando no debe mostrar ninguna diferencia (todas las líneas terminan en OK):

```
cd /home/deepsec/finanzas-hogar/docs/design/entrada && grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c
```

## Líneas de sha256sum

```
f5a78d091d4fd8e3961979dfcc0e27bd6194d2f45a442fb8addfc396a1ab3793  LEEME.md
693e3d0530da8d327451994c8fee16ef77899916793e028acaeea94448b08a46  prompt.md
670430d909b2232cff97b160039ac1e25847a02cfe053a4800f3d146f475f56c  DESIGN.pdf
a48a7693d8840639c8d72fac9ae2fb35ce5349ea6d8ece54eabe72fa3d8bcbff  DESIGN.docx
0bd9664cb765d56425ca5160f345457318b465e95aa97fe009486160afd1629f  U3-dashboard.pdf
2eb165334c330bd63bde7c0d01236cb8a2a2dccb7193fb554413a7e4bc6e3c07  U4-ajustes.pdf
6c529d29e049ac5d247618e73905041e0242ad880b1929cc96df0ff713633dcf  C1-resumen.pdf
570ff553cb58e03d096ab65c8b5f3486488633b31f8ec6088e9450386f49fe5d  C2-tokens.json
2b9429bc41836fc5a17fa74040ab519033bfdfec47b559a713de57bbaa06708d  tokens.css
8b6bb79f31d439896bcc8af0537b66c9588649e4c6eb0b0fbf02640b1cab70e4  paleta.png
2d8b5215a7a27708f3e5bd01545033ee89c9f70b9df4ebaabdb8da93cd477d97  reporte-p01.png
f23e15e788c724246e1271e3ee522dc20532fad222982f69ed41139dfd9327d2  reporte-p02.png
a16f3234097b00fe33acbdebff12f7917f09d380ffa49c99a0da958c0bb82ae8  reporte-p03.png
493a14a3d2f62d830b91fb56b3012f92b28831ce0674a5705b60d07623297101  reporte-p04.png
639cbe5b84e79730ffcc267d17eee9d7c51e6051203582170c74fefd5538e04e  reporte-p05.png
60f154879f501bdb1f807931478be8ba69817a9b060aac79f7ed32caa0d7bfcd  reporte-p06.png
c553dcc0cdc46dce0bf7539cbbf7d7f7062835fd367df0f0ee481b9b077067f3  reporte-p07.png
2743264a00438579fd9a7a3c2699dca67ac2c731ebee2509767f294d5f728412  reporte-p08.png
2bd286c7bc2be0abbb9df99a7e17cd6c5dde83a2a4a0e0cdbf2fd53feb19577a  reporte-p09.png
60ae8cb3a34ef620b7ca6033fca29a3cc3a1ac718a8e2eed0bbcf25fcd292edf  reporte-p10.png
80126d7b9425a89a75b7fa1e918e708c56c86dad4f4fcc66a6159f65993c0021  reporte-p11.png
```
