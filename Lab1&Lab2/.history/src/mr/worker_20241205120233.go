package mr

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
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

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
keepFlag:=true
	for keepFlag{
		task := GetTask()
		if task.Type == "map" {
			fmt.Println("worker: map task", task.Id)
			results := mapf(task.Data, task.Id)
			for _, kv := range results {
				reduceId := ihash(kv.Key) % NReduce
				fmt.Println("worker: emit", kv, "to reduce", reduceId)
				emit(kv, reduceId)
			}
			fmt.Println("worker: done with map task", task.Id)
		} else if task.Type == "reduce" {
			fmt.Println("worker: reduce task", task.Id)
			values := make([]string, NMap)
			for i := 0; i < NMap; i++ {
				key := fmt.Sprintf("%d_%d", task.Id, i)
				fmt.Println("worker: fetch", key)
				value := fetch(key)
				values[i] = value
			}
			result := reducef(task.Data, values)
			fmt.Println("worker: emit", result, "to coordinator")
			emitReduce(task.Id, result)
			fmt.Println("worker: done with reduce task", task.Id)
		} else if task.Type == "stop":
			fmt.Println("worker: stop task")
			keepFlag=false
		} else {
			fmt.Println("worker: unknown task type", task.Type)
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
