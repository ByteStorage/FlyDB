// CopyRight: GolQI
// Created by GolQI on 2023/4/1

package engine

import (
	"github.com/ByteStorage/FlyDB/config"
	data2 "github.com/ByteStorage/FlyDB/engine/data"
	"github.com/ByteStorage/FlyDB/engine/index"
	"github.com/ByteStorage/FlyDB/lib/const"
	s "github.com/ByteStorage/FlyDB/lib/proto/dbs"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// DB represents a FlyDB database instance,
// a high-performance storage engine built on the bitcask model.
// FlyDB utilizes a log-structured storage approach,
// which optimizes data operations by efficiently managing writes, reads, and deletes.
//
// FlyDB shines in its ability to swiftly handle data operations,
// offering exceptional performance. With just a single disk IO,
// FlyDB completes each operation efficiently,
// resulting in lightning-fast response times.
//
// One of the standout features of FlyDB is its efficient memory utilization.
// It employs an in-memory index that stores key-value mappings,
// facilitating rapid data access. This indexing mechanism enables FlyDB
// bto quickly locate data within the storage structure. However,
// it's important to note that the total size of data that can be stored is
// limited by the available memory capacity.
//
// When your system's memory can accommodate a significant portion
// of the key-value pairs, FlyDB emerges as an excellent choice of storage engine.
// It not only delivers outstanding speed and responsiveness
// but also minimizes disk IO overhead, allowing for seamless data operations.
//
// FlyDB's design philosophy revolves around balancing performance
// and memory usage. By optimizing disk IO and leveraging in-memory indexing,
// FlyDB provides a powerful and efficient storage solution for applications
// that prioritize speed and responsiveness.
type DB struct {
	// gRPC dbs
	s.FlyDBServiceServer
	options    config.Options
	lock       *sync.RWMutex
	fileIds    []int                      // File id, which can only be used when the index is loaded
	activeFile *data2.DataFile            // The current active data file that can be used for writing
	olderFiles map[uint32]*data2.DataFile // Old data file that can only be read
	index      index.Indexer              // Memory index
	transSeqNo uint64                     // Transaction sequence number, globally increasing
	isMerging  bool                       // Whether are merging
}

// NewDB open a new db instance
func NewDB(options config.Options) (*DB, error) {
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
		olderFiles: make(map[uint32]*data2.DataFile),
		index:      index.NewIndexer(options.IndexType, options.DirPath),
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

func checkOptions(options config.Options) error {
	if options.DirPath == "" {
		return _const.ErrOptionDirPathIsEmpty
	}
	if options.DataFileSize <= 0 {
		return _const.ErrOptionDataFileSizeNotPositive
	}
	return nil
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
		return _const.ErrKeyIsEmpty
	}

	// check LogRecord
	logRecord := &data2.LogRecord{
		Key:   encodeLogRecordKeyWithSeq(key, nonTransactionSeqNo),
		Value: value,
		Type:  data2.LogRecordNormal,
	}

	// append log record
	pos, err := db.appendLogRecordWithLock(logRecord)
	if err != nil {
		return err
	}

	// update index
	if ok := db.index.Put(key, pos); !ok {
		return _const.ErrIndexUpdateFailed
	}

	return nil
}

// appendLogRecord ethod added lock logic split,
// to avoid batch write resulting in deadlock problems
func (db *DB) appendLogRecordWithLock(logRecord *data2.LogRecord) (*data2.LogRecordPst, error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	return db.appendLogRecord(logRecord)
}

// appendLogRecord Append data to a file
func (db *DB) appendLogRecord(logRecord *data2.LogRecord) (*data2.LogRecordPst, error) {
	// Check whether the active data file exists
	// Initializes the data file if empty
	if db.activeFile == nil {
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	// Write data coding
	encRecord, size := data2.EncodeLogRecord(logRecord)
	if db.activeFile.WriteOff+size > db.options.DataFileSize {
		// Persisting data files to ensure that existing data is persisted to disk
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}

		// Converts the current active file to the old data file
		db.olderFiles[db.activeFile.FileID] = db.activeFile

		// Open a new active file
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	writeOff := db.activeFile.WriteOff
	if err := db.activeFile.Write(encRecord); err != nil {
		return nil, err
	}

	// Determines whether to initialize based on user configuration
	if db.options.SyncWrite {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
	}

	// Build in-memory index information
	pst := &data2.LogRecordPst{
		Fid:    db.activeFile.FileID,
		Offset: writeOff,
	}
	return pst, nil

}

// Set the current active file
// Hold a mutex before accessing this method
func (db *DB) setActiveDataFile() error {
	var initialFileID uint32 = 0
	if db.activeFile != nil {
		initialFileID = db.activeFile.FileID + 1
	}

	// Open a new data file
	dataFile, err := data2.OpenDataFile(db.options.DirPath, initialFileID,
		db.options.DataFileSize, db.options.FIOType)
	if err != nil {
		return err
	}
	db.activeFile = dataFile
	return nil
}

// Get Read data according to the key
func (db *DB) Get(key []byte) ([]byte, error) {
	zap.L().Info("get", zap.ByteString("key", key))
	db.lock.RLock()
	defer db.lock.RUnlock()

	// Determine the validity of the key
	if len(key) == 0 {
		return nil, _const.ErrKeyIsEmpty
	}

	// Retrieves the index of the key from the memory data structure
	logRecordPst := db.index.Get(key)
	// If key is not in the memory index, it does not exist
	if logRecordPst == nil {
		return nil, _const.ErrKeyNotFound
	}

	// Gets the value from the data file
	return db.getValueByPosition(logRecordPst)
}

// GetListKeys Gets all the keys in the database
func (db *DB) GetListKeys() [][]byte {
	// Retrieve an iterator for the index
	iterator := db.index.Iterator(false)

	// Create a slice to store the keys
	keys := make([][]byte, db.index.Size())

	var idx int
	// Iterate over the index
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		// Retrieve the key from the current iterator position
		keys[idx] = iterator.Key()
		idx++
	}

	// Return the list of keys
	return keys
}

// Fold get all the data and perform the operation specified by the user.
// The function returns false to exit
func (db *DB) Fold(f func(key []byte, value []byte) bool) error {
	// Acquire a read lock to ensure data consistency
	db.lock.RLock()
	defer db.lock.RUnlock()

	// Retrieve an iterator for the index
	iterator := db.index.Iterator(false)

	// Iterate over the index
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		// Retrieve the value associated with the current key
		value, err := db.getValueByPosition(iterator.Value())
		if err != nil {
			return err
		}

		// Invoke the provided function with the key and value
		// If the function returns false, stop folding and exit the loop
		if !f(iterator.Key(), value) {
			break
		}
	}

	// Return nil to indicate successful folding
	return nil
}

