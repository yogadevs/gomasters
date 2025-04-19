package cron

import (
	"sync"
	"time"
)

type Planner struct {
	mu          sync.Mutex
	cronTasks   cronTasks
	taskChannel chan cronTask
	timer       *time.Timer
}

func NewPlanner() *Planner {
	p := &Planner{
		taskChannel: make(chan cronTask, 1),
	}

	go p.run()
	return p
}

func (p *Planner) Add(task Task, t time.Time) {
	p.taskChannel <- cronTask{task: task, t: t}
}

func (p *Planner) run() {
	for {
		timerCh := p.setNextTimer()

		select {
		case newTask := <-p.taskChannel:
			p.addTask(newTask)
		case <-timerCh:
			p.runExpired()
		}
	}
}

func (p *Planner) addTask(newTask cronTask) {
	p.mu.Lock()
	p.cronTasks.Push(newTask)
	p.mu.Unlock()
}

func (p *Planner) setNextTimer() <-chan time.Time {
	nextTime, ok := p.getNextTaskTime()
	p.stopTimer()

	if !ok {
		return nil
	}

	duration := time.Until(nextTime)
	if duration < 0 {
		duration = 0
	}

	p.timer = time.NewTimer(duration)
	return p.timer.C
}

func (p *Planner) runExpired() {
	p.timer = nil

	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for p.cronTasks.Len() > 0 {
		nextTask, ok := p.cronTasks.Peek()
		if !ok || nextTask.t.After(now) {
			break
		}

		taskToRun, _ := p.cronTasks.Pop()

		p.mu.Unlock()
		taskToRun.task.Exec()
		p.mu.Lock()
	}
}

func (p *Planner) getNextTaskTime() (time.Time, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	nextTask, ok := p.cronTasks.Peek()
	if !ok {
		return time.Time{}, false
	}
	return nextTask.t, true
}

func (p *Planner) stopTimer() {
	if p.timer != nil {
		if !p.timer.Stop() {
			select {
			case <-p.timer.C:
			default:
			}
		}
		p.timer = nil
	}
}
