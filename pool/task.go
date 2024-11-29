package pool

import (
	"errors"
	"log"
	"math/rand"
)

// ExampleTask 示例任务结构
type ExampleTask struct {
	Name string
}

// Run 执行任务逻辑
func (e *ExampleTask) Run() error {
	log.Printf("Executing task: %s", e.Name)
	// 模拟成功或失败
	if rand.Intn(2) == 0 {
		return errors.New("simulated task error")
	}
	return nil
}

// NewExampleTask 创建一个新的 ExampleTask
func NewExampleTask(name string) *ExampleTask {
	return &ExampleTask{Name: name}
}
