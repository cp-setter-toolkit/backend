package isolate

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/thepluck/cp-setter-toolkit/helper/errors"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox"
)

var (
	// Can be overridden by the ISOLATE_CMD environment variable.
	isolateCmd = "isolate"
)

func init() {
	if cmd, ok := os.LookupEnv("ISOLATE_CMD"); ok {
		isolateCmd = cmd
	}
}

type Isolate struct {
	id string
}

// Id returns the id of the sandbox.
func (s *Isolate) Id() string {
	return s.id
}

// Panic if the isolate command is not available.
func checkIsolate() {
	output, err := exec.Command(isolateCmd, "--version").CombinedOutput()
	if err != nil {
		panic(errors.Wrap(err, "trying to run isolate"))
	}
	if !bytes.Contains(output, []byte("The process isolator")) {
		panic(("Wrong isolate command found. Override the ISOLATE_CMD environment variable to set the correct path."))
	}
}

// NewIsolate creates a new Isolate sandbox.
// It panics if the isolate command is not available.
// Set the logger to nil to disable logging.
func New(id string) *Isolate {
	checkIsolate()
	return &Isolate{id: id}
}

// Run runs the given command in the isolate sandbox.
func (s *Isolate) Run(input *sandbox.Input) (*sandbox.Output, error) {
	// Make sure the sandbox is cleaned up after running.
	defer s.cleanup()

	// Initialize the sandbox.
	dirBytes, err := exec.Command(isolateCmd, "--cg", "-b", s.id, "--init").Output()
	if err != nil {
		return nil, errors.Wrap(err, "initializing sandbox")
	}
	dir := filepath.Join(strings.TrimSpace(string(dirBytes)), "box")

	// Copy the files to the sandbox.
	if err := input.CopyTo(dir); err != nil {
		return nil, err
	}

	// Create a metadata file.
	metaFile, err := NewMetaFile(os.TempDir())
	if err != nil {
		return nil, err
	}
	defer os.Remove(metaFile)

	cmd := s.buildCmd(dir, metaFile, input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command.
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
			return nil, errors.WithStack(err)
		}
	}

	// Parse the metadata file.
	output := &sandbox.Output{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}
	if err := parseMetaFile(metaFile, output); err != nil {
		return nil, err
	}

	// Copy the files back.
	output.Files, err = input.CopyFrom(dir)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func parseMetaFile(metaFile string, output *sandbox.Output) error {
	meta, err := Read(metaFile)
	if err != nil {
		return err
	}

	if output.Status, err = meta.String("status"); err == nil {
		output.Success = false
	} else {
		output.Success = true
	}

	if output.MemoryUsed, err = meta.Int("cg-mem"); err != nil {
		return err
	}
	if t, err := meta.Float64("time"); err != nil {
		return err
	} else {
		output.TimeElapsed = time.Duration(t * float64(time.Second))
	}
	return nil
}

func (s *Isolate) buildCmd(dir, metaFile string, input *sandbox.Input) *exec.Cmd {
	// Time limit is in seconds.
	timeLimit := input.TimeLimit.Seconds()

	cmd := exec.Command(
		"isolate",
		"--cg",
		"-b", s.id,
		"--run",
		"-M", metaFile,
		"-t", fmt.Sprintf("%.3f", timeLimit), // Time limit
		"-w", fmt.Sprintf("%.3f", 2*timeLimit+1), // Wall time limit
		"-x", "1.0", // Extra time
		"-f", "262144", // File size limit
		"-p", // Allow multiple processes
		"-s", // Be silent
		"--env=ONLINE_JUDGE=true",
		"--full-env",
		fmt.Sprintf("--cg-mem=%d", input.MemoryLimit), // Total memory limit
		"--",
		input.Command,
	)
	cmd.Args = append(cmd.Args, input.Args...)
	cmd.Dir = dir

	// Pipe the stdin.
	cmd.Stdin = bytes.NewBuffer(input.Stdin)

	return cmd
}

func (s *Isolate) cleanup() {
	_ = exec.Command(isolateCmd, "--cg", "-b", s.id, "--cleanup").Run()
}
