package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"virtual-wallet/internal/handlers"
	"virtual-wallet/internal/repository"
	"virtual-wallet/internal/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("HOST_IP_ADDRESS"), os.Getenv("HOST_PORT"), os.Getenv("DB_NAME"))

	db, errDb := sql.Open("postgres", dsn)
	if errDb != nil {
		panic(errDb)
	}

	defer func(db *sql.DB) {
		errDbClose := db.Close()
		if errDbClose != nil {
			panic(errDbClose)
		}
	}(db)

	errPing := db.Ping()
	if errPing != nil {
		panic(errPing)
	}

	mux := http.NewServeMux()
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(os.Getenv("JWT_SECRET"), userRepo)
	userHandler := handlers.NewUserHandler(os.Getenv("JWT_SECRET"), userSvc)

	walletRepo := repository.NewWalletRepository(db)
	walletSvc := service.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletSvc)

	mux.HandleFunc("POST /api/register/", userHandler.RegisterUser)
	mux.HandleFunc("POST /api/login/", userHandler.LoginUser)

	mux.HandleFunc("GET /api/wallets/", userHandler.AuthMiddleware(walletHandler.GetWalletsByProfileID))
	mux.HandleFunc("POST /api/wallet/create", userHandler.AuthMiddleware(walletHandler.CreateWallet))

	fmt.Println("Server started")

	errHTTP := http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), mux)

	if errHTTP != nil {
		panic(errHTTP)
	}

}
