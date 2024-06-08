package cpp

import (
	"context"
	"io"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/cp-setter-toolkit/cp-setter-toolkit/pkg/language"
	"github.com/cp-setter-toolkit/cp-setter-toolkit/pkg/memory"
	"github.com/cp-setter-toolkit/cp-setter-toolkit/pkg/sandbox"
	"github.com/spf13/afero"
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

func (l Cpp) Compile(ctx context.Context, sb sandbox.Sandbox, files []sandbox.File, stderr io.Writer) (sandbox.File, error) {
	args := slices.Clone(l.args)
	for _, file := range files {
		if err := sandbox.CopyFile(sb, file); err != nil {
			return nil, err
		}
		if !strings.HasSuffix(file.Name(), ".h") {
			args = append(args, path.Base(file.Name()))
		}
	}

	config := sandbox.RunConfig{
		Stdout:     stderr,
		Stderr:     stderr,
		MemLimit:   256 * memory.MiB,
		TimeLimit:  10 * time.Second,
		MaxProcs:   200,
		InheritEnv: true,
		WorkDir:    sb.Pwd(),
	}

	if _, err := sb.Run(ctx, config, "/usr/bin/g++", args...); err != nil {
		return nil, err
	}
	data, err := afero.ReadFile(sb, "a.out")
	if err != nil {
		return nil, err
	}
	return sandbox.NewBuffer("a.out", data), nil
}

func (l Cpp) Run(ctx context.Context, sb sandbox.Sandbox, config sandbox.RunConfig, file sandbox.File) (*sandbox.Status, error) {
	return sb.RunFile(ctx, config, file)
}

var DefaultCompileArgs = []string{"-static", "-DONLINE_JUDGE", "-O2"}

var Cpp11 = NewCpp("cpp11", "C++11", append(DefaultCompileArgs, "-std=c++11")...)
var Cpp14 = NewCpp("cpp14", "C++14", append(DefaultCompileArgs, "-std=c++14")...)
var Cpp17 = NewCpp("cpp17", "C++17", append(DefaultCompileArgs, "-std=c++17")...)
var Cpp20 = NewCpp("cpp20", "C++20", append(DefaultCompileArgs, "-std=c++20")...)

func init() {
	language.DefaultRegistry.Register(Cpp11)
	language.DefaultRegistry.Register(Cpp14)
	language.DefaultRegistry.Register(Cpp17)
	language.DefaultRegistry.Register(Cpp20)
}
