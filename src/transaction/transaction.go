package transaction

import (
	"blacklizardcode/sine/auth"
	"blacklizardcode/sine/database"
	"blacklizardcode/sine/webserver"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func InitTransactionRoutes() {
	{
		transaction := webserver.Router.Group("/transaction")
		transaction.Use(auth.AuthMiddleWare())
		transaction.POST("/transfer", transferHandler)
		transaction.GET("/transactioninfo", transactionInfoHandler)
		transaction.GET("/transactionhistory", transactionHistoryHandler)
	}
}

func transferHandler(c *gin.Context) {
	type transferStruct struct {
		To_account string  `json:"to_account" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
	}

	jwtCookie, err := c.Cookie("jwt")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}

	var transferInfo transferStruct
	err = c.ShouldBindBodyWithJSON(&transferInfo)
	if err != nil {
		c.Status(http.StatusBadRequest)
		slog.Debug(err.Error())
		return
	}

	if transferInfo.Amount <= 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	claims, err := auth.JwtToJwtClaims(jwtCookie)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}
	subject, err := claims.GetSubject()
	userid, err := strconv.Atoi(subject)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}

	var senderBalance float64
	database.DB.QueryRow(context.Background(), "SELECT userid, balance FROM users WHERE userid=$1", userid).Scan(&senderBalance)

	if transferInfo.Amount > senderBalance {
		c.Status(http.StatusBadRequest)
		return
	}

	var recieverUserId int
	database.DB.QueryRow(context.Background(), "SELECT userid FROM users WHERE username=$1", transferInfo.To_account).Scan(&recieverUserId)

	if userid == recieverUserId {
		c.Status(http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Debug(err.Error())
	}

	defer tx.Rollback(context.Background())

	var transferId int
	err = tx.QueryRow(context.Background(), `INSERT INTO transactions (from_account, to_account, amount) VALUES ($1, $2, $3) RETURNING transaction_id`, userid, recieverUserId, transferInfo.Amount).Scan(&transferId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Debug(err.Error())
		return
	}

	_, err = tx.Exec(context.Background(), `UPDATE users SET balance = balance - $1 WHERE userid=$2`, transferInfo.Amount, userid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Debug(err.Error())
		return
	}

	_, err = tx.Exec(context.Background(), `UPDATE users SET balance = balance + $1 WHERE userid=$2`, transferInfo.Amount, recieverUserId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.Debug(err.Error())
		return
	}

	tx.Commit(context.Background())

	c.JSON(http.StatusOK, gin.H{"transfer_id": transferId})
}

func transactionInfoHandler(c *gin.Context) {

	jwtCookie, err := c.Cookie("jwt")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}

	claims, err := auth.JwtToJwtClaims(jwtCookie)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}
	subject, err := claims.GetSubject()
	userid, err := strconv.Atoi(subject)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}

	// Get transaction_id from query parameter
	txIdStr := c.Query("transaction_id")
	if txIdStr == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	var transaction_id int
	var from_account int
	var to_account int
	var amount float64
	var timestamp time.Time

	err = database.DB.QueryRow(context.Background(), "SELECT transaction_id, from_account, to_account, amount, timestamp FROM transactions WHERE transaction_id=$1", txIdStr).Scan(&transaction_id, &from_account, &to_account, &amount, &timestamp)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if userid != from_account && userid != to_account {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": transaction_id,
		"from_account":   from_account,
		"to_account":     to_account,
		"amount":         amount,
		"timestamp":      timestamp,
	})
}

func transactionHistoryHandler(c *gin.Context) {
	type transaction struct {
		TransactionID int       `json:"transaction_id"`
		FromAccount   int       `json:"from_account"`
		ToAccount     int       `json:"to_account"`
		Amount        float64   `json:"amount"`
		Timestamp     time.Time `json:"timestamp"`
	}

	strLimit := c.DefaultQuery("limit", "10")
	strOffset := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		slog.Debug(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(strOffset)
	if err != nil {
		slog.Debug(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	jwtCookie, err := c.Cookie("jwt")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}

	claims, err := auth.JwtToJwtClaims(jwtCookie)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}
	subject, err := claims.GetSubject()
	userId, err := strconv.Atoi(subject)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		slog.Debug(err.Error())
		return
	}

	rows, err := database.DB.Query(context.Background(), "SELECT transaction_id, from_account, to_account, amount, timestamp FROM transactions WHERE from_account=$1 OR to_account=$1 LIMIT $2 OFFSET $3", userId, limit, offset)
	if err != nil {
		slog.Debug(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	transactions := []transaction{}

	for rows.Next() {
		var t transaction
		if err := rows.Scan(&t.TransactionID, &t.FromAccount, &t.ToAccount, &t.Amount, &t.Timestamp); err != nil {
			slog.Debug(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}

	rows.Close()

	c.JSON(200, gin.H{
		"transactions": transactions,
	})
}
