6.5840 实验 1：MapReduce
协作政策 // 提交实验室 // 设置 Go // 指导 // Piazza

介绍
在本实验中，您将构建一个 MapReduce 系统。 您将实现一个调用应用程序 Map 和 Reduce 的工作进程 函数并处理读取和写入文件， 以及将任务分配给 工人并应对失败的工人。 您将构建类似于 MapReduce 论文的内容。 （注意：本实验使用“协调器”而不是论文的“master”。

开始
您需要设置 Go 来执行实验。

使用 git（一个版本控制系统）获取初始实验室软件。 要了解有关 git 的更多信息，请查看 Pro Git 书籍或 git 用户手册。

$ git clone git://g.csail.mit.edu/6.5840-golabs-2024 6.5840
$ cd 6.5840
$ ls
Makefile src
$
我们在 src/main/mrsequential.go 中为您提供了一个简单的顺序 mapreduce 实现。它运行映射并在 时间。我们还 为您提供几个 MapReduce 应用程序：word-count 在 mrapps/wc.go 中，以及文本索引器 在 mrapps/indexer.go 中。您可以运行 字数顺序如下：

$ cd ~/6.5840
$ cd src/main
$ go build -buildmode=plugin ../mrapps/wc.go
$ rm mr-out*
$ go run mrsequential.go wc.so pg*.txt
$ more mr-out-0
A 509
ABOUT 2
ACT 8
...
mrsequential.go 将其输出保留在文件 mr-out-0 中。 输入来自名为 pg-xxx.txt 的文本文件。

随意从 mrsequential.go 借用代码。 您还应该查看 mrapps/wc.go 以了解 MapReduce 应用程序代码如下所示。

对于此实验室和所有其他实验室，我们可能会对以下代码进行更新： 为您提供。为了确保您可以轻松获取这些更新 使用 git pull 合并它们，最好将代码留在 在原始文件中提供。您可以将 在实验室文章中指导;只是不要移动它。你可以把你的 在新文件中拥有新功能。

您的工作（中等/困难)
您的工作是实现一个分布式 MapReduce，它由 两个程序，协调器和工人。将会有 只有一个协调进程，以及一个或多个在 平行。在实际系统中，worker 将在一堆 不同的计算机，但对于本练习，您将在一台计算机上运行它们。 工作人员将通过 RPC 与协调器对话。每个 worker 进程将 在 Loop 中，询问 任务的协调器从一个或多个文件中读取任务的输入， 执行任务，将任务的输出写入一个 或更多文件，然后再次向协调器请求 new 任务。协调器应注意 worker 是否尚未完成 其任务在合理的时间内完成（对于本实验，请使用 10 个 秒），并将相同的任务分配给不同的 worker。
我们为您提供了一些代码，以便您开始。的 “main” 例程 协调器和 worker 位于 main/mrcoordinator.go 和 main/mrworker.go 中; 请勿更改这些文件。您应该将实现放在 mr/coordinator.go、mr/worker.go 和 mr/rpc.go 中。

下面介绍如何在字数统计 MapReduce 上运行代码 应用。首先，确保 word-count 插件是 全新构建：

$ go build -buildmode=plugin ../mrapps/wc.go
在主目录中，运行 coordinator。 mrcoordinator.go 的 pg-*.txt 参数是 输入文件;每个文件对应一个 “split”，并且是 input 添加到一个 Map 任务中。
$ rm mr-out*
$ go run mrcoordinator.go pg-*.txt
在一个或多个其他窗口中，运行一些 worker：

$ go run mrworker.go wc.so
当 worker 和 coordinator 完成后，查看输出 在 mr-out-* 中。完成实验后， 输出文件的排序并集应与 Sequential 匹配 output 中，如下所示：
$ cat mr-out-* | sort | more
A 509
ABOUT 2
ACT 8
...
我们在 main/test-mr.sh 中为您提供了一个测试脚本。测试 检查 wc 和索引器 MapReduce 应用程序 当将 pg-xxx.txt 文件指定为 输入。这些测试还会检查您的实现是否运行 Map 和 并行减少任务，并且您的实施可以从中恢复 在运行任务时崩溃的工作程序。

如果您现在运行测试脚本，它将挂起，因为协调器从未完成：

$ cd ~/6.5840/src/main
$ bash test-mr.sh
*** Starting wc test.
您可以在 mr/coordinator.go 的 Done 函数中将 ret ：= false 更改为 true，以便协调器立即退出。然后：

