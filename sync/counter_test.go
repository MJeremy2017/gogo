package sync

import (
	"testing"
	"sync"
)

func TestCounter(t *testing.T) {
	t.Run("Increase counter 3 times", func(t *testing.T) {
		counter := Counter{}
		counter.Incr()
		counter.Incr()
		counter.Incr()

		got := counter.Value()
		want := 3

		assertCounter(t, got, want)
	})

	t.Run("Increase counter concurrently", func(t *testing.T) {
		counter := Counter{}
		count := 1000

		var wg sync.WaitGroup
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				counter.Incr()
				wg.Done()
			}()
		}
		wg.Wait()

		got := counter.Value()
		want := count
		assertCounter(t, got, want)
	})

}

func assertCounter(t testing.TB, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}