package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
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
		wg         sync.WaitGroup        // Waits for all workers to finish
		errCounter int64                 // Counts the number of errors (atomically)
		taskChan   = make(chan Task)     // Channel to send tasks to workers
		stopChan   = make(chan struct{}) // Used to signal workers to stop early
	)

	wg.Add(n) // Allocate n goroutines as workers
	for i := 0; i < n; i++ {
		go worker(i, taskChan, stopChan, &errCounter, maxErrors, &wg)
	}

	// Send tasks to the task channel
	for taskIndex, task := range tasks {
		select {
		case <-stopChan: // If stopCnan is closed (read error on attempt to read from stopChan) then it's "the signal":
			// Error threshold reached, stop sending tasks
			if Debug {
				fmt.Println("[Sender] Stop signal received on iteration ", taskIndex, ", ending task dispatch")
			}
		case taskChan <- task: // Task sent to workers OK
			if Debug {
				fmt.Println("[Sender] Task #", taskIndex, " sent to channel")
			}
		}
	}

	close(taskChan) // All tasks sent — close the channel
	if Debug {
		fmt.Println("[Sender] All tasks sent, taskChan closed")
		fmt.Println("........................................")
	}

	wg.Wait() // Wait for all workers to finish

	if Debug {
		fmt.Println("........................................")
		fmt.Println("err counter = ", atomic.LoadInt64(&errCounter))
	}

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
	maxErrorsCount int64,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	if Debug {
		fmt.Printf("[Worker %d] Started\n", id)
	}

	for task := range taskChan { // Process 'task' && automatically exit when 'taskChan' is closed.
		if err := task(); err != nil { // Execute the task & process the result
			newErrCount := atomic.AddInt64(errCounter, 1)
			// If we've reached the max allowed errors, stop everyone
			if newErrCount == maxErrorsCount {
				select {
				case <-stopChan:
					// Already signaled
				default:
					close(stopChan) // Fire the stop signal
				}
			}

			/* if Debug && newErrCount > maxErrorsCount {
				fmt.Printf("[Worker %d] Executed task but signal stop was called\n", id)
			} */
		}
	}
}
