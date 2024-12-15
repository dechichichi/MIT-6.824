package mr

import (
	"os"
	"strconv"
)

// 一个任务 应该包括：
// 任务类型 任务ID 使用Reduce数量  任务本体
type Task struct {
	TaskType   TaskType
	TaskID     int
	ReducerNum int
	fileslice  []string
}

type TaskType int

// 一个任务阶段包括
// 分配阶段 枚举阶段
type Phase int

type State int

// 枚举任务的类型
const (
	MapTask TaskType = iota //itoa=0
	ReduceTask
	WaitTask
	ExitTask
)

const ()

func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
