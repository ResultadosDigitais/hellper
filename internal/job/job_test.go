package job

import (
	"testing"
	"time"
)

func TestStartJob(t *testing.T) {
	completed := false
	f := func(j Job) {
		completed = true
	}

	job := NewTrackable(time.Millisecond, f, true)

	select {
	case <-job.Completed:
		if !completed {
			t.Fatal("Job not executed with success")
		}
	case <-time.After(time.Millisecond * 10):
		if completed {
			t.Fatal("Job time out")
		} else {
			t.Fatal("Job time out and not completed")
		}
	}
}

func TestRecurrencyJob(t *testing.T) {
	job := NewTrackable(time.Millisecond, func(j Job) {}, true)

	timeOne := <-job.Completed
	timeTwo := <-job.Completed

	if timeOne == timeTwo {
		t.Fatal("Job not executed with recurrency")
	}
}

func TestStopJob(t *testing.T) {
	completed := false
	job := NewTrackable(time.Minute, func(j Job) { completed = false }, true)
	job.Recurrence = time.Millisecond
	Stop(&job)
	completed = true

	time.Sleep(time.Millisecond * 100)
	if !completed {
		t.Fatal("Job not stopped")
	}
}
