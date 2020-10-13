package main

import (
	"fmt"
	"go-core/1-intro/pkg/fibo"
)

func main() {
	nums := []int{7, 13, 20, -8, -9}
	for _, n := range nums {
		fmt.Printf("Число Фибоначчи номер %d это %d\n", n, fibo.Fibo(n))
	}
}
