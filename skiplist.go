package skiplist

import (
	"math/rand"
	"os"
)

type SkipListLevel struct {
	forward *SkipListNode
}

type SkipListNode struct {
	key      uint32
	value    interface{}
	backword *SkipListNode
	level    []SkipListLevel
}

type SkipList struct {
	header *SkipListNode
	tail   *SkipListNode
	length int
	level  int
}

const (
	MaxLevel  = 32
	SkipListP = 0.25
)

var levelRand *rand.Rand

func init() {
	levelRand = rand.New(rand.NewSource(int64(os.Getpid())))
}

func NewSkipListNode(level int, key uint32, value interface{}) *SkipListNode {
	node := &SkipListNode{key: key, value: value, level: make([]SkipListLevel, level)}
	return node
}

func NewSkipList() *SkipList {
	header := NewSkipListNode(MaxLevel, 0, nil)
	skiplist := &SkipList{header: header, tail: nil, length: 0, level: 1}
	if skiplist != nil {
		for i := 0; i < MaxLevel; i++ {
			skiplist.header.level[i].forward = nil
		}
	}
	return skiplist
}

func randomLevel() int {
	level := 1

	for levelRand.Float32() < 0.25 && level < MaxLevel {
		level += 1
	}
	return level
}

func (l *SkipList) Insert(key uint32, value interface{}) {
	update := [MaxLevel]*SkipListNode{}
	x := l.header
	for i := l.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.key < key {
			x = x.level[i].forward
		}
		update[i] = x
	}

	x = x.level[0].forward
	if x != nil && x.key == key {
		x.value = value
	} else {
		level := randomLevel()
		if level > l.level {
			for i := l.level; i < level; i++ {
				update[i] = l.header
			}
			// update skiplist max level
			l.level = level
		}

		x = NewSkipListNode(level, key, value)
		for i := 0; i < level; i++ {
			x.level[i].forward = update[i].level[i].forward
			update[i].level[i].forward = x
		}
	}

	if update[0] == l.header {
		x.backword = nil
	} else {
		x.backword = update[0]
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backword = x
	} else {
		// insert to the last position of the list
		l.tail = x
	}

	l.length += 1
}

func (l *SkipList) SearchNode(key uint32) *SkipListNode {
	x := l.header
	for i := l.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.key < key {
			x = x.level[i].forward
		}
	}

	x = x.level[0].forward
	if x != nil && x.key == key {
		return x
	}

	return nil
}

func (l *SkipList) Search(key uint32) interface{} {
	node := l.SearchNode(key)
	if node != nil {
		return node.value
	}
	return nil
}

func (l *SkipList) LowerBoundNode(key uint32) *SkipListNode {
	x := l.header
	for i := l.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.key < key {
			x = x.level[i].forward
		}
	}
	x = x.level[0].forward
	return x
}

func (l *SkipList) DeleteNode(key uint32) *SkipListNode {
	update := [MaxLevel]*SkipListNode{}
	x := l.header
	for i := l.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.key < key {
			x = x.level[i].forward
		}
		update[i] = x
	}
	x = x.level[0].forward
	if x != nil && x.key == key {
		for i := 0; i < l.level; i++ {
			if update[i].level[i].forward != x {
				break
			} else {
				update[i].level[i].forward = x.level[i].forward
			}
		}

		if x.level[0].forward == nil {
			l.tail = x.backword
		} else {
			x.level[0].forward.backword = x.backword
		}

		// if x is the hightest node
		for l.level > 1 && l.header.level[l.level-1].forward == nil {
			l.level -= 1
		}
		l.length -= 1

		return x
	} else {
		return nil
	}
}
