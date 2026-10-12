package semilla

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"finanzas-hogar/internal/dominio"
)

// Variables de entorno de las contraseñas iniciales (nunca se escriben en el repositorio ni en el reporte).
//
//	FINANZAS_SEED_PASSWORD_AXEL, FINANZAS_SEED_PASSWORD_XIMENA   una por cuenta
//	FINANZAS_SEED_PASSWORD_FILE                                   archivo FUERA del repositorio, modo 0600, con
//	                                                              renglones «usuario=contraseña» (axel=…, ximena=…)
//
// Axel las cambia después en Ajustes, «Cuenta y roles» (P-26.01, operación changeOwnPassword): la contraseña
// inicial solo sirve para entrar la primera vez.
const (
	EnvPasswordAxel   = "FINANZAS_SEED_PASSWORD_AXEL"
	EnvPasswordXimena = "FINANZAS_SEED_PASSWORD_XIMENA"
	EnvPasswordFile   = "FINANZAS_SEED_PASSWORD_FILE"
)

var seedUsers = []struct{ code, username, name string }{{"axel", "axel", "Axel"}, {"ximena", "ximena", "Ximena"}}

// PasswordsFromEnv lee las contraseñas iniciales del entorno o del archivo fuera del repositorio.
// Falta una = no se devuelve error aquí: solo se exige si hay que crear esa cuenta.
func PasswordsFromEnv() (map[string]string, error) {
	out := map[string]string{}
	if f := os.Getenv(EnvPasswordFile); f != "" {
		m, err := readPasswordFile(f)
		if err != nil {
			return nil, err
		}
		out = m
	}
	for user, env := range map[string]string{"axel": EnvPasswordAxel, "ximena": EnvPasswordXimena} {
		if v := os.Getenv(env); v != "" {
			out[user] = v
		}
	}
	return out, nil
}

func readPasswordFile(path string) (map[string]string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	st, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("semilla: no puedo leer el archivo de contraseñas iniciales: %w", err)
	}
	if st.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("semilla: el archivo de contraseñas iniciales debe ser solo tuyo (chmod 600 %s)", abs)
	}
	for dir := filepath.Dir(abs); ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return nil, errors.New("semilla: el archivo de contraseñas iniciales no puede estar dentro de un repositorio git; ponlo fuera (por ejemplo en ~/.config/finanzas-hogar/)")
		}
		if filepath.Dir(dir) == dir {
			break
		}
	}
	f, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	out := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		u, p, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(u) == "" || p == "" {
			return nil, errors.New("semilla: el archivo de contraseñas iniciales debe tener renglones «usuario=contraseña»")
		}
		out[strings.ToLower(strings.TrimSpace(u))] = p
	}
	return out, sc.Err()
}

