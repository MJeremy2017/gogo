package main

import (
	"fmt"
	"log"
)


func AssertEqual[T comparable](a, b T) {
	if a == b {
		log.Printf("Passed, %+v equals %+v", a, b)
	} else {
		log.Fatalf("Failed, %+v does not equal %+v", a, b)
	}
}

type Stack[T any] struct {
	values []T
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	lastIndex := len(s.values) - 1
	value := s.values[lastIndex]
	s.values = s.values[:lastIndex]
	return value, true
}


func main() {
	AssertEqual(1, 1)
	AssertEqual("A", "A")
	
	// int stack
	stack := new(Stack[int])
	fmt.Printf("%v \n", stack.IsEmpty())
	stack.Push(12)
	stack.Push(2131)
	fmt.Printf("%v\n", stack.values)
	
	// string stack
	stackStr := new(Stack[string])
	stackStr.Push("hehe")
	fmt.Printf("%+v\n", stackStr)
}
