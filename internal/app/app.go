package app

import (
	"database/sql"
	"fmt"

	"github.com/fs1g17/MiniQ/internal/api"
	"github.com/fs1g17/MiniQ/internal/queue"
	"github.com/fs1g17/MiniQ/internal/store"
	"github.com/joho/godotenv"
)

type App struct {
	MiniQ   *queue.MiniQ
	Handler *api.QueueHandler
}

func setup() *sql.DB {
	godotenv.Load()
	fmt.Println(store.GetConnectionString())
	pgDB, err := store.Open()
	if err != nil {
		panic("not connected to db")
	}
	return pgDB
}

func NewApp() *App {
	pgDB := setup()

	jobStore := store.NewJobStore(pgDB)

	miniQ := queue.CreateMiniQ(jobStore)

	handler := api.NewQueueHandler(miniQ)
	return &App{
		MiniQ:   miniQ,
		Handler: handler,
	}
}
