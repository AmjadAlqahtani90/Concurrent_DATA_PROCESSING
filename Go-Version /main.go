package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	totalTasks  = 10
	workerCount = 3
)

// TaskQueue concurrency-safe queue using channel.
type TaskQueue struct {
	tasks chan int
}

func NewTaskQueue(size int) *TaskQueue {
	return &TaskQueue{tasks: make(chan int, size)}
}
// input
func (q *TaskQueue) addTask(task int) {
	q.tasks <- task
}

// output
func (q *TaskQueue) getTask() (task int, ok bool) {
	task, ok = <-q.tasks
	return
}

func (q *TaskQueue) close() {
	close(q.tasks)
}

type ResultStore struct {
	mu      sync.Mutex
	results []string
}

func (s *ResultStore) add(result string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results = append(s.results, result)
}

func (s *ResultStore) getResults() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.results...)
}

func worker(id int, queue *TaskQueue, store *ResultStore, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Worker %d started.", id)
	defer log.Printf("Worker %d finished.", id)

	for {
		task, ok := queue.getTask()
		if !ok { 
			return
		}

		log.Printf("Worker %d processing Task %d", id, task)
		time.Sleep(500 * time.Millisecond)
		result := task * 2

		store.add(fmt.Sprintf("Task %d -> Result %d", task, result))
		log.Printf("Worker %d completed Task %d", id, task)
	}
}

func saveResults(filename string, results []string) (returnErr error) {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create results file: %w", err)
	}
	defer func() {
		if err := file.Close(); returnErr == nil && err != nil {
			returnErr = fmt.Errorf("close results file: %w", err)
		}
	}()

	writer := bufio.NewWriter(file)
	for _, result := range results {
		if _, err := fmt.Fprintln(writer, result); err != nil {
			return fmt.Errorf("write result: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush results file: %w", err)
	}
	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	queue := NewTaskQueue(totalTasks)
	store := &ResultStore{}

	for task := 1; task <= totalTasks; task++ {
		queue.addTask(task)
	}
	queue.close()

	var wg sync.WaitGroup
	for id := 1; id <= workerCount; id++ {
		wg.Add(1)
		go worker(id, queue, store, &wg)
	}
	wg.Wait()

	results := store.getResults()
	if len(results) != totalTasks {
		log.Fatalf("Expected %d results but received %d", totalTasks, len(results))
	}

	fmt.Println("\nFinal Results:")
	for _, result := range results {
		fmt.Println(result)
	}

	if err := saveResults("results.txt", results); err != nil {
		log.Printf("Error saving results: %v", err)
		os.Exit(1)
	}
	fmt.Println("\nResults saved to results.txt")
	fmt.Println("Data Processing System finished.")
}
