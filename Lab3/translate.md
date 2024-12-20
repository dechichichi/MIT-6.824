6.5840 - 2024 年春季
6.5840 实验 3：木筏
协作政策 // 提交实验室 // 设置 Go // 指导 // Piazza

介绍
这是一系列实验中的第一个实验，您将在其中构建一个 容错键/值存储系统。在这个 lab 中，您将实现 Raft，这是一种复制的状态机协议。 在下一个实验中，您将在 筏。然后，您将 “分片” 您的 Service 多个复制的状态机以获得更高的性能。

复制的服务出错 tolerance 通过存储其状态 （即 data） 的完整副本来实现 在多个副本服务器上。 复制允许 即使某些 其服务器遇到故障（崩溃、损坏或片状 网络）。挑战在于，故障可能会导致 replicas 来保存数据的不同副本。

Raft 将客户端请求组织成一个序列，称为 日志，并确保所有副本服务器都看到相同的日志。 每个副本执行客户端请求 按日志顺序，将它们应用于服务状态的本地副本。 由于所有活动副本 看到相同的日志内容，它们都执行相同的请求 以相同的顺序，因此继续具有相同的服务 州。如果服务器出现故障但后来恢复了，Raft 会负责 使其日志保持最新状态。Raft 将继续以 只要至少大多数服务器都处于活动状态并且可以 互相交谈。如果没有这样的多数，Raft 将 没有进展，但会尽快从中断的地方继续 大多数人可以再次通信。

在本实验中，您将把 Raft 实现为 Go 对象类型 替换为关联的方法，旨在用作 更大的服务。一组 Raft 实例通过 RPC 来维护复制的日志。你的 Raft 接口将 还支持无限序列的编号命令 称为日志条目。条目使用索引进行编号 数字。具有给定索引的日志条目最终将 被承诺。此时，你的 Raft 应该发送日志 条目以执行它。

您应该遵循扩展 Raft 论文中的设计， 特别注意图 2. 您将实现论文中的大部分内容，包括保存 persistent 状态并在节点发生故障后读取它，以及 然后重新启动。您将不会实现 cluster 成员变更（第 6 节）。

本实验分为四个部分。您必须在 相应的到期日。

开始
如果已完成实验 1，则已拥有实验的副本 源代码。 如果没有， 你可以找到通过 Git 获取源码的指导 在实验室 1 说明中。

我们为你提供骨架代码 src/raft/raft.go。我们还 提供一组测试，您应该使用这些测试来驱动 实施工作，我们将使用这些信息对您提交的 实验室。测试位于 src/raft/test_test.go 中。

当我们对您的提交进行评分时，我们将在没有 -race 标志的情况下运行测试。 但是，您应该检查您的代码没有 races， 通过运行 在开发解决方案时使用 -race 标志进行测试。

要启动并运行，请执行以下命令。 不要忘记 git pull 来获取最新的软件。

$ cd ~/6.5840
$ git pull
...
$ cd src/raft
$ go test
Test (3A): initial election ...
--- FAIL: TestInitialElection3A (5.04s)
        config.go:326: expected one leader, got none
Test (3A): election after network failure ...
--- FAIL: TestReElection3A (5.03s)
        config.go:326: expected one leader, got none
...
$
代码
通过在 raft/raft.go 中添加代码来实现 Raft。在该文件中，您将找到 框架代码，以及如何发送和接收的示例 RPC 的。
您的实现必须支持以下接口，该接口 测试人员和（最终）你的 Key/Value 服务器将使用。 你可以在 raft.go 的评论中找到更多详细信息。

// create a new Raft server instance:
rf := Make(peers, me, persister, applyCh)

// start agreement on a new log entry:
rf.Start(command interface{}) (index, term, isleader)

// ask a Raft for its current term, and whether it thinks it is leader
rf.GetState() (term, isLeader)

