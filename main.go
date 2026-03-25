package main

import (
	"database/sql"
	"fmt"

	"github.com/fs1g17/MiniQ/store"
	"github.com/joho/godotenv"
)

func setup() *sql.DB {
	godotenv.Load()
	fmt.Println(store.GetConnectionString())
	pgDB, err := store.Open()
	if err != nil {
		panic("not connected to db")
	}
	return pgDB
}

func main() {
	pgDB := setup()

	type MyData struct {
		A int
		B int
	}

	jobStore := store.NewJobStore(pgDB)
	job := store.Job{
		Data: store.AnyData{"A": 1, "B": 2},
	}
	err := jobStore.InsertJob(&job)
	if err != nil {
		fmt.Println("Failed to insert job")
		return
	}

	job2, err := jobStore.GetJob(job.ID)
	if err != nil {
		fmt.Println("failed to get job")
		return
	}

	fmt.Printf("ID %d\n", job2.ID)
	fmt.Printf("Status %d\n", job2.Status)
	fmt.Printf("Data %v\n", job2.Data)
	fmt.Printf("Attempts %d\n", job2.Attempts)
}

// func main2() {
// 	setup()

// 	type MyData struct {
// 		A int
// 		B int
// 	}

// 	log := make(chan string)

// 	miniQ := queue.CreateMiniQ[MyData](log)
// 	miniQ.AddWorker(func(data MyData) error {
// 		log <- fmt.Sprint(data)
// 		return nil
// 	})
// 	miniQ.AddWorker(func(data MyData) error {
// 		log <- fmt.Sprint(data)
// 		return nil
// 	})

// 	miniQ.AddJob(queue.NewJob(MyData{A: 1, B: 2}))
// 	miniQ.AddJob(queue.NewJob(MyData{A: 3, B: 4}))

// 	go func() {
// 		time.Sleep(12 * time.Second)
// 		miniQ.AddJob(queue.NewJob(MyData{A: 5, B: 6}))
// 	}()

// 	for {
// 		msg := <-log
// 		fmt.Println("LOG:", msg)
// 	}
// }
