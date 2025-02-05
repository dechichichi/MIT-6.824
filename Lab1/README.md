要跑动这个代码 第一步应该编译生成.so文件
 go build -buildmode=plugin   之后就有main.so文件

 在 Go 语言中，每个 .go 文件都有自己的作用域，即使它们属于同一个包。因此，您需要在 main.go 文件中也定义或导入 KeyValue 类型。

 go run main.go wc.so pg*.txt 这个命令在Go语言的MapReduce框架中代表了启动一个MapReduce作业，其中：

go run main.go：这部分表示使用Go语言的运行命令来执行名为main.go的程序。

wc.so：这是一个动态链接库（在Unix系统中，共享库文件通常以.so为后缀），它包含了MapReduce作业中使用的Map和Reduce函数的实现。在这个上下文中，wc.so很可能是wc.go文件编译后生成的，wc.go实现了统计单词数量（Word Count）的Map和Reduce逻辑。

pg*.txt：这是一个文件匹配模式，表示所有以pg开头并以.txt结尾的文本文件。这些文件将作为MapReduce作业的输入数据。在MapReduce框架中，这些文件通常被分割成多个chunk，每个chunk作为一个Map任务的输入。

综上所述，整个命令的意思是：运行main.go程序，加载wc.so插件来处理所有匹配pg*.txt模式的文本文件，执行Word Count MapReduce作业。

wc.go  文本索引器