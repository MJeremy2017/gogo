package blogposts_test

import (
	"testing"
	"testing/fstest"
	"blogposts"
	"reflect"
)


func TestNewBlogPosts(t *testing.T) {
	wants := []blogposts.Post{
	blogposts.Post{
		Title: "Post1", 
		Description: "Description1",
		Tags: []string{"tdd", "go"},
		Body: `Hello
World!`,
	},
	blogposts.Post{
		Title: "Post2", 
		Description: "Description2",
		Tags: []string{"rust", "borrow-checker"},
		Body: `B
L
M`,
	},
}
	// both string and `` works
	const (
		firstBody = `Title: Post1
Description: Description1
Tags: tdd, go
---
Hello
World!
`
		secondBody = `Title: Post2
Description: Description2
Tags: rust, borrow-checker
---
B
L
M
`
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
	for i := 0; i < len(posts); i++ {
		got := posts[i]
		want := wants[i]
		assertPost(t, got, want)
	}

}

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
