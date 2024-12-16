package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

//call 函数通过1234端口传入args和reply的内存地址，调用rpcname（Coordinator.函数名），
//通过反射机制"远程"调用Coordinator的该函数，Coordinator通过内存地址读取入参写出结果。
//worker.go里面的Worker方法调用CallExample，先运行Coordinator，再运行worker，
//看看worker端打印返回来的经过Coordinator加工过的数字

//Master节点的RPC服务端，负责分配任务给worker节点，并监控worker节点的状态，当所有worker节点完成任务后，Master节点会汇总结果并返回给客户端。
//MapReduce的基本思路是启动一个coordinator分配多个worker做map任务

type Coordinator struct {
	// Your definitions here.
	ReduceNum int
	Task      Task
	DistPhase Phase
	Mutex     sync.Mutex // 锁
}

func (c *Coordinator) Handler(files string, nReduce int) error {
	//任务分配
	c.Task.Filename = files
	c.Task.ReducerNum = nReduce
	Worker(mapf, reducef)
	var t time.Duration // 用于记录经过的时间
	// 状态监控
	for callDone() == false && t < 10*time.Second {
		time.Sleep(100 * time.Millisecond)
		t += 100 * time.Millisecond
	}
	//结果收集
	c.DistPhase = AllDone
	//异常处理
	return nil
}

func MakeCoordinator(files []string, nReduce int) *Coordinator {
	if nReduce <= 0 {
		panic(fmt.Sprintf("nReduce must be positive, not %d", nReduce))
	}
	c := Coordinator{}
	for i := 0; i < nReduce && i < len(files); i++ { // 确保不会超出files的长度
		// 对于每个文件，启动一个协程来处理
		go c.Handler(files[i], i)
	}
	go c.Server() // 启动 RPC 服务器
	return &c
}

func mapf(filename string, contents string) []KeyValue {
	// function to detect word separators.
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	// split contents into an array of words.
	words := strings.FieldsFunc(contents, ff)

	kva := []KeyValue{}
	for _, w := range words {
		kv := KeyValue{w, "1"}
		kva = append(kva, kv)
	}
	return kva
}

func reducef(key string, values []string) string {
	return strconv.Itoa(len(values))
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) Server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	sockname := coordinatorSock()
	os.Remove(sockname) // 删除旧的套接字文件
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// RPC 方法，用于获取任务
func (c *Coordinator) GetTask(args *TaskArgs, reply *Task) error {
	// 这里应该是分配任务的逻辑
	// 例如，初始化 reply 结构体并返回
	reply.TaskType = MapTask
	reply.TaskID = 1
	reply.ReducerNum = c.ReduceNum
	reply.Filename = "example.txt"
	return nil
}

// RPC 方法，用于标记任务完成
func (c *Coordinator) Done() bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if c.DistPhase == AllDone {
		fmt.Printf("All workers done\n")
		return true // 应该返回true，表示所有工作都已完成
	} else {
		return false
	}
}
