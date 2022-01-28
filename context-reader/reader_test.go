package some


import (
	"testing"
	"strings"
)

func TestContextAwareReader(t *testing.T) {
	t.Run("normal reader functions", func(t *testing.T) {
		rdr := strings.NewReader("123456")
		got := make([]byte, 3)
		_, err := rdr.Read(got)

		if err != nil {
			t.Fatal(err)
		}

		assertBufferHas(t, got, "123")

		_, err = rdr.Read(got)
		if err != nil {
			t.Fatal(err)
		}

		assertBufferHas(t, got, "456")
	})
}


func assertBufferHas(t testing.TB, buf []byte, want string) {
	t.Helper()
	got := string(buf)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}