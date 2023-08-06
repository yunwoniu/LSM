package lsm

import (
	"os"
	"sortTree"
	"sync"
)

type LSM struct {
	w *wal
	memTable *sortTree.SortTree
	rLock sync.RWMutex
}

type Options struct{

}

func Open(dir string,opt *Options)(*LSM,error){
	os.MkdirAll(dir,0666)
	var wal wal
	mem,err:=wal.Init(dir)
	if err != nil{
		return nil,err
	}
	lsm:=&LSM{
		w: &wal,
		memTable: mem,
		rLock: sync.RWMutex{},
	}
	return lsm,nil
}
func (l *LSM)Get(key string)(string,bool,error)  {
	return "",false,nil
}

func (l *LSM)Set(key,value string)error{

	return nil
}

func (l *LSM)Delete(key string)error{
	return nil
}




