package cron

import (
	"sync"
	"time"
)

type mockTask struct {
	wg *sync.WaitGroup
}

func (m *mockTask) Exec() {
	if m.wg != nil {
		m.wg.Done()
	}
}

func waitForWaitGroup(wg *sync.WaitGroup, timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}
