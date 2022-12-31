package hn

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
)

type IndexedItem struct {
	index int
	item  Item
}

type task struct {
	index int
	id    int
}

// Client is an API client used to interact with the Hacker News API
type Client struct {
	// unexported fields...
	apiBase string
	cache   *Cache
}

func NewCacheClient(expirationSec int64) *Client {
	return &Client{
		apiBase: apiBase,
		cache:   NewCache(expirationSec),
	}
}

type Cache struct {
	stories       []Item
	createdAt     time.Time
	expirationSec int64
}

func (c *Cache) GetStories() ([]Item, error) {
	if c.createdAt.IsZero() {
		return nil, fmt.Errorf("zero creation time %v", c.createdAt)
	}
	lapse := int64(time.Now().Sub(c.createdAt).Seconds())
	if lapse > c.expirationSec {
		return nil, fmt.Errorf("cache expired %d > %d", lapse, c.expirationSec)
	}
	return c.stories, nil
}

func (c *Cache) UpdateStories(stories []Item) {
	c.stories = stories
	c.createdAt = time.Now()
	log.Printf("stories cached at %v", c.createdAt)
}

func NewCache(expirationSec int64) *Cache {
	return &Cache{expirationSec: expirationSec}
}

// Making the Client zero value useful without forcing users to do something
// like `NewClient()`
func (c *Client) defaultify() {
	if c.apiBase == "" {
		c.apiBase = apiBase
	}
}

func (c *Client) GetTopStories(numStories int) ([]Item, error) {
	// Law of Demeter
	stories, err := c.getCache().GetStories()
	if err == nil {
		log.Printf("using cached stories retrieved at %v \n", c.cache.createdAt)
		return stories, nil
	}
	log.Println(err)
	ids, err := c.TopItems()
	if err != nil {
		return nil, err
	}
	i := 0
	var result []Item
	for len(result) < numStories {
		need := (numStories - len(result)) * 5 / 4
		items := c.GetOrderedBatchItems(ids[i : i+need])
		stories := c.FilterStories(items)
		for _, story := range stories {
			result = append(result, story)
		}
		i += need
	}
	c.getCache().UpdateStories(result)
	if len(result) < numStories {
		return result, nil
	}
	return result[:numStories], nil
}

func (c *Client) getCache() *Cache {
	return c.cache
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

// TODO check solution
// GetOrderedBatchItems grab items asynchronously and return the items in its original order
func (c *Client) GetOrderedBatchItems(ids []int) []Item {
	const numGo = 10
	taskChan := make(chan task, len(ids))
	size := len(ids)
	ch := make(chan IndexedItem, size)
	items := make([]Item, size)

	for i := 0; i < numGo; i++ {
		go c.asyncFetchItem(taskChan, ch)
	}

	for i, id := range ids {
		taskChan <- task{i, id}
	}
	close(taskChan)

	cnt := 0
	for it := range ch {
		items[it.index] = it.item
		cnt += 1
		if cnt == len(ids) {
			break
		}
	}
	return items
}

func (c *Client) asyncFetchItem(taskChan <-chan task, ch chan IndexedItem) {
	// `i` is the index of item `id`
	for task := range taskChan {
		resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", c.apiBase, task.id))
		if err != nil {
			ch <- IndexedItem{task.index, Item{}}
			return
		}
		defer resp.Body.Close()

		var item Item
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&item)
		if err != nil {
			ch <- IndexedItem{task.index, Item{}}
			return
		}
		ch <- IndexedItem{task.index, item}
	}
}

// FilterStories filters items and get stories (This function retains the original orders)
// Have a Type of "story". This filters out all job postings and other types of items.
// Have a URL instead of Text. This filters out things like Ask HN questions and other discussions.
func (c *Client) FilterStories(items []Item) []Item {
	var isStory = func(item Item) bool {
		return item.Type == "story" && item.URL != ""
	}

	var res []Item
	for _, item := range items {
		if isStory(item) {
			res = append(res, item)
		}
	}
	return res
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
