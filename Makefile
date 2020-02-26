# Go related variables.
GOBASE := $(shell pwd)
export GOPATH := $(GOBASE)
GOBIN := $(GOBASE)/bin
GOPKG := $(GOBASE)/pkg
GOLOG := $(GOBASE)/log
SRCDIR:= $(GOBASE)/src
APPDIR := $(GOBASE)/src/db
Vendor := $(GOBASE)/src/vendor
APPNAME := leveldb


start: clean fmt build 
	@echo " > starting...."
	$(GOBIN)/$(APPNAME) 

build: fmt 
	@echo " > building...."
	go build -o $(GOBIN)/$(APPNAME) $(APPDIR)

fmt: 
	@echo " > fmt...."
	go fmt $(APPDIR)

clean:
	@echo " > clean...."
	rm -rf $(GOBIN)
	rm -rf $(GOPKG)
	rm -rf $(GOLOG)/*


test:
	go test -v src/skiplist/*.go
	go test -v src/memtable/*.go
	go test -v src/sstable/block/*.go
	rm -rf src/sstable/000123.ldb
	go test -v src/sstable/*.go
	rm -rf /Users/zyzhao/Learn/leveldb/testDB/*
	go test -v src/db/*.go