$ bash test-mr.sh
*** Starting wc test.
sort: No such file or directory
cmp: EOF on mr-wc-all
--- wc output is not the same as mr-correct-wc.txt
--- wc test: FAIL
$
测试脚本希望看到名为 mr-out-X 的文件（一个 对于每个 reduce 任务。mr/coordinator.go 和 mr/worker.go 的空实现不会生成这些文件（也不会执行大部分 任何其他），因此测试失败。

完成后，测试脚本输出应如下所示：

$ bash test-mr.sh
*** Starting wc test.
--- wc test: PASS
*** Starting indexer test.
--- indexer test: PASS
*** Starting map parallelism test.
--- map parallelism test: PASS
*** Starting reduce parallelism test.
--- reduce parallelism test: PASS
*** Starting job count test.
--- job count test: PASS
*** Starting early exit test.
--- early exit test: PASS
*** Starting crash test.
--- crash test: PASS
*** PASSED ALL TESTS
$
您可能会看到 Go RPC 包中的一些错误，如下所示

2019/12/16 13:27:09 rpc.Register: method "Done" has 1 input parameters; needs exactly three
忽略这些消息;将协调器注册为 RPC 服务器会检查其所有 方法适用于 RPC（有 3 个输入）;我们知道 Done 不是通过 RPC 调用的。
此外，根据您终止工作进程的策略，您可能会看到一些形式的错误

2024/02/11 16:21:32 dialing:dial unix /var/tmp/5840-mr-501: connect: connection refused
每个测试看到少量这些消息是正常的;当 worker 在 之后无法联系协调器 RPC 服务器时，就会出现它们 协调器已退出。

