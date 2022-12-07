package link

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParser_ParseLinks(t *testing.T) {
	tests := []struct {
		name     string
		htmlData string
		want     []Link
	}{
		// TODO: Add test cases.
		{
			name: "can-get-link-from-a-single-html-data",
			htmlData: `
<a href="/dog">
  <span>Something in a span</span>
  Text not in a span
  <b>Bold text!</b>
</a>
`,
			want: []Link{
				{
					Href: "/dog",
					Text: "Something in a span Text not in a span Bold text!",
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.htmlData)
			p := NewParser(buf)
			if got := p.ParseLinks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLinks() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
