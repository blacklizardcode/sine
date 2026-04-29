package users

import (
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/webserver"
	"context"
	"log/slog"
	"net/http"

	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes() error{
	
	{
		users := webserver.Router.Group("/users")
		users.POST("/register", registerHandler)
	}
	
	return nil
}

func registerHandler(c *gin.Context) {
	type registerForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var json registerForm

	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.Status(http.StatusBadRequest)
		slog.Error("%s", err.Error())
		return
	}

	passwordHash := sha256.New()
	_, err = passwordHash.Write([]byte(json.Password))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error("%s", err.Error())
		return
	}

	// finalize hash and encode as hex string
	passwordSum := passwordHash.Sum(nil)
	passwordHex := hex.EncodeToString(passwordSum)

	_, err = database.DB.Exec(context.Background(), "insert into users (username, password) values ($1, $2)", json.Username, passwordHex)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error("%s", err.Error())
		return
	}
}