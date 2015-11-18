package skiplist

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

func TestCreateNode(t *testing.T) {
	node := NewSkipListNode(2, 1000, "abcde")
	if len(node.level) != 2 || cap(node.level) != 2 ||
		node.key != 1000 || node.value != "abcde" || node.backword != nil {
		t.Error("new node failed!")
	}
}

func TestCreateSkipList(t *testing.T) {
	list := NewSkipList()
	if list.header == nil || list.tail != nil || list.level != 1 || list.length != 0 || list.header.backword != nil {
		t.Error("create list failed")
	}
	for i := 0; i < MaxLevel; i++ {
		if list.header.level[i].forward != nil {
			t.Error("initial list header failed")
		}
	}
}

func TestRandomLevel(t *testing.T) {
	levels := []int{}
	for i := 0; i < 10000; i++ {
		l := randomLevel()
		if l > MaxLevel {
			t.Error("randomLevel error, level larger then max")
		}
		levels = append(levels, l)
	}
	sort.Ints(levels)
	if levels[0] == levels[len(levels)-1] {
		t.Error("not random level")
	}

	levelmap := map[int]int{}
	for _, l := range levels {
		levelmap[l] = levelmap[l] + 1
	}

	levelset := []int{}
	for k, _ := range levelmap {
		levelset = append(levelset, k)
	}
	sort.Ints(levelset)

	sum := 0
	statistic := map[int]int{}
	for i := len(levelset) - 1; i >= 0; i-- {
		sum += levelmap[levelset[i]]
		statistic[levelset[i]] = sum
	}

	for j, i := range levelset {
		fmt.Println(i, ":", statistic[i], "\t", float32(statistic[i])/float32(len(levels)), "\t", math.Pow(0.25, float64(j)))
	}
}

func TestInsert(t *testing.T) {
	list := NewSkipList()
	for i := 0; i < 10; i++ {
		key := 100 - uint32(i)
		value := fmt.Sprintf("test-%d", key)
		list.Insert(key, value)
	}

	x := list.header
	if x.backword != nil {
		t.Error("header backword must nil")
	}
	for x.level[0].forward != nil {
		x = x.level[0].forward
		fmt.Println(x.key, x.value)
	}

	fmt.Println("======node info======")
	fmt.Printf("key\tvalue\tbackword\tlevel\tmemory address\n")
	x = list.header
	var backwordAddr *SkipListNode
	for x != nil {
		fmt.Printf("%v\t%p\n", x, x)
		if x != list.header {
			if x.backword != backwordAddr {
				t.Error("backword pointer error")
			}
			backwordAddr = x
		} else {
			fmt.Println("---------------------------------------------------------")
		}
		x = x.level[0].forward
	}
}
