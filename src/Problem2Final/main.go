/*

interface JobScheduler {
  void schedule(Runnable job, long delayMs);
}


Runnable {
	name string
	start int64
	delay int64
	index int
}


priority index=0
pq[i]

Producer

pq -> 1,2,3,4

Consumer

pq[0] -> time.Now().Unix()>= (pq[0].start+pq[0].delay)
fmt.Println(pq[0].name + )


(pq[i].start+pq[i].delay)

*/
package Confluent

import (
	"container/heap"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Job struct {
	name string
	start int64
	delay int64
}

func NewJob(name string, start, end int64) *Job{
	job := new(Job)
	job.name=name
	job.start=start
	job.delay=end
	return job
}

type Schedule []*Job

func (pq Schedule) Len() int {
	return pq.Len()
}

func (pq Schedule) Less(i, j int) bool{
	return (pq[i].start+pq[i].delay)<(pq[j].start+pq[j].delay)
}

func (pq Schedule) Swap(i, j int) {
	pq[i],pq[j]=pq[j],pq[i]
}

func (pq *Schedule) Pop() interface{} {
	old := *pq
	len := old.Len()
	if(len==0) {
		return nil
	}
	elem := old[len-1]
	*pq=old[0:len-1]
	return elem
}

func (pq *Schedule) Push(pushJob interface{}) {
	job := pushJob.(*Job)
	pq.Push(job)
}

var cond sync.Cond
var allDone bool

func Producer(pq *Schedule, wg *sync.WaitGroup) {
	defer wg.Done()
	for i:=0; i<10; i++ {
    heap.Push(pq, NewJob("cmd"+strconv.Itoa(i), time.Now().Unix(), (int64(i*1000))))
		cond.Signal()
	}
	// cond.Broadcast()
	allDone=true
}



func Consumer(pq *Schedule, wg *sync.WaitGroup) {
  defer wg.Done()
	for {
	 cond.L.Lock()
	 fmt.Println(pq.Len())
	 if(pq.Len()<=0) {
		 if(allDone) {
			 break
		 } else {
			 continue
		 }
	 }
	 top := (*pq)[0]
	 for(time.Now().Unix() < (top.start+top.delay)) {
	 	 ticket := time.NewTicker(time.Duration(top.start+top.delay-time.Now().Unix())*time.Second)
		 go func() {
			 <- ticket.C
			 cond.Signal()
		 }()
		cond.Wait()
		top = (*pq)[0]
	}
	 elem := heap.Pop(pq).(*Job)
	 cond.L.Unlock()
	 fmt.Println(elem.name)
 }
}

func Run() {
  pq := new(Schedule)
  heap.Init(pq)

  cond.L = new(sync.Mutex)

  wg := new(sync.WaitGroup)
  wg.Add(2)

  go Producer(pq, wg)
  go Consumer(pq, wg)

  wg.Wait()
}
