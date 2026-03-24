package main

import (
	"fmt"
	"time"

	"github.com/fs1g17/MiniQ/migrations"
	"github.com/fs1g17/MiniQ/queue"
	"github.com/fs1g17/MiniQ/store"
	"github.com/joho/godotenv"
)

func setup() {
	godotenv.Load()
	fmt.Println(store.GetConnectionString())
	pgDB, err := store.Open()
	if err != nil {
		panic("not connected to db")
	}

	err = store.MigrateFs(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
}

func main() {
	setup()

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

	job0 := "job0"
	miniQ.AddJob(queue.NewJob(MyData{A: 1, B: 2}, &job0))
	job1 := "job1"
	miniQ.AddJob(queue.NewJob(MyData{A: 3, B: 4}, &job1))

	go func() {
		time.Sleep(12 * time.Second)
		job2 := "job2"
		miniQ.AddJob(queue.NewJob(MyData{A: 5, B: 6}, &job2))
	}()

	for {
		msg := <-log
		fmt.Println("LOG:", msg)
	}
}
