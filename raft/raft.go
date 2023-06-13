package cluster

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/ByteStorage/FlyDB/lib/dirtree"
	"github.com/ByteStorage/FlyDB/lib/proto"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"strconv"
)

type IndexerType = int8

const (
	DefaultDbDir             = "/tmp/flydb"
	Btree        IndexerType = iota + 1
)

var DefaultOptions = config.Options{
	DirPath:      DefaultDbDir,
	DataFileSize: 256 * 1024 * 1024, // 256MB
	SyncWrite:    false,
	IndexType:    Btree,
}

// Cluster define the cluster of db
type Cluster struct {
	//Master List
	Master []Master
	//Slave List
	Slave []Slave
	//Master List Leader
	Leader string
}

type Master struct {
	//grpc server
	proto.MasterGrpcServiceServer
	//ID
	ID string
	//Addr
	Addr string
	//Master List
	Peers []string
	//Slave List
	Slave []Slave
	//Heartbeat
	Heartbeat map[string]string
	//Filename to node,key is filename,value is node
	FilenameToNode map[string]string
	//Dir Tree
	DirTree *dirtree.DirTree
	//Raft
	Raft raft.Raft
	//Raft Log
	RaftLog *boltdb.BoltStore
	c       *Cluster
}

type Slave struct {
	//grpc server
	proto.SlaveGrpcServiceServer
	//ID
	ID string
	//Addr
	Addr string
	//Master List Leader
	Leader string
	//Slave List
	Peers []string
	//DB
	DB *engine.DB
	//Slave Message
	SlaveMessage SlaveMessage
	//work pool
	WorkPool chan struct{}
}

type SlaveMessage struct {
	UsedDisk   uint64
	UsedMem    uint64
	CpuPercent float32
	TotalMem   uint64
	TotalDisk  uint64
}

// FSMSnapshot use to store the snapshot of the FSM
type FSMSnapshot struct {
}

// NewRaftCluster create a new raft db cluster
func NewRaftCluster(masterList []string, slaveList []string) *Cluster {
	masters := make([]Master, len(masterList))
	slaves := make([]Slave, len(slaveList))
	for i, slave := range slaveList {
		db, err := engine.NewDB(DefaultOptions)
		if err != nil {
			panic(err)
		}
		slaves[i] = Slave{
			ID:    strconv.Itoa(i),
			Addr:  slave,
			Peers: slaveList,
			DB:    db,
		}
	}
	for i, master := range masterList {
		masters[i] = Master{
			ID:             strconv.Itoa(i),
			Addr:           master,
			Peers:          masterList,
			Slave:          slaves,
			Heartbeat:      make(map[string]string),
			FilenameToNode: make(map[string]string),
			DirTree:        dirtree.NewDirTree(),
		}
	}
	c := &Cluster{
		Master: masters,
		Slave:  slaves,
	}
	c.startMasters()
	c.startSlaves()
	panic("implement me")
}

func (c *Cluster) startMasters() {
	for _, m := range c.Master {
		m.c = c
		//start grpc server
		m.StartGrpcServer()
		//start raft
		m.NewRaft()
		//wait for leader
		m.WaitForLeader()
		//add slave or delete slave
		go m.ListenSlave(c.Slave)
		//listen user request, by wal
		go m.ListenRequest()
	}
}

func (c *Cluster) startSlaves() {
	for _, s := range c.Slave {
		//start grpc server
		s.StartGrpcServer()
		//register to master
		s.RegisterToMaster()
		//heartbeat
		go s.SendHeartbeat()
		//listen leader
		go s.ListenLeader()
		//update slave message
		go s.UpdateSlaveMessage()
	}

}
