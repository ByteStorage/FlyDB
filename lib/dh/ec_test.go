package dh

import (
	cluster "github.com/ByteStorage/FlyDB/lib/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEcEncoder(t *testing.T) {
	dataNum := 4
	parityNum := 2

	encoder, err := NewEcEncoder(dataNum, parityNum)
	assert.NoError(t, err)
	assert.NotNil(t, encoder)
}

func TestEcEncoder_AssignData(t *testing.T) {
	// Create a new EcEncoder
	encoder, err := NewEcEncoder(2, 2)
	assert.NoError(t, err)

	//Prepare test data
	data := []byte("test data")
	slaveAddrList := []cluster.Slave{
		{Addr: "slave1"},
		{Addr: "slave2"},
		{Addr: "slave3"},
		{Addr: "slave4"},
	}

	//// Assign data
	assignment, err := encoder.AssignData(data, slaveAddrList)
	assert.NoError(t, err)
	//
	//// Verify the assignment
	expectedAssignment := map[string][]byte{
		"slave1": assignment["slave1"],
		"slave2": assignment["slave2"],
		"slave3": assignment["slave3"],
		"slave4": assignment["slave4"],
	}
	assert.Equal(t, expectedAssignment, assignment)
}
