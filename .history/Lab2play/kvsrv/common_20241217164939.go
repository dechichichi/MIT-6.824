package kvsrv

// Put or Append
type PutAppendArgs struct {
	Key   string
	Value string
	//拒绝重复请求
	Flag bool
	// You'll have to add definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
}

//带有已删除消息的键/值服务器（简单)
//现在，您应该修改您的解决方案，以便在遇到 dropped 时继续 消息（例如 RPC 请求和 RPC 回复）。
//如果消息丢失，则客户端的 ck.server.Call（） 将返回 false（更准确地说，Call（） 等待 对于超时间隔的回复消息，并返回 false
//如果在该时间内没有收到回复）。 您将面临的一个问题是，Clerk 可能必须多次发送 RPC，直到它 成功。
//但是，每次调用 Clerk.Put（） 或 Clerk.Append（） 都应该 导致一次执行，因此您必须确保 重新发送不会导致服务器执行 请求两次。
//将代码添加到 Clerk 以在未收到回复时重试， 如果操作需要，则向 server.go 过滤重复项 它。这些说明包括指导 在重复检测时。

//您需要唯一标识客户端操作，以确保 键/值服务器只执行每个 API 一次。
//您必须仔细考虑服务器必须处于什么状态 maintain 用于处理重复的 Get（）、Put（）、 和 Append（） 请求（如果有）。
//您的重复检测方案应快速释放服务器内存， 例如，通过让每个 RPC 暗示客户端已经看到了 回复。可以假设客户端将 一次只能打电话给 Clerk。

type PutAppendReply struct {
	Value string
}

type GetArgs struct {
	Key string
	Seq int64
	// You'll have to add definitions here.
}

type GetReply struct {
	Value string
	Seq   int64
}
