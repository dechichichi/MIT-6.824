package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Coordinator struct {
	Tasks      []Task
	TaskStatus map[int]State // 任务状态
	Mutex      sync.Mutex    // 锁
}

func MakeCoordinator(files []string, nReduce int) *Coordinator {
	if nReduce <= 0 {
		panic(fmt.Sprintf("nReduce must be positive, not %d", nReduce))
	}
	c := Coordinator{
		Tasks:      make([]Task, 0),
		TaskStatus: make(map[int]State),
	}
	for i, file := range files {
		task := Task{TaskType: MapTask, TaskID: i, ReduceNum: nReduce, FileName: file}
		c.Tasks = append(c.Tasks, task)
		c.TaskStatus[task.TaskID] = Waiting
	}
	task := Task{TaskType: ReduceTask, TaskID: len(files) + 1, ReduceNum: nReduce}
	c.Tasks = append(c.Tasks, task)
	c.TaskStatus[task.TaskID] = Waiting
	c.server()
	return &c
}

func (c *Coordinator) Done() bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for _, status := range c.TaskStatus {
		if status != AllDone {
			return false
		}
	}
	fmt.Printf("All workers done\n")
	return true
}

func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (c *Coordinator) GetTasks(args *TaskArgs, reply *[]Task) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// 初始化回复的任务切片
	*reply = []Task{}

	// 遍历所有任务，查找处于 Waiting 状态的任务
	for _, task := range c.Tasks {
		if c.TaskStatus[task.TaskID] == Waiting {
			// 将任务添加到回复的任务切片中
			*reply = append(*reply, task)
			// 更新任务状态为 Working
			c.TaskStatus[task.TaskID] = Working
		}
	}

	// 如果没有找到任何任务，返回 nil
	if len(*reply) == 0 {
		return nil
	}

	return nil
}

func (c *Coordinator) DoneTask(args *TaskArgs, reply *Task) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	taskID := args.TaskID
	c.TaskStatus[taskID] = AllDone

	return nil
}
