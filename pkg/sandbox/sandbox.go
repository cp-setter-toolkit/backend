package sandbox

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/cp-setter-toolkit/cp-setter-toolkit/pkg/memory"
	"github.com/spf13/afero"
)

// ErrorSandboxNotInitialized is returned when the sandbox is not initialized.
var ErrorSandboxNotInitialized = errors.New("sandbox not initialized")

type BindingOpt string

// https://www.ucw.cz/moe/isolate.1.html#_directory_rules
const (
	RW     BindingOpt = "rw"
	DEV    BindingOpt = "dev"
	NOEXEC BindingOpt = "noexec"
	MAYBE  BindingOpt = "maybe"
	FS     BindingOpt = "fs"
	TMP    BindingOpt = "tmp"
	NOREC  BindingOpt = "norec"
)

type DirBinding struct {
	Inside  string
	Outside string
	Options []BindingOpt
}

// RunConfig is the configuration for running a program in a sandbox.
// For more information, see the man page of `isolate`.
type RunConfig struct {
	RunId string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	MemLimit  memory.Amount
	TimeLimit time.Duration
	MaxProcs  int

	InheritEnv bool
	Env        []string

	Bindings []DirBinding
	WorkDir  string

	Args []string
}

// Sandbox is used to run programs in a sandbox.
type Sandbox interface {
	Name() string
	Pwd() string
	Init(ctx context.Context) error
	afero.Fs
	Run(ctx context.Context, config RunConfig, name string, args ...string) (*Status, error)
	RunFile(ctx context.Context, config RunConfig, file File, args ...string) (*Status, error)
	Cleanup(ctx context.Context) error
}

// Provider is used to manage sandboxes.
type Provider interface {
	Pop() (Sandbox, error)
	Push(sb Sandbox)
}

// File is a named io.ReadCloser.
type File interface {
	Name() string
	io.ReadCloser
}

func CopyFile(fs afero.Fs, src File) error {
	dst, err := fs.Create(src.Name())
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return src.Close()
}
