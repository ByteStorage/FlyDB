package engine

import (
	data2 "github.com/ByteStorage/FlyDB/engine/data"
	"github.com/ByteStorage/FlyDB/lib/const"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
)

// merge folder name
var (
	mergeDirName = "dbmerge"
	mergeFinaKey = "mergeFina.finished"
)

// Merge Clear the invalid data and generate the hint index file
func (db *DB) Merge() error {
	// If the database is empty, it is returned directly
	if db.activeFile == nil {
		return nil
	}
	db.lock.Lock()
	// If the merge is in progress, return directly
	if db.isMerging {
		db.lock.Unlock()
		return _const.ErrMergeIsProgress
	}
	db.isMerging = true
	defer func() {
		db.isMerging = false
	}()

	// Persist the currently active file
	if err := db.activeFile.Sync(); err != nil {
		db.lock.Unlock()
		return err
	}

	// Converts the currently active file to the old data file
	db.olderFiles[db.activeFile.FileID] = db.activeFile
	// Open a new active file
	if err := db.setActiveDataFile(); err != nil {
		db.lock.Unlock()
		return nil
	}

	// Records files that have not participated in the merge recently
	noMergeFileId := db.activeFile.FileID

	// Retrieve all files that need to be merged
	var mergeFiles []*data2.DataFile
	for _, files := range db.olderFiles {
		mergeFiles = append(mergeFiles, files)
	}
	db.lock.Unlock()

	// Sort merge files from smallest to largest
	sort.Slice(mergeFiles, func(i, j int) bool {
		return mergeFiles[i].FileID < mergeFiles[j].FileID
	})

	mergePath := db.getMergePath()
	// If the directory exists, it has been merged and needs to be deleted
	if _, err := os.Stat(mergePath); err == nil {
		if err := os.RemoveAll(mergePath); err != nil {
			return err
		}
	}
	// Creating a merge Directory
	if err := os.MkdirAll(mergePath, os.ModePerm); err != nil {
		return err
	}

	// Open a temporary new instance and modify the configuration item
	mergeOptions := db.options
	mergeOptions.DirPath = mergePath
	mergeOptions.SyncWrite = false
	mergeDB, err := NewDB(mergeOptions)
	if err != nil {
		return err
	}

	// Open the hint file storage index
	hintFile, err := data2.OpenHintFile(mergePath, db.options.DataFileSize, db.options.FIOType)
	if err != nil {
		return err
	}
	// Walk through each data file
	for _, files := range mergeFiles {
		var offset int64 = 0
		for {
			logRecord, size, err := files.ReadLogRecord(offset)

			// Check if there was an error while reading the log record
			if err != nil {
				// If the error is io.EOF, it means the end of the file was reached, so we break out of the loop
				if err == io.EOF {
					break
				}
				return err
			}

			// Parse the key
			realKey, _ := parseLogRecordKeyAndSeq(logRecord.Key)
			logRecordPst := db.index.Get(realKey)
			// Compare with the index position in memory, and rewrite if valid
			if logRecordPst != nil && logRecordPst.Fid == files.FileID && logRecordPst.Offset == offset {
				// Clear transaction flag
				logRecord.Key = encodeLogRecordKeyWithSeq(realKey, nonTransactionSeqNo)
				recordPst, err := mergeDB.appendLogRecord(logRecord)
				if err != nil {
					return err
				}

				// Writes the current location index to the hint file
				if err := hintFile.WriteHintRecord(realKey, recordPst); err != nil {
					return err
				}
			}
			// Incremental offest
			offset += size
		}
	}

	// persistence
	if err := hintFile.Sync(); err != nil {
		return err
	}
	if err := mergeDB.Sync(); err != nil {
		return err
	}

	// Write a file that identifies the merge completion
	mergeFinaFile, err := data2.OpenMergeFinaFile(mergePath, db.options.DataFileSize, db.options.FIOType)
	if err != nil {
		return err
	}

	mergeFinaRecord := &data2.LogRecord{
		Key:   []byte(mergeFinaKey),
		Value: []byte(strconv.Itoa(int(noMergeFileId))),
	}

	encRecord, _ := data2.EncodeLogRecord(mergeFinaRecord)
	if err := mergeFinaFile.Write(encRecord); err != nil {
		return err
	}

	// persistence
	if err := mergeFinaFile.Sync(); err != nil {
		return err
	}

	return nil

}

