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

	miniQ := queue.CreateMiniQ(jobStore, log)
	miniQ.AddWorker(func(data store.AnyData) error {
		log <- fmt.Sprint(data)
		return nil
	})
	miniQ.AddWorker(func(data store.AnyData) error {
		log <- fmt.Sprint(data)
		return nil
	})

	err := miniQ.AddJob(&store.AnyData{"A": 1, "B": 2})
	if err != nil {
		fmt.Errorf("Failed to add job %w", err)
		return
	}
	err = miniQ.AddJob(&store.AnyData{"A": 3, "B": 4})
	if err != nil {
		fmt.Errorf("Failed to add job %w", err)
		return
	}

	go func() {
		time.Sleep(12 * time.Second)

		err = miniQ.AddJob(&store.AnyData{"A": 5, "B": 6})
		if err != nil {
			fmt.Errorf("Failed to add job %w", err)
			return
		}
	}()

	for {
		msg := <-log
		fmt.Println("LOG:", msg)
	}
}
