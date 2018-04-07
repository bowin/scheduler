package scheduler

import (
	"fmt"
)

type TestTask struct {
	ID        int
	index     int
	timestamp int64
	isCycle   bool
}

func (t TestTask) DueTimestamp() int64 {
	return t.timestamp
}
func (t TestTask) GetID() int {
	return t.ID
}
func (t TestTask) IsCycle() bool {
	return t.isCycle
}
func (t TestTask) Run() {
	fmt.Printf("run task %v \n", t)
}
func (t TestTask) SetIndex(i int) {
	t.index = i
}
