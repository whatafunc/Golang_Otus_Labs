package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task func() error

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

const Debug = false // Debug logging

// Start up to n workers to process the tasks slice.
func Run(tasks []Task, n, m int) error {
	if n <= 0 { // Workers Count should not be 0
		return nil
	}

	if m < 0 { // Any Error should stop All operations
		m = 0
	}

	maxErrors := int64(m)

	var (
		wg              sync.WaitGroup        // Waits for all workers to finish
		errCounter      int64                 // Counts the number of errors (atomically)
		taskChan        = make(chan Task)     // Channel to send tasks to workers
		stopChan        = make(chan struct{}) // Used to signal workers to stop early
		dispatchedTasks int64                 // Counter of successfully dispatched Tasks[]
		completedTasks  int64                 // Counter of successfully completed Tasks[]
	)

	if maxErrors > 0 && maxErrors < int64(len(tasks)) {
		wg.Add(1) // main routine to MOnitor the status of the Errors the number value
		go func() {
			defer wg.Done()
			for {
				if atomic.LoadInt64(&errCounter) >= maxErrors {
					select {
					case <-stopChan: // already closed
					default:
						close(stopChan)
					}
					return
				}

				// Normal stop: all dispatched tasks completed
				if atomic.LoadInt64(&completedTasks) == atomic.LoadInt64(&dispatchedTasks) {
					return // graceful exit
				}

				time.Sleep(10 * time.Millisecond) // Prevent busy looping
			}
		}()
	}

	wg.Add(n) // Allocate n goroutines as workers
	for i := 0; i < n; i++ {
		go worker(i, taskChan, stopChan, &errCounter, &wg, &completedTasks)
	}

	// Send tasks to the task channel
	for taskIndex, task := range tasks {
		// Fast path check before select (atomic)
		if maxErrors > 0 && atomic.LoadInt64(&errCounter) >= maxErrors {
			break
		}

		select {
		case <-stopChan: // If stopCnan is closed (read error on attempt to read from stopChan) then it's "the signal":
			// Error threshold reached, stop sending tasks
			_ = taskIndex
		case taskChan <- task: // Task sent to workers OK
			atomic.AddInt64(&dispatchedTasks, 1)
		}
	}

	close(taskChan) // All tasks sent — close the channel

	wg.Wait() // Wait for all workers to finish

	if errCounter > maxErrors { // Return error if the errors limit was reached
		return ErrErrorsLimitExceeded
	}
	return nil
}

// Each worker is a separate concurrent functionality excuted.
// Stops early if m errors occur across all tasks.
func worker(
	id int,
	taskChan <-chan Task,
	stopChan chan struct{},
	errCounter *int64,
	wg *sync.WaitGroup,
	completedTasks *int64,
) {
	defer wg.Done()

	if Debug {
		fmt.Printf("[Worker %d] Started\n", id)
	}

	for task := range taskChan { // Process 'task' && automatically exit when 'taskChan' is closed.
		select {
		case <-stopChan:
			return
		default:
			if err := task(); err != nil {
				atomic.AddInt64(errCounter, 1)
			}
			atomic.AddInt64(completedTasks, 1)
		}
	}
}
