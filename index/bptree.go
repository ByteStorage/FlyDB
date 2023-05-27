package index

import (
	"github.com/qishenonly/flydb/data"
	"go.etcd.io/bbolt"
	"path/filepath"
)

var _ Indexer = (*BPlusTree)(nil)

const bPlusTreeIndexFileName = "bptree-index"

var indexBucketName = []byte("flydb-buckte-index")

// BPlusTree B+ Tree Index
// go.etcd.io/bbolt  This is the library that encapsulates b+ tree
// Again, if you need to look at the source code for b+ trees,
// The following link is a good place to start
// https://github.com/etcd-io/bbolt
type BPlusTree struct {
	tree *bbolt.DB
}

// NewBPlusTree Initializes the B+ tree index
func NewBPlusTree(dirPath string) *BPlusTree {
	bptree, err := bbolt.Open(filepath.Join(dirPath, bPlusTreeIndexFileName), 0644, nil)
	if err != nil {
		panic(ErrOpenBPTreeFailed)
	}

	// Create the corresponding bucket
	// The update method can be thought of as a single transaction,
	// and all operations within the method are committed as a single transaction.
	// There is a bucket parameter in the transaction,
	// which can be interpreted as partitioning the data.
	// After creating a bucket, a bucket is returned.
	// The returned bucket can be used to Put, Get and other methods.
	if err := bptree.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(indexBucketName)
		return err
	}); err != nil {
		panic(ErrCreateBucketFailed)
	}

	return &BPlusTree{
		tree: bptree,
	}
}

// Put Inserts a key-value pair into the B+ tree index
// The two arguments to the Put method are required to be byte arrays
// The first argument is the key, and the second argument is the value
// The key is the primary key of the data,
// and the value is the offset of the data in the data file
func (bptree *BPlusTree) Put(key []byte, pst *data.LogRecordPst) bool {
	if err := bptree.tree.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(indexBucketName)
		return bucket.Put(key, data.EncodeLogRecordPst(pst))
	}); err != nil {
		panic(ErrPutValueFailed)
	}
	return true
}

// Get Gets the value corresponding to the key from the B+ tree index
// The argument to the Get method is required to be a byte array
// The argument is the key, and the return value is the value corresponding to the key
// If the key does not exist, nil is returned
func (bptree *BPlusTree) Get(key []byte) *data.LogRecordPst {
	var pst *data.LogRecordPst
	// The view method allows only reads, not inserts and deletes.
	if err := bptree.tree.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(indexBucketName)
		value := bucket.Get(key)
		if len(value) != 0 {
			pst = data.DecodeLogRecordPst(value)
		}
		return nil
	}); err != nil {
		panic(ErrGetValueFailed)
	}
	return pst
}

// Delete Deletes the key-value pair corresponding to the key from the B+ tree index
// The argument to the Delete method is required to be a byte array
// The argument is the key, and the return value is a bool value
// If the key does not exist, false is returned
func (bptree *BPlusTree) Delete(key []byte) bool {
	var ok bool
	if err := bptree.tree.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(indexBucketName)
		if value := bucket.Get(key); len(value) != 0 {
			ok = true
			return bucket.Delete(key)
		}
		return nil
	}); err != nil {
		panic(ErrDeleteValueFailed)
	}
	return ok
}

// Size Gets the number of key-value pairs in the B+ tree index
// The return value is an int value
// If the index is empty, 0 is returned
func (bptree *BPlusTree) Size() int {
	var size int
	if err := bptree.tree.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(indexBucketName)
		size = bucket.Stats().KeyN
		return nil
	}); err != nil {
		panic(ErrGetIndexSizeFailed)
	}
	return size
}

func (bptree *BPlusTree) Iterator(reverse bool) Iterator {
	//TODO implement me
	panic("implement me")
}

type bptreeIterator struct {
	tx      *bbolt.Tx
	cursor  *bbolt.Cursor
	reverse bool
}

var _ Iterator = (*bptreeIterator)(nil)

// newBptreeIterator Initializes the B+ tree index iterator
// The two arguments to the newBptreeIterator method are required to be byte arrays
// The first argument is the B+ tree index,
// and the second argument is the traversal direction of the iterator
// The return value is an iterator
func newBptreeIterator(tree *bbolt.DB, reverse bool) *bptreeIterator {
	tx, err := tree.Begin(false)
	if err != nil {
		panic(ErrBeginTxFailed)
	}
	return &bptreeIterator{
		tx:      tx,
		cursor:  tx.Bucket(indexBucketName).Cursor(),
		reverse: reverse,
	}
}

func (b bptreeIterator) Rewind() {
	//TODO implement me
	panic("implement me")
}

func (b bptreeIterator) Seek(key []byte) {
	//TODO implement me
	panic("implement me")
}

func (b bptreeIterator) Next() {
	//TODO implement me
	panic("implement me")
}

func (b bptreeIterator) Valid() bool {
	//TODO implement me
	panic("implement me")
}

func (b bptreeIterator) Key() []byte {
	//TODO implement me
	panic("implement me")
}

func (b bptreeIterator) Value() *data.LogRecordPst {
	//TODO implement me
	panic("implement me")
}

func (b bptreeIterator) Close() {
	//TODO implement me
	panic("implement me")
}
