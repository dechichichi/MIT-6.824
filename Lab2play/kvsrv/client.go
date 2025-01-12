package kvsrv

import (
	"crypto/rand"
	"math/big"
	"time"

	"kv/labrpc"
)

type Clerk struct {
	server *labrpc.ClientEnd
	seq    int64
	// You will have to modify this struct.
}

func (ck *Clerk) nextSeq() int64 {
	ck.seq++
	return ck.seq
}

//现在，您应该修改您的解决方案，以便在遇到 dropped 时继续 消息（例如 RPC 请求和 RPC 回复）。
//如果消息丢失，则客户端的 ck.server.Call（） 将返回 false（更准确地说，Call（） 等待
//对于超时间隔的回复消息，并返回 false 如果在该时间内没有收到回复）。
//您将面临的一个问题是，Clerk 可能必须多次发送 RPC，直到它 成功。
//但是，每次调用 Clerk.Put（） 或 Clerk.Append（） 都应该 导致一次执行，
//因此您必须确保 重新发送不会导致服务器执行 请求两次。
//将代码添加到 Clerk 以在未收到回复时重试， 如果操作需要，则向 server.go 过滤重复项 它。这些说明包括指导 在重复检测时。

//您需要唯一标识客户端操作，以确保 键/值服务器只执行每个 API 一次。
//您必须仔细考虑服务器必须处于什么状态 maintain 用于处理重复的 Get（）、Put（）、 和 Append（） 请求（如果有）。
//您的重复检测方案应快速释放服务器内存， 例如，通过让每个 RPC 暗示客户端已经看到了 回复。可以假设客户端将 一次只能打电话给 Clerk。

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(server *labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.server = server
	// You'll have to add code here.
	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
// 获取一个键的当前值。
// 如果键不存在，则返回""。
// 面对所有其他错误时，会不断尝试。
//
// 您可以使用如下代码发送RPC调用：
// ok := ck.server.Call("KVServer.Get", &args, &reply)
//
// args和reply的类型（包括它们是否是指针）必须与RPC处理函数的
// 参数的声明类型匹配。并且reply必须作为指针传递。
func (ck *Clerk) Get(key string) string {
	// 创建一个channel来接收结果
	replyChan := make(chan string)
	args := GetArgs{
		Key: key,
		Seq: ck.nextSeq(), // 生成序列号
	}
	var reply GetReply // 必须提供一个有效的回复结构体

	// 启动goroutine来执行Call方法
	go func() {
		// 调用Call方法并获取结果
		ok := ck.server.Call("KVServer.Get", &args, &reply)

		if ok {
			// 如果调用成功，发送回复的内容到channel
			replyChan <- reply.Value
		}
	}()

	// 设置超时
	timeout := time.After(5 * time.Second)

	// 等待结果或超时
	select {
	case reply := <-replyChan:
		// 如果从channel接收到结果，返回结果
		return reply
	case <-timeout:
		// 如果超时，返回空字符串
		return ""
	}
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer."+op, &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) string {
	args := PutAppendArgs{
		Key:   key,
		Value: value,
		Seq:   ck.nextSeq(), // 生成序列号
	}
	replyChan := make(chan string)
	var reply GetReply // 必须提供一个有效的回复结构体
	go func() {

		// 发送RPC调用
		ok := ck.server.Call("KVServer."+op, &args, &reply)

		if ok {
			replyChan <- "true"
		}
	}()
	// 设置超时
	timeout := time.After(5 * time.Second)

	// 等待结果或超时
	select {
	case reply := <-replyChan:
		// 如果从channel接收到结果，返回结果
		return reply
	case <-timeout:
		// 如果超时，返回空字符串
		return ""
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

// Append value to key's value and return that value
func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}
