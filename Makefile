# Makefile · Finanzas Hogar
# Verificación reproducible de los cimientos (WF-0 «Cimientos»), de la puerta 1 (flujo P1), de la puerta 1b (flujo P1c)
# y de WF-1 «Motores y sistema de diseño».
# Detalle en README.md. Las migraciones se toman por comodín (backend/db/migrations/*.sql, hoy 00001 a 00003).
#
#   make verify-wf0   SQL en base temporal, go build y go vet, colorimetría y validador OpenAPI
#   make verify-c1    generador de C1 (build.py --check, incluidos los extractos) y verify_c1.py contra la base migrada
#   make verify-s1    backend/db/tools/s1_verifica.py: esquema 1.3 (todas las migraciones), oráculo, negativas, Down y goose
#   make verify-u1    docs/design/tools/verificar_u1.py: U1 1.2 contra model.json, data.json, C1 y las migraciones
#   make verify-criterios  docs/criterios/tools/verif_final.sh: conteos (y conteo.py contra INDICE.md), trazabilidad, citas y cifras canon (18 archivos)
#   make verify-wf1   go test con cobertura K25, oráculos S2 a S5 con compara.py, recetario, tokens C2, contraste, U4 y paquete
#   make verify-all   verify-wf0 + verify-c1 + verify-s1 + verify-u1 + verify-criterios + verify-wf1
#   make clean        borra los temporales de este Makefile
#
# Requisitos: GNU Make y bash; sqlite3 >= 3.37 (tablas STRICT); Go 1.27 o posterior (go.mod: go 1.27);
# Python 3.11 o posterior; uv y uvx con openapi-spec-validator 0.9.0 (y pyyaml y jsonschema
# para verify-c1 y verify-criterios; jsonschema para verify-wf1) en su caché; Node 20+ con docs/design/tools/tokens/node_modules
# (npm ci) y pdftotext para verify-wf1. Sin caché, la primera vez: make verify-wf0 UVX_FLAGS= UV_FLAGS=
# Opcional: goose v3.28.0 en el PATH o con GOOSE=/ruta/a/goose para el ciclo up/down/up.
# Los temporales van a datos/tmp/verify-wf0 y se borran al terminar bien; si algo falla,
# quedan para revisarlos y «make clean» los quita.

SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help
.NOTPARALLEL:

