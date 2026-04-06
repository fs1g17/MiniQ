package main

import "github.com/fs1g17/MiniQ/internal/client"

func main() {
	client.PollForJob("http://localhost:8080")
}
