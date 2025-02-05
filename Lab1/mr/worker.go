package mr

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"strings"
)

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func DoMapTask(mapf func(string, string) []KeyValue, response *Task) {
	filename := response.FileName
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ReadFile failed:", err)
		return
	}
	content, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		fmt.Println("ReadFile failed:", err)
		return
	}
	//得到一个kv结构体数组
	KeyValueList := mapf(filename, string(content))

	//缓存
	rn := response.ReduceNum
	fileCache := make(map[string][]KeyValue)
	for _, kv := range KeyValueList {
		k := ihash(kv.Key) % rn
		oname := "mr-tmp-" + strconv.Itoa(k)
		fileCache[oname] = append(fileCache[oname], kv)
	}

	// 写入文件
	for oname, kvs := range fileCache {
		ofile, err := os.Create(oname)
		if err != nil {
			fmt.Println("Create file failed:", err)
			continue
		}
		defer ofile.Close()
		for _, kv := range kvs {
			fmt.Fprintf(ofile, "%v %v\n", kv.Key, kv.Value)
		}
	}
}

func DoReduceTask(reducef func(string, []string) string, response *Task) {
	// 读取所有与当前 Reduce 任务相关的中间文件
	var intermediate []KeyValue
	for k := 0; k < response.ReduceNum; k++ {
		tmpFileName := fmt.Sprintf("mr-tmp-%d", k)
		file, err := os.Open(tmpFileName)
		if err != nil {
			log.Printf("Failed to open file %s: %v", tmpFileName, err)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) != 2 {
				log.Printf("Invalid line in file %s: %s", tmpFileName, line)
				continue
			}
			key, value := fields[0], fields[1]
			intermediate = append(intermediate, KeyValue{Key: key, Value: value})
		}
	}

	// 对中间结果进行排序
	sort.Sort(ByKey(intermediate))

	// 生成当前 Reduce 任务的输出文件名
	outputFileName := fmt.Sprintf("mr-out-%d", response.TaskID)
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Printf("Failed to create output file %s: %v", outputFileName, err)
		return
	}
	defer outputFile.Close()

	// 归并操作
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := make([]string, j-i)
		for k := i; k < j; k++ {
			values[k-i] = intermediate[k].Value
		}
		output := reducef(intermediate[i].Key, values)
		fmt.Fprintf(outputFile, "%v %v\n", intermediate[i].Key, output)
		i = j
	}

	log.Printf("Reduce task %d completed, output file: %s", response.TaskID, outputFileName)
}

func Worker(mapf func(string, string) []KeyValue, reducef func(string, []string) string) bool {
	var tasks []Task
	args := &TaskArgs{}
	reply := &[]Task{}

	for {
		// 调用 Coordinator 的 GetTasks 方法获取任务
		if success := call("Coordinator.GetTasks", args, reply); !success {
			log.Fatalf("Worker: GetTasks Failed")
		}

		tasks = *reply
		if len(tasks) == 0 {
			// 如果没有任务，说明所有任务都已完成
			log.Println("No more tasks to process. All tasks are done.")
			break
		}

		for _, task := range tasks {
			switch task.TaskType {
			case MapTask:
				DoMapTask(mapf, &task)
				if !callDone(task.TaskID) {
					log.Printf("Failed to mark task %d as done", task.TaskID)
				}
			case ReduceTask:
				DoReduceTask(reducef, &task)
				if !callDone(task.TaskID) {
					log.Printf("Failed to mark task %d as done", task.TaskID)
				}
			case Done:
				// Done 任务类型表示当前任务列表中的任务都已处理完毕
				log.Printf("Task %d is done", task.TaskID)
			}
		}
	}
	return true
}

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()
	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func callDone(taskID int) bool {
	args := TaskArgs{TaskID: taskID}
	reply := Task{}
	ok := call("Coordinator.DoneTask", &args, &reply)
	if ok {
		log.Printf("DoneTask response: %+v", reply)
		return true
	} else {
		log.Println("DoneTask failed")
		return false
	}
}
