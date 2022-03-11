package easondb

import (
	"gitee.com/gjing1st/easondb/index"
	"gitee.com/gjing1st/easondb/model"
	"gitee.com/gjing1st/easondb/utils"
	"sync"
)

type StrIndex struct {
	mutex     sync.RWMutex
	indexList *index.SkipList
}

// Set
// @description: 字符串存储
// @param: key
// @param: value
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 11:11
// @success:
func (db *EasonDB) Set(key, value interface{}) (err error) {
	keyByte, valueByte := utils.Encode(key, value)
	if err = db.validKeyValue(keyByte, valueByte); err != nil {
		return
	}
	if err = db.set(keyByte, valueByte); err != nil {
		return
	}

	return
}
// set
func (db *EasonDB) set(key, value []byte) (err error) {
	db.strIndex.mutex.Lock()
	defer db.strIndex.mutex.Unlock()
	//实体
	entity := model.NewEntity(key, value, String, StringSet)
	//数据持久化
	if err = db.store(entity); err != nil {
		return
	}
	index := &index.Indexer{
		FileId:     db.activeFileId,
		EntitySize: entity.Size(),
		Offset:     db.activeFile.Offset - int64(entity.Size()),
		Entity: entity,
	}
	//创建索引
	err = db.buildIndex(index)

	return
}

func NewStrIndex() *StrIndex {
	return &StrIndex{
		indexList: index.NewSkipList(),
	}
}

// Get
// @description: 查找key对应的值
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 16:35
// @success:
func (db *EasonDB) Get(key interface{}) []byte {
	keyByte := utils.EncodeKey(key)
	node := db.strIndex.indexList.Get(keyByte)
	if node == nil {
		return nil
	}
	index := node.Value().(*index.Indexer)

	return index.Entity.Value
}
