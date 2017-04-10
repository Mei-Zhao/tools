package main

import (
	"fmt"
	"os/exec"
	"os"
)

func main() {
	execDirAbsPath, _ := os.Getwd()
	fmt.Println("执行程序所在目录的绝对路径　　　　　　　:", execDirAbsPath)
	cmd := exec.Command("/bin/sh", "./check_mongodb_rollback.sh")

	bytes, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output: ", err)
		return
	}

	fmt.Println(string(bytes))
}