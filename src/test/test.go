package main

import (
	"fmt"
)

func main() {
	var a map[string]interface{}
	b := make(map[string]interface{})
	c := map[string]interface{}{"1":2}
	a["1"]=2
	b["1"]=3
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)

}
