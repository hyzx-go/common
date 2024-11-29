package pool

import (
	"log"
	"sync"
)

// GoroutinePool 线程池结构
type GoroutinePool struct {
	workerCount int
	taskQueue   chan *Task
	wg          sync.WaitGroup
	stopChan    chan struct{}
}

var (
	poolInstance *GoroutinePool
	once         sync.Once
)

// NewGoroutinePool 创建线程池
func NewGoroutinePool(workerCount, queueSize int) *GoroutinePool {
	pool := &GoroutinePool{
		workerCount: workerCount,
		taskQueue:   make(chan *Task, queueSize),
		stopChan:    make(chan struct{}),
	}
	pool.startWorkers()
	return pool
}

// InitPool 初始化线程池（单例模式）
func InitPool(workerCount, queueSize int) {
	once.Do(func() {
		poolInstance = NewGoroutinePool(workerCount, queueSize)
	})
}

// GetPool 获取全局线程池实例
func GetPool() *GoroutinePool {
	if poolInstance == nil {
		log.Fatal("GoroutinePool not initialized. Call InitPool first.")
	}
	return poolInstance
}

// Submit 提交任务到线程池
func (p *GoroutinePool) Submit(task *Task) {
	p.wg.Add(1)
	p.taskQueue <- task
}

// Shutdown 释放线程池
func (p *GoroutinePool) Shutdown() {
	close(p.stopChan)
	p.wg.Wait()
	close(p.taskQueue)
	log.Println("GoroutinePool shutdown completed")
}

// startWorkers 启动协程
func (p *GoroutinePool) startWorkers() {
	for i := 0; i < p.workerCount; i++ {
		go func() {
			for {
				select {
				case task := <-p.taskQueue:
					p.executeTask(task)
				case <-p.stopChan:
					return
				}
			}
		}()
	}
}

// executeTask 执行任务
func (p *GoroutinePool) executeTask(task *Task) {
	defer p.wg.Done()
	task.Run()
}
