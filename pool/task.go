package pool

import (
	"log"
	"time"
)

// Task 定义一个任务及其重试逻辑
type Task struct {
	fn        func() error
	retry     int
	retryWait time.Duration
}

// NewTask 创建任务
func NewTask(fn func() error, retry int, retryWait time.Duration) *Task {
	return &Task{
		fn:        fn,
		retry:     retry,
		retryWait: retryWait,
	}
}

// Run 执行任务并处理重试逻辑
func (t *Task) Run() {
	for {
		err := t.fn()
		if err != nil && t.retry > 0 {
			t.retry--
			log.Printf("Task failed, retrying in %s... (%d retries left)\n", t.retryWait, t.retry)
			time.Sleep(t.retryWait)
			t.retryWait *= 2 // 指数退避策略
		} else {
			if err != nil {
				log.Printf("Task failed permanently: %v\n", err)
			}
			break
		}
	}
}
