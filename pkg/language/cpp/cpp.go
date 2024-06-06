package cpp

import (
	"context"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/cp-setter-toolkit/backend/pkg/memory"
	"github.com/cp-setter-toolkit/backend/pkg/sandbox"
)

type Cpp struct {
	id   string
	name string
	args []string
}

func NewCpp(id string, name string, args ...string) *Cpp {
	return &Cpp{id, name, args}
}

func (l Cpp) Id() string {
	return l.id
}

func (l Cpp) Name() string {
	return l.name
}

func (l Cpp) FileExtensions() []string {
	return []string{".cpp", ".cc", ".cxx", ".c++"}
}

func (l Cpp) Compile(ctx context.Context, sb sandbox.Sandbox, files []sandbox.File, stderr io.Writer) (*sandbox.File, error) {
	args := slices.Clone(l.args)
	for _, f := range files {
		if err := sandbox.CreateFile(sb, f); err != nil {
			return nil, err
		}
		if !strings.HasSuffix(f.Name, ".h") {
			args = append(args, f.Name)
		}
	}

	config := sandbox.RunConfig{
		Stdout:      stderr,
		Stderr:      stderr,
		MemoryLimit: 256 * memory.MiB,
		TimeLimit:   10 * time.Second,
		MaxProcs:    200,
		InheritEnv:  true,
		WorkingDir:  sb.Pwd(),
	}

	if _, err := sb.Run(ctx, config, "/usr/bin/g++", args...); err != nil {
		return nil, err
	}
	return sandbox.ExtractContent(sb, "a.out")
}

func (l Cpp) Run(ctx context.Context, sb sandbox.Sandbox, bin sandbox.File, stdin io.Reader, stdout io.Writer, tl time.Duration, ml memory.Amount) (*sandbox.Status, error) {
	return sandbox.RunBinary(ctx, sb, bin, stdin, stdout, tl, ml)
}
