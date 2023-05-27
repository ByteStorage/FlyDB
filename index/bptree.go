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
// The two arguments to the Put method
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
	currKey []byte
	currVal []byte
}

var _ Iterator = (*bptreeIterator)(nil)

// newBptreeIterator Initializes the B+ tree index iterator
// The two arguments to the newBptreeIterator method
// The first argument is the B+ tree index,
// and the second argument is the traversal direction of the iterator
// The return value is an iterator
func newBptreeIterator(tree *bbolt.DB, reverse bool) *bptreeIterator {
	tx, err := tree.Begin(false)
	if err != nil {
		panic(ErrBeginTxFailed)
	}
	b := &bptreeIterator{
		tx:      tx,
		cursor:  tx.Bucket(indexBucketName).Cursor(),
		reverse: reverse,
	}
	b.Rewind()
	return b
}

// Rewind Resets the iterator to the beginning of the B+ tree index
func (b *bptreeIterator) Rewind() {
	if b.reverse {
		b.currKey, b.currVal = b.cursor.Last()
	} else {
		b.currKey, b.currVal = b.cursor.First()
	}
}

// Seek Positions the iterator to the first key
// that is greater or equal to the specified key
// The argument to the Seek method is required to be a byte array
// The argument is the key
// If the key does not exist, the iterator is positioned
// to the first key that is greater than the specified key
func (b *bptreeIterator) Seek(key []byte) {
	b.currKey, b.currVal = b.cursor.Seek(key)
}

// Next Positions the iterator to the next key
// If the iterator is positioned at the last key,
// the iterator is positioned to the end of the B+ tree index
// The Next method does not return a value
func (b *bptreeIterator) Next() {
	if b.reverse {
		b.currKey, b.currVal = b.cursor.Prev()
	} else {
		b.currKey, b.currVal = b.cursor.Next()
	}
}

// Valid Determines whether the iterator is positioned at a valid key
// The return value is a bool value
// If the iterator is positioned at a valid key, true is returned
// Otherwise, false is returned
func (b *bptreeIterator) Valid() bool {
	return b.currKey != nil
}

// Key Gets the key at the current iterator position
// The return value is a byte array
// If the iterator is not positioned at a valid key, nil is returned
// The Key method does not return an error
func (b *bptreeIterator) Key() []byte {
	return b.currKey
}

// Value Gets the value at the current iterator position
// The return value is LogRecordPst struct
// If the iterator is not positioned at a valid key, nil is returned
// The Value method does not return an error
func (b *bptreeIterator) Value() *data.LogRecordPst {
	return data.DecodeLogRecordPst(b.currVal)
}

// Close Closes the iterator
// The Close method does not return a value or an error
// The Close method rolls back the transaction
func (b bptreeIterator) Close() {
	if err := b.tx.Rollback(); err != nil {
		panic(ErrRollbackTxFailed)
	}
}
