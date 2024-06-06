package sandbox

import (
	"context"
	"io"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strconv"
)

var IsolatePath = "../../../isolate"
var IsolateMetafilePath = "isolate.meta"
var TempDir = "/tmp"

type Isolate struct {
	id int
	OsFs
	Logger *slog.Logger
	inited bool
}

func (sb *Isolate) Id() string {
	return "isolate" + strconv.Itoa(sb.id)
}

func NewIsolate(id int, logger *slog.Logger) *Isolate {
	sb := &Isolate{
		id:     id,
		Logger: logger,
	}
	if sb.Logger == nil {
		sb.Logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	sb.Logger = sb.Logger.With(slog.String("sandbox", sb.Id()))
	return sb
}

func (sb *Isolate) Init(ctx context.Context) error {
	if err := sb.Cleanup(ctx); err != nil {
		return err
	}

	cmd := []string{"isolate", "--cg", "-b", strconv.Itoa(sb.id), "--init"}
	sb.Logger.Info("🏗️\trunning init", "cmd", cmd)
	sb.inited = true
	sb.OsFs = NewOsFs(filepath.Join(TempDir, sb.Id(), "box"))
	return exec.Command(cmd[0], cmd[1:]...).Run()
}

func (sb *Isolate) Cleanup(_ context.Context) error {
	cmd := []string{"isolate", "--cg", "-b", strconv.Itoa(sb.id), "--cleanup"}

	sb.Logger.Info("🧹\trunning cleanup", "cmd", cmd)
	sb.inited = false
	sb.OsFs = OsFs{}
	return exec.Command(cmd[0], cmd[1:]...).Run()
}
