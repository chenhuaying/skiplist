package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"os"
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
	r := rand.New(rand.NewSource(int64(os.Getpid())))
	list := NewSkipList()
	for i := 0; i < 10; i++ {
		//key := 100 - uint32(i)
		key := r.Uint32()
		value := fmt.Sprintf("test-%d", key)
		list.Insert(key, value)
	}

	x := list.header
	if x.backword != nil {
		t.Error("header backword must nil")
	}

	if list.tail == nil || list.tail == list.header {
		t.Error("wrong list tail")
	}

	fmt.Println("======node info======")
	fmt.Printf("key\tvalue\tbackword\tlevel\tmemory address\n")
	x = list.header
	seq := 0
	var backwordAddr *SkipListNode
	for x != nil {
		fmt.Printf("%d\t%v\t%p\n", seq, x, x)
		if x != list.header {
			if x.backword != backwordAddr {
				t.Error("backword pointer error")
			}
			backwordAddr = x
		} else {
			fmt.Println("---------------------------------------------------------")
		}
		x = x.level[0].forward
		seq += 1
	}
}

func TestSearch(t *testing.T) {
	r := rand.New(rand.NewSource(int64(os.Getpid())))
	list := NewSkipList()
	for i := 0; i < 10; i++ {
		//key := 100 - uint32(i)
		key := r.Uint32()
		value := fmt.Sprintf("test-%d", key)
		list.Insert(key, value)
	}

	list.Insert(100, "test-100")
	list.Insert(307278502, "test-307278502")
	list.Insert(2654365283, "test-2654365283")

	value1 := list.Search(100)
	if value1.(string) != "test-100" {
		t.Error("search 100 failed")
	}

	value2 := list.Search(307278502)
	if value2.(string) != "test-307278502" {
		t.Error("search 307278502 failed")
	}

	value3 := list.Search(2654365283)
	if value3.(string) != "test-2654365283" {
		t.Error("search 2654365283 failed")
	}

	value4 := list.Search(9000)
	// this may not be true, random key may right to be 9000
	if value4 != nil {
		t.Error("find non exist key, is right?")
	}

	value5 := list.Search(0)
	if value5 != nil {
		t.Error("find non exist key, is right?")
	}
}

func initRandomTestList() *SkipList {
	r := rand.New(rand.NewSource(int64(os.Getpid())))
	list := NewSkipList()
	for i := 0; i < 10; i++ {
		//key := 100 - uint32(i)
		key := r.Uint32()
		value := fmt.Sprintf("test-%d", key)
		list.Insert(key, value)
	}

	list.Insert(100, "test-100")
	list.Insert(307278502, "test-307278502")
	list.Insert(2654365283, "test-2654365283")
	return list
}

func initTestList() *SkipList {
	list := NewSkipList()
	for i := 0; i < 10; i++ {
		key := 10 + uint32(i)
		value := fmt.Sprintf("test-%d", key)
		list.Insert(key, value)
	}

	return list
}

func showSkipList(list *SkipList, t *testing.T) {
	fmt.Println("======list info======")
	fmt.Printf("%p %p %d %d\n", list.header, list.tail, list.level, list.length)
	fmt.Println("======node info======")
	fmt.Printf("key\tvalue\tbackword\tlevel\tmemory address\n")
	x := list.header
	seq := 0
	var backwordAddr *SkipListNode
	for x != nil {
		fmt.Printf("%d\t%v\t%p\n", seq, x, x)
		if x != list.header {
			if x.backword != backwordAddr {
				t.Error("backword pointer error")
			}
			backwordAddr = x
		} else {
			fmt.Println("---------------------------------------------------------")
		}
		x = x.level[0].forward
		seq += 1
	}
}

func TestDelete(t *testing.T) {
	list := initRandomTestList()
	fmt.Println("##########################before delete#########################")
	showSkipList(list, t)
	fmt.Println("##########################before delete#########################")
	node1 := list.DeleteNode(100)
	fmt.Printf("delte node: \t%v\t%p\n", node1, node1)
	if node1.value != "test-100" {
		t.Error("delete wrong node")
	}

	node2 := list.DeleteNode(307278502)
	fmt.Printf("delete node: \t%v\t%p\n", node2, node2)
	if node2.value != "test-307278502" {
		t.Error("delete wrong node")
	}

	node3 := list.DeleteNode(2654365283)
	fmt.Printf("delete node: \t%v\t%p\n", node3, node3)
	if node3.value != "test-2654365283" {
		t.Error("delete wrong node")
	}

	// this case will be fail, when random key is the MaxInt32
	node4 := list.DeleteNode(math.MaxInt32)
	if node4 != nil {
		t.Error("delete node not in the list, is right?")
	}

	// this case will be fail, when random key is the 10
	node5 := list.DeleteNode(10)
	if node5 != nil {
		t.Error("delete node not in the list, is right?")
	}

	node6 := list.tail
	node7 := list.DeleteNode(node6.key)
	if node6 != node7 {
		t.Error("delete wrong tail node")
	}

	showSkipList(list, t)
}

func TestDeleteAllNode(t *testing.T) {
	list := initTestList()
	x := list.header.level[0].forward
	for x != nil {
		list.DeleteNode(x.key)
		x = x.level[0].forward
	}

	if list.level != 1 || list.tail != nil || list.length != 0 {
		t.Error("delete all node failed")
	}
	nonNode := list.DeleteNode(0)
	if nonNode != nil {
		t.Error("delete empty node failed")
	}
	showSkipList(list, t)

	list2 := initTestList()
	x = list2.tail
	for x != nil {
		list2.DeleteNode(x.key)
		x = x.backword
	}
	if list2.level != 1 || list2.tail != nil || list2.length != 0 {
		t.Error("delete all node failed")
	}
	nonNode = list2.DeleteNode(0)
	if nonNode != nil {
		t.Error("delete empty node failed")
	}
	showSkipList(list2, t)
}

func TestLowerBoundNode(t *testing.T) {
	list := initTestList()
	list.Insert(3, "test-3")
	list.Insert(25, "test-25")
	list.Insert(35, "test-35")

	node1 := list.LowerBoundNode(0)
	if node1.key != 3 {
		t.Error("lower bound before first failed")
	}

	node2 := list.LowerBoundNode(35)
	if node2.key != 35 {
		t.Error("LowerBoundNode last failed")
	}

	node3 := list.LowerBoundNode(40)
	if node3 != nil {
		t.Error("LowerBoundNode after last failed")
	}

	node4 := list.LowerBoundNode(3)
	if node4.key != 3 {
		t.Error("LowerBoundNode first failed")
	}

	node6 := list.LowerBoundNode(12)
	if node6.key != 12 {
		t.Error("LowerBoundNode 12 failed")
	}

	node7 := list.LowerBoundNode(20)
	if node7.key != 25 {
		t.Error("LowerBoundNode 20 failed")
	}
}
