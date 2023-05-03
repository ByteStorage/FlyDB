package flydb

import (
	"github.com/qishenonly/flydb/data"
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
	mergeDB, err := Open(mergeOptions)
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
	// 返回 merge 文件路径
	return filepath.Join(parentDir + mergeDirName)
}
