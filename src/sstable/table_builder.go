package sstable

import (
	"os"
	"internal_key"
	"sstable/block"
)

const (
	MAX_BLOCK_SIZE = 4 * 1024
)

type TableBuilder struct {
	file               *os.File
	offset             uint32
	numEntries         int32
	dataBlockBuilder   block.BlockBuilder
	indexBlockBuilder  block.BlockBuilder
	pendingIndexHandle IndexBlockHandle
	status             error
}

func NewTableBuilder(fileName string) *TableBuilder {
	var builder TableBuilder
	var err error
	builder.file, err = os.Create(fileName)
	if err != nil {
		return nil
	}
	return &builder
}

func (builder *TableBuilder) FileSize() uint32 {
	return builder.offset
}

func (builder *TableBuilder) Add(internalKey *internal_key.InternalKey) {
	if builder.status != nil {
		return
	}
	// todo : filter block

	builder.pendingIndexHandle.InternalKey = internalKey

	builder.numEntries++
	builder.dataBlockBuilder.Add(internalKey)
	if builder.dataBlockBuilder.CurrentSizeEstimate() > MAX_BLOCK_SIZE {
		builder.flush()
	}
}
func (builder *TableBuilder) flush() {
	if builder.dataBlockBuilder.Empty() {
		return
	}

	orgKey := builder.pendingIndexHandle.InternalKey
	builder.pendingIndexHandle.InternalKey = internal_key.NewInternalKey(orgKey.Seq, orgKey.Type, orgKey.UserKey, nil)
	builder.pendingIndexHandle.SetBlockHandle(builder.writeblock(&builder.dataBlockBuilder))
	builder.indexBlockBuilder.Add(builder.pendingIndexHandle.InternalKey)
}

func (builder *TableBuilder) Finish() error {
	// write data block
	builder.flush()
	// todo : filter block

	// write index block
	var footer Footer
	footer.IndexHandle = builder.writeblock(&builder.indexBlockBuilder)
	// write footer block
	footer.EncodeTo(builder.file)
	builder.file.Close()
	return nil
}

func (builder *TableBuilder) writeblock(blockBuilder *block.BlockBuilder) BlockHandle {
	content := blockBuilder.Finish()
	// todo : compress, crc
	var blockHandle BlockHandle
	blockHandle.Offset = builder.offset
	blockHandle.Size = uint32(len(content))
	builder.offset += uint32(len(content))
	_, builder.status = builder.file.Write(content)
	builder.file.Sync()
	blockBuilder.Reset()
	return blockHandle
}
