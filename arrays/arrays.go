package arrays

func Sum(arr []int) int {
	s := 0
	for _, v := range arr {
		s += v
	}
	return s
}

func SumAll(arrs ...[]int) []int {
	res := make([]int, len(arrs))
	for i, arr := range arrs {
		res[i] = Sum(arr)
	}
	return res
}