// personasYCuentas: Axel (titular), Ximena (co-titular), Karen (externa, sin cuenta) y las 2 cuentas.
func (l *loader) personasYCuentas() error {
	people := []struct {
		code, name, part string
		direct           int
	}{
		{"axel", "Axel", "holder", 0}, {"ximena", "Ximena", "co_holder", 0}, {"karen", "Karen", "external", 1}}
	for _, p := range people {
		id, _, err := l.ensure("personas", "person", `SELECT id FROM person WHERE code = ?`, a(p.code),
			`INSERT INTO person (code, display_name, participation, pays_share_directly) VALUES (?, ?, ?, ?)`,
			p.code, p.name, p.part, p.direct)
		if err != nil {
			return err
		}
		l.ids["person:"+p.code] = id
	}
	minLen, err := dominio.SettingInt(l.ctx, l.tx, "auth.password_min_length")
	if err != nil {
		return err
	}
	var missing []string
	for _, u := range seedUsers {
		if _, ok, err := l.find(`SELECT id FROM user_account WHERE username = ?`, u.username); err != nil {
			return err
		} else if ok {
			continue
		}
		pw, ok := l.opt.Passwords[u.username]
		if !ok || pw == "" {
			missing = append(missing, u.username)
			continue
		}
		if int64(len([]rune(pw))) < minLen {
			return fmt.Errorf("semilla: la contraseña inicial de %s es más corta que el mínimo de %d caracteres", u.username, minLen)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("semilla: faltan las contraseñas iniciales de %s. Defínelas en %s y %s, o en un archivo fuera del repositorio con %s (ver docs/estado/wf3/semilla.md)",
			strings.Join(missing, " y "), EnvPasswordAxel, EnvPasswordXimena, EnvPasswordFile)
	}
	for _, u := range seedUsers {
		pid := l.ids["person:"+u.code]
		if id, ok, err := l.find(`SELECT id FROM user_account WHERE username = ?`, u.username); err != nil {
			return err
		} else if ok {
			l.rep.count("cuentas", false)
			l.ids["user:"+u.username] = id
			continue
		}
		hash, err := dominio.HashPassword(l.opt.Passwords[u.username], l.argo)
		if err != nil {
			return err
		}
		id, err := l.insert("user_account", `INSERT INTO user_account (person_id, username, password_hash, password_changed_at) VALUES (?, ?, ?, NULL)`,
			pid, u.username, hash)
		if err != nil {
			return err
		}
		// audit_log guarda el renglón completo: el hash argon2id es una salida de un solo sentido, nunca la contraseña.
		l.rep.count("cuentas", true)
		l.ids["user:"+u.username] = id
	}
	return nil
}

// reglasYCategorias: Hogar (por defecto, 100 % al titular), Karen (100 %) y Conjunto 50/50 (Karen 5000 pb, resto al titular)
// según RN-03; categorías Despensa (groceries, reparto Conjunto 50/50) y Comida fuera (eating_out, decisión B).
// Las demás categorías las captura Axel en Ajustes (P-26.04 «Sin valor: lo captura Axel»).
func (l *loader) reglasYCategorias() error {
	axel := l.id("person:axel")
	rules := []struct {
		key, name string
		def       int
		karenBP   int
		note      string
	}{
		{"hogar", "Hogar", 1, 0, "Todo el monto es del hogar Axel–Ximena (a nombre del titular). Regla por defecto (RN-03)."},
		{"karen", "Karen", 0, 10000, "Todo el monto es de Karen."},
		{"conjunto", "Conjunto 50/50", 0, 5000, "Karen 50 % y el resto del hogar. Regla de la categoría Despensa (decisión B)."},
	}
	for _, r := range rules {
		id, created, err := l.ensure("reglas de reparto", "cost_share_rule",
			`SELECT id FROM cost_share_rule WHERE is_template = 1 AND name = ? COLLATE NOCASE`, a(r.name),
			`INSERT INTO cost_share_rule (name, remainder_person_id, is_template, is_default, notes) VALUES (?, ?, 1, ?, ?)`,
			r.name, axel, r.def, r.note)
		if err != nil {
			return err
		}
		l.ids["rule:"+r.key] = id
		if created && r.karenBP > 0 {
			if _, err := l.insert("cost_share_rule_part", `INSERT INTO cost_share_rule_part (rule_id, person_id, pct_bp) VALUES (?, ?, ?)`,
				id, l.id("person:karen"), r.karenBP); err != nil {
				return err
			}
		}
	}
	cats := []struct {
		code, name string
		rule       string
		order      int
	}{{"groceries", "Despensa", "conjunto", 1}, {"eating_out", "Comida fuera", "", 2}}
	for _, c := range cats {
		var rule any
		if c.rule != "" {
			rule = l.id("rule:" + c.rule)
		}
		id, _, err := l.ensure("categorías", "category", `SELECT id FROM category WHERE code = ?`, a(c.code),
			`INSERT INTO category (parent_id, code, name, kind, default_cost_share_rule_id, sort_order) VALUES (NULL, ?, ?, 'variable', ?, ?)`,
			c.code, c.name, rule, c.order)
		if err != nil {
			return err
		}
		l.ids["category:"+c.code] = id
	}
	return nil
}
