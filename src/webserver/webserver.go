package webserver

import (
	"log/slog"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

var Router *gin.Engine

func InitWebServer() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.Default()
	Router.Use(sloggin.New(slog.Default()))
	Router.Use(gin.Recovery())
}

func RunWebServer() error {
	err := Router.Run()
	return err
}