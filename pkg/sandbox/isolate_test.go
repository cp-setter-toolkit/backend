package sandbox_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/cp-setter-toolkit/cp-setter-toolkit/pkg/sandbox"
	"github.com/stretchr/testify/assert"
)

func TestIsolate_Run(t *testing.T) {
	testcases := []struct {
		name     string
		config   sandbox.RunConfig
		cmd      []string
		expected sandbox.Verdict
	}{
		{
			name: "sh_echo",
			config: sandbox.RunConfig{
				RunId: "sh_echo",
			},
			cmd:      []string{"/bin/sh", "-c", "echo \"nigga\""},
			expected: sandbox.VerdictOK,
		},
		{
			name: "sh_yes",
			config: sandbox.RunConfig{
				RunId:     "sh_yes",
				TimeLimit: 50 * time.Millisecond,
			},
			cmd:      []string{"/bin/sh", "-c", "yes"},
			expected: sandbox.VerdictTL,
		},
	}

	sb, err := sandbox.NewIsolate(777, slog.Default())
	assert.NoError(t, err, "failed to create sandbox")
	for _, tc := range testcases {
		err := sb.Init(context.Background())
		assert.NoError(t, err, "failed to init sandbox")

		stat, err := sb.Run(context.Background(), tc.config, tc.cmd[0], tc.cmd[1:]...)
		assert.NoError(t, err, "failed to run command")

		assert.Equal(t, tc.expected, stat.Verdict)
		err = sb.Cleanup(context.Background())
		assert.NoError(t, err, "failed to cleanup sandbox")
	}
}
