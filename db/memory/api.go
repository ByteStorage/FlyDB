package memory

type Api interface {
	CreateColumnFamily(name string) error
	DropColumnFamily(name string) error
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Keys() ([][]byte, error)
	Close() error
	ListColumnFamilies() ([]string, error)
	PutCF(cf string, key []byte, value []byte) error
	GetCF(cf string, key []byte) ([]byte, error)
	DeleteCF(cf string, key []byte) error
	KeysCF(cf string) ([][]byte, error)
}

// Put: We can divide one column to many column families. Each column family has its own memTable and SSTable.
// Like MySQL divide one table to many partitions. Each partition has its own value and index.
