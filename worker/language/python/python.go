package python

import (
	"time"

	"github.com/thepluck/cp-setter-toolkit/helper/errors"
	"github.com/thepluck/cp-setter-toolkit/worker/language"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox"
)

type Python struct {
	name string
	cmd  string
}

func NewPython(name string, cmd string) *Python {
	return &Python{name, cmd}
}

func (l Python) Name() string {
	return l.name
}

func (l Python) FileExtensions() []string {
	return []string{".py"}
}

func (l Python) Compile(s sandbox.Sandbox, files []sandbox.File) ([]sandbox.File, error) {
	input := &sandbox.Input{
		Command:     l.cmd,
		Args:        []string{"-m", "compileall", ".", "-q", "-b"},
		Files:       files,
		TimeLimit:   5 * time.Second,
		MemoryLimit: 262144,
	}

	for _, f := range files {
		input.OutputFiles = append(input.OutputFiles, f.Name+"c")
	}

	output, err := s.Run(input)
	if err != nil {
		return nil, errors.Wrap(err, "sandbox error")
	}

	if !output.Success {
		return nil, errors.Errorf("compilation failed: %s", output.Stderr)
	}

	return output.Files, nil
}

func (l Python) Run(s sandbox.Sandbox, files []sandbox.File, tl time.Duration, ml int, stdin []byte) (*sandbox.Output, error) {
	input := &sandbox.Input{
		Command:     l.cmd,
		Args:        []string{files[0].Name},
		Files:       files,
		TimeLimit:   tl,
		MemoryLimit: ml,
		Stdin:       stdin,
	}

	return s.Run(input)
}

var Python27 = NewPython("Python 2.7", "/usr/local/bin/python2.7")
var Python310 = NewPython("Python 3.10", "/usr/local/bin/python3.10")
var PyPy27 = NewPython("PyPy 2.7", "/usr/local/pypy2.7-v7.3.16-linux64/bin/pypy2.7")
var PyPy310 = NewPython("PyPy 3.10", "/usr/local/pypy3.10-v7.3.16-linux64/bin/pypy3.10")

func init() {
	language.DefaultRegistry.Register(Python27)
	language.DefaultRegistry.Register(Python310)
	language.DefaultRegistry.Register(PyPy27)
	language.DefaultRegistry.Register(PyPy310)
}
