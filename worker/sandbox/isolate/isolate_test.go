package isolate_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox"
	"github.com/thepluck/cp-setter-toolkit/worker/sandbox/isolate"
)

func TestIsolate_Run(t *testing.T) {
	testcases := []struct {
		name   string
		input  sandbox.Input
		status string
	}{
		{
			name: "sh_echo",
			input: sandbox.Input{
				Command:     "/bin/sh",
				Args:        []string{"-c", "echo \"nigga\""},
				TimeLimit:   50 * time.Millisecond,
				MemoryLimit: 262144,
				Stdin:       []byte(""),
			},
			status: "",
		},
		{
			name: "sh_yes",
			input: sandbox.Input{
				Command:     "/bin/sh",
				Args:        []string{"-c", "yes"},
				TimeLimit:   50 * time.Millisecond,
				MemoryLimit: 262144,
				Stdin:       []byte(""),
			},
			status: "TO",
		},
	}

	s := isolate.New("777")
	for _, tc := range testcases {
		output, err := s.Run(&tc.input)
		assert.NoError(t, err, "failed to run command")
		assert.Equal(t, tc.status, output.Status)
	}
}
