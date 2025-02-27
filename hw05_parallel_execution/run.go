package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	lenTask := len(tasks)
	if lenTask < n {
		n = lenTask
	}

	taskChan := make(chan Task)
	// quitChan := make(chan struct{}, 1)

	var wg sync.WaitGroup
	// var once sync.Once

	var errCount int32
	// var shutdown int32

	// runGourutines := 0

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if atomic.LoadInt32(&errCount) >= int32(m) {
					return
				}
				if err := task(); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
		taskChan <- task
	}
	close(taskChan)

	wg.Wait()

	if atomic.LoadInt32(&errCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
