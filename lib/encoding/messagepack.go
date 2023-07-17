package encoding

import (
	"bytes"
	"errors"
	"github.com/hashicorp/go-msgpack/codec"
	"reflect"
)

// MessagePackCodec struct, holds references to MessagePack handler and byte slice,
// along with Encoder and Decoder, and a typeMap for storing reflect.Type
type MessagePackCodec struct {
	MsgPack *codec.MsgpackHandle
	b       *[]byte
	enc     *codec.Encoder
	dec     *codec.Decoder
}

// MessagePackCodecEncoder struct derives from MessagePackCodec
// it manages IDs and counts of the encoded objects.
type MessagePackCodecEncoder struct {
	MessagePackCodec // Embedded MessagePackCodec

	// nextId is used probably for tracking ID of the next object to encode.
	nextId uint

	// objects represents the count of objects that have been encoded.
	objects int
}

// MessagePackCodecDecoder struct, holds a reference to a MessagePackCodec instance.
type MessagePackCodecDecoder struct {
	MessagePackCodec
}

// InitMessagePack function initializes MessagePackCodec struct and returns it.
func InitMessagePack() MessagePackCodec {
	return MessagePackCodec{
		MsgPack: &codec.MsgpackHandle{},
	}
}

// NewMessagePackEncoder function creates new MessagePackCodecEncoder and initializes it.
func NewMessagePackEncoder() *MessagePackCodecEncoder {
	msgPack := &codec.MsgpackHandle{}
	b := make([]byte, 0)
	return &MessagePackCodecEncoder{
		MessagePackCodec: MessagePackCodec{
			MsgPack: &codec.MsgpackHandle{},
			b:       &b,
			enc:     codec.NewEncoderBytes(&b, msgPack),
		},
	}
}

// NewMessagePackDecoder function takes in a byte slice, and returns a pointer to newly created
// and initialized MessagePackCodecDecoder
func NewMessagePackDecoder(b []byte) *MessagePackCodecDecoder {
	msgPack := &codec.MsgpackHandle{}
	return &MessagePackCodecDecoder{
		MessagePackCodec: MessagePackCodec{
			MsgPack: &codec.MsgpackHandle{},
			b:       &b,
			dec:     codec.NewDecoderBytes(b, msgPack),
		},
	}
}

// Encode method for MessagePackCodec. It encodes the input value into a byte slice using MessagePack.
// Returns encoded byte slice or error.
func (m *MessagePackCodec) Encode(msg interface{}) ([]byte, error) {
	var b []byte
	err := codec.NewEncoderBytes(&b, m.MsgPack).Encode(msg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Encode is a method for MessagePackCodecEncoder.
// It takes in msg of type interface{} as input, that is to be encoded.
// Returns an error if encountered during encoding.
func (m *MessagePackCodecEncoder) Encode(msg interface{}) error {
	return m.enc.Encode(msg)
}

// Bytes is a method for MessagePackCodecEncoder.
// It returns a byte slice pointer b.
func (m *MessagePackCodecEncoder) Bytes() []byte {
	return *m.b
}

// Decode is a method on MessagePackCodecDecoder that decodes MessagePack data
// into the provided interface; returns an error if any decoding issues occur.
func (m *MessagePackCodecDecoder) Decode(msg interface{}) error {
	if m.dec == nil {
		return errors.New("decoder not initialized")
	}
	return m.dec.Decode(msg)
}

// Decode on MessagePackCodec type, using a byte slice as input.
func (m *MessagePackCodec) Decode(in []byte, out interface{}) error {
	// Create new decoder using the byte slice and MessagePack handle.
	dec := codec.NewDecoderBytes(in, m.MsgPack)

	// Attempt to decode the byte slice into the desired output structure.
	return dec.Decode(out)
}

// AddExtension method allows for setting custom encoders/decoders for specific reflect.Types.
func (m *MessagePackCodec) AddExtension(
	t reflect.Type,
	id byte,
	encoder func(reflect.Value) ([]byte, error),
	decoder func(reflect.Value, []byte) error) error {

	return m.MsgPack.AddExt(t, id, encoder, decoder)
}

// EncodeMessagePack function encodes a given object into MessagePack format.
func EncodeMessagePack(msg interface{}) ([]byte, error) {
	// Directly initialize the byte slice and encoder.
	b := make([]byte, 0)
	enc := codec.NewEncoderBytes(&b, &codec.MsgpackHandle{})

	// Attempt to encode the message.
	if err := enc.Encode(msg); err != nil {
		return nil, err
	}

	// Return the encoded byte slice.
	return b, nil
}

// DecodeMessagePack function decodes a byte slice of MessagePack data into a given object.
func DecodeMessagePack(in []byte, out interface{}) error {
	dec := codec.NewDecoder(bytes.NewBuffer(in), &codec.MsgpackHandle{})
	return dec.Decode(out)
}

// EncodeString Functions for encoding and decoding strings to and from byte slices.
func EncodeString(s string) ([]byte, error) {
	// Check if string length is within correct bounds.
	if len(s) > 0x7F {
		return nil, errors.New("invalid string length")
	}

	// Create a byte slice of appropriate length.
	b := make([]byte, len(s)+1)
	b[0] = byte(len(s))

	// Copy the string into the byte slice.
	copy(b[1:], s)

	// Return the byte slice.
	return b, nil
}

// DecodeString is a function that takes an input byte slice and attempts to decode it to obtain a string.
// Return parameters are an integer, a string and an error. Integer denotes the length of the byte slice
// representation of the string including length-field. The second return parameter is the decoded string.
// DecodeString raises an error if the length of byte slice is less than the expected string length plus
// one (considering the string length field) or if the provided byte slice is empty.
// If successful, returns length of byte representation of string, the decoded string and a nil error.
func DecodeString(b []byte) (int, string, error) {
	// Check that byte slice is not empty.
	if len(b) == 0 {
		return 0, "", errors.New("invalid length")
	}

	// Determine the length of the string.
	l := int(b[0])
	if len(b) < (l + 1) {
		return 0, "", errors.New("invalid length")
	}

	// Create a byte slice of the appropriate length and copy the string into it.
	s := make([]byte, l)
	copy(s, b[1:l+1])

	// Return the length of the string and the string itself.
	return l + 1, string(s), nil
}
