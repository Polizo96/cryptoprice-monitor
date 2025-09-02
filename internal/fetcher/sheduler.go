package fetcher

import "time"

type Scheduler struct {
	interval time.Duration
	stop     chan struct{}
}

func NewScheduler(intervalSeconds int) *Scheduler {
	return &Scheduler{
		interval: time.Duration(intervalSeconds) * time.Second,
		stop:     make(chan struct{}),
	}
}

func (s *Scheduler) Start(task func()) {
	ticker := time.NewTicker(s.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				task()
			case <-s.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}
