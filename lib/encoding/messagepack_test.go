package encoding

import (
	"bytes"
	"github.com/hashicorp/go-msgpack/codec"
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
func TestInitMessagePack(t *testing.T) {
	msgPack := InitMessagePack()

	assert.NotNil(t, msgPack.MsgPack)
}

func TestNewMessagePackEncoder(t *testing.T) {
	encoder := NewMessagePackEncoder()

	assert.NotNil(t, encoder.enc)
	assert.NotNil(t, encoder.b)
}

func TestNewMessagePackDecoder(t *testing.T) {
	exampleBytes := []byte("example")
	decoder := NewMessagePackDecoder(exampleBytes)

	assert.NotNil(t, decoder.b)
	assert.NotNil(t, decoder.dec)
}

func TestEncode(t *testing.T) {
	msg := "example message"
	msgPack := InitMessagePack()

	encoded, err := msgPack.Encode(msg)

	assert.NotNil(t, encoded)
	assert.Nil(t, err)
}

func TestMsgPackDecoder_Decode(t *testing.T) {
	msg := "example message"
	msgPack := InitMessagePack()
	encoded, _ := msgPack.Encode(msg)
	decoder := NewMessagePackDecoder(encoded)

	var decoded string
	err := decoder.Decode(&decoded)

	assert.Nil(t, err)
	assert.Equal(t, msg, decoded)
}

func TestMsgPackDecoder_Decode_ErrDecoderNotInitialized(t *testing.T) {
	msgPack := InitMessagePack()
	encoded, _ := msgPack.Encode("example message")
	decoder := &MessagePackCodecDecoder{
		MessagePackCodec: MessagePackCodec{
			MsgPack: &codec.MsgpackHandle{},
			b:       &encoded,
		},
	}

	var decoded string
	err := decoder.Decode(&decoded)

	assert.NotNil(t, err)
	assert.Equal(t, "decoder not initialized", err.Error())
}

func TestEncodeString(t *testing.T) {
	s := "example string"

	encoded, err := EncodeString(s)

	assert.NotNil(t, encoded)
	assert.Nil(t, err)
}

func TestEncodeString_ErrStringLength(t *testing.T) {
	s := bytes.Repeat([]byte("a"), 0x80) // 128-byte long string

	encoded, err := EncodeString(string(s))

	assert.Nil(t, encoded)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid string length", err.Error())
}
