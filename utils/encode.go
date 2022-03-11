package utils

import (
	"github.com/gogf/gf/util/gconv"
)

// Encode
// @description: key和value转字节
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 11:11
// @success:
func Encode(key,value interface{})(keyByte,valueByte []byte){
	keyByte = gconv.Bytes(key)
	valueByte = gconv.Bytes(value)
	return
}

// EncodeKey
// @description: 将key转为字节
// @param:
// @author: GJing
// @email: gjing1st@gmail.com
// @date: 2021/12/17 11:10
// @success:
func EncodeKey(key interface{})(keyByte []byte){
	keyByte = gconv.Bytes(key)
	return
}