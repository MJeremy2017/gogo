package blogposts

import (
	"io/fs"
	"io"
	"bufio"
	"strings"
	"fmt"
	"bytes"
)

type Post struct {
	Title 		string
	Description string
	Tags 		[]string
	Body 		string
}

const (
	titleSeparator = "Title: "
	descriptionSeparator = "Description: "
	tagsSeparator = "Tags: "
)


func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	for _, file := range dir {
		post, err := getPost(fileSystem, file.Name())
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, fileName string) (Post, error) {
	// open file
	postFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func newPost(postFile io.Reader) (Post, error) {
	// read data
	scanner := bufio.NewScanner(postFile)

	readLine := func(prefix string) string {
		scanner.Scan()
		txt := scanner.Text()
		return strings.TrimPrefix(txt, prefix)
	}
	title := readLine(titleSeparator)
	description := readLine(descriptionSeparator)
	tags := strings.Split(readLine(tagsSeparator), ", ")

	scanner.Scan() // skip one line

	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}
	body := strings.TrimSuffix(buf.String(), "\n")

	post := Post{
		Title: title, 
		Description: description, 
		Tags: tags,
		Body: body,
	}
	return post, nil
}