ROOT       := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
TMP_ROOT   := $(ROOT)/datos/tmp
TMP        := $(TMP_ROOT)/verify-wf0
MIGRATIONS := $(sort $(wildcard $(ROOT)/backend/db/migrations/*.sql))
COLOR_PY   := $(ROOT)/docs/design/tools/colorimetria.py
COLOR_OUT  := docs/design/colorimetria.md docs/design/contraste.csv docs/contratos/C2-colores.json docs/design/paleta.html
C1GEN      := $(ROOT)/backend/internal/api/c1gen
CRITERIOS_SH ?= $(ROOT)/docs/criterios/tools/verif_final.sh

PYTHON     ?= python3
SQLITE     ?= sqlite3
GO         ?= go
GOFMT      ?= gofmt
UV         ?= uv
UVX        ?= uvx
UV_FLAGS   ?= --offline
UVX_FLAGS  ?= --offline
OPENAPI_VALIDATOR ?= openapi-spec-validator==0.9.0
GOOSE      ?= $(shell command -v goose 2>/dev/null)

# Programas awk que extraen la sección Up o Down de una migración goose.
AWK_UP     := '/^-- \+goose Up/{s=1;next} /^-- \+goose Down/{s=0} s'
AWK_DOWN   := '/^-- \+goose Down/{s=1;next} /^-- \+goose Up/{s=0} s'
RM_TMP     := rm -rf "$(TMP)"; rmdir "$(TMP_ROOT)" 2>/dev/null || true

.PHONY: help verify-wf0 verify-sql verify-go verify-color verify-openapi verify-c1 verify-s1 verify-u1 verify-criterios verify-all clean \
        verify-wf1 verify-wf1-cobertura verify-wf1-oraculos verify-wf1-recetario verify-wf1-tokens verify-wf1-u4 verify-wf1-paquete

help:
	@echo "Finanzas Hogar · objetivos de make"
	@echo "  make verify-wf0   SQL en base temporal, go build y go vet, colorimetría y validador OpenAPI"
	@echo "  make verify-c1    generador de C1 (build.py --check, incluidos los extractos) y verify_c1.py contra la base migrada"
	@echo "  make verify-s1    backend/db/tools/s1_verifica.py (esquema 1.3, oráculo, negativas, Down y goose)"
	@echo "  make verify-u1    docs/design/tools/verificar_u1.py (U1 1.2)"
	@echo "  make verify-criterios  docs/criterios/tools/verif_final.sh (18 archivos de criterios)"
	@echo "  make verify-wf1   go test -cover (K25), oráculos S2 a S5, recetario, tokens C2, contraste, U4 unido y paquete"
	@echo "  make verify-all   verify-wf0 + verify-c1 + verify-s1 + verify-u1 + verify-criterios + verify-wf1"
	@echo "  make clean        borra datos/tmp/verify-wf0, backend/bin y las carpetas __pycache__"

verify-wf0: verify-sql verify-go verify-color verify-openapi
	@$(RM_TMP)
	@echo "== WF-0: verificación completa OK (SQL, Go, colorimetría y OpenAPI)"

verify-sql:
	@echo "== SQL: migraciones de backend/db/migrations en una base temporal"
	@test -n "$(MIGRATIONS)" || { echo "   FALLA: no hay migraciones en backend/db/migrations"; exit 1; }
	@mkdir -p "$(TMP)"; rm -f "$(TMP)"/sql-*.db
	@for f in $(MIGRATIONS); do \
	  awk $(AWK_UP) "$$f" | $(SQLITE) -bail "$(TMP)/sql-up.db"; \
	  echo "   Up aplicado con sqlite3 $$($(SQLITE) --version | cut -d' ' -f1): $${f#$(ROOT)/}"; \
	done
	@db="$(TMP)/sql-up.db"; \
	t=$$($(SQLITE) "$$db" "SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"); \
	s=$$($(SQLITE) "$$db" "SELECT count(*) FROM pragma_table_list WHERE schema='main' AND type='table' AND strict=1 AND name NOT LIKE 'sqlite_%'"); \
	i=$$($(SQLITE) "$$db" "SELECT count(*) FROM sqlite_master WHERE type='index'"); \
	g=$$($(SQLITE) "$$db" "SELECT count(*) FROM sqlite_master WHERE type='trigger'"); \
	fk=$$($(SQLITE) "$$db" "SELECT count(*) FROM pragma_foreign_key_check"); \
	ic=$$($(SQLITE) "$$db" "PRAGMA integrity_check"); \
	echo "   tablas=$$t (STRICT=$$s) · índices=$$i · disparadores=$$g · foreign_key_check=$$fk filas · integrity_check=$$ic"; \
	test "$$t" -gt 0 -a "$$t" = "$$s" -a "$$fk" = 0 -a "$$ic" = ok || { echo "   FALLA: el esquema no cumple STRICT, llaves foráneas o integridad"; exit 1; }
	@cp "$(TMP)/sql-up.db" "$(TMP)/sql-down.db"; \
	for f in $$(printf '%s\n' $(MIGRATIONS) | sort -r); do \
	  awk $(AWK_DOWN) "$$f" | $(SQLITE) -bail "$(TMP)/sql-down.db"; \
	done; \
	n=$$($(SQLITE) "$(TMP)/sql-down.db" "SELECT count(*) FROM sqlite_master WHERE name NOT LIKE 'sqlite_%'"); \
	echo "   Down en orden inverso ($$(printf '%s\n' $(MIGRATIONS) | sort -r | sed 's#.*/##; s/_.*//' | paste -sd, -)) deja $$n objetos"; \
	test "$$n" = 0 || { echo "   FALLA: Down no deja la base vacía"; exit 1; }; \
	for f in $(MIGRATIONS); do awk $(AWK_UP) "$$f" | $(SQLITE) -bail "$(TMP)/sql-down.db"; done; \
	a=$$($(SQLITE) "$(TMP)/sql-up.db" "SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"); \
	b=$$($(SQLITE) "$(TMP)/sql-down.db" "SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"); \
	echo "   Up otra vez: $$b tablas (primera vez: $$a)"; \
	test "$$a" = "$$b" || { echo "   FALLA: Up después de Down no reproduce el esquema"; exit 1; }
	@if [ -n "$(GOOSE)" ]; then \
	  g="$(TMP)/sql-goose.db"; log="$(TMP)/goose.log"; dir="$(ROOT)/backend/db/migrations"; \
	  { $(GOOSE) -dir "$$dir" sqlite3 "$$g" up && $(GOOSE) -dir "$$dir" sqlite3 "$$g" down-to 0 && $(GOOSE) -dir "$$dir" sqlite3 "$$g" up; } > "$$log" 2>&1 || { cat "$$log"; exit 1; }; \
	  echo "   goose $$($(GOOSE) --version 2>&1 | awk '{print $$3}'): up, down-to 0 y up otra vez sin errores"; \
	else \
	  echo "   goose no está en el PATH: se omite su ciclo (opcional; GOOSE=/ruta/a/goose)"; \
	fi

verify-go:
	@echo "== Go: go build, go vet y gofmt en backend/"
	@cd "$(ROOT)/backend" && echo "   $$($(GO) version)" && $(GO) build ./... && $(GO) vet ./... && echo "   go build ./... y go vet ./...: sin errores"
	@cd "$(ROOT)/backend" && out="$$($(GOFMT) -l .)" && if [ -n "$$out" ]; then echo "   FALLA gofmt: $$out"; exit 1; fi && echo "   gofmt -l .: sin archivos por formatear"

verify-color:
	@echo "== Colorimetría: colorimetria.py regenera en un árbol espejo y se compara con el repositorio"
	@e="$(TMP)/espejo"; rm -rf "$$e"; mkdir -p "$$e/docs/design/tools" "$$e/docs/contratos"; \
	cp "$(COLOR_PY)" "$$e/docs/design/tools/"; \
	$(PYTHON) -I "$$e/docs/design/tools/colorimetria.py" --sin-png > "$(TMP)/colorimetria.log" 2>&1 || { cat "$(TMP)/colorimetria.log"; exit 1; }; \
	grep -E 'contraste.csv:|VERIFICACI' "$(TMP)/colorimetria.log" | sed 's/^/   /' || true; \
	for f in $(COLOR_OUT); do \
	  if cmp -s "$(ROOT)/$$f" "$$e/$$f"; then echo "   idéntico byte a byte: $$f"; else echo "   FALLA, difiere del regenerado: $$f"; exit 1; fi; \
	done
	@$(PYTHON) -I -c "import csv, sys; r = list(csv.DictReader(open(sys.argv[1], encoding='utf-8'))); n = sum(1 for x in r if x['cumple'].strip().lower() != 'true'); print('   contraste.csv del repositorio: %d filas, %d con cumple=false' % (len(r), n)); sys.exit(1 if n or not r else 0)" "$(ROOT)/docs/design/contraste.csv"

verify-openapi:
	@echo "== OpenAPI: $(OPENAPI_VALIDATOR) sobre docs/contratos/openapi.yaml"
	@cd "$(ROOT)" && $(UVX) $(UVX_FLAGS) --from $(OPENAPI_VALIDATOR) openapi-spec-validator docs/contratos/openapi.yaml | sed 's/^/   /'

verify-c1:
	@echo "== C1: build.py --check (salidas idénticas al generador) y verify_c1.py --db"
	@cd "$(ROOT)" && $(PYTHON) -I "$(C1GEN)/build.py" --check | grep -E 'IDÉNTICO|DIFIERE|RESULTADO' | sed 's/^/   /'
	@mkdir -p "$(TMP)"; rm -f "$(TMP)/c1.db"; \
	for f in $(MIGRATIONS); do awk $(AWK_UP) "$$f" | $(SQLITE) -bail "$(TMP)/c1.db"; done; \
	cd "$(ROOT)" && $(UV) run $(UV_FLAGS) --no-project --with pyyaml --with jsonschema $(PYTHON) -I "$(C1GEN)/verify_c1.py" --db "$(TMP)/c1.db" > "$(TMP)/verify_c1.log" 2>&1 || { cat "$(TMP)/verify_c1.log"; exit 1; }; \
	echo "   comprobaciones OK: $$(grep -cE '^OK ' "$(TMP)/verify_c1.log")"; \
	grep -E '^FALLA|RESULTADO' "$(TMP)/verify_c1.log" | sed 's/^/   /'
	@$(RM_TMP)

verify-s1:
	@echo "== S1: backend/db/tools/s1_verifica.py (esquema 1.3: todas las migraciones, oráculo, negativas, Down y goose)"
	@mkdir -p "$(TMP)"; \
	cd "$(ROOT)" && GOOSE="$(GOOSE)" $(PYTHON) -I backend/db/tools/s1_verifica.py > "$(TMP)/s1_verifica.log" 2>&1 || { cat "$(TMP)/s1_verifica.log"; exit 1; }; \
	cut -c1-110 "$(TMP)/s1_verifica.log" | sed 's/^/   /'; \
	for f in $(MIGRATIONS); do b=$${f##*/}; grep -qF "$$b" "$(TMP)/s1_verifica.log" || { echo "   FALLA: s1_verifica.py no aplicó $$b"; exit 1; }; done; \
	echo "   s1_verifica.py aplicó las $(words $(MIGRATIONS)) migraciones del comodín"
	@$(RM_TMP)

verify-u1:
	@echo "== U1: docs/design/tools/verificar_u1.py (U1 1.2 contra model.json, data.json, C1 y las migraciones)"
	@mkdir -p "$(TMP)"; log="$(TMP)/verificar_u1.log"; \
	cd "$(ROOT)" && $(PYTHON) -I docs/design/tools/verificar_u1.py > "$$log" 2>&1 || { cat "$$log"; exit 1; }; \
	echo "   comprobaciones OK: $$(grep -cE '^OK ' "$$log") · INFO: $$(grep -cE '^INFO ' "$$log") · PENDIENTE: $$(grep -cE '^PENDIENTE' "$$log")"; \
	grep -E '^(FALLA|PENDIENTE)|PROVISIONAL línea|^RESULTADO' "$$log" | cut -c1-110 | sed 's/^/   /'
	@$(RM_TMP)

# Salidas de verif_final.sh que deben aparecer tal cual; «sin fuente» de citas.py es informativo
# (hallazgo bajo B-04 de docs/estado/p1/REVISION-criterios-r3.md: fragmentos citados en otro orden).
verify-criterios:
	@echo "== Criterios: docs/criterios/tools/verif_final.sh (18 archivos; escribe solo su base temporal bajo datos/tmp)"
	@mkdir -p "$(TMP)"; log="$(TMP)/criterios.log"; \
	cd "$(ROOT)" && bash "$(CRITERIOS_SH)" > "$$log" 2>&1 || { cat "$$log"; exit 1; }; \
	grep -E '^TOTAL [0-9]+ \(|^archivos: |^criterios definidos: |^citas verificadas: |^TOTAL comprobadas: |^cifras comparadas: |^N-xxx sin operacion|^matriz C1|^INDICE contra conteo propio' "$$log" | sed -E 's/^(TOTAL [0-9]+ \(alta [0-9]+ \/ media [0-9]+ \/ baja [0-9]+\))$$/\1 · paso 1, cuenta.py: S1 a S8 y WF4/' | cut -c1-110 | sed 's/^/   /'; \
	falla=0; \
	for e in 'incompletos: ninguno' 'secciones 1-5 faltantes: []' 'archivos con Verifica != criterios: [] | con Evidencia != criterios: []' \
	  'criterios sin linea en su seccion 5 (con rangos expandidos): 0' 'citados que no existen: {}' '| con diferencias: 0' \
	  'K-01..K-08 contra data.json: OK' 'UP OK' 'foreign_key_check filas: 0 | integrity_check: ok' 'cifras comparadas: 12 | diferencias: 0' \
	  'N-xxx sin operacion ni Local: 0 []' 'P-xx (con secciones P-26.xx y Global) sin operacion ni Local: 0 []' 'N-xxx de U1 ausentes en la matriz C1: []' \
	  'INDICE contra conteo propio: 0 diferencias'; do \
	  grep -qF -- "$$e" "$$log" || { echo "   FALLA, no aparece en la salida: $$e"; falla=1; }; \
	done; \
	for n in 5 10; do \
	  v=$$(awk -v s="== $$n. " 'index($$0, s) == 1 { getline; print; exit }' "$$log"); \
	  test "$$v" = 0 || { echo "   FALLA, la comprobación $$n de verif_final.sh da «$$v» (se espera 0)"; falla=1; }; \
	done; \
	test "$$falla" = 0 || { cat "$$log"; exit 1; }; \
	echo "   residuos «N-xx» y rutas /tmp en los criterios: 0 y 0; las 11 comprobaciones dan lo esperado"
	@$(RM_TMP)

# ---- WF-1 «Motores y sistema de diseño» (integrador-wf1, 11-oct-2026) ----
# Cobertura K25: umbral de COVER_MIN % en cada paquete con pruebas de los árboles dominio, motor_fin, motor_cred,
# recetas y precios (los programas main sin pruebas de esos árboles se informan); el resto se mide sin umbral.
# Los motores se corren DESDE LA RAÍZ del repositorio: las fixtures citan datos/semilla/... con rutas relativas a ella
# (desde backend/ el motor responde invalid_argument en 14 fixtures de S4 y S5 y compara.py da diferencias).
COVER_MIN   ?= 85
K25_RE      := ^finanzas-hogar/internal/(dominio|motor_fin|motor_cred|recetas|precios)(/|$$)
K25_RAICES  := dominio motor_fin motor_cred recetas precios
MOTORES     := s2 s3 s4 s5
MODEL_SHA   := 1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4
VERIFY3_SHA := d86087f3797e08ceabd2c721bf423115345b318ae7fe5fa13b3400cd87b95041
TOKENS_DIR  := $(ROOT)/docs/design/tools/tokens
NODE        ?= node

# Recetario: JSON Schema, 60 recetas sin ids ni nombres repetidos y orden = recetario-a + recetario-b.
define RECETARIO_PY
import json, sys
import jsonschema
s = json.load(open('datos/semilla/recetario.schema.json', encoding='utf-8'))
d = json.load(open('datos/semilla/recetario.json', encoding='utf-8'))
jsonschema.Draft202012Validator.check_schema(s)
errores = list(jsonschema.Draft202012Validator(s).iter_errors(d))
r = d['recetas']
ids = [x['id'] for x in r]
nombres = {x['nombre'].casefold() for x in r}
partes = [x['id'] for p in ('a', 'b') for x in json.load(open('datos/semilla/recetario-%s.json' % p, encoding='utf-8'))['recetas']]
sin_kcal = sum(1 for x in r if x.get('kcal_por_porcion') is None)
print('   JSON Schema: %d errores · recetas: %d · ids únicos: %d · nombres únicos: %d · kcal_por_porcion null: %d · orden = a + b: %s'
      % (len(errores), len(r), len(set(ids)), len(nombres), sin_kcal, ids == partes))
sys.exit(0 if not errores and len(r) == 60 and len(set(ids)) == 60 and len(nombres) == 60 and ids == partes else 1)
endef
export RECETARIO_PY

# U4 unido: encabezado e índice del integrador + parte a + línea en blanco + parte b, sin cambios; 14 secciones P-26.
define U4_PY
import re, sys
t = open('docs/design/U4-ajustes.md', encoding='utf-8').read()
a = open('docs/design/U4-ajustes-a.md', encoding='utf-8').read()
b = open('docs/design/U4-ajustes-b.md', encoding='utf-8').read()
ok = t.endswith(a + '\n' + b)
n = len(re.findall(r'(?m)^## P-26\.', t))
print('   U4-ajustes.md termina con parte a + parte b sin cambios: %s · secciones P-26: %d de 14' % (ok, n))
sys.exit(0 if ok and n == 14 else 1)
endef
export U4_PY

verify-wf1: verify-go verify-wf1-cobertura verify-wf1-oraculos verify-wf1-recetario verify-wf1-tokens verify-wf1-u4 verify-wf1-paquete
	@$(RM_TMP)
	@echo "== WF-1: verificación completa OK (Go y cobertura K25, oráculos S2 a S5, recetario, tokens y contraste, U4 y paquete)"

verify-wf1-cobertura:
	@echo "== Go: go test ./... -count=1 -cover (umbral K25: $(COVER_MIN) % en dominio, motor_fin, motor_cred, recetas y precios)"
	@mkdir -p "$(TMP)"; log="$(TMP)/gotest.log"; lista="$(TMP)/gopkgs.txt"; \
	cd "$(ROOT)/backend" && $(GO) test ./... -count=1 -cover > "$$log" 2>&1 || { cat "$$log"; exit 1; }; \
	$(GO) list -f '{{.ImportPath}} {{.Name}}' ./... > "$$lista"; \
	awk -v re='$(K25_RE)' -v min=$(COVER_MIN) -v lista="$$lista" -v raices='$(K25_RAICES)' ' \
	  BEGIN { while ((getline l < lista) > 0) { split(l, x, " "); nom[x[1]] = x[2] } } \
	  { pkg = ""; cov = ""; ok = ($$1 == "ok"); \
	    for (i = 1; i <= NF; i++) { if ($$i ~ /^finanzas-hogar/) pkg = $$i; if ($$i == "coverage:") { cov = $$(i + 1); sub(/%/, "", cov) } } \
	    if (pkg == "" || cov == "") next; \
	    if (pkg !~ re) { printf "   sin umbral  %-46s %5.1f %%%s\n", pkg, cov, (ok ? "" : " (sin pruebas)"); next } \
	    if (ok) { medido[pkg] = 1; if (cov + 0 >= min) printf "   K25 OK      %-46s %5.1f %%\n", pkg, cov; else { printf "   K25 FALLA   %-46s %5.1f %% (< %d %%)\n", pkg, cov, min; f = 1 } } \
	    else if (nom[pkg] == "main") printf "   K25 INFO    %-46s programa main sin pruebas (%s %%)\n", pkg, cov; \
	    else { printf "   K25 FALLA   %-46s paquete sin pruebas\n", pkg; f = 1 } } \
	  END { n = split(raices, r, " "); for (i = 1; i <= n; i++) { p = "finanzas-hogar/internal/" r[i]; \
	          if (!(p in nom)) printf "   K25 INFO    %-46s todavía no existe (no es de WF-1)\n", p; \
	          else if (!(p in medido)) { printf "   K25 FALLA   %-46s sin medición\n", p; f = 1 } } \
	        exit f }' "$$log"

verify-wf1-oraculos:
	@echo "== Oráculos S2 a S5: run_fixtures reproduce oraculo/salidas, motor Go desde la raíz y compara.py con 0 diferencias"
	@mkdir -p "$(TMP)"; cd "$(ROOT)/backend" && $(GO) build -o "$(TMP)/motores" ./cmd/motores
	@cd "$(ROOT)" && for m in $(MOTORES); do \
	  case $$m in s2) o=fin_run_fixtures.py;; s3) o=cred_run_fixtures.py;; s4) o=rec_run_fixtures.py;; s5) o=pre_run_fixtures.py;; esac; \
	  $(PYTHON) -I "oraculo/$$o" --out "$(TMP)/oraculo-$$m.json" > "$(TMP)/oraculo-$$m.log" 2>&1 || { cat "$(TMP)/oraculo-$$m.log"; exit 1; }; \
	  cmp -s "$(TMP)/oraculo-$$m.json" "oraculo/salidas/$$m.json" || { echo "   FALLA: oraculo/salidas/$$m.json difiere de lo que regenera $$o"; exit 1; }; \
	  "$(TMP)/motores" $$m dump --fixtures "oraculo/fixtures/$$m" --out "$(TMP)/motor-$$m.json" || { echo "   FALLA: motores $$m dump"; exit 1; }; \
	  $(PYTHON) -I oraculo/compara.py $$m --motor "$(TMP)/motor-$$m.json" > "$(TMP)/compara-$$m.log" 2>&1 || { grep -v '^DIF' "$(TMP)/compara-$$m.log"; grep -m20 '^DIF' "$(TMP)/compara-$$m.log" || true; exit 1; }; \
	  grep -qx 'diferencias: 0' "$(TMP)/compara-$$m.log" || { cat "$(TMP)/compara-$$m.log"; exit 1; }; \
	  echo "   $$m · $$(grep -m1 -E 'fixtures' "$(TMP)/oraculo-$$m.log" | cut -c1-60) · salidas/$$m.json idéntica"; \
	  grep -E '^(oráculo contra motor|model.json \(ramas|data.json|provisionales)' "$(TMP)/compara-$$m.log" | paste -sd'|' | sed 's/|/ · /g; s/^/     /'; \
	done
	@cd "$(ROOT)" && $(PYTHON) -I oraculo/fin_verify.py | sed 's/^/   /'; \
	$(PYTHON) -I oraculo/pre_verify.py > "$(TMP)/pre_verify.log" 2>&1 && grep -qx 'TODO COINCIDE' "$(TMP)/pre_verify.log" || { cat "$(TMP)/pre_verify.log"; exit 1; }; \
	echo "   pre_verify.py: $$(grep -c '^OK ' "$(TMP)/pre_verify.log") comprobaciones OK · TODO COINCIDE"
	@cd "$(ROOT)" && m=$$(sha256sum oraculo/model.json | cut -c1-64); v=$$(sha256sum oraculo/verify3.json | cut -c1-64); \
	test "$$m" = "$(MODEL_SHA)" -a "$$v" = "$(VERIFY3_SHA)" || { echo "   FALLA: oraculo/model.json u oraculo/verify3.json cambió (sha256)"; exit 1; }; \
	$(PYTHON) -I oraculo/model.py > "$(TMP)/model.log" 2>&1 || { cat "$(TMP)/model.log"; exit 1; }; \
	grep -q 'igual a oraculo/model.json: True' "$(TMP)/model.log" || { echo "   FALLA: model.py no reproduce model.json"; exit 1; }; \
	$(PYTHON) -I oraculo/verify3.py > /dev/null 2>&1 && cmp -s datos/tmp/verify3.json oraculo/verify3.json || { echo "   FALLA: verify3.py no reproduce verify3.json"; exit 1; }; \
	rm -f datos/tmp/model.json datos/tmp/verify3.json; \
	echo "   model.json y verify3.json intactos (sha256 $${m:0:12}… y $${v:0:12}…) y reproducidos por model.py y verify3.py"

