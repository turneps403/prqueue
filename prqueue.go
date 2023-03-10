// Package prqueue implements a priority queue with generics and ability to send comparator-func as a constructor argument.
//
// For documentation expand section above.
package prqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
)

var ErrEmptyQueue = errors.New("priority queue is empty")

type pqs[T any] struct {
	mu   sync.RWMutex
	list []T
	cmp  func(a, b T) bool
}

func New[T any](cmp func(a, b T) bool, c ...int) (pq *pqs[T]) {
	pq = &pqs[T]{cmp: cmp}
	if len(c) != 0 {
		pq.list = make([]T, 0, c[0])
	} else {
		pq.list = make([]T, 0)
	}
	return
}

func (pq *pqs[T]) Add(el T) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(pq, el)
}

func (pq *pqs[T]) Poll() (el T, er error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if len(pq.list) > 0 {
		el = heap.Pop(pq).(T)
		return
	}
	er = ErrEmptyQueue
	return
}

func (pq *pqs[T]) Peek() (el T, er error) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	if len(pq.list) > 0 {
		el = pq.list[0]
		return
	}
	er = ErrEmptyQueue
	return
}

func (pq *pqs[T]) IsEmpty() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return len(pq.list) == 0
}

func (pq *pqs[T]) String() string {
	return fmt.Sprintf("%v", pq.list)
}

//================================

func (pq *pqs[T]) Push(e any) {
	pq.list = append(pq.list, e.(T))
}

func (pq *pqs[T]) Len() int {
	return len(pq.list)
}

func (pq *pqs[T]) Less(i, j int) bool {
	return pq.cmp(pq.list[i], pq.list[j])
}

func (pq *pqs[T]) Swap(i, j int) {
	pq.list[i], pq.list[j] = pq.list[j], pq.list[i]
}

func (pq *pqs[T]) Pop() (e any) {
	lidx := len(pq.list) - 1
	e = pq.list[lidx]
	var tmp T
	pq.list[lidx] = tmp
	pq.list = pq.list[:lidx]
	return
}
