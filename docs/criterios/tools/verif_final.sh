#!/bin/bash
# Verificacion de docs/criterios (solo lectura salvo la base temporal propia de s109_sql.sh bajo datos/tmp).
# Uso: bash docs/criterios/tools/verif_final.sh   (desde cualquier carpeta)
T="$(cd "$(dirname "$0")" && pwd)"
cd "$T" || exit 1
echo "== fecha: $(date '+%F %T')"
echo "== 1. conteo por archivo y severidad (cuenta.py)"; python3 -I cuenta.py
echo "== 2. estructura y trazabilidad de los 18 archivos (traza.py)"; python3 -I traza.py
echo "== 3. referencias cruzadas entre IDs (refs.py)"; python3 -I refs.py
echo "== 4. citas de la seccion 5 contra plan/decisiones/protocolo (citas.py)"; python3 -I citas.py
echo "== 5. residuos de N-xx de 2 equis fuera de la explicacion del INDICE"; grep -n -E "N-xx([^x]|$)" /home/deepsec/finanzas-hogar/docs/criterios/*.md | grep -v "INDICE.md" | wc -l
echo "== 6. cifras canon R-01..R-85 (S2 y S3) y K-01..K-08 (S5) (canon.py)"; python3 -I canon.py
echo "== 7. S1-09 (a): SQL en base temporal (s109_sql.sh)"; bash ./s109_sql.sh
echo "== 8. S1-09 (b): equivalencia calendario vs ventanas del oraculo (s109_equiv.py)"; python3 -I s109_equiv.py
echo "== 9. C1-02/U1-03: consistencia de la matriz (matriz_c102.py)"; uv run --no-project --with pyyaml python3 -I matriz_c102.py
echo "== 10. rutas /tmp en los criterios S y WF4 (esperado 0)"; grep -n '[^s]/tmp/' /home/deepsec/finanzas-hogar/docs/criterios/{S1,S2,S3,S4,S5,S6,S7,S8,WF4-operacion}.md | wc -l
echo "== 11. conteo propio por archivo, versiones y coincidencia con INDICE.md (conteo.py)"; python3 -I conteo.py | tail -2
