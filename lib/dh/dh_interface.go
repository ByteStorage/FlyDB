package dh

type DataHandler interface {
	//AssignData assign data to slave list, return a map with single slave address as key and data as value
	AssignData(data []byte, slaveAddrList []string) (map[string][]byte, error)
}
