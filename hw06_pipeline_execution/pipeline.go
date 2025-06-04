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

	for _, stage := range stages {
		// Wrap each stage with cancellation support & chan of `data` numbers
		current = wrapStageWithCancel(stage, done, current)
	}
	return current
}

// drains and discards all values from `in`.
func drain(in In) {
	for val := range in {
		_ = val // to comply with the linter
	}
}

func wrapStageWithCancel(stage Stage, done In, in In) Out {
	filteredIn := make(Bi)

	// Filter `in` with respect to `done`
	go func() {
		defer close(filteredIn)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case filteredIn <- val:
				}
			}
		}
	}()

	stageOut := stage(filteredIn)

	out := make(Bi) // Process a stage unless it completed or cancelled
	go func() {
		defer close(out)
		defer drain(stageOut)
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
