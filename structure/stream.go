package structure

import (
	"encoding/json"
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"time"
)

// StreamMessage represents a message in a stream
type StreamMessage struct {
	// ID is the unique ID of the message
	Id string `json:"id"`
	// Fields is the message payload
	Fields map[string]interface{} `json:"fields"`
}

// StreamGroup represents a consumer group in a stream
type StreamGroup struct {
	// Name is the name of the group
	Name string `json:"name"`
	// LastDeliveredID is the last delivered message ID
	LastGeneratedID string `json:"last_generated_id"`
	// LastDeliveredTime is the last delivered message time
	LastDeliveredTime time.Time
	// PendingMessages is the list of pending messages
	PendingMessages map[string]*StreamMessage
}

// Streams represents a stream with messages and consumer groups
type Streams struct {
	// Name is the name of the stream
	Name string `json:"name"`
	// Messages is the list of messages in the stream
	Messages []*StreamMessage `json:"messages"`
	// Groups is the list of consumer groups in the stream
	Groups map[string]*StreamGroup `json:"groups"`
	// LastMessage is the last message in the stream
	LastMessage time.Time
}

type StreamStructure struct {
	db      *engine.DB
	streams *Streams
}

// StreamStructure represents a stream structure
func NewStreamStructure(options config.Options) (*StreamStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &StreamStructure{db: db}, nil
}

var (
	// ErrInvalidArgs is returned when the arguments are invalid
	ErrInvalidArgs = errors.New("id or fields cannot be empty")
	// ErrExistID is returned when the message ID already exists
	ErrExistID = errors.New("message ID already exists")
	// ErrInvalidCount is returned when the count is invalid
	ErrInvalidCount = errors.New("invalid count")
	// ErrInvalidStream is returned when the stream is invalid
	ErrAmountOfData = errors.New("The number of queries exceeds the amount of data in the stream")
)

func (s *StreamStructure) XAdd(name, id string, fields map[string]interface{}) (bool, error) {
	// Check if the arguments are valid
	if len(id) == 0 || len(fields) == 0 || fields == nil {
		return false, ErrInvalidArgs
	}

	// init stream if not exist
	if s.streams == nil {
		// Create a new stream
		s.streams = &Streams{
			Name:        name,
			Messages:    []*StreamMessage{},
			Groups:      make(map[string]*StreamGroup),
			LastMessage: time.Time{},
		}
	} else {
		// Check if the stream name is the same
		if s.streams.Name != name {
			s.streams = &Streams{
				Name:        name,
				Messages:    []*StreamMessage{},
				Groups:      make(map[string]*StreamGroup),
				LastMessage: time.Time{},
			}
		}
	}

	// Check if the message ID already exists
	_, err := s.db.Get([]byte(id))
	if err == nil {
		return false, ErrExistID
	}

	// Create a new message
	message := &StreamMessage{
		Id:     id,
		Fields: fields,
	}

	// Append the message to the stream
	s.streams.Messages = append(s.streams.Messages, message)

	// Set the last message time
	s.streams.LastMessage = time.Now()

	// Encode the streams
	encodedStreams, err := s.encodeStreams(s.streams)
	if err != nil {
		return false, err
	}

	// Set the stream
	if err = s.db.Put([]byte(s.streams.Name), encodedStreams); err != nil {
		return false, err
	}

	return true, nil
}

func (s *StreamStructure) XRead(name string, count int) ([]StreamMessage, error) {
	if count <= 0 {
		return nil, ErrInvalidCount
	}

	// Get the stream
	encodedStreams, err := s.db.Get([]byte(name))
	if err != nil {
		return nil, err
	}

	// Decode the streams
	if err = s.decodeStreams(encodedStreams, s.streams); err != nil {
		return nil, err
	}

	// Get the messages
	messages := s.streams.Messages

	// Create a new slice of StreamMessage
	var result []StreamMessage

	// Get the messages
	if len(messages) >= count {
		messages = messages[:count]
		// Convert []*StreamMessage to []StreamMessage
		for _, msg := range messages {
			result = append(result, *msg)
		}
	} else {
		return nil, ErrAmountOfData
	}

	return result, nil
}

func (s *StreamStructure) encodeStreams(ss *Streams) ([]byte, error) {
	// Encode the streams
	data, err := json.Marshal(ss)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *StreamStructure) decodeStreams(ss []byte, ss2 *Streams) error {
	// Decode the streams
	if err := json.Unmarshal(ss, ss2); err != nil {
		return err
	}
	return nil
}
