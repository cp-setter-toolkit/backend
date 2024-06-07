package sandbox

type ChanProvider struct {
	sandboxes chan Sandbox
}

func NewChanProvider() *ChanProvider {
	return &ChanProvider{make(chan Sandbox, 100)}
}

func (prov *ChanProvider) Pop() (Sandbox, error) {
	return <-prov.sandboxes, nil
}

func (prov *ChanProvider) Push(sb Sandbox) {
	prov.sandboxes <- sb
}
