package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cp-setter-toolkit/cp-setter-toolkit/pkg/memory"
	"github.com/spf13/afero"
)

var IsolateRoot = "/var/local/lib/isolate"
var IsolateMetafilePattern = "isolate-metafile*"

type Isolate struct {
	id int
	afero.Fs
	Logger *slog.Logger
	inited bool
}

func (sb *Isolate) Name() string {
	return "isolate" + strconv.Itoa(sb.id)
}

func (sb *Isolate) Pwd() string {
	return filepath.Join(IsolateRoot, strconv.Itoa(sb.id), "box")
}

func NewIsolate(id int, logger *slog.Logger) (*Isolate, error) {
	sb := &Isolate{
		id:     id,
		Logger: logger,
	}
	if sb.Logger == nil {
		sb.Logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	sb.Logger = sb.Logger.With(slog.String("sandbox", sb.Name()))
	return sb, nil
}

func (sb *Isolate) Init(ctx context.Context) error {
	// Cleanup the sandbox if it was not cleaned up properly
	if err := sb.Cleanup(ctx); err != nil {
		return err
	}

	cmd := []string{"isolate", "--cg", "-b", strconv.Itoa(sb.id), "--init"}
	sb.Logger.Info("🏗️\trunning init", "cmd", cmd)
	sb.inited = true
	sb.Fs = afero.NewBasePathFs(afero.NewOsFs(), sb.Pwd())
	return exec.Command(cmd[0], cmd[1:]...).Run()
}

func (sb *Isolate) buildArgs(config RunConfig) ([]string, error) {
	args := []string{"--cg", "-b", strconv.Itoa(sb.id)}
	args = append(args, "--stack=268435")
	if config.MemLimit > 0 {
		args = append(args, fmt.Sprintf("--cg-mem=%d", config.MemLimit/memory.KB))
	}
	if config.TimeLimit > 0 {
		sec := float64(config.TimeLimit/time.Millisecond) / 1000
		args = append(args, "--time="+strconv.FormatFloat(sec, 'f', 3, 64))
	}
	if config.MaxProcs > 0 {
		args = append(args, fmt.Sprintf("--processes=%d", config.MaxProcs))
	} else {
		args = append(args, "--processes")
	}
	if config.InheritEnv {
		args = append(args, "--full-env")
	}
	for _, env := range config.Env {
		args = append(args, "--env="+env)
	}
	for _, db := range config.Bindings {
		arg := fmt.Sprintf("--dir=%s=%s", db.Inside, db.Outside)
		for _, opt := range db.Options {
			arg += ":" + string(opt)
		}
		args = append(args, arg)
	}
	args = append(args, config.Args...)
	return args, nil
}

func (sb *Isolate) Run(ctx context.Context, config RunConfig, name string, args ...string) (*Status, error) {
	if !sb.inited {
		return IEStatus(), ErrorSandboxNotInitialized
	}
	if config.RunId != "" {
		sb.Logger = sb.Logger.With(slog.String("run-id", config.RunId))
	}

	isoArgs, err := sb.buildArgs(config)
	if err != nil {
		return IEStatus(), fmt.Errorf("failed to build isolate args: %w", err)
	}

	metafile, err := os.CreateTemp("", IsolateMetafilePattern)
	if err != nil {
		return IEStatus(), fmt.Errorf("failed to create metafile: %w", err)
	}
	defer metafile.Close()
	defer os.Remove(filepath.Join(os.TempDir(), metafile.Name()))
	isoArgs = append(isoArgs, "--meta="+metafile.Name())

	isoArgs = append(isoArgs, "--run", "-s", "--", name)
	isoArgs = append(isoArgs, args...)

	sb.Logger.Info("🛠️\tbuilt args", "args", isoArgs)

	cmd := exec.Command("isolate", isoArgs...)
	cmd.Stdin = config.Stdin
	cmd.Stdout = config.Stdout
	cmd.Stderr = config.Stderr
	cmd.Dir = config.WorkDir
	_ = cmd.Run()

	stat := Status{
		Verdict: VerdictOK,
	}
	sc := bufio.NewScanner(metafile)
	for sc.Scan() {
		lst := strings.Split(sc.Text(), ":")
		switch lst[0] {
		case "max-rss":
		case "cg-mem":
			mem, _ := strconv.Atoi(lst[1])
			stat.Memory += memory.Amount(mem) * memory.KiB
		case "time":
			tmp, _ := strconv.ParseFloat(lst[1], 32)
			stat.Time = time.Duration(tmp*1000) * time.Millisecond
		case "status":
			switch lst[1] {
			case "TO":
				stat.Verdict = VerdictTL
			case "RE":
				stat.Verdict = VerdictRE
			case "SG":
				stat.Verdict = VerdictRE
			case "XX":
				stat.Verdict = VerdictIE
			}
		case "exitcode":
			stat.ExitCode, _ = strconv.Atoi(lst[1])
		}
	}

	if err = sc.Err(); err != nil {
		return IEStatus(), fmt.Errorf("failed to read metafile: %w", err)
	}

	sb.Logger.Info("🏁\trun finished", "status", stat)
	return &stat, nil
}

func (sb *Isolate) RunFile(ctx context.Context, config RunConfig, file File, args ...string) (*Status, error) {
	if err := CopyFile(sb, file); err != nil {
		return IEStatus(), fmt.Errorf("failed to create file in sandbox: %w", err)
	}

	if err := sb.Chmod(file.Name(), 0755); err != nil {
		return IEStatus(), fmt.Errorf("failed to chmod file in sandbox: %w", err)
	}

	return sb.Run(ctx, config, file.Name(), args...)
}

func (sb *Isolate) Cleanup(ctx context.Context) error {
	cmd := []string{"isolate", "--cg", "-b", strconv.Itoa(sb.id), "--cleanup"}

	sb.Logger.Info("🧹\trunning cleanup", "cmd", cmd)
	sb.inited = false
	sb.Fs = nil
	return exec.Command(cmd[0], cmd[1:]...).Run()
}
