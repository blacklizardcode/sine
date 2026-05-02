package main

import (
	"blacklizardcode/sine/account"
	"blacklizardcode/sine/auth"
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/transaction"
	"blacklizardcode/sine/webserver"

	"log/slog"
	"os"

	"github.com/SladkyCitron/slogcolor"
)



func main() {
	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions)))
	slog.SetLogLoggerLevel(slog.LevelInfo)

	slog.Info("starting")

	// initializes values: run first
	database.InitDB()
	webserver.InitWebServer()
	

	// routes: run after InitWebServer and InitDB
	auth.InitAuthRoutes()
	account.InitAccountRoutes()
	transaction.InitTransactionRoutes()

	// blocking: run last
	webserver.RunWebServer()
}
