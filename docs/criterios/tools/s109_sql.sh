#!/bin/bash
# Verificacion SQL de S1-09 (texto nuevo): meses calendario y quincenas oct-2026..mar-2027 en base temporal bajo datos/tmp
set -u
R=/home/deepsec/finanzas-hogar
D=$R/datos/tmp/correccion-criterios-r1
DB=$D/cr1.db
mkdir -p $D && rm -f $DB
awk '/^-- \+goose Up/{f=1;next} /^-- \+goose Down/{f=0} f' $R/backend/db/migrations/00001_esquema.sql > $D/up.sql
sqlite3 -bail $DB < $D/up.sql && echo "UP OK"
echo "tablas: $(sqlite3 $DB "select count(*) from sqlite_master where type='table' and name not like 'sqlite_%'")"
python3 -I - <<'PY'
import sqlite3, calendar, datetime
DB='/home/deepsec/finanzas-hogar/datos/tmp/correccion-criterios-r1/cr1.db'
con=sqlite3.connect(DB); con.execute('PRAGMA foreign_keys=ON')
# 1) la ventana del oraculo (10 -> 9) como mes calendario debe rechazarse (antes de insertar los meses, sin conflicto de unicidad)
try:
    con.execute("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES('month',NULL,2026,10,NULL,'2026-10-10','2026-11-09',NULL)")
    print('VENTANA 10->9 ACEPTADA (inesperado)')
except sqlite3.IntegrityError as e:
    print('ventana 10->9 como mes: RECHAZADA por CHECK ->', str(e).split(chr(10))[0][:60])
meses=[(2026,10),(2026,11),(2026,12),(2027,1),(2027,2),(2027,3)]
esperado=[]
for (y,m) in meses:
    ult=calendar.monthrange(y,m)[1]
    ini=f'{y:04d}-{m:02d}-01'; fin=f'{y:04d}-{m:02d}-{ult:02d}'
    cur=con.execute("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES('month',NULL,?,?,NULL,?,?,NULL)",(y,m,ini,fin))
    pid=cur.lastrowid
    q1=(f'{y:04d}-{m:02d}-01',f'{y:04d}-{m:02d}-15',f'{y:04d}-{m:02d}-15')
    q2=(f'{y:04d}-{m:02d}-16',fin,min(f'{y:04d}-{m:02d}-30',fin))
    con.execute("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES('fortnight',?,?,?,1,?,?,?)",(pid,y,m)+q1)
    con.execute("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES('fortnight',?,?,?,2,?,?,?)",(pid,y,m)+q2)
    esperado.append(('month',y,m,None,ini,fin,None)); esperado.append(('fortnight',y,m,1)+q1); esperado.append(('fortnight',y,m,2)+q2)
con.commit()
filas=con.execute("select kind,year,month,half,start_date,end_date,pay_date from period order by year,month,kind='fortnight',half").fetchall()
print('filas insertadas:', len(filas), '(6 meses + 12 quincenas = 18)')
print('recuperadas iguales a lo insertado:', sorted(filas)==sorted(esperado))
feb=con.execute("select start_date,end_date,pay_date from period where kind='fortnight' and year=2027 and month=2 and half=2").fetchone()
print('febrero 2027, quincena 2:', feb, '-> paga el 28:', feb[2]=='2027-02-28')
# K8: quincena 1 = dia 1 al 15; quincena 2 = dia 16 al ultimo dia del mes (el 31 cae en la segunda)
lim={}
for (y,m,h) in ((2026,10,1),(2026,10,2),(2026,12,2),(2027,2,2)):
    r=con.execute("select start_date,end_date from period where kind='fortnight' and year=? and month=? and half=?",(y,m,h)).fetchone()
    lim[(y,m,h)]=r
print('K8 oct-2026 q1:', lim[(2026,10,1)], '| q2:', lim[(2026,10,2)], '| dic-2026 q2:', lim[(2026,12,2)], '| feb-2027 q2:', lim[(2027,2,2)])
print('K8 el 31-oct-2026 cae en la quincena 2:', lim[(2026,10,2)]==('2026-10-16','2026-10-31'), '| el 28-feb-2027 cierra febrero:', lim[(2027,2,2)]==('2027-02-16','2027-02-28'))
bad=con.execute("select count(*) from period where kind='fortnight' and ((half=1 and (substr(start_date,9,2)!='01' or substr(end_date,9,2)!='15')) or (half=2 and (substr(start_date,9,2)!='16' or end_date!=(select m.end_date from period m where m.id=period.parent_id))))").fetchone()[0]
print('K8 quincenas fuera de limites (esperado 0):', bad)
# K8: semana de lunes a domingo (SQLite weekday 0 = domingo)
sem=con.execute("select date('2026-10-14','weekday 0','-6 days'), date('2026-10-14','weekday 0'), date('2026-10-18','weekday 0','-6 days'), date('2026-10-19','weekday 0','-6 days')").fetchone()
print('K8 semana de lunes a domingo, miercoles 14-oct-2026:', sem[0], 'a', sem[1], '| domingo 18 pertenece a la misma:', sem[2]==sem[0], '| lunes 19 abre otra:', sem[3]=='2026-10-19')
print('foreign_key_check filas:', len(con.execute('PRAGMA foreign_key_check').fetchall()), '| integrity_check:', con.execute('PRAGMA integrity_check').fetchone()[0])
con.close()
PY
rm -f $DB $D/up.sql; rmdir $D; rmdir $R/datos/tmp 2>/dev/null
echo "mi carpeta temporal borrada: $([ -e $D ] && echo NO || echo si) | datos/ contiene: $(ls $R/datos | tr '\n' ' ')"
