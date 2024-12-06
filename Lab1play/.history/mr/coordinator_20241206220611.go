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
	"unicode"
)

//Master节点的RPC服务端，负责分配任务给worker节点，并监控worker节点的状态，当所有worker节点完成任务后，Master节点会汇总结果并返回给客户端。
//MapReduce的基本思路是启动一个coordinator分配多个worker做map任务

type Coordinator struct {
	// Your definitions here.

}

func (c *Coordinator) handler(files string, nReduce string) error {
	Worker(mapf, reducef)
	return nil
}
func mapf(files string, n string) []KeyValue {
	// function to detect word separators.
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	// split contents into an array of words.
	words := strings.FieldsFunc(contents, ff)

	kva := []mr.KeyValue{}
	for _, w := range words {
		kv := mr.KeyValue{w, "1"}
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

func (c *Coordinator) Done() bool {
	ret := false
	//
	return ret
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
