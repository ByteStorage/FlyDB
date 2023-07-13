package structure

import (
	"bytes"
	"container/heap"
	"encoding/binary"
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	"reflect"
)

// ZSetStructure is a structure for ZSet or SortedSet
type ZSetStructure struct {
	db *engine.DB
}
type ZSetNodes []*ZSetNode // implements heap.Interface and holds ZSetNode.

type ZSetNode struct {
	Value    string // The value of the item; arbitrary.
	Priority int    // The priority of the item in the queue.
	Index    int    // The index of the item in the heap.
}

func NewZSetStructure(options config.Options) (*ZSetStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}

	return &ZSetStructure{db: db}, nil
}
func (zs *ZSetStructure) ZAdd(key string, score int, value string) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	keyBytes := stringToBytesWithKey(key)
	_, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return err
	}
	return nil
}

func (pq ZSetNodes) Len() int { return len(pq) }

func (pq ZSetNodes) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

func (pq ZSetNodes) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *ZSetNodes) Push(x any) {
	n := len(*pq)
	item := x.(*ZSetNode)
	item.Index = n
	*pq = append(*pq, item)
	heap.Fix(pq, n)
}

func (pq *ZSetNodes) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *ZSetNodes) update(item *ZSetNode, value string, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}
func (pq *ZSetNodes) Bytes() ([]byte, error) {
	msgPack := encoding.InitMessagePack()

	err := msgPack.AddExtension(reflect.TypeOf(ZSetNode{}), 1, zSetNodesEncoder, zSetNodesDecoder)
	if err != nil {
		return nil, err
	}
	return msgPack.Encode(pq)
}
func (pq *ZSetNodes) FromBytes(bytes []byte) error {
	msgPack := encoding.InitMessagePack()
	err := msgPack.AddExtension(reflect.TypeOf(ZSetNode{}), 1, nil, zSetNodesDecoder)
	if err != nil {
		return err
	}
	return msgPack.Decode(bytes, pq)
}

func (l *ZSetStructure) getZSetFromDB(key []byte) (*ZSetNodes, error) {
	// Get data corresponding to the key from the database
	dbData, err := l.db.Get(key)

	// Since the key might not exist, we need to handle ErrKeyNotFound separately as it is a valid case
	if err != nil && err != _const.ErrKeyNotFound {
		return nil, err
	}
	var zSetValue ZSetNodes
	// Deserialize the data into a list
	err = encoding.DecodeMessagePack(dbData, zSetValue)
	if err != nil {
		return nil, err
	}
	return &zSetValue, nil
}
func (l *ZSetStructure) setZSetToDB(key []byte, zSetValue ZSetNodes) error {
	// Deserialize the data into a list
	val, err := encoding.EncodeMessagePack(zSetValue)
	if err != nil {
		return err
	}
	err = l.db.Put(key, val)
	if err != nil {
		return err
	}
	return nil
}

func zSetNodesDecoder(value reflect.Value, i []byte) error {
	bs := ZSetNode{}
	var bytesRead int
	num, s, err := encoding.DecodeString(i)
	if err != nil {
		return err
	}
	bytesRead += num
	bs.Value = s
	val, num := binary.Varint(i[bytesRead:])
	bytesRead += num
	bs.Index = int(val)
	val, num = binary.Varint(i[bytesRead:])
	bytesRead += num
	bs.Priority = int(val)
	value.Set(reflect.ValueOf(bs))
	return nil
}
func zSetNodesEncoder(value reflect.Value) ([]byte, error) {
	zsn := value.Interface().(ZSetNode)
	if zsn.Value == "" {
		return nil, errors.New("empty zset")
	}
	buf := bytes.NewBuffer(nil)
	es, err := encoding.EncodeString(zsn.Value)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(es)
	if err != nil {
		return nil, err
	}
	b := make([]byte, binary.MaxVarintLen64)
	written := 0
	written += binary.PutVarint(b[:], int64(zsn.Index))
	written += binary.PutVarint(b[written:], int64(zsn.Priority))
	buf.Write(b[:written])
	return buf.Bytes(), nil
}
