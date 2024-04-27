package main

import (
	"fmt"
	"sync"

	"github.com/aayustark007/concurrent_ds_and_algo/queue"
)

const (
	NTHREADS = 1000
	NTASK    = 1_000_000
)

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

func main() {
	fmt.Printf("Start processing with %d tasks on %d workers", NTASK, NTHREADS)

	taskQ := queue.NewNonBlockingQueue[Task]()
	resultQ := queue.NewNonBlockingQueue[string]()

	tasks := SetupTasks(NTASK)

	var wg sync.WaitGroup
	SetupWorkerThreads(NTHREADS, taskQ, resultQ, &wg)

	for _, task := range tasks {
		taskQ.Enqueue(task)
	}

	for {
		result, ok := resultQ.Dequeue()
		if !ok {
			fmt.Println("Receiver: no more messages")
			break
		}
		fmt.Println(result)
	}

	wg.Wait()
}
