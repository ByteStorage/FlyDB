package dh

import cluster "github.com/ByteStorage/flydb/raft"

type ReplicationEncoder struct {
	//TODO:需要获取Slave集群的磁盘，CPU信息，根据这些信息实现合适的分配策略，有其他思路亦可提出来
}

func NewReplicationEncoder([]cluster.Slave) {
	panic("implement me")
}

func (r *ReplicationEncoder) AssignData(data []byte, slaveAddrList []string) (map[string][]byte, error) {
	panic("implement me")
}