// getValueByPosition Get the corresponding value based on the location index information
func (db *DB) getValueByPosition(logRecordPst *data2.LogRecordPst) ([]byte, error) {
	// Find the corresponding data file according to the file id
	var dataFile *data2.DataFile
	if logRecordPst.Fid == db.activeFile.FileID {
		dataFile = db.activeFile
	} else {
		dataFile = db.olderFiles[logRecordPst.Fid]
	}

	// The data file is empty
	if dataFile == nil {
		return nil, _const.ErrDataFailNotFound
	}

	// The corresponding data is read according to the offset
	logRecord, _, err := dataFile.ReadLogRecord(logRecordPst.Offset)
	if err != nil {
		return nil, nil
	}
	if logRecord.Type == data2.LogRecordDeleted {
		return nil, _const.ErrKeyNotFound
	}

	return logRecord.Value, nil
}

// Delete data according to the key
func (db *DB) Delete(key []byte) error {
	zap.L().Info("delete", zap.ByteString("key", key))

	// Determine the validity of the key
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}

	// Check whether the key exists. If it does not exist, return it
	if pst := db.index.Get(key); pst == nil {
		return nil
	}

	// Construct a logRecord to indicate that it was deleted
	logRecord := &data2.LogRecord{
		Key:  encodeLogRecordKeyWithSeq(key, nonTransactionSeqNo),
		Type: data2.LogRecordDeleted,
	}

	// Write to the data file
	_, err := db.appendLogRecordWithLock(logRecord)
	if err != nil {
		return err
	}

	// Removes key from memory index
	ok := db.index.Delete(key)
	if !ok {
		return _const.ErrIndexUpdateFailed
	}
	return nil
}

// Load the data file from disk
func (db *DB) loadDataFiles() error {
	dirEntry, err := os.ReadDir(db.options.DirPath)
	if err != nil {
		return nil
	}

	var fileIds []int
	// Walk through all the files in the directory, finding all files ending in '.data'
	for _, entry := range dirEntry {
		if strings.HasSuffix(entry.Name(), data2.DataFileSuffix) {
			splitNames := strings.Split(entry.Name(), ".")
			fileID, err := strconv.Atoi(splitNames[0])
			// The data directory may be corrupted
			if err != nil {
				return _const.ErrDataDirectoryCorrupted
			}

			fileIds = append(fileIds, fileID)
		}
	}

	// Sort file ids and load them from smallest to largest
	sort.Ints(fileIds)
	db.fileIds = fileIds

	// Walk through each file id and open the corresponding data file
	for i, fid := range fileIds {
		dataFile, err := data2.OpenDataFile(db.options.DirPath, uint32(fid), db.options.DataFileSize, db.options.FIOType)
		if err != nil {
			return err
		}
		if i == len(fileIds)-1 {
			// The last id is the largest, indicating that the current file is active
			db.activeFile = dataFile
		} else {
			// Note It is an old data file
			db.olderFiles[uint32(fid)] = dataFile
		}
	}
	return nil
}

