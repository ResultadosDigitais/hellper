package job

import (
	"time"
)

type Job struct {
	Recurrence time.Duration
	Completed  chan time.Time
	Stop       chan time.Time
}

func NewTrackable(recurrence time.Duration, f func(Job), trackable bool) Job {
	job := Job{
		Recurrence: recurrence,
		Completed:  make(chan time.Time),
		Stop:       make(chan time.Time),
	}

	go func() {
		for {
			select {
			case <-job.Stop:
				return
			case <-time.After(recurrence):
				f(job)
				if trackable {
					job.Completed <- time.Now()
				}
				continue
			}
		}
	}()

	return job
}

func New(recurrence time.Duration, f func(Job)) Job {
	return NewTrackable(recurrence, f, false)
}

func Stop(job *Job) {
	job.Stop <- time.Now()
}
