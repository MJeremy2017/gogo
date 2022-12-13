package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBrowseLinks(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BrowseLinks(tt.args.url)
		})
	}
}

func TestGetHtmlPage(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHtmlPage(tt.args.url)
			if !tt.wantErr(t, err, fmt.Sprintf("GetHtmlPage(%v)", tt.args.url)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetHtmlPage(%v)", tt.args.url)
		})
	}
}
