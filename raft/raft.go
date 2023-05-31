package master

import (
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/qishenonly/flydb/lib/dirtree"
)

// Cluster define the cluster of db
type Cluster struct {
	//Master List
	Master []string
	//Slave List
	Slave []string
	//Raft
	Raft raft.Raft
	//Raft Log
	RaftLog *boltdb.BoltStore
	//Dir Tree
	DirTree *dirtree.DirTree
	//Heartbeat
	Heartbeat map[string]string
	//Leader
	Leader string
	//Filename to node,key is filename,value is node
	FilenameToNode map[string]string
	//Current Node
	CurrentNode string
	//ID
	ID string
}

// FSMSnapshot use to store the snapshot of the FSM
type FSMSnapshot struct {
	//Master List
	Master []string
	//Slave List
	Slave []string
	//Raft
	Raft raft.Raft
	//Raft Log
	RaftLog *boltdb.BoltStore
	//Dir Tree
	DirTree *dirtree.DirTree
	//Heartbeat
	Heartbeat map[string]string
	//Leader
	Leader string
	//Filename to node,key is filename,value is node
	FilenameToNode map[string]string
}

// NewRaftCluster create a new raft db cluster
func NewRaftCluster(masterList []string, slaveList []string) *Cluster {
	c := &Cluster{
		Master: masterList,
		Slave:  slaveList,
	}
	c.startMasters()
	c.startSlaves()
	panic("implement me")
}

func (c *Cluster) startMasters() {
	//启动grpc服务

	//启动raft服务

	//等待leader选举

	//监听leader变化

	//对slave进行下线或者上线

	//监听用户请求，通过wal日志处理

}

func (c *Cluster) startSlaves() {
	//启动grpc服务

	//向master注册

	//发起心跳

	//监听leader变化

}
