package arrays

func Sum(arr []int) int {
	s := 0
	for _, v := range arr {
		s += v
	}
	return s
}

func SumAll(arrs ...[]int) []int {
	var res []int
	for _, arr := range arrs {
		res = append(res, Sum(arr))
	}
	return res
}

func SumAllTails(arrs ...[]int) []int {
	var res []int
	for _, arr := range arrs {
		if len(arr) == 0 {
			res = append(res, 0)
		} else {
			res = append(res, Sum(arr[1:]))
		}
		
	}
	return res
}