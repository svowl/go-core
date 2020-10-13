package main

import (
	"flag"
	"fmt"
	"go-core/1-intro/pkg/fibo"
)

func main() {
	n := flag.Int("n", 0, "Specify position as -n <integer>")
	flag.Parse()
	fmt.Printf("Fibonacci number for position %d is %d\n", *n, fibo.Fibo(*n))
}
