package sandbox

import (
	"io"
	"os"
	"path/filepath"
)

// OsFs implements the Fs interface with underlying base path.
// Should be only created with NewOsFs.
type OsFs struct {
	base string
}

func NewOsFs(base string) OsFs {
	return OsFs{base}
}

func (fs OsFs) Pwd() string {
	return fs.base
}

func (fs OsFs) Path(name string) string {
	return filepath.Join(fs.base, name)
}

func (fs OsFs) Create(name string) (io.WriteCloser, error) {
	return os.Create(fs.Path(name))
}

func (fs OsFs) Open(name string) (io.ReadCloser, error) {
	return os.Open(fs.Path(name))
}

func (fs OsFs) MakeExecutable(name string) error {
	return os.Chmod(fs.Path(name), 0755)
}
