package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DBMeta struct {
	ActiveWriteOffset int64 `json:"active_write_offset"`
}

const (
	dbMetaFile = "db.meta" //存放数据库活跃文件偏移量
)

// Store
// @description: 存储当前活跃文件偏移量
// @param: path 数据库所在目录
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 16:13
// @success:
func (m *DBMeta) Store(path string) error {
	filePath := filepath.Join(path, dbMetaFile)
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, b, FilePerm)
	return err
}

// ReadMeta
// @description: 读取当前文件偏移量
// @param: path 数据库所在目录
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 16:13
// @success:
func ReadMeta(path string) (*DBMeta, error) {
	filePath := filepath.Join(path, dbMetaFile)
	m := &DBMeta{}
	file, err := os.OpenFile(filePath, os.O_RDONLY, FilePerm)
	if err != nil {
		return m, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return m, err
	}
	json.Unmarshal(b, &m)
	return m, err
}
