package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	file1, err := os.Open("/Users/zhaomei/qbox/tools/src/file/a")
	if err != nil {
		log.Fatal("os.Open 1: ", err)
	}
	//defer file1.Close()
	wg.Add(1)
	go func() {
		file2, err := os.OpenFile("/Users/zhaomei/qbox/tools/src/file/a", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)
		if err != nil {
			log.Fatal("os.Open 2: ",err)
		}
		fmt.Println("file2 open")
		time.Sleep(1* time.Minute)
		fmt.Println("write")
		_, err = file2.WriteString("hello")
		if err != nil {
			log.Fatal("file2.WriteString: ", err)
		}
		info2, err := file2.Stat()
		if err != nil {
			log.Fatal("fil1.Stat: ", err)
		}
		fmt.Println("size2 ", info2.Size())
		defer file2.Close()
		defer wg.Done()
	}()

	info, err := file1.Stat()
	if err != nil {
		log.Fatal("fil1.Stat: ", err)
	}
	fmt.Println("size1 ", info.Size())
	time.Sleep(30*time.Second)
	file1.Close()
	fmt.Println("file1 close")
	err = os.Remove("/Users/zhaomei/qbox/tools/src/file/a")
	if err != nil {
		log.Fatal(" os.Remove err:", err)
	}
	log.Println("remove")
	wg.Wait()
	return
}