verify-wf1-recetario:
	@echo "== Recetario: rec_genera_recetario.py --verifica, JSON Schema, unión a + b y pruebas Go del cargador"
	@cd "$(ROOT)" && $(PYTHON) -I oraculo/rec_genera_recetario.py --verifica | sed 's/^/   /'
	@cd "$(ROOT)" && $(UV) run $(UV_FLAGS) --no-project --with jsonschema $(PYTHON) -I -c "$$RECETARIO_PY"
	@cd "$(ROOT)/backend" && $(GO) test ./internal/recetas -count=1 -run 'TestRecetarioValido|TestRecetarioProduccion|TestEscaladoRecetario|TestKcalRecetasAutor' | sed 's/^/   go test: /'

verify-wf1-tokens:
	@echo "== Tokens C2 y contraste: generar idéntico, build reproducible (sincronia.mjs), verificar_c2 y contraste_componentes --check"
	@test -d "$(TOKENS_DIR)/node_modules" || { echo "   FALLA: falta docs/design/tools/tokens/node_modules (cd docs/design/tools/tokens && npm ci)"; exit 1; }
	@mkdir -p "$(TMP)"; c2="$(ROOT)/docs/contratos/C2-tokens.json"; cp "$$c2" "$(TMP)/C2-tokens.antes.json"; \
	cd "$(TOKENS_DIR)" && $(PYTHON) -I generar_c2_tokens.py > "$(TMP)/c2gen.log" 2>&1 || { cp "$(TMP)/C2-tokens.antes.json" "$$c2"; cat "$(TMP)/c2gen.log"; exit 1; }; \
	if cmp -s "$$c2" "$(TMP)/C2-tokens.antes.json"; then echo "   generar_c2_tokens.py: C2-tokens.json idéntico byte a byte"; \
	else cp "$(TMP)/C2-tokens.antes.json" "$$c2"; echo "   FALLA: generar_c2_tokens.py no reproduce C2-tokens.json (se restauró el versionado)"; exit 1; fi
	@cd "$(TOKENS_DIR)" && $(PYTHON) -I verificar_c2.py > "$(TMP)/c2verify.log" 2>&1 || { cat "$(TMP)/c2verify.log"; exit 1; }; \
	grep -E '^RESULTADO' "$(TMP)/c2verify.log" | sed 's/^/   verificar_c2.py: /'
	@cd "$(TOKENS_DIR)" && $(NODE) sincronia.mjs > "$(TMP)/c2sync.log" 2>&1 || { cat "$(TMP)/c2sync.log"; exit 1; }; \
	grep -E '^OK ' "$(TMP)/c2sync.log" | cut -c1-70 | sed 's/^/   sincronia.mjs: /'; rm -rf "$(TMP_ROOT)/c2-sync"
	@cd "$(ROOT)" && $(PYTHON) -I docs/design/tools/contraste_componentes.py --check > "$(TMP)/contraste.log" 2>&1 || { cat "$(TMP)/contraste.log"; exit 1; }; \
	grep -E 'IDENTICO|bajo umbral' "$(TMP)/contraste.log" | cut -c1-110 | sed 's/^/   contraste_componentes.py: /'

