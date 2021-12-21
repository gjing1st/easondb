package model

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	FilePerm         = 0644
	DBFileFormatName = "%09d.db"
	RedisModeDir     = "redis"
	TableModeDir     = "table"
)

var (
	ErrFileName = errors.New("easondb.model.file:数据文件名称有误(不是整型)")
)

type DBFile struct {
	Id     uint32
	Path   string
	File   *os.File
	Offset int64
}

// OpenDBFile
// @description: 打开数据文件
// @param: path 数据库文件所在目录路径
// @param: fileId 数据文件id
// @param: blockSize 数据大小
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 15:35
// @success:
func OpenDBFile(path string, fileId uint32) (*DBFile, error) {
	filePath := filepath.Join(path, fmt.Sprintf(DBFileFormatName, fileId))
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, FilePerm)
	if err != nil {
		return nil, err
	}
	return &DBFile{
		fileId,
		filePath,
		file,
		0,
	}, nil
}

// LoadFile
// @description: 加载目录下所有数据文件
// @param: dbDir 数据文件所在目录
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/20 14:12
// @success:
func LoadFile(dbDir string) (files map[uint32]*DBFile, activeFileId uint32, err error) {
	fileList, err := ioutil.ReadDir(dbDir)
	if err != nil {
		return nil, 0, err
	}
	var fileIds []int
	for _, file := range fileList {
		fileNameFull := file.Name()           //文件名
		fileExt := filepath.Ext(fileNameFull) //文件后缀
		if fileExt == ".db" {
			fileName := fileNameFull[0 : len(fileNameFull)-len(fileExt)] //文件前缀
			fileId, err := strconv.Atoi(fileName)                        //转整型
			if err != nil {
				return nil, 0, ErrFileName
			}
			fileIds = append(fileIds, fileId)
		}
	}
	sort.Ints(fileIds)
	if len(fileIds) > 0 {
		files = make(map[uint32]*DBFile)
		activeFileId = uint32(fileIds[len(fileIds)-1])
		for i := 0; i < len(fileIds)-1; i++ {
			file, err := OpenDBFile(dbDir, uint32(fileIds[i]))
			if err != nil {
				return nil, 0, err
			}
			files[uint32(i)] = file
		}
	}
	return
}

// Close
// @description: 关闭数据库文件
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 15:24
// @success:
func (df *DBFile) Close() (err error) {
	if df.File != nil {
		err = df.File.Close()
	}
	return
}

// Write
// @description: 保存到文件相当于redis的aof
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 15:48
// @success:
func (df *DBFile) Write(e *Entity) error {
	offset := df.Offset
	data, err := e.EntityToBytes()
	if err != nil {
		return err
	}
	if _, err := df.File.WriteAt(data, offset); err != nil {
		return err
	}
	df.Offset += int64(e.Size())
	return nil
}

// Read
// @description: 读取数据文件封装实体
// @param: offset 文件偏移量
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 16:08
// @success:
func (df *DBFile) Read(offset int64) (e *Entity, err error) {
	var buf []byte
	if buf, err = df.ReadBuf(offset, entityHeaderSize); err != nil {
		return
	}
	if e, err = BytesToEntity(buf); err != nil {
		return
	}
	offset += entityHeaderSize
	if e.Header.KeySize > 0 {
		if e.Key, err = df.ReadBuf(offset, int64(e.Header.KeySize)); err != nil {
			return
		}
	}
	offset += int64(e.Header.KeySize)
	if e.Header.ValueSize > 0 {
		if e.Value, err = df.ReadBuf(offset, int64(e.Header.ValueSize)); err != nil {
			return
		}
	}
	if len(e.Value) > 0 {
		crc := crc32.ChecksumIEEE(e.Value)
		if crc != e.Header.crc32 {
			err = ErrInvalidCrc
			return
		}
	}

	return
}

// ReadBuf
// @description: 读取数据文件指定长度字节
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 15:54
// @success:
func (df *DBFile) ReadBuf(offset, l int64) ([]byte, error) {
	buf := make([]byte, l)
	_, err := df.File.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
