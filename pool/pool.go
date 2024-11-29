package pool

import (
	"log"
	"runtime/debug"
	"sync"
	"time"
)

// Task 定义任务接口
type Task interface {
	Run() error
}

// RetryTask 包含重试机制的任务结构
type RetryTask struct {
	Task       Task
	RetryCount int
	RetryDelay time.Duration
}

// GoroutinePool 线程池结构
type GoroutinePool struct {
	taskQueue   chan *RetryTask
	wg          sync.WaitGroup
	closeSignal chan struct{}
}

// NewGoroutinePool 创建 GoroutinePool
func NewGoroutinePool(workerCount, taskQueueSize int) *GoroutinePool {
	pool := &GoroutinePool{
		taskQueue:   make(chan *RetryTask, taskQueueSize),
		closeSignal: make(chan struct{}),
	}

	// 启动 worker
	for i := 0; i < workerCount; i++ {
		go pool.worker()
	}
	return pool
}

// Submit 提交任务到线程池
func (p *GoroutinePool) Submit(task Task, retryCount int, retryDelay time.Duration) {
	p.wg.Add(1)
	p.taskQueue <- &RetryTask{
		Task:       task,
		RetryCount: retryCount,
		RetryDelay: retryDelay,
	}
}

// worker 执行任务的协程
func (p *GoroutinePool) worker() {
	for {
		select {
		case <-p.closeSignal:
			return
		case retryTask := <-p.taskQueue:
			p.executeTask(retryTask)
		}
	}
}

// executeTask 执行任务并支持重试和异常处理
func (p *GoroutinePool) executeTask(retryTask *RetryTask) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in task: %v\nStack trace:\n%s", r, debug.Stack())
		}
		p.wg.Done()
	}()

	for attempt := 0; attempt <= retryTask.RetryCount; attempt++ {
		err := retryTask.Task.Run()
		if err != nil {
			log.Printf("Task failed on attempt %d: %v", attempt+1, err)
			if attempt < retryTask.RetryCount {
				delay := retryTask.RetryDelay * (1 << attempt) // 指数退避
				log.Printf("Retrying task in %s...", delay)
				time.Sleep(delay)
			} else {
				log.Printf("Task failed after maximum retries.")
			}
		} else {
			return
		}
	}
}

// Close 关闭线程池并等待所有任务完成
func (p *GoroutinePool) Close() {
	close(p.closeSignal)
	p.wg.Wait()
	close(p.taskQueue)
}
