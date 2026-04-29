package main

import (
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/webserver"
	"blacklizardcode/sine/users"
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
	err := users.InitUserRoutes()
	if err != nil {
		slog.Error("%s", err.Error())
		return
	}

	// blocking: run last
	webserver.RunWebServer()
}