// each time a new entry is committed to the log, each Raft peer
// should send an ApplyMsg to the service (or tester).
type ApplyMsg
服务调用 Make（peers，me,...） 来创建一个 Raft 对等节点。peers 参数是网络标识符数组 的 Raft 对等节点（包括这个节点），用于 RPC。me 参数是此 Peer 节点在 Peer 节点中的索引 数组。start（command） 要求 Raft 开始处理 将命令附加到复制的日志中。Start（） 应立即返回，而无需等待日志附加 以完成。该服务希望您的实现将每个新提交的日志条目的 ApplyMsg 发送到 applyCh 通道参数到 Make（）。

raft.go 包含发送 RPC 的示例代码 （sendRequestVote（）） 处理传入的 RPC （RequestVote（）） 的 Vote 请求。 你的 Raft 对等节点应该使用 labrpc Go 交换 RPC package （src/labrpc 中的源代码）。 测试器可以告诉 labrpc 延迟 RPC， 对它们重新排序，并丢弃它们以模拟各种网络故障。 虽然您可以临时修改 labrpc， 确保你的 Raft 与原始 labrpc 一起工作， 因为这是我们用来测试和评分您的实验室的方法。 您的 Raft 实例只能与 RPC 交互;例如 不允许他们使用共享的 Go 变量进行通信 或文件。

后续实验基于此实验构建，因此请务必提供 您自己有足够的时间编写可靠的代码。

第 3A 部分：领导人选举（温和派)
实现 Raft leader 选举和心跳（AppendEntries RPC 中没有 log entries） 的 S第 3A 部分的目标是 选举单一领导人，以便领导人继续担任领导人 如果没有失败，并且如果 旧主服务器发生故障，或者进出旧主服务器 （Old Leader） 的数据包丢失。 运行 go test -run 3A 来测试你的 3A 代码。

你不能轻松地直接运行 Raft 实现;相反，您应该 通过测试器运行它，即 go test -run 3A 。
按照论文的图 2 进行操作。此时您关心发送 以及接收 RequestVote RPC，则与 选举，以及与领导人选举相关的国家，
添加 图 2 状态以进行领导者选举 添加到 raft.go 中的 Raft 结构体中。 您还需要定义一个 struct 来保存有关每个日志条目的信息。
填写 RequestVoteArgs 和 RequestVoteReply 结构。修改 Make（） 以创建一个后台协程，该协程将启动 leader 通过发送 RequestVote RPC 来定期进行选举 从另一个同行那里听到了一段时间。 实现 RequestVote（） RPC 处理程序，以便服务器投票支持一个 另一个。
要实现检测信号，请定义一个 AppendEntries RPC 结构（尽管你不能 还需要所有参数），并让 leader 发送 他们定期出局。编写 AppendEntries RPC 处理程序方法。
测试器要求 leader 发送的检测信号 RPC 不超过 每秒 10 次。
测试者要求你的 Raft 在 5 分钟内选出一个新的 leader 秒数（如果大多数 Peer 节点可以 仍然通信）。
该论文的第 5.2 节提到了 150 次范围内的选举超时 设置为 300 毫秒。这样的范围只有在领导者 发送检测信号的频率远高于每 150 次发送一次 毫秒（例如，每 10 毫秒一次）。因为测试器限制您每 其次，您必须使用更大的 elect timeout 比论文的 150 到 300 毫秒，但不要太大，因为那样你 可能无法在 5 秒内选出领导者。
您可能会发现 Go 的 rand 很有用。
您需要编写定期执行操作的代码，或者 时间延误后。最简单的方法是创建 一个带有调用 time 的循环的 goroutine。sleep（）; 请参阅 Make（） 为此目的创建的 ticker（） 协程。 不要占用 Go 的时间。计时器或时间。Ticker，这 很难正确使用。
如果您的代码无法通过测试， 再次阅读论文的图 2;领导者的完整逻辑 选举分布在图的多个部分。
不要忘记实现 GetState（）。
测试者调用 Raft 的 rf。Kill（） 时 永久关闭实例。您可以使用 rf.killed（） 检查 Kill（） 是否已被调用。 您可能希望在所有循环中执行此操作，以避免出现 dead Raft 实例会打印令人困惑的消息。
Go RPC 只发送名称以大写字母开头的结构体字段。 子结构还必须具有大写的字段名称（例如，日志记录的字段 在数组中）。labgob 包会警告你这一点; 不要忽略警告。
本实验中最具挑战性的部分可能是调试。花一些 时间使您的实现易于调试。指 Guidance （指南） 页面，了解调试提示。
在提交第 3A 部分之前，请确保通过 3A 测试，以便 您会看到如下内容：

