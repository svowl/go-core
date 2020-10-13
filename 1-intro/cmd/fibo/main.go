package main

import (
	"fmt"
	"go-core/1-intro/pkg/fibo"
)

func main() {
	n := 7
	fmt.Printf("Fibonacci number for position %d is %d", n, fibo.Fibo(n))
}
