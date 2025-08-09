package main

import (
	"fmt"
	"log"
	"net/http"

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
// @description This is
// @host localhost:8080
// @BasePath /
func main() {
	conf, err := config.Init()
	if err != nil {
		log.Fatalf("error init config: %v\n", err)
	}

	handler := initApp(conf)

	router := initRouter(handler)

	hostPort := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)

	log.Printf("Server starting on %s", hostPort)

	err = http.ListenAndServe(hostPort, router)
	if err != nil {
		log.Fatalf("Start server error: %s", err.Error())
	}
}

// initApp инициализирует конфигурацию, подключение к базе данных и сервисы приложения
func initApp(conf *config.Config) *controller.Handler {
	db := repository.PGRepo{}

	err := db.ConnectDB(conf.Database.Host, conf.Database.User, conf.Database.Password, conf.Database.Name, conf.Database.Port)
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
