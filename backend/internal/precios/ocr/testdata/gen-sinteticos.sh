#!/bin/sh
# Genera 2 tickets SINTÉTICOS (no son fotos reales ni imitan el formato de una
# tienda) para probar la herramienta de medición. NO cuentan para S5-06 (K24).
# Uso: sh gen-sinteticos.sh <carpeta de salida>
set -eu
out="${1:?carpeta de salida}"
font=/usr/share/fonts/liberation/LiberationMono-Regular.ttf
mk() { # nombre texto
  magick -size 760x560 xc:'#f4f1e8' -font "$font" -pointsize 24 -fill '#1a1a1a' \
    -annotate +24+40 "$2" -background '#cfcfcf' -rotate 1.6 \
    -attenuate 0.25 +noise Gaussian -blur 0x0.7 -colorspace sRGB "$out/$1"
}
mk sintetico-01.png 'TIENDA SINTETICA DE PRUEBA
*** TICKET DE PRUEBA NO REAL ***
LECHE ENTERA 1 L                28.50
2 ARROZ BLANCO 1 KG             58.00
0.750 KG JITOMATE               22.13
HUEVO BLANCO 12 PZA             45.90
3 LATA ATUN AGUA                66.00
FRIJOL NEGRO 900 G              31.00
SUBTOTAL                       251.53
TOTAL                          251.53'
mk sintetico-02.png 'TIENDA SINTETICA DE PRUEBA 2
*** TICKET DE PRUEBA NO REAL ***
CAFE SOLUBLE 200 G         $1,234.50 A
7501234567890 REFRESCO COLA 2L  32.00 B
1.250 KG PECHUGA POLLO         118.75
2 JABON DE TOCADOR              24.00
ACEITE VEGETAL 900 ML           47.80
TOTAL                        $1,457.05
EFECTIVO                     $1,500.00
CAMBIO                          42.95'
