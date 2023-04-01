package flydb

import (
	"flydb/data"
	"flydb/index"
	"sync"
)

// DB bitcask 存储引擎实例
type DB struct {
	options    Options
	lock       *sync.RWMutex
	activeFile *data.DataFile            //当前的活跃数据文件，可以用于写入
	olderFiles map[uint32]*data.DataFile //旧的数据文件，只能用于读
	index      index.Indexer             //内存索引
}

// Put 写入key/value， key不能为空
func (db *DB) Put(key []byte, value []byte) error {
	//判断key是否有效
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	//构造 LogRecord 结构体
	logRecord := &data.LogRecord{
		Key:   key,
		Value: value,
		Type:  data.LogRecordNormal,
	}

	//追加写入到当前活跃文件当中
	pos, err := db.appendLogRecord(logRecord)
	if err != nil {
		return err
	}

	//更新内存索引
	if ok := db.index.Put(key, pos); !ok {
		return ErrIndexUpdataFailed
	}

	return nil
}

// appendLogRecord 追加数据写入到文件当中
func (db *DB) appendLogRecord(logRecord *data.LogRecord) (*data.LogRecordPst, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	//判断当前活跃数据文件是否存在
	//如果为空则初始化数据文件
	if db.activeFile == nil {
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	//写入数据编码
	encRecord, size := data.EncodeLogRecord(logRecord)
	if db.activeFile.WriteOff+size > db.options.DataFileSize {
		//持久化数据文件，保证已有的数据持久到磁盘当中
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}

		//将当前活跃文件转化成旧的数据文件
		db.olderFiles[db.activeFile.FileID] = db.activeFile

		//打开新的活跃文件
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	writeOff := db.activeFile.WriteOff
	if err := db.activeFile.Write(encRecord); err != nil {
		return nil, err
	}

	//根据用户配置决定是否初始化
	if db.options.SyncWrite {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
	}

	//构建内存索引信息
	pst := &data.LogRecordPst{
		Fid:    db.activeFile.FileID,
		Offset: writeOff,
	}
	return pst, nil

}

// 设置当前活跃文件
// 在访问此方法前必须持有互斥锁
func (db *DB) setActiveDataFile() error {
	var initialFileID uint32 = 0
	if db.activeFile != nil {
		initialFileID = db.activeFile.FileID + 1
	}

	//打开新的数据文件
	dataFile, err := data.OpenDataFile(db.options.DirPath, initialFileID)
	if err != nil {
		return err
	}
	db.activeFile = dataFile
	return nil
}

// Get 根据 key 读取数据
func (db *DB) Get(key []byte) ([]byte, error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	
	//判断 key 的有效性
	if len(key) == 0 {
		return nil, ErrKeyIsEmpty
	}

	//从内存数据结构中取出 key 对应的索引信息
	logRecordPst := db.index.Get(key)
	//如果 key 不在内存索引中， 说明 key 不存在
	if logRecordPst == nil {
		return nil, ErrKeyNotFound
	}

	//根据文件 id 找到对应的数据文件
	var dataFile *data.DataFile
	if logRecordPst.Fid == db.activeFile.FileID {
		dataFile = db.activeFile
	} else {
		dataFile = db.olderFiles[logRecordPst.Fid]
	}

	//数据文件为空
	if dataFile == nil {
		return nil, ErrDataFailNotFound
	}

	//根据偏移读取对应的数据
	logRecord, err := dataFile.ReadLogRecord(logRecordPst.Offset)
	if err != nil {
		return nil, nil
	}
	if logRecord.Type == data.LogRecordDeleted {
		return nil, ErrKeyNotFound
	}

	return logRecord.Value, nil
}
