package iterations

func Repeat(ch string, times int) string {
	var repeat string
	for i := 0; i < times; i++ {
		repeat += ch
	}
	return repeat
}