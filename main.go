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

	messages := make(chan string)

	worker0 := queue.Worker[MyData]{
		ID:     0,
		Status: queue.Idle,
		Work: func(d MyData) error {
			time.Sleep(5 * time.Second)
			t := time.Now()
			fmt.Println(t)
			fmt.Println(d)
			return nil
		},
		Queue:   &q,
		Channel: messages,
	}
	worker1 := queue.Worker[MyData]{
		ID:     1,
		Status: queue.Idle,
		Work: func(d MyData) error {
			time.Sleep(5 * time.Second)
			t := time.Now()
			fmt.Println(t)
			fmt.Println(d)
			return nil
		},
		Queue:   &q,
		Channel: messages,
	}

	q.Workers = append(q.Workers, &worker0, &worker1)

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

	go func() {
		time.Sleep(8 * time.Second)
		q.Enqueue(&queue.Job[MyData]{
			Name:   "job3",
			Status: queue.Queued,
			Data:   MyData{A: 5, B: 6},
		})
	}()

	for {
		msg := <-messages
		fmt.Println(msg)
	}
}
