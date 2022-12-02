package main

import (
	"adventure/parser"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

var testChapter = parser.Chapter{
	Title: "title",
	Story: []string{"a black bird"},
	Options: []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	}{
		{"text", "arc"},
	},
}
var testStory = parser.Story{"key1": testChapter}

func Test_getChapter(t *testing.T) {
	type args struct {
		story parser.Story
		key   string
	}
	tests := []struct {
		name    string
		args    args
		want    parser.Chapter
		wantErr bool
	}{
		{
			name: "can-get-chapter-with-valid-key",
			args: args{
				story: testStory,
				key:   "key1",
			},
			want:    testChapter,
			wantErr: false,
		},
		{
			name: "return-error-with-invalid-key",
			args: args{
				story: testStory,
				key:   "key2",
			},
			want:    parser.Chapter{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getChapter(tt.args.story, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChapter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getChapter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storyHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name             string
		args             args
		expectedTitle    string
		expectedLocation string
	}{
		{
			name: "can-return-correct-chapter-with-valid-path",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/new-york", nil),
			},
			expectedTitle:    "<h2>Visiting New York</h2>",
			expectedLocation: "/new-york",
		},
		{
			name: "redirect-to-home-with-invalid-path",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/abc", nil),
			},
			expectedTitle:    "<h2></h2>",
			expectedLocation: "/home",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.args.w
			r := tt.args.r
			storyHandler(w, r)

			match := extractTitle(w)
			assert.Equal(t, tt.expectedTitle, match)
			assert.Equal(t, tt.expectedLocation, w.Result().Header.Get("Location"))
		})
	}
}

func extractTitle(w *httptest.ResponseRecorder) string {
	re := regexp.MustCompile(`<h2>(.*)</h2>`)
	return re.FindString(w.Body.String())
}
