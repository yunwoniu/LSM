package sstTable

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"lsm"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sortTree"
	"strconv"
	"strings"
)

const sstSplit = "_"

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

func (s *TableTree) MemToSst(dir string, kvs []*sortTree.Kv, level, index int) error {
	sstName := fmt.Sprintf("%d%s%d", level, sstSplit, index)
	fd, err := os.Create(path.Join(dir, sstName))
	if err != nil {
		log.Println(err)
		return err
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
	binary.Write(headBuf, binary.LittleEndian, []byte(smallKey))
	binary.Write(headBuf, binary.LittleEndian, bigKeyLen)
	binary.Write(headBuf, binary.LittleEndian, []byte(bigKey))
	binary.Write(headBuf, binary.LittleEndian, indexMapLen)
	binary.Write(headBuf, binary.LittleEndian, indexMapData)
	_, err = fd.Write(headBuf.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = fd.Write(kvsBuf.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}

	sst := &sst{
		level:    level,
		index:    index,
		smallKey: kvs[0].Key,
		bigKey:   kvs[len(kvs)-1].Key,
		indexMap: indexMap,
	}
	s.levels[level] = append(s.levels[level], sst)
	return nil
}

/*
	------------+-----------+-------------+------------+------------+--------------+-------------+--------+------+   ------------+
	 headLenght |smallKeyLen|  smallKey   | bigKeyLen  |  bigKey    |  indexMapLen | indexMap    | kvLen  | kv   |***| kvLen | kv
	------------+-----------+-------------+------------+------------+--------------+-------------+--------+------+   ------------+
		8Bytes  |  8Bytes   | smallKeyLen |  8Bytes    | bigKeyLen  |    8Bytes    | indexMapLen | 8Bytes |kvLen |***| 8Byte |kvLen
	------------+-----------+-------------+------------+------------+--------------+-------------+-------+------+   -------------+
*/

func fileToSst(path string) (*sst, error) {
	name := filepath.Base(path)
	sstName := strings.TrimSuffix(name, lsm.Sst)
	sps := strings.Split(sstName, sstSplit)
	if len(sps) != 2 {
		return nil, errors.New(fmt.Sprintf("sst name:%s error", sstName))
	}
	level, err := strconv.Atoi(sps[0])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sst name:%s error", sstName))
	}
	index, err := strconv.Atoi(sps[1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sst name:%s error", sstName))
	}

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	headLenght := 0
	err = binary.Read(fd, binary.LittleEndian, headLenght)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	headData := make([]byte, headLenght-8)
	_, err = fd.Read(headData)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	buf := &bytes.Buffer{}
	buf.Write(headData)
	smallKeyLen := uint64(0)
	binary.Read(buf, binary.LittleEndian, smallKeyLen)
	smallKey := make([]byte, smallKeyLen)
	binary.Read(buf, binary.LittleEndian, smallKey)
	bigKeyLen := uint64(0)
	binary.Read(buf, binary.LittleEndian, bigKeyLen)
	bigKey := make([]byte, bigKeyLen)
	binary.Read(buf, binary.LittleEndian, bigKey)
	indexMapLen := uint64(0)
	indexMapData := make([]byte, indexMapLen)
	binary.Read(buf, binary.LittleEndian, indexMapData)
	indexMap := make(map[string]int)
	json.Unmarshal(indexMapData, &indexMap)
	return &sst{
		level:    level,
		index:    index,
		smallKey: string(smallKey),
		bigKey:   string(bigKey),
		f:        fd,
		indexMap: indexMap,
	}, nil
}

type TableTree struct {
	levels        [][]*sst
	l0MaxSize     int
	levelSizeRate int
	levelMaxSize  []int
	levelLen      int
}

func (s *TableTree) GetLevelMaxIndex(level int) int {
	return len(s.levels[level])
}

func LoadTableTree(dir string, opt *lsm.Options) (*TableTree, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	levels := make([][]*sst, opt.LevelLen)
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), lsm.Sst) || entry.IsDir() {
			continue
		}
		sst, err := fileToSst(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		levels[sst.level] = append(levels[sst.level], sst)
	}
	levelSize := opt.L0MaxSize
	levelMaxSize := make([]int, opt.LevelLen)
	for index, level := range levels {
		levelMaxSize[index] = levelSize
		levelSize = levelSize * opt.LevelSizeRate
		sort.Slice(level, func(i, j int) bool {
			if level[i].index < level[j].index {
				return true
			}
			return false
		})
	}
	return &TableTree{
		levels:        levels,
		levelMaxSize:  levelMaxSize,
		l0MaxSize:     opt.L0MaxSize,
		levelSizeRate: opt.LevelSizeRate,
		levelLen:      opt.LevelLen,
	}, nil
}
