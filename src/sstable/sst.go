package sstable

import (
	"fmt"
	"log"
	"lsm"
	"os"
	"path"
	"sortTree"
)

const sstSplit = "_"

type TableTree struct {
	levels        [][]sst
	l0MaxSize     int
	levelSizeRate int
	levelMaxSize  []int
	levelLen      int
}

type sst struct {
	level    int
	index    int
	f        *os.File
	smallKey string
	bigKey   string
}

type Mate struct {
	SmallKey string
	BigKey string
	KeyIndex map[string]int
}

func MemToSst(dir string, kvs []*sortTree.Kv, level, index int) (*sst, error) {
	sstName := fmt.Sprintf("%d%s%d", level, sstSplit, index)
	fd, err := os.Create(path.Join(dir, sstName))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer fd.Close()

	mate:=&Mate{
		SmallKey: kvs[0].Key,
		BigKey:   kvs[len(kvs)-1].Key,
	}
	keyIndex := make(map[string]int)
	for _,kv:= range kvs{
		_,buf:=kv.Marshal()

	}



	return &sst{
		level:    level,
		index:    index,
		smallKey: kvs[0].Key,
		bigKey:   kvs[len(kvs)-1].Key,
	}, nil
}

func LoadTableTree(dir string, opt lsm.Options) (*TableTree, error) {

}
