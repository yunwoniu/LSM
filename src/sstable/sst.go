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

/*
	------------+-----------+-------------+------------+------------+--------------+-------------+--------+------+   ------------+
	 headLenght |smallKeyLen|  smallKey   | bigKeyLen  |  bigKey    |  indexMapLen | indexMap    | kvLen  | kv   |***| kvLen | kv
	------------+-----------+-------------+------------+------------+--------------+-------------+--------+------+   ------------+
		8Bytes  |  8Bytes   | smallKeyLen |  8Bytes    | bigKeyLen  |    8Bytes    | indexMapLen | 8Bytes |kvLen |***| 8Byte |kvLen
	------------+-----------+-------------+------------+------------+--------------+-------------+-------+------+   -------------+
*/

func memToSst(dir string, kvs []*sortTree.Kv, level, index int) (*sst, error) {
	sstName := fmt.Sprintf("%d%s%d", level, sstSplit, index)
	fd, err := os.Create(path.Join(dir, sstName))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer fd.Close()

	return &sst{
		level:    level,
		index:    index,
		smallKey: kvs[0].Key,
		bigKey:   kvs[len(kvs)-1].Key,
	}, nil
}

func LoadTableTree(dir string, opt lsm.Options) (*TableTree, error) {

}
