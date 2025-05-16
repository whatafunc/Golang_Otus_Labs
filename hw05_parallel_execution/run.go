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
// Stops early if m errors occur across all tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 { // Workers Count should not be 0
		return nil
	}

	if m < 0 { // Error should stop All operations
		m = 0
	}

	// fix linter issue:
	// int (which might be 64-bit) → int32 conversion could silently truncate large numbers.
	//nolint:gosec // m will never exceed int32 in our context
	maxErrors := int32(m)

	var (
		wg         sync.WaitGroup        // Waits for all workers to finish
		errCounter int32                 // Counts the number of errors (atomically)
		taskChan   = make(chan Task)     // Channel to send tasks to workers
		stopChan   = make(chan struct{}) // Used to signal workers to stop early
	)

	wg.Add(n) // Allocate n goroutines as workers
	for i := 0; i < n; i++ {
		go worker(i, taskChan, stopChan, &errCounter, maxErrors, &wg)
	}

	// Send tasks to the task channel
	go func() {
		for taskIndex, task := range tasks {
			select {
			case <-stopChan: // If stopCnan is closed (read error on attempt to read from stopChan) then it's "the signal":
				// Error threshold reached, stop sending tasks
				if Debug {
					fmt.Println("[Sender] Stop signal received on iteration ", taskIndex, ", ending task dispatch")
				}
				return
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
	}()

	wg.Wait() // Wait for all workers to finish

	if Debug {
		fmt.Println("........................................")
		fmt.Println("err counter = ", atomic.LoadInt32(&errCounter))
	}

	if atomic.LoadInt32(&errCounter) > maxErrors { // Return error if the errors limit was reached
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(
	id int,
	taskChan <-chan Task,
	stopChan chan struct{},
	errCounter *int32,
	maxErrorsCount int32,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	if Debug {
		fmt.Printf("[Worker %d] Started\n", id)
	}

	for {
		select {
		case <-stopChan:
			// Stop signal received — exit worker since stopChan got closed.
			if Debug {
				fmt.Printf("[Worker %d] Received stop signal since stopChan got a record\n", id)
			}
			return
		case task, ok := <-taskChan:
			if !ok {
				// taskChan is closed, no more tasks
				if Debug {
					fmt.Printf("[Worker %d] Task channel closed, so no more workes — exiting\n", id)
				}
				return
			}
			if Debug {
				fmt.Printf("[Worker %d] Executing task from taskChan...\n", id)
			}

			// Execute the task & process the result
			if err := task(); err != nil {
				newErrCount := atomic.AddInt32(errCounter, 1)
				// If we've reached the max allowed errors, stop everyone
				if newErrCount >= maxErrorsCount {
					select {
					case <-stopChan:
						// Already signaled
					default:
						close(stopChan)
					}
				}
			}
		}
	}
}
