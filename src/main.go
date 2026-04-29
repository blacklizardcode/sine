package main

import (
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/webserver"
	"blacklizardcode/sine/auth"
	"log/slog"
	"github.com/SladkyCitron/slogcolor"
	"os"
)



func main() {
	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions)))

	slog.Info("starting")

	// initializes values: run first
	database.InitDB()
	webserver.InitWebServer()

	// routes: run after InitWebServer and InitDB
	err := auth.InitUserRoutes()
	if err != nil {
		slog.Error("%s", err.Error())
		return
	}

	// blocking: run last
	webserver.RunWebServer()
}
