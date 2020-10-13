package main

import (
	"fmt"
	"go-core/1-intro/pkg/fibo"
)

func main() {
	nums := []int{7, 13, 20, -8, -9}
	for _, n := range nums {
		fmt.Printf("Fibonacci number for position %d is %d\n", n, fibo.Fibo(n))
	}
}
