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

func TestSumAllTails(t *testing.T) {
	assertEqual := func(t testing.TB, got []int, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expect %d got %d", want, got)
		}
	}

	t.Run("general test", func(t *testing.T) {
		got := SumAllTails([]int{1, 2, 3}, []int{0, 2})
		want := []int{5, 2}
		assertEqual(t, got, want)

	})

	t.Run("test empty slice", func(t *testing.T) {
		got := SumAllTails([]int{1, 2, 3}, []int{})
		want := []int{5, 0}
		assertEqual(t, got, want)
	})

}