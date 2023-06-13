package policy

import (
	"errors"
	"math/rand"
)

type RandomPolicy struct {
}

func NewRandomPolicy() *RandomPolicy {
	return &RandomPolicy{}
}

func (r *RandomPolicy) AssignSlave(num int, slaveAddrList []string) ([]string, error) {
	length := len(slaveAddrList)
	if length < num {
		return nil, errors.New("slaveAddrList length is less than num")
	}
	res := make([]string, 0)
	for i := 0; i < num; i++ {
		rand := rand.Intn(length)
		res = append(res, slaveAddrList[rand])
	}
	return res, nil
}
