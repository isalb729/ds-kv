package main

import (
	"context"
	"fmt"
	"time"
)

func main()  {
	//ctx := context.Background()
	ctx, _ := context.WithTimeout(context.Background(), 2 * time.Second)
	fmt.Println(ctx.Deadline())
	time.Sleep(3 * time.Second)
	ctx2, _ := context.WithTimeout(ctx, 5 * time.Second)
	fmt.Println(<-ctx2.Done())
	fmt.Println(ctx2.Deadline())
	fmt.Println(<-ctx.Done())
	fmt.Println(ctx.Err())
}
