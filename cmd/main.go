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
	"time"

	_ "github.com/Ararat25/subscription-aggregation-service/docs"
	"github.com/Ararat25/subscription-aggregation-service/internal/config"
	"github.com/Ararat25/subscription-aggregation-service/internal/controller"
	middle "github.com/Ararat25/subscription-aggregation-service/internal/middleware"
	"github.com/Ararat25/subscription-aggregation-service/internal/model"
	"github.com/Ararat25/subscription-aggregation-service/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	handler := initApp(ctx, conf)

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

	log.Printf("Server starting on %s", hostPort)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

// initApp инициализирует подключение к базе данных и сервисы приложения
func initApp(ctx context.Context, conf *config.Config) *controller.Handler {
	db := repository.PGRepo{}

	err := db.ConnectDB(ctx, conf.Database.Host, conf.Database.User, conf.Database.Password, conf.Database.Name, conf.Database.Port)
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}

	log.Println("Successful connection to the database")

	authService := model.NewAggregationService(&db)

	handler := controller.NewHandler(authService)

	return handler
}

// initRouter настраивает маршруты и middleware для сервера
func initRouter(handler *controller.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middle.JsonHeader)

	r.Get("/api/v1/docs/*", httpSwagger.WrapHandler)
	r.Post("/api/v1/subscription", handler.CreateSubscription)
	r.Get("/api/v1/subscription/{id}", handler.ReadSubscription)
	r.Put("/api/v1/subscription/update", handler.UpdateSubscription)
	r.Delete("/api/v1/subscription/delete/{id}", handler.DeleteSubscription)
	r.Get("/api/v1/subscriptions", handler.ListSubscriptions)
	r.Get("/api/v1/subscriptions/cost", handler.TotalCost)

	return r
}
