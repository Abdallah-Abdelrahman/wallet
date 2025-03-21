package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"wallet/internal/database"
	"wallet/internal/services"
)

type Server struct {
	port int

	db                 database.Service
	AccountService     services.AccountService
	TransactionService services.TransactionService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()

	NewServer := &Server{
		port:               port,
		db:                 db,
		AccountService:     services.NewAccountService(db.GetDB()),
		TransactionService: services.NewTransactionService(db.GetDB()),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
