package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"repository/internal/controller"
	"repository/internal/db"
	"repository/internal/repository"
	"repository/internal/service"
	"syscall"
	"time"

	_ "repository/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

// @title Repository
// @version 1.0
// @description	Сервер с бд, для операций CRUD
// @host localhost:8080
// @BasePath /api/users
// @in header
// @name Authorization

func main() {
	dbConn, err := db.InitDBAndMigrate()
	if err != nil {
		log.Fatalf("Failed to init DB and run migrations: %v", err)
	}
	defer dbConn.Close()

	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo)

	r := chi.NewRouter()
	controller.RegisterUserRoutes(r, userService)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("Error creating listner: %v", err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("server is starting at :8080")
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
