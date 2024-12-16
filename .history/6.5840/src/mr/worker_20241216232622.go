package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"time"
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
	filename := response.Filename
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ReadFile failed:", err)
		return
	}
	content, err := ioutil.ReadAll(file)
	//得到一个kv结构体数组
	KeyValueList := mapf(filename, string(content))

	rn := response.ReducerNum
	HashKVMap := make(map[int][]KeyValue, rn)
	for _, kv := range KeyValueList {
		hash := ihash(kv.Key) % rn
		HashKVMap[hash] = append(HashKVMap[hash], kv)
	}
	for i := 0; i < rn; i++ {
		oname := "mr-tmp-" + strconv.Itoa(response.TaskID) + "-" + strconv.Itoa(i)
		ofile, _ := os.Create(oname)
		enc := json.NewEncoder(ofile)
		for _, kv := range HashKVMap[i] {
			enc.Encode(kv)
		}
		ofile.Close()
	}
}

func DoReduceTask(reducef func(string, []string) string, response *Task) {
	rn := response.ReducerNum
	for k := 0; k < rn; k++ {
		sort.Sort(ByKey(response.intermediate))
		oname := "mr-out-" + string(k)
		ofile, _ := os.Create(oname)
		i := 0
		for i < len(response.intermediate) {
			j := i + 1
			for j < len(response.intermediate) && response.intermediate[j].Key == response.intermediate[i].Key {
				j++
			}
			values := []string{}
			for k := i; k < j; k++ {
				values = append(values, response.intermediate[k].Value)
			}
			output := reducef(response.intermediate[i].Key, values)

			// this is the correct format for each line of Reduce output.
			fmt.Fprintf(ofile, "%v %v\n", response.intermediate[i].Key, output)
			i = j
		}
	}

}

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) bool {
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
		case ReduceTask:
			{
				DoReduceTask(reducef, &task)
				callDone()
			}
		}
	}

	return true
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
	reply.TaskType = MapTask
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

func callDone() bool {
	args := TaskArgs{}
	reply := Task{}
	ok := call("Coordinator.DoneTask", args, &reply)
	if ok {
		fmt.Println("DoneTask:", reply)
		return true
	} else {
		fmt.Println("DoneTask failed")
		return false
	}
}
