package block

import (
	"fmt"
	"testing"

	"internal_key"
)

func Test_Block(t *testing.T) {
	var builder BlockBuilder

	item := internal_key.NewInternalKey(1, internal_key.TypeValue, []byte("123"), []byte("1234"))
	builder.Add(item)
	item = internal_key.NewInternalKey(2, internal_key.TypeValue, []byte("124"), []byte("1245"))
	builder.Add(item)
	item = internal_key.NewInternalKey(3, internal_key.TypeValue, []byte("125"), []byte("0245"))
	builder.Add(item)
	p := builder.Finish()

	block := NewBlock(p)
	it := block.NewIterator()

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
