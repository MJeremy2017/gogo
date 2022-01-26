package osexec


import (
	"testing"
	"strings"
)


func TestGetData(t *testing.T) {
	input := strings.NewReader(
		`
<payload>
    <message>Happy Friend!</message>
</payload>`)

	got := GetData(input)
	want := "HAPPY FRIEND!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}


func TestGetDataIntegration(t *testing.T) {
	got := GetData(getXMLFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}