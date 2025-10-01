package app

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"oolio.com/kart/configs"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
)

func StartApplication() {
	defer panicRecovery()

	configs.Logger.Info("Starting the application")

	initializeConfigs()

	router := initializeRoutes()

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

func panicRecovery() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		configs.Logger.Warn("Application panic recovered",
			zap.Any("error", r),
			zap.String("stack", string(stack)),
		)
	}
}

func initializeConfigs() {
	err := configs.InitApplicationConfig()
	if err != nil {
		configs.Logger.Fatal("Failed to load application config", zap.Error(err))
	}

	_ = configs.InitLogger(configs.LogLevel)
}

func gracefulShutdown(server *http.Server) {
	defer func() {
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
