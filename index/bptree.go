package index

import (
	"github.com/qishenonly/flydb"
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

func NewBPlusTree(dirPath string) *BPlusTree {
	bptree, err := bbolt.Open(filepath.Join(dirPath, bPlusTreeIndexFileName), 0644, nil)
	if err != nil {
		panic(flydb.ErrOpenBPTreeFailed)
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
		panic(flydb.ErrCreateBucketFailed)
	}

	return &BPlusTree{
		tree: bptree,
	}
}

func (bptree *BPlusTree) Put(key []byte, pst *data.LogRecordPst) bool {
	//TODO implement me
	panic("implement me")
}

func (bptree *BPlusTree) Get(key []byte) *data.LogRecordPst {
	//TODO implement me
	panic("implement me")
}

func (bptree *BPlusTree) Delete(key []byte) bool {
	//TODO implement me
	panic("implement me")
}

func (bptree *BPlusTree) Size() int {
	//TODO implement me
	panic("implement me")
}

func (bptree *BPlusTree) Iterator(reverse bool) Iterator {
	//TODO implement me
	panic("implement me")
}
