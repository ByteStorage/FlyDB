package cluster

import (
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/qishenonly/flydb"
	"github.com/qishenonly/flydb/lib/dirtree"
	"github.com/qishenonly/flydb/lib/proto"
	"strconv"
)

type IndexerType = int8

const (
	DefaultDbDir             = "/tmp/flydb"
	Btree        IndexerType = iota + 1
)

var DefaultOptions = flydb.Options{
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
	DB *flydb.DB
}

// FSMSnapshot use to store the snapshot of the FSM
type FSMSnapshot struct {
}

// NewRaftCluster create a new raft db cluster
func NewRaftCluster(masterList []string, slaveList []string) *Cluster {
	masters := make([]Master, len(masterList))
	slaves := make([]Slave, len(slaveList))
	for i, slave := range slaveList {
		db, err := flydb.NewFlyDB(DefaultOptions)
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
		//启动grpc服务
		m.StartGrpcServer()
		//启动raft服务
		m.NewRaft()
		//等待leader选举
		m.WaitForLeader()
		//对slave进行下线或者上线
		go m.ListenSlave(c.Slave)
		//监听用户请求，通过wal日志处理
		go m.ListenRequest()
	}
}

func (c *Cluster) startSlaves() {
	for _, s := range c.Slave {
		//启动grpc服务
		s.StartGrpcServer()
		//向master注册
		s.RegisterToMaster()
		//发起心跳
		go s.SendHeartbeat()
		//监听leader变化
		go s.ListenLeader()
	}

}
