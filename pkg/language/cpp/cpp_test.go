package cpp_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/cp-setter-toolkit/backend/pkg/language/cpp"
	"github.com/cp-setter-toolkit/backend/pkg/memory"
	"github.com/cp-setter-toolkit/backend/pkg/sandbox"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	sb, err := sandbox.NewIsolate(888, nil)
	assert.NoError(t, err, "failed to create sandbox")
	fs := afero.NewBasePathFs(afero.NewOsFs(), "./testdata")
	exampleFile, err := fs.Open("example.h")
	assert.NoError(t, err, "failed to open example.h")
	graderFile, err := fs.Open("grader.cpp")
	assert.NoError(t, err, "failed to open grader.cpp")
	solutionFile, err := fs.Open("solution.cpp")
	assert.NoError(t, err, "failed to open solution.cpp")
	testcases := []struct {
		name     string
		files    []sandbox.File
		expected string
	}{
		{
			name: "full",
			files: []sandbox.File{
				exampleFile, graderFile,
				solutionFile,
			},
			expected: "Hello, World!\nThe answer is 42.\n",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := sb.Init(context.TODO())
			assert.NoError(t, err, "failed to init sandbox")
			file, err := cpp.Cpp17.Compile(context.TODO(), sb, tc.files, io.Discard)
			assert.NoError(t, err, "failed to compile")

			stdin := bytes.NewBufferString("42\n")
			stdout := &bytes.Buffer{}
			config := sandbox.RunConfig{
				Stdin: stdin,
				Stdout: stdout,
				Stderr: io.Discard,
				MemLimit: 512 * memory.MiB,
				TimeLimit: 1 * time.Second,
			}
			stat, err := cpp.Cpp17.Run(context.TODO(), sb, config, file)
			assert.NoError(t, err, "failed to run command")

			assert.Equal(t, sandbox.VerdictOK, stat.Verdict)
			assert.Equal(t, tc.expected, stdout.String())
		})
	}
}