$ go test -run 3A
Test (3A): initial election ...
  ... Passed --   3.5  3   58   16840    0
Test (3A): election after network failure ...
  ... Passed --   5.4  3  118   25269    0
Test (3A): multiple elections ...
  ... Passed --   7.3  7  624  138014    0
PASS
ok  	6.5840/raft	16.265s
$
每个 “Passed” 行包含五个数字;这些时间是 test 以秒为单位，Raft peer 的数量、 测试期间发送的 RPC 数量，则 RPC 消息和日志条目数 Raft 报告已提交。您的数字将与那些不同 显示在这里。如果您愿意，可以忽略这些数字，但它们可能会有所帮助 您对 implementation 发送的 RPC 数量进行健全性检查。 对于所有实验 3、4 和 5，评分脚本都将使您的 如果所有测试花费的时间超过 600 秒，则解决方案 （go test） ），或者如果任何单个测试花费的时间超过 120 秒。

当我们对您的提交进行评分时，我们将在没有 -race 标志的情况下运行测试。但是，您应该确保您的代码 始终通过 -race 标志的测试。

第 3B 部分：对数（硬)
实现 leader 和 follower 代码以添加新的日志条目， ，以便 go test -run 3B 测试通过。

运行 git pull 以获取最新的实验室软件。
您的第一个目标应该是通过 TestBasicAgree3B（）。 首先实现 Start（），然后编写代码 通过 AppendEntries RPC 发送和接收新的日志条目， 如图 2 所示。发送每个新提交的条目 在每个对等体的 applyCh 上。
您需要实施选举 限制（论文中的第 5.4.1 节）。
您的代码可能具有重复检查某些事件的循环。 不要有这些循环 持续执行而不暂停，因为 会减慢您的实现速度，使其无法通过测试。 使用 Go 的条件变量 或插入时间。睡眠（10 * 时间。毫秒）。
为未来的实验室做个好事，编写（或重写）代码 这很干净。如需想法，请重新访问我们的 指南页面，其中包含有关如何 开发和调试您的代码。
如果你没有通过测试，请查看 test_test.go 和 config.go 来 了解正在测试的内容。config.go 也 说明了测试人员如何使用 Raft API。
如果代码运行速度太慢，则即将到来的实验的测试可能会失败。 您可以检查您的解决方案使用多少实时时间和 CPU 时间 time 命令。下面是典型的输出：

$ time go test -run 3B
Test (3B): basic agreement ...
  ... Passed --   0.9  3   16    4572    3
Test (3B): RPC byte count ...
  ... Passed --   1.7  3   48  114536   11
Test (3B): agreement after follower reconnects ...
  ... Passed --   3.6  3   78   22131    7
Test (3B): no agreement if too many followers disconnect ...
  ... Passed --   3.8  5  172   40935    3
Test (3B): concurrent Start()s ...
  ... Passed --   1.1  3   24    7379    6
Test (3B): rejoin of partitioned leader ...
  ... Passed --   5.1  3  152   37021    4
Test (3B): leader backs up quickly over incorrect follower logs ...
  ... Passed --  17.2  5 2080 1587388  102
Test (3B): RPC counts aren't too high ...
  ... Passed --   2.2  3   60   20119   12
PASS
ok  	6.5840/raft	35.557s

