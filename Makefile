# Makefile · Finanzas Hogar
# Verificación reproducible de los cimientos (WF-0 «Cimientos»). Detalle en README.md.
#
#   make verify-wf0   SQL en base temporal, go build y go vet, colorimetría y validador OpenAPI
#   make verify-c1    generador de C1 (build.py --check) y verify_c1.py contra la base migrada
#   make verify-s1    s1_verifica.py (Anexo B del reporte de S1): ida y vuelta del oráculo y negativas
#   make verify-all   verify-wf0 + verify-c1 + verify-s1
#   make clean        borra los temporales de este Makefile
#
# Requisitos: GNU Make y bash; sqlite3 >= 3.37 (tablas STRICT); Go 1.26 o posterior;
# Python 3.11 o posterior; uv y uvx con openapi-spec-validator 0.9.0 (y pyyaml y jsonschema
# para verify-c1) en su caché. Sin caché, la primera vez: make verify-wf0 UVX_FLAGS= UV_FLAGS=
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

.PHONY: help verify-wf0 verify-sql verify-go verify-color verify-openapi verify-c1 verify-s1 verify-all clean

help:
	@echo "Finanzas Hogar · objetivos de make"
	@echo "  make verify-wf0   SQL en base temporal, go build y go vet, colorimetría y validador OpenAPI"
	@echo "  make verify-c1    generador de C1 (build.py --check) y verify_c1.py contra la base migrada"
	@echo "  make verify-s1    s1_verifica.py del reporte de S1 (ida y vuelta del oráculo y negativas)"
	@echo "  make verify-all   verify-wf0 + verify-c1 + verify-s1"
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

# Código que extrae y ejecuta s1_verifica.py desde su anexo (mientras no viva en una ruta propia).
define S1_VERIFICA_PY
import re
t = open('docs/estado/wf0/S1-modelo-datos.md', encoding='utf-8').read()
m = re.search(r'<!-- inicio s1_verifica.py -->\n```python\n(.*?)\n```\n<!-- fin s1_verifica.py -->', t, re.S)
exec(compile(m.group(1), 's1_verifica.py', 'exec'))
endef
export S1_VERIFICA_PY

verify-s1:
	@echo "== S1: s1_verifica.py (Anexo B de docs/estado/wf0/S1-modelo-datos.md)"
	@mkdir -p "$(TMP)"; \
	cd "$(ROOT)" && GOOSE="$(GOOSE)" $(PYTHON) -I -c "$$S1_VERIFICA_PY" > "$(TMP)/s1_verifica.log" 2>&1 || { cat "$(TMP)/s1_verifica.log"; exit 1; }; \
	cut -c1-110 "$(TMP)/s1_verifica.log" | sed 's/^/   /'; \
	grep -oE 'suma pagos d1=[0-9]+ · parte Axel d1=[0-9]+' "$(TMP)/s1_verifica.log" | sed 's/^/   regresión Kueski (centavos): /' || true
	@$(RM_TMP)

verify-all: verify-wf0 verify-c1 verify-s1
	@$(RM_TMP)
	@echo "== Verificación ampliada OK (WF-0, C1 y S1)"

clean:
	@rm -rf "$(TMP)" "$(ROOT)/backend/bin"
	@find "$(ROOT)/backend" "$(ROOT)/docs" "$(ROOT)/oraculo" -type d -name __pycache__ -prune -exec rm -rf {} + 2>/dev/null || true
	@rmdir "$(TMP_ROOT)" 2>/dev/null || true
	@if [ -d "$(TMP_ROOT)" ]; then echo "Aviso: datos/tmp conserva archivos que no son de este Makefile (no se borraron):"; ls -A "$(TMP_ROOT)"; else echo "Limpio: sin datos/tmp, sin backend/bin y sin __pycache__."; fi
