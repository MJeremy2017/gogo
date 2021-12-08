package main

import (
	"blogposts"
	"log"
	"os"
)


func main() {
	posts, err := blogposts.NewPostsFromFS(os.DirFS("cmd/posts"))
    if err != nil {
        log.Fatal(err)
    }
    log.Println(posts)
}