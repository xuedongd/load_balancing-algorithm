package balance

import (
	"strconv"
)

type Instance struct {
	host string
	port int
}

func NewInstance(host string, port int) (inst *Instance) {
	inst = &Instance{
		host: host,
		port: port,
	}
	return
}

func (inst *Instance) String() string {
	return inst.host + ":" + strconv.Itoa(inst.port)
}
