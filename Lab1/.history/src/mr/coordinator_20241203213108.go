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
	// Your definitions here.

}

// Your code here -- RPC handlers for the worker to call.

func (c *Coordinator) handler(files []string) error {
	reply.Y = args.X + 1
	return nil
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	if nReduce <= 0 {
		panic(fmt.Sprintf("nReduce must be positive, not %d", nReduce))
	}
	c := Coordinator{}
	filesPerReduce := len(files) / nReduce
	remainingFiles := len(files) % nReduce
	// 使用WaitGroup来等待所有协程完成
	var wg sync.WaitGroup

	// 分配文件到每个reduce任务
	for i := 0; i < nReduce; i++ {
		// 计算当前reduce任务的文件范围
		start := i * filesPerReduce
		end := (i + 1) * filesPerReduce
		if i < remainingFiles { // 如果有剩余文件，分配给前几个reduce任务
			end++
		}

		// 为每个reduce任务启动一个协程
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			c.handler(files[start:end])
		}(i)
	}

	c.server()
	return &c
}
