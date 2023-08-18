package structure

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/db/engine"
	"time"
)

// StreamMessage represents a message in a stream.
// It holds information about the unique ID of the message and the message payload.
type StreamMessage struct {
	// Id is the unique ID of the message.
	Id string `json:"id"`
	// Fields is the payload of the message.
	Fields map[string]interface{} `json:"fields"`
}

// StreamGroup represents a consumer group in a stream.
// It holds information about the group name, the last generated message ID,
// the last delivered message time, and the list of pending messages.
type StreamGroup struct {
	// Name is the name of the group.
	Name string `json:"name"`
	// LastGeneratedID is the ID of the last generated message in the group.
	LastGeneratedID string `json:"last_generated_id"`
	// LastDeliveredTime is the timestamp of the last delivered message in the group.
	LastDeliveredTime time.Time `json:"last_delivered_time"`
	// PendingMessages is the list of pending messages in the group.
	PendingMessages map[string]*StreamMessage `json:"pending_messages"`
}

// Streams represents a stream with messages and consumer groups.
// It holds information about the stream name, messages, consumer groups,
// and the last message in the stream.
type Streams struct {
	// Name is the name of the stream.
	Name string `json:"name"`
	// Messages is the list of messages in the stream.
	Messages []*StreamMessage `json:"messages"`
	// Groups is the list of consumer groups in the stream.
	Groups map[string]*StreamGroup `json:"groups"`
	// LastMessage is the timestamp of the last message in the stream.
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

// XAdd adds a new message to a stream.
// If the stream does not exist, it will be created.
// If the message ID already exists, it will return false.
//
// Parameters:
//
//	name: The name of the stream.
//	id: The ID of the message.
//	fields: The fields (attributes) of the message.
//	        It should be a map[string]interface{} where the keys represent the field names
//	        and the values represent the field values.
//
// Returns:
//
//	bool: Indicates whether the message was successfully added to the stream.
//	error: An error if any occurred during the operation, or nil on success.
//
// Note:
// - The ID and fields parameters should be valid and non-empty.
// - The name parameter represents the name of the stream where the message will be added.
// - The fields parameter contains the attributes of the message.
// - If the stream with the specified name does not exist, a new stream will be created.
// - If the message ID already exists in the stream, the function will return false and ErrExistID.
// - Otherwise, the message will be added to the stream, and the stream will be stored in the database.
func (s *StreamStructure) XAdd(name, id string, fields map[string]interface{}) (bool, error) {
	// Check if the arguments are valid
	if len(id) == 0 || len(fields) == 0 || fields == nil {
		return false, ErrInvalidXArgs
	}

	// Get the stream
	value, _ := s.db.Get([]byte(name))
	if value == nil {
		// Create a new stream
		s.streams = &Streams{
			Name:        name,
			Messages:    []*StreamMessage{},
			Groups:      map[string]*StreamGroup{},
			LastMessage: time.Now(),
		}

		// Create a new message
		message := &StreamMessage{
			Id:     id,
			Fields: fields,
		}
		// Add the message to the stream
		s.streams.Messages = append(s.streams.Messages, message)

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
	} else {
		// Decode the streams
		err := s.decodeStreams(value, s.streams)
		if err != nil {
			return false, err
		}

		// Check if the message ID already exists
		for _, msg := range s.streams.Messages {
			if msg.Id == id {
				return false, ErrExistID
			}
		}

		// Create a new message
		message := &StreamMessage{
			Id:     id,
			Fields: fields,
		}

		// Add the message to the stream
		s.streams.Messages = append(s.streams.Messages, message)

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
}

// XRead reads messages from a stream.
// Returns a slice of StreamMessage.
//
// Parameters:
//
//	name: The name of the stream.
//	count: The maximum number of messages to read.
//
// Returns:
//
//	[]StreamMessage: A slice of StreamMessage containing the read messages.
//	error: An error if any occurred during the operation, or nil on success.
//
// Note:
//   - The name parameter represents the name of the stream to read from.
//   - The count parameter specifies the maximum number of messages to read.
//   - If count is less than or equal to 0, it will return ErrInvalidCount.
//   - It retrieves the stream from the database using the specified name.
//   - If the stream does not exist, it will return ErrKeyNotFound.
//   - It decodes the stream data and stores it in the internal s.streams field.
//   - It retrieves the messages from s.streams.Messages.
//   - If the number of messages in the stream is greater than or equal to count,
//     it returns a slice of StreamMessage with the first count messages.
//   - If the number of messages in the stream is less than count, it returns
//     ErrAmountOfData, indicating that there is not enough data in the stream.
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

// XDel deletes a message from a stream.
// Returns true if the message was deleted.
// Returns false if the message was not deleted.
// Returns the number of messages in the stream after deletion.
//
// Parameters:
//
//	name: The name of the stream.
//	ids: The ID of the message to delete.
//
// Returns:
//
//	bool: Indicates whether the message was successfully deleted.
//	int: The number of messages in the stream after deletion.
//	error: An error if any occurred during the operation, or nil on success.
//
// Note:
// - The name parameter represents the name of the stream to delete from.
// - The ids parameter specifies the ID of the message to delete.
// - It retrieves the stream from the database using the specified name.
// - If the stream does not exist, it will return ErrKeyNotFound.
// - It decodes the stream data and stores it in the internal s.streams field.
// - It iterates through the messages in the stream and removes the message with the specified ID.
// - It updates the s.streams.Messages slice to remove the deleted message.
// - It encodes the updated stream data.
// - It stores the updated stream in the database.
// - It returns true if the message was deleted successfully.
// - It returns false if the message was not found in the stream.
// - It returns the number of messages in the stream after deletion.
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
	found := false
	for _, msg := range messages {
		if msg.Id != ids {
			result = append(result, msg)
		} else {
			found = true
		}
	}

	if !found {
		return false, len(s.streams.Messages), nil
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

// XLen returns the number of elements in a given stream.
// It takes the name of the stream as an argument and returns the number of elements in the stream.
//
// Parameters:
//
//	name: The name of the stream.
//
// Returns:
//
//	int: The number of elements (messages) in the stream.
//	error: An error if any occurred during the operation, or nil on success.
//
// Note:
// - The name parameter represents the name of the stream to get the length of.
// - It retrieves the stream from the database using the specified name.
// - If the stream does not exist, it will return ErrKeyNotFound.
// - It decodes the stream data and stores it in the internal s.streams field.
// - It retrieves the messages from s.streams.Messages.
// - It returns the number of messages in the stream.
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

// XRange returns the messages in the stream.
// It takes the name of the stream, the start index, and the stop index as arguments,
// and returns a slice of StreamMessage containing the messages within the specified range.
//
// Parameters:
//
//	name: The name of the stream.
//	start: The start index of the range (inclusive).
//	stop: The stop index of the range (exclusive).
//
// Returns:
//
//	[]StreamMessage: A slice of StreamMessage containing the messages within the specified range.
//	error: An error if any occurred during the operation, or nil on success.
//
// Note:
//   - The name parameter represents the name of the stream to get the messages from.
//   - The start parameter specifies the start index of the range (inclusive).
//   - The stop parameter specifies the stop index of the range (exclusive).
//   - It retrieves the stream from the database using the specified name.
//   - If the stream does not exist, it will return ErrKeyNotFound.
//   - It decodes the stream data and stores it in the internal s.streams field.
//   - It retrieves the messages from s.streams.Messages.
//   - If the number of messages in the stream is greater than or equal to stop,
//     it returns a slice of StreamMessage within the specified range.
//   - If the number of messages in the stream is less than stop, it returns
//     ErrAmountOfData, indicating that there is not enough data in the stream.
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

// XRevRange returns the messages in the stream in reverse order.
// It takes the name of the stream, the start index, and the stop index as arguments,
// and returns a slice of StreamMessage containing the messages within the
// specified range in reverse order.
//
// Parameters:
//
//	name: The name of the stream.
//	start: The start index of the range (inclusive).
//	stop: The stop index of the range (exclusive).
//
// Returns:
//
//	 []StreamMessage: A slice of StreamMessage containing the messages
//					  within the specified range in reverse order.
//	 error: An error if any occurred during the operation, or nil on success.
//
// Note:
//   - The name parameter represents the name of the stream to get the messages from.
//   - The start parameter specifies the start index of the range (inclusive).
//   - The stop parameter specifies the stop index of the range (exclusive).
//   - It retrieves the stream from the database using the specified name.
//   - If the stream does not exist, it will return ErrKeyNotFound.
//   - It decodes the stream data and stores it in the internal s.streams field.
//   - It retrieves the messages from s.streams.Messages.
//   - If the number of messages in the stream is greater than or equal to stop,
//     it returns a slice of StreamMessage within the specified range in reverse order.
//   - If the number of messages in the stream is less than stop, it returns
//     ErrAmountOfData, indicating that there is not enough data in the stream.
//   - The returned slice of StreamMessage is reversed compared to the original order.
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
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// XTrim trims the stream to a certain size.
// It takes the name of the stream and the maximum length as arguments,
// and returns the number of messages in the stream after the trim operation.
//
// Parameters:
//
//	name: The name of the stream.
//	maxLen: The maximum length to trim the stream to.
//
// Returns:
//
//	int: The number of messages in the stream after the trim operation.
//	error: An error if any occurred during the operation, or nil on success.
//
// Note:
//   - The name parameter represents the name of the stream to trim.
//   - The maxLen parameter specifies the maximum length to trim the stream to.
//   - It retrieves the stream from the database using the specified name.
//   - If the stream does not exist, it will return ErrKeyNotFound.
//   - It decodes the stream data and stores it in the internal s.streams field.
//   - It retrieves the messages from s.streams.Messages.
//   - If the number of messages in the stream is greater than or equal to maxLen,
//     it trims the stream to the specified length.
//   - If the number of messages in the stream is less than maxLen, it returns
//     ErrAmountOfData, indicating that there is not enough data in the stream.
//   - It sets the trimmed messages back to s.streams.Messages.
//   - It encodes the streams and updates the stream in the database.
//   - It returns the number of messages in the stream after the trim operation.
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
		result = append(result, messages...)
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

// XGroup creates a new consumer group.
// It takes the name of the stream, the name of the group, and the ID of the message as arguments,
// and returns a boolean indicating whether the group was created successfully or not, and an error if any.
//
// Parameters:
//
//	name: The name of the stream.
//	group: The name of the consumer group.
//	id: The ID of the message.
//
// Returns:
//
//	bool: A boolean indicating whether the group was created successfully or not.
//	error: An error if any occurred during the operation, or nil on success.
//
// Note:
// - The name parameter represents the name of the stream to create the consumer group on.
// - The group parameter specifies the name of the consumer group to create.
// - The id parameter specifies the ID of the message.
// - It retrieves the stream from the database using the specified name.
// - If the stream does not exist, it will return ErrKeyNotFound.
// - It decodes the stream data and stores it in the internal s.streams field.
// - It creates a new StreamGroup with the specified name, last generated ID, and last delivered time.
// - It sets the StreamGroup in s.streams.Groups map with the group name as the key.
// - It encodes the streams and updates the stream in the database.
// - It returns true if the group was created successfully, and false otherwise.
// - If any error occurs during the operation, it will be returned along with false.
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

// encodeStreams encodes the streams.
// It takes the Streams object as an argument and returns the encoded data as a byte slice.
//
// Parameters:
//
//	ss: The Streams object to encode.
//
// Returns:
//
//	[]byte: The encoded data as a byte slice.
//	error: An error if any occurred during the encoding process, or nil on success.
//
// Note:
// - The ss parameter represents the Streams object to be encoded.
// - It uses the json.Marshal function to encode the Streams object into JSON format.
// - If any error occurs during the encoding process, it will be returned along with nil byte slice.
// - On success, it returns the encoded data as a byte slice.
func (s *StreamStructure) encodeStreams(ss *Streams) ([]byte, error) {
	// Encode the streams
	data, err := json.Marshal(ss)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// decodeStreams decodes the streams.
// It takes the encoded data as a byte slice and the target Streams object as arguments.
//
// Parameters:
//
//	ss: The encoded data as a byte slice.
//	ss2: The target Streams object to store the decoded data.
//
// Returns:
//
//	error: An error if any occurred during the decoding process, or nil on success.
//
// Note:
// - The ss parameter represents the encoded data to be decoded.
// - The ss2 parameter represents the target Streams object to store the decoded data.
// - It uses the json.Unmarshal function to decode the encoded data into the Streams object.
// - If any error occurs during the decoding process, it will be returned.
// - On success, it returns nil.
func (s *StreamStructure) decodeStreams(ss []byte, ss2 *Streams) error {
	// Decode the streams
	if err := json.Unmarshal(ss, ss2); err != nil {
		fmt.Println("err", err)
		return err
	}
	return nil
}

// encodeStreamGroup encodes the stream group.
// It takes the StreamGroup object as an argument and returns the encoded data as a byte slice.
//
// Parameters:
//
//	sg: The StreamGroup object to encode.
//
// Returns:
//
//	[]byte: The encoded data as a byte slice.
//	error: An error if any occurred during the encoding process, or nil on success.
//
// Note:
// - The sg parameter represents the StreamGroup object to be encoded.
// - It uses the json.Marshal function to encode the StreamGroup object into JSON format.
// - If any error occurs during the encoding process, it will be returned along with nil byte slice.
// - On success, it returns the encoded data as a byte slice.
func (s *StreamStructure) encodeStreamGroup(sg *StreamGroup) ([]byte, error) {
	// Encode the stream group
	data, err := json.Marshal(sg)
	if err != nil {
		return nil, err
	}
	return data, nil
}
