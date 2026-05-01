package auth

import (
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/webserver"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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

type JwtClaims struct {
	sub string
	exp string
	iat string
}

func InitAuthRoutes() {

	{
		users := webserver.Router.Group("/auth")
		users.POST("/register", registerHandler)
		users.POST("/login", loginHandler)
	}

}

func registerHandler(c *gin.Context) {

	var json userForm

	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.Status(http.StatusBadRequest)
		slog.Error("failed to bind JSON", "error", err)
		return
	}

	passwordHash := sha256.New()
	_, err = passwordHash.Write([]byte(json.Password))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error("failed to hash password", "error", err)
		return
	}

	passwordSum := passwordHash.Sum(nil)
	passwordHex := hex.EncodeToString(passwordSum)

	_, err = database.DB.Exec(context.Background(), "insert into users (username, password) values ($1, $2)", json.Username, passwordHex)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error("failed to insert user", "error", err)
		return
	}
}

func loginHandler(c *gin.Context) {
	var json userForm
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error("failed to bind JSON", "error", err)
		return
	}

	passwordHash := sha256.New()
	_, err = passwordHash.Write([]byte(json.Password))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error("failed to hash password", "error", err)
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

	var userId int
	err = database.DB.QueryRow(context.Background(), "SELECT userid from users WHERE username=$1", json.Username).Scan(&userId)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": strconv.Itoa(userId),
			"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			"iat": jwt.NewNumericDate(time.Now()),
		})
	signedKey, err := t.SignedString(jwtKey)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Error("failed to sign JWT", "error", err)
		return
	}

	// set a host-only cookie with correct max-age (seconds) and httpOnly
	c.SetCookie("jwt", signedKey, 24*3600, "/", "", false, true)
	c.Status(http.StatusOK)
}

func AuthMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
    	jwtCookie, err := c.Cookie("jwt")
		if err != nil {
			slog.Error("missing jwt cookie", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

    	// parse and validate JWT
    	parsed, err := jwt.Parse(jwtCookie, func(token *jwt.Token) (interface{}, error) {
    		// ensure token uses expected signing method
    		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
    			return nil, jwt.ErrTokenSignatureInvalid
    		}
    		return jwtKey, nil
    	})
		if err != nil || !parsed.Valid {
			slog.Error("failed to parse/validate jwt", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}


    	// Pre-handler phase
    	c.Next()
	}
}

func JwtToJwtClaims(jwtString string) (jwt.Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(jwtString, jwt.MapClaims{})

		
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("")
	}


	return claims, nil
}