package pool

import (
	"time"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	pool *GoroutinePool
}

// NewScheduler 创建调度器
func NewScheduler(pool *GoroutinePool) *Scheduler {
	return &Scheduler{pool: pool}
}

// Schedule 定期执行任务
func (s *Scheduler) Schedule(interval time.Duration, fn func() error) {
	go func() {
		for {
			select {
			case <-s.pool.stopChan:
				return
			case <-time.After(interval):
				task := NewTask(fn, 0, 0)
				s.pool.Submit(task)
			}
		}
	}()
}
