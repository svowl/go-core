package fibo

import "testing"

func TestFibo(t *testing.T) {
	testData := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765, 10946, 17711}
	for n, exp := range testData {
		got := Fibo(n)
		if got != exp {
			t.Fatalf("Для номера %d получено число %d, ожидается %d", n, got, exp)
		}
	}

	testData = []int{0, 1, -1, 2, -3, 5, -8, 13, -21, 34, -55, 89, -144, 233, -377, 610, -987, 1597, -2584, 4181, -6765, 10946, -17711}
	for n, exp := range testData {
		got := Fibo(-n)
		if got != exp {
			t.Fatalf("Для номера %d получено число %d, ожидается %d", -n, got, exp)
		}
	}
}
