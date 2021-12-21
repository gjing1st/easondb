package model

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

var (
	ErrEntitySize = errors.New("easondb.model.entity:实体类型长度错误")
	ErrBytesSize  = errors.New("easondb.model.entity:字节长度错误")
	ErrInvalidCrc = errors.New("easondb.model.entry: 验证crc32失败")
)

type Entity struct {
	Header *Header
	Key    []byte
	Value  []byte
}

type Header struct {
	KeySize       uint32 `json:"key_size"`
	ValueSize     uint32 `json:"value_size"`
	Type          uint16
	OperationType uint16
	crc32         uint32
	ExpireTime    uint32
}

const (
	// 实体头信息字节长度 KeySize，ValueSize，crc32 各占4字节。Type和OperationType各占2个字节
	entityHeaderSize = 16
)

// NewEntity
// @description: 实例化实体
// @param: key 字节
// @param: value 字节
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 14:17
// @success:
func NewEntity(key, value []byte, ty, operationType uint16) *Entity {
	return &Entity{
		Header: &Header{
			uint32(len(key)),
			uint32(len(value)),
			ty,
			operationType,
			0,
			0,
		},
		Key:   key,
		Value: value,
	}
}

// Size
// @description: 实体的字节长度
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 15:01
// @success:
func (e *Entity) Size() uint32 {
	return entityHeaderSize + e.Header.KeySize + e.Header.ValueSize
}

// EntityToBytes
// @description: 实体转字节，用于持久化
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 15:00
// @success:
func (e *Entity) EntityToBytes() ([]byte, error) {
	if e.Size() <= entityHeaderSize {
		return nil, ErrEntitySize
	}
	bytes := make([]byte, e.Size())
	//头信息填充
	binary.BigEndian.PutUint32(bytes[:4], e.Header.KeySize)
	binary.BigEndian.PutUint32(bytes[4:8], e.Header.ValueSize)
	binary.BigEndian.PutUint16(bytes[8:10], e.Header.Type)
	binary.BigEndian.PutUint16(bytes[10:12], e.Header.OperationType)
	crc := crc32.ChecksumIEEE(e.Value)
	binary.BigEndian.PutUint32(bytes[12:16], crc)
	//k-v
	copy(bytes[entityHeaderSize:len(e.Key)+entityHeaderSize], e.Key)
	copy(bytes[len(e.Key)+entityHeaderSize:], e.Value)
	return bytes, nil
}

// BytesToEntity
// @description: 字节转实体，用于读取文件到内存
// @param: b 数据库文件中读取的字节
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 15:00
// @success:
func BytesToEntity(b []byte) (*Entity, error) {
	if len(b) < entityHeaderSize {
		return nil, ErrBytesSize
	}
	ks := binary.BigEndian.Uint32(b[:4])
	vs := binary.BigEndian.Uint32(b[4:8])
	ty := binary.BigEndian.Uint16(b[8:10])
	opty := binary.BigEndian.Uint16(b[10:12])
	crc := binary.BigEndian.Uint32(b[12:16])
	key := make([]byte, ks)
	if uint32(len(b)) > entityHeaderSize+ks {
		key = b[entityHeaderSize : entityHeaderSize+ks]
	}
	value := make([]byte, vs)
	if uint32(len(b)) > entityHeaderSize+ks+vs {
		value = b[entityHeaderSize+ks:]
	}
	//key := b[entityHeaderSize : entityHeaderSize+ks]
	//value := b[entityHeaderSize+ks:]
	return &Entity{
		Header: &Header{
			ks,
			vs,
			ty,
			opty,
			crc,
			0,
		},
		Key:   key,
		Value: value,
	}, nil
}