real	0m35.899s
user	0m2.556s
sys	0m1.458s
$
“ok 6.5840/raft 35.557s”意味着 Go 测量了 3B 所花费的时间 测试为 35.557 秒的实际 （挂钟） 时间。“用户 0m2.556s“表示代码消耗了 2.556 秒的 CPU 时间，或者 实际执行指令所花费的时间（而不是等待或 sleeping）。如果您的解决方案使用的实时时间远超过一分钟 对于 3B 测试，或者 CPU 时间远超过 5 秒，您可以运行 后来陷入困境。查找休眠或等待 RPC 所花费的时间 timeouts、在不休眠或等待条件的情况下运行的循环或 频道消息或发送的大量 RPC。
第 3C 部分：持久性（硬)
如果基于 Raft 的服务器重启，它应该会恢复服务 它停止的地方。这需要 Raft 保持持久状态，在重启后仍然存在。这 论文的图 2 提到了哪种状态应该是持久的。

一个真正的实现会写 Raft 的持久化状态保存到 disk，并且会读取 state from disk 的 intent 值。您的实现不会使用 磁盘;相反，它将保存和恢复持久状态 从 Persister 对象（参见 persister.go）。 调用 Raft.Make（） 的人会提供一个 Persister，它最初保存 Raft 最近的持久化状态（如果 any） 的 S SRaft 应该从该 Persister 初始化其 state，并使用它来保存其持久化 state 每次状态更改时。使用 Persister 的 ReadRaftState（） 和 Save（） 方法。

通过添加代码完成 raft.go 中的 persist（） 和 readPersist（） 函数，以保存和恢复持久化状态。您需要对 （或“序列化”）状态作为字节数组，以便将其传递给 Persister 的。使用 labgob 编码器; 请参阅 persist（） 和 readPersist（） 中的注释。labgob 类似于 Go 的 gob 编码器，但 如果出现以下情况，则打印错误消息 您尝试使用小写字段名称对结构进行编码。 现在，将 nil 作为第二个参数传递给 persister。Save（） 的 在 您的实现会更改 persistent state。 完成此操作后， 如果您的 implementation 的其余部分是正确的， 您应该通过所有 3C 测试。

You will probably need the optimization that backs up nextIndex by more than one entry at a time. Look at the extended Raft paper starting at the bottom of page 7 and top of page 8 (marked by a gray line). The paper is vague about the details; you will need to fill in the gaps. One possibility is to have a rejection message include:

    XTerm:  term in the conflicting entry (if any)
    XIndex: index of first entry with that term (if any)
    XLen:   log length
Then the leader's logic can be something like: A few other hints:
  Case 1: leader doesn't have XTerm:
    nextIndex = XIndex
  Case 2: leader has XTerm:
    nextIndex = leader's last entry for XTerm
  Case 3: follower's log is too short:
    nextIndex = XLen
运行 git pull 以获取最新的实验室软件。
3C 测试比 3A 或 3B 测试要求更高，并且失败 可能是由 3A 或 3B 代码中的问题引起的。
您的代码应通过所有 3C 测试（如下所示），以及 3A 和 3B 测试。

$ go test -run 3C
Test (3C): basic persistence ...
  ... Passed --   5.0  3   86   22849    6
Test (3C): more persistence ...
  ... Passed --  17.6  5  952  218854   16
Test (3C): partitioned leader and one follower crash, leader restarts ...
  ... Passed --   2.0  3   34    8937    4
Test (3C): Figure 8 ...
  ... Passed --  31.2  5  580  130675   32
Test (3C): unreliable agreement ...
  ... Passed --   1.7  5 1044  366392  246
Test (3C): Figure 8 (unreliable) ...
  ... Passed --  33.6  5 10700 33695245  308
Test (3C): churn ...
  ... Passed --  16.1  5 8864 44771259 1544
Test (3C): unreliable churn ...
  ... Passed --  16.5  5 4220 6414632  906
PASS
ok  	6.5840/raft	123.564s
$
最好在之前多次运行测试 submitting 并检查每次运行是否打印 PASS。

