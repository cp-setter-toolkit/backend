package cpp_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thepluck/cp-setter-toolkit/worker/language/cpp"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox/isolate"
)

func TestRun(t *testing.T) {
	s := isolate.New("888")
	exampleFile, err := os.ReadFile("./testdata/example.h")
	assert.NoError(t, err, "failed to open example.h")
	graderFile, err := os.ReadFile("./testdata/grader.cpp")
	assert.NoError(t, err, "failed to open grader.cpp")
	solutionFile, err := os.ReadFile("./testdata/solution.cpp")
	assert.NoError(t, err, "failed to open solution.cpp")
	testcases := []struct {
		name     string
		files    []sandbox.File
		expected []byte
	}{
		{
			name: "full",
			files: []sandbox.File{
				{Name: "solution.cpp", Data: solutionFile},
				{Name: "example.h", Data: exampleFile},
				{Name: "grader.cpp", Data: graderFile},
			},
			expected: []byte("Hello, World!\nThe answer is 42.\n"),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			files, err := cpp.Cpp17.Compile(s, tc.files)
			assert.NoError(t, err, "failed to compile")

			stdin := []byte("42\n")
			output, err := cpp.Cpp17.Run(s, files, 5*time.Second, 262144, stdin)
			assert.NoError(t, err, "failed to run")
			assert.Equal(t, tc.expected, output.Stdout, "unexpected output")
		})
	}
}

func TestCompileError(t *testing.T) {
	s := isolate.New("999")
	exampleFile, err := os.ReadFile("./testdata/example.h")
	assert.NoError(t, err, "failed to open example.h")
	graderFile, err := os.ReadFile("./testdata/grader.cpp")
	assert.NoError(t, err, "failed to open grader.cpp")
	testcases := []struct {
		name     string
		files    []sandbox.File
		expected []byte
	}{
		{
			name: "full",
			files: []sandbox.File{
				{Name: "example.h", Data: exampleFile},
				{Name: "grader.cpp", Data: graderFile},
			},
			expected: []byte("Hello, World!\nThe answer is 42.\n"),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := cpp.Cpp17.Compile(s, tc.files)
			assert.Error(t, err, "expected to fail")
		})
	}
}
