package structure

type DataStructure = byte

const (
	// String is a string data structure
	String DataStructure = iota + 1
	// Hash is a hash data structure
	Hash
	// List is a list data structure
	List DataStructure = iota + 1
	// Set is a set data structure
	Set
	// ZSet is a zset data structure
	ZSet
	// bitmap is a bitmap data structure
	Bitmap
	// Stream is a stream data structure
	Stream
	// Expire is a expire data structure
	Expire
)
