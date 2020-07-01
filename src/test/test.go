package main

import "log"

func main() {
	a := make(map[string]bool)
	ch := make(chan bool, 3)

	go func() {
		a["1"] = false
		ch <- true
	}()
	<-ch
	log.Println(a)

}
