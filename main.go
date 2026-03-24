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

	log := make(chan string)

	miniQ := queue.CreateMiniQ[MyData](log)
	miniQ.AddWorker(func(data MyData) error {
		log <- fmt.Sprint(data)
		return nil
	})
	miniQ.AddWorker(func(data MyData) error {
		log <- fmt.Sprint(data)
		return nil
	})

	miniQ.AddJob(&queue.Job[MyData]{
		Name:   "job0",
		Status: queue.Queued,
		Data:   MyData{A: 1, B: 2},
	})
	miniQ.AddJob(&queue.Job[MyData]{
		Name:   "job1",
		Status: queue.Queued,
		Data:   MyData{A: 3, B: 4},
	})

	go func() {
		time.Sleep(12 * time.Second)
		miniQ.AddJob(&queue.Job[MyData]{
			Name:   "job2",
			Status: queue.Queued,
			Data:   MyData{A: 5, B: 6},
		})
	}()

	for {
		msg := <-log
		fmt.Println("LOG:", msg)
	}
}
