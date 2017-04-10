package main

import (
	"fmt"
	"io/ioutil"
	"bytes"
	"compress/gzip"
)

func main() {

	size := 1
	datasSlice := make([]byte, 536870912)
	for i := 0; i < size; i++ {
		datasSlice, _ = ioutil.ReadFile("/Users/zhaomei/qbox/tools/src/file/a")
		fmt.Println("raw size:", len(datasSlice))
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()

	for i := 0; i < size; i++ {
		w.Write(datasSlice)
	}
	w.Flush()
	fmt.Println("gzip size:", len(b.Bytes()))

	r, _ := gzip.NewReader(&b)
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)
	fmt.Println("ungzip size:", len(undatas), string(undatas))

}
