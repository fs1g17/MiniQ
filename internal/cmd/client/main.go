package main

import (
	"log"
	"time"

	"github.com/fs1g17/MiniQ/internal/client"
)

func main() {
	client.PollForJob("http://localhost:8080", func(jobResponse client.JobResponse) {
		log.Printf("Received: %+v", jobResponse)
		log.Println("working on job...")
		time.Sleep(5 * time.Second)
		log.Println("finished working on job")
	})
}
