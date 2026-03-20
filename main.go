package main

import (
	"fmt"
	"time"

	"github.com/fs1g17/MiniQ/queue"
)

func main() {
	type MyData struct {
		A int
		B int
	}

	q := queue.Queue[MyData]{}

	job1 := queue.Job[MyData]{
		Name:   "job1",
		Status: queue.Queued,
		Data:   MyData{A: 1, B: 2},
	}
	job2 := queue.Job[MyData]{
		Name:   "job2",
		Status: queue.Queued,
		Data:   MyData{A: 3, B: 4},
	}

	q.Enqueue(&job1)
	q.Enqueue(&job2)

	messages := make(chan string)

	worker0 := queue.Worker[MyData]{
		ID: 0,
		Work: func(d MyData) error {
			t := time.Now()
			fmt.Println(t)
			fmt.Println(d)
			return nil
		},
		Queue:   &q,
		Channel: messages,
	}
	worker1 := queue.Worker[MyData]{
		ID: 1,
		Work: func(d MyData) error {
			t := time.Now()
			fmt.Println(t)
			fmt.Println(d)
			return nil
		},
		Queue:   &q,
		Channel: messages,
	}

	go worker0.Perform()
	go worker1.Perform()

	for {
		msg := <-messages
		fmt.Println(msg)
	}
}
