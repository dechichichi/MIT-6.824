package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
)

//Master节点的RPC服务端，负责分配任务给worker节点，并监控worker节点的状态，当所有worker节点完成任务后，Master节点会汇总结果并返回给客户端。
//MapReduce的基本思路是启动一个coordinator分配多个worker做map任务

type Coordinator struct {
	// Your definitions here.

}

func (c *Coordinator) handler(files string, nReduce int) error {
	//这个函数用来启动worker节点的
	var KeyValueList []KeyValue
	// Your code here.
	// 1. 解析文件，得到KeyValueList
	parts := strings.Split(files, " ")
	n := string(nReduce)
	// 2. 启动worker节点，并将KeyValueList发送给worker节点
	// 3. 等待worker节点完成任务，并汇总结果
	// 4. 返回结果
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
	for i := 0; i < nReduce; i++ {
		//对于每个文件，启动一个协程来处理
		go c.handler(files[i])
		if files[i] == "" {
			break
		}
	}
	c.server()
	return &c
}
