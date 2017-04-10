package main

import(
	"fmt"
	"strings"
)

type Node struct  {
	content string
	len int

}

func main () {
	//nodes := putLen()
	//fmt.Println("%v\n", nodes)
	key := "xxx-test.jpg"
	separator := ".-"
	n := strings.LastIndexAny(key, separator)
	fmt.Println("n: ", n)
	fmt.Println("key: ",key[:n])
	style := key[n+1:]
	fmt.Println("style", style)
}

func putLen () (nodes []Node) {
	node1 := Node{content:"aaa"}
	node2 := Node{content:"bbb"}
	nodes = append(nodes, node1)
	nodes = append(nodes, node2)
	node1.len = 5
	node2.len = 6
	return
}