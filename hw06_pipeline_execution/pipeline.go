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
		// Wrap each stage with cancellation support
		current = wrapStageWithCancel(stage, done, current)
	}
	return current
}

// drain starts a goroutine that simply drains and discards all values from `in`.
func drain(in In) {
	go func() {
		for v := range in {
			_ = v
		}
	}()
}

func wrapStageWithCancel(stage Stage, done In, in In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		stageOut := stage(in)

		for {
			select {
			case <-done:
				drain(stageOut) // empty the unused values on `Stop`` signal
				return
			case val, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case <-done:
					drain(stageOut) // empty the unused values on `Stop`` signal
					return
				case out <- val:
				}
			}
		}
	}()
	return out
}