$ for i in {0..10}; do go test; done
第 3D 部分：原木压实（硬)
就目前的情况而言，重新启动的服务器会重放 complete Raft log 以恢复其状态。然而，事实并非如此 对于长时间运行的服务来说，记住完整的 Raft 日志是可行的 永远。相反，您将修改 Raft 以配合 持久存储其状态的 “快照” 时，位于 此时，Raft 会丢弃快照之前的日志条目。这 结果是持久性数据量较小，重启速度更快。 但是，现在追随者可能会远远落后于此 领导者已经丢弃了它需要追赶的日志条目;这 然后，leader 必须发送快照和日志，从 快照。扩展的 Raft 论文的第 7 节概述了该方案;您将不得不设计细节。

您的 Raft 必须提供以下功能，该服务 可以使用其状态的序列化快照进行调用：

快照 （索引 int， 快照 []字节）

在 Lab 3D 中，测试人员会定期调用 Snapshot（）。在实验 4 中，您将 编写一个调用 Snapshot（） 的 key/value 服务器;快照 将包含键/值对的完整表。 服务层在每个对等体（不是 就在领导者上）。

index 参数指示 反映在快照中。Raft 应该在 那个点。你需要修改你的 Raft 代码才能在 仅存储日志的尾部。

您需要实现 允许 Raft 领导者告诉滞后的 Raft peer 的 paper 将其 state 替换为 snapshot。你可能需要考虑一下 通过 InstallSnapshot 应如何与状态和规则交互 在图 2 中。

当 follower 的 Raft 代码收到 InstallSnapshot RPC 时，它可以 使用 applyCh 将快照发送到 一个 ApplyMsg。ApplyMsg 结构体定义已经 包含您需要的字段（以及测试人员期望的字段）。拿 请注意这些快照只会推进服务的状态，而不会 使其向后移动。

如果服务器崩溃，则必须从持久化数据重新启动。您的木筏 应该同时保留 Raft state 和相应的快照。 使用第二个参数来持久化。Save（） 保存快照。 如果没有快照，则传递 nil 作为第二个 论点。

当服务器重新启动时，应用程序层会读取持久化的 snapshot 并恢复其保存状态。

实现 Snapshot（） 和 InstallSnapshot RPC，以及 对 Raft 的更改以支持这些（例如，使用 修剪的原木）。当您的解决方案通过 3D 测试时，它就是完整的 （以及之前的所有 Lab 3 测试）。

git pull 以确保您拥有最新的软件。
一个好的起点是修改您的代码，使其 能够仅存储日志的一部分 从某个索引 X 开始。最初，您可以将 X 设置为零，并将 运行 3B/3C 测试。 然后让 Snapshot（index） 在索引之前丢弃日志， 并将 X 设置为等于 index。如果一切顺利，您应该 现在通过第一次 3D 测试。
下一步：如果 leader 没有，则让 leader 发送 InstallSnapshot RPC 具有使 follower 保持最新状态所需的日志条目。
在单个 InstallSnapshot RPC 中发送整个快照。 不要为 拆分快照。
Raft 必须以允许 Go 垃圾回收器释放和重用 记忆;这要求没有可访问的引用 （指针） 添加到丢弃的日志条目中。
使用全套 没有 -race 的实验室 3 测试 （3A+3B+3C+3D） 是 6 分钟的实时时间和 1 分钟 CPU 时间分钟。使用 -race 运行时，实际约为 10 分钟 时间和 2 分钟的 CPU 时间。
您的代码应该通过所有 3D 测试（如下所示），以及 3A、3B 和 3C 测试。

$ go test -run 3D
Test (3D): snapshots basic ...
  ... Passed --  11.6  3  176   61716  192
Test (3D): install snapshots (disconnect) ...
  ... Passed --  64.2  3  878  320610  336
Test (3D): install snapshots (disconnect+unreliable) ...
  ... Passed --  81.1  3 1059  375850  341
Test (3D): install snapshots (crash) ...
  ... Passed --  53.5  3  601  256638  339
Test (3D): install snapshots (unreliable+crash) ...
  ... Passed --  63.5  3  687  288294  336
Test (3D): crash and restart all servers ...
  ... Passed --  19.5  3  268   81352   58
PASS
ok      6.5840/raft      293.456s