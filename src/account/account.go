package account

import (
	"blacklizardcode/sine/auth"
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/webserver"
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InitAccountRoutes() {
	{
		account := webserver.Router.Group("/account")
		account.Use(auth.AuthMiddleWare())
		account.GET("balance", balanceHandler)
	}
}

func balanceHandler(c *gin.Context) {
	jwtCookie, err := c.Cookie("jwt")
	if err != nil {
		slog.Error("missing jwt cookie", "error", err)
		c.Status(http.StatusUnauthorized)
		return
	}

	claims, err := auth.JwtToJwtClaims(jwtCookie)
	if err != nil {
		slog.Error("failed to parse jwt claims", "error", err)
		c.Status(http.StatusUnauthorized)
		return
	}

	subject, err := claims.GetSubject()
	if err != nil {
		slog.Error("failed to get subject from claims", "error", err)
		c.Status(http.StatusUnauthorized)
		return
	}
	userid, err := strconv.Atoi(subject)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var balance float64
	database.DB.QueryRow(context.Background(), "SELECT balance FROM users WHERE userid=$1", userid).Scan(&balance)

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
