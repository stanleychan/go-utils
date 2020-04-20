package utils

import (
	"fmt"
	"sync"
)

type LinkedListObject interface{}

type LinkedListNode struct {
	Data LinkedListObject
	Next *LinkedListNode
}

type LinkedList struct {
	mutex *sync.RWMutex
	Head  *LinkedListNode
	Tail  *LinkedListNode
	Size  uint
	Cap   uint
}

func (list *LinkedList) Init(cap uint) {
	list.Size = 0
	list.Head = nil
	list.Tail = nil
	list.Cap = cap
	list.mutex = new(sync.RWMutex)
}

func (list *LinkedList) InitUnLimited() {
	list.Init(0)
}

func (list *LinkedList) Append(node *LinkedListNode) bool {
	if node == nil {
		return false
	}
	list.mutex.Lock()
	defer list.mutex.Unlock()
	if list.Size == 0 {
		list.Head = node
		list.Tail = node
		list.Size = 1
		return true
	}

	// circle list
	if list.Cap > 0 && list.Cap == list.Size {
		list.Delete(0)
	}

	tail := list.Tail
	tail.Next = node
	list.Tail = node
	list.Size += 1
	return true
}

func (list *LinkedList) Insert(index uint, node *LinkedListNode) bool {
	if node == nil {
		return false
	}

	if index > list.Size {
		return false
	}

	list.mutex.Lock()
	defer list.mutex.Unlock()

	if index == 0 {
		node.Next = list.Head
		list.Head = node
		list.Size += 1
		return true
	}
	var i uint
	ptr := list.Head
	for i = 1; i < index; i++ {
		ptr = ptr.Next
	}
	next := ptr.Next
	ptr.Next = node
	node.Next = next
	list.Size += 1
	return true
}

func (list *LinkedList) Delete(index uint) bool {
	if list == nil || list.Size == 0 || index > list.Size-1 {
		return false
	}

	list.mutex.Lock()
	defer list.mutex.Unlock()

	if index == 0 {
		head := list.Head.Next
		list.Head = head
		if list.Size == 1 {
			list.Tail = nil
		}
		list.Size -= 1
		return true
	}

	ptr := list.Head
	var i uint
	for i = 1; i < index; i++ {
		ptr = ptr.Next
	}
	next := ptr.Next

	ptr.Next = next.Next
	if index == list.Size-1 {
		list.Tail = ptr
	}
	list.Size -= 1
	return true
}

func (list *LinkedList) Get(index uint) *LinkedListNode {
	if list == nil || list.Size == 0 || index > list.Size-1 {
		return nil
	}

	list.mutex.RLock()
	defer list.mutex.RUnlock()

	if index == 0 {
		return list.Head
	}
	node := list.Head
	var i uint
	for i = 0; i < index; i++ {
		node = node.Next
	}
	return node
}

func (list *LinkedList) Display() {
	if list == nil || list.Size == 0 {
		fmt.Println("this single list is nil")
		return
	}
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	fmt.Printf("this single list size is %d \n", list.Size)
	ptr := list.Head
	var i uint
	for i = 0; i < list.Size; i++ {
		fmt.Printf("No%3d data is %v\n", i+1, ptr.Data)
		ptr = ptr.Next
	}
}
