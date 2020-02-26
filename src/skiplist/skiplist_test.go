package skiplist


import (
	"fmt"
	"math/rand"
	"testing"
	"utils"
)

func Test_Insert(t *testing.T) {
	skiplist := NewSkipList(utils.IntComparator)
	for i := 0; i < 10; i++ {
		x := rand.Int() % 10
		skiplist.Insert(x)
	}


	fmt.Println("---------")

	it := skiplist.NewIterator()
	for it.SeekToFirst(); it.Valid(); it.Next() {
		fmt.Println(it.Key())
	}
	fmt.Println("---------")

	for it.SeekToLast(); it.Valid(); it.Prev() {
		fmt.Println(it.Key())
	}

}