package db

import (
	"sync"

	"io/ioutil"
	"os"
	"strconv"
	"fmt"

	"internal_key"
	"memtable"
	"version"
)

type Db struct {
	path                  string
	mu                    sync.Mutex
	cond                  *sync.Cond
	mem                   *memtable.MemTable
	imm                   *memtable.MemTable
	current               *version.Version
	bgCompactionScheduled bool
}

func Open(path string) *Db {
	var db Db
	db.path = path
	db.mem = memtable.NewMemtable()
	db.imm = nil
	db.bgCompactionScheduled = false
	db.cond = sync.NewCond(&db.mu)
	num := db.ReadCurrentFile()
	if num > 0 {
		v, err := version.LoadVersion(path, num)
		if err != nil {
			return nil
		}
		db.current = v
	} else {
		db.current = version.NewVersion(path)
	}

	return &db
}

func (db *Db) Close() {
	db.mu.Lock()
	for db.bgCompactionScheduled {
		db.cond.Wait()
	}
	db.mu.Unlock()
	db.imm = db.mem
	db.mem = memtable.NewMemtable()
	db.bgCompactionScheduled = true 
	db.scheduleCompaction(true)
}

func (db *Db) Put(key, value []byte) error {
	// May temporarily unlock and wait.
	seq := db.makeRoomForWrite()
	db.mem.Add(seq, internal_key.TypeValue, key, value)
	return nil
}

func (db *Db) Delete(key []byte) error {
	seq := db.makeRoomForWrite()
	db.mem.Add(seq, internal_key.TypeDeletion, key, nil)
	return nil
}

func (db *Db) Get(key []byte) ([]byte, error) {
	db.mu.Lock()
	mem := db.mem
	imm := db.mem
	current := db.current
	db.mu.Unlock()
	value, err := mem.Get(key)
	if err != internal_key.ErrNotFound {
		return value, err
	}

	if imm != nil {
		value, err := imm.Get(key)
		if err != internal_key.ErrNotFound {
			return value, err
		}
	}

	value, err = current.Get(key)
	return value, err
}



func (db *Db) makeRoomForWrite() uint64 {
	// 这里可以多个同时进入 并不会造成同步
	// 详情见 cond的作用
	db.mu.Lock()
	defer db.mu.Unlock()

	for ;; {
		if db.mem.ApproximateMemoryUsage() <= internal_key.Write_buffer_size {
			return db.current.NextSeq()
		} else if db.imm != nil {
			//  Current memtable full; waiting
			db.cond.Wait()
		} else {
			// Attempt to switch to a new memtable and trigger compaction of old
			// todo : switch log
			db.imm = db.mem
			db.mem = memtable.NewMemtable()
			db.bgCompactionScheduled = true 
			go db.scheduleCompaction(false) 
		}
	}

	return 0
}





func (db *Db) scheduleCompaction(justmajor bool) {
	imm := db.imm
	version := db.current.Copy()
	// minor compaction
	if imm != nil {
		version.WriteLevel0Table(imm)
	}
	// major compaction
	for version.DoCompactionWork() && !justmajor{
		version.Log()
	}

	descriptorNumber, _ := version.Save()
	db.SetCurrentFile(descriptorNumber)
	db.imm = nil
	db.current = version
	db.bgCompactionScheduled = false
	db.cond.Broadcast()
}

func (db *Db) SetCurrentFile(descriptorNumber uint64) {
	tmp := internal_key.TempFileName(db.path, descriptorNumber)
	ioutil.WriteFile(tmp, []byte(fmt.Sprintf("%d", descriptorNumber)), 0600)
	os.Rename(tmp, internal_key.CurrentFileName(db.path))
}





func (db *Db) ReadCurrentFile() uint64 {
	b, err := ioutil.ReadFile(internal_key.CurrentFileName(db.path))
	if err != nil {
		return 0
	}
	descriptorNumber, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0
	}
	return descriptorNumber
}
