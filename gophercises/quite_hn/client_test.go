package hn

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup() (string, func()) {
	mux := http.NewServeMux()
	mux.HandleFunc("/topstories.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "[0,1,2,3,4]")
	})
	mux.HandleFunc("/item/", func(w http.ResponseWriter, r *http.Request) {
		id := getId(r)
		_, _ = fmt.Fprint(w, fmt.Sprintf("{\"by\":\"test_user\",\"descendants\":10,\"id\": %s,\"kids\":[16732999,16729637,16729517,16729595],\"score\":34,\"time\":1522599083,\"title\":\"Test Story Title\",\"type\":\"story\",\"url\":\"https://www.test-story.com\"}", id))
	})
	server := httptest.NewServer(mux)
	return server.URL, func() {
		server.Close()
	}
}

func getId(r *http.Request) string {
	s := strings.Split(r.URL.Path, "/")
	return strings.TrimRight(s[len(s)-1], ".json")
}

func TestClient_TopItems(t *testing.T) {
	baseURL, teardown := setup()
	defer teardown()

	c := Client{
		apiBase: baseURL,
	}
	ids, err := c.TopItems()
	if err != nil {
		t.Errorf("client.TopItems() received an error: %s", err.Error())
	}
	if len(ids) != 5 {
		t.Errorf("len(ids): want %d, got %d", 5, len(ids))
	}
}

func TestClient_defaultify(t *testing.T) {
	var c Client
	c.defaultify()
	if c.apiBase != apiBase {
		t.Errorf("c.apiBase: want %s, got %s", apiBase, c.apiBase)
	}
}

func TestClient_GetItem(t *testing.T) {
	baseURL, teardown := setup()
	defer teardown()

	c := Client{
		apiBase: baseURL,
	}
	item, err := c.GetItem(1)
	if err != nil {
		t.Errorf("client.GetItem() received an error: %s", err.Error())
	}
	// If this stuff errors it means our JSON is incorrect, which is unlikely, so
	// we can just check one field and consider that enough
	if item.By != "test_user" {
		t.Errorf("item.By: want %s, got %s", "test_user", item.By)
	}
}

func TestClient_GetBatchItems(t *testing.T) {
	baseURL, teardown := setup()
	defer teardown()

	c := Client{
		apiBase: baseURL,
	}
	wantIds := []int{1, 2, 3, 4, 5}
	items, err := c.GetOrderedBatchItems(wantIds)
	if err != nil {
		t.Errorf("client.BatchItems() received an error: %s", err.Error())
	}

	var gotIds []int
	for _, item := range items {
		gotIds = append(gotIds, item.ID)
	}
	assert.Equal(t, len(wantIds), len(gotIds))
	assert.Equal(t, wantIds, gotIds)
}

func TestClient_FilterStories(t *testing.T) {
	c := &Client{}

	tests := []struct {
		name  string
		items []Item
		n     int
		want  []Item
	}{
		{
			"filter-stories",
			[]Item{
				{
					ID:   1,
					URL:  "abc",
					Type: "story",
				},
				{
					ID:   2,
					URL:  "ab",
					Type: "comments",
				},
			},
			2,
			[]Item{
				{
					ID:   1,
					URL:  "abc",
					Type: "story",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, c.FilterStories(tt.items, tt.n), "FilterStories(%v)", tt.items)
		})
	}
}
