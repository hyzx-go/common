package pool

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestGoroutinePool(t *testing.T) {
	// 初始化全局线程池
	InitPool(5, 10)
	pool := GetPool()

	// 提交普通任务
	task1 := NewTask(func() error {
		log.Println("Executing task1")
		return nil
	}, 0, 0)
	pool.Submit(task1)

	// 提交会失败的任务并重试
	task2 := NewTask(func() error {
		log.Println("Executing task2 - simulate failure")
		return errors.New("task2 failed")
	}, 3, 2*time.Second)
	pool.Submit(task2)

	// 测试定时任务
	scheduler := NewScheduler(pool)
	counter := 0
	scheduler.Schedule(1*time.Second, func() error {
		log.Printf("Executing scheduled task: count %d", counter)
		counter++
		if counter >= 5 {
			return errors.New("Stopping scheduled task after 5 executions")
		}
		return nil
	})

	// 等待任务完成
	time.Sleep(10 * time.Second)
	pool.Shutdown()

	// 验证
	if counter < 5 {
		t.Errorf("Scheduled task did not run enough times: %d", counter)
	}
}
