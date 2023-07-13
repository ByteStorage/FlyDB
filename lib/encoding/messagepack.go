package encoding

import (
	"bytes"
	"errors"
	"github.com/hashicorp/go-msgpack/codec"
	"reflect"
)

type MessagePackCodec struct {
	msgPack *codec.MsgpackHandle
}

func InitMessagePack() MessagePackCodec {
	return MessagePackCodec{
		msgPack: &codec.MsgpackHandle{},
	}
}
func (m MessagePackCodec) Encode(msg interface{}) ([]byte, error) {
	m.msgPack.RawToString = true
	m.msgPack.WriteExt = true
	m.msgPack.MapType = reflect.TypeOf(map[string]interface{}(nil))

	var b []byte
	enc := codec.NewEncoderBytes(&b, m.msgPack)
	err := enc.Encode(msg)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (m MessagePackCodec) Decode(in []byte, out interface{}) error {
	dev := codec.NewDecoderBytes(in, m.msgPack) // Create a new decoder with the buffer and MessagePack handle

	return dev.Decode(out) // Decode the byte slice into the provided output structure
}
func (m MessagePackCodec) AddExtension(
	t reflect.Type,
	id byte,
	encoder func(reflect.Value) ([]byte, error),
	decoder func(reflect.Value, []byte) error) error {
	return m.msgPack.AddExt(t, id, encoder, decoder)
}

// EncodeMessagePack is a function that encodes a given message using MessagePack serialization.
// It takes an interface{} parameter representing the message and returns the encoded byte slice and an error.
func EncodeMessagePack(msg interface{}) ([]byte, error) {
	var b []byte
	var mph codec.MsgpackHandle
	h := &mph
	enc := codec.NewEncoderBytes(&b, h) // Create a new encoder with the provided message and MessagePack handle

	err := enc.Encode(msg) // Encode the message using the encoder
	if err != nil {
		return nil, err
	}
	return b, nil // Return the encoded byte slice
}

// DecodeMessagePack is a function that decodes a given byte slice using MessagePack deserialization.
// It takes an input byte slice and an interface{} representing the output structure for the deserialized message.
// It returns an error if the decoding process fails.
func DecodeMessagePack(in []byte, out interface{}) error {
	buf := bytes.NewBuffer(in) // Create a new buffer with the input byte slice
	mph := codec.MsgpackHandle{}
	dev := codec.NewDecoder(buf, &mph) // Create a new decoder with the buffer and MessagePack handle

	return dev.Decode(out) // Decode the byte slice into the provided output structure
}

func EncodeString(s string) ([]byte, error) {
	if len(s) > 0x7F {
		return nil, errors.New("invalid string length")
	}
	b := make([]byte, len(s)+1)
	b[0] = byte(len(s))
	copy(b[1:], s)
	return b, nil
}

func DecodeString(b []byte) (int, string, error) {
	if len(b) == 0 {
		return 0, "", errors.New("invalid length")
	}
	l := int(b[0])
	if len(b) < (l + 1) {
		return 0, "", errors.New("invalid length")
	}
	s := make([]byte, l)
	copy(s, b[1:l+1])
	return l + 1, string(s), nil
}
