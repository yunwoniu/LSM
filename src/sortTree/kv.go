package sortTree

import "encoding/json"

type Kv struct {
	Key string
	Val string
	Delete bool
}

func (k *Kv)Marshal()([]byte,error){
	buf,err:=json.Marshal(k)
	if err != nil{
		return nil, err
	}
	return buf,nil
}

func (k *Kv)Unmarshal(buf []byte)error{
	return json.Unmarshal(buf,k)
}
