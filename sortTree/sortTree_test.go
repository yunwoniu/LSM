package sortTree

import (
	"strconv"
	"testing"
)

func Test_Insert(t *testing.T){
	tree := &SortTree{}

	for i:=0;i<100;i++{
		key:=strconv.Itoa(i)
		tree.Set(key,key)
	}
	for i:=0;i<100;i++{
		key:=strconv.Itoa(i)
		val,_:=tree.Get(key)
		if key != val{
			t.Errorf("set get error")
		}
	}
	if tree.GetCount() != 100{
		t.Errorf("tree size error")
	}
}