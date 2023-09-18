package parser

import (
	"errors"
	"fmt"
)

type Stack[T any] struct {
	data []T
}

func (S Stack[T]) PrintStack() {
	fmt.Println(S.data)
}

func (S *Stack[T]) Push(value T) {
	S.data = append(S.data, value)
}

func (S *Stack[T]) Pop() (T, error) {
	var value T
	if len(S.data) == 0 {
		return value, errors.New("The stack is empty!")
	}
	value = S.data[len(S.data)-1]
	S.data = S.data[:len(S.data)-1]
	return value, nil
}

func (S *Stack[T]) Peek() (T, error) {
	var value T
	if len(S.data) == 0 {
		return value, errors.New("The stack is empty!")
	}
	value = S.data[len(S.data)-1]
	return value, nil
}
