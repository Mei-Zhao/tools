package main

import (
	"os/exec"
	"os"
)

func main() {
	cmd := exec.Command("/bin/sh/", "-c", "time hdfs dfs -cat /bi/xs/gz/2016-10-19/REQ_UP/* | pv | pigz -d | sift -e product/003/052/3052122_std/3052122_pop_375_500_3.jpg -e FmW3h2AM6TGL_OebVjE9RuNEO7Do -e 14768390544582604 | tee xs_up_1380370752_5.log")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Run()
	cmd.Wait()
}
