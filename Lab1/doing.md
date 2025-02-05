mrcoordinator.go
启动一个Makcoordinater 将每个任务分配
传回Coordinator
等待Coordinator的Done()==true



coordinator.go 协调器
MakeCoordinator():
对于每个文件 启动一个handler 处理
进行Coordinator的server
返回这个Coordinator给main函数

handler():
添加任务到任务队列

server():
注册并且发布任务

Done():
任务结束后返回true  

worker.go 工作区域
LoadPlugin():
加载插件 获取Map和Reduce函数

worker.go 工作区域
Worker():
根据任务类型进行任务
