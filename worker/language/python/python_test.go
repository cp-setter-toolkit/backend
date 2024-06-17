package python_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thepluck/cp-setter-toolkit/worker/language/python"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox/isolate"
)

func TestRun(t *testing.T) {
	s := isolate.New("666")
	headerFile, err := os.ReadFile("./testdata/header.py")
	assert.NoError(t, err, "failed to open header.py")
	solutionFile, err := os.ReadFile("./testdata/solution.py")
	assert.NoError(t, err, "failed to open solution.py")
	testcases := []struct {
		name     string
		files    []sandbox.File
		expected []byte
	}{
		{
			name: "full",
			files: []sandbox.File{
				{Name: "solution.py", Data: solutionFile},
				{Name: "header.py", Data: headerFile},
			},
			expected: []byte("Hello, World!\nThe answer is 42.\n"),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			files, err := python.PyPy310.Compile(s, tc.files)
			assert.NoError(t, err, "failed to compile")

			stdin := []byte("42\n")
			output, err := python.PyPy310.Run(s, files, 5*time.Second, 262144, stdin)
			assert.NoError(t, err, "failed to run")
			assert.Equal(t, tc.expected, output.Stdout, "unexpected output")
		})
	}
}
