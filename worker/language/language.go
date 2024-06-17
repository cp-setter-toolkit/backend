package language

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/thepluck/cp-setter-toolkit/worker/sandbox"
)

// Language is the interface for programming languages.
type Language interface {
	Name() string
	FileExtensions() []string
	Compile(s sandbox.Sandbox, files []sandbox.File) ([]sandbox.File, error)
	Run(s sandbox.Sandbox, files []sandbox.File, tl time.Duration, ml int, stdin []byte) (*sandbox.Output, error)
}

// BaseNameWithoutExt returns the base name of a path without the extension.
func BaseNameWithoutExt(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
