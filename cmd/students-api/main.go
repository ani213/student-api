package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ani213/student-api/internal/config"
	"github.com/ani213/student-api/internal/http/handler/student"
	"github.com/ani213/student-api/internal/storage/sqlite"
)

func main() {
	// Load application configuration from YAML or environment
	// Panics if required fields are missing
	cfg := config.MustLoad()

	// Storage setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String("storage_path", cfg.StoragePath), slog.String("ENV", cfg.Env))

	// Initialize the HTTP request multiplexer (router)
	router := http.NewServeMux()

	// Register a simple GET endpoint at the root path
	// router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Welcome to student API"))
	// })

	// Register the student creation handler at the "/students" path
	router.HandleFunc("POST /api/student", student.CreateStudent(storage))
	router.HandleFunc("GET /api/students", student.GetStudents(storage))
	router.HandleFunc("GET /api/student/{id}", student.GetStudentById(storage))
	// Create the HTTP server with the loaded address and request handler
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server started", slog.String("address", cfg.Addr))
	// Create a channel to listen for OS termination signals
	done := make(chan os.Signal, 1)

	// Notify the 'done' channel when an interrupt or terminate signal is received
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the HTTP server in a separate goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error:", err)
		}
	}()

	// Block the main goroutine until a signal is received
	<-done

	// Begin graceful shutdown
	slog.Info("Shutting down the server...")

	// Create a context with a timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	} else {
		slog.Info("Server shutdown successfully")
	}
}