verify-wf1-u4:
	@echo "== U4: docs/design/U4-ajustes.md = encabezado e índice + U4-ajustes-a.md + U4-ajustes-b.md"
	@cd "$(ROOT)" && $(PYTHON) -I -c "$$U4_PY"

verify-wf1-paquete:
	@echo "== Paquete para Claude Design: sha256sum -c contra MANIFIESTO.md y sin código en docs/design/entrada"
	@cd "$(ROOT)/docs/design/entrada" && n=$$(grep -cE '^[0-9a-f]{64}  ' MANIFIESTO.md) && \
	grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c --quiet - && \
	t=$$(ls -A | wc -l) && test "$$t" = $$((n + 1)) || { echo "   FALLA: MANIFIESTO.md no cubre los $$t archivos de entrada/"; exit 1; }; \
	echo "   sha256sum -c: $$n archivos sin diferencias (más MANIFIESTO.md)"
	@cd "$(ROOT)/docs/design/entrada" && otros=$$(ls -A | grep -vE '\.(pdf|png|md|json|css|docx)$$' || true); \
	test -z "$$otros" || { echo "   FALLA: tipos no permitidos: $$otros"; exit 1; }; \
	c=$$(grep -lE '<script|^#!|^(import|package|func|def) ' -- *.md *.json *.css 2>/dev/null || true); \
	test -z "$$c" || { echo "   FALLA: parece código: $$c"; exit 1; }; \
	k=0; if command -v pdftotext > /dev/null; then for f in *.pdf; do k=$$((k + $$(pdftotext "$$f" - | grep -c '/api/v1' || true))); done; fi; \
	test "$$k" = 0 || { echo "   FALLA: los PDF citan /api/v1 $$k veces"; exit 1; }; \
	echo "   solo pdf, png, md, json, css y docx; sin scripts ni código; pdftotext: 0 menciones de /api/v1"

verify-all: verify-wf0 verify-c1 verify-s1 verify-u1 verify-criterios verify-wf1
	@$(RM_TMP)
	@echo "== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)"

clean:
	@rm -rf "$(TMP)" "$(ROOT)/backend/bin"
	@find "$(ROOT)/backend" "$(ROOT)/docs" "$(ROOT)/oraculo" -type d -name __pycache__ -prune -exec rm -rf {} + 2>/dev/null || true
	@rmdir "$(TMP_ROOT)" 2>/dev/null || true
	@if [ -d "$(TMP_ROOT)" ]; then echo "Aviso: datos/tmp conserva archivos que no son de este Makefile (no se borraron):"; ls -A "$(TMP_ROOT)"; else echo "Limpio: sin datos/tmp, sin backend/bin y sin __pycache__."; fi
