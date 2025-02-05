package kvsrv

import (
	"fmt"
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
	mu     sync.Mutex
	value  map[string]string
	seqMap map[int64]bool // 用于跟踪已处理的序列号
	// Your definitions here.
}

// Get（key） 获取当前值 对于密钥。不存在的键的 Get 应返回 空字符串。
func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if _, ok := kv.seqMap[args.Seq]; ok {
		log.Printf("重复请求序列号: %v", args.Seq)
		return
	}
	kv.seqMap[args.Seq] = true
	// Your code here.
	if kv.value[args.Key] != "" {
		reply.Value = kv.value[args.Key]
	} else {
		reply.Value = ""
	}
}

// Put（key， value） 安装或替换 中特定键的值 映射中
func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	if _, ok := kv.seqMap[args.Seq]; ok {
		// 如果已处理过，返回错误或不做任何操作
		log.Printf("重复请求序列号: %v", args.Seq)
		return
	}
	// 标记序列号为已处理
	kv.seqMap[args.Seq] = true
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if kv.value[args.Key] != "" {
		kv.value[args.Key] = args.Value
		log.Printf("替换 %s: %s", args.Key, args.Value)
	} else {
		kv.value[args.Key] = args.Value
		fmt.Println("安装 %s: %s", args.Key, args.Value)
	}
	// Your code here.
}

// Append（key， arg） 将 arg 附加到 key 的值并返回旧值
func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	if _, ok := kv.seqMap[args.Seq]; ok {
		// 如果已处理过，返回错误或不做任何操作
		println("重复请求")
		return
	}
	// 标记序列号为已处理
	kv.seqMap[args.Seq] = true
	kv.mu.Lock()
	defer kv.mu.Unlock()
	reply.Value = kv.value[args.Key]
	kv.value[args.Key] = args.Value
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.value = make(map[string]string) // 初始化 value map
	kv.seqMap = make(map[int64]bool)   // 初始化 seqMap
	// You may need initialization code here.

	return kv
}
