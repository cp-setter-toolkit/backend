package sandbox

import (
	"time"

	"github.com/cp-setter-toolkit/cp-setter-toolkit/pkg/memory"
)

// Verdict is the overall result of running a program.
type Verdict int

const (
	VerdictOK Verdict = 1 << iota
	VerdictTL
	VerdictML
	VerdictRE
	VerdictCE
	VerdictIE
)

func (v Verdict) String() string {
	switch v {
	case VerdictOK:
		return "OK"
	case VerdictTL:
		return "TL"
	case VerdictML:
		return "ML"
	case VerdictRE:
		return "RE"
	case VerdictCE:
		return "CE"
	case VerdictIE:
		return "IE"
	}
	return "UNKNOWN"
}

type Status struct {
	Verdict  Verdict
	Signal   int
	Time     time.Duration
	Memory   memory.Amount
	ExitCode int
}

func IEStatus() *Status {
	return &Status{Verdict: VerdictIE}
}
