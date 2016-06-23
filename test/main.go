package main

import (
	"fmt"
	queue "github.com/RenatoGeh/goutils/queue"
	stack "github.com/RenatoGeh/goutils/stack"
)

func main() {
	queue := new(queue.Queue)
	stack := new(stack.Stack)

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	fmt.Printf("Queue has %d elements.\n", queue.Size())

	for !queue.Empty() {
		val, _ := (queue.Dequeue()).(int)
		fmt.Printf("Queue: %d\n", val)
	}

	fmt.Printf("Stack has %d elements.\n", stack.Size())

	for !stack.Empty() {
		val, _ := (stack.Pop()).(int)
		fmt.Printf("Stack: %d\n", val)
	}
}
