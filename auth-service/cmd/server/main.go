package main

import (
	"net/http"
	"os"

	"github.com/Bkgediya/feed_system/auth-service/internal/api"
	"github.com/Bkgediya/feed_system/auth-service/internal/repository"
	"github.com/Bkgediya/feed_system/auth-service/internal/service"
	"github.com/Bkgediya/feed_system/auth-service/pkg/db"
	"github.com/Bkgediya/feed_system/auth-service/pkg/logger"
)

func main() {
	// load db url
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://feeduser:feedpass@localhost:5432/feeddb?sslmode=disable"
	}

	l := logger.New()
	pg, err := db.NewPostgres(dbURL)
	if err != nil {
		l.Fatal("db connect failed", err)
	}
	defer pg.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		l.Fatal("JWT_SECRET not found", err)
	}

	userRepo := repository.NewUserRepo(pg)
	authSvc := service.NewAuthService(userRepo, jwtSecret)

	handler := api.NewHandler(authSvc, l)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	mux.Handle("/v1/signup", handler.Signup())
	mux.Handle("/v1/login", handler.Login())
	// Protected example:
	mux.Handle("/v1/me", handler.AuthMiddleware(handler.Me()))

	l.Info("Auth service starting on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		l.Fatal("server failed", err)
	}
}
