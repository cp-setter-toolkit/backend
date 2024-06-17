package cpp

import (
	"strings"
	"time"

	"github.com/thepluck/cp-setter-toolkit/helper/errors"
	"github.com/thepluck/cp-setter-toolkit/worker/language"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox"
)

type Cpp struct {
	name string
	args []string
}

func NewCpp(name string, args ...string) *Cpp {
	return &Cpp{name, args}
}

func (l Cpp) Name() string {
	return l.name
}

func (l Cpp) FileExtensions() []string {
	return []string{".cpp", ".cc", ".cxx", ".c++"}
}

func (l Cpp) Compile(s sandbox.Sandbox, files []sandbox.File) ([]sandbox.File, error) {
	basename := language.BaseNameWithoutExt(files[0].Name)
	input := &sandbox.Input{
		Command:     "/usr/bin/g++",
		Args:        append(l.args, "-o", basename),
		Files:       files,
		TimeLimit:   5 * time.Second,
		MemoryLimit: 262144,
		OutputFiles: []string{basename},
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name, ".h") {
			input.Args = append(input.Args, f.Name)
		}
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

func (l Cpp) Run(s sandbox.Sandbox, files []sandbox.File, tl time.Duration, ml int, stdin []byte) (*sandbox.Output, error) {
	input := &sandbox.Input{
		Command:     "./" + language.BaseNameWithoutExt(files[0].Name),
		Files:       files[1:], // Exclude the compiled file
		Excutables:  files[:1], // Include the compiled file
		TimeLimit:   tl,
		MemoryLimit: ml,
		Stdin:       stdin,
	}

	return s.Run(input)
}

var DefaultCompileArgs = []string{"-static", "-O2"}

var Cpp11 = NewCpp("C++11", append(DefaultCompileArgs, "-std=c++11")...)
var Cpp14 = NewCpp("C++14", append(DefaultCompileArgs, "-std=c++14")...)
var Cpp17 = NewCpp("C++17", append(DefaultCompileArgs, "-std=c++17")...)
var Cpp20 = NewCpp("C++20", append(DefaultCompileArgs, "-std=c++20")...)

func init() {
	language.DefaultRegistry.Register(Cpp11)
	language.DefaultRegistry.Register(Cpp14)
	language.DefaultRegistry.Register(Cpp17)
	language.DefaultRegistry.Register(Cpp20)
}
