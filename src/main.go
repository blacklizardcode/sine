package main

import (
	"blacklizardcode/sine/account"
	"blacklizardcode/sine/auth"
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/webserver"
	"blacklizardcode/sine/auth"
	"log/slog"
	"os"

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
	auth.InitAuthRoutes()
	account.InitAccountRoutes()

	// blocking: run last
	webserver.RunWebServer()
}
