package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
)

type IndexedItem struct {
	index int
	item  Item
}

// Client is an API client used to interact with the Hacker News API
type Client struct {
	// unexported fields...
	apiBase string
}

// Making the Client zero value useful without forcing users to do something
// like `NewClient()`
func (c *Client) defaultify() {
	if c.apiBase == "" {
		c.apiBase = apiBase
	}
}

// TopItems returns the ids of roughly 450 top items in decreasing order. These
// should map directly to the top 450 things you would see on HN if you visited
// their site and kept going to the next page.
//
// TopItems does not filter out job listings or anything else, as the type of
// each item is unknown without further API calls.
func (c *Client) TopItems() ([]int, error) {
	c.defaultify()
	resp, err := http.Get(fmt.Sprintf("%s/topstories.json", c.apiBase))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var ids []int
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetItem will return the Item defined by the provided ID.
func (c *Client) GetItem(id int) (Item, error) {
	c.defaultify()
	var item Item
	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", c.apiBase, id))
	if err != nil {
		return item, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}

// TODO use sync.Mutex to counter the number of stories async
// TODO add caching

// GetOrderedBatchItems grab items asynchronously and return the items in its original order
func (c *Client) GetOrderedBatchItems(ids []int) ([]Item, error) {
	size := len(ids)
	ch := make(chan IndexedItem, size)
	items := make([]Item, size)
	for i, id := range ids {
		go c.asyncFetchItem(i, id, ch)
	}
	cnt := 0
	for it := range ch {
		items[it.index] = it.item
		cnt += 1
		if cnt == len(ids) {
			break
		}
	}
	return items, nil
}

// GetNOrderedBatchItems grab items asynchronously and return the top N items in its original order from ids
func (c *Client) GetNOrderedBatchItems(ids []int, n int) ([]Item, error) {
	size := len(ids)
	ch := make(chan IndexedItem, size)
	items := make([]Item, size)
	i := 0
	succeedCnt := make(chan int)
	for i < size {
		fmt.Println("got i", i)
		go func(index int) {
			sc := <-succeedCnt
			fmt.Println("sccc", sc, n)
			if sc >= n {
				close(ch)
				return
			}
			resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", c.apiBase, ids[index]))
			if err != nil {
				ch <- IndexedItem{index, Item{}}
				return
			}
			defer resp.Body.Close()

			var item Item
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&item)
			if err != nil {
				ch <- IndexedItem{index, Item{}}
				return
			}
			ch <- IndexedItem{index, item}
			succeedCnt <- sc + 1
		}(i)
		i += 1
	}
	succeedCnt <- 0
	total := 0
	for it := range ch {
		items[it.index] = it.item
		total += 1
		if total >= size {
			fmt.Println("break", total)
			break
		}
	}
	var filteredItems []Item
	for _, item := range items {
		if item.ID != 0 {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems, nil
}

func (c *Client) asyncFetchItem(i, id int, ch chan IndexedItem) {
	// `i` is the index of item `id`
	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", c.apiBase, id))
	if err != nil {
		ch <- IndexedItem{i, Item{}}
		return
	}
	defer resp.Body.Close()

	var item Item
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		ch <- IndexedItem{i, Item{}}
		return
	}
	ch <- IndexedItem{i, item}
}

// FilterStories filters items and get stories (This function retains the original orders)
// Have a Type of "story". This filters out all job postings and other types of items.
// Have a URL instead of Text. This filters out things like Ask HN questions and other discussions.
func (c *Client) FilterStories(items []Item, n int) []Item {
	var isStory = func(item Item) bool {
		return item.Type == "story" && item.URL != ""
	}

	var res []Item
	for _, item := range items {
		if isStory(item) {
			res = append(res, item)
		}
	}
	if len(res) < n {
		return res
	}
	return res[:n]
}

// Item represents a single item returned by the HN API. This can have a type
// of "story", "comment", or "job" (and probably more values), and one of the
// URL or Text fields will be set, but not both.
//
// For the purpose of this exercise, we only care about items where the
// type is "story", and the URL is set.
type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	// Only one of these should exist
	Text string `json:"text"`
	URL  string `json:"url"`
}
