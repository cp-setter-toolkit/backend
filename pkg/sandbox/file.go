package sandbox

import (
	"bytes"
	"context"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/cp-setter-toolkit/backend/pkg/memory"
)

// CreateFile creates a file inside sandbox with the given content.
func CreateFile(fs Fs, c File) error {
	if err := syscall.Unlink(fs.Path(c.Name)); err != nil && !os.IsNotExist(err) {
		return err
	}

	f, err := fs.Create(c.Name)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, c.Source); err != nil {
		return err
	}
	return c.Source.Close()
}

// ExtractContent extracts content from a file inside sandbox.
func ExtractContent(fs Fs, name string) (*File, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &File{name, io.NopCloser(bytes.NewBuffer(c))}, nil
}

// RunBinary runs a binary inside sandbox.
func RunBinary(ctx context.Context, sb Sandbox, bin File, stdin io.Reader, stdout io.Writer, tl time.Duration, ml memory.Amount) (*Status, error) {
	stat := Status{Verdict: VerdictIE}

	if err := CreateFile(sb, bin); err != nil {
		return &stat, err
	}

	if err := sb.MakeExecutable(bin.Name); err != nil {
		return &stat, err
	}

	config := RunConfig{
		Stdin: 	stdin,
		Stdout: stdout,
		TimeLimit: tl,
		MemoryLimit: ml,
	}
	return sb.Run(ctx, config, bin.Name)
}
