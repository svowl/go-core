package fibo

// Fibo is calculates Fubonacchi number by its position n
// −21	13	−8	5	−3	2	−1	1	0	1	1	2	3	5	8	13	21
func Fibo(n int) int {
	if n == 0 {
		return 0
	}

	sign := 1
	if n < 0 {
		sign = -1
	}

	num := [3]int{0, 1, 1}

	for i := 2; i <= sign*n; i++ {
		num[2] = num[0] + num[1]
		num[0], num[1] = num[1], num[2]
	}

	res := num[2]

	if sign < 0 && n%2 == 0 {
		res = -1 * res
	}

	return res
}
