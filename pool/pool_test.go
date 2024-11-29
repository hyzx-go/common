package pool

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// MockTask 定义用于测试的任务
type MockTask struct {
	FailTimes    int           // 需要失败的次数
	RetryCounter *int32        // 重试计数
	RunCounter   *int32        // 总运行计数
	RunDuration  time.Duration // 模拟任务运行时间
}

// Run 实现 Task 接口，包含失败逻辑
func (m *MockTask) Run() error {
	atomic.AddInt32(m.RunCounter, 1)
	time.Sleep(m.RunDuration) // 模拟任务耗时

	if atomic.AddInt32(m.RetryCounter, -1) >= 0 {
		return errors.New("mock task failed")
	}
	return nil
}

// TestGoroutinePool 测试线程池的任务执行和重试机制
func TestGoroutinePool(t *testing.T) {
	const (
		workerCount      = 3
		taskQueueSize    = 10
		retryCount       = 3
		retryDelay       = 100 * time.Millisecond
		mockFailTimes    = 2
		expectedRunTimes = mockFailTimes + 1 // 最终任务运行次数 = 失败次数 + 成功次数
	)

	// 初始化线程池
	pool := NewGoroutinePool(workerCount, taskQueueSize)
	defer pool.Close()

	var retryCounter int32 = mockFailTimes
	var runCounter int32 = 0

	// 创建 MockTask
	mockTask := &MockTask{
		FailTimes:    mockFailTimes,
		RetryCounter: &retryCounter,
		RunCounter:   &runCounter,
		RunDuration:  50 * time.Millisecond,
	}

	// 提交任务
	pool.Submit(mockTask, retryCount, retryDelay)

	// 等待任务完成
	pool.wg.Wait()

	// 验证任务总运行次数
	if atomic.LoadInt32(&runCounter) != expectedRunTimes {
		t.Errorf("expected task to run %d times, but got %d", expectedRunTimes, runCounter)
	}
}

// TestTaskTimeout 测试任务超时
func TestTaskTimeout(t *testing.T) {
	const (
		workerCount   = 2
		taskQueueSize = 5
		retryCount    = 1
		retryDelay    = 100 * time.Millisecond
		timeout       = 100 * time.Millisecond
	)

	pool := NewGoroutinePool(workerCount, taskQueueSize)
	defer pool.Close()

	var completed int32 = 0

	// 创建一个超时任务
	timeoutTask := &MockTask{
		RetryCounter: new(int32), // 不重试，直接模拟超时
		RunCounter:   &completed,
		RunDuration:  timeout * 2, // 模拟任务耗时大于超时阈值
	}

	// 提交任务
	pool.Submit(timeoutTask, retryCount, retryDelay)

	// 等待任务完成
	pool.wg.Wait()

	// 验证任务是否被执行
	if atomic.LoadInt32(&completed) == 0 {
		t.Errorf("task should have been attempted at least once")
	}
}

// TestPoolClose 测试线程池关闭逻辑
func TestPoolClose(t *testing.T) {
	pool := NewGoroutinePool(2, 5)

	var executed int32 = 0

	// 提交多个任务
	for i := 0; i < 5; i++ {
		task := &MockTask{
			RetryCounter: new(int32),
			RunCounter:   &executed,
			RunDuration:  50 * time.Millisecond,
		}
		pool.Submit(task, 0, 0)
	}

	// 关闭线程池
	pool.Close()

	// 验证所有任务已完成
	if atomic.LoadInt32(&executed) != 5 {
		t.Errorf("expected 5 tasks to be executed, but got %d", executed)
	}
}

// TestMiddlewareIntegration 测试线程池与 Gin 中间件的集成
func TestMiddlewareIntegration(t *testing.T) {
	// 这里可以模拟 Gin 的请求，验证中间件与线程池的集成是否正常
	// 可使用 httptest 包生成请求，检查是否正确提交到线程池并完成任务
}
