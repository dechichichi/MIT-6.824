6.5840 实验 2：键/值服务器
协作政策 // 提交实验室 // 设置 Go // 指导 // Piazza

介绍
在本实验中，您将为单台计算机构建一个键/值服务器 这可确保每个操作都只执行一次 网络故障和操作 是线性化的。后 Labs 将复制像这样的服务器来处理服务器崩溃。

客户端可以向键/值服务器发送三种不同的 RPC：Put（key， value）、Append（key， arg） 和 Get（key）。服务器维护 键/值对的内存中映射。键和值是字符串。Put（key， value） 安装或替换 中特定键的值 映射中，Append（key， arg） 将 arg 附加到 key 的值并返回旧值，而 Get（key） 获取当前值 对于密钥。不存在的键的 Get 应返回 空字符串。Append 到不存在的键应该起作用 就像现有值是一个长度为零的字符串一样。 每个客户端都通过 具有 Put/Append/Get 方法的 Clerk。文员管理 RPC 与服务器的交互。

您的服务器必须安排应用程序调用 到 Clerk Get/Put/Append 方法都是线性化的。如果 客户端请求不是并发的， 每个客户端 Get/Put/Append 调用都应遵守 前面的 调用。对于并发调用，返回值和最终状态必须为 这与在某些 次序。如果调用在时间上重叠，则调用是并发的：例如，如果 客户端 X 调用 Clerk.Put（），客户端 Y 调用 Clerk.Append（），然后客户端 X 的调用返回。一个 call 必须观察之前完成的所有调用 呼叫开始。

线性化对于应用程序来说很方便，因为它是 您从处理请求的单个服务器看到的行为 一个时间。例如，如果一个客户端从 server 的 NEW UPDATE 请求，随后启动了从其他 客户可以保证看到该更新的效果。提供 对于单个服务器来说，线性化相对容易。

开始
我们在 src/kvsrv 中为您提供框架代码和测试。您将 需要修改 kvsrv/client.go、kvsrv/server.go 和 kvsrv/common.go。

要启动并运行，请执行以下命令。 不要忘记 git pull 来获取最新的软件。

$ cd ~/6.5840
$ git pull
...
$ cd src/kvsrv
$ go test
...
$
没有网络故障的键值服务器（简单)
您的首要任务是实施一个在没有 drop 时有效的解决方案 消息。

您需要将 RPC 发送代码添加到 Clerk Put/Append/Get 方法，并在 server.go 中实现 Put、Append（） 和 Get（） RPC 处理程序。

您完成此任务时 通过 测试套件：“一个客户端”和“多个客户端”。

使用 go test -race 检查您的代码是否没有争用。
带有已删除消息的键/值服务器（简单)
现在，您应该修改您的解决方案，以便在遇到 dropped 时继续 消息（例如 RPC 请求和 RPC 回复）。 如果消息丢失，则客户端的 ck.server.Call（） 将返回 false（更准确地说，Call（） 等待 对于超时间隔的回复消息，并返回 false 如果在该时间内没有收到回复）。 您将面临的一个问题是，Clerk 可能必须多次发送 RPC，直到它 成功。但是，每次调用 Clerk.Put（） 或 Clerk.Append（） 都应该 导致一次执行，因此您必须确保 重新发送不会导致服务器执行 请求两次。
将代码添加到 Clerk 以在未收到回复时重试， 如果操作需要，则向 server.go 过滤重复项 它。这些说明包括指导 在重复检测时。

您需要唯一标识客户端操作，以确保 键/值服务器只执行每个 API 一次。
您必须仔细考虑服务器必须处于什么状态 maintain 用于处理重复的 Get（）、Put（）、 和 Append（） 请求（如果有）。
您的重复检测方案应快速释放服务器内存， 例如，通过让每个 RPC 暗示客户端已经看到了 回复。可以假设客户端将 一次只能打电话给 Clerk。
您的代码现在应该通过所有测试，如下所示：

$ go test
Test: one client ...
  ... Passed -- t  3.8 nrpc 31135 ops 31135
Test: many clients ...
  ... Passed -- t  4.7 nrpc 102853 ops 102853
Test: unreliable net, many clients ...
  ... Passed -- t  4.1 nrpc   580 ops  496
Test: concurrent append to same key, unreliable ...
  ... Passed -- t  0.6 nrpc    61 ops   52
Test: memory use get ...
  ... Passed -- t  0.4 nrpc     4 ops    0
Test: memory use put ...
  ... Passed -- t  0.2 nrpc     2 ops    0
Test: memory use append ...
  ... Passed -- t  0.4 nrpc     2 ops    0
Test: memory use many puts ...
  ... Passed -- t 11.5 nrpc 100000 ops    0
Test: memory use many gets ...
  ... Passed -- t 12.2 nrpc 100001 ops    0
PASS
ok      6.5840/kvsrv    39.000s
每个 Passed 之后的数字都是以秒为单位的实时数字， 发送的 RPC 数量（包括客户端 RPC），以及 执行的键/值操作数 （Clerk Get/Put/Append 调用）。