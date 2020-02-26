package memtable

import (
	"fmt"
	"testing"
	"internal_key"
)

func Test_MemTable(t *testing.T) {
	memTable := NewMemtable()
	memTable.Add(1234567, internal_key.TypeValue, []byte("aadsa34a"), []byte("bb23b3423"))
	value, _ := memTable.Get([]byte("aadsa34a"))
	fmt.Println(string(value))
	fmt.Println(memTable.ApproximateMemoryUsage())
}
