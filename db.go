package flydb

import (
	"github.com/qishenonly/flydb/data"
	"github.com/qishenonly/flydb/index"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// DB bitcask 存储引擎实例
type DB struct {
	options    Options
	lock       *sync.RWMutex
	fileIds    []int                     //文件id，只能在加载索引的时候使用
	activeFile *data.DataFile            //当前的活跃数据文件，可以用于写入
	olderFiles map[uint32]*data.DataFile //旧的数据文件，只能用于读
	index      index.Indexer             //内存索引
	transSeqNo uint64                    //事务序列号，全局递增
	isMerging  bool                      //是否正在 merge
}

// NewFlyDB open a new db instance
func NewFlyDB(options Options) (*DB, error) {
	zap.L().Info("open db", zap.Any("options", options))
	// check options first
	if err := checkOptions(options); err != nil {
		return nil, err
	}

	// check data dir, if not exist, create it
	if _, err := os.Stat(options.DirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(options.DirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// init db instance
	db := &DB{
		options:    options,
		lock:       new(sync.RWMutex),
		olderFiles: make(map[uint32]*data.DataFile),
		index:      index.NewIndexer(options.IndexType),
	}

	// load merge files
	if err := db.loadMergeFiles(); err != nil {
		return nil, err
	}

	// load data files
	if err := db.loadDataFiles(); err != nil {
		return nil, err
	}

	// load index from hint file
	if err := db.loadIndexFromHintFile(); err != nil {
		return nil, err
	}

	// load index from data files
	if err := db.loadIndexFromDataFiles(); err != nil {
		return nil, err
	}

	return db, nil
}

// Close the db instance
func (db *DB) Close() error {
	zap.L().Info("close db", zap.Any("options", db.options))
	if db.activeFile == nil {
		return nil
	}
	db.lock.Lock()
	defer db.lock.Unlock()

	// close active file
	if err := db.activeFile.Close(); err != nil {
		return err
	}
	// close older files
	for _, file := range db.olderFiles {
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Sync the db instance
func (db *DB) Sync() error {
	zap.L().Info("sync db", zap.Any("options", db.options))
	if db.activeFile == nil {
		return nil
	}
	db.lock.Lock()
	defer db.lock.Unlock()
	return db.activeFile.Sync()
}

// Put write a key-value pair to db, and the key must be not empty
func (db *DB) Put(key []byte, value []byte) error {
	zap.L().Info("put", zap.ByteString("key", key), zap.ByteString("value", value))
	// check key
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	// check LogRecord
	logRecord := &data.LogRecord{
		Key:   encodeLogRecordKeyWithSeq(key, nonTransactionSeqNo),
		Value: value,
		Type:  data.LogRecordNormal,
	}

	// append log record
	pos, err := db.appendLogRecordWithLock(logRecord)
	if err != nil {
		return err
	}

	// update index
	if ok := db.index.Put(key, pos); !ok {
		return ErrIndexUpdateFailed
	}

	return nil
}

// appendLogRecord方法加锁逻辑拆分，避免批量写入时导致死锁问题
func (db *DB) appendLogRecordWithLock(logRecord *data.LogRecord) (*data.LogRecordPst, error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	return db.appendLogRecord(logRecord)
}

// appendLogRecord 追加数据写入到文件当中
func (db *DB) appendLogRecord(logRecord *data.LogRecord) (*data.LogRecordPst, error) {
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
	zap.L().Info("get", zap.ByteString("key", key))
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

	// 从数据文件中获取 value
	return db.getValueByPosition(logRecordPst)
}

// GetListKeys 获取数据库中所有的 key
func (db *DB) GetListKeys() [][]byte {
	iterator := db.index.Iterator(false)
	keys := make([][]byte, db.index.Size())
	var idx int
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		keys[idx] = iterator.Key()
		idx++
	}
	return keys
}

// Fold 获取所有的数据，并执行用户指定的操作，函数返回 false 退出
func (db *DB) Fold(f func(key []byte, value []byte) bool) error {
	db.lock.RLock()
	defer db.lock.RUnlock()

	iterator := db.index.Iterator(false)
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		value, err := db.getValueByPosition(iterator.Value())
		if err != nil {
			return err
		}
		if !f(iterator.Key(), value) {
			break
		}
	}
	return nil
}

// getValueByPosition 根据位置索引信息获取对应的 value
func (db *DB) getValueByPosition(logRecordPst *data.LogRecordPst) ([]byte, error) {
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
	logRecord, _, err := dataFile.ReadLogRecord(logRecordPst.Offset)
	if err != nil {
		return nil, nil
	}
	if logRecord.Type == data.LogRecordDeleted {
		return nil, ErrKeyNotFound
	}

	return logRecord.Value, nil
}

func (db *DB) Delete(key []byte) error {
	zap.L().Info("delete", zap.ByteString("key", key))
	// 判断 key 的有效性
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	// 检查 key 是否存在，如果不存在直接返回
	if pst := db.index.Get(key); pst == nil {
		return nil
	}

	// 构造 logRecord，标识其是被删除的
	logRecord := &data.LogRecord{
		Key:  encodeLogRecordKeyWithSeq(key, nonTransactionSeqNo),
		Type: data.LogRecordDeleted,
	}

	// 写入到数据文件中
	_, err := db.appendLogRecordWithLock(logRecord)
	if err != nil {
		return err
	}

	//从内存索引中将 key 删除
	ok := db.index.Delete(key)
	if !ok {
		return ErrIndexUpdateFailed
	}
	return nil
}

func checkOptions(options Options) error {
	if options.DirPath == "" {
		return ErrOptionDirPathIsEmpty
	}
	if options.DataFileSize <= 0 {
		return ErrOptionDataFileSizeNotPositive
	}
	return nil
}

// 从磁盘中加载数据文件
func (db *DB) loadDataFiles() error {
	dirEntry, err := os.ReadDir(db.options.DirPath)
	if err != nil {
		return nil
	}

	var fileIds []int
	// 遍历目录中的所有文件，找到所有以 .data 结尾的文件
	for _, entry := range dirEntry {
		if strings.HasSuffix(entry.Name(), data.DataFileSuffix) {
			splitNames := strings.Split(entry.Name(), ".")
			fileID, err := strconv.Atoi(splitNames[0])
			//数据目录可能损坏
			if err != nil {
				return ErrDataDirectoryCorrupted
			}

			fileIds = append(fileIds, fileID)
		}
	}

	// 对文件 id 进行排序， 从小到大依次加载
	sort.Ints(fileIds)
	db.fileIds = fileIds

	// 遍历每个文件 id， 打开对应的数据文件
	for i, fid := range fileIds {
		dataFile, err := data.OpenDataFile(db.options.DirPath, uint32(fid))
		if err != nil {
			return err
		}
		if i == len(fileIds)-1 {
			//最后一个id是最大的，说明当前文件是活跃文件
			db.activeFile = dataFile
		} else {
			//说明是旧的数据文件
			db.olderFiles[uint32(fid)] = dataFile
		}
	}
	return nil
}

// 从数据文件中加载索引
// 遍历文件中的所有记录，并更新到内存索引中
func (db *DB) loadIndexFromDataFiles() error {
	// 没有文件，说明数据库是空的
	if len(db.fileIds) == 0 {
		return nil
	}

	// 查看是否发生过 merge
	var hasMerge bool = false
	var nonMergeFileId uint32 = 0
	mergeFileName := filepath.Join(db.options.DirPath, data.MergeFinaFileSuffix)
	// 存在文件，则取出没有参与过 merge 的文件 id
	if _, err := os.Stat(mergeFileName); err == nil {
		fileId, err := db.getRecentlyNonMergeFileId(db.options.DirPath)
		if err != nil {
			return err
		}
		nonMergeFileId = fileId
		hasMerge = true
	}

	updataIndex := func(key []byte, typ data.LogRecrdType, pst *data.LogRecordPst) {
		var ok bool
		if typ == data.LogRecordDeleted {
			ok = db.index.Delete(key)
		} else {
			ok = db.index.Put(key, pst)
		}
		if !ok {
			panic(ErrIndexUpdateFailed)
		}
	}

	// 暂存事务数据
	transactionRecords := make(map[uint64][]*data.TransactionRecord)
	var currentSeqNo = nonTransactionSeqNo

	// 遍历所有文件id，处理文件中的记录
	for i, fid := range db.fileIds {
		var fileID = uint32(fid)
		// 如果比最近未参与 merge 的文件 id 更小，则说明已经从 hint 文件中加载过了
		if hasMerge && fileID < nonMergeFileId {
			continue
		}

		var dataFile *data.DataFile
		if fileID == db.activeFile.FileID {
			dataFile = db.activeFile
		} else {
			dataFile = db.olderFiles[fileID]
		}

		// 获取数据
		var offset int64 = 0
		for {
			logRecord, size, err := dataFile.ReadLogRecord(offset)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			// 构造索引内存并保存
			logRecordPst := &data.LogRecordPst{
				Fid:    fileID,
				Offset: offset,
			}

			// 解析 key，拿到事务序列号
			realKey, seqNo := parseLogRecordKeyAndSeq(logRecord.Key)
			if seqNo == nonTransactionSeqNo {
				// 非事务操作
				updataIndex(realKey, logRecord.Type, logRecordPst)
			} else {
				// 事务完成， 对应的 seqNo 数据可以更新到内存索引中
				if logRecord.Type == data.LogRecordTransFinished {
					for _, transRecord := range transactionRecords[seqNo] {
						updataIndex(transRecord.Record.Key, transRecord.Record.Type, transRecord.Pos)
					}
					delete(transactionRecords, seqNo)
				} else {
					// batch 中提交，不知道事务是否已完成，先暂存
					logRecord.Key = realKey
					transactionRecords[seqNo] = append(transactionRecords[seqNo], &data.TransactionRecord{
						Record: logRecord,
						Pos:    logRecordPst,
					})
				}
			}

			// 更新事务序列号
			if seqNo > currentSeqNo {
				currentSeqNo = seqNo
			}

			// 递增offset，下一次从新的位置读取
			offset += size
		}

		// 如果是当前活跃文件，更新这个文件的writeOff
		if i == len(db.fileIds)-1 {
			db.activeFile.WriteOff = offset
		}
	}

	// 更新事务序列号到数据库字段
	db.transSeqNo = currentSeqNo

	return nil
}
