package main

import (
	"fmt"
	"net/http"

	"github.com/fs1g17/MiniQ/internal/app"
	"github.com/labstack/echo/v5"
)

func main() {
	app := app.NewApp()

	e := echo.New()

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/addJob", app.Handler.HandleAddJob)
	e.POST("/completeJob", app.Handler.HandleCompleteJob)
	e.GET("/getJob", app.Handler.HandleGetJob)
	e.GET("/pollJob", app.Handler.HandlePollJob)
	e.GET("/getJobs", app.Handler.HandleGetJobs)

	if err := e.Start(":8080"); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
