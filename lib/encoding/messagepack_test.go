package encoding

import (
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeMessagePack(t *testing.T) {
	var data = &raft.Log{Index: 10, Data: []byte("helloWorld!")}
	var decData raft.Log
	d, err := EncodeMessagePack(data)
	assert.NoError(t, err)
	assert.NotNil(t, d)
	err = DecodeMessagePack(d, &decData)
	assert.NoError(t, err)
	assert.Equal(t, decData, *data)
}
