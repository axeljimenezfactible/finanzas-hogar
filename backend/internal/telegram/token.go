package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TokenSource da el token del bot en el momento del envío (K30: nunca se copia ni se guarda).
type TokenSource func() (string, error)

// TokenPaths son las rutas de la clave dentro de openclaw.json, en orden. La primera es la literal de K30.
// PROVISIONAL (bloqueo docs/estado/bloqueos/S7-ruta-token.md): la segunda es la que usa la herramienta de
// Axel para el mismo bot (canal telegram); se quita o se queda según responda Axel.
var TokenPaths = [][]string{
	{"bots", "multix", "botToken"},
	{"channels", "telegram", "bots", "multix", "botToken"},
}

// DefaultConfigPath es ~/.openclaw/openclaw.json del usuario del servicio (ACL de lectura, K30).
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".openclaw", "openclaw.json")
	}
	return filepath.Join(home, ".openclaw", "openclaw.json")
}

// FileToken lee el token de path en cada envío. Los errores nunca llevan el contenido del archivo.
func FileToken(path string) TokenSource {
	return func() (string, error) {
		raw, err := os.ReadFile(path)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return "", fmt.Errorf("telegram: no existe %s; el bot MultiX no está configurado en esta máquina", path)
		case errors.Is(err, fs.ErrPermission):
			return "", fmt.Errorf("telegram: sin permiso de lectura sobre %s; falta la ACL del usuario del servicio (K30)", path)
		case err != nil:
			return "", fmt.Errorf("telegram: no se pudo leer %s", path)
		}
		var root map[string]any
		if json.Unmarshal(raw, &root) != nil {
			return "", fmt.Errorf("telegram: %s no es JSON válido", path)
		}
		for _, p := range TokenPaths {
			var cur any = root
			for _, k := range p {
				m, ok := cur.(map[string]any)
				if !ok {
					cur = nil
					break
				}
				cur = m[k]
			}
			if s, ok := cur.(string); ok && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s), nil
			}
		}
		return "", fmt.Errorf("telegram: %s no trae la clave bots.multix.botToken", path)
	}
}
