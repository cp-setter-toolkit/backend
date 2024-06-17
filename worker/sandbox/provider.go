package sandbox

type SandboxQueue struct {
	sandboxes chan Sandbox
}

func NewSandboxQueue() *SandboxQueue {
	return &SandboxQueue{make(chan Sandbox, 100)}
}

func (sq *SandboxQueue) Pop() (Sandbox, error) {
	return <-sq.sandboxes, nil
}

func (sq *SandboxQueue) Push(s Sandbox) {
	sq.sandboxes <- s
}
