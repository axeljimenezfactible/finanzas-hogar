package ocr

import (
	"bytes"
	"context"
	"os/exec"
)

// ExecRunner ejecuta programas locales con os/exec, sin shell. Excluido de la
// cobertura por contrato (S5-06: «sin contar el exec»).
type ExecRunner struct{}

// Run implementa Runner.
func (ExecRunner) Run(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	return out.Bytes(), errb.Bytes(), err
}
