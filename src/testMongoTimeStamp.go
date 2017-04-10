package main

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"time"
	"encoding/binary"
	//"net/http"
	//"log"

)

func main () {
	lastIndex := int64(6356057139970048422)
	var time bson.MongoTimestamp
	time = bson.MongoTimestamp(lastIndex)
	fmt.Printf("time: %#v\n",time)
	lastIndex2 := MongoTimestampToTime(lastIndex)
	fmt.Printf("time2: %#v\n",lastIndex2)
	compare(lastIndex2)
}

func MongoTimestampToTime(ts int64) int64 {
	b := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(b, uint64(ts))
	return time.Unix(int64(binary.BigEndian.Uint32(b[:4])), 0).UnixNano()/100
}

func compare(mongoTime int64){
	//putTime延迟self.DelayTime天确认
	diffNowPut := time.Since(time.Unix(mongoTime/1e7, 0))
	shouldDelayTime := time.Duration(2*24) * time.Hour
	if diffNowPut < shouldDelayTime {
		fmt.Printf("延迟不足%v天, sleep %v\n", 2, shouldDelayTime-diffNowPut)
		time.Sleep(shouldDelayTime - diffNowPut)
	}
}
