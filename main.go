package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fs1g17/MiniQ/queue"
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
	log := make(chan string)

	jobStore := store.NewJobStore(pgDB)

	miniQ := queue.CreateMiniQ(log)
	miniQ.AddWorker(func(data store.AnyData) error {
		log <- fmt.Sprint(data)
		return nil
	})
	miniQ.AddWorker(func(data store.AnyData) error {
		log <- fmt.Sprint(data)
		return nil
	})

	job0 := &store.Job{
		Data: store.AnyData{"A": 1, "B": 2},
	}
	job1 := &store.Job{
		Data: store.AnyData{"A": 3, "B": 4},
	}
	err := jobStore.InsertJob(job0)
	if err != nil {
		fmt.Println("Failed to add job")
		return
	}
	err = jobStore.InsertJob(job1)
	if err != nil {
		fmt.Println("Failed to add job")
		return
	}

	miniQ.AddJob(job0)
	miniQ.AddJob(job1)

	go func() {
		time.Sleep(12 * time.Second)
		job2 := &store.Job{
			Data: store.AnyData{"A": 5, "B": 6},
		}
		err = jobStore.InsertJob(job2)
		if err != nil {
			fmt.Println("Failed to add job")
			return
		}

		miniQ.AddJob(job2)
	}()

	for {
		msg := <-log
		fmt.Println("LOG:", msg)
	}
}
