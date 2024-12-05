要跑动这个代码 第一步应该编译生成.so文件
 go build -buildmode=plugin   之后就有main.so文件

 在 Go 语言中，每个 .go 文件都有自己的作用域，即使它们属于同一个包。因此，您需要在 main.go 文件中也定义或导入 KeyValue 类型。