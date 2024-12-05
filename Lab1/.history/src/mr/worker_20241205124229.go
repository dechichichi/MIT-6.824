package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"time"
)

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
	filename := response.Filename
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ReadFile failed:", err)
		return
	}
	content, err := ioutil.ReadAll(file)
	//得到一个kv结构体数组
	KeyValueList := mapf(filename, string(content))

	rn:=response.ReducerNum
	HashKVMap := make(map[int][]KeyValue,rn)

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	keepFlag := true
	for keepFlag {
		task := GetTask()
		switch task.TaskType {
		case MapTask:
			{
				DoMapTask(mapf, &task)
				callDone()
			}
		case WaittingTask:
			{
				fmt.Println("Waitting.......")
				time.Sleep(1 * time.Second)
			}
		case ExitTask:
			{
				fmt.Println("Exit.......")
				keepFlag = false
			}
		}
	}
}
func GetTask() Task {
	args := TaskArgs{}
	reply := Task{}
	ok := call("Coordinator.GetTask", args, &reply)
	if ok {
		fmt.Println("GetTask:", reply)
	} else {
		fmt.Println("GetTask failed")
	}
	return reply
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

func callDone() {
	args := TaskArgs{}
	reply := Task{}
	ok := call("Coordinator.DoneTask", args, &reply)
	if ok {
		fmt.Println("DoneTask:", reply)
	} else {
		fmt.Println("DoneTask failed")
	}
	return
}
