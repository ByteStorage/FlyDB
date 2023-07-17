package encoding

import (
	"bytes"
	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"reflect"
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
		MessagePackCodec: &MessagePackCodec{
			MsgPack: NewMsgPackHandle(),
			b:       &encoded,
		},
	}

	var decoded string
	err := decoder.Decode(&decoded)

	assert.NotNil(t, err)
	assert.Equal(t, "decoder not initialized", err.Error())
}

func TestEncodeString_ErrStringLength(t *testing.T) {
	s := bytes.Repeat([]byte("a"), 0x80) // 128-byte long string

	encoded, err := EncodeString(string(s))

	assert.Nil(t, encoded)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid string length", err.Error())
}

func TestMessagePackCodec_Encode(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
	}
	assert := assert.New(t)

	codec := &MessagePackCodec{
		MsgPack: NewMsgPackHandle(),
	}

	t.Run("successful encoding", func(t *testing.T) {
		testStruct := &TestStruct{
			Field1: "Test",
			Field2: 1,
		}
		outStruct := &TestStruct{}

		b, err := codec.Encode(testStruct)
		assert.Nil(err, "Error should be nil")
		err = codec.Decode(b, outStruct)
		assert.NoError(err)
		assert.EqualValues(testStruct, outStruct)
	})

}
func TestEncodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
		hasErr   bool
	}{
		{"hello", []byte{0x05, 'h', 'e', 'l', 'l', 'o'}, false},
		{"world", []byte{0x05, 'w', 'o', 'r', 'l', 'd'}, false},
		{string(make([]byte, 0x80)), nil, true},
	}

	for _, tt := range tests {
		result, err := EncodeString(tt.input)
		if (err != nil) != tt.hasErr {
			t.Errorf("EncodeString(%q) error = %v, wantErr %v", tt.input, err, tt.hasErr)
			continue
		}
		if !tt.hasErr && string(result) != string(tt.expected) {
			t.Errorf("EncodeString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input     []byte
		expectedL int
		expected  string
		hasErr    bool
	}{
		{[]byte{0x05, 'h', 'e', 'l', 'l', 'o'}, 6, "hello", false},
		{[]byte{0x05, 'w', 'o', 'r', 'l', 'd'}, 6, "world", false},
		{[]byte{0x05, 'w', 'o', 'r'}, 0, "", true},
		{[]byte{}, 0, "", true},
	}

	for _, tt := range tests {
		length, result, err := DecodeString(tt.input)
		if (err != nil) != tt.hasErr {
			t.Errorf("DecodeString(%v) error = %v, wantErr %v", tt.input, err, tt.hasErr)
			continue
		}
		if !tt.hasErr && (result != tt.expected || length != tt.expectedL) {
			t.Errorf("DecodeString(%v) = %v,%q, want %v,%q", tt.input, length, result, tt.expectedL, tt.expected)
		}
	}
}
func TestMessagePackCodecEncoder_Encode(t *testing.T) {
	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[int]int{})

	encoder := NewMessagePackEncoder()

	err := encoder.Encode(map[int]int{1: 2})
	assert.NoError(t, err)

	err = encoder.Encode(map[int]int{3: 4})
	assert.NoError(t, err)
}

func TestAddExtension(t *testing.T) {

	type CustomType struct {
		Name string
	}

	// global extention info.
	const extensionID byte = 1

	encoder := func(rv reflect.Value) ([]byte, error) {
		ct := rv.Interface().(CustomType)
		return []byte(ct.Name), nil
	}
	decoder := func(rv reflect.Value, b []byte) error {
		rv.Set(reflect.ValueOf(CustomType{Name: string(b)}))
		return nil
	}

	m := NewMessagePackEncoder()
	err := m.AddExtension(reflect.TypeOf(CustomType{}), extensionID, encoder, decoder)
	if err != nil {
		t.Fatalf("Failed adding extension: %v", err)
	}
	data := CustomType{Name: "test"}
	dataVerify := CustomType{}
	err = m.enc.Encode(&data)
	assert.NoError(t, err)
	assert.NotNil(t, m.Bytes())
	d := NewMessagePackDecoder(m.Bytes())
	err = d.AddExtension(reflect.TypeOf(CustomType{}), extensionID, encoder, decoder)
	assert.NoError(t, err)
	err = d.Decode(&dataVerify)
	assert.NoError(t, err)
	assert.EqualValues(t, dataVerify, data)
}
