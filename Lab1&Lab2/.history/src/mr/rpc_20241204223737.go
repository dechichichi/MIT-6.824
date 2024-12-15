package mr

import (
	"os"
	"strconv"
)

//一种通过网络从远程计算机程序上请求服务，而不需要了解底层网络技术的协议

// 一个任务 应该包括：
// 任务类型 任务ID 使用Reduce数量  任务本体
type Task struct {
	TaskType   TaskType
	TaskID     int
	ReducerNum int
	fileslice  []string
}

// 一个任务阶段包括
// 分配阶段 枚举阶段
type TaskType int

type Phase int

type State int

// 枚举任务的类型
const (
	MapTask TaskType = iota //itoa=0
	ReduceTask
	WaitTask
	ExitTask
)

// 枚举阶段类型
const (
	MapPhase Phase = iota
	ReducePhase
	AllDone
)

// 任务状态类型
const (
	Working State = iota
	Waiting
	Done
)

func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
