package queue

import "sync"

type Queue struct {
	items   []interface{}
	maxSize int
	lock    sync.Mutex
}

func NewQueue(maxSize int) *Queue {
	return &Queue{
		items:   make([]interface{}, 0),
		maxSize: maxSize,
	}
}

func (q *Queue) Enqueue(item interface{}) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) >= q.maxSize {
		q.items = q.items[1:]
	}

	q.items = append(q.items, item)
	return true
}

func (q *Queue) Dequeue() (interface{}, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) == 0 {
		return nil, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.items)
}

func (q *Queue) IsEmpty() bool {
	return q.Len() == 0
}

func (q *Queue) IsFull() bool {
	return q.Len() >= q.maxSize
}

func (q *Queue) Clear() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = make([]interface{}, 0)
}

func (q *Queue) GetItems() []interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.items
}
