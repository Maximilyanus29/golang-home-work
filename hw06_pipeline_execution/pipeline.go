package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	current := in

	for _, stage := range stages {
		out := make(Bi)

		go func(in In, out Bi, stage Stage, done In) {
			defer close(out)
			res := stage(in)

			for {
				select {
				case <-done:
					return
				case resValue, ok := <-res:
					if !ok {
						return
					}
					out <- resValue
				}
			}
		}(current, out, stage, done)

		current = out
	}

	return Out(current)
}
