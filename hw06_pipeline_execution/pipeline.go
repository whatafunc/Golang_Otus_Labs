package hw06pipelineexecution

import "fmt"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 { // If no stages, just return the input channel
		return in
	}

	current := in // Connect stages in a pipeline
	fmt.Println("in has value:", current)
	for _, stage := range stages {
		// Wrap each stage with cancellation support & chan of `data` numbers

		//v1 - check done signal in a separate function for both read and output
		//current = wrapStageWithCancel(stage, done, current)

		//v2 - first filter input on done signal then send to processing
		// filteredIn := make(Bi)

		// // Filter `in` with respect to `done`
		// go func(in In, filteredIn Bi) {
		// 	defer close(filteredIn)
		// 	for {
		// 		select {
		// 		case <-done:
		// 			return
		// 		case val, ok := <-in:
		// 			fmt.Println("no cancel of reading input digit: val = ", val)
		// 			if !ok {
		// 				return
		// 			}
		// 			select {
		// 			case <-done:
		// 				return
		// 			case filteredIn <- val:
		// 			}
		// 		}
		// 	}
		// }(current, filteredIn)

		//v3 - no validation while reading
		current = wrapStageWithCancel(done, stage(current))

	}
	return current
}

// drains and discards all values from `in`.
func drain(in In) {
	for val := range in {
		_ = val // to comply with the linter
	}
}

func wrapStageWithCancel(done In, stageOut In) Out {

	out := make(Bi) // Process a stage unless it completed or cancelled
	go func() {
		defer close(out)
		defer drain(stageOut)
		for {
			select {
			case <-done:
				return
			case val, ok := <-stageOut:
				fmt.Printf("Processing value: %v \n", val)
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- val:
					fmt.Printf("sent value: %v \n", val)

				}
			}
		}
	}()

	return out
}
