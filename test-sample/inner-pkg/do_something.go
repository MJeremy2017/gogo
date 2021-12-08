package innerpkg

import (
	"fmt"
)

func Add(x, y int) int {
	fmt.Println("Inner Add")
	return x+y
}