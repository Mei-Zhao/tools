package main

import (
	"io/ioutil"
	"os"
	"fmt"

	"github.com/qiniu/errors"
)

func main () {
	uploadPath := "/Users/zhaomei/qbox/tools/src"
	fi, err := os.Stat(uploadPath)
	if err != nil {
		err = errors.Info(err, "os.Stat", uploadPath).Detail(err)
		fmt.Println("err", err)
	}
	fmt.Println("os.Stat fi",fi.Name() )
	if !fi.IsDir() {
		err = errors.New("not dir")
		err = errors.Info(err, uploadPath).Detail(err)
		fmt.Println("err", err)
	}
	folders, err :=ioutil.ReadDir(uploadPath)
	if err != nil {
		err = errors.Info(err, "os.Stat", uploadPath).Detail(err)
		fmt.Println("err", err)
	}
	//fmt.Println("folders", folders)
	for i := 0;i < len(folders); i++ {
		fmt.Println("fileName: ", folders[i].Name())
	}

}
