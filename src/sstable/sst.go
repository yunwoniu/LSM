package sstable

import (
	"lsm"
	"os"
	"strconv"
)

type TableTree struct {
	levels        [][]sst
	l0MaxSize     int
	levelSizeRate int
	levelMaxSize  []int
	levelLen      int
}

type sst struct {
	level int
	index int
	f     *os.File
}

func LoadTableTree(dir string, opt lsm.Options) (*TableTree, error) {

}


