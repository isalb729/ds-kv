package utils

import "math"

func GetPrimes(num, start int) []int {
	if num == 0 {
		return []int{}
	}
	var primeList []int
	if start <= 2 {
		primeList = []int{2}
		start = 3
	} else if start % 2 == 0 {
		start += 1
	}
	for i := start; len(primeList) < num; i += 2 {
		if IsPrime(i) {
			primeList = append(primeList, i)
		}
	}
	return primeList
}

func IsPrime(num int) bool {
	if num <= 1 {
		return false
	}
	if num == 2 {
		return true
	}
	for i := 2; i <= int(math.Sqrt(float64(num))); i++ {
		if num % i == 0 {
			return false
		}
	}
	return true
}
