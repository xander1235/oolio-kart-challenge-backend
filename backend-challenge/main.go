package main

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"oolio.com/kart/configs"
	"oolio.com/kart/repositories"
	"oolio.com/kart/routes"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
)

func main() {
	startApplication()
}

// startApplication entry point for the application
func startApplication() {
	defer panicRecovery()

	configs.Logger.Info("Starting the application")

	initializeConfigs()

	router := routes.InitializeRoutes()

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(configs.Port),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		configs.Logger.Info("starting server", zap.Int("port", configs.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			configs.Logger.Fatal("server error", zap.Error(err))
		}
	}()

	gracefulShutdown(server)
}

// panicRecovery recovers from panics and logs the error
func panicRecovery() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		configs.Logger.Warn("Application panic recovered",
			zap.Any("error", r),
			zap.String("stack", string(stack)),
		)
	}
}

// initializeConfigs initializes the application configuration
func initializeConfigs() {
	err := configs.InitApplicationConfig()
	if err != nil {
		configs.Logger.Fatal("Failed to load application config", zap.Error(err))
	}

	_ = configs.InitLogger(configs.LogLevel)

	if err := repositories.Initialize(context.Background()); err != nil {
		configs.Logger.Warn("Failed to initialize database connection pool", zap.Error(err))
		return
	}
}

// gracefulShutdown gracefully shuts down the server
func gracefulShutdown(server *http.Server) {
	defer func() {
		repositories.Close()
		configs.Logger.Info("Database connection pool closed")

		err := configs.Logger.Sync()
		if err != nil {
			configs.Logger.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	configs.Logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		configs.Logger.Fatal("forced server shutdown", zap.Error(err))
	}

	configs.Logger.Info("server exited gracefully")
}
