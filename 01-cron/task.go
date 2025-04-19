package cron

import "time"

type Task interface {
	Exec()
}

type cronTask struct {
	task Task
	t    time.Time
}
