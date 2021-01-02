package balance

import (
	"fmt"
)

type BalancerMgr struct {
	allBalance map[string]Balancer
}

var mgr = BalancerMgr{
	allBalance: make(map[string]Balancer),
}

func (p BalancerMgr) registerBalancer(name string, balancer Balancer) {
	p.allBalance[name] = balancer
}

func RegisterBalancer(name string, balancer Balancer) {
	mgr.registerBalancer(name, balancer)
}

func DoBalance(name string, insts []*Instance) (inst *Instance, err error) {
	balancer, ok := mgr.allBalance[name]
	if !ok {
		err = fmt.Errorf("Not found %s balancer\n", name)
		return
	}
	fmt.Printf("use %s balancer\n", name)
	inst, err = balancer.DoBalance(insts)
	return
}
