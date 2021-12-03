package blogposts

import (
	"io/fs"
	"io"
)

type Post struct {
	Title string
}


func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	for _, d := range dir {
		post, err := getPost(fileSystem, d)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, f fs.DirEntry) (Post, error) {
	// open file
	postFile, err := fileSystem.Open(f.Name())
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()

	// read data
	data, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(data)[7:]}
	return post, nil
}