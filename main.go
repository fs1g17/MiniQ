package main

import (
	"fmt"

	"github.com/fs1g17/MiniQ/queue"
)

func main() {
	q := queue.Queue{}

	job1 := queue.Job{Name: "job1"}
	job2 := queue.Job{Name: "job2"}

	q.Enqueue(&job1)
	q.Enqueue(&job2)

	dq1 := q.Dequeue()
	fmt.Println(dq1.Name)

	dq2 := q.Dequeue()
	fmt.Println(dq2.Name)
}
