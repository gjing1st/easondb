package easondb

import (
	"errors"
	"gitee.com/gjing1st/easondb/index"
	"gitee.com/gjing1st/easondb/model"
	"sync"
)

var (
	ErrEmptyKey     = errors.New("easondb: key不能为空")
	ErrKeyTooLong   = errors.New("easondb: key过长")
	ErrEmptyValue   = errors.New("easondb: value不能为空")
	ErrValueTooLong = errors.New("easondb: value过长")
)

type EasonDB struct {
	activeFile    *model.DBFile
	activeFileId  uint32
	archivedFiles ArchivedFiles
	strIndex      *StrIndex
	mutex         sync.Mutex
	config        Config
	meta          *model.DBMeta
}
type ArchivedFiles map[uint32]*model.DBFile

// validKeyValue
// @description: 检测key和value
// @param: key  value
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 11:28
// @success: nil
func (db *EasonDB) validKeyValue(key []byte, value ...[]byte) error {
	keySize := uint32(len(key))
	if keySize == 0 {
		return ErrEmptyKey
	}
	if keySize > db.config.MaxKeySize {
		return ErrKeyTooLong
	}
	for _, v := range value {
		valueSize := uint32(len(v))
		if valueSize == 0 {
			return ErrEmptyValue
		}
		if valueSize > db.config.MaxValueSize {
			return ErrValueTooLong
		}
	}

	return nil
}

// store
// @description: 数据持久化
// @param: entity 实体
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/20 16:16
func (db *EasonDB) store(entity *model.Entity) (err error) {
	config := db.config
	//如果数据文件空间不够，则关闭该文件，并新打开一个文件
	if db.activeFile.Offset+int64(entity.Size()) > config.BlockSize {
		if err = db.activeFile.Close(); err != nil {
			return
		}
		//活跃文件变归档文件
		db.archivedFiles[db.activeFileId] = db.activeFile
		db.activeFileId++
		dbFile, err := model.OpenDBFile(config.DbPath, db.activeFileId)
		if err != nil {
			return err
		}
		//新建文件为活跃文件
		db.activeFile = dbFile
		db.meta.ActiveWriteOffset = 0

	}
	if err = db.activeFile.Write(entity); err != nil {
		return
	}
	db.meta.ActiveWriteOffset = db.activeFile.Offset
	return
}

//
// @description: 创建索引
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 16:24
func (db *EasonDB) buildIndex(index *index.Indexer) error {

	switch index.Entity.Header.Type {
	case String:
		db.buildStringIndex(index)
	case List:

	}
	return nil
}

//
// @description: 创建字符串索引
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 16:30
// @success:
func (db *EasonDB) buildStringIndex(index *index.Indexer) {
	//now := uint32(time.Now().Unix())

	switch index.Entity.Header.OperationType {
	case StringSet:
		db.strIndex.indexList.Put(index.Entity.Key, index)
	case StringDel:
	}

}

// Run
// @description: 运行数据库
// @param: config 数据库配置
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/20 11:17
// @success:
func Run() (db *EasonDB, err error) {
	//func Run() (db *EasonDB, err error) {
	//var c Config
	//if config == c {
	//	config = DefaultConfig()
	//}
	activeFileId := uint32(0)
	activeFile := &model.DBFile{}
	var archivedFiles ArchivedFiles

	switch DbConfig.DataMode {
	case Redis:
		archivedFiles, activeFileId, err = model.LoadFile(DbConfig.DbPath)
		if err != nil {
			return nil, err
		}
		activeFile, err = model.OpenDBFile(DbConfig.DbPath, activeFileId)
	case Mysql:
		//TODO 运行模式

	}
	meta, _ := model.ReadMeta(DbConfig.DbPath)
	activeFile.Offset = meta.ActiveWriteOffset
	db = &EasonDB{
		activeFile:    activeFile,
		activeFileId:  activeFileId,
		archivedFiles: archivedFiles,
		strIndex:      NewStrIndex(),
		config:        *DbConfig,
		meta:          meta,
	}

	if err = db.readIndexFromFile(); err != nil {
		return nil, err
	}
	return db, err
}

// Close
// @description: 数据库关闭
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/20 11:25
// @success:
func (db *EasonDB) Close() (err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	err = db.meta.Store(db.config.DbPath)
	return
}
