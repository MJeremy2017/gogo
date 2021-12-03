package blogposts_test

import (
	"testing"
	"testing/fstest"
	"blogposts"
)


func TestNewBlogPosts(t *testing.T) {
	// file path -> meta data
	fs := fstest.MapFS{
		"hello world.md": {Data: []byte("Title: Post1")},
		"hello world2.md": {Data: []byte("Title: Post2")},
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
	want := blogposts.Post{Title: "Post1"}

	if got != want {
		t.Errorf("got %+v want %+v", got, want)
	}
}
