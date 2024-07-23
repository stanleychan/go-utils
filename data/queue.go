package data

type Queue struct {
	list *LinkedList
}

func (q *Queue) Init(cap uint) {
	q.list = new(LinkedList)
	q.list.Init(0)
}

func (q *Queue) Size() uint {
	return q.list.Size
}

func (q *Queue) Enqueue(data interface{}) bool {
	return q.list.Append(&LinkedListNode{Data: data})
}

func (q *Queue) Dequeue() interface{} {
	node := q.list.Get(0)
	if node == nil {
		return nil
	}
	q.list.Delete(0)
	return node.Data
}

func (q *Queue) Peek() interface{} {
	node := q.list.Get(0)
	if node == nil {
		return nil
	}
	return node.Data
}
