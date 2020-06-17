package main

import (
	"flag"
	"fmt"
)

func main()  {
	_ = flag.String("a", "", "")
	flag.Parse()
	b := flag.String("b", "", "")
	flag.Parse()
	fmt.Println(b)
}

