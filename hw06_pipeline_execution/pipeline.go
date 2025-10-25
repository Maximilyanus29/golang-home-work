package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	var cur Out
	cur = in

	go func() {
		for {
			_, ok := <-done
			if !ok {
				for range in { //nolint
				}

				return
			}
		}
	}()

	if done != nil {
		_, ok := <-done
		if !ok {
			for range in { //nolint
			}

			return cur
		}
	}

	for _, stage := range stages {
		cur = stage(createWrapForStageValue(cur, done, &wg))
	}

	wg.Wait()

	return cur
}

func createWrapForStageValue(in In, done In, wg *sync.WaitGroup) Out {
	out := make(Bi)

	if done == nil {
		return in
	}

	wg.Add(1)
	go func(in In, out Bi, done In) {
		defer func() {
			wg.Done()
			close(out)
		}()

		for {
			select {
			case <-done:
				for range in { //nolint
				}

				return
			case invalue, ok := <-in:
				if !ok {
					return
				}

				out <- invalue
			}
		}
	}(in, out, done)

	return Out(out)
}
