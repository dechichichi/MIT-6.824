package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

//您的首要任务是实施一个在没有 drop 时有效的解决方案 消息。

//您需要将 RPC 发送代码添加到 Clerk Put/Append/Get 方法，并在 server.go 中实现 Put、Append（） 和 Get（） RPC 处理程序。

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.Mutex

	// Your definitions here.
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	// You may need initialization code here.

	return kv
}