func (db *DB) getMergePath() string {
	// Gets the database parent directory
	parentDir := path.Dir(path.Clean(db.options.DirPath))
	// DB base path
	basePath := path.Base(db.options.DirPath)
	// Return the merge file path
	return filepath.Join(parentDir, basePath+mergeDirName)
}

// Load the merge data directory
func (db *DB) loadMergeFiles() error {
	mergePath := db.getMergePath()
	// Return the merge directory if it does not exist
	if _, err := os.Stat(mergePath); os.IsNotExist(err) {
		return nil
	}
	defer func() {
		_ = os.RemoveAll(mergePath)
	}()

	dirs, err := os.ReadDir(mergePath)

	// Check if there was an error while reading the directory
	if err != nil {
		return err
	}

	// Find the file that identifies the merge and determine whether the merge is complete
	var mergeFinished bool
	var mergeFileNames []string

	// Iterate over the directories
	for _, dir := range dirs {
		// Check if the directory name matches the merge finish file suffix
		if dir.Name() == data2.MergeFinaFileSuffix {
			mergeFinished = true
		}

		// Append the directory name to the mergeFileNames slice
		mergeFileNames = append(mergeFileNames, dir.Name())
	}

	// If not, return directly
	if !mergeFinished {
		return nil
	}

	nonMergeFileID, err := db.getRecentlyNonMergeFileId(mergePath)

	// Check if there was an error while retrieving the recently non-merge file ID
	if err != nil {
		return err
	}

	// Delete old data files
	var fileID uint32 = 0
	for ; fileID < nonMergeFileID; fileID++ {
		fileName := data2.GetDataFileName(db.options.DirPath, fileID)

		// Check if the file exists
		if _, err := os.Stat(fileName); err == nil {
			// Remove the file
			if err := os.Remove(fileName); err != nil {
				return err
			}
		}
	}

	// Move the new data file to the data directory
	for _, fileName := range mergeFileNames {
		mergeSrcPath := filepath.Join(mergePath, fileName)
		dataSrcPath := filepath.Join(db.options.DirPath, fileName)

		// Rename the file from mergeSrcPath to dataSrcPath
		if err := os.Rename(mergeSrcPath, dataSrcPath); err != nil {
			return err
		}
	}

	return nil

}

// Gets the id of the file that did not participate in the merge recently
func (db *DB) getRecentlyNonMergeFileId(dirPath string) (uint32, error) {
	mergeFinaFile, err := data2.OpenMergeFinaFile(dirPath, db.options.DataFileSize, db.options.FIOType)
	if err != nil {
		return 0, err
	}

	// Read the log record at offset 0 from mergeFinaFile
	record, _, err := mergeFinaFile.ReadLogRecord(0)
	if err != nil {
		return 0, err
	}

	// Convert the value of the log record to an integer
	nonMergeFileID, err := strconv.Atoi(string(record.Value))
	if err != nil {
		return 0, err
	}

	return uint32(nonMergeFileID), nil

}

// Load the index from the hint file
func (db *DB) loadIndexFromHintFile() error {
	// Check whether the hint file exists
	hintFileName := filepath.Join(db.options.DirPath, data2.HintFileSuffix)
	if _, err := os.Stat(hintFileName); os.IsNotExist(err) {
		return nil
	}

	// Open hint file
	hintFile, err := data2.OpenHintFile(db.options.DirPath, db.options.DataFileSize, db.options.FIOType)
	if err != nil {
		return err
	}

	// Read the index in the file
	var offset int64 = 0
	for {
		logRecord, size, err := hintFile.ReadLogRecord(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Decode to get the actual index location
		pst := data2.DecodeLogRecordPst(logRecord.Value)
		db.index.Put(logRecord.Key, pst)
		offset += size
	}
	return nil
}
