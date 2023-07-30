package lsm

import "sync"

type LSM struct {
	w *wal
	rwLock sync.RWMutex
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




