package handlers

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_parseYAMLtoMap(t *testing.T) {
	type args struct {
		yml []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "test-can-parse-yaml",
			args: args{
				yml: []byte(
					`
- path: /some-path
  url: https://www.some-url.com/demo
`)},
			want: map[string]string{"/some-path": "https://www.some-url.com/demo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseYAMLtoMap(tt.args.yml)
			assert.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseYAMLtoMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseJSONtoMap(t *testing.T) {
	type args struct {
		js []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "test-can-parse-json",
			args: args{
				js: []byte(
					`
[{"path": "/some-path", "url": "https://www.some-url.com/demo"}]
`)},
			want: map[string]string{"/some-path": "https://www.some-url.com/demo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJSONtoMap(tt.args.js)
			assert.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseJSONtoMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
