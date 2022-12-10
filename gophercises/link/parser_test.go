package link

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_ParseLinks(t *testing.T) {
	tests := []struct {
		name     string
		htmlData string
		want     []Link
		hasErr   bool
	}{
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
			hasErr: false,
		},
		{
			name: "can-get-link-from-a-single-html-data-2",
			htmlData: `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>
`,
			want: []Link{
				{
					Href: "/other-page",
					Text: "A link to another page",
				}},
			hasErr: false,
		},
		{
			name: "can-get-link-from-multiple-html-data",
			htmlData: `
<html>
<head>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
<h1>Social stuffs</h1>
<div>
    <a href="https://www.twitter.com/joncalhoun">
        Check me out on twitter
        <i class="fa fa-twitter" aria-hidden="true"></i>
    </a>
    <a href="https://github.com/gophercises">
        Gophercises is on <strong>Github</strong>!
    </a>
</div>
</body>
</html>
`,
			want: []Link{
				{
					Href: "https://www.twitter.com/joncalhoun",
					Text: "Check me out on twitter",
				},
				{
					Href: "https://github.com/gophercises",
					Text: "Gophercises is on Github !",
				}},
			hasErr: false,
		},
		{
			name: "can-ignore-comment",
			htmlData: `
<html>
<body>
<a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
</body>
</html>
`,
			want: []Link{
				{
					Href: "/dog-cat",
					Text: "dog cat",
				}},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.htmlData)
			p := NewParser(buf)
			got, err := p.ParseLinks()
			if tt.hasErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got, fmt.Sprintf("ParseLinks() = %+v, want %+v", got, tt.want))
			}
		})
	}
}
