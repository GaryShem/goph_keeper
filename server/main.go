package main

import (
	"goph_keeper/goph_server/internal/app"
	"goph_keeper/goph_server/internal/logging"
)

func main() {
	if err := app.Run(); err != nil {
		logging.Log().Fatal("server failed: %v", err)
	}
}
