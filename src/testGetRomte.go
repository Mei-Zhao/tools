package main

import (
	"io/ioutil"
	"crypto/md5"

	"github.com/qiniu/rpc.v1"
	"github.com/qiniu/xlog.v1"
	"github.com/qiniu/errors"
	"os"
)

func main() {
	xl := xlog.NewDummy()
	url := os.Args[1]

	fetchRemote(xl, url)
}

func fetchRemote(xl *xlog.Logger, URL string) (data, md5sum []byte, err error) {

	resp, err := rpc.DefaultClient.Get(xl, URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = rpc.ResponseError(resp)
		return
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Info(err, "ioutil.ReadAll")
		return
	}
	md5sum = calcMd5sum(data)

	return
}


func calcMd5sum(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}