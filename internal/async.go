package srm

import (
	"sync"
)

type taskError struct {
	Id  int
	Err error
}

/*
Runs given tasks concurrently using a WaitGroup

Returns the a map where the index of the ran task is the key
and the value is the returned error itself
*/
func RunConcurrently(tasks ...func() error) map[int]error {
	var wg sync.WaitGroup
	errch := make(chan taskError)

	// Adds tasks in waitgroup
	for i, action := range tasks {
		wg.Add(1)

		go func(i int, action func() error) {
			defer wg.Done()
			err := action()

			if err != nil {
				errch <- taskError{Id: i, Err: err}
			}
		}(i, action)
	}

	// Wait for the group to be done and closes the error channel
	go func() {
		wg.Wait()
		close(errch)
	}()

	// Creates the return map
	errs := make(map[int]error)
	for err := range errch {
		errs[err.Id] = err.Err
	}

	return errs
}
