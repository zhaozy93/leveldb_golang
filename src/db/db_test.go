package db

import (
	"fmt"
	"testing"
)

func Test_DB(test *testing.T){
	db := Open("/Users/zyzhao/Learn/leveldb/testDB")
	fmt.Println("存入 123,  456")
	db.Put([]byte("123"), []byte("456"))
	value, err := db.Get([]byte("123"))
	fmt.Println("get 123")
	fmt.Println(string(value), err)
	fmt.Println("get 1234")
	value, err = db.Get([]byte("1234"))
	fmt.Println(string(value), err)

	fmt.Println("删除 123")
	db.Delete([]byte("123"))
	value, err = db.Get([]byte("123"))
	fmt.Println("get 123")
	fmt.Println(string(value), err)
	fmt.Println("db close.....")
	db.Close()
	fmt.Println("db reopen.....")
	db = Open("/Users/zyzhao/Learn/leveldb/testDB")
	fmt.Println("get 123")
	value, err = db.Get([]byte("123"))
	fmt.Println(string(value), err)
}
