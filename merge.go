package flydb

import (
	"github.com/ByteStorage/flydb/data"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
)

// merge 文件夹名称
var (
	mergeDirName = "dbmerge"
	mergeFinaKey = "mergeFina.finished"
)

// Merge 清理无效数据，生成 hint 索引文件
func (db *DB) Merge() error {
	// 如果数据库为空，则直接返回
	if db.activeFile == nil {
		return nil
	}
	db.lock.Lock()
	// 如果 merge 正在进行中，则直接返回
	if db.isMerging {
		db.lock.Unlock()
		return ErrMergeIsProgress
	}
	db.isMerging = true
	defer func() {
		db.isMerging = false
	}()

	// 持久化当前活跃文件
	if err := db.activeFile.Sync(); err != nil {
		db.lock.Unlock()
		return err
	}

	// 将当前活跃文件转换为旧的数据文件
	db.olderFiles[db.activeFile.FileID] = db.activeFile
	// 打开一个新的活跃文件
	if err := db.setActiveDataFile(); err != nil {
		db.lock.Unlock()
		return nil
	}

	// 记录最近没有参与 merge 的文件
	noMergeFileId := db.activeFile.FileID

	// 取出所有需要 merge 的文件
	var mergeFiles []*data.DataFile
	for _, files := range db.olderFiles {
		mergeFiles = append(mergeFiles, files)
	}
	db.lock.Unlock()

	// 将 merge 文件从小到大排序
	sort.Slice(mergeFiles, func(i, j int) bool {
		return mergeFiles[i].FileID < mergeFiles[j].FileID
	})

	mergePath := db.getMergePath()
	// 如果目录存在，就说明 merge 过，需要删除
	if _, err := os.Stat(mergePath); err == nil {
		if err := os.RemoveAll(mergePath); err != nil {
			return err
		}
	}
	// 新建 merge 目录
	if err := os.MkdirAll(mergePath, os.ModePerm); err != nil {
		return err
	}

	// 打开一个临时的新的实例，并修改配置项
	mergeOptions := db.options
	mergeOptions.DirPath = mergePath
	mergeOptions.SyncWrite = false
	mergeDB, err := NewFlyDB(mergeOptions)
	if err != nil {
		return err
	}

	// 打开 hint 文件存储索引
	hintFile, err := data.OpenHintFile(mergePath)
	if err != nil {
		return err
	}
	// 遍历每个数据文件
	for _, files := range mergeFiles {
		var offset int64 = 0
		for {
			logRecord, size, err := files.ReadLogRecord(offset)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			// 解析拿到的 key
			realKey, _ := parseLogRecordKeyAndSeq(logRecord.Key)
			logRecordPst := db.index.Get(realKey)
			// 和内存中的索引位置进行比较，有效则重写
			if logRecordPst != nil && logRecordPst.Fid == files.FileID && logRecordPst.Offset == offset {
				// 清除事务标记
				logRecord.Key = encodeLogRecordKeyWithSeq(realKey, nonTransactionSeqNo)
				recordPst, err := mergeDB.appendLogRecord(logRecord)
				if err != nil {
					return err
				}

				// 将当前位置索引写到 hint 文件中
				if err := hintFile.WriteHintRecord(realKey, recordPst); err != nil {
					return err
				}
			}
			// 递增 offest
			offset += size
		}
	}

	// 持久化
	if err := hintFile.Sync(); err != nil {
		return err
	}
	if err := mergeDB.Sync(); err != nil {
		return err
	}

	// 写标识 merge 完成的文件
	mergeFinaFile, err := data.OpenMergeFinaFile(mergePath)
	if err != nil {
		return err
	}

	mergeFinaRecord := &data.LogRecord{
		Key:   []byte(mergeFinaKey),
		Value: []byte(strconv.Itoa(int(noMergeFileId))),
	}

	encRecord, _ := data.EncodeLogRecord(mergeFinaRecord)
	if err := mergeFinaFile.Write(encRecord); err != nil {
		return err
	}

	// 持久化
	if err := mergeFinaFile.Sync(); err != nil {
		return err
	}

	return nil

}

func (db *DB) getMergePath() string {
	// 获取数据库父级目录
	parentDir := path.Dir(path.Clean(db.options.DirPath))
	// DB base path
	basePath := path.Base(db.options.DirPath)
	// 返回 merge 文件路径
	return filepath.Join(parentDir, basePath+mergeDirName)
}

// 加载 merge 数据目录
func (db *DB) loadMergeFiles() error {
	mergePath := db.getMergePath()
	// merge 目录不存在的话直接返回
	if _, err := os.Stat(mergePath); os.IsNotExist(err) {
		return nil
	}
	defer func() {
		_ = os.RemoveAll(mergePath)
	}()

	dirs, err := os.ReadDir(mergePath)
	if err != nil {
		return err
	}

	// 查找标识 merge 的文件，判断 merge 是否完成
	var mergeFinished bool
	var mergeFileNames []string
	for _, dir := range dirs {
		if dir.Name() == data.MergeFinaFileSuffix {
			mergeFinished = true
		}
		mergeFileNames = append(mergeFileNames, dir.Name())
	}
	// 没有则直接返回
	if !mergeFinished {
		return nil
	}

	nonMergeFileID, err := db.getRecentlyNonMergeFileId(mergePath)
	if err != nil {
		return err
	}

	// 删除旧的数据文件
	var fileID uint32 = 0
	for ; fileID < nonMergeFileID; fileID++ {
		fileName := data.GetDataFileName(db.options.DirPath, fileID)
		if _, err := os.Stat(fileName); err == nil {
			if err := os.Remove(fileName); err != nil {
				return err
			}
		}
	}

	// 移动新的数据文件到数据目录中
	for _, fileName := range mergeFileNames {
		mergeSrcPath := filepath.Join(mergePath, fileName)
		dataSrcPath := filepath.Join(db.options.DirPath, fileName)
		if err := os.Rename(mergeSrcPath, dataSrcPath); err != nil {
			return err
		}
	}
	return nil
}

// 获取最近没有参与 merge 的文件 id
func (db *DB) getRecentlyNonMergeFileId(dirPath string) (uint32, error) {
	mergeFinaFile, err := data.OpenMergeFinaFile(dirPath)
	if err != nil {
		return 0, err
	}
	record, _, err := mergeFinaFile.ReadLogRecord(0)
	if err != nil {
		return 0, err
	}
	nonMergeFileID, err := strconv.Atoi(string(record.Value))
	if err != nil {
		return 0, err
	}
	return uint32(nonMergeFileID), nil
}

// 从 hint 文件中加载索引
func (db *DB) loadIndexFromHintFile() error {
	// 判断 hint 文件是否存在
	hintFileName := filepath.Join(db.options.DirPath, data.HintFileSuffix)
	if _, err := os.Stat(hintFileName); os.IsNotExist(err) {
		return nil
	}

	// 打开 hint 文件
	hintFile, err := data.OpenHintFile(db.options.DirPath)
	if err != nil {
		return err
	}

	// 读取文件中的索引
	var offset int64 = 0
	for {
		logRecord, size, err := hintFile.ReadLogRecord(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// 解码拿到实际的索引位置
		pst := data.DecodeLogRecordPst(logRecord.Value)
		db.index.Put(logRecord.Key, pst)
		offset += size
	}
	return nil
}
