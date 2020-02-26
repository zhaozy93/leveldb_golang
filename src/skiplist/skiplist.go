package skiplist


import (
	"math/rand"
	"sync"
	"utils"
)


const (
	kMaxHeight = 12
	kBranching = 4
)

type SkipList struct {
	maxHeight  int 
	head       *Node   
	comparator utils.Comparator
	locker     sync.RWMutex
}

func NewSkipList(comp utils.Comparator) (list *SkipList){
	list = &SkipList{}
	list.head = newNode(nil, kMaxHeight)
	list.maxHeight = 0
	list.comparator = comp
	return
}


func (list *SkipList)randomHeight() (height int){
	for height < kMaxHeight && (rand.Intn(kBranching) == 0){
		height++
	}
	return
}

func (list *SkipList) Insert(key interface{}){
	list.locker.Lock()
	defer list.locker.Unlock()
	_, prev := list.findGreaterOrEqual(key)
	height := list.randomHeight()
	if height > list.maxHeight {
		for i := list.maxHeight+1; i <= height; i++ {
			prev[i] = list.head
		}
		list.maxHeight = height
	}
	x := newNode(key, height)
	for i := 0; i <= height; i++ {
		x.setNext(i, prev[i].getNext(i))
		prev[i].setNext(i, x)
	}
}


func (list *SkipList) Contains(key interface{}) bool {
	list.locker.RLock()
	defer list.locker.RUnlock()
	x, _ := list.findGreaterOrEqual(key)
	if x != nil && list.comparator(x.key, key) == 0 {
		return true
	}
	return false
}


func (list *SkipList) findGreaterOrEqual(key interface{}) (*Node, [kMaxHeight]*Node) {
	var prev [kMaxHeight]*Node
	x := list.head
	level := list.maxHeight
	for ;; {
		next := x.getNext(level)
		if list.keyIsAfterNode(key, next) {
			x = next
		} else {
			prev[level] = x
			if level == 0 {
				return next, prev
			} else {
				// Switch to next list
				level--
			}
		}
	}
	return nil, prev
}

func (list *SkipList) keyIsAfterNode(key interface{}, n *Node) bool {
	return (n != nil) && (list.comparator(n.key, key) < 0)
}


func (list *SkipList) NewIterator() *Iterator {
	var it Iterator
	it.list = list
	return &it
}

func (list *SkipList) findLessThan(key interface{}) *Node {
	x := list.head
	level := list.maxHeight
	for ;; {
		next := x.getNext(level)
		if next == nil || list.comparator(next.key, key) >= 0 {
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = next
		}
	}
	return nil
}

func (list *SkipList) findlast() *Node {
	x := list.head
	level := list.maxHeight
	for ;; {
		next := x.getNext(level)
		if next == nil {
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = next
		}
	}
	return nil
}


