package sandbox

import (
	"os"
	"path/filepath"
	"time"

	"github.com/thepluck/cp-setter-toolkit/helper/errors"
)

// Sandbox is an interface for running programs in a sandbox.
type Sandbox interface {
	Id() string
	Run(input *Input) (*Output, error)
}

// Input is the input to the sandbox.
type Input struct {
	Command     string        // Command to run
	Args        []string      // Additional arguments
	Files       []File        // Additional files to copy
	Excutables  []File        // Additional executables to copy
	TimeLimit   time.Duration // Time limit
	MemoryLimit int           // Memory limit in KBs
	Stdin       []byte        // Input to the program
	OutputFiles []string      // Files to copy back
}

// Output is the output of the sandbox.
type Output struct {
	Success     bool          // Whether the program ran successfully
	TimeElapsed time.Duration // Time taken by the program
	MemoryUsed  int           // Memory used by the program
	Stdout      []byte        // Output of the program
	Stderr      []byte        // Error output of the program
	Status      string        // Status of the program
	Files       []File        // Files copied back
}

// File contains the name and data of a file.
type File struct {
	Name string
	Data []byte
}

// CopyTo copies the files of Input to a directory.
func (i *Input) CopyTo(dir string) error {
	for _, f := range i.Files {
		if err := os.WriteFile(filepath.Join(dir, f.Name), f.Data, 0666); err != nil {
			return errors.Wrapf(err, "copying file %s", f.Name)
		}
	}
	for _, f := range i.Excutables {
		if err := os.WriteFile(filepath.Join(dir, f.Name), f.Data, 0777); err != nil {
			return errors.Wrapf(err, "copying file %s", f.Name)
		}
	}
	return nil
}

// CopyFrom copies the files from a directory.
func (i *Input) CopyFrom(dir string) ([]File, error) {
	files := make([]File, len(i.OutputFiles))
	for j, name := range i.OutputFiles {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, errors.Wrapf(err, "reading file %s", name)
		}
		files[j] = File{Name: name, Data: data}
	}
	return files, nil
}
