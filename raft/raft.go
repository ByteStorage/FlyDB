package master

import (
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/qishenonly/flydb/utils"
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
	DirTree *utils.DirTree
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
	DirTree *utils.DirTree
	//Heartbeat
	Heartbeat map[string]string
	//Leader
	Leader string
	//Filename to node,key is filename,value is node
	FilenameToNode map[string]string
}

// NewRaftCluster create a new raft db cluster
func NewRaftCluster(masterList []string, slaveList []string) *Cluster {
	startMasters(masterList)
	startSlaves(slaveList)
	panic("implement me")
}

func startMasters(addr []string) {

}

func startSlaves(addr []string) {

}
