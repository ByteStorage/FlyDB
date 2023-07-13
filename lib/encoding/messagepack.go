package encoding

import (
	"bytes"
	"github.com/hashicorp/go-msgpack/codec"
)

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
