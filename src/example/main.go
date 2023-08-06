package main

import "log"
import "lsm"

func main()  {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	db,err := lsm.Open("/tmp/lsm",nil)
	if err != nil{
		log.Println(err)
		return
	}
	db.Set("a","b")
	db.Get("a")
}
