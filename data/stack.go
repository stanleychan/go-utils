package data

import (
	"sync"
)

type Item = interface{}

type Stack struct {
	items []Item
	lock  sync.RWMutex
}

func (s *Stack) New() *Stack {
	s.items = make([]Item, 0)
	return s
}

func (s *Stack) Push(data interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = append(s.items, data)
}

func (s *Stack) Pop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.items) < 1 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	return item
}

func (s *Stack) Top() interface{} {
	if len(s.items) < 1 {
		return nil
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}
