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

func wrapStageWithCancel(stage Stage, done In, current In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		stageOut := stage(current)

		for {
			select {
			case <-done: // check before starting to process the next value (func)
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case <-done: // check again but before sending a value to the next stage
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}
