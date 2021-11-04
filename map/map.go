package main

import (
	"errors"
)

var ErrKeyNotFound = errors.New("word not found")

type Dictionary map[string]string

func (d Dictionary) Search(key string) (string, error) {
	value, ok := d[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (d Dictionary) Add(key string, value string) {
	d[key] = value
}