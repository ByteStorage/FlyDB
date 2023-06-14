package policy

type Policy interface {
	// AssignSlave provide multiple strategies for num node elections, such as polling and randomization
	AssignSlave(num int, slaveAddrList []string) ([]string, error)
}
