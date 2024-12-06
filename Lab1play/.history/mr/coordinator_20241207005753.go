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
	TaskType  TaskType
	DistPhase Phase
	Mutex     sync.Mutex // 锁
}

func (c *Coordinator) handler(files string, nReduce string) error {
	//handle函数应该做什么：
	//任务分配：handle函数需要接收来自worker节点的请求，并根据当前的任务分配情况，给worker节点分配一个或多个任务（map或reduce任务）。

	//状态监控：handle函数需要监控worker节点的状态，确保它们正在执行分配给它们的任务。

	//结果收集：当worker节点完成任务后，handle函数需要收集这些节点的结果，并进行适当的处理（例如，对于map任务的输出，可能需要将结果传递给reduce任务）。

	//异常处理：如果worker节点失败或无法完成任务，handle函数需要能够处理这些异常情况，可能包括重新分配任务给其他worker节点。
	Worker(mapf, reducef)
	return nil
}

func MakeCoordinator(files []string, nReduce int) *Coordinator {
	if nReduce <= 0 {
		panic(fmt.Sprintf("nReduce must be positive, not %d", nReduce))
	}
	c := Coordinator{}
	for i := 0; i < nReduce; i++ {
		//对于每个文件，启动一个协程来处理
		go c.handler(files[i], string(i))
		if files[i] == "" {
			break
		}
	}
	c.server()
	return &c
}

// 已完工
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
