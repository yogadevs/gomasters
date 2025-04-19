package cron

// implementing the wheel just for fun ...

type cronTasks []cronTask

func (q *cronTasks) Len() int {
	return len(*q)
}

func (q *cronTasks) Push(task cronTask) {
	*q = append(*q, task)
	i := len(*q) - 1

	for i > 0 {
		p := (i - 1) / 2
		if !(*q)[i].t.Before((*q)[p].t) {
			break
		}
		(*q)[i], (*q)[p] = (*q)[p], (*q)[i]
		i = p
	}
}

func (q *cronTasks) Pop() (cronTask, bool) {
	n := len(*q)
	if n == 0 {
		return cronTask{}, false
	}

	taskToReturn := (*q)[0]
	(*q)[0] = (*q)[n-1]
	*q = (*q)[:n-1]

	if len(*q) > 0 {
		q.sink(0)
	}

	return taskToReturn, true
}

func (q *cronTasks) Peek() (cronTask, bool) {
	if len(*q) == 0 {
		return cronTask{}, false
	}
	return (*q)[0], true
}

func (q *cronTasks) sink(i int) {
	n := len(*q)
	for {
		l := 2*i + 1
		r := 2*i + 2
		s := i

		if l < n && (*q)[l].t.Before((*q)[s].t) {
			s = l
		}

		if r < n && (*q)[r].t.Before((*q)[s].t) {
			s = r
		}

		if s == i {
			break
		}

		(*q)[i], (*q)[s] = (*q)[s], (*q)[i]
		i = s
	}
}
