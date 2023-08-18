package memory

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/engine"
	"os"
)

type Options struct {
	options  config.Options
	chanSize int
}

type DbWal struct {
	log *SequentialLogger
	db  *engine.DB
	ch  chan func()
	dir string
	mem *MemTable
}

func NewDbWal(option Options) (*DbWal, error) {
	db, err := engine.NewDB(option.options)
	if err != nil {
		return nil, err
	}

	log, err := NewSequentialLogger(option.options.DirPath + "/wal.log")
	if err != nil {
		return nil, err
	}

	d := &DbWal{
		log: log,
		db:  db,
		ch:  make(chan func(), 1000000),
		mem: NewMemTable(),
		dir: option.options.DirPath,
	}

	go d.asyncWorker()
	return d, nil
}

func (d *DbWal) asyncWorker() {
	for task := range d.ch {
		task()
	}
}

func (d *DbWal) PutByWal(key []byte, value []byte) error {

	//keySize := len(key)
	//valueSize := len(value)
	//logMsg := make([]byte, 4+4+keySize+valueSize) // Assuming 4 bytes for each size
	//binary.BigEndian.PutUint32(logMsg[0:4], uint32(keySize))
	//binary.BigEndian.PutUint32(logMsg[4:8], uint32(valueSize))
	//copy(logMsg[8:8+keySize], key)
	//copy(logMsg[8+keySize:], value)
	//err := d.log.Write(string(logMsg))
	//if err != nil {
	//	return err
	//}
	//
	//// memory update
	//d.mem.Put(key, value)
	//
	//// sync update db
	//d.ch <- func() {
	//	if err := d.db.Put(key, value); err != nil { // Add error handling
	//		fmt.Printf("Error updating DB: %v\n", err) // Log the error
	//	}
	//}
	return nil
}

func (d *DbWal) GetByWal(key []byte) ([]byte, error) {
	return d.db.Get(key)
}

func (d *DbWal) Close() error {
	return d.log.Flush()
}

func (d *DbWal) Clean() {
	if d.db != nil {
		_ = d.db.Close()
		err := os.RemoveAll(d.dir)
		if err != nil {
			panic(err)
		}
	}
}
