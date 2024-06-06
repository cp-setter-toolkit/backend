package sandbox

type ChanProvider struct {
    sandboxes chan Sandbox
}

func NewChanProvider() *ChanProvider {
    return &ChanProvider{make(chan Sandbox, 100)}
}

func (p *ChanProvider) Pop() (Sandbox, error) {
    return <-p.sandboxes, nil
}

func (p *ChanProvider) Push(s Sandbox) {
    p.sandboxes <- s
}
