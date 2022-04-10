package Problem2

import (
	"container/heap"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var lock sync.Mutex
var allDone bool

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return (pq[i].start+pq[i].delay) < (pq[j].start+pq[j].delay)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Task)
	item.index=n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	if(n==0) {
		return nil
	}
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(task *Task, name string, start int64, delay int64) {
	task.name = name
	task.start = start
	task.delay = delay
	heap.Fix(pq, task.index)
}

func Run() {
	pq := new(PriorityQueue)
	heap.Init(pq)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go Producer(pq, wg)
	for i:=0; i<5; i++ {
		wg.Add(1)
		go Consumer(pq, wg)
	}
	wg.Wait()
}

func Producer(pq *PriorityQueue, wg *sync.WaitGroup) {
	defer wg.Done()
	for i:=0; i<10; i++ {
		lock.Lock()
		heap.Push(pq, NewTaskWithoutIndex(
			"cmd"+strconv.Itoa(i), []string{}, time.Now().Unix(), int64(10-1*i),
		))
		lock.Unlock()
	}
	allDone = true
}

func Consumer(pq *PriorityQueue, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		lock.Lock()
		if(pq.Len() <= 0) {
			lock.Unlock()
			if(allDone) {
				break
			} else {
				continue
			}
		}
		top := (*pq)[0]
		if(time.Now().Unix() < (top.start+top.delay)) {
			lock.Unlock()
			continue
		}
		item := heap.Pop(pq).(*Task)
		lock.Unlock()
		fmt.Println(item.name + " " + strconv.FormatInt(item.start + item.delay, 10) + " " + strconv.FormatInt(time.Now().Unix(), 10) + "\n")
		time.Sleep(10*time.Second)
	}
}
