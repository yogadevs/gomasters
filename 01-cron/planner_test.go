package cron

import (
	"sync"
	"testing"
	"time"
)

func TestPlanner_SingleTask(t *testing.T) {
	planner := NewPlanner()

	var wg sync.WaitGroup
	mock := &mockTask{wg: &wg}
	wg.Add(1)

	planner.Add(mock, time.Now().Add(100*time.Millisecond))

	if !waitForWaitGroup(&wg, 200*time.Millisecond) {
		t.Errorf("not in time")
	}
}

func TestPlanner_MultipleTasks(t *testing.T) {
	planner := NewPlanner()

	var wg sync.WaitGroup
	wg.Add(5)

	now := time.Now()
	planner.Add(&mockTask{wg: &wg}, now.Add(20*time.Millisecond))
	planner.Add(&mockTask{wg: &wg}, now.Add(10*time.Millisecond))
	planner.Add(&mockTask{wg: &wg}, now.Add(40*time.Millisecond))
	planner.Add(&mockTask{wg: &wg}, now.Add(50*time.Millisecond))
	planner.Add(&mockTask{wg: &wg}, now.Add(60*time.Millisecond))

	if !waitForWaitGroup(&wg, 100*time.Millisecond) {
		t.Errorf("not in time")
	}
}
