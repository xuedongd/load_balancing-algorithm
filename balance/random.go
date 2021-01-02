package balance

import (
	"errors"
	"math/rand"
)

func init() {
	RegisterBalancer("random", &randomBalancer{})
}

type randomBalancer struct {
}

func (p *randomBalancer) DoBalance(insts []*Instance, key ...string) (inst *Instance, err error) {
	if len(insts) == 0 {
		err = errors.New("NO Instance!")
		return
	}
	lens := len(insts)
	index := rand.Intn(lens)
	inst = insts[index]
	return
}
