package lsm

import (
	"os"
	"sortTree"
	"sync"
)

type LSM struct {
	w        *wal
	memTable *sortTree.SortTree
	rLock    sync.RWMutex
}

type Options struct {
	L0MaxSize     int
	LevelSizeRate int
	LevelLen      int
}

func Open(dir string, opt *Options) (*LSM, error) {
	os.MkdirAll(dir, 0666)
	var wal wal
	mem, err := wal.Init(dir)
	if err != nil {
		return nil, err
	}
	lsm := &LSM{
		w:        &wal,
		memTable: mem,
		rLock:    sync.RWMutex{},
	}
	return lsm, nil
}
func (l *LSM) Get(key string) (string, bool, error) {
	val, is := l.memTable.Get(key)
	if is {
		return val, true, nil
	}

	return "", false, nil
}

func (l *LSM) Set(key, value string) error {
	l.rLock.Lock()
	defer l.rLock.Unlock()
	err := l.w.Write(key, value, false)
	if err != nil {
		return err
	}
	l.memTable.Set(key, value)
	return nil
}

func (l *LSM) Delete(key string) error {
	l.rLock.Lock()
	defer l.rLock.Unlock()
	err := l.w.Write(key, "", true)
	if err != nil {
		return err
	}
	l.memTable.Delete(key)
	return nil
}
