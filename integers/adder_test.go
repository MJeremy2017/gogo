package integers

import "testing"
import "fmt"

func TestAdd(t *testing.T) {
	assertError := func(t testing.TB, got int, want int) {
		t.Helper()
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	}

	t.Run("test add", func(t *testing.T) {
		got := Add(2, 2)
		want := 4
		if got != want {
			assertError(t, got, want)
		}
	})
}

func ExampleAdd() {
	sum := Add(1, 4)
	fmt.Println(sum)
	// Output: 5
}