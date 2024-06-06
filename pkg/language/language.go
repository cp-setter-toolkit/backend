package language

import (
	"context"
	"io"
	"time"

	"github.com/cp-setter-toolkit/backend/pkg/memory"
	"github.com/cp-setter-toolkit/backend/pkg/sandbox"
)

// Language is the interface for programming languages.
type Language interface {
	Id() string
	Name() string
	FileExtensions() []string
	Compile(ctx context.Context, sb sandbox.Sandbox, files []sandbox.File, stderr io.Writer) (*sandbox.File, error)
	Run(ctx context.Context, sb sandbox.Sandbox, bin sandbox.File, stdin io.Reader, stdout io.Writer, tl time.Duration, ml memory.Amount) (*sandbox.Status, error)
}

// Registry is a collection of languages.
type Registry interface {
	Register(lang Language)
	Get(id string) (Language, error)
}
