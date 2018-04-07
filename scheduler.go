package scheduler

import (
	"container/heap"
	"fmt"
	"time"
)

// TaskStore 用于批量加载tasks
type TaskStore interface {
	BuckLoad() Tasks
}

// Task 任务需要实现的接口
type Task interface {
	DueTimestamp() int64
	IsCycle() bool
	Run()
	GetID() int
	SetIndex(i int)
}

// Tasks 调度器所有的tasks
type Tasks []Task

func (pq Tasks) Len() int { return len(pq) }

func (pq Tasks) Less(i, j int) bool {
	return pq[i].DueTimestamp() < pq[j].DueTimestamp()
}

func (pq Tasks) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].SetIndex(i)
	pq[j].SetIndex(j)
}

func (pq *Tasks) Push(x interface{}) {
	n := len(*pq)
	t := x.(Task)
	t.SetIndex(n)
	*pq = append(*pq, t)
}

func (pq *Tasks) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.SetIndex(-1)
	*pq = old[0 : n-1]
	return &item
}

// Schedule 调度器接口
type Schedule interface {
	BulkImportTasks(s *TaskStore)
	AddTask(t Task)
	CancelTask()
	Schedule()
	RunTask()
}

// Scheduler 调度器
type Scheduler struct {
	taskChannel chan Task
	Tasks
}

// Init 调度器初始化
func (s *Scheduler) Init(store *TaskStore) {
	if store == nil {
		return
	}
	s.Tasks = (*store).BuckLoad()
	heap.Init(&s.Tasks)
}

func (s *Scheduler) getOneTask() *Task {
	t := heap.Pop(&s.Tasks)
	return t.(*Task)
}

// AddTask 加入任务
func (s *Scheduler) AddTask(t Task) {
	heap.Push(&s.Tasks, t)
}

// CancelTask 取消任务
func (s *Scheduler) CancelTask(id int) {
	if len(s.Tasks) == 0 {
		return
	}
	index := s.findTask(id, 0)
	if index > -1 {
		heap.Remove(&s.Tasks, index)
	}
}
func (s *Scheduler) findTask(id int, fromIndex int) int {
	if fromIndex >= len(s.Tasks) {
		return -1
	}
	if s.Tasks[fromIndex].GetID() == id {
		return fromIndex
	}
	left := s.findTask(id, 2*fromIndex+1)
	if left > -1 {
		return left
	}
	right := s.findTask(id, 2*fromIndex+2)
	if right > -1 {
		return right
	}
	return -1
}

// RunTask 运行任务
func (s *Scheduler) RunTask() {
	task := s.getOneTask()
	if task == nil {
		return
	}
	(*task).Run()
}

// Schedule 开始调度
func (s *Scheduler) Schedule() {
	ticker := time.Tick(time.Second)
	for {
		select {
		case t := <-s.taskChannel:
			s.AddTask(t)
		case <-ticker:
			fmt.Print("1s")
			s.RunTask()
		}
	}
}

// NewSchedule 新建调度器
func NewSchedule() Scheduler {
	return Scheduler{
		taskChannel: make(chan Task, 10),
		Tasks:       make(Tasks, 0),
	}
}
