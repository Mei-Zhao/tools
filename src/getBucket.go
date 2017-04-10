package main

import (
	"os"
	"log"
	"io/ioutil"

	"github.com/qiniu/api/auth/digest"
)

var (
	mac = &digest.Mac{"ppca4hFBYQ_ykozmLUcSIJi8eLnYhFahE0OF5MoZ", []byte("kc6oDxKD3TYoRq3lUoS41-e4qtNYWzBSQZmigm7K")}
	cli = digest.NewClient(mac, nil)
)

func main() {
	var bucket string
	bucket = os.Args[1]
	path := "/bucket/" + bucket
	host:= "http://127.0.0.1:10220"

	resp, err := cli.Get(host + path)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}


}
