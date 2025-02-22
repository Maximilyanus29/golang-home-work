package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		m = len(tasks) + 1
	}

	taskChan := make(chan Task)
	quitChan := make(chan struct{}, 1)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var quitOnce sync.Once

	var errCount int

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-quitChan:
					return
				case task, ok := <-taskChan:
					if !ok {
						return
					}
					if err := task(); err != nil {
						mu.Lock()
						errCount++
						if errCount >= m {
							mu.Unlock()
							quitOnce.Do(func() {
								close(quitChan)
							})
							return
						}
						mu.Unlock()
					}
				}
			}
		}()
	}

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil

}
