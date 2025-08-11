package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Ararat25/subscription-aggregation-service/docs"
	"github.com/Ararat25/subscription-aggregation-service/internal/config"
	"github.com/Ararat25/subscription-aggregation-service/internal/controller"
	"github.com/Ararat25/subscription-aggregation-service/internal/logger"
	middle "github.com/Ararat25/subscription-aggregation-service/internal/middleware"
	"github.com/Ararat25/subscription-aggregation-service/internal/model"
	"github.com/Ararat25/subscription-aggregation-service/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Subscription aggregation service API
// @version 1.0
// @description API для управления подписками и расчета их стоимости.
// @BasePath /api/v1
func main() {
	conf, err := config.Init()
	if err != nil {
		log.Fatalf("error init config: %v\n", err)
	}

	err = logger.Init(conf.Server.Logging, conf.Server.LogPath)
	if err != nil {
		log.Fatalf("error init logger: %v\n", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Log.Error("failed to sync logger", zap.Error(err))
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := repository.PGRepo{}
	err = db.ConnectDB(ctx, conf.Database.Host, conf.Database.User, conf.Database.Password, conf.Database.Name, conf.Database.Port)
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer func() {
		if err := db.Close(ctx); err != nil {
			logger.Log.Error("failed to close db", zap.Error(err))
		}
	}()

	logger.Log.Info("Successful connection to the database",
		zap.String("host", conf.Database.Host),
		zap.Int("port", conf.Database.Port),
	)

	handler := initApp(&db)
	router := initRouter(handler)
	runApp(ctx, conf, router)
}

// runApp запускает сервер с graceful shutdown
func runApp(ctx context.Context, conf *config.Config, router *chi.Mux) {
	hostPort := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)

	server := &http.Server{
		Addr:    hostPort,
		Handler: router,
	}

	logger.Log.Info("Server starting...", zap.String("addr", hostPort))
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("Server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Log.Info("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.Server.Timeout)
	defer cancel()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Log.Fatal("Server shutdown error", zap.Error(err))
	}

	logger.Log.Info("Server stopped")
}

// initApp инициализирует сервисы приложения
func initApp(db repository.Repo) *controller.Handler {
	authService := model.NewAggregationService(db)

	handler := controller.NewHandler(authService)

	return handler
}

// initRouter настраивает маршруты и middleware для сервера
func initRouter(handler *controller.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middle.ZapLogger(logger.Log))
	r.Use(middleware.Recoverer)
	r.Use(middle.JsonHeader)

	r.Get("/api/v1/doc/*", httpSwagger.WrapHandler)
	r.Post("/api/v1/subscription", handler.CreateSubscription)
	r.Get("/api/v1/subscription/{id}", handler.ReadSubscription)
	r.Put("/api/v1/subscription/update", handler.UpdateSubscription)
	r.Delete("/api/v1/subscription/delete/{id}", handler.DeleteSubscription)
	r.Get("/api/v1/subscriptions", handler.ListSubscriptions)
	r.Get("/api/v1/subscriptions/cost", handler.TotalCost)

	return r
}