// Load the index from the data file
// Iterate over all the records in the file and update them to the memory index
func (db *DB) loadIndexFromDataFiles() error {
	// If there is no file, the database is empty
	if len(db.fileIds) == 0 {
		return nil
	}

	// Check whether the merge occurred
	var hasMerge bool = false
	var nonMergeFileId uint32 = 0
	mergeFileName := filepath.Join(db.options.DirPath, data2.MergeFinaFileSuffix)
	// If a file exists, retrieve the id of the file that did not participate in the merge
	if _, err := os.Stat(mergeFileName); err == nil {
		// Check if the merge file exists
		// If it exists, determine the ID of the most recently non-merged file
		fileId, err := db.getRecentlyNonMergeFileId(db.options.DirPath)
		if err != nil {
			return err
		}
		nonMergeFileId = fileId
		hasMerge = true
	}

	// Define a function to update the in-memory index
	updataIndex := func(key []byte, typ data2.LogRecrdType, pst *data2.LogRecordPst) {
		var ok bool
		if typ == data2.LogRecordDeleted {
			// If the log record type is 'deleted', delete the key from the index
			ok = db.index.Delete(key)
		} else {
			// Otherwise, update the key with the new position in the index
			ok = db.index.Put(key, pst)
		}
		if !ok {
			// Panic if the index update fails
			panic(_const.ErrIndexUpdateFailed)
		}
	}

	// Temporary transaction data
	transactionRecords := make(map[uint64][]*data2.TransactionRecord)
	var currentSeqNo = nonTransactionSeqNo

	// Iterate through all file ids, processing records in the file
	for i, fid := range db.fileIds {
		var fileID = uint32(fid)
		// If the id is smaller than that of the file that did not participate in the merge recently,
		// the hint file has been loaded
		if hasMerge && fileID < nonMergeFileId {
			continue
		}

		var dataFile *data2.DataFile
		if fileID == db.activeFile.FileID {
			dataFile = db.activeFile
		} else {
			dataFile = db.olderFiles[fileID]
		}

		// Obtain data
		var offset int64 = 0
		for {
			logRecord, size, err := dataFile.ReadLogRecord(offset)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			// Construct index memory and save it
			logRecordPst := &data2.LogRecordPst{
				Fid:    fileID,
				Offset: offset,
			}

			// Parse the key and get the transaction sequence number
			realKey, seqNo := parseLogRecordKeyAndSeq(logRecord.Key)
			if seqNo == nonTransactionSeqNo {
				// Non-transactional operation
				updataIndex(realKey, logRecord.Type, logRecordPst)
			} else {
				// When the transaction completes, update the corresponding seqNo data in the in-memory index
				if logRecord.Type == data2.LogRecordTransFinished {
					// Update the in-memory index with the transaction records
					for _, transRecord := range transactionRecords[seqNo] {
						updataIndex(transRecord.Record.Key, transRecord.Record.Type, transRecord.Pos)
					}
					// Remove the transaction records from the map
					delete(transactionRecords, seqNo)
				} else {
					// Batch submission, unsure if the transaction has been completed, temporarily store the records
					logRecord.Key = realKey
					transactionRecords[seqNo] = append(transactionRecords[seqNo], &data2.TransactionRecord{
						Record: logRecord,
						Pos:    logRecordPst,
					})
				}
			}

			// Update transaction sequence number
			if seqNo > currentSeqNo {
				currentSeqNo = seqNo
			}

			// Increments offset, next read from a new position
			offset += size
		}

		// If it is a current active file, update writeOff for this file
		if i == len(db.fileIds)-1 {
			db.activeFile.WriteOff = offset
		}
	}

	// Update the transaction sequence number to the database field
	db.transSeqNo = currentSeqNo

	return nil
}

// Clean the DB data directory after the test is complete
func (db *DB) Clean() {
	if db != nil {
		_ = db.Close()
		err := os.RemoveAll(db.options.DirPath)
		if err != nil {
			panic(err)
		}
	}
}
