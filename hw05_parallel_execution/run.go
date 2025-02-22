package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		m = len(tasks) + 1
	}

	taskChan := make(chan Task)
	quitChan := make(chan struct{})

	var wg sync.WaitGroup
	var mu sync.Mutex

	var errCount int

	go func() {
		if errCount >= m {
			close(quitChan)
			return
		}
	}()

	// Запускаем n воркеров
	for i := 0; i < n; i++ {
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
