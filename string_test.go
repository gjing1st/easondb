package easondb

import (
	"fmt"
	"log"
	"testing"
)

var dbPath = "E:\\db_test\\easondb"

func InitDb() *EasonDB {
	config := DefaultConfig()
	config.DbPath = dbPath
	config.BlockSize =  1024

	db, err := Run(config)
	//db, err := Run(Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestEasonDB_Set(t *testing.T) {
	db := InitDb()
	defer db.Close()
	fmt.Println("db.offset",db.activeFile.Offset)

	//db.Set([]byte("123"),[]byte("456"))
	//db.Set([]byte("我是"), []byte("中国人1"))
	//db.Set("你好", "世界2")
	//db.Set("hello", "easondb")
	db.Set("33","44")
	db.Set("斗罗","大陆")
	//if err != nil {
	//	log.Fatal("write data error ", err)
	//}
	s := db.Get("123")
	s1 := db.Get("我是")
	s2 := db.Get("你好")
	s3 := db.Get("hello")
	s4 := db.Get("斗罗")
	s5 := db.Get("33")

	fmt.Println("-------------", string(s))
	fmt.Println("-------------", string(s1))
	fmt.Println("-------------", string(s2))
	fmt.Println("-------------", string(s3))
	fmt.Println("-------------", string(s4),"--",string(s5))
	fmt.Println("db.offset",db.activeFile.Offset)

}
