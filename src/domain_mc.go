package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/qiniu/xlog.v1"

	"qbox.us/memcache.v1"
	"qbox.us/api/one/domain"
)

var mcs = `[
		{

                        "host":"10.34.37.21:1141"
                },
                {

                        "host":"10.34.37.22:1141"
                }
	]
`


type memcacheValue struct {
	Phy              string `json:"phy" bson:"phy"`
	Tbl              string `json:"tbl" bson:"tbl"`
	Uid              uint32 `json:"uid" bson:"uid"`
	Itbl             uint32 `json:"itbl" bson:"itbl"`
	Refresh          bool   `json:"refresh" bson:"refresh"`
	Global           bool   `json:"global" bson:"global"`
	Domain           string `json:"domain" bson:"domain"`
	domain.AntiLeech `json:"antileech,omitempty" bson:"antileech,omitempty"`
}

func main() {
	var conf []memcache.Conn
	err := json.Unmarshal([]byte(mcs), &conf)
	if err != nil {
		log.Fatal(err)
	}
	mcS := memcache.New(conf)
	key := os.Args[1]
	xl := xlog.NewDummy()
	ret := memcacheValue{}
	err = mcS.Get(xl, key, &ret)
	if err != nil {
		log.Println(key, err)
		return
	}
	b, err := json.Marshal(ret)
	if err != nil {
		log.Println("marshal json failed", err)
		return
	}
	fmt.Printf("%v %v\n", key, string(b))
	var del string
	fmt.Print("\ndel?(y/n): ")
	_, err = fmt.Scan(&del)
	if err != nil {
		log.Println(err)
		return
	}
	if del != "y" && del != "Y" {
		return
	}
	err = mcS.Del(xl, key)
	if err == nil {
		fmt.Println("del success")
	} else {
		fmt.Println("del failed", err)
	}
}