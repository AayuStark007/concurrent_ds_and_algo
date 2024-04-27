// Workers simulates a number of tasks and worker threads communicating via
// a message-passing channel. The channel is basically our Queue implementation.
// This examples serves as a demonstration of using our Queue as a Concurrent Channel.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aayustark007/concurrent_ds_and_algo/queue"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: workers [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type Task struct {
	Id      uint
	Payload string
}

type Worker struct {
	Id uint
}

func CreateTask(id uint, payload string) Task {
	return Task{
		id,
		payload,
	}
}

func CreateWorker(id uint) Worker {
	return Worker{
		id,
	}
}

func (w *Worker) ProcessTask(task Task) string {
	// time.Sleep(300 * time.Millisecond)
	return fmt.Sprintf("[Worker: {%d}] Processed Task %d::%s", w.Id, task.Id, task.Payload)
}

func SetupTasks(size uint) []Task {
	tasks := make([]Task, size)
	for i := uint(0); i < size; i++ {
		payload := fmt.Sprintf("TaskID: %d", i)
		tasks[i] = CreateTask(i, payload)
	}

	return tasks
}

func SetupWorkerThreads(numThreads uint, receiver queue.Queue[Task], sender queue.Queue[string], wg *sync.WaitGroup) {
	for i := uint(0); i < NTHREADS; i++ {
		worker := CreateWorker(i)
		wg.Add(1)

		go func(worker Worker) {
			defer wg.Done()
			for {
				task, ok := receiver.Dequeue()
				if !ok {
					return
				}

				result := worker.ProcessTask(task)
				sender.Enqueue(result)
			}
		}(worker)
	}
}

const (
	NTHREADS = 500
	NTASK    = 1_000_000
)

var (
	blocking = flag.Bool("blocking", false, "Use Blocking Queue")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("workers: ")

	// Parse Flags
	flag.Usage = usage
	flag.Parse()

	// begin
	var taskQ queue.Queue[Task]
	var resultQ queue.Queue[string]

	if *blocking {
		log.Printf("Start processing with %d tasks on %d workers using blocking channels\n", NTASK, NTHREADS)
		taskQ = queue.NewBlockingQueue[Task]()
		resultQ = queue.NewBlockingQueue[string]()
	} else {
		log.Printf("Start processing with %d tasks on %d workers using non-blocking channels\n", NTASK, NTHREADS)
		taskQ = queue.NewNonBlockingQueue[Task]()
		resultQ = queue.NewNonBlockingQueue[string]()
	}

	tasks := SetupTasks(NTASK)
	start := time.Now()
	defer func() {
		end := time.Since(start)
		log.Printf("Took: %dus | %dms", end.Microseconds(), end.Milliseconds())
	}()

	var wg sync.WaitGroup
	SetupWorkerThreads(NTHREADS, taskQ, resultQ, &wg)

	for _, task := range tasks {
		taskQ.Enqueue(task)
	}

	resultCount := 0
	for {
		_, ok := resultQ.Dequeue()
		if !ok {
			if resultCount == len(tasks) {
				break
			}
			continue
		}
		// log.Println(result)
		resultCount++
	}

	wg.Wait()
}
