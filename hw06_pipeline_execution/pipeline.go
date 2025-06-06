package hw06pipelineexecution

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
	current = makeCancellableStage(done, current)
	for _, stage := range stages {
		// make stage with cancellation support & chan of `data` numbers
		current = makeCancellableStage(done, stage(current))
	}
	return current
}

// drains and discards all values from `in`.
func drain(in In) {
	for val := range in {
		_ = val // to comply with the linter
	}
}

func makeCancellableStage(done In, stageOut In) Out {
	out := make(Bi) // Process a stage unless it completed or cancelled
	go func() {
		defer drain(stageOut) // 2
		defer close(out)      // 1

		for {
			select {
			case <-done:
				return
			case val, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- val:
				}
			}
		}
	}()

	return out
}
