// Package adjuntos guarda los adjuntos de los gastos (ticket, factura, estado de cuenta) en disco local: nombre
// generado por el servidor (nunca el del usuario), tipo detectado por el contenido, tamaño máximo vigente y rutas
// que no pueden salir de la carpeta de adjuntos (S6-10).
package adjuntos

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"finanzas-hogar/internal/dominio"
)

// BytesPerMB es el factor de attachment.max_size_mb. PROVISIONAL: C1 dice «10 MB» sin el factor en bytes (S6-10);
// se usa la propuesta de S1, 1 MB = 1,048,576 bytes (bloqueo S6-adjuntos).
const BytesPerMB int64 = 1 << 20

// Dir es la subcarpeta de adjuntos dentro de la carpeta de datos; attachment.file_path es relativa a la de datos.
const Dir = "adjuntos"

// AllowedTypes son los tipos admitidos y su extensión. PROVISIONAL (K31, puerta 4): los que declara C1 en
// uploadExpenseAttachment (JPG, PNG, WebP y PDF; U1 N-064). El tipo se detecta por el contenido, no por el nombre.
var AllowedTypes = map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/webp": ".webp", "application/pdf": ".pdf"}

var (
	ErrTooLarge = errors.New("adjuntos: el archivo pasa del tamaño máximo")
	ErrType     = errors.New("adjuntos: tipo de archivo no admitido")
	ErrPath     = errors.New("adjuntos: ruta fuera de la carpeta de adjuntos")
	ErrName     = errors.New("adjuntos: nombre de archivo no válido")
)

// Store es la carpeta de datos del sistema (api.Config.DataDir).
type Store struct{ DataDir string }

// Saved es un archivo ya guardado.
type Saved struct {
	RelPath  string // adjuntos/AAAA/MM/<32 hex><ext>, relativa a DataDir
	MimeType string
	Size     int64
	SHA256   string
}

// CheckName rechaza nombres originales con separadores de ruta, «..», «.» o NUL: el nombre solo se guarda como dato
// (original_name), nunca se usa para la ruta.
func CheckName(name string) error {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\\x00") || len(name) > 255 {
		return ErrName
	}
	return nil
}

// Save copia r a un archivo nuevo con nombre aleatorio; si pasa de maxBytes o el tipo no está admitido, no deja nada.
func (s Store) Save(r io.Reader, maxBytes int64, now time.Time) (Saved, error) {
	var out Saved
	sub := path.Join(Dir, now.In(dominio.MX).Format("2006/01"))
	dir := filepath.Join(s.DataDir, filepath.FromSlash(sub))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return out, err
	}
	tmp, err := os.CreateTemp(dir, ".subida-*")
	if err != nil {
		return out, err
	}
	keep := false
	defer func() {
		tmp.Close()
		if !keep {
			os.Remove(tmp.Name())
		}
	}()
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(tmp, h), io.LimitReader(r, maxBytes+1))
	if err != nil {
		return out, err
	}
	if n > maxBytes {
		return out, ErrTooLarge
	}
	head := make([]byte, 512)
	k, _ := tmp.ReadAt(head, 0)
	mime, _, _ := strings.Cut(http.DetectContentType(head[:k]), ";")
	ext, ok := AllowedTypes[mime]
	if !ok || n == 0 {
		return out, ErrType
	}
	var rnd [16]byte
	if _, err := rand.Read(rnd[:]); err != nil {
		return out, err
	}
	name := hex.EncodeToString(rnd[:]) + ext
	if err := tmp.Close(); err != nil {
		return out, err
	}
	if err := os.Rename(tmp.Name(), filepath.Join(dir, name)); err != nil {
		return out, err
	}
	keep = true
	return Saved{RelPath: path.Join(sub, name), MimeType: mime, Size: n, SHA256: hex.EncodeToString(h.Sum(nil))}, nil
}

// Resolve convierte attachment.file_path en ruta absoluta y comprueba que quede dentro de la carpeta de adjuntos y
// sea un archivo normal (sin enlaces simbólicos).
func (s Store) Resolve(rel string) (string, error) {
	if rel == "" || strings.ContainsAny(rel, "\\\x00") || path.IsAbs(rel) || path.Clean(rel) != rel || !strings.HasPrefix(rel, Dir+"/") {
		return "", ErrPath
	}
	root, err := filepath.Abs(filepath.Join(s.DataDir, Dir))
	if err != nil {
		return "", err
	}
	full, err := filepath.Abs(filepath.Join(s.DataDir, filepath.FromSlash(rel)))
	if err != nil {
		return "", err
	}
	if r, err := filepath.Rel(root, full); err != nil || r == "." || strings.HasPrefix(r, "..") {
		return "", ErrPath
	}
	st, err := os.Lstat(full)
	if err != nil {
		return "", err
	}
	if !st.Mode().IsRegular() {
		return "", ErrPath
	}
	return full, nil
}

// Remove borra el archivo (si ya no está, no es error).
func (s Store) Remove(rel string) error {
	full, err := s.Resolve(rel)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.Remove(full)
}
