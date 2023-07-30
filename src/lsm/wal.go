package lsm

import (
	"encoding/binary"
	"log"
	"os"
	"path/filepath"
	"sortTree"
	"sync"
)

type wal struct {
	f *os.File
	lock sync.Mutex
}



func (w *wal)Init(dir string)(*sortTree.SortTree ,error){
	w.lock = sync.Mutex{}
	f,err := os.OpenFile(filepath.Join(dir,"wal.log"),os.O_CREATE|os.O_SYNC|os.O_RDWR,0666)
	if err != nil{
		log.Println(err)
		return nil,err
	}
	w.f = f
	fi,err := f.Stat()
	if err != nil{
		log.Println(err)
		return nil,err
	}
	buf:=make([]byte,fi.Size())
	_,err = f.Read(buf)
	if err != nil{
		log.Println(err)
		return nil,err
	}
	index:=int(0)
	tree := &sortTree.SortTree{}
	for index < len(buf){
		lenght:=binary.LittleEndian.Uint64(buf[index:index+8])
		index+=8
		kv := &sortTree.Kv{}
		err := kv.Unmarshal(buf[index:index+int(lenght)])
		if err != nil{
			return nil,err
		}
		if kv.Delete {
			tree.Delete(kv.Key)
		}else{
			tree.Set(kv.Key,kv.Val)
		}
		index+=int(lenght)
	}
	return tree,nil
}

func (w *wal)Write(key,value string,isDetele bool)error{
	w.lock.Lock()
	defer w.lock.Unlock()
	kv := &sortTree.Kv{
		Key: key,
		Val: value,
		Delete: isDetele,
	}
	buf,err := kv.Marshal()
	if err!= nil{
		log.Println(err)
		return err
	}
	lenght:=len(buf)
	err = binary.Write(w.f,binary.LittleEndian,lenght)
	if err != nil{
		log.Println(err)
		return err
	}
	err = binary.Write(w.f,binary.LittleEndian,buf)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}


