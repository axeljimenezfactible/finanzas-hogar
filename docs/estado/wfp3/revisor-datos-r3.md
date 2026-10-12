# revisor-datos-r3 · Reporte

Agente: `revisor-datos-r3` (revisor independiente, ronda 3 de 3 de los datos sintéticos de la Pista 3). Fecha: 11-oct-2026, noche. Informe completo: `docs/estado/wfp3/REVISION-datos-r3.md`.

**Veredicto: NO APROBADO** (0 altos, 1 medio, 4 bajos). M-1 y M-2 de la ronda 2: cerrados. Medio nuevo: `seis_mayores` va por mayor plan de octubre y U1 N-044 pide «las 6 categorías con más gasto del rango» (el conjunto cambia: servicios en lugar de restaurantes).

## Hecho
- Leí (por partes o por script) el JSON, el LEEME 1.2, el verificador, `correccion-datos-r2.md`, `REVISION-datos-r2.md`, P3-08 y P3-12, K45, U1 §2.1 y las N-xxx de §3 que piden series, U3 §3 y §4, U4 línea 624 y la línea 669 del generador.
- Repetí el verificador, la regeneración determinista, el grep de P3-12, el grep de credenciales, 22 sumas y porcentajes propios, la comparación de fechas, montos, marcas y nombres con `data.json` y `model.json`, el formato de dinero de P3-08 y 2 pruebas negativas del verificador.
- Escribí `REVISION-datos-r3.md` y este reporte. No corregí nada.

## Verificación (comandos y salidas)
```
python3 -I docs/design/tools/ds/verificar_sinteticos.py                         -> 57 OK, RESULTADO: OK
generar_sinteticos.py -> datos/tmp/revdatos-r3/regen.json; cmp con el JSON     -> REGEN-IDENTICA
grep -rnoE "<21 cifras de P3-12, con y sin coma>" datos-sinteticos.json | wc -l -> 0  (maquetas/ completo: 0)
las 21 cifras en centavos sin separador                                         -> 0
grep de patrones de credencial de P3-12 en maquetas/                            -> 0; correos en el JSON: 0
chk1..chk4 (22 sumas y porcentajes)                                             -> 21 OK, 1 FALLA (S4: seis_mayores contra U1 N-044)
fechas coincidentes con model.json/data.json (comando de P3-12)                 -> 15, todas día 15/30/28-feb; historial (87 fechas): 0
copia de data.json                                                              -> 0 nombres idénticos; 5 marcas comunes sin precio igual al real
prueba negativa «en 3 gastos» el 20-oct                                         -> FALLA 34, RESULTADO: FALLA
prueba negativa seis_mayores por gasto real                                     -> RESULTADO: OK (el verificador no revisa el criterio)
P3-08: 848 pares _fmt/_cents                                                    -> 0 diferencias; 0 $ mal formados
```

## Supuestos
- «Gasto del rango» de U1 N-044 = Σ de la parte del hogar devengada de los 6 meses que trae el archivo (mayo a octubre), que es el rango por defecto («del primer mes con datos al mes en curso»). Con solo octubre el conjunto también se puede revisar; no lo usé porque U1 dice «del rango».
- Califiqué M-1 como media (inconsistencia entre archivos que cambia una vista) y no como alta (no es un valor de C2/DESIGN.md, ni una cifra real, ni un problema de forma); el autor la declaró como supuesto, aunque con una cita de U1 que no es textual.
- Las coincidencias de valores redondos ($14.00 a $105.00, $1,000.00 a $9,000.00) entre conceptos distintos no las cuento como copia, igual que en la ronda 2.

## Bloqueos
Ninguno.

## Pendientes
- M-1: corregir el criterio en `generar_sinteticos.py` (línea 669), el verificador (comprobación 36) y el LEEME.
- B-1 a B-4 (bajos): N-114 y N-194 sin serie lista; fuente de la regla de 12 meses; bajas B-1 a B-5 de r2 todavía abiertas; CSV de datos semilla «pendiente» con sus datos ya visibles.
- Esta era la ronda 3 de 3: el orquestador decide si se aplica la corrección de M-1 fuera de rondas o si se escala a Axel.

## Consumo aproximado
~34 acciones de herramienta; ~140 KB leídos (salidas de scripts, LEEME 13 KB, REVISION r2 14 KB, correccion r2 5 KB, fragmentos de U1, U3, U4, criterios y decisiones); ~19 KB escritos (este reporte, el informe y 4 scripts temporales en `datos/tmp/revdatos-r3/`, ya borrados).
