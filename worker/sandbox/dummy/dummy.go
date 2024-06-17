package dummy

import (
	"errors"
	"log"

	"github.com/thepluck/cp-setter-toolkit/worker/sandbox"
)

var warned bool

type Dummy struct {
	id string
}

func NewDummy(id string) *Dummy {
	if !warned {
		log.Println("Warning: `Dummy` sandbox is selected. WE ARE NOT RESPONSIBLE FOR ANY DAMAGE CAUSED BY THIS SANDBOX.")
	}
	return &Dummy{id: id}
}

func (s *Dummy) Id() string {
	return s.id
}

func (s *Dummy) Run(input *sandbox.Input) (*sandbox.Output, error) {
	return nil, errors.New("not implemented")
}
