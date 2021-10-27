package arrays

import "testing"
import "reflect"

func TestSum(t *testing.T) {
	t.Run("sum with any size", func(t *testing.T) {
		arr := []int{1, 2, 3}
		got := Sum(arr)
		want := 6

		if got != want {
			t.Errorf("Expect %d got %d given %v", want, got, arr)
		}
	})

}

func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2, 3}, []int{0, 2})
	want := []int{6, 2}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expect %d got %d", want, got)
	}
}