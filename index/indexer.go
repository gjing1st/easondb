package index

import (
	"github.com/gjing1st/easondb/model"
)

type Indexer struct {
	Entity     *model.Entity
	FileId     uint32
	EntitySize uint32
	Offset     int64
}
