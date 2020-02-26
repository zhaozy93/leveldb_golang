package sstable

import (
	"fmt"

	"testing"

	"internal_key"
)

func Test_Table(t *testing.T) {
	builder := NewTableBuilder("./000123.ldb")
	item := internal_key.NewInternalKey(1, internal_key.TypeValue, []byte("123"), []byte("1234"))
	builder.Add(item)
	item = internal_key.NewInternalKey(2, internal_key.TypeValue, []byte("124"), []byte("1245"))
	builder.Add(item)
	item = internal_key.NewInternalKey(3, internal_key.TypeValue, []byte("125"), []byte("0245"))
	builder.Add(item)
	builder.Finish()

	fmt.Println("-----")
	table, err := SsTabelOpen("./000123.ldb")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(table.index)
	fmt.Println(table.footer)
	it := table.NewIterator()

	it.SeekToFirst()
	for ;it.Valid();{
		fmt.Println(it.InternalKey())
		it.Next()
	}

	it.Seek([]byte("1244"))
	if it.Valid() {
		if string(it.InternalKey().UserKey) != "125" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
