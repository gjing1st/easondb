package index

import (
	"bytes"
	"math"
	"math/rand"
	"time"
)

const (
	maxLevel            = 18
	probability float64 = 1 / math.E
)

type (
	// Node 跳表节点
	Node struct {
		next []*Element
	}
	// Element 跳表存储元素定义
	Element struct {
		Node
		key   []byte
		value interface{}
	}
	// SkipList 跳表定义
	SkipList struct {
		Node
		maxLevel       int
		Len            int
		randSource     rand.Source
		probability    float64
		probTable      []float64
		prevNodesCache []*Node
	}
)

// Put 方法存储一个元素至跳表中，如果key已经存在，则会更新其对应的value
//因此此跳表的实现暂不支持相同的key
func (skl *SkipList) Put(key []byte, value interface{})  {
	var element *Element
	prev := skl.prevNodes(key)
	element = prev[0].next[0]
	if element != nil && bytes.Compare(element.key, key) == 0 {
		element.value = value
		return
	}
	//else if bytes.Compare(element.key, key) > 0 {
	//	return element
	//}
	element = &Element{
		Node: Node{
			next: make([]*Element, skl.randomLevel()),
		},
		key:   key,
		value: value,
	}
	for i := range element.next {
		element.next[i] = prev[i].next[i]
		prev[i].next[i] = element

	}
	skl.Len++
	return
}

//找到key对应的前一个节点索引的信息
func (skl *SkipList) prevNodes(key []byte) []*Node {
	var prev = &skl.Node
	var next *Element

	prevs := skl.prevNodesCache
	for i := skl.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && bytes.Compare(key, next.key) > 0 {
			prev = &next.Node
			next = prev.next[i]
		}
		prevs[i] = prev
	}
	return prevs
}

//生成索引随机层数
func (skl *SkipList) randomLevel() (level int) {
	r := float64(skl.randSource.Int63()) / (1 << 63)

	level = 1
	for level < skl.maxLevel && r < skl.probTable[level] {
		level++
	}
	return
}

func NewSkipList() *SkipList {
	return &SkipList{
		Node: Node{
			next: make([]*Element, maxLevel),
		},
		prevNodesCache: make([]*Node, maxLevel),
		maxLevel:       maxLevel,
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
		probability:    probability,
		probTable:      probabilityTable(probability, maxLevel),
	}
}

func probabilityTable(probability float64, maxLevel int) (table []float64) {
	for i := 1; i <= maxLevel; i++ {
		prob := math.Pow(probability, float64(i-1))
		table = append(table, prob)
	}
	return table
}

// Get 方法根据 key 查找对应的 Element 元素
//未找到则返回nil
func (skl *SkipList) Get(key []byte) *Element {
	var prev = &skl.Node
	var next *Element

	for i := skl.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && bytes.Compare(key, next.key) > 0 {
			prev = &next.Node
			next = prev.next[i]
		}
	}
	if next != nil && bytes.Compare(key, next.key) == 0 {
		return next
	}
	return nil
}

func (e *Element) Value() interface{} {
	return e.value
}
