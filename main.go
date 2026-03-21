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

	messages := make(chan string)

	miniQ := queue.CreateMiniQ[MyData]()
	miniQ.AddWorker(func(data MyData) error {
		messages <- fmt.Sprint(data)
		return nil
	}, messages)
	miniQ.AddWorker(func(data MyData) error {
		messages <- fmt.Sprint(data)
		return nil
	}, messages)

	miniQ.AddJob(&queue.Job[MyData]{
		Name:   "job1",
		Status: queue.Queued,
		Data:   MyData{A: 1, B: 2},
	})
	miniQ.AddJob(&queue.Job[MyData]{
		Name:   "job2",
		Status: queue.Queued,
		Data:   MyData{A: 3, B: 4},
	})

	go func() {
		time.Sleep(8 * time.Second)
		miniQ.AddJob(&queue.Job[MyData]{
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
