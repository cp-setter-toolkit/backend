package python

import (
	"context"
	"io"

	"github.com/cp-setter-toolkit/backend/pkg/language"
	"github.com/cp-setter-toolkit/backend/pkg/sandbox"
)

type Python struct {
	cmd   string
	name string
}

func NewPython(cmd string, name string) *Python {
	return &Python{cmd, name}
}

func (l Python) Id() string {
	return l.cmd
}

func (l Python) Name() string {
	return l.name
}

func (l Python) FileExtensions() []string {
	return []string{".py"}
}

func (l Python) Compile(ctx context.Context, sb sandbox.Sandbox, files []sandbox.File, stderr io.Writer) (sandbox.File, error) {
	return files[0], nil

}

func (l Python) Run(ctx context.Context, sb sandbox.Sandbox, config sandbox.RunConfig, file sandbox.File) (*sandbox.Status, error) {
	if err := sandbox.CopyFile(sb, file); err != nil {
		return nil, err
	}

	return sb.Run(ctx, config, "/usr/bin/" + l.cmd, file.Name())
}

var Python27 = NewPython("python2.7", "Python 2.7")
var Python310 = NewPython("python3.10", "Python 3.10")
var PyPy27 = NewPython("pypy", "PyPy 2.7")
var PyPy310 = NewPython("pypy3", "PyPy 3.10")

func init() {
	language.DefaultRegistry.Register(NewPython("python", "Python"))
}
