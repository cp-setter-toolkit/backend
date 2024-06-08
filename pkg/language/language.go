package language

import (
	"context"
	"io"

	"github.com/cp-setter-toolkit/cp-setter-toolkit/pkg/sandbox"
)

// Language is the interface for programming languages.
type Language interface {
	Id() string
	Name() string
	FileExtensions() []string
	Compile(ctx context.Context, sb sandbox.Sandbox, files []sandbox.File, stderr io.Writer) (sandbox.File, error)
	Run(ctx context.Context, sb sandbox.Sandbox, config sandbox.RunConfig, file sandbox.File) (*sandbox.Status, error)
}

// Registry is a collection of languages.
type Registry interface {
	Register(lang Language)
	Get(id string) (Language, error)
}
