#!/usr/bin/env python3
"""Verificación independiente de las cifras canon de S5 (K-01 a K-08) directo de datos/semilla/data.json, con Decimal.

No usa pre_calculo.py: suma el reporte renglón por renglón, muestra el cálculo de +138.50 y de +100.50 y luego lo contrasta con
oraculo/salidas/s5.json (caso s5-comparativa-canon). Uso: python3 -I oraculo/pre_verify.py   Código 0 si todo coincide.
"""
import json
import sys
from decimal import Decimal
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
rep = json.load(open(RAIZ / 'datos' / 'semilla' / 'data.json', encoding='utf-8'), parse_float=Decimal)
idn, sus, t = rep['items_identicos'], rep['sustitutos'], rep['totales']
c = lambda x: int(Decimal(x) * 100)
ok = True


def chequea(nombre, calculado, esperado):
    global ok
    bien = calculado == esperado
    ok &= bien
    print(f"{'OK ' if bien else 'DIF'} {nombre}: calculado={calculado} esperado={esperado}")


ib = sum(c(i['bodega_unit']) * i['cantidad'] for i in idn)
iw = sum(c(i['walmart_unit']) * i['cantidad'] for i in idn)
sb = sum(c(s['bodega_total']) for s in sus)
sw = sum(c(s['walmart_total']) for s in sus)
lata = c(rep['atun_sexta_pieza_bodega'])
print(f'Idénticos: Bodega {ib / 100:.2f} (Σ de 33 renglones: precio unitario x cantidad), Walmart {iw / 100:.2f}')
print(f'Sustitutos: Bodega {sb / 100:.2f} (66.00 + 99.00 + 135.00), Walmart {sw / 100:.2f} (84.00 + 105.00 + 170.00); sexta lata de atún de Bodega {lata / 100:.2f}')
tb, tw = ib + sb + lata, iw + sw
print(f'Total Bodega = {ib / 100:.2f} + {sb / 100:.2f} + {lata / 100:.2f} = {tb / 100:.2f}')
print(f'Total Walmart = {iw / 100:.2f} + {sw / 100:.2f} = {tw / 100:.2f}')
print(f'VALOR CANON 1 (carrito contra carrito): Walmart - Bodega = {tw / 100:.2f} - {tb / 100:.2f} = +{(tw - tb) / 100:.2f}')
print(f'VALOR CANON 2 (mismos artículos, mismas cantidades): {iw / 100:.2f} - {ib / 100:.2f} = +{(iw - ib) / 100:.2f}')
print(f'Explicación: +{(iw - ib) / 100:.2f} (idénticos) + {(sw - sb) / 100:.2f} (sustitutos {sw / 100:.2f} - {sb / 100:.2f}) - {lata / 100:.2f} (lata que solo compró Bodega) = +{(iw - ib + sw - sb - lata) / 100:.2f}')
chequea('K-01 Bodega', tb, c(t['bodega']))
chequea('K-02 Walmart', tw, c(t['walmart']))
chequea('K-03 diferencia total', tw - tb, 13850)
chequea('K-03 suma de la explicación', iw - ib + sw - sb - lata, 13850)
chequea('K-04 idénticos Bodega', ib, c(t['bodega_identicos']))
chequea('K-05 idénticos Walmart', iw, c(t['walmart_identicos']))
chequea('K-06 diferencia idénticos', iw - ib, 10050)
pi = sum(i['cantidad'] for i in idn)
pb = pi + sum(s['cant_b'] for s in sus) + 1
pw = pi + sum(s['cant_w'] for s in sus)
chequea('K-07 piezas Bodega/Walmart/idénticos', (pb, pw, pi), (t['articulos_bodega'], t['articulos_walmart'], t['articulos_identicos']))
chequea('K-08 renglones idénticos/sustitutos', (len(idn), len(sus)), (33, 3))
g = (sum(i['bodega_unit'] < i['walmart_unit'] for i in idn), sum(i['bodega_unit'] == i['walmart_unit'] for i in idn), sum(i['bodega_unit'] > i['walmart_unit'] for i in idn))
chequea('ganadores bodega/empates/walmart', g, (22, 8, 3))
sal = json.load(open(RAIZ / 'oraculo' / 'salidas' / 's5.json', encoding='utf-8'))['casos']['s5-comparativa-canon']['canon-totales']
b, w = sal['stores']
chequea('oráculo: totales', (b['total_cents'], w['total_cents'], b['identical_total_cents'], w['identical_total_cents']), (tb, tw, ib, iw))
chequea('oráculo: piezas', (b['piece_count'], w['piece_count'], sal['identical_pieces']), (pb, pw, pi))
chequea('oráculo: renglones', (b['identical_count'], b['substitute_count']), (33, 3))
chequea('oráculo: ganadores', ([x['count'] for x in sal['winners']], sal['ties']), ([22, 3], 8))
chequea('oráculo: Walmart +138.50', w['difference_vs_cheapest_cents'], 13850)
print('TODO COINCIDE' if ok else 'HAY DIFERENCIAS')
sys.exit(0 if ok else 1)
