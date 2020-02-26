package memtable

import(
	"skiplist"
	"internal_key"
)

type MemTable struct{
	table *skiplist.SkipList
	memoryUsage uint64
}


func NewMemtable() (m *MemTable){
	m = &MemTable{ skiplist.NewSkipList(internal_key.InternalKeyComparator), 0 }
	return 
}


func (memTable *MemTable) NewIterator() *Iterator {
	return &Iterator{listIter: memTable.table.NewIterator()}
}

func (memTable *MemTable) Add(seq uint64, valueType internal_key.ValueType, key, value []byte) {
	internalKey := internal_key.NewInternalKey(seq, valueType, key, value)

	memTable.memoryUsage += uint64(16 + len(key) + len(value))
	memTable.table.Insert(internalKey)
}

func (memTable *MemTable) Get(key []byte) ([]byte, error) {
	lookupKey := internal_key.LookupKey(key)

	it := memTable.table.NewIterator()
	it.Seek(lookupKey)
	if it.Valid() {
		internalKey := it.Key().(*internal_key.InternalKey)
		if internal_key.UserKeyComparator(key, internalKey.UserKey) == 0 {
			// 判断valueType
			if internalKey.Type == internal_key.TypeValue {
				return internalKey.UserValue, nil
			} else {
				return nil, internal_key.ErrDeletion
			}
		}
	}
	return nil, internal_key.ErrNotFound
}

func (memTable *MemTable) ApproximateMemoryUsage() uint64 {
	return memTable.memoryUsage
}
