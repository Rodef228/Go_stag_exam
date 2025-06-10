package main

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	workerID int
	tasks    chan string
	quit     chan struct{}
	wg       sync.WaitGroup
	mu       sync.Mutex
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		tasks: make(chan string),
		quit:  make(chan struct{}),
	}
}

func (wp *WorkerPool) AddWorker() {
	wp.mu.Lock()
	wp.workerID++
	id := wp.workerID
	wp.mu.Unlock()

	wp.wg.Add(1)
	go func() {
		defer wp.wg.Done()
		for {
			select {
			case task, ok := <-wp.tasks:
				if !ok {
					return
				}
				fmt.Printf("Worker %d: %s\n", id, task)
			case <-wp.quit:
				return
			}
		}
	}()
}

func (wp *WorkerPool) RemoveWorker() {
	wp.quit <- struct{}{}
	wp.wg.Done()
}

func (wp *WorkerPool) Submit(task string) {
	wp.tasks <- task
}

func (wp *WorkerPool) Close() {
	close(wp.tasks)
	close(wp.quit)
	wp.wg.Wait()
}

func main() {
	pool := NewWorkerPool()

	for i := 0; i < 3; i++ {
		pool.AddWorker()
	}

	go func() {
		for i := 0; i < 10; i++ {
			pool.Submit(fmt.Sprintf("Task %d", i))
			time.Sleep(10 * time.Millisecond)
		}
	}()

	time.Sleep(500 * time.Millisecond)
	pool.AddWorker()
	time.Sleep(500 * time.Millisecond)
	pool.RemoveWorker()

	time.Sleep(1 * time.Second)
	pool.Close()
}
