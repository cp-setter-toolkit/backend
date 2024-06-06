package sandbox

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/cp-setter-toolkit/backend/pkg/memory"
)

// ErrorSandboxNotInitialized is returned when the sandbox is not initialized.
var ErrorSandboxNotInitialized = errors.New("sandbox not initialized")

// Fs is the interface for file system operations.
type Fs interface {
	Pwd() string
	Path(name string) string
	Create(name string) (io.WriteCloser, error)
	MakeExecutable(name string) error
	Open(name string) (io.ReadCloser, error)
}

type DirMapping struct {
	Inside  string
	Outside string
}

type DirOpt string

const (
	AllowReadWrite  DirOpt = "rw"
	DisallowExecute DirOpt = "noexec"
)

// RunConfig is the configuration for running a program in a sandbox.
// For more information, see the man page of `isolate`.
type RunConfig struct {
	RunId string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	MemoryLimit memory.Amount
	TimeLimit   time.Duration
	MaxProcs    int

	InheritEnv bool
	Env        []string

	DirMappings []DirMapping
	DirOpts     []DirOpt
	WorkingDir  string

	Args []string
}

// Sandbox is used to run programs in a sandbox.
type Sandbox interface {
	Id() string
	Init(ctx context.Context) error
	Fs
	Run(ctx context.Context, config RunConfig, cmd string, args ...string) (*Status, error)
	Cleanup(ctx context.Context) error
}

// Provider is used to manage sandboxes.
type Provider interface {
	Pop() (Sandbox, error)
	Push(sb Sandbox)
}

// File is a named source of data.
type File struct {
	Name   string
	Source io.ReadCloser
}
