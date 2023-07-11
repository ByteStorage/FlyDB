package structure

import (
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
