package blogposts_test

import (
	"testing"
	"testing/fstest"
	"blogposts"
	"reflect"
)


func TestNewBlogPosts(t *testing.T) {
	// both string and `` works
	const (
		firstBody = "Title: Post1\nDescription: Description1"
		secondBody = `Title: Post2
		Description: Description2`
	)

	// file path -> meta data
	fs := fstest.MapFS{
		"hello world.md": {Data: []byte(firstBody)},
		"hello world2.md": {Data: []byte(secondBody)},
	}

	posts, err := blogposts.NewPostsFromFS(fs)

	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != len(fs) {
		t.Errorf("got %d posts, want %d posts", len(posts), len(fs))
	}

	// test content matches
	got := posts[0]
	want := blogposts.Post{Title: "Post1", Description: "Description1"}
	assertPost(t, got, want)

}

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
