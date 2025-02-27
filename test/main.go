package main

import "sync"

func main() {
	var wg sync.WaitGroup

	quitCh := make(chan struct{})

	// _, ok = <-quitCh
	// print(ok)
	wg.Add(3)

	go func() {
		wg.Done()
		close(quitCh)

	}()

	go func() {
		wg.Done()

		_, ok := <-quitCh
		print(ok)
	}()

	go func() {
		wg.Done()

		_, ok := <-quitCh
		print(ok)
	}()

	wg.Wait()

}
