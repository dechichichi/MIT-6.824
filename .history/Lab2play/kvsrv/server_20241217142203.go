package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

//您的首要任务是实施一个在没有 drop 时有效的解决方案 消息。

//您需要将 RPC 发送代码添加到 Clerk Put/Append/Get 方法，并在 server.go 中实现
//Put、Append（） 和 Get（） RPC 处理程序。

// 客户端可以向键/值服务器发送三种不同的 RPC：Put（key， value）、Append（key， arg） 和 Get（key）。
// 服务器维护 键/值对的内存中映射。键和值是字符串。
// Put（key， value） 安装或替换 中特定键的值 映射中，
// Append（key， arg） 将 arg 附加到 key 的值并返回旧值，
// 而 Get（key） 获取当前值 对于密钥。不存在的键的 Get 应返回 空字符串。
// Append 到不存在的键应该起作用 就像现有值是一个长度为零的字符串一样。
// 每个客户端都通过 具有 Put/Append/Get 方法的 Clerk。文员管理 RPC 与服务器的交互。
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
	kv.mu.Lock()
	defer kv.mu.Unlock()
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
