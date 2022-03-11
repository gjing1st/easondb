package easondb

import (
	"gitee.com/gjing1st/easondb/index"
	"gitee.com/gjing1st/easondb/model"
	"io"
)

//
// @description: 从数据文件读取索引
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 16:34
// @success:
func (db *EasonDB) readIndexFromFile() (err error) {
	//归档文件
	if db.archivedFiles != nil {
		for _, file := range db.archivedFiles {
			err = db.readFile(file)
			if err != nil {
				return
			}
		}
	}
	//活跃文件
	err = db.readFile(db.activeFile)
	return err
}

// readFile
// @description: 读取数据文件
// @param: df 数据文件
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/20 15:34
func (db *EasonDB) readFile(df *model.DBFile) error {
	var offset int64
	for offset <= db.config.BlockSize {
		if e, err := df.Read(offset); err == nil {
			index := &index.Indexer{
				Entity:     e,
				FileId:     0,
				EntitySize: e.Size(),
				Offset:     offset,
			}
			offset += int64(e.Size())
			if err = db.buildIndex(index); err != nil {
				return err
			}
		} else {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}
