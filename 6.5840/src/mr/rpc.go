package mr

import (
	"os"
	"strconv"
)

//rpc:一种通过网络从远程计算机程序上请求服务，而不需要了解底层网络技术的协议

// 一个任务 应该包括：
// 任务类型 任务ID 使用Reduce数量  任务本体
type Task struct {
	ReduceNum int
	FileName  string
	TaskType  TaskType
	TaskID    int
	State     State
}

// 这个参数无实际意义，只是为了让RPC调用时能够传参
type TaskArgs struct {
	TaskID int
}

// 一个任务阶段包括
// 分配阶段 枚举阶段
type TaskType int
type State int

// 任务类型
const (
	MapTask TaskType = iota
	ReduceTask
	Done
)

// 所有任务状态类型
const (
	Working State = iota
	Waiting
	AllDone
)

func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
