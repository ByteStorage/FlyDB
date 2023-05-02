package flydb

import (
	"github.com/qishenonly/flydb/data"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
)

// merge 文件夹名称
var mergeDirName = "dbmerge"

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

	// 遍历每个数据文件
	for _, files := range mergeFiles {
		var offset int64 = 0
		for {
			logRecord, _, err := files.ReadLogRecord(offset)
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
				_, err := mergeDB.appendLogRecord(logRecord)
				if err != nil {
					return err
				}

			}

		}
	}

	return nil
}

func (db *DB) getMergePath() string {
	// 获取数据库父级目录
	parentDir := path.Dir(path.Clean(db.options.DirPath))
	// 返回 merge 文件路径
	return filepath.Join(parentDir + mergeDirName)
}