一些规则：
map 阶段应该将中间键划分为 nReduce reduce 任务的存储桶， 其中 nReduce 是 reduce 任务的数量 —— main/mrcoordinator.go 传递给 MakeCoordinator（） 的参数。 每个映射器都应为 消耗由 减少 任务。
worker 实现应将第 X 个 reduce 任务 mr-out-X 中的任务。
mr-out-X 文件应包含每个 Reduce 的一行 function 输出。该行应使用 Go “%v %v” 格式生成，并使用键和值调用。在 main/mrsequential.go 中查看注释为 “this is the correct format” 的行。 如果您的 implementation 与此格式有太多偏差，则测试脚本将失败。
您可以修改 mr/worker.go、mr/coordinator.go 和 mr/rpc.go。 您可以临时修改其他文件以进行测试，但请确保您的代码正常工作 与原始版本;我们将使用原始版本进行测试。
worker 应该将中间的 Map 输出放在当前 目录中，您的 worker 稍后可以在其中读取它们作为 Reduce 任务的输入。
main/mrcoordinator.go 期望 mr/coordinator.go 实现一个 Done（） 方法，当 MapReduce 作业完全完成时，该方法返回 true; 此时，mrcoordinator.go 将退出。
当作业完全完成后，工作进程应退出。 实现这一点的一种简单方法是使用 call（） 的返回值： 如果 worker 未能联系 coordinator，它可以假定 coordinator 已退出 因为 Job 已经完成，所以 worker 也可以终止。取决于您 design 中，您可能还会发现有一个 “please exit” 伪任务很有帮助 协调者可以给 worker 的。
提示
Guidance （指南） 页面有一些 有关开发和调试的提示。
一种开始的方法是修改 mr/worker.go 的 Worker（） 以向协调器发送 RPC 请求任务。然后 修改协调器以使用尚未启动的文件名进行响应 map 任务。然后修改 worker 以读取该文件并调用 application Map 函数，如 mrsequential.go 中所示。
应用程序 Map 和 Reduce 函数在运行时加载 使用 Go 插件包，从名称以 .so 结尾的文件中。
如果你更改了 mr/ 目录中的任何内容，你将 可能必须重新构建您使用的任何 MapReduce 插件，使用 像 go build -buildmode=plugin ../mrapps/wc.go
此实验室依赖于共享文件系统的工作程序。 当所有 worker 都在同一台机器上运行时，这很简单，但需要一个全局的 文件系统（如果 worker 在不同的机器上运行）。
中间文件的合理命名约定是 mr-X-Y， 其中 X 是 Map 任务编号，Y 是减少任务编号。
worker 的 map 任务代码将需要一种方法来存储中间 键/值对，以便正确读回 在 减少 任务 期间。一种可能性是使用 Go 的 encoding/json 包。自 将 JSON 格式的键/值对写入打开的文件：并读回此类文件：
  enc := json.NewEncoder(file)
  for _, kv := ... {
    err := enc.Encode(&kv)
  dec := json.NewDecoder(file)
  for {
    var kv KeyValue
    if err := dec.Decode(&kv); err != nil {
      break
    }
    kva = append(kva, kv)
  }
worker 的 map 部分可以使用 ihash（key） 函数 （在 worker.go 中）为给定的 key 选择 reduce 任务。
您可以从 mrsequential.go 中窃取一些代码进行读取 映射输入文件，用于对 Map 和 Reduce，用于将 Reduce 输出存储在文件中。
协调器作为 RPC 服务器，将是并发的;别忘了 以锁定共享数据。
使用 Go 的 race 检测器，与 go run -race 一起使用。test-mr.sh 在开头有一条注释，告诉您 如何使用 -race 运行它。 当我们对您的实验室进行评分时，我们不会使用 race 检测器。不过，如果您的代码有 races，则 即使没有 RACE 检测器。
worker 有时需要等待，例如 reduce can't start 直到最后一张地图结束。一种可能性是工人 定期向协调员询问 work， sleeping 随着时间的推移。Sleep（） 的请求。另一种可能性 是协调器中的相关 RPC 处理程序具有一个循环，该循环 等待，要么随着时间的推移。Sleep（） 或 sync 的 Sleep（） 方法。康德。去 在自己的线程中运行每个 RPC 的处理程序，因此一个 handler is waiting不需要阻止 coordinator 处理其他 RPC 的。
协调器无法可靠地区分崩溃的 worker， 还活着但由于某种原因而停滞的 worker， 以及正在执行但速度太慢而无法使用的 worker。 你能做的最好的事情就是让协调器等待 一段时间后，然后放弃并重新下发任务给 不同的 worker。对于此实验室，请让协调器等待 十秒;之后，协调器应假定 worker 具有 死亡（当然，它可能没有）。
如果您选择实现备份任务（第 3.6 节），请注意，我们会测试您的代码不会 在 worker 执行任务时调度无关的任务而不会崩溃。备份任务应仅 安排在相对较长的时间 （例如 10 秒） 之后。
要测试崩溃恢复，您可以使用 mrapps/crash.go 应用程序插件。它会在 Map 和 Reduce 函数中随机退出。
为了确保没有人在存在 崩溃时，MapReduce 论文提到了使用临时文件的技巧 并在它完全写入后自动重命名它。您可以使用 ioutil.TempFile（或 os.CreateTemp（如果您正在运行 Go 1.17 或更高版本）创建临时文件和 os。Rename 以原子方式重命名它。
test-mr.sh 在子目录 mr-tmp 中运行其所有进程，因此如果 出现问题，并且您想查看中间文件或输出文件，请查看那里。随意 临时修改 test-mr.sh 以在测试失败后退出，因此脚本不会 继续测试（并覆盖输出文件）。
test-mr-many.sh 连续运行 test-mr.sh 多次， 您可能希望这样做以发现低概率的 bug。 它将要运行的次数作为参数 测试。您不应并行运行多个 test-mr.sh 实例，因为 Coordinator 将重用相同的套接字，从而导致冲突。
Go RPC 只发送名称以大写字母开头的结构体字段。 子结构还必须具有大写的字段名称。
调用 RPC call（） 函数时， reply struct 应包含所有默认值。RPC 调用 应如下所示：在调用之前不设置任何 reply 字段。如果你 传递具有非默认字段的回复结构，RPC 系统可能会静默返回不正确的值。
  reply := SomeType{}
  call(..., &reply)
无学分挑战练习
实现您自己的 MapReduce 应用程序（参见 mrapps/* 中的示例），例如 Distributed Grep（MapReduce 论文的第 2.3 节）。

让您的 MapReduce 协调器和工作程序在单独的计算机上运行，就像在实践中一样。 您需要将 RPC 设置为通过 TCP/IP 而不是 Unix 套接字进行通信（请参阅 Coordinator.server（） 中注释掉的行），并使用共享文件读/写文件 系统。例如，您可以通过 ssh 连接到 MIT 的多台 Athena 集群计算机，这些计算机使用 AFS 共享文件;或者，您可以租用几个 AWS 实例并使用 S3 进行存储。