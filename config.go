package easondb

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// DataMode 数据库模式
type DataMode uint8

const (
	Redis DataMode = iota //redis模式
	Mysql                 //mysql表模式
)

const (
	// DefaultDirPath 默认数据库目录
	DefaultDbPath       = "E:\\db_test\\easondb"
	DefaultBlockSize    = 100 * 1024 * 1024
	DefaultMaxKeySize   = uint32(1024 * 1024)
	DefaultMaxValueSize = uint32(20 * 1024 * 1024)
)

type Config struct {
	DbPath       string   `json:"db_path"`
	DataMode     DataMode `json:"data_mode"`
	BlockSize    int64    `json:"block_size"`
	MaxKeySize   uint32   `json:"max_key_size"`
	MaxValueSize uint32   `json:"max_value_size"`
	Sync         bool
}

var DbConfig *Config

func init() {
	DbConfig = &Config{
		DbPath:       DefaultDbPath,
		DataMode:     Redis,
		BlockSize:    DefaultBlockSize,
		MaxKeySize:   DefaultMaxKeySize,
		MaxValueSize: DefaultMaxValueSize,
		Sync:         false,
	}
	DbConfig.Reload()
}

func (c *Config) Reload() {
	data, err := ioutil.ReadFile("config/db.json")
	if err != nil {
		log.Println("使用默认配置")
		return
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Println("使用默认配置")
	}
}

func (c *Config) SetDbPath(path string) {
	c.DbPath = path
}
func DefaultConfig() Config {
	return Config{
		DbPath:       DefaultDbPath,
		DataMode:     Redis,
		BlockSize:    DefaultBlockSize,
		MaxKeySize:   DefaultMaxKeySize,
		MaxValueSize: DefaultMaxValueSize,
		Sync:         false,
	}
}
