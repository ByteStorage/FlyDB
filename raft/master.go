package master

import (
	"time"
)

func (m *Master) ListenSlave(slaves []Slave) {
	//TODO implement me
	panic("implement me")
}

func (m *Master) WaitForLeader() {
	timeTick := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-timeTick.C:
			ch := m.Raft.LeaderCh()
			select {
			case isLeader := <-ch:
				if isLeader {
					return
				}
			default:
				continue
			}
		}
	}
}

func (m *Master) NewRaft() {
	//判断是新启动的master还是重启的master
}

func (m *Master) StartGrpcServer() {

}

func (m *Master) ListenRequest() {

}
