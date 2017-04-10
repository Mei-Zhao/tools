package main

import (
	"github.com/qiniu/xlog.v1"
	"log"
	"os"
	"qbox.us/api/kmqcli"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//var name string
	//if len(os.Args) == 2 {
	//	name = os.Args[1]
	//}
	host := "http://192.168.66.26:6532"
	hosts := make([]string, 1)
	hosts = append(hosts, host)

	cfg := &kmqcli.Config{
		AccessKey: "ppca4hFBYQ_ykozmLUcSIJi8eLnYhFahE0OF5MoZ",
		SecretKey: "kc6oDxKD3TYoRq3lUoS41-e4qtNYWzBSQZmigm7K",
		//AccessKey: "4_odedBxmrAHiu4Y0Qp0HPG0NANCf6VAsAjWL_k9",
		//SecretKey: "SrRuUVfDX6drVRvpyN8mv8Vcm9XnMZzlbDfvVfMe",
		Hosts:     hosts,
		TryTimes:  2,
	}
	cli := kmqcli.New(cfg)
	xl := xlog.NewDummy()
	uid := uint32(1351985460)
	name := "OPLOGTEST"
	code, partions, err := cli.GetQueuePartitions(uid, name, xl)
	if err == nil && code == 200 {
		for _, partion := range partions {
			xl.Println("partion", partion)
		}
	} else {
		xl.Println("err", err, "code", code)
	}
}
