package scrape

import (
	"net/url"
	"reflect"
	"testing"
)

func TestScraper_FindCategory(t *testing.T) {
	tests := []struct {
		name string
		want map[string]url.URL
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scraper{}
			if got := s.FindCategory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
