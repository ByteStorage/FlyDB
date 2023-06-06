package master

import (
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/qishenonly/flydb"
	"github.com/qishenonly/flydb/lib/dirtree"
	"strconv"
)

// Cluster define the cluster of db
type Cluster struct {
	//Master List
	Master []Master
	//Slave List
	Slave []Slave
	//Raft
	Raft raft.Raft
	//Raft Log
	RaftLog *boltdb.BoltStore
}

type Master struct {
	//ID
	ID string
	//Addr
	Addr string
	//Master List Leader
	Leader string
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
}

type Slave struct {
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
