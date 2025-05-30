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

func wrapStageWithCancel(stage Stage, done In, in In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		stageOut := stage(in)

		for {
			select {
			case <-done:
				// Drain the input channel to unblock the stage goroutine
				go func() {
					for i := range stageOut {
						// fmt.Println("drain out: ", i)
						_ = i // fight with linter
					}
				}()
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				// time.Sleep(111)
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}
