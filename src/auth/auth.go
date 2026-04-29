package auth

import (
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/webserver"
	"context"
	"log/slog"
	"net/http"
	"time"

	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("smtn")

type userForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

func InitUserRoutes() error{
	
	{
		users := webserver.Router.Group("/auth")
		users.POST("/register", registerHandler)
		users.POST("/login", loginHandler)
	}
	
	return nil
}

func registerHandler(c *gin.Context) {
	

	var json userForm

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

func loginHandler(c *gin.Context) {
	var json userForm 
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.Status(http.StatusInternalServerError)
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
	passwordSum := passwordHash.Sum(nil)
	passwordHex := hex.EncodeToString(passwordSum)

	var savedPasswordHash string
	_ = database.DB.QueryRow(context.Background(), "SELECT password FROM users WHERE username=$1", json.Username).Scan(&savedPasswordHash)
	if savedPasswordHash != passwordHex {
		c.Status(http.StatusUnauthorized)
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{
			"sub": json.Username,
			"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			"iat": jwt.NewNumericDate(time.Now()),
	})
	signedKey, err := t.SignedString(jwtKey)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error("%s", err.Error())
		return
	}

	c.SetCookie("jwt", signedKey, int(time.Hour) * 24, "/", "localhost", false, false)
	c.Status(http.StatusOK)
}