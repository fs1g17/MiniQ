# MiniQ

This is my little Queue project. The purpose of MiniQ is to learn about Go.

## How to Run

### PostgreSQL Container

You need to have Docker installed to run PostgreSQL docker container, which can be spun up with:

- `make up`

### Go Server

- `go mod download`
- `go mod tidy`
- `make server`

### Go Client

emulates a worker that operates on a job:

- `make client`

### Running tests with race condition

`CGO_ENABLED=1 go test -race -count=1 ./internal/queue/`
this `-count=1` forces a fresh run

## How it Works

The idea was simple: a small Go server that acts as a job queue, backed by PostgreSQL for persistence. There are only two endpoints — POST /addJob and GET /pollJob — and the whole thing is held together with long polling and a map of channels.

When a worker calls /pollJob, the server creates a channel for that connection and registers it in a clientsmap. The request just hangs there. The channel isn't carrying the job itself — it's only used to ping the goroutine once something gets added to the queue. Think of it as a doorbell, not a mailbox.

On the other side, POST /addJob expects a JSON body with a single data field typed as map[string]any — deliberately loose, to keep the server generic. When the endpoint is hit, the job gets written to a PostgreSQL jobstable and enqueued in memory. The queue uses a mutex to keep enqueue and dequeue operations thread-safe. Once that's done, the server loops over the registered clients and pings them — the first one to respond gets the job.

If a job is already waiting when a worker calls /pollJob, it gets dequeued immediately: the job status is updated to in_progressin Postgres and returned with a 200. If the queue is empty, the request stays open. If it times out before anything arrives, the worker just retries — reconnecting and waiting again. This keeps things reactive without any polling delay: the right worker gets the job the moment it's available.

It's a toy, but it forced me to actually think about goroutines, mutexes, and channel signaling rather than just reading about them. Getting the first-available-worker assignment right took a few attempts — the initial version would occasionally assign the same job to two workers when they both happened to poll at the same time.
