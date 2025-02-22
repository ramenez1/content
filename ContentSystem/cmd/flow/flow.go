package main

import (
	"encoding/json"
	"fmt"
	flow "github.com/s8sg/goflow/flow/v1"
	goflow "github.com/s8sg/goflow/v1"
	"math/rand"
	"strconv"
)

func Input(data []byte, option map[string][]string) ([]byte, error) {
	var input map[string]int
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	outputInt := input["input"]
	return []byte(strconv.Itoa(outputInt)), nil
}

func AddOne(data []byte, option map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num + rand.Intn(10) + 1
	fmt.Println("AddOne:", outputInt)
	return []byte(strconv.Itoa(outputInt)), nil
}

func AddTwo(data []byte, option map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num + rand.Intn(10) + 1
	fmt.Println("AddTwo:", outputInt)
	return []byte(strconv.Itoa(outputInt)), nil
}

func Aggregator(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("Aggregate:", string(data))
	return data, nil
}

// Expand10 扩大十倍
func Expand10(data []byte, option map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num * 10
	fmt.Println("Expand10:", outputInt)
	return []byte(strconv.Itoa(outputInt)), nil
}

// Expand100 扩大100倍
func Expand100(data []byte, option map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num * 100
	fmt.Println("Expand100:", outputInt)
	return []byte(strconv.Itoa(outputInt)), nil
}

func Output(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("Output:", string(data))
	return []byte("ok"), nil
}

func MyFlow(workflow *flow.Workflow, context *flow.Context) error {
	dag := workflow.Dag()

	//构建节点
	dag.Node("input", Input)
	dag.Node("add-one", AddOne)
	dag.Node("add-two", AddTwo)
	dag.Node("aggregate", Aggregator, flow.Aggregator(func(m map[string][]byte) ([]byte, error) {
		a, _ := strconv.Atoi(string(m["add-one"]))
		b, _ := strconv.Atoi(string(m["add-two"]))
		num := a + b
		fmt.Println("Aggregate:", num)
		return []byte(strconv.Itoa(num)), nil
	}))

	branches := dag.ConditionalBranch("judge", []string{"moreThan", "lessThan"}, func(bytes []byte) []string {
		num, _ := strconv.Atoi(string(bytes))
		fmt.Println("ConditionalBranch = ", num)
		if num > 10 {
			return []string{"moreThan"}
		}
		return []string{"lessThan"}
	}, flow.Aggregator(func(m map[string][]byte) ([]byte, error) {
		if v, ok := m["moreThan"]; ok {
			return v, nil
		}
		if v, ok := m["lessThan"]; ok {
			return v, nil
		}
		return nil, nil
	}))
	branches["moreThan"].Node("expand10", Expand10)
	branches["lessThan"].Node("expand100", Expand100)
	dag.Node("output", Output)

	//构建边（依赖关系）
	dag.Edge("input", "add-one")
	dag.Edge("input", "add-two")
	dag.Edge("add-one", "aggregate")
	dag.Edge("add-two", "aggregate")
	dag.Edge("aggregate", "judge")
	dag.Edge("judge", "output")

	return nil
}

func main() {
	fs := goflow.FlowService{
		Port:              8080,
		RedisURL:          "localhost:6379",
		WorkerConcurrency: 5,
	}
	err := fs.Register("add-flow", MyFlow)
	if err != nil {
		return
	}
	err = fs.Start()
	if err != nil {
		panic(err)
	}
}
