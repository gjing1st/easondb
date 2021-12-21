package easondb

type DataType = uint16

const (
	String DataType = iota
	Hash
	List
	Set
)

const (
	StringSet uint16 = iota
	StringDel
)
