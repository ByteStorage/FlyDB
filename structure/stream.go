package structure

import (
	"encoding/json"
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"sort"
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
	LastDeliveredTime time.Time `json:"last_delivered_time"`
	// PendingMessages is the list of pending messages
	PendingMessages map[string]*StreamMessage `json:"pending_messages"`
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
	// ErrInvalidXArgs is returned when the arguments are invalid
	ErrInvalidXArgs = errors.New("id or fields cannot be empty")
	// ErrExistID is returned when the message ID already exists
	ErrExistID = errors.New("message ID already exists")
	// ErrInvalidCount is returned when the count is invalid
	ErrInvalidCount = errors.New("invalid count")
	// ErrInvalidStream is returned when the stream is invalid
	ErrAmountOfData = errors.New("The number of queries exceeds the amount of data in the stream")
)

// XAdd adds a new message to a stream
// If the stream does not exist, it will be created
// If the message ID already exists, it will return false
func (s *StreamStructure) XAdd(name, id string, fields map[string]interface{}) (bool, error) {
	// Check if the arguments are valid
	if len(id) == 0 || len(fields) == 0 || fields == nil {
		return false, ErrInvalidXArgs
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

// XRead reads messages from a stream
// Returns a slice of StreamMessage
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

// XDel deletes a message from a stream
// Returns true if the message was deleted
// Returns false if the message was not deleted
// Returns the number of messages in the stream
func (s *StreamStructure) XDel(name string, ids string) (bool, int, error) {
	// Get the stream
	encodedStreams, err := s.db.Get([]byte(name))
	if err != nil {
		return false, len(s.streams.Messages), err
	}

	// Decode the streams
	if err = s.decodeStreams(encodedStreams, s.streams); err != nil {
		return false, len(s.streams.Messages), err
	}

	// Get the messages
	messages := s.streams.Messages

	// Create a new slice of StreamMessage
	var result []*StreamMessage

	// Get the messages
	for _, msg := range messages {
		if msg.Id != ids {
			result = append(result, msg)
		}
	}

	// Set the messages
	s.streams.Messages = result

	// Encode the streams
	encodedStreams, err = s.encodeStreams(s.streams)
	if err != nil {
		return false, len(s.streams.Messages), err
	}

	// Set the stream
	if err = s.db.Put([]byte(s.streams.Name), encodedStreams); err != nil {
		return false, len(s.streams.Messages), err
	}

	return true, len(s.streams.Messages), nil
}

// XLen returns the number of elements in a given stream
// with the name of the stream as an argument
// and the number of elements in the stream as the return value
func (s *StreamStructure) XLen(name string) (int, error) {
	// Get the stream
	encodedStreams, err := s.db.Get([]byte(name))
	if err != nil {
		return 0, err
	}

	// Decode the streams
	if err = s.decodeStreams(encodedStreams, s.streams); err != nil {
		return 0, err
	}

	// Get the messages
	messages := s.streams.Messages

	return len(messages), nil
}

// XRange returns the messages in the stream
// with the []StreamMessage as the return value
func (s *StreamStructure) XRange(name string, start, stop int) ([]StreamMessage, error) {
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
	if len(messages) >= stop {
		messages = messages[start:stop]
		// Convert []*StreamMessage to []StreamMessage
		for _, msg := range messages {
			result = append(result, *msg)
		}
	} else {
		return nil, ErrAmountOfData
	}

	return result, nil
}

// XRevRange returns the messages in the stream
// with the []StreamMessage as the return value
func (s *StreamStructure) XRevRange(name string, start, stop int) ([]StreamMessage, error) {
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
	if len(messages) >= stop {
		messages = messages[start:stop]
		// Convert []*StreamMessage to []StreamMessage
		for _, msg := range messages {
			result = append(result, *msg)
		}
	} else {
		return nil, ErrAmountOfData
	}

	// Reverse the slice
	sort.Slice(result, func(i, j int) bool {
		return i > j
	})

	return result, nil
}

// XTrim trims the stream to a certain size
// with the number of messages in the stream as the return value
func (s *StreamStructure) XTrim(name string, maxLen int) (int, error) {
	// Get the stream
	encodedStreams, err := s.db.Get([]byte(name))
	if err != nil {
		return 0, err
	}

	// Decode the streams
	if err = s.decodeStreams(encodedStreams, s.streams); err != nil {
		return 0, err
	}

	// Get the messages
	messages := s.streams.Messages

	// Create a new slice of StreamMessage
	var result []*StreamMessage

	// Get the messages
	if len(messages) >= maxLen {
		messages = messages[:maxLen]
		// Convert []*StreamMessage to []StreamMessage
		for _, msg := range messages {
			result = append(result, msg)
		}
	} else {
		return 0, ErrAmountOfData
	}

	// Set the messages
	s.streams.Messages = result

	// Encode the streams
	encodedStreams, err = s.encodeStreams(s.streams)
	if err != nil {
		return 0, err
	}

	// Set the stream
	if err = s.db.Put([]byte(s.streams.Name), encodedStreams); err != nil {
		return 0, err
	}

	return len(s.streams.Messages), nil
}

// XGroup creates a new consumer group
// with the name of the stream, the name of the group, and the id of the message as arguments
// and a boolean as the return value
func (s *StreamStructure) XGroup(name, group, id string) (bool, error) {
	// Get the stream
	encodedStreams, err := s.db.Get([]byte(name))
	if err != nil {
		return false, err
	}

	// Decode the streams
	if err = s.decodeStreams(encodedStreams, s.streams); err != nil {
		return false, err
	}

	// Create a new stream group
	sg := &StreamGroup{
		Name:              group,
		LastGeneratedID:   id,
		LastDeliveredTime: time.Now(),
		PendingMessages:   make(map[string]*StreamMessage),
	}

	// Set the stream group
	s.streams.Groups[group] = sg

	// Encode the streams
	encodedStreams, err = s.encodeStreams(s.streams)
	if err != nil {
		return false, err
	}

	// Set the stream
	if err = s.db.Put([]byte(s.streams.Name), encodedStreams); err != nil {
		return false, err
	}

	return true, nil
}

// encodeStreams encodes the streams
// with the []byte as the return value
func (s *StreamStructure) encodeStreams(ss *Streams) ([]byte, error) {
	// Encode the streams
	data, err := json.Marshal(ss)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// decodeStreams decodes the streams
func (s *StreamStructure) decodeStreams(ss []byte, ss2 *Streams) error {
	// Decode the streams
	if err := json.Unmarshal(ss, ss2); err != nil {
		return err
	}
	return nil
}

// encodeStreamGroup encodes the stream group
// with the []byte as the return value
func (s *StreamStructure) encodeStreamGroup(sg *StreamGroup) ([]byte, error) {
	// Encode the stream group
	data, err := json.Marshal(sg)
	if err != nil {
		return nil, err
	}
	return data, nil
}
