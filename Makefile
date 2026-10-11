# Makefile · Finanzas Hogar
# Verificación reproducible de los cimientos (WF-0 «Cimientos») y de la puerta 1 (flujo P1). Detalle en README.md.
#
#   make verify-wf0   SQL en base temporal, go build y go vet, colorimetría y validador OpenAPI
#   make verify-c1    generador de C1 (build.py --check, incluidos los extractos) y verify_c1.py contra la base migrada
#   make verify-s1    backend/db/tools/s1_verifica.py: esquema 1.2 (todas las migraciones), oráculo, negativas, Down y goose
#   make verify-u1    docs/design/tools/verificar_u1.py: U1 1.1 contra model.json, data.json, C1 y las migraciones
#   make verify-criterios  docs/criterios/tools/verif_final.sh: conteos, trazabilidad, citas y cifras canon (18 archivos)
#   make verify-all   verify-wf0 + verify-c1 + verify-s1 + verify-u1 + verify-criterios
#   make clean        borra los temporales de este Makefile
#
# Requisitos: GNU Make y bash; sqlite3 >= 3.37 (tablas STRICT); Go 1.27 o posterior (go.mod: go 1.27);
# Python 3.11 o posterior; uv y uvx con openapi-spec-validator 0.9.0 (y pyyaml y jsonschema
# para verify-c1 y verify-criterios) en su caché. Sin caché, la primera vez: make verify-wf0 UVX_FLAGS= UV_FLAGS=
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

.PHONY: help verify-wf0 verify-sql verify-go verify-color verify-openapi verify-c1 verify-s1 verify-u1 verify-criterios verify-all clean

help:
	@echo "Finanzas Hogar · objetivos de make"
	@echo "  make verify-wf0   SQL en base temporal, go build y go vet, colorimetría y validador OpenAPI"
	@echo "  make verify-c1    generador de C1 (build.py --check, incluidos los extractos) y verify_c1.py contra la base migrada"
	@echo "  make verify-s1    backend/db/tools/s1_verifica.py (esquema 1.2, oráculo, negativas, Down y goose)"
	@echo "  make verify-u1    docs/design/tools/verificar_u1.py (U1 1.1)"
	@echo "  make verify-criterios  docs/criterios/tools/verif_final.sh (18 archivos de criterios)"
	@echo "  make verify-all   verify-wf0 + verify-c1 + verify-s1 + verify-u1 + verify-criterios"
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
	echo "   Down en orden inverso deja $$n objetos"; \
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
	@echo "== S1: backend/db/tools/s1_verifica.py (esquema 1.2: todas las migraciones, oráculo, negativas, Down y goose)"
	@mkdir -p "$(TMP)"; \
	cd "$(ROOT)" && GOOSE="$(GOOSE)" $(PYTHON) -I backend/db/tools/s1_verifica.py > "$(TMP)/s1_verifica.log" 2>&1 || { cat "$(TMP)/s1_verifica.log"; exit 1; }; \
	cut -c1-110 "$(TMP)/s1_verifica.log" | sed 's/^/   /'
	@$(RM_TMP)

verify-u1:
	@echo "== U1: docs/design/tools/verificar_u1.py (U1 1.1 contra model.json, data.json, C1 y las migraciones)"
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
	grep -E '^TOTAL [0-9]+ \(|^archivos: |^criterios definidos: |^citas verificadas: |^TOTAL comprobadas: |^cifras comparadas: |^N-xxx sin operacion|^matriz C1' "$$log" | cut -c1-110 | sed 's/^/   /'; \
	falla=0; \
	for e in 'incompletos: ninguno' 'secciones 1-5 faltantes: []' 'archivos con Verifica != criterios: [] | con Evidencia != criterios: []' \
	  'criterios sin linea en su seccion 5 (con rangos expandidos): 0' 'citados que no existen: {}' '| con diferencias: 0' \
	  'K-01..K-08 contra data.json: OK' 'UP OK' 'foreign_key_check filas: 0 | integrity_check: ok' 'cifras comparadas: 12 | diferencias: 0' \
	  'N-xxx sin operacion ni Local: 0 []' 'P-xx (con secciones P-26.xx y Global) sin operacion ni Local: 0 []' 'N-xxx de U1 ausentes en la matriz C1: []'; do \
	  grep -qF -- "$$e" "$$log" || { echo "   FALLA, no aparece en la salida: $$e"; falla=1; }; \
	done; \
	for n in 5 10; do \
	  v=$$(awk -v s="== $$n. " 'index($$0, s) == 1 { getline; print; exit }' "$$log"); \
	  test "$$v" = 0 || { echo "   FALLA, la comprobación $$n de verif_final.sh da «$$v» (se espera 0)"; falla=1; }; \
	done; \
	test "$$falla" = 0 || { cat "$$log"; exit 1; }; \
	echo "   residuos «N-xx» y rutas /tmp en los criterios: 0 y 0; las 10 comprobaciones dan lo esperado"
	@$(RM_TMP)

verify-all: verify-wf0 verify-c1 verify-s1 verify-u1 verify-criterios
	@$(RM_TMP)
	@echo "== Verificación ampliada OK (WF-0, C1, S1, U1 y criterios)"

clean:
	@rm -rf "$(TMP)" "$(ROOT)/backend/bin"
	@find "$(ROOT)/backend" "$(ROOT)/docs" "$(ROOT)/oraculo" -type d -name __pycache__ -prune -exec rm -rf {} + 2>/dev/null || true
	@rmdir "$(TMP_ROOT)" 2>/dev/null || true
	@if [ -d "$(TMP_ROOT)" ]; then echo "Aviso: datos/tmp conserva archivos que no son de este Makefile (no se borraron):"; ls -A "$(TMP_ROOT)"; else echo "Limpio: sin datos/tmp, sin backend/bin y sin __pycache__."; fi
