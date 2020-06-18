package main

import (
	"fmt"
)


func jch(key uint64, num int64) int64 {
	b, j := int64(-1), int64(0)
	for ; j < num; {
		b = j
		key = key * 2862933555777941757 + 1
		j = (b + 1) * int64(float64(1 << 31) / float64(key >> 33 + 1))
	}
	return b
}

func main()  {
	keys := []uint64{5, 9, 10, 16, 18, 20, 21}
	for i := int64(0); i < 16; i++ {
		fmt.Printf("%d servers\n", i)
		for _, v := range keys {
			fmt.Print(jch(v, i), ", ")
		}
		fmt.Println()
	}

}
