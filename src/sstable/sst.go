package sstable

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	indexMap map[string]int
}

/*
	------------+-----------+-------------+------------+------------+--------------+-------------+--------+------+   ------------+
	 headLenght |smallKeyLen|  smallKey   | bigKeyLen  |  bigKey    |  indexMapLen | indexMap    | kvLen  | kv   |***| kvLen | kv
	------------+-----------+-------------+------------+------------+--------------+-------------+--------+------+   ------------+
		8Bytes  |  8Bytes   | smallKeyLen |  8Bytes    | bigKeyLen  |    8Bytes    | indexMapLen | 8Bytes |kvLen |***| 8Byte |kvLen
	------------+-----------+-------------+------------+------------+--------------+-------------+-------+------+   -------------+
*/

func MemToSst(dir string, kvs []*sortTree.Kv, level, index int) (*sst, error) {
	sstName := fmt.Sprintf("%d%s%d", level, sstSplit, index)
	fd, err := os.Create(path.Join(dir, sstName))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer fd.Close()

	headLenght := int(0)
	smallKey := kvs[0].Key
	smallKeyLen := len(smallKey)
	bigKey := kvs[len(kvs)-1].Key
	bigKeyLen := len(bigKey)
	indexMap := make(map[string]int)
	offset := 0
	kvsBuf := &bytes.Buffer{}
	for _, kv := range kvs {
		buf, _ := kv.Marshal()
		kvLen := len(buf)
		binary.Write(kvsBuf, binary.LittleEndian, kvLen)
		offset += 8
		indexMap[kv.Key] = offset
		offset += kvLen
		binary.Write(kvsBuf, binary.LittleEndian, buf)
	}
	indexMapData, _ := json.Marshal(indexMap)
	indexMapLen := len(indexMapData)
	headLenght = 8 + 8 + smallKeyLen + 8 + bigKeyLen + 8 + indexMapLen
	headBuf := &bytes.Buffer{}
	binary.Write(headBuf, binary.LittleEndian, headLenght)
	binary.Write(headBuf, binary.LittleEndian, smallKeyLen)
	binary.Write(headBuf, binary.LittleEndian, smallKey)
	binary.Write(headBuf, binary.LittleEndian, bigKeyLen)
	binary.Write(headBuf, binary.LittleEndian, bigKey)
	binary.Write(headBuf, binary.LittleEndian, indexMapLen)
	binary.Write(headBuf, binary.LittleEndian, indexMapData)
	_, err = fd.Write(headBuf.Bytes())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = fd.Write(kvsBuf.Bytes())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &sst{
		level:    level,
		index:    index,
		smallKey: kvs[0].Key,
		bigKey:   kvs[len(kvs)-1].Key,
		indexMap: indexMap,
	}, nil
}

func LoadTableTree(dir string, opt lsm.Options) (*TableTree, error) {

}
