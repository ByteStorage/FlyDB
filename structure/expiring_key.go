package structure

import (
	"encoding/binary"
	"github.com/ByteStorage/FlyDB/db/engine"
	"time"

	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
)

type ExpireStructure struct {
	db *engine.DB
}

// NewExpireStructure returns a new ExpireStructure
// It will return a nil ExpireStructure if the database cannot be opened
// or the database cannot be created
// The database will be created if it does not exist
func NewExpireStructure(options config.Options) (*ExpireStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &ExpireStructure{db: db}, nil
}

func (ek *ExpireStructure) EXPIRE(key string, seconds int64) error {
	deadtime := ek.getCurrentMiliUnixTimeStamp() + seconds*1000
	return ek.setExpireToDB(key, deadtime)
}

func (ek *ExpireStructure) PEXPIRE(key string, milliseconds int64) error {
	deadtime := ek.getCurrentMiliUnixTimeStamp() + milliseconds
	return ek.setExpireToDB(key, deadtime)
}

func (ek *ExpireStructure) EXPIREAT(key string, timestamp int64) error {
	return ek.setExpireToDB(key, timestamp*1000)
}

func (ek *ExpireStructure) PEXPIREAT(key string, milliseconds_timestamp int64) error {
	return ek.setExpireToDB(key, milliseconds_timestamp)
}

func (ek *ExpireStructure) TTL(key string) (int64, error) {
	// Get the deadtime
	deadtime, err := ek.getExpireFromDB(key, false)
	if err != nil {
		return 0, err
	}
	//key not exists
	if deadtime == -2 {
		return -2, nil
	}

	//forever
	if deadtime == -1 {
		return -1, nil
	}

	result := (deadtime - ek.getCurrentMiliUnixTimeStamp()) / 1000
	if result <= 0 {
		return -2, nil
	}
	return result, nil
}

func (ek *ExpireStructure) PTTL(key string) (int64, error) {
	// Get the deadtime
	deadtime, err := ek.getExpireFromDB(key, false)
	if err != nil {
		return 0, err
	}
	//key not exists
	if deadtime == -2 {
		return -2, nil
	}
	//forever
	if deadtime == -1 {
		return -1, nil
	}
	result := deadtime - ek.getCurrentMiliUnixTimeStamp()
	if result <= 0 {
		return -2, nil
	}
	return result, nil
}

func (ek *ExpireStructure) PERSIST(key string, seconds int64) error {
	return ek.setExpireToDB(key, -1)
}

func (ek *ExpireStructure) EXPIREBY(key string, duration int64, tag ...string) error {
	if len(tag) == 0 {
		return ek.setExpireToDB(key, duration*1000)
	} else if len(tag) == 1 && tag[0] == "ex" {
		return ek.setExpireToDB(key, ek.getCurrentMiliUnixTimeStamp()+duration*1000)
	} else {
		return ErrInvalidArgs
	}
}

func (ek *ExpireStructure) PEXPIREBY(key string, duration int64, tag ...string) error {
	if len(tag) == 0 {
		return ek.setExpireToDB(key, duration)
	} else if len(tag) == 1 && tag[0] == "ex" {
		return ek.setExpireToDB(key, ek.getCurrentMiliUnixTimeStamp()+duration)
	} else {
		return ErrInvalidArgs
	}
}

func (ek *ExpireStructure) EXPIREBYAT(key string, timestamp int64, tag ...string) error {
	if len(tag) == 0 {
		return ek.setExpireToDB(key, timestamp*1000)
	} else if len(tag) == 1 && tag[0] == "ex" {
		return ek.setExpireToDB(key, ek.getCurrentMiliUnixTimeStamp()+timestamp*1000)
	} else {
		return ErrInvalidArgs
	}
}

func (ek *ExpireStructure) PEXPIREBYAT(key string, timestamp int64, tag ...string) error {
	if len(tag) == 0 {
		return ek.setExpireToDB(key, timestamp)
	} else if len(tag) == 1 && tag[0] == "ex" {
		return ek.setExpireToDB(key, ek.getCurrentMiliUnixTimeStamp()+timestamp)
	} else {
		return ErrInvalidArgs
	}
}

func (ek *ExpireStructure) getCurrentMiliUnixTimeStamp() int64 {
	now := time.Now()
	return now.UnixNano() / int64(time.Millisecond)
}

// getExpireFromDB retrieves data from the database.
func (ek *ExpireStructure) getExpireFromDB(key string, isKeyCanNotExist bool) (int64, error) {
	// Get data corresponding to the key from the database
	dbData, err := ek.db.Get([]byte(key))

	// Since the key might not exist, we need to handle ErrKeyNotFound separately as it is a valid case
	if err != nil && err != _const.ErrKeyNotFound {
		return 0, err
	}

	// Deserialize the data into a list
	deadtime, err := ek.decodeExpire(dbData)
	if err != nil {
		if len(dbData) != 0 {
			return 0, err
		} else {
			deadtime = -2
		}
	}
	return deadtime, nil
}

// setExpireToDB stores the data into the database.
func (ek *ExpireStructure) setExpireToDB(key string, deadtime int64) error {
	// Serialize into binary array
	encValue, err := ek.encodeExpire(deadtime)
	if err != nil {
		return err
	}
	// Store in the database
	return ek.db.Put([]byte(key), encValue)
}

// encodeExpire encodes the value
// format: [type][data]
func (ek *ExpireStructure) encodeExpire(deadtime int64) ([]byte, error) {
	// Create a byte slice buf with a length of binary.MaxVarintLen64
	// to hold the encoded value and additional data.
	buf := make([]byte, binary.MaxVarintLen64)

	// Set the first element of buf to represent the data structure type as String.
	buf[0] = Expire

	// Use the variable bufIndex to keep track of the current index position in buf,
	// starting from 1 to indicate the number of bytes written so far.
	var bufIndex = 1

	// Encode the expiration time expire as a variable-length integer
	// and write it to the sub-slice of byte slice buf starting
	// from the current index position bufIndex.
	bufIndex += binary.PutVarint(buf[bufIndex:], deadtime)

	return buf[:bufIndex], nil
}

// decodeList decodes the value
func (ek *ExpireStructure) decodeExpire(value []byte) (int64, error) {
	// Check the length of the value
	if len(value) < 1 {
		return 0, ErrInvalidValue
	}

	// Check the type of the value
	if value[0] != Expire {
		return 0, ErrInvalidType
	}

	// Use the variable bufIndex to keep track of the current index position in value,
	// starting from 1 to indicate the number of bytes read so far.
	var bufIndex = 1

	// Decode the expiration time expire from the sub-slice of byte slice value
	// starting from the current index position bufIndex.
	deadtime, n := binary.Varint(value[bufIndex:])

	// Check the number of bytes read
	if n <= 0 {
		return 0, ErrInvalidValue
	}

	// Return the original value value
	return deadtime, nil
}

func (ek *ExpireStructure) Stop() error {
	err := ek.db.Close()
	return err
}
