package main

import (
  "container/heap"
  "fmt"
  "sync"
  "time"
)


// Logic:
// Create a priority queue based on expiry
// Implement Len, Less, Push, Pop and Swap methods of heap
// Insert elements into priority queue using Push()
// Spawn goroutines to check and spawn jobs
// checkJobExpiry checks the first job from the priority queue and checks
// if the unix timestamp has expired. If so, broadcasts on the condition
// variable.
// Go routines waiting on condition variables check for the condition
// and Pop() the job from priority queue and run the job.
// heap.Init() satisfies heap invariant.

const (
  NUMPARALLELJOBS int = 100
)

var lock sync.Mutex
var cond *sync.Cond
var runJob bool

type Job struct {
  cmdName string
  args []string
  expiry int64

  index int
}

type PriorityJobs []*Job

func (pj PriorityJobs) Len() int {
  return len(pj)
}

func (pj PriorityJobs) Less(i, j int) bool {
  return pj[i].expiry < pj[j].expiry
}

func (pj *PriorityJobs) Pop() interface{} {
  item := (*pj)[len(*pj) -1]
  item.index = -1
  *pj = (*pj)[0:len(*pj)-1]
  return item
}

func (pj *PriorityJobs) Push(x interface{}) {
  n := len(*pj)
  item := x.(*Job)
  item.index = n
  *pj = append(*pj, item)
}

func (pj PriorityJobs) Swap(i,j int) {
  pj[i], pj[j] = pj[j], pj[i]
  pj[i].index = i
  pj[j].index = j
}

// Check if there a job every 1 minute and signal the goroutine
func (pj PriorityJobs) checkJobExpiry(wg sync.WaitGroup) {
  defer wg.Done()
  for {
    lock.Lock()
    // Check if the condition expired
    t := time.Now().Unix()
    if t > pj[0].expiry {
      // Timer expired
      runJob = true
      cond.Broadcast()
    }
    lock.Unlock()
    time.Sleep(1 * time.Minute)
  }
}

func (pj PriorityJobs) runJobRoutine(wg sync.WaitGroup) {
  defer wg.Done()
  for {
    lock.Lock()
    // Check if the condition occurred
    for runJob != true {
      cond.Wait()
    }
    item := pj.Pop()
    func () {
      // Do nothing
      fmt.Println(item.(Job).cmdName, item.(Job).args)
    }()
    runJob = false
    lock.Unlock()
  }
}

func main() {
  cond = sync.NewCond(&lock)
  var wg sync.WaitGroup
  var pj PriorityJobs

  heap.Init(&pj)

  wg.Add(1)
  go pj.checkJobExpiry(wg)
  for i:=0; i< NUMPARALLELJOBS; i++ {
    go pj.runJobRoutine(wg)
    wg.Add(1)
  }

  wg.Wait()
}
